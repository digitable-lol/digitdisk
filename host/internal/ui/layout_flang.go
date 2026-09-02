// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build flangui

// Раскладка экрана, посчитанная flang. Пара к layout_stock.go: те же имена,
// те же подписи, те же места вызова, — но за каждым стоит функция библиотеки
// flang-tui, доказанная на завершение и сверенная с рукописным Go.
//
// Что сюда НЕ уехало и почему: краска. Библиотека возвращает числа и голые
// строки — сколько клеток закрасить, какой знак у замера, какой номер у цвета,
// — а последовательности цвета навешивает хозяин, ровно как раньше. Полоса,
// раскрашенная по клетке вместо двух прогонов, была бы уже другой строкой, а
// уговор этой ветки — вывод обязан остаться прежним.
//
// Собирается: go build -tags flangui -o digitdisk ./host

package ui

import "digitdisk/internal/uiflang"

func barCells(frac float64, n int) int { return uiflang.BarCells(frac, n) }

func clipTo(s string, n int, ellipsis, reset string) string {
	return uiflang.Clip(s, n, ellipsis, reset)
}

func plainWidth(s string) int { return uiflang.PlainWidth(s) }

func layoutSection(body []string, h, scroll int) []string {
	return uiflang.Section(body, h, scroll)
}

func clampScroll(body []string, scroll int) int {
	return uiflang.ClampScroll(scroll, 0, len(body))
}

func bodyHeightOf(rows, chrome int) int { return uiflang.BodyHeight(rows, chrome) }

func nextTab(open, count int) int { return uiflang.NextTab(open, count) }

func prevTab(open, count int) int { return uiflang.PrevTab(open, count) }

func pickTab(open, n, count int) int { return uiflang.PickTab(open, n, count) }

func pushSample(h []float64, v float64, limit int) []float64 {
	return uiflang.Push(h, v, limit)
}

func sparkTail(h []float64, n int) []float64 { return uiflang.SparkTail(h, n) }

func sparkGlyph(v float64) string { return uiflang.SparkGlyph(v) }

func sparkGap(n int) string { return uiflang.SparkGap(n) }

func levelOf(frac float64) int { return uiflang.Level(frac) }

func cube256(c RGB) int { return uiflang.Cube256(int(c.R), int(c.G), int(c.B)) }

func seqTrue(bg bool, c RGB) string {
	return uiflang.SeqTrue(bg, int(c.R), int(c.G), int(c.B))
}

func seq256(bg bool, c RGB) string {
	return uiflang.Seq256(bg, int(c.R), int(c.G), int(c.B))
}

func seq16(bg bool, ansi int) string { return uiflang.Seq16(bg, ansi) }

func pctOf(part, whole uint64) float64 { return uiflang.PctOf(part, whole) }

// layoutFailures reports how many calls the printed layer refused, and the
// first reason. Zero is the expected answer on any terminal a person owns; a
// non-zero answer means a вшитый предел вычислителя was hit and the screen
// fell back on the Go it was supposed to replace — which would make «вывод
// совпал» say nothing at all. TestSnapshotDump therefore checks it.
func layoutFailures() (int, error) { return uiflang.Failures() }

// layoutSource names, for the reader and for `digitdisk version`, who computed
// the layout of this build.
const layoutSource = "flang"
