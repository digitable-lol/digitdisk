// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package sysinfo assembles a snapshot of the running Linux system out of
// /proc and /sys.  Anything the kernel does not expose is left empty (nil,
// zero-length, or a null in JSON) and named in Status.Missing — the snapshot
// never invents a value.
package sysinfo

import "digitdisk/internal/procfs"

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
	// Missing names the sources that could not be read, with the reason.
	Missing map[string]string `json:"missing,omitempty"`
}

// Host is identity and lifetime.
type Host struct {
	Hostname      string  `json:"hostname"`
	KernelRelease string  `json:"kernel_release"`
	KernelVersion string  `json:"kernel_version"`
	Machine       string  `json:"machine"`
	Distro        string  `json:"distro"`
	DistroID      string  `json:"distro_id"`
	DistroVersion string  `json:"distro_version"`
	BootTime      int64   `json:"boot_time_unix"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	UptimeHuman   string  `json:"uptime_human"`
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
	Total       int    `json:"total"`
	Running     int    `json:"running"`
	Blocked     int    `json:"blocked"`
	Threads     int    `json:"threads"`
	Unreadable  int    `json:"unreadable"`
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
