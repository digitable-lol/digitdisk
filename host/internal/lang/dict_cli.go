// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

// Английские половины того, что говорит сам инструмент о себе: справка, версия,
// перечень подкоманд и ключей, отказы разбора командной строки, предупреждение
// сборки без ядра и всё, что печатает выбор языка.
//
// Ключ — русская строка ровно так, как она написана в исходнике, включая
// пробелы выравнивания: справка печатается колонками, и колонка, разъехавшаяся
// на втором языке, — это та же полуработа, что непереведённая строка.
func init() {
	add(map[string]string{
		// ── справка: шапка и подкоманды ──────────────────────────────
		"digitdisk — снимок системы, разбор дерева каталогов и уборка.": "digitdisk — a system snapshot, a directory-tree walk and a cleanup.",
		"digitdisk [подкоманда] [ключи]      без подкоманды — %s":       "digitdisk [subcommand] [keys]       without a subcommand — %s",
		"Подкоманды:": "Subcommands:",
		"--help, --version  эта справка; версия, сборка и решающий слой": "--help, --version  this help; version, build and decision layer",

		// доводы подкоманд: то, что человек подставляет своим
		"<путь>":    "<path>",
		"<корзина>": "<trash>",
		"<команда>": "<command>",

		// однострочные пояснения подкоманд (cli.Commands)
		"снимок системы: ЦП, память, диски, сеть; в терминале — живой экран": "a system snapshot: CPU, memory, disks, network; live screen in a terminal",
		"обход дерева: каталоги по размеру и самые крупные файлы":            "a tree walk: directories by size and the largest files",
		"план уборки; переносит в корзину только с --apply":                  "a cleanup plan; moves to the trash only with --apply",
		"вернуть корзину на прежние места":                                   "put a trash back where it came from",
		"стереть корзину: необратимо, требует --confirm N":                   "erase a trash: irreversible, requires --confirm N",
		"справочник известных мест и что из него есть на этой машине":        "the known-places directory and what of it is on this machine",
		"чем кончались прошлые уборки под этим корнем":                       "how past cleanups under this root ended",
		"запустить команду и показать, во что она обошлась":                  "run a command and show what it cost",

		// ── справка: ключи, по строке на каждый ──────────────────────
		"Ключи:": "Keys:",
		"  -c <команда>      короткий вид run: всё после -c принадлежит команде,":         "  -c <command>      the short spelling of run: all after -c is the command's,",
		"                    свои ключи ставятся до неё; одна строка с пробелами":         "                    ours go before it; one word holding spaces and shell",
		"                    и метасимволами исполняется оболочкой ($SHELL)":              "                    metacharacters is handed to the shell ($SHELL)",
		"  --json            машиночитаемый вывод; принимают все подкоманды":              "  --json            machine-readable output; every subcommand takes it",
		"  --lang ЯЗЫК       ru или en на этот запуск; принимают все подкоманды":          "  --lang LANG       ru or en for this run; every subcommand takes it",
		"  --top N           строк в списках: status 10, analyze и clean 15,":             "  --top N           rows in the lists: status 10, analyze and clean 15,",
		"                    places 40, history 20; 0 — без предела":                      "                    places 40, history 20; 0 — no limit",
		"  --why             status: что не измерено и почему":                            "  --why             status: what was not measured, and why",
		"  --sample MS       status: окно замера загрузки ЦП, по умолчанию 200":           "  --sample MS       status: CPU-busy sampling window, 200 by default",
		"  --live            status, analyze: живой экран; без терминала — ошибка":        "  --live            status, analyze: the live screen; error without a terminal",
		"  --plain           status, analyze: печать без экрана; run: без строки":         "  --plain           status, analyze: print without the screen; run: no",
		"                    состояния — и то и другое даже в терминале":                  "                    status line — either way, even in a terminal",
		"  --interval MS     status: период обновления живого экрана, 2000;":              "  --interval MS     status: live-screen refresh period, 2000;",
		"                    run: период обновления строки состояния, 1000":               "                    run: status-line refresh period, 1000",
		"  --cross-device    analyze, clean: заходить на другие файловые системы":         "  --cross-device    analyze, clean: cross into other filesystems",
		"  --max-depth N     analyze, clean: предел глубины обхода; 0 — без предела":      "  --max-depth N     analyze, clean: walk depth limit; 0 — no limit",
		"  --places ФАЙЛ     analyze, clean, places: свой справочник известных мест":      "  --places FILE     analyze, clean, places: your own known-places file",
		"  --no-places       analyze, clean: судить одними приметами, без справочника":    "  --no-places       analyze, clean: judge by signs alone, no directory",
		"  --no-fold         analyze: судить о каждом файле внутри node_modules и":        "  --no-fold         analyze: judge every file inside node_modules and",
		"                    подобных каталогов; по умолчанию такой каталог":              "                    the like; by default such a directory is counted",
		"                    считается целиком, а приговор выносится один — о нём":        "                    whole and gets a single verdict, about itself",
		"  --no-measure      places: не считать размеры, только назвать места":            "  --no-measure      places: do not measure sizes, only name the places",
		"  --apply           clean: перенести в корзину, а не только показать план":       "  --apply           clean: move into the trash, not just show the plan",
		"  --trash КАТ       clean: другая корзина; обязана лежать внутри корня":          "  --trash DIR       clean: another trash; must lie inside the root",
		"  --protect ЧТО     clean, analyze: не трогать путь или «разряд:кэш»; повторяем": "  --protect WHAT    clean, analyze: spare a path or «разряд:кэш»; may repeat",
		"  --protect-file Ф  clean, analyze: защитный список файлом":                      "  --protect-file F  clean, analyze: the protect list as a file",
		"  --dry-run         restore: показать, что вернулось бы, и не возвращать":        "  --dry-run         restore: show what would come back, restore nothing",
		"  --confirm N       purge: подтвердить стирание ровно N файлов":                  "  --confirm N       purge: confirm erasing exactly N files",

		"Экран status: раздел КОМАНДЫ не только называет — ↑ ↓ и 1…8 выбирают,":       "status screen: the COMMANDS section does more than name — ↑ ↓ and 1…8 choose,",
		"  Enter запускает; analyze и clean спросят путь, предложив текущий каталог.": "  Enter runs; analyze and clean ask for a path, offering the current directory.",
		"  Разделы: ← → (их одиннадцать), 1…9 сразу, ↑ ↓ PgUp/PgDn прокрутка,":        "  Sections: ← → (eleven of them), 1…9 at once, ↑ ↓ PgUp/PgDn scroll,",
		"  p пауза, r замер, l язык, ? сразу к КОМАНДАМ, q выход.":                    "  p pause, r sample, l language, ? straight to COMMANDS, q quit.",

		// Живой экран запустил подкоманду и ждёт, пока прочитают её вывод.
		"— Enter возвращает на экран состояния": "— Enter goes back to the status screen",

		// Куда экран посылает за подкомандой, которую сам не запускает
		// (cli.Command.Instead). Тоже данные: их печатает раздел КОМАНДЫ.
		"Enter на корзине в разделе ЖУРНАЛ экрана analyze":             "Enter on a trash in the JOURNAL section of the analyze screen",
		"из оболочки — это единственный необратимый шаг":               "from the shell only — the one irreversible step",
		"раздел ЖУРНАЛ экрана analyze, по обойденному корню":           "the JOURNAL section of the analyze screen, for the walked root",
		"из оболочки: у чужой команды свои ключи, экран их не наберёт": "from the shell: somebody else's command has its own flags, and the screen cannot type them",

		// Обёртка: что показывает строка состояния и чего она не трогает.
		"Строка состояния run: последняя строка терминала, вывод команды её не":      "run status line: the last row of the terminal; the command's output never",
		"  задевает; в трубу и в файл её нет вовсе. Полноэкранная программа (vim,":   "  touches it, and into a pipe or a file there is none at all. A full-screen",
		"  ssh, less) забирает терминал — строка уходит и возвращается сама.":        "  program (vim, ssh, less) takes the terminal — the line goes and comes back.",
		"  Код возврата и сигнал — команды, не обёртки. Сводка идёт в поток ошибок.": "  The exit code and signal are the command's own. The summary goes to stderr.",

		// Экран обхода. Он длиннее строкой, потому что делает больше: по
		// найденному ходят и из него убирают, и обе эти клавиши надо назвать
		// там же, где называются остальные.
		"Экран analyze: пока идёт обход — растущие числа и чем наполняется дерево;":       "analyze screen: while the walk runs — climbing numbers and what fills the tree;",
		"  перечень помечен ПРЕДВАРИТЕЛЬНО, пока он догадка, q прерывает обход.":          "  the list is marked ПРЕДВАРИТЕЛЬНО while it is a guess, q stops the walk.",
		"  После обхода: Tab и 1…8 разделы, ↑ ↓ строка, → внутрь каталога, ← назад,":      "  After it: Tab and 1…8 sections, ↑ ↓ row, → into a directory, ← back out,",
		"  Пробел отметить каталог, «.» текущий, c план уборки и подтверждение числом,":   "  Space marks a directory, «.» this one, c plans the cleaning and asks a count,",
		"  o обойти другой каталог, Enter в ЖУРНАЛЕ вернуть корзину, l язык, ? клавиши.":  "  o walks another directory, Enter in JOURNAL puts a trash back, l lang, ? keys.",
		"  Убирается ровно то, что убрала бы подкоманда clean, и так же — через план.":    "  Exactly what the clean subcommand would remove goes, and the same way: by plan.",
		"Палитра: DIGITDISK_PALETTE=carbon|paper|signal, NO_COLOR и TERM=dumb уважаются.": "Palette: DIGITDISK_PALETTE=carbon|paper|signal; NO_COLOR and TERM=dumb honoured.",

		"Язык: спрашивается один раз и помнится в ~/.digitable/digitdisk/settings.conf;": "Language: asked once and kept in ~/.digitable/digitdisk/settings.conf;",
		"  DIGITDISK_LANG=ru|en и --lang перекрывают его, --json не переводится.":        "  DIGITDISK_LANG=ru|en and --lang override it; --json is not translated.",

		"Подробно: man digitdisk": "In full: man digitdisk",

		// ── пояснения ключей, которые печатает сам flag ──────────────
		"язык вывода на этот запуск: ru или en":                                         "output language for this run: ru or en",
		"машиночитаемый вывод":                                                          "machine-readable output",
		"что не измерено и почему":                                                      "what was not measured, and why",
		"сколько процессов в каждом списке":                                             "how many processes in each list",
		"окно замера загрузки ЦП, мс":                                                   "CPU-busy sampling window, ms",
		"живой экран, даже если о терминале не спрашивали":                              "the live screen, even where no terminal was asked about",
		"печать одним снимком, без живого экрана":                                       "print one snapshot, no live screen",
		"период обновления живого экрана, мс":                                           "live-screen refresh period, ms",
		"сколько строк в списках":                                                       "how many rows in the lists",
		"заходить на другие файловые системы":                                           "cross into other filesystems",
		"предел глубины обхода, 0 — без предела":                                        "walk depth limit, 0 — no limit",
		"свой справочник известных мест":                                                "your own known-places directory",
		"судить одними приметами, без справочника":                                      "judge by signs alone, without the directory",
		"судить о каждом файле внутри node_modules и подобных, а не о каталоге целиком": "judge every file inside node_modules and the like, not the directory whole",
		"перенести в корзину, а не только показать план":                                "move into the trash, not just show the plan",
		"сколько строк в перечнях, 0 — без предела":                                     "how many rows in the lists, 0 — no limit",
		"корзина (по умолчанию <корень>/%s); обязана лежать внутри корня":               "the trash (by default <root>/%s); must lie inside the root",
		"защитный список файлом":                                                        "the protect list as a file",
		"не трогать: путь или «разряд:кэш»; можно повторять":                            "spare: a path or «разряд:кэш»; may repeat",
		"показать, что вернулось бы, и не возвращать":                                   "show what would come back, restore nothing",
		"подтвердить стирание ровно N файлов":                                           "confirm erasing exactly N files",
		"сколько найденных мест печатать, 0 — без предела":                              "how many found places to print, 0 — no limit",
		"не считать размеры, только назвать места":                                      "do not measure sizes, only name the places",
		"сколько корзин печатать, 0 — без предела":                                      "how many trashes to print, 0 — no limit",

		// ── отказы разбора командной строки ──────────────────────────
		"неизвестная подкоманда %q":                                                             "unknown subcommand %q",
		"подкоманда %q объявлена в internal/cli, но не разобрана в main":                        "subcommand %q is declared in internal/cli and not dispatched in main",
		"подкоманда status не принимает путей (лишнее: %q)":                                     "the status subcommand takes no paths (spare: %q)",
		"подкоманда places не принимает путей (лишнее: %q)":                                     "the places subcommand takes no paths (spare: %q)",
		"нужен ровно один путь для обхода, получено %d":                                         "exactly one path to walk is needed, got %d",
		"нужен ровно один путь для уборки, получено %d":                                         "exactly one path to clean is needed, got %d",
		"нужна ровно одна корзина, получено %d":                                                 "exactly one trash is needed, got %d",
		"нужен ровно один путь — корень уборки, хранилище корзин или одна корзина; получено %d": "exactly one path is needed — a cleanup root, a trash store or one trash; got %d",
		"%s; без --live тот же обход печатается текстом":                                        "%s; without --live the same walk is printed as text",
		"обойти молча и напечатать отчёт, без экрана":                                           "walk silently and print the report, with no screen",
		"%s; без --live тот же снимок печатается текстом":                                       "%s; without --live the same snapshot is printed as text",
		"справочник %s: %s": "known places %s: %s",
		"решающий слой %q справочника не принимает — %d мест не применено": "the decision layer %q takes no directory — %d places were not applied",

		// ── версия ───────────────────────────────────────────────────
		"сборка          %s, %s":                "build           %s, %s",
		"инструментарий  %s %s/%s":              "toolchain       %s %s/%s",
		"решающий слой   %s, договор версии %d": "decision layer  %s, contract version %d",
		"язык            %s (%s)":               "language        %s (%s)",
		"неизвестен":                            "unknown",
		"неизвестно":                            "unknown",
		"(дерево с правками)":                   "(tree with edits)",
		"— собрано без признака flangcore":      "— built without the flangcore tag",

		// ── сборка без решающего ядра ────────────────────────────────
		"ВНИМАНИЕ: digitdisk собран БЕЗ решающего ядра.":                        "WARNING: digitdisk is built WITHOUT the decision core.",
		"Разбора не будет: каждая запись вернётся как «Неизвестное/НеТрогать»,": "There will be no analysis: every record comes back «Неизвестное/НеТрогать»,",
		"и пустой список «предложено убрать» НЕ означает, что убирать нечего.":  "and an empty «proposed for removal» does NOT mean there is nothing to remove.",
		"Соберите с ядром:  go build -tags flangcore -o digitdisk ./host":       "Build with the core:  go build -tags flangcore -o digitdisk ./host",
		"Или поставьте выпуск: brew install digitable-lol/tap/digitdisk":        "Or install a release: brew install digitable-lol/tap/digitdisk",
		"(выпуск всегда идёт с ядром — scripts/build-release.sh)":               "(a release always ships with the core — scripts/build-release.sh)",

		// ── выбор языка и настройки ──────────────────────────────────
		"язык сохранён: %s":     "language stored: %s",
		"язык не сохранён (%s)": "language not stored (%s)",
		"язык не сохранён (%s) — в этот раз вывод по-выбранному, а вопрос повторится":     "language not stored (%s) — this run is in the language you chose, and the question will come back",
		"язык %q не из двух (ru, en) — вывод на языке машины":                             "language %q is not one of the two (ru, en) — output in the machine's language",
		"DIGITDISK_LANG=%q — не из двух (ru, en), вывод на языке машины":                  "DIGITDISK_LANG=%q — not one of the two (ru, en), output in the machine's language",
		"настройки не прочитаны (%s) — вывод на языке машины":                             "settings unread (%s) — output in the machine's language",
		"настройки переехали в %s — %s ещё читается, перенесите его, когда удобно":        "settings have moved to %s — %s is still read, move it when it suits you",
		"# Настройки digitdisk. Правится руками: ключ=значение, одна строка — один ключ.": "# digitdisk settings.  Edit by hand: key=value, one key per line.",
		"# Этот файл завёл сам digitdisk, когда спросил про язык. Убрать его можно.":      "# digitdisk made this file when it asked about the language.  You may delete it.",

		// ── единицы времени ──────────────────────────────────────────
		"д":      "d",
		"дн":     "d",
		"%d с":   "%d s",
		"%d мин": "%d min",
		"%d ч":   "%d h",
		"%d дн":  "%d d",
		"%d мс":  "%d ms",
		"%s с":   "%s s",
	})
}

// Английские половины отказов пакета настроек. Они здесь, а не в dict_clean.go,
// потому что настройки — часть разговора инструмента о самом себе, как справка
// и версия.
func init() {
	add(map[string]string{
		"%s, строка %d: ждалось «ключ=значение», а написано %q": "%s, line %d: expected «key=value», got %q",
		"%s, строка %d: язык %q не из двух (ru, en)":            "%s, line %d: language %q is not one of the two (ru, en)",
	})
}

func init() {
	add(map[string]string{
		"сказано один раз, отметка в %s": "said once; the mark is in %s",
	})
}
