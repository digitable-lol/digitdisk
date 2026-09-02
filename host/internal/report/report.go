// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package report renders snapshots and walk results for a human reader.
// Absent data is printed as "—" so an empty field is never mistaken for a
// measured zero.
package report

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"digitdisk/internal/core"
	"digitdisk/internal/procfs"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
)

// Bytes renders a byte count with a binary-prefix unit.
func Bytes(n int64) string {
	neg := ""
	if n < 0 {
		neg, n = "-", -n
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%s%d Б", neg, n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit && exp < 4; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%s%.1f %s", neg, float64(n)/float64(div), [...]string{"КиБ", "МиБ", "ГиБ", "ТиБ", "ПиБ"}[exp])
}

// UBytes is Bytes for unsigned counts.
func UBytes(n uint64) string { return Bytes(int64(n)) }

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
func Status(w io.Writer, st sysinfo.Status) {
	p := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	p("СИСТЕМА")
	p("  узел          %s", dash(st.Host.Hostname))
	p("  дистрибутив   %s", dash(st.Host.Distro))
	if st.Host.Model != "" {
		p("  модель        %s", st.Host.Model)
	}
	p("  ядро          %s (%s)", dash(st.Host.KernelRelease), dash(st.Host.Machine))
	if st.Host.UptimeSeconds > 0 {
		boot := "—"
		if st.Host.BootTime > 0 {
			boot = time.Unix(st.Host.BootTime, 0).Format("2006-01-02 15:04")
		}
		p("  время работы  %s (с %s)", dash(st.Host.UptimeHuman), boot)
	} else {
		p("  время работы  —")
	}

	p("")
	p("ЗАГРУЗКА")
	p("  средняя       %.2f / %.2f / %.2f  (1/5/15 мин)", st.Load.One, st.Load.Five, st.Load.Fifteen)
	if st.Load.CPUCount > 0 {
		p("  ядер          %d", st.Load.CPUCount)
	} else {
		p("  ядер          —")
	}
	switch why, unmeasured := st.Unmeasured(sysinfo.FactCPUBusy); {
	case st.Load.BusyPercent != nil:
		p("  занято ЦП     %.1f%% (замер %d мс)", *st.Load.BusyPercent, st.Load.SampleMillis)
	case unmeasured:
		p("  занято ЦП     — (%s)", why)
	default:
		p("  занято ЦП     — (замер не делался)")
	}

	p("")
	p("ПАМЯТЬ")
	m := st.Memory
	if m.Total == 0 || !m.Has(procfs.FieldTotal) {
		p("  —")
	} else {
		p("  всего         %s", UBytes(m.Total))
		if m.Has(procfs.FieldUsed) {
			p("  занято        %s (%.1f%%)  [всего − доступно]", UBytes(m.Used), pct(m.Used, m.Total))
		} else {
			p("  занято        —")
		}
		p("  свободно      %s", measured(m, procfs.FieldFree, m.Free))
		p("  кэш/буферы    %s  (в т.ч. разделяемая %s)",
			measured(m, procfs.FieldBuffCache, m.BuffCache), measured(m, procfs.FieldShared, m.Shared))
		p("  доступно      %s", measured(m, procfs.FieldAvailable, m.Available))
		switch {
		case !m.Has(procfs.FieldSwapTotal):
			p("  своп          —")
		case m.SwapTotal > 0:
			p("  своп          %s из %s занято", measured(m, procfs.FieldSwapUsed, m.SwapUsed), UBytes(m.SwapTotal))
		default:
			p("  своп          нет")
		}
	}

	p("")
	pr := st.Processes
	p("ПРОЦЕССЫ  всего %d, потоков %s, выполняется %d, заблокировано %s%s",
		pr.Total, counted(st, sysinfo.FactThreads, pr.Threads), pr.Running,
		counted(st, sysinfo.FactBlocked, pr.Blocked), plural(pr.Unreadable))
	if len(pr.TopByMemory) > 0 {
		p("  десятка по памяти:")
		for _, x := range pr.TopByMemory {
			p("    %7d %-12s %10s  %s", x.PID, cut(dash(x.User), 12), Bytes(x.RSSBytes), cut(dash(firstNonEmpty(x.Cmdline, x.Comm)), 64))
		}
	}
	if len(pr.TopByCPU) > 0 {
		p("  десятка по процессору:")
		for _, x := range pr.TopByCPU {
			cpu := "    —"
			if x.CPUPercent != nil {
				cpu = fmt.Sprintf("%5.1f%%", *x.CPUPercent)
			}
			p("    %7d %-12s %9s  %s", x.PID, cut(dash(x.User), 12), cpu, cut(dash(firstNonEmpty(x.Cmdline, x.Comm)), 64))
		}
	}

	p("")
	p("ДИСКИ")
	if len(st.Disks) == 0 {
		p("  —")
	} else {
		p("  %-28s %-22s %10s %10s %10s %6s", "точка монтирования", "устройство", "размер", "занято", "свободно", "занят")
		for _, d := range st.Disks {
			if d.Error != "" {
				p("  %-28s %-22s  ошибка: %s", cut(d.MountPoint, 28), cut(d.Source, 22), d.Error)
				continue
			}
			p("  %-28s %-22s %10s %10s %10s %5.1f%%", cut(d.MountPoint, 28), cut(d.Source, 22),
				UBytes(d.TotalBytes), UBytes(d.UsedBytes), UBytes(d.AvailableBytes), d.UsePercent)
		}
	}

	p("")
	p("СЕТЬ")
	if len(st.Network) == 0 {
		p("  —")
	} else {
		p("  %-12s %-8s %12s %12s %10s %10s", "интерфейс", "состоян.", "принято", "передано", "пак. вх.", "пак. исх.")
		_, noCounters := st.Unmeasured(sysinfo.FactNetCounters)
		for _, n := range st.Network {
			if noCounters {
				p("  %-12s %-8s %12s %12s %10s %10s", cut(n.Name, 12), cut(dash(n.OperState), 8), "—", "—", "—", "—")
				continue
			}
			p("  %-12s %-8s %12s %12s %10d %10d", cut(n.Name, 12), cut(dash(n.OperState), 8),
				UBytes(n.RxBytes), UBytes(n.TxBytes), n.RxPackets, n.TxPackets)
		}
	}

	p("")
	p("ТЕМПЕРАТУРА")
	if len(st.Sensors) == 0 {
		if why, ok := st.Unmeasured(sysinfo.FactSensors); ok {
			p("  — (%s)", why)
		} else {
			p("  —")
		}
	} else {
		for _, s := range st.Sensors {
			extra := ""
			if s.CritC > 0 {
				extra = fmt.Sprintf("  (критич. %.1f°C)", s.CritC)
			}
			p("  %-18s %-16s %6.1f°C%s", cut(s.Chip, 18), cut(s.Label, 16), s.Celsius, extra)
		}
	}

	// Sorted, because this section is read to compare two runs or two
	// machines, and Go's map order would shuffle it on every run.
	if len(st.Missing) > 0 {
		p("")
		p("НЕ ИЗМЕРЕНО (и почему)")
		names := make([]string, 0, len(st.Missing))
		for k := range st.Missing {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, k := range names {
			p("  %s: %s", k, st.Missing[k])
		}
	}
}

// Analyze prints the result of a tree walk.
func Analyze(w io.Writer, r scan.Result) {
	p := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	p("ОБХОД  %s", r.Root)
	p("  записей       %d  (файлов %d, каталогов %d, ссылок %d, прочего %d)",
		r.Entries, r.Files, r.Dirs, r.Links, r.Others)
	p("  объём         %s (%d Б, видимый размер: du --apparent-size, он же du -sb)", Bytes(r.TotalBytes), r.TotalBytes)
	p("                файлы %s, ссылки %s; сверх того сами каталоги %s (в объём не входят, как и у du)",
		Bytes(r.FileBytes), Bytes(r.LinkBytes), Bytes(r.DirBytes))
	if r.HardlinkDupes > 0 {
		p("  жёсткие ссылки %d повторных имён не засчитано (%s)", r.HardlinkDupes, Bytes(r.HardlinkBytes))
	}
	s := r.Skipped
	p("  пропущено     %d  (нет доступа %d, исчезло %d, иные ошибки %d, граница ФС %d, предел глубины %d)",
		s.Total(), s.PermissionDenied, s.Vanished, s.OtherErrors, s.DeviceBoundaries, s.DepthLimited)
	p("  время         %.2f с", r.DurationSeconds)

	p("")
	p("РЕШАЮЩИЙ СЛОЙ  %s, договор версии %d", r.Decider, r.ContractVersion)
	if !r.DeciderReady {
		p("  ВНИМАНИЕ: настоящий разбор не выполнялся — все записи возвращены как")
		p("  Неизвестное/НеТрогать. Разряды ниже показывают работу стыковки, не анализ.")
	}

	p("")
	p("ПО РАЗРЯДАМ")
	for _, c := range core.Classes {
		b := r.ByClass[c]
		p("  %-14s %8d %12s", c, b.Count, Bytes(b.Bytes))
	}
	p("")
	p("ПО ПРИГОВОРАМ")
	for _, v := range core.Verdicts {
		b := r.ByVerdict[v]
		p("  %-14s %8d %12s", v, b.Count, Bytes(b.Bytes))
	}

	p("")
	rem := r.ByVerdict[core.VerdictRemovable]
	p("ПРЕДЛОЖЕНО УБРАТЬ  %d записей, %s", rem.Count, Bytes(rem.Bytes))
	p("  (analyze только считает; убирает `digitdisk clean`, и тоже не сразу)")
	if len(r.Removable) == 0 {
		p("  — нечего")
	} else {
		for _, e := range r.Removable {
			p("  %10s  %-12s %5.0f дн  %s", Bytes(e.Size), e.Class, e.AgeDays, cut(e.Path, 80))
		}
	}

	if len(r.Largest) > 0 {
		p("")
		p("САМОЕ КРУПНОЕ")
		for _, e := range r.Largest {
			p("  %10s  %-9s %5.0f дн  %s", Bytes(e.Size), e.Kind, e.AgeDays, cut(e.Path, 80))
		}
	}
}

// measured renders a memory field: its value when the system published it,
// "—" when it did not.  A field nobody measured must not print as "0 Б", which
// is a number a reader would believe.
func measured(m procfs.Memory, field string, v uint64) string {
	if !m.Has(field) {
		return "—"
	}
	return UBytes(v)
}

// counted renders a counter the running system may not publish at all.
func counted(st sysinfo.Status, fact string, v int) string {
	if _, unmeasured := st.Unmeasured(fact); unmeasured {
		return "—"
	}
	return strconv.Itoa(v)
}

func pct(a, b uint64) float64 {
	if b == 0 {
		return 0
	}
	return 100 * float64(a) / float64(b)
}

func plural(unreadable int) string {
	if unreadable == 0 {
		return ""
	}
	return fmt.Sprintf(", не прочитано %d", unreadable)
}

func firstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}
