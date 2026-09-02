// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"fmt"
	"io"
	"strings"

	"digitdisk/internal/clean"
	"digitdisk/internal/core"
	"digitdisk/internal/lang"
)

// CleanPlan prints what `clean --apply` would move, and nothing about what it
// would free — moving into the корзина frees no space, and a plan that
// promised a saving at that point would be promising the wrong step's result.
//
// # Why the lists are cut and the counts are not
//
// A plan over a real home directory is hundreds of lines long, and a wall of
// them is not a plan anybody reads: the first thing a person wants is the
// total and the biggest few.  So every LIST here stops at top and says what it
// left out, while every COUNT — the total, the bytes, the breakdown by разряд —
// is computed over the whole plan and never changes with top.  A summary that
// moved when the screen got shorter would be a summary of the screen.
//
// top ≤ 0 prints everything, the same "0 — без предела" this CLI already uses
// for --max-depth.  `--json` is not cut at all: see cmdClean.
func CleanPlan(w io.Writer, l lang.Lang, p clean.Plan, top int) {
	pr := func(line string) { fmt.Fprintln(w, line) }

	pr(l.F("ПЛАН УБОРКИ  %s", p.Root))
	pr(l.F("  решающий слой  %s, договор версии %d", l.Word(p.Decider), p.ContractVersion))
	if !p.DeciderReady {
		pr("")
		pr(l.F("  ВНИМАНИЕ: слой — заглушка. Она никому не выносит «%s», поэтому список ниже",
			l.Word(string(core.VerdictRemovable))))
		pr(l.T("  пуст не потому, что убирать нечего, а потому, что никто не решал."))
		pr(l.T("  Настоящий разбор: go build -tags flangcore -o digitdisk ./host"))
	}
	if p.PlacesOrigin != "" {
		pr(l.F("  справочник     %d мест, %s", p.PlacesCount, l.Word(p.PlacesOrigin)))
	}
	if len(p.ProtectOrigins) > 0 {
		pr(l.F("  защитный спис. %s", strings.Join(p.ProtectOrigins, ", ")))
	}
	pr(l.F("  обойдено       %s записей, %s", l.Num(int64(p.Walk.Entries)), l.Bytes(p.Walk.TotalBytes)))
	s := p.Walk.Skipped
	pr(l.F("  пропущено      %s  (нет доступа %s, исчезло %s, иные ошибки %s, граница ФС %s, предел глубины %s)",
		l.Num(int64(s.Total())), l.Num(int64(s.PermissionDenied)), l.Num(int64(s.Vanished)),
		l.Num(int64(s.OtherErrors)), l.Num(int64(s.DeviceBoundaries)), l.Num(int64(s.DepthLimited))))
	if p.PrunedTrash > 0 {
		pr(l.F("  своя корзина   %s записей не обходилось (%s)", l.Num(int64(p.PrunedTrash)), clean.TrashName))
	}

	pr("")
	pr(l.F("К ПЕРЕНОСУ В КОРЗИНУ  %s файлов, %s", l.Num(int64(len(p.Items))), l.Bytes(p.Bytes)))
	pr(l.F("  корзина        %s/<метка времени>", p.Trash))
	if p.HardlinkItems > 0 {
		pr(l.F("  из них %s — жёсткие ссылки: у их содержимого есть второе имя, и стирание",
			l.Num(int64(p.HardlinkItems))))
		pr(l.F("  этого имени места не освободит. Освободится стиранием: %s", l.Bytes(p.FreeableBytes)))
	}
	if len(p.ByClass) > 0 {
		pr("")
		pr(l.T("  по разрядам (весь план, ключ --top на этот счёт не влияет):"))
		for _, c := range p.ByClass {
			pr(l.F("  %-11s %6s файлов  %10s", l.Word(string(c.Class)), l.Num(int64(c.Count)), l.Bytes(c.Bytes)))
		}
	}
	if len(p.Items) == 0 {
		pr(l.F("  — нечего: ядро не пометило «%s» ни одного файла", l.Word(string(core.VerdictRemovable))))
	} else {
		pr("")
		pr(fmt.Sprintf("  %10s  %-11s %-30s %s", l.T("размер"), l.T("разряд"), l.T("почему"), l.T("путь")))
		shown, hidden, hiddenBytes := 0, 0, int64(0)
		for _, it := range p.Items {
			if top > 0 && shown == top {
				hidden++
				hiddenBytes += it.Size
				continue
			}
			pr(fmt.Sprintf("  %10s  %-11s %-30s %s", l.Bytes(it.Size), l.Word(string(it.Class)),
				it.Why(l), cut(it.Path, 70)))
			if where := it.Where(); where != "" {
				pr(l.F("  %10s  %-11s место: %s", "", "", where))
			}
			shown++
		}
		if hidden > 0 {
			pr(l.F("  …и ещё %s файлов на %s — весь список: --top 0, или --json",
				l.Num(int64(hidden)), l.Bytes(hiddenBytes)))
		}
	}

	if len(p.Protected) > 0 {
		pr("")
		pr(l.F("ЗАЩИЩЕНО  %s файлов, %s: ядро назвало их «%s», защитный список запретил",
			l.Num(int64(len(p.Protected))), l.Bytes(p.ProtectedBytes), l.Word(string(core.VerdictRemovable))))
		pr(l.T("  Это не расхождение слоёв, а ваше же распоряжение — оно и выполнено."))
		shown, hidden := 0, 0
		for _, pt := range p.Protected {
			if top > 0 && shown == top {
				hidden++
				continue
			}
			pr(fmt.Sprintf("  %10s  %-11s %s", l.Bytes(pt.Size), l.Word(string(pt.Class)), cut(pt.Path, 62)))
			// Правило называется словом читателя, а значение, файл и
			// строка — теми словами, которыми их написал сам человек:
			// переводить чужой файл никто не подряжался.
			pr(fmt.Sprintf("              %s", pt.Rule.In(l)))
			shown++
		}
		if hidden > 0 {
			pr(l.F("  …и ещё %s защищённых файлов", l.Num(int64(hidden))))
		}
	}

	pr("")
	if len(p.Refused) == 0 {
		pr(l.T("ОТКАЗОВ НЕТ: хозяин согласен с ядром по каждой записи"))
	} else {
		pr(l.F("ОТКАЗАНО  %s записей: ядро назвало их «%s», хозяин не тронет",
			l.Num(int64(len(p.Refused))), l.Word(string(core.VerdictRemovable))))
		pr(l.T("  Это расхождение двух слоёв. Оно — факт о правилах, и его должен увидеть человек."))
		shown, hidden := 0, 0
		for _, r := range p.Refused {
			if top > 0 && shown == top {
				hidden++
				continue
			}
			pr(fmt.Sprintf("  %-11s %s", l.Word(string(r.Class)), cut(r.Path, 76)))
			pr(fmt.Sprintf("               %s", r.Reason.In(l)))
			shown++
		}
		if hidden > 0 {
			pr(l.F("  …и ещё %s отказов", l.Num(int64(hidden))))
		}
	}

	pr("")
	pr(l.T("НИЧЕГО НЕ ТРОНУТО. Это план."))
	pr(l.F("  перенести в корзину:  digitdisk clean %s --apply", p.Root))
	pr(l.T("  перенос обратим: digitdisk restore <корзина>; стирает только digitdisk purge"))
}

