// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build linux

package sysinfo

import (
	"os"
	"path/filepath"
	"testing"
)

// fakeRoot builds a miniature /proc, /sys and /etc out of captured text, so
// the collector is exercised without touching the running kernel.
func fakeRoot(t *testing.T) (dir, mountPoint string) {
	t.Helper()
	dir = t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("proc/sys/kernel/hostname", "testnode\n")
	write("proc/sys/kernel/osrelease", "9.9.9-test\n")
	write("proc/version", "Linux version 9.9.9-test\n")
	write("proc/uptime", "90061.00 1000.00\n")
	write("proc/loadavg", "1.50 2.50 3.50 7/900 4242\n")
	write("proc/stat", "cpu  100 0 100 800 0 0 0 0 0 0\ncpu0 50 0 50 400 0 0 0 0 0 0\ncpu1 50 0 50 400 0 0 0 0 0 0\nbtime 1700000000\nprocesses 5\nprocs_running 2\nprocs_blocked 1\n")
	write("proc/meminfo", "MemTotal:       1024 kB\nMemFree:         256 kB\nMemAvailable:    512 kB\nBuffers:          64 kB\nCached:          128 kB\nShmem:            16 kB\nSReclaimable:     32 kB\nSwapTotal:       200 kB\nSwapFree:        150 kB\n")
	write("proc/net/dev", "Inter-|   Receive |  Transmit\n face |bytes packets errs drop fifo frame compressed multicast|bytes packets errs drop fifo colls carrier compressed\n  eth9: 10 2 0 0 0 0 0 0 20 4 0 0 0 0 0 0\n")
	// the mount point must exist for statfs to succeed
	mountPoint = filepath.Join(dir, "mnt")
	if err := os.MkdirAll(mountPoint, 0o755); err != nil {
		t.Fatal(err)
	}
	write("proc/mounts", "/dev/fake "+mountPoint+" ext4 rw,relatime 0 0\n/dev/gone /definitely/not/here ext4 ro 0 0\n")

	write("proc/1234/stat", "1234 (pretend) S 1 1 1 0 -1 0 0 0 0 0 5 6 0 0 20 0 3 0 1000 4096000 250 0 0 0 0 0 0 0 0 0 0 0 0 0\n")
	write("proc/1234/cmdline", "pretend\x00--flag\x00")
	write("proc/1234/status", "Name:\tpretend\nUid:\t1001\t1001\t1001\t1001\n")
	write("proc/not-a-pid/stat", "junk\n")

	write("sys/class/net/eth9/address", "aa:bb:cc:dd:ee:ff\n")
	write("sys/class/net/eth9/operstate", "up\n")
	write("sys/class/net/eth9/mtu", "9000\n")
	write("sys/class/hwmon/hwmon0/name", "faketemp\n")
	write("sys/class/hwmon/hwmon0/temp1_input", "42500\n")
	write("sys/class/hwmon/hwmon0/temp1_label", "Package\n")
	write("sys/class/hwmon/hwmon0/temp1_crit", "100000\n")

	write("etc/os-release", "PRETTY_NAME=\"Testix 1.0\"\nID=testix\nVERSION_ID=\"1.0\"\n")
	return dir, mountPoint
}

func newTestCollector(t *testing.T) (Collector, string) {
	dir, mp := fakeRoot(t)
	return Collector{
		Proc: filepath.Join(dir, "proc"),
		Sys:  filepath.Join(dir, "sys"),
		Etc:  filepath.Join(dir, "etc"),
		Top:  5,
	}, mp
}

