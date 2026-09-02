// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"os"
	"sort"
	"strconv"
	"time"

	"digitdisk/internal/procfs"
)

// procTree follows the whole tree of a command through /proc, by parent
// links.
//
// # Why not read every process every time
//
// `npm ci` and `make -j` are trees of dozens of processes, and the machine
// they run on may hold four thousand others.  Reading /proc/<pid>/stat for
// all of them takes about a tenth of a second here, five times a second would
// be half a core, and an обёртка that eats half a core is measuring itself.
//
// So the tree is followed rather than searched.  Every pid that already
// existed when the command started is a NON-member by construction — the
// command did not exist yet — so the first замер reads nothing but the
// command itself.  After that only two things are read: the members (for
// their CPU and memory) and the pids that APPEARED since the last замер (to
// see whose children they are).  Listing /proc is a directory read and costs
// about four milliseconds even here; the reads are the expensive part, and
// there are as many of them as the tree is wide.
//
// A pid that disappears is forgotten, so a pid number reused later is read
// again as new.  That is what keeps the parent links honest over a long
// build: pid numbers wrap, and a stale «not ours» would quietly drop a whole
// branch.
//
// # What this cannot see, and says so
//
//   - A process that lives and dies between two замера is not counted at all.
//   - The CPU of a process that died is counted up to the last замер that saw
//     it, not to its death.
//   - A process reparented away from the tree (double fork) leaves it.
//   - RSS summed over a tree counts a shared page once per process holding
//     it, so the пик is an upper bound, not the memory the tree really held.
//
// None of that applies to a контрольная группа, which is why one is preferred
// when the machine gives us one; and the сводка names which of the two
// answered, so an approximate number is never read as an exact one.
type procTree struct {
	dir  string
	root int

	// known is every pid that exists and has been classified: true for a
	// member of the tree, false for somebody else's process.
	known map[int]bool
	// cpu is the ticks each living member had at the last замер.
	cpu map[int]uint64
	// gone is the ticks of the members that have since exited.  Their last
	// reading is kept: a build whose compilers all finished must not report
	// less CPU than it did a second ago.
	gone uint64

	seen      int
	peakProcs int
	peakBytes uint64

	lastTicks uint64
	lastAt    time.Time
	started   bool
}

// userHZ is what /proc counts CPU in.  It is 100 for the procfs interface
// whatever the kernel was built with — the same constant internal/procfs
// names in its own comment.
const userHZ = 100

// newProcTree prepares to follow root.  Everything alive right now is
// somebody else's: root has not been started yet.
func newProcTree(dir string) *procTree {
	t := &procTree{dir: dir, known: map[int]bool{}, cpu: map[int]uint64{}}
	for _, pid := range t.list() {
		t.known[pid] = false
	}
	return t
}

// start names the process whose tree this is.
func (t *procTree) start(pid int, now time.Time) {
	t.root = pid
	t.known[pid] = true
	t.seen = 1
	t.lastAt = now
	t.started = true
}

// list reads the pids /proc holds right now.  A directory read, no files.
func (t *procTree) list() []int {
	d, err := os.Open(t.dir)
	if err != nil {
		return nil
	}
	names, err := d.Readdirnames(-1)
	d.Close()
	if err != nil {
		return nil
	}
	pids := make([]int, 0, len(names))
	for _, name := range names {
		if pid, err := strconv.Atoi(name); err == nil {
			pids = append(pids, pid)
		}
	}
	return pids
}

// stat reads one process.  A pid that has just exited reads as «not there»,
// which is not an error: this is a directory of processes, and processes end.
func (t *procTree) stat(pid int) (procfs.ProcStat, bool) {
	b, err := os.ReadFile(t.dir + "/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return procfs.ProcStat{}, false
	}
	return procfs.ParseProcStat(string(b))
}

// sample reads the tree once.
func (t *procTree) sample(now time.Time) reading {
	if !t.started {
		return reading{}
	}
	live := t.list()
	if len(live) == 0 {
		return t.last(now)
	}
	here := make(map[int]bool, len(live))
	for _, pid := range live {
		here[pid] = true
	}

	// Processes that ended keep the CPU they had been seen to use.
	for pid, member := range t.known {
		if here[pid] {
			continue
		}
		if member {
			t.gone += t.cpu[pid]
			delete(t.cpu, pid)
		}
		delete(t.known, pid)
	}

	// Processes that appeared are classified by their parent.  A chain of
	// them may have appeared at once — a shell that forked a make that
	// forked a compiler — so the pass repeats while it keeps finding
	// members, and stops as soon as it stops finding any.
	fresh := map[int]procfs.ProcStat{}
	for _, pid := range live {
		if _, ok := t.known[pid]; ok {
			continue
		}
		if st, ok := t.stat(pid); ok {
			fresh[pid] = st
		} else {
			// Gone already: it is nobody's business now, and it must
			// not be remembered as «not ours» — the pid may come
			// back as one of ours.
			continue
		}
	}
	for grew := true; grew; {
		grew = false
		for pid, st := range fresh {
			if _, done := t.known[pid]; done {
				continue
			}
			if parent, ok := t.known[st.PPID]; ok && parent {
				t.known[pid] = true
				t.seen++
				grew = true
			}
		}
	}
	for pid := range fresh {
		if _, done := t.known[pid]; !done {
			t.known[pid] = false
		}
	}

	// The members are read: ticks and pages.
	var ticks, bytes uint64
	procs := 0
	page := uint64(os.Getpagesize())
	for pid, member := range t.known {
		if !member {
			continue
		}
		st, ok := fresh[pid]
		if !ok {
			if st, ok = t.stat(pid); !ok {
				continue
			}
		}
		procs++
		used := st.UTime + st.STime
		t.cpu[pid] = used
		ticks += used
		if st.RSSPages > 0 {
			bytes += uint64(st.RSSPages) * page
		}
	}
	ticks += t.gone

	if procs > t.peakProcs {
		t.peakProcs = procs
	}
	if bytes > t.peakBytes {
		t.peakBytes = bytes
	}
	r := reading{Known: true, Processes: procs, Bytes: bytes, Peak: t.peakBytes}
	r.CPUSeconds = float64(ticks) / userHZ
	if d := now.Sub(t.lastAt).Seconds(); d > 0 && ticks >= t.lastTicks {
		r.Percent = float64(ticks-t.lastTicks) / userHZ / d * 100
	}
	t.lastTicks, t.lastAt = ticks, now
	return r
}

// last is the reading when /proc could not be listed at all.
func (t *procTree) last(now time.Time) reading {
	return reading{Known: true, CPUSeconds: float64(t.lastTicks) / userHZ}
}

// members lists the pids of the tree, for the one thing that needs pids: the
// question about video memory, which the driver's program answers by pid.
func (t *procTree) members() []int {
	out := make([]int, 0, len(t.known))
	for pid, member := range t.known {
		if member {
			out = append(out, pid)
		}
	}
	sort.Ints(out)
	return out
}
