// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build flangcore

package coreflang

import (
	"testing"

	"digitdisk/internal/core"
)

// Второй вопрос доезжает до НАСТОЯЩЕГО ядра и возвращается словом, а не
// догадкой хозяина.
//
// Правила проверяются там, где живут — 183 примера в
// core/disk-inventory.flang, — а здесь проверяется мост: что «Природа
// находки» вызвана, что имя варианта опознано, и что строгость приезжает из
// «Строгости», а не из копии шкалы на этой стороне.
func TestTheSecondQuestionReachesTheRealCore(t *testing.T) {
	b := New()
	keep := core.Decision{Class: core.ClassUnknown, Verdict: core.VerdictKeep}
	rubbish := core.Decision{Class: core.ClassCache, Verdict: core.VerdictRemovable}

	for _, c := range []struct {
		name string
		rec  core.Record
		d    core.Decision
		want core.Nature
		step int
	}{
		{"исходник владельца", core.Record{Path: "/home/u/проект/сервер.go", Kind: core.KindFile, Accessible: true}, keep, core.NatureSource, 2},
		{"записка рядом", core.Record{Path: "/home/u/проект/README.md", Kind: core.KindFile, Accessible: true}, keep, core.NatureSource, 2},
		{"старый кэш", core.Record{Path: "/home/u/.cache/pip/a.whl", Kind: core.KindFile, AgeDays: 30, Accessible: true}, rubbish, core.NatureTrash, 1},
		{"объект git", core.Record{Path: "/home/u/п/.git/objects/ab/cdef", Kind: core.KindFile, Accessible: true}, keep, core.NatureVCS, 3},
		{"слой образа", core.Record{Path: "/home/u/.local/share/containers/storage/overlay/c0359c5b1a2b3c4d5e6f708192a3b4c5d6e7f8091a2b3c4d5e6f708192a3fb3e/diff/x", Kind: core.KindFile, Accessible: true}, keep, core.NatureStore, 3},
		{"фильм", core.Record{Path: "/home/u/видео/фильм.mkv", Kind: core.KindFile, Accessible: true}, keep, core.NaturePersonal, 2},
		{"молодой кэш", core.Record{Path: "/home/u/.cache/pip/b.whl", Kind: core.KindFile, AgeDays: 1, Accessible: true},
			core.Decision{Class: core.ClassCache, Verdict: core.VerdictAsk}, core.NatureFresh, 2},
	} {
		got := b.Nature(c.rec, c.d)
		if got != c.want {
			t.Errorf("%s: ядро назвало %q, ожидалось %q", c.name, got, c.want)
			continue
		}
		if step := b.Strictness(got); step != c.step {
			t.Errorf("%s: строгость %d, ожидалась %d", c.name, step, c.step)
		}
	}
	if n, err := b.Failures(); n > 0 {
		t.Errorf("слой отказал %d раз, первый — %v", n, err)
	}
}

// Мусором ядро зовёт ТОЛЬКО то, что само же убирает. Это и есть постусловие
// «Природа обоснована», и мост обязан его унести целым: иначе забой мог бы
// спросить одной клавишей про исходники.
func TestOnlyWhatTheCoreRemovesIsCalledRubbish(t *testing.T) {
	b := New()
	for _, v := range []core.Verdict{core.VerdictAsk, core.VerdictKeep} {
		got := b.Nature(
			core.Record{Path: "/home/u/.cache/pip/a.whl", Kind: core.KindFile, AgeDays: 1, Accessible: true},
			core.Decision{Class: core.ClassCache, Verdict: v})
		if got == core.NatureTrash {
			t.Errorf("приговор %q, а природа «%s» — мусором зовётся неубираемое", v, got)
		}
	}
}
