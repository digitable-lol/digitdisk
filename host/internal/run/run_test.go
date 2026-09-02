// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/lang"
)

// Здесь проверяется то, чем обёртка отличается от вызова команды напрямую, —
// и то, чем она не должна отличаться.
//
// Главное обещание одно: вывод команды не тронут. Всё остальное — как считаны
// дерево процессов и во что превращены его числа.

// TestВыводКомандыНеТронут — то самое побайтовое обещание: `digitdisk run X >
// файл` кладёт в файл ровно то, что положил бы туда X.
func TestВыводКомандыНеТронут(t *testing.T) {
	const script = `printf 'первая\nвторая\n'; printf 'без перевода строки'`

	direct := filepath.Join(t.TempDir(), "прямо")
	wrapped := filepath.Join(t.TempDir(), "под-обёрткой")

	runShell(t, script, direct, nil)
	res, err := Run(Options{
		Args:  []string{"sh", "-c", script},
		Out:   create(t, wrapped),
		Err:   devNull(t),
		Plain: true,
		Lang:  lang.RU,
	})
	if err != nil {
		t.Fatalf("команда не запустилась: %v", err)
	}
	if res.Code != 0 {
		t.Errorf("код %d, ожидался 0", res.Code)
	}
	a, b := read(t, direct), read(t, wrapped)
	if a != b {
		t.Errorf("вывод под обёрткой отличается:\nбез неё %q\nпод ней %q", a, b)
	}
	if strings.Contains(b, "\x1b") {
		t.Error("в выводе команды оказались управляющие последовательности")
	}
}

// TestКодВозвратаЧужой: обёртка не сочиняет код возврата.
func TestКодВозвратаЧужой(t *testing.T) {
	for _, code := range []int{0, 1, 3, 42} {
		res, err := Run(Options{
			Args:  []string{"sh", "-c", "exit " + itoa(code)},
			Out:   devNull(t),
			Err:   devNull(t),
			Plain: true,
			Lang:  lang.RU,
		})
		if err != nil {
			t.Fatalf("команда не запустилась: %v", err)
		}
		if res.Code != code {
			t.Errorf("код %d, ожидался %d", res.Code, code)
		}
		if res.Signal != "" {
			t.Errorf("сигнал %q там, где его не было", res.Signal)
		}
	}
}

// TestСигналНазванИменем: команда, убитую сигналом, нельзя показать как
// команду с кодом возврата — у неё его нет.
func TestСигналНазванИменем(t *testing.T) {
	res, err := Run(Options{
		Args:  []string{"sh", "-c", "kill -TERM $$"},
		Out:   devNull(t),
		Err:   devNull(t),
		Plain: true,
		Lang:  lang.RU,
	})
	if err != nil {
		t.Fatalf("команда не запустилась: %v", err)
	}
	if res.Signal != "SIGTERM" {
		t.Errorf("сигнал %q, ожидался SIGTERM", res.Signal)
	}
	if res.Code != 143 {
		t.Errorf("код %d, ожидался 143 — то, что печатает оболочка", res.Code)
	}
}

// TestНетТакойКомандыЭто127 — числа, которые обещает страница руководства.
func TestНетТакойКомандыЭто127(t *testing.T) {
	_, err := Run(Options{
		Args:  []string{"такой-команды-нет-и-не-будет"},
		Out:   devNull(t),
		Err:   devNull(t),
		Plain: true,
		Lang:  lang.RU,
	})
	if err == nil {
		t.Fatal("несуществующая команда запустилась")
	}
	if code := StartCode(err); code != 127 {
		t.Errorf("код %d, ожидался 127", code)
	}
	if StartCode(os.ErrNotExist) != 1 {
		t.Error("чужая ошибка получила код запуска")
	}
}

// TestИтогиСчитаныИНеПридуманы: у команды, которая точно поработала, есть
// процессорное время, и оно от ядра.
func TestИтогиСчитаныИНеПридуманы(t *testing.T) {
	res, err := Run(Options{
		Args:  []string{"sh", "-c", "i=0; while [ $i -lt 200000 ]; do i=$((i+1)); done"},
		Out:   devNull(t),
		Err:   devNull(t),
		Plain: true,
		Lang:  lang.RU,
	})
	if err != nil {
		t.Fatalf("команда не запустилась: %v", err)
	}
	if !res.CPUExact {
		t.Error("процессорное время не помечено точным, а оно от wait4")
	}
	if res.CPUSeconds <= 0 {
		t.Errorf("процессорное время %v — цикл на двести тысяч итераций столько не стоит", res.CPUSeconds)
	}
	if res.Seconds <= 0 {
		t.Error("прошло нисколько времени")
	}
	if res.Accounting == "" {
		t.Error("учёт не назван, а сводка обязана его называть")
	}
	t.Logf("учёт %s, %.3f с, ЦП %.3f с, пик %d Б", res.Accounting, res.Seconds, res.CPUSeconds, res.PeakBytes)
}

