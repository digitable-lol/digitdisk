// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

// The macOS half of the collector.  Everything here asks the kernel through
// sysctl(3), getfsstat(2) and the routing socket; the decoding of what comes
// back lives in internal/darwinsys and is built and tested everywhere.
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
//	память всего            sysctl hw.memsize
//	своп                    sysctl vm.swapusage (struct xsw_usage)
//	процессы                sysctl kern.proc.all (struct kinfo_proc)
//	диски                   getfsstat(2) (struct statfs)
//	интерфейсы, адреса      net.Interfaces (стандартная библиотека)
//	счётчики интерфейсов    sysctl NET_RT_IFLIST2 (struct if_msghdr2/if_data64)
//
// WHAT MACOS DOES NOT GIVE US WITHOUT cgo, AND IS THEREFORE LEFT EMPTY.  Each
// of these is named in Status.Missing with its reason, and the report prints
// the reason where the number would have been:
//
//   - the CPU busy share and the per-process CPU time: host_statistics with
//     HOST_CPU_LOAD_INFO, and the task port of each process.  Mach calls, not
//     system calls: reaching them means either cgo or issuing Mach traps by
//     hand, which Apple does not support from outside libSystem.
//   - free / active / inactive / wired / compressed memory:
//     host_statistics64 with HOST_VM_INFO64, a Mach call again.  The total and
//     the swap come from sysctl and are real; the breakdown is empty.
//   - resident memory per process and its command line: libproc
//     (proc_pidinfo, proc_pidpath) and KERN_PROCARGS2, which needs a pid
//     argument the standard library's Sysctl cannot pass.
//   - temperatures: the SMC, reached through IOKit, a framework rather than a
//     system call.
//
// The price of taking cgo for any of them is the same one every time, and it
// is the reason the answer is no: the release build is cross-compiled and
// byte-for-byte repeatable with CGO_ENABLED=0 (scripts/build-release.sh), and
// cgo ends both properties at once.
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
// sysctl, not through files, so the seam the Linux collector opens for its
// tests is not available here.  What is testable — every structure the kernel
// answers with — is tested in internal/darwinsys instead.
type Collector struct {
	// SampleWindow is accepted so `--sample` reads the same on every
	// system, and is not slept through here: the only fact it would time,
	// the CPU busy share, comes from a Mach call this build cannot make.
	SampleWindow time.Duration
	// Top is how many processes to list per ranking.
	Top int
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
// hopeful.  The alternatives are issuing the sysctl trap by hand, which Apple
// does not support from outside libSystem, and cgo.
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
	st := Status{TakenAt: time.Now().Format(time.RFC3339), Missing: map[string]string{}}

	c.host(&st)
	c.load(&st)
	c.processes(&st)
	c.memory(&st)

	disks, err := c.Disks()
	if err != nil {
		st.Missing["getfsstat"] = err.Error()
	}
	st.Disks = disks
	st.Network = c.network(&st)

	st.Sensors = nil
	st.Missing[FactSensors] = "температуру на маке отдаёт SMC через IOKit — фреймворк, а не системный вызов; из Go без cgo не читается"

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
		st.Missing[name] = err.Error()
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
		st.Missing["kern.boottime"] = err.Error()
		return
	}
	boot, ok := darwinsys.ParseTimeval(b)
	if !ok {
		st.Missing["kern.boottime"] = "ответ не разбирается как struct timeval — время работы не считаем"
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
		st.Missing["vm.loadavg"] = err.Error()
	} else if one, five, fifteen, ok := darwinsys.ParseLoadAvg(b); ok {
		st.Load.One, st.Load.Five, st.Load.Fifteen = one, five, fifteen
	} else {
		st.Missing["vm.loadavg"] = "ответ не разбирается как struct loadavg — средние не показываем"
	}
	// /proc/loadavg carries three more numbers macOS does not publish.
	st.Missing[FactLoadEntities] = "число задач в очереди и последний pid — это ключи /proc/loadavg; на маке их нет"

	for _, name := range []string{"hw.logicalcpu", "hw.ncpu"} {
		if n, err := syscall.SysctlUint32(name); err == nil && n > 0 {
			st.Load.CPUCount = int(n)
			break
		}
	}
	st.Missing[FactCPUBusy] = "доля занятого процессора живёт в host_statistics(HOST_CPU_LOAD_INFO) — это вызов Mach, а не системный; из Go без cgo не делается"
}

