// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

import (
	"fmt"
	"sort"
	"sync"
)

// dict is the whole dictionary: Russian wording → English wording.  Russian
// needs no map because the key is the Russian, and an identity map is a second
// copy of the same text waiting to disagree with the first.
var dict = map[string]string{}

// add registers one area's worth of wordings.  Every dict_*.go file calls it
// from its init, so a new area is a new file and touches nothing shared.
//
// A wording registered twice with two different English halves is a defect
// caught here and not in the output: the two callers would then get whichever
// init ran last, and which one that is depends on file names.
func add(pairs map[string]string) {
	for ru, en := range pairs {
		if was, dup := dict[ru]; dup && was != en {
			panic(fmt.Sprintf("lang: «%s» переведено дважды и по-разному: %q и %q", ru, was, en))
		}
		dict[ru] = en
	}
}

// misses records wordings asked for and not found.  A miss prints the Russian
// — half a report is better than a crash in a tool a person runs to look at
// their disk — but it is a defect, and complete_test.go is what makes sure the
// list stays empty.
var (
	missMu sync.Mutex
	misses = map[string]bool{}
)

// T translates one wording.  In Russian it is the wording itself.
func (l Lang) T(ru string) string {
	if l != EN {
		return ru
	}
	if en, ok := dict[ru]; ok {
		return en
	}
	missMu.Lock()
	misses[ru] = true
	missMu.Unlock()
	return ru
}

// F translates a wording and fills it in.  The verbs of the English half are
// the verbs of the Russian half, in the same order — complete_test.go checks
// exactly that, because a swapped %s and %d is a crash in the other language
// and in no test that only ever ran in one.
func (l Lang) F(ru string, a ...any) string { return fmt.Sprintf(l.T(ru), a...) }

// Missing lists the wordings asked for that the dictionary does not have.  It
// is what a test asserts is empty after exercising the reports.
func Missing() []string {
	missMu.Lock()
	defer missMu.Unlock()
	out := make([]string, 0, len(misses))
	for ru := range misses {
		out = append(out, ru)
	}
	sort.Strings(out)
	return out
}

// Known reports whether the dictionary has an English half for this wording.
func Known(ru string) bool { _, ok := dict[ru]; return ok }

// Size is how many wordings have a pair.  It is printed by the check so the
// number is a number somebody can watch move.
func Size() int { return len(dict) }

// Wordings lists every Russian wording that has a pair, sorted.
func Wordings() []string {
	out := make([]string, 0, len(dict))
	for ru := range dict {
		out = append(out, ru)
	}
	sort.Strings(out)
	return out
}

// vocab is the part of the dictionary whose Russian half is never written in
// the source as a literal: it is a VALUE that arrives at run time — разряд
// «Кэш» from the решающий слой, вид «Каталог» from the обход, the name of a
// fact nobody measured.  Those words are identifiers on the way in and words
// on the way out, and only the way out is translated.
//
// They are registered apart from the wordings so that complete_test.go can
// tell the two cases apart: a wording must be found in a T call somewhere, a
// vocabulary word must not be, and neither list may hold what belongs to the
// other.
var vocab = map[string]bool{}

// addVocab registers value words.  They go into the same dict — one lookup,
// one place to be missing from — and into vocab as well, which is only a note
// about where they come from.
func addVocab(pairs map[string]string) {
	add(pairs)
	for ru := range pairs {
		vocab[ru] = true
	}
}

// IsVocab reports whether this wording is a value word rather than a literal
// of the source.
func IsVocab(ru string) bool { return vocab[ru] }

// Word translates a value that arrived from the решающий слой, from the
// справочник or from a collector — разряд, приговор, вид, якорь, имя
// неизмеренного.  The value itself is not touched: this is the word the screen
// shows for it, and the identifier goes on travelling in the JSON unchanged.
//
// A value with no entry is shown as it came.  That is not a hole the check
// tolerates — every разряд and every приговор of the договор has an entry, and
// core_test.go in this package walks the договор to prove it — but a справочник
// a person wrote themselves may hold a word this dictionary never saw, and a
// tool that refused to print it would be refusing to show a person their own
// file.
func (l Lang) Word(ru string) string {
	if l != EN {
		return ru
	}
	if en, ok := dict[ru]; ok {
		return en
	}
	return ru
}
