// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"fmt"
	"strings"
	"time"

	"digitdisk/internal/lang"
)

// Что показывается живьём и что говорится в конце.
//
// The two are different on purpose.  The строка состояния is a glance: it
// shows what is true right now and it has one line to do it in, so it drops
// what does not fit rather than wrapping.  The сводка is the answer a person
// came for — how long, how much CPU, what the peak was — and it says which
// numbers are exact and which were sampled, because «6,1 ГиБ» and «около 6,1
// ГиБ» are different claims.

// barText builds the live line, as wide as the terminal and no wider.
//
// Nothing is abbreviated to fit a narrow terminal beyond a point: it is
// better to show three true fields than six unreadable ones, so the line is
// tried whole, then short, then shortened by dropping from the right.
func barText(l lang.Lang, width int, d time.Duration, r reading, g gpuReading) string {
	if width <= 0 {
		width = 80
	}
	long := segments(l, d, r, g, true)
	short := segments(l, d, r, g, false)
	if s := strings.Join(long, "   "); width >= runes(s) {
		return s
	}
	for n := len(short); n > 0; n-- {
		if s := strings.Join(short[:n], "   "); width >= runes(s) {
			return s
		}
	}
	return clip(strings.Join(short, "   "), width)
}

// segments is the line in pieces, most important first.
func segments(l lang.Lang, d time.Duration, r reading, g gpuReading, full bool) []string {
	out := []string{clock(d)}
	if r.Known {
		cpu := l.F("ЦП %s", l.Pct(r.Percent, 0))
		if full && d > 2*time.Second {
			cpu = l.F("ЦП %s, средн. %s", l.Pct(r.Percent, 0), l.Pct(average(r, d), 0))
		}
		out = append(out, cpu)
		if r.Bytes > 0 {
			mem := l.F("память %s", l.UBytes(r.Bytes))
			if full && r.Peak > r.Bytes {
				mem = l.F("память %s, пик %s", l.UBytes(r.Bytes), l.UBytes(r.Peak))
			}
			out = append(out, mem)
		}
		if r.Processes > 0 {
			out = append(out, l.F("процессов %d", r.Processes))
		}
	}
	if g.Known && g.Bytes > 0 {
		out = append(out, l.F("видеопамять %s", l.UBytes(g.Bytes)))
	}
	return out
}

func average(r reading, d time.Duration) float64 {
	if d <= 0 {
		return 0
	}
	return r.CPUSeconds / d.Seconds() * 100
}

// Lines is the сводка: two lines, three when the driver's program was asked
// about a video card.
//
// It goes to standard ERROR, always.  Standard output belongs to the command
// — a сводка printed there would end up inside the log the command was piped
// into, and `digitdisk run make | tee log` would stop being `make | tee log`.
func Lines(l lang.Lang, r Result) []string {
	var out []string
	if r.Signal != "" {
		out = append(out, l.F("команда «%s»: убита сигналом %s, %s", r.Command, r.Signal, span(l, seconds(r.Seconds))))
	} else {
		out = append(out, l.F("команда «%s»: код %d, %s", r.Command, r.Code, span(l, seconds(r.Seconds))))
	}

	cpu := span(l, seconds(r.CPUSeconds))
	peak := l.UBytes(r.PeakBytes)
	switch {
	case r.PeakBytes > 0 && !r.PeakExact:
		peak = l.F("около %s", l.UBytes(r.PeakBytes))
	case r.PeakBytes == 0 && r.PeakOneBytes > 0:
		peak = l.F("%s (самый крупный процесс)", l.UBytes(r.PeakOneBytes))
	case r.PeakBytes == 0:
		// Nothing was seen and nothing is invented: a command shorter
		// than one замер leaves no memory reading behind, and its
		// ru_maxrss is our own footprint rather than its.
		peak = l.T("не измерен — команда короче замера")
	}
	how := accounting(l, r)
	if r.Processes > 0 {
		out = append(out, l.F("процессорное время %s (в среднем %s), пик памяти %s, процессов %d; учёт — %s",
			cpu, l.Pct(r.CPUPercent, 0), peak, r.Processes, how))
	} else {
		out = append(out, l.F("процессорное время %s (в среднем %s), пик памяти %s; учёт — %s",
			cpu, l.Pct(r.CPUPercent, 0), peak, how))
	}

	if r.GPU != nil {
		out = append(out, l.F("видеокарта: пик памяти процессов %s (%s); загрузку по процессам драйвер не публикует",
			l.UBytes(r.GPU.PeakBytes), r.GPU.Source))
	}
	return out
}

// accounting is the one phrase that says how much the numbers beside it are
// worth.
func accounting(l lang.Lang, r Result) string {
	switch r.Accounting {
	case ByGroup:
		return l.T("своя контрольная группа, считает ядро")
	case ByProc:
		return l.F("обход /proc раз в %d мс — пик памяти приблизителен", r.SampleMS)
	default:
		return l.T("итог ядра о команде и о том, чего она дождалась")
	}
}

// span is a length of time in words: seconds while there are only seconds,
// then minutes, then hours.
func span(l lang.Lang, d time.Duration) string {
	switch {
	case d < 10*time.Millisecond:
		// Below ten milliseconds whole milliseconds are all zeroes, and
		// «0 мс» about a command that did run reads as a failure to
		// measure rather than as a small number.
		return l.F("%s мс", l.Dec(float64(d.Microseconds())/1000, 1))
	case d < time.Minute:
		return l.Millis(d)
	case d < time.Hour:
		return l.F("%d мин %d с", int(d.Minutes()), int(d.Seconds())%60)
	default:
		return l.F("%d ч %d мин", int(d.Hours()), int(d.Minutes())%60)
	}
}

// clock is a length of time in digits, for the bar, where every column counts
// and nobody needs a unit to read 1:23.
func clock(d time.Duration) string {
	s := int(d.Seconds())
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func seconds(f float64) time.Duration { return time.Duration(f * float64(time.Second)) }

func runes(s string) int { return len([]rune(s)) }
