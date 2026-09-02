// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"

	"digitdisk/internal/core"
	"digitdisk/internal/lang"
)

// Метки выпуска. Их проставляет компоновщик из scripts/build-release.sh:
//
//	-X main.version=$(cat VERSION) -X main.commit=<хеш> -X main.date=<время коммита>
//
// Значения по умолчанию — не заглушка ради красоты, а утверждение: этот
// двоичный файл собран не выпуском. Врать «0.1.0» о сборке из рабочего дерева
// нельзя — тогда отчёт об ошибке назовёт версию, которой не соответствует ни
// один коммит.
var (
	version = "dev"
	commit  = ""
	date    = ""
)

// buildStamp возвращает хеш сборки и её время. Если компоновщик их не
// проставил (обычный `go build`), они берутся из записи о системе контроля
// версий, которую Go кладёт в двоичный файл сам. Поэтому хеш сборки есть
// всегда, а не только у выпуска.
func buildStamp() (hash, when string, dirty bool) {
	hash, when = commit, date
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return hash, when, false
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if hash == "" {
				hash = s.Value
			}
		case "vcs.time":
			if when == "" {
				when = s.Value
			}
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}
	return hash, when, dirty
}

// shortHash укорачивает хеш до семи знаков — столько печатает git по умолчанию,
// и столько человек переносит в отчёт об ошибке без ошибок.
func shortHash(h string) string {
	if len(h) > 7 {
		return h[:7]
	}
	return h
}

// printVersion печатает версию, хеш сборки и то, каким решающим слоем собран
// этот двоичный файл. Последнее — не украшение: сборка без признака
// `flangcore` считает, но не решает, и человек обязан видеть это до того, как
// удивится пустому разбору по разрядам.
func printVersion(w io.Writer, l lang.Lang) {
	hash, when, dirty := buildStamp()
	if hash == "" {
		hash = l.T("неизвестен")
	} else {
		hash = shortHash(hash)
		if dirty {
			hash += " " + l.T("(дерево с правками)")
		}
	}
	if when == "" {
		when = l.T("неизвестно")
	}

	d := chosenDecider(l)
	layer := l.Word(d.Name())
	if !d.Ready() {
		layer += " " + l.T("— собрано без признака flangcore")
	}

	fmt.Fprintf(w, "digitdisk %s\n", version)
	fmt.Fprintf(w, l.T("сборка          %s, %s")+"\n", hash, when)
	fmt.Fprintf(w, l.T("инструментарий  %s %s/%s")+"\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(w, l.T("решающий слой   %s, договор версии %d")+"\n", layer, core.ContractVersion)
	// The language of this run and who chose it.  A person who wonders why
	// the tool is speaking the language it is speaking asks this, and the
	// answer names the file to edit.
	fmt.Fprintf(w, l.T("язык            %s (%s)")+"\n", l.Name(), l.Word(string(langChoice.Source)))
}
