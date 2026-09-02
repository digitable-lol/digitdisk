// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"digitdisk/internal/lang"
)

// Journal is the record of one корзина: what was moved into it, from where,
// how big it was, when, and what happened to it afterwards.  It is the file
// `restore` and `purge` read; the text printed on the screen is for a person
// and is not the record.
//
// It is written twice.  Once before anything moves, with every intended
// путь → в_корзине pair and no outcomes — so that a crash leaves something
// restorable.  Once after, with the outcomes filled in.
//
// A journal records ONE OF TWO THINGS and says which in «способ».  A journal
// of a перенос has a корзина with the files in it and every entry restorable;
// a journal of a стирание has a корзина holding nothing but the journal
// itself, because the files it names are gone.  The two are never told apart
// by guessing from the counts: `restore` and `purge` read «способ» and refuse
// a journal of the other kind out loud.
type Journal struct {
	Version         int       `json:"версия_журнала"`
	Way             Way       `json:"способ"`
	Tool            string    `json:"инструмент"`
	ContractVersion int       `json:"версия_договора"`
	Decider         string    `json:"решающий_слой"`
	Root            string    `json:"корень"`
	Box             string    `json:"корзина"`
	StartedAt       string    `json:"начато"`
	FinishedAt      string    `json:"закончено,omitempty"`
	RestoredAt      string    `json:"возврат,omitempty"`
	PurgedAt        string    `json:"стирание,omitempty"`
	Items           []Item    `json:"записи"`
	Refused         []Refusal `json:"отказано,omitempty"`

	// path is where this journal was read from or written to.  It is not
	// serialised: a file that carries its own location is wrong the moment
	// the directory is renamed.
	path string
}

// Way is «способ» — what the journal's корзина did with the files it names.
// It travels in JSON and in the journal file as the Russian word, for the
// reason the разряд and the приговор do: it is a name of a record, not a line
// of output.  The screen translates the WORD it shows for it and leaves this
// alone.
type Way string

const (
	// WayMove — перенос: the files are in the корзина and restore puts
	// them back.  An empty «способ» means this too, so a journal written
	// by an older digitdisk keeps its meaning.
	WayMove Way = "перенос"
	// WayErase — стирание: the files were removed and there is nothing to
	// put back.  The корзина of such a journal holds the journal alone.
	WayErase Way = "стирание"
)

// Erasure reports whether this journal records a стирание.
func (j *Journal) Erasure() bool { return j.Way == WayErase }

// Path is the file this journal lives in.
func (j *Journal) Path() string { return j.path }

// Moved counts the entries that reached the корзина and are still there.
func (j *Journal) Moved() (n int, bytes int64) {
	for _, it := range j.Items {
		if it.MovedAt != "" && it.RestoredAt == "" && it.PurgedAt == "" {
			n++
			bytes += it.Size
		}
	}
	return n, bytes
}

// Failed counts the entries that were planned but refused at the moment of
// moving — the file had changed, vanished, or was not what it had been.
func (j *Journal) Failed() []Item {
	var out []Item
	for _, it := range j.Items {
		if it.MovedAt == "" && !it.Failed.Empty() {
			out = append(out, it)
		}
	}
	return out
}

// Restored counts the entries put back.
func (j *Journal) Restored() (n int, bytes int64) {
	for _, it := range j.Items {
		if it.RestoredAt != "" {
			n++
			bytes += it.Size
		}
	}
	return n, bytes
}

// Purged counts the entries erased.
func (j *Journal) Purged() (n int, bytes int64) {
	for _, it := range j.Items {
		if it.PurgedAt != "" {
			n++
			bytes += it.Size
		}
	}
	return n, bytes
}

// write serialises the journal into the корзина through the корень's os.Root,
// then fsyncs it.  Without the fsync the record could be the thing lost in the
// crash it exists to survive.
func (j *Journal) write(root *os.Root, journalRel string) error {
	body, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')

	if err := root.MkdirAll(path.Dir(journalRel), 0o700); err != nil {
		return err
	}
	f, err := root.Create(journalRel)
	if err != nil {
		return err
	}
	if _, err := f.Write(body); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// ReadJournal loads the journal of one корзина and checks that the корзина is
// where the journal says it is.  A корзина that has been moved is refused
// rather than guessed at: its entries name absolute paths to restore to, and
// a копия of a корзина somewhere else would restore into the original tree —
// which is very probably not what the person holding the copy meant.
func ReadJournal(boxPath string) (*Journal, error) {
	abs, err := filepath.Abs(boxPath)
	if err != nil {
		return nil, err
	}
	abs = filepath.Clean(abs)

	file := filepath.Join(abs, JournalName)
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, lang.Errorf("в %s нет %s — это не корзина digitdisk", abs, JournalName)
		}
		return nil, err
	}
	defer f.Close()

	body, err := io.ReadAll(io.LimitReader(f, 512<<20))
	if err != nil {
		return nil, err
	}
	var j Journal
	if err := json.Unmarshal(body, &j); err != nil {
		return nil, lang.Errorf("%s не читается как журнал: %s", file, err)
	}
	if j.Version != JournalVersion {
		return nil, lang.Errorf("%s: версия журнала %d, а этот digitdisk понимает %d — работать с непонятым журналом опаснее, чем отказаться",
			file, j.Version, JournalVersion)
	}
	if j.Box != abs {
		return nil, lang.Errorf(`корзина лежит в %s, а журнал записан для %s.
Возврат кладёт файлы по абсолютным путям, записанным при переносе;
из перемещённой копии он писал бы в исходное дерево, а не туда, где копия`,
			abs, j.Box)
	}
	if j.Root == "" {
		return nil, lang.Errorf("%s: в журнале нет корня", file)
	}
	j.path = file
	return &j, nil
}

// boxRel is the корзина's path relative to the корень.
func (j *Journal) boxRel() (string, error) {
	rel, err := filepath.Rel(j.Root, j.Box)
	if err != nil {
		return "", err
	}
	if rel == "." || strings.HasPrefix(rel, "..") {
		return "", lang.Errorf("корзина %s вне корня %s", j.Box, j.Root)
	}
	return filepath.ToSlash(rel), nil
}
