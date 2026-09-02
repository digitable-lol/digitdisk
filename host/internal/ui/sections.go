// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"digitdisk/internal/report"
	"digitdisk/internal/sysinfo"
)

// A section is one page of the screen.  The sections are the sections of the
// printed report, in the printed order, with an opening page that gathers the
// gauges; no reading appears here that `digitdisk status` does not print.
type section struct {
	title  string
	render func(*screen) []string
}

var sections = []section{
	{"ОБЗОР", (*screen).overview},
	{"СИСТЕМА", (*screen).system},
	{"ЗАГРУЗКА", (*screen).load},
	{"ПАМЯТЬ", (*screen).memory},
	{"ПРОЦЕССЫ", (*screen).processes},
	{"ДИСКИ", (*screen).disks},
	{"СЕТЬ", (*screen).network},
	{"ТЕМПЕРАТУРА", (*screen).sensors},
	{"НЕ ПРОЧИТАНО", (*screen).missing},
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
	return []string{"", s.note("замер идёт…")}
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
		out = append(out, t.gauge("ЦП занято", 12, *st.Load.BusyPercent/100, percent(*st.Load.BusyPercent/100), bw))
	} else {
		out = append(out, t.gaugeUnmeasured("ЦП занято", 12, "замер не делался", bw))
	}
	out = append(out, s.sparkLine(12, s.cpuHist))

	if st.Memory.Total > 0 {
		f := pctOf(st.Memory.Used, st.Memory.Total)
		out = append(out, t.gauge("Память", 12, f, s.reading(
			fmt.Sprintf("%s из %s  (%s)", report.UBytes(st.Memory.Used), report.UBytes(st.Memory.Total), percent(f)),
			percent(f)), bw))
		out = append(out, s.sparkLine(12, s.memHist))
	} else {
		out = append(out, t.gaugeUnmeasured("Память", 12, dash, bw))
	}

	if st.Memory.SwapTotal > 0 {
		f := pctOf(st.Memory.SwapUsed, st.Memory.SwapTotal)
		out = append(out, t.gauge("Своп", 12, f, fmt.Sprintf("%s из %s",
			report.UBytes(st.Memory.SwapUsed), report.UBytes(st.Memory.SwapTotal)), bw))
	} else {
		out = append(out, t.gaugeUnmeasured("Своп", 12, "нет", bw))
	}

	out = append(out, "")
	out = append(out, s.kv("средняя", fmt.Sprintf("%.2f / %.2f / %.2f   (1/5/15 мин, ядер %s)",
		st.Load.One, st.Load.Five, st.Load.Fifteen, count(st.Load.CPUCount))))
	out = append(out, s.kv("процессы", strings.Join(report.ProcessCounts(st), ", ")))
	out = append(out, s.kv("время работы", text(st.Host.UptimeHuman)))

	if len(st.Disks) > 0 {
		out = append(out, "", s.caption("ДИСКИ"))
		shown := st.Disks
		if len(shown) > 6 {
			shown = shown[:6]
		}
		for _, d := range shown {
			if d.Error != "" {
				out = append(out, t.gaugeUnmeasured(d.MountPoint, 20, "ошибка: "+d.Error, bw))
				continue
			}
			out = append(out, t.gauge(d.MountPoint, 20, d.UsePercent/100, s.reading(
				fmt.Sprintf("%s свободно из %s", report.UBytes(d.AvailableBytes), report.UBytes(d.TotalBytes)),
				report.UBytes(d.AvailableBytes)+" своб."), bw))
		}
		if len(st.Disks) > len(shown) {
			out = append(out, s.note(fmt.Sprintf("…и ещё %d — раздел ДИСКИ", len(st.Disks)-len(shown))))
		}
	}

	if len(st.Missing) > 0 {
		out = append(out, "", s.note(fmt.Sprintf("не прочитано источников: %d — раздел НЕ ПРОЧИТАНО", len(st.Missing))))
	}
	return out
}

