// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package gpuinfo

// Which PROCESSES hold video memory — the one per-process question about a
// video card that has an honest answer.
//
// The NVIDIA driver publishes, through its own program, how much memory each
// compute process holds.  It does NOT publish how much of the card's time
// each process took: `utilization.gpu` is the whole card's share, and pinning
// it on one command would be a lie the moment somebody else's job is running
// beside ours — which, on a shared machine with a card in it, is the normal
// case rather than the exception.  So memory by process is read here and load
// by process is not read at all, and the сводка says as much in words instead
// of printing a number that means something else.
//
// The rule of nvsmi.go holds here without change: nothing runs unless the
// reader turned Tool on, and everything that arrives this way is marked as
// having come from outside.

import (
	"strconv"
	"strings"
)

// A ComputeApp is one process the driver says is holding video memory.
type ComputeApp struct {
	PID   int
	Bytes uint64
}

// computeAppsQuery is the documented machine-readable form: pid and memory,
// no header, no units — the unit is fixed by the query and is MiB.
var computeAppsQuery = []string{
	"--query-compute-apps=pid,used_gpu_memory",
	"--format=csv,noheader,nounits",
}

// ComputeApps asks which processes hold video memory.  ok is false when the
// reader may not run the program, when the program is not there, and when it
// answered something that is not a list of processes.
func (r Reader) ComputeApps() (apps []ComputeApp, ok bool) {
	if !r.Tool {
		return nil, false
	}
	out, err := r.run("nvidia-smi", computeAppsQuery...)
	if err != nil {
		return nil, false
	}
	return ParseComputeApps(string(out)), true
}

// ParseComputeApps decodes «pid, MiB» lines.  A line that is not two numbers
// is dropped rather than guessed at: «[N/A]» is what the program writes when
// it cannot tell, and [N/A] must never become a zero.
func ParseComputeApps(text string) []ComputeApp {
	var out []ComputeApp
	for _, line := range strings.Split(text, "\n") {
		f := strings.Split(line, ",")
		if len(f) < 2 {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(f[0]))
		if err != nil || pid <= 0 {
			continue
		}
		mib, err := strconv.ParseFloat(strings.TrimSpace(f[1]), 64)
		if err != nil || mib < 0 {
			continue
		}
		out = append(out, ComputeApp{PID: pid, Bytes: uint64(mib) * 1024 * 1024})
	}
	return out
}

// HasNVIDIA reports whether any of these cards is one nvidia-smi could speak
// about.  A machine without one never starts the program.
func HasNVIDIA(cards []Card) bool {
	for _, c := range cards {
		if c.Driver == "nvidia" || c.VendorID == "10de" {
			return true
		}
	}
	return false
}
