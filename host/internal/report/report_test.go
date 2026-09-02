// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"bytes"
	"strings"
	"testing"

	"digitdisk/internal/core"
	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/lang"
	"digitdisk/internal/procfs"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
)

// Everything a person reads is read in one of two languages, so everything
// worth checking is checked in both.  A report that is right in Russian and
// broken in English is broken, and a test that only ever ran in Russian is
// exactly how it would stay broken: half the wordings would print their key.
var both = []lang.Lang{lang.RU, lang.EN}

func status(l lang.Lang, st sysinfo.Status) string {
	var b bytes.Buffer
	Status(&b, l, st)
	return b.String()
}

// TestBytes is about a comma and a point, and they are not decoration:
// «12,3 ГиБ» read by an English reader is twelve and three, and "12.3 GiB"
// read by a Russian one is the same mistake mirrored.  The unit is spelled by
// the language too — КиБ is the Russian standard's spelling, not a
// transliteration of KiB.
func TestBytes(t *testing.T) {
	cases := []struct {
		in     int64
		ru, en string
	}{
		{0, "0 Б", "0 B"},
		{512, "512 Б", "512 B"},
		{1024, "1,0 КиБ", "1.0 KiB"},
		{1536, "1,5 КиБ", "1.5 KiB"},
		{1 << 20, "1,0 МиБ", "1.0 MiB"},
		{1 << 30, "1,0 ГиБ", "1.0 GiB"},
		{1 << 40, "1,0 ТиБ", "1.0 TiB"},
		{-2048, "-2,0 КиБ", "-2.0 KiB"},
	}
	for _, c := range cases {
		if got := lang.RU.Bytes(c.in); got != c.ru {
			t.Errorf("RU.Bytes(%d) = %q, ждали %q", c.in, got, c.ru)
		}
		if got := lang.EN.Bytes(c.in); got != c.en {
			t.Errorf("EN.Bytes(%d) = %q, want %q", c.in, got, c.en)
		}
		if c.in >= 0 {
			if got := lang.RU.UBytes(uint64(c.in)); got != c.ru {
				t.Errorf("RU.UBytes(%d) = %q, ждали %q", c.in, got, c.ru)
			}
		}
	}
}

func TestStatusPrintsDashesForMissingData(t *testing.T) {
	for _, l := range both {
		out := status(l, sysinfo.Status{})
		// The English half is written out here rather than looked up:
		// a check that asked the dictionary what the dictionary says
		// would pass with an empty dictionary.
		host, cpu := "  узел          —", "  занято ЦП     —"
		if l == lang.EN {
			host, cpu = "  host          —", "  CPU busy      —"
		}
		if !strings.Contains(out, host) {
			t.Errorf("[%s] an absent hostname must print as —, got:\n%s", l, out)
		}
		if !strings.Contains(out, cpu) {
			t.Errorf("[%s] an unsampled CPU must print a dash rather than 0%%:\n%s", l, out)
		}
		if strings.Contains(out, l.Pct(0, 1)) && !strings.Contains(out, l.F("  средняя       %s / %s / %s  (1/5/15 мин)", l.Dec(0, 2), l.Dec(0, 2), l.Dec(0, 2))) {
			t.Errorf("[%s] empty status invented a percentage:\n%s", l, out)
		}
	}
}

func TestAnalyzeWarnsWhenDeciderIsAStub(t *testing.T) {
	for _, l := range both {
		var buf bytes.Buffer
		Analyze(&buf, l, scan.Result{
			Root:      "/x",
			Decider:   core.Default().Name(),
			ByClass:   map[core.Class]scan.Bucket{},
			ByVerdict: map[core.Verdict]scan.Bucket{},
		})
		out := buf.String()
		warn := "ВНИМАНИЕ"
		if l == lang.EN {
			warn = "WARNING"
		}
		if !strings.Contains(out, warn) {
			t.Errorf("[%s] a report backed by the stub must say so:\n%s", l, out)
		}
		// The разряд is an identifier of the решающий слой and is never
		// translated as a value — but the WORD the screen shows for it
		// has to be there, in the reader's language.
		for _, c := range core.Classes {
			if !strings.Contains(out, l.Word(string(c))) {
				t.Errorf("[%s] разряд %q missing from the report", l, c)
			}
		}
	}
}