func TestCollectFromFakeRoot(t *testing.T) {
	c, mountPoint := newTestCollector(t)
	st := c.Collect()

	if st.Host.Hostname != "testnode" || st.Host.KernelRelease != "9.9.9-test" {
		t.Errorf("host = %+v", st.Host)
	}
	if st.Host.Distro != "Testix 1.0" || st.Host.DistroID != "testix" {
		t.Errorf("distro = %q/%q", st.Host.Distro, st.Host.DistroID)
	}
	if st.Host.UptimeSeconds != 90061 || st.Host.UptimeHuman != "1д 01:01" {
		t.Errorf("uptime = %v %q", st.Host.UptimeSeconds, st.Host.UptimeHuman)
	}
	if st.Host.BootTime != 1700000000 {
		t.Errorf("boot time = %d", st.Host.BootTime)
	}
	if st.Load.One != 1.5 || st.Load.CPUCount != 2 {
		t.Errorf("load = %+v", st.Load)
	}
	if st.Load.BusyPercent != nil {
		t.Errorf("with no sample window the busy share must stay null, got %v", *st.Load.BusyPercent)
	}
	if st.Memory.Total != 1024*1024 || st.Memory.Used != (1024-512)*1024 {
		t.Errorf("memory = %+v", st.Memory)
	}

	if st.Processes.Total != 1 {
		t.Fatalf("processes = %d, want 1 (the non-numeric directory is not a pid)", st.Processes.Total)
	}
	if st.Processes.Running != 2 || st.Processes.Blocked != 1 {
		t.Errorf("running/blocked = %d/%d", st.Processes.Running, st.Processes.Blocked)
	}
	p := st.Processes.TopByMemory[0]
	if p.PID != 1234 || p.Comm != "pretend" || p.Cmdline != "pretend --flag" {
		t.Errorf("process = %+v", p)
	}
	if p.UID != 1001 || p.Threads != 3 {
		t.Errorf("uid/threads = %d/%d", p.UID, p.Threads)
	}
	if want := int64(250) * int64(os.Getpagesize()); p.RSSBytes != want {
		t.Errorf("rss = %d, want %d", p.RSSBytes, want)
	}
	if p.CPUPercent != nil {
		t.Errorf("with no sample window CPU%% must stay null, got %v", *p.CPUPercent)
	}

	if len(st.Disks) != 2 {
		t.Fatalf("disks = %d, want 2", len(st.Disks))
	}
	var mounted, broken *Disk
	for i := range st.Disks {
		if st.Disks[i].MountPoint == mountPoint {
			mounted = &st.Disks[i]
		} else {
			broken = &st.Disks[i]
		}
	}
	if mounted == nil || mounted.TotalBytes == 0 {
		t.Errorf("the existing mount point produced no statfs numbers: %+v", st.Disks)
	}
	if broken == nil || broken.Error == "" {
		t.Errorf("an unstattable mount must carry its error, not fake numbers: %+v", broken)
	}
	if broken != nil && !broken.ReadOnly {
		t.Errorf("the ro option was not picked up: %+v", broken)
	}

	if len(st.Network) != 1 {
		t.Fatalf("interfaces = %d, want 1", len(st.Network))
	}
	n := st.Network[0]
	if n.Name != "eth9" || n.RxBytes != 10 || n.TxBytes != 20 || n.MAC != "aa:bb:cc:dd:ee:ff" || n.MTU != 9000 || n.OperState != "up" {
		t.Errorf("interface = %+v", n)
	}
	if n.SpeedMbit != 0 {
		t.Errorf("an absent speed file must stay zero, got %d", n.SpeedMbit)
	}

	if len(st.Sensors) != 1 {
		t.Fatalf("sensors = %d, want 1", len(st.Sensors))
	}
	if s := st.Sensors[0]; s.Chip != "faketemp" || s.Label != "Package" || s.Celsius != 42.5 || s.CritC != 100 {
		t.Errorf("sensor = %+v", s)
	}
}

func TestCollectMissingSourcesAreNamedNotInvented(t *testing.T) {
	c := Collector{Proc: "/nonexistent/proc", Sys: "/nonexistent/sys", Etc: "/nonexistent/etc", Top: 3}
	st := c.Collect()
	for _, key := range []string{"meminfo", "stat", "loadavg", "uptime", "os-release", "mounts"} {
		if _, ok := st.Missing[key]; !ok {
			t.Errorf("source %q went missing without being reported: %v", key, st.Missing)
		}
	}
	if st.Memory.Total != 0 || st.Processes.Total != 0 || len(st.Disks) != 0 {
		t.Errorf("nothing readable, yet the snapshot has content: %+v", st)
	}
}

func TestSensorsAbsentIsEmptyNotError(t *testing.T) {
	c := Collector{Sys: filepath.Join(t.TempDir(), "sys")}
	if got := c.Sensors(); len(got) != 0 {
		t.Errorf("a machine without hwmon must report no sensors, got %+v", got)
	}
}
