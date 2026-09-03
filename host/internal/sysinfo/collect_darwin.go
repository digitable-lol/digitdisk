// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

// The macOS half of the collector.  Everything here asks the kernel — through
// sysctl(3), getfsstat(2), the routing socket, the libSystem functions wrapped
// in internal/libsystem, and the IOKit and CoreFoundation ones wrapped in
// internal/iokit.  The decoding of what comes back lives in internal/darwinsys
// and internal/iokit, and both are built and tested everywhere, including on
// machines that are not Macs.
//
// WHERE EVERY FACT COMES FROM.  The interfaces are the documented ones, and
// nothing below is read out of another project's source:
//
//	узел                    sysctl kern.hostname (os.Hostname as a fallback)
//	ядро, версия            sysctl kern.ostype, kern.osrelease, kern.version
//	выпуск системы          sysctl kern.osproductversion, kern.osversion
//	машина, модель          sysctl hw.machine, hw.model
//	время работы            sysctl kern.boottime (struct timeval)
//	средняя загрузка        sysctl vm.loadavg (struct loadavg)
//	ядер                    sysctl hw.logicalcpu, hw.ncpu
//	память всего, страница  sysctl hw.memsize, hw.pagesize
//	своп                    sysctl vm.swapusage (struct xsw_usage)
//	занято ЦП               host_statistics(HOST_CPU_LOAD_INFO)
//	разбивка памяти         host_statistics64(HOST_VM_INFO64)
//	процессы                sysctl kern.proc.all (struct kinfo_proc)
//	память и потоки, время  proc_pidinfo(PROC_PIDTASKINFO)
//	командные строки        sysctl KERN_PROCARGS2 (по MIB, с pid)
//	диски                   getfsstat(2) (struct statfs)
//	интерфейсы, адреса      net.Interfaces (стандартная библиотека)
//	счётчики интерфейсов    sysctl NET_RT_IFLIST2 (struct if_msghdr2/if_data64)
//	видеокарты              IOServiceGetMatchingServices по классам
//	                        IOAccelerator и IOPCIDevice, затем
//	                        IORegistryEntryCreateCFProperties
//
// WHAT IS STILL NOT MEASURED, AND WHY.  Two kinds, and they are not the same
// kind:
//
//   - Closed by permission, not by the language.  A process belonging to
//     another user answers neither proc_pidinfo nor KERN_PROCARGS2 unless the
//     caller is the administrator — the kernel checks the uid and refuses.
//     Running the tool under sudo fills those rows in; nothing else will.
//   - Not published by the system in a form anybody may rely on.  Die
//     temperature comes from the SMC through IOKit, and Apple documents no
//     interface to it: what circulates is a reverse-engineered structure.  A
//     number read that way would be a guess wearing a unit, so there is none.
//
// Both are named in Status.Missing, in one short phrase each.  The report
// prints the names and nothing else; `--why` prints the phrases.
package sysinfo

import (
	"net"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"digitdisk/internal/darwinsys"
	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/iokit"
	"digitdisk/internal/lang"
	"digitdisk/internal/libsystem"
	"digitdisk/internal/procfs"
)

// Mount flags from <sys/mount.h>.  The standard library's darwin syscall
// package exports neither, so they are written out here with their source.
const (
	mntNoWait = 2 // MNT_NOWAIT: answer from what is already known
	mntRDOnly = 1 // MNT_RDONLY
)

// Collector reads a snapshot of this Mac.
//
// It has no filesystem roots to point at: macOS publishes these facts through
// sysctl and libSystem, not through files, so the seam the Linux collector
// opens for its tests is not available here.  What is testable — every
// structure the kernel answers with — is tested in internal/darwinsys instead,
// and what can only be checked on a Mac is checked on the Mac, at run time,
// before any of it is published.
type Collector struct {
	// SampleWindow is how long the CPU delta sample runs, for the machine
	// as a whole and for each process.  Zero disables sampling.
	SampleWindow time.Duration
	// Top is how many processes to list per ranking.
	Top int
	// GPUTool would allow running the graphics driver's own program for the
	// numbers no file publishes.  It is here so the command line has the
	// same shape on both systems; macOS publishes no video card at all, so
	// nothing reads it.
	GPUTool bool
}

