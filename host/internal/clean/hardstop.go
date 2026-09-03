// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"os"
	"path/filepath"
	"strings"

	"digitdisk/internal/lang"
	"digitdisk/internal/protect"
	"digitdisk/internal/settings"
)

// # Твёрдые запреты
//
// Забой стирает то, на что человек указал, и приговор ядра ему не указ: ядро
// там спрашивается за СЛОВО, а не за разрешение (см. core.Naturer).  Ровно
// поэтому список мест, куда указывать нельзя вовсе, обязан быть коротким,
// явным и объяснимым — иначе он превратится во второй приговор, только
// написанный на коленке и без доказательств.
//
// В списке пять пунктов, и все пять об одном: НЕ «это жалко», а «после этого
// не станет машины или не станет инструмента, который сейчас работает».
//
//  1. корень файловой системы;
//  2. системный каталог целиком — /usr, /etc, /var и прочие поимённо;
//  3. домашний каталог целиком;
//  4. каталог самого digitdisk и его настройки;
//  5. корзина digitdisk — её опустошает `purge`, у которого есть журнал.
//
// ЧЕГО ЗДЕСЬ НЕТ И НЕ БУДЕТ.  Здесь нет запрета «это не мусор»: именно он и
// был бедой, ради которой всё переделывалось.  Здесь нет и «это чужое» или
// «это большое» — на такое отвечает природа и строгость вопроса, а не отказ.
//
// ГРАНИЦА ЗАПРЕТА — ЦЕЛОЕ, А НЕ ЛИСТ.  Запрещён каталог САМ и всякий, кто его
// в себе содержит; то, что лежит ВНУТРИ, не запрещено.  `/var` снести нельзя,
// `/var/tmp/моя-сборка` — можно, и это не послабление: машину ломает вынос
// каталога целиком, а отдельный файл внутри него — дело хозяина машины, и
// разговаривает с ним про этот файл природа, а не запрет.
//
// НА КАЖДЫЙ ОТКАЗ — ПРИЧИНА И ОБХОД.  Отказ без «почему» читается как каприз
// инструмента, а отказ без «как быть, если я всё же прав» заставляет искать
// обход мимо инструмента — то есть `rm -rf` вслепую, ровно то, от чего digitdisk
// и заводился.

// Stop is one твёрдый запрет: the path it refused, why, and the way past it.
type Stop struct {
	Path   string      `json:"путь"`
	Why    lang.Phrase `json:"почему"`
	Around lang.Phrase `json:"как_обойти"`
}

// Empty reports whether nothing was refused.
func (s Stop) Empty() bool { return s.Why.Empty() }

// Err is the refusal as it is shown: the path, the reason, and the way past
// it.  All three, always, in that order — an отказ without «почему» reads as a
// caprice of the tool, and an отказ without «как быть, если я всё же прав»
// sends a person looking for a way ROUND the tool, which is `rm -rf` typed
// blind and is exactly what digitdisk exists instead of.
func (s Stop) Err() error {
	return lang.Errorf(`СТИРАТЬ ОТСЮДА НЕЛЬЗЯ: %s

%s.

Как быть, если вы всё же правы: %s`, s.Path, s.Why, s.Around)
}

// ProtectStop is the защитный список refusing the ground, in the same shape as
// a твёрдый запрет.
//
// It is NOT one of the five: those are about the machine and the tool, this is
// a standing order from whoever holds the disk (see internal/protect for why
// the two live apart).  It wears the same shape because a person meeting a
// refusal wants the same three things either way — what was refused, why, and
// what to do if they meant it — and a screen that laid out one of them nicely
// and the other as a single long line would be saying that one of them matters
// less.
func ProtectStop(path string, rule protect.Rule) Stop {
	return Stop{Path: path,
		Why:    lang.Say("этот путь держит защитный список: %s", lang.Raw(rule.String())),
		Around: lang.Say("снимите эту строку из защитного списка или уберите ключ --protect: список — ваше собственное распоряжение, и digitdisk не отменяет его сам")}
}

// StopOptions is what the hard stops need to know about this machine.  Every
// field is injectable so that a test can name a fake home and a fake tool
// directory instead of erasing under the real ones.
type StopOptions struct {
	// Home is the person's home directory; empty asks the system.
	Home string
	// Tool is the directory the running digitdisk lives in; empty asks the
	// system for the executable's path.
	Tool string
	// Getenv reads the environment; nil means os.Getenv.
	Getenv func(string) string
}

// systemRoots are the каталоги системы, by name and with the reason each is
// named.  The list is deliberately flat and readable: a rule nobody can recite
// is a rule nobody can check.
//
// Both systems are in one list on purpose.  A Linux tree does not contain
// /System and a mac does not contain /proc, so naming both costs nothing and
// keeps one list instead of two that drift.
var systemRoots = []string{
	// POSIX и Linux
	"/bin", "/sbin", "/lib", "/lib32", "/lib64", "/libx32",
	"/usr", "/etc", "/var", "/boot", "/dev", "/proc", "/sys", "/run",
	"/opt", "/srv", "/mnt", "/media", "/home", "/root",
	// macOS
	"/System", "/Library", "/Applications", "/private", "/Users", "/Volumes", "/cores", "/Network",
}

