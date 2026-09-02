// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"fmt"
	"os"
	"strings"

	"digitdisk/internal/ui"
)

// bar is the строка состояния: one line at the bottom of the terminal that
// belongs to us, while every other line goes on belonging to the command.
//
// # Why a scroll region and not a carriage return
//
// The command writes to the same terminal we do.  A bar drawn with \r at the
// current cursor position would land in the middle of the command's own
// output and be scrolled away by it — which is what makes most such wrappers
// unusable with anything that prints.  So the terminal is asked to keep the
// last line out of the scrolling area (DECSTBM): the command's output scrolls
// in the rows above and never touches the row we draw on.
//
// The region is re-declared on every draw rather than once.  A command that
// resets the terminal — `clear`, `tput reset`, a program that exits badly —
// takes our region with it, and re-declaring costs six bytes and heals it.
//
// # Three cases where the bar must not exist
//
//	no terminal   Standard error is a pipe or a file: NOTHING is written.
//	              `digitdisk run make | tee log` must give what `make | tee
//	              log` gives, byte for byte, and it does — the command's
//	              descriptors are ours, unchanged, and we never write to
//	              standard output at all.
//	--plain       Asked not to.
//	raw mode      The command took the terminal for itself — vim, less, ssh
//	              to a shell, anything full-screen.  The region goes back to
//	              the whole screen and the line is erased, and it comes back
//	              when the command gives the terminal back.
type bar struct {
	tty        *os.File
	off        bool
	colour     bool
	rows, cols int
	set        bool
}

// ANSI, spelled once.  DECSC/DECRC (ESC 7, ESC 8) save and restore the cursor
// around everything we do, because the cursor belongs to the command.
const (
	save      = "\x1b7"
	restore   = "\x1b8"
	upOneLine = "\x1b[A"
	eraseLine = "\x1b[K"
	invert    = "\x1b[7m"
	plainAgn  = "\x1b[0m"
)

// newBar reserves the last line — before the command starts, while the cursor
// is still where the shell left it.  Making room later, in the middle of a
// half-written line of the command's own output, is what turns «Compiling…»
// into « done».
func newBar(err *os.File, plain bool, colour bool) *bar {
	b := &bar{tty: err}
	if plain || err == nil || !ui.Available(err) {
		b.off = true
		return b
	}
	st, ok := ttyLook(err)
	if !ok || !st.Cooked || st.Rows < 3 || st.Cols < 20 {
		b.off = true
		return b
	}
	b.rows, b.cols, b.colour = st.Rows, st.Cols, colour
	b.setup()
	return b
}

// width is how much room a line of text has.
func (b *bar) width() int {
	if b.off {
		return 0
	}
	return b.cols
}

// setup makes room for the bar and declares the region.
func (b *bar) setup() {
	// A newline first: if the cursor is on the last row this scrolls
	// everything up one line, and the row we are about to reserve is then
	// empty rather than holding somebody's text.  The cursor is put back
	// where it was.
	b.write(fmt.Sprintf("\n%s%s\x1b[1;%dr%s", upOneLine, save, b.rows-1, restore))
	b.set = true
}

// draw puts one line of text on the reserved row.
func (b *bar) draw(text string) {
	if b.off || text == "" {
		return
	}
	st, ok := ttyLook(b.tty)
	if !ok {
		// The terminal went away mid-run.  Nothing more is written to it.
		b.off = true
		return
	}
	if !st.Cooked {
		b.hide()
		return
	}
	if st.Rows != b.rows || st.Cols != b.cols {
		b.rows, b.cols = st.Rows, st.Cols
		b.set = false
	}
	if b.rows < 3 || b.cols < 20 {
		b.hide()
		return
	}
	if !b.set {
		b.setup()
	}
	var w strings.Builder
	w.WriteString(save)
	fmt.Fprintf(&w, "\x1b[1;%dr", b.rows-1)
	fmt.Fprintf(&w, "\x1b[%d;1H%s", b.rows, eraseLine)
	if b.colour {
		w.WriteString(invert)
	}
	w.WriteString(clip(text, b.cols))
	if b.colour {
		w.WriteString(plainAgn)
	}
	w.WriteString(restore)
	b.write(w.String())
}

// hide gives the whole screen back and erases the line, and leaves the bar
// able to come back.
func (b *bar) hide() {
	if b.off || !b.set {
		return
	}
	b.set = false
	b.write(fmt.Sprintf("%s\x1b[1;%dr\x1b[%d;1H%s%s", save, b.rows, b.rows, eraseLine, restore))
}

// wake is called after the terminal may have been reset behind our back —
// after Ctrl-Z and back, for one.  The next draw declares everything again.
func (b *bar) wake() { b.set = false }

// resize is called on SIGWINCH; the next draw asks the terminal how big it is
// now.
func (b *bar) resize() { b.set = false }

// close puts the terminal back the way it was found.
func (b *bar) close() { b.hide() }

// write is best-effort: a terminal that has gone away is not an error worth
// telling anybody about, and it is certainly not worth failing the command
// over.
func (b *bar) write(s string) {
	if b.tty != nil {
		_, _ = b.tty.WriteString(s)
	}
}

// clip cuts a line to fit the terminal, counting runes rather than bytes: a
// Russian line cut by bytes is cut in the middle of a letter.
func clip(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// colourOK reports whether this terminal may be coloured at all.  NO_COLOR is
// honoured the way internal/ui honours it: the bar still appears, it is simply
// not inverted.
func colourOK() bool {
	if _, off := os.LookupEnv("NO_COLOR"); off {
		return false
	}
	t := os.Getenv("TERM")
	return t != "" && t != "dumb"
}
