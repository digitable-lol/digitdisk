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
	"digitdisk/internal/lang"
)

// Здесь сверяется не поведение, а то, что три перечня — код, справка и
// страница руководства — говорят об одном и том же наборе команд и ключей.
// Сверка идёт В ОБЕ СТОРОНЫ: ключ, заведённый в коде и не попавший на
// страницу, и ключ, обещанный страницей и не заведённый в коде, — обе беды
// одинаково молчаливы, и обе ловятся здесь.

// manPages — обе страницы руководства. Их две, и сверка идёт по обеим:
// страница, отставшая от кода, одинаково молчалива на любом языке, а страница,
// отставшая от ВТОРОЙ страницы, — это ещё и две разные правды об одном
// инструменте. Ставятся они так:
//
//	digitdisk.en.1 → share/man/man1/digitdisk.1     — её находит `man digitdisk`
//	digitdisk.1    → share/man/ru/man1/digitdisk.1  — её находит `LANG=ru_RU… man digitdisk`
var manPages = []string{"../digitdisk.1", "../digitdisk.en.1"}

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
func manFlags(t *testing.T, manPage string) map[string]bool {
	t.Helper()
	out := longFlags(readPage(t, manPage))
	for _, m := range regexp.MustCompile(`(?:^|[ .])Fl (-?[A-Za-z][a-z0-9-]*)`).FindAllStringSubmatch(readPage(t, manPage), -1) {
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
	help := longFlags(cli.Usage(lang.RU) + cli.Usage(lang.EN))
	if lost := missing(code, help); len(lost) > 0 {
		t.Errorf("ключи заведены в коде, но не названы в справке: %v", lost)
	}
	if extra := missing(help, code); len(extra) > 0 {
		t.Errorf("справка обещает ключи, которых нет в коде: %v", extra)
	}
	t.Logf("ключей в коде %d, в справке %d: %v", len(code), len(help), keys(code))
}

func TestКлючиКодаИСтраницыСовпадают(t *testing.T) {
	for _, manPage := range manPages {
		ключиСтраницы(t, manPage)
	}
}

func ключиСтраницы(t *testing.T, manPage string) {
	code := registered(t)
	page := manFlags(t, manPage)
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
	for _, l := range []lang.Lang{lang.RU, lang.EN} {
		if help := cli.Usage(l); !strings.Contains(help, "digitdisk") {
			t.Fatalf("справка на %s пуста", l)
		}
	}
	for _, manPage := range manPages {
		подкомандыСтраницы(t, manPage)
	}
}

func подкомандыСтраницы(t *testing.T, manPage string) {
	page := []byte(readPage(t, manPage))
	help := cli.Usage(lang.RU) + cli.Usage(lang.EN)
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
	re := regexp.MustCompile("(?m)^\\.It Cm ([a-z]+)")
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
	if !strings.Contains(cli.Usage(lang.RU), "без подкоманды — "+cli.Default) {
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

// readPage читает страницу руководства и роняет прогон, если её нет: страница,
// которой нет, — это не «перевод ещё не написан», а обещание формулы, которое
// некому выполнить.
func readPage(t *testing.T, manPage string) string {
	t.Helper()
	body, err := os.ReadFile(manPage)
	if err != nil {
		t.Fatalf("%s не читается: %v", manPage, err)
	}
	return string(body)
}

// TestОбеСтраницыОписываютОдноИТоЖе сверяет две редакции руководства между
// собой, а не только каждую с кодом.
//
// Расхождение в фактах хуже отсутствующего перевода: читатель английской
// страницы не знает, что русская обещает другое, и узнаёт об этом от машины,
// которая повела себя иначе. Поэтому сверяется состав: те же подкоманды, те же
// ключи, те же файлы.
func TestОбеСтраницыОписываютОдноИТоЖе(t *testing.T) {
	ru, en := readPage(t, manPages[0]), readPage(t, manPages[1])

	ruFlags, enFlags := manFlags(t, manPages[0]), manFlags(t, manPages[1])
	if lost := missing(ruFlags, enFlags); len(lost) > 0 {
		t.Errorf("русская страница называет ключи, которых нет в английской: %v", lost)
	}
	if lost := missing(enFlags, ruFlags); len(lost) > 0 {
		t.Errorf("английская страница называет ключи, которых нет в русской: %v", lost)
	}

	// Имя подкоманды на строке .It — с доводом за ним или без: обе
	// страницы пишут его одинаково, но довод у них на своём языке.
	subs := regexp.MustCompile("(?m)^\\.It Cm ([a-z]+)")
	names := func(page string) map[string]bool {
		out := map[string]bool{}
		for _, m := range subs.FindAllStringSubmatch(page, -1) {
			out[m[1]] = true
		}
		return out
	}
	ruNames, enNames := names(ru), names(en)
	if lost := missing(ruNames, enNames); len(lost) > 0 {
		t.Errorf("подкоманды есть на русской странице и нет на английской: %v", lost)
	}
	if lost := missing(enNames, ruNames); len(lost) > 0 {
		t.Errorf("подкоманды есть на английской странице и нет на русской: %v", lost)
	}
	t.Logf("страниц %d, подкоманд на каждой %d, ключей на каждой %d", len(manPages), len(ruNames), len(ruFlags))
}

// TestКлючЯзыкаЗаведёнИНазван — тот самый ключ, ради которого всё это.
func TestКлючЯзыкаЗаведёнИНазван(t *testing.T) {
	if !registered(t)["--lang"] {
		t.Fatal("ключ --lang не заведён ни в одной подкоманде")
	}
	for _, l := range []lang.Lang{lang.RU, lang.EN} {
		if !strings.Contains(cli.Usage(l), "--lang") {
			t.Errorf("справка на %s не называет --lang", l)
		}
	}
	for _, manPage := range manPages {
		if !strings.Contains(readPage(t, manPage), "-lang") {
			t.Errorf("%s не называет --lang", manPage)
		}
	}
}

// TestСправкаЕстьНаОбоихЯзыках: пустой или одинаковый вывод на двух языках —
// это непереведённая справка, притворившаяся переведённой.
func TestСправкаЕстьНаОбоихЯзыках(t *testing.T) {
	ru, en := cli.Usage(lang.RU), cli.Usage(lang.EN)
	if ru == en {
		t.Fatal("справка на двух языках вышла одинаковой")
	}
	for _, c := range cli.Commands {
		if !strings.Contains(en, c.Name) {
			t.Errorf("английская справка не называет подкоманду %q", c.Name)
		}
	}
	ruLines, enLines := strings.Count(ru, "\n"), strings.Count(en, "\n")
	if ruLines != enLines {
		t.Errorf("строк в справке: по-русски %d, по-английски %d — редакции разъехались", ruLines, enLines)
	}
	t.Logf("строк в справке %d на каждом языке", ruLines)
}
