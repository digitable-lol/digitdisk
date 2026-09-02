// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !flangcore

package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"digitdisk/internal/core"
	"digitdisk/internal/lang"
)

// Сборка без признака `flangcore` получает заглушку: она отвечает
// «Неизвестное/НеТрогать» на КАЖДУЮ запись и ничего не решает.
//
// ПОЧЕМУ ЗДЕСЬ ШУМ, А НЕ МОЛЧАНИЕ. Выпуск собирается с признаком —
// `scripts/build-release.sh` зовёт `go build -tags flangcore`, и у того, кто
// поставил digitdisk из архива или через brew, работает настоящее ядро на
// flang. Заглушку получает другой человек: тот, кто собрал дерево руками
// обычным `go build ./host`. Он при этом видит полноценный отчёт, в котором
// все до единой записи попали в «Неизвестное», и решает, что на диске просто
// нет ничего интересного. Тихо отданный слой, который на всё отвечает
// «не знаю», хуже отказа — поэтому он больше не тихий.
//
// Шум идёт в поток ОШИБОК, а не в вывод: `--json` обязан остаться
// машиночитаемым, а текстовый отчёт — прежним. Печатается один раз за прогон и
// ровно там, где заглушку собираются применить, — `status` решающего слоя не
// трогает и молчит.
var предупреждён sync.Once

// warningLines is the text of the box, line by line and without the box.
//
// The frame is drawn around it rather than typed into it because the two
// languages do not have the same line lengths, and a box whose right edge is
// typed by hand is a box that comes out ragged in the language nobody was
// looking at when they typed it.
func warningLines(l lang.Lang) []string {
	return []string{
		l.T("ВНИМАНИЕ: digitdisk собран БЕЗ решающего ядра."),
		"",
		l.T("Разбора не будет: каждая запись вернётся как «Неизвестное/НеТрогать»,"),
		l.T("и пустой список «предложено убрать» НЕ означает, что убирать нечего."),
		"",
		l.T("Соберите с ядром:  go build -tags flangcore -o digitdisk ./host"),
		l.T("Или поставьте выпуск: brew install digitable-lol/tap/digitdisk"),
		l.T("(выпуск всегда идёт с ядром — scripts/build-release.sh)"),
	}
}

// boxed frames the lines, padded to the longest of them.
func boxed(lines []string) string {
	width := 0
	for _, s := range lines {
		if n := len([]rune(s)); n > width {
			width = n
		}
	}
	var b strings.Builder
	b.WriteString("\n┌─" + strings.Repeat("─", width) + "─┐\n")
	for _, s := range lines {
		b.WriteString("│ " + s + strings.Repeat(" ", width-len([]rune(s))) + " │\n")
	}
	b.WriteString("└─" + strings.Repeat("─", width) + "─┘\n")
	return b.String()
}

// chosenDecider returns the placeholder decision layer and says so out loud.
// See coreflang/bridge.go for how to build against the real one.
func chosenDecider(l lang.Lang) core.Decider {
	предупреждён.Do(func() {
		fmt.Fprint(os.Stderr, boxed(warningLines(l)))
	})
	return core.Default()
}
