// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"digitdisk/internal/gpuinfo"
)

// Разбор того, что публикуют ядро и чужая программа. Всё здесь — чистые
// функции над текстом: живая машина в этих проверках не участвует, потому что
// проверяется не она, а понимание её ответов.

// ── дерево процессов через /proc ────────────────────────────────────────────

// procDir строит поддельный /proc: каталог на процесс, файл stat в каждом.
//
// Процессы кончаются здесь не удалением каталога, а ПЕРЕЕЗДОМ обхода на
// другой поддельный /proc, где их уже нет. Это и ближе к правде — /proc не
// каталог на диске, а окно в таблицу процессов, — и не заводит в дерево
// вызова, стирающего файлы: такие вызовы живут в host/internal/clean, и
// только там.
func procDir(t *testing.T, procs [][4]int) string {
	t.Helper()
	dir := t.TempDir()
	for _, p := range procs {
		writeProc(t, dir, p[0], p[1], p[2], p[3])
	}
	return dir
}

// writeProc кладёт один процесс: pid, ppid, тики процессора, страницы памяти.
func writeProc(t *testing.T, dir string, pid, ppid, ticks, pages int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, fmt.Sprint(pid)), 0o755); err != nil {
		t.Fatal(err)
	}
	// Поля после «)» — те же и в том же порядке, что в proc(5).
	line := fmt.Sprintf("%d (проба) S %d 0 0 0 -1 0 0 0 0 0 %d 0 0 0 20 0 1 0 100 0 %d\n",
		pid, ppid, ticks, pages)
	if err := os.WriteFile(filepath.Join(dir, fmt.Sprint(pid), "stat"), []byte(line), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestДеревоБерётСвоихИНеБерётЧужих — вся суть обхода: считается дерево
// команды, а не машина.
func TestДеревоБерётСвоихИНеБерётЧужих(t *testing.T) {
	// До запуска в машине уже живут два чужих процесса.
	tree := newProcTree(procDir(t, [][4]int{{1, 0, 500, 100}, {2, 1, 500, 100}}))

	// Команда запускается: 100 — она сама, 101 — её потомок, 102 — потомок
	// потомка, 103 — чужой процесс, появившийся тогда же.
	tree.dir = procDir(t, [][4]int{
		{1, 0, 500, 100}, {2, 1, 500, 100},
		{100, 2, 10, 50}, {101, 100, 20, 60}, {102, 101, 30, 70}, {103, 1, 40, 80},
	})
	tree.start(100, time.Now())

	page := uint64(os.Getpagesize())
	r := tree.sample(time.Now().Add(time.Second))
	if r.Processes != 3 {
		t.Errorf("процессов %d, ожидалось три — команда и двое её потомков", r.Processes)
	}
	if want := (50 + 60 + 70) * page; r.Bytes != want {
		t.Errorf("память %d, ожидалось %d — только своя", r.Bytes, want)
	}
	if want := float64(10+20+30) / userHZ; r.CPUSeconds != want {
		t.Errorf("процессорное время %v, ожидалось %v", r.CPUSeconds, want)
	}

	// Потомок кончился: его время остаётся в итоге, его память — нет.
	tree.dir = procDir(t, [][4]int{
		{1, 0, 500, 100}, {2, 1, 500, 100},
		{100, 2, 10, 50}, {101, 100, 20, 60}, {103, 1, 40, 80},
	})
	r = tree.sample(time.Now().Add(2 * time.Second))
	if r.Processes != 2 {
		t.Errorf("после конца потомка процессов %d, ожидалось два", r.Processes)
	}
	if want := float64(10+20+30) / userHZ; r.CPUSeconds != want {
		t.Errorf("процессорное время потерялось вместе с процессом: %v, ожидалось %v", r.CPUSeconds, want)
	}
	if want := (50 + 60) * page; r.Bytes != want {
		t.Errorf("память %d, ожидалось %d", r.Bytes, want)
	}
	if r.Peak != (50+60+70)*page {
		t.Errorf("пик %d — он обязан помнить, как было", r.Peak)
	}
	if tree.seen != 3 {
		t.Errorf("видено процессов %d, ожидалось три", tree.seen)
	}
}

// TestНомерПроцессаПереиспользуется: pid, который умер и выдан заново, — это
// НОВЫЙ процесс, и он может оказаться нашим. Помнить его чужим — значит тихо
// потерять целую ветку дерева на длинной сборке.
func TestНомерПроцессаПереиспользуется(t *testing.T) {
	tree := newProcTree(procDir(t, [][4]int{{1, 0, 0, 0}, {55, 1, 100, 10}}))
	tree.dir = procDir(t, [][4]int{{1, 0, 0, 0}, {55, 1, 100, 10}, {100, 1, 0, 10}})
	tree.start(100, time.Now())
	tree.sample(time.Now())

	// 55 был чужим и кончился.
	tree.dir = procDir(t, [][4]int{{1, 0, 0, 0}, {100, 1, 0, 10}})
	tree.sample(time.Now())

	// Номер выдан заново, и теперь он потомок команды.
	tree.dir = procDir(t, [][4]int{{1, 0, 0, 0}, {100, 1, 0, 10}, {55, 100, 7, 10}})
	r := tree.sample(time.Now())
	if r.Processes != 2 {
		t.Errorf("процессов %d, ожидалось два: номер выдан заново и он наш", r.Processes)
	}
}

// ── контрольная группа ──────────────────────────────────────────────────────

func TestРазборCPUStat(t *testing.T) {
	const text = "usage_usec 1234567\nuser_usec 1000000\nsystem_usec 234567\nnr_periods 0\n"
	got, ok := parseCPUUsage(text)
	if !ok {
		t.Fatal("usage_usec не найден")
	}
	if want := 1234567 * time.Microsecond; got != want {
		t.Errorf("%v, ожидалось %v", got, want)
	}
	if _, ok := parseCPUUsage("user_usec 1\n"); ok {
		t.Error("нашлось то, чего нет")
	}
}

// TestMaxНеЧисло: «max» в файле контрольной группы — это отсутствие предела, а
// не ноль, и принять его за замер значило бы напечатать бессмыслицу.
func TestMaxНеЧисло(t *testing.T) {
	for _, c := range []struct {
		text string
		want uint64
		ok   bool
	}{
		{"1048576\n", 1048576, true},
		{"max\n", 0, false},
		{"\n", 0, false},
		{"не число\n", 0, false},
	} {
		got, ok := parseNumber(c.text)
		if got != c.want || ok != c.ok {
			t.Errorf("parseNumber(%q) = %d,%v; ожидалось %d,%v", c.text, got, ok, c.want, c.ok)
		}
	}
}

func TestСвояГруппаНаходитсяВПервойСтроке(t *testing.T) {
	if got := ownCgroup("0::/user.slice/user-1000.slice/session-3.scope\n"); got != "/user.slice/user-1000.slice/session-3.scope" {
		t.Errorf("своя группа %q", got)
	}
	// Машина только с cgroup v1 отвечает иначе, и на ней контрольной
	// группы не будет: пусть лучше обход /proc, чем счёт по чужим правилам.
	if got := ownCgroup("11:cpu:/user\n10:memory:/user\n"); got != "" {
		t.Errorf("на v1 нашлась группа v2: %q", got)
	}
}

func TestПроцессыГруппыСчитаютсяПоСтрокам(t *testing.T) {
	if n := countPids("100\n101\n102\n"); n != 3 {
		t.Errorf("процессов %d, ожидалось три", n)
	}
	if n := countPids(""); n != 0 {
		t.Errorf("в пустом файле нашлось %d процессов", n)
	}
}

// ── чужая программа про видеокарту ──────────────────────────────────────────

// TestПамятьПоПроцессамРазбирается — единственный вопрос о карте, у которого
// есть честный ответ.
func TestПамятьПоПроцессамРазбирается(t *testing.T) {
	const answer = "3921, 1024\n4102, 512\n4200, [N/A]\nмусор\n"
	apps := gpuinfo.ParseComputeApps(answer)
	if len(apps) != 2 {
		t.Fatalf("процессов %d, ожидалось два: третий не измерен, четвёртый не строка", len(apps))
	}
	if apps[0].PID != 3921 || apps[0].Bytes != 1024*1024*1024 {
		t.Errorf("первый процесс разобран как %+v", apps[0])
	}
	if apps[1].Bytes != 512*1024*1024 {
		t.Errorf("второй процесс разобран как %+v", apps[1])
	}
}

// TestБезКлючаНичегоНеЗапускается — правило всего дерева: чужая программа
// зовётся только по ключу.
func TestБезКлючаНичегоНеЗапускается(t *testing.T) {
	ran := false
	r := gpuinfo.Reader{Run: func(string, ...string) ([]byte, error) {
		ran = true
		return []byte("1, 1\n"), nil
	}}
	if _, ok := r.ComputeApps(); ok {
		t.Error("ответ пришёл без ключа")
	}
	if ran {
		t.Fatal("чужая программа запущена без ключа --gpu-tool")
	}
	r.Tool = true
	apps, ok := r.ComputeApps()
	if !ok || len(apps) != 1 {
		t.Errorf("с ключом ответ не пришёл: %v, %v", apps, ok)
	}
	if !ran {
		t.Error("с ключом программа не запускалась")
	}
}

func TestКартаНеNVIDIAНеСпрашивается(t *testing.T) {
	if gpuinfo.HasNVIDIA([]gpuinfo.Card{{Driver: "mgag200"}, {Driver: "amdgpu", VendorID: "1002"}}) {
		t.Error("на машине без карты NVIDIA собрались звать nvidia-smi")
	}
	if !gpuinfo.HasNVIDIA([]gpuinfo.Card{{Driver: "mgag200"}, {VendorID: "10de"}}) {
		t.Error("карта NVIDIA не узнана по коду поставщика")
	}
}
