// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package sysinfo assembles a snapshot of the running system.
//
// The snapshot has one shape everywhere; where the facts come from does not.
// Linux reads /proc and /sys (collect_linux.go), macOS asks sysctl and
// getfsstat (collect_darwin.go).  Everything else in this package — the types
// below, the durations, the interface addresses — is the same code on both.
//
// Anything the system does not publish is left empty: a nil pointer (JSON
// null), an empty list, a false flag in Memory.Present.  It is then named in
// Status.Missing with the reason.  A zero that was measured and a field nobody
// could measure must never look alike, and on a system that publishes less
// than Linux does that rule is most of the work.
package sysinfo

import (
	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/procfs"
)

// Status is the whole snapshot printed by `digitdisk status`.
type Status struct {
	TakenAt   string        `json:"taken_at"`
	Host      Host          `json:"host"`
	Load      Load          `json:"load"`
	Memory    procfs.Memory `json:"memory"`
	Processes Processes     `json:"processes"`
	Disks     []Disk        `json:"disks"`
	Network   []Iface       `json:"network"`
	Sensors   []Sensor      `json:"sensors"`
	// GPUs are the video cards, each with whatever its driver publishes
	// about it.  A card whose driver publishes nothing is still here: the
	// machine has it, and an absent line would say it does not.
	GPUs []gpuinfo.Card `json:"gpus"`
	// Missing names what is not in the snapshot and why: a source that
	// could not be read (keyed by its name, "meminfo", "mounts"), and a
	// fact this system does not publish at all (keyed by one of the Fact
	// constants below).  The reason is meant for a reader, and the report
	// prints it where the number would have been.
	Missing map[string]string `json:"missing,omitempty"`
}

// Host is identity and lifetime.
type Host struct {
	Hostname      string `json:"hostname"`
	KernelRelease string `json:"kernel_release"`
	KernelVersion string `json:"kernel_version"`
	Machine       string `json:"machine"`
	// Model is the machine's own name for itself, when it has one:
	// "MacBookPro18,3" from hw.model, "Dell Inc. PowerEdge C6525" from the
	// firmware tables on Linux.
	Model         string  `json:"model,omitempty"`
	Distro        string  `json:"distro"`
	DistroID      string  `json:"distro_id"`
	DistroVersion string  `json:"distro_version"`
	BootTime      int64   `json:"boot_time_unix"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	UptimeHuman   string  `json:"uptime_human"`
	// CPUModel is what the processor calls itself.
	CPUModel string `json:"cpu_model,omitempty"`
	// Bits is the word size of this build — 64 on every release target,
	// and the thing to look at when a binary turns up somewhere else.
	Bits int `json:"bits,omitempty"`
	// User, Shell, Desktop and Terminal are the session this snapshot was
	// taken in, not properties of the machine.  A server has no desktop and
	// says so with an empty field.
	User     string `json:"user,omitempty"`
	Shell    string `json:"shell,omitempty"`
	Desktop  string `json:"desktop,omitempty"`
	Terminal string `json:"terminal,omitempty"`
}

// Load is scheduler pressure plus a short CPU-busy sample.
type Load struct {
	procfs.LoadAvg
	CPUCount int `json:"cpu_count"`
	// BusyPercent is the share of non-idle CPU time over SampleMillis,
	// across all cores (0..100).  A nil pointer (JSON null) means the
	// sample was not taken — never confuse that with a measured zero.
	BusyPercent  *float64 `json:"busy_percent"`
	SampleMillis int64    `json:"sample_millis"`
	// Cores is the same share measured for each processor over the same
	// window.  It is empty when the system does not publish the
	// per-processor counters, or when the shares failed the check against
	// BusyPercent — see coresAgree.
	Cores []Core `json:"cores,omitempty"`
}

// Proc is one process in the snapshot.
type Proc struct {
	PID        int    `json:"pid"`
	PPID       int    `json:"ppid"`
	User       string `json:"user"`
	UID        int    `json:"uid"`
	State      string `json:"state"`
	Comm       string `json:"comm"`
	Cmdline    string `json:"cmdline"`
	Threads    int    `json:"threads"`
	RSSBytes   int64  `json:"rss_bytes"`
	VSizeBytes uint64 `json:"vsize_bytes"`
	// CPUPercent is nil (JSON null) when the process was not present for
	// the whole sample window, so it is never a measured zero.
	CPUPercent *float64 `json:"cpu_percent"`
}

// Processes is the process section.
type Processes struct {
	Total      int `json:"total"`
	Running    int `json:"running"`
	Blocked    int `json:"blocked"`
	Threads    int `json:"threads"`
	Unreadable int `json:"unreadable"`
	// WithDetail is how many of Total the memory, thread and processor
	// figures actually cover.  On Linux that is every process; on macOS the
	// kernel answers only about processes belonging to the same user, so a
	// snapshot taken without administrator rights covers a part of the list
	// and says how large a part.  Zero means the figures were not collected
	// at all, and the report prints dashes rather than a partial sum.
	WithDetail  int    `json:"with_detail"`
	TopByMemory []Proc `json:"top_by_memory"`
	TopByCPU    []Proc `json:"top_by_cpu"`
}

// Disk is one mounted filesystem with statfs numbers.
type Disk struct {
	Source         string  `json:"source"`
	MountPoint     string  `json:"mount_point"`
	FSType         string  `json:"fs_type"`
	ReadOnly       bool    `json:"read_only"`
	TotalBytes     uint64  `json:"total_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	AvailableBytes uint64  `json:"available_bytes"`
	UsePercent     float64 `json:"use_percent"`
	InodesTotal    uint64  `json:"inodes_total"`
	InodesFree     uint64  `json:"inodes_free"`
	Error          string  `json:"error,omitempty"`
}

// Iface is one network interface: counters from /proc/net/dev, the rest from
// /sys/class/net.
type Iface struct {
	procfs.NetCounters
	MAC       string   `json:"mac,omitempty"`
	OperState string   `json:"oper_state,omitempty"`
	MTU       int      `json:"mtu,omitempty"`
	SpeedMbit int      `json:"speed_mbit,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
}

// Sensor is one temperature reading from /sys/class/hwmon.
type Sensor struct {
	Chip    string  `json:"chip"`
	Label   string  `json:"label"`
	Celsius float64 `json:"celsius"`
	HighC   float64 `json:"high_celsius,omitempty"`
	CritC   float64 `json:"crit_celsius,omitempty"`
}
