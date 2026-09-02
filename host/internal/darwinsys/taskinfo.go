// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

// TaskInfoFlavor is PROC_PIDTASKINFO, the flavor of proc_pidinfo(3) that
// answers with struct proc_taskinfo (<sys/proc_info.h>).
const TaskInfoFlavor = 4

// TaskInfoSize is sizeof(struct proc_taskinfo): six 64-bit counters followed
// by twelve 32-bit ones.
//
// The kernel is told this number and refuses a smaller buffer, then returns
// how many bytes it wrote.  That makes the size the first proof of the layout:
// a wrong constant here does not produce wrong numbers, it produces no numbers.
const TaskInfoSize = 6*8 + 12*4

// Byte offsets inside struct proc_taskinfo, in declaration order.
const (
	offTIVirtual  = 0  // pti_virtual_size,  uint64_t
	offTIResident = 8  // pti_resident_size, uint64_t
	offTIUser     = 16 // pti_total_user,    uint64_t, nanoseconds
	offTISystem   = 24 // pti_total_system,  uint64_t, nanoseconds
	offTIThreads  = 84 // pti_threadnum,     int32_t  (48 + 9*4)
	offTIRunning  = 88 // pti_numrunning,    int32_t  (48 + 10*4)
)

// TaskInfo is what one process costs: memory now, processor time so far, and
// how many threads it is spread over.
type TaskInfo struct {
	VirtualBytes  uint64
	ResidentBytes uint64
	// CPUNanos is user plus system time consumed since the process started.
	CPUNanos uint64
	Threads  int
	// Running is how many of those threads are on a processor right now.
	//
	// It is the only honest answer macOS has to "is this process running":
	// the scheduler state in a process record says SRUN for every process
	// that is merely alive, because sleeping is a property of a thread and
	// not of a process.  Counting processes by that field would report the
	// whole list as running, which is what it used to do.
	Running int
}

// ParseTaskInfo decodes a proc_pidinfo answer against the machine's memory
// size, which the caller already knows from hw.memsize.
//
// The checks are the arithmetic the kernel cannot break: a process cannot hold
// more memory than the machine has, its address space is never smaller than
// the part of it that is resident, and it has at least one thread and not an
// absurd number of them.  A layout read one field out of place puts a
// nanosecond count — a number in the billions — where a thread count belongs,
// and fails here rather than in the report.
//
// A zero resident size is allowed: a process on its way out really does hold
// no pages, and refusing that would drop a row that is telling the truth.
func ParseTaskInfo(b []byte, memTotal uint64) (TaskInfo, bool) {
	if len(b) < TaskInfoSize || memTotal == 0 {
		return TaskInfo{}, false
	}
	t := TaskInfo{
		VirtualBytes:  u64(b, offTIVirtual),
		ResidentBytes: u64(b, offTIResident),
		CPUNanos:      u64(b, offTIUser) + u64(b, offTISystem),
		Threads:       int(i32(b, offTIThreads)),
		Running:       int(i32(b, offTIRunning)),
	}
	if t.ResidentBytes > memTotal || t.VirtualBytes < t.ResidentBytes {
		return TaskInfo{}, false
	}
	if t.Threads < 1 || t.Threads > maxThreads {
		return TaskInfo{}, false
	}
	// A task cannot have more threads on a processor than it has threads.
	if t.Running < 0 || t.Running > t.Threads {
		return TaskInfo{}, false
	}
	return t, true
}

// maxThreads is a ceiling no process on a workstation reaches, and one no
// misread field stays under: a thread count taken from the wrong offset is a
// pointer or a nanosecond counter, not a number in the thousands.
const maxThreads = 1 << 16
