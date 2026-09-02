// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

// What the screen shows about the machine itself: the page of what it is
// (СИСТЕМА), the page of its video cards (ВИДЕОКАРТЫ), and the block of the
// ЗАГРУЗКА page that takes its processors one by one.
//
// Nothing here reads a source of its own — the pages draw sysinfo.Status, the
// same snapshot the printed report prints, and where the snapshot has no
// number they draw the dash the report draws.

import (
	"fmt"
	"strings"
	"time"

	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/report"
	"digitdisk/internal/sysinfo"
)

// system is the page a person opens to see which machine this is: the mark of
// the system, and beside it the facts a person recognises a machine by.
func (s *screen) system() []string {
	if !s.haveSt {
		return s.waiting()
	}
	h := s.st.Host
	mark := emblemFor(h.DistroID, h.Distro)
	fields := s.systemFields()

	const labelWidth = 14
	art := mark.art
	width := emblemWidth(art)
	paint := func(x string) string { return s.t.Fg(mark.role(s.t.P), x) }

	// Side by side where the terminal is wide enough for the mark and a
	// readable value beside it; stacked where it is not.  A mark squeezed
	// against the values is worse than a mark above them.
	if s.cols < width+labelWidth+24 {
		out := []string{""}
		for _, line := range art {
			out = append(out, "  "+paint(line))
		}
		out = append(out, "")
		for _, f := range fields {
			out = append(out, s.kv(f[0], f[1]))
		}
		return out
	}

	out := []string{""}
	rows := len(art)
	if len(fields) > rows {
		rows = len(fields)
	}
	for i := 0; i < rows; i++ {
		var r row
		r.plain("  ")
		if i < len(art) {
			r.add(art[i], paint)
		} else {
			r.plain(strings.Repeat(" ", width))
		}
		r.plain("  ")
		if i < len(fields) {
			r.add(fit(fields[i][0], labelWidth)+"  ", func(x string) string { return s.t.Fg(s.t.P.Muted, x) })
			r.add(fit(fields[i][1], maxInt(1, s.cols-r.w-1)), func(x string) string {
				return s.t.Fg(s.t.P.Foreground, strings.TrimRight(x, " "))
			})
		}
		out = append(out, r.String())
	}
	return out
}

// systemFields is what stands beside the mark, in the order a person reads
// it: who and where first, then what the machine is made of.
func (s *screen) systemFields() [][2]string {
	st := s.st
	h := st.Host

	node := text(h.Hostname)
	if h.User != "" && h.Hostname != "" {
		node = h.User + "@" + h.Hostname
	}
	kernel := text(h.KernelRelease)
	if h.Machine != "" {
		kernel += " (" + h.Machine + ")"
	}
	if h.Bits > 0 {
		kernel += fmt.Sprintf(", %d-разрядная", h.Bits)
	}
	uptime := text(h.UptimeHuman)
	if h.UptimeSeconds > 0 && h.BootTime > 0 {
		uptime += " (с " + time.Unix(h.BootTime, 0).Format("2006-01-02 15:04") + ")"
	}
	cpu := text(h.CPUModel)
	if st.Load.CPUCount > 0 {
		cpu += fmt.Sprintf(" × %d", st.Load.CPUCount)
	}
	memory := dash
	if st.Memory.Total > 0 {
		f := pctOf(st.Memory.Used, st.Memory.Total)
		memory = fmt.Sprintf("%s из %s (%s)", report.UBytes(st.Memory.Used),
			report.UBytes(st.Memory.Total), percent(f))
	}

	return [][2]string{
		{"узел", node},
		{"дистрибутив", text(h.Distro)},
		{"модель", text(h.Model)},
		{"ядро", kernel},
		{"оболочка", text(h.Shell)},
		{"рабочий стол", text(h.Desktop)},
		{"терминал", text(h.Terminal)},
		{"время работы", uptime},
		{"процессор", cpu},
		{"память", memory},
		{"видеокарты", text(cardNames(st.GPUs))},
		{"снимок взят", text(st.TakenAt)},
	}
}

// cardNames is the one-line answer to "what graphics does this machine have".
func cardNames(cards []gpuinfo.Card) string {
	names := make([]string, 0, len(cards))
	for _, c := range cards {
		names = append(names, c.Name)
	}
	return strings.Join(names, " · ")
}

