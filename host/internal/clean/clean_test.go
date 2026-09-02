// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/core"
)

// judge is a stand-in for the decision layer.  It is deliberately blunter than
// the flang rules — it says «МожноУбрать» to whatever the test asks it to —
// because these tests are about the host: what it does with a verdict, not how
// the verdict is reached.  The rules themselves are proved in core/ against
// 12 000 inputs and three implementations, and repeating them here in Go would
// only mean a fourth implementation to disagree with.
type judge struct {
	removable map[string]bool
	classOf   map[string]core.Class
}

func (j judge) Decide(r core.Record) core.Decision {
	class := core.ClassUnknown
	if c, ok := j.classOf[filepath.Base(r.Path)]; ok {
		class = c
	}
	if j.removable[filepath.Base(r.Path)] {
		return core.Decision{Class: class, Verdict: core.VerdictRemovable, Weight: float64(r.Size)}
	}
	return core.Decision{Class: class, Verdict: core.VerdictKeep}
}

func (judge) Name() string { return "испытательный слой" }
func (judge) Ready() bool  { return true }
func (judge) Threshold(c core.Class) (float64, bool) {
	if c == core.ClassCache {
		return 7, true
	}
	return 0, false
}

// tree builds a small tree and returns its root.
func tree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	mk := func(rel, body string) {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mk("cache/old.bin", "мусор постарше")
	mk("cache/deep/older.bin", "мусор поглубже")
	mk("docs/letter.txt", "важное письмо")
	return root
}

func plan(t *testing.T, root string, removable ...string) Plan {
	t.Helper()
	set := map[string]bool{}
	for _, r := range removable {
		set[r] = true
	}
	p, err := Make(Options{
		Root:    root,
		Decider: judge{removable: set, classOf: map[string]core.Class{"old.bin": core.ClassCache, "older.bin": core.ClassCache}},
		Now:     time.Now(),
		Version: "испытание",
	})
	if err != nil {
		t.Fatalf("план не построился: %v", err)
	}
	return p
}

// TestPlanTouchesNothing is the property the whole design rests on: asking
// what would be removed must never remove anything, and must not even create
// the корзина.
func TestPlanTouchesNothing(t *testing.T) {
	root := tree(t)
	before := snapshot(t, root)

	p := plan(t, root, "old.bin", "older.bin")
	if len(p.Items) != 2 {
		t.Fatalf("к переносу %d файлов, ожидалось 2", len(p.Items))
	}
	if _, err := os.Lstat(p.Trash); !os.IsNotExist(err) {
		t.Fatalf("план создал корзину %s — план обязан быть вопросом без последствий", p.Trash)
	}
	if after := snapshot(t, root); after != before {
		t.Fatalf("план изменил дерево:\nбыло:\n%s\nстало:\n%s", before, after)
	}
}

// TestOnlyRemovableIsTaken: a file the layer did not mark stays put, however
// much its name looks like rubbish.
func TestOnlyRemovableIsTaken(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin") // older.bin is NOT marked
	if len(p.Items) != 1 || filepath.Base(p.Items[0].Path) != "old.bin" {
		t.Fatalf("в план попало %+v, ожидался ровно old.bin", p.Items)
	}
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := j.Moved(); n != 1 {
		t.Fatalf("перенесено %d, ожидался 1", n)
	}
	for _, keep := range []string{"cache/deep/older.bin", "docs/letter.txt"} {
		if _, err := os.Lstat(filepath.Join(root, keep)); err != nil {
			t.Fatalf("%s исчез, а его никто не помечал: %v", keep, err)
		}
	}
}

// TestApplyRestoreRoundTrip: what was moved comes back byte for byte, in the
// same place.  Reversibility is the promise; this is the proof of it.
func TestApplyRestoreRoundTrip(t *testing.T) {
	root := tree(t)
	before := snapshot(t, root)

	p := plan(t, root, "old.bin", "older.bin")
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := j.Moved(); n != 2 {
		t.Fatalf("перенесено %d, ожидалось 2", n)
	}
	for _, gone := range []string{"cache/old.bin", "cache/deep/older.bin"} {
		if _, err := os.Lstat(filepath.Join(root, gone)); !os.IsNotExist(err) {
			t.Fatalf("%s не перенесён", gone)
		}
	}
	if _, err := os.Stat(filepath.Join(j.Box, JournalName)); err != nil {
		t.Fatalf("журнала нет: %v", err)
	}

	back, err := ReadJournal(j.Box)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Restore(back, false, time.Now()); err != nil {
		t.Fatal(err)
	}
	if n, _ := back.Restored(); n != 2 {
		t.Fatalf("возвращено %d, ожидалось 2", n)
	}

	// The корзина itself is the only difference the round trip may leave.
	if after := snapshot(t, root); after != before {
		t.Fatalf("возврат не восстановил дерево:\nбыло:\n%s\nстало:\n%s", before, after)
	}
}

