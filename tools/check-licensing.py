#!/usr/bin/env python3
# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
"""Лицензионный сторож digitdisk.

Дерево написано с нуля и целиком выходит под BSD-2-Clause; состав показателей
навеян mole, но ни строки его кода здесь нет, а имя и логотип mole не
используются (NOTICE, TRADEMARK.md апстрима). Проверки ниже держат все три
утверждения правдой — и держат зелёными двух сторожей стека, ради которых
лицензия и выбрана: «No GPL in the tree» в соседнем репозитории стека и запрет
копилефта в платном архиве Workbench.

  1. Копилефта в дереве нет. Это отказ с зубами: один заимствованный файл под
     GPL релицензирует всё дерево, и обнаруживается это обычно после
     публикации. Проза о копилефте разрешена — NOTICE обязан объяснить, чего мы
     не взяли и почему; запрещён файл, который *под* копилефтом, а не тот, что
     его называет.

  2. Наши исходники несут `SPDX-License-Identifier: BSD-2-Clause` и строку
     копирайта Marat Zimnurov. Набор выводится из дерева — все отслеживаемые
     исходники минус напечатанное, — поэтому новый файл без заголовка падает, а
     не становится тихо первым безымянным.

  3. Имя апстрима не присвоено. `TRADEMARK.md` mole запрещает форкам имя «Mole»
     и логотип; мы не форк, но требование соблюдаем буквально, и проверка
     делает нарушение невозможным: к моменту, когда его заметит человек, оно
     уже опубликовано.

  4. LICENSE держит формулу BSD-2-Clause и правообладателя, NOTICE записывает
     происхождение, а README живёт парой редакций. Проверки 1–3 что-то значат
     только против написанного текста.

Использование:  python3 tools/check-licensing.py
Выходит с кодом 1 и называет каждый файл, который не прошёл.
"""

import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent

# Напечатанное из flang. Освобождено от требования заголовка — печать его не
# ставит, а дописывать руками в вывод печати запрещено (AGENTS.md). От проверки
# на копилефт НЕ освобождено: печать чужого кода была бы тем же нарушением.
GENERATED = ("core/out-go/", "core/out-c/")

SOURCE_SUFFIXES = {".go", ".flang", ".py", ".sh", ".mjs"}

# Проза вправе обсуждать копилефт: NOTICE объясняет, почему кода mole здесь
# нет, а README называет лицензию, которую мы не взяли. Ни один из этих файлов
# не может нести лицензированный код; запрещён файл, который *является*
# копилефтным, а не тот, который его называет.
#
# Сам сторож освобождён по более резкой причине: он держит искомые строки
# внутри себя и совпадает со своим же поиском по построению. Без исключения он
# падал бы на себе в первом же прогоне после добавления в дерево.
SCAN_EXEMPT_SUFFIXES = {".md"}
SCAN_EXEMPT_NAMES = {"NOTICE", ".gitignore", "check-licensing.py"}

COPYLEFT_MARKERS = (
    "GNU General Public License", "GNU GENERAL PUBLIC LICENSE",
    "GNU Affero General Public License", "GNU Lesser General Public License",
    "GPL-2.0", "GPL-3.0", "AGPL-3.0", "LGPL-2.1", "LGPL-3.0",
)

# Знаки апстрима. Ловится присвоение имени, а не упоминание проекта: упоминания
# живут в NOTICE и README, которые из этой проверки исключены теми же списками.
UPSTREAM_MARKS = (
    re.compile(r"\bMole\b"),
    re.compile(r"\btw93\b"),
    re.compile(r"mole\.fit"),
)

failures = []
passed = []


def tracked_files():
    out = subprocess.run(
        ["git", "ls-files"],
        cwd=ROOT, capture_output=True, text=True, check=True,
    ).stdout
    return [line for line in out.splitlines() if line]


def read(relative):
    return (ROOT / relative).read_text(encoding="utf-8", errors="replace")


def header_of(relative):
    """Первые строки файла — там, где заголовок есть или где его нет.

    Намеренно не «первые 2 КБ»: собственная строка документации этого файла
    цитирует `SPDX-License-Identifier: BSD-2-Clause`, объясняя правило, и
    сканер по 2 КБ принял бы цитату за заголовок. Заголовок живёт наверху или
    он не заголовок.
    """
    return "\n".join(read(relative).split("\n")[:15])


def scanned(relative):
    path = Path(relative)
    return not (path.suffix in SCAN_EXEMPT_SUFFIXES or path.name in SCAN_EXEMPT_NAMES)