// coresBlock is the lower half of the ЗАГРУЗКА page: the processors one by
// one, under the share they add up to.
//
// It is not a page of its own on purpose.  The section strip is the printed
// report's, and in the report the per-processor figures are a line of
// ЗАГРУЗКА, not a section — so a tab here would be a tab with nothing behind
// it in the report, and the keys 1…9 would stop reaching the last section.
//
// A machine with four processors and a machine with two hundred and fifty-six
// are the same block here, and they cannot be drawn the same way: four fit as
// gauges with room to spare, and two hundred and fifty-six as gauges are a
// list nobody scrolls through.  So it draws whichever of the two fits the
// terminal it is on — the gauges while they fit, the map of cells when they
// stop fitting.
func (s *screen) coresBlock() []string {
	if !s.haveSt {
		return nil
	}
	t, l := s.t, s.st.Load
	out := []string{""}
	sum, ok := s.st.Cores()
	if !ok {
		out = append(out,
			t.gaugeUnmeasured("по ядрам", 12, dash, s.barWidth()),
			s.note("почему — digitdisk status --why"),
		)
		return out
	}

	out = append(out, s.kv("по ядрам", fmt.Sprintf("замерено %d из %s  (окно %d мс)",
		sum.Measured, count(l.CPUCount), l.SampleMillis)))
	out = append(out, s.kv("разброс", fmt.Sprintf("мин %.1f%% · медиана %.1f%% · макс %.1f%% (ядро %d)",
		sum.Min, sum.Median, sum.Max, sum.Busiest)))
	out = append(out, s.kv("под нагрузкой", fmt.Sprintf("%d ядер из %d заняты больше чем наполовину", sum.Loaded, sum.Total)))
	out = append(out, "")

	if per, rows := s.barGrid(sum.Total); rows > 0 {
		out = append(out, s.caption("ПО ЯДРАМ"))
		out = append(out, s.coreBars(per)...)
		return out
	}

	out = append(out, s.caption("КАРТА ЯДЕР"))
	out = append(out, s.coreMap()...)
	out = append(out, "")
	out = append(out, s.caption("САМЫЕ ЗАНЯТЫЕ"))
	out = append(out, s.busiestCores(6)...)
	return out
}

// coreCellWidth is one processor's cell in the grid of gauges: its number, a
// short bar and the share.
const coreCellWidth = 18

// barGrid decides whether the gauges fit, and in how many columns.  rows is
// zero when they do not fit and the map is to be drawn instead.
func (s *screen) barGrid(n int) (perRow, rows int) {
	perRow = (s.cols - 2) / coreCellWidth
	if perRow < 1 {
		return 0, 0
	}
	rows = (n + perRow - 1) / perRow
	// Five lines of the page are already spent on the numbers above the
	// grid, and a grid that does not fit under them is a grid a reader has
	// to scroll through core by core.
	if rows > maxInt(1, s.bodyHeight()-6) {
		return 0, 0
	}
	return perRow, rows
}

