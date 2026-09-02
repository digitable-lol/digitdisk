// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/cli"
	"digitdisk/internal/lang"
	"digitdisk/internal/procfs"
	"digitdisk/internal/sysinfo"
)

// sgr matches one colour sequence, so tests can look at what is drawn rather
// than at how it is painted.
var sgr = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func plain(s string) string { return sgr.ReplaceAllString(s, "") }

func TestFitPadsAndCutsByRunes(t *testing.T) {
	// Russian is what the screen speaks: a byte count would cut a letter in
	// half and leave a broken rune on the terminal.
	if got := fit("память", 10); got != "память    " {
		t.Errorf("fit padded to %q", got)
	}
	if got := fit("температура", 6); got != "темпе…" {
		t.Errorf("fit cut to %q, want темпе…", got)
	}
	if got := runes(fit("температура", 6)); got != 6 {
		t.Errorf("cut left %d cells, want 6", got)
	}
	if got := fit("ok", 0); got != "" {
		t.Errorf("zero width gave %q", got)
	}
}

func TestRightAligns(t *testing.T) {
	if got := right("42", 5); got != "   42" {
		t.Errorf("right = %q", got)
	}
	if got := runes(right("длинновато", 4)); got != 4 {
		t.Errorf("right overran to %d cells", got)
	}
}

func TestPlainWidthIgnoresEscapes(t *testing.T) {
	th := NewTheme(Carbon)
	painted := th.Fg(th.P.Accent, "живой")
	if strings.Contains(painted, "\x1b") == false {
		t.Skip("терминал этого прогона без цвета")
	}
	if got := plainWidth(painted); got != 5 {
		t.Errorf("plainWidth = %d, want 5 (экранные последовательности ширины не имеют)", got)
	}
}

func TestClipKeepsLinesInsideTheTerminal(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	var r row
	r.add("занято ", func(x string) string { return th.Fg(th.P.Muted, x) })
	r.add("220.3 ГиБ из 499.1 ГиБ", func(x string) string { return th.Fg(th.P.Foreground, x) })

	// A line wider than the terminal wraps, pushes every line below it down
	// and takes the whole layout with it.  Clipping is what stops that.
	for _, n := range []int{1, 5, 12, 28, 200} {
		got := th.clip(r.String(), n)
		if w := plainWidth(got); w > n {
			t.Errorf("clip(%d) оставил %d ячеек: %q", n, w, plain(got))
		}
	}
	if w := plainWidth(th.clip(r.String(), 200)); w != r.w {
		t.Errorf("clip шире строки её укоротил")
	}
	if !strings.Contains(plain(th.clip(r.String(), 12)), "…") {
		t.Error("обрезанная строка не помечена многоточием")
	}
}

func TestClipDoesNotBreakAnEscapeSequence(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	line := th.Fg(th.P.Accent, "абв") + th.Fg(th.P.Red, "где")
	for n := 1; n <= 8; n++ {
		got := th.clip(line, n)
		// Every ESC in the result must be a whole sequence: a half-written
		// one is printed as rubbish by the terminal.
		if regexp.MustCompile(`\x1b\[[0-9;]*$`).MatchString(got) {
			t.Errorf("clip(%d) оборвал последовательность: %q", n, got)
		}
		if !strings.HasSuffix(got, "m") && n < 6 {
			t.Errorf("clip(%d) не вернул цвет холста: %q", n, got)
		}
	}
}

func TestBarAndSparkAreExactlyAsWideAsAsked(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	for _, w := range []int{1, 8, 20, 40} {
		if got := plainWidth(th.bar(0.42, w)); got != w {
			t.Errorf("bar(%d) шириной %d", w, got)
		}
		if got := plainWidth(th.emptyBar(w)); got != w {
			t.Errorf("emptyBar(%d) шириной %d", w, got)
		}
		if got := plainWidth(th.spark([]float64{0.1, 0.9, -1}, w)); got != w {
			t.Errorf("spark(%d) шириной %d", w, got)
		}
	}
	// A share outside 0..1 must not widen the gauge past its cell count.
	if got := plainWidth(th.bar(9, 10)); got != 10 {
		t.Errorf("bar с долей 9 шириной %d", got)
	}
	if got := plainWidth(th.bar(-3, 10)); got != 10 {
		t.Errorf("bar с долей -3 шириной %d", got)
	}
}

