// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "testing"

// record builds one kinfo_proc as the headers lay it out.  The offsets are
// written here as literal numbers on purpose — see the note in
// darwinsys_test.go.
func record(pid, ppid, uid int, stat byte, comm string, pctcpu uint32) []byte {
	b := make([]byte, 648) // sizeof(struct kinfo_proc)
	b[36] = stat           // extern_proc.p_stat
	put32(b, 40, uint32(pid))
	put32(b, 88, pctcpu)      // extern_proc.p_pctcpu
	copy(b[243:243+17], comm) // extern_proc.p_comm
	put32(b, 420, uint32(uid))
	put32(b, 560, uint32(ppid))
	return b
}

func TestParseProcs(t *testing.T) {
	var buf []byte
	buf = append(buf, record(1, 0, 0, 3, "launchd", 0)...)
	buf = append(buf, record(4242, 1, 501, 2, "digitdisk", 2048)...)
	buf = append(buf, record(4243, 4242, 501, 5, "gone", 0)...)

	procs, ok := ParseProcs(buf)
	if !ok {
		t.Fatal("three whole records were refused")
	}
	if len(procs) != 3 {
		t.Fatalf("processes = %d, want 3", len(procs))
	}
	if p := procs[0]; p.PID != 1 || p.PPID != 0 || p.UID != 0 || p.Comm != "launchd" || p.State != StateSleeping {
		t.Errorf("first process = %+v", p)
	}
	if p := procs[1]; p.PID != 4242 || p.PPID != 1 || p.UID != 501 || p.Comm != "digitdisk" || p.State != StateRunning {
		t.Errorf("second process = %+v", p)
	}
	if procs[1].PctCPU != 100 {
		t.Errorf("p_pctcpu = %v, want 100%% of a core for one FSCALE", procs[1].PctCPU)
	}
	if !AnyCPU(procs) {
		t.Error("a list with a busy process must report CPU as published")
	}
	if AnyCPU([]Proc{{PID: 1}, {PID: 2}}) {
		t.Error("a list of zeros is a kernel that does not fill p_pctcpu, not an idle machine")
	}
	if procs[2].State != StateZombie {
		t.Errorf("state letter = %q, want %q", procs[2].State, StateZombie)
	}
}

func TestParseProcsAcceptsTheTrimmedBuffer(t *testing.T) {
	buf := record(7, 1, 0, 2, "kernel_task", 0)
	trimmed := buf[:len(buf)-1] // as syscall.Sysctl hands it over
	procs, ok := ParseProcs(trimmed)
	if !ok || len(procs) != 1 || procs[0].PID != 7 {
		t.Fatalf("trimmed buffer = %+v, ok=%v", procs, ok)
	}
}

func TestParseProcsRefusesAPartialRecord(t *testing.T) {
	buf := append(record(1, 0, 0, 3, "launchd", 0), make([]byte, 100)...)
	if _, ok := ParseProcs(buf); ok {
		t.Error("a buffer that does not divide into records must be refused")
	}
	if _, ok := ParseProcs(nil); ok {
		t.Error("an empty buffer must be refused")
	}
}

func TestVerifyCatchesAShiftedLayout(t *testing.T) {
	procs, _ := ParseProcs(record(4242, 1, 501, 2, "digitdisk", 0))
	if !Verify(procs, 4242, 1, 501) {
		t.Error("the running process with its own parent and user must verify")
	}
	if Verify(procs, 4242, 999, 501) {
		t.Error("the right pid with the wrong parent is a shifted layout, not a match")
	}
	if Verify(procs, 4242, 1, 0) {
		t.Error("the right pid with the wrong user is a shifted layout, not a match")
	}
	if Verify(procs, 999, 1, 501) {
		t.Error("a list without the asking process cannot be this machine's")
	}
	if Verify(nil, 1, 0, 0) {
		t.Error("an empty list verifies nothing")
	}

	// A record whose fields sit one 4-byte slot later than we expect: the
	// pid still lands on a plausible small number, and only the second
	// field gives the shift away.
	shifted := make([]byte, 648)
	put32(shifted, 44, 4242)
	put32(shifted, 424, 501)
	put32(shifted, 564, 1)
	got, ok := ParseProcs(shifted)
	if !ok {
		t.Fatal("the buffer is still record-sized")
	}
	if Verify(got, 4242, 1, 501) {
		t.Error("a shifted layout must not verify")
	}
}
