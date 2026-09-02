// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package run starts somebody else's command, watches what it costs while it
// works, and says what it cost when it is over.
//
// # The command is the point, not us
//
// Everything here is arranged around one promise: a command under the wrapper
// behaves exactly as it behaves without it.  Three consequences, and they are
// the reason for most of the code below.
//
//   - Nothing of ours is ever written to standard output.  The child is
//     handed our own descriptors — not pipes — so `digitdisk run make | tee
//     log` writes the bytes `make | tee log` writes, in the same order, with
//     nothing added.  The строка состояния goes to the terminal behind
//     standard ERROR, and only when there is one; the сводка goes there too.
//   - The код возврата is the command's own, and a command killed by a
//     signal is reported as killed rather than as an exit code somebody
//     invented.  Signals from the terminal reach the command by themselves:
//     it is left in our process group on purpose, so Ctrl-C goes to it and
//     not only to us, and the wrapper absorbs its own copy instead of dying
//     and leaving an orphan behind.
//   - The measurement never gets in the way.  A замер that took too long
//     backs itself off (see period), and a command that puts the terminal
//     into raw mode — vim, ssh, less, anything full-screen — takes the
//     строка состояния down with it until it gives the terminal back.
//
// # What is measured, and how honestly
//
// A build is a tree of processes, not one process, and measuring the direct
// child alone would answer «almost nothing».  Three ways to count a tree, in
// the order they are preferred:
//
//	контрольная группа  cgroup v2, own group: exact and cheap — the kernel
//	                    counts CPU and peak memory for everything inside,
//	                    including what died between two замера.
//	обход /proc         the tree is followed by parent links: approximate,
//	                    and the approximation is named in the сводка.
//	итог ядра           wait4(2) rusage: exact CPU time and exact peak of the
//	                    largest single process, but nothing live.  Every
//	                    platform has it, and it is the floor under the other
//	                    two.
//
// The сводка always says which of the three answered, because «6 ГиБ пик» and
// «около 6 ГиБ пик» are different claims and the reader is entitled to know
// which one they are holding.
package run

import (
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"digitdisk/internal/lang"
)

// Options is one wrapped command.
type Options struct {
	// Args is what to run: a program and its arguments.  When Shell is set,
	// Args holds exactly one string and it is a line for that shell.
	Args  []string
	Shell string

	// In, Out and Err are handed to the command as they are.  Nothing is
	// copied through us: a pipe of our own would change how the command
	// buffers its output, and that is a change to its behaviour.
	In, Out, Err *os.File

	// Interval is how often the строка состояния is redrawn.  It is not how
	// often the tree is read — see period.
	Interval time.Duration

	// Plain takes the строка состояния away even in a terminal.
	Plain bool

	// GPUTool allows asking the driver's own program about video memory.
	// Nothing runs it unless this is set: see gpu.go.
	GPUTool bool

	Lang lang.Lang

	// Now is the clock, for tests.
	Now func() time.Time
}

// Result is what the command cost.  The JSON is the same in either language:
// keys and machine values are bytes, not text.
type Result struct {
	Command      string  `json:"команда"`
	Shell        string  `json:"оболочка,omitempty"`
	Code         int     `json:"код"`
	Signal       string  `json:"сигнал,omitempty"`
	SignalNumber int     `json:"сигнал_номер,omitempty"`
	Seconds      float64 `json:"секунд"`
	CPUSeconds   float64 `json:"процессорных_секунд"`
	CPUExact     bool    `json:"процессорное_время_точно"`
	CPUPercent   float64 `json:"процессор_средний_процент"`
	PeakBytes    uint64  `json:"пик_памяти_байт"`
	PeakExact    bool    `json:"пик_памяти_точно"`
	PeakOneBytes uint64  `json:"пик_одного_процесса_байт"`
	Processes    int     `json:"процессов_в_пике"`
	Seen         int     `json:"процессов_видено"`
	Accounting   string  `json:"учёт"`
	SampleMS     int64   `json:"опрос_мс"`
	GPU          *GPU    `json:"видеокарта,omitempty"`
}

// period is how often the tree is read.  It is shorter than the redraw so a
// пик that lasts a second is not missed between two frames, and it is not
// shorter than this because reading a tree costs something on a machine with
// four thousand processes.
const period = 200 * time.Millisecond

