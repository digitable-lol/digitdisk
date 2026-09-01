#!/usr/bin/env bash
# SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
# SPDX-License-Identifier: BSD-2-Clause
#
# Сборка выпуска digitdisk: двоичные файлы под Linux x86-64 и arm64, архивы,
# контрольные суммы и готовая к копированию формула Homebrew.
#
#   scripts/build-release.sh [версия]
#
# Версия берётся из довода, иначе из переменной DIGITDISK_VERSION, иначе из
# файла VERSION в корне дерева. Ведущее «v» снимается: тег v0.1.0 и версия
# 0.1.0 — одно и то же.
#
# ПОЧЕМУ СБОРКА ПОВТОРИМА. Три условия, и все три проверяются здесь, а не
# обещаются в README:
#
#   1. Дерево чистое, и хеш коммита — единственный источник меток. Грязное
#      дерево останавливает сборку: архив, которому не соответствует ни один
#      коммит, нечем проверить.
#   2. Время берётся у коммита (SOURCE_DATE_EPOCH), а не у часов машины, —
#      иначе один и тот же коммит давал бы разные архивы каждую минуту.
#   3. Пути обрезаны (-trimpath), запись о системе контроля версий выключена
#      (-buildvcs=false, метки и так проставляются компоновщиком), CGO выключен,
#      инструментарий берётся местный (GOTOOLCHAIN=local) — подкачка другого
#      компилятора посреди выпуска сломала бы повторимость молча.
#
# Каждая цель собирается ДВАЖДЫ в разные каталоги, и отпечатки сверяются. Это
# и есть доказательство: расхождение останавливает выпуск. Отключается
# доводом --skip-repro-check, когда нужна просто быстрая сборка.
#
# Ядро на flang собирается ВНУТРЬ двоичного файла признаком `flangcore`
# (host/decider_flang.go). Без него получилась бы заглушка, которая считает,
# но не решает, — такой файл не выпускается: сборка падает на самопроверке.

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

TARGETS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64")
DIST="$ROOT/dist"
REPRO_CHECK=1
VERSION_ARG=""

for arg in "$@"; do
	case "$arg" in
	--skip-repro-check) REPRO_CHECK=0 ;;
	-h | --help)
		sed -n '5,12p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
		exit 0
		;;
	-*)
		echo "build-release: неизвестный ключ $arg" >&2
		exit 2
		;;
	*) VERSION_ARG="$arg" ;;
	esac
done

# --- версия, коммит, время -------------------------------------------------

VERSION="${VERSION_ARG:-${DIGITDISK_VERSION:-}}"
if [ -z "$VERSION" ]; then
	[ -f VERSION ] || {
		echo "build-release: нет файла VERSION и версия не задана" >&2
		exit 1
	}
	VERSION="$(tr -d ' \t\n\r' <VERSION)"
fi
VERSION="${VERSION#v}"
[ -n "$VERSION" ] || {
	echo "build-release: версия пуста" >&2
	exit 1
}

git rev-parse --is-inside-work-tree >/dev/null 2>&1 || {
	echo "build-release: это не клон git — хеш сборки взять неоткуда" >&2
	exit 1
}

if [ -n "$(git status --porcelain)" ]; then
	if [ "${DIGITDISK_ALLOW_DIRTY:-0}" = "1" ]; then
		echo "build-release: ВНИМАНИЕ, дерево с правками, DIGITDISK_ALLOW_DIRTY=1" >&2
	else
		echo "build-release: дерево с правками — выпуск собирается только из чистого." >&2
		echo "               git status --porcelain должен быть пуст." >&2
		exit 1
	fi
fi

COMMIT="$(git rev-parse HEAD)"
COMMIT_DATE="$(git show -s --format=%cI HEAD)"
SOURCE_DATE_EPOCH="$(git show -s --format=%ct HEAD)"
export SOURCE_DATE_EPOCH

GO_VERSION="$(go env GOVERSION)"

echo "digitdisk $VERSION"
echo "  коммит          $COMMIT"
echo "  время коммита   $COMMIT_DATE (SOURCE_DATE_EPOCH=$SOURCE_DATE_EPOCH)"
echo "  инструментарий  $GO_VERSION"
echo "  цели            ${TARGETS[*]}"
echo "  сверка повтором $([ "$REPRO_CHECK" = 1 ] && echo включена || echo выключена)"
echo

# --- сборка ----------------------------------------------------------------

LDFLAGS="-s -w -buildid="
LDFLAGS="$LDFLAGS -X main.version=$VERSION"
LDFLAGS="$LDFLAGS -X main.commit=$COMMIT"
LDFLAGS="$LDFLAGS -X main.date=$COMMIT_DATE"

build_one() {
	# build_one <goos> <goarch> <путь к файлу>
	(
		cd "$ROOT/host"
		env -u GOFLAGS \
			GOTOOLCHAIN=local CGO_ENABLED=0 GOOS="$1" GOARCH="$2" \
			go build -tags flangcore -trimpath -buildvcs=false \
			-ldflags "$LDFLAGS" -o "$3" .
	)
}

rm -rf "$DIST"
mkdir -p "$DIST"