func TestSparkDrawsAGapForAnUnmeasuredSample(t *testing.T) {
	th := Theme{P: Carbon, d: depthNone}
	// -1 is "never measured".  It must not be drawn as the lowest bar,
	// which would read as a measured zero.
	got := plain(th.spark([]float64{-1, -1}, 2))
	if strings.ContainsAny(got, string(sparkFaces)) {
		t.Errorf("незамеренная доля нарисована столбиком: %q", got)
	}
}

func TestLevelThresholds(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	cases := []struct {
		frac float64
		want slot
	}{{0, Carbon.Green}, {0.74, Carbon.Green}, {0.75, Carbon.Yellow}, {0.89, Carbon.Yellow}, {0.9, Carbon.Red}, {1, Carbon.Red}}
	for _, c := range cases {
		if got := th.level(c.frac); got != c.want {
			t.Errorf("level(%.2f) = %v, want %v", c.frac, got, c.want)
		}
	}
}

func TestDetectDepthHonoursTheEnvironment(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want depth
	}{
		{"NO_COLOR молчит", map[string]string{"NO_COLOR": "1", "TERM": "xterm-256color", "COLORTERM": "truecolor"}, depthNone},
		{"TERM=dumb молчит", map[string]string{"TERM": "dumb"}, depthNone},
		{"пустой TERM молчит", map[string]string{"TERM": ""}, depthNone},
		{"COLORTERM=truecolor", map[string]string{"TERM": "xterm", "COLORTERM": "truecolor"}, depthTrue},
		{"256 в имени TERM", map[string]string{"TERM": "screen-256color"}, depth256},
		{"обычный TERM", map[string]string{"TERM": "xterm"}, depth16},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, k := range []string{"NO_COLOR", "TERM", "COLORTERM"} {
				t.Setenv(k, "")
				os.Unsetenv(k)
			}
			for k, v := range c.env {
				t.Setenv(k, v)
			}
			if _, ok := c.env["NO_COLOR"]; !ok {
				os.Unsetenv("NO_COLOR")
			}
			if got := detectDepth(); got != c.want {
				t.Errorf("detectDepth = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUsableTERM(t *testing.T) {
	t.Setenv("TERM", "dumb")
	if UsableTERM() {
		t.Error("TERM=dumb признан пригодным")
	}
	t.Setenv("TERM", "")
	if UsableTERM() {
		t.Error("пустой TERM признан пригодным")
	}
	t.Setenv("TERM", "xterm-256color")
	if !UsableTERM() {
		t.Error("xterm-256color признан непригодным")
	}
}

func TestAvailableSaysNoToAPipeAndToAFile(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	if Available(w) {
		t.Error("живой экран разрешён в трубу — вывод `digitdisk status | cat` был бы испорчен")
	}

	f, err := os.CreateTemp(t.TempDir(), "out")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if Available(f) {
		t.Error("живой экран разрешён в файл")
	}
	if IsTerminal(nil) {
		t.Error("nil признан терминалом")
	}

	// /dev/null carries the character-device bit and is the reason the bit
	// alone is not the test.
	if dn, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0); err == nil {
		defer dn.Close()
		if Available(dn) {
			t.Error("живой экран разрешён в /dev/null")
		}
	}
}

func TestPaletteByNameFallsBackToCarbon(t *testing.T) {
	if PaletteByName("paper").Name != Paper.Name {
		t.Error("paper не выбралась")
	}
	if PaletteByName(" SIGNAL ").Name != Signal.Name {
		t.Error("signal не выбралась без учёта регистра и пробелов")
	}
	for _, n := range []string{"", "нетакой", "carbon"} {
		if PaletteByName(n).Name != Carbon.Name {
			t.Errorf("%q не дало carbon", n)
		}
	}
}

func TestCarbonMatchesTheBrandCanon(t *testing.T) {
	// The values come from the stack's palette file, Digitable Focus Carbon.
	// Written down here so a stray edit to the accent is a failing test and
	// not a screen that quietly stops looking like Digitable.
	for _, c := range []struct {
		role string
		got  RGB
		want RGB
	}{
		{"background", Carbon.Background.c, RGB{0x05, 0x08, 0x0D}},
		{"foreground", Carbon.Foreground.c, RGB{0xF5, 0xF7, 0xFA}},
		{"accent", Carbon.Accent.c, RGB{0x00, 0xE5, 0xE5}},
		{"accentSoft", Carbon.AccentSoft.c, RGB{0x00, 0xD8, 0xFF}},
		{"muted", Carbon.Muted.c, RGB{0x9B, 0xAA, 0xB8}},
		{"green", Carbon.Green.c, RGB{0x7C, 0xFF, 0x6B}},
		{"yellow", Carbon.Yellow.c, RGB{0xFF, 0xC2, 0x47}},
		{"red", Carbon.Red.c, RGB{0xFF, 0x5B, 0x5B}},
	} {
		if c.got != c.want {
			t.Errorf("%s = %v, канон %v", c.role, c.got, c.want)
		}
	}
	// The ANSI fallback follows the stack mapping: cyan = accent,
	// brightCyan = accentSoft, magenta = purple, white = muted.
	if Carbon.Accent.ansi != 36 || Carbon.AccentSoft.ansi != 96 ||
		Carbon.Purple.ansi != 35 || Carbon.Muted.ansi != 37 {
		t.Error("раскладка на 16 цветов разошлась с каноном стека")
	}
}

func TestTrueColorSequenceIsTheBrandColour(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	if got := th.seq(Carbon.Accent, false); got != "\x1b[38;2;0;229;229m" {
		t.Errorf("accent = %q", got)
	}
	if got := th.seq(Carbon.Background, true); got != "\x1b[48;2;5;8;13m" {
		t.Errorf("фон = %q", got)
	}
	th16 := Theme{P: Carbon, d: depth16}
	if got := th16.seq(Carbon.Accent, true); got != "\x1b[46m" {
		t.Errorf("фон в 16 цветах = %q", got)
	}
}

func TestResetLaysTheCanvasBackDown(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	// ESC[0m alone would drop the brand background together with the colour,
	// and ESC[K would then erase the rest of the line in the terminal's own.
	if !strings.HasSuffix(th.Fg(Carbon.Accent, "x"), th.Canvas()) {
		t.Error("после окрашенного куска холст не восстановлен")
	}
	if (Theme{P: Carbon, d: depthNone}).Canvas() != "" {
		t.Error("бесцветный режим всё-таки красит фон")
	}
}

func TestCube256StaysInRange(t *testing.T) {
	for _, c := range []RGB{{0, 0, 0}, {255, 255, 255}, {0, 229, 229}, {5, 8, 13}, {155, 170, 184}} {
		if n := cube256(c); n < 16 || n > 255 {
			t.Errorf("cube256(%v) = %d вне палитры", c, n)
		}
	}
}

// filled is a snapshot with something in every section.
func filled() sysinfo.Status {
	busy, cpu := 42.5, 7.5
	return sysinfo.Status{
		TakenAt: "2026-09-02T00:00:00Z",
		Host: sysinfo.Host{
			Hostname: "узел", Distro: "Debian GNU/Linux 13", KernelRelease: "6.12.0",
			Machine: "x86_64", KernelVersion: "#1 SMP", UptimeSeconds: 90061,
			UptimeHuman: "1д 01:01", BootTime: 1756000000,
		},
		Load: sysinfo.Load{
			LoadAvg:  procfs.LoadAvg{One: 1, Five: 2, Fifteen: 3, Runnable: 2, TotalEntities: 300},
			CPUCount: 8, BusyPercent: &busy, SampleMillis: 200,
		},
		Memory: procfs.Memory{
			Total: 16 << 30, Free: 4 << 30, Available: 8 << 30, BuffCache: 3 << 30,
			Shared: 1 << 30, Used: 8 << 30, SwapTotal: 2 << 30, SwapUsed: 1 << 30,
		},
		Processes: sysinfo.Processes{
			Total: 300, Running: 2, Blocked: 1, Threads: 900, Unreadable: 3,
			TopByMemory: []sysinfo.Proc{{PID: 1, User: "root", Cmdline: "/sbin/init", RSSBytes: 1 << 20}},
			TopByCPU:    []sysinfo.Proc{{PID: 2, User: "пользователь", Comm: "kthreadd", CPUPercent: &cpu}, {PID: 3, Comm: "без замера"}},
		},
		Disks: []sysinfo.Disk{
			{Source: "/dev/sda1", MountPoint: "/", FSType: "ext4", TotalBytes: 1 << 40, UsedBytes: 9 << 36, AvailableBytes: 1 << 39, UsePercent: 56.25, InodesTotal: 100, InodesFree: 40},
			{Source: "нет", MountPoint: "/mnt/сломано", Error: "нет доступа"},
			{Source: "ro", MountPoint: "/boot/efi", FSType: "vfat", ReadOnly: true, TotalBytes: 1 << 20, AvailableBytes: 1 << 19, UsePercent: 50},
		},
		Network: []sysinfo.Iface{
			{NetCounters: procfs.NetCounters{Name: "eth0", RxBytes: 1 << 30, TxBytes: 1 << 29, RxPackets: 10, TxPackets: 20, RxErrors: 1, TxDropped: 2},
				MAC: "00:11:22:33:44:55", OperState: "up", MTU: 1500, SpeedMbit: 1000, Addresses: []string{"10.0.0.2/24"}},
			{NetCounters: procfs.NetCounters{Name: "lo"}, OperState: "down"},
		},
		Sensors: []sysinfo.Sensor{
			{Chip: "coretemp", Label: "Package id 0", Celsius: 55.5, CritC: 100},
			{Chip: "acpi", Celsius: 40},
		},
		Missing: map[string]lang.Phrase{
			"я": lang.Raw("нет доступа"), "б": lang.Raw("нет файла"), "а": lang.Raw("не разобрано"),
		},
	}
}

// newTestScreen builds a screen without a terminal: the sections are pure
// functions of a snapshot and are tested as such.
func newTestScreen(st sysinfo.Status, have bool, cols int) *screen {
	s := &screen{t: Theme{P: Carbon, d: depthTrue}, rows: 30, cols: cols, l: lang.RU}
	s.st, s.haveSt = st, have
	s.taken, s.took = time.Now(), 1900*time.Millisecond
	s.o.Interval = 2 * time.Second
	s.cpuHist = []float64{0.1, 0.5, -1, 0.9}
	s.memHist = []float64{0.3, 0.4}
	return s
}

func TestEverySectionDrawsWithinTheTerminal(t *testing.T) {
	for _, st := range []struct {
		name string
		val  sysinfo.Status
		have bool
	}{
		{"полный снимок", filled(), true},
		{"пустой снимок", sysinfo.Status{}, true},
		{"замер ещё идёт", sysinfo.Status{}, false},
	} {
		for _, cols := range []int{40, 60, 80, 120, 200} {
			s := newTestScreen(st.val, st.have, cols)
			for i, sec := range sections {
				s.tab, s.scroll = i, 0
				for _, line := range s.frame() {
					if w := plainWidth(s.t.clip(line, cols)); w > cols {
						t.Errorf("%s / %s / %d колонок: строка в %d ячеек", st.name, sec.title(s.l), cols, w)
					}
				}
			}
		}
	}
}

func TestEmptySnapshotNeverInventsAZero(t *testing.T) {
	s := newTestScreen(sysinfo.Status{}, true, 100)
	for i, sec := range sections {
		s.tab, s.scroll = i, 0
		body := plain(strings.Join(sec.render(s), "\n"))
		// «0,0%» по-русски и «0.0%» по-английски — один и тот же
		// придуманный ноль, записанный двумя разделителями.
		if strings.Contains(body, "0.0%") || strings.Contains(body, "0,0%") || strings.Contains(body, "0 Б из 0 Б") {
			t.Errorf("раздел %s выдал измеренный ноль там, где ничего не читали:\n%s", sec.title(s.l), body)
		}
	}
	// The CPU share was never sampled, so the gauge must say so.
	s.tab = 2
	if got := plain(strings.Join(s.load(), "\n")); !strings.Contains(got, "замер не делался") {
		t.Errorf("ЗАГРУЗКА без замера не сказала об этом:\n%s", got)
	}
	// Memory that was never read is a dash, not an empty bar at zero.
	s.tab = 3
	if got := plain(strings.Join(s.memory(), "\n")); !strings.Contains(got, dash) {
		t.Errorf("ПАМЯТЬ без чисел не показала прочерк:\n%s", got)
	}
}

func TestUnmeasuredProcessCPUIsADash(t *testing.T) {
	s := newTestScreen(filled(), true, 120)
	got := plain(strings.Join(s.processes(), "\n"))
	if !strings.Contains(got, "без замера") {
		t.Fatalf("процесс без доли ЦП не показан:\n%s", got)
	}
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "без замера") && !strings.Contains(line, dash) {
			t.Errorf("процесс без замера получил число вместо прочерка: %q", line)
		}
	}
}

