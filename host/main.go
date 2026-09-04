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
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/cli"
	"digitdisk/internal/core"
	"digitdisk/internal/lang"
	"digitdisk/internal/places"
	"digitdisk/internal/protect"
	"digitdisk/internal/report"
	"digitdisk/internal/scan"
	"digitdisk/internal/settings"
	"digitdisk/internal/sysinfo"
	"digitdisk/internal/ui"

	// The package is named for what it does; the import is renamed because
	// this file already has a run of its own — the one that dispatches a
	// command line.
	wrap "digitdisk/internal/run"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

// handlers is the other half of cli.Commands: a name there and a function
// here.  A map rather than a switch so the two lists can be compared by a
// test instead of by eye.
//
// It is filled in init and not where it is declared because the живой экран
// dispatches through it too — cmdStatus draws the screen, the screen names a
// подкоманда, runFromScreen looks it up here — and a map literal naming
// cmdStatus while cmdStatus reaches the map is a cycle Go refuses at compile
// time.  Filling it in init says the same thing and is allowed to be circular,
// which this is: one dispatch table, entered twice.
var handlers map[string]func([]string) error

func init() {
	handlers = map[string]func([]string) error{
		"status":  cmdStatus,
		"analyze": cmdAnalyze,
		"clean":   cmdClean,
		"restore": cmdRestore,
		"purge":   cmdPurge,
		"places":  cmdPlaces,
		"history": cmdHistory,
		"run":     cmdRun,
	}
}

// run dispatches one command line and returns the code the process exits
// with.  Three cases have to be kept apart: подкоманду не назвали (do the
// default), назвали ключ (the default's own flags), and назвали слово,
// которого нет — refused with code 2, never answered by quietly doing
// something else.
func run(args []string) int {
	// --lang стоит перед подкомандой не реже, чем после неё: «digitdisk
	// --lang en clean ~» читается как одно распоряжение, а не как ключ
	// подкоманды status. Поэтому он снимается с начала строки здесь, до
	// того как первое слово решит, какая это подкоманда, — и после
	// снятия «digitdisk --lang en clean ~» это ровно «digitdisk clean ~».
	args, globalLang := stripGlobalLang(args)
	if globalLang != "" {
		chooseLang(globalLang, false)
	}
	if len(args) > 0 {
		switch {
		case cli.Is(cli.HelpArgs, args[0]):
			fmt.Print(cli.Usage(langOnly(args)))
			return 0
		case cli.Is(cli.VersionArgs, args[0]):
			printVersion(os.Stdout, langOnly(args))
			return 0
		}
	}

	name, rest := cli.Default, args
	switch i := runSplit(args); {
	case i >= 0:
		// `digitdisk -c make -j8` is `digitdisk run make -j8`.  The -c
		// itself is dropped and the two halves are joined: what stood
		// before it is ours, what stands after it is the command's, and
		// нашими ключами команда не распоряжается.
		name = "run"
		rest = append(append([]string{}, args[:i]...), args[i+1:]...)
	case len(args) > 0 && !strings.HasPrefix(args[0], "-"):
		if !cli.Known(args[0]) {
			l := langOnly(args)
			fmt.Fprintf(os.Stderr, "digitdisk: "+l.T("неизвестная подкоманда %q")+"\n\n%s", args[0], cli.Usage(l))
			return 2
		}
		name, rest = args[0], args[1:]
	}

	h, ok := handlers[name]
	if !ok {
		// cli.Commands and handlers are the same list; a name in one and
		// not in the other is a defect, and it says so out loud.
		fmt.Fprintf(os.Stderr, "digitdisk: "+langOnly(args).T("подкоманда %q объявлена в internal/cli, но не разобрана в main")+"\n", name)
		return 2
	}
	if err := h(rest); err != nil {
		// A wrapped command's own код возврата travels out as an error
		// nobody prints: `digitdisk run false` must be `false` in a
		// script, and a wrapper that turned every ending into 1 could
		// not be put into one.
		var st exitStatus
		if errors.As(err, &st) {
			if st.err != nil {
				fmt.Fprintf(os.Stderr, "digitdisk: %s\n", lang.InLang(st.err, langChoice.Lang))
			}
			return st.code
		}
		// The refusal is rendered here and nowhere else: it was built four
		// levels down, where nobody knows what language the reader has.
		fmt.Fprintf(os.Stderr, "digitdisk: %s\n", lang.InLang(err, langChoice.Lang))
		return 1
	}
	return 0
}

