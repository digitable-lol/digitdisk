# Заготовка записи digitdisk в каталоге инструментов портала

Заплатка, а не правка: ничего из описанного здесь **не применено** к общему
чекауту `courses` — там одновременно работают другие сессии.

## Провенанс номеров строк

Все пути и строки — по стволу **`origin/master` репозитория `courses`, коммит
`96a5bb0a` (2026-08-31)**. Читать так:

```bash
git -C /home/b/projects/courses show origin/master:data/tool-catalog.toml
```

**Рабочий чекаут `/home/b/projects/courses` стоит на ветке
`qa/ml-claims-into-master`, и в ней этих файлов нет вовсе** — ни
`data/tool-catalog.toml`, ни `scripts/check-portal-tools.mjs`. Локальный
`master` там тоже отстал. Чекаут `/home/b/projects/courses-tools` содержит
файлы, но его `master` (`1ecf45c9`, 2026-08-19) старше ствола: набор иконок в
нём 36 против 38, «Уробороса» в каталоге нет. Правку класть на `origin/master`.

Образец, по которому эта заготовка построена, — коммит `52e09cd5` «Уроборос
доехал до master»: последняя запись инструмента, доехавшая до ствола целиком.

## Файлов правится шесть, а не один и не четыре

| # | Файл | Что делаем | Обязательно? |
|---|---|---|---|
| 1 | `data/tool-catalog.toml` | новая запись `[[tools]]` + два числа в шапке | да |
| 2 | `i18n/ru.toml` | блок `[tools.tool.digitdisk]` с `name`, `hint` | да |
| 3 | `i18n/en.toml` | тот же блок по-английски | да |
| 4 | `layouts/digitdisk/list.html` | новый файл: тип страницы у Hugo = имя раздела | да |
| 5 | `content/digitdisk/_index.md` (+ `_index.en.md`) | страница, на которую ведёт `url` | да |
| 6 | `themes/github-style/static/css/portal-tools.css` | акцент карточки | **нет** |

Шестой файл в брифе назван обязательным — он не обязателен. Запасной цвет
стоит в `portal-tools.css:240` (`--tool-accent: var(--course-cyan)`), а
именованные акценты — строки 243–249, и `ouroboros` среди них **отсутствует**:
инструмент доехал до ствола без своего акцента и живёт на запасном. Правка
файла — вопрос вкуса, а не прохождения сторожа.

## 1. `data/tool-catalog.toml`

Запись кладётся после `digitwm` (строки 92–96) — порядок в файле
повествовательный, а не алфавитный (`# Порядок инструментов остался прежним (он
же порядок повествования)`, строки 23–26):

```toml
# digitdisk стоит последним в инженерной среде: он читает машину, на которой
# всё остальное работает, и ничего на ней не меняет. Английского раздела пока
# нет, поэтому localized = false — как у digitwm.
[[tools]]
id = "digitdisk"
url = "/digitdisk/"
icon = "database"
localized = false
```

Плюс два числа в шапке файла, которые ведутся руками:

- строка 7: `# ЗДЕСЬ НЕ ВСЕ ОДИННАДЦАТЬ ИНСТРУМЕНТОВ` → `ДВЕНАДЦАТЬ`;
- строка 15: `# У шести продуктовых — FTS, «Архитектор», Workbench, «Уроборос»,
  Digit, digitwm` → `У семи продуктовых — …, digitwm, digitdisk`.

`accent` в записи не объявляется: поле простояло в файле один день и не
читалось ни одним шаблоном (строки 38–43 того же файла).

## 2–3. Оба словаря

`i18n/ru.toml`, вслед за блоком `[tools.tool.digitwm]` (строки 1485–1487):

```toml
[tools.tool.digitdisk]
name = "digitdisk"
hint = "куда ушло место и как себя чувствует машина"
```

`i18n/en.toml`, вслед за `[tools.tool.digitwm]` (строки 1366–1368):

```toml
[tools.tool.digitdisk]
name = "digitdisk"
hint = "where the disk went and how the machine feels"
```

Английский блок обязателен, даже при `localized = false`: сторож требует
`name` и `hint` в **обоих** словарях (`scripts/check-portal-tools.mjs:156-160`).
`eyebrow` и `menu` не заводятся — `menu` есть только у прикладных инструментов
(строки 182–199), а лишний ключ `tools.tool.*` уронит прогон на проверке
мёртвых ключей (строки 206–214).

## 4–5. Раздел, иначе ссылка ведёт в пустоту

`layouts/digitdisk/list.html` — по образцу `layouts/ouroboros/list.html`:

```go-html-template
{{ define "content" }}
{{ partial "post.html" . }}
{{ end }}
```

Без него раздел уедет на `layouts/_default/list.html` — постраничный список
материалов, который у раздела из одной страницы напечатает пустоту.

`content/digitdisk/_index.md` — сама страница. Её требует не сторож, а набор
BDD: `features/tools-hub.feature:31-34`, сценарий «Каталог печатает все
инструменты портала», шаг «каждая ссылка элемента `#tools-catalog .tool-card`
ведёт на существующую страницу».

## Что именно требует сторож `scripts/check-portal-tools.mjs`

