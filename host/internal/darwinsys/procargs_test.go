// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "testing"

// procArgs builds a KERN_PROCARGS2 block: the argument count, the executed
// path, alignment padding, the arguments, then the environment the decoder
// must stop before.
func procArgs(argc int32, exe string, argv, env []string) []byte {
	b := make([]byte, 4)
	put32(b, 0, uint32(argc))
	b = append(b, exe...)
	b = append(b, 0, 0, 0) // the kernel pads the path out
	for _, a := range argv {
		b = append(b, a...)
		b = append(b, 0)
	}
	for _, e := range env {
		b = append(b, e...)
		b = append(b, 0)
	}
	return b
}

func TestParseProcArgs2(t *testing.T) {
	b := procArgs(3, "/usr/local/bin/digitdisk",
		[]string{"digitdisk", "status", "--top=5"},
		[]string{"HOME=/Users/kто-то", "SECRET=не для отчёта"})

	exe, argv, ok := ParseProcArgs2(b)
	if !ok {
		t.Fatal("a well-formed block must decode")
	}
	if exe != "/usr/local/bin/digitdisk" {
		t.Errorf("exe = %q", exe)
	}
	if len(argv) != 3 || argv[0] != "digitdisk" || argv[1] != "status" || argv[2] != "--top=5" {
		t.Errorf("argv = %q", argv)
	}
	for _, a := range argv {
		if a == "SECRET=не для отчёта" {
			t.Error("the environment must not be read: a report is not a place for it")
		}
	}
}

func TestParseProcArgs2RefusesABlockItCannotFinish(t *testing.T) {
	// The count says four arguments and three are there.  Publishing the
	// three would print a command line nobody ran.
	short := procArgs(4, "/bin/sh", []string{"sh", "-c", "true"}, nil)
	if _, _, ok := ParseProcArgs2(short); ok {
		t.Error("a block with fewer arguments than it promises must be refused")
	}
	if _, _, ok := ParseProcArgs2([]byte{1, 0}); ok {
		t.Error("a block too short to hold a count must be refused")
	}
	if _, _, ok := ParseProcArgs2(procArgs(-1, "/bin/sh", nil, nil)); ok {
		t.Error("a negative argument count must be refused")
	}
	// A path with no terminator is not a path.
	if _, _, ok := ParseProcArgs2([]byte{1, 0, 0, 0, 'x', 'y'}); ok {
		t.Error("a block with no terminator must be refused")
	}
}

func TestParseProcArgs2ReadsAProcessWithNoArguments(t *testing.T) {
	exe, argv, ok := ParseProcArgs2(procArgs(0, "/sbin/launchd", nil, []string{"PATH=/usr/bin"}))
	if !ok || exe != "/sbin/launchd" || len(argv) != 0 {
		t.Errorf("exe=%q argv=%q ok=%v — a process may legitimately carry no arguments", exe, argv, ok)
	}
}

func TestSameArgvIsTheSelfCheck(t *testing.T) {
	if !SameArgv([]string{"digitdisk", "status"}, []string{"digitdisk", "status"}) {
		t.Error("the same list must compare equal")
	}
	if SameArgv([]string{"digitdisk"}, []string{"digitdisk", "status"}) {
		t.Error("a shorter list is not the same list")
	}
	if SameArgv([]string{"digitdisk", "analyze"}, []string{"digitdisk", "status"}) {
		t.Error("a different argument is not the same list")
	}
}
