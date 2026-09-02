// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"digitdisk/internal/clean"
	"digitdisk/internal/cli"
	"digitdisk/internal/lang"
)

// The screen is a place to work from, not a display case.  A walk starts here,
// what it found is walked here, and what the decision layer marked
// «МожноУбрать» is removed from here.
//
// REMOVAL FROM THE SCREEN GOES DOWN THE SAME ROAD AS `clean`, ALL OF IT.
// There is no second road, and this file adds no rule of its own:
//
//   - What may go is what clean.Make put in a plan, which is what the decision
//     layer gave the verdict «МожноУбрать» and the host's own guard let past.
//     Marking a directory narrows the ground the plan is made on; it can never
//     add a path the layer did not mark.  A «Спросить» or «НеТрогать» item
//     inside a marked directory is not in the plan and cannot be put there by
//     marking harder.
//   - Nothing moves before the plan has been shown: how many files, how many
//     bytes, broken down by разряд, into which корзина.
//   - Nothing moves before the exact number of files has been typed, the way
//     `purge --confirm N` asks for it.  A confirmation that can be given
//     without looking confirms nothing, so it is not one keystroke.
//   - The move is rename(2) into a корзина inside the корень.  It frees
//     nothing, and the line that says so stands next to the number of bytes,
//     not in a footnote.
//   - Putting it back is here too, and it asks in the same way.
//
// ERASING FOR GOOD IS ALSO HERE, on забой, and it goes down THAT SAME ROAD.
// The корзина is what it does not use; the plan is not what it does not use.
// Every line above holds for it word for word — what may go is what clean.Make
// put in a plan, marking narrows and never widens, nothing happens before the
// plan is on the screen — and it adds exactly two things of its own:
//
//   - the plan says «стереть насовсем» and «корзины не будет» where the other
//     one says «перенести» and names a корзина, because a screen that shows the
//     same words for two different fates is a screen that lies once;
//   - HOW HARD IT ASKS DEPENDS ON HOW MUCH IS GOING, and on nothing else.  A
//     handful of small files, all of them listed in front of the reader: one
//     key, «y».  Anything above that: the exact number of files, typed, the way
//     `purge --confirm N` asks.  Where "above that" is, is not this screen's
//     opinion — see eraseAtOnce below.
//
// `purge` itself is still not here, and that has not changed: it empties a
// корзина, the screen has no корзина open, and the journal section names the
// command for the reader who wants one.

// overlay is what the screen is asking about, drawn over the section body.
type overlay int

const (
	overlayNone  overlay = iota
	overlayPath          // where to walk
	overlayPlan          // what would be moved, and the number to type
	overlayErase         // what would be erased for good, and how it is confirmed
	overlayBack          // what would be put back, and the number to type
	overlayNote          // one answer to read and dismiss
	overlayKeys          // what the keys do
)

// How much may go on one keystroke.
//
// The question does not scroll: frame() zeroes w.scroll while a question is up,
// so whatever does not fit under the title is not shown at all.  A list the
// reader cannot see through is a list they cannot check, and a confirmation
// given without checking confirms nothing.  Both halves of the rule follow from
// that one fact:
//
//   - eraseAtOnce is how many paths the question prints IN FULL.  Seven is what
//     is left under the head of the question in a window of 24 rows by 80 — the
//     size tools/sverka-ui.sh draws the screen at, and the smallest anybody
//     really works in.  It is measured, not chosen:
//     TestSevenIsWhatFitsInTheSmallestWindow counts the lines and fails if the
//     question or the window ever changes shape.  Above seven the list is cut
//     with «…и ещё N», and a cut list means the number gets typed.
//   - Size is not visible as length: one file of forty gigabytes fits on one
//     line.  So the volume has a ceiling too, and THE SCREEN DOES NOT PICK IT.
//     It is «Порог крупного», the size at which the decision layer stops
//     calling a file ordinary, carried on the plan (clean.Plan.LargeBytes, see
//     core.Sizer).  A layer that does not state one leaves the screen with no
//     idea what large means, and then the number is typed every time.
const eraseAtOnce = 7

// askChrome is how many lines askLines draws around a question's own: a blank,
// the title, a blank above, and a blank and the hint below.  fitsOnScreen adds
// it back to ask whether the whole question is visible.
const askChrome = 5

// ask is the state of a question the screen is asking.
type ask struct {
	kind    overlay
	title   string
	lines   []string
	input   string   // what has been typed
	hint    string   // what to type
	choices []string // completion candidates, for the path prompt
	err     string
	want    int    // the number that has to be typed back
	box     string // корзина the question is about
	plan    *clean.Plan
	journal *clean.Journal

	// quick says this question may be answered with one key — «y» — because
	// what it is asking about is small enough to be shown whole and small
	// enough by the decision layer's own reckoning.  It is only ever set on
	// overlayErase; see eraseAtOnce.
	quick bool
}

