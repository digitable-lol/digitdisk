// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

/**
 * Где лежит flang.
 *
 * Ядро печатается и сверяется чужим инструментом, и путь к нему не может быть
 * зашит: у каждого он свой. Порядок поиска — от явного к угадыванию, и если не
 * нашлось, отказ называет, что именно сделать, а не «module not found».
 *
 *   1. переменная окружения FLANG_HOME — корень клона flang;
 *   2. соседний каталог ../../flang рядом с этим репозиторием.
 *
 * Проверяется наличие bin/flang.mjs, а не самого каталога: пустой клон и
 * отсутствующий клон должны различаться на первом же шаге.
 */
import { existsSync } from "node:fs"
import { dirname, join, resolve } from "node:path"
import { fileURLToPath } from "node:url"

const здесь = dirname(fileURLToPath(import.meta.url))

/* Явно названный корень не подменяется соседним: молчаливая подмена означала
   бы, что сверка прошла не на том дереве, которое назвали, и об этом никто не
   узнал бы. Поэтому FLANG_HOME, если задан, — единственный кандидат. */
const кандидаты = process.env.FLANG_HOME
  ? [process.env.FLANG_HOME]
  : [resolve(здесь, "..", "..", "flang")]

function годен(корень) {
  return existsSync(join(корень, "flang", "bin", "flang.mjs"))
}

const найденный = кандидаты.map((к) => resolve(к)).find(годен)

if (!найденный) {
  const перечень = кандидаты.map((к) => `  ${resolve(к)}`).join("\n")
  throw new Error(
    "flang не найден. Искал bin/flang.mjs в:\n" +
      перечень +
      "\nУкажите корень клона flang переменной FLANG_HOME " +
      "(например: FLANG_HOME=~/projects/flang make -C core).",
  )
}

export const КОРЕНЬ = найденный
export const БИН = join(найденный, "flang", "bin", "flang.mjs")

/** Загружает модуль flang по пути внутри его дерева. */
export function модуль(относительный) {
  return import(join(найденный, относительный))
}