func (c Collector) memory(st *Status) {
	m := procfs.Memory{Present: map[string]bool{}}
	if b, err := sysctlRaw("hw.memsize"); err != nil {
		st.Missing["hw.memsize"] = err.Error()
	} else if v, ok := darwinsys.ParseUint64(b); ok {
		m.Total = v
		m.Present[procfs.FieldTotal] = true
	} else {
		st.Missing["hw.memsize"] = "ответ не восьмибайтовый — объём памяти не показываем"
	}

	if b, err := sysctlRaw("vm.swapusage"); err != nil {
		st.Missing["vm.swapusage"] = err.Error()
	} else if s, ok := darwinsys.ParseSwapUsage(b); ok {
		m.SwapTotal, m.SwapFree, m.SwapUsed = s.Total, s.Avail, s.Used
		m.Present[procfs.FieldSwapTotal] = true
		m.Present[procfs.FieldSwapFree] = true
		m.Present[procfs.FieldSwapUsed] = true
	} else {
		st.Missing["vm.swapusage"] = "ответ не разбирается как struct xsw_usage — своп не показываем"
	}
	st.Memory = m

	st.Missing[FactMemoryPages] = "свободная, активная, неактивная, сжатая и проводная память живёт в host_statistics64(HOST_VM_INFO64) — вызов Mach; из Go без cgo не делается. Всего и своп — измерены"
}

func (c Collector) processes(st *Status) {
	st.Missing[FactProcessRSS] = "занятая процессом память живёт в libproc (proc_pidinfo) и в порте задачи Mach — из Go без cgo не читается, поэтому десятки по памяти нет"
	st.Missing[FactProcessArgs] = "командная строка живёт в sysctl KERN_PROCARGS2, а он требует довода-pid, которого syscall.Sysctl передать не умеет; показываем короткое имя из kinfo_proc"
	st.Missing[FactThreads] = "число потоков процесса есть только в task_info Mach — в kinfo_proc его нет"
	st.Missing[FactBlocked] = "непрерываемый сон в состояниях kinfo_proc не различается: маковский SSTOP — это остановленный, а не заблокированный"

	b, err := sysctlGrowing("kern.proc.all", 3)
	if err != nil {
		st.Missing["kern.proc.all"] = err.Error()
		return
	}
	procs, ok := darwinsys.ParseProcs(b)
	if !ok {
		st.Missing["kern.proc.all"] = "ответ не делится на записи kinfo_proc — список процессов не публикуем"
		return
	}
	if !darwinsys.Verify(procs, os.Getpid(), os.Getppid(), os.Getuid()) {
		st.Missing["kern.proc.all"] = "самопроверка записи kinfo_proc не сошлась: свой pid, родитель и пользователь читаются не там, где мы их ждём. Числа из такой записи не публикуем"
		return
	}

	st.Processes.Total = len(procs)
	for _, p := range procs {
		if p.State == darwinsys.StateRunning {
			st.Processes.Running++
		}
	}

	// p_pctcpu is a leftover of BSD accounting, and a kernel that no longer
	// keeps it leaves it zero for everybody.  A column of zeros is not a
	// measurement, so it is published only when somebody in the list is
	// busy — and then only if the numbers stay inside what the machine can
	// physically do.
	if !darwinsys.AnyCPU(procs) {
		st.Missing[FactProcessCPU] = "ядро не заполняет p_pctcpu в kinfo_proc (ноль у всех процессов) — десятки по процессору нет"
		return
	}
	ceiling := 100 * float64(max(st.Load.CPUCount, 1))
	out := make([]Proc, 0, len(procs))
	for _, p := range procs {
		if p.PctCPU > ceiling+1 {
			st.Missing[FactProcessCPU] = "p_pctcpu показал больше, чем машина может занять всеми ядрами, — такой записи не верим и десятку по процессору не публикуем"
			return
		}
		pct := p.PctCPU
		out = append(out, Proc{
			PID: p.PID, PPID: p.PPID, UID: p.UID, State: p.State,
			Comm: p.Comm, CPUPercent: &pct,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if a, b := cpuOrZero(out[i]), cpuOrZero(out[j]); a != b {
			return a > b
		}
		return out[i].PID < out[j].PID
	})
	st.Processes.TopByCPU = head(out, c.Top)
	// Names are looked up for the listed processes only.  In a build
	// without cgo os/user reads /etc/passwd, which on macOS holds the
	// system accounts and not the human ones, so this often stays empty —
	// and an empty name prints as "—", not as a wrong one.
	for i := range st.Processes.TopByCPU {
		if u, err := user.LookupId(strconv.Itoa(st.Processes.TopByCPU[i].UID)); err == nil {
			st.Processes.TopByCPU[i].User = u.Username
		}
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
		st.Missing["net.Interfaces"] = err.Error()
		return nil
	}
	mtus := make(map[int]int, len(ifaces))
	for _, in := range ifaces {
		mtus[in.Index] = in.MTU
	}

	counters := map[int]darwinsys.IfCounters{}
	if b, err := syscall.RouteRIB(syscall.NET_RT_IFLIST2, 0); err != nil {
		st.Missing[FactNetCounters] = "список интерфейсов маршрутного сокета не прочитался: " + err.Error()
	} else if got := darwinsys.ParseIfList2(b); darwinsys.VerifyIfList(got, mtus) {
		counters = got
		st.Missing[FactNetTxDrops] = "в struct if_data64 нет счётчика отброшенных исходящих пакетов — tx_dropped остаётся пустым"
	} else {
		st.Missing[FactNetCounters] = "разбор struct if_data64 не сошёлся с MTU, который система называет сама, — счётчиков не публикуем"
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
