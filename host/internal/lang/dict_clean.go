// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

// Английские половины отказов и заметок нижних слоёв: сторож уборки, перенос в
// корзину, возврат и стирание, разбор справочника известных мест и защитного
// списка, причины неизмеренного.
//
// Ключ — русская строка ровно так, как она написана в исходнике: она же едет в
// JSON и в журнал, и потому не меняется ни на байт. Значения, которые человек
// пишет в справочник и в защитный список сам («кэш», «дом», «всюду», «путь»),
// в перечислениях оставлены по-русски: это не слова отчёта, а то, что надо
// набрать в файле, и переведённый список подсказывал бы набрать не то.
func init() {
	add(map[string]string{
		// Отметки экрана обхода сужают землю плана; путь вне корня в неё не
		// попадает, и об этом говорится, а не молчится.
		"ни один из отмеченных путей не лежит внутри %s": "not one of the marked paths lies inside %s",

		// ── отпечаток: чем файл перестал быть собой ──────────────────
		"это уже другой файл (был узел %d:%d, стал %d:%d)":    "this is a different file now (node was %d:%d, is %d:%d)",
		"размер изменился (был %d Б, стал %d Б)":              "the size changed (was %d B, is %d B)",
		"в файл писали после обхода (время изменения другое)": "the file was written to after the walk (its modification time differs)",
		"права изменились (были %v, стали %v)":                "the permissions changed (were %v, are %v)",

		// ── почему запись попала в план ──────────────────────────────
		"разряд %s: возраст %s ≥ порога %s": "class %s: age %s ≥ threshold %s",
		"разряд %s, приговор %s":            "class %s, verdict %s",

		// ── сторож: почему приговор ядра всё-таки не исполнен ────────
		"приговор %s — убирается только «%s»":                                                    "verdict %s — only «%s» is removed",
		"файл не удалось прочитать при обходе — недоступное не трогаем":                          "the file could not be read during the walk — the inaccessible is left alone",
		"это каталог: снос каталога рекурсивен, и ядро его «МожноУбрать» не выдаёт (правило П3)": "this is a directory: removing one is recursive, and the core gives it no «МожноУбрать» (rule П3)",
		"это символическая ссылка: ядро её не убирает (правило П2)":                              "this is a symbolic link: the core does not remove one (rule П2)",
		"это %s, а не обычный файл":                                                              "this is a %s, not a regular file",
		"не обычный файл (%v)":                                                "not a regular file (%v)",
		"путь не считается от корня: %v":                                      "the path does not count from the root: %v",
		"путь вне указанного корня":                                           "the path is outside the given root",
		"путь внутри корзины digitdisk — своё же убирает `purge`, не `clean`": "the path is inside digitdisk's trash — its own is erased by `purge`, not by `clean`",

		// ── корзина и корень ─────────────────────────────────────────
		`корзина %s лежит вне корня %s.
Корзина обязана быть внутри корня: тогда перенос — это rename(2), то есть
мгновенно и без копирования, и все обращения идут через os.Root, который
из корня не выпускает даже по символической ссылке. Корзина на другой
файловой системе превратила бы перенос в копирование: цена обратимости
стала бы равна объёму уборки, а обрыв на середине оставил бы полфайла`: `the trash %s lies outside the root %s.
The trash must be inside the root: the move is then rename(2), that is
instant and without copying, and every access goes through os.Root, which
does not leave the root even by a symbolic link. A trash on another
filesystem would turn the move into copying: the price of reversibility
would become the size of the cleanup, and a break midway would leave half a file`,
		"%s лежит вне корня %s":        "%s lies outside the root %s",
		"корень %s не открывается: %s": "the root %s does not open: %s",

		// ── перенос в корзину ────────────────────────────────────────
		`решающий слой — %s.
Он не выносит приговора «%s» ни одной записи, поэтому убирать нечего и не по чему.
Собери хозяина с признаком flangcore: go build -tags flangcore -o digitdisk ./host`: `the decision layer is %s.
It gives the verdict «%s» to no record, so there is nothing to clean and nothing to clean by.
Build the host with the flangcore tag: go build -tags flangcore -o digitdisk ./host`,
		"убирать нечего: ядро не пометило «%s» ни одного файла под %s": "nothing to clean: the core marked «%s» on no file under %s",
		"корзина %s не создаётся: %s":                                  "the trash %s cannot be created: %s",
		"журнал %s не записывается: %s":                                "the journal %s cannot be written: %s",
		`файлы перенесены, но журнал %s не переписан: %s.
Первая запись журнала на месте и возврат по ней работает`: `the files were moved, but the journal %s was not rewritten: %s.
The first journal entry is in place and restoring by it works`,
		"исчез между обходом и переносом":                          "vanished between the walk and the move",
		"не читается: %v":                                          "does not read: %v",
		"перестал быть обычным файлом (стал %v)":                   "is no longer a regular file (became %v)",
		"%s — не убран":                                            "%s — not cleaned",
		"каталог в корзине не создаётся: %v":                       "the directory in the trash cannot be created: %v",
		"в корзине уже есть %s — перезаписывать не будем":          "the trash already holds %s — we will not overwrite it",
		"корзина не проверяется: %v":                               "the trash cannot be checked: %v",
		"не переносится: %v":                                       "does not move: %v",
		"перенесён, но в корзине не находится: %v":                 "moved, but is not found in the trash: %v",
		"в корзине оказался не тот файл (узел %d:%d вместо %d:%d)": "the trash holds the wrong file (node %d:%d instead of %d:%d)",

		// ── стирание насовсем: тот же план, но без корзины ───────────
		"стирать нечего: ядро не пометило «%s» ни одного файла под %s": "nothing to erase: the core marked «%s» on no file under %s",
		`файлы стёрты, но журнал %s не переписан: %s.
Первая запись журнала на месте, и она называет всё, что стиралось`: `the files were erased, but the journal %s was not rewritten: %s.
The first journal entry is in place, and it names everything that was being erased`,
		"исчез между обходом и стиранием": "vanished between the walk and the erasure",
		"%s — не стёрт":                   "%s — not erased",
		"%s — журнал стирания, а не корзины: файлы стёрты насовсем, возвращать нечего": "%s is the journal of an erasure, not of a trash: the files are gone for good, there is nothing to put back",
		"%s — журнал стирания, а не корзины: эти файлы уже стёрты, стирать нечего":     "%s is the journal of an erasure, not of a trash: these files are already gone, there is nothing to erase",

		// ── журнал корзины ───────────────────────────────────────────
		"в %s нет %s — это не корзина digitdisk": "%s holds no %s — this is not a digitdisk trash",
		"%s не читается как журнал: %s":          "%s does not read as a journal: %s",
		"%s: версия журнала %d, а этот digitdisk понимает %d — работать с непонятым журналом опаснее, чем отказаться": "%s: the journal is version %d, and this digitdisk understands %d — working with a journal it has not understood is more dangerous than refusing",
		`корзина лежит в %s, а журнал записан для %s.
Возврат кладёт файлы по абсолютным путям, записанным при переносе;
из перемещённой копии он писал бы в исходное дерево, а не туда, где копия`: `the trash lies in %s, and the journal was written for %s.
Restoring puts the files back by the absolute paths recorded at the move;
from a moved copy it would write into the original tree, not where the copy is`,
		"%s: в журнале нет корня": "%s: the journal names no root",
		"корзина %s вне корня %s": "the trash %s is outside the root %s",
		"%s — не каталог":         "%s is not a directory",

		// ── стирание ─────────────────────────────────────────────────
		"--confirm не может быть отрицательным": "--confirm cannot be negative",
		"в корзине %d файлов (%d Б), а --confirm назвал %d — ничего не стёрто.\nЗапусти `digitdisk purge %s` без ключа: он напечатает, что будет стёрто, и число для --confirm": "the trash holds %d files (%d B), and --confirm named %d — nothing was erased.\nRun `digitdisk purge %s` with no key: it prints what would be erased, and the number for --confirm",
		"корзина %s пуста — стирать нечего":                     "the trash %s is empty — nothing to erase",
		"стёрто %d файлов, но журнал %s не переписан: %s":       "%d files erased, but the journal %s was not rewritten: %s",
		"в корзине его уже нет":                                 "it is no longer in the trash",
		"в корзине не читается: %v":                             "does not read in the trash: %v",
		"в корзине это не обычный файл (%v) — стирать не будем": "in the trash this is not a regular file (%v) — we will not erase it",
		"в корзине его правили: %s — не стёрт":                  "it was edited in the trash: %s — not erased",
		"не стирается: %v":                                      "does not erase: %v",

		// ── возврат ──────────────────────────────────────────────────
		"не переносился — возвращать нечего":                             "was never moved — nothing to put back",
		"уже возвращён %s":                                               "already put back %s",
		"стёрт %s — возврат невозможен":                                  "erased %s — putting it back is impossible",
		"файлы возвращены, но журнал %s не переписан: %s":                "the files were put back, but the journal %s was not rewritten: %s",
		"в корзине его нет: перенос не дошёл или корзину правили руками": "it is not in the trash: the move never finished or the trash was edited by hand",
		"в корзине это уже не обычный файл (%v)":                         "in the trash this is no longer a regular file (%v)",
		"в корзине его правили: %s — не возвращён":                       "it was edited in the trash: %s — not put back",
		"на прежнем месте уже что-то есть — перезаписывать не будем":     "something is at the former place already — we will not overwrite it",
		"прежнее место не проверяется: %v":                               "the former place cannot be checked: %v",
		"прежний каталог не создаётся: %v":                               "the former directory cannot be created: %v",
		"не возвращается: %v":                                            "does not come back: %v",

		// ── справочник известных мест ────────────────────────────────
		"домашний каталог не определяется, а справочник считает места от него: %s": "the home directory cannot be determined, and the known-places directory counts places from it: %s",
		"%s, строка %d: %s": "%s, line %d: %s",
		"полей %d, а надо 7: разряд|якорь|система|путь|переменная|имя|источник":                                  "%d fields, and 7 are needed: class|anchor|system|path|variable|name|source",
		"разряд %q справочнику не позволен; можно: кэш, журнал, сборка, загрузка":                                "the class %q is not allowed in the directory; you may write: кэш, журнал, сборка, загрузка",
		"система %q: можно linux, macos, все":                                                                    "the system %q: you may write linux, macos, все",
		"у места нет имени: непонятно, чей это каталог":                                                          "the place has no name: whose directory this is stays unclear",
		"у места %q нет источника; в этом справочнике строка без ссылки на документацию инструмента — не строка": "the place %q has no source; in this directory a row without a link to the tool's documentation is not a row",
		"путь пуст": "the path is empty",
		"путь начинается с косой, а якорь %q — это можно только с якорем «корень»":                "the path starts with a slash while the anchor is %q — that is allowed only with the anchor «корень»",
		"переменная %s задана как %q — не абсолютный путь; место пропущено бы молча, а это отказ": "the variable %s is set to %q — not an absolute path; the place would be skipped silently, and that is a refusal",
		"у якоря «всюду» переменной быть не может: место не привязано ни к какому основанию":      "the anchor «всюду» can have no variable: the place is tied to no base at all",
		"якорь %q: можно дом, кэш, данные, настройки, корень, всюду":                              "the anchor %q: you may write дом, кэш, данные, настройки, корень, всюду",
		"другая система":                      "another system",
		"на любой глубине — одного места нет": "at any depth — there is no single place",
		"нет":        "none",
		"не каталог": "not a directory",

		// ── защитный список ──────────────────────────────────────────
		"%s, строка %d: полей %d, а надо не меньше двух: вид|значение[|почему]": "%s, line %d: %d fields, and at least two are needed: kind|value[|why]",
		"разряд %q не известен; есть: %s":                                       "the class %q is unknown; there are: %s",
		"путь начинается с ~, а домашний каталог не определяется":               "the path starts with ~, and the home directory cannot be determined",
		"вид %q: можно «путь» или «разряд»":                                     "the kind %q: you may write «путь» or «разряд»",

		// ── мост к решающему слою ────────────────────────────────────
		"справочник: разряд %q решающему слою не известен": "known places: the class %q is unknown to the decision layer",
		"справочник: якорь %q решающему слою не известен":  "known places: the anchor %q is unknown to the decision layer",
		"решающий слой не смог проверить справочник: %s":   "the decision layer could not check the known places: %s",
		"решающий слой отверг справочник: у какого-то места цепь не ограничена косыми с обеих сторон, а без них сверка перестала бы быть сверкой по составляющим": "the decision layer rejected the known places: some place's chain is not bounded by slashes on both sides, and without them the check would stop being a check by components",

		// ── неизмеренное: linux ──────────────────────────────────────
		"в %s датчиков нет": "there are no sensors in %s",

		// ── неизмеренное: macOS ──────────────────────────────────────
		"macOS не публикует показания датчиков, а угадывать их формат нельзя":                                "macOS does not publish sensor readings, and guessing their format is not allowed",
		"ответ не читается как момент времени — время работы не считаем":                                     "the answer does not read as an instant — uptime is not counted",
		"ответ не читается как средние загрузки — не показываем":                                             "the answer does not read as load averages — not shown",
		"macOS не публикует длину очереди планировщика":                                                      "macOS does not publish the length of the scheduler queue",
		"ответ не восьмибайтовый — объём памяти не показываем":                                               "the answer is not eight bytes — the memory size is not shown",
		"ответ не читается как сведения о свопе — не показываем":                                             "the answer does not read as swap information — not shown",
		"без размера страницы разбивку памяти не пересчитать в байты":                                        "without the page size the memory breakdown cannot be counted in bytes",
		"без объёма памяти разбивку не с чем сверить":                                                        "without the memory size there is nothing to check the breakdown against",
		"ядро не дало разбивку памяти: %s":                                                                   "the kernel gave no memory breakdown: %s",
		"разбивка памяти не сошлась с объёмом памяти машины — не публикуем":                                  "the memory breakdown did not agree with the machine's memory size — not published",
		"окно замера нулевое — доля занятого процессора не измерялась":                                       "the measurement window is zero — the busy share of the processor was not measured",
		"окно замера нулевое — процессорное время процессов не измерялось":                                   "the measurement window is zero — the processes' CPU time was not measured",
		"ядро не дало счётчики процессорного времени":                                                        "the kernel gave no CPU time counters",
		"за окно замера счётчики процессора не сдвинулись":                                                   "over the measurement window the processor counters did not move",
		"ответ не делится на записи процессов — список не публикуем":                                         "the answer does not divide into process records — the list is not published",
		"самопроверка записи о процессе не сошлась — числа из неё не публикуем":                              "the self-check of a process record did not agree — its numbers are not published",
		"самопроверка памяти процессов не сошлась — их память и потоки не публикуем":                         "the self-check of process memory did not agree — their memory and threads are not published",
		"сколько процессов работает прямо сейчас, видно только по их потокам":                                "how many processes are running right now is visible only through their threads",
		"macOS не различает заблокированные процессы среди спящих":                                           "macOS does not tell blocked processes apart from sleeping ones",
		"память, потоки и командные строки чужих процессов видны только администратору — запустите под sudo": "the memory, threads and command lines of other users' processes are visible only to the administrator — run under sudo",
		"без памяти процессов их процессорное время тоже не публикуем":                                       "without the processes' memory their CPU time is not published either",
		"ни один процесс не прожил всё окно замера":                                                          "no process lived through the whole measurement window",
		"ядро не назвало предельный размер блока аргументов":                                                 "the kernel did not name the limit on the argument block size",
		"самопроверка командной строки не сошлась — чужих командных строк не публикуем":                      "the self-check of a command line did not agree — other users' command lines are not published",
		"список интерфейсов не прочитался: %s":                                                               "the list of interfaces did not read: %s",
		"macOS не считает отброшенные исходящие пакеты":                                                      "macOS does not count dropped outgoing packets",
		"самопроверка счётчиков интерфейсов не сошлась — не публикуем":                                       "the self-check of interface counters did not agree — not published",
	})
}
