// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Здесь сверяется не поведение, а полнота: КАЖДАЯ строка, которую видит
// человек, обязана иметь пару на втором языке.
//
// Полупереведённый вывод хуже непереведённого: непереведённый читается целиком
// на одном языке, а наполовину переведённый заставляет читателя гадать, что
// именно ему не показали. Поэтому лекарство здесь то же, которым уже лечится
// расхождение справки, кода и страницы руководства, — прогон, читающий сам
// исходник:
//
//	Пары        каждый литерал в T/F/Say/Errorf есть в словаре;
//	Утечки      в пакетах вывода нет кириллического литерала мимо словаря;
//	Мёртвое     в словаре нет статьи, которую никто не спрашивает;
//	Заполнители  %-глаголы русской и английской половин совпадают;
//	Договор     у каждого разряда, приговора, вида и якоря есть слово.
//
// Проверить, что прогон ловит: убрать одну статью из dict_*.go (падают Пары),
// снять обёртку l.T с одной строки в internal/report (падают Утечки).

// hostDir is the root of the Go module: this package sits two levels down.
const hostDir = "../.."

// outputPackages are the packages whose whole job is to put text in front of a
// person.  A Cyrillic literal there is a line somebody will read, so it has to
// go through the dictionary; everywhere else the text travels as a Phrase or
// as a вокабула and is caught by TestПарыЕстьУКаждойСтроки instead.
var outputPackages = []string{
	".", // main.go, version.go, decider_*.go
	"internal/report",
	"internal/ui",
	"internal/cli",
	"internal/run",
}

// wrappers are the calls that put a Russian wording into the dictionary's
// hands.  The first argument of each is a wording and must have a pair.
var wrappers = map[string]bool{"T": true, "F": true, "Say": true, "Errorf": true}

// isWrapper reports whether this call hands its first argument to the
// dictionary.  fmt.Errorf spells its name the same way lang.Errorf does and
// means the opposite — a message nobody translated — so the receiver is
// checked, and a fmt.Errorf holding Russian is caught by TestОтказыПереводятся
// instead.
func isWrapper(sel *ast.SelectorExpr) bool {
	if !wrappers[sel.Sel.Name] {
		return false
	}
	if id, ok := sel.X.(*ast.Ident); ok && (id.Name == "fmt" || id.Name == "errors") {
		return false
	}
	return true
}

// collectWordings reads every non-test .go file under the host and returns the
// wordings handed to the wrappers, each with the files it was found in.
func collectWordings(t *testing.T) map[string][]string {
	t.Helper()
	found := map[string][]string{}
	forEachGoFile(t, hostDir, func(path string, file *ast.File, fset *token.FileSet) {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !isWrapper(sel) {
				return true
			}
			lit, ok := stringLit(call.Args[0])
			if !ok {
				return true
			}
			found[lit] = append(found[lit], rel(path))
			return true
		})
	})
	return found
}

// collectLiterals returns every Cyrillic string literal of the host, wrapped or
// not.  A wording may reach the dictionary two ways: through a T/F call, or as
// DATA that something else translates later — the gloss of a подкоманда in
// cli.Commands is looked up when the справка is built, not where it is
// written.  Both are lines a person reads, and both must have a pair.
func collectLiterals(t *testing.T) map[string][]string {
	t.Helper()
	found := map[string][]string{}
	forEachGoFile(t, hostDir, func(path string, file *ast.File, fset *token.FileSet) {
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING || isStructTag(lit) {
				return true
			}
			// Не только кириллические: статья словаря бывает и
			// вовсе без русских букв — «--protect %q: %s» одинаково
			// выглядит на обоих языках, но пару имеет и мёртвой не
			// является.
			text, _ := stringLit(lit)
			found[text] = append(found[text], rel(path))
			return true
		})
	})
	return found
}

func TestПарыЕстьУКаждойСтроки(t *testing.T) {
	found := collectWordings(t)
	var lost []string
	for ru := range found {
		if !Known(ru) {
			lost = append(lost, ru)
		}
	}
	sort.Strings(lost)
	for _, ru := range lost {
		t.Errorf("строка без пары на втором языке: %q (%s)", ru, strings.Join(found[ru], ", "))
	}
	if len(found) < 200 {
		t.Fatalf("в исходнике нашлось %d строк — разбор сломался", len(found))
	}
	t.Logf("строк в исходнике %d, статей в словаре %d, вокабул %d", len(found), Size(), len(vocab))
}