func TestMissingSectionIsSortedSoTheScreenDoesNotShuffle(t *testing.T) {
	// Go walks a map in a different order every time; without sorting the
	// section would reshuffle itself once a second under the reader's eyes.
	s := newTestScreen(filled(), true, 100)
	first := plain(strings.Join(s.missing(), "\n"))
	for i := 0; i < 20; i++ {
		if got := plain(strings.Join(s.missing(), "\n")); got != first {
			t.Fatalf("раздел НЕ ПРОЧИТАНО переставился между отрисовками:\n%s\n---\n%s", first, got)
		}
	}
	order := strings.Index(first, "а") < strings.Index(first, "б")
	if !order {
		t.Error("источники не отсортированы")
	}
}

func TestSectionsAreTheSectionsOfThePrintedReport(t *testing.T) {
	// The screen is a second view of the same report, not a second tool.
	want := []string{"ОБЗОР", "СИСТЕМА", "ЗАГРУЗКА", "ПАМЯТЬ", "ПРОЦЕССЫ", "ДИСКИ", "СЕТЬ",
		"ТЕМПЕРАТУРА", "ВИДЕОКАРТЫ", "НЕ ПРОЧИТАНО"}
	if len(sections) != len(want) {
		t.Fatalf("разделов %d, ждали %d", len(sections), len(want))
	}
	for i, w := range want {
		if got := sections[i].title(lang.RU); got != w {
			t.Errorf("раздел %d = %q, ждали %q", i, got, w)
		}
	}
	// The keys did not change, so what the digits reach did not either: 1…9
	// are the first nine sections.  The section added since — ВИДЕОКАРТЫ —
	// went in as the ninth, after the last one the report prints, so every
	// section the report has still has a digit; and the tenth is НЕ
	// ПРОЧИТАНО, which is named on the opening page and is one ← away from
	// it.  A section list longer than this one would leave a reading behind
	// a key nobody can guess.
	if len(sections) > 10 {
		t.Error("разделов больше десяти — до последних клавишами 1…9 уже не добраться")
	}
	if sections[len(sections)-1].title(lang.RU) != "НЕ ПРОЧИТАНО" {
		t.Error("без цифры остался не тот раздел: последним должен быть НЕ ПРОЧИТАНО")
	}
}

