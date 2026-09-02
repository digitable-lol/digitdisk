#!/bin/sh
# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
#
# Сверка двух сборок digitdisk: раскладку экрана считает рукописный Go против
# того же дерева, где её считает печатанный flang (признак `flangui`).
#
# ЗАЧЕМ ОТДЕЛЬНЫЙ СКРИПТ, А НЕ ПРОВЕРКА. Одна сборка держит одну реализацию:
# сравнить их изнутри `go test` нельзя, нужны ДВА прогона одного дерева с
# разными признаками. Поэтому кадры снимает дампер
# host/internal/ui/snapshot_test.go — по разу под каждым признаком, — а
# сравнивает их этот скрипт.
#
# ПОЧЕМУ НЕ ЭТАЛОННЫЙ ФАЙЛ В ДЕРЕВЕ. Эталон кадров протух бы на первом же новом
# разделе, и его правили бы вслепую. Здесь сравниваются два прогона ОДНОГО
# дерева: что бы ни нарисовал экран, обе сборки обязаны нарисовать это одинаково.
#
#   tools/sverka-ui.sh            снимки, вывод в трубу и --json
#   tools/sverka-ui.sh --замер    плюс цена кадра и нажатия в обеих сборках
#
# Нужны: go. Компилятор flang НЕ нужен — печать лежит в дереве.
set -eu

ROOT=$(cd "$(dirname "$0")/.." && pwd)
WORK="${WORK:-${TMPDIR:-/tmp}/digitdisk-sverka-ui}"
BENCH=no
for arg in "$@"; do
	case "$arg" in
	--замер | --bench) BENCH=yes ;;
	*)
		echo "неизвестный ключ: $arg" >&2
		exit 2
		;;
	esac
done

rm -rf "$WORK"
mkdir -p "$WORK"
BAD=0
HOMEDIR="$WORK/home"
mkdir -p "$HOMEDIR"

echo "== сборка обеих"
(cd "$ROOT/host" && go build -tags flangcore -o "$WORK/dd-go" .)
(cd "$ROOT/host" && go build -tags flangcore,flangui -o "$WORK/dd-flang" .)
echo "   собраны обе"

echo
echo "== снимки кадров: 4 глубины цвета × ширины 40/80/120/200 × 2 языка × все разделы"
(cd "$ROOT/host" && DIGITDISK_SNAPSHOT="$WORK/snap-go.txt" go test -count=1 -run TestSnapshotDump ./internal/ui/ >/dev/null)
(cd "$ROOT/host" && DIGITDISK_SNAPSHOT="$WORK/snap-flang.txt" go test -count=1 -tags flangui -run TestSnapshotDump ./internal/ui/ >/dev/null)
if cmp -s "$WORK/snap-go.txt" "$WORK/snap-flang.txt"; then
	echo "   СОВПАЛО байт в байт: строк $(wc -l <"$WORK/snap-go.txt"), байт $(wc -c <"$WORK/snap-go.txt")"
else
	echo "   РАЗОШЛОСЬ:"
	diff "$WORK/snap-go.txt" "$WORK/snap-flang.txt" | head -20
	BAD=$((BAD + 1))
fi

echo
echo "== вывод в трубу и --json на неподвижном дереве"
TREE="$WORK/tree"
mkdir -p "$TREE/обычный" "$TREE/пустой" "$TREE/путь/к/🙂/файлу" "$TREE/.cache/пример"
i=1
while [ "$i" -le 20 ]; do
	dd if=/dev/zero of="$TREE/обычный/файл-$i.bin" bs=1024 count="$i" 2>/dev/null
	i=$((i + 1))
done
dd if=/dev/zero of="$TREE/путь/к/🙂/файлу/большой🙂.bin" bs=1024 count=900 2>/dev/null
dd if=/dev/zero of="$TREE/.cache/пример/кэш.tmp" bs=1024 count=120 2>/dev/null
find "$TREE" -exec touch -t 202601020304.05 {} +

