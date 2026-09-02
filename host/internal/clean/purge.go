// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"digitdisk/internal/lang"
)

// Purge erases the contents of a корзина.  This is the one operation in
// digitdisk that cannot be undone, and it is gated three ways:
//
//   - it is its own subcommand, so it is never a flag someone else's script
//     picked up by accident;
//   - it needs --confirm N, where N is the exact number of files the journal
//     says are in the корзина.  The number is not printed by the failure, so
//     it cannot be typed without first running `purge` with no key and reading
//     the plan.  A confirmation you can satisfy without looking confirms
//     nothing;
//   - it erases only paths listed in that journal, checked one by one against
//     the отпечаток taken when they were moved.
//
// os.Remove is called on files, one at a time, and never on a directory that
// still has anything in it.  The recursive delete of the standard library
// appears nowhere in this tree and tools/check-licensing.py fails the build if
// it ever does — the identifier alone is enough to fail it, comments included,
// which is why this sentence does not name it.  A recursive delete is the call
// that turns a one-line mistake into a lost directory, and the way not to have
// that mistake is not to have the call.
func Purge(j *Journal, confirm int, now time.Time) (*Journal, error) {
	if now.IsZero() {
		now = time.Now()
	}
	inBox, bytes := j.Moved()
	if confirm < 0 {
		return nil, lang.Errorf("--confirm не может быть отрицательным")
	}
	if confirm != inBox {
		return nil, lang.Errorf("в корзине %d файлов (%d Б), а --confirm назвал %d — ничего не стёрто.\nЗапусти `digitdisk purge %s` без ключа: он напечатает, что будет стёрто, и число для --confirm",
			inBox, bytes, confirm, j.Box)
	}
	if inBox == 0 {
		return nil, lang.Errorf("корзина %s пуста — стирать нечего", j.Box)
	}

	root, err := openRoot(j.Root)
	if err != nil {
		return nil, err
	}
	defer root.Close()

	boxRel, err := j.boxRel()
	if err != nil {
		return nil, err
	}

	erased := 0
	for i := range j.Items {
		it := &j.Items[i]
		if it.MovedAt == "" || it.RestoredAt != "" || it.PurgedAt != "" {
			continue
		}
		it.Failed = lang.Phrase{}
		if err := purgeOne(root, it, now); err != nil {
			it.Failed = phraseOf(err)
			continue
		}
		erased++
	}

	// Directories the корзина made for the files are removed bottom-up with
	// os.Remove, which refuses a directory that is not empty.  Anything a
	// person put into the корзина themselves therefore survives, along with
	// the directory holding it — and the journal below it survives always.
	pruneEmptyDirs(root, j, boxRel)

	if erased > 0 {
		j.PurgedAt = now.UTC().Format(time.RFC3339Nano)
	}
	if err := j.write(root, path.Join(boxRel, JournalName)); err != nil {
		return j, lang.Errorf("стёрто %d файлов, но журнал %s не переписан: %s", erased, j.path, err)
	}
	return j, nil
}

func purgeOne(root *os.Root, it *Item, now time.Time) error {
	info, err := root.Lstat(it.TrashRel)
	if err != nil {
		if os.IsNotExist(err) {
			return lang.Errorf("в корзине его уже нет")
		}
		return lang.Errorf("в корзине не читается: %v", err)
	}
	if !info.Mode().IsRegular() {
		return lang.Errorf("в корзине это не обычный файл (%v) — стирать не будем", info.Mode())
	}
	want := it.Before
	if it.After != nil {
		want = *it.After
	}
	if got := identityOf(info); !want.Same(got) {
		return lang.Errorf("в корзине его правили: %s — не стёрт", want.Differs(got))
	}
	if err := root.Remove(it.TrashRel); err != nil {
		return lang.Errorf("не стирается: %v", err)
	}
	it.PurgedAt = now.UTC().Format(time.RFC3339Nano)
	return nil
}

// pruneEmptyDirs removes the now-empty directories the корзина created, from
// the deepest up.  Failures are silent by design: a directory that will not go
// is a directory with something in it, which is a reason to keep it.
func pruneEmptyDirs(root *os.Root, j *Journal, boxRel string) {
	seen := map[string]bool{}
	var dirs []string
	filesRoot := path.Join(boxRel, FilesDir)
	for _, it := range j.Items {
		if it.TrashRel == "" {
			continue
		}
		for d := path.Dir(it.TrashRel); strings.HasPrefix(d, filesRoot); d = path.Dir(d) {
			if !seen[d] {
				seen[d] = true
				dirs = append(dirs, d)
			}
			if d == filesRoot {
				break
			}
		}
	}
	sort.Slice(dirs, func(a, b int) bool { return len(dirs[a]) > len(dirs[b]) })
	for _, d := range dirs {
		_ = root.Remove(d)
	}
}
