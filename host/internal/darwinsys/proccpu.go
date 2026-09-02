// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

// The per-processor tick counters, which is the same reading as
// HOST_CPU_LOAD_INFO taken one processor at a time.
//
// host_processor_info answers with an array of struct processor_cpu_load_info
// — the very structure host_statistics answers with for the whole machine —
// laid end to end, one for each processor, plus the number of processors it
// wrote about.  So the decoding is the decoding already written for the
// machine-wide flavor, applied in a loop; what is new is only the length
// check, and that is where a wrong idea about the layout is caught.

// Flavor and element count of the per-processor statistics
// (<mach/processor_info.h>).  The count is per processor.
const (
	FlavorProcessorCPULoad = 2 // PROCESSOR_CPU_LOAD_INFO
	ProcessorCPULoadCount  = 4 // PROCESSOR_CPU_LOAD_INFO_COUNT
)

// ParseProcessorTicks decodes the array against the processor count the
// kernel returned beside it.
//
// The length has to be exactly the count times the size of one structure.
// That is the whole check available here and it is a real one: the caller
// passes the byte length the kernel itself reported, so an array of some
// other shape — a kernel that grew the structure, a call that answered about
// a different number of processors than it filled in — is refused instead of
// being read as if every processor after the first were shifted.
func ParseProcessorTicks(b []byte, cpus int) ([]CPUTicks, bool) {
	if cpus <= 0 || len(b) != cpus*CPUTicksSize {
		return nil, false
	}
	out := make([]CPUTicks, 0, cpus)
	for i := 0; i < cpus; i++ {
		t, ok := ParseCPUTicks(b[i*CPUTicksSize : (i+1)*CPUTicksSize])
		if !ok {
			return nil, false
		}
		out = append(out, t)
	}
	return out, true
}