// A snapshot from a system that publishes less than Linux does: the fields
// exist because the shape is shared, and each one says out loud that nobody
// measured it.  This is the rule the whole port stands on — an absent
// measurement must never print as a zero.
func TestStatusPrintsDashesForFactsTheSystemDoesNotPublish(t *testing.T) {
	// The reasons are Raw phrases: nobody translates the words of the
	// system, so the same sentence is looked for in both runs.
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
		Missing: map[string]lang.Phrase{
			sysinfo.FactCPUBusy:      lang.Raw("счётчики не дались"),
			sysinfo.FactThreads:      lang.Raw("самопроверка не сошлась"),
			sysinfo.FactBlocked:      lang.Raw("система не различает такие процессы"),
			sysinfo.FactNetCounters:  lang.Raw("самопроверка не сошлась"),
			sysinfo.FactSensors:      lang.Raw("система не публикует показания"),
			sysinfo.FactMemoryPages:  lang.Raw("самопроверка не сошлась"),
			sysinfo.FactNetTxDrops:   lang.Raw("система не считает такие пакеты"),
			sysinfo.FactLoadEntities: lang.Raw("система не публикует очередь"),
		},
	}

	for _, l := range both {
		out := status(l, st)
		for _, want := range []string{
			l.F("  всего         %s", l.UBytes(32<<30)), // measured, and printed
			l.T("  занято        —"),                    // not measured, and said so
			l.F("  свободно      %s", "—"),
			l.F("  доступно      %s", "—"),
			l.F("  кэш/буферы    %s", "—"),
			l.F("  своп          %s из %s занято", l.UBytes(1<<30), l.UBytes(2<<30)),
			l.T("  занято ЦП     —"),
			l.F("ПРОЦЕССЫ  %s", l.F("всего %d", 412)+", "+l.F("выполняется %d", 3)),
			strings.TrimSuffix(l.T("НЕ ИЗМЕРЕНО  %s"), "%s"),
		} {
			if !strings.Contains(out, want) {
				t.Errorf("[%s] the report must contain %q:\n%s", l, want, out)
			}
		}
		// A field nobody measured must never look like a measured zero.
		zero := strings.TrimSuffix(l.T("  занято        —"), "—") + l.UBytes(0)
		if strings.Contains(out, zero) || strings.Contains(out, l.F("потоков %d", 0)) {
			t.Errorf("[%s] an unmeasured field was printed as a measured zero:\n%s", l, out)
		}
		// The point of the whole rework: a reader looking for numbers is
		// not made to read an explanation of the ones that are not there.
		for _, reason := range []string{"счётчики не дались", "самопроверка не сошлась",
			"система не публикует показания", "система не различает такие процессы"} {
			if strings.Contains(out, reason) {
				t.Errorf("[%s] the report printed a reason instead of leaving it to --why: %q\n%s", l, reason, out)
			}
		}
		// A fact with no column in the report has no line in it either.
		if strings.Contains(out, l.Word(sysinfo.FactNetTxDrops)) || strings.Contains(out, l.Word(sysinfo.FactLoadEntities)) {
			t.Errorf("[%s] a fact the report never shows must not be named in it:\n%s", l, out)
		}
		if !strings.Contains(out, l.Word(sysinfo.FactSensors)) || !strings.Contains(out, l.Word(sysinfo.FactCPUBusy)) {
			t.Errorf("[%s] a fact the report would have shown must be named:\n%s", l, out)
		}
		// The interface is real and is listed; its counters are not measured.
		if !strings.Contains(out, "en0") {
			t.Errorf("[%s] a known interface must still be listed:\n%s", l, out)
		}
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "en0") && strings.Contains(line, l.UBytes(0)) {
				t.Errorf("[%s] unmeasured counters printed as zero bytes: %q", l, line)
			}
		}
	}
}

func TestStatusOrdersTheUnmeasuredLine(t *testing.T) {
	st := sysinfo.Status{Missing: map[string]lang.Phrase{
		"я": lang.Raw("1"), "а": lang.Raw("2"), "м": lang.Raw("3"),
	}}
	for _, l := range both {
		if out := status(l, st); !strings.Contains(out, l.F("НЕ ИЗМЕРЕНО  %s", "а, м, я")) {
			t.Errorf("[%s] the line must be ordered, so two runs can be compared:\n%s", l, out)
		}
	}
}

