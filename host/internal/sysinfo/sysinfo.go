// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package sysinfo

import (
	"fmt"
	"time"
)

// Keys of Status.Missing for facts a system does not publish at all, as
// opposed to sources that failed to read.  They are constants because both
// sides need the same string: the collector writes the reason, and the report
// looks it up to print "—" and that reason instead of a number.
//
// Linux publishes all of these and sets none of them.  macOS sets several,
// which is the honest half of porting: the fields exist in the snapshot
// because the shape is shared, and they say out loud that nobody measured
// them.
const (
	FactCPUBusy      = "загрузка ЦП"
	FactLoadEntities = "очередь планировщика"
	FactThreads      = "потоки процессов"
	FactBlocked      = "заблокированные процессы"
	FactProcessCPU   = "процессорное время процессов"
	FactProcessRSS   = "память процессов"
	FactProcessArgs  = "командные строки процессов"
	FactNetCounters  = "счётчики интерфейсов"
	FactNetTxDrops   = "отброшенные исходящие пакеты"
	FactSensors      = "датчики температуры"
	FactMemoryPages  = "разбивка памяти по страницам"
)

// Unmeasured reports whether the named fact is absent from this snapshot, and
// why.  An empty reason with ok = true is still an absence: the caller prints
// a dash either way.
func (s Status) Unmeasured(fact string) (why string, ok bool) {
	why, ok = s.Missing[fact]
	return why, ok
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