// TestRestoreRefusesToOverwrite: if something has taken the old name, the file
// stays in the корзина and the caller is told.
func TestRestoreRefusesToOverwrite(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin")
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "cache/old.bin"), []byte("новое содержимое"), 0o644); err != nil {
		t.Fatal(err)
	}
	back, err := ReadJournal(j.Box)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Restore(back, false, time.Now()); err != nil {
		t.Fatal(err)
	}
	if n, _ := back.Restored(); n != 0 {
		t.Fatalf("возвращено %d поверх занятого места, ожидалось 0", n)
	}
	body, err := os.ReadFile(filepath.Join(root, "cache/old.bin"))
	if err != nil || string(body) != "новое содержимое" {
		t.Fatalf("возврат затёр то, что появилось на прежнем месте: %q, %v", body, err)
	}
}

// TestChangedFileIsRefused is falsifier Ф-«отказ вместо угадывания»: a file
// written to between the plan and the move is not moved.
func TestChangedFileIsRefused(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")

	// Someone writes to one of them after the walk and before the move.
	target := filepath.Join(root, "cache/old.bin")
	if err := os.WriteFile(target, []byte("это уже другое содержимое"), 0o644); err != nil {
		t.Fatal(err)
	}

	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := j.Moved(); n != 1 {
		t.Fatalf("перенесено %d, ожидался 1 (второй обязан быть отвергнут)", n)
	}
	failed := j.Failed()
	if len(failed) != 1 || filepath.Base(failed[0].Path) != "old.bin" {
		t.Fatalf("отвергнуто %+v, ожидался ровно old.bin", failed)
	}
	if !strings.Contains(failed[0].Failed.String(), "размер") && !strings.Contains(failed[0].Failed.String(), "писали") {
		t.Fatalf("отказ не называет, что изменилось: %q", failed[0].Failed)
	}
	if _, err := os.Lstat(target); err != nil {
		t.Fatalf("изменённый файл всё-таки убран: %v", err)
	}
}

// TestVanishedFileIsRefused: a file removed by someone else between the walk
// and the move is reported, not silently counted as cleaned.
func TestVanishedFileIsRefused(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin")
	if err := os.Remove(filepath.Join(root, "cache/old.bin")); err != nil {
		t.Fatal(err)
	}
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := j.Moved(); n != 0 {
		t.Fatalf("перенесено %d исчезнувших файлов", n)
	}
	failed := j.Failed()
	if len(failed) != 1 || !strings.Contains(failed[0].Failed.String(), "исчез") {
		t.Fatalf("исчезновение не названо: %+v", failed)
	}
}

// TestGuardRefusesDirectoryAndLink: even if the layer were to say
// «МожноУбрать» to a directory or a symlink — which the flang rules П2 and П3
// forbid — the host refuses and reports the disagreement.
func TestGuardRefusesDirectoryAndLink(t *testing.T) {
	root := tree(t)
	if err := os.Symlink(filepath.Join(root, "docs/letter.txt"), filepath.Join(root, "cache/link")); err != nil {
		t.Skipf("символические ссылки недоступны: %v", err)
	}
	p := plan(t, root, "deep", "link") // "deep" is a directory, "link" a symlink
	if len(p.Items) != 0 {
		t.Fatalf("в план попало %+v, ожидалась пустота", p.Items)
	}
	if len(p.Refused) != 2 {
		t.Fatalf("отказов %d, ожидалось 2: %+v", len(p.Refused), p.Refused)
	}
	var sawDir, sawLink bool
	for _, r := range p.Refused {
		if strings.Contains(r.Reason.String(), "каталог") {
			sawDir = true
		}
		if strings.Contains(r.Reason.String(), "ссылка") {
			sawLink = true
		}
	}
	if !sawDir || !sawLink {
		t.Fatalf("отказы не названы по существу: %+v", p.Refused)
	}
	if _, err := os.Lstat(filepath.Join(root, "cache/deep")); err != nil {
		t.Fatalf("каталог всё-таки тронут: %v", err)
	}
}

