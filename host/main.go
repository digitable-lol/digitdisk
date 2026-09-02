// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Command digitdisk looks at a system and at a directory tree, reports what it
// found, and — when told to, in as many words — removes files the decision
// layer marked «МожноУбрать».
//
// status and analyze read and nothing else.  clean, restore and purge are the
// three steps of removal, and they are three because one would be a mistake
// nobody could take back:
//
//	clean <путь>          план: что, сколько, почему.  Ничего не тронуто.
//	clean <путь> --apply  перенос в корзину внутри корня.  Обратимо.
//	restore <корзина>     возврат на прежние места.
//	purge <корзина> --confirm N  стирание.  Необратимо.
//
// What may be removed is decided entirely by the layer in core/: exactly the
// paths it gives the приговор «МожноУбрать» and nothing that merely resembles
// one.  See internal/clean for the guards around that.
//
// Subcommands:
//
//	digitdisk status [--json] [--top N] [--sample MS] [--live|--plain] [--interval MS]
//	digitdisk analyze <путь> [--json] [--top N] [--cross-device] [--max-depth N]
//	digitdisk clean <путь> [--json] [--apply] [--trash DIR] [--cross-device] [--max-depth N]
//	digitdisk restore <корзина> [--json] [--dry-run]
//	digitdisk purge <корзина> [--json] [--confirm N]
//	digitdisk --version
//
// `status` draws a live screen when it is talking to a terminal and prints the
// same snapshot as text when it is not, so a pipe, a file and a script see
// exactly what they always saw.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/report"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
	"digitdisk/internal/ui"
)