// exitStatus carries a код возврата out of a подкоманда without a word being
// printed about it.  Only run uses it, and only because the ending it has to
// report is not its own.
//
// Its Error text is never read by anybody: run() answers exitStatus before it
// prints anything, and the message a person sees, when there is one, is the
// wrapped err — which is a lang.Error and speaks their language.
type exitStatus struct {
	code int
	err  error
}

func (e exitStatus) Error() string { return "exit status " + strconv.Itoa(e.code) }

func (e exitStatus) Unwrap() error { return e.err }

// runSplit finds the -c that says «everything after this is the command».
//
// The rule is the one env(1), nice(1) and time(1) settled on long ago, and it
// is the only rule under which `digitdisk -c make -j8` can work at all: OUR
// keys stand before the command, the command's keys stand after it, and -c is
// the line between them.  Otherwise --json in `digitdisk -c ls --json` would
// belong to whichever of the two asked for it more loudly.
//
// The search stops at the first word that is not a key of ours: `digitdisk
// clean ~ -c` is clean's business, and `digitdisk фигня -c ls` is a mistake to
// be refused rather than a command to be run.
func runSplit(args []string) int {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if cli.Is(cli.RunArgs, a) {
			return i
		}
		if !strings.HasPrefix(a, "-") {
			return -1
		}
		if takesValue(a) {
			i++
		}
	}
	return -1
}

// takesValue reports whether one of run's own keys is followed by its value
// as a separate word, so that `digitdisk --interval 500 -c make` finds the -c
// behind the 500.
func takesValue(a string) bool {
	name := strings.TrimLeft(a, "-")
	if strings.ContainsRune(name, '=') {
		return false
	}
	return name == langFlagName || name == "interval"
}

func cmdStatus(args []string) error {
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	why := fs.Bool("why", false, l.T("что не измерено и почему"))
	top := fs.Int("top", 10, l.T("сколько процессов в каждом списке"))
	sample := fs.Int("sample", 200, l.T("окно замера загрузки ЦП, мс"))
	gpuTool := fs.Bool("gpu-tool", false, l.T("спросить о видеокартах программу их драйвера (nvidia-smi)"))
	live := fs.Bool("live", false, l.T("живой экран, даже если о терминале не спрашивали"))
	plain := fs.Bool("plain", false, l.T("печать одним снимком, без живого экрана"))
	interval := fs.Int("interval", 2000, l.T("период обновления живого экрана, мс"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return lang.Errorf("подкоманда status не принимает путей (лишнее: %q)", rest[0])
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
		report.Why(os.Stdout, l, c.Collect())
		return nil
	}

	// The screen is the default only when there is a terminal to draw it on.
	// A pipe, a file, /dev/null, TERM=dumb and an empty TERM all fall through
	// to the printed report, which is what they have always received.
	if *live || (!*plain && ui.Available(os.Stdout)) {
		err := statusScreen(c.Collect, time.Duration(*interval)*time.Millisecond)
		switch {
		case err == nil:
			return nil
		case !errors.Is(err, ui.ErrNoTerminal):
			return err
		case *live:
			return lang.Errorf("%s; без --live тот же снимок печатается текстом", err)
		}
		// The terminal went away between the question and the answer.  Print.
	}

	report.Status(os.Stdout, l, c.Collect())
	return nil
}

// statusScreen draws the live screen, runs what the reader chose from its
// КОМАНДЫ section, and draws it again.
//
// THIS LOOP IS THE WHOLE OF «вернуться, не выходя из программы».  The process
// never restarts: what changes hands is the terminal.  Every screen this tool
// draws puts the terminal into raw mode and takes the signals for its own
// lifetime — so the экран состояния closes and hands the terminal back exactly
// as it found it, the подкоманда runs (for analyze that is the walk screen,
// which owns the terminal in its turn), and the экран состояния opens again.
// The keyboard itself is not handed back and forth: internal/ui reads it once
// for the whole process, and that is what keeps a keypress from falling
// between two screens.
//
// The language chosen along the way is carried across in langChoice, so a
// reader who switched to English inside analyze comes back to an English экран
// состояния.
//
// What a подкоманда PRINTED stays on the terminal until the reader has read
// it: the screen would otherwise cover its own output with its first frame.
func statusScreen(collect func() sysinfo.Status, interval time.Duration) error {
	for {
		req, err := ui.Run(ui.Options{
			Out:      os.Stdout,
			Interval: interval,
			Palette:  ui.PaletteByName(os.Getenv("DIGITDISK_PALETTE")),
			Collect:  collect,
			Lang:     langChoice.Lang,
			Remember: rememberLang,
		})
		if err != nil || req == nil {
			return err
		}
		printed, err := runFromScreen(*req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "digitdisk: %s\n", lang.InLang(err, langChoice.Lang))
			printed = true
		}
		if printed {
			waitForReader(langChoice.Lang)
		}
	}
}