files = tracked_files()
if not files:
    print("FAIL  git ls-files ничего не вернул — это не клон?", file=sys.stderr)
    sys.exit(1)

# 1. Ни одного копилефтного файла под контролем версий.
copyleft_hits = []
for relative in files:
    if not scanned(relative):
        continue
    try:
        text = read(relative)
    except (OSError, UnicodeError):
        continue
    for marker in COPYLEFT_MARKERS:
        if marker in text:
            copyleft_hits.append(
                f"{relative}: содержит «{marker}» — копилефтный файл под контролем "
                f"версий релицензирует всё дерево и роняет оба сторожа стека. "
                f"Код mole сюда не копируется: идея — можно, код — нельзя."
            )
            break

if copyleft_hits:
    failures.extend(copyleft_hits)
else:
    passed.append(f"копилефта нет ни в одном из {len(files)} отслеживаемых файлов")

# 2. Наши исходники несут заголовок.
ours = [
    f for f in files
    if Path(f).suffix in SOURCE_SUFFIXES and not f.startswith(GENERATED)
]

missing_header = []
for relative in ours:
    head = header_of(relative)
    if "SPDX-License-Identifier: BSD-2-Clause" not in head:
        missing_header.append(
            f"{relative}: в первых строках нет «SPDX-License-Identifier: BSD-2-Clause» "
            f"(файл нашего авторства — NOTICE ставит весь код под BSD-2-Clause)"
        )
    elif "SPDX-FileCopyrightText:" not in head or "Marat Zimnurov" not in head:
        missing_header.append(
            f"{relative}: в первых строках нет «SPDX-FileCopyrightText: … Marat Zimnurov»"
        )

if missing_header:
    failures.extend(missing_header)
else:
    passed.append(f"все {len(ours)} наших исходника несут заголовок SPDX BSD-2-Clause")

# 3. Имя апстрима не присвоено.
claimed = []
for relative in files:
    if not scanned(relative):
        continue
    try:
        text = read(relative)
    except (OSError, UnicodeError):
        continue
    for mark in UPSTREAM_MARKS:
        if mark.search(text):
            claimed.append(
                f"{relative}: несёт знак апстрима «{mark.pattern}» — TRADEMARK.md mole "
                f"запрещает имя «Mole» и логотип форкам, а мы даже не форк. "
                f"Имя инструмента — digitdisk."
            )
            break

if claimed:
    failures.extend(claimed)
else:
    passed.append("ни один исходник не присваивает имя и знаки апстрима")

# 4. LICENSE, NOTICE и пара README говорят то, что предполагают проверки выше.
license_text = read("LICENSE")
if "Redistribution and use in source and binary forms" not in license_text:
    failures.append("LICENSE: формула BSD-2-Clause пропала")
elif "BSD 2-Clause License" not in license_text.split("\n")[0]:
    failures.append(
        "LICENSE: первая строка не «BSD 2-Clause License» — приписка сверху "
        "превращает известную лицензию в «не определена» для детекторов"
    )
elif "Marat Zimnurov" not in license_text:
    failures.append("LICENSE: нет строки копирайта Marat Zimnurov")
else:
    passed.append("LICENSE держит дословную формулу BSD-2-Clause и правообладателя")

if not (ROOT / "NOTICE").exists():
    failures.append("NOTICE: отсутствует — происхождение показателей не записано")
else:
    notice = read("NOTICE")
    for required in ("mole", "GPL-3.0", "TRADEMARK", "BSD-2-Clause"):
        if required not in notice:
            failures.append(f"NOTICE: не упоминает {required}")
            break
    else:
        passed.append("NOTICE записывает происхождение идеи и лицензию нашего кода")

pair = {"README.md": "README.ru.md", "README.ru.md": "README.md"}
broken_pair = []
for page, other in pair.items():
    if not (ROOT / page).exists():
        broken_pair.append(f"{page}: редакции нет, а вторая на неё ссылается")
        continue
    first = read(page).split("\n")[0]
    if f"({other})" not in first:
        broken_pair.append(
            f"{page}: первая строка не ведёт на {other} — редакция, до которой "
            f"нельзя дойти со второй, расходится с ней молча"
        )

if broken_pair:
    failures.extend(broken_pair)
else:
    passed.append("обе редакции README на месте и переключаются друг на друга первой строкой")

for name in passed:
    print(f"ok    {name}")
for failure in failures:
    print(f"FAIL  {failure}", file=sys.stderr)

if failures:
    print(
        f"\n{len(failures)} licensing obligation(s) not met. See NOTICE.",
        file=sys.stderr,
    )
    sys.exit(1)

print(f"\n{len(passed)} checked, 0 failed.")
