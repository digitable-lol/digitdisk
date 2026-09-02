// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package report renders snapshots and walk results for a human reader.
// Absent data is printed as "—" so an empty field is never mistaken for a
// measured zero.
//
// Every line here goes through the dictionary: the wording through l.T and
// l.F, the numbers, the units and the dates through the formatters of the
// lang package.  Nothing in this package chooses a fact and nothing in it
// decides anything — it chooses the words and the spelling of the numbers,
// which is exactly the part of the output that depends on the reader.
package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"digitdisk/internal/core"
	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/lang"
	"digitdisk/internal/procfs"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
)

// dashMark is the mark of a reading nobody took.  It is the same sign in both
// languages — a dash is not a word — and it is a constant so that the four
// renderers of a video card and every column beside them cannot drift apart.
const dashMark = "—"

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func cut(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}

// Status prints the system snapshot.
func Status(w io.Writer, l lang.Lang, st sysinfo.Status) {
	p := func(line string) { fmt.Fprintln(w, line) }

	p(l.T("СИСТЕМА"))
	p(l.F("  узел          %s", dash(st.Host.Hostname)))
	p(l.F("  дистрибутив   %s", dash(st.Host.Distro)))
	if st.Host.Model != "" {
		p(l.F("  модель        %s", st.Host.Model))
	}
	if st.Host.CPUModel != "" {
		if st.Load.CPUCount > 0 {
			p(l.F("  процессор     %s × %d", st.Host.CPUModel, st.Load.CPUCount))
		} else {
			p(l.F("  процессор     %s", st.Host.CPUModel))
		}
	}
	p(l.F("  ядро          %s (%s)", dash(st.Host.KernelRelease), dash(st.Host.Machine)))
	if st.Host.UptimeSeconds > 0 {
		boot := "—"
		if st.Host.BootTime > 0 {
			boot = l.DateTime(time.Unix(st.Host.BootTime, 0))
		}
		p(l.F("  время работы  %s (с %s)", dash(l.Uptime(st.Host.UptimeSeconds)), boot))
	} else {
		p(l.T("  время работы  —"))
	}

	p("")
	p(l.T("ЗАГРУЗКА"))
	p(l.F("  средняя       %s / %s / %s  (1/5/15 мин)",
		l.Dec(st.Load.One, 2), l.Dec(st.Load.Five, 2), l.Dec(st.Load.Fifteen, 2)))
	if st.Load.CPUCount > 0 {
		p(l.F("  ядер          %d", st.Load.CPUCount))
	} else {
		p(l.T("  ядер          —"))
	}
	if st.Load.BusyPercent != nil {
		p(l.F("  занято ЦП     %s (замер %s)", l.Pct(*st.Load.BusyPercent, 1),
			l.Millis(time.Duration(st.Load.SampleMillis)*time.Millisecond)))
	} else {
		p(l.T("  занято ЦП     —"))
	}
	// The same share, processor by processor, in one line: on a machine with
	// two hundred and fifty-six of them the list itself belongs on the
	// screen, and what a printed report can say is how far apart they are.
	// One busy core among two hundred idle ones is a machine at eight per
	// cent by the line above and a machine with a problem by this one.
	if cores, ok := st.Cores(); ok {
		p(l.F("  по ядрам      мин %s / медиана %s / макс %s (ядро %d); занято больше половины %d из %d",
			l.Pct(cores.Min, 1), l.Pct(cores.Median, 1), l.Pct(cores.Max, 1), cores.Busiest, cores.Loaded, cores.Total))
	}

	p("")
	p(l.T("ПАМЯТЬ"))
	m := st.Memory
	if m.Total == 0 || !m.Has(procfs.FieldTotal) {
		p("  —")
	} else {
		p(l.F("  всего         %s", l.UBytes(m.Total)))
		if m.Has(procfs.FieldUsed) {
			p(l.F("  занято        %s (%s)  [всего − доступно]", l.UBytes(m.Used), l.Pct(pct(m.Used, m.Total), 1)))
		} else {
			p(l.T("  занято        —"))
		}
		p(l.F("  свободно      %s", measured(l, m, procfs.FieldFree, m.Free)))
		// The shared part is a /proc/meminfo key with no counterpart
		// elsewhere.  Where it was not measured the phrase is dropped
		// rather than left holding a dash inside a sentence.
		if m.Has(procfs.FieldShared) {
			p(l.F("  кэш/буферы    %s  (в т.ч. разделяемая %s)",
				measured(l, m, procfs.FieldBuffCache, m.BuffCache), l.UBytes(m.Shared)))
		} else {
			p(l.F("  кэш/буферы    %s", measured(l, m, procfs.FieldBuffCache, m.BuffCache)))
		}
		p(l.F("  доступно      %s", measured(l, m, procfs.FieldAvailable, m.Available)))
		if v, ok := m.Raw[procfs.RawWired]; ok {
			p(l.F("  закреплённая  %s", l.UBytes(v)))
		}
		if v, ok := m.Raw[procfs.RawCompressed]; ok {
			p(l.F("  сжатая        %s", l.UBytes(v)))
		}
		switch {
		case !m.Has(procfs.FieldSwapTotal):
			p(l.T("  своп          —"))
		case m.SwapTotal > 0:
			p(l.F("  своп          %s из %s занято",
				measured(l, m, procfs.FieldSwapUsed, m.SwapUsed), l.UBytes(m.SwapTotal)))
		default:
			p(l.T("  своп          нет"))
		}
	}

	p("")
	pr := st.Processes
	p(l.F("ПРОЦЕССЫ  %s", strings.Join(ProcessCounts(l, st), ", ")))
	if len(pr.TopByMemory) > 0 {
		p(l.T("  десятка по памяти:"))
		for _, x := range pr.TopByMemory {
			p(fmt.Sprintf("    %7d %-12s %10s  %s", x.PID, cut(dash(x.User), 12),
				l.Bytes(x.RSSBytes), cut(dash(firstNonEmpty(x.Cmdline, x.Comm)), 64)))
		}
	}
	if len(pr.TopByCPU) > 0 {
		p(l.T("  десятка по процессору:"))
		for _, x := range pr.TopByCPU {
			cpu := "    —"
			if x.CPUPercent != nil {
				cpu = l.Pct(*x.CPUPercent, 1)
			}
			p(fmt.Sprintf("    %7d %-12s %9s  %s", x.PID, cut(dash(x.User), 12), cpu,
				cut(dash(firstNonEmpty(x.Cmdline, x.Comm)), 64)))
		}
	}

	p("")
	p(l.T("ДИСКИ"))
	if len(st.Disks) == 0 {
		p("  —")
	} else {
		p(fmt.Sprintf("  %-28s %-22s %10s %10s %10s %6s", l.T("точка монтирования"),
			l.T("устройство"), l.T("размер"), l.T("занято"), l.T("свободно"), l.T("занят")))
		for _, d := range st.Disks {
			if d.Error != "" {
				p(l.F("  %-28s %-22s  ошибка: %s", cut(d.MountPoint, 28), cut(d.Source, 22), d.Error))
				continue
			}
			p(fmt.Sprintf("  %-28s %-22s %10s %10s %10s %6s", cut(d.MountPoint, 28), cut(d.Source, 22),
				l.UBytes(d.TotalBytes), l.UBytes(d.UsedBytes), l.UBytes(d.AvailableBytes),
				l.Pct(d.UsePercent, 1)))
		}
	}

	p("")
	p(l.T("СЕТЬ"))
	if len(st.Network) == 0 {
		p("  —")
	} else {
		p(fmt.Sprintf("  %-12s %-8s %12s %12s %10s %10s", l.T("интерфейс"), l.T("состоян."),
			l.T("принято"), l.T("передано"), l.T("пак. вх."), l.T("пак. исх.")))
		_, noCounters := st.Unmeasured(sysinfo.FactNetCounters)
		for _, n := range st.Network {
			if noCounters {
				p(fmt.Sprintf("  %-12s %-8s %12s %12s %10s %10s",
					cut(n.Name, 12), cut(dash(n.OperState), 8), "—", "—", "—", "—"))
				continue
			}
			p(fmt.Sprintf("  %-12s %-8s %12s %12s %10s %10s", cut(n.Name, 12), cut(dash(n.OperState), 8),
				l.UBytes(n.RxBytes), l.UBytes(n.TxBytes), l.UNum(n.RxPackets), l.UNum(n.TxPackets)))
		}
	}

	p("")
	p(l.T("ТЕМПЕРАТУРА"))
	if len(st.Sensors) == 0 {
		p("  —")
	} else {
		for _, s := range st.Sensors {
			extra := ""
			if s.CritC > 0 {
				extra = l.F("  (критич. %s°C)", l.Dec(s.CritC, 1))
			}
			p(fmt.Sprintf("  %-18s %-16s %6s°C%s", cut(s.Chip, 18), cut(s.Label, 16), l.Dec(s.Celsius, 1), extra))
		}
	}

	p("")
	p(l.T("ВИДЕОКАРТЫ"))
	if len(st.GPUs) == 0 {
		p("  " + dashMark)
	} else {
		p(fmt.Sprintf("  %-34s %8s %22s %9s %9s", l.T("карта"), l.T("занято"), l.T("память"), l.T("темп."), l.T("мощность")))
		for _, g := range st.GPUs {
			p(fmt.Sprintf("  %-34s %8s %22s %9s %9s", cut(dash(g.Name), 34), share(l, g.BusyPercent),
				videoMemory(l, g), celsius(l, g.Celsius), watts(l, g.Watts)))
			// Where the numbers came from, on the line under them — a
			// reading a program outside digitdisk supplied is not the
			// same claim as a reading the driver publishes in a file.
			if origin := cardOrigin(l, g); origin != "" {
				p("      " + origin)
			}
		}
	}

	// What is missing is named and left at that.  A reader who wants the
	// reasons asks for them; a reader who wants the numbers should not have
	// to read past a paragraph of them to reach the next reading.
	if names := st.UnmeasuredNames(); len(names) > 0 {
		words := make([]string, len(names))
		for i, name := range names {
			words[i] = l.Word(name)
		}
		p("")
		p(l.F("НЕ ИЗМЕРЕНО  %s", strings.Join(words, ", ")))
		p(l.T("             почему — digitdisk status --why"))
	}
}

