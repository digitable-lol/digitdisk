// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build linux

// The Linux half of the collector: everything here reads /proc and /sys.  The
// parsers it hands the text to live in internal/procfs and are built
// everywhere, because a parser of captured text is testable on any machine;
// only this file, which names the paths and calls the Linux kernel, is behind
// the build tag.

package sysinfo

import (
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"digitdisk/internal/procfs"
)

// Collector reads a snapshot from a set of filesystem roots.  The roots are
// fields so tests can point them at a captured tree instead of the live
// kernel.
//
// Every system has its own Collector with its own fields; what they share is
// the two knobs the command line sets (SampleWindow, Top), New, and Collect.
type Collector struct {
	Proc string // usually /proc
	Sys  string // usually /sys
	Etc  string // usually /etc
	// SampleWindow is how long the CPU delta sample runs.  Zero disables
	// sampling (CPU percentages then come back as -1).
	SampleWindow time.Duration
	// Top is how many processes to list per ranking.
	Top int
}

// New returns a Collector pointed at the live system.
func New() Collector {
	return Collector{Proc: "/proc", Sys: "/sys", Etc: "/etc", SampleWindow: 200 * time.Millisecond, Top: 10}
}

func (c Collector) read(parts ...string) (string, error) {
	b, err := os.ReadFile(filepath.Join(parts...))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// trimmed reads a one-line pseudo-file, returning "" when unreadable.
func (c Collector) trimmed(parts ...string) string {
	s, err := c.read(parts...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(s)
}

// Collect gathers the whole snapshot.  Sources that fail are recorded in
// Status.Missing rather than aborting the run.
func (c Collector) Collect() Status {
	st := Status{TakenAt: time.Now().Format(time.RFC3339), Missing: map[string]string{}}
	miss := func(what string, err error) {
		if err != nil {
			st.Missing[what] = err.Error()
		}
	}

	// --- host -------------------------------------------------------
	st.Host.Hostname = c.trimmed(c.Proc, "sys/kernel/hostname")
	st.Host.KernelRelease = c.trimmed(c.Proc, "sys/kernel/osrelease")
	st.Host.KernelVersion = c.trimmed(c.Proc, "version")
	var un syscall.Utsname
	if err := syscall.Uname(&un); err == nil {
		st.Host.Machine = utsString(un.Machine[:])
	} else {
		miss("uname", err)
	}
	if text, err := c.read(c.Etc, "os-release"); err == nil {
		kv := procfs.ParseOSRelease(text)
		st.Host.Distro = kv["PRETTY_NAME"]
		if st.Host.Distro == "" {
			st.Host.Distro = kv["NAME"]
		}
		st.Host.DistroID = kv["ID"]
		st.Host.DistroVersion = kv["VERSION_ID"]
	} else {
		miss("os-release", err)
	}
	if text, err := c.read(c.Proc, "uptime"); err == nil {
		if up, _, ok := procfs.ParseUptime(text); ok {
			st.Host.UptimeSeconds = up
			st.Host.UptimeHuman = HumanDuration(up)
		}
	} else {
		miss("uptime", err)
	}

	// --- load, cpu --------------------------------------------------
	statText, statErr := c.read(c.Proc, "stat")
	miss("stat", statErr)
	first := procfs.ParseStat(statText)
	st.Host.BootTime = first.BootTime
	st.Load.CPUCount = len(first.PerCPU)
	if text, err := c.read(c.Proc, "loadavg"); err == nil {
		if la, ok := procfs.ParseLoadAvg(text); ok {
			st.Load.LoadAvg = la
		}
	} else {
		miss("loadavg", err)
	}

	// --- processes: two passes around the sample window --------------
	pass1, _ := c.snapProcesses()
	if c.SampleWindow > 0 {
		time.Sleep(c.SampleWindow)
	}
	pass2, unreadable := c.snapProcesses()
	elapsed := c.SampleWindow.Seconds()

	if statErr == nil && c.SampleWindow > 0 {
		if text, err := c.read(c.Proc, "stat"); err == nil {
			second := procfs.ParseStat(text)
			dTotal := float64(second.CPU.Total()) - float64(first.CPU.Total())
			dBusy := float64(second.CPU.Busy()) - float64(first.CPU.Busy())
			if dTotal > 0 {
				busy := 100 * dBusy / dTotal
				st.Load.BusyPercent = &busy
			}
			st.Load.SampleMillis = c.SampleWindow.Milliseconds()
			st.Processes.Running = second.ProcsRunning
			st.Processes.Blocked = second.ProcsBlocked
		}
	} else {
		st.Processes.Running = first.ProcsRunning
		st.Processes.Blocked = first.ProcsBlocked
	}

	procs := make([]Proc, 0, len(pass2))
	for pid, p := range pass2 {
		if elapsed > 0 {
			if before, ok := pass1[pid]; ok && before.start == p.start {
				window := p.at.Sub(before.at).Seconds()
				dTicks := float64(p.ticks) - float64(before.ticks)
				// USER_HZ is 100 for the procfs interface, so
				// ticks/100 is seconds of CPU time; 100% is one
				// fully busy core.
				if window > 0 && dTicks >= 0 {
					pct := dTicks / 100 * 100 / window
					p.proc.CPUPercent = &pct
				}
			}
		}
		procs = append(procs, p.proc)
		st.Processes.Threads += p.proc.Threads
	}
	st.Processes.Total = len(procs)
	st.Processes.Unreadable = unreadable
	// Every process /proc let us open was read whole: Linux publishes a
	// process's memory and threads to anybody who can list it.
	st.Processes.WithDetail = len(procs)

	byMem := append([]Proc(nil), procs...)
	sort.Slice(byMem, func(i, j int) bool {
		if byMem[i].RSSBytes != byMem[j].RSSBytes {
			return byMem[i].RSSBytes > byMem[j].RSSBytes
		}
		return byMem[i].PID < byMem[j].PID
	})
	byCPU := append([]Proc(nil), procs...)
	sort.Slice(byCPU, func(i, j int) bool {
		a, b := cpuOrZero(byCPU[i]), cpuOrZero(byCPU[j])
		if a != b {
			return a > b
		}
		return byCPU[i].PID < byCPU[j].PID
	})
	st.Processes.TopByMemory = head(byMem, c.Top)
	st.Processes.TopByCPU = head(byCPU, c.Top)

	// --- memory ------------------------------------------------------
	if text, err := c.read(c.Proc, "meminfo"); err == nil {
		st.Memory = procfs.ParseMeminfo(text)
	} else {
		miss("meminfo", err)
	}

	// --- disks, network, sensors -------------------------------------
	disks, err := c.Disks()
	miss("mounts", err)
	st.Disks = disks
	st.Network, err = c.Network()
	miss("net/dev", err)
	st.Sensors = c.Sensors()
	if len(st.Sensors) == 0 {
		// An absent sensor is a fact about the machine, not a failed
		// read — but "no temperatures" and "temperatures not looked
		// for" are still different things, so it is named.
		st.Missing[FactSensors] = "в " + filepath.Join(c.Sys, "class/hwmon") + " датчиков нет"
	}

	if len(st.Missing) == 0 {
		st.Missing = nil
	}
	return st
}

type procSample struct {
	proc  Proc
	ticks uint64
	start uint64
	// at is when this process's stat line was actually read.  A pass over
	// several thousand PIDs is not instantaneous, so the CPU delta must be
	// divided by each process's own interval, not by the nominal window.
	at time.Time
}

// snapProcesses walks /proc for numeric entries.  Processes that vanish or
// deny access mid-walk are counted, never fatal.
func (c Collector) snapProcesses() (map[int]procSample, int) {
	out := make(map[int]procSample, 512)
	unreadable := 0
	entries, err := os.ReadDir(c.Proc)
	if err != nil {
		return out, unreadable
	}
	pageSize := int64(os.Getpagesize())
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid <= 0 {
			continue
		}
		text, err := c.read(c.Proc, e.Name(), "stat")
		if err != nil {
			unreadable++
			continue
		}
		ps, ok := procfs.ParseProcStat(text)
		if !ok {
			unreadable++
			continue
		}
		p := Proc{
			PID:        ps.PID,
			PPID:       ps.PPID,
			State:      ps.State,
			Comm:       ps.Comm,
			Threads:    ps.NumThreads,
			RSSBytes:   ps.RSSPages * pageSize,
			VSizeBytes: ps.VSize,
			UID:        -1,
		}
		if raw, err := os.ReadFile(filepath.Join(c.Proc, e.Name(), "cmdline")); err == nil {
			p.Cmdline = procfs.ParseCmdline(raw)
		}
		if text, err := c.read(c.Proc, e.Name(), "status"); err == nil {
			if uid, ok := procfs.ParseStatusUID(text); ok {
				p.UID = uid
				if u, err := user.LookupId(strconv.Itoa(uid)); err == nil {
					p.User = u.Username
				}
			}
		}
		out[pid] = procSample{proc: p, ticks: ps.UTime + ps.STime, start: ps.StartTime, at: time.Now()}
	}
	return out, unreadable
}

// Disks lists mounted filesystems that report a non-zero size, which is the
// same rule df(1) applies to hide purely synthetic mounts.
func (c Collector) Disks() ([]Disk, error) {
	text, err := c.read(c.Proc, "mounts")
	if err != nil {
		return nil, err
	}
	var out []Disk
	for _, m := range procfs.ParseMounts(text) {
		d := Disk{Source: m.Source, MountPoint: m.Point, FSType: m.FSType}
		for _, opt := range strings.Split(m.Options, ",") {
			if opt == "ro" {
				d.ReadOnly = true
			}
		}
		var fs syscall.Statfs_t
		if err := syscall.Statfs(m.Point, &fs); err != nil {
			d.Error = err.Error()
			out = append(out, d)
			continue
		}
		if fs.Blocks == 0 {
			continue
		}
		bs := uint64(fs.Bsize)
		d.TotalBytes = fs.Blocks * bs
		d.UsedBytes = (fs.Blocks - fs.Bfree) * bs
		d.AvailableBytes = fs.Bavail * bs
		if d.UsedBytes+d.AvailableBytes > 0 {
			d.UsePercent = 100 * float64(d.UsedBytes) / float64(d.UsedBytes+d.AvailableBytes)
		}
		d.InodesTotal = fs.Files
		d.InodesFree = fs.Ffree
		out = append(out, d)
	}
	return out, nil
}

// Network merges /proc/net/dev counters with the per-interface attributes
// exposed under /sys/class/net.
func (c Collector) Network() ([]Iface, error) {
	text, err := c.read(c.Proc, "net/dev")
	if err != nil {
		return nil, err
	}
	addrs := interfaceAddresses()
	var out []Iface
	for _, nc := range procfs.ParseNetDev(text) {
		i := Iface{NetCounters: nc}
		base := filepath.Join(c.Sys, "class/net", nc.Name)
		i.MAC = c.trimmed(base, "address")
		i.OperState = c.trimmed(base, "operstate")
		if v := c.trimmed(base, "mtu"); v != "" {
			i.MTU, _ = strconv.Atoi(v)
		}
		if v := c.trimmed(base, "speed"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				i.SpeedMbit = n
			}
		}
		i.Addresses = addrs[nc.Name]
		out = append(out, i)
	}
	return out, nil
}

