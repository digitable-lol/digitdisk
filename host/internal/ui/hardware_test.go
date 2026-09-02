// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package ui

import (
	"fmt"
	"strings"
	"testing"

	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/sysinfo"
)

// withCores is the filled snapshot plus n processors, the busiest of which is
// the one named.
func withCores(n int, busiest int) sysinfo.Status {
	st := filled()
	st.Load.CPUCount = n
	st.Load.Cores = make([]sysinfo.Core, 0, n)
	for i := 0; i < n; i++ {
		share := 3.0
		if i == busiest {
			share = 100
		}
		v := share
		st.Load.Cores = append(st.Load.Cores, sysinfo.Core{Index: i, Name: fmt.Sprintf("cpu%d", i), BusyPercent: &v})
	}
	return st
}

// The marks are drawn here, and two things about them are load-bearing: they
// are printable ASCII, so a terminal without our fonts still draws them, and
// they are all one width, so the fields beside them keep their column.
func TestEmblemsAreASCIIAndOfOneWidth(t *testing.T) {
	if len(emblems) == 0 {
		t.Fatal("значков нет вовсе")
	}
	for name, e := range emblems {
		if len(e.art) == 0 {
			t.Errorf("значок %q пуст", name)
			continue
		}
		width := emblemWidth(e.art)
		if width > emblemMaxWidth {
			t.Errorf("значок %q шириной %d — шире предела %d", name, width, emblemMaxWidth)
		}
		for i, line := range e.art {
			if runes(line) != width {
				t.Errorf("значок %q: строка %d шириной %d, а первая %d", name, i, runes(line), width)
			}
			for _, r := range line {
				if r > 126 || r < 32 {
					t.Errorf("значок %q содержит знак %q вне печатной латиницы — на LANG=C он развалится", name, r)
				}
			}
		}
	}
}

func TestEmblemForKnowsFamiliesAndFallsBack(t *testing.T) {
	// A rebuild gets the mark of what it is a rebuild of.
	if got := emblemFor("rocky", "Rocky Linux 9"); emblemWidth(got.art) == 0 {
		t.Error("у rocky нет значка")
	}
	// A system nobody has drawn still gets a mark rather than a hole.
	unknown := emblemFor("своя-система", "Своя система 1.0")
	if len(unknown.art) == 0 {
		t.Fatal("незнакомый дистрибутив остался без значка")
	}
	// And the general mark is exactly the general one.
	if unknown.art[0] != emblems[""].art[0] {
		t.Error("незнакомому дистрибутиву достался не общий значок")
	}
	// The pretty name is the second chance: /etc/os-release without ID=.
	if got := emblemFor("", "Ubuntu 26.04 LTS"); got.art[0] == emblems[""].art[0] {
		t.Error("Ubuntu опознана только по ID, хотя её видно и по имени")
	}
}

func TestSystemPagePutsTheMarkBesideTheFactsAndStacksWhenNarrow(t *testing.T) {
	st := filled()
	st.Host.User = "гость"
	st.Host.CPUModel = "AMD EPYC 7742 64-Core Processor"
	st.Host.Model = "Dell Inc. PowerEdge C6525"

	wide := plain(strings.Join(newTestScreen(st, true, 120).system(), "\n"))
	if !strings.Contains(wide, "гость@узел") {
		t.Errorf("узел напечатан не как пользователь@машина:\n%s", wide)
	}
	if !strings.Contains(wide, "AMD EPYC") || !strings.Contains(wide, "PowerEdge") {
		t.Errorf("процессор или модель машины не попали на страницу:\n%s", wide)
	}
	// Whether the two stand side by side or one above the other is a matter
	// of geometry, and geometry is what the test asks about: side by side
	// the page is as tall as the taller of the two, stacked it is as tall as
	// both together.
	art := len(emblemFor(st.Host.DistroID, st.Host.Distro).art)
	fields := len(newTestScreen(st, true, 120).systemFields())
	tall := art
	if fields > tall {
		tall = fields
	}
	if got := len(newTestScreen(st, true, 120).system()); got != tall+1 {
		t.Errorf("на 120 колонках страница в %d строк, а рядом стоящие дали бы %d", got, tall+1)
	}
	if got := len(newTestScreen(st, true, 40).system()); got != art+fields+2 {
		t.Errorf("на 40 колонках страница в %d строк, а сложенные в столбик дали бы %d", got, art+fields+2)
	}
}

// The whole point of the map: two hundred and fifty-six processors have to be
// visible at once on an ordinary terminal.
func TestTwoHundredFiftySixCoresFitOnEightyColumns(t *testing.T) {
	s := newTestScreen(withCores(256, 20), true, 80)
	lines := s.coresBlock()
	body := plain(strings.Join(lines, "\n"))
	if !strings.Contains(body, "КАРТА ЯДЕР") {
		t.Fatalf("на 80 колонках 256 ядер нарисованы не картой:\n%s", body)
	}
	cells := 0
	for _, line := range strings.Split(body, "\n") {
		if !strings.Contains(line, "–") {
			continue
		}
		for _, r := range line {
			if strings.ContainsRune(string(sparkFaces)+"·", r) {
				cells++
			}
		}
	}
	if cells != 256 {
		t.Errorf("на карте %d ячеек, а ядер 256", cells)
	}
	if len(lines) > 20 {
		t.Errorf("карта заняла %d строк — на экране в 30 строк это уже прокрутка", len(lines))
	}
	// The busiest processor is named by its number, not left to be found by
	// eye among two hundred and fifty-six cells.
	if !strings.Contains(body, "ядро 20") {
		t.Errorf("самое занятое ядро не названо:\n%s", body)
	}
}