// The reasons did not vanish; they moved to where somebody who wants them asks
// for them.  That is the whole trade: the report keeps the names, --why keeps
// the sentences, and the JSON keeps both.
func TestWhyCarriesTheReasonsTheReportLeavesOut(t *testing.T) {
	st := sysinfo.Status{Missing: map[string]lang.Phrase{
		sysinfo.FactSensors:    lang.Raw("система не публикует показания датчиков"),
		sysinfo.FactNetTxDrops: lang.Raw("система не считает такие пакеты"),
	}}
	for _, l := range both {
		var buf bytes.Buffer
		Why(&buf, l, st)
		out := buf.String()
		for _, want := range []string{
			l.Word(sysinfo.FactSensors), "система не публикует показания датчиков",
			l.Word(sysinfo.FactNetTxDrops), "система не считает такие пакеты",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("[%s] --why must carry %q:\n%s", l, want, out)
			}
		}
		// The order is the order of the имена, not of their translations:
		// two runs of the same command must be comparable line by line.
		if i, j := strings.Index(out, l.Word(sysinfo.FactSensors)), strings.Index(out, l.Word(sysinfo.FactNetTxDrops)); i > j {
			t.Errorf("[%s] --why must be ordered by name:\n%s", l, out)
		}
	}
}

func TestWhySaysSoWhenNothingIsMissing(t *testing.T) {
	for _, l := range both {
		var buf bytes.Buffer
		Why(&buf, l, sysinfo.Status{})
		if out := buf.String(); !strings.Contains(out, l.T("НЕ ИЗМЕРЕНО  ничего: снимок полон")) {
			t.Errorf("[%s] a complete snapshot must say so:\n%s", l, out)
		}
	}
	// And it says it in English when English is asked for: a wording that
	// came back in Russian here would mean the dictionary was never asked.
	var buf bytes.Buffer
	Why(&buf, lang.EN, sysinfo.Status{})
	if out := buf.String(); !strings.Contains(out, "NOT MEASURED nothing: the snapshot is complete") {
		t.Errorf("English --why printed the Russian wording:\n%s", out)
	}
}

// A partial count must say how partial it is: on macOS a snapshot taken
// without administrator rights covers the caller's own processes only.
func TestStatusSaysHowManyProcessesTheThreadCountCovers(t *testing.T) {
	st := sysinfo.Status{Processes: sysinfo.Processes{
		Total: 906, Running: 4, Threads: 1240, WithDetail: 214,
	}}
	for _, l := range both {
		out := status(l, st)
		if !strings.Contains(out, l.F("замерено по %d процессам", 214)) {
			t.Errorf("[%s] a partial count must name its coverage:\n%s", l, out)
		}
	}
	if out := status(lang.EN, st); !strings.Contains(out, "measured over 214 processes") {
		t.Errorf("the coverage must be said in English too:\n%s", out)
	}

	st.Processes.WithDetail = st.Processes.Total
	for _, l := range both {
		out := status(l, st)
		if !strings.Contains(out, l.F("потоков %d", 1240)+",") ||
			strings.Contains(out, strings.TrimSuffix(l.T("замерено по %d процессам"), "%d процессам")) {
			t.Errorf("[%s] a complete count needs no qualifier:\n%s", l, out)
		}
	}
}

// TestАнглийскийСнимокБезКириллицы is the check no single wording can pass on
// its own: a snapshot rendered for an English reader must not hold one
// Cyrillic letter that this tool wrote.  Every value in it is ASCII on
// purpose, so anything Cyrillic left in the output is a wording that never
// reached the dictionary — the half-translated report the whole rework exists
// to prevent.
func TestАнглийскийСнимокБезКириллицы(t *testing.T) {
	st := sysinfo.Status{
		Host: sysinfo.Host{
			Hostname: "kolobok", Distro: "NetBSD 10.1", KernelRelease: "10.1",
			Machine: "riscv64", Model: "VisionFive 2",
			BootTime: 1756000000, UptimeSeconds: 90061,
		},
		Load: sysinfo.Load{CPUCount: 4, SampleMillis: 200},
		Memory: procfs.Memory{
			Total: 8 << 30, Used: 3 << 30, Free: 1 << 30, Available: 5 << 30,
			BuffCache: 1 << 30, SwapTotal: 2 << 30, SwapUsed: 1 << 30,
			Present: map[string]bool{
				procfs.FieldTotal: true, procfs.FieldUsed: true, procfs.FieldFree: true,
				procfs.FieldAvailable: true, procfs.FieldBuffCache: true,
				procfs.FieldSwapTotal: true, procfs.FieldSwapUsed: true,
			},
		},
		Processes: sysinfo.Processes{
			Total: 412, Running: 3, Threads: 900, WithDetail: 200, Unreadable: 2,
			TopByMemory: []sysinfo.Proc{{PID: 1, User: "root", Comm: "init", RSSBytes: 1 << 20}},
		},
		Disks: []sysinfo.Disk{{
			Source: "/dev/ld0a", MountPoint: "/", TotalBytes: 100 << 30,
			UsedBytes: 40 << 30, AvailableBytes: 60 << 30, UsePercent: 40,
		}},
		Network: []sysinfo.Iface{{
			NetCounters: procfs.NetCounters{Name: "en0", RxBytes: 1 << 20, TxBytes: 1 << 20, RxPackets: 1200, TxPackets: 900},
			OperState:   "up",
		}},
		Sensors: []sysinfo.Sensor{{Chip: "cpu0", Label: "core", Celsius: 41.5, CritC: 95}},
		Missing: map[string]lang.Phrase{sysinfo.FactSensors: lang.Raw("no readings")},
	}
	out := status(lang.EN, st)
	for i, line := range strings.Split(out, "\n") {
		for _, r := range line {
			if r >= 0x0400 && r <= 0x04FF {
				t.Fatalf("строка %d английского снимка осталась по-русски: %q", i+1, line)
			}
		}
	}
}

