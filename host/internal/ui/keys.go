// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"bufio"
	"os"
	"sync"
	"unicode/utf8"
)

// ОДНА КЛАВИАТУРА НА ВЕСЬ ПРОЦЕСС.
//
// Экранов теперь два и они сменяют друг друга: с экрана состояния
// запускается analyze, из него — возврат обратно. Пока каждый экран открывал
// свой /dev/tty и заводил своё чтение, смена стоила читателю НАЖАТИЙ, и не
// одного: чтение предыдущего экрана остаётся висеть в read(2) на том же
// терминале — закрытие файла его не будит, — и первым получает всё, что
// набрано дальше. Экран, которому это предназначалось, не получает ничего и
// выглядит зависшим.
//
// Поэтому терминал открывается один раз и читается одной горутиной, а экраны
// берут её канал по очереди. Сырой режим по-прежнему у каждого свой: его
// ставят и снимают вокруг отрисовки, и это к чтению отношения не имеет.
//
// Набранное между экранами не пропадает, а лежит в канале: читатель, набравший
// вперёд, получает своё нажатие на следующем экране, а не в никуда.
var (
	kbdOnce sync.Once
	kbdTTY  *os.File
	kbdKeys chan key
	kbdErr  error
)

// keyboard hands out the process's one terminal and the keys read from it.
func keyboard() (*os.File, <-chan key, error) {
	kbdOnce.Do(func() {
		f, err := openInput()
		if err != nil {
			kbdErr = err
			return
		}
		kbdTTY, kbdKeys = f, make(chan key, 64)
		go readKeys(f, kbdKeys)
	})
	return kbdTTY, kbdKeys, kbdErr
}

// WaitKey holds printed output on the terminal until the reader has read it,
// and answers on the first key they press.
//
// It reads through the same one keyboard the screens read through, and that is
// the whole reason it lives here rather than in the caller: a second reader of
// the same terminal loses keys to the first — see above.  With no terminal to
// read it returns at once, which is right: nobody is watching.
func WaitKey(out *os.File, prompt string) {
	_, keys, err := keyboard()
	if err != nil {
		return
	}
	if prompt != "" {
		_, _ = out.WriteString(prompt)
	}
	<-keys
}

// kind is what a keypress turned out to be.
type kind int

const (
	keyRune kind = iota
	keyEsc
	keyCtrlC
	keyUp
	keyDown
	keyLeft
	keyRight
	keyTab
	keyShiftTab
	keyPgUp
	keyPgDn
	keyEnter
	keyBack
	keyKill
)

type key struct {
	kind kind
	r    rune
}

// readKeys turns the bytes of a raw terminal into keypresses until the
// terminal closes.  Runes are decoded whole, so a Russian letter is one key
// and not two bytes.
func readKeys(tty *os.File, out chan<- key) {
	defer close(out)
	r := bufio.NewReaderSize(tty, 64)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return
		}
		switch {
		case b == 0x03: // Ctrl-C
			out <- key{kind: keyCtrlC}
		case b == 0x09:
			out <- key{kind: keyTab}
		case b == 0x0d, b == 0x0a:
			out <- key{kind: keyEnter}
		case b == 0x7f, b == 0x08:
			out <- key{kind: keyBack}
		case b == 0x15: // Ctrl-U, the way a shell clears a line
			out <- key{kind: keyKill}
		case b == 0x1b:
			out <- escape(r)
		case b < 0x20:
			// other control bytes carry nothing the screen answers
		case b < utf8.RuneSelf:
			out <- key{kind: keyRune, r: rune(b)}
		default:
			if err := r.UnreadByte(); err != nil {
				return
			}
			ru, _, err := r.ReadRune()
			if err != nil {
				return
			}
			out <- key{kind: keyRune, r: ru}
		}
	}
}

// escape reads what follows ESC.  A lone ESC — nothing buffered behind it — is
// the reader asking to leave.
func escape(r *bufio.Reader) key {
	if r.Buffered() == 0 {
		return key{kind: keyEsc}
	}
	b, err := r.ReadByte()
	if err != nil {
		return key{kind: keyEsc}
	}
	if b == 'O' { // the application-cursor spelling of the arrows
		if c, err := r.ReadByte(); err == nil {
			return arrow(c)
		}
		return key{kind: keyEsc}
	}
	if b != '[' {
		return key{kind: keyEsc}
	}
	// CSI: digits and semicolons, then one final letter.
	var params []byte
	for {
		c, err := r.ReadByte()
		if err != nil {
			return key{kind: keyEsc}
		}
		if c >= '0' && c <= '9' || c == ';' {
			params = append(params, c)
			continue
		}
		switch c {
		case 'A', 'B', 'C', 'D':
			return arrow(c)
		case 'Z':
			return key{kind: keyShiftTab}
		case '~':
			switch string(params) {
			case "5":
				return key{kind: keyPgUp}
			case "6":
				return key{kind: keyPgDn}
			}
		}
		return key{kind: keyRune, r: 0}
	}
}

func arrow(c byte) key {
	switch c {
	case 'A':
		return key{kind: keyUp}
	case 'B':
		return key{kind: keyDown}
	case 'C':
		return key{kind: keyRight}
	case 'D':
		return key{kind: keyLeft}
	}
	return key{kind: keyRune, r: 0}
}