func TestНиОднаСтрокаНеПечатаетсяМимоСловаря(t *testing.T) {
	leaks := 0
	for _, pkg := range outputPackages {
		dir := filepath.Join(hostDir, pkg)
		forEachGoFile(t, dir, func(path string, file *ast.File, fset *token.FileSet) {
			if filepath.Dir(path) != filepath.Clean(dir) {
				return // only this package, not the ones under it
			}
			wrapped := map[ast.Node]bool{}
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok || len(call.Args) == 0 {
					return true
				}
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok && isWrapper(sel) {
					wrapped[call.Args[0]] = true
				}
				return true
			})
			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING || wrapped[n] {
					return true
				}
				text, _ := stringLit(lit)
				if !hasCyrillic(text) || isStructTag(lit) {
					return true
				}
				if dataFiles[rel(path)] && Known(text) {
					// In one file the wordings ARE data: the
					// glosses of cli.Commands and the lines of the
					// key list are written there and looked up when
					// the справка is built, not where they stand.
					// A pair is all that can be asked of them, and
					// it is asked.  Everywhere else a Cyrillic
					// literal must go through T/F at the point it
					// is printed: having a pair is not the same as
					// using it, and an unwrapped literal prints
					// Russian to an English reader while every
					// count in this file stays green.
					return true
				}
				leaks++
				t.Errorf("%s:%d: кириллица мимо словаря — печатается как есть, на любом языке: %q",
					rel(path), fset.Position(lit.Pos()).Line, text)
				return true
			})
		})
	}
	if leaks == 0 {
		t.Logf("в пакетах вывода (%s) кириллицы мимо словаря нет", strings.Join(outputPackages, ", "))
	}
}

func TestВСловареНетМёртвыхСтатей(t *testing.T) {
	found := collectLiterals(t)
	var dead []string
	for _, ru := range Wordings() {
		if IsVocab(ru) {
			continue
		}
		if _, used := found[ru]; !used {
			dead = append(dead, ru)
		}
	}
	sort.Strings(dead)
	for _, ru := range dead {
		t.Errorf("статья словаря, которую никто не спрашивает: %q", ru)
	}
}

func TestЗаполнителиСовпадаютВОбеихПоловинах(t *testing.T) {
	checked := 0
	for _, ru := range Wordings() {
		en := EN.T(ru)
		a, b := Verbs(ru), Verbs(en)
		if len(a) != len(b) {
			t.Errorf("заполнителей не поровну: %q → %v, %q → %v", ru, a, en, b)
			continue
		}
		for i := range a {
			if a[i] != b[i] {
				t.Errorf("заполнитель %d разный: %q → %v, %q → %v", i+1, ru, a, en, b)
				break
			}
		}
		checked++
	}
	t.Logf("сверено статей: %d", checked)
}

// forEachGoFile walks a directory tree and parses every .go file that is not a
// test.  Tests are left out on purpose: what a test prints nobody reads in
// English.
func forEachGoFile(t *testing.T, root string, do func(string, *ast.File, *token.FileSet)) {
	t.Helper()
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// The dictionary itself is all Russian and all English by
		// definition; scanning it would only find itself.
		if strings.Contains(filepath.ToSlash(path), "internal/lang/dict_") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("%s не разбирается: %v", path, err)
		}
		do(path, file, fset)
		return nil
	})
	if err != nil {
		t.Fatalf("обход %s: %v", root, err)
	}
}

func stringLit(e ast.Expr) (string, bool) {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", false
	}
	s := lit.Value
	if strings.HasPrefix(s, "`") {
		return strings.Trim(s, "`"), true
	}
	return unquote(s), true
}