// coreBars draws every processor as its own small gauge, laid out in columns.
func (s *screen) coreBars(perRow int) []string {
	t := s.t
	bw := coreCellWidth - 10
	var out []string
	var r row
	for i, c := range s.st.Load.Cores {
		if i > 0 && i%perRow == 0 {
			out = append(out, r.String())
			r = row{}
		}
		if i%perRow == 0 {
			r.plain("  ")
		}
		r.add(right(fmt.Sprint(c.Index), 3)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		if c.BusyPercent == nil {
			r.plain(t.emptyBar(bw))
			r.w += bw
			r.add(right(dash, 5)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
			continue
		}
		f := *c.BusyPercent / 100
		r.plain(t.bar(f, bw))
		r.w += bw
		r.add(right(fmt.Sprintf("%.0f%%", *c.BusyPercent), 5)+" ", func(x string) string { return t.Fg(t.level(f), x) })
	}
	return append(out, r.String())
}

// coreMap draws one cell per processor: the same eight faces the history line
// uses, in the same colours, so a reader who has learned one has learned the
// other.  Two hundred and fifty-six processors are four rows of it.
func (s *screen) coreMap() []string {
	t := s.t
	cores := s.st.Load.Cores
	per := s.mapWidth()
	var out []string
	for start := 0; start < len(cores); start += per {
		end := start + per
		if end > len(cores) {
			end = len(cores)
		}
		var r row
		r.add(right(fmt.Sprintf("%d–%d", cores[start].Index, cores[end-1].Index), 9)+"  ",
			func(x string) string { return t.Fg(t.P.Subtle, x) })
		shares := make([]float64, 0, end-start)
		for _, c := range cores[start:end] {
			if c.BusyPercent == nil {
				shares = append(shares, -1)
				continue
			}
			shares = append(shares, *c.BusyPercent/100)
		}
		r.plain(t.spark(shares, len(shares)))
		r.w += len(shares)
		out = append(out, r.String())
	}
	return append(out, s.note("ячейка — ядро: ▁ пусто, █ занято целиком"))
}

// mapWidth is how many cells a row of the map holds.  It is a multiple of
// eight so a reader can count along it, and it never runs past the terminal.
func (s *screen) mapWidth() int {
	w := s.cols - 13
	w -= w % 8
	if w < 8 {
		w = 8
	}
	return w
}

// busiestCores lists the processors doing the most work, as full gauges: the
// map says how the load is spread, and this says where it is.
func (s *screen) busiestCores(n int) []string {
	cores := append([]sysinfo.Core(nil), s.st.Load.Cores...)
	sortCores(cores)
	if len(cores) > n {
		cores = cores[:n]
	}
	bw := s.barWidth()
	out := make([]string, 0, len(cores))
	for _, c := range cores {
		if c.BusyPercent == nil {
			out = append(out, s.t.gaugeUnmeasured("ядро "+fmt.Sprint(c.Index), 12, dash, bw))
			continue
		}
		out = append(out, s.t.gauge("ядро "+fmt.Sprint(c.Index), 12, *c.BusyPercent/100,
			percent(*c.BusyPercent/100), bw))
	}
	return out
}

// sortCores puts the busiest first, and an unmeasured processor last: a
// processor nobody measured is not a processor doing nothing.
func sortCores(cores []sysinfo.Core) {
	share := func(c sysinfo.Core) float64 {
		if c.BusyPercent == nil {
			return -1
		}
		return *c.BusyPercent
	}
	for i := 1; i < len(cores); i++ {
		for j := i; j > 0 && share(cores[j]) > share(cores[j-1]); j-- {
			cores[j], cores[j-1] = cores[j-1], cores[j]
		}
	}
}

// gpus is the page of the video cards.  Each card gets the same four readings
// — load, memory, temperature, power — and each of them is a gauge or a dash,
// never a zero standing in for a driver that publishes nothing.
func (s *screen) gpus() []string {
	if !s.haveSt {
		return s.waiting()
	}
	if len(s.st.GPUs) == 0 {
		return []string{"", s.note("— (видеокарт не нашлось)"), s.note("почему — digitdisk status --why")}
	}
	t := s.t
	bw := s.barWidth()
	out := []string{""}
	for i, c := range s.st.GPUs {
		if i > 0 {
			out = append(out, "")
		}
		out = append(out, s.caption(fmt.Sprintf("%d · %s", i, c.Name)))
		if c.BusyPercent != nil {
			f := *c.BusyPercent / 100
			out = append(out, t.gauge("занято", 12, f, percent(f), bw))
		} else {
			out = append(out, t.gaugeUnmeasured("занято", 12, dash, bw))
		}
		switch {
		case c.MemoryTotalBytes != nil && c.MemoryUsedBytes != nil:
			f := pctOf(*c.MemoryUsedBytes, *c.MemoryTotalBytes)
			out = append(out, t.gauge("память", 12, f, s.reading(
				fmt.Sprintf("%s из %s  (%s)", report.UBytes(*c.MemoryUsedBytes),
					report.UBytes(*c.MemoryTotalBytes), percent(f)), percent(f)), bw))
		case c.MemoryTotalBytes != nil:
			out = append(out, t.gaugeUnmeasured("память", 12,
				"всего "+report.UBytes(*c.MemoryTotalBytes)+", занято "+dash, bw))
		default:
			out = append(out, t.gaugeUnmeasured("память", 12, dash, bw))
		}
		out = append(out, s.kv("температура", cardTemp(c)))
		if line := cardPower(c); line != "" {
			out = append(out, s.kv("питание", line))
		}
		out = append(out, s.note(cardOrigin(c)))
	}
	return out
}

// cardTemp and cardPower are the readings that are not shares: what the card
// is heated to, and what it is drawing while it does it.  Either can be
// missing on its own, so they are two lines and not one sentence with a hole
// in the middle.
func cardTemp(c gpuinfo.Card) string {
	if c.Celsius == nil {
		return dash
	}
	return fmt.Sprintf("%.1f °C", *c.Celsius)
}

func cardPower(c gpuinfo.Card) string {
	var parts []string
	if c.Watts != nil {
		parts = append(parts, fmt.Sprintf("%.1f Вт", *c.Watts))
	}
	if c.MHz != nil {
		parts = append(parts, fmt.Sprintf("%.0f МГц", *c.MHz))
	}
	return strings.Join(parts, "  ·  ")
}

// cardOrigin says where the card's numbers came from.  It is not an
// explanation of anything missing — those live behind `--why` — it is the
// provenance of what is shown, and a number that came out of somebody else's
// program says so on the line under itself.
func cardOrigin(c gpuinfo.Card) string {
	parts := []string{}
	if c.Bus != "" {
		parts = append(parts, "шина "+c.Bus)
	}
	if c.Driver != "" {
		parts = append(parts, "драйвер "+c.Driver)
	}
	if c.Outside {
		parts = append(parts, "числа от чужой программы "+c.Source)
	} else if c.Source != "" {
		parts = append(parts, "числа из "+c.Source)
	}
	return strings.Join(parts, "  ·  ")
}

// gpuGauges is the short form of the same page for the opening screen: one
// line per card, load and memory on it, and no more than three cards.
func (s *screen) gpuGauges(limit int) []string {
	t := s.t
	bw := s.barWidth()
	var out []string
	for i, c := range s.st.GPUs {
		if i >= limit {
			out = append(out, s.note(fmt.Sprintf("…и ещё %d — раздел ВИДЕОКАРТЫ", len(s.st.GPUs)-i)))
			break
		}
		name := fit(c.Name, 20)
		reading := []string{}
		if c.MemoryTotalBytes != nil && c.MemoryUsedBytes != nil {
			reading = append(reading, fmt.Sprintf("%s из %s", report.UBytes(*c.MemoryUsedBytes),
				report.UBytes(*c.MemoryTotalBytes)))
		}
		if c.Celsius != nil {
			reading = append(reading, fmt.Sprintf("%.0f °C", *c.Celsius))
		}
		if c.BusyPercent == nil {
			out = append(out, t.gaugeUnmeasured(name, 20, strings.Join(append([]string{dash}, reading...), "  ·  "), bw))
			continue
		}
		f := *c.BusyPercent / 100
		out = append(out, t.gauge(name, 20, f, strings.Join(append([]string{percent(f)}, reading...), "  ·  "), bw))
	}
	return out
}

// coreComb is the map of the processors squeezed onto one line, for the
// opening screen: the shape of the load across the machine, under the gauge
// that gives its total.
func (s *screen) coreComb(nameWidth int) string {
	t := s.t
	var r row
	r.add("  "+fit("по ядрам", nameWidth)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	w := s.barWidth()
	shares := make([]float64, 0, len(s.st.Load.Cores))
	for _, c := range s.st.Load.Cores {
		if c.BusyPercent == nil {
			shares = append(shares, -1)
			continue
		}
		shares = append(shares, *c.BusyPercent/100)
	}
	// More processors than cells: each cell then stands for a group of
	// them, and it takes the largest share in the group — a screen that
	// averaged them would hide the one processor that is on fire.
	r.plain(t.spark(squeeze(shares, w), w))
	r.w += w
	if len(shares) > w {
		r.add(fmt.Sprintf("  по %d ядер в ячейке", (len(shares)+w-1)/w),
			func(x string) string { return t.Fg(t.P.Subtle, x) })
	}
	return r.String()
}

// squeeze folds a long list of shares into n cells, each cell the largest of
// the shares it covers.
func squeeze(shares []float64, n int) []float64 {
	if n <= 0 || len(shares) <= n {
		return shares
	}
	per := (len(shares) + n - 1) / n
	out := make([]float64, 0, n)
	for i := 0; i < len(shares); i += per {
		end := i + per
		if end > len(shares) {
			end = len(shares)
		}
		best := -1.0
		for _, v := range shares[i:end] {
			if v > best {
				best = v
			}
		}
		out = append(out, best)
	}
	return out
}