// Sensors reads every temperature input under /sys/class/hwmon.  Missing or
// unreadable chips are skipped silently: absent sensors are normal.
func (c Collector) Sensors() []Sensor {
	dirs, err := filepath.Glob(filepath.Join(c.Sys, "class/hwmon/hwmon*"))
	if err != nil {
		return nil
	}
	var out []Sensor
	for _, dir := range dirs {
		chip := c.trimmed(dir, "name")
		if chip == "" {
			chip = filepath.Base(dir)
		}
		inputs, _ := filepath.Glob(filepath.Join(dir, "temp*_input"))
		sort.Strings(inputs)
		for _, in := range inputs {
			milli := c.trimmed(in)
			if milli == "" {
				continue
			}
			n, err := strconv.ParseFloat(milli, 64)
			if err != nil {
				continue
			}
			prefix := strings.TrimSuffix(in, "_input")
			s := Sensor{Chip: chip, Celsius: n / 1000}
			s.Label = c.trimmed(prefix + "_label")
			if s.Label == "" {
				s.Label = filepath.Base(prefix)
			}
			if v := c.trimmed(prefix + "_max"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					s.HighC = f / 1000
				}
			}
			if v := c.trimmed(prefix + "_crit"); v != "" {
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					s.CritC = f / 1000
				}
			}
			out = append(out, s)
		}
	}
	return out
}

func utsString[T ~int8 | ~uint8](b []T) string {
	out := make([]byte, 0, len(b))
	for _, c := range b {
		if c == 0 {
			break
		}
		out = append(out, byte(c))
	}
	return string(out)
}
