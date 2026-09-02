// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package places reads the справочник известных мест — the list of concrete
// directories that particular tools use for caches, logs and derived data —
// and turns it into the values the decision layer judges.
//
// # Why this is data and not code
//
// The rules in core/ know what a кэш IS: a path with a component called
// .cache, Caches, cache.  That is enough to recognise a cache in general and
// not enough to recognise npm's, whose store is ~/.npm/_cacache — no component
// named cache anywhere in it.  The missing knowledge is a LIST OF PLACES, and a
// list is data: it changes when a tool changes, not when a rule changes.
//
// So the list lives in places.conf next to this file, is embedded into the
// binary as the default, and is replaced whole by --places FILE or by
// ~/.config/digitdisk/places.conf.  Adding a place is a line in a file.  No
// rule in core/ is touched, and none of its постусловия move.
//
// # How a place reaches the verdict
//
// Each row resolves to a "цепь" — the place's absolute path with a slash at
// both ends, `/home/u/.npm/_cacache/` — and an anchor saying whether that chain
// must start the path or may appear anywhere in it.  The core matches the
// chain against the path with the slashes included, which is what makes the
// comparison a comparison of COMPONENTS: `/home/u/x.npm/_cacache/` does not
// contain `/home/u/.npm/_cacache/`, because there is no slash before `.npm`.
// The host does the resolving; the core does the judging, and refuses a
// справочник whose chains are not slash-bounded («Справочник ограничен»).
//
// # What the file may not say
//
// A place may only be кэш, журнал, сборка or загрузка.  «Крупное» is decided by
// size and «Неизвестное» means "no place matched"; letting a file assert either
// would be a lie about where the разряд came from, and the core rejects it.
package places

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"digitdisk/internal/core"
)

// builtin is the справочник shipped inside the binary.  A user file replaces
// it whole rather than adding to it: a справочник assembled from two halves
// would make "why was this proposed" a question about merge order.
//
//go:embed places.conf
var builtin string

// BuiltinName is what reports call the embedded справочник.
const BuiltinName = "встроенный"

// UserFile is where a user's own справочник is looked for, under the config
// directory ($XDG_CONFIG_HOME, else ~/.config).
const UserFile = "digitdisk/places.conf"

// Entry is one row of the справочник, resolved for this machine.
type Entry struct {
	Class  core.Class `json:"разряд"`
	Anchor string     `json:"якорь"`
	OS     string     `json:"система"`
	Path   string     `json:"путь"`
	Env    string     `json:"переменная,omitempty"`
	Name   string     `json:"имя"`
	Source string     `json:"источник"`
	Line   int        `json:"строка"`

	// Applies says the row is for this operating system.  Rows for the other
	// system are kept — `digitdisk places` shows the whole file, including
	// what would be looked at on a Mac — but never reach the decision layer.
	Applies bool `json:"на_этой_системе"`

	// Resolved is the absolute path this row points at here, empty for a
	// «всюду» row, which has no one place.
	Resolved string `json:"здесь,omitempty"`

	// Chain is what the core matches: slash-bounded, /home/u/.npm/_cacache/.
	Chain string `json:"цепь"`

	// Relocated says the row's основание came from its environment variable
	// rather than from the anchor.
	Relocated bool `json:"перенесено,omitempty"`
}

// Place returns the value the decision layer takes.
func (e Entry) Place() core.Place {
	anchor := core.AnchorRoot
	if e.Anchor == anchorAnywhere {
		anchor = core.AnchorAnywhere
	}
	return core.Place{Class: e.Class, Anchor: anchor, Chain: e.Chain}
}

// Directory is a whole справочник as loaded.
type Directory struct {
	Origin  string  `json:"откуда"`
	Entries []Entry `json:"места"`

	// applicable is the rows for this system, kept because Match is called
	// once for every item that makes a plan and rebuilding the slice there
	// would allocate the whole справочник per item.
	applicable []Entry
}

// Applicable returns the rows that apply to this operating system, in file
// order — the order is part of the meaning, since the first matching place
// wins.
func (d *Directory) Applicable() []Entry {
	if d.applicable == nil {
		out := make([]Entry, 0, len(d.Entries))
		for _, e := range d.Entries {
			if e.Applies {
				out = append(out, e)
			}
		}
		d.applicable = out
	}
	return d.applicable
}

// Places returns the value list for the decision layer.
func (d *Directory) Places() []core.Place {
	rows := d.Applicable()
	out := make([]core.Place, 0, len(rows))
	for _, e := range rows {
		out = append(out, e.Place())
	}
	return out
}

