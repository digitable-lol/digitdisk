# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
#
# Ядро digitdisk: исходник на flang → AST с постусловиями → печать в Go и в C.
#
# ПОЧЕМУ ПЕЧАТЬ ИДЁТ ИЗ opis-diska.ast.json, А НЕ ИЗ opis-diska.flang.
# Постусловия функции (flang/SPEC.md §5) — часть контракта AST, но поверхность
# языка формы для них не имеет: `parseFunction` (flang/src/parser.mjs:1115)
# знает у функции только «принимает», «возвращает», «убывает», «пример» и тело,
# а слово «свойство» в теле даёт FLANG_PARSE. Поэтому инварианты И1 и И2
# приписываются к AST (core/postconditions.mjs) — и оттуда они попадают и в
# интерпретатор, и в напечатанный Go, и в напечатанный C одним и тем же кодом
# отказа DIGITDISK_RULE. Печать без постусловий (прямо из .flang) тоже работает
# и даёт те же значения — просто без сторожа.
#
# Каталоги out-go/ и out-c/ ПЕЧАТАЮТСЯ, а не правятся руками: `make -C core`
# восстанавливает их целиком.

FLANG ?= node /home/b/projects/flang/flang/bin/flang.mjs
ИСХОДНИК = opis-diska.flang
AST = opis-diska.ast.json

.PHONY: all проверка печать сверка чисто

all: проверка печать сверка

$(AST): $(ИСХОДНИК) postconditions.mjs
	node postconditions.mjs

проверка: $(AST)
	$(FLANG) check $(ИСХОДНИК)
	$(FLANG) test $(ИСХОДНИК)
	$(FLANG) check $(AST)
	$(FLANG) test $(AST)

печать: $(AST)
	rm -rf out-go out-c
	$(FLANG) emit --target go $(AST) --out out-go
	$(FLANG) emit --target c $(AST) --out out-c
	cd out-go && gofmt -l . && go vet ./... && go build ./... && go build -o flang_cli ./cli
	$(MAKE) -C out-c

сверка:
	node diff.mjs --описей 400 --си 2400

чисто:
	rm -rf out-go out-c $(AST)
