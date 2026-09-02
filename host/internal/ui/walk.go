// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/core"
	"digitdisk/internal/lang"
	"digitdisk/internal/scan"
)

// The screen of `digitdisk analyze`, in two halves.
//
// WHILE THE WALK RUNS it shows the walk: entries and bytes going up, the
// directory being read right now, and which children of the root are filling
// with what.  A walk of four million entries takes minutes, and for those
// minutes the walk IS what is happening — the numbers are not decoration on
// top of the work, they are the work.
//
// WHEN IT ENDS the same screen becomes a way to walk what was found: shares as
// bars, into a directory and back out.
//
// Two rules hold the whole file together:
//
//   - Nothing here is a second measurement.  Every number is the walk's own
//     arithmetic, and the result the report prints afterwards is the same
//     result, unchanged.
//   - A number that is not final says so.  Until the walk ends the largest
//     directory is a guess — the next directory read can overturn it — and the
//     list carries the word ПРЕДВАРИТЕЛЬНО until it stops being one.
type walkMode int

const (
	modeWalking walkMode = iota
	modeSettling
	modeBrowsing
)

// ErrWalkStopped is returned when the reader closed the screen before the walk
// finished.  A walk cut short has no result to report, and reporting the part
// of it that did run as if it were the whole would be the one lie this screen
// must not tell.
var ErrWalkStopped = lang.Errorf("обход прерван")

// WalkOptions is everything the walk screen needs.  Walk is handed in rather
// than called for: the screen owns no traversal of its own and shows exactly
// the walk `digitdisk analyze` would have run without it.
type WalkOptions struct {
	Out      *os.File
	Root     string // the tree to walk; empty asks for one
	Palette  Palette
	Interval time.Duration // how often the screen redraws; default below

	// Lang is the language the screen opens in, and Remember stores the one
	// the reader switches to — the same pair `status` takes, for the same
	// reason: one key changes the language of everything on the screen, and
	// the next run speaks it too.
	Lang     lang.Lang
	Remember func(lang.Lang) lang.Phrase

	// Walk runs the traversal and returns the result the report will print.
	//
	// stop is the screen asking to be let go.  It is answered on every
	// entry, and it exists because this screen is no longer the last thing
	// the process does: a reader who presses q on a walk of four million
	// entries goes back to the экран состояния, and a walk nobody is
	// watching any more would go on reading the disk behind it for the
	// minute it had left.  The caller wires it to scan.Prune, and the
	// truncated result is thrown away — ErrWalkStopped says there is none.
	Walk func(root string, watch func(scan.Step), stop func() bool) (scan.Result, error)

	// After is what the screen does the moment the walk finishes.  It is
	// how `clean`, chosen from the экран состояния, arrives here: the walk
	// runs and the план уборки opens by itself, exactly as the «c» key
	// opens it, with the приговор ядра and the число файлов in front of it
	// as always.
	After After

	// The rest is what the screen can do besides look.  Every one of them is
	// the same call the matching subcommand makes, handed in rather than
	// built here: the screen is a place to act from, never a second set of
	// rules.  A nil one is a screen that cannot do that thing, and says so.
	Plan    func(root string, only []string) (*clean.Plan, error)
	Apply   func(p *clean.Plan) (*clean.Journal, error)
	Restore func(box string, dryRun bool) (*clean.Journal, error)
	History func(root string) (*clean.History, error)
	Places  func() (origin string, found []PlaceRow, err error)
}

// After is what the walk screen does the moment its walk finishes.
type After int

const (
	// AfterNothing — ничего: экран ждёт читателя.
	AfterNothing After = iota
	// AfterPlan — открыть план уборки, как это делает клавиша «c».
	AfterPlan
)

// PlaceRow is one known place as the screen shows it.  The справочник itself
// stays in internal/places; the screen is handed rows, not a parser.
type PlaceRow struct{ Class, Name, Path string }

// drawEvery is how often the screen is redrawn while the walk runs.  Four
// times a second is fast enough that the counters читаются as moving and slow
// enough that the drawing is nothing against the walk: see the measurement in
// walk_test.go (TestLiveOrderIsCalmOnARealTree) and the numbers in the report.
const drawEvery = 250 * time.Millisecond

// calmMargin is how much bigger a directory must get before it is allowed to
// climb over the one above it in the live list.  Without it the list reorders
// itself several times a second while two neighbours trade places by a
// kilobyte, and a list that reshuffles under the eye cannot be read.  Five per
// cent is the smallest step that is still visible as a change of order.
const calmMargin = 0.05

type walkDone struct {
	res  scan.Result
	err  error
	tree *wnode
	snap walkSnap
}