// unquote undoes the escapes a Go string literal may carry.  strconv.Unquote
// would do it, and does not, because a literal with a stray escape is a
// literal this check must still see rather than drop.
func unquote(s string) string {
	s = strings.TrimPrefix(strings.TrimSuffix(s, `"`), `"`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\\`, `\`)
	return s
}

func hasCyrillic(s string) bool {
	for _, r := range s {
		if r >= 0x0400 && r <= 0x04FF {
			return true
		}
	}
	return false
}

// isStructTag reports whether this literal is a field tag — `json:"путь"` —
// which names a field for a machine and is never printed.
func isStructTag(lit *ast.BasicLit) bool {
	return strings.HasPrefix(lit.Value, "`") && strings.Contains(lit.Value, `json:"`)
}

func rel(path string) string {
	p, err := filepath.Rel(hostDir, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(p)
}

// TestОтказыПереводятся ловит сообщение об отказе, собранное мимо словаря.
//
// `fmt.Errorf("нужен ровно один путь…")` печатается человеку так же, как любая
// строка отчёта, но выглядит в исходнике как внутренняя ошибка и потому легко
// остаётся непереведённым. Русский текст в fmt.Errorf и errors.New — дефект;
// то же сообщение через lang.Errorf едет как Phrase, печатается на языке
// читателя, а в журнал и в проверки errors.Is попадает по-русски, как прежде.
func TestОтказыПереводятся(t *testing.T) {
	forEachGoFile(t, hostDir, func(path string, file *ast.File, fset *token.FileSet) {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || (id.Name != "fmt" && id.Name != "errors") {
				return true
			}
			if sel.Sel.Name != "Errorf" && sel.Sel.Name != "New" {
				return true
			}
			text, _ := stringLit(call.Args[0])
			if !hasCyrillic(text) {
				return true
			}
			t.Errorf("%s:%d: отказ мимо словаря: %s.%s(%q) — нужен lang.Errorf",
				rel(path), fset.Position(call.Pos()).Line, id.Name, sel.Sel.Name, text)
			return true
		})
	})
}

// TestДоговорПереведёнСловами доказывает, где проходит граница между именем в
// ядре и словом на экране.
//
// Разряд «Кэш» и приговор «МожноУбрать» — идентификаторы решающего слоя на
// flang: они доказаны там, ездят в JSON и в журнале как есть, и ни одна буква
// из них здесь не меняется. Меняется ТОЛЬКО слово, которым их называет экран,
// и живёт это слово на стороне хозяина — вот в этой вокабуле. Прогон требует,
// чтобы у каждого имени договора слово было: разряд без английского слова —
// это полупереведённый отчёт.
func TestДоговорПереведёнСловами(t *testing.T) {
	names := contractNames(t)
	for _, name := range names {
		if !Known(name) {
			t.Errorf("имя договора без слова на втором языке: %q", name)
			continue
		}
		if !IsVocab(name) {
			t.Errorf("имя договора %q заведено как строка исходника, а оно значение", name)
		}
	}
	t.Logf("имён договора переведено словами: %d", len(names))
}

// contractNames reads the имена договора out of the source that declares them
// rather than repeating them here.  A repeated list is a list that goes stale:
// a new разряд would be added in core/contract.go, this file would go on
// naming five, and the check would pass while the screen printed a Russian
// word to an English reader.
//
// This package cannot import those packages — internal/sysinfo imports this
// one — so the constants are read as text, which is also what makes the check
// honest: it sees exactly what the source says.
func contractNames(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, file := range []string{
		"internal/core/contract.go",    // разряд, приговор, вид, якорь
		"internal/protect/protect.go",  // вид правила защитного списка
		"internal/sysinfo/sysinfo.go",  // имена неизмеренного
		"internal/sysinfo/hardware.go", // имена неизмеренного про железо
	} {
		path := filepath.Join(hostDir, file)
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("%s не разбирается: %v", file, err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				return true
			}
			for _, spec := range decl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, v := range vs.Values {
					if s, ok := stringLit(v); ok && hasCyrillic(s) {
						out = append(out, s)
					}
				}
			}
			return true
		})
	}
	sort.Strings(out)
	if len(out) < 15 {
		t.Fatalf("имён договора нашлось %d — разбор сломался", len(out))
	}
	return out
}

// dataFiles are the files where a Russian wording is DATA rather than a line
// being printed: internal/cli holds the one list of подкоманды and ключи, and
// the справка looks each wording up when it builds itself.  There the check
// asks for a pair and no more; everywhere else it asks for the wrapper too.
var dataFiles = map[string]bool{
	"internal/cli/cli.go": true,
}
