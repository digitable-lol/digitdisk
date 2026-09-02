// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"os"
	"strconv"
	"syscall"
	"time"
)

// meter counts what a command's whole tree of processes costs.
//
// On Linux there are two ways, and the better one is tried first: a
// контрольная группа of our own, where the kernel counts, and a walk over
// /proc, where we count.  Which one answered is not a detail — it is printed
// in the сводка, because one of them is exact and the other is a series of
// glances.
type meter struct {
	cg   *cgroup
	fd   *os.File
	tree *procTree

	begin   time.Time
	lastAt  time.Time
	lastCPU time.Duration
}

// newMeter prepares the accounting BEFORE the command starts.  Both ways need
// that: a control group has to exist before there is a process to put in it,
// and the walk over /proc needs to know which pids were already there.
func newMeter() *meter {
	m := &meter{}
	if own := ownCgroup(readFile("/proc/self/cgroup")); own != "" {
		if c := makeCgroup(cgroupMount, own, "digitdisk-"+strconv.Itoa(os.Getpid())); c != nil {
			if fd, err := os.Open(c.dir); err == nil {
				m.cg, m.fd = c, fd
				return m
			}
			c.close()
		}
	}
	m.tree = newProcTree("/proc")
	return m
}

// attr is how the command must be created.  clone3(2) with CLONE_INTO_CGROUP
// puts the child INTO the group at the moment it is born, so there is no
// window in which it could fork a child outside the accounting — writing a
// pid into cgroup.procs after the fact leaves exactly that window.
func (m *meter) attr() *syscall.SysProcAttr {
	if m.fd == nil {
		return nil
	}
	return &syscall.SysProcAttr{UseCgroupFD: true, CgroupFD: int(m.fd.Fd())}
}

// dropGroup gives up on the control group.  It is called when the kernel
// refused to start the command inside it — an old kernel without clone3, a
// group that turned out not to be ours — and the walk over /proc takes over.
// Still before the command runs, so the walk's «everything alive now is
// somebody else's» is as true as it was.
func (m *meter) dropGroup() {
	if m.fd != nil {
		m.fd.Close()
		m.fd = nil
	}
	if m.cg != nil {
		m.cg.close()
		m.cg = nil
	}
	if m.tree == nil {
		m.tree = newProcTree("/proc")
	}
}

func (m *meter) started(pid int, now time.Time) {
	m.begin, m.lastAt = now, now
	if m.tree != nil {
		m.tree.start(pid, now)
	}
}

func (m *meter) sample(now time.Time) reading {
	if m.cg != nil {
		r := m.cg.sample(now, m.lastAt, m.lastCPU)
		m.lastAt = now
		m.lastCPU = time.Duration(r.CPUSeconds * float64(time.Second))
		return r
	}
	if m.tree != nil {
		return m.tree.sample(now)
	}
	return reading{}
}

// members lists the processes of the tree right now.
func (m *meter) members() []int {
	if m.cg != nil {
		return m.cg.pids()
	}
	if m.tree != nil {
		return m.tree.members()
	}
	return nil
}

func (m *meter) finish(ps *os.ProcessState) totals {
	cpu, one, haveRusage := rusageTotals(ps)
	if one <= ownPeak() {
		one = 0 // indistinguishable from our own footprint: see ownPeak
	}
	if m.cg != nil {
		t := m.cg.totals()
		t.PeakOne = one
		if !t.CPUExact && haveRusage {
			t.CPU, t.CPUExact = cpu, true
		}
		return t
	}
	t := totals{How: ByProc, Peak: 0, PeakOne: one}
	if m.tree != nil {
		t.Peak = m.tree.peakBytes
		t.Procs = m.tree.peakProcs
		t.Seen = m.tree.seen
		t.CPU = time.Duration(float64(m.tree.lastTicks) / userHZ * float64(time.Second))
	}
	if haveRusage {
		// The kernel's own total beats anything we watched: it holds
		// every process the command waited for, including the ones that
		// began and ended between two замера.
		t.CPU, t.CPUExact = cpu, true
	}
	return t
}

func (m *meter) close() {
	if m.fd != nil {
		m.fd.Close()
		m.fd = nil
	}
	if m.cg != nil {
		m.cg.close()
		m.cg = nil
	}
}

func readFile(name string) string {
	b, err := os.ReadFile(name)
	if err != nil {
		return ""
	}
	return string(b)
}
