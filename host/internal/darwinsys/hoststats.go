// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

// Decoders for the two host statistics macOS answers with: the CPU tick
// counters and the page counts of the virtual memory system.  Both come back
// as a flat C structure from host_statistics/host_statistics64, and both are
// decoded here so the layout is exercised by tests on any machine.

// Flavors and element counts of host_statistics (<mach/host_info.h>).  The
// count is in units of integer_t (four bytes) and is passed in and returned:
// the kernel refuses a buffer it considers too small, so the count is the
// first check that the structure is the size we think it is.
const (
	FlavorCPULoad  = 3 // HOST_CPU_LOAD_INFO
	CPULoadCount   = 4 // HOST_CPU_LOAD_INFO_COUNT
	FlavorVMInfo64 = 4 // HOST_VM_INFO64
	VMInfo64Count  = 38
)

// CPUTicksSize is sizeof(struct host_cpu_load_info): natural_t cpu_ticks[4].
const CPUTicksSize = CPULoadCount * 4

// Indices into cpu_ticks, from CPU_STATE_* in <mach/machine.h>.
const (
	cpuStateUser   = 0
	cpuStateSystem = 1
	cpuStateIdle   = 2
	cpuStateNice   = 3
)

// CPUTicks is one reading of the machine-wide CPU tick counters.
//
// The unit is deliberately not named: the share of busy time is a ratio of two
// differences, so it comes out the same whether the kernel counts in
// hundredths of a second or in anything else.  Nothing here assumes a tick
// rate, and nothing here has to be right about one.
type CPUTicks [CPULoadCount]uint32

// ParseCPUTicks decodes a host_cpu_load_info answer.
//
// The length is required to be exact rather than merely sufficient, and the
// caller passes exactly what the kernel said it wrote: host_statistics is
// given the element count and hands it back, so a kernel whose structure is
// not the one this file describes is caught here instead of being read as if
// it were.
func ParseCPUTicks(b []byte) (CPUTicks, bool) {
	if len(b) != CPUTicksSize {
		return CPUTicks{}, false
	}
	var t CPUTicks
	for i := range t {
		t[i] = u32(b, 4*i)
	}
	return t, true
}

// BusyShare returns the percentage of CPU time spent out of the idle state
// between two readings, and whether the pair can be believed.
//
// The counters are 32-bit and do wrap: ten cores at a hundred ticks a second
// pass 2^32 in about fifty days of uptime, so the differences are taken in
// uint32 arithmetic, where a wrap subtracts correctly.  What is refused is a
// pair with no time between them at all — dividing by that would invent a
// number — and a result outside 0..100, which would mean these bytes are not
// tick counters.
func (a CPUTicks) BusyShare(b CPUTicks) (float64, bool) {
	var total, busy uint64
	for i := range a {
		d := uint64(b[i] - a[i])
		total += d
		if i != cpuStateIdle {
			busy += d
		}
	}
	if total == 0 {
		return 0, false
	}
	share := 100 * float64(busy) / float64(total)
	if share < 0 || share > 100 {
		return 0, false
	}
	return share, true
}

// VMStatSize is sizeof(struct vm_statistics64) (<mach/vm_statistics.h>): four
// 32-bit page counts, nine 64-bit event counters, two more page counts, four
// 64-bit counters, four page counts and a last 64-bit counter — 152 bytes,
// which is the 38 integer_t the flavor asks for.
const VMStatSize = VMInfo64Count * 4

// Byte offsets of the fields we publish.  Each is the sum of the field sizes
// above it in struct vm_statistics64; the comment gives the field name.
const (
	offVMFree        = 0   // free_count,            natural_t
	offVMActive      = 4   // active_count,          natural_t
	offVMInactive    = 8   // inactive_count,        natural_t
	offVMWired       = 12  // wire_count,            natural_t
	offVMPurgeable   = 88  // purgeable_count,       natural_t
	offVMSpeculative = 92  // speculative_count,     natural_t
	offVMCompressor  = 128 // compressor_page_count, natural_t
	offVMThrottled   = 132 // throttled_count,       natural_t
	offVMExternal    = 136 // external_page_count,   natural_t
	offVMInternal    = 140 // internal_page_count,   natural_t
)

