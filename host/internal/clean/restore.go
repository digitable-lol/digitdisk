// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"fmt"
	"os"
	"path"
	"time"
)

// Restore puts a корзина back where it came from.
//
// It runs without a --apply key, and that asymmetry is on purpose: `clean`
// needs a key because it takes files away, `restore` does not because it gives
// them back.  Friction belongs on the destructive direction.  `restore` writes
// only to paths its own journal says it took a file from, refuses to overwrite
// anything that has appeared there since, and refuses a file in the корзина
// that is no longer the file that was put there.
//
// dryRun answers "what would come back" without moving anything, for the same
// reason `clean` has a plan: a person should be able to look first.
func Restore(j *Journal, dryRun bool, now time.Time) (*Journal, error) {
	if now.IsZero() {
		now = time.Now()
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

	touched := 0
	for i := range j.Items {
		it := &j.Items[i]
		if it.MovedAt == "" {
			it.Failed = "не переносился — возвращать нечего"
			continue
		}
		if it.RestoredAt != "" {
			it.Failed = "уже возвращён " + it.RestoredAt
			continue
		}
		if it.PurgedAt != "" {
			it.Failed = "стёрт " + it.PurgedAt + " — возврат невозможен"
			continue
		}
		it.Failed = ""
		if err := restoreOne(root, it, dryRun, now); err != nil {
			it.Failed = err.Error()
			continue
		}
		touched++
	}

	if dryRun {
		return j, nil
	}
	if touched > 0 {
		j.RestoredAt = now.UTC().Format(time.RFC3339Nano)
	}
	if err := j.write(root, path.Join(boxRel, JournalName)); err != nil {
		return j, fmt.Errorf("файлы возвращены, но журнал %s не переписан: %w", j.path, err)
	}
	return j, nil
}

func restoreOne(root *os.Root, it *Item, dryRun bool, now time.Time) error {
	info, err := root.Lstat(it.TrashRel)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("в корзине его нет: перенос не дошёл или корзину правили руками")
		}
		return fmt.Errorf("в корзине не читается: %v", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("в корзине это уже не обычный файл (%v)", info.Mode())
	}
	// The file in the корзина must still be the file that was put there.
	// A корзина is an ordinary directory and nothing stops a person from
	// editing what is in it; restoring an edited file to the original path
	// under the original name would be the tool losing track of the truth.
	want := it.Before
	if it.After != nil {
		want = *it.After
	}
	if got := identityOf(info); !want.Same(got) {
		return fmt.Errorf("в корзине его правили: %s — не возвращён", want.Differs(got))
	}

	if _, err := root.Lstat(it.Rel); err == nil {
		return fmt.Errorf("на прежнем месте уже что-то есть — перезаписывать не будем")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("прежнее место не проверяется: %v", err)
	}

	if dryRun {
		return nil
	}
	if err := root.MkdirAll(path.Dir(it.Rel), 0o700); err != nil {
		return fmt.Errorf("прежний каталог не создаётся: %v", err)
	}
	if err := root.Rename(it.TrashRel, it.Rel); err != nil {
		return fmt.Errorf("не возвращается: %v", err)
	}
	it.RestoredAt = now.UTC().Format(time.RFC3339Nano)
	return nil
}
