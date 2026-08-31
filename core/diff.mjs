#!/usr/bin/env node
// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause
/**
 * Дифференциальный прогон ядра «Опись диска».
 *
 * Одна и та же программа считается ДВАЖДЫ — интерпретатором flang и
 * напечатанным кодом, — и значения сверяются побайтово. Значения ездят
 * размеченным JSON (числа строкой), потому что JSON не знает ни NaN, ни
 * Infinity, ни знака нуля, а `Object.is` их различает.
 *
 * Сетка: размеры 0…10^11, возраст 0…4000 дней, все три вида, доля недоступных.
 * Пути собираются из примет всех шести разрядов и из путей без примет, чтобы
 * каждое правило срабатывало.
 *
 *   node core/diff.mjs [--описей N] [--семя N]
 *
 * Расхождения печатаются числом — заявить «ноль» без прогона нельзя.
 */
import { execFileSync } from "node:child_process"
import { existsSync } from "node:fs"
import { dirname, join } from "node:path"
import { fileURLToPath } from "node:url"

import { errorCode, evaluateFlang } from "/home/b/projects/flang/flang/src/compat.mjs"
import { variant } from "/home/b/projects/flang/flang/src/interpret.mjs"
import { loadProgram } from "/home/b/projects/flang/flang/bin/flang.mjs"

const здесь = dirname(fileURLToPath(import.meta.url))

const ключи = new Map()
for (let i = 2; i < process.argv.length; i += 2) ключи.set(process.argv[i], process.argv[i + 1])
const ОПИСЕЙ = Number(ключи.get("--описей") ?? 400)
const СЕМЯ = Number(ключи.get("--семя") ?? 20260831)
const ПО_C = Number(ключи.get("--си") ?? 10)

/* ───────────────────────── сетка входов ───────────────────────── */

/** mulberry32: прогон обязан воспроизводиться, значит генератор свой. */
function генератор(семя) {
  let состояние = семя >>> 0
  return () => {
    состояние = (состояние + 0x6d2b79f5) >>> 0
    let t = состояние
    t = Math.imul(t ^ (t >>> 15), t | 1)
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61)
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}
const случайное = генератор(СЕМЯ)
const выбрать = (список) => список[Math.floor(случайное() * список.length)]
const целое = (предел) => Math.floor(случайное() * (предел + 1))

const ПРИМЕТЫ = [
  "/home/u/.cache/pip/", "/var/cache/apt/", "/Users/u/Library/Caches/", "/tmp/",
  "/var/log/", "/home/u/logs/", "/home/u/app.",
  "/home/u/p/node_modules/", "/home/u/p/target/", "/home/u/p/build/", "/home/u/p/_build/", "/home/u/.gradle/",
  "/home/u/Downloads/", "/home/u/Загрузки/",
  "/home/u/docs/", "/home/u/видео/", "/srv/данные/", "/",
]
const ХВОСТЫ = ["a.bin", "b.log", "письмо.txt", "фильм.mkv", "iso.img", "x", "деталь.log", ""]
const ВИДЫ = ["Файл", "Каталог", "Ссылка"]

/* Пороги правил — 7, 30, 180 дней и гибибайт. Равномерная сетка попадает в
   них раз в тысячи входов, то есть почти НЕ проверяет границу: сдвиг порога на
   единицу такой прогон не ловит (проверено отрицательным контролем). Поэтому
   в двух случаях из пяти значение берётся из окрестности порога. */
const ГРАНИЦЫ_ВОЗРАСТА = [0, 1, 6, 7, 8, 29, 30, 31, 179, 180, 181, 3999, 4000]
const ГРАНИЦЫ_РАЗМЕРА = [0, 1, 1073741823, 1073741824, 1073741825, 100000000000]

function находка() {
  /* Размеры до 10^11 и с перекосом к мелочи: и порог крупного, и ноль обязаны
     попадаться, а не встречаться раз в тысячу входов. */
  const шкала = выбрать([0, 1, 1e3, 1e6, 1073741824, 1e11])
  const размер = случайное() < 0.4 ? выбрать(ГРАНИЦЫ_РАЗМЕРА) : шкала === 0 ? 0 : Math.floor(случайное() * шкала)
  return {
    путь: `${выбрать(ПРИМЕТЫ)}${выбрать(ХВОСТЫ)}`,
    размер,
    возраст_дней: случайное() < 0.4 ? выбрать(ГРАНИЦЫ_ВОЗРАСТА) : целое(4000),
    вид: variant(выбрать(ВИДЫ), {}),
    доступен: случайное() > 0.2,
  }
}

function опись() {
  const длина = целое(12)
  return Array.from({ length: длина }, находка)
}

/* ───────────────────────── значения на проводе ───────────────────────── */

const этоВариант = (значение) =>
  typeof значение === "object" && значение !== null && !Array.isArray(значение) &&
  typeof значение.variant === "string" && typeof значение.fields === "object" && значение.fields !== null

function закодировать(значение) {
  if (значение === null || значение === undefined) return null
  if (typeof значение === "boolean") return значение
  if (typeof значение === "number") return { n: Object.is(значение, -0) ? "-0" : String(значение) }
  if (typeof значение === "string") return { s: значение }
  if (Array.isArray(значение)) return { l: значение.map(закодировать) }
  if (этоВариант(значение)) {
    return { v: значение.variant, f: Object.entries(значение.fields).map(([имя, поле]) => [имя, закодировать(поле)]) }
  }
  if (typeof значение === "object") return { r: Object.entries(значение).map(([имя, поле]) => [имя, закодировать(поле)]) }
  throw new Error(`нечего кодировать: ${typeof значение}`)
}

