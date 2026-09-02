// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"bufio"
	"io"
	"testing"

	"digitdisk/internal/lang"
)

// Цена раскладки: одна и та же проверка собирается дважды — без признака
// `flangui` её считает рукописный Go, с признаком — печатанная библиотека
// flang-tui, — и разница между прогонами и есть цена печатанного варианта.
//
//	go test -run '^$' -bench Layout -benchmem ./internal/ui/
//	go test -run '^$' -bench Layout -benchmem -tags flangui ./internal/ui/
//
// Меряется ровно то, что человек ждёт: КАДР — то, что делает draw() на каждой
// перерисовке, и НАЖАТИЕ — ответ на клавишу вместе с перерисовкой, которая за
// ним следует. Второе важнее первого: кадр по таймеру человек не ждёт, а
// стрелку — ждёт.

func benchScreen(rows, cols int) *screen {
	s := &screen{t: Theme{P: Carbon, d: depth256}, rows: rows, cols: cols, l: lang.RU}
	s.st, s.haveSt = filled(), true
	s.o.Interval = 2000000000 // 2s, без импорта time
	s.cpuHist = []float64{0.1, 0.5, -1, 0.9, 0.95, 0.76, 0, 1}
	s.memHist = []float64{0.3, 0.4}
	s.out = bufio.NewWriterSize(io.Discard, 1<<16)
	return s
}

// benchFrame — то, что делает draw(): уложить кадр и обрезать каждую строку
// по ширине окна.
func benchFrame(s *screen) {
	for _, line := range s.frame() {
		_ = s.t.clip(line, s.cols)
	}
}

func BenchmarkLayoutFrame80x24(b *testing.B) {
	s := benchScreen(24, 80)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchFrame(s)
	}
}

func BenchmarkLayoutFrame200x50(b *testing.B) {
	s := benchScreen(50, 200)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchFrame(s)
	}
}

// BenchmarkLayoutKey — задержка на нажатие: клавиша плюс перерисовка, которая
// за ней следует. Стрелка вправо переключает раздел, то есть заставляет
// пересобрать тело целиком — это самый дорогой из обычных ответов.
func BenchmarkLayoutKey80x24(b *testing.B) {
	s := benchScreen(24, 80)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.handle(key{kind: keyRight}, nil)
		s.draw()
	}
}

func BenchmarkLayoutKey200x50(b *testing.B) {
	s := benchScreen(50, 200)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.handle(key{kind: keyRight}, nil)
		s.draw()
	}
}

// BenchmarkLayoutScrollKey — прокрутка на одну строку: тело раздела то же,
// меняется только окно. Так проверяется, что дорога именно раскладка, а не
// сборка тела.
func BenchmarkLayoutScrollKey200x50(b *testing.B) {
	s := benchScreen(50, 200)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.handle(key{kind: keyDown}, nil)
		s.draw()
	}
}
