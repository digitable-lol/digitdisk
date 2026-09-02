# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
#
# Ядро digitdisk: исходник на flang → проверка, примеры, печать в Go и в C.
#
# ── ЧТО ИЗМЕНИЛОСЬ 1 СЕНТЯБРЯ 2026 ──────────────────────────────────────────
# Печать шла из opis-diska.ast.json, а не из исходника, потому что «поверхность
# flang формы для постусловия функции НЕ ИМЕЕТ» и постусловия приходилось
# приписывать к AST скриптом на JavaScript (core/postconditions.mjs).
# На стволе flang форма ЕСТЬ — слово `обеспечивает`, — и проверяет её сам
# компилятор. Постусловия переехали в исходник, а скрипт, AST и переходник
# core/flang-home.mjs сняты: печать идёт прямо из .flang.
#
# Снят и core/diff.mjs: сверка интерпретатора с напечатанным кодом теперь
# написана на самом языке — core/differential.flang, план для `flang io`.
#
# ── ЧЕМ СОБИРАЕТСЯ ──────────────────────────────────────────────────────────
# Компилятор flang — ОДИН двоичный файл, которому нужен только `cc`; ни Node,
# ни JavaScript тут больше нет. Взять его можно тремя дорогами, и все три дают
# один и тот же файл под одним именем:
#
#   brew install flang        asdf plugin add flang        make -C bootstrap
#
# Каталоги out-go/ и out-c/ ПЕЧАТАЮТСЯ, а не правятся руками: `make -C core`
# восстанавливает их целиком.

# Команда flang. По умолчанию — та, что стоит в системе; FLANG=путь/к/flang
# переопределяет её, а FLANG_HOME указывает на клон flang с собранным семенем.
ifdef FLANG_HOME
FLANG ?= $(FLANG_HOME)/bootstrap/flang
else
FLANG ?= flang
endif

ИСХОДНИК = disk-inventory.flang
СВЕРКА = differential.flang

.PHONY: all проверка печать сверка чисто

all: проверка печать сверка

# `check` — разбор, типы, завершаемость и постусловия; `test` — примеры.
# Порознь, потому что порознь и падают: у них разные причины.
проверка:
	$(FLANG) check $(ИСХОДНИК)
	$(FLANG) test $(ИСХОДНИК)

печать:
	rm -rf out-go out-c
	$(FLANG) emit --target go $(ИСХОДНИК) --out out-go
	$(FLANG) emit --target c $(ИСХОДНИК) --out out-c
	cd out-go && gofmt -l . && go vet ./... && go build ./... && go build -o flang_cli ./cli
	$(MAKE) -C out-c

# Сверка идёт `flang io`: план запускает напечатанные двоичные и сверяет их
# ответы с ответами интерпретатора. Расхождения печатаются числом.
сверка: печать
	$(FLANG) check $(СВЕРКА)
	$(FLANG) test $(СВЕРКА)
	$(FLANG) io $(СВЕРКА)

чисто:
	rm -rf out-go out-c