// Why prints every absence in the snapshot with its reason.  This is the whole
// of `digitdisk status --why`, and the only place the reasons are printed.
func Why(w io.Writer, l lang.Lang, st sysinfo.Status) {
	p := func(line string) { fmt.Fprintln(w, line) }
	all := st.UnmeasuredAll()
	if len(all) == 0 {
		p(l.T("НЕ ИЗМЕРЕНО  ничего: снимок полон"))
		return
	}
	p(l.F("НЕ ИЗМЕРЕНО  %d", len(all)))
	for _, pair := range all {
		p("")
		p("  " + l.Word(pair.Name))
		p("      " + dash(pair.Why.In(l)))
	}
}

// ProcessCounts builds the process line as a sentence of the counts that were
// actually taken.  A count nobody could take is left out of the sentence and
// named at the end of the report instead: a dash in the middle of a phrase
// reads as a number, and this is not a table where a column must be held open.
//
// It is exported because the live screen prints the same sentence, and the two
// must not drift apart: a number that is honest on paper and invented on the
// screen is worse than either.
func ProcessCounts(l lang.Lang, st sysinfo.Status) []string {
	pr := st.Processes
	out := []string{l.F("всего %d", pr.Total)}
	partial := false
	if _, unmeasured := st.Unmeasured(sysinfo.FactThreads); !unmeasured && pr.WithDetail > 0 {
		out = append(out, l.F("потоков %d", pr.Threads))
		partial = partial || pr.WithDetail < pr.Total
	}
	if _, unmeasured := st.Unmeasured(sysinfo.FactRunning); !unmeasured {
		out = append(out, l.F("выполняется %d", pr.Running))
	}
	if _, unmeasured := st.Unmeasured(sysinfo.FactBlocked); !unmeasured {
		out = append(out, l.F("заблокировано %d", pr.Blocked))
	}
	if pr.Unreadable > 0 {
		out = append(out, l.F("не прочитано %d", pr.Unreadable))
	}
	// The coverage is its own phrase rather than a bracket on one number:
	// it qualifies every count that came from reading the processes one by
	// one, and a bracket would look like it qualified only the last of them.
	if partial {
		out = append(out, l.F("замерено по %d процессам", pr.WithDetail))
	}
	return out
}

