// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

// Английские половины всего, что говорит живой экран: названия разделов,
// подписи столбцов и датчиков, заметки под ними, состояние замера в подвале и
// подсказка клавиш.
//
// Ширина здесь — часть перевода. Экран рисуется в колонках: подпись живёт в
// поле на 12, 16 или 26 ячеек, полоса начинается сразу за ней, а подсказка
// клавиш выбирается по остатку ширины. Английская половина, которая длиннее
// русской, не переводит строку, а ломает вёрстку — поэтому она такой же длины
// или короче. Строка подсказки переведена отдельно на каждую ширину: то же
// самое, сказанное короче, — это другая строка, а не обрезанная первая.
func init() {
	add(map[string]string{
		// Разделы. Они же — заголовки в отчёте, и полоса разделов
		// целиком должна помещаться в ту же ширину, что и по-русски.
		"ОБЗОР":        "OVERVIEW",
		"СИСТЕМА":      "SYSTEM",
		"ЗАГРУЗКА":     "LOAD",
		"ПАМЯТЬ":       "MEMORY",
		"ПРОЦЕССЫ":     "PROCESSES",
		"ДИСКИ":        "DISKS",
		"СЕТЬ":         "NETWORK",
		"ТЕМПЕРАТУРА":  "TEMPERATURE",
		"НЕ ПРОЧИТАНО": "NOT READ",

		// Подвал: состояние замера, счёт строк и клавиши.
		"ПАУЗА":                "PAUSED",
		"ЗАМЕР":                "SAMPLING",
		"ЖИВОЙ":                "LIVE",
		"q выход ":             "q quit ",
		"  строки %d–%d из %d": "  lines %d–%d of %d",
		"  замер %s назад · длился %s · каждые %s": "  sampled %s ago · took %s · every %s",
		"  замер %s назад · каждые %s":             "  sampled %s ago · every %s",
		"  замер %s назад":                         "  sampled %s ago",

		"← → разделы · ↑ ↓ прокрутка · p пауза · r замер · l язык · 1 КОМАНДЫ · q выход ": "← → sections · ↑ ↓ scroll · p pause · r sample · l lang · 1 COMMANDS · q quit ",
		"← → разделы · p пауза · r замер · l язык · 1 КОМАНДЫ · q выход ":                 "← → sections · p pause · r sample · l lang · 1 COMMANDS · q quit ",
		"← → · p · r · l язык · 1 КОМАНДЫ · q выход ":                                     "← → · p · r · l lang · 1 COMMANDS · q quit ",
		"1 КОМАНДЫ · q выход ": "1 COMMANDS · q quit ",

		// Подвал раздела КОМАНДЫ: там клавиши другие — выбор и запуск.
		"↑ ↓ и 1…7 выбрать · Enter запустить · ← → разделы · l язык · q выход ": "↑ ↓ 1…7 choose · Enter runs · ← → sections · l lang · q quit ",
		"↑ ↓ выбрать · Enter запустить · ← → разделы · q выход ":                "↑ ↓ choose · Enter runs · ← → sections · q quit ",
		"↑ ↓ · Enter запустить · q выход ":                                      "↑ ↓ · Enter runs · q quit ",
		"Enter запустить · q выход ":                                            "Enter runs · q quit ",

		// Обзор.
		"замер идёт…":       "sampling…",
		"ЦП занято":         "CPU busy",
		"замер не делался":  "not sampled",
		"Память":            "Memory",
		"Своп":              "Swap",
		"нет":               "none",
		"история":           "history",
		"%s из %s  (%s)":    "%s of %s  (%s)",
		"%s из %s":          "%s of %s",
		"%s свободно из %s": "%s free of %s",
		"%s своб.":          "%s free",
		"%s / %s / %s   (1/5/15 мин, ядер %s)":              "%s / %s / %s   (1/5/15 min, %s cores)",
		"…и ещё %d — раздел ДИСКИ":                          "…and %d more — the DISKS section",
		"не прочитано источников: %d — раздел НЕ ПРОЧИТАНО": "sources not read: %d — the NOT READ section",

		// Система.
		"узел":         "host",
		"дистрибутив":  "distro",
		"ядро":         "kernel",
		"время работы": "uptime",
		"снимок взят":  "snapshot at",
		"%s (с %s)":    "%s (since %s)",

		// Загрузка.
		"средняя":                    "load avg",
		"ядер":                       "cores",
		"в очереди":                  "runnable",
		"занято ЦП":                  "CPU busy",
		"%d из %d":                   "%d of %d",
		"%s (замер %d мс)":           "%s (sampled %d ms)",
		"%s / %s / %s  (1/5/15 мин)": "%s / %s / %s  (1/5/15 min)",

		// Память.
		"занято":     "used",
		"всего":      "total",
		"свободно":   "free",
		"доступно":   "available",
		"кэш/буферы": "cache/buffers",
		"своп":       "swap",
		"%s  (в т.ч. разделяемая %s)":              "%s  (shared %s of it)",
		"%s из %s занято":                          "%s of %s used",
		"занято = всего − доступно, как в free(1)": "used = total − available, as in free(1)",

		// Процессы.
		"процессы":      "processes",
		"владелец":      "owner",
		"команда":       "command",
		"память":        "memory",
		"ЦП":            "CPU",
		"ПО ПАМЯТИ":     "BY MEMORY",
		"ПО ПРОЦЕССОРУ": "BY CPU",

		// Диски.
		"точка монтирования":           "mount point",
		"размер":                       "size",
		"занят":                        "use",
		"ошибка: %s":                   "error: %s",
		"  ·  только чтение":           "  ·  read-only",
		"  ·  inode %s свободно из %s": "  ·  inodes %s free of %s",

		// Сеть.
		"интерфейс": "interface",
		"состоян.":  "state",
		"принято":   "received",
		"передано":  "sent",
		"пак. вх.":  "pkts in",
		"пак. исх.": "pkts out",
		"%d Мбит/с": "%d Mbit/s",
		"ошибок %d/%d, потеряно %d/%d": "errors %d/%d, dropped %d/%d",

		// Температура.
		"— (датчиков не нашлось)": "— (no sensors found)",
		"из 100 °C":     "of 100 °C",
		"критич. %s °C": "critical %s °C",

		// Не прочитано.
		"прочитано всё, чего ждали": "everything expected was read",
		"источник назван вместе с причиной; нулём его отсутствие не притворяется": "the source is named with its reason; its absence never poses as a zero",

		// Раздел КОМАНДЫ: он же список, он же место, откуда работают.
		"КОМАНДЫ": "COMMANDS",
		"l — язык вывода (ru ⇄ en), выбор запоминается":                                "l — output language (ru ⇄ en), the choice is kept",
		"Ключи: digitdisk --help.  Подробно: man digitdisk":                            "Flags: digitdisk --help.  In full: man digitdisk",
		"Enter — новый замер: этот экран и есть status":                                "Enter — a fresh sample: this screen IS status",
		"Enter — выполнить `digitdisk %s` и вернуться сюда":                            "Enter — run `digitdisk %s` and come back here",
		"Enter — спросить путь (предложит текущий каталог) и выполнить `digitdisk %s`": "Enter — ask for a path (the current directory offered) and run `digitdisk %s`",
		"не отсюда: %s": "not from here: %s",

		// Отказы живого экрана.
		"вывод не в терминал: живой экран невозможен": "output is not a terminal: no live screen",
		"живому экрану не передан сборщик снимка":     "the live screen was given no snapshot collector",
	})
}
