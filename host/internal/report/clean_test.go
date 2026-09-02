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

func render(t *testing.T, p clean.Plan, top int) string {
	t.Helper()
	var b bytes.Buffer
	CleanPlan(&b, p, top)
	return b.String()
}

// TestTopCutsTheListAndNotTheCount is the whole promise of --top: a plan over a
// real home directory is hundreds of thousands of lines long, and the answer to
// that is a shorter LIST, never a smaller COUNT.
func TestTopCutsTheListAndNotTheCount(t *testing.T) {
	p := bigPlan(400)
	full := render(t, p, 0)
	short := render(t, p, 15)

	if strings.Count(short, "\n") >= strings.Count(full, "\n") {
		t.Fatalf("--top 15 не укоротил вывод: %d строк против %d",
			strings.Count(short, "\n"), strings.Count(full, "\n"))
	}
	for _, must := range []string{
		"К ПЕРЕНОСУ В КОРЗИНУ  400 файлов",
		"Кэш            400 файлов",
	} {
		if !strings.Contains(short, must) {
			t.Errorf("укороченный вывод потерял счёт: нет строки %q", must)
		}
		if !strings.Contains(full, must) {
			t.Errorf("полный вывод потерял счёт: нет строки %q", must)
		}
	}
	if !strings.Contains(short, "и ещё 385 файлов") {
		t.Errorf("укороченный вывод не сказал, чего не показал:\n%s", short)
	}
	if strings.Contains(full, "и ещё") {
		t.Errorf("--top 0 обязан печатать всё и ни о чём не умалчивать")
	}
	if got := strings.Count(short, ".bin"); got != 15 {
		t.Errorf("строк перечня %d, ждали 15", got)
	}
	if got := strings.Count(full, ".bin"); got != 400 {
		t.Errorf("строк перечня при --top 0: %d, ждали 400", got)
	}
}

// TestProtectedIsItsOwnSection: защита и отказ — разные вещи, и в отчёте они
// разные разделы. Смешать их значило бы спрятать расхождение слоёв среди
// распоряжений человека.
func TestProtectedIsItsOwnSection(t *testing.T) {
	p := bigPlan(2)
	p.Protected = []clean.Protected{{Path: "/home/u/projects/a", Class: core.ClassCache, Size: 10}}
	p.ProtectedBytes = 10
	p.Refused = []clean.Refusal{{Path: "/home/u/b", Class: core.ClassCache, Verdict: core.VerdictRemovable, Reason: "это каталог"}}
	out := render(t, p, 0)

	if !strings.Contains(out, "ЗАЩИЩЕНО  1 файлов") {
		t.Errorf("нет раздела ЗАЩИЩЕНО:\n%s", out)
	}
	if !strings.Contains(out, "ОТКАЗАНО  1 записей") {
		t.Errorf("нет раздела ОТКАЗАНО:\n%s", out)
	}
	if strings.Index(out, "ЗАЩИЩЕНО") > strings.Index(out, "ОТКАЗАНО") {
		t.Error("ОТКАЗАНО обязано идти последним: это расхождение слоёв, и его читают внимательнее")
	}
}