// A machine with few processors gets the gauges, one per processor.
func TestFewCoresAreDrawnAsGauges(t *testing.T) {
	s := newTestScreen(withCores(8, 3), true, 120)
	body := plain(strings.Join(s.coresBlock(), "\n"))
	if !strings.Contains(body, "ПО ЯДРАМ") {
		t.Fatalf("восемь ядер не нарисованы полосами:\n%s", body)
	}
	for i := 0; i < 8; i++ {
		if !strings.Contains(body, fmt.Sprint(i)) {
			t.Errorf("ядра %d нет на странице:\n%s", i, body)
		}
	}
	if !strings.Contains(body, "100%") {
		t.Errorf("занятое ядро не показано:\n%s", body)
	}
}

// A processor that was not measured is a dash, never a zero.
func TestUnmeasuredCoreIsNotAZero(t *testing.T) {
	st := withCores(4, 0)
	st.Load.Cores[2].BusyPercent = nil
	body := plain(strings.Join(newTestScreen(st, true, 120).coresBlock(), "\n"))
	if !strings.Contains(body, dash) {
		t.Errorf("незамеренное ядро не помечено прочерком:\n%s", body)
	}
	if strings.Contains(body, "  2 ────────   0%") {
		t.Errorf("незамеренное ядро нарисовано нулём:\n%s", body)
	}
}

// The whole snapshot without per-core figures still draws, and says nothing
// it cannot back up.
func TestCoresBlockWithoutMeasurements(t *testing.T) {
	body := plain(strings.Join(newTestScreen(filled(), true, 100).coresBlock(), "\n"))
	if strings.Contains(body, "%") {
		t.Errorf("без замера по ядрам напечатана доля:\n%s", body)
	}
	if !strings.Contains(body, dash) {
		t.Errorf("без замера по ядрам нет прочерка:\n%s", body)
	}
}

func gpuStatus() sysinfo.Status {
	st := filled()
	busy, celsius := 99.0, 85.0
	used, total := uint64(29)<<30, uint64(48)<<30
	st.GPUs = []gpuinfo.Card{
		{Node: "card0", Bus: "0000:02:00.0", Driver: "nvidia", Name: "NVIDIA RTX 6000 Ada Generation",
			BusyPercent: &busy, Celsius: &celsius, MemoryUsedBytes: &used, MemoryTotalBytes: &total,
			Source: "nvidia-smi", Outside: true},
		{Node: "card1", Bus: "0000:62:00.0", Driver: "mgag200", Name: "Matrox G200eW3",
			Source: "/sys/class/drm/card1/device"},
	}
	return st
}

func TestGPUPageShowsBothCardsAndWhereTheNumbersCameFrom(t *testing.T) {
	body := plain(strings.Join(newTestScreen(gpuStatus(), true, 120).gpus(), "\n"))
	for _, want := range []string{"NVIDIA RTX 6000 Ada Generation", "Matrox G200eW3",
		"99.0%", "85.0 °C", "драйвер nvidia", "чужой программы nvidia-smi", "/sys/class/drm/card1/device"} {
		if !strings.Contains(body, want) {
			t.Errorf("на странице видеокарт нет %q:\n%s", want, body)
		}
	}
	// The silent card gets dashes, and not one number it did not publish.
	tail := body[strings.Index(body, "Matrox"):]
	if strings.Contains(tail, "0.0%") || strings.Contains(tail, "0 Б") {
		t.Errorf("немая карта получила нули:\n%s", tail)
	}
	if strings.Count(tail, dash) < 3 {
		t.Errorf("немая карта показана без прочерков:\n%s", tail)
	}
}

func TestOverviewCarriesTheCardsAndTheCoreComb(t *testing.T) {
	st := gpuStatus()
	st.Load.Cores = withCores(64, 5).Load.Cores
	body := plain(strings.Join(newTestScreen(st, true, 120).overview(), "\n"))
	if !strings.Contains(body, "ВИДЕОКАРТЫ") {
		t.Errorf("на обзоре нет видеокарт:\n%s", body)
	}
	if !strings.Contains(body, "по ядрам") {
		t.Errorf("на обзоре нет гребёнки по ядрам:\n%s", body)
	}
}

// Folding many processors into few cells must not hide the busy one: the cell
// takes the largest share it covers, not the average.
func TestSqueezeKeepsTheBusiestCore(t *testing.T) {
	shares := make([]float64, 64)
	shares[37] = 1
	got := squeeze(shares, 8)
	if len(got) != 8 {
		t.Fatalf("ячеек %d, ждали 8", len(got))
	}
	max := 0.0
	for _, v := range got {
		if v > max {
			max = v
		}
	}
	if max != 1 {
		t.Errorf("занятое ядро потерялось при сжатии: %v", got)
	}
	if same := squeeze(shares[:4], 8); len(same) != 4 {
		t.Errorf("короткий список сжат до %d", len(same))
	}
}
