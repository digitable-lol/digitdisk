// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !flangui

// Раскладка экрана, посчитанная рукописным Go. Это УМОЛЧАНИЕ: без признака
// сборки `flangui` работает всё, что здесь.
//
// Второй файл пары — layout_flang.go — отдаёт ровно эти же решения печатанной
// библиотеке flang-tui (ui-flang/). Функции здесь и там ОДНИ И ТЕ ЖЕ по имени
// и по подписи, и вызываются они из одних и тех же мест: в widgets.go,
// screen.go и theme.go стоит по одной строке-переходнику, а не по две ветки.
//
// Тела ниже — не новый код: они переехали сюда из widgets.go, screen.go и
// theme.go без единой правки, чтобы вариант по умолчанию остался буква в букву
// тем же, каким был.

package ui

import (
	"fmt"
	"strings"
)

// barCells is how many of n cells a gauge fills for a share.
func barCells(frac float64, n int) int {
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	return int(frac*float64(n) + 0.5)
}

// clipTo cuts a painted line to n printing cells.
func clipTo(s string, n int, ellipsis, reset string) string {
	if plainWidth(s) <= n {
		return s
	}
	rs := []rune(s)
	var b strings.Builder
	w, keep := 0, n-1
	for i := 0; i < len(rs); {
		if rs[i] != 0x1b {
			if w >= keep {
				break
			}
			b.WriteRune(rs[i])
			w++
			i++
			continue
		}
		j := skipEscape(rs, i)
		b.WriteString(string(rs[i:j]))
		i = j
	}
	return b.String() + ellipsis + reset
}

// plainWidth counts the cells a painted string will occupy, which is its runes
// less everything inside an escape sequence.
func plainWidth(s string) int {
	rs := []rune(s)
	w := 0
	for i := 0; i < len(rs); {
		if rs[i] == 0x1b {
			i = skipEscape(rs, i)
			continue
		}
		w++
		i++
	}
	return w
}

// skipEscape returns the index just past the escape sequence starting at i.
func skipEscape(rs []rune, i int) int {
	j := i + 1
	if j < len(rs) && (rs[j] == '[' || rs[j] == ']' || rs[j] == '?') {
		j++
		for j < len(rs) && !isFinal(rs[j]) {
			j++
		}
		if j < len(rs) {
			j++
		}
	} else if j < len(rs) {
		j++
	}
	return j
}

// isFinal reports whether r ends a CSI sequence.
func isFinal(r rune) bool {
	return r >= 0x40 && r <= 0x7E && r != '[' && r != ';' && !(r >= '0' && r <= '9')
}

// layoutSection lays body into exactly h lines starting at scroll, padding the
// tail with empty lines so the rule and the footer never ride up.
func layoutSection(body []string, h, scroll int) []string {
	shown := body[clampScroll(body, scroll):]
	if len(shown) > h {
		shown = shown[:h]
	}
	out := make([]string, 0, h)
	out = append(out, shown...)
	for i := len(shown); i < h; i++ {
		out = append(out, "")
	}
	return out
}

// clampScroll holds a scroll position inside the content.
func clampScroll(body []string, scroll int) int {
	if scroll > len(body)-1 {
		scroll = len(body) - 1
	}
	if scroll < 0 {
		scroll = 0
	}
	return scroll
}

// bodyHeightOf is the screen less the two-line head, the rule and the footer.
func bodyHeightOf(rows, chrome int) int {
	h := rows - chrome
	if h < 1 {
		h = 1
	}
	return h
}

// nextTab, prevTab and pickTab move between the sections.
func nextTab(open, count int) int { return (open + 1) % count }

func prevTab(open, count int) int { return (open + count - 1) % count }

func pickTab(open, n, count int) int {
	if n >= 1 && n-1 < count {
		return n - 1
	}
	return open
}

// pushSample appends one sample and keeps at most limit of them.
func pushSample(h []float64, v float64, limit int) []float64 {
	h = append(h, v)
	if len(h) > limit {
		h = h[len(h)-limit:]
	}
	return h
}

// sparkTail is the last n samples of a history, and no more.
func sparkTail(h []float64, n int) []float64 {
	if len(h) > n {
		return h[len(h)-n:]
	}
	return h
}

// sparkGlyph is the one character standing for a share.
func sparkGlyph(v float64) string {
	if v < 0 {
		return string(sparkGapFace)
	}
	if v > 1 {
		v = 1
	}
	return string(sparkFaces[int(v*float64(len(sparkFaces)-1))])
}

// sparkGap is n gap characters.
func sparkGap(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(string(sparkGapFace), n)
}

// levelOf reads a share the way the screen colours it: 0 calm, 1 attention,
// 2 alarm.
func levelOf(frac float64) int {
	switch {
	case frac >= 0.90:
		return 2
	case frac >= 0.75:
		return 1
	}
	return 0
}

// cube256 places an RGB colour in the xterm-256 palette: the 6×6×6 cube, or
// the grey ramp when the three channels sit close together.
func cube256(c RGB) int {
	r, g, b := int(c.R), int(c.G), int(c.B)
	if abs(r-g) < 12 && abs(g-b) < 12 && abs(r-b) < 12 {
		grey := (r + g + b) / 3
		switch {
		case grey < 8:
			return 16
		case grey > 248:
			return 231
		}
		return 232 + (grey-8)*23/240
	}
	q := func(v int) int { return (v*5 + 127) / 255 }
	return 16 + 36*q(r) + 6*q(g) + q(b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// seqTrue, seq256 and seq16 are the escape sequences for the three colour
// depths the screen knows.
func seqTrue(bg bool, c RGB) string {
	lead := 38
	if bg {
		lead = 48
	}
	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", lead, c.R, c.G, c.B)
}

func seq256(bg bool, c RGB) string {
	lead := 38
	if bg {
		lead = 48
	}
	return fmt.Sprintf("\x1b[%d;5;%dm", lead, cube256(c))
}

func seq16(bg bool, ansi int) string {
	if bg {
		ansi += 10 // 30..37 -> 40..47, 90..97 -> 100..107
	}
	return fmt.Sprintf("\x1b[%dm", ansi)
}

// pctOf is a share of a whole, guarded against a zero whole.
func pctOf(part, whole uint64) float64 {
	if whole == 0 {
		return 0
	}
	return float64(part) / float64(whole)
}

// layoutSource names, for the reader and for `digitdisk version`, who computed
// the layout of this build.
const layoutSource = "Go"

// layoutFailures is the pair of the same name in layout_flang.go.  Рукописному
// Go отказывать не по чему: пределов вычислителя у него нет.
func layoutFailures() (int, error) { return 0, nil }
