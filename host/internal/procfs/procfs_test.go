// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package procfs

import "testing"

// All samples below are captured text, deliberately not read from the running
// kernel, so the parsers are tested against a fixed input.

const meminfoSample = `MemTotal:       523390352 kB
MemFree:        176926112 kB
MemAvailable:   361842484 kB
Buffers:         7450536 kB
Cached:         193526464 kB
SwapCached:      2699960 kB
Shmem:           2990384 kB
SReclaimable:   11624824 kB
SwapTotal:      662700020 kB
SwapFree:       619129668 kB
HugePages_Total:       0
DirectMap4k:     1234 kB
Nonsense:        12 MB
`

func TestParseMeminfo(t *testing.T) {
	m := ParseMeminfo(meminfoSample)
	if want := uint64(523390352) * 1024; m.Total != want {
		t.Errorf("Total = %d, want %d", m.Total, want)
	}
	if want := uint64(176926112) * 1024; m.Free != want {
		t.Errorf("Free = %d, want %d", m.Free, want)
	}
	// used, the way procps-ng 4 free(1) computes it
	if want := uint64(523390352-361842484) * 1024; m.Used != want {
		t.Errorf("Used = %d, want %d", m.Used, want)
	}
	if want := uint64(7450536+193526464+11624824) * 1024; m.BuffCache != want {
		t.Errorf("BuffCache = %d, want %d", m.BuffCache, want)
	}
	if want := uint64(2990384) * 1024; m.Shared != want {
		t.Errorf("Shared = %d, want %d", m.Shared, want)
	}
	if want := uint64(662700020-619129668) * 1024; m.SwapUsed != want {
		t.Errorf("SwapUsed = %d, want %d", m.SwapUsed, want)
	}
	if m.Raw["HugePages_Total"] != 0 {
		t.Errorf("unit-less key must keep its value unscaled")
	}
	if _, ok := m.Raw["Nonsense"]; ok {
		t.Errorf("unknown unit MB must be skipped, not guessed")
	}
}

func TestParseMeminfoEmpty(t *testing.T) {
	m := ParseMeminfo("")
	if m.Total != 0 || m.Used != 0 || len(m.Raw) != 0 {
		t.Errorf("empty input must yield an empty Memory, got %+v", m)
	}
}

func TestParseLoadAvg(t *testing.T) {
	la, ok := ParseLoadAvg("43.68 51.47 58.79 236/8282 3072717\n")
	if !ok {
		t.Fatal("ParseLoadAvg reported failure on a valid line")
	}
	if la.One != 43.68 || la.Five != 51.47 || la.Fifteen != 58.79 {
		t.Errorf("averages = %v", la)
	}
	if la.Runnable != 236 || la.TotalEntities != 8282 || la.LastPID != 3072717 {
		t.Errorf("entity counts = %+v", la)
	}
	if _, ok := ParseLoadAvg("garbage"); ok {
		t.Errorf("garbage must not parse")
	}
}

func TestParseUptime(t *testing.T) {
	up, idle, ok := ParseUptime("2100448.78 480407629.73\n")
	if !ok || up != 2100448.78 || idle != 480407629.73 {
		t.Errorf("ParseUptime = %v %v %v", up, idle, ok)
	}
	if _, _, ok := ParseUptime(""); ok {
		t.Errorf("empty input must not parse")
	}
}

const statSample = `cpu  2267347153 893978355 2378569559 48040762999 7372698 0 628242 0 0 0
cpu0 10330513 3425525 7606495 162999600 33304 0 192273 0 0 0
cpu1 8890129 2635804 9991724 184081302 43608 0 45993 0 0 0
intr 1 2 3
btime 1786114390
processes 379883810
procs_running 245
procs_blocked 3
`

func TestParseStat(t *testing.T) {
	s := ParseStat(statSample)
	if s.CPU.User != 2267347153 || s.CPU.Idle != 48040762999 || s.CPU.SoftIRQ != 628242 {
		t.Errorf("aggregate cpu line = %+v", s.CPU)
	}
	if len(s.PerCPU) != 2 {
		t.Errorf("per-cpu lines = %d, want 2", len(s.PerCPU))
	}
	if s.BootTime != 1786114390 || s.ProcsRunning != 245 || s.ProcsBlocked != 3 {
		t.Errorf("scalars = %+v", s)
	}
	if s.ForksSinceBoot != 379883810 {
		t.Errorf("forks = %d", s.ForksSinceBoot)
	}
	wantTotal := uint64(2267347153 + 893978355 + 2378569559 + 48040762999 + 7372698 + 0 + 628242 + 0)
	if s.CPU.Total() != wantTotal {
		t.Errorf("Total = %d, want %d", s.CPU.Total(), wantTotal)
	}
	if s.CPU.Busy() != wantTotal-48040762999-7372698 {
		t.Errorf("Busy = %d", s.CPU.Busy())
	}
}

func TestParseProcStat(t *testing.T) {
	// comm deliberately contains a space and both kinds of parenthesis:
	// the split must anchor on the LAST ')'.
	line := "1234 (weird ) name) S 1 1234 1234 0 -1 4194560 100 200 0 0 " +
		"11 22 0 0 20 0 7 0 987654 123456789 4096 18446744073709551615 1 2 3\n"
	p, ok := ParseProcStat(line)
	if !ok {
		t.Fatal("ParseProcStat reported failure")
	}
	if p.PID != 1234 {
		t.Errorf("PID = %d", p.PID)
	}
	if p.Comm != "weird ) name" {
		t.Errorf("Comm = %q, want %q", p.Comm, "weird ) name")
	}
	if p.State != "S" || p.PPID != 1 {
		t.Errorf("State/PPID = %q/%d", p.State, p.PPID)
	}
	if p.UTime != 11 || p.STime != 22 {
		t.Errorf("times = %d/%d, want 11/22", p.UTime, p.STime)
	}
	if p.NumThreads != 7 {
		t.Errorf("threads = %d, want 7", p.NumThreads)
	}
	if p.StartTime != 987654 {
		t.Errorf("start = %d", p.StartTime)
	}
	if p.VSize != 123456789 || p.RSSPages != 4096 {
		t.Errorf("vsize/rss = %d/%d", p.VSize, p.RSSPages)
	}
}