// sparkLine draws the history of one share directly under its gauge, in the
// same columns, so the bar and its past line up.
func (s *screen) sparkLine(nameWidth int, history []float64) string {
	t := s.t
	var r row
	r.add("  "+fit("история", nameWidth)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
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

func (s *screen) system() []string {
	if !s.haveSt {
		return s.waiting()
	}
	h := s.st.Host
	out := []string{""}
	out = append(out,
		s.kv("узел", text(h.Hostname)),
		s.kv("дистрибутив", text(h.Distro)),
		s.kv("ядро", fmt.Sprintf("%s (%s)", text(h.KernelRelease), text(h.Machine))),
	)
	if strings.TrimSpace(h.KernelVersion) != "" {
		out = append(out, s.kv("сборка ядра", h.KernelVersion))
	}
	if h.UptimeSeconds > 0 {
		boot := dash
		if h.BootTime > 0 {
			boot = time.Unix(h.BootTime, 0).Format("2006-01-02 15:04")
		}
		out = append(out, s.kv("время работы", fmt.Sprintf("%s (с %s)", text(h.UptimeHuman), boot)))
	} else {
		out = append(out, s.kv("время работы", dash))
	}
	out = append(out, s.kv("снимок взят", text(s.st.TakenAt)))
	return out
}

func (s *screen) load() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t, l := s.t, s.st.Load
	bw := s.barWidth()
	out := []string{""}
	out = append(out,
		s.kv("средняя", fmt.Sprintf("%.2f / %.2f / %.2f  (1/5/15 мин)", l.One, l.Five, l.Fifteen)),
		s.kv("ядер", count(l.CPUCount)),
	)
	if l.TotalEntities > 0 {
		out = append(out, s.kv("в очереди", fmt.Sprintf("%d из %d", l.Runnable, l.TotalEntities)))
	}
	out = append(out, "")
	if l.BusyPercent != nil {
		f := *l.BusyPercent / 100
		out = append(out, t.gauge("занято ЦП", 12, f, fmt.Sprintf("%s (замер %d мс)", percent(f), l.SampleMillis), bw))
	} else {
		out = append(out, t.gaugeUnmeasured("занято ЦП", 12, "замер не делался", bw))
	}
	out = append(out, s.sparkLine(12, s.cpuHist))
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
		t.gauge("занято", 12, f, s.reading(
			fmt.Sprintf("%s из %s  (%s)", report.UBytes(m.Used), report.UBytes(m.Total), percent(f)),
			percent(f)), bw),
		s.note("занято = всего − доступно, как в free(1)"),
		"",
		s.kv("всего", report.UBytes(m.Total)),
		s.kv("свободно", report.UBytes(m.Free)),
		s.kv("доступно", report.UBytes(m.Available)),
		s.kv("кэш/буферы", fmt.Sprintf("%s  (в т.ч. разделяемая %s)", report.UBytes(m.BuffCache), report.UBytes(m.Shared))),
	)
	out = append(out, "")
	if m.SwapTotal > 0 {
		sf := pctOf(m.SwapUsed, m.SwapTotal)
		out = append(out, t.gauge("своп", 12, sf, fmt.Sprintf("%s из %s занято",
			report.UBytes(m.SwapUsed), report.UBytes(m.SwapTotal)), bw))
	} else {
		out = append(out, t.gaugeUnmeasured("своп", 12, "нет", bw))
	}
	out = append(out, "", s.sparkLine(12, s.memHist))
	return out
}

