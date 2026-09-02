// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build flangui

// Package uiflang binds the screen layout to flang-tui, the library of total
// functions printed as Go packages into digitdisk/ui-flang/out-go.
//
// It sits behind the "flangui" build tag on the same terms as coreflang sits
// behind "flangcore": the printed packages are separate Go modules generated
// from another repository, and the host's default build must not depend on
// them.  Without the tag the host uses the hand-written Go in package ui and
// nothing here is compiled.
//
// To build with the printed layout:
//
//	go build -tags flangui -o digitdisk ./host
//
// ── ЧТО СЮДА ЕДЕТ И ЧТО НЕ ЕДЕТ ──────────────────────────────────────────────
//
// Библиотека считает РАСКЛАДКУ: сколько клеток закрасить, где обрезать строку,
// какие строки раздела показать, куда уехала прокрутка, какой знак у замера,
// какой номер у цвета.  Краска остаётся у хозяина: `bar` красит два прогона
// клеток двумя цветами, `spark` красит каждый знак своим, и раскрашенная
// полоса — это уже не то, что сверялось.  Поэтому отсюда уезжают ЧИСЛА и
// ГОЛЫЕ строки, а последовательности цвета вешает на них пакет ui.
//
// ── ПОЧЕМУ КОНТЕКСТ СОЗДАЁТСЯ НА КАЖДЫЙ ВЫЗОВ ────────────────────────────────
//
// `rt.Ctx` считает шаги НАКОПИТЕЛЬНО и сам их не сбрасывает: `ctx.Steps++` на
// каждом шаге, отказ на `MaxSteps` (1 000 000).  Кадр 200×50 стоит 51 шаг, то
// есть один общий контекст на весь прогон умер бы примерно через 19 600
// перерисовок — часов через пять работы `digitdisk analyze` на открытом экране.
// Контекст — маленькая запись из пяти чисел; создать его на вызов дешевле, чем
// объяснять человеку, почему экран перестал рисоваться к вечеру.
package uiflang

import (
	"errors"
	"sync"

	fc "flangcolour/flang"
	ff "flangformat/flang"
	fh "flanghistory/flang"
	fs "flangscreen/flang"
	fsc "flangscroll/flang"
	ft "flangtabs/flang"

	rtc "flangcolour/flangrt"
	rtf "flangformat/flangrt"
	rth "flanghistory/flangrt"
	rts "flangscreen/flangrt"
	rtsc "flangscroll/flangrt"
	rtt "flangtabs/flangrt"
)

// Функции библиотеки тотальны: отказать они могут только по вшитым пределам
// вычислителя — глубине вызовов (10 000) и числу шагов, — и оба предела лежат
// далеко за размерами настоящего терминала.  Замеренные границы:
//
//	полоса        отказывает при ширине  ≥ 10 000 клеток
//	раздел        отказывает при высоте  ≥ 10 000 строк
//	обрезка       считает и 20 000 знаков в строке
//	график        считает и 20 000 замеров
//
// Отказ не роняет инструмент: вызов возвращает безопасное значение, а счётчик
// ниже помнит, сколько раз это случилось и почему.  Пустой экран лучше
// упавшего, но молчание хуже обоих — Failures() говорит об этом вслух.
var (
	mu       sync.Mutex
	failures int
	firstErr error
)

func note(err error) {
	mu.Lock()
	defer mu.Unlock()
	failures++
	if firstErr == nil {
		firstErr = err
	}
}

// Failures reports how many calls the printed layer refused, and the first
// reason.  Zero is the expected answer on any terminal a person owns.
func Failures() (int, error) {
	mu.Lock()
	defer mu.Unlock()
	return failures, firstErr
}

// ───────────────────────────── Screen ─────────────────────────────

// BarCells is how many of width cells a gauge fills for a share.  «Клеток
// полосы»: int(clamp(frac)*width + 0.5), the rounding digitdisk has always
// done.  The host paints the two runs; the split point is decided here.
func BarCells(frac float64, width int) int {
	v, err := fs.KletokPolosy(fs.NewContext(), rts.Number(frac), rts.Number(float64(width)))
	if err != nil {
		note(err)
		return 0
	}
	return int(v.Num)
}