// New returns a Collector pointed at the live system.
func New() Collector {
	return Collector{SampleWindow: 200 * time.Millisecond, Top: 10}
}

// sysctlRaw returns the raw bytes of a sysctl node.
//
// syscall.Sysctl is the only sysctl the standard library exports, and it is
// written for string-valued nodes: it strips one trailing NUL.  A numeric or
// structure node therefore arrives up to one byte short, and darwinsys.Padded
// puts that byte back — see its comment for why that is exact rather than
// hopeful.
func sysctlRaw(name string) ([]byte, error) {
	s, err := syscall.Sysctl(name)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// sysctlGrowing reads a node whose size changes while it is being read.  The
// process table is exactly that: the size is asked for first and the bytes
// second, and a process that starts in between makes the kernel refuse.  Trying
// again is the answer; trying forever is not.
func sysctlGrowing(name string, attempts int) ([]byte, error) {
	var err error
	for i := 0; i < attempts; i++ {
		var b []byte
		if b, err = sysctlRaw(name); err == nil {
			return b, nil
		}
	}
	return nil, err
}

// Collect gathers the whole snapshot.  A source that fails is recorded in
// Status.Missing rather than aborting the run, and a fact this system does not
// publish at all is recorded there under its Fact constant.
func (c Collector) Collect() Status {
	st := Status{TakenAt: time.Now().Format(time.RFC3339), Missing: map[string]lang.Phrase{}}

	c.host(&st)
	c.load(&st)
	c.memory(&st)
	c.sample(&st)

	disks, err := c.Disks()
	if err != nil {
		st.Missing["getfsstat"] = lang.FromError(err)
	}
	st.Disks = disks
	st.Network = c.network(&st)

	st.Sensors = nil
	st.Missing[FactSensors] = lang.Say("macOS не публикует показания датчиков, а угадывать их формат нельзя")

	// The video cards are NOT the same case as the sensors, and the
	// difference was measured rather than assumed.  What a Mac knows about
	// its graphics lives in the IORegistry, and IOKit answers with Core
	// Foundation objects rather than with numbers — but that turned out to
	// be a difficulty of decoding, not a wall: the frameworks are called
	// the same way internal/libsystem calls libSystem, without cgo, and
	// internal/iokit does it.  See TestДверьОткрывается there, which
	// refuses the answer unless what the registry says about the model of
	// this machine matches what sysctl says.
	//
	// What comes back is a name, a driver and, on a card that sits on PCI,
	// the bus address and the identifiers.  What does NOT come back is
	// counters: the registry publishes neither the busy share nor the
	// memory in use, and a card that shares the machine's memory publishes
	// no memory of its own either.  Those are said as absences and not
	// filled in from anywhere.
	st.GPUs = iokit.Cards()
	if len(st.GPUs) == 0 {
		st.Missing[FactGPUs] = lang.Say("в реестре устройств нет ни одной записи о видеокарте")
	} else if !anyGPUNumbers(st.GPUs) {
		st.Missing[FactGPUNumbers] = lang.Say("реестр устройств macOS называет видеокарту и её драйвер, но счётчиков загрузки и памяти в нём нет")
	}

	if len(st.Missing) == 0 {
		st.Missing = nil
	}
	return st
}

// text reads a string-valued sysctl, naming the node in Missing when it is not
// there.  An absent node is a fact about this system, not a failure of the run.
func text(st *Status, name string) string {
	v, err := syscall.Sysctl(name)
	if err != nil {
		st.Missing[name] = lang.FromError(err)
		return ""
	}
	return strings.TrimSpace(v)
}

func (c Collector) host(st *Status) {
	st.Host.Hostname = text(st, "kern.hostname")
	if st.Host.Hostname == "" {
		if name, err := os.Hostname(); err == nil {
			st.Host.Hostname = name
			delete(st.Missing, "kern.hostname")
		}
	}
	st.Host.KernelRelease = text(st, "kern.osrelease")
	st.Host.KernelVersion = text(st, "kern.version")
	st.Host.Machine = text(st, "hw.machine")
	st.Host.Model = text(st, "hw.model")
	if v, err := syscall.Sysctl("machdep.cpu.brand_string"); err == nil && strings.TrimSpace(v) != "" {
		st.Host.CPUModel = strings.TrimSpace(v)
	} else {
		st.Missing[FactCPUModel] = lang.Say("узел machdep.cpu.brand_string ничего не ответил — процессор себя не назвал")
	}
	environment(st)

	// The product name is not a sysctl: kern.ostype answers "Darwin", which
	// is the kernel, not the system the user knows.  "macOS" is written out
	// here because this file only ever runs on it; the version and the
	// build next to it are measured, not assumed.
	if version := text(st, "kern.osproductversion"); version != "" {
		st.Host.DistroID = "macos"
		st.Host.DistroVersion = version
		st.Host.Distro = "macOS " + version
		if build := text(st, "kern.osversion"); build != "" {
			st.Host.Distro += " (" + build + ")"
		}
	}

	b, err := sysctlRaw("kern.boottime")
	if err != nil {
		st.Missing["kern.boottime"] = lang.FromError(err)
		return
	}
	boot, ok := darwinsys.ParseTimeval(b)
	if !ok {
		st.Missing["kern.boottime"] = lang.Say("ответ не читается как момент времени — время работы не считаем")
		return
	}
	st.Host.BootTime = boot.Unix()
	if up := time.Since(boot).Seconds(); up > 0 {
		st.Host.UptimeSeconds = up
		st.Host.UptimeHuman = HumanDuration(up)
	}
}

func (c Collector) load(st *Status) {
	if b, err := sysctlRaw("vm.loadavg"); err != nil {
		st.Missing["vm.loadavg"] = lang.FromError(err)
	} else if one, five, fifteen, ok := darwinsys.ParseLoadAvg(b); ok {
		st.Load.One, st.Load.Five, st.Load.Fifteen = one, five, fifteen
	} else {
		st.Missing["vm.loadavg"] = lang.Say("ответ не читается как средние загрузки — не показываем")
	}
	// /proc/loadavg carries three more numbers macOS does not publish.  The
	// report never prints them, so this is a note for the JSON only.
	st.Missing[FactLoadEntities] = lang.Say("macOS не публикует длину очереди планировщика")

	for _, name := range []string{"hw.logicalcpu", "hw.ncpu"} {
		if n, err := syscall.SysctlUint32(name); err == nil && n > 0 {
			st.Load.CPUCount = int(n)
			break
		}
	}
}

func (c Collector) memory(st *Status) {
	m := procfs.Memory{Present: map[string]bool{}, Raw: map[string]uint64{}}
	if b, err := sysctlRaw("hw.memsize"); err != nil {
		st.Missing["hw.memsize"] = lang.FromError(err)
	} else if v, ok := darwinsys.ParseUint64(b); ok {
		m.Total = v
		m.Present[procfs.FieldTotal] = true
	} else {
		st.Missing["hw.memsize"] = lang.Say("ответ не восьмибайтовый — объём памяти не показываем")
	}

	if b, err := sysctlRaw("vm.swapusage"); err != nil {
		st.Missing["vm.swapusage"] = lang.FromError(err)
	} else if s, ok := darwinsys.ParseSwapUsage(b); ok {
		m.SwapTotal, m.SwapFree, m.SwapUsed = s.Total, s.Avail, s.Used
		m.Present[procfs.FieldSwapTotal] = true
		m.Present[procfs.FieldSwapFree] = true
		m.Present[procfs.FieldSwapUsed] = true
	} else {
		st.Missing["vm.swapusage"] = lang.Say("ответ не читается как сведения о свопе — не показываем")
	}

	c.pages(st, &m)
	st.Memory = m
}

// pages fills the breakdown from host_statistics64, and is where the page
// counts become bytes.
//
// Nothing is written into the snapshot until darwinsys.ParseVMStat64 has
// checked the answer against the machine's own page count, which came from a
// different source — hw.memsize and hw.pagesize.  Two sources agreeing is the
// whole warrant for printing any of these numbers.
func (c Collector) pages(st *Status, m *procfs.Memory) {
	notes := func(why lang.Phrase) { st.Missing[FactMemoryPages] = why }

	pageSize, err := syscall.SysctlUint32("hw.pagesize")
	if err != nil || pageSize == 0 {
		notes(lang.Say("без размера страницы разбивку памяти не пересчитать в байты"))
		return
	}
	if m.Total == 0 {
		notes(lang.Say("без объёма памяти разбивку не с чем сверить"))
		return
	}

	b, err := libsystem.HostStatistics64(darwinsys.FlavorVMInfo64, darwinsys.VMInfo64Count)
	if err != nil {
		notes(lang.Say("ядро не дало разбивку памяти: %s", err))
		return
	}
	v, ok := darwinsys.ParseVMStat64(b, m.Total/uint64(pageSize))
	if !ok {
		notes(lang.Say("разбивка памяти не сошлась с объёмом памяти машины — не публикуем"))
		return
	}

	p := uint64(pageSize)
	// The kernel folds the read-ahead pages into its free count, so what is
	// really free is what is left after taking them back out.
	m.Free = v.TrulyFree() * p
	m.Present[procfs.FieldFree] = true
	// External pages are the ones backed by a file: the page cache, and the
	// nearest thing macOS has to the buff/cache column of a Linux report.
	m.BuffCache = v.External * p
	m.Present[procfs.FieldBuffCache] = true
	// Available and used are the two sides of one statement: memory that is
	// neither free nor file cache is in use.  It is the same reading free(1)
	// gives on Linux, and the report prints the arithmetic next to the
	// number so nobody has to take it on trust.
	m.Available = m.Free + m.BuffCache
	m.Present[procfs.FieldAvailable] = true
	if m.Available <= m.Total {
		m.Used = m.Total - m.Available
		m.Present[procfs.FieldUsed] = true
	}
	// Wired and compressed are what a Mac owner is used to seeing, and both
	// are counts the kernel keeps itself.
	m.Raw[procfs.RawWired] = v.Wired * p
	m.Raw[procfs.RawCompressed] = v.Compressor * p
	m.Raw[procfs.RawActive] = v.Active * p
	m.Raw[procfs.RawInactive] = v.Inactive * p
	m.Raw[procfs.RawSpeculative] = v.Speculative * p
	m.Raw[procfs.RawPurgeable] = v.Purgeable * p
	m.Raw[procfs.RawAnonymous] = v.Internal * p
}

// procSample is one process as one pass saw it.
type procSample struct {
	proc Proc
	// cpuNanos is processor time consumed since the process started.  It is
	// only set when proc_pidinfo answered, which for a process belonging to
	// somebody else it will not.
	cpuNanos uint64
	detailed bool
	// running says the process had at least one thread on a processor when
	// it was read.
	running bool
	// at is when this process was read.  A pass over nine hundred processes
	// is not instantaneous, so a CPU delta is divided by each process's own
	// interval rather than by the nominal window.
	at time.Time
}

// sample takes the whole time-dependent half of the snapshot: the machine's
// busy share and every process, both measured across one window.
func (c Collector) sample(st *Status) {
	if c.SampleWindow <= 0 {
		// Everything measured across a window is measured across no
		// window at all, which is not a measurement.
		st.Missing[FactCPUBusy] = lang.Say("окно замера нулевое — доля занятого процессора не измерялась")
		st.Missing[FactCores] = lang.Say("окно замера нулевое — доля занятого времени каждого ядра не измерялась")
		st.Missing[FactProcessCPU] = lang.Say("окно замера нулевое — процессорное время процессов не измерялось")
		if after, ok := c.snapProcesses(st, true); ok {
			c.rank(st, nil, after)
		}
		return
	}

	ticksBefore, haveTicks := c.cpuReading()
	started := time.Now()
	before, _ := c.snapProcesses(st, false)
	time.Sleep(c.SampleWindow)
	after, gotProcs := c.snapProcesses(st, true)

	c.busyShare(st, ticksBefore, haveTicks, started)

	if gotProcs {
		c.rank(st, before, after)
	}
}

// busyShare closes the CPU sample, waiting longer if the counters have not
// moved yet.
//
// A machine whose processors did nothing at all in two hundred milliseconds
// has counters that did not move, and there is nothing to divide by.  That is
// not hypothetical: on a virtual Mac the tick counters advance in steps
// coarser than the sample window, and a single window lands between two of
// them about half the time.  Waiting out another window or two turns a dash
// that means "the machine was idle" into the number the reader asked for,
// and the window that is finally reported is the one that was really used.
func (c Collector) busyShare(st *Status, before cpuReading, have bool, started time.Time) {
	if !have {
		st.Missing[FactCPUBusy] = lang.Say("ядро не дало счётчики процессорного времени")
		st.Missing[FactCores] = lang.Say("ядро не дало счётчики процессорного времени")
		return
	}
	for attempt := 0; ; attempt++ {
		if now, ok := c.cpuReading(); ok {
			if share, ok := before.whole.BusyShare(now.whole); ok {
				st.Load.BusyPercent = &share
				st.Load.SampleMillis = time.Since(started).Milliseconds()
				// The same pair of readings carries the same share
				// for each processor separately.
				c.cores(st, before, now)
				return
			}
		}
		if attempt >= busyRetries {
			st.Missing[FactCPUBusy] = lang.Say("за окно замера счётчики процессора не сдвинулись")
			st.Missing[FactCores] = lang.Say("за окно замера счётчики процессоров не сдвинулись")
			return
		}
		time.Sleep(c.SampleWindow)
	}
}

// cores publishes the share of each processor separately — but only after two
// checks, because a misread array here would put every processor's number in
// the wrong row and nothing on the screen would look wrong.
//
// The first check is the count: the kernel says how many processors it wrote
// about, and that has to be the number the machine says it has.  The second is
// arithmetic: the machine-wide counters are the sum of the per-processor ones,
// so the mean of the shares has to be the share already measured for the
// machine.  Either one failing means these bytes are not what we think they
// are, and then there are no per-processor numbers at all.
func (c Collector) cores(st *Status, before, now cpuReading) {
	if len(before.cores) == 0 || len(now.cores) != len(before.cores) {
		st.Missing[FactCores] = lang.Say("ядро не дало счётчики по каждому процессору отдельно")
		return
	}
	if st.Load.CPUCount > 0 && len(now.cores) != st.Load.CPUCount {
		st.Missing[FactCores] = lang.Say("процессоров в ответе ядра не столько, сколько машина насчитала у себя, — по ядрам не публикуем")
		return
	}
	out := make([]Core, 0, len(now.cores))
	for i := range now.cores {
		core := Core{Index: i, Name: "cpu" + strconv.Itoa(i)}
		if share, ok := before.cores[i].BusyShare(now.cores[i]); ok {
			s := share
			core.BusyPercent = &s
		}
		out = append(out, core)
	}
	if ok, why := coresAgree(out, st.Load.BusyPercent); !ok {
		st.Missing[FactCores] = why
		return
	}
	st.Load.Cores = out
}

// busyRetries is how many extra windows the CPU sample may wait for the
// counters to move.  Three is enough for the coarsest step seen on a virtual
// Mac and still bounded: with the default window the whole snapshot cannot
// grow past a second because of it.
const busyRetries = 3

// cpuReading is one look at the processor tick counters: the machine as a
// whole, and each of its processors.  The two come from two calls, but from
// the same counters — the machine-wide flavor is the sum of the per-processor
// one — so they are taken together and used together.
type cpuReading struct {
	whole darwinsys.CPUTicks
	cores []darwinsys.CPUTicks
}

// cpuReading reads the tick counters.  The per-processor array is allowed to
// be missing: the machine-wide share is the older and the more important of
// the two, and it must not be lost because the newer call refused.
func (c Collector) cpuReading() (cpuReading, bool) {
	b, err := libsystem.HostStatistics(darwinsys.FlavorCPULoad, darwinsys.CPULoadCount)
	if err != nil {
		return cpuReading{}, false
	}
	whole, ok := darwinsys.ParseCPUTicks(b)
	if !ok {
		return cpuReading{}, false
	}
	out := cpuReading{whole: whole}
	if cpus, data, err := libsystem.ProcessorInfo(darwinsys.FlavorProcessorCPULoad); err == nil {
		if cores, ok := darwinsys.ParseProcessorTicks(data, cpus); ok {
			out.cores = cores
		}
	}
	return out, true
}

// snapProcesses reads the process table once.  The bool says whether anything
// came back at all; report says whether the reasons for what is missing should
// be written into the snapshot, so the first of two passes stays silent.
func (c Collector) snapProcesses(st *Status, report bool) (map[int]procSample, bool) {
	notes := func(key string, why lang.Phrase) {
		if report {
			st.Missing[key] = why
		}
	}

	b, err := sysctlGrowing("kern.proc.all", 3)
	if err != nil {
		notes("kern.proc.all", lang.FromError(err))
		return nil, false
	}
	procs, ok := darwinsys.ParseProcs(b)
	if !ok {
		notes("kern.proc.all", lang.Say("ответ не делится на записи процессов — список не публикуем"))
		return nil, false
	}
	if !darwinsys.Verify(procs, os.Getpid(), os.Getppid(), os.Getuid()) {
		notes("kern.proc.all", lang.Say("самопроверка записи о процессе не сошлась — числа из неё не публикуем"))
		return nil, false
	}

	// The self-check for the second source.  proc_pidinfo is asked about
	// this very process first, and its answer is checked against what this
	// process knows about itself: it is running, so it holds pages and has
	// at least one thread.  Until that agrees, no process's memory is
	// published — a table of misread numbers is worse than an empty one.
	memTotal := st.Memory.Total
	detailOK := memTotal > 0 && c.selfTaskInfo(memTotal)
	if !detailOK && report {
		st.Missing[FactProcessRSS] = lang.Say("самопроверка памяти процессов не сошлась — их память и потоки не публикуем")
		st.Missing[FactThreads] = lang.Say("самопроверка памяти процессов не сошлась — их память и потоки не публикуем")
	}

	out := make(map[int]procSample, len(procs))
	denied := 0
	for _, p := range procs {
		s := procSample{
			at: time.Now(),
			proc: Proc{
				PID: p.PID, PPID: p.PPID, UID: p.UID,
				State: p.State, Comm: p.Comm,
			},
		}
		if detailOK {
			if ti, ok := taskInfo(p.PID, memTotal); ok {
				s.proc.RSSBytes = int64(ti.ResidentBytes)
				s.proc.VSizeBytes = ti.VirtualBytes
				s.proc.Threads = ti.Threads
				s.cpuNanos = ti.CPUNanos
				s.running = ti.Running > 0
				s.detailed = true
			} else {
				denied++
			}
		}
		out[p.PID] = s
	}

	if report {
		st.Processes.Total = len(procs)
		// How many processes are running is counted in rank, from the
		// thread counts: the scheduler state in a process record says
		// "runnable" for every process that is merely alive, and counting
		// by it reported the whole list as running.
		if !detailOK {
			st.Missing[FactRunning] = lang.Say("сколько процессов работает прямо сейчас, видно только по их потокам")
		}
		// A process belonging to another user is not a process we failed
		// to read: the kernel refused on purpose, and it refuses the same
		// way to every tool that is not the administrator.
		if detailOK && denied > 0 {
			st.Missing[FactOtherUsers] = lang.Say("память, потоки и командные строки чужих процессов видны только администратору — запустите под sudo")
		}
		// The scheduler state macOS keeps in a process record does not
		// separate uninterruptible sleep from ordinary sleep.
		st.Missing[FactBlocked] = lang.Say("macOS не различает заблокированные процессы среди спящих")
	}
	return out, true
}

// selfTaskInfo is the run-time proof that struct proc_taskinfo is laid out the
// way darwinsys describes it.  This process is running, so it holds resident
// pages and has threads; a decoder reading the wrong offsets does not produce
// a small plausible pair of numbers on a live process, it produces nonsense
// that ParseTaskInfo refuses.
func (c Collector) selfTaskInfo(memTotal uint64) bool {
	ti, ok := taskInfo(os.Getpid(), memTotal)
	return ok && ti.ResidentBytes > 0 && ti.Threads >= 1
}

func taskInfo(pid int, memTotal uint64) (darwinsys.TaskInfo, bool) {
	b, err := libsystem.ProcTaskInfo(pid, darwinsys.TaskInfoFlavor, darwinsys.TaskInfoSize)
	if err != nil {
		return darwinsys.TaskInfo{}, false
	}
	return darwinsys.ParseTaskInfo(b, memTotal)
}

// rank turns two passes into the two lists the report prints.
func (c Collector) rank(st *Status, before, after map[int]procSample) {
	procs := make([]Proc, 0, len(after))
	detailed := 0
	for pid, s := range after {
		if s.detailed {
			detailed++
			st.Processes.Threads += s.proc.Threads
			if s.running {
				st.Processes.Running++
			}
		}
		if b, ok := before[pid]; ok && b.detailed && s.detailed {
			window := s.at.Sub(b.at).Seconds()
			if window > 0 && s.cpuNanos >= b.cpuNanos {
				pct := float64(s.cpuNanos-b.cpuNanos) / 1e9 * 100 / window
				s.proc.CPUPercent = &pct
			}
		}
		procs = append(procs, s.proc)
	}
	st.Processes.WithDetail = detailed
	if detailed == 0 {
		if _, named := st.Unmeasured(FactProcessCPU); !named {
			st.Missing[FactProcessCPU] = lang.Say("без памяти процессов их процессорное время тоже не публикуем")
		}
		return
	}

	byMem := make([]Proc, 0, len(procs))
	for _, p := range procs {
		if p.RSSBytes > 0 || p.Threads > 0 {
			byMem = append(byMem, p)
		}
	}
	sort.Slice(byMem, func(i, j int) bool {
		if byMem[i].RSSBytes != byMem[j].RSSBytes {
			return byMem[i].RSSBytes > byMem[j].RSSBytes
		}
		return byMem[i].PID < byMem[j].PID
	})
	st.Processes.TopByMemory = head(byMem, c.Top)

	byCPU := append([]Proc(nil), procs...)
	sort.Slice(byCPU, func(i, j int) bool {
		if a, b := cpuOrZero(byCPU[i]), cpuOrZero(byCPU[j]); a != b {
			return a > b
		}
		return byCPU[i].PID < byCPU[j].PID
	})
	if len(byCPU) > 0 && byCPU[0].CPUPercent == nil {
		if _, named := st.Unmeasured(FactProcessCPU); !named {
			st.Missing[FactProcessCPU] = lang.Say("ни один процесс не прожил всё окно замера")
		}
		st.Processes.TopByCPU = nil
	} else {
		st.Processes.TopByCPU = head(byCPU, c.Top)
	}

	// The command line reader checks itself once, so it is built once and
	// handed to both lists.
	args := argsReader(st)
	describe(args, st.Processes.TopByMemory)
	describe(args, st.Processes.TopByCPU)
}

// describe fills in the command line and the user name of the processes that
// made it into a list — and only those.  Asking the kernel for nine hundred
// argument blocks to print twenty of them would cost a megabyte of copying per
// process for nothing.
func describe(args func(int) (string, bool), list []Proc) {
	for i := range list {
		if u, err := user.LookupId(strconv.Itoa(list[i].UID)); err == nil {
			list[i].User = u.Username
		}
		if args != nil {
			if line, ok := args(list[i].PID); ok {
				list[i].Cmdline = line
			}
		}
	}
}

// argsReader returns a function that reads one process's command line, or nil
// when command lines are not to be published at all.
//
// The self-check is the strongest one in this file, and it is worth saying why
// it is strong: the argument block of this very process is decoded and then
// compared, string by string, with os.Args — which the Go runtime received
// from the same kernel by a completely different road.  Two roads agreeing on
// every word of a real command line is not something a misread layout does.
//
// The check runs once per snapshot, and its result is cached in the closure.
func argsReader(st *Status) func(int) (string, bool) {
	argmax, err := syscall.SysctlUint32("kern.argmax")
	if err != nil || argmax == 0 {
		st.Missing[FactProcessArgs] = lang.Say("ядро не назвало предельный размер блока аргументов")
		return nil
	}
	read := func(pid int) ([]string, bool) {
		b, err := libsystem.SysctlRaw(darwinsys.ArgsMIB(pid), int(argmax))
		if err != nil {
			return nil, false
		}
		_, argv, ok := darwinsys.ParseProcArgs2(b)
		return argv, ok
	}
	mine, ok := read(os.Getpid())
	if !ok || !darwinsys.SameArgv(mine, os.Args) {
		st.Missing[FactProcessArgs] = lang.Say("самопроверка командной строки не сошлась — чужих командных строк не публикуем")
		return nil
	}
	return func(pid int) (string, bool) {
		argv, ok := read(pid)
		if !ok || len(argv) == 0 {
			return "", false
		}
		return strings.Join(argv, " "), true
	}
}

// Disks lists mounted filesystems that report a non-zero size, which is the
// same rule df(1) applies to hide purely synthetic mounts, and the same rule
// the Linux collector applies.
func (c Collector) Disks() ([]Disk, error) {
	n, err := syscall.Getfsstat(nil, mntNoWait)
	if err != nil {
		return nil, err
	}
	// Slack for mounts that appear between the counting call and this one.
	buf := make([]syscall.Statfs_t, n+8)
	n, err = syscall.Getfsstat(buf, mntNoWait)
	if err != nil {
		return nil, err
	}
	if n > len(buf) {
		n = len(buf)
	}
	var out []Disk
	for i := 0; i < n; i++ {
		fs := buf[i]
		if fs.Blocks == 0 {
			continue
		}
		d := Disk{
			Source:      charString(fs.Mntfromname[:]),
			MountPoint:  charString(fs.Mntonname[:]),
			FSType:      charString(fs.Fstypename[:]),
			ReadOnly:    fs.Flags&mntRDOnly != 0,
			InodesTotal: fs.Files,
			InodesFree:  fs.Ffree,
		}
		bs := uint64(fs.Bsize)
		d.TotalBytes = fs.Blocks * bs
		d.UsedBytes = (fs.Blocks - fs.Bfree) * bs
		d.AvailableBytes = fs.Bavail * bs
		if d.UsedBytes+d.AvailableBytes > 0 {
			d.UsePercent = 100 * float64(d.UsedBytes) / float64(d.UsedBytes+d.AvailableBytes)
		}
		out = append(out, d)
	}
	return out, nil
}

// network lists the interfaces with their addresses, and joins the counters to
// them by interface index when the counters can be trusted.
func (c Collector) network(st *Status) []Iface {
	ifaces, err := net.Interfaces()
	if err != nil {
		st.Missing["net.Interfaces"] = lang.FromError(err)
		return nil
	}
	mtus := make(map[int]int, len(ifaces))
	for _, in := range ifaces {
		mtus[in.Index] = in.MTU
	}

	counters := map[int]darwinsys.IfCounters{}
	if b, err := syscall.RouteRIB(syscall.NET_RT_IFLIST2, 0); err != nil {
		st.Missing[FactNetCounters] = lang.Say("список интерфейсов не прочитался: %s", err)
	} else if got := darwinsys.ParseIfList2(b); darwinsys.VerifyIfList(got, mtus) {
		counters = got
		st.Missing[FactNetTxDrops] = lang.Say("macOS не считает отброшенные исходящие пакеты")
	} else {
		st.Missing[FactNetCounters] = lang.Say("самопроверка счётчиков интерфейсов не сошлась — не публикуем")
	}

	addrs := interfaceAddresses()
	out := make([]Iface, 0, len(ifaces))
	for _, in := range ifaces {
		i := Iface{NetCounters: procfs.NetCounters{Name: in.Name}}
		if n, ok := counters[in.Index]; ok {
			i.RxBytes, i.RxPackets, i.RxErrors, i.RxDropped = n.RxBytes, n.RxPackets, n.RxErrors, n.RxDropped
			i.TxBytes, i.TxPackets, i.TxErrors = n.TxBytes, n.TxPackets, n.TxErrors
		}
		if len(in.HardwareAddr) > 0 {
			i.MAC = in.HardwareAddr.String()
		}
		i.MTU = in.MTU
		i.OperState = operState(in.Flags)
		i.Addresses = addrs[in.Name]
		out = append(out, i)
	}
	return out
}

// operState renders what macOS publishes about an interface the way Linux's
// operstate reads.  The two are not the same thing and are not pretended to
// be: macOS has flags, so a raised IFF_UP is "up" and its absence is "down",
// and the "unknown" and "dormant" a Linux driver can report never appear.
func operState(f net.Flags) string {
	if f&net.FlagUp != 0 {
		return "up"
	}
	return "down"
}

// charString reads a NUL-terminated char[] field of a C structure.
func charString(c []int8) string {
	b := make([]byte, 0, len(c))
	for _, ch := range c {
		if ch == 0 {
			break
		}
		b = append(b, byte(ch))
	}
	return string(b)
}

// anyGPUNumbers answers whether any card came back with a number beside its
// name.  It is asked before the absence is declared, so that a Mac which one
// day does publish a counter stops being told it has none.
func anyGPUNumbers(cards []gpuinfo.Card) bool {
	for _, c := range cards {
		if c.MemoryTotalBytes != nil || c.MemoryUsedBytes != nil || c.BusyPercent != nil {
			return true
		}
	}
	return false
}
