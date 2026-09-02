// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package places

import (
	"strings"
	"testing"

	"digitdisk/internal/core"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

// TestBuiltinLoads is the test that keeps the shipped file honest: every row
// parses, every row has a source, and every chain is bounded by slashes — the
// property the decision layer refuses a справочник for missing.
func TestBuiltinLoads(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		d, err := Load(Options{Home: "/home/u", GOOS: goos, Getenv: env(nil), Config: t.TempDir()})
		if err != nil {
			t.Fatalf("%s: встроенный справочник не читается: %v", goos, err)
		}
		if len(d.Entries) < 50 {
			t.Errorf("%s: мест %d — справочник, который ничего не знает, не справочник", goos, len(d.Entries))
		}
		if len(d.Applicable()) == 0 {
			t.Errorf("%s: ни одно место не применимо", goos)
		}
		for _, e := range d.Entries {
			if !strings.HasPrefix(e.Source, "https://") {
				t.Errorf("место %q: источник %q — не ссылка на документацию", e.Name, e.Source)
			}
			if e.Chain == "" || !strings.HasPrefix(e.Chain, "/") || !strings.HasSuffix(e.Chain, "/") {
				t.Errorf("место %q: цепь %q не ограничена косыми", e.Name, e.Chain)
			}
			if len(e.Chain) < 2 {
				t.Errorf("место %q: цепь пуста", e.Name)
			}
			switch e.Class {
			case core.ClassCache, core.ClassLog, core.ClassBuild, core.ClassDownload:
			default:
				t.Errorf("место %q: разряд %q справочнику запрещён", e.Name, e.Class)
			}
		}
	}
}

// Здесь стоял тест «ни одна строка справочника не ссылается на чужой
// чистильщик». Он списан, и по хорошей причине: чтобы проверить запрет, тест
// нёс запрещённые имена литералами — и лицензионный сторож уронил на нём
// прогон, `flang io tools/licensing.flang`, назвав файл и строку. Проверка уже
// есть, она шире (весь тронутый диапазон, а не один файл) и работает; вторая её
// копия только заводила бы в дерево ровно то, что запрещает.
//
// TestSourcesAreDocumentation держит ту часть обещания, которую сторож не
// проверяет: источником обязана быть страница документации, а не что попало.
func TestSourcesAreDocumentation(t *testing.T) {
	d, err := Load(Options{Home: "/home/u", GOOS: "linux", Getenv: env(nil), Config: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]int{}
	for _, e := range d.Entries {
		if !strings.HasPrefix(e.Source, "https://") || len(e.Source) < len("https://x/") {
			t.Errorf("строка %d: источник %q — не страница документации", e.Line, e.Source)
		}
		seen[e.Resolved+"|"+e.Chain]++
	}
	for key, n := range seen {
		if n > 1 && key != "|" {
			t.Errorf("место %q названо %d раз — первое подошедшее победит, второе никогда", key, n)
		}
	}
}

func TestAnchors(t *testing.T) {
	opt := Options{Home: "/home/u", GOOS: "linux", Getenv: env(map[string]string{
		"XDG_CACHE_HOME": "/var/tmp/cache",
	})}
	body := strings.Join([]string{
		"кэш|дом|все|.npm//_cacache|npm_config_cache|npm|https://x/",
		"кэш|кэш|все|go-build|GOCACHE|go|https://x/",
		"кэш|корень|linux|/var/cache/apt/archives||apt|https://x/",
		"сборка|всюду|все|.turbo||turbo|https://x/",
		"кэш|данные|все|NuGet/http-cache||nuget|https://x/",
	}, "\n")
	d, err := parse(body, "проба", opt)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"/home/u/.npm/_cacache/",
		"/var/tmp/cache/go-build/",
		"/var/cache/apt/archives/",
		"/.turbo/",
		"/home/u/.local/share/NuGet/http-cache/",
	}
	for i, w := range want {
		if d.Entries[i].Chain != w {
			t.Errorf("место %d: цепь %q, ждали %q", i, d.Entries[i].Chain, w)
		}
	}
}