func TestHeaderCarriesTheBrand(t *testing.T) {
	s := newTestScreen(filled(), true, 120)
	head := plain(s.header())
	for _, want := range []string{"◇", "digitdisk", "Digitable", "узел"} {
		if !strings.Contains(head, want) {
			t.Errorf("в шапке нет %q: %q", want, head)
		}
	}
	if !strings.Contains(s.header(), "\x1b[48;2;0;229;229m") {
		t.Error("метка шапки не выкрашена фирменным акцентом")
	}
}

func TestFooterAlwaysSaysHowToLeave(t *testing.T) {
	// A full-screen program that does not say how to quit is a trap.
	for _, cols := range []int{20, 30, 40, 60, 80, 120, 200} {
		s := newTestScreen(filled(), true, cols)
		if got := plain(s.footer("")); !strings.Contains(got, "q") {
			t.Errorf("%d колонок: подвал не назвал выход: %q", cols, got)
		}
	}
}

func TestKeysMoveBetweenSectionsAndOut(t *testing.T) {
	s := newTestScreen(filled(), true, 100)
	ch := make(chan sample, 1)

	if s.handle(key{kind: keyRight}, ch); s.tab != 1 {
		t.Errorf("вправо дало раздел %d", s.tab)
	}
	if s.handle(key{kind: keyLeft}, ch); s.tab != 0 {
		t.Errorf("влево дало раздел %d", s.tab)
	}
	// Left from the first section wraps to the last rather than going nowhere.
	if s.handle(key{kind: keyLeft}, ch); s.tab != len(sections)-1 {
		t.Errorf("влево с первого не завернуло: раздел %d", s.tab)
	}
	if s.handle(key{kind: keyRight}, ch); s.tab != 0 {
		t.Errorf("вправо с последнего не завернуло: раздел %d", s.tab)
	}
	if s.handle(key{kind: keyRune, r: '5'}, ch); s.tab != 4 {
		t.Errorf("клавиша 5 дала раздел %d", s.tab)
	}
	if s.handle(key{kind: keyRune, r: 'p'}, ch); !s.paused {
		t.Error("p не поставила паузу")
	}
	if s.handle(key{kind: keyRune, r: ' '}, ch); s.paused {
		t.Error("пробел не снял паузу")
	}
	// Scrolling never goes above the first line.
	s.handle(key{kind: keyUp}, ch)
	s.handle(key{kind: keyPgUp}, ch)
	if s.scroll != 0 {
		t.Errorf("прокрутка ушла в минус: %d", s.scroll)
	}
	for _, k := range []key{{kind: keyRune, r: 'q'}, {kind: keyRune, r: 'й'}, {kind: keyEsc}, {kind: keyCtrlC}} {
		if !s.handle(k, ch) {
			t.Errorf("клавиша %+v не закрыла экран", k)
		}
	}
	if s.handle(key{kind: keyRune, r: 'x'}, ch) {
		t.Error("посторонняя клавиша закрыла экран")
	}
}