// The printed report gained a section and two lines, and nothing else moved.
// This test is the guard on that: the labels a script or a person has been
// reading since the first release are still spelled the way they were.
func TestOldLabelsOfTheStatusReportDidNotMove(t *testing.T) {
	var buf bytes.Buffer
	Status(&buf, lang.RU, sysinfo.Status{})
	out := buf.String()
	for _, want := range []string{
		"СИСТЕМА", "  узел          ", "  дистрибутив   ", "  ядро          ", "  время работы  ",
		"ЗАГРУЗКА", "  средняя       ", "  ядер          ", "  занято ЦП     ",
		"ПАМЯТЬ", "ПРОЦЕССЫ", "ДИСКИ", "СЕТЬ", "ТЕМПЕРАТУРА",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("из отчёта пропало %q:\n%s", want, out)
		}
	}
	// The new section is printed even when the machine has no card, the way
	// ДИСКИ and СЕТЬ are: a section that disappears looks like a tool that
	// did not look.
	if !strings.Contains(out, "ВИДЕОКАРТЫ\n  —") {
		t.Errorf("раздел видеокарт не напечатан прочерком:\n%s", out)
	}
	// And a machine nobody measured must not gain a per-core line.
	if strings.Contains(out, "по ядрам") {
		t.Errorf("строка по ядрам напечатана без замера:\n%s", out)
	}
}

func TestStatusPrintsCardsWithDashesAndTheirOrigin(t *testing.T) {
	busy, celsius := 99.0, 85.0
	used, total := uint64(29)<<30, uint64(48)<<30
	share := 12.5
	st := sysinfo.Status{
		Load: sysinfo.Load{
			CPUCount: 4, BusyPercent: &share, SampleMillis: 200,
			Cores: []sysinfo.Core{
				{Index: 0, BusyPercent: &busy},
				{Index: 1, BusyPercent: &share},
				{Index: 2},
			},
		},
		Host: sysinfo.Host{CPUModel: "Придуманный процессор 9000"},
		GPUs: []gpuinfo.Card{
			{Name: "NVIDIA RTX 6000 Ada Generation", Bus: "0000:02:00.0", Driver: "nvidia",
				BusyPercent: &busy, Celsius: &celsius, MemoryUsedBytes: &used, MemoryTotalBytes: &total,
				Source: "nvidia-smi", Outside: true},
			{Name: "Matrox G200eW3", Bus: "0000:62:00.0", Driver: "mgag200",
				Source: "/sys/class/drm/card1/device"},
		},
	}
	var buf bytes.Buffer
	Status(&buf, lang.RU, st)
	out := buf.String()

	for _, want := range []string{
		"процессор     Придуманный процессор 9000 × 4",
		"по ядрам      мин 12,5% / медиана 55,8% / макс 99,0% (ядро 0); занято больше половины 1 из 3",
		"NVIDIA RTX 6000 Ada Generation",
		"99,0%",
		"29,0 ГиБ из 48,0 ГиБ",
		"85,0°C",
		"числа от чужой программы nvidia-smi",
		"числа из /sys/class/drm/card1/device",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("в отчёте нет %q:\n%s", want, out)
		}
	}
	// The silent card is four dashes, not four zeros.
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "Matrox") && (strings.Contains(line, "0.0%") || strings.Contains(line, "0 Б")) {
			t.Errorf("немая карта напечатана нулями: %q", line)
		}
	}
}
