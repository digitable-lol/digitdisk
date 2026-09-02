// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package sysinfo

import (
	"os"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/darwinsys"
	"digitdisk/internal/libsystem"
	"digitdisk/internal/procfs"
	"syscall"
)

// THE ONLY TESTS IN THIS TREE THAT NEED A MAC.
//
// Everything in internal/darwinsys is a decoder of bytes and is tested
// anywhere, including on the Linux machine this was written on.  What no Linux
// machine can answer is whether the layouts those decoders describe are the
// layouts a Mac really hands back.  These tests ask a Mac.
//
// They are meant to run on a macOS runner in continuous integration and on the
// owner's own machine, and they fail loudly rather than skip: a snapshot that
// silently stopped measuring half of itself is exactly the state this whole
// change exists to end.

func TestMachineStatisticsAreMeasured(t *testing.T) {
	st := New().Collect()

	if st.Load.CPUCount <= 0 {
		t.Errorf("cores = %d", st.Load.CPUCount)
	}
	if st.Load.BusyPercent == nil {
		t.Fatalf("the busy share was not measured; not measured: %v", st.Missing)
	}
	if b := *st.Load.BusyPercent; b < 0 || b > 100 {
		t.Errorf("busy share = %.2f%%, which is not a share", b)
	}

	m := st.Memory
	if m.Total == 0 {
		t.Fatal("the machine reported no memory at all")
	}
	for _, field := range []string{
		procfs.FieldTotal, procfs.FieldFree, procfs.FieldBuffCache,
		procfs.FieldAvailable, procfs.FieldUsed,
	} {
		if !m.Has(field) {
			t.Errorf("memory field %q was not measured; not measured: %v", field, st.Missing)
		}
	}
	if m.Free > m.Total || m.Available > m.Total || m.Used > m.Total {
		t.Errorf("a memory field exceeds the machine's memory: %+v", m)
	}
	if _, ok := m.Raw[procfs.RawWired]; !ok {
		t.Error("wired memory was not measured")
	}
}

func TestProcessesAreMeasured(t *testing.T) {
	st := New().Collect()

	if st.Processes.Total < 2 {
		t.Fatalf("processes = %d on a running system", st.Processes.Total)
	}
	if st.Processes.WithDetail == 0 {
		t.Fatalf("no process was read in detail; not measured: %v", st.Missing)
	}
	// The process that is asking is on a processor, so at least one is.
	if st.Processes.Running < 1 || st.Processes.Running > st.Processes.WithDetail {
		t.Errorf("running = %d of %d processes read in detail; not measured: %v",
			st.Processes.Running, st.Processes.WithDetail, st.Missing)
	}
	if st.Processes.Threads < st.Processes.WithDetail {
		t.Errorf("threads = %d over %d processes: every process has at least one",
			st.Processes.Threads, st.Processes.WithDetail)
	}
	if len(st.Processes.TopByMemory) == 0 {
		t.Fatal("nothing in the memory ranking")
	}
	if st.Processes.TopByMemory[0].RSSBytes <= 0 {
		t.Errorf("the largest process holds %d bytes", st.Processes.TopByMemory[0].RSSBytes)
	}
	if len(st.Processes.TopByCPU) == 0 {
		t.Fatal("nothing in the processor ranking")
	}
	if st.Processes.TopByCPU[0].CPUPercent == nil {
		t.Errorf("the busiest process has no measured share; not measured: %v", st.Missing)
	}
	// The command line of this very process is the one row we can check
	// against something we already know.
	var seen bool
	for _, p := range st.Processes.TopByMemory {
		if p.PID == os.Getpid() && strings.Contains(p.Cmdline, os.Args[0]) {
			seen = true
		}
	}
	if !seen && !st.Missing[FactProcessArgs].Empty() {
		t.Errorf("command lines were not published: %s", st.Missing[FactProcessArgs])
	}
}

// The self-checks, run directly rather than through the collector, so a
// failure names which layout is wrong.

func TestTaskInfoLayoutOnThisMac(t *testing.T) {
	memsize := machineMemory(t)
	b, err := libsystem.ProcTaskInfo(os.Getpid(), darwinsys.TaskInfoFlavor, darwinsys.TaskInfoSize)
	if err != nil {
		t.Fatalf("proc_pidinfo on our own process: %v", err)
	}
	if len(b) != darwinsys.TaskInfoSize {
		t.Fatalf("proc_pidinfo wrote %d bytes, struct proc_taskinfo is %d",
			len(b), darwinsys.TaskInfoSize)
	}
	ti, ok := darwinsys.ParseTaskInfo(b, memsize)
	if !ok {
		t.Fatalf("our own task info did not pass its own checks: % x", b)
	}
	if ti.ResidentBytes == 0 || ti.Threads < 1 {
		t.Errorf("a running test holds %d bytes in %d threads", ti.ResidentBytes, ti.Threads)
	}
}

