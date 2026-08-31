# Сгенерировано flang (бэкенд Go). Флаги здесь — часть контракта бэкенда:
# напечатанный код обязан проходить go vet и gofmt без единого замечания.
GO ?= go

all: flang_cli

build:
	$(GO) build ./...

flang_cli:
	$(GO) build -o flang_cli ./cli

vet:
	$(GO) vet ./...

fmt:
	gofmt -l .

clean:
	rm -f flang_cli

.PHONY: all build vet fmt clean