// Analyze prints the result of a tree walk.
func Analyze(w io.Writer, l lang.Lang, r scan.Result) {
	p := func(line string) { fmt.Fprintln(w, line) }

	p(l.F("ОБХОД  %s", r.Root))
	p(l.F("  записей       %s  (файлов %s, каталогов %s, ссылок %s, прочего %s)",
		l.Num(int64(r.Entries)), l.Num(int64(r.Files)), l.Num(int64(r.Dirs)),
		l.Num(int64(r.Links)), l.Num(int64(r.Others))))
	p(l.F("  объём         %s (%s, видимый размер: du --apparent-size, он же du -sb)",
		l.Bytes(r.TotalBytes), l.RawBytes(r.TotalBytes)))
	p(l.F("                файлы %s, ссылки %s; сверх того сами каталоги %s (в объём не входят, как и у du)",
		l.Bytes(r.FileBytes), l.Bytes(r.LinkBytes), l.Bytes(r.DirBytes)))
	if r.HardlinkDupes > 0 {
		p(l.F("  жёсткие ссылки %s повторных имён не засчитано (%s)",
			l.Num(int64(r.HardlinkDupes)), l.Bytes(r.HardlinkBytes)))
	}
	s := r.Skipped
	p(l.F("  пропущено     %s  (нет доступа %s, исчезло %s, иные ошибки %s, граница ФС %s, предел глубины %s)",
		l.Num(int64(s.Total())), l.Num(int64(s.PermissionDenied)), l.Num(int64(s.Vanished)),
		l.Num(int64(s.OtherErrors)), l.Num(int64(s.DeviceBoundaries)), l.Num(int64(s.DepthLimited))))
	p(l.F("  время         %s с", l.Dec(r.DurationSeconds, 2)))

	p("")
	p(l.F("РЕШАЮЩИЙ СЛОЙ  %s, договор версии %d", l.Word(r.Decider), r.ContractVersion))
	if !r.DeciderReady {
		p(l.T("  ВНИМАНИЕ: настоящий разбор не выполнялся — все записи возвращены как"))
		p(l.F("  %s/%s. Разряды ниже показывают работу стыковки, не анализ.",
			l.Word(string(core.ClassUnknown)), l.Word(string(core.VerdictKeep))))
	}

	p("")
	p(l.T("ПО РАЗРЯДАМ"))
	for _, c := range core.Classes {
		b := r.ByClass[c]
		p(fmt.Sprintf("  %-14s %8s %12s", l.Word(string(c)), l.Num(int64(b.Count)), l.Bytes(b.Bytes)))
	}
	p("")
	p(l.T("ПО ПРИГОВОРАМ"))
	for _, v := range core.Verdicts {
		b := r.ByVerdict[v]
		p(fmt.Sprintf("  %-14s %8s %12s", l.Word(string(v)), l.Num(int64(b.Count)), l.Bytes(b.Bytes)))
	}

	p("")
	rem := r.ByVerdict[core.VerdictRemovable]
	p(l.F("ПРЕДЛОЖЕНО УБРАТЬ  %s записей, %s", l.Num(int64(rem.Count)), l.Bytes(rem.Bytes)))
	p(l.T("  (analyze только считает; убирает `digitdisk clean`, и тоже не сразу)"))
	if len(r.Removable) == 0 {
		p(l.T("  — нечего"))
	} else {
		for _, e := range r.Removable {
			p(fmt.Sprintf("  %10s  %-12s %8s  %s", l.Bytes(e.Size), l.Word(string(e.Class)),
				l.Days(e.AgeDays), cut(e.Path, 80)))
		}
	}

	if len(r.Largest) > 0 {
		p("")
		p(l.T("САМОЕ КРУПНОЕ"))
		for _, e := range r.Largest {
			p(fmt.Sprintf("  %10s  %-9s %8s  %s", l.Bytes(e.Size), l.Word(string(e.Kind)),
				l.Days(e.AgeDays), cut(e.Path, 80)))
		}
	}
}

