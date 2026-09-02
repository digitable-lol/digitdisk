// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// A cgroup is a control group of our own: the exact way to count a tree of
// processes on Linux.
//
// The kernel does the counting, so nothing is missed — a compiler that lived
// and died between two замера is in cpu.stat all the same, and memory.peak is
// the peak of the group TOGETHER, with a page shared by ten processes counted
// once.  That is the difference between «6,1 ГиБ» and «около 6 ГиБ», and it is
// why this is tried first.
//
// It is tried, not required.  A control group can only be made where the
// machine has given this user a subtree to make it in: a desktop session
// under systemd usually has one, a container or a bare session usually does
// not, and asking for privileges to get one is not something a tool that
// looks at disks may do.  Everything here fails softly and the walk over
// /proc takes over.
type cgroup struct {
	dir      string
	memPeak  bool // the kernel publishes memory.peak here
	sampled  uint64
	procPeak int
}

// cgroupMount is where cgroup v2 is mounted on a system that has it.
const cgroupMount = "/sys/fs/cgroup"

// makeCgroup builds a group for this run, or returns nil.
//
// The search goes outward from our own group: the group we are in, then its
// parent, and so on to the mount.  The first directory that lets us create a
// subdirectory AND gives that subdirectory the two counters we came for wins.
// A directory that lets us create but cannot count is cleaned up behind us —
// an empty group left in cgroupfs is litter nobody knows to remove.
func makeCgroup(mount, own, name string) *cgroup {
	if mount == "" || own == "" {
		return nil
	}
	for dir := path.Join(mount, own); strings.HasPrefix(dir, mount); dir = path.Dir(dir) {
		made := path.Join(dir, name)
		if err := os.Mkdir(made, 0o755); err != nil {
			if dir == mount {
				return nil
			}
			continue
		}
		if c := ready(made); c != nil {
			return c
		}
		// It counts nothing; ask the parent to distribute the two
		// controllers we need, and look again.
		_ = os.WriteFile(path.Join(dir, "cgroup.subtree_control"), []byte("+cpu +memory"), 0o644)
		if c := ready(made); c != nil {
			return c
		}
		_ = rmdir(made)
		if dir == mount {
			return nil
		}
	}
	return nil
}

// ready reports whether a freshly made group actually counts.  A group whose
// parent never enabled the controllers has the files of a group and none of
// the numbers, and taking it would mean printing zeros as measurements.
func ready(dir string) *cgroup {
	if _, err := os.Stat(path.Join(dir, "cpu.stat")); err != nil {
		return nil
	}
	if _, err := os.Stat(path.Join(dir, "memory.current")); err != nil {
		return nil
	}
	_, peak := readNumber(path.Join(dir, "memory.peak"))
	return &cgroup{dir: dir, memPeak: peak}
}

// ownCgroup reads which group this process is in.  cgroup v2 writes one line,
// «0::/user.slice/…»; a v1-only machine writes several and none of them start
// with «0::», and this returns nothing for it — v1 cannot do what we came for.
func ownCgroup(text string) string {
	for _, line := range strings.Split(text, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "0::"); ok {
			return rest
		}
	}
	return ""
}

// parseCPUUsage pulls usage_usec out of cpu.stat: the CPU time everything in
// the group has used, alive or not.
func parseCPUUsage(text string) (time.Duration, bool) {
	for _, line := range strings.Split(text, "\n") {
		if rest, ok := strings.CutPrefix(line, "usage_usec "); ok {
			n, err := strconv.ParseUint(strings.TrimSpace(rest), 10, 64)
			if err != nil {
				return 0, false
			}
			return time.Duration(n) * time.Microsecond, true
		}
	}
	return 0, false
}

// parseNumber reads a one-number cgroup file.  «max» is not a number and not
// a zero: it is the absence of a limit, and a caller that treated it as a
// reading would print nonsense.
func parseNumber(text string) (uint64, bool) {
	s := strings.TrimSpace(text)
	if s == "" || s == "max" {
		return 0, false
	}
	n, err := strconv.ParseUint(s, 10, 64)
	return n, err == nil
}

