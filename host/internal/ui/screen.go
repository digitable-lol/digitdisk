// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"digitdisk/internal/sysinfo"
)

// Options is everything the screen needs.  Collect is handed in rather than
// called for: the screen owns no source of its own and takes exactly the
// snapshot the printed report takes.
type Options struct {
	Out      *os.File
	Interval time.Duration
	Palette  Palette
	Collect  func() sysinfo.Status
}

// ErrNoTerminal is returned when the screen was asked for where it cannot be
// drawn.  The caller prints text instead, or says so and stops.
var ErrNoTerminal = errors.New("вывод не в терминал: живой экран невозможен")

// Available reports whether the live screen may be drawn on out.  It is the
// whole of the default rule: a pipe, a file, /dev/null, an empty TERM and the
// dumb terminal all answer no, and `digitdisk status` then prints text exactly
// as it always has.
func Available(out *os.File) bool { return UsableTERM() && IsTerminal(out) }

const (
	altOn     = "\x1b[?1049h"
	altOff    = "\x1b[?1049l"
	hideCur   = "\x1b[?25l"
	showCur   = "\x1b[?25h"
	home      = "\x1b[H"
	clearLine = "\x1b[K"
	clearDown = "\x1b[J"
	minCols   = 40
	minRows   = 8
	histLen   = 240
)

type screen struct {
	o    Options
	t    Theme
	out  *bufio.Writer
	tty  *os.File
	rows int
	cols int

	st      sysinfo.Status
	haveSt  bool
	taken   time.Time
	took    time.Duration
	tab     int
	scroll  int
	paused  bool
	busying bool
	// menu shows the list of subcommands over the body.  The screen is
	// `status`, which reads and writes nothing; the list therefore only
	// names commands and never runs one — see commandsPage.
	menu bool

	cpuHist []float64
	memHist []float64
}

// Run draws the live screen until the reader asks to leave.  The terminal is
// handed back the way it was found on every path out, including a signal.
func Run(o Options) error {
	if o.Out == nil {
		o.Out = os.Stdout
	}
	if o.Interval <= 0 {
		o.Interval = 2 * time.Second
	}
	if o.Collect == nil {
		return errors.New("живому экрану не передан сборщик снимка")
	}
	if !Available(o.Out) {
		return ErrNoTerminal
	}

	tty, closeTTY, err := openInput()
	if err != nil {
		return ErrNoTerminal
	}
	defer closeTTY()

	restore, err := Raw(tty)
	if err != nil {
		return ErrNoTerminal
	}

	s := &screen{o: o, t: NewTheme(o.Palette), out: bufio.NewWriterSize(o.Out, 1<<16), tty: tty}
	s.rows, s.cols, _ = Size(o.Out)
	if s.rows < minRows {
		s.rows = 24
	}
	if s.cols < minCols {
		s.cols = 80
	}

	// The terminal is put back exactly once, whichever way this function
	// leaves — a return, a signal, or a panic on the way through.
	leave := func() {
		fmt.Fprint(o.Out, showCur+altOff)
		restore()
	}
	defer leave()

	sig := make(chan os.Signal, 4)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)
	defer signal.Stop(sig)

	fmt.Fprint(o.Out, altOn+hideCur)

	keys := make(chan key, 16)
	go readKeys(tty, keys)

	snaps := make(chan sample, 1)
	s.request(snaps)
	s.draw()

	// The next measurement is armed once the last one has landed, not on a
	// free-running tick.  On a busy machine a snapshot can take longer than
	// the period, and a ticker would then start one the moment the last
	// finished and hold the processor down for as long as the screen is up.
	tick := time.NewTimer(o.Interval)
	defer tick.Stop()
	rearm := func() {
		if !tick.Stop() {
			select {
			case <-tick.C:
			default:
			}
		}
		tick.Reset(o.Interval)
	}
	blink := time.NewTicker(time.Second)
	defer blink.Stop()

	for {
		select {
		case sm := <-snaps:
			s.accept(sm)
			rearm()
			s.draw()
		case <-tick.C:
			if s.paused {
				rearm()
				continue
			}
			s.request(snaps)
		case <-blink.C:
			s.draw() // the clock and the "measured N s ago" line move on
		case k := <-keys:
			if s.handle(k, snaps) {
				return nil
			}
			s.draw()
		case sg := <-sig:
			if sg == syscall.SIGWINCH {
				if r, c, ok := Size(o.Out); ok {
					s.rows, s.cols = r, c
					s.scroll = 0
				}
				s.draw()
				continue
			}
			return nil
		}
	}
}

// openInput finds a terminal to read keys from.  /dev/tty is asked first so
// the screen still answers the keyboard when standard input is taken by
// something else; standard input itself is the fallback.
func openInput() (*os.File, func(), error) {
	if f, err := os.Open("/dev/tty"); err == nil {
		return f, func() { f.Close() }, nil
	}
	if IsTerminal(os.Stdin) {
		return os.Stdin, func() {}, nil
	}
	return nil, nil, ErrNoTerminal
}

