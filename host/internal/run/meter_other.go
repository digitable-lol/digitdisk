// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !linux

package run

import (
	"os"
	"syscall"
	"time"
)

// meter outside Linux is the kernel's own total and nothing else.
//
// wait4(2) gives the exact CPU time of the command and everything it waited
// for, and the peak of its largest single process; it gives nothing while the
// command runs, and nothing about the tree together.  Both facts are said out
// loud in the сводка rather than filled in with a plausible number: there is
// no /proc here and no control group, and inventing a tree from `ps` output
// five times a second would be a worse answer more expensively.
type meter struct{}

func newMeter() *meter { return &meter{} }

func (m *meter) attr() *syscall.SysProcAttr { return nil }

func (m *meter) dropGroup() {}

func (m *meter) started(pid int, now time.Time) {}

func (m *meter) sample(now time.Time) reading { return reading{} }

func (m *meter) members() []int { return nil }

func (m *meter) finish(ps *os.ProcessState) totals {
	cpu, one, ok := rusageTotals(ps)
	if one <= ownPeak() {
		one = 0 // indistinguishable from our own footprint: see ownPeak
	}
	return totals{How: ByKernel, CPU: cpu, CPUExact: ok, PeakOne: one}
}

func (m *meter) close() {}
