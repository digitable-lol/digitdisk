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

// Erase removes the files a plan lists.  For good: there is no корзина behind
// this, and `restore` has nothing to work with afterwards.
//
// IT IS THE SAME PLAN.  Erase takes what Make produced and nothing else — the
// приговор «МожноУбрать» of the decision layer, narrowed by the отметки, minus
// what the защитный список kept and minus what guard() refused.  It adds no
// rule, widens no ground and knows no path the plan does not name.  A caller
// who wants to erase something the layer did not mark has no way to say so
// here, and that is the whole point of this file being three dozen lines long.
//
// What differs from Apply, and only this:
//
//   - the step on each file is root.Remove instead of root.Rename, so the file
//     is gone rather than renamed;
//   - the корзина that is made holds the JOURNAL ALONE.  Nothing is put in it,
//     because there is nothing left to put;
//   - the journal is stamped «способ: стирание», so it can never be read back
//     as a корзина somebody could empty.
//
// Everything else is Apply's order of operations, kept because the reasons for
// it are the same:
//
//  1. the journal goes down in full BEFORE the first file is touched, so a
//     crash in the middle leaves a record of what was about to go;
//  2. each file is lstat-ed again through the корень's os.Root and compared
//     with the отпечаток taken during the walk — a file that changed since it
//     was judged is refused, not erased;
//  3. only a regular file is removed, one at a time, by a call that cannot
//     touch a directory.
//
// Why a journal for something that cannot be undone.  Because the person whose
// files these were has one question left afterwards — WHAT WENT — and the only
// honest answer is a list written down at the time.  A tool that erases and
// keeps no record leaves them with a smaller disk and no way to find out what
// used to be on it.
func Erase(p Plan, opt Options) (*Journal, error) {
	if !p.DeciderReady {
		return nil, lang.Errorf(`решающий слой — %s.
Он не выносит приговора «%s» ни одной записи, поэтому убирать нечего и не по чему.
Собери хозяина с признаком flangcore: go build -tags flangcore -o digitdisk ./host`,
			p.Decider, "МожноУбрать")
	}
	if len(p.Items) == 0 {
		return nil, lang.Errorf("стирать нечего: ядро не пометило «%s» ни одного файла под %s", "МожноУбрать", p.Root)
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
		Way:             WayErase,
		Tool:            "digitdisk " + firstNonEmpty(opt.Version, "dev"),
		ContractVersion: p.ContractVersion,
		Decider:         p.Decider,
		Root:            p.Root,
		Box:             boxAbs,
		StartedAt:       now.UTC().Format(time.RFC3339Nano),
		Items:           append([]Item(nil), p.Items...),
		Refused:         p.Refused,
	}

	journalRel := path.Join(boxRelPath, JournalName)
	if err := root.MkdirAll(boxRelPath, 0o700); err != nil {
		return nil, lang.Errorf("корзина %s не создаётся: %s", boxAbs, err)
	}
	// Written before the first file goes, and for a стирание that order
	// matters more than it does for a перенос: after a crash the корзина of a
	// перенос still holds the files, and here the list is all there is.
	if err := j.write(root, journalRel); err != nil {
		return nil, lang.Errorf("журнал %s не записывается: %s", filepath.Join(boxAbs, JournalName), err)
	}
	j.path = filepath.Join(boxAbs, JournalName)

	erased := 0
	for i := range j.Items {
		it := &j.Items[i]
		if err := eraseOne(root, it, now); err != nil {
			it.Failed = phraseOf(err)
			continue
		}
		erased++
	}

	j.FinishedAt = time.Now().UTC().Format(time.RFC3339Nano)
	if erased > 0 {
		j.PurgedAt = now.UTC().Format(time.RFC3339Nano)
	}
	if err := j.write(root, journalRel); err != nil {
		return j, lang.Errorf(`файлы стёрты, но журнал %s не переписан: %s.
Первая запись журнала на месте, и она называет всё, что стиралось`, j.path, err)
	}
	return j, nil
}

// eraseOne removes one file, refusing at the first sign that it is not the file
// the plan judged.
func eraseOne(root *os.Root, it *Item, now time.Time) error {
	info, err := root.Lstat(it.Rel)
	if err != nil {
		if os.IsNotExist(err) {
			return lang.Errorf("исчез между обходом и стиранием")
		}
		return lang.Errorf("не читается: %v", err)
	}
	if !info.Mode().IsRegular() {
		return lang.Errorf("перестал быть обычным файлом (стал %v)", info.Mode())
	}
	if current := identityOf(info); !it.Before.Same(current) {
		return lang.Errorf("%s — не стёрт", it.Before.Differs(current))
	}
	if err := root.Remove(it.Rel); err != nil {
		return lang.Errorf("не стирается: %v", err)
	}
	it.PurgedAt = now.UTC().Format(time.RFC3339Nano)
	return nil
}
