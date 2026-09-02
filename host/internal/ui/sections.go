// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"digitdisk/internal/cli"
	"digitdisk/internal/lang"
	"digitdisk/internal/report"
	"digitdisk/internal/sysinfo"
)

// A section is one page of the screen.  The sections are the sections of the
// printed report, in the printed order, with an opening page that gathers the
// gauges; no reading appears here that `digitdisk status` does not print.
type section struct {
	// title is the name in the strip at the top, asked for in the language
	// the screen is being drawn in.
	//
	// It is a call and not a string on purpose.  A field holding «ОБЗОР»
	// would be Russian sitting outside the dictionary — the reader of the
	// source could not tell that this name is ever translated, and the
	// check that walks the source would be right to call it a leak.  Written
	// as l.T("ОБЗОР") the wording stands where it is looked up, and the name
	// is translated at the moment it is shown rather than at start-up, so
	// the key that switches the language switches the strip with it.
	title  func(lang.Lang) string
	render func(*screen) []string
}

var sections = []section{
	{func(l lang.Lang) string { return l.T("ОБЗОР") }, (*screen).overview},
	{func(l lang.Lang) string { return l.T("СИСТЕМА") }, (*screen).system},
	{func(l lang.Lang) string { return l.T("ЗАГРУЗКА") }, (*screen).load},
	{func(l lang.Lang) string { return l.T("ПАМЯТЬ") }, (*screen).memory},
	{func(l lang.Lang) string { return l.T("ПРОЦЕССЫ") }, (*screen).processes},
	{func(l lang.Lang) string { return l.T("ДИСКИ") }, (*screen).disks},
	{func(l lang.Lang) string { return l.T("СЕТЬ") }, (*screen).network},
	{func(l lang.Lang) string { return l.T("ТЕМПЕРАТУРА") }, (*screen).sensors},
	{func(l lang.Lang) string { return l.T("ВИДЕОКАРТЫ") }, (*screen).gpus},
	{func(l lang.Lang) string { return l.T("НЕ ПРОЧИТАНО") }, (*screen).missing},
}

// dash is the mark of a reading the system did not publish.  It is the mark
// the printed report uses, and it is never a zero.
const dash = "—"

func (s *screen) kv(label, value string) string {
	t := s.t
	var r row
	r.add("  "+fit(label, 16)+"  ", func(x string) string { return t.Fg(t.P.Muted, x) })
	r.add(fit(value, maxInt(1, s.cols-r.w-1)), func(x string) string {
		return t.Fg(t.P.Foreground, strings.TrimRight(x, " "))
	})
	return r.String()
}

func (s *screen) caption(text string) string {
	return s.t.Bold(s.t.P.AccentSoft, "  "+text)
}

func (s *screen) note(text string) string {
	return s.t.Fg(s.t.P.Subtle, "  "+fit(text, maxInt(1, s.cols-3)))
}

func (s *screen) waiting() []string {
	return []string{"", s.note(s.l.T("замер идёт…"))}
}

func text(v string) string {
	if strings.TrimSpace(v) == "" {
		return dash
	}
	return v
}

// barWidth is the width of a gauge, kept between a useful minimum and a width
// that does not swallow the screen.
// reading picks how much of a measurement is spelled out: the whole sentence
// where there is room, the share alone where there is not.  The number does
// not change with the width — only how much of its wording survives.
func (s *screen) reading(full, short string) string {
	if s.cols < 80 {
		return short
	}
	return full
}

func (s *screen) barWidth() int {
	w := s.cols/3 - 6
	if w < 10 {
		w = 10
	}
	if w > 40 {
		w = 40
	}
	return w
}

