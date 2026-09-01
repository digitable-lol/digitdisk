// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package darwinsys decodes the binary structures macOS hands back through
// sysctl(3), the way package procfs parses the text Linux hands back through
// /proc.
//
// Nothing here calls a kernel.  Every function takes bytes and returns values,
// so the layouts are exercised by tests on any machine — including the Linux
// one this tree is written on, where no macOS is available to run against.
// The code that does call sysctl lives in internal/sysinfo behind
// //go:build darwin and is a thin shell around these functions.
//
// # Layouts are asserted, not assumed
//
// The offsets below are read off the XNU headers named at each structure.  We
// have no Mac to check them against, so every decoder validates what it can —
// the size of the buffer, a version byte, a field whose value the caller
// already knows from somewhere else — and returns ok = false instead of
// numbers it cannot vouch for.  A false answer costs an empty field in the
// report, which the report prints as "—"; a wrong answer would be a number
// that looks measured and is not.  That trade is the whole point of this
// package.
//
// # Byte order
//
// Little-endian, with no switch on the architecture: macOS runs on amd64 and
// arm64 only, and both are little-endian.  A big-endian macOS would need this
// package changed, and would announce itself by failing every validation here.
package darwinsys

import "encoding/binary"

// Padded returns b as exactly size bytes, restoring the zero bytes a caller
// may have lost on the way here, and reports whether that was possible.
//
// Callers lose bytes because syscall.Sysctl is written for string-valued
// nodes: it returns the raw bytes of a sysctl and strips ONE trailing NUL.  A
// numeric node whose top byte happens to be zero — hw.memsize on any machine
// with less than 64 PiB of RAM, kern.boottime with its four bytes of padding —
// therefore arrives one byte short.  Appending zeros restores the value
// exactly on a little-endian machine, and macOS has no other kind.
//
// Two bytes short is not the trim; it is a different structure, and is
// refused.  A longer buffer is taken by its prefix: a kernel that grew a
// structure by appending fields still answers the fields we read.
func Padded(b []byte, size int) ([]byte, bool) {
	switch {
	case size <= 0:
		return nil, false
	case len(b) >= size:
		return b[:size], true
	case len(b) == size-1:
		out := make([]byte, size)
		copy(out, b)
		return out, true
	}
	return nil, false
}

// u16, u32, u64 and i32 read one little-endian field at a byte offset.  The
// caller has already checked the length; these panic on a short buffer rather
// than return a zero that would look like a measurement.
func u16(b []byte, off int) uint16 { return binary.LittleEndian.Uint16(b[off:]) }
func u32(b []byte, off int) uint32 { return binary.LittleEndian.Uint32(b[off:]) }
func u64(b []byte, off int) uint64 { return binary.LittleEndian.Uint64(b[off:]) }
func i32(b []byte, off int) int32  { return int32(binary.LittleEndian.Uint32(b[off:])) }

// cstring reads a NUL-terminated char[] field, which the kernel does not
// promise to terminate when the name fills the field exactly.
func cstring(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
