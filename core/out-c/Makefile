# Сгенерировано flang (бэкенд C). Флаги здесь — часть контракта бэкенда:
# сгенерированный код обязан собираться без единого предупреждения.
CC ?= cc
CFLAGS ?= -std=c99 -Wall -Wextra -Werror -pedantic -O2
LDLIBS ?= -lm

OBJECTS = flang_runtime.o opis_diska.o

all: libopis_diska.a flang_cli

libopis_diska.a: $(OBJECTS)
	ar rcs $@ $(OBJECTS)

flang_cli: flang_cli.o $(OBJECTS)
	$(CC) $(CFLAGS) -o $@ flang_cli.o $(OBJECTS) $(LDLIBS)

clean:
	rm -f $(OBJECTS) flang_cli.o flang_cli libopis_diska.a

.PHONY: all clean
