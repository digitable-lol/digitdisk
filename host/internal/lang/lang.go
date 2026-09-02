// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package lang is the one place digitdisk keeps the language of its output.
//
// # Two languages, one set of facts
//
// Everything a person reads — the sections of the report, the labels, the
// units, разряды and приговоры on the screen, the справка, the refusals,
// --why, the list of commands — exists in Russian and in English.  Everything
// a program reads — the JSON of every подкоманда, the ключи of the журнал, the
// имена of разряды and приговоров inside the договор — does not: a script that
// parses `digitdisk clean --json` must not care what language the person who
// ran it speaks.  The line between the two runs exactly here.  A value that
// travels in JSON is written by Phrase, which renders Russian for the machine
// and the chosen language for the screen, so translating a refusal never moves
// a byte of the JSON.
//
// # The key is the Russian text
//
// A message is looked up by its Russian wording, not by an invented
// identifier: `l.T("СИСТЕМА")`.  Russian is then the identity translation and
// cannot rot, English is one map entry away, and the check that every line has
// a pair is a check over the source itself — see complete_test.go, which reads
// every .go file of the host, finds every T/F call, and fails on a wording
// with no English beside it and on a Cyrillic literal printed past the
// dictionary.
//
// # Nothing here decides anything
//
// Разряд «Кэш» and приговор «МожноУбрать» are identifiers of the layer in
// core/, proved there and named there.  This package never translates the
// values it receives from the core; it translates the WORD THE SCREEN SHOWS
// for them (Class, Removable) and leaves the identifier alone.  See Class and
// Verdict below: the mapping is one-way, host-side, and display-only.
package lang

import (
	"os"
	"strings"
)

// Lang is one of the two languages digitdisk speaks.
type Lang string

const (
	// RU is Russian: the language the tool was written in, and therefore
	// the language its dictionary is keyed by.
	RU Lang = "ru"
	// EN is English.
	EN Lang = "en"
)

// Default is the language of a run nobody chose one for and no locale
// answered for.
//
// It is English, and the reason is POSIX rather than taste: an unset locale,
// "C" and "POSIX" all name the portable locale, whose messages are English by
// definition.  A machine that says nothing about its language is not saying
// "Russian" — it is saying "the portable one" — and answering it in Russian
// would be a guess about the reader.  A person who wants Russian either has a
// ru locale, or is asked once and says so (see Settings).
const Default = EN

// All lists both languages in the order the first-run question offers them.
var All = []Lang{RU, EN}

// Valid reports whether l is one of the two.
func (l Lang) Valid() bool { return l == RU || l == EN }

// Other is the language l is not.  The live screen switches with it.
func (l Lang) Other() Lang {
	if l == RU {
		return EN
	}
	return RU
}

// Name is what the language calls itself, for the question and the screen.
func (l Lang) Name() string {
	if l == RU {
		return "русский"
	}
	return "English"
}

// Parse reads a language out of what a person or a file wrote.  It takes the
// spellings a person would try — the code, the language's own name, the digit
// the first-run question offers — and refuses everything else rather than
// guessing.
func Parse(s string) (Lang, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "ru", "rus", "russian", "русский", "рус", "ру", "1":
		return RU, true
	case "en", "eng", "english", "английский", "англ", "2":
		return EN, true
	}
	return "", false
}

// Locale is what the environment says about the reader's language, and the
// order is the one POSIX gives: LC_ALL overrides everything, LC_MESSAGES is
// the category that governs messages, LANG is the fallback for both.
//
// ok is false when nothing was set or when what was set names no language this
// tool speaks — «C», «POSIX», «de_DE.UTF-8».  The caller then uses Default,
// and, on a terminal, asks.
func Locale(env func(string) string) (Lang, bool) {
	if env == nil {
		env = os.Getenv
	}
	for _, name := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		v := strings.TrimSpace(env(name))
		if v == "" {
			continue
		}
		// ru_RU.UTF-8 → ru; en_GB → en; C.UTF-8 → C.
		code := v
		if i := strings.IndexAny(code, "_.@"); i >= 0 {
			code = code[:i]
		}
		switch strings.ToLower(code) {
		case "ru":
			return RU, true
		case "en":
			return EN, true
		case "c", "posix":
			// The portable locale is an answer, and the answer is
			// Default — but it is not a language the reader chose,
			// so it is reported as "nothing was said".
			return Default, false
		}
		// A locale naming some third language is not this tool's, and
		// the next variable is not consulted: LC_ALL=de_DE means the
		// reader set German on purpose, and LANG behind it is stale.
		return Default, false
	}
	return Default, false
}
