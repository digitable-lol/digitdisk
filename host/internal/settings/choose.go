// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package settings

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"digitdisk/internal/lang"
)

// Source says who decided the language of this run.  It is printed by
// --version and read by the tests, because "which of the five answers won" is
// the whole of this file and the only thing that can be wrong about it.
type Source string

const (
	// FromFlag: --lang on this command line.  One run, nothing stored.
	FromFlag Source = "--lang"
	// FromEnv: DIGITDISK_LANG in the environment.  One session, nothing
	// stored — this is how the two editions of a page are photographed.
	FromEnv Source = "DIGITDISK_LANG"
	// FromFile: the answer stored in settings.conf.
	FromFile Source = "settings.conf"
	// FromAsk: the person answered the question just now.
	FromAsk Source = "спрошено"
	// FromLocale: nobody chose, and the locale of the machine answered.
	FromLocale Source = "локаль"
	// FromDefault: nobody chose and the locale said nothing.
	FromDefault Source = "умолчание"
)

// Choice is the language of one run and how it came to be.
type Choice struct {
	Lang   lang.Lang
	Source Source
	// Saved is the path settings.conf was written to, empty when nothing
	// was written — which is every run but the one that asked.
	Saved string
	// Notes are the lines to put in front of a person, on stderr, in the
	// language that was chosen.  Writing in somebody's home directory and
	// failing to write in it are both news; neither is an error.
	Notes []lang.Phrase
}

// Ask is how a Choice may talk to a person.  It is nil-safe: a run with no
// terminal — a pipe, a script, a CI job, `--json` — leaves it zero, and then
// nothing is asked and nothing is written.
//
// Both files are required and both must be the terminal.  A question written
// to a terminal whose answer would come from a pipe is a question that hangs
// forever, and a tool that hangs in somebody's build is worse than a tool in
// the wrong language.
type Ask struct {
	In  io.Reader
	Out io.Writer
	// May is the caller's answer to "is this a person at a terminal".  The
	// caller asks the terminal, not this package: what a terminal is
	// belongs to internal/ui, and internal/ui is a long way above here.
	May bool
}

// Decide answers "which language is this run written in", in the order a
// person would expect: what they typed now, what they put in the environment,
// what they chose before, what they answer when asked, what the machine's
// locale says, and — when all of that is silent — Default.
//
// It writes settings.conf in exactly one case: the question was asked and
// answered.  A locale is not an answer, a flag is not an answer, and neither
// leaves a file behind in somebody's home directory.
func Decide(o Options, flag string, ask Ask) Choice {
	if flag != "" {
		if l, ok := lang.Parse(flag); ok {
			return Choice{Lang: l, Source: FromFlag}
		}
		return Choice{
			Lang:   fallback(o),
			Source: FromLocale,
			Notes:  []lang.Phrase{lang.Say("язык %q не из двух (ru, en) — вывод на языке машины", flag)},
		}
	}
	if v := strings.TrimSpace(o.env("DIGITDISK_LANG")); v != "" {
		if l, ok := lang.Parse(v); ok {
			return Choice{Lang: l, Source: FromEnv}
		}
		return Choice{
			Lang:   fallback(o),
			Source: FromLocale,
			Notes:  []lang.Phrase{lang.Say("DIGITDISK_LANG=%q — не из двух (ru, en), вывод на языке машины", v)},
		}
	}

	stored, err := Load(o)
	switch {
	case err != nil:
		// A settings file that will not parse is named and stepped over.
		// Refusing to look at a disk because a preference file has a typo
		// in it would be the tool holding the person's own machine
		// hostage over its own convenience.
		return Choice{
			Lang:   fallback(o),
			Source: FromLocale,
			Notes:  []lang.Phrase{lang.Say("настройки не прочитаны (%s) — вывод на языке машины", err)},
		}
	case stored.Lang.Valid():
		return Choice{Lang: stored.Lang, Source: FromFile}
	}

	hint, known := lang.Locale(o.Getenv)
	if !ask.May || ask.In == nil || ask.Out == nil {
		src := FromLocale
		if !known {
			src = FromDefault
		}
		return Choice{Lang: hint, Source: src}
	}

	chosen, ok := question(ask, hint)
	if !ok {
		// End of input, or an answer that is neither: the run goes on in
		// the language the machine suggested, and nothing is stored —
		// nobody chose.
		src := FromLocale
		if !known {
			src = FromDefault
		}
		return Choice{Lang: hint, Source: src}
	}

	c := Choice{Lang: chosen, Source: FromAsk}
	stored.Lang = chosen
	path, err := Save(o, stored)
	if err != nil {
		c.Notes = append(c.Notes, lang.Say("язык не сохранён (%s) — в этот раз вывод по-выбранному, а вопрос повторится", err))
		return c
	}
	c.Saved = path
	c.Notes = append(c.Notes, lang.Say("язык сохранён: %s", path))
	return c
}

// Remember stores a language chosen while the tool was already running — the
// live screen's language key.  It answers with the line to show the person,
// which is the whole point: a program that silently rewrites a file in a home
// directory is the thing this tool is for cleaning up after.
func Remember(o Options, l lang.Lang) lang.Phrase {
	stored, err := Load(o)
	if err != nil {
		stored = Settings{}
	}
	stored.Lang = l
	path, err := Save(o, stored)
	if err != nil {
		return lang.Say("язык не сохранён (%s)", err)
	}
	return lang.Say("язык сохранён: %s", path)
}

func fallback(o Options) lang.Lang {
	l, _ := lang.Locale(o.Getenv)
	return l
}

// question asks, in both languages at once, because at the moment it is asked
// there is no answer yet to the question of which one the reader has.  The
// locale picks which of the two is offered by pressing Enter, and that is all
// the locale is allowed to do: suggest.
func question(ask Ask, hint lang.Lang) (lang.Lang, bool) {
	numbered := make([]string, 0, len(lang.All))
	def := 1
	for i, l := range lang.All {
		numbered = append(numbered, l.Name())
		if l == hint {
			def = i + 1
		}
	}
	say := func(format string, a ...any) { _, _ = fmt.Fprintf(ask.Out, format, a...) }
	say("\n")
	say("digitdisk: язык вывода ещё не выбран / output language is not set yet.\n")
	for i, name := range numbered {
		mark := "  "
		if i+1 == def {
			mark = "→ "
		}
		say("  %s%d) %s\n", mark, i+1, name)
	}
	say("  выбор / choice [%d]: ", def)

	line, err := bufio.NewReader(ask.In).ReadString('\n')
	if err != nil && strings.TrimSpace(line) == "" {
		say("\n")
		return "", false
	}
	answer := strings.TrimSpace(line)
	if answer == "" {
		return lang.All[def-1], true
	}
	l, ok := lang.Parse(answer)
	if !ok {
		say("digitdisk: %q — не 1 и не 2 / neither 1 nor 2.\n", answer)
		return "", false
	}
	return l, true
}
