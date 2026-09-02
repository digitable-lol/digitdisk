// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"os"
	"syscall"
	"unsafe"
)

// ttyState is what the строка состояния has to know about the terminal
// before it draws anything: how big it is, and whether the command has taken
// it over.
type ttyState struct {
	Rows, Cols int
	// Cooked is false when the terminal is in raw mode — which is what vim,
	// less, ssh to an interactive shell and every other full-screen program
	// does the moment it starts.  The bar goes away while it is false: two
	// programs drawing on one screen is a mess neither of them can fix.
	Cooked bool
}

// ttyLook asks the terminal itself, with two ioctls and no child process.
//
// It is asked at every замер, five times a second, so that a full-screen
// program is noticed within a fifth of a second of starting rather than at
// the next redraw.  On Linux that costs a microsecond; the portable spelling
// — running stty(1) — costs a fork, and is what the other platforms do at the
// slower redraw rate instead.
func ttyLook(f *os.File) (ttyState, bool) {
	if f == nil {
		return ttyState{}, false
	}
	fd := f.Fd()
	var ws struct{ Rows, Cols, X, Y uint16 }
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws))); e != 0 {
		return ttyState{}, false
	}
	var t syscall.Termios
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(&t))); e != 0 {
		return ttyState{}, false
	}
	return ttyState{
		Rows:   int(ws.Rows),
		Cols:   int(ws.Cols),
		Cooked: t.Lflag&syscall.ICANON != 0 && t.Lflag&syscall.ECHO != 0,
	}, true
}