// RunWalk draws the walk and then lets the reader walk what it found — and
// act on it.  The terminal is handed back the way it was found on every path
// out.  The second return says whether there is a result to print: a screen
// closed at the path prompt walked nothing, and printing an empty report would
// be a report about nothing.
func RunWalk(o WalkOptions) (scan.Result, bool, error) {
	if o.Out == nil {
		o.Out = os.Stdout
	}
	if o.Interval <= 0 {
		o.Interval = drawEvery
	}
	if o.Walk == nil {
		return scan.Result{}, false, lang.Errorf("живому экрану не передан обход")
	}
	if !Available(o.Out) {
		return scan.Result{}, false, ErrNoTerminal
	}

	tty, keys, err := keyboard()
	if err != nil {
		return scan.Result{}, false, ErrNoTerminal
	}

	restore, err := Raw(tty)
	if err != nil {
		return scan.Result{}, false, ErrNoTerminal
	}

	if !o.Lang.Valid() {
		o.Lang = lang.Default
	}
	w := &walkScreen{o: o, l: o.Lang, t: NewTheme(o.Palette), out: bufio.NewWriterSize(o.Out, 1<<16),
		started: time.Now(), jobs: make(chan jobResult, 1)}
	w.rows, w.cols, _ = Size(o.Out)
	if w.rows < minRows {
		w.rows = 24
	}
	if w.cols < minCols {
		w.cols = 80
	}

	leave := func() {
		fmt.Fprint(o.Out, showCur+altOff)
		restore()
	}
	defer leave()
	// A walk nobody is watching is let go on the way out.  Before this
	// screen could be closed and reopened it did not matter — the process
	// ended a line later; now the reader goes back to the экран состояния,
	// and a walk left running behind it would read the disk for a minute
	// under a screen that says nothing about it.
	defer w.release()

	sig := make(chan os.Signal, 4)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)
	defer signal.Stop(sig)

	fmt.Fprint(o.Out, altOn+hideCur)

	w.after = o.After
	if o.Root != "" {
		w.startWalk(o.Root)
	} else {
		// `analyze` without a path, in a terminal, asks for one instead of
		// refusing: the screen is where a walk is started from now.
		w.mode = modeBrowsing
		w.openPath()
	}

	tick := time.NewTicker(o.Interval)
	defer tick.Stop()
	w.draw()

	for {
		select {
		case sn := <-w.feedOut():
			w.accept(sn)
			w.draw()
		case <-w.settlingCh():
			// The walk is over and the totals are being added up the tree.
			// On a tree of millions of entries that is a second or two, and
			// a screen that says nothing for a second reads as a hang.
			if w.mode == modeWalking {
				w.mode = modeSettling
			}
			w.settling = nil
			w.draw()
		case d := <-w.doneCh():
			w.done = nil
			if d.err != nil {
				// The first walk failing is the command failing: `analyze
				// /нет-такого` must say so and stop.  A later one failing is
				// an answer to a question asked on the screen, and the screen
				// stays open to be asked another.
				if !w.haveRes {
					return d.res, false, d.err
				}
				w.mode = modeBrowsing
				w.tell(w.l.T("ОБХОД НЕ ВЫШЕЛ"), []string{lang.InLang(d.err, w.l)})
				w.draw()
				continue
			}
			w.finish(d)
			w.draw()
		case <-tick.C:
			if w.mode == modeWalking && w.feed != nil {
				// The next snapshot is asked for, not taken: the walk
				// hands one over when it next passes a thousand entries,
				// and until then nothing interrupts it.
				w.feed.want.Store(true)
			}
			w.draw()
		case r := <-w.jobs:
			w.accepted(r)
			w.draw()
		case k, ok := <-keys:
			if !ok {
				// The terminal closed under the screen; see the
				// same case in Run.
				if w.mode == modeWalking || w.mode == modeSettling {
					return scan.Result{}, false, ErrWalkStopped
				}
				return w.res, w.haveRes, nil
			}
			if w.handle(k) {
				if w.mode == modeWalking || w.mode == modeSettling {
					return scan.Result{}, false, ErrWalkStopped
				}
				return w.res, w.haveRes, nil
			}
			w.draw()
		case sg := <-sig:
			if sg == syscall.SIGWINCH {
				if r, c, ok := Size(o.Out); ok {
					w.rows, w.cols = r, c
					w.scroll = 0
				}
				w.draw()
				continue
			}
			if w.mode == modeWalking || w.mode == modeSettling {
				return scan.Result{}, false, ErrWalkStopped
			}
			return w.res, w.haveRes, nil
		}
	}
}

// ownFilesRow is what the files lying directly in a directory are called when
// they stand in the list beside its subdirectories.
func ownFilesRow(l lang.Lang) string { return l.T("· файлы этого каталога") }

// browseRow is one line of a directory: a subdirectory, or the files lying
// directly in it gathered into one.
type browseRow struct {
	name    string
	bytes   int64
	entries int32
	node    *wnode // nil for the row that stands for the files themselves
	marked  bool
	// own marks the row that stands for the files lying directly in the
	// directory.  It is a flag and not the name, because the name is
	// translated and code that recognised a row by its Russian text would
	// stop recognising it the moment the reader switched languages.
	own bool
}

type walkScreen struct {
	o    WalkOptions
	t    Theme
	out  *bufio.Writer
	rows int
	cols int

	// l is the language every line of this screen is written in.
	l lang.Lang
	// said is the one-line answer to the last key that did something
	// outside the screen — storing the language — and when it was said.
	said   lang.Phrase
	saidAt time.Time

	mode     walkMode
	started  time.Time
	snap     walkSnap
	res      scan.Result
	haveRes  bool
	tree     *wnode
	walkRoot string

	// the walk in flight; replaced whole when another one starts
	feed     *walkFeed
	done     chan walkDone
	settling chan struct{}
	// stop is how the walk in flight is let go when nobody is watching it
	// any more — see WalkOptions.Walk.
	stop *atomic.Bool
	// after is WalkOptions.After, spent once: the план уборки opens on the
	// walk the reader asked for and not on every walk they take afterwards.
	after After

	// order is the list as it is SHOWN — see calmMargin.  swaps counts how
	// often it was reordered, which is what the measurement reads.
	order []string
	swaps int
	rate  []float64
	last  walkSnap

	tab    int
	scroll int

	stack   []*wnode // where in the tree the reader stands
	sel     []int
	atNode  *wnode
	curRows []browseRow

	// what the reader has ticked, and what the screen is asking about
	marks map[string]bool
	ask   *ask
	busy  bool
	jobs  chan jobResult

	// read on demand, forgotten when the disk changes under them
	history     *clean.History
	boxSel      int
	placesDir   string
	placesFound []PlaceRow
}

// feedOut, doneCh and settlingCh are the channels of the walk in flight.  A
// nil channel blocks for ever in a select, which is exactly right when there
// is no walk: the screen waits for a key instead.
func (w *walkScreen) feedOut() <-chan walkSnap {
	if w.feed == nil {
		return nil
	}
	return w.feed.out
}

func (w *walkScreen) doneCh() <-chan walkDone { return w.done }

// release tells the walk in flight, if there is one, that nobody is reading it
// any more.
func (w *walkScreen) release() {
	if w.stop != nil {
		w.stop.Store(true)
	}
}

func (w *walkScreen) settlingCh() <-chan struct{} { return w.settling }

// startWalk begins a traversal and throws away everything the last one left:
// a screen showing one tree's numbers under another tree's name would be the
// worst kind of wrong.
func (w *walkScreen) startWalk(root string) {
	// The walk that was running is let go first: two walks reading the disk
	// for one screen is twice the work and half the speed of the one being
	// watched.
	w.release()
	w.o.Root, w.walkRoot = root, root
	w.mode = modeWalking
	w.started = time.Now()
	w.snap, w.last, w.rate, w.order = walkSnap{}, walkSnap{}, nil, nil
	w.tree, w.stack, w.sel = nil, nil, nil
	w.atNode, w.curRows, w.marks = nil, nil, nil
	w.history, w.boxSel = nil, 0
	w.tab, w.scroll = 0, 0

	feed := newWalkFeed(root)
	done := make(chan walkDone, 1)
	settling := make(chan struct{}, 1)
	stop := new(atomic.Bool)
	w.feed, w.done, w.settling, w.stop = feed, done, settling, stop
	walk := w.o.Walk
	go func() {
		res, err := walk(root, feed.step, stop.Load)
		// The snapshot and the settling both belong to the goroutine that
		// owns the collector, so nothing here is ever read from two places.
		settling <- struct{}{}
		snap := feed.snapshot()
		var tree *wnode
		if err == nil {
			tree = feed.settle()
		}
		done <- walkDone{res: res, err: err, tree: tree, snap: snap}
	}()
}