func TestEscapeSequencesBecomeKeys(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	// Arrows, page keys, shift-tab, a Russian letter and a bare Escape.
	if _, err := w.Write([]byte("\x1b[A\x1b[B\x1b[C\x1b[D\x1bOA\x1b[Z\x1b[5~\x1b[6~й\x09\x03")); err != nil {
		t.Fatal(err)
	}
	w.Close()

	ch := make(chan key, 32)
	readKeys(r, ch)
	var got []kind
	for k := range ch {
		got = append(got, k.kind)
	}
	want := []kind{keyUp, keyDown, keyRight, keyLeft, keyUp, keyShiftTab, keyPgUp, keyPgDn, keyRune, keyTab, keyCtrlC}
	if len(got) != len(want) {
		t.Fatalf("разобрано %v, ждали %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("клавиша %d = %v, ждали %v", i, got[i], want[i])
		}
	}
}

func TestRussianLetterArrivesWhole(t *testing.T) {
	r, w, _ := os.Pipe()
	w.Write([]byte("йщ"))
	w.Close()
	ch := make(chan key, 8)
	readKeys(r, ch)
	var runes []rune
	for k := range ch {
		runes = append(runes, k.r)
	}
	if string(runes) != "йщ" {
		t.Errorf("двухбайтовые буквы разобраны как %q", string(runes))
	}
}