// A sample is one snapshot together with how long it took to take.  The
// duration is shown on the screen because on a machine with many thousands of
// processes it is the honest answer to "why is this not instant".
type sample struct {
	st   sysinfo.Status
	took time.Duration
}

// request takes the next snapshot away from the drawing loop, so a slow
// source never freezes the keyboard.
func (s *screen) request(out chan<- sample) {
	if s.busying {
		return
	}
	s.busying = true
	collect := s.o.Collect
	go func() {
		start := time.Now()
		st := collect()
		out <- sample{st: st, took: time.Since(start)}
	}()
}

func (s *screen) accept(sm sample) {
	st := sm.st
	s.busying = false
	s.st, s.haveSt, s.taken, s.took = st, true, time.Now(), sm.took

	// The history is the samples already shown, kept as they scroll by.  An
	// unmeasured share is stored as -1 and drawn as a gap, never as a zero.
	cpu := -1.0
	if st.Load.BusyPercent != nil {
		cpu = *st.Load.BusyPercent / 100
	}
	s.cpuHist = push(s.cpuHist, cpu)
	mem := -1.0
	if st.Memory.Total > 0 {
		mem = pctOf(st.Memory.Used, st.Memory.Total)
	}
	s.memHist = push(s.memHist, mem)
}

func push(h []float64, v float64) []float64 {
	h = append(h, v)
	if len(h) > histLen {
		h = h[len(h)-histLen:]
	}
	return h
}

// handle answers one key.  It reports whether the screen should close.
func (s *screen) handle(k key, snaps chan sample) bool {
	switch k.kind {
	case keyRune:
		switch k.r {
		case 'q', 'Q', 'й', 'Й':
			return true
		case 'p', 'P', 'з', 'З', ' ':
			s.paused = !s.paused
		case 'r', 'R', 'к', 'К':
			s.request(snaps)
		case 'j':
			s.scroll++
		case 'k':
			s.scroll--
		case 'g':
			s.scroll = 0
		case '?':
			s.menu = !s.menu
			s.scroll = 0
		}
		if k.r >= '1' && k.r <= '9' {
			if n := int(k.r - '1'); n < len(sections) {
				s.tab, s.scroll, s.menu = n, 0, false
			}
		}
	case keyEsc:
		// Esc backs out of the list first: a key that both closes an
		// overlay and quits the program would quit it by surprise.
		if s.menu {
			s.menu, s.scroll = false, 0
			return false
		}
		return true
	case keyCtrlC:
		return true
	case keyRight, keyTab:
		s.tab = (s.tab + 1) % len(sections)
		s.scroll, s.menu = 0, false
	case keyLeft, keyShiftTab:
		s.tab = (s.tab + len(sections) - 1) % len(sections)
		s.scroll, s.menu = 0, false
	case keyDown:
		s.scroll++
	case keyUp:
		s.scroll--
	case keyPgDn:
		s.scroll += s.bodyHeight()
	case keyPgUp:
		s.scroll -= s.bodyHeight()
	}
	if s.scroll < 0 {
		s.scroll = 0
	}
	return false
}

// bodyHeight is how many lines the section itself gets: the screen less the
// two-line head, the rule and the footer.
func (s *screen) bodyHeight() int {
	h := s.rows - 5
	if h < 1 {
		h = 1
	}
	return h
}

func (s *screen) draw() {
	var b strings.Builder
	b.WriteString(home)
	b.WriteString(s.t.Canvas())

	lines := s.frame()
	for i := 0; i < s.rows; i++ {
		b.WriteString(s.t.Canvas())
		if i < len(lines) {
			b.WriteString(s.t.clip(lines[i], s.cols))
		}
		b.WriteString(clearLine)
		if i < s.rows-1 {
			b.WriteString("\r\n")
		}
	}
	b.WriteString(clearDown)
	s.out.WriteString(b.String())
	s.out.Flush()
}

// frame lays out the whole screen: the brand head, the sections, the body of
// the chosen one, and the keys.
func (s *screen) frame() []string {
	t := s.t
	out := make([]string, 0, s.rows)
	out = append(out, s.header())
	out = append(out, s.tabs())
	out = append(out, t.Fg(t.P.Border, strings.Repeat("─", s.cols)))

	body := sections[s.tab].render(s)
	if s.menu {
		body = s.commandsPage()
	}
	h := s.bodyHeight()
	if s.scroll > len(body)-1 {
		s.scroll = len(body) - 1
	}
	if s.scroll < 0 {
		s.scroll = 0
	}
	shown := body[s.scroll:]
	if len(shown) > h {
		shown = shown[:h]
	}
	out = append(out, shown...)
	for i := len(shown); i < h; i++ {
		out = append(out, "")
	}

	more := ""
	if len(body) > h {
		more = fmt.Sprintf("  строки %d–%d из %d", s.scroll+1, s.scroll+len(shown), len(body))
	}
	out = append(out, t.Fg(t.P.Border, strings.Repeat("─", s.cols)))
	out = append(out, s.footer(more))
	return out
}

