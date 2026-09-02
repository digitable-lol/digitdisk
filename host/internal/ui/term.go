// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package ui draws the live screen of `digitdisk status`.  It shows the same
// readings the printed report shows and reads no source of its own: the screen
// is a second view of sysinfo.Status, never a second collector.
//
// The terminal is driven by hand with ANSI sequences and stty(1), so the tool
// keeps its two build promises — no cgo and no third-party code — on all four
// release targets.  Nothing here is required to work: every entry point
// degrades to "this is not a terminal", and the caller then prints text.
package ui

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// IsTerminal reports whether f is a terminal we may draw on.
//
// The character-device bit alone is not enough: /dev/null carries it too, and
// `digitdisk status >/dev/null` must stay printed text.  The second question —
// does this file descriptor answer a terminal query — is asked with `stty
// size`, which fails on /dev/null and on a pipe.
func IsTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	_, _, ok := Size(f)
	return ok
}

// UsableTERM reports whether $TERM names a terminal that can place the cursor.
// An empty TERM and the "dumb" terminal both mean it cannot, and the screen
// must not be attempted there.
func UsableTERM() bool {
	t := os.Getenv("TERM")
	return t != "" && t != "dumb"
}

// Size asks the terminal behind f for its size in rows and columns.  ok is
// false when f is not a terminal, when stty is missing, or when it answers
// something that is not two numbers.
func Size(f *os.File) (rows, cols int, ok bool) {
	out, err := sttyOn(f, "size")
	if err != nil {
		return 0, 0, false
	}
	fields := strings.Fields(string(out))
	if len(fields) != 2 {
		return 0, 0, false
	}
	r, err1 := strconv.Atoi(fields[0])
	c, err2 := strconv.Atoi(fields[1])
	if err1 != nil || err2 != nil || r <= 0 || c <= 0 {
		return 0, 0, false
	}
	return r, c, true
}

// Raw puts the terminal behind f into raw mode with echo off and returns the
// call that puts it back.  The previous settings are saved with `stty -g` and
// handed back to the same stty, so no assumption is made about the format of
// that string — it travels from one stty to the same stty.
//
// The restore function is safe to call more than once and never panics: a
// terminal left in raw mode is the worst outcome this package can produce, so
// restoring is attempted on every path, including a signal and a panic.
func Raw(f *os.File) (restore func(), err error) {
	saved, err := sttyOn(f, "-g")
	if err != nil {
		return nil, err
	}
	state := strings.TrimSpace(string(saved))
	if _, err := sttyOn(f, "raw", "-echo"); err != nil {
		return nil, err
	}
	done := false
	return func() {
		if done {
			return
		}
		done = true
		_, _ = sttyOn(f, state)
	}, nil
}

// sttyOn runs stty against the terminal behind f.  stty acts on its own
// standard input, which is the portable spelling: GNU stty spells the same
// thing `-F file` and BSD stty spells it `-f file`, and neither is needed when
// the descriptor is simply handed over.
func sttyOn(f *os.File, args ...string) ([]byte, error) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = f
	cmd.Stderr = nil
	return cmd.Output()
}
