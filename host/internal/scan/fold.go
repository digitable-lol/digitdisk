// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package scan

import (
	"io/fs"
	"path/filepath"
)

// ЗАЧЕМ ЭТОТ СПИСОК ЕСТЬ, ЧИСЛОМ.
//
// Замер на дереве портала courses (82 923 записи, 1,4 ГБ):
//
//	обход без решающего слоя          0,81 с
//	обход с решающим слоем на flang  17,5 с
//	из них один node_modules          9,9 с (50 384 записи)
//
// То есть девять десятых времени уходит на то, чтобы вынести пятьдесят тысяч
// отдельных приговоров о файлах внутри каталога, который человек удаляет
// целиком или не трогает вовсе. Свёртка спрашивает решающий слой ОДИН раз — о
// самом каталоге — и заряжает его приговором всё поддерево.
//
// ЧТО СВЁРТКА НЕ ДЕЛАЕТ. Она не пропускает каталог: байты сосчитаны все до
// одного, TotalBytes по-прежнему сходится с `du -sb`, жёсткие ссылки
// по-прежнему считаются один раз. Пропуск — это `--exclude` у чужих
// чистильщиков, и он врёт о размере; здесь размер честный, дешевле только
// приговор.
//
// ПРАВИЛО ОТБОРА, и оно строгое: сюда попадает имя, которое (1) принадлежит
// производным данным, восстановимым одной командой, (2) достаточно
// характерно, чтобы не совпасть со своим каталогом человека, и (3) обычно
// содержит тысячи записей. Поэтому здесь нет `build`, `dist`, `target` и
// `venv`: первые три — обычные имена рабочих каталогов, а `venv` слишком
// часто зовут иначе. Их свернуло бы имя, а не природа.
var foldNames = map[string]string{
	"node_modules":  "npm/pnpm/yarn восстанавливают его из локфайла: npm ci",
	"__pycache__":   "Python перепишет байткод при следующем импорте",
	".mypy_cache":   "кэш проверки типов, mypy соберёт заново",
	".pytest_cache": "кэш прогона тестов, pytest соберёт заново",
	".ruff_cache":   "кэш линтера, ruff соберёт заново",
	".gradle":       "кэш сборки Gradle в проекте",
	".tox":          "окружения tox, восстанавливаются прогоном",
	".next":         "кэш и вывод сборки Next.js",
	".nuxt":         "кэш и вывод сборки Nuxt",
	".turbo":        "кэш Turborepo",
	".parcel-cache": "кэш Parcel",
	".svelte-kit":   "вывод сборки SvelteKit",
}

// FoldName returns the reason a directory of this name is folded, and whether
// it is folded at all.  Matching is by the directory's own name, not by path:
// a node_modules three levels down is the same node_modules.
func FoldName(name string) (string, bool) {
	reason, ok := foldNames[name]
	return reason, ok
}

// FoldNames lists every folded name with its reason.  Sorted output is the
// caller's business; the map is small and the caller usually prints it once.
func FoldNames() map[string]string {
	out := make(map[string]string, len(foldNames))
	for name, reason := range foldNames {
		out[name] = reason
	}
	return out
}

// FoldByName is the Fold function the tool uses by default: fold a directory
// whose own name is in the list.
func FoldByName(path string, _ fs.FileInfo) bool {
	_, ok := foldNames[filepath.Base(path)]
	return ok
}