// Applied prints the outcome of a move into the корзина.
func Applied(w io.Writer, l lang.Lang, j *clean.Journal) {
	pr := func(line string) { fmt.Fprintln(w, line) }

	moved, bytes := j.Moved()
	failed := j.Failed()

	pr(l.F("ПЕРЕНЕСЕНО В КОРЗИНУ  %s файлов, %s", l.Num(int64(moved)), l.Bytes(bytes)))
	pr(l.F("  корзина  %s", j.Box))
	pr(l.F("  журнал   %s", j.Path()))
	pr("")
	pr(l.T("  Место НЕ освобождено: файлы лежат на той же файловой системе под другим"))
	pr(l.T("  именем. Освобождает его только `digitdisk purge`, и это необратимо."))

	if len(failed) > 0 {
		pr("")
		pr(l.F("НЕ ПЕРЕНЕСЕНО  %s файлов: между обходом и переносом они изменились",
			l.Num(int64(len(failed)))))
		for _, it := range failed {
			pr(fmt.Sprintf("  %10s  %s", l.Bytes(it.Size), cut(it.Path, 70)))
			pr(fmt.Sprintf("              %s", it.Failed.In(l)))
		}
	}

	pr("")
	pr(l.F("  вернуть всё:  digitdisk restore %s", j.Box))
	// The number after --confirm is typed by a person into a shell: it is
	// written without grouping on purpose, because «1 234» pasted into a
	// command line is not a number any shell will take.
	pr(l.F("  стереть:      digitdisk purge %s --confirm %d", j.Box, moved))
}