// runFromScreen runs one подкоманда chosen on the живой экран and says whether
// it left anything on the terminal to be read.
//
// There is no second dispatch here: analyze and clean are the SAME call
// cmdAnalyze makes, with the same справочник, the same защитный список and the
// same решающий слой, and everything else goes through the handlers map — the
// one main dispatches the command line with.  Which подкоманды may arrive at
// all is decided in internal/cli beside the подкоманда itself, so `purge`
// cannot reach here: the screen never offers to start it.
func runFromScreen(req ui.Request) (printed bool, err error) {
	switch req.Command {
	case "analyze":
		return analyze(nil, ui.AfterNothing)
	case "clean":
		// Уборка с экрана — это обход, а затем ПЛАН, и ни одним шагом
		// меньше: приговор ядра, разбивка по разрядам, корзина и точное
		// число файлов, набранное руками. Печатный `clean` без --apply
		// делает ровно это же и тоже ничего не двигает.
		return analyze(nil, ui.AfterPlan)
	}
	h, ok := handlers[req.Command]
	if !ok {
		return false, lang.Errorf("подкоманда %q объявлена в internal/cli, но не разобрана в main", req.Command)
	}
	return true, h(nil)
}

// waitForReader holds printed output on the terminal until the reader says
// they have read it.  Without it the экран состояния would draw its first
// frame over the answer the reader just asked for.
//
// The waiting is done by internal/ui and not by a second reader of /dev/tty
// opened here, for the reason that package spells out at length: one terminal
// answers one reader, and a second one silently takes the keys of the first.
func waitForReader(l lang.Lang) {
	ui.WaitKey(os.Stdout, "\n"+l.T("— Enter возвращает на экран состояния")+" ")
}

func cmdAnalyze(args []string) error {
	_, err := analyze(args, ui.AfterNothing)
	return err
}