func TestScrollStopsAtTheEndOfTheSection(t *testing.T) {
	s := newTestScreen(filled(), true, 100)
	s.tab = 5 // ДИСКИ
	s.scroll = 10000
	lines := s.frame()
	if len(lines) != s.rows {
		t.Errorf("кадр в %d строк при экране в %d", len(lines), s.rows)
	}
	if s.scroll >= len(s.disks()) {
		t.Errorf("прокрутка %d ушла за конец раздела в %d строк", s.scroll, len(s.disks()))
	}
}

func TestFrameIsExactlyTheHeightOfTheTerminal(t *testing.T) {
	for _, rows := range []int{8, 12, 24, 60} {
		s := newTestScreen(filled(), true, 100)
		s.rows = rows
		for i := range sections {
			s.tab, s.scroll = i, 0
			if got := len(s.frame()); got != rows {
				t.Errorf("раздел %s на %d строк дал кадр в %d", sections[i].title(s.l), rows, got)
			}
		}
	}
}

func TestRunRefusesWhereThereIsNoTerminal(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()
	err := Run(Options{Out: w, Collect: func() sysinfo.Status { return filled() }})
	if err != ErrNoTerminal {
		t.Errorf("Run в трубу вернул %v, ждали ErrNoTerminal", err)
	}
	t.Setenv("TERM", "dumb")
	if err := Run(Options{Out: os.Stdout, Collect: func() sysinfo.Status { return filled() }}); err != ErrNoTerminal {
		t.Errorf("Run при TERM=dumb вернул %v, ждали ErrNoTerminal", err)
	}
}

func TestRunNeedsACollector(t *testing.T) {
	if err := Run(Options{Out: os.Stdout}); err == nil {
		t.Error("Run без сборщика не пожаловался")
	}
}

func TestSensorGaugeSaysWhatItIsMeasuredAgainst(t *testing.T) {
	s := newTestScreen(filled(), true, 120)
	got := plain(strings.Join(s.sensors(), "\n"))
	if !strings.Contains(got, "критич. 100,0 °C") {
		t.Errorf("датчик с критической точкой не назвал её:\n%s", got)
	}
	if !strings.Contains(got, "из 100 °C") {
		t.Errorf("датчик без критической точки не назвал, против чего меряется:\n%s", got)
	}
}