const usage = `digitdisk — снимок системы, разбор дерева каталогов и уборка.

Использование:
  digitdisk status  [--json] [--top N] [--sample MS] [--why] [--live|--plain] [--interval MS]
      Снимок системы: ядро и выпуск, время работы, загрузка, память, процессы,
      диски, сеть, температура. Источники платформенные — /proc и /sys на
      Linux, sysctl и libSystem на macOS. Чего система не публикует, печатается
      прочерком, а не нулём, и называется одной строкой в конце; почему — по
      ключу --why.
      В терминале — живой экран, который обновляется сам; в трубу, в файл и
      под --json — та же печать текстом, что и всегда.

  digitdisk analyze <путь> [--json] [--top N] [--cross-device] [--max-depth N]
      Обход дерева через lstat: символические ссылки не раскрываются,
      границу файловой системы без --cross-device не пересекаем,
      недоступное считается и пропускается.

  digitdisk clean <путь> [--json] [--apply] [--trash КАТ] [--cross-device] [--max-depth N]
      Уборка. БЕЗ КЛЮЧА --apply НИЧЕГО НЕ ТРОГАЕТ: печатает план — какие файлы,
      сколько байт и по какому правилу ядра помечены «МожноУбрать». С --apply
      переносит их в корзину <корень>/.digitdisk-trash/<метка времени>/ и пишет
      журнал. Перенос — это rename(2): мгновенно, обратимо и НИЧЕГО НЕ
      ОСВОБОЖДАЕТ, файлы остаются на диске под другим именем.

  digitdisk restore <корзина> [--json] [--dry-run]
      Возврат корзины на прежние места по её журналу. Ключа не требует: ключи
      стоят на разрушении, а не на его отмене. Ничего не перезаписывает.

  digitdisk purge <корзина> [--json] [--confirm N]
      Стирание корзины — единственное необратимое действие. Без --confirm
      печатает план и число, которое надо назвать. С --confirm N, где N —
      ровно столько файлов, сколько в корзине, стирает их по одному.

  digitdisk --version
      Версия, хеш сборки, инструментарий и решающий слой этого двоичного файла.

Ключи:
  --json           машиночитаемый вывод
  --why            вместо снимка — что не измерено и почему (status)
  --top N          сколько строк в списках (по умолчанию 10 / 15)
  --sample MS      окно замера загрузки ЦП, мс (status, по умолчанию 200)
  --cross-device   заходить на смонтированные другие файловые системы
  --max-depth N    предел глубины обхода (0 — без предела)
  --apply          clean: перенести в корзину, а не только показать план
  --trash КАТ      clean: другая корзина; обязана лежать внутри корня
  --dry-run        restore: показать, что вернулось бы, и не возвращать
  --confirm N      purge: подтвердить стирание ровно N файлов
  --live           живой экран; без терминала — ошибка, а не тихая печать
  --plain          печать одним снимком, даже когда вывод в терминал
  --interval MS    период обновления живого экрана, мс (по умолчанию 2000)

Живой экран (status):
  ← →, Tab      предыдущий и следующий раздел      1…9  раздел сразу
  ↑ ↓, PgUp/Dn  прокрутка длинного раздела         p    пауза
  r             замер сейчас                       q    выход

  Палитра — Digitable Focus. DIGITDISK_PALETTE=carbon|paper|signal выбирает
  вариант (по умолчанию carbon); NO_COLOR и TERM=dumb уважаются.

Убирается ровно то, чему решающий слой вынес приговор «МожноУбрать», — не
похожее на него и не совпавшее с маской. status и analyze не пишут ничего.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "status":
		err = cmdStatus(os.Args[2:])
	case "analyze":
		err = cmdAnalyze(os.Args[2:])
	case "clean":
		err = cmdClean(os.Args[2:])
	case "restore":
		err = cmdRestore(os.Args[2:])
	case "purge":
		err = cmdPurge(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Print(usage)
		return
	case "-V", "--version", "version":
		printVersion(os.Stdout)
		return
	default:
		fmt.Fprintf(os.Stderr, "digitdisk: неизвестная подкоманда %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "digitdisk: %v\n", err)
		os.Exit(1)
	}
}

func cmdStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	why := fs.Bool("why", false, "что не измерено и почему")
	top := fs.Int("top", 10, "сколько процессов в каждом списке")
	sample := fs.Int("sample", 200, "окно замера загрузки ЦП, мс")
	live := fs.Bool("live", false, "живой экран, даже если о терминале не спрашивали")
	plain := fs.Bool("plain", false, "печать одним снимком, без живого экрана")
	interval := fs.Int("interval", 2000, "период обновления живого экрана, мс")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("подкоманда status не принимает путей (лишнее: %q)", rest[0])
	}

	c := sysinfo.New()
	c.Top = *top
	c.SampleWindow = time.Duration(*sample) * time.Millisecond

	// A machine reader is answered first and never gets the screen: --json is
	// how scripts call this tool, and its output must not depend on where it
	// is pointed.
	if *asJSON {
		return writeJSON(c.Collect())
	}
	// --why answers one question and prints nothing else, so it comes before
	// the screen: somebody asking why a number is missing wants the answer,
	// not a dashboard.
	if *why {
		report.Why(os.Stdout, c.Collect())
		return nil
	}

	// The screen is the default only when there is a terminal to draw it on.
	// A pipe, a file, /dev/null, TERM=dumb and an empty TERM all fall through
	// to the printed report, which is what they have always received.
	if *live || (!*plain && ui.Available(os.Stdout)) {
		err := ui.Run(ui.Options{
			Out:      os.Stdout,
			Interval: time.Duration(*interval) * time.Millisecond,
			Palette:  ui.PaletteByName(os.Getenv("DIGITDISK_PALETTE")),
			Collect:  c.Collect,
		})
		switch {
		case err == nil:
			return nil
		case !errors.Is(err, ui.ErrNoTerminal):
			return err
		case *live:
			return fmt.Errorf("%w; без --live тот же снимок печатается текстом", err)
		}
		// The terminal went away between the question and the answer.  Print.
	}

	report.Status(os.Stdout, c.Collect())
	return nil
}

func cmdAnalyze(args []string) error {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	top := fs.Int("top", 15, "сколько строк в списках")
	cross := fs.Bool("cross-device", false, "заходить на другие файловые системы")
	maxDepth := fs.Int("max-depth", 0, "предел глубины обхода, 0 — без предела")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужен ровно один путь для обхода, получено %d", len(rest))
	}

	res, err := scan.Walk(scan.Options{
		Root:        rest[0],
		CrossDevice: *cross,
		MaxDepth:    *maxDepth,
		Top:         *top,
		Decider:     chosenDecider(),
		Now:         time.Now(),
	})
	if err != nil {
		return err
	}

	if *asJSON {
		return writeJSON(res)
	}
	report.Analyze(os.Stdout, res)
	return nil
}

// cmdClean prints the plan, and moves files into the корзина only when told
// to.  The default has to be the harmless one: a person who runs `clean` to
// find out what it would do must not find out by having it done.
func cmdClean(args []string) error {
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	apply := fs.Bool("apply", false, "перенести в корзину, а не только показать план")
	trash := fs.String("trash", "", "корзина (по умолчанию <корень>/"+clean.TrashName+"); обязана лежать внутри корня")
	cross := fs.Bool("cross-device", false, "заходить на другие файловые системы")
	maxDepth := fs.Int("max-depth", 0, "предел глубины обхода, 0 — без предела")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужен ровно один путь для уборки, получено %d", len(rest))
	}

	plan, err := clean.Make(clean.Options{
		Root:        rest[0],
		Trash:       *trash,
		CrossDevice: *cross,
		MaxDepth:    *maxDepth,
		Decider:     chosenDecider(),
		Now:         time.Now(),
		Version:     version,
	})
	if err != nil {
		return err
	}

	if !*apply {
		if *asJSON {
			return writeJSON(plan)
		}
		report.CleanPlan(os.Stdout, plan)
		return nil
	}

	j, err := clean.Apply(plan, clean.Options{Now: time.Now(), Version: version})
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(j)
	}
	report.Applied(os.Stdout, j)
	return nil
}

// cmdRestore puts a корзина back.  It acts without a key by design: see the
// comment on clean.Restore.
func cmdRestore(args []string) error {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	dry := fs.Bool("dry-run", false, "показать, что вернулось бы, и не возвращать")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужна ровно одна корзина, получено %d", len(rest))
	}

	j, err := clean.ReadJournal(rest[0])
	if err != nil {
		return err
	}
	j, err = clean.Restore(j, *dry, time.Now())
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(j)
	}
	report.Restored(os.Stdout, j, *dry)
	return nil
}

// cmdPurge erases a корзина, and only with the count named.
func cmdPurge(args []string) error {
	fs := flag.NewFlagSet("purge", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	confirm := fs.Int("confirm", -1, "подтвердить стирание ровно N файлов")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужна ровно одна корзина, получено %d", len(rest))
	}

	j, err := clean.ReadJournal(rest[0])
	if err != nil {
		return err
	}
	if *confirm < 0 {
		if *asJSON {
			return writeJSON(j)
		}
		report.PurgePlan(os.Stdout, j)
		return nil
	}

	j, err = clean.Purge(j, *confirm, time.Now())
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(j)
	}
	report.Purged(os.Stdout, j)
	return nil
}

// parseFlags parses args allowing flags to appear after positional arguments,
// which the standard flag package stops at.  Parsing resumes after each
// positional, so a flag that takes a separate value ("--top 5") still works.
func parseFlags(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		if fs.NArg() == 0 {
			return positional, nil
		}
		positional = append(positional, fs.Arg(0))
		args = fs.Args()[1:]
	}
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}
