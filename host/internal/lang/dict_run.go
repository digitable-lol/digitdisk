// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

// Английские половины того, что говорит обёртка чужой команды: строка
// состояния, сводка и отказы запуска.
//
// Строка состояния печатается в одну строку терминала и режется по ширине,
// поэтому английская половина не длиннее русской там, где это возможно:
// строка, которая на втором языке не влезает, теряет последнее поле — и
// человек не узнает, что оно было.
func init() {
	add(map[string]string{
		// ── строка состояния ────────────────────────────────────────
		"ЦП %s":             "CPU %s",
		"ЦП %s, средн. %s":  "CPU %s, avg %s",
		"память %s":         "memory %s",
		"память %s, пик %s": "memory %s, peak %s",
		"процессов %d":      "%d processes",
		"видеопамять %s":    "video memory %s",

		// ── сводка ──────────────────────────────────────────────────
		"команда «%s»: код %d, %s":            "command «%s»: code %d, %s",
		"команда «%s»: убита сигналом %s, %s": "command «%s»: killed by %s, %s",
		"процессорное время %s (в среднем %s), пик памяти %s, процессов %d; учёт — %s": "CPU time %s (%s on average), peak memory %s, processes %d; accounting — %s",
		"процессорное время %s (в среднем %s), пик памяти %s; учёт — %s":               "CPU time %s (%s on average), peak memory %s; accounting — %s",
		"около %s": "about %s",
		"%s (самый крупный процесс)": "%s (the largest single process)",
		"%s мс": "%s ms",
		"не измерен — команда короче замера": "not measured — the command was shorter than one sample",
		"%d мин %d с": "%d min %d s",
		"%d ч %d мин": "%d h %d min",

		// Чем считали. Это не украшение сводки: «6,1 ГиБ» и «около
		// 6,1 ГиБ» — разные утверждения, и читатель вправе знать, какое
		// из двух он держит в руках.
		"своя контрольная группа, считает ядро":              "a control group of our own, counted by the kernel",
		"обход /proc раз в %d мс — пик памяти приблизителен": "a /proc walk every %d ms — the memory peak is approximate",
		"итог ядра о команде и о том, чего она дождалась":    "the kernel's total for the command and what it waited for",

		// Про видеокарту говорится ровно то, что можно узнать честно:
		// память по процессам драйвер публикует, долю времени карты по
		// процессам — нет, и вместо выдуманного числа стоит слово.
		"видеокарта: пик памяти процессов %s (%s); загрузку по процессам драйвер не публикует": "video card: peak memory of its processes %s (%s); the driver publishes no per-process load",

		// ── отказы и пояснения ключей ───────────────────────────────
		"нужна команда: digitdisk run <команда> [доводы]":        "a command is needed: digitdisk run <command> [arguments]",
		"команда %q не найдена":                                  "command %q not found",
		"команда %q не запускается: %s":                          "command %q will not run: %s",
		"без строки состояния, даже в терминале":                 "no status line, even in a terminal",
		"период обновления строки состояния, мс":                 "status-line refresh period, ms",
		"спросить о видеопамяти программу драйвера (nvidia-smi)": "ask the driver's own program (nvidia-smi) about video memory",
	})
}