// header is the brand line: the diamond of Digitable, the tool's name, the
// machine it is looking at, and the clock.
func (s *screen) header() string {
	t := s.t
	var r row
	r.add(" ◇ digitdisk ", func(x string) string { return t.Chip(t.P.Accent, x) })
	r.plain(" ")

	if s.haveSt {
		h := s.st.Host
		parts := []string{}
		for _, p := range []string{h.Hostname, h.Distro, h.KernelRelease} {
			if strings.TrimSpace(p) != "" {
				parts = append(parts, p)
			}
		}
		if len(parts) > 0 {
			r.add(fit(strings.Join(parts, " · "), maxInt(0, s.cols-r.w-24)), func(x string) string {
				return t.Fg(t.P.Muted, strings.TrimRight(x, " "))
			})
		}
	}

	tail := "Digitable  " + time.Now().Format("15:04:05") + " "
	r.pad(maxInt(r.w, s.cols-runes(tail)))
	r.add("Digitable  ", func(x string) string { return t.Fg(t.P.AccentSoft, x) })
	r.add(time.Now().Format("15:04:05")+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	return r.String()
}

// tabs is the strip of sections.  A narrow terminal gets the current section
// alone with arrows, rather than a strip cut off in the middle of a word.
func (s *screen) tabs() string {
	t := s.t
	// The strip is drawn in full where it fits, tightened by one space where
	// it nearly fits, and only then given up for the current section alone.
	// A terminal of ninety-odd columns is common, and it is worth a space.
	for _, gap := range []int{2, 1} {
		width := 1
		for _, sec := range sections {
			width += runes(sec.title) + 2 + gap
		}
		if width+1 > s.cols {
			continue
		}
		var r row
		r.plain(" ")
		for i, sec := range sections {
			if i == s.tab {
				r.add(" "+sec.title+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
			} else {
				r.add(" "+sec.title+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
			}
			r.plain(strings.Repeat(" ", gap))
		}
		return r.String()
	}

	cur := sections[s.tab]
	var r row
	r.add(" ‹ ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	r.add(" "+cur.title+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
	r.add(fmt.Sprintf(" ›  %d/%d", s.tab+1, len(sections)), func(x string) string { return t.Fg(t.P.Subtle, x) })
	return r.String()
}

// footer is the state of the screen and the keys that move it.
func (s *screen) footer(more string) string {
	t := s.t
	var r row
	switch {
	case s.paused:
		r.add(" ПАУЗА ", func(x string) string { return t.Chip(t.P.Yellow, x) })
	case s.busying || !s.haveSt:
		r.add(" ЗАМЕР ", func(x string) string { return t.Chip(t.P.AccentSoft, x) })
	default:
		r.add(" ЖИВОЙ ", func(x string) string { return t.Chip(t.P.Green, x) })
	}
	// The keys come first in the budget: a full-screen program that does not
	// say how to leave it is a trap, so the state of the measurement is what
	// gives way on a narrow terminal, never the way out.
	const exit = "q выход "
	if s.haveSt {
		states := []string{
			fmt.Sprintf("  замер %s назад · длился %s · каждые %s", since(s.taken), lasted(s.took), every(s.o.Interval)),
			fmt.Sprintf("  замер %s назад · каждые %s", since(s.taken), every(s.o.Interval)),
			fmt.Sprintf("  замер %s назад", since(s.taken)),
			"",
		}
		budget := s.cols - r.w - runes(exit) - 2
		for _, state := range states {
			if runes(state) <= budget {
				r.add(state, func(x string) string { return t.Fg(t.P.Subtle, x) })
				break
			}
		}
	}
	if more != "" && r.w+runes(more)+runes(exit)+2 <= s.cols {
		r.add(more, func(x string) string { return t.Fg(t.P.Subtle, x) })
	}

	hints := []string{
		"← → разделы · ↑ ↓ прокрутка · p пауза · r замер · ? команды · q выход ",
		"← → разделы · p пауза · r замер · ? команды · q выход ",
		"← → · p · r · ? команды · q выход ",
		"← → · p · r · q выход ",
		exit,
	}
	if s.menu {
		hints = []string{
			"↑ ↓ прокрутка · ? или Esc назад · q выход ",
			"? назад · q выход ",
			exit,
		}
	}
	for _, hint := range hints {
		if r.w+runes(hint)+2 <= s.cols {
			r.pad(s.cols - runes(hint))
			r.add(hint, func(x string) string { return t.Fg(t.P.Muted, x) })
			break
		}
	}
	return r.String()
}

func since(at time.Time) string {
	d := time.Since(at)
	if d < time.Second {
		return "0 с"
	}
	if d < time.Minute {
		return fmt.Sprintf("%d с", int(d.Seconds()))
	}
	return fmt.Sprintf("%d мин", int(d.Minutes()))
}

// lasted says how long a measurement ran, in the units that suit its size.
func lasted(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d мс", d.Milliseconds())
	}
	return fmt.Sprintf("%.1f с", d.Seconds())
}

func every(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d мс", d.Milliseconds())
	}
	return fmt.Sprintf("%g с", d.Seconds())
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