for target in "${TARGETS[@]}"; do
	goos="${target%%/*}"
	goarch="${target##*/}"
	name="digitdisk-$VERSION-$goos-$goarch"
	stage="$DIST/$name"
	mkdir -p "$stage"

	printf '%-28s' "сборка $goos/$goarch"
	build_one "$goos" "$goarch" "$stage/digitdisk"
	sum="$(sha256sum "$stage/digitdisk" | cut -d' ' -f1)"

	if [ "$REPRO_CHECK" = 1 ]; then
		second="$DIST/.repro/$name"
		mkdir -p "$second"
		build_one "$goos" "$goarch" "$second/digitdisk"
		sum2="$(sha256sum "$second/digitdisk" | cut -d' ' -f1)"
		if [ "$sum" != "$sum2" ]; then
			echo
			echo "build-release: сборка НЕ повторима: $sum != $sum2" >&2
			exit 1
		fi
		rm -rf "$second"
	fi
	echo "$sum  повторена"

	# В архив едет то, без чего двоичный файл нельзя ни законно раздать, ни
	# понять: лицензия, происхождение, обе редакции описания, версия.
	cp LICENSE NOTICE README.md README.ru.md VERSION "$stage/"

	touch -d "@$SOURCE_DATE_EPOCH" "$stage"/* "$stage"
	tar --sort=name --owner=0 --group=0 --numeric-owner \
		--mtime="@$SOURCE_DATE_EPOCH" \
		--pax-option=exthdr.name=%d/PaxHeaders/%f,delete=atime,delete=ctime \
		-C "$DIST" -cf - "$name" | gzip -9 -n >"$DIST/$name.tar.gz"
	rm -rf "$stage"
done
rmdir "$DIST/.repro" 2>/dev/null || true

# --- контрольные суммы -----------------------------------------------------

(cd "$DIST" && sha256sum ./*.tar.gz | sed 's| \./| |' >SHA256SUMS)
echo
echo "контрольные суммы:"
sed 's/^/  /' "$DIST/SHA256SUMS"

# --- формула Homebrew ------------------------------------------------------

SHA_AMD64="$(grep -- "-linux-amd64.tar.gz\$" "$DIST/SHA256SUMS" | cut -d' ' -f1)"
SHA_ARM64="$(grep -- "-linux-arm64.tar.gz\$" "$DIST/SHA256SUMS" | cut -d' ' -f1)"
SHA_MAC_AMD64="$(grep -- "-darwin-amd64.tar.gz\$" "$DIST/SHA256SUMS" | cut -d' ' -f1)"
SHA_MAC_ARM64="$(grep -- "-darwin-arm64.tar.gz\$" "$DIST/SHA256SUMS" | cut -d' ' -f1)"

mkdir -p "$DIST/homebrew"
sed -e "s/VERSION_PLACEHOLDER/$VERSION/g" \
	-e "s/SHA256_LINUX_AMD64_PLACEHOLDER/$SHA_AMD64/" \
	-e "s/SHA256_LINUX_ARM64_PLACEHOLDER/$SHA_ARM64/" \
	-e "s/SHA256_MACOS_AMD64_PLACEHOLDER/$SHA_MAC_AMD64/" \
	-e "s/SHA256_MACOS_ARM64_PLACEHOLDER/$SHA_MAC_ARM64/" \
	packaging/homebrew/digitdisk.rb >"$DIST/homebrew/digitdisk.rb"

if grep -q PLACEHOLDER "$DIST/homebrew/digitdisk.rb"; then
	echo "build-release: в формуле остались незаполненные места" >&2
	exit 1
fi
echo
echo "формула:  dist/homebrew/digitdisk.rb  → Formula/digitdisk.rb в digitable-lol/homebrew-tap"

# --- самопроверка ----------------------------------------------------------
#
# Проверяется не «файл собрался», а «файл делает своё дело»: называет свою
# версию, называет свой хеш и решает настоящим ядром. Двоичный файл с
# заглушкой вместо ядра молча давал бы пустой разбор по разрядам, и заметил бы
# это только поставивший.

host_os="$(go env GOHOSTOS)"
host_arch="$(go env GOHOSTARCH)"
native="digitdisk-$VERSION-$host_os-$host_arch.tar.gz"
if [ ! -f "$DIST/$native" ]; then
	echo
	echo "самопроверка пропущена: под $host_os/$host_arch выпуск не собирается"
	exit 0
fi

probe="$(mktemp -d)"
trap 'rm -rf "$probe"' EXIT
tar -C "$probe" -xzf "$DIST/$native"
bin="$probe/digitdisk-$VERSION-$host_os-$host_arch/digitdisk"

echo
echo "самопроверка ($host_os/$host_arch):"
out="$("$bin" --version)"
echo "$out" | sed 's/^/  /'

echo "$out" | grep -q "^digitdisk $VERSION\$" || {
	echo "build-release: --version не назвал версию $VERSION" >&2
	exit 1
}
echo "$out" | grep -q "${COMMIT:0:7}" || {
	echo "build-release: --version не назвал хеш сборки ${COMMIT:0:7}" >&2
	exit 1
}
echo "$out" | grep -q "flang" || {
	echo "build-release: собрано без ядра на flang — признак flangcore не сработал" >&2
	exit 1
}
"$bin" analyze "$ROOT/tools" --json >/dev/null || {
	echo "build-release: analyze не отработал" >&2
	exit 1
}
"$bin" status --json --sample 10 >/dev/null || {
	echo "build-release: status не отработал" >&2
	exit 1
}
echo "  analyze и status отработали, ядро на flang на месте"
echo
echo "готово: $DIST"