// Restored prints the outcome of a возврат.
func Restored(w io.Writer, l lang.Lang, j *clean.Journal, dryRun bool) {
	pr := func(line string) { fmt.Fprintln(w, line) }

	if dryRun {
		would := 0
		var bytes int64
		for _, it := range j.Items {
			if it.MovedAt != "" && it.RestoredAt == "" && it.PurgedAt == "" && it.Failed.Empty() {
				would++
				bytes += it.Size
			}
		}
		pr(l.F("ВЕРНУЛОСЬ БЫ  %s файлов, %s — из %s", l.Num(int64(would)), l.Bytes(bytes), j.Box))
	} else {
		n, bytes := j.Restored()
		pr(l.F("ВОЗВРАЩЕНО  %s файлов, %s — на прежние места под %s",
			l.Num(int64(n)), l.Bytes(bytes), j.Root))
		pr(l.F("  журнал  %s", j.Path()))
	}

	var stuck []clean.Item
	for _, it := range j.Items {
		if !it.Failed.Empty() {
			stuck = append(stuck, it)
		}
	}
	if len(stuck) > 0 {
		pr("")
		pr(l.F("НЕ ВОЗВРАЩЕНО  %s записей", l.Num(int64(len(stuck)))))
		for _, it := range stuck {
			pr(fmt.Sprintf("  %10s  %s", l.Bytes(it.Size), cut(it.Path, 70)))
			pr(fmt.Sprintf("              %s", it.Failed.In(l)))
		}
	}
}

// PurgePlan prints what erasing this корзина would destroy, and the number
// --confirm has to name.
func PurgePlan(w io.Writer, l lang.Lang, j *clean.Journal) {
	pr := func(line string) { fmt.Fprintln(w, line) }

	n, bytes := j.Moved()
	pr(l.F("ПЛАН СТИРАНИЯ  %s", j.Box))
	pr(l.F("  перенесено сюда  %s из %s", l.StampDate(j.StartedAt), j.Root))
	pr(l.F("  в корзине        %s файлов, %s", l.Num(int64(n)), l.Bytes(bytes)))
	if r, _ := j.Restored(); r > 0 {
		pr(l.F("  уже возвращено   %s — их стирать нечем", l.Num(int64(r))))
	}
	if p, _ := j.Purged(); p > 0 {
		pr(l.F("  уже стёрто       %s", l.Num(int64(p))))
	}
	pr("")
	if n == 0 {
		pr(l.T("СТИРАТЬ НЕЧЕГО."))
		return
	}
	pr(fmt.Sprintf("  %10s  %s", l.T("размер"), l.T("откуда взят")))
	shown := 0
	for _, it := range j.Items {
		if it.MovedAt == "" || it.RestoredAt != "" || it.PurgedAt != "" {
			continue
		}
		if shown == 20 {
			pr(l.F("  … и ещё %s", l.Num(int64(n-shown))))
			break
		}
		pr(fmt.Sprintf("  %10s  %s", l.Bytes(it.Size), cut(it.Path, 74)))
		shown++
	}
	pr("")
	pr(l.T("НИЧЕГО НЕ СТЁРТО. Это план, и стирание необратимо."))
	pr(l.F("  стереть:  digitdisk purge %s --confirm %d", j.Box, n))
	pr(l.F("  вернуть:  digitdisk restore %s", j.Box))
}

// Purged prints the outcome of an erase.
func Purged(w io.Writer, l lang.Lang, j *clean.Journal) {
	pr := func(line string) { fmt.Fprintln(w, line) }

	n, bytes := j.Purged()
	pr(l.F("СТЁРТО  %s файлов, %s", l.Num(int64(n)), l.Bytes(bytes)))
	pr(l.F("  корзина  %s", j.Box))
	pr(l.F("  журнал   %s — он остаётся: это запись о том, чего больше нет", j.Path()))

	var stuck []clean.Item
	for _, it := range j.Items {
		if !it.Failed.Empty() {
			stuck = append(stuck, it)
		}
	}
	if len(stuck) > 0 {
		pr("")
		pr(l.F("НЕ СТЁРТО  %s записей", l.Num(int64(len(stuck)))))
		for _, it := range stuck {
			pr(fmt.Sprintf("  %10s  %s", l.Bytes(it.Size), cut(it.Path, 70)))
			pr(fmt.Sprintf("              %s", it.Failed.In(l)))
		}
	}
}
