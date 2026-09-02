// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"digitdisk/internal/core"
	"digitdisk/internal/protect"
)

func guarded(t *testing.T, root string, args []string, removable ...string) Plan {
	t.Helper()
	set := map[string]bool{}
	for _, r := range removable {
		set[r] = true
	}
	list, err := protect.Load(protect.Options{
		Args: args, Home: root, Config: t.TempDir(), Getenv: func(string) string { return "" },
	})
	if err != nil {
		t.Fatal(err)
	}
	p, err := Make(Options{
		Root:    root,
		Decider: judge{removable: set, classOf: map[string]core.Class{"old.bin": core.ClassCache, "older.bin": core.ClassCache}},
		Now:     time.Now(),
		Version: "испытание",
		Protect: list,
	})
	if err != nil {
		t.Fatalf("план не построился: %v", err)
	}
	return p
}

// TestProtectedPathLeavesThePlan is the whole point of the защитный список: a
// path the layer approved and the operator forbade never reaches the work
// list, and says so out loud instead of going missing.
func TestProtectedPathLeavesThePlan(t *testing.T) {
	root := tree(t)
	p := guarded(t, root, []string{filepath.Join(root, "cache", "deep")}, "old.bin", "older.bin")

	if len(p.Items) != 1 || filepath.Base(p.Items[0].Path) != "old.bin" {
		t.Fatalf("к переносу %d файлов, ожидался один — old.bin", len(p.Items))
	}
	if len(p.Protected) != 1 {
		t.Fatalf("защищено %d записей, ожидалась одна", len(p.Protected))
	}
	got := p.Protected[0]
	if filepath.Base(got.Path) != "older.bin" {
		t.Errorf("защищено %q, ожидался older.bin", got.Path)
	}
	if got.Rule.Kind != protect.KindPath {
		t.Errorf("правило %+v — ожидалось правило по пути", got.Rule)
	}
	if p.ProtectedBytes != got.Size || got.Size == 0 {
		t.Errorf("байты защищённого не сосчитаны: %d при размере %d", p.ProtectedBytes, got.Size)
	}
	if len(p.Refused) != 0 {
		t.Errorf("защита попала в отказы: %+v — это разные вещи", p.Refused)
	}
}

// TestProtectedClassLeavesThePlan: «разряд:кэш» empties a plan made entirely
// of caches, and the ledger says why.
func TestProtectedClassLeavesThePlan(t *testing.T) {
	root := tree(t)
	p := guarded(t, root, []string{"разряд:Кэш"}, "old.bin", "older.bin")
	if len(p.Items) != 0 {
		t.Fatalf("к переносу %d файлов, ожидался пустой план", len(p.Items))
	}
	if len(p.Protected) != 2 {
		t.Fatalf("защищено %d записей, ожидались две", len(p.Protected))
	}
	for _, pr := range p.Protected {
		if pr.Rule.Kind != protect.KindClass {
			t.Errorf("правило %+v — ожидалось правило по разряду", pr.Rule)
		}
	}
}

// TestProtectionOnlySubtracts is the property that lets the защитный список
// live in the host instead of in the proved rules: it can never add a file to
// the plan, only take one away.
func TestProtectionOnlySubtracts(t *testing.T) {
	root := tree(t)
	bare := plan(t, root, "old.bin", "older.bin")
	for _, args := range [][]string{
		{"разряд:Кэш"},
		{filepath.Join(root, "cache")},
		{"cache"},
		{"/nowhere"},
	} {
		p := guarded(t, root, args, "old.bin", "older.bin")
		if len(p.Items) > len(bare.Items) {
			t.Errorf("защита %v ДОБАВИЛА файлов в план: %d против %d", args, len(p.Items), len(bare.Items))
		}
		if len(p.Items)+len(p.Protected) != len(bare.Items) {
			t.Errorf("защита %v потеряла записи: %d + %d != %d", args, len(p.Items), len(p.Protected), len(bare.Items))
		}
	}
}

// TestByClassCountsEverything: the breakdown is over the whole plan, so a
// report that shortens its lists still tells the truth about the disk.
func TestByClassCountsEverything(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin", "older.bin")
	if len(p.ByClass) != 1 {
		t.Fatalf("разрядов в сводке %d, ожидался один", len(p.ByClass))
	}
	sum := p.ByClass[0]
	if sum.Class != core.ClassCache || sum.Count != 2 || sum.Bytes != p.Bytes {
		t.Errorf("сводка %+v не сходится с планом (%d файлов, %d Б)", sum, len(p.Items), p.Bytes)
	}
}

// TestHistoryReadsTheJournals: history remembers nothing of its own — every
// number comes out of the корзины themselves.
func TestHistoryReadsTheJournals(t *testing.T) {
	root := tree(t)
	p := plan(t, root, "old.bin")
	j, err := Apply(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}

	h, err := ReadHistory(root)
	if err != nil {
		t.Fatal(err)
	}
	if h.Boxes != 1 {
		t.Fatalf("корзин %d, ожидалась одна", h.Boxes)
	}
	e := h.Entries[0]
	if e.Moved != 1 || e.MovedBytes == 0 {
		t.Errorf("в корзине %d файлов на %d Б — ожидался один непустой", e.Moved, e.MovedBytes)
	}
	if e.FreedBytes != 0 {
		t.Errorf("освобождено %d Б — перенос в корзину не освобождает ничего", e.FreedBytes)
	}
	if !e.Restorable() {
		t.Error("корзина объявлена невозвратной, а из неё ещё ничего не брали")
	}

	// The same корзина named directly is the same answer.
	one, err := ReadHistory(j.Box)
	if err != nil {
		t.Fatal(err)
	}
	if one.Boxes != 1 || one.Entries[0].Moved != e.Moved {
		t.Errorf("корзина, названная прямо, прочиталась иначе: %+v", one.Entries)
	}

	// After a purge the same journal reports freed bytes, and nothing else
	// changed its story.
	if _, err := Purge(j, 1, time.Now()); err != nil {
		t.Fatal(err)
	}
	h, err = ReadHistory(root)
	if err != nil {
		t.Fatal(err)
	}
	if h.Entries[0].FreedBytes != e.MovedBytes {
		t.Errorf("после стирания освобождено %d Б, ожидалось %d", h.Entries[0].FreedBytes, e.MovedBytes)
	}
	if h.Entries[0].Restorable() {
		t.Error("стёртое объявлено возвратным")
	}
}

// TestHistoryOfACleanTreeIsAnAnswer: no корзина is a fact, not an error.
func TestHistoryOfACleanTreeIsAnAnswer(t *testing.T) {
	root := t.TempDir()
	h, err := ReadHistory(root)
	if err != nil {
		t.Fatalf("чистое дерево дало ошибку: %v", err)
	}
	if h.Boxes != 0 || len(h.Entries) != 0 {
		t.Errorf("в чистом дереве нашлись корзины: %+v", h.Entries)
	}
	if _, err := os.Lstat(filepath.Join(root, TrashName)); !os.IsNotExist(err) {
		t.Error("history создала хранилище корзин — она обязана только читать")
	}
}
