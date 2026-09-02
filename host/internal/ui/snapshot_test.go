// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/cli"
	"digitdisk/internal/lang"
)

// Снимок экрана на неподвижном замере: тем же кадром, что рисует draw(), но
// без терминала и без часов.
//
// ЗАЧЕМ. Ветка `flangui` отдаёт раскладку печатанной библиотеке flang-tui, и
// уговор ветки — вывод обязан остаться прежним. Проверить это на живом экране
// нельзя: там часы, доли ЦП и «замер N назад» меняются между запусками, и
// расхождение реализаций утонуло бы в расхождении времени. Поэтому кадр
// снимается с НЕПОДВИЖНОГО снимка системы (filled()), а два места, которые
// всё-таки читают часы, — время в шапке и «замер N назад» в подвале —
// заменяются на постоянные строки прямо здесь, а не маской снаружи.
//
// ЭТО ДАМПЕР, А НЕ ЭТАЛОННЫЙ ФАЙЛ. Эталон в дереве протух бы на первом же
// новом разделе; сверка идёт между ДВУМЯ СБОРКАМИ одного и того же дерева:
//
//	tools/sverka-ui.sh
//
// Без переменной DIGITDISK_SNAPSHOT проверка не делает ничего.
func TestSnapshotDump(t *testing.T) {
	path := os.Getenv("DIGITDISK_SNAPSHOT")
	if path == "" {
		t.Skip("снимок не заказан: DIGITDISK_SNAPSHOT не задан")
	}

	var b strings.Builder
	depths := []struct {
		name string
		d    depth
	}{{"истинный", depthTrue}, {"256", depth256}, {"16", depth16}, {"без цвета", depthNone}}

	for _, dep := range depths {
		for _, cols := range []int{40, 80, 120, 200} {
			for _, rows := range []int{24, 50} {
				for _, l := range []lang.Lang{lang.RU, lang.EN} {
					for tab := 0; tab < len(sections); tab++ {
						for _, scroll := range []int{0, 1, 7, 999} {
							// Курсор списка команд — часть кадра: он
							// красит строку и двигает прокрутку.
							for _, pick := range []int{0, 3, len(cli.Commands) - 1} {
								s := snapScreen(dep.d, rows, cols, l)
								s.tab, s.scroll, s.pick = tab, scroll, pick
								fmt.Fprintf(&b, "== цвет=%s кол=%d строк=%d язык=%v вкладка=%d прокрутка=%d строка=%d\n",
									dep.name, cols, rows, l, tab, scroll, pick)
								// draw() обрезает каждую строку по ширине окна —
								// это часть раскладки, и она входит в снимок.
								for i, line := range s.frame() {
									fmt.Fprintf(&b, "%3d %q\n", i, stopClock(s.t.clip(line, s.cols)))
								}
								fmt.Fprintf(&b, "прокрутка после кадра %d\n", s.scroll)
							}
						}
					}
				}
			}
		}
	}

	// Ф1: данных, которых на сетке не было. Очень длинные имена, эмодзи в
	// пути, пустые разделы, голая управляющая последовательность.
	for _, cols := range []int{40, 80, 120, 200} {
		s := snapScreen(depthTrue, 24, cols, lang.RU)
		fmt.Fprintf(&b, "== край кол=%d\n", cols)
		for i, line := range awkwardLines() {
			fmt.Fprintf(&b, "%3d %q\n", i, s.t.clip(line, s.cols))
			fmt.Fprintf(&b, "%3d ширина %d\n", i, plainWidth(line))
		}
		for _, h := range [][]float64{nil, {}, {-1}, {0, 0.751, 0.9, 1.0001, -1}} {
			fmt.Fprintf(&b, "график %q\n", s.t.spark(h, cols/4))
			fmt.Fprintf(&b, "полоса %q\n", s.t.bar(0.7314, cols/4))
		}
		for _, body := range [][]string{nil, {}, {""}, {"одна"}} {
			fmt.Fprintf(&b, "раздел %q\n", layoutSection(body, 5, 0))
		}
	}

	// Снимки совпали бы и тогда, когда печатанная раскладка молча отказала и
	// откатилась на тот же Go. Проверка держит утверждение честным: за весь
	// обход НИ ОДНОГО отката быть не должно.
	if n, err := layoutFailures(); n != 0 {
		t.Fatalf("раскладка отказывала %d раз, первый раз: %v", n, err)
	}

	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("снимок не записан: %v", err)
	}
	t.Logf("снимок записан: %s, байт %d", path, len(b.String()))
}

// snapScreen builds a screen whose every reading is fixed, including the two
// that would otherwise be read from the clock.
func snapScreen(d depth, rows, cols int, l lang.Lang) *screen {
	s := &screen{t: Theme{P: Carbon, d: d}, rows: rows, cols: cols, l: l}
	s.st, s.haveSt = filled(), true
	// taken лежит РОВНО на «сейчас», а took постоянен: подвал тогда печатает
	// одну и ту же строку в любом запуске.
	s.taken, s.took = time.Now(), 1900*time.Millisecond
	s.o.Interval = 2 * time.Second
	s.cpuHist = []float64{0.1, 0.5, -1, 0.9, 0.95, 0.76, 0, 1}
	s.memHist = []float64{0.3, 0.4}
	return s
}

// awkwardLines is the data the grid never had: длинные имена, эмодзи в пути,
// склеенные последовательности, пустая строка.
func awkwardLines() []string {
	return []string{
		"",
		"\x1b[0m",
		"\x1b[38;2;1;2;3m\x1b[48;5;17m\x1b[1m",
		strings.Repeat("оченьдлинноеимякаталога", 40),
		"/home/пользователь/" + strings.Repeat("вложенный/", 30) + "файл.txt",
		"путь/к/🙂/файлу",
		"👩‍💻 👨‍👩‍👧‍👦 🇷🇺 é ё",
		"\x1b[31m" + strings.Repeat("🙂", 100) + "\x1b[0m",
		"\x1b]8;;http://пример/очень/длинная/ссылка\x07подпись\x1b]8;;\x07",
		strings.Repeat("\x1b[32mз\x1b[0m", 200),
		"обрыв последовательности \x1b[38;2;9;9",
		"\x1b",
	}
}

// clockFace is the one thing on the screen that is neither a reading nor a
// layout: the wall clock in the head line.  Two builds cannot be started at
// the same second, so it is replaced by a constant before the frames are
// compared.  Nothing else is touched — every other number on the screen comes
// from the fixed snapshot and must match byte for byte.
var clockFace = regexp.MustCompile(`\d\d:\d\d:\d\d`)

func stopClock(s string) string { return clockFace.ReplaceAllString(s, "ЧЧ:ММ:СС") }