// Match reports which row of the справочник claims a path, and whether any
// does.  It repeats, in Go, the comparison the core makes — and is used only
// to NAME the place in a report, never to decide anything.  Whichever row it
// names, the разряд on the item came from the core.
func (d *Directory) Match(path string) (Entry, bool) {
	trail := path + "/"
	for _, e := range d.Applicable() {
		if e.Anchor == anchorAnywhere {
			if strings.Contains(trail, e.Chain) {
				return e, true
			}
			continue
		}
		if strings.HasPrefix(trail, e.Chain) {
			return e, true
		}
	}
	return Entry{}, false
}

const (
	anchorHome     = "дом"
	anchorCache    = "кэш"
	anchorData     = "данные"
	anchorConfig   = "настройки"
	anchorRoot     = "корень"
	anchorAnywhere = "всюду"
)

var classByWord = map[string]core.Class{
	"кэш":      core.ClassCache,
	"журнал":   core.ClassLog,
	"сборка":   core.ClassBuild,
	"загрузка": core.ClassDownload,
}

// Options configures loading.
type Options struct {
	// File, when set, is the справочник to read instead of the user file and
	// the embedded default.
	File string
	// Off loads nothing: the decision layer then judges by приметы alone,
	// exactly as it did before the справочник existed.
	Off bool

	Home   string              // defaults to the user's home directory
	GOOS   string              // defaults to runtime.GOOS
	Getenv func(string) string // defaults to os.Getenv
	Config string              // config dir; defaults to os.UserConfigDir
}

// Load reads the справочник: the file named by --places if there is one, else
// the user's own file if it exists, else the copy embedded in the binary.
func Load(opt Options) (*Directory, error) {
	if opt.Off {
		return &Directory{Origin: "выключен ключом --no-places"}, nil
	}
	if opt.Getenv == nil {
		opt.Getenv = os.Getenv
	}
	if opt.GOOS == "" {
		opt.GOOS = runtime.GOOS
	}
	if opt.Home == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("домашний каталог не определяется, а справочник считает места от него: %w", err)
		}
		opt.Home = home
	}

	if opt.File != "" {
		body, err := os.ReadFile(opt.File)
		if err != nil {
			return nil, err
		}
		return parse(string(body), opt.File, opt)
	}

	if user := userPath(opt); user != "" {
		body, err := os.ReadFile(user)
		switch {
		case err == nil:
			return parse(string(body), user, opt)
		case !os.IsNotExist(err):
			return nil, err
		}
	}
	return parse(builtin, BuiltinName, opt)
}

func userPath(opt Options) string {
	dir := opt.Config
	if dir == "" {
		if v := opt.Getenv("XDG_CONFIG_HOME"); v != "" {
			dir = v
		} else {
			dir = filepath.Join(opt.Home, ".config")
		}
	}
	return filepath.Join(dir, filepath.FromSlash(UserFile))
}

// parse turns the file into resolved entries, refusing a malformed row rather
// than skipping it: a справочник quietly missing a line is a справочник nobody
// can check against the file they edited.
func parse(body, origin string, opt Options) (*Directory, error) {
	d := &Directory{Origin: origin, Entries: []Entry{}}
	sc := bufio.NewScanner(strings.NewReader(body))
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		e, err := row(text, line, opt)
		if err != nil {
			return nil, fmt.Errorf("%s, строка %d: %w", origin, line, err)
		}
		d.Entries = append(d.Entries, e)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return d, nil
}

