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

		// Подсказки клавиш, по одной на каждую ширину. Забой назван в
		// каждой, где вообще названы клавиши, и назван тем, что делает:
		// «стереть насовсем» — erase for good. Клавиша без корзины не
		// может в подвале выглядеть мягче, чем в вопросе, который она
		// открывает, ни на одном из двух языков.
		"Tab разделы · ↑ ↓ строка · → внутрь · Пробел отметить · c убрать · забой стереть насовсем · l язык · ? клавиши · q выход ": "Tab sections · ↑ ↓ row · → into · Space mark · c clean · Backspace erase for good · l lang · ? keys · q quit ",
		"Tab · ↑ ↓ · → внутрь · Пробел отметить · c убрать · забой стереть насовсем · l язык · ? · q выход ":                        "Tab · ↑ ↓ · → into · Space mark · c clean · Backspace erase for good · l lang · ? · q quit ",
		"Tab · ↑ ↓ · → ← · c убрать · забой стереть насовсем · l · ? · q выход ":                                                    "Tab · ↑ ↓ · → ← · c clean · Backspace erase for good · l · ? · q quit ",
		"c убрать · забой стереть насовсем · ? · q выход ":                                                                          "c clean · Backspace erase for good · ? · q quit ",
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
		"крупнейший файл прямо здесь: %s (%s)":                                               "largest file right here: %s (%s)",
		"Пробел — отметить каталог, «.» — этот, «c» — план уборки, забой — стереть насовсем": "Space marks a directory, «.» this one, «c» plans the cleaning, Backspace erases for good",
		"отмечено каталогов: %d — «c» в корзину, забой стереть насовсем":                     "directories marked: %d — «c» to the trash, Backspace erases for good",

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
		"предложен текущий каталог — Enter соглашается с ним.":                     "the current directory is offered — Enter agrees with it.",
		"обходится всё дерево под ним: на домашнем каталоге это миллионы":          "the whole tree under it is walked: on a home directory that is millions",
		"записей и минуты. Числа идут с первой секунды, q прерывает обход.":        "of entries and minutes. Numbers move from the first second, q stops it.",

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
		"путей %s":  "paths %s",
		"журнал пишется до того, как сдвинется первый файл.":               "the journal is written before the first file moves.",
		"перенесено      %s файлов, %s":                                    "moved         %s files, %s",
		"место НЕ освобождено: файлы лежат в корзине под другими именами.": "no space is freed: the files lie in the trash under other names.",
		"вернуть — раздел ЖУРНАЛ, клавиша Enter на этой корзине.":          "to put back: the JOURNAL section, Enter on this trash.",
		"стереть насовсем — только отдельной командой:":                    "erasing for good is a separate command only:",
		"не тронуто %d — файл изменился между обходом и переносом:":        "%d not touched — the file changed between the walk and the move:",

		// Стирание насовсем — то же, что уборка, кроме корзины, и по-английски
		// разница обязана быть слышна там же, где по-русски: «убрать» против
		// «стереть насовсем» — clean против erase for good. Слово «erase» без
		// «for good» уже занято корзиной (`purge`), и одного его мало.
		"СТЕРЕТЬ НАСОВСЕМ":               "ERASE FOR GOOD",
		"СЧИТАЕТСЯ СТИРАНИЕ":             "WORKING OUT THE ERASURE",
		"СТИРАНИЕ ИДЁТ":                  "ERASING",
		"СТЁРТО НАСОВСЕМ":                "ERASED FOR GOOD",
		"СТИРАТЬ НЕЧЕГО":                 "NOTHING TO ERASE",
		"СТИРАНИЕ НЕ ВЫШЛО":              "THE ERASURE DID NOT HAPPEN",
		"этот экран собран без стирания": "this screen was built without erasing",
		"уйдёт всё, что там лежит, — ядро спрашивается о том, что это,": "everything lying there will go — the core is asked WHAT it is,",
		"а не о том, можно ли. пока считается, не тронут ни один файл.": "not whether it may. while it is being worked out, not one file is touched.",
		"дерево обходится заново: считается то, что на диске СЕЙЧАС.":   "the tree is walked again: what is on the disk NOW is what is counted.",
		"под курсором: %s": "under the cursor: %s",
		"здесь нет строки, на которую наведён курсор":                                 "there is no row here for the cursor to stand on",
		"Стирается каталог, а не строка «%s».":                                        "A directory is erased, not the row «%s».",
		"исчезнет насовсем  %s путей: %s файлов, %s каталогов, %s":                    "vanishing for good  %s paths: %s files, %s directories, %s",
		"освободится %s — на этот раз по-настоящему; корзины не будет, возврата нет.": "%s will be freed — really freed, this time; there will be no trash and no way back.",
		"исчезнет:":           "what will vanish:",
		"к стиранию 0 путей.": "0 paths to erase.",
		"защитный список оставил на месте %d (%s):":                                          "the protect list left %d in place (%s):",
		"хозяин отказался трогать %d:":                                                       "the host refused %d:",
		"не тронуто %d — ссылки и недоступное; их забой не стирает":                          "%d not touched — links and the unreadable; Backspace does not erase those",
		"y — стереть насовсем (%d); Esc — отменить, ничего не тронуть":                       "y erases for good (%d); Esc cancels and touches nothing",
		"наберите число путей (%d) и Enter — стереть насовсем; Ctrl-U стереть, Esc отменить": "type the number of paths (%d) and Enter to erase for good; Ctrl-U clears, Esc cancels",
		"стирание подтверждает «y», а не Enter":                                              "an erasure is confirmed by «y», not by Enter",
		"путей %d — столько одной клавишей не стирается, назовите число":                     "%d paths — that many are not erased on one key, name the number",
		"журнал пишется до того, как исчезнет первый файл.":                                  "the journal is written before the first file vanishes.",
		"стёрто          %s файлов и %s каталогов, %s":                                       "erased        %s files and %s directories, %s",
		"место освобождено: этого на диске больше нет, и возврата нет.":                      "the space is freed: this is off the disk, and there is no way back.",
		"журнал: %s": "journal: %s",
		"он называет всё, что исчезло, и стоит в разделе ЖУРНАЛ как стирание.":                   "it names everything that vanished and stands in the JOURNAL section as an erasure.",
		"не тронуто %d — файл изменился между обходом и стиранием:":                              "%d not touched — the file changed between the walk and the erasure:",
		"забой в разделе ДЕРЕВО стирает то же самое насовсем, минуя корзину.":                    "Backspace in the TREE section erases the same for good, past the trash.",
		"стёрто насовсем %d — возврата нет":                                                      "%d erased for good — no way back",
		"запись «стёрто насовсем» вернуть нельзя: файлов нет, есть только список того, что было": "an «erased for good» entry cannot be restored: the files are gone, only the list of them is left",

		// Забой по воле человека: слова второго вопроса и трёхступенчатое
		// подтверждение.  «Строгость» — имя решающего слоя, и переводится
		// только слово, которым её называет экран.
		"ядро зовёт это: %s": "the core calls this: %s",
		"%s %d":              "%s %d",
		"не названо":         "unnamed",
		"ядро не сказало, что это: сборка без решающего слоя.":                           "the core did not say what this is: built without a decision layer.",
		"ЭТО НЕ МУСОР: под присмотром системы версий — пропадёт история, а не файл.":     "THIS IS NOT RUBBISH: under version control — the history goes, not a file.",
		"ЭТО НЕ МУСОР: хранилище по содержимому — сломается целость, а не файл.":         "THIS IS NOT RUBBISH: a content-addressed store — its integrity breaks, not a file.",
		"ЭТО НЕ МУСОР: исходники — то, что писали руками, и заново их никто не сделает.": "THIS IS NOT RUBBISH: source — written by hand, and nobody will write it again.",
		"ЭТО НЕ МУСОР: ядро не знает, что это, и убирать не советовало.":                 "THIS IS NOT RUBBISH: the core does not know what it is and never advised removing it.",
		"ЭТО СВЕЖЕЕ: разряд мусорный, но ядро ещё не считает это остывшим.":              "THIS IS FRESH: a rubbish class, but the core does not call it cold yet.",
		"ЭТО НЕ МУСОР: ядро назвать это не смогло — считайте, что оно нужное.":           "THIS IS NOT RUBBISH: the core could not name it — take it as needed.",
		"СТЕРЕТЬ": "ERASE",
		"наберите число путей (%d) и Enter, потом слово %s; Ctrl-U стереть, Esc отменить": "type the number of paths (%d) and Enter, then the word %s; Ctrl-U clears, Esc cancels",
		"наберите слово %s и Enter — стереть насовсем; Esc отменить":                      "type the word %s and Enter to erase for good; Esc cancels",
		"набрано %q, а нужно слово %s — ничего не тронуто":                                "%q was typed, and the word wanted is %s — nothing was touched",
		"названо %q, а путей %d — ничего не тронуто":                                      "named %q, but the paths are %d — nothing was touched",
		"это не мусор — одной клавишей такое не стирается, назовите число путей":          "this is not rubbish — one key does not erase it, name the number of paths",
		"число путей: ": "number of paths: ",
		"слово: ":       "word: ",
		"здесь нет ни одного обычного файла и ни одного каталога, который можно снять.":          "there is not one ordinary file and not one directory here that can be taken.",
		"забой стирает то, на что указали, каким бы ядро его ни считало, — но пусто есть пусто.": "Backspace erases what you pointed at, whatever the core calls it — but empty is empty.",
		"каталогов осталось %d — в них ещё что-то лежит:":                                        "%d directories are still there — something is still lying in them:",

		// Твёрдые запреты: заголовок и подпись к обходу.  Сами причины и обходы
		// едут lang.Phrase и переводятся своими статьями в dict_clean.go.
		"СТИРАТЬ ОТСЮДА НЕЛЬЗЯ":           "NOTHING IS ERASED FROM HERE",
		"Как быть, если вы всё же правы:": "What to do if you are right after all:",

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
		"язык экрана: русский или English": "the language of the screen: Russian or English",
		"эта справка":                      "this help",
		"←":                                "←",
		"забой":                            "Backspace",
		"СТЕРЕТЬ НАСОВСЕМ отмеченное, а без отметок — строку под курсором":                   "ERASE FOR GOOD what is marked, or the row under the cursor when nothing is",
		"стирает то, на что указали, мусор это или нет: ядро предупреждает, но не запрещает": "erases what you pointed at, rubbish or not: the core warns, it does not forbid",
		"Пробел": "Space",
		"любая клавиша — закрыть":            "any key closes it",
		"выход; отчёт печатается как всегда": "quit; the report is printed as always",
		"уборка с экрана идёт тем же путём, что `digitdisk clean`: приговор ядра, план, подтверждение числом.":                          "cleaning from the screen goes the way `digitdisk clean` goes: the core's verdict, a plan, a confirmation by count.",
		"уборка «c» ищет мусор САМА: что убрать, решает приговор ядра, и «НеТрогать» она не берёт.":                                     "cleaning with «c» finds the rubbish ITSELF: what goes is the core's verdict, and «НеТрогать» it does not take.",
		"забой стирает ВОТ ЭТО: всё, что лежит под отмеченным, и сами каталоги, без корзины и без возврата.":                            "Backspace erases THIS: everything under what is marked, and the directories too, with no trash and no way back.",
		"ядро при забое спрашивается о том, ЧТО ЭТО (исходники, кэш, хранилище) — и предупреждает, а не отказывает.":                    "on Backspace the core is asked WHAT THIS IS (source, cache, a store) — and it warns, it does not refuse.",
		"мусор — одна «y»; не мусор, длинный список или крупное — число путей; хранилище и git — число и слово.":                        "rubbish: one «y»; not rubbish, a cut list or something large: the number of paths; a store or git: the number and a word.",
		"твёрдо отказано только там, где ломается машина или сам инструмент: корень, система, дом целиком, каталог digitdisk, корзина.": "the hard refusals are only where the machine or the tool itself breaks: the root, the system, a whole home, digitdisk's own place, a trash.",
		"корзину целиком стирает отдельная команда `digitdisk purge`; с экрана она не запускается.":                                     "a whole trash is erased by the separate command `digitdisk purge`; no screen starts it.",

		// Отказы экрана.
		"обход прерван":                  "the walk was stopped",
		"живому экрану не передан обход": "the live screen was given no walk",
	})
}
