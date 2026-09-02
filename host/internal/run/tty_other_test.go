// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !linux

package run

import "testing"

// Обе записи stty(1) в ходу, и обе надо понимать: GNU пишет «rows 24; columns
// 80», BSD и macOS — «24 rows; 80 columns». Ошибиться здесь значит не найти
// размер терминала и не нарисовать строку состояния вовсе.
func TestРазборSttyВОбеихЗаписях(t *testing.T) {
	const gnu = "speed 38400 baud; rows 24; columns 80; line = 0;\n" +
		"intr = ^C; -brkint -imaxbel iutf8\nicanon echo echoe\n"
	const bsd = "speed 38400 baud; 24 rows; 80 columns;\n" +
		"lflags: icanon isig iexten echo echoe\n"
	for name, text := range map[string]string{"GNU": gnu, "BSD": bsd} {
		st, ok := parseStty(text)
		if !ok {
			t.Errorf("%s: размер не разобран", name)
			continue
		}
		if st.Rows != 24 || st.Cols != 80 {
			t.Errorf("%s: %d×%d, ожидалось 24×80", name, st.Rows, st.Cols)
		}
		if !st.Cooked {
			t.Errorf("%s: терминал сочтён сырым, а он обычный", name)
		}
	}
	raw, ok := parseStty("speed 38400 baud; rows 24; columns 80;\n-icanon -echo\n")
	if !ok || raw.Cooked {
		t.Error("сырой режим не узнан — строка состояния подралась бы с полноэкранной программой")
	}
}
