// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"os"
	"strconv"
	"syscall"
	"time"
)

// Три названия учёта. They travel in the JSON as they are written here — a
// script that asks «was this measured exactly» compares bytes, not words —
// and the screen shows a translated word for each (see accountingWord).
const (
	// ByGroup is a control group of our own: the kernel counts the tree.
	ByGroup = "cgroup"
	// ByProc is the walk over /proc by parent links.
	ByProc = "proc"
	// ByKernel is wait4(2) rusage and nothing else: exact totals, no live
	// reading.  It is what every platform has.
	ByKernel = "rusage"
)

// A reading is what the бар shows: the tree as it is right now.
type reading struct {
	Known      bool
	Processes  int
	Bytes      uint64  // memory the tree holds now
	Peak       uint64  // the most it has held
	Percent    float64 // CPU over the last замер
	CPUSeconds float64 // CPU since the start, as far as we saw it
}

// totals is what the сводка prints.  Exact and approximate are kept apart on
// purpose: a number that was sampled must not be printed as if the kernel had
// said it.
type totals struct {
	CPU       time.Duration
	CPUExact  bool
	Peak      uint64
	PeakExact bool
	PeakOne   uint64 // largest single process, from the kernel
	Procs     int    // processes at the busiest moment
	Seen      int    // distinct processes noticed
	How       string
}

// rusageTotals reads what the kernel says about the child and everything the
// child waited for.
//
// This is the floor under every platform: ru_utime and ru_stime of a reaped
// child include the children it reaped itself, so the CPU time of a build is
// exact here even when nothing else is available.  Two things it does not
// say, and neither is guessed at elsewhere: a process that was orphaned
// instead of waited for is not in it, and ru_maxrss is the peak of the
// LARGEST SINGLE process, never the peak of the tree together.
func rusageTotals(ps *os.ProcessState) (cpu time.Duration, maxOne uint64, ok bool) {
	if ps == nil {
		return 0, 0, false
	}
	ru, ok := ps.SysUsage().(*syscall.Rusage)
	if !ok || ru == nil {
		return 0, 0, false
	}
	cpu = timeval(ru.Utime) + timeval(ru.Stime)
	return cpu, maxRSSBytes(ru), true
}

// ownPeak is the peak resident memory of digitdisk itself.
//
// It is not a curiosity: a child is forked from us before it execs, and until
// it execs its resident memory is OURS — the pages are shared and counted
// again.  So ru_maxrss of a command that ended quickly is our own footprint
// and not the command's, and a peak below it cannot be told apart from that.
// The tool says «не измерен» there rather than printing its own size as if it
// were somebody else's.
func ownPeak() uint64 {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return 0
	}
	return maxRSSBytes(&ru)
}

func timeval(t syscall.Timeval) time.Duration {
	return time.Duration(t.Sec)*time.Second + time.Duration(t.Usec)*time.Microsecond
}

// signalNames are the signals a wrapped command actually dies of, spelled the
// way a person searches for them.  Go's own name for SIGINT is «interrupt»,
// which is true and unfindable; a number nobody recognises is worse.  Anything
// not here is printed as its number by the caller.
var signalNames = map[syscall.Signal]string{
	syscall.SIGHUP:  "SIGHUP",
	syscall.SIGINT:  "SIGINT",
	syscall.SIGQUIT: "SIGQUIT",
	syscall.SIGILL:  "SIGILL",
	syscall.SIGABRT: "SIGABRT",
	syscall.SIGFPE:  "SIGFPE",
	syscall.SIGKILL: "SIGKILL",
	syscall.SIGSEGV: "SIGSEGV",
	syscall.SIGPIPE: "SIGPIPE",
	syscall.SIGALRM: "SIGALRM",
	syscall.SIGTERM: "SIGTERM",
	syscall.SIGBUS:  "SIGBUS",
	syscall.SIGUSR1: "SIGUSR1",
	syscall.SIGUSR2: "SIGUSR2",
	syscall.SIGXCPU: "SIGXCPU",
	syscall.SIGXFSZ: "SIGXFSZ",
}

// signalName is the name of the signal that killed the command.
func signalName(s syscall.Signal) string {
	if name, ok := signalNames[s]; ok {
		return name
	}
	return "SIG" + strconv.Itoa(int(s))
}