func (s *screen) processes() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t, pr := s.t, s.st.Processes
	out := []string{"", s.kv("процессы", strings.Join(report.ProcessCounts(s.st), ", "))}

	cmdWidth := maxInt(10, s.cols-38)
	head := func(what string) string {
		var r row
		r.add("    "+right("PID", 7)+"  "+fit("владелец", 12)+" "+right(what, 9)+"  "+fit("команда", cmdWidth),
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
		out = append(out, "", s.caption("ПО ПАМЯТИ"), head("память"))
		for _, p := range pr.TopByMemory {
			out = append(out, line(p, report.Bytes(p.RSSBytes)))
		}
	}
	if len(pr.TopByCPU) > 0 {
		out = append(out, "", s.caption("ПО ПРОЦЕССОРУ"), head("ЦП"))
		for _, p := range pr.TopByCPU {
			v := dash
			if p.CPUPercent != nil {
				v = fmt.Sprintf("%.1f%%", *p.CPUPercent)
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
	h.add("  "+fit("точка монтирования", mount)+" "+right("размер", 10)+" "+right("занято", 10)+" "+
		right("свободно", 10)+"  "+fit("", bw)+" "+right("занят", 6),
		func(x string) string { return t.Fg(t.P.Border, x) })
	out = append(out, h.String())
	for _, d := range s.st.Disks {
		var r row
		r.add("  "+fit(d.MountPoint, mount)+" ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		if d.Error != "" {
			r.add("ошибка: "+d.Error, func(x string) string { return t.Fg(t.P.Red, x) })
			out = append(out, r.String())
			continue
		}
		f := d.UsePercent / 100
		r.add(right(report.UBytes(d.TotalBytes), 10)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(report.UBytes(d.UsedBytes), 10)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(report.UBytes(d.AvailableBytes), 10)+"  ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		r.plain(t.bar(f, bw))
		r.w += bw
		r.add(" "+right(fmt.Sprintf("%.1f%%", d.UsePercent), 6), func(x string) string { return t.Fg(t.level(f), x) })
		out = append(out, r.String())
		out = append(out, t.Fg(t.P.Subtle, "  "+fit("  "+text(d.Source)+"  ·  "+text(d.FSType)+readOnly(d.ReadOnly)+inodes(d),
			maxInt(1, s.cols-3))))
	}
	return out
}

func readOnly(ro bool) string {
	if ro {
		return "  ·  только чтение"
	}
	return ""
}

func inodes(d sysinfo.Disk) string {
	if d.InodesTotal == 0 {
		return ""
	}
	return fmt.Sprintf("  ·  inode %d свободно из %d", d.InodesFree, d.InodesTotal)
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
	h.add("  "+fit("интерфейс", 14)+" "+fit("состоян.", 9)+" "+right("принято", 12)+" "+
		right("передано", 12)+" "+right("пак. вх.", 11)+" "+right("пак. исх.", 11),
		func(x string) string { return t.Fg(t.P.Border, x) })
	out = append(out, h.String())
	for _, n := range s.st.Network {
		var r row
		r.add("  "+fit(n.Name, 14)+" ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		r.add(fit(text(n.OperState), 9)+" ", func(x string) string { return t.Fg(stateColour(t, n.OperState), x) })
		r.add(right(report.UBytes(n.RxBytes), 12)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(report.UBytes(n.TxBytes), 12)+" ", func(x string) string { return t.Fg(t.P.Muted, x) })
		r.add(right(fmt.Sprint(n.RxPackets), 11)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		r.add(right(fmt.Sprint(n.TxPackets), 11), func(x string) string { return t.Fg(t.P.Subtle, x) })
		out = append(out, r.String())
		if extra := ifaceExtra(n); extra != "" {
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

func ifaceExtra(n sysinfo.Iface) string {
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
		parts = append(parts, fmt.Sprintf("%d Мбит/с", n.SpeedMbit))
	}
	if n.RxErrors+n.TxErrors+n.RxDropped+n.TxDropped > 0 {
		parts = append(parts, fmt.Sprintf("ошибок %d/%d, потеряно %d/%d",
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
		return []string{"", s.note("— (датчиков не нашлось)")}
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
		limit, against := 100.0, "из 100 °C"
		if sn.CritC > 0 {
			limit, against = sn.CritC, fmt.Sprintf("критич. %.1f °C", sn.CritC)
		}
		out = append(out, t.gauge(name, 26, sn.Celsius/limit, fmt.Sprintf("%.1f °C  %s", sn.Celsius, against), bw))
	}
	return out
}

func (s *screen) missing() []string {
	if !s.haveSt {
		return s.waiting()
	}
	t := s.t
	if len(s.st.Missing) == 0 {
		return []string{"", s.note("прочитано всё, чего ждали")}
	}
	// Go walks a map in a different order every time.  Sorting keeps the
	// section still between redraws instead of reshuffling once a second.
	keys := make([]string, 0, len(s.st.Missing))
	for k := range s.st.Missing {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := []string{"", s.note("источник назван вместе с причиной; нулём его отсутствие не притворяется"), ""}
	for _, k := range keys {
		var r row
		r.add("  "+fit(k, 24)+"  ", func(x string) string { return t.Fg(t.P.Yellow, x) })
		r.add(fit(s.st.Missing[k], maxInt(1, s.cols-r.w-1)), func(x string) string {
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
