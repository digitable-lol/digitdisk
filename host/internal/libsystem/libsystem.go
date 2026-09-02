// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package libsystem calls the documented C functions of libSystem that macOS
// publishes its own state through, and returns their answers as bytes.
//
// # Why this package exists
//
// Some of what a Mac knows about itself is not behind sysctl(3).  The share of
// busy processor time and the page counts of the virtual memory system come
// from host_statistics; a process's resident memory, thread count and consumed
// processor time come from proc_pidinfo; and the sysctl node that answers with
// a process's command line takes the pid inside its name, which the standard
// library's Sysctl cannot express.  All four are ordinary C functions exported
// by /usr/lib/libSystem.B.dylib, and none of them is private.
//
// # Why it needs neither cgo nor a Mac to build
//
// The obvious way to call a C function from Go is cgo, and cgo would end the
// release: the four binaries are cross-compiled on one Linux machine and
// checked byte for byte against a second build of themselves, and cgo ends
// both properties at once.  So the calls are made the way the Go standard
// library itself makes them on macOS — the library records a dynamic import of
// the symbol, an assembly stub jumps to it, and the call goes out through
// syscall.syscall6.  The Go linker writes the import into the Mach-O file and
// the system loader binds it to libSystem at start-up, exactly as it binds the
// imports the runtime already needs.
//
// The result is that GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build still
// produces the whole program, on a machine that is not a Mac, and the release
// script is untouched.
//
// # What this package does not do
//
// It does not decode.  Every function here returns the bytes the kernel wrote
// and how many of them there are; the structures are decoded in package
// darwinsys, which has no build tag and is therefore tested on any machine.
// Nothing is published from those bytes until a decoder has checked them
// against something already known — see darwinsys for what each check is.
package libsystem
