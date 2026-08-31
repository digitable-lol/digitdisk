// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Command digitdisk looks at a Linux system and at a directory tree, and
// reports what it found.  It never deletes, moves, or writes anything: the
// tool only reads /proc, /sys, and the tree it is pointed at.
//
// Subcommands:
//
//	digitdisk status [--json] [--top N] [--sample MS]
//	digitdisk analyze <путь> [--json] [--top N] [--cross-device] [--max-depth N]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"digitdisk/internal/report"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
)

const usage = `digitdisk — снимок системы и разбор дерева каталогов (только чтение).

Использование:
  digitdisk status  [--json] [--top N] [--sample MS]
      Снимок системы из /proc и /sys: ядро и дистрибутив, время работы,
      загрузка, память, процессы, диски, сеть, температура.

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
	st := c.Collect()

	if *asJSON {
		return writeJSON(st)
	}
	report.Status(os.Stdout, st)
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