// TestTrashOutsideRootIsRefused: the корзина may not leave the корень, and the
// refusal names the reason — the price of a cross-filesystem copy.
func TestTrashOutsideRootIsRefused(t *testing.T) {
	root := tree(t)
	outside := t.TempDir()
	_, err := Make(Options{Root: root, Trash: outside, Decider: judge{}, Now: time.Now()})
	if err == nil {
		t.Fatal("корзина вне корня принята")
	}
	if !strings.Contains(err.Error(), "вне корня") {
		t.Fatalf("отказ не объясняет себя: %v", err)
	}
}

// TestTrashIsNotWalked: a second run must not find the first run's корзина and
// clean it into a корзина inside the корзина.
func TestTrashIsNotWalked(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")
	if _, err := Apply(p, Options{Now: time.Now(), Version: "испытание"}); err != nil {
		t.Fatal(err)
	}
	again := plan(t, root, "old.bin", "older.bin")
	if len(again.Items) != 0 {
		t.Fatalf("второй проход нашёл в корзине %+v", again.Items)
	}
	if again.PrunedTrash != 1 {
		t.Fatalf("корзина пропущена %d раз, ожидался 1", again.PrunedTrash)
	}
}

// TestPurgeNeedsTheRightCount: the confirmation is the count, and a wrong
// count erases nothing.
func TestPurgeNeedsTheRightCount(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}

	for _, wrong := range []int{0, 1, 3} {
		box, err := ReadJournal(j.Box)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Purge(box, wrong, time.Now()); err == nil {
			t.Fatalf("--confirm %d принят, а в корзине 2 файла", wrong)
		}
	}
	// Nothing was erased by any of those attempts.
	box, err := ReadJournal(j.Box)
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := box.Moved(); n != 2 {
		t.Fatalf("в корзине осталось %d файлов после неудачных подтверждений", n)
	}

	if _, err := Purge(box, 2, time.Now()); err != nil {
		t.Fatal(err)
	}
	if n, _ := box.Purged(); n != 2 {
		t.Fatalf("стёрто %d, ожидалось 2", n)
	}
	for _, it := range box.Items {
		if _, err := os.Lstat(filepath.Join(root, it.TrashRel)); !os.IsNotExist(err) {
			t.Fatalf("%s не стёрт", it.TrashRel)
		}
	}
	// The journal survives the erase: it is the record of what is gone.
	if _, err := os.Stat(filepath.Join(j.Box, JournalName)); err != nil {
		t.Fatalf("журнал не пережил стирание: %v", err)
	}
}

// TestJournalIsWrittenBeforeTheMove: a crash between the first file moving and
// the last must still leave a корзина that can be emptied back.
func TestJournalIsWrittenBeforeTheMove(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := ReadJournal(j.Box)
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range loaded.Items {
		if it.TrashRel == "" {
			t.Fatalf("в журнале нет места в корзине для %s — после обрыва его не найти", it.Path)
		}
		if it.Path == "" || it.Rel == "" {
			t.Fatalf("в журнале нет прежнего пути: %+v", it)
		}
	}
}

// TestMovedJournalIsRefused: a корзина copied elsewhere would restore into the
// original tree, which is not what the person holding the copy meant.
func TestMovedJournalIsRefused(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin")
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	elsewhere := filepath.Join(t.TempDir(), "копия")
	if err := os.MkdirAll(elsewhere, 0o755); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(j.Box, JournalName))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(elsewhere, JournalName), body, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadJournal(elsewhere); err == nil {
		t.Fatal("журнал перемещённой корзины принят")
	}
}

