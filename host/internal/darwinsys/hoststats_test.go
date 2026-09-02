// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "testing"

// The offsets below are written out as literal numbers on purpose — see the
// note at the top of darwinsys_test.go.  A constant that drifts must break a
// test, not agree with itself.

func TestBusyShareIsAShareAndNeedsNoTickRate(t *testing.T) {
	// cpu_ticks[USER, SYSTEM, IDLE, NICE]
	before := CPUTicks{1000, 500, 8000, 0}
	after := CPUTicks{1100, 550, 8850, 0}
	// busy 150 of 1000 elapsed ticks
	got, ok := before.BusyShare(after)
	if !ok {
		t.Fatal("two readings a window apart must produce a share")
	}
	if got != 15 {
		t.Errorf("busy = %.2f%%, want 15%%", got)
	}

	// The same machine sampled ten times as often reads the same share:
	// nothing here depends on how fast the kernel counts.
	tenth, ok := before.BusyShare(CPUTicks{1010, 505, 8085, 0})
	if !ok || tenth != 15 {
		t.Errorf("a tenth of the window = %.2f%% (ok=%v), want the same 15%%", tenth, ok)
	}
}

func TestBusyShareSurvivesTheCounterWrapping(t *testing.T) {
	// Ten cores at a hundred ticks a second pass 2^32 in about fifty days,
	// so a real machine does reach this.
	before := CPUTicks{0xFFFFFFF0, 0, 100, 0}
	after := CPUTicks{9, 0, 125, 0} // user wrapped past zero, idle did not
	got, ok := before.BusyShare(after)
	if !ok {
		t.Fatal("a wrapped counter is still a measurement")
	}
	if got != 50 {
		t.Errorf("busy = %.2f%%, want 50%% (25 busy ticks of 50)", got)
	}
}

func TestBusyShareRefusesAPairWithNoTimeBetween(t *testing.T) {
	same := CPUTicks{1, 2, 3, 4}
	if _, ok := same.BusyShare(same); ok {
		t.Error("two identical readings measure nothing and must be refused")
	}
}

// eventCount stands in for the 64-bit counters between the page counts —
// faults, pageins, lookups.  A Mac that has been awake for a day carries tens
// of billions there, which is the property the shift test leans on: no count
// of physical pages can be that large.
const eventCount = 20_000_000_000

// vmStat builds a struct vm_statistics64 with the page counts at their
// declared offsets and the event counters where a running machine keeps them.
func vmStat(free, active, inactive, wired, purgeable, speculative, compressor, throttled, external, internal uint32) []byte {
	b := make([]byte, 152)
	put32(b, 0, free)
	put32(b, 4, active)
	put32(b, 8, inactive)
	put32(b, 12, wired)
	for _, off := range []int{16, 24, 32, 40, 48, 56, 64, 72, 80} {
		put64(b, off, eventCount) // zero_fill_count … purges
	}
	put32(b, 88, purgeable)
	put32(b, 92, speculative)
	for _, off := range []int{96, 104, 112, 120} {
		put64(b, off, eventCount) // decompressions … swapouts
	}
	put32(b, 128, compressor)
	put32(b, 132, throttled)
	put32(b, 136, external)
	put32(b, 140, internal)
	put64(b, 144, eventCount)
	return b
}

// live is a reading copied field for field off a real Mac — an Intel machine
// with 14 GiB of memory and 4 KiB pages, from the run that first showed the
// buckets overlap.  Building the test out of a measurement rather than out of
// invented numbers is the point: the identity below is the machine's, not ours.
func live() (b []byte, totalPages uint64) {
	return vmStat(
		1_428_074, // free (the speculative pages are already inside it)
		869_212,   // active
		1_023_874, // inactive
		346_775,   // wired
		20_782,    // purgeable
		872_498,   // speculative
		0,         // compressor
		0,         // throttled
		2_221_439, // external, file-backed
		544_145,   // internal, anonymous
	), 14 * 1024 * 1024 * 1024 / 4096
}