func TestDiskErrorIsShownAsAnError(t *testing.T) {
	s := newTestScreen(filled(), true, 120)
	got := plain(strings.Join(s.disks(), "\n"))
	if !strings.Contains(got, "ошибка: нет доступа") {
		t.Errorf("недоступная точка монтирования показана без ошибки:\n%s", got)
	}
	if !strings.Contains(got, "только чтение") {
		t.Errorf("том только для чтения не помечен:\n%s", got)
	}
}

// The screen prints the same process sentence the report prints, and for the
// same reason: a count nobody could take must not appear on it as a zero.
// macOS is where this bites — the kernel answers about the caller's own
// processes only, unless the caller is the administrator.
func TestScreenDoesNotInventProcessCounts(t *testing.T) {
	st := sysinfo.Status{
		Processes: sysinfo.Processes{Total: 906, Running: 5, Threads: 1487, WithDetail: 214},
		Missing:   map[string]lang.Phrase{sysinfo.FactBlocked: lang.Raw("система не различает такие процессы")},
	}
	s := newTestScreen(st, true, 120)
	for _, got := range []string{
		plain(strings.Join(s.overview(), "\n")),
		plain(strings.Join(s.processes(), "\n")),
	} {
		if strings.Contains(got, "заблокировано 0") {
			t.Errorf("незамеренный счётчик напечатан нулём:\n%s", got)
		}
		if !strings.Contains(got, "замерено по 214 процессам") {
			t.Errorf("неполный счёт не назвал охват:\n%s", got)
		}
	}
}

// Список команд на экране — тот же список, из которого строятся справка и
// страница руководства. Если он разойдётся, разойдутся все три.
func TestCommandsPageNamesEveryCommand(t *testing.T) {
	s := newTestScreen(filled(), true, 100)
	s.menu = true
	page := strings.Join(s.frame(), "\n")
	for _, c := range cli.Commands {
		if !strings.Contains(plain(page), c.Name) {
			t.Errorf("экран не называет подкоманду %q", c.Name)
		}
	}
	if !strings.Contains(plain(page), "ничего не запускает") {
		t.Error("экран не говорит, что список только называет команды")
	}
	for _, cols := range []int{40, 60, 80, 120, 200} {
		s := newTestScreen(filled(), true, cols)
		s.menu = true
		for _, line := range s.frame() {
			if w := plainWidth(s.t.clip(line, cols)); w > cols {
				t.Errorf("список команд на %d колонках: строка в %d ячеек", cols, w)
			}
		}
	}
}

// «?» открывает и закрывает список, Esc из списка возвращает на экран, а не
// закрывает программу.
func TestQuestionMarkOpensTheCommandsAndEscBacksOut(t *testing.T) {
	s := newTestScreen(filled(), true, 100)
	ch := make(chan sample, 1)
	if s.handle(key{kind: keyRune, r: '?'}, ch); !s.menu {
		t.Fatal("«?» не открыла список команд")
	}
	if s.handle(key{kind: keyEsc}, ch) {
		t.Error("Esc из списка закрыл экран, а должен был вернуть на снимок")
	}
	if s.menu {
		t.Error("Esc не закрыл список")
	}
	s.handle(key{kind: keyRune, r: '?'}, ch)
	if s.handle(key{kind: keyRune, r: '?'}, ch); s.menu {
		t.Error("повторная «?» не закрыла список")
	}
	s.handle(key{kind: keyRune, r: '?'}, ch)
	if s.handle(key{kind: keyRight}, ch); s.menu {
		t.Error("переход к разделу оставил список открытым")
	}
	if !s.handle(key{kind: keyEsc}, ch) {
		t.Error("Esc вне списка не закрыл экран")
	}
}

