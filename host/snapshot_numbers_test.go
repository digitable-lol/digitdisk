// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"digitdisk/internal/gpuinfo"
	"digitdisk/internal/lang"
	"digitdisk/internal/procfs"
	"digitdisk/internal/report"
	"digitdisk/internal/sysinfo"
)

// ЗДЕСЬ СВЕРЯЕТСЯ НЕ ЧИСЛО, А СЛОВО, КОТОРЫМ СНИМОК ЭТО ЧИСЛО НАЗЫВАЕТ.
//
// Прогон .github/workflows/check.yml поднимает маковский бегунок, снимает
// снимок настоящей машины и ищет в нём образцы строк. Что числа на маке
// СНИМАЮТСЯ, доказывают тесты, которым нужен мак
// (internal/sysinfo/sysinfo_darwin_test.go); здесь доказывается второе — что
// образцы прогона совпадают с тем, что печать на самом деле печатает.
//
// БЕЗ ЭТОГО ТЕСТА ОБРАЗЕЦ И ПЕЧАТЬ РАЗОШЛИСЬ МОЛЧА И СТОЯЛИ ТАК ТРОЕ СУТОК.
// Образцы завели, когда снимок был только русским и печатал долю как «10.7%».
// Потом инструмент заговорил на двух языках: слова стали зависеть от языка, а
// разделитель дробной части — тоже («10,7%» по-русски). Образец «занято ЦП
// 10.7%» с тех пор не совпадал НИ С ОДНИМ языком, и шаг прогона докладывал
// «НЕТ ЧИСЛА» о числах, стоявших в той же строке журнала. Ствол был красен на
// каждом коммите, и ни одно настоящее падение под этим уже не было бы видно.
//
// Поэтому образцы читаются ИЗ ФАЙЛА ПРОГОНА, а не переписываются сюда:
// переписанная копия разошлась бы точно так же. Так же читается страница
// руководства в main_test.go — файл вне модуля, сверяемый с кодом.
//
// ЧЕГО ЭТОТ ТЕСТ НЕ ДОКАЗЫВАЕТ. Он печатает снимок, СОБРАННЫЙ ЗДЕСЬ, а не
// снятый с машины: он не может сказать, снимет ли мак эти числа. Это работа
// маковского бегунка и тестов с признаком darwin. Здесь — только про слова.

// workflowFile — файл прогона. Он вне модуля, и кэш `go test` его не
// отслеживает, поэтому шаг «Проверки хозяина» зовёт тесты с `-count=1`.
const workflowFile = "../.github/workflows/check.yml"

// sample — одна строка `need`/`named` из шага «Снимок содержит числа».
type sample struct {
	kind string // need — строка обязана нести число; named — назвать непокрытое
	lang lang.Lang
	re   *regexp.Regexp
	name string
}

var sampleLine = regexp.MustCompile(`(?m)^\s+(need|named) (ru|en) '([^']*)' '([^']*)'$`)

func samples(t *testing.T) []sample {
	t.Helper()
	src, err := os.ReadFile(workflowFile)
	if err != nil {
		t.Fatalf("файл прогона не читается: %v", err)
	}
	found := sampleLine.FindAllStringSubmatch(string(src), -1)
	out := make([]sample, 0, len(found))
	for _, m := range found {
		l, ok := lang.Parse(m[2])
		if !ok {
			t.Fatalf("в прогоне язык %q, а их два", m[2])
		}
		// (?m) — потому что в прогоне образец отдаётся `grep -E`, а тот
		// смотрит КАЖДУЮ строку отдельно. Без этого «^» в Go значил бы
		// начало всего снимка, и образец, привязанный к началу строки,
		// не сошёлся бы здесь, сходясь в прогоне.
		re, err := regexp.Compile("(?m)" + m[3])
		if err != nil {
			t.Fatalf("образец %q не разбирается: %v", m[3], err)
		}
		out = append(out, sample{kind: m[1], lang: l, re: re, name: m[4]})
	}
	// Разбор, который ничего не нашёл, тихо зеленеет — а это и есть та
	// беда, ради которой файл написан. Порог грубый нарочно: он ловит
	// сломавшийся разбор, а не считает образцы.
	if len(out) < 20 {
		t.Fatalf("в %s нашлось %d образцов — разбор сломался", workflowFile, len(out))
	}
	var needs, nameds int
	for _, s := range out {
		if s.kind == "need" {
			needs++
		} else {
			nameds++
		}
	}
	if needs == 0 || nameds == 0 {
		t.Fatalf("нашлось need=%d, named=%d — обе половины обязаны быть", needs, nameds)
	}
	return out
}