| Строка | Требование |
|---|---|
| 154 | `url` начинается и кончается на `/` — путь от корня сайта |
| 155 | `icon` — имя из закрытого набора `layouts/partials/icon.html` |
| 156–160 | `tools.tool.<id>.name` и `.hint` есть **в обоих** словарях |
| 178–214 | обратная сторона: ключ `tools.tool.*` без инструмента роняет прогон |
| 270–279 | `localized` сверяется со сборкой: `public/en<url>/index.html` обязан существовать ровно тогда, когда `localized = true` |
| 132–142 | строку в панели шапки писать руками запрещено — она собирается тем же партиалом |

`eyebrow`, `accent`, страница в `content/tools/` и запись в
`portal-tools.css` сторожем **не требуются**.

## Иконка: набор закрыт, диска в нём нет

`layouts/partials/icon.html:15`, 38 имён (фигуры — `layouts/partials/
icon-sprite.html`, 39 символов вместе с запасным):

```
home file play tag users search menu close info arrow-right arrow-left
arrow-down chevron-down clock eye list message external bookmark braces layout
database network sparkles messages shield book sun moon sliders settings
maximize cube panel-left coffee headphones tree bell
```

Заняты: `braces` (FTS), `network` («Архитектор»), `cube` (Workbench), `eye`
(«Уроборос»), `sparkles` (Digit), `panel-left` (digitwm), `layout` (Квадрат
Декарта), `file` (резюме), `tree` (S.C.O.R.E.), `messages` (STAR), `sliders`
(сильные стороны).

**Ближайшее свободное — `database`**: стопка цилиндров читается как хранилище,
и это единственная фигура набора про место, а не про действие. Чего в наборе
нет и чего не хватает именно этому инструменту: накопителя (диск/HDD), шкалы
заполнения (gauge, круговая диаграмма) и датчика живых показателей. Расширение
набора — правка двух файлов разом: список имён в `icon.html:15` и сама фигура в
`icon-sprite.html`; на неизвестное имя шаблон молча подставит запасной символ
(`icon.html:16`), а сторож упадёт на строке 155 раньше.

## Оба лицензионных гейта: что нужно и где ловушка

**`digitwm/tools/check-licensing.py`** ходит только по своему дереву; портала и
digitdisk он не касается. Его роль здесь — образец: наш `tools/check-licensing.py`
устроен по нему.

**`products/workbench/scripts/license-gate.mjs`** (пути ниже — от корня
`courses`, номера строк — того же коммита ствола):

1. **Строки 173–190, `FACTS`.** `digitdisk` там нет — значит его заявленная
   лицензия не проверяется ВООБЩЕ, ни в манифесте, ни в прозе. Ровно об этой
   дыре предупреждает комментарий на строках 165–172: каталог утилит стоял в
   манифесте под одним именем, в реестре под другим, «и заявленная лицензия не
   проверялась ни разу». Добавить: `digitdisk: "BSD-2-Clause",` в `FACTS` и
   `digitdisk: ["digitdisk"],` в `PROSE_ALIASES` (строки 420–434).
2. **Ловушка прозы (проверка 2г, строки 570–627).** Честный текст карточки
   скажет «ядро на flang» и «идея от mole под GPL-3.0». Если это окажется в
   ОДНОМ утверждении — абзаце, строке таблицы, записи TOML, — гейт упадёт:
   `flang` в `FACTS` стоит как `BSD-2-Clause`, и соседство с `GPL-3.0`
   читается как ложное заявление о лицензии flang. Разносить по разным
   предложениям или не называть лицензию mole в прозе портала вовсе.
3. **Проверка 1, строки 88–89 и 108–158.** Всё, что уезжает покупателю
   (`templates/`, `themes/`, `toolchain/`, `compendium/`, `README.md`,
   `LICENSE-PERSONAL.md`, `PRODUCT.json`), не имеет права называть копилефт без
   оговорки из `PROSE_MARKERS` (строки 64–71): «не форк», «в архив не входит»,
   «не поставляется», «отдельным процессом», `"shipped": false`. Файлы
   digitdisk эту оговорку уже несут: `NOTICE` — «не форк… в дерево не входит и
   покупателю не поставляется», `README.md` — «not a fork of mole»,
   `README.ru.md` и `LICENSE-RU.md` — «это не форк mole».
4. **Строки 249–266.** Если digitdisk попадёт в `PRODUCT.json` как компонент,
   его `license` обязан быть `"BSD-2-Clause"`; `"shipped": false` не нужен —
   требование висит только на копилефтных компонентах.
5. **Строки 636–645, `SPDX_TARGETS`.** Строку
   `[join(SIBLINGS, "digitdisk/LICENSE"), "BSD-2-Clause", null],` добавлять
   безопасно: тело нашего `LICENSE` дословно совпадает с `flang/LICENSE`
   (`diff <(tail -n +4 LICENSE) <(tail -n +4 ../flang/LICENSE)` — пусто), а
   `flang/LICENSE` уже стоит в этом списке и сверку с эталоном проходит.
   Отличается только строка копирайта, а `normalise()` (строки 649–667) её
   вычёркивает.

## Порядок применения

1. Правка шести файлов выше на ветке от `origin/master`.
2. `node scripts/check-portal-tools.mjs` — сторож каталога.
3. `node scripts/check-i18n-keys.mjs` — мёртвые и недостающие ключи.
4. Сборка Hugo, затем `npx cucumber-js features/tools-hub.feature` — сценарий
   про существующие ссылки работает по собранному дереву, а проверка
   `localized` (строки 270–279 сторожа) — по `public/en`.
5. `node products/workbench/scripts/license-gate.mjs` — если digitdisk
   упоминается в прозе Workbench или в `PRODUCT.json`.
