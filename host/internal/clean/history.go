// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry is one корзина as history sees it: when it was filled, from what root,
// what is still in it, what has already gone back or been erased, and the two
// commands that act on it.
//
// Nothing here is remembered by digitdisk between runs.  It is read, every
// time, out of the journal.json files the корзины carry — the same files
// restore and purge obey.  A separate history database would be a second
// account of the same events, and the two would disagree the first time
// somebody moved a корзина with mv.
type Entry struct {
	Box       string `json:"корзина"`
	Root      string `json:"корень"`
	Tool      string `json:"инструмент,omitempty"`
	Decider   string `json:"решающий_слой,omitempty"`
	StartedAt string `json:"начато,omitempty"`

	Planned  int `json:"запланировано"`
	Moved    int `json:"в_корзине"`
	Restored int `json:"возвращено"`
	Purged   int `json:"стёрто"`
	Failed   int `json:"не_сделано"`

	MovedBytes    int64 `json:"в_корзине_байт"`
	RestoredBytes int64 `json:"возвращено_байт"`
	PurgedBytes   int64 `json:"стёрто_байт"`

	// FreedBytes is what erasing actually gave back: the bytes of the files
	// purged.  Moving into a корзина frees nothing, so a корзина that has
	// only been filled reports zero here however large it is.
	FreedBytes int64 `json:"освобождено_байт"`

	// Problem is why this корзина could not be read.  A корзина that is
	// there and unreadable is a fact worth printing, not a row to drop.
	Problem string `json:"беда,omitempty"`
}

// Restorable reports whether anything in this корзина can still be put back.
func (e Entry) Restorable() bool { return e.Moved > 0 }

// History is every корзина found under one root, newest first.
type History struct {
	Root    string  `json:"корень"`
	Trash   string  `json:"хранилище_корзин"`
	Entries []Entry `json:"корзины"`

	Boxes         int   `json:"корзин"`
	MovedBytes    int64 `json:"в_корзинах_байт"`
	FreedBytes    int64 `json:"освобождено_байт"`
	RestoredBytes int64 `json:"возвращено_байт"`
}

// ReadHistory lists the корзины under a root.  The root may be the tree that
// was cleaned (its .digitdisk-trash is read), the хранилище itself, or one
// корзина — all three are things a person types, and refusing two of them
// would only teach them to type the third.
func ReadHistory(root string) (*History, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	abs = filepath.Clean(abs)

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s — не каталог", abs)
	}

	// One корзина named directly.
	if _, err := os.Stat(filepath.Join(abs, JournalName)); err == nil {
		h := &History{Root: abs, Trash: filepath.Dir(abs), Entries: []Entry{}}
		h.add(entryOf(abs))
		return h, nil
	}

	store := abs
	if filepath.Base(abs) != TrashName {
		store = filepath.Join(abs, TrashName)
	}
	h := &History{Root: abs, Trash: store, Entries: []Entry{}}

	names, err := os.ReadDir(store)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return h, nil // no корзина has ever been made here; that is an answer
		}
		return nil, err
	}
	for _, de := range names {
		if !de.IsDir() {
			continue
		}
		h.add(entryOf(filepath.Join(store, de.Name())))
	}
	// Newest first: the корзина a person is looking for is almost always the
	// last one they made.
	sort.Slice(h.Entries, func(i, j int) bool {
		if h.Entries[i].StartedAt != h.Entries[j].StartedAt {
			return h.Entries[i].StartedAt > h.Entries[j].StartedAt
		}
		return h.Entries[i].Box > h.Entries[j].Box
	})
	return h, nil
}

func (h *History) add(e Entry) {
	h.Entries = append(h.Entries, e)
	h.Boxes++
	h.MovedBytes += e.MovedBytes
	h.FreedBytes += e.FreedBytes
	h.RestoredBytes += e.RestoredBytes
}

func entryOf(box string) Entry {
	e := Entry{Box: box}
	j, err := ReadJournal(box)
	if err != nil {
		e.Problem = err.Error()
		return e
	}
	e.Root, e.Tool, e.Decider, e.StartedAt = j.Root, j.Tool, j.Decider, j.StartedAt
	e.Planned = len(j.Items)
	e.Moved, e.MovedBytes = j.Moved()
	e.Restored, e.RestoredBytes = j.Restored()
	e.Purged, e.PurgedBytes = j.Purged()
	e.Failed = len(j.Failed())
	e.FreedBytes = e.PurgedBytes
	return e
}

// Age is how long ago a корзина was filled, or -1 when the journal's timestamp
// cannot be read.  Printing "неизвестно" beats printing a duration measured
// from the zero time.
func (e Entry) Age(now time.Time) time.Duration {
	t, err := time.Parse(time.RFC3339Nano, e.StartedAt)
	if err != nil {
		if t, err = time.Parse(time.RFC3339, e.StartedAt); err != nil {
			return -1
		}
	}
	d := now.Sub(t)
	if d < 0 {
		return 0
	}
	return d
}
