// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/core"
	"digitdisk/internal/lang"
	"digitdisk/internal/scan"
)

// tree builds a small tree with known sizes and returns its root.
func testTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write := func(rel string, n int) {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, make([]byte, n), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("большой/файл.bin", 40000)
	write("большой/ещё.bin", 10000)
	write("большой/глубже/мелочь", 500)
	write("средний/один", 9000)
	write("мелкий/два", 100)
	write("сам-по-себе", 7000)
	return root
}

// walkWith runs a walk over root with the collector attached, the way analyze
// does when it has a screen.
func walkWith(t *testing.T, root string) (scan.Result, *walkFeed) {
	t.Helper()
	f := newWalkFeed(root)
	res, err := scan.Walk(scan.Options{Root: root, Top: 5, Now: time.Now(), Watch: f.step})
	if err != nil {
		t.Fatal(err)
	}
	f.settle()
	return res, f
}

// The one number the screen must never get wrong: the total it shows growing
// has to be the total the report finishes with.  They are the same arithmetic
// only because Watch carries the CHARGED bytes and not the entry's own size.
func TestFeedTotalsAreTheWalkTotals(t *testing.T) {
	root := testTree(t)
	res, f := walkWith(t, root)

	if f.bytes != res.TotalBytes {
		t.Errorf("объём на экране %d, в отчёте %d", f.bytes, res.TotalBytes)
	}
	if f.entries != int64(res.Entries) {
		t.Errorf("записей на экране %d, в отчёте %d", f.entries, res.Entries)
	}
	if f.files != int64(res.Files) || f.dirs != int64(res.Dirs) {
		t.Errorf("файлов/каталогов %d/%d, в отчёте %d/%d", f.files, f.dirs, res.Files, res.Dirs)
	}
	if f.tree.bytes != res.TotalBytes {
		t.Errorf("дерево сложилось в %d, обход насчитал %d", f.tree.bytes, res.TotalBytes)
	}
	if int(f.tree.entries) != res.Entries {
		t.Errorf("дерево держит %d записей, обход насчитал %d", f.tree.entries, res.Entries)
	}
}

// A hard link met a second time is charged nothing, and the tree must not
// charge it either: a screen that counted st_size would show more bytes than
// the report prints, and the two would never meet.
func TestFeedCountsAHardLinkOnce(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "один")
	if err := os.WriteFile(first, make([]byte, 5000), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(first, filepath.Join(root, "второй")); err != nil {
		t.Skipf("жёсткие ссылки недоступны: %v", err)
	}
	res, f := walkWith(t, root)
	if res.HardlinkDupes == 0 {
		t.Fatal("обход не заметил повторного имени")
	}
	if f.bytes != res.TotalBytes {
		t.Errorf("экран насчитал %d, отчёт %d — повторное имя посчитано дважды", f.bytes, res.TotalBytes)
	}
	if f.tree.bytes != res.TotalBytes {
		t.Errorf("дерево насчитало %d, отчёт %d", f.tree.bytes, res.TotalBytes)
	}
}

// The tree has to add up: what lies in the children plus what lies directly in
// a directory is that directory, at every level.
func TestTreeAddsUpAtEveryLevel(t *testing.T) {
	root := testTree(t)
	_, f := walkWith(t, root)

	var check func(*wnode)
	check = func(n *wnode) {
		sum, cnt := n.own, n.ownEntries
		for _, k := range n.kids {
			sum += k.bytes
			cnt += k.entries
			check(k)
		}
		if sum != n.bytes || cnt != n.entries {
			t.Errorf("%s: сложение даёт %d/%d, в узле %d/%d", n.path(), sum, cnt, n.bytes, n.entries)
		}
	}
	check(f.tree)

	kids := map[string]*wnode{}
	for _, k := range f.tree.kids {
		kids[k.name] = k
	}
	big, ok := kids["большой"]
	if !ok {
		t.Fatalf("каталог «большой» не попал в дерево: %v", kids)
	}
	if big.bytes != 50500 {
		t.Errorf("«большой» = %d Б, ждали 50500", big.bytes)
	}
	if big.topName != "файл.bin" || big.topBytes != 40000 {
		t.Errorf("крупнейший файл в «большой» — %s (%d)", big.topName, big.topBytes)
	}
	// The largest child of the root is the largest child of the root, and the
	// list the screen draws is that order.
	if got := f.tree.children()[0].name; got != "большой" {
		t.Errorf("первым в списке %q", got)
	}
}

// The root is not a child of itself.
func TestRootIsNotItsOwnChild(t *testing.T) {
	root := testTree(t)
	_, f := walkWith(t, root)
	for _, k := range f.tree.kids {
		if k.name == filepath.Base(root) && k.bytes == 0 {
			t.Errorf("корень попал в собственный перечень: %q", k.name)
		}
	}
}

// A relative root and "/" spell their paths differently, and the collector has
// to find the parent of an entry in both.
func TestFeedFindsTheParentWhateverTheRootLooksLike(t *testing.T) {
	cases := []struct{ path, root, want string }{
		{"/усы/лапы/хвост", "/усы", "/усы/лапы"},
		{"/усы", "/усы", "/усы"},
		{"/лапы", "/", "/"},
		{"каталог/файл", ".", "каталог"},
		{"файл", ".", "."},
	}
	for _, c := range cases {
		if got := parentOf(c.path, c.root); got != c.want {
			t.Errorf("parentOf(%q, %q) = %q, ждали %q", c.path, c.root, got, c.want)
		}
	}
}