// Экран рисуется на обоих языках и на обоих остаётся в своих колонках.
//
// Английская подпись длиннее русской — не опечатка в словаре, а поехавшая
// вёрстка: полоса разделов не помещается в строку, заголовок таблицы уезжает
// из-под чисел. Поэтому оба языка проходят ту же проверку ширины, что и один.
func TestScreenDrawsInBothLanguages(t *testing.T) {
	for _, l := range []lang.Lang{lang.RU, lang.EN} {
		for _, cols := range []int{40, 60, 80, 120, 200} {
			s := newTestScreen(filled(), true, cols)
			s.l = l
			for i, sec := range sections {
				s.tab, s.scroll = i, 0
				for _, line := range s.frame() {
					if w := plainWidth(s.t.clip(line, cols)); w > cols {
						t.Errorf("%s / %s / %d колонок: строка в %d ячеек", l, sec.title(l), cols, w)
					}
				}
			}
			s.tab, s.scroll, s.menu = 0, 0, true
			for _, line := range s.frame() {
				if w := plainWidth(s.t.clip(line, cols)); w > cols {
					t.Errorf("%s / список команд / %d колонок: строка в %d ячеек", l, cols, w)
				}
			}
		}
	}

	// Английский экран должен быть английским: имена разделов и подвал
	// приходят из словаря, а не остаются русскими.
	if got := sections[0].title(lang.EN); got != "OVERVIEW" {
		t.Errorf("первый раздел по-английски = %q", got)
	}
	s := newTestScreen(filled(), true, 120)
	s.l = lang.EN
	if got := plain(s.footer("")); !strings.Contains(got, "q quit") {
		t.Errorf("английский подвал не назвал выход: %q", got)
	}
	// Подписи — не значения: имя узла, команда процесса и причина отказа
	// приходят из снимка и по-английски выглядят так же, как по-русски.
	// Поэтому ищутся слова, которые бывают только подписью.
	labels := []string{"занято", "свободно", "доступно", "владелец", "интерфейс",
		"точка монтирования", "история", "средняя", "прочитано всё"}
	for _, sec := range sections {
		body := plain(strings.Join(sec.render(s), "\n"))
		for _, ru := range labels {
			if strings.Contains(body, ru) {
				t.Errorf("раздел %s по-английски оставил подпись %q:\n%s", sec.title(lang.RU), ru, body)
			}
		}
	}
	if got := plain(strings.Join(s.memory(), "\n")); !strings.Contains(got, "available") {
		t.Errorf("английская MEMORY без английских подписей:\n%s", got)
	}
}

// Клавиша языка переключает экран и просит запомнить выбор — ровно один раз на
// нажатие, потому что запись в домашний каталог не должна случаться дважды за
// одно движение пальцем.
func TestLanguageKeySwitchesTheScreenAndIsRemembered(t *testing.T) {
	s := newTestScreen(filled(), true, 100)
	ch := make(chan sample, 1)
	var asked []lang.Lang
	s.o.Remember = func(l lang.Lang) lang.Phrase {
		asked = append(asked, l)
		return lang.Raw("язык сохранён")
	}
	// Обе раскладки: «l» и «д» — одна и та же клавиша.
	for _, k := range []rune{'l', 'L', 'д', 'Д'} {
		was := s.l
		if s.handle(key{kind: keyRune, r: k}, ch) {
			t.Fatalf("клавиша %q закрыла экран", string(k))
		}
		if s.l != was.Other() {
			t.Errorf("клавиша %q дала язык %q, а из %q ждали %q", string(k), s.l, was, was.Other())
		}
	}
	if len(asked) != 4 {
		t.Fatalf("Remember позвали %d раз(а), ждали 4 — по разу на нажатие", len(asked))
	}
	if asked[0] != lang.EN || asked[1] != lang.RU {
		t.Errorf("Remember просили запомнить %v", asked)
	}

	// Ответ хранилища живёт в подвале несколько секунд: заметка, которая
	// никогда не уходит, перестаёт читаться.
	if got := plain(strings.Join(s.frame(), "\n")); !strings.Contains(got, "язык сохранён") {
		t.Errorf("ответ о сохранении не показан:\n%s", got)
	}
	s.saidAt = time.Now().Add(-time.Minute)
	if got := plain(strings.Join(s.frame(), "\n")); strings.Contains(got, "язык сохранён") {
		t.Error("ответ о сохранении остался в подвале навсегда")
	}
}

// Запоминать выбор может быть некуда — тогда клавиша всё равно переключает
// язык, и экран об этом молчит, а не падает.
func TestLanguageKeyWorksWithNowhereToRememberIt(t *testing.T) {
	s := newTestScreen(filled(), true, 100)
	s.o.Remember = nil
	ch := make(chan sample, 1)
	if s.handle(key{kind: keyRune, r: 'l'}, ch) {
		t.Fatal("клавиша языка закрыла экран")
	}
	if s.l != lang.EN {
		t.Errorf("без Remember язык не переключился: %q", s.l)
	}
	if !s.said.Empty() {
		t.Error("без Remember в подвале появилась запись о сохранении")
	}
	if got := len(s.frame()); got != s.rows {
		t.Errorf("после переключения языка кадр стал %d строк при экране в %d", got, s.rows)
	}
}
