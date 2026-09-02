# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
#
# Раскладка экрана на flang: печать библиотеки flang-tui в Go для хозяина.
#
# ── ЧТО ЗДЕСЬ ЛЕЖИТ ─────────────────────────────────────────────────────────
# flang-tui/  — ПОДМОДУЛЬ, клон github.com/digitable-lol/flang-tui на закреплённом
#               отпечатке. Исходник библиотеки НЕ скопирован сюда, а подключён:
#               `git submodule status` называет и адрес, и отпечаток, а
#               `git -C flang-tui log` показывает, откуда взялась каждая строка.
# out-go/     — ПЕЧАТАЕТСЯ, а не правится руками. Тот же уговор, что у
#               `core/out-go` (AGENTS.md): коммитится напечатанным, чтобы сборка
#               из чистого клона требовала только Go, и восстанавливается
#               целиком одной командой `make -C ui-flang печать`.
#
# ── ЗАЧЕМ ПЕРЕИМЕНОВАНИЕ МОДУЛЯ ─────────────────────────────────────────────
# У всех целей печати имя модуля Go одно — «flangprogram», — а библиотек здесь
# шесть, и в одной сборке они столкнулись бы именами. Поэтому после печати имя
# переписывается на «flang<имя>». Это ЕДИНСТВЕННАЯ правка напечатанного, и
# делает её `sed`, а не рука; ровно то же делает `make печать` у самой
# библиотеки.
#
# ── КАК ОБНОВИТЬ БИБЛИОТЕКУ ─────────────────────────────────────────────────
#   git -C ui-flang/flang-tui fetch && git -C ui-flang/flang-tui checkout <отпечаток>
#   make -C ui-flang            # проверка, печать, сборка
#   git add ui-flang            # новый отпечаток подмодуля и новая печать вместе

ifdef FLANG_HOME
FLANG ?= $(FLANG_HOME)/bootstrap/flang
else
FLANG ?= flang
endif

БИБЛИОТЕКА = flang-tui
МОДУЛИ = screen colour format history tabs scroll

.PHONY: all проверка печать сверка чисто

all: проверка печать

# `check` — разбор, типы, завершаемость и постусловия; `test` — примеры.
# Порознь, потому что порознь и падают: у них разные причины.
проверка:
	@for m in $(МОДУЛИ); do \
	  echo "== $$m"; \
	  $(FLANG) check $(БИБЛИОТЕКА)/flang/$$m.flang || exit 1; \
	  $(FLANG) test  $(БИБЛИОТЕКА)/flang/$$m.flang || exit 1; \
	done

печать:
	rm -rf out-go
	@for m in $(МОДУЛИ); do \
	  echo "== печать $$m"; \
	  $(FLANG) emit $(БИБЛИОТЕКА)/flang/$$m.flang --target go --out out-go/$$m || exit 1; \
	  grep -rl flangprogram out-go/$$m | while read f; do \
	    sed -i.bak "s|flangprogram|flang$$m|g" "$$f" && rm -f "$$f.bak"; \
	  done; \
	  rm -rf out-go/$$m/cli out-go/$$m/Makefile; \
	  ( cd out-go/$$m && test -z "$$(gofmt -l .)" && go vet ./... && go build ./... ) || exit 1; \
	done

# Сверка: перепечатать во временный каталог и сравнить с деревом. Расхождение
# значит, что кто-то правил вывод компилятора руками.
сверка:
	@rm -rf .перепечатка && mkdir -p .перепечатка
	@for m in $(МОДУЛИ); do \
	  $(FLANG) emit $(БИБЛИОТЕКА)/flang/$$m.flang --target go --out .перепечатка/$$m || exit 1; \
	  grep -rl flangprogram .перепечатка/$$m | while read f; do \
	    sed -i.bak "s|flangprogram|flang$$m|g" "$$f" && rm -f "$$f.bak"; \
	  done; \
	  rm -rf .перепечатка/$$m/cli .перепечатка/$$m/Makefile; \
	  diff -r out-go/$$m .перепечатка/$$m || exit 1; \
	done
	@rm -rf .перепечатка
	@echo "печать совпадает с деревом"

чисто:
	rm -rf out-go .перепечатка
