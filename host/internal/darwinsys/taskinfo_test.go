// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "testing"

// taskInfo builds a struct proc_taskinfo: six 64-bit counters, then twelve
// 32-bit ones.  Offsets are written out as literal numbers on purpose.
func taskInfo(virtual, resident, user, system uint64, threads uint32) []byte {
	return taskInfoRunning(virtual, resident, user, system, threads, 1)
}

func taskInfoRunning(virtual, resident, user, system uint64, threads, running uint32) []byte {
	b := make([]byte, 96)
	put64(b, 0, virtual)   // pti_virtual_size
	put64(b, 8, resident)  // pti_resident_size
	put64(b, 16, user)     // pti_total_user
	put64(b, 24, system)   // pti_total_system
	put64(b, 32, user)     // pti_threads_user
	put64(b, 40, system)   // pti_threads_system
	put32(b, 48, 1)        // pti_policy
	put32(b, 52, 900_000)  // pti_faults
	put32(b, 56, 400)      // pti_pageins
	put32(b, 60, 70_000)   // pti_cow_faults
	put32(b, 64, 120_000)  // pti_messages_sent
	put32(b, 68, 120_000)  // pti_messages_received
	put32(b, 72, 800_000)  // pti_syscalls_mach
	put32(b, 76, 300_000)  // pti_syscalls_unix
	put32(b, 80, 2_000_00) // pti_csw
	put32(b, 84, threads)  // pti_threadnum
	put32(b, 88, running)  // pti_numrunning
	put32(b, 92, 31)       // pti_priority
	return b
}

const machineMemory = 32 * 1024 * 1024 * 1024

func TestParseTaskInfo(t *testing.T) {
	b := taskInfo(420*1024*1024*1024, 1500*1024*1024, 12_000_000_000, 3_000_000_000, 27)
	ti, ok := ParseTaskInfo(b, machineMemory)
	if !ok {
		t.Fatal("a plausible process must be accepted")
	}
	if ti.ResidentBytes != 1500*1024*1024 {
		t.Errorf("resident = %d, want %d", ti.ResidentBytes, uint64(1500*1024*1024))
	}
	if ti.CPUNanos != 15_000_000_000 {
		t.Errorf("cpu = %d ns, want user+system = 15000000000", ti.CPUNanos)
	}
	if ti.Threads != 27 {
		t.Errorf("threads = %d, want 27", ti.Threads)
	}
	if ti.Running != 1 {
		t.Errorf("running threads = %d, want 1", ti.Running)
	}
}

func TestParseTaskInfoRefusesMoreRunningThreadsThanThreads(t *testing.T) {
	b := taskInfoRunning(1<<40, 4096, 0, 0, 4, 5)
	if _, ok := ParseTaskInfo(b, machineMemory); ok {
		t.Error("a task cannot have more threads on a processor than it has")
	}
}

func TestParseTaskInfoAcceptsAProcessHoldingNothing(t *testing.T) {
	// A process on its way out really does hold no pages, and dropping its
	// row would hide a process that exists.
	if _, ok := ParseTaskInfo(taskInfo(1024, 0, 0, 0, 1), machineMemory); !ok {
		t.Error("zero resident memory is a measurement, not a failure")
	}
}

func TestParseTaskInfoRefusesWhatTheKernelCannotMean(t *testing.T) {
	cases := []struct {
		why string
		b   []byte
	}{
		{"more resident memory than the machine has",
			taskInfo(1<<50, machineMemory+1, 0, 0, 4)},
		{"an address space smaller than the resident part of it",
			taskInfo(1024, 4096, 0, 0, 4)},
		{"no threads at all",
			taskInfo(1<<40, 4096, 0, 0, 0)},
		{"a thread count read out of a nanosecond counter",
			taskInfo(1<<40, 4096, 0, 0, 3_000_000_000)},
		{"a structure shorter than the kernel writes",
			taskInfo(1<<40, 4096, 0, 0, 4)[:95]},
	}
	for _, c := range cases {
		if _, ok := ParseTaskInfo(c.b, machineMemory); ok {
			t.Errorf("must be refused: %s", c.why)
		}
	}
	if _, ok := ParseTaskInfo(taskInfo(1<<40, 4096, 0, 0, 4), 0); ok {
		t.Error("without the machine's memory size there is nothing to check against")
	}
}
