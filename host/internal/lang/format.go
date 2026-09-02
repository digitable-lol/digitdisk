// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Numbers and dates are written the way the language writes them, and that is
// not decoration.  «12,3 ГиБ» read as English is twelve and three, and
// «02.09.2026» read as English is the second of September or the ninth of
// February depending on who is reading — a report whose numbers change meaning
// with the reader is worse than a report in the wrong language, because the
// wrong language is obvious and a wrong number is not.
//
// Russian: comma for the decimal, non-breaking space between the thousands,
// day.month.year.  English: point, comma, ISO year-month-day — ISO rather than
// month/day/year because the tool is read on machines, not only in one
// country, and 2026-09-02 is the one spelling nobody misreads.

// nbsp is the Russian thousands separator.  It is one rune, so every %-Ns
// column in the reports keeps its width.
const nbsp = " "

// Dec renders a fractional number with prec digits after the separator.
func (l Lang) Dec(v float64, prec int) string {
	s := strconv.FormatFloat(v, 'f', prec, 64)
	if l == RU {
		s = strings.Replace(s, ".", ",", 1)
	}
	return s
}

// Pct renders a percentage, sign included.
func (l Lang) Pct(v float64, prec int) string { return l.Dec(v, prec) + "%" }

// Num renders a whole number with its thousands grouped.
func (l Lang) Num(n int64) string {
	neg := ""
	if n < 0 {
		neg, n = "-", -n
	}
	sep := nbsp
	if l == EN {
		sep = ","
	}
	digits := strconv.FormatInt(n, 10)
	var b strings.Builder
	for i, r := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteString(sep)
		}
		b.WriteRune(r)
	}
	return neg + b.String()
}

// UNum is Num for unsigned counts.
func (l Lang) UNum(n uint64) string { return l.Num(int64(n)) }

// byteUnits are the binary prefixes, and they are translated because their
// spelling is: КиБ is not a transliteration of KiB, it is what the Russian
// standard writes.
var byteUnits = map[Lang][6]string{
	RU: {"Б", "КиБ", "МиБ", "ГиБ", "ТиБ", "ПиБ"},
	EN: {"B", "KiB", "MiB", "GiB", "TiB", "PiB"},
}

// Bytes renders a byte count with a binary-prefix unit.
func (l Lang) Bytes(n int64) string {
	units := byteUnits[l]
	neg := ""
	if n < 0 {
		neg, n = "-", -n
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%s%d %s", neg, n, units[0])
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit && exp < 4; v /= unit {
		div *= unit
		exp++
	}
	return neg + l.Dec(float64(n)/float64(div), 1) + " " + units[exp+1]
}

// UBytes is Bytes for unsigned counts.
func (l Lang) UBytes(n uint64) string { return l.Bytes(int64(n)) }

// RawBytes is the exact count with its unit, for the places a report prints
// the number itself beside the rounded one.
func (l Lang) RawBytes(n int64) string { return l.Num(n) + " " + byteUnits[l][0] }

// DateTime renders an instant: day and minute, no seconds — the reports that
// use it are about when something happened, not about how long it took.
func (l Lang) DateTime(t time.Time) string {
	if l == RU {
		return t.Format("02.01.2006 15:04")
	}
	return t.Format("2006-01-02 15:04")
}

// Date renders a day.
func (l Lang) Date(t time.Time) string {
	if l == RU {
		return t.Format("02.01.2006")
	}
	return t.Format("2006-01-02")
}

// StampDate turns the ISO-8601 stamp the журнал keeps into the date this
// language writes.  A stamp that is not a stamp comes back untouched: the
// журнал is a file on disk and may hold whatever an earlier version wrote.
func (l Lang) StampDate(stamp string) string {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15-04-05Z"} {
		if t, err := time.Parse(layout, stamp); err == nil {
			return l.Date(t)
		}
	}
	if len(stamp) >= 10 {
		if t, err := time.Parse("2006-01-02", stamp[:10]); err == nil {
			return l.Date(t)
		}
	}
	return stamp
}

// Duration units are words, and short ones: they sit inside table cells whose
// width is the same in both languages.

// Days renders a whole number of days with its unit.
func (l Lang) Days(n float64) string { return l.Dec(n, 0) + " " + l.T("дн") }

// Since renders how long ago something was, the way a person says it aloud.
func (l Lang) Since(d time.Duration) string {
	switch {
	case d < time.Minute:
		return l.F("%d с", int(d.Seconds()))
	case d < time.Hour:
		return l.F("%d мин", int(d.Minutes()))
	case d < 48*time.Hour:
		return l.F("%d ч", int(d.Hours()))
	default:
		return l.F("%d дн", int(d.Hours()/24))
	}
}

// Millis renders a short duration.
func (l Lang) Millis(d time.Duration) string {
	if d < time.Second {
		return l.F("%d мс", d.Milliseconds())
	}
	return l.F("%s с", l.Dec(d.Seconds(), 1))
}

// Every renders the period of the live screen.
func (l Lang) Every(d time.Duration) string {
	if d < time.Second {
		return l.F("%d мс", d.Milliseconds())
	}
	return l.F("%s с", strings.TrimRight(strings.TrimRight(l.Dec(d.Seconds(), 3), "0"), ".,"))
}

// Uptime renders a lifetime as "5д 03:14" — and as "5d 03:14" in English.
func (l Lang) Uptime(seconds float64) string {
	if seconds <= 0 {
		return ""
	}
	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%d%s %02d:%02d", days, l.T("д"), hours, mins)
	}
	return fmt.Sprintf("%02d:%02d", hours, mins)
}