// overview gathers the gauges of the whole snapshot on one page: what a glance
// at the machine is meant to answer.
func (s *screen) overview() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t, st := s.t, s.st
	bw := s.barWidth()
	out := []string{""}

	if st.Load.BusyPercent != nil {
		out = append(out, t.gauge(s.l.T("ЦП занято"), 12, *st.Load.BusyPercent/100, s.l.Pct(*st.Load.BusyPercent, 1), bw))
	} else {
		out = append(out, t.gaugeUnmeasured(s.l.T("ЦП занято"), 12, s.l.T("замер не делался"), bw))
	}
	out = append(out, s.sparkLine(12, s.cpuHist))
	// The same share spread across the processors, on one line: a machine
	// with one core on fire and two hundred idle reads as eight per cent
	// busy, and this is where that becomes visible.
	if len(st.Load.Cores) > 0 {
		out = append(out, s.coreComb(12))
	}

	if st.Memory.Total > 0 {
		f := pctOf(st.Memory.Used, st.Memory.Total)
		out = append(out, t.gauge(s.l.T("Память"), 12, f, s.reading(
			s.l.F("%s из %s  (%s)", s.l.UBytes(st.Memory.Used), s.l.UBytes(st.Memory.Total), s.l.Pct(f*100, 1)),
			s.l.Pct(f*100, 1)), bw))
		out = append(out, s.sparkLine(12, s.memHist))
	} else {
		out = append(out, t.gaugeUnmeasured(s.l.T("Память"), 12, dash, bw))
	}

	if st.Memory.SwapTotal > 0 {
		f := pctOf(st.Memory.SwapUsed, st.Memory.SwapTotal)
		out = append(out, t.gauge(s.l.T("Своп"), 12, f, s.l.F("%s из %s",
			s.l.UBytes(st.Memory.SwapUsed), s.l.UBytes(st.Memory.SwapTotal)), bw))
	} else {
		out = append(out, t.gaugeUnmeasured(s.l.T("Своп"), 12, s.l.T("нет"), bw))
	}

	if len(st.GPUs) > 0 {
		out = append(out, "", s.caption("ВИДЕОКАРТЫ"))
		out = append(out, s.gpuGauges(3)...)
	}

	out = append(out, "")
	out = append(out, s.kv(s.l.T("средняя"), s.l.F("%s / %s / %s   (1/5/15 мин, ядер %s)",
		s.l.Dec(st.Load.One, 2), s.l.Dec(st.Load.Five, 2), s.l.Dec(st.Load.Fifteen, 2), count(st.Load.CPUCount))))
	out = append(out, s.kv(s.l.T("процессы"), strings.Join(report.ProcessCounts(s.l, st), ", ")))
	out = append(out, s.kv(s.l.T("время работы"), text(s.l.Uptime(st.Host.UptimeSeconds))))

	if len(st.Disks) > 0 {
		out = append(out, "", s.caption(s.l.T("ДИСКИ")))
		shown := st.Disks
		if len(shown) > 6 {
			shown = shown[:6]
		}
		for _, d := range shown {
			if d.Error != "" {
				out = append(out, t.gaugeUnmeasured(d.MountPoint, 20, s.l.F("ошибка: %s", d.Error), bw))
				continue
			}
			out = append(out, t.gauge(d.MountPoint, 20, d.UsePercent/100, s.reading(
				s.l.F("%s свободно из %s", s.l.UBytes(d.AvailableBytes), s.l.UBytes(d.TotalBytes)),
				s.l.F("%s своб.", s.l.UBytes(d.AvailableBytes))), bw))
		}
		if len(st.Disks) > len(shown) {
			out = append(out, s.note(s.l.F("…и ещё %d — раздел ДИСКИ", len(st.Disks)-len(shown))))
		}
	}

	if len(st.Missing) > 0 {
		out = append(out, "", s.note(s.l.F("не прочитано источников: %d — раздел НЕ ПРОЧИТАНО", len(st.Missing))))
	}
	return out
}

