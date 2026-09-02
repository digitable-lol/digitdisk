// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"digitdisk/internal/lang"
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
	// Lang is the language the screen is drawn in when it opens.  The
	// reader may change it with one key while it is up.
	Lang lang.Lang
	// Remember stores a language the reader chose here, and answers with
	// the line to show them — «язык сохранён: /path», or why it was not.
	// It is handed in rather than called for, so that this package knows
	// nothing about home directories and a test can watch what it was
	// asked to store.  A nil Remember means the key still switches the
	// screen and nothing is written anywhere.
	Remember func(lang.Lang) lang.Phrase
}

// ErrNoTerminal is returned when the screen was asked for where it cannot be
// drawn.  The caller prints text instead, or says so and stops.
//
// It is a lang.Error rather than a plain one because it is shown to a person:
// the wording travels as a Phrase and is rendered where it is printed, while
// `err == ErrNoTerminal` goes on working — the value is still this one.
var ErrNoTerminal = lang.Errorf("вывод не в терминал: живой экран невозможен")

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

	st     sysinfo.Status
	haveSt bool
	// l is the language every line of this screen is written in.
	l lang.Lang
	// said is the one-line answer to the last key that did something
	// outside the screen — storing the language — and when it was said.
	// It sits in the footer for a few seconds and then goes away: a
	// notice that never leaves is a notice that stops being read.
	said    lang.Phrase
	saidAt  time.Time
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
		return lang.Errorf("живому экрану не передан сборщик снимка")
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

	s := &screen{o: o, l: o.Lang, t: NewTheme(o.Palette), out: bufio.NewWriterSize(o.Out, 1<<16), tty: tty}
	if !s.l.Valid() {
		s.l = lang.Default
	}
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
		case 'l', 'L', 'д', 'Д':
			// The language of the whole screen, switched where the
			// reader is looking at it, and stored so that the next
			// run — and `digitdisk clean` tomorrow — speaks the same
			// language.  «l» for language and «д» for the same key on
			// a Russian layout, beside the other letters this screen
			// answers in both alphabets.
			s.l = s.l.Other()
			if s.o.Remember != nil {
				s.said, s.saidAt = s.o.Remember(s.l), time.Now()
			}
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
		more = s.l.F("  строки %d–%d из %d", s.scroll+1, s.scroll+len(shown), len(body))
	}
	// What was just stored in the reader's home directory outranks the
	// line count for a few seconds: writing a file is news, and «строки
	// 1–20 из 40» is not.
	if !s.said.Empty() && time.Since(s.saidAt) < 6*time.Second {
		more = "  " + s.said.In(s.l)
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
			width += runes(sec.title(s.l)) + 2 + gap
		}
		if width+1 > s.cols {
			continue
		}
		var r row
		r.plain(" ")
		for i, sec := range sections {
			if i == s.tab {
				r.add(" "+sec.title(s.l)+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
			} else {
				r.add(" "+sec.title(s.l)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
			}
			r.plain(strings.Repeat(" ", gap))
		}
		return r.String()
	}

	cur := sections[s.tab]
	var r row
	r.add(" ‹ ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	r.add(" "+cur.title(s.l)+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
	r.add(fmt.Sprintf(" ›  %d/%d", s.tab+1, len(sections)), func(x string) string { return t.Fg(t.P.Subtle, x) })
	return r.String()
}

// footer is the state of the screen and the keys that move it.
func (s *screen) footer(more string) string {
	t := s.t
	var r row
	switch {
	case s.paused:
		r.add(" "+s.l.T("ПАУЗА")+" ", func(x string) string { return t.Chip(t.P.Yellow, x) })
	case s.busying || !s.haveSt:
		r.add(" "+s.l.T("ЗАМЕР")+" ", func(x string) string { return t.Chip(t.P.AccentSoft, x) })
	default:
		r.add(" "+s.l.T("ЖИВОЙ")+" ", func(x string) string { return t.Chip(t.P.Green, x) })
	}
	// The keys come first in the budget: a full-screen program that does not
	// say how to leave it is a trap, so the state of the measurement is what
	// gives way on a narrow terminal, never the way out.
	exit := s.l.T("q выход ")
	if s.haveSt {
		ago, took, period := s.l.Since(time.Since(s.taken)), s.l.Millis(s.took), s.l.Every(s.o.Interval)
		states := []string{
			s.l.F("  замер %s назад · длился %s · каждые %s", ago, took, period),
			s.l.F("  замер %s назад · каждые %s", ago, period),
			s.l.F("  замер %s назад", ago),
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

	// Every width of the hint is its own wording: an English line is not a
	// Russian line with the words swapped, and one that had to be cut to fit
	// would be cut in a different place.
	hints := []string{
		s.l.T("← → разделы · ↑ ↓ прокрутка · p пауза · r замер · l язык · ? команды · q выход "),
		s.l.T("← → разделы · p пауза · r замер · l язык · ? команды · q выход "),
		s.l.T("← → · p · r · l язык · ? команды · q выход "),
		s.l.T("← → · p · r · l · q выход "),
		exit,
	}
	if s.menu {
		hints = []string{
			s.l.T("↑ ↓ прокрутка · ? или Esc назад · q выход "),
			s.l.T("? назад · q выход "),
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