func (w *walkScreen) accept(sn walkSnap) {
	// The history is bytes per second between the snapshots already shown.
	if !w.last.At.IsZero() {
		if d := sn.At.Sub(w.last.At).Seconds(); d > 0 {
			w.rate = push(w.rate, float64(sn.Entries-w.last.Entries)/d)
		}
	}
	w.last = sn
	w.snap = sn
	w.reorder(sn.Tops)
}

// finish turns the walking screen into the browsing one.
func (w *walkScreen) finish(d walkDone) {
	w.snap = d.snap
	w.res = d.res
	w.tree = d.tree
	w.haveRes = true
	w.reorder(d.snap.Tops)
	w.mode = modeBrowsing
	if w.tree != nil {
		w.stack = []*wnode{w.tree}
		w.sel = []int{0}
	}
	// `clean`, chosen on the экран состояния, is a walk and then the plan.
	// Spent once: the reader who walks somewhere else afterwards asked for
	// a walk, not for another plan.
	if w.after == AfterPlan {
		w.after = AfterNothing
		w.proposeClean()
	}
}

// reorder settles the shown order of the live list.  A directory climbs over
// the one above it only when it is calmMargin bigger, so two neighbours within
// a hair of each other stop trading places every quarter second.
func (w *walkScreen) reorder(tops []topRow) { w.reorderBy(tops, calmMargin) }

// reorderBy is reorder with the margin named, so the measurement can run the
// same list through it twice and count what the margin saved.
func (w *walkScreen) reorderBy(tops []topRow, margin float64) {
	size := make(map[string]int64, len(tops))
	for _, t := range tops {
		size[t.Name] = t.Bytes
	}
	kept := w.order[:0]
	for _, n := range w.order {
		if _, ok := size[n]; ok {
			kept = append(kept, n)
		}
	}
	w.order = kept
	known := make(map[string]bool, len(w.order))
	for _, n := range w.order {
		known[n] = true
	}
	for _, t := range tops {
		if !known[t.Name] {
			w.order = append(w.order, t.Name)
		}
	}
	for pass := 0; pass < len(w.order); pass++ {
		moved := false
		for i := len(w.order) - 1; i > 0; i-- {
			hi, lo := size[w.order[i-1]], size[w.order[i]]
			if float64(lo) > float64(hi)*(1+margin) {
				w.order[i-1], w.order[i] = w.order[i], w.order[i-1]
				w.swaps++
				moved = true
			}
		}
		if !moved {
			break
		}
	}
}

// shownTops is the live list in the order the screen is showing it.
func (w *walkScreen) shownTops() []topRow {
	by := make(map[string]topRow, len(w.snap.Tops))
	for _, t := range w.snap.Tops {
		by[t.Name] = t
	}
	out := make([]topRow, 0, len(w.order))
	for _, n := range w.order {
		out = append(out, by[n])
	}
	return out
}

func (w *walkScreen) handle(k key) bool {
	// A question on the screen takes the keyboard whole: half-answering it
	// while the section underneath also moves is how a confirmation gets
	// given by accident.
	if w.ask != nil {
		return w.askKey(k)
	}
	if w.mode != modeBrowsing {
		// While the walk runs there is nothing to walk around in yet, so the
		// keyboard answers only what makes sense during it: scrolling the
		// list, the language, the keys themselves, and leaving.
		switch k.kind {
		case keyEsc, keyCtrlC:
			return true
		case keyDown:
			w.scroll++
		case keyUp:
			w.scroll--
		case keyPgDn:
			w.scroll += w.bodyHeight()
		case keyPgUp:
			w.scroll -= w.bodyHeight()
		case keyRune:
			switch k.r {
			case 'q', 'Q', 'й', 'Й':
				return true
			case 'j':
				w.scroll++
			case 'k':
				w.scroll--
			case 'g':
				w.scroll = 0
			case 'l', 'L', 'д', 'Д':
				w.switchLang()
			case '?':
				w.ask = &ask{kind: overlayKeys, title: w.l.T("КЛАВИШИ И КОМАНДЫ"),
					hint: w.l.T("любая клавиша — закрыть")}
			}
		}
		if w.scroll < 0 {
			w.scroll = 0
		}
		return false
	}

	inTree := walkSections[w.tab].id == treeTitle
	inJournal := walkSections[w.tab].id == journalTitle
	switch k.kind {
	case keyRune:
		switch k.r {
		case 'q', 'Q', 'й', 'Й':
			return true
		case 'j':
			w.move(1)
		case 'k':
			w.move(-1)
		case 'g':
			w.top()
		case 'G':
			w.bottom()
		case ' ':
			if inTree {
				w.toggleMark()
			}
		case 'l', 'L', 'д', 'Д':
			// The language of the whole screen, on the same key `status`
			// gives it.  It is why this screen has no vim «l» for going
			// into a directory: one letter cannot be two things, and the
			// arrows and Enter do that job anyway.
			w.switchLang()
		case '.':
			if inTree {
				w.markHere()
			}
		case 'c', 'C', 'с', 'С':
			w.proposeClean()
		case 'o', 'O', 'щ', 'Щ':
			w.openPath()
		case '?':
			w.ask = &ask{kind: overlayKeys, title: w.l.T("КЛАВИШИ И КОМАНДЫ"),
				hint: w.l.T("любая клавиша — закрыть")}
		}
		if k.r >= '1' && k.r <= '9' {
			if n := int(k.r - '1'); n < len(walkSections) {
				w.tab, w.scroll = n, 0
			}
		}
	case keyEnter:
		if inJournal {
			w.restoreSelected()
		} else {
			w.enter()
		}
	case keyBack:
		w.up()
	case keyEsc, keyCtrlC:
		return true
	case keyTab:
		w.tab = (w.tab + 1) % len(walkSections)
		w.scroll = 0
	case keyShiftTab:
		w.tab = (w.tab + len(walkSections) - 1) % len(walkSections)
		w.scroll = 0
	case keyRight:
		// Inside ДЕРЕВО the arrows walk the tree, because that is what the
		// reader came to the section for; Tab and the digits still change
		// section, and the footer says so.
		if inTree {
			w.enter()
		} else {
			w.tab = (w.tab + 1) % len(walkSections)
			w.scroll = 0
		}
	case keyLeft:
		if inTree {
			w.up()
		} else {
			w.tab = (w.tab + len(walkSections) - 1) % len(walkSections)
			w.scroll = 0
		}
	case keyDown:
		w.move(1)
	case keyUp:
		w.move(-1)
	case keyPgDn:
		w.move(w.bodyHeight())
	case keyPgUp:
		w.move(-w.bodyHeight())
	}
	if w.scroll < 0 {
		w.scroll = 0
	}
	return false
}

