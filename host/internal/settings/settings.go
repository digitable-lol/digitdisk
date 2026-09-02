// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package settings is digitdisk's own home in the reader's home directory:
// ~/.digitable/digitdisk/.
//
// # One home, not two
//
// Before this file the tool had one home already — ~/.config/digitdisk/, where
// places.conf and protect.conf live — and the language would have made a
// second.  Two homes are two places to look for one answer, and the tool would
// then have to explain in every document which of them holds what.  So there
// is one, it is ~/.digitable/, and it is the family's: digitdisk keeps its
// settings in ~/.digitable/digitdisk/, and the tools beside it will keep
// theirs beside it.
//
// The old home is not broken and not deleted.  Three files are affected —
// settings.conf (new), places.conf and protect.conf (moved) — and the two that
// moved are still READ where they were: a person whose ~/.config/digitdisk
// works today goes on working, is told once where the tool now looks, and
// moves the file when it suits them.  Nothing is copied on their behalf: a
// tool that writes into a person's home unasked is the thing this tool exists
// to clean up after.
//
// # Writing here is an action, and it is announced
//
// The only thing digitdisk writes here is the language, and only after a
// person answered the question with their own hands.  It then says what it
// wrote and where, once, in one line.  If the home cannot be written — a
// read-only mount, a directory owned by somebody else — the run goes on in the
// language that was chosen and says plainly that it was not saved.  Refusing
// to work because a preference could not be stored would be answering a small
// problem with a big one.
package settings

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"digitdisk/internal/lang"
)

// FamilyDir is the family's home under the reader's home directory.
const FamilyDir = ".digitable"

// ToolDir is digitdisk's own room in it.
const ToolDir = "digitdisk"

// FileName is the settings file itself.
const FileName = "settings.conf"

// LegacyDir is where the tool used to look, and still does when the new home
// holds nothing.
const LegacyDir = "digitdisk"

// Options is everything this package needs from the outside.  Nothing is read
// from the process directly, so a test names its own home and its own
// environment and never touches the machine it runs on.
type Options struct {
	Home   string              // defaults to os.UserHomeDir
	Getenv func(string) string // defaults to os.Getenv
}

func (o Options) env(name string) string {
	if o.Getenv == nil {
		return os.Getenv(name)
	}
	return o.Getenv(name)
}

func (o Options) home() (string, error) {
	if o.Home != "" {
		return o.Home, nil
	}
	return os.UserHomeDir()
}

// Dir is ~/.digitable/digitdisk.
func Dir(o Options) (string, error) {
	home, err := o.home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, FamilyDir, ToolDir), nil
}

// LegacyPath is where a file used to be looked for: $XDG_CONFIG_HOME/digitdisk
// or ~/.config/digitdisk.
func LegacyPath(o Options, name string) (string, error) {
	dir := o.env("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := o.home()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, LegacyDir, name), nil
}

// Find looks for one of the tool's files by name — "places.conf",
// "protect.conf" — in the new home first and in the old one after.  It returns
// the path that exists, whether that path is the old home, and false when
// neither holds it.
//
// A file present in both is read from the new home: the person moved it and
// left a copy behind, and the copy they moved is the one they meant.
func Find(o Options, name string) (path string, legacy bool, ok bool) {
	if dir, err := Dir(o); err == nil {
		p := filepath.Join(dir, name)
		if exists(p) {
			return p, false, true
		}
	}
	if p, err := LegacyPath(o, name); err == nil {
		if exists(p) {
			return p, true, true
		}
	}
	return "", false, false
}

func exists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

// MovedNote is the one line that tells a person their file is being read from
// the old home.  It names both paths, because "the settings moved" without the
// two paths is a sentence nobody can act on.
func MovedNote(old, now string) lang.Phrase {
	return lang.Say("настройки переехали в %s — %s ещё читается, перенесите его, когда удобно", now, old)
}

// Settings is the content of settings.conf: what a person chose and what they
// have already been told.
type Settings struct {
	// Lang is the chosen language, empty when nobody has chosen.
	Lang lang.Lang
	// MoveAnnounced records that the reader has been told the home moved,
	// so that the line is said once and not at every run.
	MoveAnnounced bool
	// Path is where this came from, empty when it came from nowhere.
	Path string
	// Unknown keeps the lines this version does not understand, so that a
	// newer digitdisk's setting is not erased by an older one rewriting
	// the file.
	Unknown []string
}

// keys are spelled in both languages on the way in, because the file is read
// by a person as well as by this code, and a Russian key in an English
// person's settings is the same discourtesy as a Russian report.
func settingKey(name string) string {
	switch strings.ToLower(name) {
	case "язык", "lang", "language":
		return "lang"
	case "переезд_объявлен", "move_announced":
		return "move_announced"
	}
	return ""
}

// Load reads settings.conf.  A missing file is not an error: it is the state
// of a machine where nobody has chosen yet, and that is the state the question
// exists for.
func Load(o Options) (Settings, error) {
	var s Settings
	dir, err := Dir(o)
	if err != nil {
		return s, err
	}
	path := filepath.Join(dir, FileName)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return s, err
	}
	defer f.Close()
	s.Path = path

	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		name, value, found := strings.Cut(text, "=")
		if !found {
			return s, lang.Errorf("%s, строка %d: ждалось «ключ=значение», а написано %q", path, line, text)
		}
		name, value = strings.TrimSpace(name), strings.TrimSpace(value)
		switch settingKey(name) {
		case "lang":
			l, ok := lang.Parse(value)
			if !ok {
				return s, lang.Errorf("%s, строка %d: язык %q не из двух (ru, en)", path, line, value)
			}
			s.Lang = l
		case "move_announced":
			s.MoveAnnounced = value == "да" || value == "yes" || value == "true"
		default:
			s.Unknown = append(s.Unknown, text)
		}
	}
	return s, sc.Err()
}

// Save writes settings.conf, making the home if it is not there.  It returns
// the path it wrote, and an error the caller is expected to SAY rather than
// die of.
func Save(o Options, s Settings) (string, error) {
	dir, err := Dir(o)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, FileName)

	l := s.Lang
	if !l.Valid() {
		l = lang.Default
	}
	var b strings.Builder
	b.WriteString(l.T("# Настройки digitdisk. Правится руками: ключ=значение, одна строка — один ключ.") + "\n")
	b.WriteString(l.T("# Этот файл завёл сам digitdisk, когда спросил про язык. Убрать его можно.") + "\n")
	b.WriteString("\n")
	if l == lang.RU {
		b.WriteString("язык=" + string(l) + "\n")
		if s.MoveAnnounced {
			b.WriteString("переезд_объявлен=да\n")
		}
	} else {
		b.WriteString("lang=" + string(l) + "\n")
		if s.MoveAnnounced {
			b.WriteString("move_announced=yes\n")
		}
	}
	for _, line := range s.Unknown {
		b.WriteString(line + "\n")
	}

	// Written in one call and in one piece.  The tidier spelling — write
	// beside it, rename over it — is not available here and that is on
	// purpose: Rename and Remove live in internal/clean and nowhere else,
	// where the приговор of the ядро, the отпечаток and the журнал stand
	// around them, and tools/licensing.flang fails the run over a call to
	// either anywhere else in the tree.  A settings file is a few dozen
	// bytes written by one call; the risk that buys the exception is not
	// worth the hole it would open.
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
