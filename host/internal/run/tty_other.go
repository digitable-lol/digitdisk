// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !linux

package run

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ttyState is what the строка состояния has to know about the terminal.  See
// tty_linux.go for what each field means; only the way of asking differs.
type ttyState struct {
	Rows, Cols int
	Cooked     bool
}

// ttyLook asks stty(1), the way internal/ui already asks it for a size.
//
// One child process per question, so the caller asks at the redraw rate and
// not at the замер rate.  stty acts on its own standard input, which is the
// spelling both GNU and BSD understand without a flag for the file.
func ttyLook(f *os.File) (ttyState, bool) {
	if f == nil {
		return ttyState{}, false
	}
	cmd := exec.Command("stty", "-a")
	cmd.Stdin = f
	out, err := cmd.Output()
	if err != nil {
		return ttyState{}, false
	}
	return parseStty(string(out))
}

// parseStty reads the size and the two flags out of `stty -a`.
//
// Both spellings are accepted, because both are in use: GNU writes «rows 24;
// columns 80», BSD and macOS write «24 rows; 80 columns».  A flag that is off
// is written with a leading minus by both, and that is the whole of what is
// needed here: a terminal in raw mode has neither icanon nor echo.
func parseStty(text string) (ttyState, bool) {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == ';' || r == ','
	})
	st := ttyState{Cooked: true}
	for i, f := range fields {
		switch f {
		case "rows":
			st.Rows = numberNear(fields, i)
		case "columns":
			st.Cols = numberNear(fields, i)
		case "-icanon", "-echo":
			st.Cooked = false
		}
	}
	return st, st.Rows > 0 && st.Cols > 0
}

// numberNear takes the number written beside a word, on whichever side this
// stty put it.
func numberNear(fields []string, i int) int {
	if i+1 < len(fields) {
		if n, err := strconv.Atoi(fields[i+1]); err == nil {
			return n
		}
	}
	if i > 0 {
		if n, err := strconv.Atoi(fields[i-1]); err == nil {
			return n
		}
	}
	return 0
}