// TestEnvRelocates is why the переменная column exists: a developer who moved
// a cache with the tool's own environment variable must get their real place,
// not the default one that is no longer there.
func TestEnvRelocates(t *testing.T) {
	opt := Options{Home: "/home/u", GOOS: "linux", Getenv: env(map[string]string{
		"npm_config_cache": "/data/npm",
		"GOCACHE":          "/data/gocache",
	})}
	body := strings.Join([]string{
		"кэш|дом|все|.npm//_cacache|npm_config_cache|npm|https://x/",
		"кэш|кэш|все|go-build|GOCACHE|go|https://x/",
	}, "\n")
	d, err := parse(body, "проба", opt)
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Entries[0].Chain; got != "/data/npm/_cacache/" {
		t.Errorf("основание не перенеслось: %q", got)
	}
	if got := d.Entries[1].Chain; got != "/data/gocache/" {
		t.Errorf("весь путь не перенёсся: %q", got)
	}
	for _, e := range d.Entries {
		if !e.Relocated {
			t.Errorf("место %q не отмечено как перенесённое", e.Name)
		}
	}
}

// TestMatchIsByComponents is the September trap, restated for the справочник:
// a place must never be recognised inside a longer component.
func TestMatchIsByComponents(t *testing.T) {
	d, err := parse("кэш|дом|все|.npm//_cacache|npm_config_cache|npm|https://x/\nсборка|всюду|все|node_modules/.vite||vite|https://x/",
		"проба", Options{Home: "/home/u", GOOS: "linux", Getenv: env(nil)})
	if err != nil {
		t.Fatal(err)
	}
	yes := []string{
		"/home/u/.npm/_cacache",
		"/home/u/.npm/_cacache/index-v5/aa",
		"/home/u/p/node_modules/.vite/deps/x.js",
	}
	no := []string{
		"/home/u/x.npm/_cacache/aa",
		"/home/u/.npm/_cacacheX/aa",
		"/srv/backup/home/u/.npm/_cacache/aa",
		"/home/u/p/my_node_modules/.vite/x",
		"/home/u/.npm/_logs/x.log",
	}
	for _, p := range yes {
		if _, ok := d.Match(p); !ok {
			t.Errorf("%q должен попасть в известное место", p)
		}
	}
	for _, p := range no {
		if e, ok := d.Match(p); ok {
			t.Errorf("%q не место, а совпало с %q", p, e.Name)
		}
	}
}

func TestBadRowsAreRefused(t *testing.T) {
	opt := Options{Home: "/home/u", GOOS: "linux", Getenv: env(nil)}
	for _, bad := range []string{
		"крупное|дом|все|x||имя|https://x/",
		"неизвестное|дом|все|x||имя|https://x/",
		"кэш|нетакой|все|x||имя|https://x/",
		"кэш|дом|плутон|x||имя|https://x/",
		"кэш|дом|все|x|||https://x/",
		"кэш|дом|все|x||имя|",
		"кэш|дом|все|||имя|https://x/",
		"кэш|дом|все|x||имя",
		"кэш|всюду|все|x|VAR|имя|https://x/",
		"кэш|дом|все|/abs||имя|https://x/",
	} {
		if _, err := parse(bad, "проба", opt); err == nil {
			t.Errorf("строка %q принята, а должна быть отвергнута", bad)
		}
	}
}

// TestOffLoadsNothing states what --no-places means: the layer goes back to
// judging by приметы alone, which is what it did before the справочник.
func TestOffLoadsNothing(t *testing.T) {
	d, err := Load(Options{Off: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Entries) != 0 || len(d.Places()) != 0 {
		t.Error("--no-places обязан не давать ни одного места")
	}
}