// openPath asks where to walk.  It starts from the root of the last walk, so
// the frequent case — look next door — is a few keys and not a whole path.
//
// WITH NOTHING WALKED YET the текущий каталог is offered instead.  That is the
// case `digitdisk analyze` without a path lands in, and it is the case the
// КОМАНДЫ section of the экран состояния lands in: the reader chose analyze
// and has typed no path, so the answer that is right most of the time is put
// in the line where they can SEE it before agreeing to it, and agreeing is one
// keystroke.
//
// A default that is one keystroke from a walk of millions of entries is a
// default that has to say so, and it does: the lines under the title name what
// Enter will cost and name the key that stops it.  Nothing else was needed —
// the walk shows its own numbers from the first second and q взводит выход из
// него, so the reader who agrees by accident finds out in one second and not
// in ten minutes.
func (w *walkScreen) openPath() {
	start := w.o.Root
	var lines []string
	if start == "" {
		if wd, err := os.Getwd(); err == nil {
			start = wd
			lines = []string{
				w.l.T("предложен текущий каталог — Enter соглашается с ним."),
				w.l.T("обходится всё дерево под ним: на домашнем каталоге это миллионы"),
				w.l.T("записей и минуты. Числа идут с первой секунды, q прерывает обход."),
			}
		}
	}
	if start != "" && !strings.HasSuffix(start, string(filepath.Separator)) {
		start += string(filepath.Separator)
	}
	w.ask = &ask{kind: overlayPath, title: w.l.T("ОБОЙТИ КАТАЛОГ"), input: start, lines: lines,
		hint: w.l.T("Tab — дополнить, Ctrl-U — стереть строку, Enter — обойти, Esc — отменить")}
	w.complete(false)
}

