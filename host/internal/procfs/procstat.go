// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package procfs

import (
	"strconv"
	"strings"
)

// ProcStat is the subset of /proc/<pid>/stat digitdisk needs.  Ticks are in
// USER_HZ (always 100 for the procfs interface), RSS is in pages.
type ProcStat struct {
	PID        int    `json:"pid"`
	Comm       string `json:"comm"`
	State      string `json:"state"`
	PPID       int    `json:"ppid"`
	UTime      uint64 `json:"utime_ticks"`
	STime      uint64 `json:"stime_ticks"`
	NumThreads int    `json:"threads"`
	StartTime  uint64 `json:"start_ticks"`
	VSize      uint64 `json:"vsize_bytes"`
	RSSPages   int64  `json:"rss_pages"`
}

// ParseProcStat parses one /proc/<pid>/stat line.  The second field (comm) is
// wrapped in parentheses and may itself contain spaces and parentheses, so the
// split is anchored on the LAST ')' in the line.
func ParseProcStat(text string) (ProcStat, bool) {
	text = strings.TrimRight(text, "\n")
	open := strings.IndexByte(text, '(')
	end := strings.LastIndexByte(text, ')')
	if open < 0 || end < 0 || end < open {
		return ProcStat{}, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(text[:open]))
	if err != nil {
		return ProcStat{}, false
	}
	p := ProcStat{PID: pid, Comm: text[open+1 : end]}

	// f[0] is field 3 of the manual page (state), f[i] is field i+3.
	f := strings.Fields(text[end+1:])
	if len(f) < 22 {
		return ProcStat{}, false
	}
	u := func(i int) uint64 {
		n, _ := strconv.ParseUint(f[i], 10, 64)
		return n
	}
	p.State = f[0]
	p.PPID, _ = strconv.Atoi(f[1])
	p.UTime = u(11)                           // field 14
	p.STime = u(12)                           // field 15
	p.NumThreads, _ = strconv.Atoi(f[17])     // field 20
	p.StartTime = u(19)                       // field 22
	p.VSize = u(20)                           // field 23
	rss, _ := strconv.ParseInt(f[21], 10, 64) // field 24
	p.RSSPages = rss
	return p, true
}

// ParseCmdline turns the NUL-separated /proc/<pid>/cmdline blob into a single
// display string.  An empty blob (kernel thread) yields "".
func ParseCmdline(raw []byte) string {
	parts := strings.FieldsFunc(string(raw), func(r rune) bool { return r == 0 })
	return strings.Join(parts, " ")
}

// ParseStatusUID pulls the real UID out of /proc/<pid>/status ("Uid: r e s fs").
func ParseStatusUID(text string) (int, bool) {
	for _, line := range strings.Split(text, "\n") {
		if rest, ok := strings.CutPrefix(line, "Uid:"); ok {
			f := strings.Fields(rest)
			if len(f) == 0 {
				return 0, false
			}
			uid, err := strconv.Atoi(f[0])
			return uid, err == nil
		}
	}
	return 0, false
}
