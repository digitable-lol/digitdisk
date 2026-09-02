// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"digitdisk/internal/clean"
	"digitdisk/internal/core"
	"digitdisk/internal/lang"
)

func bigPlan(n int) clean.Plan {
	p := clean.Plan{
		Root: "/home/u", Trash: "/home/u/.digitdisk-trash",
		Decider: "испытательный слой", DeciderReady: true, ContractVersion: core.ContractVersion,
		Items: []clean.Item{}, Refused: []clean.Refusal{}, Protected: []clean.Protected{},
	}
	for i := 0; i < n; i++ {
		p.Items = append(p.Items, clean.Item{
			Path: fmt.Sprintf("/home/u/.cache/x/%04d.bin", i), Size: 1000,
			Class: core.ClassCache, Verdict: core.VerdictRemovable, Weight: 1000,
		})
		p.Bytes += 1000
	}
	p.FreeableBytes = p.Bytes
	p.ByClass = []clean.ClassSum{{Class: core.ClassCache, Count: n, Bytes: p.Bytes}}
	return p
}

func render(t *testing.T, l lang.Lang, p clean.Plan, top int) string {
	t.Helper()
	var b bytes.Buffer
	CleanPlan(&b, l, p, top)
	return b.String()
}

// TestTopCutsTheListAndNotTheCount is the whole promise of --top: a plan over a
// real home directory is hundreds of thousands of lines long, and the answer to
// that is a shorter LIST, never a smaller COUNT.
//
// It runs in both languages because the counts are written by the language
// too — the separator between the thousands is a comma in one and a space in
// the other — and a summary that came out right in Russian and wrong in
// English would still be a wrong summary.
func TestTopCutsTheListAndNotTheCount(t *testing.T) {
	for _, l := range both {
		p := bigPlan(400)
		full := render(t, l, p, 0)
		short := render(t, l, p, 15)

		if strings.Count(short, "\n") >= strings.Count(full, "\n") {
			t.Fatalf("[%s] --top 15 не укоротил вывод: %d строк против %d",
				l, strings.Count(short, "\n"), strings.Count(full, "\n"))
		}
		for _, must := range []string{
			l.F("К ПЕРЕНОСУ В КОРЗИНУ  %s файлов, %s", l.Num(400), l.Bytes(400000)),
			l.F("  %-11s %6s файлов  %10s", l.Word(string(core.ClassCache)), l.Num(400), l.Bytes(400000)),
		} {
			if !strings.Contains(short, must) {
				t.Errorf("[%s] укороченный вывод потерял счёт: нет строки %q", l, must)
			}
			if !strings.Contains(full, must) {
				t.Errorf("[%s] полный вывод потерял счёт: нет строки %q", l, must)
			}
		}
		if want := l.F("  …и ещё %s файлов на %s — весь список: --top 0, или --json",
			l.Num(385), l.Bytes(385000)); !strings.Contains(short, want) {
			t.Errorf("[%s] укороченный вывод не сказал, чего не показал: нет %q\n%s", l, want, short)
		}
		// "…and 385 more" in either language starts the same way; a full
		// listing must not start it at all.
		more := strings.SplitN(l.T("  …и ещё %s файлов на %s — весь список: --top 0, или --json"), "%s", 2)[0]
		if strings.Contains(full, more) {
			t.Errorf("[%s] --top 0 обязан печатать всё и ни о чём не умалчивать", l)
		}
		if got := strings.Count(short, ".bin"); got != 15 {
			t.Errorf("[%s] строк перечня %d, ждали 15", l, got)
		}
		if got := strings.Count(full, ".bin"); got != 400 {
			t.Errorf("[%s] строк перечня при --top 0: %d, ждали 400", l, got)
		}
	}
}

// TestProtectedIsItsOwnSection: защита и отказ — разные вещи, и в отчёте они
// разные разделы. Смешать их значило бы спрятать расхождение слоёв среди
// распоряжений человека.
func TestProtectedIsItsOwnSection(t *testing.T) {
	for _, l := range both {
		p := bigPlan(2)
		p.Protected = []clean.Protected{{Path: "/home/u/projects/a", Class: core.ClassCache, Size: 10}}
		p.ProtectedBytes = 10
		p.Refused = []clean.Refusal{{
			Path: "/home/u/b", Class: core.ClassCache, Verdict: core.VerdictRemovable,
			Reason: lang.Say("это каталог: снос каталога рекурсивен, и ядро его «МожноУбрать» не выдаёт (правило П3)"),
		}}
		out := render(t, l, p, 0)

		protectedHead := strings.SplitN(l.T("ЗАЩИЩЕНО  %s файлов, %s: ядро назвало их «%s», защитный список запретил"), " ", 2)[0]
		refusedHead := strings.SplitN(l.T("ОТКАЗАНО  %s записей: ядро назвало их «%s», хозяин не тронет"), " ", 2)[0]
		if !strings.Contains(out, l.F("ЗАЩИЩЕНО  %s файлов, %s: ядро назвало их «%s», защитный список запретил",
			l.Num(1), l.Bytes(10), l.Word(string(core.VerdictRemovable)))) {
			t.Errorf("[%s] нет раздела ЗАЩИЩЕНО:\n%s", l, out)
		}
		if !strings.Contains(out, l.F("ОТКАЗАНО  %s записей: ядро назвало их «%s», хозяин не тронет",
			l.Num(1), l.Word(string(core.VerdictRemovable)))) {
			t.Errorf("[%s] нет раздела ОТКАЗАНО:\n%s", l, out)
		}
		if strings.Index(out, protectedHead) > strings.Index(out, refusedHead) {
			t.Errorf("[%s] ОТКАЗАНО обязано идти последним: это расхождение слоёв, и его читают внимательнее", l)
		}
		// The refusal is a Phrase: it travels to the JSON in Russian and
		// reaches the screen in the reader's language.
		if !strings.Contains(out, p.Refused[0].Reason.In(l)) {
			t.Errorf("[%s] отказ напечатан не на языке читателя:\n%s", l, out)
		}
	}
}

// TestАнглийскийПланБезКириллицы: план, показанный английскому читателю, не
// должен содержать ни одной русской буквы, которую написал этот инструмент.
// Все значения здесь — пути и имя слоя — латиницей нарочно, поэтому любая
// кириллица в выводе означает строку, не дошедшую до словаря.
func TestАнглийскийПланБезКириллицы(t *testing.T) {
	p := bigPlan(3)
	p.Decider = "test layer"
	p.PlacesOrigin = "/etc/digitdisk/places.conf"
	p.PlacesCount = 12
	p.ProtectOrigins = []string{"/home/u/.config/digitdisk/protect.conf"}
	p.HardlinkItems = 1
	p.PrunedTrash = 2
	p.Walk.Entries = 100
	p.Walk.TotalBytes = 1 << 20
	p.Items[0].Place = "cache of the package manager"
	p.Refused = []clean.Refusal{{
		Path: "/home/u/b", Class: core.ClassCache, Verdict: core.VerdictRemovable,
		Reason: lang.Raw("is a directory"),
	}}
	out := render(t, lang.EN, p, 0)
	for i, line := range strings.Split(out, "\n") {
		for _, r := range line {
			if r >= 0x0400 && r <= 0x04FF {
				t.Fatalf("строка %d английского плана осталась по-русски: %q", i+1, line)
			}
		}
	}
}