// complete fills in what is unambiguous and lists what is not.  Tab does it
// out loud (fill in and list); typing does it quietly (list only), so the
// candidates under the line follow the letters as they are typed.
func (w *walkScreen) complete(fill bool) {
	a := w.ask
	if a == nil || a.kind != overlayPath {
		return
	}
	a.err = ""
	dir, prefix := filepath.Split(a.input)
	// The head is kept exactly as it was typed — an empty one is read as the
	// working directory but must not be written back as "./", or completing
	// «tre» would answer «./tree/» and quietly rewrite what the reader typed.
	read := dir
	if read == "" {
		read = "."
	}
	entries, err := os.ReadDir(read)
	if err != nil {
		a.choices = nil
		if fill {
			a.err = w.l.F("каталог не читается: %s", lang.InLang(err, w.l))
		}
		return
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			continue // a walk starts at a directory
		}
		if strings.HasPrefix(e.Name(), prefix) {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	a.choices = names
	if !fill || len(names) == 0 {
		return
	}
	// The longest head every candidate shares is safe to fill in: it is
	// what the reader would have typed anyway.
	common := names[0]
	for _, n := range names[1:] {
		common = sharedHead(common, n)
	}
	if len(common) > len(prefix) {
		a.input = dir + common
	}
	if len(names) == 1 {
		a.input = dir + names[0] + string(filepath.Separator)
		w.complete(false)
	}
}

// sharedHead is the longest beginning two names have in common, counted in
// runes: cutting a Russian directory name by bytes would leave half a letter
// in the line the reader is typing.
func sharedHead(a, b string) string {
	ra, rb := []rune(a), []rune(b)
	i := 0
	for i < len(ra) && i < len(rb) && ra[i] == rb[i] {
		i++
	}
	return string(ra[:i])
}

// marked is the paths the reader has ticked, in the order they were ticked.
func (w *walkScreen) markedPaths() []string {
	out := make([]string, 0, len(w.marks))
	for p := range w.marks {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// toggleMark ticks the row under the cursor.  Only a directory can be ticked:
// the row that stands for the files themselves is not a place, and ticking it
// would mean something different from every other row on the screen.
func (w *walkScreen) toggleMark() {
	rows := w.here()
	if len(w.sel) == 0 || len(rows) == 0 {
		return
	}
	i := w.sel[len(w.sel)-1]
	if i < 0 || i >= len(rows) || rows[i].node == nil {
		w.tell(w.l.T("ОТМЕТКА"), []string{
			w.l.F("Отмечается каталог, а не строка «%s».", ownFilesRow(w.l)),
			w.l.T("Она стоит за файлы, которые лежат прямо здесь; чтобы взять их,"),
			w.l.T("отметьте сам этот каталог — «.» отмечает тот, в котором стоите.")})
		return
	}
	p := rows[i].node.path()
	w.atNode = nil // the rows carry the ticks, so they are built again
	if w.marks == nil {
		w.marks = map[string]bool{}
	}
	if w.marks[p] {
		delete(w.marks, p)
	} else {
		w.marks[p] = true
	}
}

// markHere ticks the directory the reader is standing in.
func (w *walkScreen) markHere() {
	if len(w.stack) == 0 {
		return
	}
	p := w.stack[len(w.stack)-1].path()
	w.atNode = nil
	if w.marks == nil {
		w.marks = map[string]bool{}
	}
	if w.marks[p] {
		delete(w.marks, p)
	} else {
		w.marks[p] = true
	}
}

// tell puts one answer on the screen for the reader to dismiss.
func (w *walkScreen) tell(title string, lines []string) {
	w.busy = false
	w.ask = &ask{kind: overlayNote, title: title, lines: lines,
		hint: w.l.T("любая клавиша — закрыть")}
}

// jobResult is what a piece of work done off the drawing loop came back with.
type jobResult struct {
	kind    overlay
	title   string
	where   string
	plan    *clean.Plan
	journal *clean.Journal
	box     string
	err     string
}

// work runs one call off the drawing loop and puts a "this is happening"
// notice up meanwhile.  Making a plan means walking the tree again, and on a
// tree of millions of entries that is a minute: a screen frozen for a minute
// is a screen that looks broken, and a keyboard that answers during it is a
// keyboard that can confirm something by accident.  Both are answered here —
// the notice moves, and nothing but Ctrl-C is listened to until the work is
// back.
func (w *walkScreen) work(title string, lines []string, f func() jobResult) {
	if w.busy {
		return
	}
	w.busy = true
	w.ask = &ask{kind: overlayNote, title: title, lines: lines, hint: w.l.T("идёт работа…")}
	out := w.jobs
	go func() { out <- f() }()
}

// accepted takes the answer of a finished job.
func (w *walkScreen) accepted(r jobResult) {
	w.busy = false
	if r.err != "" {
		w.tell(r.title, []string{r.err})
		return
	}
	switch r.kind {
	case overlayPlan:
		p := r.plan
		w.ask = &ask{kind: overlayPlan, title: w.l.T("ПЛАН УБОРКИ"), plan: p, want: len(p.Items),
			lines: w.planLines(p, r.where),
			hint:  w.l.F("наберите число файлов (%d) и Enter — перенести; Ctrl-U стереть, Esc отменить", len(p.Items))}
		if len(p.Items) == 0 {
			w.ask.kind = overlayNote
			w.ask.hint = w.l.T("любая клавиша — закрыть")
		}
	case overlayErase:
		p := r.plan
		if len(p.Items) == 0 {
			w.tell(w.l.T("СТИРАТЬ НЕЧЕГО"), w.eraseNothing(p, r.where))
			return
		}
		lines := w.eraseLines(p, r.where)
		a := &ask{kind: overlayErase, title: w.l.T("СТЕРЕТЬ НАСОВСЕМ"), plan: p,
			want: len(p.Items), lines: lines, quick: w.quickEnough(p, lines)}
		if a.quick {
			a.hint = w.l.F("y — стереть насовсем (%d); Esc — отменить, ничего не тронуть", a.want)
		} else {
			a.hint = w.l.F("наберите число файлов (%d) и Enter — стереть насовсем; Ctrl-U стереть, Esc отменить", a.want)
		}
		w.ask = a
	case overlayBack:
		j := r.journal
		// A dry run moves nothing and therefore stamps nothing: what would
		// come back is what was moved, is still in the корзина, and the run
		// found no reason to refuse.  This is the count the printed
		// `restore --dry-run` names, arrived at the same way.
		n, bytes := wouldReturn(j)
		if n == 0 {
			w.tell(w.l.T("ВОЗВРАЩАТЬ НЕЧЕГО"),
				[]string{w.l.T("в этой корзине нет файлов, которые можно вернуть"), "", r.box})
			return
		}
		w.ask = &ask{kind: overlayBack, title: w.l.T("ВОЗВРАТ ИЗ КОРЗИНЫ"), box: r.box, want: n,
			lines: []string{
				r.box,
				"",
				w.l.F("вернётся на прежние места   %s файлов, %s", w.l.Num(int64(n)), w.l.Bytes(bytes)),
				w.l.T("ничего не перезаписывается: файл, вернувшийся на занятое место, останется в корзине."),
			},
			hint: w.l.F("наберите число файлов (%d) и Enter — вернуть; Ctrl-U стереть, Esc отменить", n)}
	case overlayNote:
		w.marks, w.history, w.atNode = nil, nil, nil
		w.tell(r.title, strings.Split(r.where, "\n"))
	}
}

// proposeClean asks the decision layer what it would remove under what is
// marked, and shows the plan.  Nothing is touched: clean.Make opens no file
// for writing and makes no корзина.
func (w *walkScreen) proposeClean() {
	if w.o.Plan == nil {
		w.tell(w.l.T("УБОРКА"), []string{w.l.T("этот экран собран без уборки")})
		return
	}
	only := w.markedPaths()
	where := w.l.F("отмечено каталогов: %d", len(only))
	if len(only) == 0 {
		where = w.l.T("отмечено ничего — план по всему обойденному дереву")
	}
	root, plan, l := w.walkRoot, w.o.Plan, w.l
	w.work(w.l.T("СТРОИТСЯ ПЛАН УБОРКИ"), []string{
		where,
		"",
		w.l.T("дерево обходится заново, и заново судится решающим слоем:"),
		w.l.T("план строится по тому, что на диске сейчас, а не по тому, что было."),
		w.l.T("ничего не открывается на запись и корзина не создаётся."),
	}, func() jobResult {
		p, err := plan(root, only)
		if err != nil {
			return jobResult{title: l.T("УБОРКА НЕ ВЫШЛА"), err: lang.InLang(err, l)}
		}
		return jobResult{kind: overlayPlan, plan: p, where: where}
	})
}

// quickEnough decides whether this erasure may go on one keystroke.  All three
// have to hold, and every one of them is about what the reader can actually
// check before answering — see eraseAtOnce for why.
func (w *walkScreen) quickEnough(p *clean.Plan, lines []string) bool {
	switch {
	case len(p.Items) > eraseAtOnce:
		return false // the list is cut, and a cut list is not read
	case len(lines)+askChrome > w.bodyHeight():
		return false // it does not fit THIS window, and the question does not scroll
	case p.LargeBytes <= 0:
		return false // the layer named no size at which large begins
	default:
		return p.Bytes <= p.LargeBytes
	}
}

// proposeErase asks the decision layer what, out of the ground the reader
// pointed at, may go — and then shows it as what it is: a removal with no
// корзина behind it.
//
// THE GROUND IS WHAT WAS MARKED, or the row under the cursor when nothing was.
// That is the whole of what забой decides.  What may go on that ground is not
// its business and never becomes its business: the same w.o.Plan the «c» key
// calls is called here, with the same корень and the same справочник and the
// same защитный список, and it answers with the same приговор ядра.  A path the
// layer did not mark cannot be erased by pointing at it harder.
func (w *walkScreen) proposeErase() {
	if w.o.Plan == nil || w.o.Erase == nil {
		w.tell(w.l.T("СТЕРЕТЬ НАСОВСЕМ"), []string{w.l.T("этот экран собран без стирания")})
		return
	}
	only := w.markedPaths()
	where := w.l.F("отмечено каталогов: %d", len(only))
	if len(only) == 0 {
		n, ok := w.rowNode()
		if !ok {
			return // rowNode has said why
		}
		only = []string{n.path()}
		where = w.l.F("под курсором: %s", tail(n.path(), 52))
	}
	root, plan, l := w.walkRoot, w.o.Plan, w.l
	w.work(w.l.T("СЧИТАЕТСЯ СТИРАНИЕ"), []string{
		where,
		"",
		w.l.T("дерево обходится заново, и заново судится решающим слоем:"),
		w.l.T("стирается только то, чему ядро вынесло «МожноУбрать» СЕЙЧАС."),
		w.l.T("пока считается, не тронут ни один файл."),
	}, func() jobResult {
		p, err := plan(root, only)
		if err != nil {
			return jobResult{title: l.T("СТИРАНИЕ НЕ ВЫШЛО"), err: lang.InLang(err, l)}
		}
		return jobResult{kind: overlayErase, plan: p, where: where}
	})
}

// rowNode is the directory the cursor stands on.  Only a directory: the row
// that stands for the files lying here is not a place, and erasing "a row"
// would mean something different from erasing every other row on the screen.
func (w *walkScreen) rowNode() (*wnode, bool) {
	rows := w.here()
	if len(w.sel) == 0 || len(rows) == 0 {
		w.tell(w.l.T("СТЕРЕТЬ НАСОВСЕМ"), []string{w.l.T("здесь нет строки, на которую наведён курсор")})
		return nil, false
	}
	i := w.sel[len(w.sel)-1]
	if i < 0 || i >= len(rows) || rows[i].node == nil {
		w.tell(w.l.T("СТЕРЕТЬ НАСОВСЕМ"), []string{
			w.l.F("Стирается каталог, а не строка «%s».", ownFilesRow(w.l)),
			w.l.T("Она стоит за файлы, которые лежат прямо здесь; чтобы взять их,"),
			w.l.T("отметьте сам этот каталог — «.» отмечает тот, в котором стоите.")})
		return nil, false
	}
	return rows[i].node, true
}

// eraseLines is the question as the reader reads it.  The two numbers that
// matter stand FIRST, before the ground and before the list: a question longer
// than the window is cut from the bottom, and «сколько файлов» must not be what
// gets cut.
func (w *walkScreen) eraseLines(p *clean.Plan, where string) []string {
	out := []string{
		w.l.F("стереть насовсем  %s файлов, %s", w.l.Num(int64(len(p.Items))), w.l.Bytes(p.Bytes)),
		w.l.F("освободится %s — на этот раз по-настоящему", w.l.Bytes(p.FreeableBytes)),
		w.l.T("корзины не будет: файлы уходят с диска, и вернуть их нечем."),
		"",
		where,
	}
	if n := len(p.Protected); n > 0 {
		out = append(out, w.l.F("защитный список оставил на месте %d (%s)", n, w.l.Bytes(p.ProtectedBytes)))
	}
	if n := len(p.Refused); n > 0 {
		out = append(out, w.l.F("хозяин отказался трогать %d — слои разошлись, см. `digitdisk clean`", n))
	}
	out = append(out, "", w.l.T("исчезнет:"))
	for i, it := range p.Items {
		if i >= eraseAtOnce {
			out = append(out, w.l.F("   …и ещё %d", len(p.Items)-i))
			break
		}
		out = append(out, fmt.Sprintf("   %10s  %-12s %s",
			w.l.Bytes(it.Size), w.l.Word(string(it.Class)), tail(it.Path, 52)))
	}
	return out
}

// eraseNothing is the answer when the ground held nothing the layer would let
// go.  It says WHICH of the three possible reasons it was, because "нечего" on
// its own is indistinguishable from a tool that did not look.
func (w *walkScreen) eraseNothing(p *clean.Plan, where string) []string {
	out := []string{where, "", w.l.T("к стиранию 0 файлов.")}
	if n := len(p.Protected); n > 0 {
		out = append(out, "", w.l.F("защитный список оставил на месте %d (%s):", n, w.l.Bytes(p.ProtectedBytes)))
		for i, pr := range p.Protected {
			if i >= 4 {
				out = append(out, w.l.F("   …и ещё %d", len(p.Protected)-i))
				break
			}
			out = append(out, "   "+tail(pr.Path, 44)+" — "+pr.Rule.In(w.l))
		}
	}
	if n := len(p.Refused); n > 0 {
		out = append(out, "", w.l.F("хозяин отказался трогать %d:", n))
		for i, r := range p.Refused {
			if i >= 4 {
				out = append(out, w.l.F("   …и ещё %d", len(p.Refused)-i))
				break
			}
			out = append(out, "   "+tail(r.Path, 44)+" — "+r.Reason.In(w.l))
		}
	}
	if len(p.Protected) == 0 && len(p.Refused) == 0 {
		out = append(out, "",
			w.l.T("решающий слой не пометил «МожноУбрать» ничего из отмеченного."),
			w.l.T("отметка не делает файл убираемым — приговор выносит ядро."))
	}
	return out
}

// applyErase erases what the question listed, once it has been confirmed.
func (w *walkScreen) applyErase() {
	a := w.ask
	if a == nil || a.plan == nil || w.o.Erase == nil {
		return
	}
	erase, plan, l := w.o.Erase, a.plan, w.l
	w.work(w.l.T("СТИРАНИЕ ИДЁТ"), []string{
		w.l.F("файлов %s", w.l.Num(int64(len(plan.Items)))),
		"",
		w.l.T("журнал пишется до того, как исчезнет первый файл."),
	}, func() jobResult {
		j, err := erase(plan)
		if err != nil {
			return jobResult{title: l.T("СТИРАНИЕ НЕ ВЫШЛО"), err: lang.InLang(err, l)}
		}
		return jobResult{kind: overlayNote, title: l.T("СТЁРТО НАСОВСЕМ"),
			where: strings.Join(erasedLines(l, j), "\n")}
	})
}

// erasedLines is what came of the erasure.
func erasedLines(l lang.Lang, j *clean.Journal) []string {
	n, bytes := j.Purged()
	lines := []string{
		l.F("стёрто          %s файлов, %s", l.Num(int64(n)), l.Bytes(bytes)),
		l.T("место освобождено: этих файлов на диске больше нет, и возврата нет."),
		"",
		l.F("журнал: %s", j.Path()),
		l.T("он называет всё, что исчезло, и стоит в разделе ЖУРНАЛ как стирание."),
	}
	if f := j.Failed(); len(f) > 0 {
		lines = append(lines, "", l.F("не тронуто %d — файл изменился между обходом и стиранием:", len(f)))
		for i, it := range f {
			if i >= 4 {
				break
			}
			lines = append(lines, "   "+tail(it.Path, 50)+" — "+it.Failed.In(l))
		}
	}
	return lines
}

// planLines is the plan as the reader reads it.
func (w *walkScreen) planLines(p *clean.Plan, where string) []string {
	out := []string{where, ""}
	out = append(out,
		w.l.F("к переносу      %s файлов, %s", w.l.Num(int64(len(p.Items))), w.l.Bytes(p.Bytes)),
		w.l.T("место не освободится: перенос — это переименование, байты остаются на диске."),
		w.l.T("освободит только `digitdisk purge <корзина> --confirm N`."),
		"")
	if len(p.ByClass) > 0 {
		out = append(out, w.l.T("по разрядам:"))
		for _, c := range p.ByClass {
			out = append(out, fmt.Sprintf("   %-14s %8s  %12s",
				w.l.Word(string(c.Class)), w.l.F("%s файлов", w.l.Num(int64(c.Count))), w.l.Bytes(c.Bytes)))
		}
		out = append(out, "")
	}
	if len(p.Items) == 0 {
		out = append(out, w.l.T("решающий слой не пометил «МожноУбрать» ничего из отмеченного."),
			w.l.T("отметка не делает файл убираемым — приговор выносит ядро."))
	}
	if n := len(p.Protected); n > 0 {
		out = append(out, w.l.F("защитный список оставил на месте %d (%s)", n, w.l.Bytes(p.ProtectedBytes)))
	}
	if n := len(p.Refused); n > 0 {
		out = append(out, w.l.F("хозяин отказался трогать %d — слои разошлись, см. `digitdisk clean`", n))
	}
	out = append(out, "", w.l.F("корзина: %s", p.Trash))
	if len(p.Items) > 0 {
		out = append(out, "", w.l.T("первые из списка:"))
		for i, it := range p.Items {
			if i >= 6 {
				out = append(out, w.l.F("   …и ещё %d", len(p.Items)-i))
				break
			}
			out = append(out, fmt.Sprintf("   %10s  %-12s %s",
				w.l.Bytes(it.Size), w.l.Word(string(it.Class)), tail(it.Path, 52)))
		}
	}
	return out
}

// applyClean moves what the plan lists, once the number has been typed back.
func (w *walkScreen) applyClean() {
	a := w.ask
	if a == nil || a.plan == nil || w.o.Apply == nil {
		return
	}
	apply, plan, l := w.o.Apply, a.plan, w.l
	w.work(w.l.T("ПЕРЕНОС В КОРЗИНУ"), []string{
		w.l.F("файлов %s", w.l.Num(int64(len(plan.Items)))),
		"",
		w.l.T("журнал пишется до того, как сдвинется первый файл."),
	}, func() jobResult {
		j, err := apply(plan)
		if err != nil {
			return jobResult{title: l.T("УБОРКА НЕ ВЫШЛА"), err: lang.InLang(err, l)}
		}
		return jobResult{kind: overlayNote, title: l.T("УБРАНО В КОРЗИНУ"),
			where: strings.Join(appliedLines(l, j), "\n")}
	})
}

// appliedLines is what came of the move.
func appliedLines(l lang.Lang, j *clean.Journal) []string {
	n, bytes := j.Moved()
	lines := []string{
		l.F("перенесено      %s файлов, %s", l.Num(int64(n)), l.Bytes(bytes)),
		l.T("место НЕ освобождено: файлы лежат в корзине под другими именами."),
		"",
		l.F("корзина: %s", j.Path()),
		l.T("вернуть — раздел ЖУРНАЛ, клавиша Enter на этой корзине."),
		l.T("стереть насовсем — только отдельной командой:"),
		fmt.Sprintf("   digitdisk purge %s --confirm %d", j.Path(), n),
	}
	if f := j.Failed(); len(f) > 0 {
		lines = append(lines, "", l.F("не тронуто %d — файл изменился между обходом и переносом:", len(f)))
		for i, it := range f {
			if i >= 4 {
				break
			}
			lines = append(lines, "   "+tail(it.Path, 50)+" — "+it.Failed.In(l))
		}
	}
	return lines
}

// proposeRestore shows what putting a корзина back would do, and asks for the
// number.  The plan comes from clean.Restore itself, run without touching
// anything.
func (w *walkScreen) proposeRestore(box string) {
	if w.o.Restore == nil {
		w.tell(w.l.T("ВОЗВРАТ"), []string{w.l.T("этот экран собран без возврата")})
		return
	}
	back, l := w.o.Restore, w.l
	w.work(w.l.T("СЧИТАЕТСЯ ВОЗВРАТ"),
		[]string{box, "", w.l.T("журнал корзины читается; ничего не двигается.")},
		func() jobResult {
			j, err := back(box, true)
			if err != nil {
				return jobResult{title: l.T("ВОЗВРАТ НЕ ВЫШЕЛ"), err: lang.InLang(err, l)}
			}
			return jobResult{kind: overlayBack, journal: j, box: box}
		})
}

func (w *walkScreen) applyRestore() {
	a := w.ask
	if a == nil || w.o.Restore == nil {
		return
	}
	back, box, l := w.o.Restore, a.box, w.l
	w.work(w.l.T("ВОЗВРАТ ИДЁТ"), []string{box}, func() jobResult {
		j, err := back(box, false)
		if err != nil {
			return jobResult{title: l.T("ВОЗВРАТ НЕ ВЫШЕЛ"), err: lang.InLang(err, l)}
		}
		n, bytes := j.Restored()
		return jobResult{kind: overlayNote, title: l.T("ВОЗВРАЩЕНО"), where: strings.Join([]string{
			l.F("вернулось на прежние места  %s файлов, %s", l.Num(int64(n)), l.Bytes(bytes)),
			"", l.F("корзина: %s", box)}, "\n")}
	})
}

// wouldReturn counts what a dry run of restore says would come back.
func wouldReturn(j *clean.Journal) (n int, bytes int64) {
	for _, it := range j.Items {
		if it.MovedAt != "" && it.RestoredAt == "" && it.PurgedAt == "" && it.Failed.Empty() {
			n++
			bytes += it.Size
		}
	}
	return n, bytes
}

// askKeys handles a keypress while a question is on the screen.  It reports
// whether the screen should close.
func (w *walkScreen) askKey(k key) bool {
	a := w.ask
	if w.busy {
		// Nothing is answered while work is out: a keypress meant for the
		// notice must not land on whatever question comes back in its place.
		return k.kind == keyCtrlC
	}
	switch a.kind {
	case overlayNote, overlayKeys:
		if k.kind == keyCtrlC {
			return true
		}
		w.ask = nil
		return false
	case overlayPath:
		switch k.kind {
		case keyEsc:
			w.ask = nil
		case keyCtrlC:
			return true
		case keyTab:
			w.complete(true)
		case keyBack:
			if n := len([]rune(a.input)); n > 0 {
				a.input = string([]rune(a.input)[:n-1])
				w.complete(false)
			}
		case keyKill:
			a.input = ""
			w.complete(false)
		case keyEnter:
			p := strings.TrimSpace(a.input)
			if p == "" {
				a.err = w.l.T("путь пуст")
				return false
			}
			// The completion leaves a separator on the end, which is right
			// while typing and wrong as a корень: every path digitdisk prints
			// afterwards would carry it.
			p = filepath.Clean(p)
			if fi, err := os.Stat(p); err != nil {
				a.err = lang.InLang(err, w.l)
				return false
			} else if !fi.IsDir() {
				a.err = w.l.T("это не каталог")
				return false
			}
			w.ask = nil
			w.startWalk(p)
		case keyRune:
			if k.r < 0x20 {
				// A key the reader pressed that this screen does not answer
				// arrives as an empty rune.  It is not a letter and must not
				// be typed into the path.
				return false
			}
			if a.input == "" && (k.r == 'q' || k.r == 'й') {
				// An empty line and «q» is somebody asking to leave, not
				// somebody starting to type a path called q.
				return true
			}
			a.input += string(k.r)
			w.complete(false)
		}
		return false
	case overlayPlan, overlayErase, overlayBack:
		switch k.kind {
		case keyEsc:
			w.ask = nil
		case keyCtrlC:
			return true
		case keyBack:
			// Here забой is a text key and nothing else: this is a line
			// being typed, and a question already on the screen is not
			// answered by the key that opened it.
			if n := len(a.input); n > 0 {
				a.input = a.input[:n-1]
			}
			a.err = ""
		case keyKill:
			a.input, a.err = "", ""
		case keyEnter:
			typed := strings.TrimSpace(a.input)
			if typed == "" && a.quick {
				a.err = w.l.T("стирание подтверждает «y», а не Enter")
				return false
			}
			if typed != fmt.Sprint(a.want) {
				a.err = w.l.F("названо %q, а файлов %d — ничего не тронуто", typed, a.want)
				a.input = ""
				return false
			}
			switch a.kind {
			case overlayPlan:
				w.applyClean()
			case overlayErase:
				w.applyErase()
			default:
				w.applyRestore()
			}
		case keyRune:
			if k.r >= '0' && k.r <= '9' {
				a.input += string(k.r)
				a.err = ""
				return false
			}
			// «y» — and «н», which is the same physical key on a Russian
			// layout, the way «q» is also «й» everywhere else on this
			// screen.  It answers only the question that offered it, and
			// only while that question is the small one.
			if a.kind == overlayErase && (k.r == 'y' || k.r == 'Y' || k.r == 'н' || k.r == 'Н') {
				if a.quick {
					w.applyErase()
					return false
				}
				a.err = w.l.F("файлов %d — столько одной клавишей не стирается, назовите число", a.want)
			}
		}
		return false
	}
	return false
}

// askLines draws the question over the body of the section.
func (w *walkScreen) askLines() []string {
	t, a := w.t, w.ask
	out := []string{"", w.t.Bold(t.P.Accent, "  "+a.title), ""}
	for _, l := range a.lines {
		out = append(out, w.t.Fg(t.P.Foreground, "  "+fit(l, maxInt(1, w.cols-3))))
	}
	if a.kind == overlayPath {
		out = append(out, "", w.t.Fg(t.P.AccentSoft, "  "+w.l.T("путь: ")+a.input+"▏"))
		if len(a.choices) > 0 {
			out = append(out, "", w.note(w.l.F("подходит каталогов: %d", len(a.choices))))
			for i, c := range a.choices {
				if i >= 8 {
					out = append(out, w.note(w.l.F("   …и ещё %d", len(a.choices)-i)))
					break
				}
				out = append(out, w.t.Fg(t.P.Muted, "     "+fit(c+"/", maxInt(1, w.cols-8))))
			}
		}
	}
	if a.kind == overlayPlan || a.kind == overlayBack || (a.kind == overlayErase && !a.quick) {
		out = append(out, "", w.t.Fg(t.P.AccentSoft, "  "+w.l.T("число файлов: ")+a.input+"▏"))
	}
	if a.err != "" {
		out = append(out, "", w.t.Fg(t.P.Red, "  "+fit(a.err, maxInt(1, w.cols-3))))
	}
	if a.hint != "" {
		out = append(out, "", w.note(a.hint))
	}
	return out
}

// keysLines is the whole keyboard, and the commands the screen stands for.
func (w *walkScreen) keysLines() []string {
	rows := [][2]string{
		{"Tab, 1…8", w.l.T("разделы")},
		{"↑ ↓, k j", w.l.T("строка; g G — начало и конец")},
		{"→, Enter", w.l.T("внутрь каталога (в разделе ДЕРЕВО)")},
		{w.l.T("←"), w.l.T("назад из каталога")},
		{w.l.T("Пробел"), w.l.T("отметить каталог; «.» — тот, в котором стоите")},
		{"c", w.l.T("план уборки по отмеченному и подтверждение")},
		{w.l.T("забой"), w.l.T("СТЕРЕТЬ НАСОВСЕМ отмеченное, а без отметок — строку под курсором")},
		{"o", w.l.T("обойти другой каталог (Tab дополняет путь)")},
		{"Enter " + w.l.T("в ЖУРНАЛЕ"), w.l.T("вернуть корзину на прежние места")},
		{"l", w.l.T("язык экрана: русский или English")},
		{"?", w.l.T("эта справка")},
		{"q", w.l.T("выход; отчёт печатается как всегда")},
	}
	out := []string{"", w.t.Bold(w.t.P.Accent, "  "+w.l.T("КЛАВИШИ И КОМАНДЫ")), ""}
	for _, r := range rows {
		var line row
		line.add("  "+fit(r[0], 16)+"  ", func(x string) string { return w.t.Fg(w.t.P.AccentSoft, x) })
		line.add(fit(r[1], maxInt(1, w.cols-line.w-1)), func(x string) string {
			return w.t.Fg(w.t.P.Foreground, strings.TrimRight(x, " "))
		})
		out = append(out, line.String())
	}
	// The commands come from internal/cli, the one list справка, страница
	// руководства and the КОМАНДЫ section of the status screen are also
	// built from.
	//
	// Both screens now START them, and both start them the same way: this
	// one is where a walk is acted on, and the other is where a walk is
	// asked for.  Removal is behind the same plan and the same count on
	// either road; `purge` is on neither.
	out = append(out, "", w.caption(w.l.T("КОМАНДЫ")), "")
	for _, c := range cli.Commands {
		var line row
		line.add("  "+fit(c.Call(w.l), 20), func(x string) string { return w.t.Fg(w.t.P.Accent, x) })
		line.add(fit(c.Gloss, maxInt(1, w.cols-line.w-1)), func(x string) string {
			return w.t.Fg(w.t.P.Muted, strings.TrimRight(x, " "))
		})
		out = append(out, line.String())
	}
	out = append(out, "",
		w.note(w.l.T("уборка с экрана идёт тем же путём, что `digitdisk clean`: приговор ядра, план, подтверждение числом.")),
		w.note(w.l.T("забой идёт тем же путём и тем же приговором ядра, но без корзины: стёртое не возвращается.")),
		w.note(w.l.T("сколько-то файлов — подтверждается одной «y»; больше или крупнее — только числом файлов.")),
		w.note(w.l.T("корзину целиком стирает отдельная команда `digitdisk purge`; с экрана она не запускается.")))
	return out
}

// --- the sections that read what has been decided and done ------------------

// places is the справочник: what digitdisk knows about particular caches, and
// which of them are on this machine.  Sizes are not measured here — measuring
// a hundred places means a hundred walks, and this screen is already walking
// one tree; `digitdisk places` measures them.
func (w *walkScreen) placesSection() []string {
	if w.o.Places == nil {
		return []string{"", w.note(w.l.T("этот экран собран без справочника"))}
	}
	if w.placesFound == nil {
		dir, found, err := w.o.Places()
		if err != nil {
			return []string{"", w.note(w.l.F("справочник не прочитан: %s", lang.InLang(err, w.l)))}
		}
		w.placesDir, w.placesFound = dir, found
		if w.placesFound == nil {
			w.placesFound = []PlaceRow{}
		}
	}
	out := []string{""}
	out = append(out, w.kv(w.l.T("справочник"), w.placesDir, ""))
	out = append(out, w.kv(w.l.T("нашлось здесь"), w.l.Num(int64(len(w.placesFound))),
		w.l.T("размеры считает `digitdisk places`")))
	out = append(out, "", w.caption(w.l.T("ЕСТЬ НА ЭТОЙ МАШИНЕ")), "")
	for _, f := range w.placesFound {
		var r row
		r.add("  "+fit(w.l.Word(f.Class), 12)+" ", func(x string) string { return w.t.Fg(w.t.P.Purple, x) })
		r.add(fit(f.Name, 22)+" ", func(x string) string { return w.t.Fg(w.t.P.Foreground, x) })
		r.add(fit(tail(f.Path, maxInt(8, w.cols-r.w-2)), maxInt(8, w.cols-r.w-2)), func(x string) string {
			return w.t.Fg(w.t.P.Muted, strings.TrimRight(x, " "))
		})
		out = append(out, r.String())
	}
	if len(w.placesFound) == 0 {
		out = append(out, w.note(w.l.T("ни одного известного места на этой машине не нашлось")))
	}
	return out
}

// journal is what past уборки did under this root, and the place they can be
// undone from.
func (w *walkScreen) journalSection() []string {
	if w.o.History == nil {
		return []string{"", w.note(w.l.T("этот экран собран без журнала"))}
	}
	if w.history == nil {
		h, err := w.o.History(w.walkRoot)
		if err != nil {
			return []string{"", w.note(w.l.F("журнал не прочитан: %s", lang.InLang(err, w.l)))}
		}
		w.history = h
		if w.boxSel >= len(h.Entries) {
			w.boxSel = 0
		}
	}
	h := w.history
	out := []string{""}
	out = append(out, w.kv(w.l.T("корзин"), w.l.Num(int64(h.Boxes)), w.l.F("хранилище %s", h.Trash)))
	out = append(out, w.kv(w.l.T("в корзинах"), w.l.Bytes(h.MovedBytes),
		w.l.T("место не освобождено, пока не purge")))
	out = append(out, w.kv(w.l.T("освобождено"), w.l.Bytes(h.FreedBytes), w.l.T("это уже стёрто и не вернётся")))
	out = append(out, "")
	if len(h.Entries) == 0 {
		return append(out, w.note(w.l.T("под этим корнем ещё не убирали")))
	}
	out = append(out, w.caption(w.l.T("КОРЗИНЫ")), "")
	for i, e := range h.Entries {
		var r row
		if i == w.boxSel {
			r.add(" ▶ ", func(x string) string { return w.t.Fg(w.t.P.Accent, x) })
		} else {
			r.plain("   ")
		}
		r.add(fit(filepath.Base(e.Box), 20)+" ", func(x string) string { return w.t.Fg(w.t.P.Foreground, x) })
		// An erasure counted nothing into a корзина, so the two columns that
		// say "how much is lying there" say what went instead.
		count, bytes := e.Moved, e.MovedBytes
		if e.Way == clean.WayErase {
			count, bytes = e.Purged, e.PurgedBytes
		}
		r.add(right(w.l.F("%s файл.", w.l.Num(int64(count))), 14)+" ", func(x string) string { return w.t.Fg(w.t.P.Muted, x) })
		r.add(right(w.l.Bytes(bytes), 11)+"  ", func(x string) string { return w.t.Fg(w.t.P.Muted, x) })
		state := w.l.T("в корзине")
		switch {
		case !e.Problem.Empty():
			state = w.l.F("беда: %s", e.Problem.In(w.l))
		case e.Way == clean.WayErase:
			// Never «стёрто N» plain: that is also what a purged корзина
			// says, and those two differ in whether anything was ever
			// restorable.  Here nothing ever was.
			state = w.l.F("стёрто насовсем %d — возврата нет", e.Purged)
		case e.Purged > 0:
			state = w.l.F("стёрто %d", e.Purged)
		case e.Restored > 0:
			state = w.l.F("возвращено %d", e.Restored)
		}
		r.add(fit(state, maxInt(1, w.cols-r.w-1)), func(x string) string {
			return w.t.Fg(w.t.P.Subtle, strings.TrimRight(x, " "))
		})
		out = append(out, r.String())
	}
	out = append(out, "",
		w.note(w.l.T("Enter — вернуть выбранную корзину на прежние места (спросит число файлов)")),
		w.note(w.l.T("запись «стёрто насовсем» вернуть нельзя: файлов нет, есть только список того, что было")),
		w.note(w.l.T("стирание корзины — только командой `digitdisk purge <корзина> --confirm N`")))
	return out
}

// journalHead is how many lines the ЖУРНАЛ section draws before its list.
const journalHead = 7
