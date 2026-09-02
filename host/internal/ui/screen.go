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

	"digitdisk/internal/cli"
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
	// pick is the line the cursor stands on in the КОМАНДЫ section: an
	// index into cli.Commands, not into anything drawn, so the cursor
	// cannot come to point at a command that is not there.
	pick int
	// req is the command the reader chose and this screen cannot run
	// itself.  Setting it closes the screen; Run hands it to the caller,
	// who runs the subcommand and opens the screen again.
	req *Request

	cpuHist []float64
	memHist []float64
}

// A Request is the подкоманда the reader chose from the КОМАНДЫ section and
// this screen does not run itself.
//
// It is handed BACK rather than run here, and that is not a limitation working
// around itself: a screen owns the raw mode of the terminal and the signal
// handler for as long as it is drawn, and two screens drawing at once own them
// twice.  So the status screen closes, hands the terminal back exactly as it
// found it, the caller runs the подкоманда — for analyze that is the walk
// screen, which draws in its turn — and the status screen opens again.  The
// PROCESS never restarts, and the language chosen on the way stays chosen;
// only the drawing changes hands.
//
// The KEYBOARD does not change hands at all: internal/ui reads the terminal
// once for the whole process (see keyboard()), and the screens take keys from
// that one channel by turns.  That is what keeps a keypress from falling
// between two screens — and it is not a supposition: with a reader per screen
// the walk screen opened after this one received nothing at all.
type Request struct {
	// Command is a name out of cli.Commands.  Nothing else travels: a
	// подкоманда that needs a path asks for it where the path is used, in
	// the прочтённая строка ввода of the walk screen, and not twice.
	Command string
}

// Run draws the live screen until the reader asks to leave or asks for a
// подкоманда.  The terminal is handed back the way it was found on every path
// out, including a signal.  A non-nil Request means the caller should run that
// подкоманда and then call Run again.
func Run(o Options) (*Request, error) {
	if o.Out == nil {
		o.Out = os.Stdout
	}
	if o.Interval <= 0 {
		o.Interval = 2 * time.Second
	}
	if o.Collect == nil {
		return nil, lang.Errorf("живому экрану не передан сборщик снимка")
	}
	if !Available(o.Out) {
		return nil, ErrNoTerminal
	}

	tty, keys, err := keyboard()
	if err != nil {
		return nil, ErrNoTerminal
	}

	restore, err := Raw(tty)
	if err != nil {
		return nil, ErrNoTerminal
	}

	s := &screen{o: o, l: o.Lang, t: NewTheme(o.Palette), out: bufio.NewWriterSize(o.Out, 1<<16), tty: tty,
		tab: openTab}
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
		case k, ok := <-keys:
			// The keyboard went away — the terminal closed under the
			// screen.  Leaving is the only honest answer: a closed
			// channel hands back a zero key for ever, and a screen
			// answering that in a loop is a program spinning on a
			// terminal nobody is holding any more.
			if !ok {
				return s.req, nil
			}
			if s.handle(k, snaps) {
				return s.req, nil
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
			return nil, nil
		}
	}
}