// HardStop answers whether забой may take this path as its ground at all.  An
// empty Stop means yes.
//
// It is asked about the GROUND — the каталог a person marked or stood on — and
// not about every file under it.  That is where it belongs: the ban is about
// what a whole place is, the walk under it is judged by природа, and asking a
// question about `/etc` a million times would only make the answer slower, not
// truer.
func HardStop(path string, opt StopOptions) Stop {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Stop{Path: path,
			Why:    lang.Say("путь не приводится к абсолютному: %s", err),
			Around: lang.Say("назовите путь заново — от корня, без «..»")}
	}
	abs = filepath.Clean(abs)

	if abs == string(filepath.Separator) {
		return Stop{Path: abs,
			Why:    lang.Say("это корень файловой системы: под ним лежит вся машина целиком"),
			Around: lang.Say("назовите каталог, а не корень — забой стирает то, на что указали, и на корень указывать нечем, кроме как всем сразу")}
	}

	for _, root := range systemRoots {
		if !covers(abs, root) {
			continue
		}
		return Stop{Path: abs,
			Why:    lang.Say("%s — системный каталог: его содержимое ставит менеджер пакетов, а не человек, и машина после его выноса не загрузится", root),
			Around: lang.Say("снимите пакет тем, кто его поставил (apt, dnf, brew), либо укажите каталог ВНУТРИ него: запрещён каталог целиком, а не то, что в нём лежит")}
	}

	home := opt.Home
	if home == "" {
		if h, err := os.UserHomeDir(); err == nil {
			home = h
		}
	}
	if home != "" {
		home = filepath.Clean(home)
		if covers(abs, home) {
			return Stop{Path: abs,
				Why:    lang.Say("%s — домашний каталог целиком: это всё, что у вас есть на этой машине", home),
				Around: lang.Say("укажите каталог внутри дома — забой возьмёт его вместе со всем содержимым; дом целиком не берётся никогда")}
		}
	}

	for _, own := range ownPlaces(opt, home) {
		if !covers(abs, own.path) {
			continue
		}
		return Stop{Path: abs, Why: own.why, Around: own.around}
	}

	for _, part := range strings.Split(filepath.ToSlash(abs), "/") {
		if part != TrashName {
			continue
		}
		return Stop{Path: abs,
			Why:    lang.Say("это корзина digitdisk (%s): в ней лежит и журнал того, что в неё попало", TrashName),
			Around: lang.Say("корзину опустошает `digitdisk purge <корзина> --confirm N` — он читает журнал и стирает ровно то, что там записано")}
	}
	return Stop{}
}

// ownPlace is one place that belongs to digitdisk itself.
type ownPlace struct {
	path   string
	why    lang.Phrase
	around lang.Phrase
}

// ownPlaces names what digitdisk must not erase because digitdisk is running
// out of it.  A tool that erases itself halfway through an erasure leaves half
// the work done and no journal saying which half.
func ownPlaces(opt StopOptions, home string) []ownPlace {
	var out []ownPlace

	tool := opt.Tool
	if tool == "" {
		if exe, err := os.Executable(); err == nil {
			if resolved, err := filepath.EvalSymlinks(exe); err == nil {
				exe = resolved
			}
			tool = filepath.Dir(exe)
		}
	}
	if tool != "" {
		out = append(out, ownPlace{path: filepath.Clean(tool),
			why:    lang.Say("здесь лежит сам digitdisk, который сейчас работает"),
			around: lang.Say("снимите инструмент так, как ставили (brew uninstall, rm самого файла), а не им самим: стирание, снёсшее себя на середине, дописать журнал уже нечем")})
	}
	if home != "" {
		dir, err := settings.Dir(settings.Options{Home: home, Getenv: opt.Getenv})
		if err == nil {
			out = append(out, ownPlace{path: filepath.Clean(dir),
				why:    lang.Say("здесь настройки digitdisk: язык, справочник мест и защитный список"),
				around: lang.Say("правьте эти файлы редактором — они текстовые; забой на них снёс бы и защитный список заодно")})
		}
		out = append(out, ownPlace{path: filepath.Join(home, ".config", "digitdisk"),
			why:    lang.Say("здесь прежний дом настроек digitdisk"),
			around: lang.Say("правьте эти файлы редактором — они текстовые; забой на них снёс бы и защитный список заодно")})
	}
	return out
}

// covers reports whether ground is the named place or contains it.  Erasing a
// place from above is erasing the place, and «я указал на родителя» is not a
// way around a ban on the child.
//
// The comparison is by whole components: `/varnish` neither is nor contains
// `/var`.
func covers(ground, place string) bool {
	if ground == place {
		return true
	}
	sep := string(filepath.Separator)
	return strings.HasPrefix(place+sep, strings.TrimSuffix(ground, sep)+sep)
}
