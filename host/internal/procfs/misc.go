// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package procfs

import (
	"strconv"
	"strings"
)

// LoadAvg mirrors /proc/loadavg.
type LoadAvg struct {
	One           float64 `json:"1min"`
	Five          float64 `json:"5min"`
	Fifteen       float64 `json:"15min"`
	Runnable      int     `json:"runnable"`
	TotalEntities int     `json:"total_entities"`
	LastPID       int     `json:"last_pid"`
}

// ParseLoadAvg parses "0.44 0.55 0.60 2/1234 56789".
func ParseLoadAvg(text string) (LoadAvg, bool) {
	f := strings.Fields(text)
	if len(f) < 3 {
		return LoadAvg{}, false
	}
	var la LoadAvg
	var err error
	if la.One, err = strconv.ParseFloat(f[0], 64); err != nil {
		return LoadAvg{}, false
	}
	if la.Five, err = strconv.ParseFloat(f[1], 64); err != nil {
		return LoadAvg{}, false
	}
	if la.Fifteen, err = strconv.ParseFloat(f[2], 64); err != nil {
		return LoadAvg{}, false
	}
	if len(f) >= 4 {
		if run, tot, ok := strings.Cut(f[3], "/"); ok {
			la.Runnable, _ = strconv.Atoi(run)
			la.TotalEntities, _ = strconv.Atoi(tot)
		}
	}
	if len(f) >= 5 {
		la.LastPID, _ = strconv.Atoi(f[4])
	}
	return la, true
}

// ParseUptime parses /proc/uptime and returns (uptime seconds, idle seconds).
func ParseUptime(text string) (up float64, idle float64, ok bool) {
	f := strings.Fields(text)
	if len(f) < 1 {
		return 0, 0, false
	}
	var err error
	if up, err = strconv.ParseFloat(f[0], 64); err != nil {
		return 0, 0, false
	}
	if len(f) >= 2 {
		idle, _ = strconv.ParseFloat(f[1], 64)
	}
	return up, idle, true
}

// CPUTimes is one "cpu"/"cpuN" line of /proc/stat, in USER_HZ ticks.
type CPUTimes struct {
	Name      string `json:"name"`
	User      uint64 `json:"user"`
	Nice      uint64 `json:"nice"`
	System    uint64 `json:"system"`
	Idle      uint64 `json:"idle"`
	IOWait    uint64 `json:"iowait"`
	IRQ       uint64 `json:"irq"`
	SoftIRQ   uint64 `json:"softirq"`
	Steal     uint64 `json:"steal"`
	Guest     uint64 `json:"guest"`
	GuestNice uint64 `json:"guest_nice"`
}

// Total is the sum of every accounted tick on the line.  Guest ticks are
// already included in User/Nice by the kernel and are therefore not re-added.
func (c CPUTimes) Total() uint64 {
	return c.User + c.Nice + c.System + c.Idle + c.IOWait + c.IRQ + c.SoftIRQ + c.Steal
}

// Busy is Total minus idle-ish ticks.
func (c CPUTimes) Busy() uint64 {
	return c.Total() - c.Idle - c.IOWait
}

// Stat is the digest of /proc/stat that digitdisk uses.
type Stat struct {
	CPU            CPUTimes   `json:"cpu"`
	PerCPU         []CPUTimes `json:"-"`
	BootTime       int64      `json:"boot_time"`
	ForksSinceBoot uint64     `json:"forks_since_boot"`
	ProcsRunning   int        `json:"procs_running"`
	ProcsBlocked   int        `json:"procs_blocked"`
}

// ParseStat parses /proc/stat.
func ParseStat(text string) Stat {
	var s Stat
	for _, line := range strings.Split(text, "\n") {
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		switch {
		case strings.HasPrefix(f[0], "cpu"):
			c := CPUTimes{Name: f[0]}
			nums := make([]uint64, 0, 10)
			for _, v := range f[1:] {
				n, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					break
				}
				nums = append(nums, n)
			}
			get := func(i int) uint64 {
				if i < len(nums) {
					return nums[i]
				}
				return 0
			}
			c.User, c.Nice, c.System, c.Idle = get(0), get(1), get(2), get(3)
			c.IOWait, c.IRQ, c.SoftIRQ, c.Steal = get(4), get(5), get(6), get(7)
			c.Guest, c.GuestNice = get(8), get(9)
			if f[0] == "cpu" {
				s.CPU = c
			} else {
				s.PerCPU = append(s.PerCPU, c)
			}
		case f[0] == "btime" && len(f) > 1:
			s.BootTime, _ = strconv.ParseInt(f[1], 10, 64)
		case f[0] == "processes" && len(f) > 1:
			s.ForksSinceBoot, _ = strconv.ParseUint(f[1], 10, 64)
		case f[0] == "procs_running" && len(f) > 1:
			s.ProcsRunning, _ = strconv.Atoi(f[1])
		case f[0] == "procs_blocked" && len(f) > 1:
			s.ProcsBlocked, _ = strconv.Atoi(f[1])
		}
	}
	return s
}

// ParseOSRelease parses the shell-ish KEY=VALUE format of /etc/os-release.
func ParseOSRelease(text string) map[string]string {
	out := make(map[string]string)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') && val[len(val)-1] == val[0] {
			val = val[1 : len(val)-1]
		}
		out[strings.TrimSpace(key)] = val
	}
	return out
}

// ParseCPUModel picks the processor's own name out of /proc/cpuinfo.
//
// The key is not the same on every architecture: x86 and its kin write "model
// name", ARM writes "Processor" or nothing at all and leaves the name to
// "Hardware", and a machine that writes none of them has no name to give.
// The first line that matches wins, because every processor in the file is
// the same one repeated.
func ParseCPUModel(text string) string {
	for _, key := range []string{"model name", "Model Name", "cpu model", "Processor", "Hardware", "cpu"} {
		for _, line := range strings.Split(text, "\n") {
			name, value, ok := strings.Cut(line, ":")
			if !ok || strings.TrimSpace(name) != key {
				continue
			}
			if v := strings.TrimSpace(value); v != "" {
				return v
			}
		}
	}
	return ""
}