// measured renders a memory field: its value when the system published it,
// "—" when it did not.  A field nobody measured must not print as "0 Б", which
// is a number a reader would believe.
func measured(l lang.Lang, m procfs.Memory, field string, v uint64) string {
	if !m.Has(field) {
		return "—"
	}
	return l.UBytes(v)
}

func pct(a, b uint64) float64 {
	if b == 0 {
		return 0
	}
	return 100 * float64(a) / float64(b)
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

// The four readings of a video card, each rendered as itself or as a dash.
// A driver that publishes no load share is not a card doing nothing, and the
// difference has to survive into the printed line.

func share(l lang.Lang, v *float64) string {
	if v == nil {
		return dashMark
	}
	return l.Pct(*v, 1)
}

func celsius(l lang.Lang, v *float64) string {
	if v == nil {
		return dashMark
	}
	return l.Dec(*v, 1) + "°C"
}

func watts(l lang.Lang, v *float64) string {
	if v == nil {
		return dashMark
	}
	return l.F("%s Вт", l.Dec(*v, 1))
}

// videoMemory is the card's memory as one column: what is taken out of what
// there is, and a dash for either half nobody published.
func videoMemory(l lang.Lang, g gpuinfo.Card) string {
	switch {
	case g.MemoryTotalBytes != nil && g.MemoryUsedBytes != nil:
		return l.F("%s из %s", l.UBytes(*g.MemoryUsedBytes), l.UBytes(*g.MemoryTotalBytes))
	case g.MemoryTotalBytes != nil:
		return l.F("%s из %s", dashMark, l.UBytes(*g.MemoryTotalBytes))
	default:
		return dashMark
	}
}

// cardOrigin says where a card's numbers came from: the bus it sits on, the
// driver that answered, and — when it was not a file that answered — the name
// of the program that was run.
func cardOrigin(l lang.Lang, g gpuinfo.Card) string {
	var parts []string
	if g.Bus != "" {
		parts = append(parts, l.F("шина %s", g.Bus))
	}
	if g.Driver != "" {
		parts = append(parts, l.F("драйвер %s", g.Driver))
	}
	switch {
	case g.Outside:
		parts = append(parts, l.F("числа от чужой программы %s", g.Source))
	case g.Source != "":
		parts = append(parts, l.F("числа из %s", g.Source))
	}
	return strings.Join(parts, "  ·  ")
}
