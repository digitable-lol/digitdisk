// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package protect holds the защитный список: the operator's standing "do not
// touch this".
//
// # Why this is the host's and not the core's
//
// The decision layer answers one question — what IS this path — and every
// answer it gives is proved.  "Do not touch my ~/projects" is not an answer to
// that question: the path may well be a кэш, and saying otherwise in the
// справочник would be writing a falsehood into the layer to get an effect.
// The защитный список is an instruction from the person holding the machine,
// and instructions belong where the host already keeps its veto — next to the
// checks in internal/clean that refuse an item the layer approved.
//
// That placement buys three things.  It cannot weaken anything: the list only
// ever moves an item from "will be moved" to "refused", never the other way,
// so no постусловие of the core is touched.  It costs nothing: the list is
// consulted only for items the layer already marked «МожноУбрать», a few
// hundred paths, not the millions the walk visits.  And it is honest in the
// report: a protected path is listed under ЗАЩИЩЕНО with the rule and the file
// line that protected it, not silently missing from the plan.
//
// # What a rule may say
//
//	путь|<путь>|<почему>      — this path and everything under it
//	разряд|<разряд>|<почему>  — every find of that разряд
//
// A path rule is matched the way the справочник is matched — by whole
// components, with the slashes on both ends — so `~/projects` protects
// `~/projects/x` and does not protect `~/projects-old`.  A path written
// without a leading `/` or `~/` protects that chain of components at any
// depth: `node_modules` protects every node_modules on the disk.
package protect

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"digitdisk/internal/core"
	"digitdisk/internal/lang"
	"digitdisk/internal/settings"
)

// FileName is the name of the защитный список.  Where it is looked for is
// settings.Find's business: ~/.digitable/digitdisk/ first, the old
// $XDG_CONFIG_HOME/digitdisk (else ~/.config/digitdisk) after.
const FileName = "protect.conf"

// Kind is what a rule matches on.
type Kind string

const (
	KindPath  Kind = "путь"
	KindClass Kind = "разряд"
)

// Rule is one line of the защитный список.
type Rule struct {
	Kind   Kind   `json:"вид"`
	Value  string `json:"значение"`
	Why    string `json:"почему,omitempty"`
	Origin string `json:"откуда"`
	Line   int    `json:"строка,omitempty"`

	// chain is the resolved, slash-bounded path for a KindPath rule.
	chain string
	// anywhere says the chain may sit at any depth.
	anywhere bool
	// class is the parsed разряд for a KindClass rule.
	class core.Class
}

// Where names the rule for a person: the file and line it came from.
func (r Rule) Where() string {
	if r.Line > 0 {
		return fmt.Sprintf("%s:%d", r.Origin, r.Line)
	}
	return r.Origin
}

// String is how a refusal names the rule that caused it.
func (r Rule) String() string {
	s := fmt.Sprintf("%s %s (%s)", r.Kind, r.Value, r.Where())
	if r.Why != "" {
		s += ": " + r.Why
	}
	return s
}

// In is String for a reader: the same rule, with the вид named in the reader's
// language.
//
// String stays as it is and is not routed through here: a rule printed into a
// журнал or compared in a test must read the same on every machine, and «путь
// ~/projects (protect.conf:3)» is what it has always read.  Only what a person
// sees on the screen changes — and only the вид changes in it, because the
// value, the file and the line are what the person wrote themselves.
func (r Rule) In(l lang.Lang) string {
	s := fmt.Sprintf("%s %s (%s)", l.Word(string(r.Kind)), r.Value, r.Where())
	if r.Why != "" {
		s += ": " + r.Why
	}
	return s
}

// List is the whole защитный список.
type List struct {
	Origins []string `json:"откуда"`
	Rules   []Rule   `json:"правила"`

	// Moved is the one line to say when this список was read from the home
	// it used to live in.  Out of the JSON on purpose: it is news for a
	// person, not a field of the plan.
	Moved lang.Phrase `json:"-"`
}

// Empty reports whether nothing is protected.
func (l *List) Empty() bool { return l == nil || len(l.Rules) == 0 }

// Covers reports the first rule that protects this path or разряд.  Path rules
// are tried before class rules: a person reading "почему не тронули" is better
// served by the exact path they wrote than by a разряд that also matched.
func (l *List) Covers(path string, class core.Class) (Rule, bool) {
	if l == nil {
		return Rule{}, false
	}
	trail := filepath.ToSlash(path) + "/"
	for _, r := range l.Rules {
		if r.Kind != KindPath {
			continue
		}
		if r.anywhere {
			if strings.Contains(trail, r.chain) {
				return r, true
			}
			continue
		}
		if strings.HasPrefix(trail, r.chain) {
			return r, true
		}
	}
	for _, r := range l.Rules {
		if r.Kind == KindClass && r.class == class {
			return r, true
		}
	}
	return Rule{}, false
}

// Options configures loading.
type Options struct {
	// File, when set, is read instead of the user's file.
	File string
	// Args are the --protect values given on the command line.  Each is
	// either a path or `разряд:<разряд>`.
	Args []string

	Home   string
	Config string
	Getenv func(string) string
}

