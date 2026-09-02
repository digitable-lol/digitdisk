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
	if !strings.Contains(out, "занято ЦП     —") {
		t.Errorf("an unsampled CPU must print a dash rather than 0%%:\n%s", out)
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
			sysinfo.FactCPUBusy:      "счётчики не дались",
			sysinfo.FactThreads:      "самопроверка не сошлась",
			sysinfo.FactBlocked:      "система не различает такие процессы",
			sysinfo.FactNetCounters:  "самопроверка не сошлась",
			sysinfo.FactSensors:      "система не публикует показания",
			sysinfo.FactMemoryPages:  "самопроверка не сошлась",
			sysinfo.FactNetTxDrops:   "система не считает такие пакеты",
			sysinfo.FactLoadEntities: "система не публикует очередь",
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
		"занято ЦП     —",
		"ПРОЦЕССЫ  всего 412, выполняется 3",
		"НЕ ИЗМЕРЕНО  ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the report must contain %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "занято        0 Б") || strings.Contains(out, "потоков 0") {
		t.Errorf("an unmeasured field was printed as a measured zero:\n%s", out)
	}
	// The point of the whole rework: a reader looking for numbers is not
	// made to read an explanation of the ones that are not there.
	for _, reason := range []string{"счётчики не дались", "самопроверка не сошлась",
		"система не публикует показания", "система не различает такие процессы"} {
		if strings.Contains(out, reason) {
			t.Errorf("the report printed a reason instead of leaving it to --why: %q\n%s", reason, out)
		}
	}
	// A fact with no column in the report has no line in it either.
	if strings.Contains(out, sysinfo.FactNetTxDrops) || strings.Contains(out, sysinfo.FactLoadEntities) {
		t.Errorf("a fact the report never shows must not be named in it:\n%s", out)
	}
	if !strings.Contains(out, sysinfo.FactSensors) || !strings.Contains(out, sysinfo.FactCPUBusy) {
		t.Errorf("a fact the report would have shown must be named:\n%s", out)
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

func TestStatusOrdersTheUnmeasuredLine(t *testing.T) {
	st := sysinfo.Status{Missing: map[string]string{"я": "1", "а": "2", "м": "3"}}
	var buf bytes.Buffer
	Status(&buf, st)
	if out := buf.String(); !strings.Contains(out, "НЕ ИЗМЕРЕНО  а, м, я") {
		t.Errorf("the line must be ordered, so two runs can be compared:\n%s", out)
	}
}

// The reasons did not vanish; they moved to where somebody who wants them asks
// for them.  That is the whole trade: the report keeps the names, --why keeps
// the sentences, and the JSON keeps both.
func TestWhyCarriesTheReasonsTheReportLeavesOut(t *testing.T) {
	st := sysinfo.Status{Missing: map[string]string{
		sysinfo.FactSensors:    "система не публикует показания датчиков",
		sysinfo.FactNetTxDrops: "система не считает такие пакеты",
	}}
	var buf bytes.Buffer
	Why(&buf, st)
	out := buf.String()
	for _, want := range []string{
		sysinfo.FactSensors, "система не публикует показания датчиков",
		sysinfo.FactNetTxDrops, "система не считает такие пакеты",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("--why must carry %q:\n%s", want, out)
		}
	}
	if i, j := strings.Index(out, sysinfo.FactSensors), strings.Index(out, sysinfo.FactNetTxDrops); i > j {
		t.Errorf("--why must be ordered by name:\n%s", out)
	}
}

func TestWhySaysSoWhenNothingIsMissing(t *testing.T) {
	var buf bytes.Buffer
	Why(&buf, sysinfo.Status{})
	if out := buf.String(); !strings.Contains(out, "снимок полон") {
		t.Errorf("a complete snapshot must say so:\n%s", out)
	}
}

// A partial count must say how partial it is: on macOS a snapshot taken
// without administrator rights covers the caller's own processes only.
func TestStatusSaysHowManyProcessesTheThreadCountCovers(t *testing.T) {
	st := sysinfo.Status{Processes: sysinfo.Processes{
		Total: 906, Running: 4, Threads: 1240, WithDetail: 214,
	}}
	var buf bytes.Buffer
	Status(&buf, st)
	if out := buf.String(); !strings.Contains(out, "замерено по 214 процессам") {
		t.Errorf("a partial count must name its coverage:\n%s", out)
	}

	st.Processes.WithDetail = st.Processes.Total
	buf.Reset()
	Status(&buf, st)
	if out := buf.String(); !strings.Contains(out, "потоков 1240,") || strings.Contains(out, "замерено по") {
		t.Errorf("a complete count needs no qualifier:\n%s", out)
	}
}