// openInput finds a terminal to read keys from.  /dev/tty is asked first so
// the screen still answers the keyboard when standard input is taken by
// something else; standard input itself is the fallback.
//
// It is called once per process — see keyboard() — and what it opens is never
// closed: the terminal outlives every screen drawn on it, and a file closed
// while a read is in flight does not wake that read, it only makes the next
// screen wonder where its keys went.
func openInput() (*os.File, error) {
	if f, err := os.Open("/dev/tty"); err == nil {
		return f, nil
	}
	if IsTerminal(os.Stdin) {
		return os.Stdin, nil
	}
	return nil, ErrNoTerminal
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

func push(h []float64, v float64) []float64 { return pushSample(h, v, histLen) }

// handle answers one key.  It reports whether the screen should close.
func (s *screen) handle(k key, snaps chan sample) bool {
	// КОМАНДЫ is a list to choose from, and there the keys that move a
	// cursor and the key that starts a command belong to it alone.  Every
	// other section is a page to read, and there the same keys scroll.
	if s.tab == menuTab {
		if done, answered := s.menuKey(k, snaps); answered {
			return done
		}
	}
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
			// «?» осталась: она вела к списку команд, и ведёт к нему
			// по-прежнему — теперь к разделу, а не к накладке.
			s.tab, s.scroll = menuTab, 0
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
			if n := int(k.r - '0'); n <= len(sections) {
				s.tab, s.scroll = pickTab(s.tab, n, len(sections)), 0
			}
		}
	case keyEsc:
		return true
	case keyCtrlC:
		return true
	case keyRight, keyTab:
		s.tab = nextTab(s.tab, len(sections))
		s.scroll = 0
	case keyLeft, keyShiftTab:
		s.tab = prevTab(s.tab, len(sections))
		s.scroll = 0
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

// menuKey answers a key while the КОМАНДЫ section is open.  It reports whether
// the screen should close and whether the key was this section's business at
// all: a key it does not claim falls through to the keys every section shares.
func (s *screen) menuKey(k key, snaps chan sample) (done, answered bool) {
	switch k.kind {
	case keyUp:
		s.movePick(-1)
		return false, true
	case keyDown:
		s.movePick(1)
		return false, true
	case keyEnter:
		return s.start(snaps), true
	case keyEsc:
		// Esc goes back to the reading the screen opened on rather than
		// out of the program: a key that both leaves a list and quits
		// would quit by surprise, and that is what it used to do here.
		s.tab, s.scroll = openTab, 0
		return false, true
	case keyRune:
		switch k.r {
		case 'j':
			s.movePick(1)
			return false, true
		case 'k':
			s.movePick(-1)
			return false, true
		}
		if k.r >= '1' && k.r <= '9' {
			if n := int(k.r - '0'); n <= len(cli.Commands) {
				s.pick = n - 1
				s.showPick()
			}
			// A digit on this page chooses a command and never a
			// section: two meanings on one key is how a reader ends
			// up somewhere they did not ask for.  ← → and Tab are
			// how this page is left.
			return false, true
		}
	}
	return false, false
}

// movePick walks the cursor along the list, stopping at its ends: a cursor
// that wraps turns «вниз до упора» into «наверх», and the reader who was
// holding the key finds themselves on `status` meaning to be on `history`.
func (s *screen) movePick(by int) {
	s.pick += by
	if s.pick < 0 {
		s.pick = 0
	}
	if s.pick > len(cli.Commands)-1 {
		s.pick = len(cli.Commands) - 1
	}
	s.showPick()
}

// showPick scrolls the section until the chosen line is on the screen.  On a
// short terminal the list runs off the bottom, and a cursor nobody can see is
// a cursor that starts the wrong command.
func (s *screen) showPick() {
	line := commandsHead + s.pick
	h := s.bodyHeight()
	if line < s.scroll {
		s.scroll = line
	}
	if line >= s.scroll+h {
		s.scroll = line - h + 1
	}
}

// start does what Enter on the chosen line means.  What that is comes from
// cli.Command.Start — beside the command, in the one list — and not from a
// switch on names here.
func (s *screen) start(snaps chan sample) bool {
	c := cli.Commands[s.pick]
	switch c.Start {
	case cli.StartHere:
		// `status` IS this screen.  Running it from here is taking a
		// fresh замер and going to the page that shows it.
		s.request(snaps)
		s.tab, s.scroll = openTab, 0
		return false
	case cli.StartRun, cli.StartPath:
		// The screen closes and the caller runs it; see Request.
		s.req = &Request{Command: c.Name}
		return true
	}
	// Everything else names where it lives instead of doing nothing: a key
	// that answers with silence reads as a key that is broken.
	s.said, s.saidAt = lang.Say("не отсюда: %s", lang.Say(c.Instead)), time.Now()
	return false
}

// bodyHeight is how many lines the section itself gets: the screen less the
// two-line head, the rule and the footer.
func (s *screen) bodyHeight() int { return bodyHeightOf(s.rows, 5) }

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
	h := s.bodyHeight()
	s.scroll = clampScroll(body, s.scroll)
	// shown считает СКОЛЬКО строк раздела видно — это число уезжает в
	// подпись «строки N–M из K». Сами строки укладывает layoutSection: она
	// добирает хвост пустыми, чтобы черта и подвал не уехали вверх.
	shown := body[s.scroll:]
	if len(shown) > h {
		shown = shown[:h]
	}
	out = append(out, layoutSection(body, h, s.scroll)...)

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

	// Narrow: the current section alone, with arrows — and КОМАНДЫ beside
	// it whatever the width.
	//
	// КОМАНДЫ IS THE ONE NAME THAT DOES NOT GIVE WAY.  Everything else on
	// this screen is a reading, and a reading the reader can walk to; that
	// section is where the work is started from, and a person opening the
	// tool for the first time has to find it without being told a key.  The
	// full strip above shows it only from about 118 columns, and eighty is
	// the common terminal — so at eighty it stands here, alone with the
	// section being read.
	var r row
	r.plain(" ")
	if s.tab == menuTab {
		r.add(" "+sections[menuTab].title(s.l)+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
		r.add(fmt.Sprintf("  1/%d", len(sections)), func(x string) string { return t.Fg(t.P.Subtle, x) })
		return r.String()
	}
	r.add(" "+sections[menuTab].title(s.l)+" ", func(x string) string { return t.Chip(t.P.AccentSoft, x) })
	r.add(" ‹ ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	r.add(" "+sections[s.tab].title(s.l)+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
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
	// Every width of the hint is its own wording: an English line is not a
	// Russian line with the words swapped, and one that had to be cut to fit
	// would be cut in a different place.
	exit := s.l.T("q выход ")
	hints := []string{
		s.l.T("← → разделы · ↑ ↓ прокрутка · p пауза · r замер · l язык · 1 КОМАНДЫ · q выход "),
		s.l.T("← → разделы · p пауза · r замер · l язык · 1 КОМАНДЫ · q выход "),
		s.l.T("← → · p · r · l язык · 1 КОМАНДЫ · q выход "),
		s.l.T("1 КОМАНДЫ · q выход "),
		exit,
	}
	if s.tab == menuTab {
		hints = []string{
			s.l.T("↑ ↓ и 1…8 выбрать · Enter запустить · ← → разделы · l язык · q выход "),
			s.l.T("↑ ↓ выбрать · Enter запустить · ← → разделы · q выход "),
			s.l.T("↑ ↓ · Enter запустить · q выход "),
			s.l.T("Enter запустить · q выход "),
			exit,
		}
	}

	// The keys come first in the budget, and «the keys» now means one more
	// than the way out: the shortest hint that still names КОМАНДЫ — the
	// second from the end, the last being the bare way out.
	//
	// It used to reserve only «q выход», and the наблюдение that made this
	// change necessary is exactly what that cost: on eighty columns the
	// state of the замер ate the room, the hint fell back to «q выход», and
	// the list of commands went unnamed on the screen it is started from.
	// Состояние замера уступает — оно и раньше уступало; уступает оно
	// теперь на одну строку раньше.
	keep := hints[maxInt(0, len(hints)-2)]
	// Состояние замера — подпись к ЧТЕНИЮ, и на списке команд его нет
	// вовсе: там строка подвала целиком уходит на клавиши, потому что
	// выбирают и запускают ими, а не датчиком «замер 2 с назад».
	if s.haveSt && s.tab != menuTab {
		ago, took, period := s.l.Since(time.Since(s.taken)), s.l.Millis(s.took), s.l.Every(s.o.Interval)
		states := []string{
			s.l.F("  замер %s назад · длился %s · каждые %s", ago, took, period),
			s.l.F("  замер %s назад · каждые %s", ago, period),
			s.l.F("  замер %s назад", ago),
			"",
		}
		budget := s.cols - r.w - runes(keep) - 2
		for _, state := range states {
			if runes(state) <= budget {
				r.add(state, func(x string) string { return t.Fg(t.P.Subtle, x) })
				break
			}
		}
	}
	if more != "" && r.w+runes(more)+runes(keep)+2 <= s.cols {
		r.add(more, func(x string) string { return t.Fg(t.P.Subtle, x) })
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