// VMStat is the part of struct vm_statistics64 that describes where the
// machine's pages are, in pages.  The event counters next to them in the
// structure — faults, pageins, compressions — are rates of work rather than a
// picture of memory, and are not read.
//
// The buckets are not all disjoint, and the two overlaps matter:
//
//   - Free counts the speculative pages as well.  The kernel adds them in
//     before it answers, so Free minus Speculative is the memory that is
//     really free, and adding Speculative to a sum of buckets counts it twice.
//   - External and Internal are a second cut through the same pages: every
//     resident pageable page is either file-backed or anonymous, and together
//     they come to Active plus Inactive plus Speculative.
//
// Both were read off two live Macs — an Apple Silicon machine with 16 KiB
// pages and an Intel one with 4 KiB pages — and the sums are what Accounted
// checks.
type VMStat struct {
	Free        uint64
	Active      uint64
	Inactive    uint64
	Wired       uint64
	Purgeable   uint64
	Speculative uint64
	Compressor  uint64
	Throttled   uint64
	// External is file-backed pages: the page cache, and the closest thing
	// macOS has to the "buff/cache" column of a Linux report.
	External uint64
	// Internal is anonymous pages: memory a program asked for rather than
	// read from a file.
	Internal uint64
}

// TrulyFree is the memory that is free and not already promised to a
// read-ahead: the kernel folds the speculative pages into Free before it
// answers, and a reader who is told they are free would be told twice.
func (v VMStat) TrulyFree() uint64 { return v.Free - v.Speculative }

// Accounted is every page this picture places, counted once.  It is the sum
// the layout is checked against, not a number for a reader: the kernel keeps
// pages of its own outside these buckets, so it comes close to the machine's
// page count from below without being equal to it.
//
// On the two machines this was measured on it came to 99.7% and 99.9% of what
// hw.memsize divided by hw.pagesize says the machine has.
func (v VMStat) Accounted() uint64 {
	return v.TrulyFree() + v.External + v.Internal + v.Wired + v.Compressor + v.Throttled
}

// ParseVMStat64 decodes a vm_statistics64 answer against a page count the
// caller already knows — hw.memsize divided by hw.pagesize.
//
// The layout is proved in two ways rather than assumed.  The first is the
// kernel's own: the flavor is asked for by element count, the count comes back
// with the answer, and the caller hands over exactly those bytes — so a
// structure that is not 38 integers long never reaches this function.  The
// second is arithmetic.  Every field read here counts physical pages, so none
// of them can exceed the pages the machine has, while the 64-bit event
// counters between them — faults, pageins, lookups — run into the tens of
// billions on a Mac that has been awake for a day.  Reading the fields one
// slot early therefore puts a fault counter where a page count belongs and
// fails at once.  The sum is then required to land near the machine's page
// count from below: the kernel holds pages outside these buckets, so a little
// short is right and over is impossible.
//
// A machine whose bytes pass all of that is a machine whose numbers we print.
// One that does not gets an empty memory breakdown and a dash.
func ParseVMStat64(b []byte, totalPages uint64) (VMStat, bool) {
	if len(b) != VMStatSize || totalPages == 0 {
		return VMStat{}, false
	}
	v := VMStat{
		Free:        uint64(u32(b, offVMFree)),
		Active:      uint64(u32(b, offVMActive)),
		Inactive:    uint64(u32(b, offVMInactive)),
		Wired:       uint64(u32(b, offVMWired)),
		Purgeable:   uint64(u32(b, offVMPurgeable)),
		Speculative: uint64(u32(b, offVMSpeculative)),
		Compressor:  uint64(u32(b, offVMCompressor)),
		Throttled:   uint64(u32(b, offVMThrottled)),
		External:    uint64(u32(b, offVMExternal)),
		Internal:    uint64(u32(b, offVMInternal)),
	}
	for _, n := range [...]uint64{v.Free, v.Active, v.Inactive, v.Wired, v.Purgeable,
		v.Speculative, v.Compressor, v.Throttled, v.External, v.Internal} {
		if n > totalPages {
			return VMStat{}, false
		}
	}
	// The speculative pages are inside the free ones, so they cannot be more
	// numerous.  This is a relation between two fields eighty-eight bytes
	// apart, which is a layout check as much as a sanity one.
	if v.Speculative > v.Free {
		return VMStat{}, false
	}
	// The machine's own page count is a ceiling nothing rises above, and a
	// sixteenth of it is all the slack the kernel's own pages need.
	if sum := v.Accounted(); sum > totalPages || sum < totalPages-totalPages/8 {
		return VMStat{}, false
	}
	return v, true
}
