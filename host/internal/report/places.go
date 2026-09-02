// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"fmt"
	"io"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/places"
)

// Places prints the справочник известных мест and what it found on this
// machine.  It is the answer to "what does this tool actually know", and it
// exists because a справочник nobody can look at is indistinguishable from a
// hardcoded list: the point of keeping it in a file is that the file can be
// read, checked against the machine and argued with.
func Places(w io.Writer, d *places.Directory, found []places.Found, top int) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	here, absent, elsewhere, anywhere := 0, 0, 0, 0
	var bytes int64
	var files int
	for _, f := range found {
		switch {
		case !f.Applies:
			elsewhere++
		case f.Resolved == "":
			anywhere++
		case f.Exists:
			here++
			bytes += f.Bytes
			files += f.Files
		default:
			absent++
		}
	}

	pr("СПРАВОЧНИК ИЗВЕСТНЫХ МЕСТ  %s", d.Origin)
	pr("  всего мест     %d", len(d.Entries))
	pr("  на этой машине %d найдено, %d нет, %d на любой глубине, %d для другой системы",
		here, absent, anywhere, elsewhere)
	pr("  занято         %s в %d файлах (только найденные каталоги)", Bytes(bytes), files)
	pr("")
	pr("  Справочник — данные: %s. Свой: --places ФАЙЛ или ~/.config/%s",
		d.Origin, places.UserFile)
	pr("  Разряд из справочника не смягчает приговора: пороги, каталоги, ссылки и")
	pr("  хранилища, адресуемые содержимым, судятся ядром ровно как прежде.")

	pr("")
	pr("  %-9s %10s %8s  %-34s %s", "разряд", "размер", "файлов", "имя", "путь")
	shown, hidden := 0, 0
	for _, f := range found {
		if !f.Applies || !f.Exists {
			continue
		}
		if top > 0 && shown == top {
			hidden++
			continue
		}
		pr("  %-9s %10s %8d  %-34s %s", f.Class, Bytes(f.Bytes), f.Files, cut(f.Name, 34), cut(f.Resolved, 60))
		shown++
	}
	if shown == 0 {
		pr("  — ни одно место справочника на этой машине не нашлось")
	}
	if hidden > 0 {
		pr("  …и ещё %d найденных мест — весь список: --top 0, или --json", hidden)
	}

	pr("")
	pr("  Чего здесь нет — тоже ответ: %d мест справочник называет, а на этой машине", absent)
	pr("  их каталогов не существует. Это не ошибка: справочник общий, а машины разные.")
}

// History prints the корзины digitdisk has left under a root: when, how much,
// and how to get it back.
func History(w io.Writer, h *clean.History, now time.Time, top int) {
	pr := func(format string, a ...any) { fmt.Fprintf(w, format+"\n", a...) }

	pr("ЖУРНАЛ УБОРКИ  %s", h.Root)
	pr("  хранилище корзин  %s", h.Trash)
	pr("  корзин            %d", h.Boxes)
	pr("  лежит в корзинах  %s — место ещё НЕ освобождено", Bytes(h.MovedBytes))
	pr("  освобождено       %s — это стёртое, и вернуть его нечем", Bytes(h.FreedBytes))
	pr("  возвращено        %s", Bytes(h.RestoredBytes))

	if h.Boxes == 0 {
		pr("")
		pr("  Уборки здесь не было: ни одной корзины под %s.", h.Trash)
		return
	}

	pr("")
	pr("  %-22s %8s %10s %8s %8s  %s", "когда", "в корзине", "байт", "возвращ.", "стёрто", "корзина")
	shown, hidden := 0, 0
	for _, e := range h.Entries {
		if top > 0 && shown == top {
			hidden++
			continue
		}
		when := dash(e.StartedAt)
		if age := e.Age(now); age >= 0 {
			when = fmt.Sprintf("%s (%s назад)", when[:min(len(when), 10)], since(age))
		}
		pr("  %-22s %8d %10s %8d %8d  %s", when, e.Moved, Bytes(e.MovedBytes), e.Restored, e.Purged, cut(e.Box, 46))
		if e.Problem != "" {
			pr("      беда: %s", e.Problem)
		}
		shown++
	}
	if hidden > 0 {
		pr("  …и ещё %d корзин — весь список: --top 0, или --json", hidden)
	}

	var newest *clean.Entry
	for i := range h.Entries {
		if h.Entries[i].Restorable() {
			newest = &h.Entries[i]
			break
		}
	}
	pr("")
	if newest == nil {
		pr("  Возвращать нечего: во всех корзинах пусто.")
		return
	}
	pr("  вернуть последнее:  digitdisk restore %s", newest.Box)
	pr("  стереть насовсем:   digitdisk purge %s --confirm %d", newest.Box, newest.Moved)
}

// since renders a duration the way a person says it out loud.
func since(d time.Duration) string {
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%d мин", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%d ч", int(d.Hours()))
	default:
		return fmt.Sprintf("%d дн", int(d.Hours()/24))
	}
}