func readNumber(file string) (uint64, bool) {
	b, err := os.ReadFile(file)
	if err != nil {
		return 0, false
	}
	return parseNumber(string(b))
}

// countPids counts the processes in the group.  pids.current is there only
// when the pids controller was distributed, and cgroup.procs is there always,
// so the second is what is counted.
func countPids(text string) int {
	n := 0
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}

// pids lists the processes of the group, for the question about video memory.
func (c *cgroup) pids() []int {
	b, err := os.ReadFile(path.Join(c.dir, "cgroup.procs"))
	if err != nil {
		return nil
	}
	var out []int
	for _, line := range strings.Split(string(b), "\n") {
		if pid, err := strconv.Atoi(strings.TrimSpace(line)); err == nil {
			out = append(out, pid)
		}
	}
	return out
}

// sample reads the three numbers the бар shows.
func (c *cgroup) sample(now time.Time, since time.Time, before time.Duration) reading {
	r := reading{Known: true}
	if b, err := os.ReadFile(path.Join(c.dir, "cpu.stat")); err == nil {
		if cpu, ok := parseCPUUsage(string(b)); ok {
			r.CPUSeconds = cpu.Seconds()
			if d := now.Sub(since).Seconds(); d > 0 && cpu >= before {
				r.Percent = (cpu - before).Seconds() / d * 100
			}
		}
	}
	if mem, ok := readNumber(path.Join(c.dir, "memory.current")); ok {
		r.Bytes = mem
		if mem > c.sampled {
			c.sampled = mem
		}
		r.Peak = c.sampled
	}
	if b, err := os.ReadFile(path.Join(c.dir, "cgroup.procs")); err == nil {
		r.Processes = countPids(string(b))
		if r.Processes > c.procPeak {
			c.procPeak = r.Processes
		}
	}
	return r
}

// totals reads what the kernel kept for the whole life of the group.
func (c *cgroup) totals() totals {
	t := totals{How: ByGroup, Procs: c.procPeak}
	if b, err := os.ReadFile(path.Join(c.dir, "cpu.stat")); err == nil {
		if cpu, ok := parseCPUUsage(string(b)); ok {
			t.CPU, t.CPUExact = cpu, true
		}
	}
	if peak, ok := readNumber(path.Join(c.dir, "memory.peak")); ok {
		t.Peak, t.PeakExact = peak, true
	} else {
		// An older kernel keeps no peak.  What was seen is what is
		// said, and it is not called exact.
		t.Peak = c.sampled
	}
	if peak, ok := readNumber(path.Join(c.dir, "pids.peak")); ok && int(peak) > t.Procs {
		t.Procs = int(peak)
	}
	return t
}

// rmdir removes an EMPTY DIRECTORY and can do nothing else.
//
// It is rmdir(2) and not os.Remove on purpose, and the purpose is the rule
// this tree lives by: стирать файлы вправе только host/internal/clean, where
// the приговор ядра, the сверка отпечатка, the корзина and the журнал stand
// around every such call.  Nothing of that applies to a directory in
// cgroupfs that this process created a moment ago and that holds no files at
// all — and rmdir cannot touch a file even by mistake, which is exactly why
// AGENTS.md names it as the way empty directories go.
func rmdir(dir string) error { return syscall.Rmdir(dir) }

// close removes the group.  A group the command left processes in cannot be
// removed, and they are not killed to make it removable: this tool does not
// end processes it did not start, and the ones it did start are the command's
// own business.
func (c *cgroup) close() {
	if c == nil || c.dir == "" {
		return
	}
	if rmdir(c.dir) == nil {
		return
	}
	time.Sleep(50 * time.Millisecond)
	_ = rmdir(c.dir)
}
