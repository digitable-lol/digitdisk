// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package protect

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"digitdisk/internal/core"
	"digitdisk/internal/lang"
)

func load(t *testing.T, body string, args ...string) *List {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "protect.conf")
	if body != "" {
		if err := os.WriteFile(file, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	} else {
		file = ""
	}
	l, err := Load(Options{File: file, Args: args, Home: "/home/u", Config: dir, Getenv: func(string) string { return "" }})
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// TestPathRuleIsByComponents restates the September trap for the защитный
// список: «не трогай ~/projects» must not also protect ~/projects-old, and
// must protect everything under ~/projects.
func TestPathRuleIsByComponents(t *testing.T) {
	l := load(t, "путь|~/projects|мои исходники\n")
	for _, p := range []string{"/home/u/projects", "/home/u/projects/a/b.txt"} {
		if _, ok := l.Covers(p, core.ClassCache); !ok {
			t.Errorf("%q должен быть защищён", p)
		}
	}
	for _, p := range []string{"/home/u/projects-old/a", "/home/u/proj", "/srv/home/u/projects/a"} {
		if _, ok := l.Covers(p, core.ClassCache); ok {
			t.Errorf("%q защищён, а не должен", p)
		}
	}
}

// TestRelativeRuleProtectsEverywhere is the other half: a rule written without
// an anchor means "this chain of components, wherever it is".
func TestRelativeRuleProtectsEverywhere(t *testing.T) {
	l := load(t, "путь|node_modules/.bin|запускаемое\n")
	if _, ok := l.Covers("/srv/p/node_modules/.bin/tsc", core.ClassBuild); !ok {
		t.Error("цепочка составляющих должна ловиться на любой глубине")
	}
	if _, ok := l.Covers("/srv/p/my_node_modules/.bin/tsc", core.ClassBuild); ok {
		t.Error("составляющая, а не подстрока")
	}
}

func TestClassRule(t *testing.T) {
	l := load(t, "разряд|журнал|журналы разбираю сам\n")
	if _, ok := l.Covers("/var/log/x.log", core.ClassLog); !ok {
		t.Error("разряд не защищён")
	}
	if _, ok := l.Covers("/var/log/x.log", core.ClassCache); ok {
		t.Error("защита разряда протекла на другой разряд")
	}
}

// TestPathBeatsClass: a person asking "why was this kept" is better served by
// the exact path they wrote than by a разряд that also matched.
func TestPathBeatsClass(t *testing.T) {
	l := load(t, "разряд|кэш|весь кэш\nпуть|~/projects|исходники\n")
	r, ok := l.Covers("/home/u/projects/x", core.ClassCache)
	if !ok || r.Kind != KindPath {
		t.Errorf("сработало правило %+v, ждали правило по пути", r)
	}
}

func TestArgs(t *testing.T) {
	l := load(t, "", "/srv/data", "разряд:сборка", "node_modules")
	if len(l.Rules) != 3 {
		t.Fatalf("правил %d, ждали 3", len(l.Rules))
	}
	if _, ok := l.Covers("/srv/data/x", core.ClassCache); !ok {
		t.Error("путь из ключа не защищает")
	}
	if _, ok := l.Covers("/x/y", core.ClassBuild); !ok {
		t.Error("разряд из ключа не защищает")
	}
	for _, r := range l.Rules {
		if r.Where() != "--protect" {
			t.Errorf("правило %q не помнит, откуда оно: %q", r.Value, r.Where())
		}
	}
}

// TestEmptyProtectsNothing states the default: an absent защитный список is
// not an error and protects nothing, because everything digitdisk touches
// already has a приговор behind it.
func TestEmptyProtectsNothing(t *testing.T) {
	l := load(t, "")
	if !l.Empty() {
		t.Error("пустой список обязан быть пустым")
	}
	if _, ok := l.Covers("/home/u/.cache/x", core.ClassCache); ok {
		t.Error("пустой список что-то защитил")
	}
	var nilList *List
	if _, ok := nilList.Covers("/x", core.ClassCache); ok {
		t.Error("отсутствующий список обязан ничего не защищать и не падать")
	}
}

func TestBadRulesAreRefused(t *testing.T) {
	dir := t.TempDir()
	for _, bad := range []string{"разряд|небывалый|\n", "нетакой|x|\n", "путь||\n", "одно поле\n"} {
		file := filepath.Join(dir, "protect.conf")
		if err := os.WriteFile(file, []byte(bad), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(Options{File: file, Home: "/home/u", Config: dir, Getenv: func(string) string { return "" }}); err == nil {
			t.Errorf("правило %q принято, а должно быть отвергнуто", bad)
		}
	}
}

// TestRuleReadsInBothLanguages: правило называет свой вид словом читателя, а
// String остаётся прежним — его видят журнал и сравнения по тексту.
func TestRuleReadsInBothLanguages(t *testing.T) {
	l := load(t, "путь|~/projects|мои исходники\n")
	r := l.Rules[0]
	if got, want := r.String(), r.In(lang.RU); got != want {
		t.Errorf("по-русски правило читается двояко: %q и %q", got, want)
	}
	if !strings.HasPrefix(r.String(), "путь ") {
		t.Errorf("String() перестал начинаться с вида по-русски: %q", r.String())
	}
	en := r.In(lang.EN)
	if !strings.HasPrefix(en, "path ") {
		t.Errorf("вид правила не назван по-английски: %q", en)
	}
	if !strings.Contains(en, "мои исходники") {
		t.Errorf("причина, написанная человеком, обязана остаться его словами: %q", en)
	}
}