// A walk with no screen attached must behave exactly as it did: Watch is an
// addition, not a change.
func TestWalkWithoutAWatcherIsUnchanged(t *testing.T) {
	root := testTree(t)
	bare, err := scan.Walk(scan.Options{Root: root, Top: 5, Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	watched, _ := walkWith(t, root)
	if bare.TotalBytes != watched.TotalBytes || bare.Entries != watched.Entries ||
		len(bare.Largest) != len(watched.Largest) {
		t.Errorf("обход со сборщиком дал другой итог:\n%+v\n%+v", bare, watched)
	}
}

// --- the screen ------------------------------------------------------------

// walkTestScreen builds a screen over a walked tree without a terminal: the
// sections are pure functions of the walk, and are tested as such.
func walkTestScreen(t *testing.T, cols int) *walkScreen {
	t.Helper()
	root := testTree(t)
	res, f := walkWith(t, root)
	w := &walkScreen{t: Theme{P: Carbon, d: depthTrue}, l: lang.RU, rows: 30, cols: cols, started: time.Now()}
	w.o.Root, w.walkRoot = root, root
	w.rate = []float64{0.2, 0.9, 0.5}
	w.finish(walkDone{res: res, tree: f.tree, snap: f.snapshot()})
	return w
}

func TestWalkScreenDrawsWithinTheTerminal(t *testing.T) {
	// The widths the status screen is checked at, and for the same reason: a
	// line wider than the terminal wraps and takes the layout with it.
	for _, cols := range []int{40, 60, 80, 120, 200} {
		w := walkTestScreen(t, cols)
		for i, sec := range walkSections {
			w.tab, w.scroll = i, 0
			for _, line := range w.frame() {
				if got := plainWidth(w.t.clip(line, cols)); got > cols {
					t.Errorf("%s / %d колонок: строка в %d ячеек", sec.title(w.l), cols, got)
				}
			}
		}
		// and the walking half, which is drawn before any of the sections
		walking := &walkScreen{t: w.t, l: w.l, rows: 30, cols: cols, started: time.Now(), snap: w.snap}
		walking.o.Root, walking.walkRoot = w.walkRoot, w.walkRoot
		walking.reorder(w.snap.Tops)
		for _, line := range walking.frame() {
			if got := plainWidth(walking.t.clip(line, cols)); got > cols {
				t.Errorf("обход / %d колонок: строка в %d ячеек", cols, got)
			}
		}
	}
}

func TestWalkFrameIsExactlyTheHeightOfTheTerminal(t *testing.T) {
	for _, rows := range []int{8, 12, 24, 60} {
		w := walkTestScreen(t, 100)
		w.rows = rows
		for i := range walkSections {
			w.tab, w.scroll = i, 0
			if got := len(w.frame()); got != rows {
				t.Errorf("раздел %s на %d строк дал кадр в %d", walkSections[i].title(w.l), rows, got)
			}
		}
		w.mode = modeWalking
		if got := len(w.frame()); got != rows {
			t.Errorf("обход на %d строк дал кадр в %d", rows, got)
		}
	}
}

// While the walk runs the largest directory is a guess, and the screen says so
// in one word.  When the walk is over the word goes.
func TestTheLiveListSaysItIsPreliminary(t *testing.T) {
	w := walkTestScreen(t, 120)
	live := &walkScreen{t: w.t, l: w.l, rows: 30, cols: 120, started: time.Now(), snap: w.snap}
	live.o.Root, live.walkRoot = w.walkRoot, w.walkRoot
	live.reorder(w.snap.Tops)
	body := plain(strings.Join(live.walking(), "\n"))
	if !strings.Contains(body, "ПРЕДВАРИТЕЛЬНО") {
		t.Errorf("промежуточный перечень не помечен:\n%s", body)
	}
	done := plain(strings.Join(w.total(), "\n"))
	if strings.Contains(done, "ПРЕДВАРИТЕЛЬНО") {
		t.Errorf("законченный обход всё ещё помечен предварительным:\n%s", done)
	}
}

// The footer always says how to leave, at every width — a full-screen program
// that does not is a trap.
func TestWalkFooterAlwaysSaysHowToLeave(t *testing.T) {
	for _, cols := range []int{20, 30, 40, 60, 80, 120, 200} {
		w := walkTestScreen(t, cols)
		if got := plain(w.footer("")); !strings.Contains(got, "q") {
			t.Errorf("%d колонок: подвал не назвал выход: %q", cols, got)
		}
		w.mode = modeWalking
		if got := plain(w.footer("")); !strings.Contains(got, "q") {
			t.Errorf("%d колонок, обход: подвал не назвал выход: %q", cols, got)
		}
	}
}

func TestBrowsingWalksIntoADirectoryAndBackOut(t *testing.T) {
	w := walkTestScreen(t, 120)
	w.tab = 1 // ДЕРЕВО
	if walkSections[w.tab].id != treeTitle {
		t.Fatalf("раздел %d — не дерево", w.tab)
	}
	rows := w.here()
	if len(rows) == 0 {
		t.Fatal("корень пуст")
	}
	if rows[0].name != "большой/" {
		t.Fatalf("первой строкой %q, ждали «большой/»", rows[0].name)
	}
	deep := len(w.stack)
	w.handle(key{kind: keyRight})
	if len(w.stack) != deep+1 {
		t.Fatal("стрелка вправо не вошла в каталог")
	}
	if got := w.stack[len(w.stack)-1].name; got != "большой" {
		t.Errorf("вошли в %q", got)
	}
	w.handle(key{kind: keyBack})
	if len(w.stack) != deep {
		t.Error("забой не вернул наверх")
	}
	w.handle(key{kind: keyEnter})
	if len(w.stack) != deep+1 {
		t.Error("Enter не вошёл внутрь")
	}
	// «l» is the language and nothing else on this screen: the vim pair h/l
	// is given up so that one letter never means two things.
	was := w.l
	w.handle(key{kind: keyRune, r: 'l'})
	if len(w.stack) != deep+1 {
		t.Error("«l» походила по дереву — она язык")
	}
	if w.l == was {
		t.Error("«l» не переключила язык")
	}
	w.handle(key{kind: keyRune, r: 'l'})
	w.handle(key{kind: keyLeft})
	if len(w.stack) != deep {
		t.Error("стрелка влево не вернула наверх")
	}
	// The root has nowhere above it, and asking must not take the screen with it.
	for i := 0; i < 5; i++ {
		w.handle(key{kind: keyLeft})
	}
	if len(w.stack) != 1 {
		t.Errorf("из корня ушли наверх: глубина %d", len(w.stack))
	}
	// The row for the files themselves is not a directory and must not open.
	w.sel[0] = len(w.here()) - 1
	for _, r := range w.here() {
		if r.node == nil {
			w.sel[0] = indexOfRow(w.here(), r.name)
		}
	}
	before := len(w.stack)
	w.handle(key{kind: keyEnter})
	if len(w.stack) != before {
		t.Error("строка «файлы этого каталога» открылась как каталог")
	}
}

func indexOfRow(rows []browseRow, name string) int {
	for i, r := range rows {
		if r.name == name {
			return i
		}
	}
	return 0
}

func TestCursorStaysInsideTheList(t *testing.T) {
	w := walkTestScreen(t, 120)
	w.tab = 1
	for i := 0; i < 100; i++ {
		w.handle(key{kind: keyDown})
	}
	if got := w.sel[0]; got != len(w.here())-1 {
		t.Errorf("курсор ушёл на %d при %d строках", got, len(w.here()))
	}
	for i := 0; i < 100; i++ {
		w.handle(key{kind: keyUp})
	}
	if w.sel[0] != 0 {
		t.Errorf("курсор ушёл выше первой строки: %d", w.sel[0])
	}
}

func TestWalkKeysMoveBetweenSectionsAndOut(t *testing.T) {
	w := walkTestScreen(t, 120)
	w.tab = 0
	w.handle(key{kind: keyTab})
	if w.tab != 1 {
		t.Errorf("Tab дал раздел %d", w.tab)
	}
	w.handle(key{kind: keyShiftTab})
	if w.tab != 0 {
		t.Errorf("Shift-Tab дал раздел %d", w.tab)
	}
	// Outside ДЕРЕВО the arrows are what they are on the status screen.
	w.handle(key{kind: keyRight})
	if w.tab != 1 {
		t.Errorf("стрелка вправо в ИТОГЕ дала раздел %d", w.tab)
	}
	w.tab = 2
	w.handle(key{kind: keyLeft})
	if w.tab != 1 {
		t.Errorf("стрелка влево дала раздел %d", w.tab)
	}
	w.handle(key{kind: keyRune, r: '5'})
	if w.tab != 4 {
		t.Errorf("клавиша 5 дала раздел %d", w.tab)
	}
	for _, k := range []key{{kind: keyRune, r: 'q'}, {kind: keyRune, r: 'й'}, {kind: keyEsc}, {kind: keyCtrlC}} {
		if !w.handle(k) {
			t.Errorf("клавиша %+v не закрыла экран", k)
		}
	}
	if w.handle(key{kind: keyRune, r: 'x'}) {
		t.Error("посторонняя клавиша закрыла экран")
	}
	// While the walk runs only leaving is answered: there is nothing yet to
	// walk around in.
	w.mode = modeWalking
	if !w.handle(key{kind: keyRune, r: 'q'}) {
		t.Error("q не прервала обход")
	}
	if w.handle(key{kind: keyTab}) {
		t.Error("Tab во время обхода закрыл экран")
	}
}

// Ф2 of the brief: the live list must not reshuffle so fast it cannot be read.
// Two directories within a hair of each other keep their order; a real
// overtaking still happens.
func TestTheLiveListIsCalmWhenTwoAreNeckAndNeck(t *testing.T) {
	w := &walkScreen{l: lang.RU}
	w.reorder([]topRow{{Name: "а", Bytes: 1000}, {Name: "б", Bytes: 990}})
	base := w.swaps
	for i := 0; i < 100; i++ {
		// б creeps past а by a per cent and back, a hundred times
		w.reorder([]topRow{{Name: "а", Bytes: 1000}, {Name: "б", Bytes: 1010}})
		w.reorder([]topRow{{Name: "а", Bytes: 1020}, {Name: "б", Bytes: 1010}})
	}
	if w.swaps != base {
		t.Errorf("перестановок %d на двух почти равных каталогах", w.swaps-base)
	}
	if w.order[0] != "а" {
		t.Errorf("порядок сорвался: %v", w.order)
	}
	// A directory that really is bigger still climbs.
	w.reorder([]topRow{{Name: "а", Bytes: 1000}, {Name: "б", Bytes: 5000}})
	if w.order[0] != "б" {
		t.Errorf("настоящий обгон не показан: %v", w.order)
	}
	// A directory the walk has not reached yet appears at the end and does not
	// displace anything until it grows.
	w.reorder([]topRow{{Name: "а", Bytes: 1000}, {Name: "б", Bytes: 5000}, {Name: "в", Bytes: 1}})
	if w.order[2] != "в" {
		t.Errorf("новый каталог влез в середину: %v", w.order)
	}
}

// Ф2, measured on a real tree rather than argued.  Point DIGITDISK_TREE at a
// large directory and the test replays the walk's own snapshots through the
// ordering twice — with the margin and without it — and prints both counts.
//
//	cd host && DIGITDISK_TREE=/srv go test ./internal/ui/ -run Calm -v
func TestLiveOrderIsCalmOnARealTree(t *testing.T) {
	root := os.Getenv("DIGITDISK_TREE")
	if root == "" {
		t.Skip("замер по настоящему дереву: DIGITDISK_TREE=<путь>")
	}
	f := newWalkFeed(root)
	done := make(chan struct{})
	go func() {
		scan.Walk(scan.Options{Root: root, Top: 1, Now: time.Now(), Watch: f.step})
		f.want.Store(false)
		close(done)
	}()

	var snaps []walkSnap
	tick := time.NewTicker(drawEvery)
	defer tick.Stop()
	running := true
	for running {
		select {
		case sn := <-f.out:
			snaps = append(snaps, sn)
		case <-tick.C:
			f.want.Store(true)
		case <-done:
			running = false
		}
	}
	if len(snaps) < 4 {
		t.Skipf("дерево обошлось за %d кадров — замерять нечего", len(snaps))
	}

	// What the eye sees is not how many positions the sort moved, but how
	// many FRAMES showed a different top of the list, and how far a name
	// travelled in one of them.  A list that changes its first ten rows in
	// every frame cannot be read; a list that changes them now and again is
	// the walk showing its work.
	const shown = 10
	replay := func(margin float64) (changed, maxJump, swaps int) {
		w := &walkScreen{l: lang.RU}
		var prev []string
		for _, sn := range snaps {
			w.reorderBy(sn.Tops, margin)
			top := w.order
			if len(top) > shown {
				top = top[:shown]
			}
			was := map[string]int{}
			for i, n := range prev {
				was[n] = i
			}
			diff := false
			for i, n := range top {
				j, seen := was[n]
				if !seen {
					continue
				}
				if j != i {
					diff = true
					if d := j - i; d > maxJump {
						maxJump = d
					}
				}
			}
			if diff {
				changed++
			}
			prev = append(prev[:0], top...)
		}
		return changed, maxJump, w.swaps
	}
	secs := float64(len(snaps)) * drawEvery.Seconds()
	lc, lj, ls := replay(0)
	cc, cj, cs := replay(calmMargin)
	t.Logf("дерево %s: кадров %d за %.1f с, записей %d", root, len(snaps), secs, snaps[len(snaps)-1].Entries)
	t.Logf("без порога:   кадров с новым порядком первой десятки %d из %d (%.2f/с), самый длинный скачок %d строк, перестановок всего %d",
		lc, len(snaps), float64(lc)/secs, lj, ls)
	t.Logf("с порогом %.0f%%: кадров с новым порядком первой десятки %d из %d (%.2f/с), самый длинный скачок %d строк, перестановок всего %d",
		calmMargin*100, cc, len(snaps), float64(cc)/secs, cj, cs)
	if cc > lc {
		t.Errorf("порог добавил движения: %d кадров против %d", cc, lc)
	}
}

// --- the same promise the status screen makes ------------------------------

func TestRunWalkRefusesWhereThereIsNoTerminal(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()
	called := false
	_, _, err := RunWalk(WalkOptions{Out: w, Root: ".", Walk: func(string, func(scan.Step), func() bool) (scan.Result, error) {
		called = true
		return scan.Result{}, nil
	}})
	if err != ErrNoTerminal {
		t.Errorf("RunWalk в трубу вернул %v, ждали ErrNoTerminal", err)
	}
	if called {
		t.Error("обход всё-таки запустился — вывод в трубу был бы не тот, что всегда")
	}
	t.Setenv("TERM", "dumb")
	if _, _, err := RunWalk(WalkOptions{Out: os.Stdout, Root: ".", Walk: func(string, func(scan.Step), func() bool) (scan.Result, error) {
		return scan.Result{}, nil
	}}); err != ErrNoTerminal {
		t.Errorf("RunWalk при TERM=dumb вернул %v, ждали ErrNoTerminal", err)
	}
}

func TestRunWalkNeedsAWalk(t *testing.T) {
	if _, _, err := RunWalk(WalkOptions{Out: os.Stdout}); err == nil {
		t.Error("RunWalk без обхода не пожаловался")
	}
}

// The sections are the sections of the printed report, in its order.
func TestWalkSectionsFollowThePrintedReport(t *testing.T) {
	want := []string{"ИТОГ", "ДЕРЕВО", "КРУПНЕЙШЕЕ", "МОЖНО УБРАТЬ", "РАЗРЯДЫ", "ПРОПУЩЕНО", "МЕСТА", "ЖУРНАЛ"}
	if len(walkSections) != len(want) {
		t.Fatalf("разделов %d, ждали %d", len(walkSections), len(want))
	}
	for i, x := range want {
		if got := walkSections[i].title(lang.RU); got != x {
			t.Errorf("раздел %d = %q, ждали %q", i, got, x)
		}
		// Every heading is a heading in both languages, and the English one
		// is no wider: the strip of eight has to fit the same terminal.
		en := walkSections[i].title(lang.EN)
		if en == x {
			t.Errorf("раздел %q не переведён", x)
		}
		if runes(en) > runes(x)+3 {
			t.Errorf("английский заголовок %q шире русского %q — полоса разделов разъедется", en, x)
		}
	}
	if len(walkSections) > 9 {
		t.Error("разделов больше девяти — клавиши 1…9 перестанут доставать до всех")
	}
}

// The screen does not remove anything and does not offer to: it names the
// command that does, which asks again before it acts.
func TestTheScreenOffersTheCommandAndNotTheDeed(t *testing.T) {
	w := walkTestScreen(t, 120)
	body := plain(strings.Join(w.removable(), "\n"))
	if !strings.Contains(body, "digitdisk clean ") {
		t.Errorf("раздел не назвал команду уборки:\n%s", body)
	}
	if !strings.Contains(body, "--apply") {
		t.Errorf("раздел не сказал, что без --apply ничего не трогается:\n%s", body)
	}
	for _, k := range []key{{kind: keyRune, r: 'd'}, {kind: keyRune, r: 'D'}, {kind: keyRune, r: 'x'}} {
		if w.handle(k) {
			t.Errorf("клавиша %q что-то сделала с экраном", k.r)
		}
	}
}

func TestTailKeepsTheEndOfAPath(t *testing.T) {
	// A path that does not fit keeps its END: the last components say where
	// the walk is, the first only repeat the root it was given.
	if got := tail("/очень/длинный/путь/сюда", 10); !strings.HasSuffix(got, "сюда") {
		t.Errorf("tail = %q, конец пути потерян", got)
	}
	if got := runes(tail("/очень/длинный/путь/сюда", 10)); got != 10 {
		t.Errorf("tail дал %d ячеек вместо 10", got)
	}
	if got := tail("коротко", 20); got != "коротко" {
		t.Errorf("короткий путь обрезан: %q", got)
	}
}

// A share of a directory is not a warning: the bar is drawn in the brand
// accent, not in the traffic-light colours a full disk gets.
func TestShareBarIsExactlyAsWideAsAsked(t *testing.T) {
	th := Theme{P: Carbon, d: depthTrue}
	for _, n := range []int{1, 8, 20, 40} {
		for _, f := range []float64{-1, 0, 0.3, 1, 9} {
			if got := plainWidth(th.shareBar(f, n)); got != n {
				t.Errorf("shareBar(%v,%d) шириной %d", f, n, got)
			}
		}
	}
}

// A tree bigger than the cap is still counted whole; only the walking-around
// stops, and the screen says so.
func TestTheCapKeepsTheTotalsRight(t *testing.T) {
	root := t.TempDir()
	// Ten directories deep, one file at the bottom of each level.
	p := root
	for i := 0; i < 10; i++ {
		p = filepath.Join(p, fmt.Sprintf("уровень%d", i))
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(p, "файл"), make([]byte, 1000), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	f := newWalkFeed(root)
	f.count = maxNodes // as if the cap were already reached
	res, err := scan.Walk(scan.Options{Root: root, Top: 1, Now: time.Now(), Watch: f.step})
	if err != nil {
		t.Fatal(err)
	}
	f.settle()
	if !f.truncated {
		t.Error("предел не отмечен")
	}
	if f.bytes != res.TotalBytes || f.tree.bytes != res.TotalBytes {
		t.Errorf("за пределом счёт разошёлся: %d / %d против %d", f.bytes, f.tree.bytes, res.TotalBytes)
	}
	if got := f.snapshot(); !got.Truncated {
		t.Error("снимок не сказал, что дерево усечено")
	}
}

var _ = core.KindFile

// --- the screen as a place to act from -------------------------------------

// cacheDecider stands in for the real decision layer in the tests of this
// package.  It is not a copy of the rules and does not pretend to be: the
// screen is being tested, not the core, and what the screen must do is obey
// whatever the layer says.  This one says «Кэш / МожноУбрать» for an old file
// under .cache and «НеТрогать» for everything else — enough to check that the
// screen removes exactly the first kind and never the second.
type cacheDecider struct{}

func (cacheDecider) Name() string { return "проверочный слой" }
func (cacheDecider) Ready() bool  { return true }
func (cacheDecider) Decide(r core.Record) core.Decision {
	if r.Kind == core.KindFile && r.AgeDays >= 7 && strings.Contains(r.Path, "/.cache/") {
		return core.Decision{Class: core.ClassCache, Verdict: core.VerdictRemovable, Weight: 1}
	}
	return core.Decision{Class: core.ClassUnknown, Verdict: core.VerdictKeep}
}

// workScreen is a browsing screen over a tree with something removable in it,
// wired to the same calls the subcommands make.
func workScreen(t *testing.T, cols int) (*walkScreen, string) {
	t.Helper()
	root := t.TempDir()
	old := time.Now().Add(-400 * 24 * time.Hour)
	mk := func(rel string, n int, aged bool) {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, make([]byte, n), 0o644); err != nil {
			t.Fatal(err)
		}
		if aged {
			if err := os.Chtimes(p, old, old); err != nil {
				t.Fatal(err)
			}
		}
	}
	mk(".cache/app/один.bin", 40000, true)
	mk(".cache/app/два.bin", 30000, true)
	mk("проекты/исходник.go", 9000, false)

	res, f := walkWith(t, root)
	w := &walkScreen{t: Theme{P: Carbon, d: depthTrue}, l: lang.RU, rows: 30, cols: cols,
		started: time.Now(), jobs: make(chan jobResult, 1)}
	w.o.Root, w.walkRoot = root, root
	w.o.Plan = func(r string, only []string) (*clean.Plan, error) {
		p, err := clean.Make(clean.Options{Root: r, Decider: cacheDecider{}, Now: time.Now(), Only: only})
		if err != nil {
			return nil, err
		}
		return &p, nil
	}
	w.o.Apply = func(p *clean.Plan) (*clean.Journal, error) {
		return clean.Apply(*p, clean.Options{Now: time.Now()})
	}
	w.o.History = func(r string) (*clean.History, error) { return clean.ReadHistory(r) }
	w.o.Restore = func(box string, dry bool) (*clean.Journal, error) {
		j, err := clean.ReadJournal(box)
		if err != nil {
			return nil, err
		}
		return clean.Restore(j, dry, time.Now())
	}
	w.finish(walkDone{res: res, tree: f.tree, snap: f.snapshot()})
	w.tab = 1 // ДЕРЕВО
	return w, root
}

// settle runs the one job the screen has out and takes its answer, so a test
// reads the same state a reader would see a moment later.
func (w *walkScreen) settleJob(t *testing.T) {
	t.Helper()
	select {
	case r := <-w.jobs:
		w.accepted(r)
	case <-time.After(20 * time.Second):
		t.Fatal("работа не вернулась")
	}
}

func TestMarkingNarrowsThePlanAndNeverWidensIt(t *testing.T) {
	w, _ := workScreen(t, 100)
	// Standing on «проекты», where the layer marked nothing.
	rows := w.here()
	for i, r := range rows {
		if strings.HasPrefix(r.name, "проекты") {
			w.sel[0] = i
		}
	}
	w.toggleMark()
	if len(w.marks) != 1 {
		t.Fatalf("отметилось %d каталогов", len(w.marks))
	}
	w.proposeClean()
	w.settleJob(t)
	body := plain(strings.Join(w.askLines(), "\n"))
	if !strings.Contains(body, "0 файлов") {
		t.Errorf("план по каталогу без «МожноУбрать» не пуст:\n%s", body)
	}
	if !strings.Contains(body, "отметка не делает файл убираемым") {
		t.Errorf("экран не сказал, чей это приговор:\n%s", body)
	}
	if w.ask.kind == overlayPlan {
		t.Error("пустой план всё-таки спрашивает подтверждение")
	}
}

func TestThePlanIsShownAndTheNumberIsAsked(t *testing.T) {
	w, root := workScreen(t, 100)
	w.proposeClean() // nothing marked: the whole tree
	w.settleJob(t)
	if w.ask == nil || w.ask.kind != overlayPlan {
		t.Fatalf("плана нет: %+v", w.ask)
	}
	if w.ask.want != 2 {
		t.Fatalf("к переносу %d файлов, ждали 2", w.ask.want)
	}
	body := plain(strings.Join(w.askLines(), "\n"))
	for _, want := range []string{"2 файлов", "по разрядам", "Кэш", "место не освободится", "purge"} {
		if !strings.Contains(body, want) {
			t.Errorf("в плане нет %q:\n%s", want, body)
		}
	}

	// A wrong number moves nothing and says so.
	w.askKey(key{kind: keyRune, r: '9'})
	w.askKey(key{kind: keyEnter})
	if w.ask.kind != overlayPlan {
		t.Fatal("неверное число всё-таки что-то сделало")
	}
	if !strings.Contains(plain(strings.Join(w.askLines(), "\n")), "ничего не тронуто") {
		t.Error("отказ не сказан вслух")
	}
	if _, err := os.Stat(filepath.Join(root, ".cache", "app", "один.bin")); err != nil {
		t.Fatalf("файл тронут после НЕВЕРНОГО числа: %v", err)
	}
	// A letter is not a number and is not typed into the answer at all.
	w.askKey(key{kind: keyRune, r: 'ф'})
	if w.ask.input != "" {
		t.Errorf("в поле числа попала буква: %q", w.ask.input)
	}

	// The right number moves them.
	w.askKey(key{kind: keyRune, r: '2'})
	w.askKey(key{kind: keyEnter})
	w.settleJob(t)
	if _, err := os.Stat(filepath.Join(root, ".cache", "app", "один.bin")); !os.IsNotExist(err) {
		t.Errorf("после подтверждения файл на месте: %v", err)
	}
	done := plain(strings.Join(w.askLines(), "\n"))
	for _, want := range []string{"перенесено", "место НЕ освобождено", "purge"} {
		if !strings.Contains(done, want) {
			t.Errorf("итог уборки не сказал %q:\n%s", want, done)
		}
	}
	if len(w.marks) != 0 {
		t.Error("отметки пережили уборку")
	}
}

func TestRestoreFromTheScreenAsksTheSameWay(t *testing.T) {
	w, root := workScreen(t, 100)
	w.proposeClean()
	w.settleJob(t)
	w.askKey(key{kind: keyRune, r: '2'})
	w.askKey(key{kind: keyEnter})
	w.settleJob(t)
	w.ask = nil

	// ЖУРНАЛ knows the корзина that has just appeared.
	w.tab = len(walkSections) - 1
	body := plain(strings.Join(w.journalSection(), "\n"))
	if !strings.Contains(body, "КОРЗИНЫ") {
		t.Fatalf("корзина не показана:\n%s", body)
	}
	if !strings.Contains(body, "место не освобождено, пока не purge") {
		t.Errorf("журнал не сказал, что место не освобождено:\n%s", body)
	}
	w.restoreSelected()
	w.settleJob(t)
	if w.ask == nil || w.ask.kind != overlayBack {
		t.Fatalf("возврат не спросил: %+v", w.ask)
	}
	if w.ask.want != 2 {
		t.Errorf("возврат насчитал %d файлов", w.ask.want)
	}
	w.askKey(key{kind: keyRune, r: '1'})
	w.askKey(key{kind: keyEnter})
	if w.ask.kind != overlayBack || !strings.Contains(w.ask.err, "ничего не тронуто") {
		t.Fatalf("неверное число при возврате не отказало: %+v", w.ask)
	}
	w.askKey(key{kind: keyRune, r: '2'})
	w.askKey(key{kind: keyEnter})
	w.settleJob(t)
	if _, err := os.Stat(filepath.Join(root, ".cache", "app", "один.bin")); err != nil {
		t.Errorf("файл не вернулся: %v", err)
	}
}

// The screen has no key that removes without a plan and a confirmation.  This
// walks every key it answers and checks the disk is untouched.
func TestNoKeyRemovesByItself(t *testing.T) {
	w, root := workScreen(t, 100)
	before := treeList(t, root)
	for _, k := range []key{
		{kind: keyRune, r: 'd'}, {kind: keyRune, r: 'D'}, {kind: keyRune, r: 'x'},
		{kind: keyRune, r: 'у'}, {kind: keyRune, r: ' '}, {kind: keyRune, r: '.'},
		{kind: keyEnter}, {kind: keyBack}, {kind: keyDown}, {kind: keyUp},
		{kind: keyRight}, {kind: keyLeft}, {kind: keyTab}, {kind: keyShiftTab},
	} {
		w.handle(k)
		if w.ask != nil && w.ask.kind == overlayPlan {
			t.Fatalf("клавиша %+v сразу дала подтверждаемый план", k)
		}
	}
	if got := treeList(t, root); got != before {
		t.Errorf("клавиши изменили дерево:\n%s\n---\n%s", before, got)
	}
	// And «c» stops at the plan: it asks, it does not act.
	w.tab = 1
	w.proposeClean()
	w.settleJob(t)
	if w.ask.kind != overlayPlan {
		t.Fatalf("«c» не остановилась на плане: %+v", w.ask)
	}
	if got := treeList(t, root); got != before {
		t.Errorf("построение плана изменило дерево:\n%s\n---\n%s", before, got)
	}
}

func treeList(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		fmt.Fprintf(&b, "%s %d\n", rel, info.Size())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestPathPromptCompletesDirectories(t *testing.T) {
	root := t.TempDir()
	for _, d := range []string{"дерево", "деревня", "другое"} {
		if err := os.Mkdir(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "дерево.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	w := &walkScreen{t: Theme{P: Carbon, d: depthTrue}, l: lang.RU, rows: 30, cols: 100}
	w.ask = &ask{kind: overlayPath, input: filepath.Join(root, "дере")}
	w.complete(true)
	// Two directories share the head «дерев», and a file is not an answer.
	if !strings.HasSuffix(w.ask.input, "дерев") {
		t.Errorf("дополнение дало %q", w.ask.input)
	}
	if len(w.ask.choices) != 2 {
		t.Errorf("предложено %v, ждали два каталога", w.ask.choices)
	}
	w.ask.input = filepath.Join(root, "деревн")
	w.complete(true)
	if !strings.HasSuffix(w.ask.input, "деревня"+string(filepath.Separator)) {
		t.Errorf("единственный каталог не дополнен целиком: %q", w.ask.input)
	}
	// A path that is not a directory is refused rather than walked.
	w.ask.input = filepath.Join(root, "дерево.txt")
	if w.askKey(key{kind: keyEnter}) {
		t.Fatal("Enter закрыл экран")
	}
	if w.ask == nil || w.ask.err == "" {
		t.Error("файл принят как корень обхода")
	}
}

// Приглашение пути предлагает ТЕКУЩИЙ каталог, когда ничего ещё не обойдено,
// и корень прошлого обхода, когда есть что обходить рядом. И то и другое
// проверяется здесь, потому что «предложить» — это единственное, чем экран
// отвечает на выбор `analyze` в списке команд.
func TestPathPromptOffersTheCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "внутри"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	w := &walkScreen{t: Theme{P: Carbon, d: depthTrue}, l: lang.RU, rows: 30, cols: 100}
	w.openPath()
	if w.ask == nil || w.ask.kind != overlayPath {
		t.Fatal("приглашение пути не открылось")
	}
	// Точный путь сравнивается по Clean: временный каталог может быть
	// символической ссылкой (/tmp → /private/tmp на macOS), и os.Getwd
	// вернёт разрешённый, а t.TempDir — исходный.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if got := filepath.Clean(w.ask.input); got != wd {
		t.Errorf("предложен %q, а текущий каталог %q", got, wd)
	}
	if !strings.HasSuffix(w.ask.input, string(filepath.Separator)) {
		t.Errorf("предложенный путь без разделителя на конце: %q — дополнение начнёт с имени каталога", w.ask.input)
	}
	if len(w.ask.choices) != 1 {
		t.Errorf("подкаталоги не предложены: %v", w.ask.choices)
	}
	// Согласие в одно нажатие обязано быть сказано вслух: обход домашнего
	// каталога — это миллионы записей, и Enter не должен быть тихим.
	said := strings.Join(w.ask.lines, " ")
	for _, want := range []string{"текущий каталог", "q прерывает обход"} {
		if !strings.Contains(said, want) {
			t.Errorf("приглашение не сказало %q: %q", want, said)
		}
	}

	// Обойдено что-то — предлагается сосед, а не текущий каталог, и
	// предупреждение снимается: путь выбран не по умолчанию.
	w.o.Root = filepath.Join(root, "внутри")
	w.openPath()
	if got := filepath.Clean(w.ask.input); got != filepath.Join(root, "внутри") {
		t.Errorf("после обхода предложено %q, ждали корень прошлого обхода", got)
	}
	if len(w.ask.lines) != 0 {
		t.Errorf("предупреждение о цене осталось на выбранном пути: %v", w.ask.lines)
	}
}

// Обход, который больше некому читать, отпускается.
//
// Пока экран был последним, что делает процесс, это не значило ничего: за
// возвратом из RunWalk шёл выход. Теперь за ним идёт экран состояния, и
// брошенный обход читал бы диск под ним ещё минуту. Замер на этой машине,
// одна и та же последовательность нажатий, обход /home/b (4 814 507 записей):
// 12,95 с пользовательского времени с отпусканием против 32,24 с без него при
// одинаковых 17,04 с по часам.
func TestWalkIsLetGoWhenNobodyIsWatching(t *testing.T) {
	let := make(chan struct{})
	w := &walkScreen{t: Theme{P: Carbon, d: depthTrue}, l: lang.RU, rows: 30, cols: 100}
	w.o.Walk = func(root string, watch func(scan.Step), stop func() bool) (scan.Result, error) {
		for !stop() {
			time.Sleep(time.Millisecond)
		}
		close(let)
		return scan.Result{}, nil
	}
	w.startWalk(t.TempDir())
	select {
	case <-let:
		t.Fatal("обход отпустили раньше, чем его бросили")
	case <-time.After(20 * time.Millisecond):
	}
	w.release()
	select {
	case <-let:
	case <-time.After(2 * time.Second):
		t.Fatal("обход не отпустили: он читает диск под экраном, который его не показывает")
	}
	// Новый обход отпускает прежний: два чтения диска на один экран — это
	// вдвое больше работы и вдвое меньше скорости у того, что видно.
	first := make(chan struct{})
	w.o.Walk = func(root string, watch func(scan.Step), stop func() bool) (scan.Result, error) {
		for !stop() {
			time.Sleep(time.Millisecond)
		}
		select {
		case <-first:
		default:
			close(first)
		}
		return scan.Result{}, nil
	}
	w.startWalk(t.TempDir())
	w.startWalk(t.TempDir())
	select {
	case <-first:
	case <-time.After(2 * time.Second):
		t.Fatal("прежний обход остался читать диск после начала нового")
	}
	w.release()
}

func TestBusyScreenAnswersNothingButCtrlC(t *testing.T) {
	w := &walkScreen{t: Theme{P: Carbon, d: depthTrue}, l: lang.RU, rows: 30, cols: 100, jobs: make(chan jobResult, 1)}
	w.busy = true
	w.ask = &ask{kind: overlayNote, title: "ИДЁТ"}
	for _, k := range []key{{kind: keyEnter}, {kind: keyRune, r: 'q'}, {kind: keyEsc}} {
		if w.handle(k) {
			t.Errorf("клавиша %+v закрыла экран во время работы", k)
		}
		if w.ask == nil {
			t.Fatalf("клавиша %+v сняла извещение о работе", k)
		}
	}
	if !w.handle(key{kind: keyCtrlC}) {
		t.Error("Ctrl-C не прервал")
	}
}

// Every question the screen asks has to fit the terminal too — a plan that
// wraps is a plan that cannot be read before it is confirmed.
func TestEveryQuestionDrawsWithinTheTerminal(t *testing.T) {
	for _, cols := range []int{40, 60, 80, 120, 200} {
		w, _ := workScreen(t, cols)
		asks := []func(){
			func() { w.openPath() },
			func() {
				w.tell("ОТВЕТ", []string{"строка раз", "строка два, длинная-предлинная, чтобы точно не влезла в сорок колонок"})
			},
			func() { w.ask = &ask{kind: overlayKeys, title: "КЛАВИШИ И КОМАНДЫ"} },
			func() {
				w.proposeClean()
				w.settleJob(t)
			},
		}
		for i, open := range asks {
			w.ask = nil
			open()
			if w.ask == nil {
				t.Fatalf("вопрос %d не открылся", i)
			}
			for _, line := range w.frame() {
				if got := plainWidth(w.t.clip(line, cols)); got > cols {
					t.Errorf("вопрос %d / %d колонок: строка в %d ячеек", i, cols, got)
				}
			}
		}
		w.ask = nil
	}
}
