// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

// The MIB of the sysctl node that answers with one process's arguments
// (<sys/sysctl.h>).  Unlike every other node the collector reads, it takes the
// pid as part of the name, which is why it needs a sysctl call that passes a
// whole MIB instead of a string.
const (
	CtlKern   = 1  // CTL_KERN
	ProcArgs2 = 49 // KERN_PROCARGS2
)

// ArgsMIB is the name of that node for one process.
func ArgsMIB(pid int) []int32 { return []int32{CtlKern, ProcArgs2, int32(pid)} }

// maxArgs is a ceiling on the argument count taken from the buffer's first
// four bytes.  It is not a limit anyone runs into — the kernel caps a whole
// argument block at a megabyte — but a misread length field is usually a huge
// number, and stopping on it is cheaper than allocating for it.
const maxArgs = 1 << 20

// ParseProcArgs2 decodes the argument block of one process.
//
// The block starts with a 32-bit argument count, then the path the kernel
// executed, then padding, then that many NUL-terminated strings, then the
// environment, which this function stops before: a report that printed
// somebody's environment would be leaking, not measuring.
//
// It answers ok only when exactly argc strings were found.  A block that runs
// out early is a block whose shape we do not understand, and half an argument
// list read out of it would look like a command line that was never run.
func ParseProcArgs2(b []byte) (exe string, argv []string, ok bool) {
	if len(b) < 4 {
		return "", nil, false
	}
	argc := int(i32(b, 0))
	if argc < 0 || argc > maxArgs {
		return "", nil, false
	}
	rest := b[4:]

	// The executable path, then the NUL bytes the kernel pads it with to
	// align what follows.
	end := indexNUL(rest)
	if end < 0 {
		return "", nil, false
	}
	exe = string(rest[:end])
	rest = rest[end:]
	for len(rest) > 0 && rest[0] == 0 {
		rest = rest[1:]
	}

	argv = make([]string, 0, argc)
	for len(argv) < argc {
		end = indexNUL(rest)
		if end < 0 {
			return "", nil, false
		}
		argv = append(argv, string(rest[:end]))
		rest = rest[end+1:]
	}
	return exe, argv, true
}

func indexNUL(b []byte) int {
	for i, c := range b {
		if c == 0 {
			return i
		}
	}
	return -1
}

// SameArgv reports whether a decoded argument list is the one the caller knows
// it was started with.
//
// This is the self-check for the whole block: run the decoder on our own
// process and compare with os.Args, which the runtime got from the same kernel
// by another road.  Two roads agreeing on every string of a real command line
// cannot happen by accident, and until they do agree the collector publishes
// nobody's command line at all.
func SameArgv(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