func runShell(t *testing.T, script, out string, env []string) {
	t.Helper()
	f := create(t, out)
	res, err := Run(Options{
		Args: []string{"sh", "-c", script}, Out: f, Err: devNull(t), Plain: true, Lang: lang.RU,
	})
	if err != nil || res.Code != 0 {
		t.Fatalf("эталонный прогон не удался: %v (код %d)", err, res.Code)
	}
}

func create(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.Create(name)
	if err != nil {
		t.Fatalf("%s не создаётся: %v", name, err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

func devNull(t *testing.T) *os.File {
	t.Helper()
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("%s не открывается: %v", os.DevNull, err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

func read(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("%s не читается: %v", name, err)
	}
	return string(b)
}

func itoa(n int) string { return strings.TrimSpace(lang.RU.F("%d", n)) }

// ── строка состояния ────────────────────────────────────────────────────────

func TestСтрокаСостоянияНеШиреТерминала(t *testing.T) {
	r := reading{Known: true, Processes: 37, Bytes: 2 << 30, Peak: 6 << 30, Percent: 340, CPUSeconds: 250}
	g := gpuReading{Known: true, Bytes: 1 << 30}
	for _, width := range []int{200, 100, 80, 60, 40, 20, 10, 5} {
		for _, l := range []lang.Lang{lang.RU, lang.EN} {
			line := barText(l, width, 83*time.Second, r, g)
			if n := len([]rune(line)); n > width {
				t.Errorf("ширина %d, язык %s: строка в %d знаков: %q", width, l, n, line)
			}
			if !strings.HasPrefix(line, "1:23") {
				t.Errorf("ширина %d, язык %s: время должно стоять первым: %q", width, l, line)
			}
		}
	}
}

func TestСтрокаСостоянияМолчитОНеизмеренном(t *testing.T) {
	line := barText(lang.RU, 80, 5*time.Second, reading{}, gpuReading{})
	if line != "0:05" {
		t.Errorf("без замеров в строке %q, а должно быть одно время", line)
	}
}

func TestЧасыИдутОтНуля(t *testing.T) {
	for _, c := range []struct {
		d    time.Duration
		want string
	}{
		{0, "0:00"},
		{9 * time.Second, "0:09"},
		{83 * time.Second, "1:23"},
		{3723 * time.Second, "1:02:03"},
	} {
		if got := clock(c.d); got != c.want {
			t.Errorf("clock(%v) = %q, ожидалось %q", c.d, got, c.want)
		}
	}
}

// TestСводкаНазываетТочностьПрямо: «6,1 ГиБ» и «около 6,1 ГиБ» — разные
// утверждения, и сводка обязана говорить, какое из них она делает.
func TestСводкаНазываетТочностьПрямо(t *testing.T) {
	exact := Lines(lang.RU, Result{Command: "make", Seconds: 10, CPUSeconds: 30,
		PeakBytes: 6 << 30, PeakExact: true, Processes: 12, Accounting: ByGroup})
	if strings.Contains(strings.Join(exact, " "), "около") {
		t.Errorf("точный пик назван приблизительным: %q", exact)
	}
	rough := Lines(lang.RU, Result{Command: "make", Seconds: 10, CPUSeconds: 30,
		PeakBytes: 6 << 30, Processes: 12, Accounting: ByProc, SampleMS: 200})
	if !strings.Contains(strings.Join(rough, " "), "около") {
		t.Errorf("приблизительный пик выдан за точный: %q", rough)
	}
	for _, l := range []lang.Lang{lang.RU, lang.EN} {
		if n := len(Lines(l, Result{Command: "make", Accounting: ByProc})); n != 2 {
			t.Errorf("%s: сводка в %d строк, а обещаны две", l, n)
		}
	}
}

// TestПроВидеокартуТолькоТо,ЧтоМожноЗнать: память по процессам — можно,
// загрузку по процессам — нельзя, и второе говорится словом.
func TestПроВидеокартуГоворитсяТолькоЗнаемое(t *testing.T) {
	lines := Lines(lang.RU, Result{Command: "x", Accounting: ByProc,
		GPU: &GPU{PeakBytes: 1 << 30, Source: "nvidia-smi"}})
	if len(lines) != 3 {
		t.Fatalf("строк %d, ожидалось три: %q", len(lines), lines)
	}
	if !strings.Contains(lines[2], "не публикует") {
		t.Errorf("про загрузку по процессам ничего не сказано: %q", lines[2])
	}
	if g := (&GPU{}).ByProcessLoad; g {
		t.Error("загрузка по процессам помечена измеренной")
	}
}