// switchLang draws the whole screen in the other language and remembers the
// choice, so the next run — and `digitdisk clean` tomorrow — speaks it too.
func (w *walkScreen) switchLang() {
	w.l = w.l.Other()
	// The rows carry translated names, and the sections that read from disk
	// carry translated labels; both are built again in the new language.
	w.atNode = nil
	if w.o.Remember != nil {
		w.said, w.saidAt = w.o.Remember(w.l), time.Now()
	}
}

// restoreSelected asks about the корзина the cursor is on.
func (w *walkScreen) restoreSelected() {
	if w.history == nil || w.boxSel < 0 || w.boxSel >= len(w.history.Entries) {
		return
	}
	w.proposeRestore(w.history.Entries[w.boxSel].Box)
}

// move steps the cursor in the two sections that have one and scrolls
// anywhere else.
func (w *walkScreen) move(by int) {
	if walkSections[w.tab].id == journalTitle && w.history != nil && len(w.history.Entries) > 0 {
		i := w.boxSel + by
		if i >= len(w.history.Entries) {
			i = len(w.history.Entries) - 1
		}
		if i < 0 {
			i = 0
		}
		w.boxSel = i
		return
	}
	if walkSections[w.tab].id != treeTitle || len(w.stack) == 0 {
		w.scroll += by
		if w.scroll < 0 {
			w.scroll = 0
		}
		return
	}
	rows := w.here()
	i := w.sel[len(w.sel)-1] + by
	if i >= len(rows) {
		i = len(rows) - 1
	}
	if i < 0 {
		i = 0
	}
	w.sel[len(w.sel)-1] = i
}

func (w *walkScreen) top() {
	if walkSections[w.tab].id == treeTitle && len(w.sel) > 0 {
		w.sel[len(w.sel)-1] = 0
		return
	}
	w.scroll = 0
}

func (w *walkScreen) bottom() {
	if walkSections[w.tab].id == treeTitle && len(w.sel) > 0 {
		w.sel[len(w.sel)-1] = len(w.here()) - 1
		if w.sel[len(w.sel)-1] < 0 {
			w.sel[len(w.sel)-1] = 0
		}
	}
}

func (w *walkScreen) enter() {
	if len(w.stack) == 0 {
		return
	}
	rows := w.here()
	i := w.sel[len(w.sel)-1]
	if i < 0 || i >= len(rows) || rows[i].node == nil {
		return
	}
	w.stack = append(w.stack, rows[i].node)
	w.sel = append(w.sel, 0)
	w.scroll = 0
}

func (w *walkScreen) up() {
	if len(w.stack) > 1 {
		w.stack = w.stack[:len(w.stack)-1]
		w.sel = w.sel[:len(w.sel)-1]
		w.scroll = 0
	}
}

// here is the rows of the directory the reader stands in, built once per
// directory rather than once per redraw.
func (w *walkScreen) here() []browseRow {
	if len(w.stack) == 0 {
		return nil
	}
	n := w.stack[len(w.stack)-1]
	if w.atNode == n {
		return w.curRows
	}
	rows := make([]browseRow, 0, len(n.kids)+1)
	for _, k := range n.children() {
		rows = append(rows, browseRow{name: k.name + "/", bytes: k.bytes,
			entries: k.entries, node: k, marked: w.marks[k.path()]})
	}
	if n.ownFiles > 0 || n.own > 0 {
		rows = append(rows, browseRow{name: ownFilesRow(w.l), bytes: n.own, entries: n.ownFiles, own: true})
	}
	// The rows are already sorted by size except for the one that stands for
	// the files; it is put where its size says it belongs.
	for i := len(rows) - 1; i > 0; i-- {
		if rows[i].bytes > rows[i-1].bytes {
			rows[i], rows[i-1] = rows[i-1], rows[i]
			continue
		}
		break
	}
	w.atNode, w.curRows = n, rows
	return rows
}

func (w *walkScreen) bodyHeight() int {
	h := w.rows - 5
	if h < 1 {
		h = 1
	}
	return h
}

func (w *walkScreen) draw() {
	var b strings.Builder
	b.WriteString(home)
	b.WriteString(w.t.Canvas())
	lines := w.frame()
	for i := 0; i < w.rows; i++ {
		b.WriteString(w.t.Canvas())
		if i < len(lines) {
			b.WriteString(w.t.clip(lines[i], w.cols))
		}
		b.WriteString(clearLine)
		if i < w.rows-1 {
			b.WriteString("\r\n")
		}
	}
	b.WriteString(clearDown)
	w.out.WriteString(b.String())
	w.out.Flush()
}

func (w *walkScreen) frame() []string {
	t := w.t
	out := make([]string, 0, w.rows)
	out = append(out, w.header())
	out = append(out, w.strip())
	out = append(out, t.Fg(t.P.Border, strings.Repeat("─", w.cols)))

	var body []string
	switch {
	case w.ask != nil && w.ask.kind == overlayKeys:
		body = w.keysLines()
	case w.ask != nil:
		body = w.askLines()
	case w.mode == modeBrowsing:
		body = walkSections[w.tab].render(w)
	default:
		body = w.walking()
	}
	h := w.bodyHeight()
	// In the two sections that have a cursor the scroll follows it;
	// everywhere else the reader scrolls it.
	if w.ask != nil {
		w.scroll = 0
	} else if w.mode == modeBrowsing {
		switch walkSections[w.tab].id {
		case treeTitle:
			w.scroll = follow(w.scroll, w.cursorLine(), h, len(body))
		case journalTitle:
			w.scroll = follow(w.scroll, journalHead+w.boxSel, h, len(body))
		}
	}
	if w.scroll > len(body)-1 {
		w.scroll = len(body) - 1
	}
	if w.scroll < 0 {
		w.scroll = 0
	}
	shown := body[w.scroll:]
	if len(shown) > h {
		shown = shown[:h]
	}
	out = append(out, shown...)
	for i := len(shown); i < h; i++ {
		out = append(out, "")
	}

	more := ""
	if len(body) > h {
		more = w.l.F("  строки %d–%d из %d", w.scroll+1, w.scroll+len(shown), len(body))
	}
	out = append(out, t.Fg(t.P.Border, strings.Repeat("─", w.cols)))
	out = append(out, w.footer(more))
	return out
}

