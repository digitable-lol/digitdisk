// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package sysinfo

// The hardware half of the snapshot: the processors one by one, the video
// cards, and the few facts about the machine a person recognises it by.
//
// Everything here is shared between the systems.  What differs — which file
// or which kernel call the numbers come from — lives in collect_linux.go and
// collect_darwin.go; what is the same is the shape of the answer and the
// arithmetic on it, and it is written once so the two cannot drift apart.

import (
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"

	"digitdisk/internal/lang"
	"digitdisk/internal/procfs"
)

// Keys of Status.Missing for the facts this file is about.  They follow the
// rule the others follow: the name is what a reader sees, and the reason is
// printed by `--why` and nowhere else.
const (
	FactCores      = "загрузка по ядрам"
	FactGPUs       = "видеокарты"
	FactGPUNumbers = "загрузка и память видеокарт"
	FactGPUPower   = "мощность видеокарт"
	FactMachine    = "модель машины"
	FactCPUModel   = "модель процессора"
	FactShell      = "оболочка"
	FactDesktop    = "рабочий стол"
)

// Core is one processor of the machine over the sample window.
//
// The share is a pointer for the same reason the machine-wide one is: a
// processor that was taken offline between the two readings has no share, and
// no share is not a share of zero.
type Core struct {
	Index       int      `json:"index"`
	Name        string   `json:"name"`
	BusyPercent *float64 `json:"busy_percent"`
}

// coresFrom turns two readings of the per-processor counters into the shares.
// The two are matched by name rather than by position: a processor that goes
// offline between the readings leaves a hole in the list, and matching by
// position would then shift every core after it by one and report the whole
// machine wrongly.
func coresFrom(before, after []procfs.CPUTimes) []Core {
	first := make(map[string]procfs.CPUTimes, len(before))
	for _, c := range before {
		first[c.Name] = c
	}
	out := make([]Core, 0, len(after))
	for _, c := range after {
		core := Core{Name: c.Name, Index: coreIndex(c.Name)}
		if b, ok := first[c.Name]; ok {
			total := float64(c.Total()) - float64(b.Total())
			busy := float64(c.Busy()) - float64(b.Busy())
			if total > 0 && busy >= 0 {
				share := 100 * busy / total
				if share > 100 {
					share = 100
				}
				core.BusyPercent = &share
			}
		}
		out = append(out, core)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Index < out[j].Index })
	return out
}

func coreIndex(name string) int {
	n, err := strconv.Atoi(strings.TrimPrefix(name, "cpu"))
	if err != nil {
		return 1 << 30
	}
	return n
}

// coresAgree is the self-check on the per-processor shares: their mean has to
// be the share the machine reported for itself, because the machine-wide
// counter is the sum of the per-processor ones.  A list that fails it is a
// list read wrongly — a shifted index, a misparsed line — and it is not
// published.
//
// The tolerance is five points of percentage.  The two readings of the
// counters are taken one after the other rather than at the same instant, so
// they never agree exactly; five points is wider than that gap and far
// narrower than any error worth catching.
func coresAgree(cores []Core, whole *float64) (bool, lang.Phrase) {
	if len(cores) == 0 {
		return false, lang.Say("ядро не публикует счётчики по каждому процессору")
	}
	sum, n := 0.0, 0
	for _, c := range cores {
		if c.BusyPercent != nil {
			sum += *c.BusyPercent
			n++
		}
	}
	if n == 0 {
		return false, lang.Say("за окно замера счётчики процессоров не сдвинулись")
	}
	if whole == nil {
		return true, lang.Phrase{}
	}
	if diff := sum/float64(n) - *whole; diff > 5 || diff < -5 {
		return false, lang.Say("среднее по ядрам разошлось с общей загрузкой машины — список ядер не публикуем")
	}
	return true, lang.Phrase{}
}

// CoreLoad is the digest of the per-core shares: what a reader can take in at
// a glance about a machine with two hundred and fifty-six of them.
type CoreLoad struct {
	Total    int
	Measured int
	Min      float64
	Median   float64
	Max      float64
	// Busiest is the index of the processor with the largest share.
	Busiest int
	// Loaded is how many processors are busy more than half the time.
	Loaded int
}

// Cores digests the per-core shares of this snapshot.  ok is false when there
// are none, and the caller then prints a dash.
func (s Status) Cores() (CoreLoad, bool) {
	shares := make([]float64, 0, len(s.Load.Cores))
	out := CoreLoad{Total: len(s.Load.Cores), Busiest: -1}
	best := -1.0
	for _, c := range s.Load.Cores {
		if c.BusyPercent == nil {
			continue
		}
		shares = append(shares, *c.BusyPercent)
		if *c.BusyPercent > best {
			best, out.Busiest = *c.BusyPercent, c.Index
		}
		if *c.BusyPercent >= 50 {
			out.Loaded++
		}
	}
	if len(shares) == 0 {
		return CoreLoad{}, false
	}
	sort.Float64s(shares)
	out.Measured = len(shares)
	out.Min, out.Max = shares[0], shares[len(shares)-1]
	out.Median = shares[len(shares)/2]
	if len(shares)%2 == 0 {
		out.Median = (shares[len(shares)/2-1] + shares[len(shares)/2]) / 2
	}
	return out, true
}

// environment fills in the few things about the session that a person expects
// to see next to the machine's name — who is logged in, which shell, which
// desktop.  They are not measurements of the machine: they are what the
// system handed this process, and where it handed nothing the field stays
// empty and the reason goes to `--why`.
func environment(st *Status) {
	h := &st.Host
	h.Bits = strconv.IntSize
	h.Terminal = os.Getenv("TERM")
	if u, err := user.Current(); err == nil {
		h.User = u.Username
	} else {
		h.User = os.Getenv("USER")
	}
	h.Shell = os.Getenv("SHELL")
	if h.Shell == "" {
		st.Missing[FactShell] = lang.Say("переменная окружения SHELL пуста — оболочку назвать нечем")
	}
	for _, name := range []string{"XDG_CURRENT_DESKTOP", "DESKTOP_SESSION", "XDG_SESSION_DESKTOP"} {
		if v := strings.TrimSpace(os.Getenv(name)); v != "" {
			h.Desktop = v
			break
		}
	}
	if h.Desktop == "" {
		st.Missing[FactDesktop] = lang.Say("переменные окружения XDG_CURRENT_DESKTOP и DESKTOP_SESSION пусты — рабочего стола в этом сеансе нет")
	}
}