// sparkLine draws the history of one share directly under its gauge, in the
// same columns, so the bar and its past line up.
func (s *screen) sparkLine(nameWidth int, history []float64) string {
	t := s.t
	var r row
	r.add("  "+fit(s.l.T("история"), nameWidth)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	w := s.barWidth()
	r.plain(t.spark(history, w))
	r.w += w
	return r.String()
}

func count(n int) string {
	if n <= 0 {
		return dash
	}
	return fmt.Sprint(n)
}

// СИСТЕМА is drawn in hardware.go: the mark of the system and the facts a
// person recognises a machine by.

// takenAt is the moment of the snapshot, written the way this language writes
// a moment.  The stamp itself is ISO-8601 and stays that way in the JSON; a
// stamp this tool cannot read is shown as it came rather than guessed at.
func (s *screen) takenAt() string {
	if t, err := time.Parse(time.RFC3339, s.st.TakenAt); err == nil {
		return s.l.DateTime(t)
	}
	return text(s.st.TakenAt)
}

func (s *screen) load() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t, ld := s.t, s.st.Load
	bw := s.barWidth()
	out := []string{""}
	out = append(out,
		s.kv(s.l.T("средняя"), s.l.F("%s / %s / %s  (1/5/15 мин)",
			s.l.Dec(ld.One, 2), s.l.Dec(ld.Five, 2), s.l.Dec(ld.Fifteen, 2))),
		s.kv(s.l.T("ядер"), count(ld.CPUCount)),
	)
	if ld.TotalEntities > 0 {
		out = append(out, s.kv(s.l.T("в очереди"), s.l.F("%d из %d", ld.Runnable, ld.TotalEntities)))
	}
	out = append(out, "")
	if ld.BusyPercent != nil {
		f := *ld.BusyPercent / 100
		out = append(out, t.gauge(s.l.T("занято ЦП"), 12, f,
			s.l.F("%s (замер %d мс)", s.l.Pct(*ld.BusyPercent, 1), ld.SampleMillis), bw))
	} else {
		out = append(out, t.gaugeUnmeasured(s.l.T("занято ЦП"), 12, s.l.T("замер не делался"), bw))
	}
	out = append(out, s.sparkLine(12, s.cpuHist))
	// The same window, processor by processor.  It lives on this page and
	// not on one of its own: it is the same measurement, and the section
	// strip is the printed report's, where the per-core figures are a line
	// of ЗАГРУЗКА — see coresBlock in hardware.go.
	out = append(out, s.coresBlock()...)
	return out
}

func (s *screen) memory() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t, m := s.t, s.st.Memory
	bw := s.barWidth()
	if m.Total == 0 {
		return []string{"", s.note(dash)}
	}
	f := pctOf(m.Used, m.Total)
	out := []string{""}
	out = append(out,
		t.gauge(s.l.T("занято"), 12, f, s.reading(
			s.l.F("%s из %s  (%s)", s.l.UBytes(m.Used), s.l.UBytes(m.Total), s.l.Pct(f*100, 1)),
			s.l.Pct(f*100, 1)), bw),
		s.note(s.l.T("занято = всего − доступно, как в free(1)")),
		"",
		s.kv(s.l.T("всего"), s.l.UBytes(m.Total)),
		s.kv(s.l.T("свободно"), s.l.UBytes(m.Free)),
		s.kv(s.l.T("доступно"), s.l.UBytes(m.Available)),
		s.kv(s.l.T("кэш/буферы"), s.l.F("%s  (в т.ч. разделяемая %s)", s.l.UBytes(m.BuffCache), s.l.UBytes(m.Shared))),
	)
	out = append(out, "")
	if m.SwapTotal > 0 {
		sf := pctOf(m.SwapUsed, m.SwapTotal)
		out = append(out, t.gauge(s.l.T("своп"), 12, sf, s.l.F("%s из %s занято",
			s.l.UBytes(m.SwapUsed), s.l.UBytes(m.SwapTotal)), bw))
	} else {
		out = append(out, t.gaugeUnmeasured(s.l.T("своп"), 12, s.l.T("нет"), bw))
	}
	out = append(out, "", s.sparkLine(12, s.memHist))
	return out
}

func (s *screen) processes() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t, pr := s.t, s.st.Processes
	out := []string{"", s.kv(s.l.T("процессы"), strings.Join(report.ProcessCounts(s.l, s.st), ", "))}

	cmdWidth := maxInt(10, s.cols-38)
	head := func(what string) string {
		var r row
		r.add("    "+right("PID", 7)+"  "+fit(s.l.T("владелец"), 12)+" "+right(what, 9)+"  "+fit(s.l.T("команда"), cmdWidth),
			func(x string) string { return t.Fg(t.P.Border, x) })
		return r.String()
	}
	line := func(p sysinfo.Proc, val string) string {
		var r row
		r.add("    "+right(fmt.Sprint(p.PID), 7)+"  ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		r.add(fit(text(p.User), 12)+" ", func(x string) string { return t.Fg(t.P.Purple, x) })
		r.add(right(val, 9)+"  ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		r.add(fit(text(firstOf(p.Cmdline, p.Comm)), cmdWidth), func(x string) string {
			return t.Fg(t.P.Muted, strings.TrimRight(x, " "))
		})
		return r.String()
	}

	if len(pr.TopByMemory) > 0 {
		out = append(out, "", s.caption(s.l.T("ПО ПАМЯТИ")), head(s.l.T("память")))
		for _, p := range pr.TopByMemory {
			out = append(out, line(p, s.l.Bytes(p.RSSBytes)))
		}
	}
	if len(pr.TopByCPU) > 0 {
		out = append(out, "", s.caption(s.l.T("ПО ПРОЦЕССОРУ")), head(s.l.T("ЦП")))
		for _, p := range pr.TopByCPU {
			v := dash
			if p.CPUPercent != nil {
				v = s.l.Pct(*p.CPUPercent, 1)
			}
			out = append(out, line(p, v))
		}
	}
	return out
}

