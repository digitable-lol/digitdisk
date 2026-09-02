// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

// Английские половины экрана обхода — того, что показывает `analyze`, пока
// дерево обходится, и того, чем оно потом становится: разделы, полосы
// наполнения, отметки, план уборки и подтверждение числом.
//
// Ширина здесь такая же часть перевода, как и на экране состояния. Подпись
// живёт в поле на 14 ячеек, состояние в подвале — в фишке, а подсказка клавиш
// выбирается по остатку ширины и переведена отдельно на каждую ширину:
// то же самое, сказанное короче, — это другая строка, а не обрезанная первая.
//
// Отдельная забота — слова, которыми экран говорит про уборку. «Перенесено»
// и «место НЕ освобождено» стоят рядом не случайно: по-английски точно так же
// нельзя дать понять, что место освободилось, — moved и NOT freed.
func init() {
	add(map[string]string{
		// Разделы. Полоса из восьми имён должна помещаться в ту же ширину,
		// что и по-русски, иначе на узком терминале она свернётся в стрелки
		// раньше, чем по-русски, и экран станет двумя разными экранами.
		"ИТОГ":         "TOTAL",
		"ДЕРЕВО":       "TREE",
		"КРУПНЕЙШЕЕ":   "LARGEST",
		"МОЖНО УБРАТЬ": "REMOVABLE",
		"РАЗРЯДЫ":      "CLASSES",
		"ПРОПУЩЕНО":    "SKIPPED",
		"МЕСТА":        "PLACES",
		"ЖУРНАЛ":       "JOURNAL",

		// Шапка и полоса состояния.
		"обход ":            "walking ",
		"обойдено ":         "walked ",
		"каталог не выбран": "no directory chosen",
		" ИДЁТ ОБХОД ":      " WALKING ",
		" СВОДИТСЯ ДЕРЕВО ": " ADDING UP THE TREE ",
		" ГОТОВО ":          " DONE ",
		" ОБХОД ":           " WALK ",
		"  %s записей · %s": "  %s entries · %s",

		// Подсказки клавиш, по одной на каждую ширину.
		"Tab разделы · ↑ ↓ строка · → внутрь · Пробел отметить · c убрать · l язык · ? клавиши · q выход ": "Tab sections · ↑ ↓ row · → into · Space mark · c clean · l lang · ? keys · q quit ",
		"Tab · ↑ ↓ · → внутрь · Пробел отметить · c убрать · l язык · ? · q выход ":                        "Tab · ↑ ↓ · → into · Space mark · c clean · l lang · ? · q quit ",
		"Tab · ↑ ↓ · → ← · c убрать · l · ? · q выход ":                                                    "Tab · ↑ ↓ · → ← · c clean · l · ? · q quit ",
		"? клавиши · q выход ":       "? keys · q quit ",
		"Esc — отменить, q — выход ": "Esc cancel, q quit ",
		"Esc · q ": "Esc · q ",
		"любая клавиша — закрыть ": "any key closes ",
		"закрыть ": "close ",
		"↑ ↓ прокрутка · l язык · q прервать обход ": "↑ ↓ scroll · l lang · q stop the walk ",
		"↑ ↓ · l · q прервать ":                      "↑ ↓ · l · q stop ",
		"q прервать ":                                "q stop ",

		// Числа обхода, пока он идёт.
		"записей":   "entries",
		"объём":     "size",
		"скорость":  "rate",
		"сейчас":    "now at",
		"записей/с": "entries/s",
		"файлов %s · каталогов %s · ссылок %s":              "files %s · dirs %s · links %s",
		"файлов %s · каталогов %s · ссылок %s · прочего %s": "files %s · dirs %s · links %s · other %s",
		"%s записей/с · %s/с":                               "%s entries/s · %s/s",
		"глубина %d":                                        "depth %d",
		"  пик %s/с":                                        "  peak %s/s",

		// Наполнение и его предварительность.
		"  ЧЕМ НАПОЛНЕНО  ":      "  WHAT FILLS IT  ",
		"ЧЕМ НАПОЛНЕНО":          "WHAT FILLS IT",
		" ПРЕДВАРИТЕЛЬНО ":       " PRELIMINARY ",
		"· файлы этого каталога": "· the files of this directory",
		"%s зап.": "%s ent.",
		"пусто":   "empty",
		"каталогов больше %s — дальше считается всё, но ходить можно не везде":    "past %s directories: everything is still counted, not everything can be walked",
		"каталогов больше %s — посчитано всё, но ходить можно не по всему дереву": "past %s directories: everything is counted, not all of the tree can be walked",

		// Итог обхода.
		"время":            "took",
		"жёсткие ссылки":   "hard links",
		"пропущено":        "skipped",
		"раздел ПРОПУЩЕНО": "the SKIPPED section",
		"видимый размер, как у du --apparent-size":                                      "apparent size, as du --apparent-size counts it",
		"%s повторных имён не засчитано (%s)":                                           "%s repeated names not counted (%s)",
		"настоящий разбор не выполнялся: разряды показывают работу стыковки, не анализ": "no real analysis was made: the classes show the wiring, not the disk",
		"решающий слой %s, договор версии %d":                                           "decision layer %s, contract version %d",

		// Ходьба по дереву.
		"каталог":           "directory",
		"в нём":             "holds",
		"%s · %s записей":   "%s · %s entries",
		"подкаталогов %s":   "subdirectories %s",
		"дерево не собрано": "the tree was not built",
		"крупнейший файл прямо здесь: %s (%s)":                     "largest file right here: %s (%s)",
		"Пробел — отметить каталог, «.» — этот, «c» — план уборки": "Space marks a directory, «.» this one, «c» plans the cleaning",
		"отмечено каталогов: %d — «c» покажет план уборки":         "directories marked: %d — «c» shows the cleaning plan",

		// Крупнейшее и можно убрать.
		"КРУПНЕЙШИЕ ЗАПИСИ — %d": "LARGEST ENTRIES — %d",
		"нечего":          "nothing",
		"предложено":      "proposed",
		"%s записей · %s": "%s entries · %s",
		"убрать можно отсюда: «c» строит план и спрашивает число файлов.":                      "cleaning starts here: «c» builds a plan and asks for the number of files.",
		"отметьте каталоги в разделе ДЕРЕВО, чтобы взять только их; без отметок — всё дерево.": "mark directories in TREE to take only those; with none marked, the whole tree.",
		"та же уборка одной командой, если экран не нужен:":                                    "the same cleaning as one command, if the screen is not wanted:",
		"без --apply она печатает план и не трогает ни одного файла":                           "without --apply it prints the plan and touches no file",

		// Пропущенное.
		"нет доступа":    "no access",
		"исчезло":        "vanished",
		"иные ошибки":    "other errors",
		"граница ФС":     "fs boundary",
		"предел глубины": "depth limit",
		"ПРИМЕРЫ":        "EXAMPLES",

		// Известные места и журнал.
		"справочник":                         "directory of places",
		"нашлось здесь":                      "found here",
		"корзин":                             "trashes",
		"в корзинах":                         "in the trashes",
		"освобождено":                        "freed",
		"КОРЗИНЫ":                            "TRASHES",
		"ЕСТЬ НА ЭТОЙ МАШИНЕ":                "PRESENT ON THIS MACHINE",
		"размеры считает `digitdisk places`": "sizes are measured by `digitdisk places`",
		"ни одного известного места на этой машине не нашлось": "not one known place was found on this machine",
		"справочник не прочитан: %s":                           "the directory of places was not read: %s",
		"журнал не прочитан: %s":                               "the journal was not read: %s",
		"хранилище %s":                                         "store %s",
		"место не освобождено, пока не purge":                  "no space is freed until purge",
		"это уже стёрто и не вернётся":                         "this is erased and will not come back",
		"под этим корнем ещё не убирали":                       "nothing has been cleaned under this root yet",
		"%s файл.":      "%s files",
		"беда: %s":      "trouble: %s",
		"стёрто %d":     "erased %d",
		"возвращено %d": "restored %d",
		"Enter — вернуть выбранную корзину на прежние места (спросит число файлов)":  "Enter puts the chosen trash back (it will ask for the number of files)",
		"стирание корзины — только командой `digitdisk purge <корзина> --confirm N`": "erasing a trash is only `digitdisk purge <trash> --confirm N`",

		// Приглашение пути.
		"ОБОЙТИ КАТАЛОГ":          "WALK A DIRECTORY",
		"путь: ":                  "path: ",
		"это не каталог":          "this is not a directory",
		"каталог не читается: %s": "the directory does not read: %s",
		"подходит каталогов: %d":  "directories matching: %d",
		"   …и ещё %d":            "   …and %d more",
		"Tab — дополнить, Ctrl-U — стереть строку, Enter — обойти, Esc — отменить": "Tab completes, Ctrl-U clears the line, Enter walks, Esc cancels",

		// Отметка.
		"ОТМЕТКА": "MARKING",
		"Отмечается каталог, а не строка «%s».":                           "A directory is marked, not the «%s» row.",
		"Она стоит за файлы, которые лежат прямо здесь; чтобы взять их,":  "It stands for the files lying right here; to take those,",
		"отметьте сам этот каталог — «.» отмечает тот, в котором стоите.": "mark this directory itself — «.» marks the one you stand in.",

		// План уборки и подтверждение.
		"ПЛАН УБОРКИ":                       "CLEANING PLAN",
		"СТРОИТСЯ ПЛАН УБОРКИ":              "BUILDING THE CLEANING PLAN",
		"ПЕРЕНОС В КОРЗИНУ":                 "MOVING TO THE TRASH",
		"УБРАНО В КОРЗИНУ":                  "MOVED TO THE TRASH",
		"УБОРКА":                            "CLEANING",
		"УБОРКА НЕ ВЫШЛА":                   "THE CLEANING DID NOT HAPPEN",
		"ОБХОД НЕ ВЫШЕЛ":                    "THE WALK DID NOT HAPPEN",
		"идёт работа…":                      "working…",
		"этот экран собран без уборки":      "this screen was built without cleaning",
		"этот экран собран без возврата":    "this screen was built without restoring",
		"этот экран собран без журнала":     "this screen was built without the journal",
		"этот экран собран без справочника": "this screen was built without the directory of places",
		"отмечено каталогов: %d":            "directories marked: %d",
		"отмечено ничего — план по всему обойденному дереву":                           "nothing marked — the plan covers the whole walked tree",
		"дерево обходится заново, и заново судится решающим слоем:":                    "the tree is walked again and judged again by the decision layer:",
		"план строится по тому, что на диске сейчас, а не по тому, что было.":          "the plan is about what is on the disk now, not about what was.",
		"ничего не открывается на запись и корзина не создаётся.":                      "nothing is opened for writing and no trash is created.",
		"к переносу      %s файлов, %s":                                                "to be moved   %s files, %s",
		"место не освободится: перенос — это переименование, байты остаются на диске.": "no space is freed: the move is a rename, the bytes stay on the disk.",
		"освободит только `digitdisk purge <корзина> --confirm N`.":                    "only `digitdisk purge <trash> --confirm N` frees them.",
		"по разрядам:": "by class:",
		"%s файлов":    "%s files",
		"решающий слой не пометил «МожноУбрать» ничего из отмеченного.":       "the decision layer marked nothing «МожноУбрать» in what was marked.",
		"отметка не делает файл убираемым — приговор выносит ядро.":           "a mark does not make a file removable — the core passes the verdict.",
		"защитный список оставил на месте %d (%s)":                            "the protect list left %d in place (%s)",
		"хозяин отказался трогать %d — слои разошлись, см. `digitdisk clean`": "the host refused %d — the layers disagree, see `digitdisk clean`",
		"корзина: %s":       "trash: %s",
		"первые из списка:": "the first of the list:",
		"наберите число файлов (%d) и Enter — перенести; Ctrl-U стереть, Esc отменить": "type the number of files (%d) and Enter to move; Ctrl-U clears, Esc cancels",
		"наберите число файлов (%d) и Enter — вернуть; Ctrl-U стереть, Esc отменить":   "type the number of files (%d) and Enter to restore; Ctrl-U clears, Esc cancels",
		"число файлов: ": "number of files: ",
		"названо %q, а файлов %d — ничего не тронуто": "named %q, but the files are %d — nothing was touched",
		"файлов %s": "files %s",
		"журнал пишется до того, как сдвинется первый файл.":               "the journal is written before the first file moves.",
		"перенесено      %s файлов, %s":                                    "moved         %s files, %s",
		"место НЕ освобождено: файлы лежат в корзине под другими именами.": "no space is freed: the files lie in the trash under other names.",
		"вернуть — раздел ЖУРНАЛ, клавиша Enter на этой корзине.":          "to put back: the JOURNAL section, Enter on this trash.",
		"стереть насовсем — только отдельной командой:":                    "erasing for good is a separate command only:",
		"не тронуто %d — файл изменился между обходом и переносом:":        "%d not touched — the file changed between the walk and the move:",

		// Возврат.
		"ВОЗВРАТ":            "RESTORE",
		"ВОЗВРАТ ИДЁТ":       "RESTORING",
		"СЧИТАЕТСЯ ВОЗВРАТ":  "WORKING OUT THE RESTORE",
		"ВОЗВРАТ ИЗ КОРЗИНЫ": "RESTORE FROM THE TRASH",
		"ВОЗВРАТ НЕ ВЫШЕЛ":   "THE RESTORE DID NOT HAPPEN",
		"ВОЗВРАЩАТЬ НЕЧЕГО":  "NOTHING TO RESTORE",
		"ВОЗВРАЩЕНО":         "RESTORED",
		"журнал корзины читается; ничего не двигается.":                                        "the trash journal is read; nothing moves.",
		"в этой корзине нет файлов, которые можно вернуть":                                     "this trash holds no file that can be put back",
		"вернётся на прежние места   %s файлов, %s":                                            "would go back to where they were   %s files, %s",
		"вернулось на прежние места  %s файлов, %s":                                            "went back to where they were  %s files, %s",
		"ничего не перезаписывается: файл, вернувшийся на занятое место, останется в корзине.": "nothing is overwritten: a file whose place is taken stays in the trash.",

		// Страница клавиш.
		"КЛАВИШИ И КОМАНДЫ":                             "KEYS AND COMMANDS",
		"разделы":                                       "sections",
		"строка; g G — начало и конец":                  "row; g G to the top and the bottom",
		"внутрь каталога (в разделе ДЕРЕВО)":            "into a directory (in the TREE section)",
		"назад из каталога":                             "back out of a directory",
		"отметить каталог; «.» — тот, в котором стоите": "mark a directory; «.» the one you stand in",
		"план уборки по отмеченному и подтверждение":    "the cleaning plan for what is marked, and its confirmation",
		"обойти другой каталог (Tab дополняет путь)":    "walk another directory (Tab completes the path)",
		"в ЖУРНАЛЕ": "in JOURNAL",
		"язык экрана: русский или English":   "the language of the screen: Russian or English",
		"эта справка":                        "this help",
		"←, забой":                           "←, Backspace",
		"Пробел":                             "Space",
		"любая клавиша — закрыть":            "any key closes it",
		"выход; отчёт печатается как всегда": "quit; the report is printed as always",
		"уборка с экрана идёт тем же путём, что `digitdisk clean`: приговор ядра, план, подтверждение числом.": "cleaning from the screen goes the way `digitdisk clean` goes: the core's verdict, a plan, a confirmation by count.",
		"стирание — только отдельной командой `digitdisk purge`: это единственный необратимый шаг.":            "erasing is `digitdisk purge` and nothing else: it is the one step that cannot be undone.",

		// Отказы экрана.
		"обход прерван":                  "the walk was stopped",
		"живому экрану не передан обход": "the live screen was given no walk",
	})
}