// Clip cuts a painted line to n printing cells.  «Обрезать по ячейкам»:
// escape sequences carry no width and pass through whole.
func Clip(s string, n int, ellipsis, reset string) string {
	v, err := fs.ObrezatPoYacheykam(fs.NewContext(), rts.Text(s), rts.Number(float64(n)), rts.Text(ellipsis), rts.Text(reset))
	if err != nil {
		note(err)
		return s
	}
	return v.Str
}

// PlainWidth counts the cells a painted string occupies.  «Ширина без
// последовательностей».
func PlainWidth(s string) int {
	v, err := fs.ShirinaBezPosledovatelnostey(fs.NewContext(), rts.Text(s))
	if err != nil {
		note(err)
		return 0
	}
	return int(v.Num)
}

// Section lays a body of lines into exactly height lines, starting at scroll
// and padding the tail with empty lines.  «Уложить раздел».
func Section(body []string, height, scroll int) []string {
	v, err := fs.UlozhitRazdel(fs.NewContext(), texts(body), rts.Number(float64(height)), rts.Number(float64(scroll)))
	if err != nil {
		note(err)
		return padTo(body, height, scroll)
	}
	out := make([]string, len(v.List))
	for i, x := range v.List {
		out[i] = x.Str
	}
	return out
}

func texts(ss []string) rts.Value {
	out := make([]rts.Value, len(ss))
	for i, s := range ss {
		out[i] = rts.Text(s)
	}
	return rts.List(out)
}

// padTo is the fallback for Section: the same window, computed here, for the
// one case the printed layer refuses — a section 10 000 lines tall.
func padTo(body []string, height, scroll int) []string {
	if scroll > len(body)-1 {
		scroll = len(body) - 1
	}
	if scroll < 0 {
		scroll = 0
	}
	out := make([]string, 0, height)
	for i := 0; i < height; i++ {
		if j := scroll + i; j < len(body) {
			out = append(out, body[j])
			continue
		}
		out = append(out, "")
	}
	return out
}

// ───────────────────────────── Scroll ─────────────────────────────

// ClampScroll holds a scroll position inside the content: never above the top,
// never past the last line.  «Прокрутить» with «Мимо» — the step the reader
// did not take — is exactly that clamp.
func ClampScroll(scroll, height, length int) int {
	v, err := fsc.Prokrutit(fsc.NewContext(), rtsc.Number(float64(scroll)), rtsc.Number(float64(height)),
		rtsc.Number(float64(length)), fsc.VariantMimo())
	if err != nil {
		note(err)
		return scroll
	}
	return int(v.Num)
}

// BodyHeight is how many lines the section itself gets: the screen less the
// chrome around it.  «Высота тела», never below one.
func BodyHeight(rows, chrome int) int {
	v, err := fsc.VysotaTela(fsc.NewContext(), rtsc.Number(float64(rows)), rtsc.Number(float64(chrome)))
	if err != nil {
		note(err)
		return 1
	}
	return int(v.Num)
}

// ───────────────────────────── Tabs ─────────────────────────────

func tabAfter(open, count int, press rtt.Value) int {
	v, err := ft.Pereklyuchit(ft.NewContext(), ft.SozdatVkladki(rtt.Number(float64(open)), rtt.Number(float64(count))), press)
	if err != nil {
		note(err)
		return open
	}
	for _, f := range v.Fields {
		if f.Name == "открыта" {
			return int(f.Value.Num)
		}
	}
	note(errors.New("flangtabs.Pereklyuchit: the answer carries no open-tab field"))
	return open
}

// NextTab and PrevTab walk the sections in a ring.  «Переключить».
func NextTab(open, count int) int { return tabAfter(open, count, ft.VariantSleduyuschaya()) }

// PrevTab is NextTab the other way.
func PrevTab(open, count int) int { return tabAfter(open, count, ft.VariantPredyduschaya()) }

// PickTab answers a digit key: n is one-based, as the reader typed it.  A
// digit past the last section leaves the open one alone.
func PickTab(open, n, count int) int {
	return tabAfter(open, count, ft.VariantNomerom(rtt.Number(float64(n))))
}

// ───────────────────────────── History ─────────────────────────────

