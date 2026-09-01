// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"bytes"
	"strings"
	"testing"

	"digitdisk/internal/core"
	"digitdisk/internal/procfs"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
)

func TestBytes(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 Б"},
		{512, "512 Б"},
		{1024, "1.0 КиБ"},
		{1536, "1.5 КиБ"},
		{1 << 20, "1.0 МиБ"},
		{1 << 30, "1.0 ГиБ"},
		{1 << 40, "1.0 ТиБ"},
		{-2048, "-2.0 КиБ"},
	}
	for _, c := range cases {
		if got := Bytes(c.in); got != c.want {
			t.Errorf("Bytes(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestStatusPrintsDashesForMissingData(t *testing.T) {
	var buf bytes.Buffer
	Status(&buf, sysinfo.Status{})
	out := buf.String()
	if !strings.Contains(out, "узел          —") {
		t.Errorf("an absent hostname must print as —, got:\n%s", out)
	}
	if !strings.Contains(out, "занято ЦП     — (замер не делался)") {
		t.Errorf("an unsampled CPU must say so rather than print 0%%:\n%s", out)
	}
	if strings.Contains(out, "0.0%") && !strings.Contains(out, "средняя       0.00") {
		t.Errorf("empty status invented a percentage:\n%s", out)
	}
}

func TestAnalyzeWarnsWhenDeciderIsAStub(t *testing.T) {
	var buf bytes.Buffer
	Analyze(&buf, scan.Result{
		Root:      "/x",
		Decider:   core.Default().Name(),
		ByClass:   map[core.Class]scan.Bucket{},
		ByVerdict: map[core.Verdict]scan.Bucket{},
	})
	out := buf.String()
	if !strings.Contains(out, "ВНИМАНИЕ") {
		t.Errorf("a report backed by the stub must say so:\n%s", out)
	}
	for _, c := range core.Classes {
		if !strings.Contains(out, string(c)) {
			t.Errorf("разряд %q missing from the report", c)
		}
	}
}

// A snapshot from a system that publishes less than Linux does: the fields
// exist because the shape is shared, and each one says out loud that nobody
// measured it.  This is the rule the whole port stands on — an absent
// measurement must never print as a zero.
func TestStatusPrintsDashesForFactsTheSystemDoesNotPublish(t *testing.T) {
	st := sysinfo.Status{
		Memory: procfs.Memory{
			Total:     32 << 30,
			SwapTotal: 2 << 30,
			SwapUsed:  1 << 30,
			Present: map[string]bool{
				procfs.FieldTotal:     true,
				procfs.FieldSwapTotal: true,
				procfs.FieldSwapUsed:  true,
			},
		},
		Processes: sysinfo.Processes{Total: 412, Running: 3},
		Network:   []sysinfo.Iface{{NetCounters: procfs.NetCounters{Name: "en0"}, OperState: "up"}},
		Missing: map[string]string{
			sysinfo.FactCPUBusy:     "нужен вызов Mach",
			sysinfo.FactThreads:     "в kinfo_proc потоков нет",
			sysinfo.FactBlocked:     "состояния не различают непрерываемый сон",
			sysinfo.FactNetCounters: "разбор не сошёлся с MTU",
			sysinfo.FactSensors:     "нужен IOKit",
			sysinfo.FactMemoryPages: "нужен host_statistics64",
		},
	}
	var buf bytes.Buffer
	Status(&buf, st)
	out := buf.String()

	for _, want := range []string{
		"всего         32.0 ГиБ", // measured, and printed
		"занято        —",        // not measured, and said so
		"свободно      —",
		"доступно      —",
		"кэш/буферы    —",
		"своп          1.0 ГиБ из 2.0 ГиБ занято",
		"занято ЦП     — (нужен вызов Mach)",
		"потоков —",
		"заблокировано —",
		"— (нужен IOKit)",
		"НЕ ИЗМЕРЕНО",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the report must contain %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "занято        0 Б") || strings.Contains(out, "потоков 0") {
		t.Errorf("an unmeasured field was printed as a measured zero:\n%s", out)
	}
	// The interface is real and is listed; its counters are not measured.
	if !strings.Contains(out, "en0") {
		t.Errorf("a known interface must still be listed:\n%s", out)
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "en0") && strings.Contains(line, "0 Б") {
			t.Errorf("unmeasured counters printed as zero bytes: %q", line)
		}
	}
}

func TestStatusOrdersTheUnmeasuredSection(t *testing.T) {
	st := sysinfo.Status{Missing: map[string]string{"я": "1", "а": "2", "м": "3"}}
	var buf bytes.Buffer
	Status(&buf, st)
	out := buf.String()
	if a, m, ya := strings.Index(out, "\n  а: "), strings.Index(out, "\n  м: "), strings.Index(out, "\n  я: "); !(a < m && m < ya) {
		t.Errorf("the section must be ordered, so two runs can be compared:\n%s", out)
	}
}
