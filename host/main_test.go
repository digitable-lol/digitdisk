// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"digitdisk/internal/cli"
)

// Здесь сверяется не поведение, а то, что три перечня — код, справка и
// страница руководства — говорят об одном и том же наборе команд и ключей.
// Сверка идёт В ОБЕ СТОРОНЫ: ключ, заведённый в коде и не попавший на
// страницу, и ключ, обещанный страницей и не заведённый в коде, — обе беды
// одинаково молчаливы, и обе ловятся здесь.

const manPage = "../digitdisk.1"

// registered находит ключи, заведённые в main.go, — там, где они и заводятся,
// а не там, где о них рассказано.
func registered(t *testing.T) map[string]bool {
	t.Helper()
	src, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatalf("main.go не читается: %v", err)
	}
	out := map[string]bool{}
	re := regexp.MustCompile(`fs\.(?:Bool|Int|Int64|Float64|String|Duration)\("([a-z][a-z0-9-]*)"`)
	for _, m := range re.FindAllStringSubmatch(string(src), -1) {
		out["--"+m[1]] = true
	}
	reVar := regexp.MustCompile(`fs\.Var\(&[A-Za-z0-9_]+, "([a-z][a-z0-9-]*)"`)
	for _, m := range reVar.FindAllStringSubmatch(string(src), -1) {
		out["--"+m[1]] = true
	}
	// Ключи, которые разбирает сам run, а не набор флагов подкоманды.
	for _, list := range [][]string{cli.HelpArgs, cli.VersionArgs} {
		for _, a := range list {
			if strings.HasPrefix(a, "--") {
				out[a] = true
			}
		}
	}
	if len(out) < 10 {
		t.Fatalf("в main.go нашлось %d ключей — разбор сломался", len(out))
	}
	return out
}

// longFlags собирает длинные ключи из обычного текста.
func longFlags(text string) map[string]bool {
	out := map[string]bool{}
	for _, m := range regexp.MustCompile(`--[a-z][a-z0-9-]*`).FindAllString(text, -1) {
		out[m] = true
	}
	return out
}

// manFlags собирает ключи со страницы руководства: и записанные макросом
// (`.Fl -json`), и набранные в примерах как есть (`--json`).
func manFlags(t *testing.T) map[string]bool {
	t.Helper()
	page, err := os.ReadFile(manPage)
	if err != nil {
		t.Fatalf("%s не читается: %v", manPage, err)
	}
	out := longFlags(string(page))
	for _, m := range regexp.MustCompile(`(?:^|[ .])Fl (-?[A-Za-z][a-z0-9-]*)`).FindAllStringSubmatch(string(page), -1) {
		out["-"+m[1]] = true
	}
	return out
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// missing — то, что есть в have и чего нет в want.
func missing(have, want map[string]bool) []string {
	var out []string
	for k := range have {
		if !want[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func TestКлючиКодаИСправкиСовпадают(t *testing.T) {
	code := registered(t)
	help := longFlags(cli.Usage())
	if lost := missing(code, help); len(lost) > 0 {
		t.Errorf("ключи заведены в коде, но не названы в справке: %v", lost)
	}
	if extra := missing(help, code); len(extra) > 0 {
		t.Errorf("справка обещает ключи, которых нет в коде: %v", extra)
	}
	t.Logf("ключей в коде %d, в справке %d: %v", len(code), len(help), keys(code))
}

func TestКлючиКодаИСтраницыСовпадают(t *testing.T) {
	code := registered(t)
	page := manFlags(t)
	// Короткие написания разбирает run, и страница обязана их назвать.
	for _, short := range []string{"-h", "-V"} {
		if !page[short] {
			t.Errorf("%s не называет короткий ключ %s", manPage, short)
		}
		delete(page, short)
	}
	if lost := missing(code, page); len(lost) > 0 {
		t.Errorf("ключи заведены в коде, но не названы в %s: %v", manPage, lost)
	}
	if extra := missing(page, code); len(extra) > 0 {
		t.Errorf("%s обещает ключи, которых нет в коде: %v", manPage, extra)
	}
	t.Logf("ключей в коде %d, на странице %d", len(code), len(page))
}

func TestПодкомандыОдниИТеЖеВезде(t *testing.T) {
	page, err := os.ReadFile(manPage)
	if err != nil {
		t.Fatalf("%s не читается: %v", manPage, err)
	}
	help := cli.Usage()
	for _, c := range cli.Commands {
		if _, ok := handlers[c.Name]; !ok {
			t.Errorf("подкоманда %q объявлена в internal/cli, но не разобрана в main", c.Name)
		}
		if !strings.Contains(help, c.Name) {
			t.Errorf("подкоманда %q не названа в справке", c.Name)
		}
		if !strings.Contains(string(page), ".Nm digitdisk Cm "+c.Name) {
			t.Errorf("подкоманды %q нет в синопсисе %s", c.Name, manPage)
		}
		if !strings.Contains(string(page), ".It Cm "+c.Name) {
			t.Errorf("подкоманды %q нет в разделе SUBCOMMANDS %s", c.Name, manPage)
		}
	}
	for name := range handlers {
		if !cli.Known(name) {
			t.Errorf("main разбирает подкоманду %q, которой нет в internal/cli", name)
		}
	}
	// Обратная сторона: страница не называет подкоманд, которых нет.
	re := regexp.MustCompile(`(?m)^\.It Cm ([a-z]+)$`)
	for _, m := range re.FindAllStringSubmatch(string(page), -1) {
		if !cli.Known(m[1]) {
			t.Errorf("%s описывает подкоманду %q, которой нет в коде", manPage, m[1])
		}
	}
	t.Logf("подкоманд %d, все на месте в коде, справке и %s", len(cli.Commands), manPage)
}

func TestГолыйВызовЭтоУмолчание(t *testing.T) {
	if !cli.Known(cli.Default) {
		t.Fatalf("умолчание %q не подкоманда", cli.Default)
	}
	if _, ok := handlers[cli.Default]; !ok {
		t.Fatalf("умолчание %q не разобрано в main", cli.Default)
	}
	if !strings.Contains(cli.Usage(), "без подкоманды — "+cli.Default) {
		t.Error("справка не называет подкоманду по умолчанию")
	}
}

func TestНеизвестнаяПодкомандаОтвергается(t *testing.T) {
	// Слово, которого нет, обязано быть отказом с кодом 2, а не тихим
	// исполнением умолчания.
	if code := runQuiet(t, []string{"фигня"}); code != 2 {
		t.Errorf("digitdisk фигня вернул %d, ожидался 2", code)
	}
	if code := runQuiet(t, []string{"--help"}); code != 0 {
		t.Errorf("digitdisk --help вернул %d, ожидался 0", code)
	}
}

// runQuiet выполняет run с отведёнными в никуда потоками вывода.
func runQuiet(t *testing.T, args []string) int {
	t.Helper()
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("%s не открывается: %v", os.DevNull, err)
	}
	defer null.Close()
	outWas, errWas := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = null, null
	defer func() { os.Stdout, os.Stderr = outWas, errWas }()
	return run(args)
}
