// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package iokit reads the IORegistry — the tree macOS keeps its devices in —
// and returns what it finds as Go values.
//
// # Why this package exists
//
// A Mac says nothing about its video cards through sysctl(3).  What it knows
// lives in the IORegistry and is published by IOKit, and IOKit answers with
// Core Foundation objects rather than with numbers.  Until this package the
// tree said so and stopped there: internal/sysinfo/collect_darwin.go declared
// the video cards unmeasurable without cgo.
//
// # Why it needs neither cgo nor a Mac to build
//
// The same door internal/libsystem opens for Mach: the library records a
// dynamic import of each symbol, an assembly stub jumps to it, and the call
// goes out through syscall.syscall6.  The only difference is the library the
// import names — /System/Library/Frameworks/IOKit.framework/… and its
// CoreFoundation twin instead of /usr/lib/libSystem.B.dylib.  The Go linker
// writes one LC_LOAD_DYLIB per library and binds each symbol to the right one;
// see TestДверьОткрывается for the measurement that says so.
//
// # What is deliberately not imported
//
// Only functions.  kIOMainPortDefault, kCFAllocatorDefault and the CFNumber
// type codes are variables and constants, not exported functions, and a
// variable is not a symbol this door opens — the same wall internal/libsystem
// met at mach_task_self().  All three are zero or a documented literal, so
// they are written here as zero and as literals, and nothing is guessed.
package iokit
