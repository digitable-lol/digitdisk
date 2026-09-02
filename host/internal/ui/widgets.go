// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"strings"
)

// A row is a line under construction.  Colour is added as each piece is
// appended, while the width is counted from the piece's plain text, so escape
// sequences never enter the arithmetic and columns stay aligned.
type row struct {
	b strings.Builder
	w int
}

// add appends text painted by paint, counting the width of the plain text.
func (r *row) add(plain string, paint func(string) string) {
	r.b.WriteString(paint(plain))
	r.w += runes(plain)
}

// plain appends unpainted text.
func (r *row) plain(s string) { r.add(s, func(x string) string { return x }) }

// pad grows the row with spaces until it is at least n cells wide.
func (r *row) pad(n int) {
	if n > r.w {
		r.b.WriteString(strings.Repeat(" ", n-r.w))
		r.w = n
	}
}

func (r *row) String() string { return r.b.String() }

func runes(s string) int { return len([]rune(s)) }

// fit pads s with spaces to n cells, or cuts it with an ellipsis.  Widths are
// counted in runes: the screen speaks Russian, and a byte count would cut a
// letter in half.
func fit(s string, n int) string {
	if n <= 0 {
		return ""
	}
	rs := []rune(s)
	switch {
	case len(rs) == n:
		return s
	case len(rs) < n:
		return s + strings.Repeat(" ", n-len(rs))
	case n == 1:
		return "…"
	}
	return string(rs[:n-1]) + "…"
}

// right pads s on the left so it ends at column n.
func right(s string, n int) string {
	if d := n - runes(s); d > 0 {
		return strings.Repeat(" ", d) + s
	}
	return fit(s, n)
}

// level is how a share of something full is to be read.  The screen shows the
// same numbers the report prints; the colour is a reading aid, not a new
// measurement, and the two thresholds are stated here once.
func (t Theme) level(frac float64) slot {
	switch {
	case frac >= 0.90:
		return t.P.Red
	case frac >= 0.75:
		return t.P.Yellow
	default:
		return t.P.Green
	}
}

// bar draws a gauge n cells wide for a share between 0 and 1.  An unmeasured
// share — a nil pointer upstream — is not a zero and must not reach here; the
// callers draw dashes for that instead.
func (t Theme) bar(frac float64, n int) string {
	if n <= 0 {
		return ""
	}
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	full := int(frac*float64(n) + 0.5)
	var r row
	r.add(strings.Repeat("█", full), func(s string) string { return t.Fg(t.level(frac), s) })
	r.add(strings.Repeat("─", n-full), func(s string) string { return t.Fg(t.P.Border, s) })
	return r.String()
}

// emptyBar draws the same gauge for a share that was never measured.
func (t Theme) emptyBar(n int) string {
	if n <= 0 {
		return ""
	}
	return t.Fg(t.P.Border, strings.Repeat("┈", n))
}

var sparkFaces = []rune("▁▂▃▄▅▆▇█")

// spark draws the recent history of one share as a single line.  The history
// is nothing but the samples already shown at the top of the screen, kept as
// they scroll by — no source is read for it.
func (t Theme) spark(history []float64, n int) string {
	if n <= 0 {
		return ""
	}
	vals := history
	if len(vals) > n {
		vals = vals[len(vals)-n:]
	}
	var r row
	if lead := n - len(vals); lead > 0 {
		r.add(strings.Repeat("·", lead), func(s string) string { return t.Fg(t.P.Border, s) })
	}
	for _, v := range vals {
		if v < 0 {
			r.add("·", func(s string) string { return t.Fg(t.P.Border, s) })
			continue
		}
		if v > 1 {
			v = 1
		}
		i := int(v * float64(len(sparkFaces)-1))
		r.add(string(sparkFaces[i]), func(s string) string { return t.Fg(t.level(v), s) })
	}
	return r.String()
}

// gauge is the whole measured line: a name, a bar, and the reading itself.
func (t Theme) gauge(name string, nameWidth int, frac float64, reading string, barWidth int) string {
	var r row
	r.add("  "+fit(name, nameWidth)+" ", func(s string) string { return t.Fg(t.P.Muted, s) })
	r.plain(t.bar(frac, barWidth))
	r.w += barWidth
	r.add("  "+reading, func(s string) string { return t.Fg(t.P.Foreground, s) })
	return r.String()
}

// gaugeUnmeasured is gauge for a reading the system never published.
func (t Theme) gaugeUnmeasured(name string, nameWidth int, note string, barWidth int) string {
	var r row
	r.add("  "+fit(name, nameWidth)+" ", func(s string) string { return t.Fg(t.P.Muted, s) })
	r.plain(t.emptyBar(barWidth))
	r.w += barWidth
	r.add("  "+note, func(s string) string { return t.Fg(t.P.Subtle, s) })
	return r.String()
}

// pctOf is a share of a whole, guarded against a zero whole.
func pctOf(part, whole uint64) float64 {
	if whole == 0 {
		return 0
	}
	return float64(part) / float64(whole)
}

// clip cuts a painted line to n printing cells.  Escape sequences carry no
// width and are copied through whole, so a line is shortened without its
// colours coming apart.  A line that overruns the terminal would otherwise
// wrap, push every line below it down, and take the layout with it — so this
// runs over every line the screen draws, not only the ones suspected of it.
func (t Theme) clip(s string, n int) string {
	if n <= 0 {
		return ""
	}
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
	return b.String() + "…" + t.reset()
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