// follow keeps the cursor line inside the window without moving the window
// more than it has to.
func follow(scroll, line, h, n int) int {
	if line < 0 {
		return scroll
	}
	if line < scroll {
		return line
	}
	if line >= scroll+h {
		return line - h + 1
	}
	if scroll > n-1 {
		return maxInt(0, n-1)
	}
	return scroll
}

// cursorLine is which line of the ДЕРЕВО body the cursor sits on.
func (w *walkScreen) cursorLine() int {
	if len(w.sel) == 0 {
		return -1
	}
	return treeHead + w.sel[len(w.sel)-1]
}

// treeHead is how many lines the ДЕРЕВО section draws before its list.
const treeHead = 4

func (w *walkScreen) header() string {
	t := w.t
	clock := time.Now().Format("15:04:05")
	what := w.l.T("обход ")
	if w.mode == modeBrowsing {
		what = w.l.T("обойдено ")
	}
	var r row
	r.add(" ◇ digitdisk ", func(x string) string { return t.Chip(t.P.Accent, x) })
	r.plain(" ")
	r.add(what, func(x string) string { return t.Fg(t.P.Subtle, x) })

	// The brand and the clock give way before the path does, and they give way
	// whole: a path cut down to one ellipsis says nothing, and a clock cut in
	// half is not a clock.  On a forty-column terminal only the path remains.
	head := r.w
	for _, tailer := range []string{"Digitable  " + clock + " ", clock + " ", ""} {
		room := w.cols - head - runes(tailer)
		if tailer != "" && room < 12 {
			continue
		}
		r.add(tail(w.headRoot(), maxInt(1, room-1)), func(x string) string { return t.Fg(t.P.Muted, x) })
		if tailer == "" {
			break
		}
		r.pad(w.cols - runes(tailer))
		if strings.HasPrefix(tailer, "Digitable") {
			r.add("Digitable  ", func(x string) string { return t.Fg(t.P.AccentSoft, x) })
			r.add(clock+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		} else {
			r.add(tailer, func(x string) string { return t.Fg(t.P.Subtle, x) })
		}
		break
	}
	return r.String()
}

// headRoot is what the head names: the tree being walked, or the fact that
// none has been chosen yet.
func (w *walkScreen) headRoot() string {
	if w.walkRoot == "" {
		return w.l.T("каталог не выбран")
	}
	return w.walkRoot
}

func firstNonBlank(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// strip is the tab bar once there is something to walk, and the state of the
// walk while it runs.
func (w *walkScreen) strip() string {
	t := w.t
	if w.mode != modeBrowsing {
		var r row
		face := spinner[int(time.Since(w.started)/(120*time.Millisecond))%len(spinner)]
		word := w.l.T(" ИДЁТ ОБХОД ")
		if w.mode == modeSettling {
			word, face = w.l.T(" СВОДИТСЯ ДЕРЕВО "), "·"
		}
		r.add(" "+face+" ", func(x string) string { return t.Fg(t.P.Accent, x) })
		r.add(word, func(x string) string { return t.Chip(t.P.AccentSoft, x) })
		r.add("  "+elapsed(w.l, time.Since(w.started)), func(x string) string {
			return t.Fg(t.P.Subtle, x)
		})
		return r.String()
	}
	for _, gap := range []int{2, 1} {
		width := 1
		for _, sec := range walkSections {
			width += runes(sec.title(w.l)) + 2 + gap
		}
		if width+1 > w.cols {
			continue
		}
		var r row
		r.plain(" ")
		for i, sec := range walkSections {
			if i == w.tab {
				r.add(" "+sec.title(w.l)+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
			} else {
				r.add(" "+sec.title(w.l)+" ", func(x string) string { return t.Fg(t.P.Subtle, x) })
			}
			r.plain(strings.Repeat(" ", gap))
		}
		return r.String()
	}
	cur := walkSections[w.tab]
	var r row
	r.add(" ‹ ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	r.add(" "+cur.title(w.l)+" ", func(x string) string { return t.Chip(t.P.Accent, x) })
	r.add(fmt.Sprintf(" ›  %d/%d", w.tab+1, len(walkSections)), func(x string) string {
		return t.Fg(t.P.Subtle, x)
	})
	return r.String()
}

func (w *walkScreen) footer(more string) string {
	t := w.t
	var r row
	if w.mode == modeBrowsing {
		r.add(w.l.T(" ГОТОВО "), func(x string) string { return t.Chip(t.P.Green, x) })
	} else {
		r.add(w.l.T(" ОБХОД "), func(x string) string { return t.Chip(t.P.AccentSoft, x) })
	}

	exit := w.l.T("q выход ")
	state := w.l.F("  %s записей · %s", w.l.Num(w.snap.Entries), w.l.Bytes(w.snap.Bytes))
	if runes(state) <= w.cols-r.w-runes(exit)-2 {
		r.add(state, func(x string) string { return t.Fg(t.P.Subtle, x) })
	}
	if !w.said.Empty() && time.Since(w.saidAt) < 6*time.Second {
		more = "  " + w.said.In(w.l)
	}
	if more != "" && r.w+runes(more)+runes(exit)+2 <= w.cols {
		r.add(more, func(x string) string { return t.Fg(t.P.Subtle, x) })
	}

	// Every width is its own line: the same thing said shorter is another
	// sentence, not the first one cut off, and it is translated as one.
	hints := []string{
		w.l.T("Tab разделы · ↑ ↓ строка · → внутрь · Пробел отметить · c убрать · l язык · ? клавиши · q выход "),
		w.l.T("Tab · ↑ ↓ · → внутрь · Пробел отметить · c убрать · l язык · ? · q выход "),
		w.l.T("Tab · ↑ ↓ · → ← · c убрать · l · ? · q выход "),
		w.l.T("? клавиши · q выход "),
		exit,
	}
	if w.ask != nil {
		hints = []string{w.l.T("Esc — отменить, q — выход "), w.l.T("Esc · q ")}
		if w.ask.kind == overlayNote || w.ask.kind == overlayKeys {
			hints = []string{w.l.T("любая клавиша — закрыть "), w.l.T("закрыть ")}
		}
	}
	if w.mode != modeBrowsing {
		hints = []string{w.l.T("↑ ↓ прокрутка · l язык · q прервать обход "),
			w.l.T("↑ ↓ · l · q прервать "), w.l.T("q прервать "), "q "}
	}
	for _, hint := range hints {
		if r.w+runes(hint)+2 <= w.cols {
			r.pad(w.cols - runes(hint))
			r.add(hint, func(x string) string { return t.Fg(t.P.Muted, x) })
			break
		}
	}
	return r.String()
}

// walking is the page the screen shows while the walk runs.
func (w *walkScreen) walking() []string {
	t, sn := w.t, w.snap
	out := []string{""}
	out = append(out, w.kv(w.l.T("записей"), w.l.Num(sn.Entries),
		w.l.F("файлов %s · каталогов %s · ссылок %s",
			w.l.Num(sn.Files), w.l.Num(sn.Dirs), w.l.Num(sn.Links))))
	out = append(out, w.kv(w.l.T("объём"), w.l.Bytes(sn.Bytes), ""))
	secs := time.Since(w.started).Seconds()
	if secs > 0.5 && sn.Entries > 0 {
		out = append(out, w.kv(w.l.T("скорость"), w.l.F("%s записей/с · %s/с",
			w.l.Num(int64(float64(sn.Entries)/secs)), w.l.Bytes(int64(float64(sn.Bytes)/secs))), ""))
	}
	if sn.Cur != "" {
		out = append(out, w.kv(w.l.T("сейчас"), tail(sn.Cur, maxInt(10, w.cols-38)),
			w.l.F("глубина %d", sn.Depth)))
	}
	out = append(out, w.sparkRow(w.l.T("записей/с"), w.rate))
	out = append(out, "")

	var head row
	head.add(w.l.T("  ЧЕМ НАПОЛНЕНО  "), func(x string) string { return t.Bold(t.P.AccentSoft, x) })
	head.add(w.l.T(" ПРЕДВАРИТЕЛЬНО "), func(x string) string { return t.Chip(t.P.Yellow, x) })
	out = append(out, head.String())
	out = append(out, w.list(w.shownTops(), -1)...)
	if sn.Truncated {
		out = append(out, "", w.note(w.l.F(
			"каталогов больше %s — дальше считается всё, но ходить можно не везде", w.l.Num(maxNodes))))
	}
	return out
}

// list draws the shares of one level: the live list while the walk runs, and
// the top of the tree once it has.  cursor is the row to mark, or -1.
func (w *walkScreen) list(tops []topRow, cursor int) []string {
	rows := make([]browseRow, 0, len(tops))
	for _, t := range tops {
		name := t.Name
		if t.Own {
			name = ownFilesRow(w.l)
		}
		rows = append(rows, browseRow{name: name, bytes: t.Bytes, entries: t.Entries, own: t.Own})
	}
	return w.rowLines(rows, cursor)
}

// rowLines is the whole picture of one directory: index, bar, share, name,
// size and how many entries lie under it.  Columns are dropped from the right
// as the terminal narrows — the name and the size are the two that never go.
func (w *walkScreen) rowLines(rows []browseRow, cursor int) []string {
	t := w.t
	var total int64
	for _, r := range rows {
		total += r.bytes
	}
	const (
		lead = 4  // marker
		num  = 4  // "12. "
		size = 10 // a rendered byte count
		cnt  = 15 // entries under it: seven digits, a space and «зап.»
		pct  = 7  // "100.0%"
	)
	bar, showPct, showCnt := 0, false, false
	switch {
	case w.cols >= 100:
		bar, showPct, showCnt = minInt(28, w.cols/3), true, true
	case w.cols >= 76:
		bar, showPct = minInt(20, w.cols/3), true
	case w.cols >= 56:
		showPct = true
	}
	name := w.cols - lead - num - size - 2
	if bar > 0 {
		name -= bar + 2
	}
	if showPct {
		name -= pct
	}
	if showCnt {
		name -= cnt
	}
	if name < 8 {
		name = 8
	}

	out := make([]string, 0, len(rows))
	for i, r := range rows {
		frac := 0.0
		if total > 0 {
			frac = float64(r.bytes) / float64(total)
		}
		var line row
		// One four-cell column carries both: where the cursor is, and what
		// has been ticked.  Two columns would cost the name two cells on a
		// forty-column terminal for nothing.
		switch {
		case i == cursor && r.marked:
			line.add(" ▶✓ ", func(x string) string { return t.Fg(t.P.Green, x) })
		case i == cursor:
			line.add(" ▶  ", func(x string) string { return t.Fg(t.P.Accent, x) })
		case r.marked:
			line.add("  ✓ ", func(x string) string { return t.Fg(t.P.Green, x) })
		default:
			line.plain("    ")
		}
		line.add(right(fmt.Sprint(i+1), 2)+". ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		if bar > 0 {
			line.plain(t.shareBar(frac, bar))
			line.w += bar
			line.plain("  ")
		}
		if showPct {
			line.add(right(w.l.Pct(frac*100, 1), pct), func(x string) string { return t.Fg(t.P.Foreground, x) })
		}
		// A directory is named in the reading colour; the row that stands
		// for the files themselves is dimmer, because it is not a place the
		// reader can walk into.
		paint := func(x string) string { return t.Fg(t.P.Foreground, x) }
		if r.own {
			paint = func(x string) string { return t.Fg(t.P.Muted, x) }
		}
		line.add("  "+fit(r.name, name), paint)
		line.add(right(w.l.Bytes(r.bytes), size), func(x string) string { return t.Fg(t.P.Foreground, x) })
		if showCnt {
			line.add(right(w.l.F("%s зап.", w.l.Num(int64(r.entries))), cnt), func(x string) string {
				return t.Fg(t.P.Subtle, x)
			})
		}
		out = append(out, line.String())
	}
	if len(out) == 0 {
		out = append(out, w.note(w.l.T("пусто")))
	}
	return out
}

// shareBar is bar without the traffic-light colouring: a share of a disk is a
// warning at ninety per cent, a share of a directory is not.
//
// How many cells the share fills is asked of barCells and not worked out here:
// that arithmetic is one of the layout decisions the tree keeps in one place,
// so that a build with the `flangui` tag hands it to the printed flang library
// together with every other bar on the screen.  A copy of the rounding would
// go on being Go while the rest of the screen stopped being it.
func (t Theme) shareBar(frac float64, n int) string {
	if n <= 0 {
		return ""
	}
	full := barCells(frac, n)
	var r row
	r.add(strings.Repeat("█", full), func(s string) string { return t.Fg(t.P.Accent, s) })
	r.add(strings.Repeat("─", n-full), func(s string) string { return t.Fg(t.P.Border, s) })
	return r.String()
}

func (w *walkScreen) kv(label, value, extra string) string {
	t := w.t
	var r row
	// Fourteen cells is the longest label the walk has («предел глубины»,
	// «жёсткие ссылки»): a label cut in half reads as a different word.
	r.add("  "+fit(label, 14)+"  ", func(x string) string { return t.Fg(t.P.Muted, x) })
	r.add(value, func(x string) string { return t.Fg(t.P.Foreground, x) })
	if extra != "" {
		r.add("   "+extra, func(x string) string { return t.Fg(t.P.Subtle, x) })
	}
	return r.String()
}

func (w *walkScreen) note(text string) string {
	return w.t.Fg(w.t.P.Subtle, "  "+fit(text, maxInt(1, w.cols-3)))
}

func (w *walkScreen) caption(text string) string {
	return w.t.Bold(w.t.P.AccentSoft, "  "+text)
}

func (w *walkScreen) sparkRow(label string, history []float64) string {
	t := w.t
	var r row
	r.add("  "+fit(label, 14)+"  ", func(x string) string { return t.Fg(t.P.Subtle, x) })
	n := minInt(40, maxInt(8, w.cols/3))
	// The spark is drawn against the fastest moment so far, so the shape is
	// the shape of the walk and not of an invented ceiling.
	top := 0.0
	for _, v := range history {
		if v > top {
			top = v
		}
	}
	scaled := make([]float64, len(history))
	for i, v := range history {
		if top > 0 {
			scaled[i] = v / top
		}
	}
	r.plain(t.spark(scaled, n))
	r.w += n
	if top > 0 {
		r.add(w.l.F("  пик %s/с", w.l.Num(int64(top))), func(x string) string { return t.Fg(t.P.Subtle, x) })
	}
	return r.String()
}

// elapsed is how long the walk has been running, in the units that suit it.
func elapsed(l lang.Lang, d time.Duration) string {
	s := int(d.Seconds())
	if s < 60 {
		return l.F("%d с", s)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// tail keeps the END of a path when it does not fit: the last components are
// what says where the walk is, the first are the root it was given.
func tail(p string, n int) string {
	rs := []rune(p)
	if len(rs) <= n {
		return p
	}
	if n <= 1 {
		return "…"
	}
	return "…" + string(rs[len(rs)-n+1:])
}

// --- the sections after the walk -------------------------------------------

// The two sections the keyboard treats specially are named by an identifier
// and not by their heading: the heading is translated, and code that compared
// it against a Russian string would stop recognising its own section the
// moment the reader pressed «l».
const (
	treeTitle    = "tree"
	journalTitle = "journal"
)

type walkSection struct {
	// id is what the keyboard knows this section by, and title is what the
	// reader sees.  title is a call rather than a string for the reason the
	// status screen gives: written as l.T("ИТОГ") the wording stands where
	// it is looked up, and it is translated at the moment it is drawn, so
	// the language key switches the strip with everything else.
	id     string
	title  func(lang.Lang) string
	render func(*walkScreen) []string
}

var walkSections = []walkSection{
	{"total", func(l lang.Lang) string { return l.T("ИТОГ") }, (*walkScreen).total},
	{treeTitle, func(l lang.Lang) string { return l.T("ДЕРЕВО") }, (*walkScreen).browse},
	{"biggest", func(l lang.Lang) string { return l.T("КРУПНЕЙШЕЕ") }, (*walkScreen).biggest},
	{"removable", func(l lang.Lang) string { return l.T("МОЖНО УБРАТЬ") }, (*walkScreen).removable},
	{"classes", func(l lang.Lang) string { return l.T("РАЗРЯДЫ") }, (*walkScreen).classes},
	{"skipped", func(l lang.Lang) string { return l.T("ПРОПУЩЕНО") }, (*walkScreen).skipped},
	{"places", func(l lang.Lang) string { return l.T("МЕСТА") }, (*walkScreen).placesSection},
	{journalTitle, func(l lang.Lang) string { return l.T("ЖУРНАЛ") }, (*walkScreen).journalSection},
}

func (w *walkScreen) total() []string {
	r := w.res
	out := []string{""}
	out = append(out, w.kv(w.l.T("записей"), w.l.Num(int64(r.Entries)),
		w.l.F("файлов %s · каталогов %s · ссылок %s · прочего %s",
			w.l.Num(int64(r.Files)), w.l.Num(int64(r.Dirs)),
			w.l.Num(int64(r.Links)), w.l.Num(int64(r.Others)))))
	out = append(out, w.kv(w.l.T("объём"), w.l.Bytes(r.TotalBytes),
		w.l.T("видимый размер, как у du --apparent-size")))
	out = append(out, w.kv(w.l.T("время"), w.l.F("%s с", w.l.Dec(r.DurationSeconds, 2)), ""))
	if r.HardlinkDupes > 0 {
		out = append(out, w.kv(w.l.T("жёсткие ссылки"), w.l.F("%s повторных имён не засчитано (%s)",
			w.l.Num(int64(r.HardlinkDupes)), w.l.Bytes(r.HardlinkBytes)), ""))
	}
	if s := r.Skipped; s.Total() > 0 {
		out = append(out, w.kv(w.l.T("пропущено"), w.l.Num(int64(s.Total())), w.l.T("раздел ПРОПУЩЕНО")))
	}
	if !r.DeciderReady {
		out = append(out, "", w.note(w.l.T("настоящий разбор не выполнялся: разряды показывают работу стыковки, не анализ")))
	}
	out = append(out, "", w.caption(w.l.T("ЧЕМ НАПОЛНЕНО")))
	out = append(out, w.list(w.finalTops(), -1)...)
	if w.snap.Truncated {
		out = append(out, "", w.note(w.l.F(
			"каталогов больше %s — посчитано всё, но ходить можно не по всему дереву", w.l.Num(maxNodes))))
	}
	return out
}

// finalTops is the top of the tree as the walk left it.
func (w *walkScreen) finalTops() []topRow {
	if w.tree == nil {
		return nil
	}
	rows := make([]topRow, 0, len(w.tree.kids)+1)
	for _, k := range w.tree.children() {
		rows = append(rows, topRow{Name: k.name + "/", Bytes: k.bytes, Entries: k.entries})
	}
	if w.tree.own > 0 {
		rows = append(rows, topRow{Own: true, Bytes: w.tree.own, Entries: w.tree.ownFiles})
	}
	for i := len(rows) - 1; i > 0; i-- {
		if rows[i].Bytes > rows[i-1].Bytes {
			rows[i], rows[i-1] = rows[i-1], rows[i]
			continue
		}
		break
	}
	return rows
}

func (w *walkScreen) browse() []string {
	if len(w.stack) == 0 {
		return []string{"", w.note(w.l.T("дерево не собрано"))}
	}
	n := w.stack[len(w.stack)-1]
	rows := w.here()
	out := []string{""}
	out = append(out, w.kv(w.l.T("каталог"), tail(n.path(), maxInt(10, w.cols-22)), ""))
	marked := w.l.T("Пробел — отметить каталог, «.» — этот, «c» — план уборки")
	if len(w.marks) > 0 {
		marked = w.l.F("отмечено каталогов: %d — «c» покажет план уборки", len(w.marks))
	}
	out = append(out, w.kv(w.l.T("в нём"), w.l.F("%s · %s записей",
		w.l.Bytes(n.bytes), w.l.Num(int64(n.entries))),
		w.l.F("подкаталогов %s", w.l.Num(int64(len(n.kids))))))
	out = append(out, w.note(marked))
	// treeHead lines have been drawn: the cursor arithmetic depends on it.
	sel := w.sel[len(w.sel)-1]
	out = append(out, w.rowLines(rows, sel)...)
	if n.topName != "" {
		out = append(out, "", w.note(w.l.F("крупнейший файл прямо здесь: %s (%s)",
			n.topName, w.l.Bytes(n.topBytes))))
	}
	return out
}

func (w *walkScreen) biggest() []string {
	if len(w.res.Largest) == 0 {
		return []string{"", w.note(w.l.T("нечего"))}
	}
	out := []string{"", w.caption(w.l.F("КРУПНЕЙШИЕ ЗАПИСИ — %d", len(w.res.Largest))), ""}
	return append(out, w.entryLines(w.res.Largest, false)...)
}

func (w *walkScreen) removable() []string {
	rem := w.res.ByVerdict[core.VerdictRemovable]
	out := []string{"", w.kv(w.l.T("предложено"), w.l.F("%s записей · %s",
		w.l.Num(int64(rem.Count)), w.l.Bytes(rem.Bytes)), "")}
	// The screen counts and nothing else, exactly as `analyze` does.  What
	// removes is `clean`, and it does not remove at once either: the command
	// is spelled out here so the next step is a decision and not a keypress.
	out = append(out, "", w.note(w.l.T("убрать можно отсюда: «c» строит план и спрашивает число файлов.")))
	out = append(out, w.note(w.l.T("отметьте каталоги в разделе ДЕРЕВО, чтобы взять только их; без отметок — всё дерево.")))
	out = append(out, "", w.note(w.l.T("та же уборка одной командой, если экран не нужен:")))
	out = append(out, w.t.Fg(w.t.P.AccentSoft, "    digitdisk clean "+w.walkRoot))
	out = append(out, w.note(w.l.T("без --apply она печатает план и не трогает ни одного файла")))
	if len(w.res.Removable) == 0 {
		return append(out, "", w.note(w.l.T("нечего")))
	}
	out = append(out, "")
	return append(out, w.entryLines(w.res.Removable, true)...)
}

// entryLines draws the walk's own rankings — the same rows the report prints.
func (w *walkScreen) entryLines(items []scan.Entry, class bool) []string {
	t := w.t
	kindWidth := 10
	if class {
		kindWidth = 13
	}
	pathWidth := maxInt(12, w.cols-kindWidth-28)
	out := make([]string, 0, len(items))
	for _, e := range items {
		var r row
		r.add("  "+right(w.l.Bytes(e.Size), 10)+"  ", func(x string) string { return t.Fg(t.P.Foreground, x) })
		// The разряд and the вид are identifiers of the решающий слой and
		// of the обход; Word gives the WORD the screen shows for them and
		// leaves the identifier alone.
		kind := w.l.Word(string(e.Kind))
		if class {
			kind = w.l.Word(string(e.Class))
		}
		r.add(fit(kind, kindWidth)+" ", func(x string) string { return t.Fg(t.P.Purple, x) })
		r.add(right(w.l.Days(e.AgeDays), 7)+"  ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		r.add(fit(tail(e.Path, pathWidth), pathWidth), func(x string) string {
			return t.Fg(t.P.Muted, strings.TrimRight(x, " "))
		})
		out = append(out, r.String())
	}
	return out
}

func (w *walkScreen) classes() []string {
	t := w.t
	out := []string{"", w.caption(w.l.T("ПО РАЗРЯДАМ"))}
	line := func(name string, b scan.Bucket, total int64) string {
		frac := 0.0
		if total > 0 {
			frac = float64(b.Bytes) / float64(total)
		}
		var r row
		r.add("  "+fit(name, 15), func(x string) string { return t.Fg(t.P.Foreground, x) })
		r.add(right(w.l.Num(int64(b.Count)), 10)+"  ", func(x string) string { return t.Fg(t.P.Subtle, x) })
		r.add(right(w.l.Bytes(b.Bytes), 11)+"  ", func(x string) string { return t.Fg(t.P.Muted, x) })
		if bw := minInt(24, maxInt(0, w.cols-48)); bw > 4 {
			r.plain(t.shareBar(frac, bw))
			r.w += bw
		}
		return r.String()
	}
	for _, c := range core.Classes {
		out = append(out, line(w.l.Word(string(c)), w.res.ByClass[c], w.res.TotalBytes))
	}
	out = append(out, "", w.caption(w.l.T("ПО ПРИГОВОРАМ")))
	for _, v := range core.Verdicts {
		out = append(out, line(w.l.Word(string(v)), w.res.ByVerdict[v], w.res.TotalBytes))
	}
	out = append(out, "", w.note(w.l.F("решающий слой %s, договор версии %d",
		w.res.Decider, w.res.ContractVersion)))
	return out
}

func (w *walkScreen) skipped() []string {
	s := w.res.Skipped
	out := []string{""}
	out = append(out, w.kv(w.l.T("всего"), w.l.Num(int64(s.Total())), ""))
	out = append(out, w.kv(w.l.T("нет доступа"), w.l.Num(int64(s.PermissionDenied)), ""))
	out = append(out, w.kv(w.l.T("исчезло"), w.l.Num(int64(s.Vanished)), ""))
	out = append(out, w.kv(w.l.T("иные ошибки"), w.l.Num(int64(s.OtherErrors)), ""))
	out = append(out, w.kv(w.l.T("граница ФС"), w.l.Num(int64(s.DeviceBoundaries)), ""))
	out = append(out, w.kv(w.l.T("предел глубины"), w.l.Num(int64(s.DepthLimited)), ""))
	if len(s.Examples) > 0 {
		out = append(out, "", w.caption(w.l.T("ПРИМЕРЫ")))
		for _, e := range s.Examples {
			out = append(out, w.note(e))
		}
	}
	return out
}