function раскодировать(узел) {
  if (узел === null) return null
  if (typeof узел === "boolean") return узел
  if (Object.hasOwn(узел, "n")) return Number(узел.n)
  if (Object.hasOwn(узел, "s")) return узел.s
  if (Object.hasOwn(узел, "l")) return узел.l.map(раскодировать)
  if (Object.hasOwn(узел, "r")) {
    const запись = {}
    for (const [имя, поле] of узел.r) запись[имя] = раскодировать(поле)
    return запись
  }
  if (Object.hasOwn(узел, "v")) {
    const поля = {}
    for (const [имя, поле] of узел.f ?? []) поля[имя] = раскодировать(поле)
    return variant(узел.v, поля)
  }
  throw new Error(`нечего декодировать: ${JSON.stringify(узел)}`)
}

function тоЖе(левое, правое) {
  if (typeof левое !== "object" || левое === null || typeof правое !== "object" || правое === null) {
    return Object.is(левое, правое)
  }
  if (Array.isArray(левое) || Array.isArray(правое)) {
    if (!Array.isArray(левое) || !Array.isArray(правое) || левое.length !== правое.length) return false
    return левое.every((элемент, номер) => тоЖе(элемент, правое[номер]))
  }
  if (этоВариант(левое) || этоВариант(правое)) {
    if (!этоВариант(левое) || !этоВариант(правое)) return false
    return левое.variant === правое.variant && тоЖе(левое.fields, правое.fields)
  }
  const слева = Object.keys(левое).sort()
  const справа = Object.keys(правое).sort()
  if (слева.length !== справа.length) return false
  if (!слева.every((имя, номер) => имя === справа[номер])) return false
  return слева.every((имя) => тоЖе(левое[имя], правое[имя]))
}

/* ───────────────────────── два движка ───────────────────────── */

const программа = await loadProgram(join(здесь, "opis-diska.ast.json"))

function интерпретатор(имя, аргументы) {
  try {
    return { ok: true, value: evaluateFlang(программа, имя, аргументы) }
  } catch (ошибка) {
    return { ok: false, code: errorCode(ошибка), message: ошибка instanceof Error ? ошибка.message : String(ошибка) }
  }
}

function прогонщик(путь, запросы) {
  const ввод = `${запросы.map((запрос) => JSON.stringify(запрос)).join("\n")}\n`
  const вывод = execFileSync(путь, { input: ввод, encoding: "utf8", maxBuffer: 512 * 1024 * 1024 })
  const строки = вывод.split("\n").filter((строка) => строка.length > 0)
  if (строки.length !== запросы.length) throw new Error("прогонщик ответил не на каждый запрос")
  return строки.map((строка) => {
    const ответ = JSON.parse(строка)
    return ответ.ok
      ? { ok: true, value: раскодировать(ответ.value) }
      : { ok: false, code: ответ.code, message: ответ.message }
  })
}

const тотЖеИсход = (левое, правое) =>
  левое.ok !== правое.ok ? false : левое.ok ? тоЖе(левое.value, правое.value) : левое.code === правое.code && левое.message === правое.message

const опиши = (исход) => (исход.ok ? JSON.stringify(закодировать(исход.value)) : `${исход.code}: ${исход.message}`)

/* ───────────────────────── прогон ───────────────────────── */

/** Один вход = одна пара (функция, аргументы). */
const входы = []
for (let номер = 0; номер < ОПИСЕЙ; номер += 1) {
  const записи = опись()
  входы.push({ fn: "Решить всё", args: { записи } })
  входы.push({ fn: "Свести", args: { записи } })
  входы.push({ fn: "Отчёт", args: { записи } })
  входы.push({ fn: "Сумма размеров убираемых", args: { записи } })
  const одна = находка()
  входы.push({ fn: "Решить находку", args: { находка: одна } })
  входы.push({ fn: "Разряд находки", args: { находка: одна } })
}

function сверить(двоичный, входы, ярлык) {
  const порядок = программа.functions
  const запросы = входы.map(({ fn, args }) => {
    const функция = порядок.find((элемент) => элемент.name === fn)
    return { fn, args: функция.params.map((параметр) => закодировать(args[параметр.name])) }
  })
  const чужие = прогонщик(двоичный, запросы)
  const расхождения = []
  входы.forEach(({ fn, args }, номер) => {
    const свой = интерпретатор(fn, args)
    if (!тотЖеИсход(свой, чужие[номер])) {
      расхождения.push({ ярлык, fn, args: JSON.stringify(закодировать(args)), интерпретатор: опиши(свой), напечатанное: опиши(чужие[номер]) })
    }
  })
  return расхождения
}

const goCli = join(здесь, "out-go", "flang_cli")
const cCli = join(здесь, "out-c", "flang_cli")
if (!existsSync(goCli)) throw new Error(`нет ${goCli}: соберите core/out-go`)

const расхожденияGo = сверить(goCli, входы, "go")
const расхожденияC = existsSync(cCli) ? сверить(cCli, входы.slice(0, ПО_C), "c") : null

for (const строка of [...расхожденияGo, ...(расхожденияC ?? [])].slice(0, 5)) {
  process.stderr.write(`${JSON.stringify(строка, null, 2)}\n`)
}

process.stdout.write(
  `${JSON.stringify({
    семя: СЕМЯ,
    описей: ОПИСЕЙ,
    входов: входы.length,
    "расхождений go": расхожденияGo.length,
    "входов по C": расхожденияC === null ? 0 : Math.min(ПО_C, входы.length),
    "расхождений c": расхожденияC === null ? "прогонщик C не собран" : расхожденияC.length,
  }, null, 2)}\n`,
)
process.exit(расхожденияGo.length + (расхожденияC?.length ?? 0) === 0 ? 0 : 1)
