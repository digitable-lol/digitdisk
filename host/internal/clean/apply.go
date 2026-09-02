// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"digitdisk/internal/lang"
)

// Apply moves every item of a plan into a fresh корзина and writes the journal.
//
// It erases nothing.  Every file it touches ends up under
// <корень>/.digitdisk-trash/<метка>/files/ with its path below the корень
// preserved, and `restore` puts it back.  Space is not freed here and the
// report says so: the bytes are on the same filesystem under a different name,
// which is exactly why the move is instant and reversible.
//
// Order of operations per file, and none of it is optional:
//
//  1. lstat the file again through the корень's os.Root;
//  2. compare it with the отпечаток taken during the walk — any difference and
//     the file is refused, not moved;
//  3. make the destination directory, refusing if something is already there;
//  4. rename(2) — atomic, so the file is never in two places or in none;
//  5. lstat the destination and check the inode is the one that was moved.
func Apply(p Plan, opt Options) (*Journal, error) {
	if !p.DeciderReady {
		return nil, lang.Errorf(`решающий слой — %s.
Он не выносит приговора «%s» ни одной записи, поэтому убирать нечего и не по чему.
Собери хозяина с признаком flangcore: go build -tags flangcore -o digitdisk ./host`,
			p.Decider, "МожноУбрать")
	}
	if len(p.Items) == 0 {
		return nil, lang.Errorf("убирать нечего: ядро не пометило «%s» ни одного файла под %s", "МожноУбрать", p.Root)
	}

	now := opt.Now
	if now.IsZero() {
		now = time.Now()
	}
	boxAbs := filepath.Join(p.Trash, stamp(now))
	boxRelPath, err := relUnder(p.Root, boxAbs)
	if err != nil {
		return nil, err
	}

	root, err := openRoot(p.Root)
	if err != nil {
		return nil, err
	}
	defer root.Close()

	j := &Journal{
		Version:         JournalVersion,
		Tool:            "digitdisk " + firstNonEmpty(opt.Version, "dev"),
		ContractVersion: p.ContractVersion,
		Decider:         p.Decider,
		Root:            p.Root,
		Box:             boxAbs,
		StartedAt:       now.UTC().Format(time.RFC3339Nano),
		Items:           append([]Item(nil), p.Items...),
		Refused:         p.Refused,
	}

	// Where each item is going, decided and recorded before anything moves.
	for i := range j.Items {
		rel, err := trashRelFor(p.Root, boxAbs, j.Items[i].Rel)
		if err != nil {
			return nil, err
		}
		j.Items[i].TrashRel = rel
	}

	journalRel := path.Join(boxRelPath, JournalName)
	if err := root.MkdirAll(boxRelPath, 0o700); err != nil {
		return nil, lang.Errorf("корзина %s не создаётся: %s", boxAbs, err)
	}
	// The journal goes down first, listing every intended move with no
	// outcome.  A crash from here on leaves a корзина `restore` can empty.
	if err := j.write(root, journalRel); err != nil {
		return nil, lang.Errorf("журнал %s не записывается: %s", filepath.Join(boxAbs, JournalName), err)
	}
	j.path = filepath.Join(boxAbs, JournalName)

	for i := range j.Items {
		it := &j.Items[i]
		if err := move(root, it, now); err != nil {
			it.Failed = phraseOf(err)
			continue
		}
	}

	j.FinishedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if err := j.write(root, journalRel); err != nil {
		return j, lang.Errorf(`файлы перенесены, но журнал %s не переписан: %s.
Первая запись журнала на месте и возврат по ней работает`, j.path, err)
	}
	return j, nil
}

// move carries out one file's move, refusing at the first sign that the file
// is not the one the plan judged.
func move(root *os.Root, it *Item, now time.Time) error {
	info, err := root.Lstat(it.Rel)
	if err != nil {
		if os.IsNotExist(err) {
			return lang.Errorf("исчез между обходом и переносом")
		}
		return lang.Errorf("не читается: %v", err)
	}
	if !info.Mode().IsRegular() {
		return lang.Errorf("перестал быть обычным файлом (стал %v)", info.Mode())
	}
	current := identityOf(info)
	if !it.Before.Same(current) {
		return lang.Errorf("%s — не убран", it.Before.Differs(current))
	}

	if err := root.MkdirAll(path.Dir(it.TrashRel), 0o700); err != nil {
		return lang.Errorf("каталог в корзине не создаётся: %v", err)
	}
	if _, err := root.Lstat(it.TrashRel); err == nil {
		return lang.Errorf("в корзине уже есть %s — перезаписывать не будем", it.TrashRel)
	} else if !os.IsNotExist(err) {
		return lang.Errorf("корзина не проверяется: %v", err)
	}

	if err := root.Rename(it.Rel, it.TrashRel); err != nil {
		return lang.Errorf("не переносится: %v", err)
	}

	after, err := root.Lstat(it.TrashRel)
	if err != nil {
		return lang.Errorf("перенесён, но в корзине не находится: %v", err)
	}
	id := identityOf(after)
	if id.Ino != it.Before.Ino || id.Dev != it.Before.Dev {
		return lang.Errorf("в корзине оказался не тот файл (узел %d:%d вместо %d:%d)",
			id.Dev, id.Ino, it.Before.Dev, it.Before.Ino)
	}
	it.After = &id
	it.MovedAt = now.UTC().Format(time.RFC3339Nano)
	return nil
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