// backoff is the share of a замер's own cost the wrapper is willing to be.  A
// tree that takes longer than this to read stretches the замер until it fits
// — an обёртка that measured a build by slowing it down would be measuring
// itself.
const backoff = 20

// maxPeriod is as far as the backoff will stretch.
const maxPeriod = 2 * time.Second

// Run starts the command and returns when it is over.  The error is returned
// only when the command could not be started at all; a command that ran and
// failed is not an error of ours, it is a Result with its own код.
func Run(o Options) (Result, error) {
	if len(o.Args) == 0 || strings.TrimSpace(o.Args[0]) == "" {
		return Result{}, lang.Errorf("нужна команда: digitdisk run <команда> [доводы]")
	}
	if o.Interval <= 0 {
		o.Interval = time.Second
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	l := o.Lang

	// Both of these happen BEFORE the command starts, and both have to.
	// The учёт has to exist before there is a process to account for — a
	// control group cannot adopt a process that has already forked twice —
	// and the строка состояния has to reserve its row while the cursor is
	// still where the shell left it: making room later, in the middle of a
	// half-written line of the command's own output, is what turns
	// «Compiling…» into « done».
	m := newMeter()
	defer m.close()
	b := newBar(o.Err, o.Plain, colourOK())
	defer b.close()

	cmd, err := start(o, m)
	if err != nil {
		return Result{}, err
	}

	res := Result{
		Command: shown(o.Args),
		Shell:   o.Shell,
	}

	g := newGPU(o.GPUTool)
	begin := o.Now()
	m.started(cmd.Process.Pid, begin)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	sigs := make(chan os.Signal, 8)
	notify(sigs)
	defer signal.Stop(sigs)

	every := period
	sample := time.NewTicker(every)
	defer sample.Stop()
	redraw := time.NewTicker(o.Interval)
	defer redraw.Stop()

	var last reading
	waiting := true
	for waiting {
		select {
		case err = <-done:
			waiting = false
		case <-sample.C:
			t0 := o.Now()
			last = m.sample(t0)
			g.sample(t0, m.members())
			// A замер that costs more than its share stretches itself.
			// The сводка then says how often it actually looked.
			if cost := o.Now().Sub(t0); cost*backoff > every {
				if every = cost * backoff; every > maxPeriod {
					every = maxPeriod
				}
				sample.Reset(every)
			}
		case <-redraw.C:
			b.draw(barText(l, b.width(), o.Now().Sub(begin), last, g.reading()))
		case s := <-sigs:
			handle(cmd, s, b)
		}
	}

	b.close()
	end := o.Now()
	t := m.finish(cmd.ProcessState)

	res.Seconds = end.Sub(begin).Seconds()
	res.CPUSeconds = t.CPU.Seconds()
	res.CPUExact = t.CPUExact
	if res.Seconds > 0 {
		res.CPUPercent = res.CPUSeconds / res.Seconds * 100
	}
	res.PeakBytes = t.Peak
	res.PeakExact = t.PeakExact
	res.PeakOneBytes = t.PeakOne
	res.Processes = t.Procs
	res.Seen = t.Seen
	res.Accounting = t.How
	res.SampleMS = every.Milliseconds()
	res.GPU = g.total()
	res.Code, res.Signal, res.SignalNumber = outcome(cmd.ProcessState, err)
	return res, nil
}

// start runs the command, and runs it a second time without the контрольная
// группа if that is what the kernel refused.
//
// A cgroup is a measuring instrument, and an instrument that stops the thing
// it measures from starting is worse than no instrument: the fallback is
// silent by design, and the сводка names the учёт that actually answered.
func start(o Options, m *meter) (*exec.Cmd, error) {
	cmd := command(o)
	cmd.SysProcAttr = m.attr()
	err := cmd.Start()
	if err == nil || cmd.SysProcAttr == nil {
		return cmd, wrapStart(o, err)
	}
	m.dropGroup()
	cmd = command(o)
	cmd.SysProcAttr = m.attr()
	return cmd, wrapStart(o, cmd.Start())
}

// shown is the command as it is written back to the person who typed it.  An
// argument holding spaces gets its quotes back: «sh -c exit 3» and «sh -c
// \"exit 3\"» are different command lines, and the сводка names the one that
// ran.
func shown(args []string) string {
	parts := make([]string, len(args))
	for i, a := range args {
		if strings.ContainsAny(a, " \t\n") {
			parts[i] = "\"" + a + "\""
		} else {
			parts[i] = a
		}
	}
	return strings.Join(parts, " ")
}

func command(o Options) *exec.Cmd {
	name, argv := o.Args[0], o.Args[1:]
	if o.Shell != "" {
		name, argv = o.Shell, []string{"-c", o.Args[0]}
	}
	cmd := exec.Command(name, argv...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = o.In, o.Out, o.Err
	return cmd
}

// wrapStart turns «cannot start» into a refusal a person can read, keeping
// the two cases a shell keeps apart: a name that is not there and a file that
// will not run.  The codes go with them — 127 and 126 are what every shell
// answers, and a script that already tests for them must not have to learn a
// new number because the command was measured.
func wrapStart(o Options, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, exec.ErrNotFound) || errors.Is(err, syscall.ENOENT) {
		return &startError{code: 127, err: lang.Errorf("команда %q не найдена", o.Args[0])}
	}
	return &startError{code: 126, err: lang.Errorf("команда %q не запускается: %s", o.Args[0], err)}
}

// A startError is a command that never ran.  It carries the number a shell
// would have answered with, and wraps the refusal itself, so the message and
// the код возврата travel together and neither is invented at the far end.
type startError struct {
	code int
	err  error
}

func (e *startError) Error() string { return e.err.Error() }

func (e *startError) Unwrap() error { return e.err }

// StartCode is the код возврата for a command that could not be started: 127
// when there is no such command, 126 when there is and it will not run.
// Anything else is an ordinary refusal and answers 1.
func StartCode(err error) int {
	var s *startError
	if errors.As(err, &s) {
		return s.code
	}
	return 1
}

// outcome reads the ending out of the wait status: an exit code, or the
// signal that killed it.  Nothing is invented here — a command killed by
// SIGKILL has no exit code of its own, and 137 is what a shell reports, not
// what the command said.
func outcome(ps *os.ProcessState, err error) (code int, name string, number int) {
	if ps == nil {
		if err != nil {
			return 1, "", 0
		}
		return 0, "", 0
	}
	if ws, ok := ps.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		s := ws.Signal()
		return 128 + int(s), signalName(s), int(s)
	}
	return ps.ExitCode(), "", 0
}

