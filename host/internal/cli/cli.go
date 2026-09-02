// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package cli holds the one list of digitdisk's subcommands and the text of
// its справка.
//
// It exists because the list had started to live in three places at once —
// the помощь, the live screen and the страница руководства — and three copies
// of a list are three chances to disagree about what the tool can do.  Here it
// is data: main dispatches on it, the screen shows it under «?», and
// scripts/check-docs.sh reads it to make sure digitdisk.1 names the same
// commands and the same keys the code registers.
package cli

import (
	"fmt"
	"strings"

	"digitdisk/internal/lang"
)

// A Command is one subcommand: its name, the argument it needs, one line
// saying what it does, and what the live screen does when this line is chosen
// from its list.  One line of gloss, because this text has to fit a справка, a
// section of the screen and a table in the man page without being rewritten
// for each.
type Command struct {
	Name  string
	Arg   string
	Gloss string
	// Start is what the screen does with this command, and Instead is where
	// the command lives when the screen does not start it.
	//
	// They stand HERE, beside the command, and not in the screen, for the
	// reason the whole file exists: «что будет, если нажать Enter на этой
	// строке» — свойство подкоманды, а не рисунка.  Держи его в ui, и через
	// месяц список экрана разойдётся со списком справки на одну строку,
	// которую никто не заметит.
	Start   Start
	Instead string
}

// Start says what the live screen may do with a subcommand chosen from its
// list of commands.
type Start int

const (
	// StartElsewhere — экран эту подкоманду не запускает, и Instead
	// говорит, где она живёт: разделом другого экрана или в оболочке.
	StartElsewhere Start = iota
	// StartHere — экран делает это, никуда не уходя: `status` и есть этот
	// экран, а запустить его отсюда значит взять новый замер.
	StartHere
	// StartRun — экран закрывается, подкоманда работает и печатает, экран
	// возвращается.
	StartRun
	// StartPath — то же, что StartRun, но сперва спрашивается путь.
	StartPath
)

// Commands is the list.  The order is the order everything shows it: reading
// first, then the three steps of removal, then the two that look at what has
// been decided and done.
// Столбцы: имя, довод, одна строка о деле, что делает экран, куда он посылает.
var Commands = []Command{
	{"status", "", "снимок системы: ЦП, память, диски, сеть; в терминале — живой экран",
		StartHere, ""},
	{"analyze", "<путь>", "обход дерева: каталоги по размеру и самые крупные файлы",
		StartPath, ""},
	{"clean", "<путь>", "план уборки; переносит в корзину только с --apply",
		StartPath, ""},
	{"restore", "<корзина>", "вернуть корзину на прежние места",
		StartElsewhere, "Enter на корзине в разделе ЖУРНАЛ экрана analyze"},
	{"purge", "<корзина>", "стереть корзину: необратимо, требует --confirm N",
		StartElsewhere, "из оболочки — это единственный необратимый шаг"},
	{"places", "", "справочник известных мест и что из него есть на этой машине",
		StartRun, ""},
	{"history", "<путь>", "чем кончались прошлые уборки под этим корнем",
		StartElsewhere, "раздел ЖУРНАЛ экрана analyze, по обойденному корню"},
}

// Default is the subcommand a bare `digitdisk` runs.  Reading is the frequent
// thing and it changes nothing, so it is the one that may happen without being
// asked for by name.
const Default = "status"

// HelpArgs and VersionArgs are what may stand where a subcommand stands and
// mean something other than «сделай умолчание».  They are lists rather than a
// switch in main for the same reason Commands is: digitdisk.1 promises exactly
// these spellings, and a test compares the two.
var (
	HelpArgs    = []string{"-h", "--help", "help"}
	VersionArgs = []string{"-V", "--version", "version"}
)

// Is reports whether arg is one of list.
func Is(list []string, arg string) bool {
	for _, s := range list {
		if s == arg {
			return true
		}
	}
	return false
}

// Known reports whether name is a subcommand.
func Known(name string) bool {
	for _, c := range Commands {
		if c.Name == name {
			return true
		}
	}
	return false
}

// Call is "clean   <путь>", padded so the glosses line up.  The name of a
// подкоманда is not translated — it is what a person types — and the довод
// beside it is, because «<путь>» is not something anybody types.
func (c Command) Call(l lang.Lang) string {
	return fmt.Sprintf("%-7s %s", c.Name, l.T(c.Arg))
}