// Push appends one sample and keeps at most limit of them.  «Дописать замер».
func Push(h []float64, v float64, limit int) []float64 {
	r, err := fh.DopisatZamer(fh.NewContext(), nums(h), rth.Number(v), rth.Number(float64(limit)))
	if err != nil {
		note(err)
		return h
	}
	out := make([]float64, len(r.List))
	for i, x := range r.List {
		out[i] = x.Num
	}
	return out
}

func nums(xs []float64) rth.Value {
	out := make([]rth.Value, len(xs))
	for i, x := range xs {
		out[i] = rth.Number(x)
	}
	return rth.List(out)
}

// SparkTail is the last n samples of a history, and no more.  «Хвост списка».
func SparkTail(h []float64, n int) []float64 {
	r, err := fh.HvostSpiska(fh.NewContext(), nums(h), rth.Number(float64(n)))
	if err != nil {
		note(err)
		return h
	}
	out := make([]float64, len(r.List))
	for i, x := range r.List {
		out[i] = x.Num
	}
	return out
}

// SparkGlyph is the one character standing for a share: a gap for an
// unmeasured one, ▁ for zero, █ for full.  «Знак замера».
func SparkGlyph(v float64) string {
	r, err := fh.ZnakZamera(fh.NewContext(), rth.Number(v))
	if err != nil {
		note(err)
		return SparkGap(1)
	}
	return r.Str
}

// SparkGap is n gap characters, the lead-in before a history shorter than the
// graph.  «Пропуски».
func SparkGap(n int) string {
	r, err := fh.Propuski(fh.NewContext(), rth.Number(float64(n)))
	if err != nil {
		note(err)
		return ""
	}
	return r.Str
}

// Level reads a share the way the screen colours it: 0 calm, 1 attention,
// 2 alarm.  «Уровень доли» — the two thresholds, 0.75 and 0.90, stated once
// in the library and not repeated here.
func Level(frac float64) int {
	r, err := fh.UrovenDoli(fh.NewContext(), rth.Number(frac))
	if err != nil {
		note(err)
		return 0
	}
	switch {
	case rth.VariantIs(r, "Тревога"):
		return 2
	case rth.VariantIs(r, "Внимание"):
		return 1
	}
	return 0
}

// ───────────────────────────── Colour ─────────────────────────────

// Cube256 places an RGB colour in the xterm-256 palette.  «Куб 256».
func Cube256(r, g, b int) int {
	v, err := fc.Kub256(fc.NewContext(), rtc.Number(float64(r)), rtc.Number(float64(g)), rtc.Number(float64(b)))
	if err != nil {
		note(err)
		return 0
	}
	return int(v.Num)
}

// SeqTrue, Seq256 and Seq16 are the escape sequences for the three colour
// depths the screen knows.  «Цвет истинный», «Цвет из палитры», «Цвет из
// шестнадцати».
func SeqTrue(bg bool, r, g, b int) string {
	v, err := fc.CvetIstinnyy(fc.NewContext(), rtc.Flag(bg), rtc.Number(float64(r)), rtc.Number(float64(g)), rtc.Number(float64(b)))
	if err != nil {
		note(err)
		return ""
	}
	return v.Str
}

// Seq256 is SeqTrue through the 256-colour cube.
func Seq256(bg bool, r, g, b int) string {
	v, err := fc.CvetIzPalitry(fc.NewContext(), rtc.Flag(bg), rtc.Number(float64(r)), rtc.Number(float64(g)), rtc.Number(float64(b)))
	if err != nil {
		note(err)
		return ""
	}
	return v.Str
}

// Seq16 is the sixteen-colour sequence for an ANSI code.
func Seq16(bg bool, ansi int) string {
	v, err := fc.CvetIzShestnadcati(fc.NewContext(), rtc.Flag(bg), rtc.Number(float64(ansi)))
	if err != nil {
		note(err)
		return ""
	}
	return v.Str
}

// ───────────────────────────── Format ─────────────────────────────

// PctOf is a share of a whole, guarded against a zero whole.  «Доля от».
func PctOf(part, whole uint64) float64 {
	v, err := ff.DolyaOt(ff.NewContext(), rtf.Number(float64(part)), rtf.Number(float64(whole)))
	if err != nil {
		note(err)
		return 0
	}
	return v.Num
}