func (s *screen) disks() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t := s.t
	if len(s.st.Disks) == 0 {
		return []string{"", s.note(dash)}
	}
	mount := maxInt(10, minInt(30, s.cols-52))
	bw := maxInt(8, minInt(20, s.cols-mount-42))
	out := []string{""}
	var h row
	h.add("  "+fit(s.l.T("точка монтирования"), mount)+" "+right(s.l.T("размер"), 10)+" "+right(s.l.T("занято"), 10)+" "+
		right(s.l.T("свободно"), 10)+"  "+fit("", bw)+" "+right(s.l.T("занят"), 6),
		func(x string) string { return t.Fg(t.P.Border, x) })
	out = append(out, h.String())
	for _, d := range s.st.Disks {
		var r row
		r.add("  "+fit(d.MountPoint, mount)+" ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		if d.Error != "" {
			r.add(s.l.F("ошибка: %s", d.Error), func(x string) string { return t.Fg(t.P.Red, x) })
			out = append(out, r.String())
			continue
		}
		f := d.UsePercent / 100
		r.add(right(s.l.UBytes(d.TotalBytes), 10)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(s.l.UBytes(d.UsedBytes), 10)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(s.l.UBytes(d.AvailableBytes), 10)+"  ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		r.plain(t.bar(f, bw))
		r.w += bw
		r.add(" "+right(s.l.Pct(d.UsePercent, 1), 6), func(x string) string { return t.Fg(t.level(f), x) })
		out = append(out, r.String())
		out = append(out, t.Fg(t.P.Subtle, "  "+fit("  "+text(d.Source)+"  ·  "+text(d.FSType)+
			readOnly(s.l, d.ReadOnly)+inodes(s.l, d), maxInt(1, s.cols-3))))
	}
	return out
}

func readOnly(l lang.Lang, ro bool) string {
	if ro {
		return l.T("  ·  только чтение")
	}
	return ""
}

func inodes(l lang.Lang, d sysinfo.Disk) string {
	if d.InodesTotal == 0 {
		return ""
	}
	return l.F("  ·  inode %s свободно из %s", l.UNum(d.InodesFree), l.UNum(d.InodesTotal))
}