func TestParseProcStatRejectsJunk(t *testing.T) {
	for _, in := range []string{"", "no parens here", "12 (short) S 1 2 3"} {
		if _, ok := ParseProcStat(in); ok {
			t.Errorf("ParseProcStat(%q) accepted junk", in)
		}
	}
}

func TestParseCmdlineAndUID(t *testing.T) {
	if got := ParseCmdline([]byte("ls\x00-la\x00/srv\x00")); got != "ls -la /srv" {
		t.Errorf("ParseCmdline = %q", got)
	}
	if got := ParseCmdline(nil); got != "" {
		t.Errorf("kernel thread must yield empty cmdline, got %q", got)
	}
	uid, ok := ParseStatusUID("Name:\tbash\nUid:\t1001\t1001\t1001\t1001\n")
	if !ok || uid != 1001 {
		t.Errorf("ParseStatusUID = %d %v", uid, ok)
	}
	if _, ok := ParseStatusUID("Name:\tbash\n"); ok {
		t.Errorf("missing Uid line must report failure")
	}
}

func TestParseMounts(t *testing.T) {
	sample := "/dev/mapper/vg99546-root / ext4 rw,relatime 0 0\n" +
		"tmpfs /mnt/with\\040space tmpfs rw,nosuid 0 0\n" +
		"short line\n"
	ms := ParseMounts(sample)
	if len(ms) != 2 {
		t.Fatalf("mounts = %d, want 2 (short line dropped)", len(ms))
	}
	if ms[0].Source != "/dev/mapper/vg99546-root" || ms[0].Point != "/" || ms[0].FSType != "ext4" {
		t.Errorf("first mount = %+v", ms[0])
	}
	if ms[1].Point != "/mnt/with space" {
		t.Errorf("octal escape not decoded: %q", ms[1].Point)
	}
}

func TestParseNetDev(t *testing.T) {
	sample := `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
    lo: 512538206382 47024813    0    0    0     0          0         0 512538206382 47024813    0    0    0     0       0          0
 eno33: 452379043716 711629386    0 184971    0     0          0   1787306 389470585000 360088024    0    5    0     0       0          0
`
	ns := ParseNetDev(sample)
	if len(ns) != 2 {
		t.Fatalf("interfaces = %d, want 2", len(ns))
	}
	if ns[0].Name != "lo" || ns[0].RxBytes != 512538206382 || ns[0].TxPackets != 47024813 {
		t.Errorf("lo = %+v", ns[0])
	}
	if ns[1].Name != "eno33" || ns[1].RxDropped != 184971 || ns[1].TxDropped != 5 {
		t.Errorf("eno33 = %+v", ns[1])
	}
}

func TestParseOSRelease(t *testing.T) {
	kv := ParseOSRelease("# comment\nPRETTY_NAME=\"Ubuntu 26.04 LTS\"\nID=ubuntu\nVERSION_ID=\"26.04\"\nBROKEN\n")
	if kv["PRETTY_NAME"] != "Ubuntu 26.04 LTS" {
		t.Errorf("PRETTY_NAME = %q", kv["PRETTY_NAME"])
	}
	if kv["ID"] != "ubuntu" || kv["VERSION_ID"] != "26.04" {
		t.Errorf("kv = %v", kv)
	}
	if _, ok := kv["BROKEN"]; ok {
		t.Errorf("line without '=' must be ignored")
	}
}

// A kernel too old for MemAvailable, which is what the Present map is for: the
// derived "used" cannot be computed from what is here, and must not come back
// as a zero that reads like a measurement.
const meminfoWithoutAvailable = `MemTotal:        8000 kB
MemFree:         1000 kB
Buffers:          100 kB
Cached:           200 kB
`

func TestParseMeminfoMarksWhatItDidNotSee(t *testing.T) {
	full := ParseMeminfo(meminfoSample)
	for _, field := range []string{FieldTotal, FieldFree, FieldAvailable, FieldUsed, FieldBuffCache, FieldSwapUsed} {
		if !full.Has(field) {
			t.Errorf("%s is in the sample and must be marked present", field)
		}
	}

	old := ParseMeminfo(meminfoWithoutAvailable)
	if !old.Has(FieldTotal) || !old.Has(FieldFree) {
		t.Errorf("what the old kernel does publish must stay present: %+v", old.Present)
	}
	if old.Has(FieldAvailable) || old.Has(FieldUsed) {
		t.Errorf("MemAvailable is absent, so available and used are not measurements: %+v", old.Present)
	}
	if old.Used != 0 {
		t.Errorf("Used = %d, want 0 — and the zero is only safe because Has says it is not a measurement", old.Used)
	}
	if old.Has(FieldSwapTotal) || old.Has(FieldSwapUsed) {
		t.Errorf("there is no swap line in this sample: %+v", old.Present)
	}
}

func TestMemoryHasTakesAHandBuiltStructAtFaceValue(t *testing.T) {
	// A Memory nobody marked up — a test's, or one from a source that
	// fills every field it has — answers true, so the nil map never turns
	// real numbers into dashes.
	m := Memory{Total: 1024}
	if !m.Has(FieldTotal) || !m.Has(FieldUsed) {
		t.Errorf("a nil Present map must be taken at face value")
	}
}