func TestParseVMStat64(t *testing.T) {
	b, total := live()
	v, ok := ParseVMStat64(b, total)
	if !ok {
		t.Fatal("a reading taken off a real machine must be accepted")
	}
	if v.Free != 1_428_074 || v.Active != 869_212 || v.Inactive != 1_023_874 || v.Wired != 346_775 {
		t.Errorf("first four counts misread: %+v", v)
	}
	if v.Purgeable != 20_782 || v.Speculative != 872_498 {
		t.Errorf("counts after the event block misread: %+v", v)
	}
	if v.Compressor != 0 || v.Throttled != 0 || v.External != 2_221_439 || v.Internal != 544_145 {
		t.Errorf("counts at the tail misread: %+v", v)
	}
	if v.TrulyFree() != 555_576 {
		t.Errorf("TrulyFree() = %d, want 555576 — the speculative pages come out", v.TrulyFree())
	}
	// The machine has 3 670 016 pages and this places 3 667 935 of them.
	if got := v.Accounted(); got != 3_667_935 || got > total || 100*got/total < 99 {
		t.Errorf("Accounted() = %d of %d pages", got, total)
	}
	// The second cut through the same pages: anonymous plus file-backed is
	// active plus inactive plus speculative.
	if v.Internal+v.External != v.Active+v.Inactive+v.Speculative {
		t.Errorf("internal+external = %d, active+inactive+speculative = %d",
			v.Internal+v.External, v.Active+v.Inactive+v.Speculative)
	}
}

func TestParseVMStat64RefusesMoreReadAheadThanFreeMemory(t *testing.T) {
	// The speculative pages are a part of the free ones and cannot outnumber
	// them.  Two fields eighty-eight bytes apart agreeing on that is a layout
	// check as much as a sanity one.
	b, total := live()
	put32(b, 92, 1_428_075) // speculative, one page past free
	if _, ok := ParseVMStat64(b, total); ok {
		t.Error("more read-ahead pages than free pages must be refused")
	}
}

func TestParseVMStat64RefusesAShiftedLayout(t *testing.T) {
	good, total := live()

	// Reading one slot early puts the low half of an event counter where a
	// page count belongs — a number in the billions on a machine with four
	// million pages.  That is the mistake the ceilings exist to catch.
	early := append(append([]byte{}, good[4:]...), 0, 0, 0, 0)
	if _, ok := ParseVMStat64(early, total); ok {
		t.Error("a layout read one field early must be refused")
	}

	// A structure of a different size never gets this far in the collector,
	// because the kernel reports how many elements it filled; the decoder
	// refuses it anyway rather than reading a prefix of something else.
	if _, ok := ParseVMStat64(append(append([]byte{}, good...), 0, 0, 0, 0), total); ok {
		t.Error("a longer structure must be refused, not read by its prefix")
	}
	if _, ok := ParseVMStat64(good[:len(good)-4], total); ok {
		t.Error("a shorter structure must be refused")
	}
}

func TestParseVMStat64RefusesCountsTheMachineCannotHold(t *testing.T) {
	_, total := live()
	tooMuch := vmStat(uint32(total)+1, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	if _, ok := ParseVMStat64(tooMuch, total); ok {
		t.Error("more free pages than the machine has must be refused")
	}
	// A sum far under the machine's size means the buckets are not the
	// buckets we think they are.
	empty := vmStat(1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	if _, ok := ParseVMStat64(empty, total); ok {
		t.Error("a reading accounting for almost no memory must be refused")
	}
	// Nine tenths is under the slack the kernel's own pages need, and must
	// still be refused: the band is there to catch a layout, not to be wide.
	short, _ := live()
	put32(short, 136, 2_221_439-500_000) // external
	if _, ok := ParseVMStat64(short, total); ok {
		t.Error("a reading half a gigabyte short of the machine must be refused")
	}
	if _, ok := ParseVMStat64(vmStat(1, 1, 1, 1, 1, 1, 1, 1, 1, 1), 0); ok {
		t.Error("without the machine's page count there is nothing to check against")
	}
	if _, ok := ParseVMStat64(make([]byte, 100), total); ok {
		t.Error("a short buffer must be refused")
	}
}

// The same identity on the other kind of Mac: Apple Silicon, 16 KiB pages,
// macOS 26 — where the structure itself grew a field at the end and the kernel
// answers a request for the older length with the older length.  These numbers
// were read off that machine.
func TestParseVMStat64OnAppleSilicon(t *testing.T) {
	const total = 7 * 1024 * 1024 * 1024 / 16384
	b := vmStat(20_886, 184_838, 173_650, 57_829, 247, 10_218, 20_017, 0, 235_672, 133_034)

	v, ok := ParseVMStat64(b, total)
	if !ok {
		t.Fatal("a reading taken off an Apple Silicon Mac must be accepted")
	}
	if v.TrulyFree() != 10_668 {
		t.Errorf("TrulyFree() = %d, want 10668", v.TrulyFree())
	}
	if got := v.Accounted(); got != 457_220 || got > total || 100*got/total < 99 {
		t.Errorf("Accounted() = %d of %d pages", got, total)
	}
	if v.Internal+v.External != v.Active+v.Inactive+v.Speculative {
		t.Errorf("internal+external = %d, active+inactive+speculative = %d",
			v.Internal+v.External, v.Active+v.Inactive+v.Speculative)
	}
}
