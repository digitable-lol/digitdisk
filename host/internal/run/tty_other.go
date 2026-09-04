// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !linux

package run

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ttyState is what the строка состояния has to know about the terminal.  See
// tty_linux.go for what each field means; only the way of asking differs.
type ttyState struct {
	Rows, Cols int
	Cooked     bool
}

// ttyLook asks stty(1), the way internal/ui already asks it for a size.
//
// One child process per question, so the caller asks at the redraw rate and
// not at the замер rate.  stty acts on its own standard input, which is the
// spelling both GNU and BSD understand without a flag for the file.
func ttyLook(f *os.File) (ttyState, bool) {
	if f == nil {
		return ttyState{}, false
	}
	cmd := exec.Command("stty", "-a")
	cmd.Stdin = f
	out, err := cmd.Output()
	if err != nil {
		return ttyState{}, false
	}
	return parseStty(string(out))
}

// parseStty reads the size and the two flags out of `stty -a`.
//
// Both spellings are accepted, because both are in use: GNU writes «rows 24;
// columns 80», BSD and macOS write «24 rows; 80 columns».  A flag that is off
// is written with a leading minus by both, and that is the whole of what is
// needed here: a terminal in raw mode has neither icanon nor echo.
func parseStty(text string) (ttyState, bool) {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == ';' || r == ','
	})
	st := ttyState{Cooked: true}
	// РАЗБИРАЕТСЯ ПАРА, А НЕ СОСЕД, и это единственный способ прочитать обе
	// записи одним проходом.
	//
	// Прежний разбор брал число «рядом» — сначала справа, потом слева, — и на
	// записи BSD отвечал 80×80 вместо 24×80:
	//
	//	speed 38400 baud; 24 rows; 80 columns;
	//	      [0]   [1]  [2] [3]  [4]  [5]
	//
	// у «rows» справа стоит 80, и оно годилось. Перевернуть порядок нельзя:
	// тогда «columns» у GNU («rows 24; columns 80») взял бы слева 24. Ни одна
	// сторона по отдельности не верна — верна пара, а какая она, видно по
	// тому, число слово или нет.
	//
	// Съеденная пара пропускается целиком: иначе число из пары GNU досталось
	// бы следующему ключу как пара BSD.
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		if f == "-icanon" || f == "-echo" {
			st.Cooked = false
			continue
		}
		if i+1 >= len(fields) {
			continue
		}
		next := fields[i+1]
		// GNU: <ключ> <число>
		if key := sizeKey(f); key != "" {
			if n, err := strconv.Atoi(next); err == nil {
				st.set(key, n)
				i++
				continue
			}
		}
		// BSD: <число> <ключ>
		if key := sizeKey(next); key != "" {
			if n, err := strconv.Atoi(f); err == nil {
				st.set(key, n)
				i++
			}
		}
	}
	return st, st.Rows > 0 && st.Cols > 0
}

// sizeKey says which of the two sizes this word names, and "" when it names
// neither.
func sizeKey(word string) string {
	switch word {
	case "rows", "columns":
		return word
	}
	return ""
}

// set writes the size the key names.  Первым записанным значением, а не
// последним: у BSD слово «rows» встречается ещё раз в строке флагов, и
// перезапись стёрла бы уже найденный размер.
func (s *ttyState) set(key string, n int) {
	switch key {
	case "rows":
		if s.Rows == 0 {
			s.Rows = n
		}
	case "columns":
		if s.Cols == 0 {
			s.Cols = n
		}
	}
}