// analyze is `digitdisk analyze`, plus what the живой экран needs of it: what
// to do when the walk finishes, and whether anything was printed.
func analyze(args []string, after ui.After) (printed bool, err error) {
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	top := fs.Int("top", 15, l.T("сколько строк в списках"))
	cross := fs.Bool("cross-device", false, l.T("заходить на другие файловые системы"))
	maxDepth := fs.Int("max-depth", 0, l.T("предел глубины обхода, 0 — без предела"))
	placesFile := fs.String("places", "", l.T("свой справочник известных мест"))
	noPlaces := fs.Bool("no-places", false, l.T("судить одними приметами, без справочника"))
	protectFile := fs.String("protect-file", "", l.T("защитный список файлом"))
	live := fs.Bool("live", false, l.T("живой экран, даже если о терминале не спрашивали"))
	plain := fs.Bool("plain", false, l.T("обойти молча и напечатать отчёт, без экрана"))
	noFold := fs.Bool("no-fold", false, l.T("судить о каждом файле внутри node_modules и подобных, а не о каталоге целиком"))
	var protectArgs stringList
	fs.Var(&protectArgs, "protect", l.T("не трогать: путь или «разряд:кэш»; можно повторять"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return false, err
	}
	// More than one path is a mistake; none of them is the screen being
	// asked for a path, and only in a terminal — see below.
	if len(rest) > 1 {
		return false, lang.Errorf("нужен ровно один путь для обхода, получено %d", len(rest))
	}

	decider := chosenDecider(l)
	dir, err := usePlaces(l, decider, places.Options{File: *placesFile, Off: *noPlaces})
	if err != nil {
		return false, err
	}
	sayMoved(l, dir.Moved)

	opt := scan.Options{
		CrossDevice: *cross,
		MaxDepth:    *maxDepth,
		Top:         *top,
		Decider:     decider,
	}
	// Свёртка тяжёлых каталогов включена по умолчанию: смысл её в скорости,
	// а платит за неё только подробность приговора внутри — размер и итог
	// остаются теми же. Кому нужна подробность, тот её просит.
	if !*noFold {
		opt.Fold = scan.FoldByName
	}
	if len(rest) == 1 {
		opt.Root = rest[0]
	}

	// A machine reader is answered first and never gets the screen, for the
	// reason status gives: --json is how scripts call this tool.
	//
	// The screen is the default only where there is a terminal to draw it on,
	// and it changes nothing about what is printed: the walk is the same walk,
	// the result is the same result, and the report goes to standard output
	// after the screen closes exactly as it always has.
	if !*asJSON && (*live || (!*plain && ui.Available(os.Stdout))) {
		guard, err := protect.Load(protect.Options{File: *protectFile, Args: protectArgs})
		if err != nil {
			return false, err
		}

		// The report printed after the screen closes is printed in the
		// language the screen closed in.  A reader who switched to English
		// and quit would otherwise be handed a Russian report by the same
		// keypress that ended an English screen.
		shown := l
		remember := func(chosen lang.Lang) lang.Phrase {
			shown = chosen
			return rememberLang(chosen)
		}
		res, have, err := ui.RunWalk(ui.WalkOptions{
			Out:      os.Stdout,
			Root:     opt.Root,
			Palette:  ui.PaletteByName(os.Getenv("DIGITDISK_PALETTE")),
			Lang:     l,
			Remember: remember,
			After:    after,
			Walk: func(root string, watch func(scan.Step), stop func() bool) (scan.Result, error) {
				o := opt
				o.Root, o.Watch, o.Now = root, watch, time.Now()
				// Экран отпустил обход — значит его никто не читает.
				// Prune роняет всё, что ниже, и обход сворачивается за
				// доли секунды вместо оставшихся минут.
				o.Prune = func(string) bool { return stop() }
				return scan.Walk(o)
			},
			// The screen acts through the same calls the subcommands make,
			// with the same справочник, the same защитный список and the same
			// decision layer.  It builds nothing of its own: a second road to
			// removal is exactly what internal/clean exists to prevent.
			Plan: func(root string, only []string) (*clean.Plan, error) {
				p, err := clean.Make(clean.Options{
					Root: root, CrossDevice: *cross, MaxDepth: *maxDepth,
					Decider: decider, Now: time.Now(), Version: version,
					Places: dir, Protect: guard, Only: only,
				})
				if err != nil {
					return nil, err
				}
				return &p, nil
			},
			Apply: func(p *clean.Plan) (*clean.Journal, error) {
				return clean.Apply(*p, clean.Options{Now: time.Now(), Version: version})
			},
			// The ground забой was pointed at, judged before the walk:
			// the твёрдые запреты and the защитный список, in that
			// order.  clean.Make asks both again — this is only so
			// that the refusal comes back in one keypress instead of
			// after a walk, and so that the screen can lay its three
			// parts out instead of printing one long line.
			HardStop: func(ground string) *clean.Stop {
				if s := clean.HardStop(ground, clean.StopOptions{}); !s.Empty() {
					return &s
				}
				if rule, ok := guard.Covers(ground, core.ClassUnknown); ok {
					s := clean.ProtectStop(ground, rule)
					return &s
				}
				return nil
			},
			// The OTHER question, and the only place it is asked from.
			// Same корень, same справочник, same защитный список, same
			// решающий слой — and ByHand, which is what makes the plan
			// «стереть вот это» instead of «найди мусор сам».  No
			// подкоманда builds one: забой is the whole of it, and a
			// person is looking at the screen when it happens.
			PlanByHand: func(root string, only []string) (*clean.Plan, error) {
				p, err := clean.Make(clean.Options{
					Root: root, CrossDevice: *cross, MaxDepth: *maxDepth,
					Decider: decider, Now: time.Now(), Version: version,
					Places: dir, Protect: guard, Only: only, ByHand: true,
				})
				if err != nil {
					return nil, err
				}
				return &p, nil
			},
			// One verb over either plan: what забой does after the
			// question is confirmed.  It is handed the plan and nothing
			// else, so the screen cannot name a path of its own here any
			// more than it can when the file goes to the корзина.
			Erase: func(p *clean.Plan) (*clean.Journal, error) {
				return clean.Erase(*p, clean.Options{Now: time.Now(), Version: version})
			},
			Restore: func(box string, dryRun bool) (*clean.Journal, error) {
				j, err := clean.ReadJournal(box)
				if err != nil {
					return nil, err
				}
				return clean.Restore(j, dryRun, time.Now())
			},
			History: func(root string) (*clean.History, error) { return clean.ReadHistory(root) },
			Places: func() (string, []ui.PlaceRow, error) {
				d, err := places.Load(places.Options{File: *placesFile})
				if err != nil {
					return "", nil, err
				}
				var rows []ui.PlaceRow
				// Sizes are not measured here: measuring a hundred places
				// means a hundred walks, and `digitdisk places` is the
				// command that does it.
				for _, f := range d.Look(nil) {
					if !f.Exists {
						continue
					}
					rows = append(rows, ui.PlaceRow{
						Class: string(f.Entry.Class), Name: f.Entry.Name, Path: f.Entry.Resolved,
					})
				}
				return d.Origin, rows, nil
			},
		})
		switch {
		case err == nil:
			if have {
				report.Analyze(os.Stdout, shown, res)
			}
			return have, nil
		case errors.Is(err, ui.ErrWalkStopped):
			return false, err
		case !errors.Is(err, ui.ErrNoTerminal):
			return false, err
		case *live:
			return false, lang.Errorf("%s; без --live тот же обход печатается текстом", err)
		}
		// The terminal went away between the question and the answer.  Walk.
	}

	if opt.Root == "" {
		return false, lang.Errorf("нужен ровно один путь для обхода, получено %d", 0)
	}
	opt.Now = time.Now()
	res, err := scan.Walk(opt)
	if err != nil {
		return false, err
	}

	if *asJSON {
		return true, writeJSON(res)
	}
	report.Analyze(os.Stdout, l, res)
	return true, nil
}

// cmdClean prints the plan, and moves files into the корзина only when told
// to.  The default has to be the harmless one: a person who runs `clean` to
// find out what it would do must not find out by having it done.
func cmdClean(args []string) error {
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	apply := fs.Bool("apply", false, l.T("перенести в корзину, а не только показать план"))
	top := fs.Int("top", 15, l.T("сколько строк в перечнях, 0 — без предела"))
	trash := fs.String("trash", "", l.F("корзина (по умолчанию <корень>/%s); обязана лежать внутри корня", clean.TrashName))
	cross := fs.Bool("cross-device", false, l.T("заходить на другие файловые системы"))
	maxDepth := fs.Int("max-depth", 0, l.T("предел глубины обхода, 0 — без предела"))
	placesFile := fs.String("places", "", l.T("свой справочник известных мест"))
	noPlaces := fs.Bool("no-places", false, l.T("судить одними приметами, без справочника"))
	protectFile := fs.String("protect-file", "", l.T("защитный список файлом"))
	var protectArgs stringList
	fs.Var(&protectArgs, "protect", l.T("не трогать: путь или «разряд:кэш»; можно повторять"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return lang.Errorf("нужен ровно один путь для уборки, получено %d", len(rest))
	}

	decider := chosenDecider(l)
	dir, err := usePlaces(l, decider, places.Options{File: *placesFile, Off: *noPlaces})
	if err != nil {
		return err
	}
	guard, err := protect.Load(protect.Options{File: *protectFile, Args: protectArgs})
	if err != nil {
		return err
	}
	sayMoved(l, dir.Moved, guard.Moved)

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
		report.CleanPlan(os.Stdout, l, plan, *top)
		return nil
	}

	j, err := clean.Apply(plan, clean.Options{Now: time.Now(), Version: version})
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(j)
	}
	report.Applied(os.Stdout, l, j)
	return nil
}

