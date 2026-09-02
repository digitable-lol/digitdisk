// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package sysinfo

import (
	"fmt"
	"sort"
	"time"

	"digitdisk/internal/lang"
)

// Keys of Status.Missing for facts a system does not publish, as opposed to
// sources that failed to read.  They are constants because both sides need the
// same string: the collector writes the reason, and the report looks it up to
// print "—" instead of a number.
//
// The name is what a reader sees; the reason is one short phrase, and it is
// not printed beside the number.  A report is a place for measurements, and a
// paragraph about why a measurement is missing is not one — so the reasons
// live behind `digitdisk status --why` and in the JSON, where somebody who
// wants them can ask.
const (
	FactCPUBusy      = "загрузка ЦП"
	FactLoadEntities = "очередь планировщика"
	FactThreads      = "потоки процессов"
	FactBlocked      = "заблокированные процессы"
	FactProcessCPU   = "процессорное время процессов"
	FactProcessRSS   = "память процессов"
	FactRunning      = "выполняющиеся процессы"
	FactProcessArgs  = "командные строки процессов"
	FactOtherUsers   = "процессы других пользователей"
	FactNetCounters  = "счётчики интерфейсов"
	FactNetTxDrops   = "отброшенные исходящие пакеты"
	FactSensors      = "датчики температуры"
	FactMemoryPages  = "разбивка памяти"
)

// jsonOnly names the facts that have a field in the JSON but no column in the
// printed report.  Naming them in the report would be noise about numbers the
// reader was never shown; dropping them from the snapshot would let a zero in
// the JSON pass for a measurement.  So they stay in Status.Missing and out of
// the short line.
var jsonOnly = map[string]bool{
	FactLoadEntities: true,
	FactNetTxDrops:   true,
	// The session's shell and desktop are drawn on the live screen beside
	// the mark of the system, and the printed report has no line for
	// either.  A headless server has no desktop, and naming that at the
	// foot of a report about disks and memory would be noise.
	FactShell:   true,
	FactDesktop: true,
}

// Unmeasured reports whether the named fact is absent from this snapshot, and
// why.  An empty reason with ok = true is still an absence: the caller prints
// a dash either way.
func (s Status) Unmeasured(fact string) (why lang.Phrase, ok bool) {
	why, ok = s.Missing[fact]
	return why, ok
}

// UnmeasuredNames lists, in order, what the report should say is missing: the
// names only, and only the ones that stand for something the report would
// otherwise have printed.
func (s Status) UnmeasuredNames() []string {
	out := make([]string, 0, len(s.Missing))
	for name := range s.Missing {
		if !jsonOnly[name] {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// An Absence is one fact the snapshot does not hold, with the reason.
//
// The name is a вокабула — the same Russian word that keys Status.Missing in
// the JSON, so a script reading the snapshot sees exactly what it always saw —
// and the screen shows the reader's word for it.  The reason is a Phrase for
// the same purpose: Russian in the file, the reader's language on the screen.
type Absence struct {
	Name string
	Why  lang.Phrase
}

// UnmeasuredAll lists every absence with its reason, sorted by name.  This is
// what `--why` prints and nothing else does.
func (s Status) UnmeasuredAll() []Absence {
	names := make([]string, 0, len(s.Missing))
	for name := range s.Missing {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]Absence, 0, len(names))
	for _, name := range names {
		out = append(out, Absence{Name: name, Why: s.Missing[name]})
	}
	return out
}

// HumanDuration renders a number of seconds as "5д 03:14".
func HumanDuration(seconds float64) string {
	if seconds <= 0 {
		return ""
	}
	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dд %02d:%02d", days, hours, mins)
	}
	return fmt.Sprintf("%02d:%02d", hours, mins)
}

// cpuOrZero sorts unmeasured processes below every measured one.
func cpuOrZero(p Proc) float64 {
	if p.CPUPercent == nil {
		return -1
	}
	return *p.CPUPercent
}

func head(p []Proc, n int) []Proc {
	if n <= 0 {
		return nil
	}
	if len(p) > n {
		p = p[:n]
	}
	return p
}
