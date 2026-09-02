// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package sysinfo

import (
	"digitdisk/internal/lang"
	"testing"

	"digitdisk/internal/procfs"
)

// ticks builds one /proc/stat line's worth of counters.
func ticks(name string, user, system, idle uint64) procfs.CPUTimes {
	return procfs.CPUTimes{Name: name, User: user, System: system, Idle: idle}
}

func TestCoresFromMatchesByNameNotByPosition(t *testing.T) {
	before := []procfs.CPUTimes{ticks("cpu0", 100, 0, 100), ticks("cpu1", 0, 0, 200), ticks("cpu2", 50, 0, 150)}
	// cpu1 went offline between the readings: the kernel simply stops
	// printing its line.  Matching by position would then report cpu2's
	// work as cpu1's.
	after := []procfs.CPUTimes{ticks("cpu0", 200, 0, 100), ticks("cpu2", 50, 0, 250)}

	cores := coresFrom(before, after)
	if len(cores) != 2 {
		t.Fatalf("ядер %d, ждали 2", len(cores))
	}
	if cores[0].Index != 0 || cores[1].Index != 2 {
		t.Errorf("ядра пронумерованы %d и %d", cores[0].Index, cores[1].Index)
	}
	if cores[0].BusyPercent == nil || *cores[0].BusyPercent != 100 {
		t.Errorf("cpu0 занято %v, ждали 100%%", cores[0].BusyPercent)
	}
	if cores[1].BusyPercent == nil || *cores[1].BusyPercent != 0 {
		t.Errorf("cpu2 занято %v, ждали 0%%", cores[1].BusyPercent)
	}
}

// A processor that appears only in the second reading has no interval, and no
// interval is not a share of zero.
func TestNewCoreHasNoShare(t *testing.T) {
	cores := coresFrom(nil, []procfs.CPUTimes{ticks("cpu0", 10, 0, 10)})
	if len(cores) != 1 {
		t.Fatalf("ядер %d", len(cores))
	}
	if cores[0].BusyPercent != nil {
		t.Errorf("у нового ядра появилась доля %v", *cores[0].BusyPercent)
	}
}

func TestCoresAgreeRefusesAListThatContradictsTheMachine(t *testing.T) {
	half := 50.0
	cores := []Core{{Index: 0, BusyPercent: &half}, {Index: 1, BusyPercent: &half}}
	if ok, _ := coresAgree(cores, &half); !ok {
		t.Error("совпадающие среднее и общая доля отвергнуты")
	}
	whole := 5.0
	ok, why := coresAgree(cores, &whole)
	if ok {
		t.Error("среднее по ядрам 50%% принято при общей доле 5%%")
	}
	if why.Empty() {
		t.Error("отказ не объяснён")
	}
	// Nothing measured at all is a refusal too, and a different one.
	if ok, why := coresAgree([]Core{{Index: 0}}, &whole); ok || why.Empty() {
		t.Errorf("пустой замер принят (ok=%v, why=%q)", ok, why)
	}
	if ok, why := coresAgree(nil, nil); ok || why.Empty() {
		t.Errorf("пустой список принят (ok=%v, why=%q)", ok, why)
	}
	// A list without a machine-wide share to check against is still usable:
	// there is simply nothing to contradict.
	if ok, _ := coresAgree(cores, nil); !ok {
		t.Error("без общей доли список ядер отвергнут")
	}
}

func TestCoresDigest(t *testing.T) {
	var st Status
	for i, v := range []float64{0, 10, 60, 100} {
		share := v
		st.Load.Cores = append(st.Load.Cores, Core{Index: i, BusyPercent: &share})
	}
	// One processor nobody measured must not drag the digest towards zero.
	st.Load.Cores = append(st.Load.Cores, Core{Index: 4})

	got, ok := st.Cores()
	if !ok {
		t.Fatal("сводка по ядрам не составилась")
	}
	if got.Total != 5 || got.Measured != 4 {
		t.Errorf("всего %d, замерено %d", got.Total, got.Measured)
	}
	if got.Min != 0 || got.Max != 100 {
		t.Errorf("мин %v, макс %v", got.Min, got.Max)
	}
	if got.Median != 35 {
		t.Errorf("медиана %v, ждали 35 (среднее двух средних значений)", got.Median)
	}
	if got.Busiest != 3 {
		t.Errorf("самое занятое ядро %d, ждали 3", got.Busiest)
	}
	if got.Loaded != 2 {
		t.Errorf("занятых больше половины %d, ждали 2", got.Loaded)
	}

	if _, ok := (Status{}).Cores(); ok {
		t.Error("сводка составилась из пустого списка ядер")
	}
}

func TestEnvironmentNamesWhatTheSessionDidNotGive(t *testing.T) {
	t.Setenv("SHELL", "")
	t.Setenv("XDG_CURRENT_DESKTOP", "")
	t.Setenv("DESKTOP_SESSION", "")
	t.Setenv("XDG_SESSION_DESKTOP", "")
	st := Status{Missing: map[string]lang.Phrase{}}
	environment(&st)
	if st.Host.Bits != 64 && st.Host.Bits != 32 {
		t.Errorf("разрядность %d", st.Host.Bits)
	}
	for _, fact := range []string{FactShell, FactDesktop} {
		if why, ok := st.Unmeasured(fact); !ok || why.Empty() {
			t.Errorf("пустое окружение не объяснено для %q", fact)
		}
	}
	// And what is missing here belongs to --json and to --why, not to the
	// short line at the foot of the printed report.
	for _, name := range st.UnmeasuredNames() {
		if name == FactShell || name == FactDesktop {
			t.Errorf("%q попал в печатный отчёт, где для него нет колонки", name)
		}
	}

	t.Setenv("SHELL", "/bin/zsh")
	t.Setenv("XDG_CURRENT_DESKTOP", "GNOME")
	st = Status{Missing: map[string]lang.Phrase{}}
	environment(&st)
	if st.Host.Shell != "/bin/zsh" || st.Host.Desktop != "GNOME" {
		t.Errorf("оболочка %q, рабочий стол %q", st.Host.Shell, st.Host.Desktop)
	}
	if len(st.Missing) != 0 {
		t.Errorf("заполненное окружение объявлено неизмеренным: %v", st.Missing)
	}
}
