// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"time"

	"digitdisk/internal/gpuinfo"
)

// GPU is what can honestly be said about a video card and this command: the
// memory the command's own processes held on it.
//
// Not the load.  The driver publishes a share of the card's time for the CARD,
// not for a process, and attributing the whole card's load to the command
// would be false whenever anything else is using the card — which on a machine
// with a card in it is the ordinary case.  The field that is missing is
// missing on purpose, and the сводка says so in words.
type GPU struct {
	PeakBytes uint64 `json:"пик_памяти_байт"`
	Processes int    `json:"процессов_с_памятью"`
	Source    string `json:"источник"`
	// ByProcessLoad is false and stays false: it marks, for a machine
	// reader, that no per-process load was measured rather than that it was
	// measured as zero.
	ByProcessLoad bool `json:"загрузка_по_процессам"`
}

// gpuAsk is how often the driver's program is run.  It is a fork and a wait,
// so it happens on its own slow clock rather than at every замер.
const gpuAsk = 2 * time.Second

type gpuWatch struct {
	on     bool
	reader gpuinfo.Reader
	asked  time.Time
	ever   bool
	held   uint64
	peak   uint64
	procs  int
}

type gpuReading struct {
	Known bool
	Bytes uint64
}

// newGPU prepares the question, and asks nothing.  Without the ключ nothing
// here runs at all; with it, the cards are read from files first, and the
// program is left alone on a machine that has no card it could answer about.
func newGPU(on bool) *gpuWatch {
	if !on {
		return &gpuWatch{}
	}
	r := gpuinfo.New()
	r.Tool = true
	if !gpuinfo.HasNVIDIA(r.Read().Cards) {
		return &gpuWatch{}
	}
	return &gpuWatch{on: true, reader: r}
}

// sample asks, at most every gpuAsk, how much video memory the command's own
// processes hold.
func (g *gpuWatch) sample(now time.Time, pids []int) {
	if !g.on || now.Sub(g.asked) < gpuAsk {
		return
	}
	g.asked = now
	apps, ok := g.reader.ComputeApps()
	if !ok {
		return
	}
	g.ever = true
	mine := map[int]bool{}
	for _, pid := range pids {
		mine[pid] = true
	}
	var sum uint64
	n := 0
	for _, a := range apps {
		if !mine[a.PID] {
			continue
		}
		sum += a.Bytes
		n++
	}
	g.held = sum
	if sum > g.peak {
		g.peak, g.procs = sum, n
	}
}

func (g *gpuWatch) reading() gpuReading { return gpuReading{Known: g.ever, Bytes: g.held} }

// total is what the сводка prints, or nothing at all when the question was
// never allowed or never answered.
func (g *gpuWatch) total() *GPU {
	if !g.on || !g.ever {
		return nil
	}
	return &GPU{PeakBytes: g.peak, Processes: g.procs, Source: "nvidia-smi"}
}
