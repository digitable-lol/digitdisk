// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"fmt"
	"io"
	"strings"

	"digitdisk/internal/clean"
	"digitdisk/internal/core"
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
func CleanPlan(w io.Writer, p clean.Plan, top int) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	pr("ПЛАН УБОРКИ  %s", p.Root)
	pr("  решающий слой  %s, договор версии %d", p.Decider, p.ContractVersion)
	if !p.DeciderReady {
		pr("")
		pr("  ВНИМАНИЕ: слой — заглушка. Она никому не выносит «%s», поэтому список ниже", core.VerdictRemovable)
		pr("  пуст не потому, что убирать нечего, а потому, что никто не решал.")
		pr("  Настоящий разбор: go build -tags flangcore -o digitdisk ./host")
	}
	if p.PlacesOrigin != "" {
		pr("  справочник     %d мест, %s", p.PlacesCount, p.PlacesOrigin)
	}
	if len(p.ProtectOrigins) > 0 {
		pr("  защитный спис. %s", strings.Join(p.ProtectOrigins, ", "))
	}
	pr("  обойдено       %d записей, %s", p.Walk.Entries, Bytes(p.Walk.TotalBytes))
	s := p.Walk.Skipped
	pr("  пропущено      %d  (нет доступа %d, исчезло %d, иные ошибки %d, граница ФС %d, предел глубины %d)",
		s.Total(), s.PermissionDenied, s.Vanished, s.OtherErrors, s.DeviceBoundaries, s.DepthLimited)
	if p.PrunedTrash > 0 {
		pr("  своя корзина   %d записей не обходилось (%s)", p.PrunedTrash, clean.TrashName)
	}

	pr("")
	pr("К ПЕРЕНОСУ В КОРЗИНУ  %d файлов, %s", len(p.Items), Bytes(p.Bytes))
	pr("  корзина        %s/<метка времени>", p.Trash)
	if p.HardlinkItems > 0 {
		pr("  из них %d — жёсткие ссылки: у их содержимого есть второе имя, и стирание", p.HardlinkItems)
		pr("  этого имени места не освободит. Освободится стиранием: %s", Bytes(p.FreeableBytes))
	}
	if len(p.ByClass) > 0 {
		pr("")
		pr("  по разрядам (весь план, ключ --top на этот счёт не влияет):")
		for _, c := range p.ByClass {
			pr("  %-11s %6d файлов  %10s", c.Class, c.Count, Bytes(c.Bytes))
		}
	}
	if len(p.Items) == 0 {
		pr("  — нечего: ядро не пометило «%s» ни одного файла", core.VerdictRemovable)
	} else {
		pr("")
		pr("  %10s  %-11s %-30s %s", "размер", "разряд", "почему", "путь")
		shown, hidden, hiddenBytes := 0, 0, int64(0)
		for _, it := range p.Items {
			if top > 0 && shown == top {
				hidden++
				hiddenBytes += it.Size
				continue
			}
			pr("  %10s  %-11s %-30s %s", Bytes(it.Size), it.Class, it.Why(), cut(it.Path, 70))
			if where := it.Where(); where != "" {
				pr("  %10s  %-11s место: %s", "", "", where)
			}
			shown++
		}
		if hidden > 0 {
			pr("  …и ещё %d файлов на %s — весь список: --top 0, или --json", hidden, Bytes(hiddenBytes))
		}
	}

	if len(p.Protected) > 0 {
		pr("")
		pr("ЗАЩИЩЕНО  %d файлов, %s: ядро назвало их «%s», защитный список запретил",
			len(p.Protected), Bytes(p.ProtectedBytes), core.VerdictRemovable)
		pr("  Это не расхождение слоёв, а ваше же распоряжение — оно и выполнено.")
		shown, hidden := 0, 0
		for _, pt := range p.Protected {
			if top > 0 && shown == top {
				hidden++
				continue
			}
			pr("  %10s  %-11s %s", Bytes(pt.Size), pt.Class, cut(pt.Path, 62))
			pr("              %s", pt.Rule)
			shown++
		}
		if hidden > 0 {
			pr("  …и ещё %d защищённых файлов", hidden)
		}
	}

	pr("")
	if len(p.Refused) == 0 {
		pr("ОТКАЗОВ НЕТ: хозяин согласен с ядром по каждой записи")
	} else {
		pr("ОТКАЗАНО  %d записей: ядро назвало их «%s», хозяин не тронет", len(p.Refused), core.VerdictRemovable)
		pr("  Это расхождение двух слоёв. Оно — факт о правилах, и его должен увидеть человек.")
		shown, hidden := 0, 0
		for _, r := range p.Refused {
			if top > 0 && shown == top {
				hidden++
				continue
			}
			pr("  %-11s %s", r.Class, cut(r.Path, 76))
			pr("               %s", r.Reason)
			shown++
		}
		if hidden > 0 {
			pr("  …и ещё %d отказов", hidden)
		}
	}

	pr("")
	pr("НИЧЕГО НЕ ТРОНУТО. Это план.")
	pr("  перенести в корзину:  digitdisk clean %s --apply", p.Root)
	pr("  перенос обратим: digitdisk restore <корзина>; стирает только digitdisk purge")
}

