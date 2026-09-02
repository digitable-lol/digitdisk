// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package report

import (
	"fmt"
	"io"
	"time"

	"digitdisk/internal/clean"
	"digitdisk/internal/lang"
	"digitdisk/internal/places"
)

// Places prints the справочник известных мест and what it found on this
// machine.  It is the answer to "what does this tool actually know", and it
// exists because a справочник nobody can look at is indistinguishable from a
// hardcoded list: the point of keeping it in a file is that the file can be
// read, checked against the machine and argued with.
func Places(w io.Writer, l lang.Lang, d *places.Directory, found []places.Found, top int) {
	pr := func(line string) { fmt.Fprintln(w, line) }

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

	pr(l.F("СПРАВОЧНИК ИЗВЕСТНЫХ МЕСТ  %s", l.Word(d.Origin)))
	pr(l.F("  всего мест     %d", len(d.Entries)))
	pr(l.F("  на этой машине %d найдено, %d нет, %d на любой глубине, %d для другой системы",
		here, absent, anywhere, elsewhere))
	pr(l.F("  занято         %s в %s файлах (только найденные каталоги)", l.Bytes(bytes), l.Num(int64(files))))
	pr("")
	pr(l.F("  Справочник — данные: %s. Свой: --places ФАЙЛ или %s",
		l.Word(d.Origin), places.UserHint))
	pr(l.T("  Разряд из справочника не смягчает приговора: пороги, каталоги, ссылки и"))
	pr(l.T("  хранилища, адресуемые содержимым, судятся ядром ровно как прежде."))

	pr("")
	pr(fmt.Sprintf("  %-9s %10s %8s  %-34s %s", l.T("разряд"), l.T("размер"), l.T("файлов"),
		l.T("имя"), l.T("путь")))
	shown, hidden := 0, 0
	for _, f := range found {
		if !f.Applies || !f.Exists {
			continue
		}
		if top > 0 && shown == top {
			hidden++
			continue
		}
		pr(fmt.Sprintf("  %-9s %10s %8s  %-34s %s", l.Word(string(f.Class)), l.Bytes(f.Bytes),
			l.Num(int64(f.Files)), cut(f.Name, 34), cut(f.Resolved, 60)))
		shown++
	}
	if shown == 0 {
		pr(l.T("  — ни одно место справочника на этой машине не нашлось"))
	}
	if hidden > 0 {
		pr(l.F("  …и ещё %d найденных мест — весь список: --top 0, или --json", hidden))
	}

	pr("")
	pr(l.F("  Чего здесь нет — тоже ответ: %d мест справочник называет, а на этой машине", absent))
	pr(l.T("  их каталогов не существует. Это не ошибка: справочник общий, а машины разные."))
}

// History prints the корзины digitdisk has left under a root: when, how much,
// and how to get it back.
func History(w io.Writer, l lang.Lang, h *clean.History, now time.Time, top int) {
	pr := func(line string) { fmt.Fprintln(w, line) }

	pr(l.F("ЖУРНАЛ УБОРКИ  %s", h.Root))
	pr(l.F("  хранилище корзин  %s", h.Trash))
	pr(l.F("  корзин            %d", h.Boxes))
	pr(l.F("  лежит в корзинах  %s — место ещё НЕ освобождено", l.Bytes(h.MovedBytes)))
	pr(l.F("  освобождено       %s — это стёртое, и вернуть его нечем", l.Bytes(h.FreedBytes)))
	pr(l.F("  возвращено        %s", l.Bytes(h.RestoredBytes)))

	if h.Boxes == 0 {
		pr("")
		pr(l.F("  Уборки здесь не было: ни одной корзины под %s.", h.Trash))
		return
	}

	pr("")
	pr(fmt.Sprintf("  %-22s %8s %10s %8s %8s  %s", l.T("когда"), l.T("в корзине"), l.T("байт"),
		l.T("возвращ."), l.T("стёрто"), l.T("корзина")))
	shown, hidden := 0, 0
	for _, e := range h.Entries {
		if top > 0 && shown == top {
			hidden++
			continue
		}
		when := dash(e.StartedAt)
		if age := e.Age(now); age >= 0 {
			when = l.F("%s (%s назад)", l.StampDate(e.StartedAt), l.Since(age))
		}
		pr(fmt.Sprintf("  %-22s %8s %10s %8s %8s  %s", when, l.Num(int64(e.Moved)),
			l.Bytes(e.MovedBytes), l.Num(int64(e.Restored)), l.Num(int64(e.Purged)), cut(e.Box, 46)))
		if !e.Problem.Empty() {
			pr(l.F("      беда: %s", e.Problem.In(l)))
		}
		shown++
	}
	if hidden > 0 {
		pr(l.F("  …и ещё %d корзин — весь список: --top 0, или --json", hidden))
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
		pr(l.T("  Возвращать нечего: во всех корзинах пусто."))
		return
	}
	pr(l.F("  вернуть последнее:  digitdisk restore %s", newest.Box))
	// --confirm names a number a person retypes, so it is printed plain.
	pr(l.F("  стереть насовсем:   digitdisk purge %s --confirm %d", newest.Box, newest.Moved))
}