func TestArgumentBlockLayoutOnThisMac(t *testing.T) {
	argmax, err := syscall.SysctlUint32("kern.argmax")
	if err != nil {
		t.Fatalf("kern.argmax: %v", err)
	}
	b, err := libsystem.SysctlRaw(darwinsys.ArgsMIB(os.Getpid()), int(argmax))
	if err != nil {
		t.Fatalf("KERN_PROCARGS2 on our own process: %v", err)
	}
	_, argv, ok := darwinsys.ParseProcArgs2(b)
	if !ok {
		t.Fatal("our own argument block did not decode")
	}
	if !darwinsys.SameArgv(argv, os.Args) {
		t.Errorf("the kernel and the runtime disagree about our own arguments:\n  kernel  %q\n  runtime %q",
			argv, os.Args)
	}
}

func TestHostStatisticsFlavorsOnThisMac(t *testing.T) {
	b, err := libsystem.HostStatistics(darwinsys.FlavorCPULoad, darwinsys.CPULoadCount)
	if err != nil {
		t.Fatalf("host_statistics(HOST_CPU_LOAD_INFO): %v", err)
	}
	first, ok := darwinsys.ParseCPUTicks(b)
	if !ok {
		t.Fatalf("the kernel answered %d bytes of tick counters", len(b))
	}
	// Busy work rather than sleep, and a generous deadline: on a virtual Mac
	// the counters advance in steps coarser than a hundred milliseconds, so a
	// test that slept once and looked once would be a coin flip.
	var last darwinsys.CPUTicks
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		spin(20 * time.Millisecond)
		if b, err = libsystem.HostStatistics(darwinsys.FlavorCPULoad, darwinsys.CPULoadCount); err != nil {
			t.Fatalf("second reading: %v", err)
		}
		last, _ = darwinsys.ParseCPUTicks(b)
		if share, ok := first.BusyShare(last); ok {
			if share < 0 || share > 100 {
				t.Errorf("busy share = %.2f%%, which is not a share", share)
			}
			break
		}
	}
	if _, ok := first.BusyShare(last); !ok {
		t.Errorf("the counters did not move in five seconds of work: %v then %v", first, last)
	}

	// The memory flavor grew a field in macOS 26.  Asking for the older
	// length must keep working: the kernel answers with the length it filled.
	b, err = libsystem.HostStatistics64(darwinsys.FlavorVMInfo64, darwinsys.VMInfo64Count)
	if err != nil {
		t.Fatalf("host_statistics64(HOST_VM_INFO64): %v", err)
	}
	if len(b) != darwinsys.VMStatSize {
		t.Fatalf("the kernel filled %d bytes, this build reads %d", len(b), darwinsys.VMStatSize)
	}
	memsize := machineMemory(t)
	pageSize, err := syscall.SysctlUint32("hw.pagesize")
	if err != nil {
		t.Skipf("hw.pagesize: %v", err)
	}
	v, ok := darwinsys.ParseVMStat64(b, memsize/uint64(pageSize))
	if !ok {
		t.Fatalf("the memory breakdown did not agree with hw.memsize: % x", b)
	}
	if got, want := v.Accounted(), memsize/uint64(pageSize); got > want || 100*got/want < 90 {
		t.Errorf("the breakdown places %d of the machine's %d pages", got, want)
	}
}

// machineMemory reads hw.memsize the way the collector does, so a test that
// fails names the same source the report would have named.
func machineMemory(t *testing.T) uint64 {
	t.Helper()
	b, err := sysctlRaw("hw.memsize")
	if err != nil {
		t.Fatalf("hw.memsize: %v", err)
	}
	v, ok := darwinsys.ParseUint64(b)
	if !ok || v == 0 {
		t.Fatalf("hw.memsize answered %d bytes that do not read as a size", len(b))
	}
	return v
}

// spin burns processor time, which is the one thing a sleep cannot do and the
// tick counters are there to count.
func spin(d time.Duration) {
	end := time.Now().Add(d)
	for i := 0; time.Now().Before(end); i++ {
		_ = i * i
	}
}