// Applied prints the outcome of a move into the корзина.
func Applied(w io.Writer, j *clean.Journal) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	moved, bytes := j.Moved()
	failed := j.Failed()

	pr("ПЕРЕНЕСЕНО В КОРЗИНУ  %d файлов, %s", moved, Bytes(bytes))
	pr("  корзина  %s", j.Box)
	pr("  журнал   %s", j.Path())
	pr("")
	pr("  Место НЕ освобождено: файлы лежат на той же файловой системе под другим")
	pr("  именем. Освобождает его только `digitdisk purge`, и это необратимо.")

	if len(failed) > 0 {
		pr("")
		pr("НЕ ПЕРЕНЕСЕНО  %d файлов: между обходом и переносом они изменились", len(failed))
		for _, it := range failed {
			pr("  %10s  %s", Bytes(it.Size), cut(it.Path, 70))
			pr("              %s", it.Failed)
		}
	}

	pr("")
	pr("  вернуть всё:  digitdisk restore %s", j.Box)
	pr("  стереть:      digitdisk purge %s --confirm %d", j.Box, moved)
}

// Restored prints the outcome of a возврат.
func Restored(w io.Writer, j *clean.Journal, dryRun bool) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	if dryRun {
		would := 0
		var bytes int64
		for _, it := range j.Items {
			if it.MovedAt != "" && it.RestoredAt == "" && it.PurgedAt == "" && it.Failed == "" {
				would++
				bytes += it.Size
			}
		}
		pr("ВЕРНУЛОСЬ БЫ  %d файлов, %s — из %s", would, Bytes(bytes), j.Box)
	} else {
		n, bytes := j.Restored()
		pr("ВОЗВРАЩЕНО  %d файлов, %s — на прежние места под %s", n, Bytes(bytes), j.Root)
		pr("  журнал  %s", j.Path())
	}

	var stuck []clean.Item
	for _, it := range j.Items {
		if it.Failed != "" {
			stuck = append(stuck, it)
		}
	}
	if len(stuck) > 0 {
		pr("")
		pr("НЕ ВОЗВРАЩЕНО  %d записей", len(stuck))
		for _, it := range stuck {
			pr("  %10s  %s", Bytes(it.Size), cut(it.Path, 70))
			pr("              %s", it.Failed)
		}
	}
}

// PurgePlan prints what erasing this корзина would destroy, and the number
// --confirm has to name.
func PurgePlan(w io.Writer, j *clean.Journal) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	n, bytes := j.Moved()
	pr("ПЛАН СТИРАНИЯ  %s", j.Box)
	pr("  перенесено сюда  %s из %s", j.StartedAt, j.Root)
	pr("  в корзине        %d файлов, %s", n, Bytes(bytes))
	if r, _ := j.Restored(); r > 0 {
		pr("  уже возвращено   %d — их стирать нечем", r)
	}
	if p, _ := j.Purged(); p > 0 {
		pr("  уже стёрто       %d", p)
	}
	pr("")
	if n == 0 {
		pr("СТИРАТЬ НЕЧЕГО.")
		return
	}
	pr("  %10s  %s", "размер", "откуда взят")
	shown := 0
	for _, it := range j.Items {
		if it.MovedAt == "" || it.RestoredAt != "" || it.PurgedAt != "" {
			continue
		}
		if shown == 20 {
			pr("  … и ещё %d", n-shown)
			break
		}
		pr("  %10s  %s", Bytes(it.Size), cut(it.Path, 74))
		shown++
	}
	pr("")
	pr("НИЧЕГО НЕ СТЁРТО. Это план, и стирание необратимо.")
	pr("  стереть:  digitdisk purge %s --confirm %d", j.Box, n)
	pr("  вернуть:  digitdisk restore %s", j.Box)
}

// Purged prints the outcome of an erase.
func Purged(w io.Writer, j *clean.Journal) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	n, bytes := j.Purged()
	pr("СТЁРТО  %d файлов, %s", n, Bytes(bytes))
	pr("  корзина  %s", j.Box)
	pr("  журнал   %s — он остаётся: это запись о том, чего больше нет", j.Path())

	var stuck []clean.Item
	for _, it := range j.Items {
		if it.Failed != "" {
			stuck = append(stuck, it)
		}
	}
	if len(stuck) > 0 {
		pr("")
		pr("НЕ СТЁРТО  %d записей", len(stuck))
		for _, it := range stuck {
			pr("  %10s  %s", Bytes(it.Size), cut(it.Path, 70))
			pr("              %s", it.Failed)
		}
	}
}
