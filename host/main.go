// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Command digitdisk looks at a system and at a directory tree, reports what it
// found, and — when told to, in as many words — removes files the decision
// layer marked «МожноУбрать».
//
//	digitdisk [подкоманда] [ключи]
//
// The subcommands and their one-line glosses are in internal/cli; the flags
// are registered in the cmd* functions below.  digitdisk.1 is the reference,
// and scripts/check-docs.sh holds it to what is registered here.
//
// A bare `digitdisk` runs status: reading is the frequent thing and it changes
// nothing.  A first argument that begins with «-» is status's own flags, so
// `digitdisk --json` and `digitdisk status --json` are one command; --help and
// --version stay themselves.  A first argument that is a word and not a
// subcommand is a mistake, and it is refused rather than guessed at.
//
// clean, restore and purge are the three steps of removal, and they are three
// because one would be a mistake nobody could take back:
//
//	clean <путь>          план: что, сколько, почему.  Ничего не тронуто.
//	clean <путь> --apply  перенос в корзину внутри корня.  Обратимо.
//	restore <корзина>     возврат на прежние места.
//	purge <корзина> --confirm N  стирание.  Необратимо.
//
// What may be removed is decided entirely by the layer in core/: exactly the
// paths it gives the приговор «МожноУбрать» and nothing that merely resembles
// one.  See internal/clean for the guards around that.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/cli"
	"digitdisk/internal/core"
	"digitdisk/internal/places"
	"digitdisk/internal/protect"
	"digitdisk/internal/report"
	"digitdisk/internal/scan"
	"digitdisk/internal/sysinfo"
	"digitdisk/internal/ui"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

// handlers is the other half of cli.Commands: a name there and a function
// here.  A map rather than a switch so the two lists can be compared by a
// test instead of by eye.
var handlers = map[string]func([]string) error{
	"status":  cmdStatus,
	"analyze": cmdAnalyze,
	"clean":   cmdClean,
	"restore": cmdRestore,
	"purge":   cmdPurge,
	"places":  cmdPlaces,
	"history": cmdHistory,
}

// run dispatches one command line and returns the code the process exits
// with.  Three cases have to be kept apart: подкоманду не назвали (do the
// default), назвали ключ (the default's own flags), and назвали слово,
// которого нет — refused with code 2, never answered by quietly doing
// something else.
func run(args []string) int {
	if len(args) > 0 {
		switch {
		case cli.Is(cli.HelpArgs, args[0]):
			fmt.Print(cli.Usage())
			return 0
		case cli.Is(cli.VersionArgs, args[0]):
			printVersion(os.Stdout)
			return 0
		}
	}

	name, rest := cli.Default, args
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		if !cli.Known(args[0]) {
			fmt.Fprintf(os.Stderr, "digitdisk: неизвестная подкоманда %q\n\n%s", args[0], cli.Usage())
			return 2
		}
		name, rest = args[0], args[1:]
	}

	h, ok := handlers[name]
	if !ok {
		// cli.Commands and handlers are the same list; a name in one and
		// not in the other is a defect, and it says so out loud.
		fmt.Fprintf(os.Stderr, "digitdisk: подкоманда %q объявлена в internal/cli, но не разобрана в main\n", name)
		return 2
	}
	if err := h(rest); err != nil {
		fmt.Fprintf(os.Stderr, "digitdisk: %v\n", err)
		return 1
	}
	return 0
}

func cmdStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	why := fs.Bool("why", false, "что не измерено и почему")
	top := fs.Int("top", 10, "сколько процессов в каждом списке")
	sample := fs.Int("sample", 200, "окно замера загрузки ЦП, мс")
	gpuTool := fs.Bool("gpu-tool", false, "спросить о видеокартах программу их драйвера (nvidia-smi)")
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
	c.GPUTool = *gpuTool

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
	placesFile := fs.String("places", "", "свой справочник известных мест")
	noPlaces := fs.Bool("no-places", false, "судить одними приметами, без справочника")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужен ровно один путь для обхода, получено %d", len(rest))
	}

	decider := chosenDecider()
	if _, err := usePlaces(decider, places.Options{File: *placesFile, Off: *noPlaces}); err != nil {
		return err
	}

	res, err := scan.Walk(scan.Options{
		Root:        rest[0],
		CrossDevice: *cross,
		MaxDepth:    *maxDepth,
		Top:         *top,
		Decider:     decider,
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
	top := fs.Int("top", 15, "сколько строк в перечнях, 0 — без предела")
	trash := fs.String("trash", "", "корзина (по умолчанию <корень>/"+clean.TrashName+"); обязана лежать внутри корня")
	cross := fs.Bool("cross-device", false, "заходить на другие файловые системы")
	maxDepth := fs.Int("max-depth", 0, "предел глубины обхода, 0 — без предела")
	placesFile := fs.String("places", "", "свой справочник известных мест")
	noPlaces := fs.Bool("no-places", false, "судить одними приметами, без справочника")
	protectFile := fs.String("protect-file", "", "защитный список файлом")
	var protectArgs stringList
	fs.Var(&protectArgs, "protect", "не трогать: путь или «разряд:кэш»; можно повторять")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужен ровно один путь для уборки, получено %d", len(rest))
	}

	decider := chosenDecider()
	dir, err := usePlaces(decider, places.Options{File: *placesFile, Off: *noPlaces})
	if err != nil {
		return err
	}
	guard, err := protect.Load(protect.Options{File: *protectFile, Args: protectArgs})
	if err != nil {
		return err
	}

	plan, err := clean.Make(clean.Options{
		Root:        rest[0],
		Trash:       *trash,
		CrossDevice: *cross,
		MaxDepth:    *maxDepth,
		Decider:     decider,
		Now:         time.Now(),
		Version:     version,
		Places:      dir,
		Protect:     guard,
	})
	if err != nil {
		return err
	}

	if !*apply {
		if *asJSON {
			// A machine gets the whole work list.  --top shortens what a
			// person reads; shortening what a script parses would make
			// `clean --json | jq` quietly wrong about the disk.
			return writeJSON(plan)
		}
		report.CleanPlan(os.Stdout, plan, *top)
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

// stringList collects a flag that may be given more than once.
type stringList []string

func (s *stringList) String() string     { return strings.Join(*s, ", ") }
func (s *stringList) Set(v string) error { *s = append(*s, v); return nil }

// usePlaces loads the справочник and hands it to the decision layer.
//
// A layer that does not take one (the заглушка, or any future layer built
// without the capability) is told so out loud rather than left to judge with a
// справочник nobody applied: the plan would then be smaller than the file
// promised, and nothing on the screen would say why.
func usePlaces(d core.Decider, opt places.Options) (*places.Directory, error) {
	dir, err := places.Load(opt)
	if err != nil {
		return nil, err
	}
	placer, ok := d.(core.Placer)
	if !ok {
		if len(dir.Entries) > 0 {
			fmt.Fprintf(os.Stderr,
				"digitdisk: решающий слой %q справочника не принимает — %d мест не применено\n",
				d.Name(), len(dir.Entries))
		}
		return dir, nil
	}
	if err := placer.UsePlaces(dir.Places()); err != nil {
		return nil, fmt.Errorf("справочник %s: %w", dir.Origin, err)
	}
	return dir, nil
}

// cmdPlaces prints the справочник and what of it exists here.  It reads and
// nothing else — the same promise status and analyze make.
func cmdPlaces(args []string) error {
	fs := flag.NewFlagSet("places", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	top := fs.Int("top", 40, "сколько найденных мест печатать, 0 — без предела")
	placesFile := fs.String("places", "", "свой справочник известных мест")
	noMeasure := fs.Bool("no-measure", false, "не считать размеры, только назвать места")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("подкоманда places не принимает путей (лишнее: %q)", rest[0])
	}

	dir, err := places.Load(places.Options{File: *placesFile})
	if err != nil {
		return err
	}

	measure := measureTree
	if *noMeasure {
		measure = nil
	}
	found := dir.Look(measure)

	if *asJSON {
		return writeJSON(struct {
			Origin string         `json:"откуда"`
			Count  int            `json:"мест"`
			Found  []places.Found `json:"места"`
		}{dir.Origin, len(dir.Entries), found})
	}
	report.Places(os.Stdout, dir, found, *top)
	return nil
}

// measureTree sums a directory the way analyze does — apparent size, hard
// links counted once — by asking the same walk, so the two commands can never
// answer differently about the same place.
func measureTree(root string) (int64, int, error) {
	res, err := scan.Walk(scan.Options{Root: root, Top: 1, Now: time.Now()})
	if err != nil {
		return 0, 0, err
	}
	return res.TotalBytes, res.Files, nil
}

// cmdHistory prints what past уборки did under a root.  It reads journals and
// writes nothing.
func cmdHistory(args []string) error {
	fs := flag.NewFlagSet("history", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "машиночитаемый вывод")
	top := fs.Int("top", 20, "сколько корзин печатать, 0 — без предела")
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return fmt.Errorf("нужен ровно один путь — корень уборки, хранилище корзин или одна корзина; получено %d", len(rest))
	}

	h, err := clean.ReadHistory(rest[0])
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(h)
	}
	report.History(os.Stdout, h, time.Now(), *top)
	return nil
}