// notify subscribes to the signals the wrapper has to think about.
//
// SIGINT and SIGQUIT are subscribed to in order to be IGNORED: the terminal
// sends them to the whole foreground process group, so the command already
// has its copy, and a wrapper that died on its own copy would leave the
// command running with nobody waiting for it.  SIGTERM is the one a person
// sends by pid, and it is the one forwarded.
func notify(ch chan<- os.Signal) {
	signal.Notify(ch,
		syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGWINCH, syscall.SIGTSTP, syscall.SIGCONT)
}

// handle acts on one signal.
func handle(cmd *exec.Cmd, s os.Signal, b *bar) {
	switch s {
	case syscall.SIGTERM:
		if cmd.Process != nil {
			_ = cmd.Process.Signal(s)
		}
	case syscall.SIGWINCH:
		b.resize()
	case syscall.SIGTSTP:
		// Ctrl-Z: the terminal stopped the whole group, and this half of
		// the group has to stop by hand — subscribing to a signal takes
		// its default away.  The строка состояния is taken down first: a
		// shell prompt inside our scroll region, with a stale bar frozen
		// under it, is a mess somebody else has to reset.
		//
		// The stop is SIGSTOP and not SIGTSTP, and that is not a detail:
		// a Go program cannot stop itself with SIGTSTP once anything has
		// asked to be notified of it — the runtime installs a handler,
		// signal.Reset does not put the default back in time, and the
		// signal is swallowed.  Checked, not assumed: with SIGTSTP the
		// process goes on running while the command under it is stopped.
		// SIGSTOP nobody can catch, which is exactly why it works.
		b.close()
		_ = syscall.Kill(os.Getpid(), syscall.SIGSTOP)
	case syscall.SIGCONT:
		b.wake()
	}
}