// TestStubDeciderCannotApply: a host built without the real decision layer
// must say so rather than quietly clean nothing and call it success.
func TestStubDeciderCannotApply(t *testing.T) {
	root := tree(t)
	p, err := Make(Options{Root: root, Decider: core.Default(), Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(p, Options{Now: time.Now()}); err == nil {
		t.Fatal("уборка с заглушкой принята")
	}
}

// snapshot renders a tree as a stable string: path and contents, with the
// корзина left out.  An empty корзина left behind is the one trace a round
// trip is allowed to leave, and its journal is a record of the round trip
// rather than a part of the tree.
func snapshot(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		for _, part := range strings.Split(rel, string(filepath.Separator)) {
			if part == TrashName {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.Mode().IsRegular() {
			body, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			b.WriteString(rel + " " + string(body) + "\n")
		} else if info.IsDir() {
			b.WriteString(rel + "/\n")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// --- стирание насовсем ------------------------------------------------------

// Стирание берёт РОВНО план и ничего больше: помеченный файл исчезает,
// непомеченный рядом — нет, и корзины после этого нет ни одной.
func TestEraseTakesExactlyThePlanAndNothingElse(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin") // older.bin и letter.txt НЕ помечены
	j, err := Erase(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := j.Purged(); n != 1 {
		t.Fatalf("стёрто %d, ожидался 1", n)
	}
	if _, err := os.Lstat(filepath.Join(root, "cache/old.bin")); !os.IsNotExist(err) {
		t.Errorf("помеченный файл пережил стирание: %v", err)
	}
	for _, keep := range []string{"cache/deep/older.bin", "docs/letter.txt"} {
		if _, err := os.Lstat(filepath.Join(root, keep)); err != nil {
			t.Errorf("%s исчез, а его никто не помечал: %v", keep, err)
		}
	}
	// Ничего не отложено «на потом»: в корзине лежит журнал и только он.
	files, err := os.ReadDir(filepath.Join(j.Box))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Name() != JournalName {
		var names []string
		for _, f := range files {
			names = append(names, f.Name())
		}
		t.Errorf("рядом с журналом стирания лежит ещё что-то: %v", names)
	}
	if n, _ := j.Moved(); n != 0 {
		t.Errorf("журнал стирания насчитал %d перенесённых", n)
	}
}

// Журнал стирания отличим от журнала уборки — и отличим не по счётчикам, а по
// собственному слову. Возврат и `purge` по нему отказывают вслух.
func TestAnEraseJournalIsNeverMistakenForATrash(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")
	j, err := Erase(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if j.Way != WayErase || !j.Erasure() {
		t.Fatalf("журнал не назвал себя стиранием: способ %q", j.Way)
	}

	// Прочитанный с диска — тот же самый.
	back, err := ReadJournal(j.Box)
	if err != nil {
		t.Fatal(err)
	}
	if !back.Erasure() {
		t.Fatalf("с диска журнал прочитался как корзина: способ %q", back.Way)
	}
	if _, err := Restore(back, true, time.Now()); err == nil {
		t.Error("возврат согласился работать по журналу стирания")
	} else if !strings.Contains(err.Error(), "журнал стирания") {
		t.Errorf("возврат отказал не по той причине: %v", err)
	}
	if _, err := Purge(back, 2, time.Now()); err == nil {
		t.Error("purge согласился стирать по журналу стирания")
	} else if !strings.Contains(err.Error(), "журнал стирания") {
		t.Errorf("purge отказал не по той причине: %v", err)
	}

	// А история видит стирание как стирание и считает место освобождённым.
	h, err := ReadHistory(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Entries) != 1 {
		t.Fatalf("корзин в истории %d", len(h.Entries))
	}
	e := h.Entries[0]
	if e.Way != WayErase {
		t.Errorf("история не назвала способ: %q", e.Way)
	}
	if e.Moved != 0 || e.Purged != 2 {
		t.Errorf("история насчитала в корзине %d и стёртых %d", e.Moved, e.Purged)
	}
	if e.FreedBytes != e.PurgedBytes || e.FreedBytes == 0 {
		t.Errorf("освобождено %d при стёртых %d байт", e.FreedBytes, e.PurgedBytes)
	}
	if e.Restorable() {
		t.Error("история считает стёртое возвратимым")
	}
}

// Файл, который изменился между обходом и стиранием, не стирается: отпечаток
// сверяется так же, как перед переносом.
func TestEraseRefusesAFileThatChangedSinceTheWalk(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")
	changed := filepath.Join(root, "cache/old.bin")
	if err := os.WriteFile(changed, []byte("совсем другое содержимое"), 0o644); err != nil {
		t.Fatal(err)
	}
	j, err := Erase(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(changed); err != nil {
		t.Errorf("изменившийся файл всё-таки стёрт: %v", err)
	}
	if n, _ := j.Purged(); n != 1 {
		t.Errorf("стёрто %d, ожидался ровно один — второй изменился", n)
	}
	if len(j.Failed()) != 1 {
		t.Fatalf("отказ не записан в журнал: %d", len(j.Failed()))
	}
	if got := j.Failed()[0].Failed.String(); !strings.Contains(got, "не стёрт") {
		t.Errorf("отказ не сказал, что файл не стёрт: %q", got)
	}
}

// Пустой план не стирает ничего и говорит это отказом, а не тишиной.
func TestEraseRefusesAnEmptyPlan(t *testing.T) {
	root := tree(t)
	before := snapshot(t, root)
	p := plan(t, root) // ничего не помечено
	if _, err := Erase(p, Options{Now: time.Now()}); err == nil {
		t.Fatal("пустой план принят к стиранию")
	}
	if after := snapshot(t, root); after != before {
		t.Errorf("пустой план изменил дерево:\n%s\n---\n%s", before, after)
	}
}