// Usage is the text of --help: what to type, and nothing about why.
func Usage(l lang.Lang) string {
	var b strings.Builder
	b.WriteString(l.T("digitdisk — снимок системы, разбор дерева каталогов и уборка.") + "\n\n")
	b.WriteString("  " + l.F("digitdisk [подкоманда] [ключи]      без подкоманды — %s", Default) + "\n\n")
	b.WriteString(l.T("Подкоманды:") + "\n")
	for _, c := range Commands {
		b.WriteString(fmt.Sprintf("  %-18s %s\n", strings.TrimRight(c.Call(l), " "), l.T(c.Gloss)))
	}
	b.WriteString("  " + l.T("--help, --version  эта справка; версия, сборка и решающий слой") + "\n")
	for _, line := range keys {
		b.WriteString(l.T(line) + "\n")
	}
	return b.String()
}

// keys lists every flag the subcommands register, one line each, naming the
// subcommands that take it.  Nothing here explains itself: a key, what it
// takes, where it works, what it does.
//
// It is a list of lines rather than one block of text because every line is
// translated on its own: an English line wraps where English wraps, and a
// block would have to break in the same places in both languages or be
// re-flowed by hand at every edit.
var keys = []string{
	"",
	"Ключи:",
	"  --json            машиночитаемый вывод; принимают все подкоманды",
	"  --lang ЯЗЫК       ru или en на этот запуск; принимают все подкоманды",
	"  --top N           строк в списках: status 10, analyze и clean 15,",
	"                    places 40, history 20; 0 — без предела",
	"  --why             status: что не измерено и почему",
	"  --sample MS       status: окно замера загрузки ЦП, по умолчанию 200",
	"  --gpu-tool        status: спросить о видеокартах чужую программу",
	"                    (nvidia-smi) — то, чего драйвер не публикует файлами",
	"  --live            status, analyze: живой экран; без терминала — ошибка",
	"  --plain           status, analyze: печать без экрана, даже в терминале",
	"  --interval MS     status: период обновления живого экрана, 2000",
	"  --cross-device    analyze, clean: заходить на другие файловые системы",
	"  --max-depth N     analyze, clean: предел глубины обхода; 0 — без предела",
	"  --places ФАЙЛ     analyze, clean, places: свой справочник известных мест",
	"  --no-places       analyze, clean: судить одними приметами, без справочника",
	"  --no-measure      places: не считать размеры, только назвать места",
	"  --apply           clean: перенести в корзину, а не только показать план",
	"  --trash КАТ       clean: другая корзина; обязана лежать внутри корня",
	"  --protect ЧТО     clean, analyze: не трогать путь или «разряд:кэш»; повторяем",
	"  --protect-file Ф  clean, analyze: защитный список файлом",
	"  --dry-run         restore: показать, что вернулось бы, и не возвращать",
	"  --confirm N       purge: подтвердить стирание ровно N файлов",
	"",
	"Экран status: раздел КОМАНДЫ не только называет — ↑ ↓ и 1…7 выбирают,",
	"  Enter запускает; analyze и clean спросят путь, предложив текущий каталог.",
	"  Разделы: ← → (их одиннадцать), 1…9 сразу, ↑ ↓ PgUp/PgDn прокрутка,",
	"  p пауза, r замер, l язык, ? сразу к КОМАНДАМ, q выход.",
	"",
	"Экран analyze: пока идёт обход — растущие числа и чем наполняется дерево;",
	"  перечень помечен ПРЕДВАРИТЕЛЬНО, пока он догадка, q прерывает обход.",
	"  После обхода: Tab и 1…8 разделы, ↑ ↓ строка, → внутрь каталога, ← назад,",
	"  Пробел отметить каталог, «.» текущий, c план уборки и подтверждение числом,",
	"  o обойти другой каталог, Enter в ЖУРНАЛЕ вернуть корзину, l язык, ? клавиши.",
	"  Убирается ровно то, что убрала бы подкоманда clean, и так же — через план.",
	"",
	"Палитра: DIGITDISK_PALETTE=carbon|paper|signal, NO_COLOR и TERM=dumb уважаются.",
	"",
	"Язык: спрашивается один раз и помнится в ~/.digitable/digitdisk/settings.conf;",
	"  DIGITDISK_LANG=ru|en и --lang перекрывают его, --json не переводится.",
	"",
	"Подробно: man digitdisk",
}