// cmdRestore puts a корзина back.  It acts without a key by design: see the
// comment on clean.Restore.
func cmdRestore(args []string) error {
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	dry := fs.Bool("dry-run", false, l.T("показать, что вернулось бы, и не возвращать"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return lang.Errorf("нужна ровно одна корзина, получено %d", len(rest))
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
	report.Restored(os.Stdout, l, j, *dry)
	return nil
}

// cmdPurge erases a корзина, and only with the count named.
func cmdPurge(args []string) error {
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("purge", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	confirm := fs.Int("confirm", -1, l.T("подтвердить стирание ровно N файлов"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return lang.Errorf("нужна ровно одна корзина, получено %d", len(rest))
	}

	j, err := clean.ReadJournal(rest[0])
	if err != nil {
		return err
	}
	if *confirm < 0 {
		if *asJSON {
			return writeJSON(j)
		}
		report.PurgePlan(os.Stdout, l, j)
		return nil
	}

	j, err = clean.Purge(j, *confirm, time.Now())
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(j)
	}
	report.Purged(os.Stdout, l, j)
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
func usePlaces(l lang.Lang, d core.Decider, opt places.Options) (*places.Directory, error) {
	dir, err := places.Load(opt)
	if err != nil {
		return nil, err
	}
	placer, ok := d.(core.Placer)
	if !ok {
		if len(dir.Entries) > 0 {
			fmt.Fprintf(os.Stderr, "digitdisk: "+l.T("решающий слой %q справочника не принимает — %d мест не применено")+"\n",
				d.Name(), len(dir.Entries))
		}
		return dir, nil
	}
	if err := placer.UsePlaces(dir.Places()); err != nil {
		return nil, lang.Errorf("справочник %s: %s", dir.Origin, err)
	}
	return dir, nil
}

// cmdPlaces prints the справочник and what of it exists here.  It reads and
// nothing else — the same promise status and analyze make.
func cmdPlaces(args []string) error {
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("places", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	top := fs.Int("top", 40, l.T("сколько найденных мест печатать, 0 — без предела"))
	placesFile := fs.String("places", "", l.T("свой справочник известных мест"))
	noMeasure := fs.Bool("no-measure", false, l.T("не считать размеры, только назвать места"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return lang.Errorf("подкоманда places не принимает путей (лишнее: %q)", rest[0])
	}

	dir, err := places.Load(places.Options{File: *placesFile})
	if err != nil {
		return err
	}
	sayMoved(l, dir.Moved)

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
	report.Places(os.Stdout, l, dir, found, *top)
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
	l := chooseLang(peekLang(args), !peekJSON(args))
	fs := flag.NewFlagSet("history", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	top := fs.Int("top", 20, l.T("сколько корзин печатать, 0 — без предела"))
	rest, err := parseFlags(fs, args)
	if err != nil {
		return err
	}
	if len(rest) != 1 {
		return lang.Errorf("нужен ровно один путь — корень уборки, хранилище корзин или одна корзина; получено %d", len(rest))
	}

	h, err := clean.ReadHistory(rest[0])
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(h)
	}
	report.History(os.Stdout, l, h, time.Now(), *top)
	return nil
}

// cmdRun starts somebody else's command and says what it cost.
//
// It is the one подкоманда that prints nothing to standard output: that
// descriptor is the command's, and everything of ours — the строка состояния
// while it runs, the сводка after it, and the сводка under --json too — goes
// to standard error.  `digitdisk run make | tee log` has to be `make | tee
// log` byte for byte, and there is no other way to promise that.
//
// Ключи ставятся ДО команды, её собственные — после: `digitdisk run --json ls
// --json` gives the first --json to us and the second to ls.  With the short
// spelling the line between the two is -c itself.
func cmdRun(args []string) error {
	mine := ourArgs(args)
	l := chooseLang(peekLang(mine), !peekJSON(mine))
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	asJSON := fs.Bool("json", false, l.T("машиночитаемый вывод"))
	_ = fs.String("lang", "", l.T("язык вывода на этот запуск: ru или en"))
	plain := fs.Bool("plain", false, l.T("без строки состояния, даже в терминале"))
	interval := fs.Int("interval", 1000, l.T("период обновления строки состояния, мс"))
	gpuTool := fs.Bool("gpu-tool", false, l.T("спросить о видеопамяти программу драйвера (nvidia-smi)"))
	// Not parseFlags: it resumes parsing after a positional argument, and
	// here the first positional argument is the command — after which every
	// word belongs to the command and to nothing else.
	if err := fs.Parse(args); err != nil {
		return err
	}
	rest := fs.Args()
	if len(rest) == 0 {
		return lang.Errorf("нужна команда: digitdisk run <команда> [доводы]")
	}

	res, err := wrap.Run(wrap.Options{
		Args:     rest,
		Shell:    shellFor(rest),
		In:       os.Stdin,
		Out:      os.Stdout,
		Err:      os.Stderr,
		Interval: time.Duration(*interval) * time.Millisecond,
		Plain:    *plain,
		GPUTool:  *gpuTool,
		Lang:     l,
	})
	if err != nil {
		// 127 and 126 are what a shell answers for «нет такой команды» and
		// «не запускается», and a wrapper that answered anything else
		// would break the scripts that already test for them.
		return exitStatus{code: wrap.StartCode(err), err: err}
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stderr)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(res); err != nil {
			return err
		}
	} else {
		for _, line := range wrap.Lines(l, res) {
			fmt.Fprintf(os.Stderr, "digitdisk: %s\n", line)
		}
	}

	if res.SignalNumber != 0 {
		return raise(res)
	}
	return exitStatus{code: res.Code}
}

// raise ends this process the way the command ended: with the same signal.
//
// A код возврата of 128+N looks the same to a shell script, but it is not the
// same thing to a shell: `digitdisk run sleep 60` interrupted from the
// keyboard must leave the shell's own «^C» behaviour intact, and that happens
// only if the process really dies of the signal.  If it does not — a signal
// that stops instead of killing, a signal somebody blocked — the number is
// what is left, and that is what is returned.
func raise(res wrap.Result) error {
	s := syscall.Signal(res.SignalNumber)
	signal.Reset(s)
	_ = syscall.Kill(os.Getpid(), s)
	return exitStatus{code: res.Code}
}

// ourArgs is the part of a run command line that belongs to digitdisk: the
// keys before the command's name, and the value of a key that takes one.
//
// It exists because --lang and --json are looked for before the flags are
// parsed — see peekLang — and a wrapped `ls --lang` must not be mistaken for
// a request to speak Latvian.
func ourArgs(args []string) []string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" || !strings.HasPrefix(a, "-") {
			return args[:i]
		}
		if takesValue(a) {
			i++
		}
	}
	return args
}

// shellFor decides whether what was typed is a line for a shell rather than
// the name of a program.
//
// `digitdisk -c 'make && make test'` is what a person writes by habit, and -c
// in every shell means «a line to interpret».  One argument holding spaces or
// metacharacters is that line and is handed to $SHELL; anything else is a
// program and its arguments, and is started directly, without a shell
// standing between the wrapper and the thing being measured.
func shellFor(args []string) string {
	if len(args) != 1 || !strings.ContainsAny(args[0], " \t\n|&;<>()$`\"'*?") {
		return ""
	}
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	return "/bin/sh"
}

// ЯЗЫК ВЫВОДА
//
// The language is settled once per run and then handed down: every report
// takes it as an argument, and nothing below main reads the environment or a
// settings file for it.  That is why `--lang en` and a saved choice and a
// locale all end in the same place, and why a test can run the whole tool in
// either language without touching the machine it runs on.
//
// The order is: --lang, then DIGITDISK_LANG, then the answer stored in
// ~/.digitable/digitdisk/settings.conf, then — if a person is sitting at a
// terminal and has never answered — the question, then the locale, then
// English.  internal/settings holds the whole of it; this is the wiring.
var (
	langOnce   sync.Once
	langChoice settings.Choice
)

// langFlagName is the key, spelled once: it is registered on every подкоманда,
// named in the справка and named on both страницы руководства, and the test in
// main_test.go compares those three lists.
const langFlagName = "lang"

// chooseLang settles the language of this run.
//
// mayAsk is the caller's promise that a question is safe here: it is false for
// `--json`, which is read by a script that would hang on a prompt, and it is
// false when either end of the conversation is not a terminal.  A run that may
// not ask never writes anything into the home directory either — see
// settings.Decide.
func chooseLang(flagValue string, mayAsk bool) lang.Lang {
	langOnce.Do(func() {
		ask := settings.Ask{}
		if mayAsk && ui.IsInteractive(os.Stdin) && ui.IsInteractive(os.Stderr) {
			ask = settings.Ask{In: os.Stdin, Out: os.Stderr, May: true}
		}
		langChoice = settings.Decide(settings.Options{}, flagValue, ask)
		for _, note := range langChoice.Notes {
			fmt.Fprintf(os.Stderr, "digitdisk: %s\n", note.In(langChoice.Lang))
		}
	})
	return langChoice.Lang
}

// langOnly settles the language without any chance of asking.  It answers
// --help, --version and the refusal of a word that is not a подкоманда: three
// answers that must not depend on whether somebody is watching.
func langOnly(args []string) lang.Lang {
	return chooseLang(peekLang(args), false)
}

// peekLang finds --lang on a command line before the flag package has looked
// at it.  --help and an unknown подкоманда are answered before any подкоманда
// parses anything, and they still have to come out in the language asked for.
func peekLang(args []string) string {
	for i, a := range args {
		switch {
		case a == "--"+langFlagName || a == "-"+langFlagName:
			if i+1 < len(args) {
				return args[i+1]
			}
		case strings.HasPrefix(a, "--"+langFlagName+"="):
			return strings.TrimPrefix(a, "--"+langFlagName+"=")
		case strings.HasPrefix(a, "-"+langFlagName+"="):
			return strings.TrimPrefix(a, "-"+langFlagName+"=")
		}
	}
	return ""
}

// rememberLang stores a language the reader switched to on the live screen.
// The screen is handed this rather than internal/settings itself: what a home
// directory is has nothing to do with drawing a terminal, and a screen that
// could not write a file would be a screen that could not be tested without
// one.
func rememberLang(l lang.Lang) lang.Phrase {
	langChoice.Lang = l
	return settings.Remember(settings.Options{}, l)
}

// peekJSON reports whether this command line asks for machine-readable output.
//
// It is read before the flags are parsed for one reason: a run that will be
// parsed by a script must never be interrupted by a question, and the question
// has to be settled before anything is printed at all.  A spelling this misses
// costs nothing — the question is asked only when both ends of the
// conversation are a terminal, which a script's are not.
func peekJSON(args []string) bool {
	for _, a := range args {
		if a == "--json" || a == "-json" ||
			strings.HasPrefix(a, "--json=") || strings.HasPrefix(a, "-json=") {
			return true
		}
	}
	return false
}

// stripGlobalLang takes --lang off the front of the command line and gives
// back the rest.
//
// It stops at the first argument that is not the language key: a --lang after
// the подкоманда belongs to the подкоманда's own flag set, which registers it,
// and taking it here as well would let `digitdisk clean --lang` mean two
// things at once.
func stripGlobalLang(args []string) (rest []string, value string) {
	for len(args) > 0 {
		a := args[0]
		switch {
		case a == "--"+langFlagName || a == "-"+langFlagName:
			if len(args) < 2 {
				return args, value
			}
			value, args = args[1], args[2:]
		case strings.HasPrefix(a, "--"+langFlagName+"="):
			value, args = strings.TrimPrefix(a, "--"+langFlagName+"="), args[1:]
		case strings.HasPrefix(a, "-"+langFlagName+"="):
			value, args = strings.TrimPrefix(a, "-"+langFlagName+"="), args[1:]
		default:
			return args, value
		}
	}
	return args, value
}

// sayMoved tells a person, once, that a settings file was read from the home
// it used to live in.
//
// Once and not every run: a notice repeated at every invocation stops being
// read, and this one has nothing new to say after the first time.  The mark
// that it has been said lives in settings.conf beside the language — and
// writing it is itself announced, because writing in somebody's home directory
// is an action and this tool exists to clean up after programs that do it
// quietly.  If the mark cannot be written, the notice comes back next time,
// which is the honest failure: nothing was remembered, so nothing is pretended
// to have been.
func sayMoved(l lang.Lang, notes ...lang.Phrase) {
	var say []lang.Phrase
	for _, n := range notes {
		if !n.Empty() {
			say = append(say, n)
		}
	}
	if len(say) == 0 {
		return
	}
	o := settings.Options{}
	stored, err := settings.Load(o)
	if err == nil && stored.MoveAnnounced {
		return
	}
	for _, n := range say {
		fmt.Fprintf(os.Stderr, "digitdisk: %s\n", n.In(l))
	}
	if err != nil {
		return
	}
	stored.MoveAnnounced = true
	if path, err := settings.Save(o, stored); err == nil {
		fmt.Fprintf(os.Stderr, "digitdisk: %s\n", lang.Say("сказано один раз, отметка в %s", path).In(l))
	}
}
