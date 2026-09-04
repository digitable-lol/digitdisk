// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/core"
)

// buildTree lays out a small tree with a known apparent size and returns the
// root plus the sum of st_size over every distinct inode in it, which is what
// du -sb reports.
func buildTree(t *testing.T) (root string, wantBytes int64, wantFiles, wantDirs, wantLinks int) {
	t.Helper()
	root = t.TempDir()

	write := func(rel string, n int) string {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, make([]byte, n), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	a := write("a.bin", 1000)
	write("sub/b.bin", 2500)
	write("sub/deep/c.bin", 7)

	// a hard link to a.bin: du charges those bytes once
	if err := os.Link(a, filepath.Join(root, "sub", "a-link.bin")); err != nil {
		t.Skipf("hard links unavailable here: %v", err)
	}
	// a symlink to a directory: must be recorded but never descended into
	if err := os.Symlink(filepath.Join(root, "sub"), filepath.Join(root, "sub-link")); err != nil {
		t.Fatal(err)
	}
	// a dangling symlink: must not be an error
	if err := os.Symlink(filepath.Join(root, "gone"), filepath.Join(root, "dangling")); err != nil {
		t.Fatal(err)
	}

	seen := map[uint64]bool{}
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		// bytes are charged once per inode (du -sb), but every name is
		// still an entry
		charge := true
		if ino := inoOf(info); nlink(info) > 1 {
			charge = !seen[ino]
			seen[ino] = true
		}
		// du charges a directory nothing for itself
		if charge && !info.IsDir() {
			wantBytes += info.Size()
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			wantLinks++
		case info.IsDir():
			wantDirs++
		default:
			wantFiles++
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return root, wantBytes, wantFiles, wantDirs, wantLinks
}

func TestWalkCountsAndSizes(t *testing.T) {
	root, wantBytes, wantFiles, wantDirs, wantLinks := buildTree(t)

	res, err := Walk(Options{Root: root, Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if res.TotalBytes != wantBytes {
		t.Errorf("TotalBytes = %d, want %d (du -sb semantics)", res.TotalBytes, wantBytes)
	}
	if res.Files != wantFiles || res.Dirs != wantDirs || res.Links != wantLinks {
		t.Errorf("files/dirs/links = %d/%d/%d, want %d/%d/%d",
			res.Files, res.Dirs, res.Links, wantFiles, wantDirs, wantLinks)
	}
	if res.HardlinkDupes != 1 {
		t.Errorf("HardlinkDupes = %d, want 1", res.HardlinkDupes)
	}
	if res.Entries != res.Files+res.Dirs+res.Links+res.Others {
		t.Errorf("Entries %d does not add up from the per-kind counters", res.Entries)
	}
	// the buckets must add up to the same total
	var sum int64
	for _, b := range res.ByClass {
		sum += b.Bytes
	}
	if sum != res.TotalBytes {
		t.Errorf("class buckets sum to %d, TotalBytes is %d", sum, res.TotalBytes)
	}
}

func TestWalkDoesNotFollowSymlinks(t *testing.T) {
	root, _, _, _, _ := buildTree(t)
	res, err := Walk(Options{Root: root, Now: time.Now(), Top: 100})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range res.Largest {
		if filepath.Base(filepath.Dir(e.Path)) == "sub-link" {
			t.Fatalf("walk descended through a symlink: %s", e.Path)
		}
	}
	if res.Skipped.OtherErrors != 0 || res.Skipped.PermissionDenied != 0 {
		t.Errorf("a dangling symlink must not produce an error: %+v", res.Skipped)
	}
}

func TestWalkSkipsUnreadableDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	root := t.TempDir()
	locked := filepath.Join(root, "locked")
	if err := os.Mkdir(locked, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(locked, "secret"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(locked, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) })

	res, err := Walk(Options{Root: root, Now: time.Now()})
	if err != nil {
		t.Fatalf("an unreadable subdirectory must not fail the walk: %v", err)
	}
	if res.Skipped.PermissionDenied != 1 {
		t.Errorf("PermissionDenied = %d, want 1 (skips: %+v)", res.Skipped.PermissionDenied, res.Skipped)
	}
	if res.Dirs != 2 {
		t.Errorf("Dirs = %d, want 2 (root and the locked directory itself)", res.Dirs)
	}
}

func TestWalkMaxDepth(t *testing.T) {
	root, _, _, _, _ := buildTree(t)
	res, err := Walk(Options{Root: root, Now: time.Now(), MaxDepth: 1})
	if err != nil {
		t.Fatal(err)
	}
	if res.Skipped.DepthLimited == 0 {
		t.Errorf("depth limit was never hit: %+v", res.Skipped)
	}
	for _, e := range res.Largest {
		if filepath.Base(e.Path) == "c.bin" {
			t.Errorf("depth limit did not stop the descent: %s", e.Path)
		}
	}
}

func TestWalkMissingRoot(t *testing.T) {
	if _, err := Walk(Options{Root: filepath.Join(t.TempDir(), "nope")}); err == nil {
		t.Errorf("a missing root must be reported as an error")
	}
}

// countingDecider records what the walk hands to the decision layer.
type countingDecider struct {
	seen  []core.Record
	ready bool
}

func (d *countingDecider) Decide(r core.Record) core.Decision {
	d.seen = append(d.seen, r)
	if r.Kind == core.KindFile && r.Size > 2000 {
		return core.Decision{Class: core.ClassLarge, Verdict: core.VerdictRemovable, Weight: 1}
	}
	return core.Decision{Class: core.ClassUnknown, Verdict: core.VerdictKeep}
}
func (d *countingDecider) Name() string { return "counting" }
func (d *countingDecider) Ready() bool  { return d.ready }

func TestWalkFeedsContractRecords(t *testing.T) {
	root, _, _, _, _ := buildTree(t)
	d := &countingDecider{ready: true}
	res, err := Walk(Options{Root: root, Now: time.Now(), Decider: d, Top: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(d.seen) != res.Entries {
		t.Errorf("decider saw %d records, walk counted %d entries", len(d.seen), res.Entries)
	}
	for _, r := range d.seen {
		switch r.Kind {
		case core.KindFile, core.KindDir, core.KindLink:
		default:
			t.Errorf("record %q carries вид %q, outside the contract", r.Path, r.Kind)
		}
		if r.AgeDays < 0 {
			t.Errorf("negative возраст_дней for %q", r.Path)
		}
		if !r.Accessible {
			t.Errorf("everything in this tree is readable, but %q came back недоступен", r.Path)
		}
	}
	if got := res.ByVerdict[core.VerdictRemovable].Count; got != 1 {
		t.Errorf("removable count = %d, want 1 (only sub/b.bin exceeds 2000 bytes)", got)
	}
	if len(res.Removable) != 1 || filepath.Base(res.Removable[0].Path) != "b.bin" {
		t.Errorf("removable list = %+v", res.Removable)
	}
}

func TestWalkSingleFileRoot(t *testing.T) {
	root := t.TempDir()
	p := filepath.Join(root, "only.bin")
	if err := os.WriteFile(p, make([]byte, 42), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Walk(Options{Root: p, Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if res.Entries != 1 || res.Files != 1 || res.TotalBytes != 42 {
		t.Errorf("single-file root = %+v", res)
	}
}

// The size rule written out as numbers rather than derived by walking the tree
// a second time: the total is the apparent size of every distinct inode, a
// directory is charged nothing for itself, a symlink is charged the length of
// its target, and a second name for one inode is charged nothing.
//
// It is here because the rule has to be the same on every system digitdisk
// runs on.  The walk is one piece of code on Linux and macOS, and this is the
// arithmetic it must produce on both — expected values that a run cannot talk
// its way out of.
func TestWalkSizeRuleInExactBytes(t *testing.T) {
	root := t.TempDir()
	write := func(rel string, n int) string {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, make([]byte, n), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	big := write("big.bin", 1000)
	write("sub/small.bin", 7)
	if err := os.Link(big, filepath.Join(root, "sub", "second-name.bin")); err != nil {
		t.Skipf("hard links unavailable here: %v", err)
	}
	// Two symlinks whose targets are of a length we choose, since a symlink
	// is charged the length of the string it holds.
	target := strings.Repeat("x", 30)
	if err := os.Symlink(target, filepath.Join(root, "link-a")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("y", filepath.Join(root, "link-b")); err != nil {
		t.Fatal(err)
	}

	const wantBytes = 1000 + 7 + 30 + 1
	res, err := Walk(Options{Root: root, Now: time.Now(), Top: 20})
	if err != nil {
		t.Fatal(err)
	}
	if res.TotalBytes != wantBytes {
		t.Errorf("TotalBytes = %d, want %d (1000 + 7 + два адреса ссылок; жёсткая ссылка и каталоги — ноль)",
			res.TotalBytes, wantBytes)
	}
	if res.FileBytes != 1007 {
		t.Errorf("FileBytes = %d, want 1007", res.FileBytes)
	}
	if res.LinkBytes != 31 {
		t.Errorf("LinkBytes = %d, want 31", res.LinkBytes)
	}
	if res.HardlinkDupes != 1 || res.HardlinkBytes != 1000 {
		t.Errorf("hard link = %d name, %d bytes; want 1 and 1000", res.HardlinkDupes, res.HardlinkBytes)
	}
	if res.Files != 3 || res.Links != 2 || res.Dirs != 2 {
		t.Errorf("files/links/dirs = %d/%d/%d, want 3/2/2", res.Files, res.Links, res.Dirs)
	}
	if res.DirBytes == 0 {
		t.Errorf("the directories' own size must still be reported, just not counted into the total")
	}
}

// TestFoldCountsEverythingAndDecidesOnce проверяет обещание свёртки целиком, а
// не по частям: числа обхода обязаны совпасть со свёрткой и без неё до
// единицы, приговор внутри свёрнутого каталога обязан быть вынесен ОДИН раз, а
// сам каталог обязан попасть в рейтинг с размером поддерева.
//
// Отрицательный контроль здесь встроен: тот же самый счётчик решений на том же
// дереве без Fold обязан быть больше единицы, иначе тест не отличил бы свёртку
// от пустого дерева.
func TestFoldCountsEverythingAndDecidesOnce(t *testing.T) {
	root := t.TempDir()
	heavy := filepath.Join(root, "node_modules")
	if err := os.MkdirAll(filepath.Join(heavy, "pkg", "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(path string, size int) {
		if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(root, "своё.txt"), 100)
	write(filepath.Join(heavy, "a.js"), 200)
	write(filepath.Join(heavy, "pkg", "b.js"), 300)
	write(filepath.Join(heavy, "pkg", "dist", "c.js"), 400)

	counting := &countingDecider{ready: true}
	folded, err := Walk(Options{Root: root, Top: 5, Now: time.Now(), Decider: counting, Fold: FoldByName})
	if err != nil {
		t.Fatal(err)
	}
	decisionsWithFold := len(counting.seen)

	counting2 := &countingDecider{ready: true}
	plain, err := Walk(Options{Root: root, Top: 5, Now: time.Now(), Decider: counting2})
	if err != nil {
		t.Fatal(err)
	}

	if folded.Entries != plain.Entries || folded.Files != plain.Files || folded.Dirs != plain.Dirs {
		t.Errorf("свёртка изменила счёт записей: %d/%d/%d против %d/%d/%d",
			folded.Entries, folded.Files, folded.Dirs, plain.Entries, plain.Files, plain.Dirs)
	}
	if folded.TotalBytes != plain.TotalBytes {
		t.Errorf("свёртка изменила объём: %d против %d", folded.TotalBytes, plain.TotalBytes)
	}
	// Приговоров со свёрткой: корень, «своё.txt» и сам node_modules — три.
	// Без свёртки их восемь: те же три плюс пять записей внутри.
	if decisionsWithFold != 3 {
		t.Errorf("решений со свёрткой %d, ждали 3 (корень, свой файл, каталог)", decisionsWithFold)
	}
	if len(counting2.seen) <= decisionsWithFold {
		t.Fatalf("отрицательный контроль не сработал: без свёртки решений %d, со свёрткой %d",
			len(counting2.seen), decisionsWithFold)
	}
	if len(folded.Folded) != 1 || filepath.Base(folded.Folded[0].Path) != "node_modules" {
		t.Fatalf("свёрнутое не названо: %+v", folded.Folded)
	}
	// Пять: a.js, pkg, pkg/b.js, pkg/dist, pkg/dist/c.js. Сам node_modules в
	// счёт своего содержимого не входит.
	if folded.Folded[0].Entries != 5 {
		t.Errorf("внутри свёрнутого записей %d, на диске 5", folded.Folded[0].Entries)
	}
	if folded.Folded[0].Bytes != 900 {
		t.Errorf("байт в свёрнутом %d, ждали 900", folded.Folded[0].Bytes)
	}
	var found bool
	for _, e := range folded.Largest {
		if filepath.Base(e.Path) == "node_modules" {
			found = true
			if e.Size != 900 {
				t.Errorf("в рейтинге у свёрнутого размер %d, ждали 900 — размер поддерева", e.Size)
			}
		}
	}
	if !found {
		t.Error("свёрнутый каталог не попал в рейтинг крупнейших")
	}
	// Корзины обязаны сойтись с итогом, иначе свёртка «потеряла» байты.
	var bucket int64
	for _, b := range folded.ByClass {
		bucket += b.Bytes
	}
	if bucket != folded.TotalBytes {
		t.Errorf("корзины по разряду дают %d, итог %d", bucket, folded.TotalBytes)
	}
}