// Load reads the защитный список: the file (named by --protect-file, else the
// user's own, else none at all) plus every --protect argument.  An absent file
// is not an error — protecting nothing is the honest default, and everything
// digitdisk would touch has a приговор behind it either way.
func Load(opt Options) (*List, error) {
	if opt.Getenv == nil {
		opt.Getenv = os.Getenv
	}
	if opt.Home == "" {
		if home, err := os.UserHomeDir(); err == nil {
			opt.Home = home
		}
	}
	l := &List{Rules: []Rule{}}

	file := opt.File
	var legacy bool
	if file == "" {
		file, legacy, _ = userPath(opt)
	}
	if file != "" {
		body, err := os.ReadFile(file)
		switch {
		case err == nil:
			if err := l.readFile(string(body), file, opt); err != nil {
				return nil, err
			}
			l.Origins = append(l.Origins, file)
			if legacy {
				l.Moved = movedNote(opt, file)
			}
		case opt.File != "":
			return nil, err
		case !os.IsNotExist(err):
			return nil, err
		}
	}

	for _, arg := range opt.Args {
		r, err := parseArg(arg, opt)
		if err != nil {
			return nil, fmt.Errorf("--protect %q: %s", arg, err)
		}
		l.Rules = append(l.Rules, r)
	}
	if len(opt.Args) > 0 {
		l.Origins = append(l.Origins, "--protect")
	}
	return l, nil
}

// userPath finds the защитный список: the new home first, the old one after.
func userPath(opt Options) (path string, legacy bool, ok bool) {
	if opt.Config != "" {
		p := filepath.Join(opt.Config, "digitdisk", FileName)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, true, true
		}
		return "", false, false
	}
	if opt.Home == "" {
		return "", false, false
	}
	return settings.Find(settings.Options{Home: opt.Home, Getenv: opt.Getenv}, FileName)
}

// movedNote is what a person is told when their защитный список came from the
// home it used to live in.
func movedNote(opt Options, from string) lang.Phrase {
	dir, err := settings.Dir(settings.Options{Home: opt.Home, Getenv: opt.Getenv})
	if err != nil {
		return lang.Phrase{}
	}
	return settings.MovedNote(from, filepath.Join(dir, FileName))
}

func (l *List) readFile(body, origin string, opt Options) error {
	sc := bufio.NewScanner(strings.NewReader(body))
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		f := strings.SplitN(text, "|", 3)
		for i := range f {
			f[i] = strings.TrimSpace(f[i])
		}
		if len(f) < 2 {
			return lang.Errorf("%s, строка %d: полей %d, а надо не меньше двух: вид|значение[|почему]", origin, line, len(f))
		}
		why := ""
		if len(f) == 3 {
			why = f[2]
		}
		r, err := rule(Kind(f[0]), f[1], why, opt)
		if err != nil {
			return lang.Errorf("%s, строка %d: %s", origin, line, err)
		}
		r.Origin, r.Line = origin, line
		l.Rules = append(l.Rules, r)
	}
	return sc.Err()
}

// parseArg reads one --protect value.  `разряд:кэш` names a разряд; anything
// else is a path.  The prefix is spelled out rather than guessed from the
// shape: a directory really can be called «кэш», and a rule that silently
// meant something else than it says is the one nobody debugs.
func parseArg(arg string, opt Options) (Rule, error) {
	if rest, ok := strings.CutPrefix(arg, string(KindClass)+":"); ok {
		r, err := rule(KindClass, strings.TrimSpace(rest), "", opt)
		if err != nil {
			return Rule{}, err
		}
		r.Origin = "--protect"
		return r, nil
	}
	r, err := rule(KindPath, arg, "", opt)
	if err != nil {
		return Rule{}, err
	}
	r.Origin = "--protect"
	return r, nil
}

func rule(kind Kind, value, why string, opt Options) (Rule, error) {
	r := Rule{Kind: kind, Value: value, Why: why}
	switch kind {
	case KindClass:
		for _, c := range core.Classes {
			if strings.EqualFold(value, string(c)) {
				r.class = c
				return r, nil
			}
		}
		return Rule{}, lang.Errorf("разряд %q не известен; есть: %s", value, joinClasses())
	case KindPath:
		if value == "" {
			return Rule{}, lang.Errorf("путь пуст")
		}
		switch {
		case strings.HasPrefix(value, "~/"):
			if opt.Home == "" {
				return Rule{}, lang.Errorf("путь начинается с ~, а домашний каталог не определяется")
			}
			r.chain = chainOf(filepath.Join(opt.Home, filepath.FromSlash(value[2:])))
		case strings.HasPrefix(value, "/"):
			abs, err := filepath.Abs(value)
			if err != nil {
				return Rule{}, err
			}
			r.chain = chainOf(abs)
		default:
			// No anchor: protect this chain of components wherever it is.
			r.anywhere = true
			r.chain = chainOf("/" + strings.Trim(filepath.ToSlash(value), "/"))
		}
		return r, nil
	}
	return Rule{}, lang.Errorf("вид %q: можно «путь» или «разряд»", kind)
}

func chainOf(abs string) string {
	return "/" + strings.Trim(filepath.ToSlash(filepath.Clean(abs)), "/") + "/"
}

func joinClasses() string {
	names := make([]string, 0, len(core.Classes))
	for _, c := range core.Classes {
		names = append(names, string(c))
	}
	return strings.Join(names, ", ")
}