func (s *screen) network() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t := s.t
	if len(s.st.Network) == 0 {
		return []string{"", s.note(dash)}
	}
	out := []string{""}
	var h row
	h.add("  "+fit(s.l.T("интерфейс"), 14)+" "+fit(s.l.T("состоян."), 9)+" "+right(s.l.T("принято"), 12)+" "+
		right(s.l.T("передано"), 12)+" "+right(s.l.T("пак. вх."), 13)+" "+right(s.l.T("пак. исх."), 13),
		func(x string) string { return t.Fg(t.P.Border, x) })
	out = append(out, h.String())
	for _, n := range s.st.Network {
		var r row
		r.add("  "+fit(n.Name, 14)+" ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		r.add(fit(text(n.OperState), 9)+" ", func(x string) string { return t.Fg(stateColour(t, n.OperState), x) })
		r.add(right(s.l.UBytes(n.RxBytes), 12)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(s.l.UBytes(n.TxBytes), 12)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(s.l.UNum(n.RxPackets), 13)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		r.add(right(s.l.UNum(n.TxPackets), 13), func(x string) string { return t.Fg(t.P.Subtle, x) })
		out = append(out, r.String())
		if extra := ifaceExtra(s.l, n); extra != "" {
			out = append(out, t.Fg(t.P.Subtle, "  "+fit("  "+extra, maxInt(1, s.cols-3))))
		}
	}
	return out
}

func stateColour(t Theme, state string) slot {
	switch strings.ToLower(state) {
	case "up":
		return t.P.Green
	case "down":
		return t.P.Red
	case "":
		return t.P.Subtle
	}
	return t.P.Yellow
}

func ifaceExtra(l lang.Lang, n sysinfo.Iface) string {
	var parts []string
	if len(n.Addresses) > 0 {
		parts = append(parts, strings.Join(n.Addresses, ", "))
	}
	if n.MAC != "" {
		parts = append(parts, n.MAC)
	}
	if n.MTU > 0 {
		parts = append(parts, fmt.Sprintf("MTU %d", n.MTU))
	}
	if n.SpeedMbit > 0 {
		parts = append(parts, l.F("%d Мбит/с", n.SpeedMbit))
	}
	if n.RxErrors+n.TxErrors+n.RxDropped+n.TxDropped > 0 {
		parts = append(parts, l.F("ошибок %d/%d, потеряно %d/%d",
			n.RxErrors, n.TxErrors, n.RxDropped, n.TxDropped))
	}
	return strings.Join(parts, "  ·  ")
}

func (s *screen) sensors() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t := s.t
	if len(s.st.Sensors) == 0 {
		return []string{"", s.note(s.l.T("— (датчиков не нашлось)"))}
	}
	bw := s.barWidth()
	out := []string{""}
	for _, sn := range s.st.Sensors {
		name := text(sn.Chip)
		if sn.Label != "" {
			name += " · " + sn.Label
		}
		// The gauge is filled against the chip's own critical point when it
		// published one, and against 100 °C when it did not.  Which of the
		// two was used is said on the line, so the bar is never a guess
		// dressed as a measurement.
		limit, against := 100.0, s.l.T("из 100 °C")
		if sn.CritC > 0 {
			limit, against = sn.CritC, s.l.F("критич. %s °C", s.l.Dec(sn.CritC, 1))
		}
		out = append(out, t.gauge(name, 26, sn.Celsius/limit,
			fmt.Sprintf("%s °C  %s", s.l.Dec(sn.Celsius, 1), against), bw))
	}
	return out
}

func (s *screen) missing() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t := s.t
	if len(s.st.Missing) == 0 {
		return []string{"", s.note(s.l.T("прочитано всё, чего ждали"))}
	}
	// Go walks a map in a different order every time.  Sorting keeps the
	// section still between redraws instead of reshuffling once a second.
	keys := make([]string, 0, len(s.st.Missing))
	for k := range s.st.Missing {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := []string{"", s.note(s.l.T("источник назван вместе с причиной; нулём его отсутствие не притворяется")), ""}
	for _, k := range keys {
		var r row
		r.add("  "+fit(s.l.Word(k), 24)+"  ", func(x string) string { return t.Fg(t.P.Yellow, x) })
		r.add(fit(s.st.Missing[k].In(s.l), maxInt(1, s.cols-r.w-1)), func(x string) string {
			return t.Fg(t.P.Muted, strings.TrimRight(x, " "))
		})
		out = append(out, r.String())
	}
	return out
}

func firstOf(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// commandsPage is the list of subcommands, shown over the body by «?».
//
// It comes from internal/cli, the same list the справка and digitdisk.1 are
// built from, so the screen cannot come to name a command the tool does not
// have or miss one it does.
//
// It SHOWS and does not RUN, on purpose.  This screen is `status`, which reads
// and writes nothing; a choice here that ran clean would put a command that
// moves files one keystroke away from a read-only view, and past the --apply
// and the --confirm that stand in front of removal everywhere else.  Two of
// the commands need a path the keyboard has not been given anyway.
func (s *screen) commandsPage() []string {
	t := s.t
	out := []string{"", s.caption(s.l.T("КОМАНДЫ")), ""}
	for _, c := range cli.Commands {
		var r row
		r.add("  "+fit(c.Call(s.l), 20), func(x string) string { return t.Bold(t.P.AccentSoft, x) })
		r.add(fit(s.l.T(c.Gloss), maxInt(1, s.cols-r.w-1)), func(x string) string {
			return t.Fg(t.P.Foreground, strings.TrimRight(x, " "))
		})
		out = append(out, r.String())
	}
	out = append(out, "",
		s.note(s.l.T("l — язык вывода (ru ⇄ en), выбор запоминается")),
		s.note(s.l.T("Экран ничего не запускает: команды набираются в оболочке.")),
		s.note(s.l.T("Ключи: digitdisk --help.  Подробно: man digitdisk")))
	return out
}