func printed(l lang.Lang, st sysinfo.Status) string {
	var b bytes.Buffer
	report.Status(&b, l, st)
	return b.String()
}

// naMac — снимок ровно такой, каким его собирает мак: замерено всё, что мак
// замеряет, и НАЗВАНО всё, чего маку взять нечем. Числа взяты из настоящего
// прогона на macos-latest, чтобы образцы проверялись на правдоподобных
// величинах, а не на удобных.
func onMac() sysinfo.Status {
	busy := 83.7
	top := func(pid int, rss int64, cpu float64, cmd string) sysinfo.Proc {
		c := cpu
		return sysinfo.Proc{
			PID: pid, PPID: 1, User: "runner", UID: 501, State: "S",
			Comm: "proc", Cmdline: cmd, Threads: 4,
			RSSBytes: rss, VSizeBytes: uint64(rss) * 8, CPUPercent: &c,
		}
	}
	ranking := []sysinfo.Proc{
		top(412, 150_786_048, 0.4, "/System/Library/CoreServices/Spotlight.app/Contents/MacOS/Spotlight"),
		top(948, 126_672_896, 0.3, "/opt/hca/hosted-compute-agent"),
		top(5529, 109_178_880, 0.2, "/Users/runner/actions-runner/extracted/bin/Runner.Worker"),
	}
	return sysinfo.Status{
		Host: sysinfo.Host{
			Hostname: "sat12.local", Distro: "macOS 26.6.2 (25G83)",
			Model: "VirtualMac2,1", CPUModel: "Apple M1 (Virtual)",
			KernelRelease: "25.6.0", Machine: "arm64",
		},
		Load: sysinfo.Load{
			LoadAvg:  procfs.LoadAvg{One: 4.19, Five: 12.23, Fifteen: 9.01},
			CPUCount: 3, BusyPercent: &busy, SampleMillis: 608,
		},
		Memory: procfs.Memory{
			Total:     7_516_192_768,
			Free:      236_874_137,
			BuffCache: 3_543_348_838,
			Available: 3_758_096_384,
			Used:      3_758_096_384,
			// Закреплённая и сжатая — маковские: на Linux этих строк в
			// снимке нет вовсе, и образцы на них проверяются только здесь.
			Raw: map[string]uint64{
				procfs.RawWired:      954_728_038,
				procfs.RawCompressed: 489_053_388,
			},
			Present: map[string]bool{
				procfs.FieldTotal: true, procfs.FieldFree: true,
				procfs.FieldBuffCache: true, procfs.FieldAvailable: true,
				procfs.FieldUsed: true,
			},
		},
		Processes: sysinfo.Processes{
			Total: 559, Running: 14, Threads: 946, WithDetail: 315,
			TopByMemory: ranking, TopByCPU: ranking,
		},
		GPUs: []gpuinfo.Card{{Name: "Apple M1 (Virtual)"}},
		// Чего на маке взять нечем. Обе причины — не наша недоделка:
		// показаний датчиков macOS не публикует вовсе, а состояние в записи
		// о процессе не отделяет заблокированный от спящего.
		Missing: map[string]lang.Phrase{
			sysinfo.FactSensors: lang.Say("macOS не публикует показания датчиков, а угадывать их формат нельзя"),
			sysinfo.FactBlocked: lang.Say("macOS не различает заблокированные процессы среди спящих"),
		},
	}
}

// nothing — снимок машины, с которой не сняли ничего. Present пуст НАРОЧНО:
// nil-карта у procfs.Memory значит «замерено всё» (см. Memory.Has), и снимок
// с nil напечатал бы нули там, где замера не было.
func nothing() sysinfo.Status {
	return sysinfo.Status{
		Memory: procfs.Memory{Present: map[string]bool{}},
		Missing: map[string]lang.Phrase{
			sysinfo.FactRunning:    lang.Say("сколько процессов работает прямо сейчас, видно только по их потокам"),
			sysinfo.FactThreads:    lang.Say("самопроверка памяти процессов не сошлась — их память и потоки не публикуем"),
			sysinfo.FactCPUBusy:    lang.Say("ядро не дало счётчики процессорного времени"),
			sysinfo.FactProcessRSS: lang.Say("самопроверка памяти процессов не сошлась — их память и потоки не публикуем"),
		},
	}
}

