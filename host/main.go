// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Command digitdisk looks at a Linux system and at a directory tree, and
// reports what it found.  It never deletes, moves, or writes anything: the
// tool only reads /proc, /sys, and the tree it is pointed at.
//
// Subcommands:
//
//	digitdisk status [--json] [--top N] [--sample MS] [--live|--plain] [--interval MS]
//	digitdisk analyze <путь> [--json] [--top N] [--cross-device] [--max-depth N]
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

	"digitdisk/internal/report"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
	"digitdisk/internal/ui"
)

const usage = `digitdisk — снимок системы и разбор дерева каталогов (только чтение).

Использование:
  digitdisk status  [--json] [--top N] [--sample MS] [--live|--plain] [--interval MS]
      Снимок системы из /proc и /sys: ядро и дистрибутив, время работы,
      загрузка, память, процессы, диски, сеть, температура.
      В терминале — живой экран, который обновляется сам; в трубу, в файл и
      под --json — та же печать текстом, что и всегда.

  digitdisk analyze <путь> [--json] [--top N] [--cross-device] [--max-depth N]
      Обход дерева через lstat: символические ссылки не раскрываются,
      границу файловой системы без --cross-device не пересекаем,
      недоступное считается и пропускается.

Ключи:
  --json           машиночитаемый вывод
  --top N          сколько строк в списках (по умолчанию 10 / 15)
  --sample MS      окно замера загрузки ЦП, мс (status, по умолчанию 200)
  --cross-device   заходить на смонтированные другие файловые системы
  --max-depth N    предел глубины обхода (0 — без предела)
  --live           живой экран; без терминала — ошибка, а не тихая печать
  --plain          печать одним снимком, даже когда вывод в терминал
  --interval MS    период обновления живого экрана, мс (по умолчанию 2000)

Живой экран (status):
  ← →, Tab      предыдущий и следующий раздел      1…9  раздел сразу
  ↑ ↓, PgUp/Dn  прокрутка длинного раздела         p    пауза
  r             замер сейчас                       q    выход

  Палитра — Digitable Focus. DIGITDISK_PALETTE=carbon|paper|signal выбирает
  вариант (по умолчанию carbon); NO_COLOR и TERM=dumb уважаются.

Удаления нет ни в каком виде: digitdisk только смотрит и считает.
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
	case "-h", "--help", "help":
		fmt.Print(usage)
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