# Время НЕ раскладка: возраст находки и длительность прогона читаются с часов и
# расходятся у одной и той же сборки с самой собой. Их и только их закрывает маска.
mask() {
	sed -E 's/"возраст_дней": [0-9.]+/"возраст_дней": ВРЕМЯ/g; s/"duration_seconds": [0-9.e-]+/"duration_seconds": ВРЕМЯ/g; s/"taken_at": "[^"]*"/"taken_at": ВРЕМЯ/g'
}
for order in "analyze $TREE" "analyze $TREE --json" "clean $TREE" "clean $TREE --json" "places" "places --json" "history $TREE" "--help" "--version"; do
	for lang in ru en; do
		# shellcheck disable=SC2086
		HOME="$HOMEDIR" "$WORK/dd-go" --lang "$lang" $order 2>&1 | mask >"$WORK/a.txt" || true
		# shellcheck disable=SC2086
		HOME="$HOMEDIR" "$WORK/dd-flang" --lang "$lang" $order 2>&1 | mask >"$WORK/b.txt" || true
		label=$(printf %s "$order" | sed "s|$TREE|ДЕРЕВО|")
		if cmp -s "$WORK/a.txt" "$WORK/b.txt"; then
			printf '   СОВПАЛО   %-32s %s (%s байт)\n' "$label" "$lang" "$(wc -c <"$WORK/a.txt")"
		else
			printf '   РАЗОШЛОСЬ %-32s %s\n' "$label" "$lang"
			diff "$WORK/a.txt" "$WORK/b.txt" | head -10
			BAD=$((BAD + 1))
		fi
	done
done

echo
echo "== живой экран в pty: список команд на ширинах 40/80/120/200"
if command -v script >/dev/null 2>&1; then
	# Часы в шапке и «замер N назад» в подвале читаются с часов, а не из
	# раскладки; две сборки нельзя запустить в одну и ту же секунду.
	frame() {
		awk 'BEGIN{RS="\033\\[H"} {f[NR]=$0} END{print f[NR-1]}' "$1" |
			sed -E 's/[0-9]{2}:[0-9]{2}:[0-9]{2}/ЧЧ:ММ:СС/g; s/замер [^·]*назад/замер ВРЕМЯ назад/g; s/длился [0-9]+[^·]*/длился ВРЕМЯ /g'
	}
	shoot() {
		(
			sleep 5
			printf '?'
			sleep 2
			printf 'q'
			sleep 1
		) | env HOME="$HOMEDIR" TERM=xterm-256color script -q -c "stty rows $2 cols $3; $1 --lang ru status" /dev/null >"$4" 2>&1 || true
	}
	for cols in 40 80 120 200; do
		rows=24
		if [ "$cols" -ge 100 ]; then rows=50; fi
		shoot "$WORK/dd-go" "$rows" "$cols" "$WORK/pty-go.raw"
		shoot "$WORK/dd-flang" "$rows" "$cols" "$WORK/pty-flang.raw"
		frame "$WORK/pty-go.raw" >"$WORK/f-go.txt"
		frame "$WORK/pty-flang.raw" >"$WORK/f-flang.txt"
		if cmp -s "$WORK/f-go.txt" "$WORK/f-flang.txt"; then
			printf '   СОВПАЛО   %dx%d (%s байт кадра)\n' "$cols" "$rows" "$(wc -c <"$WORK/f-go.txt")"
		else
			printf '   РАЗОШЛОСЬ %dx%d\n' "$cols" "$rows"
			diff "$WORK/f-go.txt" "$WORK/f-flang.txt" | head -10
			BAD=$((BAD + 1))
		fi
	done
else
	echo "   script(1) не найден — живой экран пропущен"
fi

if [ "$BENCH" = yes ]; then
	echo
	echo "== цена кадра и нажатия: рукописный Go"
	(cd "$ROOT/host" && go test -run '^$' -bench Layout -benchmem ./internal/ui/ | grep -E '^Benchmark')
	echo "== цена кадра и нажатия: печатанный flang"
	(cd "$ROOT/host" && go test -run '^$' -bench Layout -benchmem -timeout 40m -tags flangui ./internal/ui/ | grep -E '^Benchmark')
fi

echo
if [ "$BAD" -eq 0 ]; then
	echo "сверка раскладки: расхождений 0"
else
	echo "сверка раскладки: РАСХОЖДЕНИЙ $BAD"
fi
echo "рабочий каталог остался в $WORK — убрать: rm -rf $WORK"
[ "$BAD" -eq 0 ]