// TestОбразцыПрогонаСовпадаютСПечатью — положительный контроль: снимок мака
// собран, и КАЖДЫЙ образец шага прогона в нём находится.
func TestОбразцыПрогонаСовпадаютСПечатью(t *testing.T) {
	st := onMac()
	out := map[lang.Lang]string{lang.RU: printed(lang.RU, st), lang.EN: printed(lang.EN, st)}
	for _, s := range samples(t) {
		if !s.re.MatchString(out[s.lang]) {
			t.Errorf("[%s] образец %s %q (%s) не нашёлся в снимке, который его обязан нести.\n"+
				"Прогон %s ищет эту строку на живом маке; печать её больше не печатает.\n%s",
				s.lang, s.kind, s.re, s.name, workflowFile, out[s.lang])
		}
	}
}

// TestОбразцыЧиселНеСрабатываютНаНезамеренном — отрицательный контроль для
// `need`. Образец, который находится в снимке БЕЗ замеров, проверяет не число,
// а собственную формулировку: ровно так «выполняется [0-9]» сошёлся бы на
// «выполняется 0», которого никто не мерил.
func TestОбразцыЧиселНеСрабатываютНаНезамеренном(t *testing.T) {
	st := nothing()
	out := map[lang.Lang]string{lang.RU: printed(lang.RU, st), lang.EN: printed(lang.EN, st)}
	for _, s := range samples(t) {
		if s.kind != "need" {
			continue
		}
		if s.re.MatchString(out[s.lang]) {
			t.Errorf("[%s] образец need %q (%s) нашёлся в снимке, где не замерено НИЧЕГО.\n"+
				"Такой образец не отличит число от прочерка и зазеленеет на пустой машине.\n%s",
				s.lang, s.re, s.name, out[s.lang])
		}
	}
}

// TestОбразцыНепокрытогоНеСрабатываютНаПолномСнимке — отрицательный контроль
// для `named`. Строка «НЕ ИЗМЕРЕНО» печатается, только когда есть чего не
// измерить; образец, находящийся и в полном снимке, ничего не сторожит.
func TestОбразцыНепокрытогоНеСрабатываютНаПолномСнимке(t *testing.T) {
	st := onMac()
	st.Missing = nil
	out := map[lang.Lang]string{lang.RU: printed(lang.RU, st), lang.EN: printed(lang.EN, st)}
	for _, s := range samples(t) {
		if s.kind != "named" {
			continue
		}
		if s.re.MatchString(out[s.lang]) {
			t.Errorf("[%s] образец named %q (%s) нашёлся в снимке, где непокрытого нет вовсе.\n"+
				"Такой образец не проверяет, что непокрытое НАЗВАНО.\n%s",
				s.lang, s.re, s.name, out[s.lang])
		}
	}
}

// TestНепокрытоеНаМакеНазваноСловами — то, ради чего заведена половина
// `named`: чего на маке взять нечем, обязано стоять в снимке СЛОВОМ, а не
// уезжать в молчание. Причина при этом остаётся под `--why` и в снимок не
// лезет — это отдельный шаг прогона «В выводе нет инженерного дневника».
func TestНепокрытоеНаМакеНазваноСловами(t *testing.T) {
	st := onMac()
	for _, l := range []lang.Lang{lang.RU, lang.EN} {
		out := printed(l, st)
		for _, fact := range []string{sysinfo.FactSensors, sysinfo.FactBlocked} {
			word := l.Word(fact)
			if !regexp.MustCompile(`(?m)^(НЕ ИЗМЕРЕНО|NOT MEASURED) .*` + regexp.QuoteMeta(word)).MatchString(out) {
				t.Errorf("[%s] %q не названо в снимке словом:\n%s", l, word, out)
			}
			// И причина — не в снимке: она печатается по --why.
			var why bytes.Buffer
			report.Why(&why, l, st)
			if why.Len() == 0 {
				t.Errorf("[%s] --why ничего не сказал про %q", l, word)
			}
		}
	}
}