func row(text string, line int, opt Options) (Entry, error) {
	f := strings.Split(text, "|")
	if len(f) != 7 {
		return Entry{}, fmt.Errorf("полей %d, а надо 7: разряд|якорь|система|путь|переменная|имя|источник", len(f))
	}
	for i := range f {
		f[i] = strings.TrimSpace(f[i])
	}
	e := Entry{
		Anchor: f[1], OS: f[2], Path: f[3], Env: f[4], Name: f[5], Source: f[6], Line: line,
	}

	class, ok := classByWord[f[0]]
	if !ok {
		return Entry{}, fmt.Errorf("разряд %q справочнику не позволен; можно: кэш, журнал, сборка, загрузка", f[0])
	}
	e.Class = class

	switch e.OS {
	case "linux", "macos", "все":
	default:
		return Entry{}, fmt.Errorf("система %q: можно linux, macos, все", e.OS)
	}
	e.Applies = e.OS == "все" || (e.OS == "linux" && opt.GOOS == "linux") || (e.OS == "macos" && opt.GOOS == "darwin")

	if e.Name == "" {
		return Entry{}, fmt.Errorf("у места нет имени: непонятно, чей это каталог")
	}
	if e.Source == "" {
		return Entry{}, fmt.Errorf("у места %q нет источника; в этом справочнике строка без ссылки на документацию инструмента — не строка", e.Name)
	}

	base, tail, split := strings.Cut(e.Path, "//")
	if base == "" {
		return Entry{}, fmt.Errorf("путь пуст")
	}
	if strings.HasPrefix(e.Path, "/") && e.Anchor != anchorRoot {
		return Entry{}, fmt.Errorf("путь начинается с косой, а якорь %q — это можно только с якорем «корень»", e.Anchor)
	}

	if e.Env != "" {
		if v := strings.TrimSpace(opt.Getenv(e.Env)); v != "" {
			if !strings.HasPrefix(v, "/") {
				return Entry{}, fmt.Errorf("переменная %s задана как %q — не абсолютный путь; место пропущено бы молча, а это отказ", e.Env, v)
			}
			e.Relocated = true
			e.Resolved = filepath.Clean(joinTail(v, tail))
			e.Chain = chainOf(e.Resolved)
			return e, nil
		}
	}
	if !split {
		tail = ""
	}

	switch e.Anchor {
	case anchorAnywhere:
		if e.Env != "" {
			return Entry{}, fmt.Errorf("у якоря «всюду» переменной быть не может: место не привязано ни к какому основанию")
		}
		e.Chain = chainOf("/" + strings.Trim(e.Path, "/"))
		return e, nil
	case anchorRoot:
		e.Resolved = filepath.Clean("/" + strings.TrimPrefix(joinTail(base, tail), "/"))
	case anchorHome:
		e.Resolved = filepath.Clean(filepath.Join(opt.Home, filepath.FromSlash(joinTail(base, tail))))
	case anchorCache:
		e.Resolved = filepath.Clean(filepath.Join(xdg(opt, "XDG_CACHE_HOME", ".cache"), filepath.FromSlash(joinTail(base, tail))))
	case anchorData:
		e.Resolved = filepath.Clean(filepath.Join(xdg(opt, "XDG_DATA_HOME", filepath.Join(".local", "share")), filepath.FromSlash(joinTail(base, tail))))
	case anchorConfig:
		e.Resolved = filepath.Clean(filepath.Join(xdg(opt, "XDG_CONFIG_HOME", ".config"), filepath.FromSlash(joinTail(base, tail))))
	default:
		return Entry{}, fmt.Errorf("якорь %q: можно дом, кэш, данные, настройки, корень, всюду", e.Anchor)
	}
	e.Chain = chainOf(e.Resolved)
	return e, nil
}

func joinTail(base, tail string) string {
	if tail == "" {
		return base
	}
	return strings.TrimSuffix(base, "/") + "/" + tail
}

func xdg(opt Options, name, fallback string) string {
	if v := strings.TrimSpace(opt.Getenv(name)); v != "" && strings.HasPrefix(v, "/") {
		return v
	}
	return filepath.Join(opt.Home, fallback)
}

// chainOf puts the slashes on: they are what makes the core's comparison a
// comparison of whole components.
func chainOf(abs string) string {
	return "/" + strings.Trim(filepath.ToSlash(abs), "/") + "/"
}

// Found is what one place turned out to be on this machine.
type Found struct {
	Entry
	Exists bool   `json:"есть"`
	Dir    bool   `json:"каталог"`
	Bytes  int64  `json:"байт"`
	Files  int    `json:"файлов"`
	Note   string `json:"замечание,omitempty"`
}

// Look measures the places that name one directory on this machine.  A
// «всюду» row names no single place and is reported as such rather than as
// missing: it is not absent, it is everywhere or nowhere.
func (d *Directory) Look(measure func(string) (int64, int, error)) []Found {
	out := make([]Found, 0, len(d.Entries))
	for _, e := range d.Entries {
		f := Found{Entry: e}
		switch {
		case !e.Applies:
			f.Note = "другая система"
		case e.Resolved == "":
			f.Note = "на любой глубине — одного места нет"
		default:
			info, err := os.Lstat(e.Resolved)
			switch {
			case err != nil:
				f.Note = "нет"
			case !info.IsDir():
				f.Exists = true
				f.Note = "не каталог"
			default:
				f.Exists, f.Dir = true, true
				if measure != nil {
					bytes, files, err := measure(e.Resolved)
					if err != nil {
						f.Note = err.Error()
					}
					f.Bytes, f.Files = bytes, files
				}
			}
		}
		out = append(out, f)
	}
	return out
}

// ByClass counts the rows of each разряд, for a one-line summary.
func (d *Directory) ByClass() []string {
	n := map[core.Class]int{}
	for _, e := range d.Entries {
		n[e.Class]++
	}
	var out []string
	for _, c := range core.Classes {
		if n[c] > 0 {
			out = append(out, fmt.Sprintf("%s %d", c, n[c]))
		}
	}
	return out
}
