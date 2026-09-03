// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package clean

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digitdisk/internal/core"
	"digitdisk/internal/protect"
)

// naming is a stand-in decision layer that answers BOTH questions: the
// приговор (like judge, bluntly) and the природа.  It exists because these
// tests are about what the host does with the second answer, not about how the
// second answer is reached — that is proved in core/disk-inventory.flang,
// where «Природа находки» carries its examples and its постусловие.
type naming struct {
	judge
	nature map[string]core.Nature
}

func (n naming) Nature(r core.Record, d core.Decision) core.Nature {
	if got, ok := n.nature[filepath.Base(r.Path)]; ok {
		return got
	}
	if d.Verdict == core.VerdictRemovable {
		return core.NatureTrash
	}
	return core.NaturePersonal
}

func (naming) Strictness(nat core.Nature) int {
	switch nat {
	case core.NatureTrash:
		return 1
	case core.NatureStore, core.NatureVCS:
		return 3
	}
	return 2
}

// byHand builds the plan забой builds, on a ground inside the tree.  The
// твёрдые запреты are pointed at a home and a tool directory that are not this
// machine's, so that a test tree under t.TempDir() is ordinary ground; the
// запреты themselves are under test in TestHardStopsRefuseAndSayHowToGetPast.
func byHand(t *testing.T, root string, only []string, guard *protect.List, nature map[string]core.Nature) (Plan, error) {
	t.Helper()
	return Make(Options{
		Root: root,
		Decider: naming{
			judge:  judge{removable: map[string]bool{"old.bin": true, "older.bin": true}, classOf: map[string]core.Class{"old.bin": core.ClassCache, "older.bin": core.ClassCache}},
			nature: nature,
		},
		Now: time.Now(), Version: "испытание", Only: only, Protect: guard, ByHand: true,
		Stop: StopOptions{Home: t.TempDir(), Tool: t.TempDir()},
	})
}

// СЛУЧАЙ ВЛАДЕЛЬЦА, на слое хозяина. Каталог, которому ядро не давало
// «МожноУбрать», попадает в план забоя ЦЕЛИКОМ — и файлом, и собой, — а слово
// ядра о том, что это, приезжает вместе с ним.
func TestByHandTakesWhatTheVerdictKeeps(t *testing.T) {
	root := tree(t)
	p, err := byHand(t, root, []string{filepath.Join(root, "docs")}, nil,
		map[string]core.Nature{"letter.txt": core.NatureSource})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Items) != 1 || filepath.Base(p.Items[0].Path) != "letter.txt" {
		t.Fatalf("в план забоя попало %+v, ожидалось letter.txt", p.Items)
	}
	if p.Items[0].Verdict == core.VerdictRemovable {
		t.Fatal("испытание не о том: слой всё-таки пометил письмо «МожноУбрать»")
	}
	if p.Items[0].Nature != core.NatureSource {
		t.Errorf("природа не доехала до плана: %q", p.Items[0].Nature)
	}
	if len(p.Dirs) != 1 || filepath.Base(p.Dirs[0].Path) != "docs" {
		t.Fatalf("каталог не попал в план: %+v", p.Dirs)
	}
	if p.Paths() != 2 {
		t.Errorf("путей в плане %d, ожидалось 2", p.Paths())
	}
	if p.Strictness != 2 {
		t.Errorf("строгость %d, а исходники — вторая ступень", p.Strictness)
	}

	// И «c» на том же месте по-прежнему не берёт ничего: уборка не тронута.
	clean, err := Make(Options{Root: root, Decider: judge{}, Now: time.Now(),
		Only: []string{filepath.Join(root, "docs")}})
	if err != nil {
		t.Fatal(err)
	}
	if len(clean.Items) != 0 || len(clean.Dirs) != 0 {
		t.Errorf("уборка взяла то, чего ядро не помечало: %+v %+v", clean.Items, clean.Dirs)
	}
}

// Строгость плана — САМАЯ ВЫСОКАЯ из названных, а не средняя и не последняя:
// один объект хранилища среди сотни кэшей платит за весь вопрос.
func TestStrictnessIsTheHighestOverThePlan(t *testing.T) {
	root := tree(t)
	p, err := byHand(t, root, []string{root}, nil,
		map[string]core.Nature{"letter.txt": core.NatureStore})
	if err != nil {
		t.Fatal(err)
	}
	if p.Strictness != 3 {
		t.Fatalf("строгость %d, а в плане хранилище — должна быть 3", p.Strictness)
	}
	// И слой, который второго вопроса не отвечает, тоже покупает третью.
	mute, err := Make(Options{Root: root, Decider: judge{}, Now: time.Now(),
		Only: []string{root}, ByHand: true,
		Stop: StopOptions{Home: t.TempDir(), Tool: t.TempDir()}})
	if err != nil {
		t.Fatal(err)
	}
	if mute.Strictness != core.Strictest {
		t.Errorf("молчащий слой дал строгость %d, а «не знаю» — это %d", mute.Strictness, core.Strictest)
	}
}

// Забой сносит и опустевшие каталоги: «удалить папку» значит и папку. Каталог,
// в котором что-то осталось, при этом стоит и говорит почему.
func TestByHandRemovesTheDirectoriesItEmptied(t *testing.T) {
	root := tree(t)
	// В cache/deep кладём ссылку — её забой не стирает, и каталог обязан
	// остаться, потому что снимается он вызовом, отказывающим на непустом.
	if err := os.Symlink(filepath.Join(root, "docs", "letter.txt"), filepath.Join(root, "cache", "deep", "ссылка")); err != nil {
		t.Skipf("символические ссылки недоступны: %v", err)
	}
	p, err := byHand(t, root, []string{filepath.Join(root, "cache")}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Каталоги идут от глубоких к мелким: иначе родитель встретил бы
	// непустого ребёнка и остался бы стоять просто по порядку.
	if len(p.Dirs) != 2 || filepath.Base(p.Dirs[0].Path) != "deep" {
		t.Fatalf("каталоги не отсортированы вглубь: %+v", p.Dirs)
	}
	j, err := Erase(p, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(root, "cache", "deep")); err != nil {
		t.Errorf("каталог со ссылкой внутри всё-таки снят: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, "cache")); err != nil {
		t.Errorf("родитель непустого каталога всё-таки снят: %v", err)
	}
	if n := j.PurgedDirs(); n != 0 {
		t.Errorf("снято каталогов %d, а снять было нельзя ни одного", n)
	}
	if kept := j.KeptDirs(); len(kept) != 2 {
		t.Errorf("журнал не записал, почему каталоги остались: %+v", kept)
	}

	// А там, где не осталось ничего, каталог уходит.
	p2, err := byHand(t, root, []string{filepath.Join(root, "docs")}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	j2, err := Erase(p2, Options{Now: time.Now(), Version: "испытание"})
	if err != nil {
		t.Fatal(err)
	}
	if j2.PurgedDirs() != 1 {
		t.Errorf("опустевший каталог не снят: %+v", j2.Dirs)
	}
	if _, err := os.Lstat(filepath.Join(root, "docs")); !os.IsNotExist(err) {
		t.Errorf("docs пережил забой: %v", err)
	}
	if !j2.ByHand {
		t.Error("журнал не записал, что план строился по воле человека")
	}
}

// Забой не берёт того, что и не смог бы стереть: ссылку, недоступное, чужой
// корень, свою корзину. Приговор из этого списка ушёл — остальное осталось.
func TestByHandStillRefusesWhatItCannotErase(t *testing.T) {
	root := tree(t)
	if err := os.Symlink(filepath.Join(root, "docs"), filepath.Join(root, "cache", "ссылка")); err != nil {
		t.Skipf("символические ссылки недоступны: %v", err)
	}
	p, err := byHand(t, root, []string{filepath.Join(root, "cache")}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range p.Items {
		if strings.Contains(it.Path, "ссылка") {
			t.Fatalf("ссылка попала в план забоя: %+v", it)
		}
	}
	var said bool
	for _, r := range p.Refused {
		if strings.Contains(r.Path, "ссылка") && strings.Contains(r.Reason.String(), "символическая ссылка") {
			said = true
		}
	}
	if !said {
		t.Errorf("отказ по ссылке не назван вслух: %+v", p.Refused)
	}
}

// Твёрдый запрет: пять мест поимённо, каждое со СВОЕЙ причиной и со способом
// обойти. Проверяется и то, что запрет ловит попытку зайти СВЕРХУ — указать на
// родителя запрещённого места, — и то, что внутри запрещённого каталога
// стирать по-прежнему можно.
func TestHardStopsRefuseAndSayHowToGetPast(t *testing.T) {
	home := filepath.Join(t.TempDir(), "дом")
	tool := filepath.Join(t.TempDir(), "инструмент")
	opt := StopOptions{Home: home, Tool: tool, Getenv: func(string) string { return "" }}

	for _, c := range []struct {
		path string
		why  string
	}{
		{"/", "корень файловой системы"},
		{"/usr", "системный каталог"},
		{"/etc", "системный каталог"},
		{"/var", "системный каталог"},
		{home, "домашний каталог целиком"},
		{tool, "сам digitdisk"},
		{filepath.Join(home, ".digitable", "digitdisk"), "настройки digitdisk"},
		{filepath.Join(home, ".config", "digitdisk"), "прежний дом настроек"},
		{filepath.Join("/данные", TrashName), "корзина digitdisk"},
	} {
		stop := HardStop(c.path, opt)
		if stop.Empty() {
			t.Errorf("%s пропущен твёрдым запретом", c.path)
			continue
		}
		if !strings.Contains(stop.Why.String(), c.why) {
			t.Errorf("%s: причина %q не про %q", c.path, stop.Why.String(), c.why)
		}
		if stop.Around.Empty() {
			t.Errorf("%s: отказ не сказал, как быть, если человек всё же прав", c.path)
		}
		msg := stop.Err().Error()
		for _, want := range []string{"СТИРАТЬ ОТСЮДА НЕЛЬЗЯ", c.path, "Как быть"} {
			if !strings.Contains(msg, want) {
				t.Errorf("%s: в отказе нет %q:\n%s", c.path, want, msg)
			}
		}
	}

	// Сверху — тоже запрет: указать на родителя дома не значит обойти его.
	if HardStop(filepath.Dir(home), opt).Empty() {
		t.Errorf("родитель дома %s пропущен: запрет обходится сверху", filepath.Dir(home))
	}
	// А внутри — не запрет: ломает машину вынос каталога целиком, а не файла
	// в нём, и об отдельном файле говорит природа, а не запрет.
	for _, ok := range []string{"/var/tmp/моя-сборка", filepath.Join(home, "проекты"), "/usr/local/src/моё"} {
		if stop := HardStop(ok, opt); !stop.Empty() {
			t.Errorf("%s запрещён, а запрещаться должен только каталог целиком: %s", ok, stop.Why)
		}
	}
	// И «varnish» не «var»: сверка идёт по составляющим.
	if stop := HardStop("/varnish", opt); !stop.Empty() {
		t.Errorf("/varnish принят за /var: %s", stop.Why)
	}
}

// Твёрдый запрет держит и того, кто зовёт clean.Make напрямую, а не через
// экран: правило живёт в одном месте, и экран — не единственная дорога к нему.
func TestMakeByHandRefusesAHardStoppedGround(t *testing.T) {
	root := tree(t)
	home := filepath.Join(root, "docs")
	_, err := Make(Options{Root: root, Decider: judge{}, Now: time.Now(),
		Only: []string{home}, ByHand: true, Stop: StopOptions{Home: home, Tool: t.TempDir()}})
	if err == nil {
		t.Fatal("clean.Make принял запрещённую землю")
	}
	if !strings.Contains(err.Error(), "СТИРАТЬ ОТСЮДА НЕЛЬЗЯ") {
		t.Errorf("отказ не назвал себя запретом: %v", err)
	}
	// И указано ничего — тоже отказ: забой стирает то, на что указали.
	if _, err := Make(Options{Root: root, Decider: judge{}, Now: time.Now(), ByHand: true,
		Stop: StopOptions{Home: t.TempDir(), Tool: t.TempDir()}}); err == nil {
		t.Fatal("забой без земли принят")
	}
}

// Защитный список на самой земле отказывает целиком и называет правило.
func TestProtectOnTheGroundRefusesTheWholeErasure(t *testing.T) {
	root := tree(t)
	guard, err := protect.Load(protect.Options{Args: []string{filepath.Join(root, "docs")}, Home: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	before := snapshot(t, root)
	_, err = byHand(t, root, []string{filepath.Join(root, "docs")}, guard, nil)
	if err == nil {
		t.Fatal("защищённая земля принята к стиранию")
	}
	for _, want := range []string{"СТИРАТЬ ОТСЮДА НЕЛЬЗЯ", "защитный список", "--protect"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("отказ не сказал %q: %v", want, err)
		}
	}
	if after := snapshot(t, root); after != before {
		t.Errorf("отказ тронул дерево:\n%s\n---\n%s", before, after)
	}
}

// Корень обхода забой не сносит никогда: на нём открыт os.Root, и снявший его
// стирает пол под собой.
func TestByHandNeverTakesTheWalkRootItself(t *testing.T) {
	root := tree(t)
	p, err := byHand(t, root, []string{root}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range p.Dirs {
		if d.Path == root {
			t.Fatalf("корень обхода попал в план: %+v", d)
		}
	}
	if _, err := Erase(p, Options{Now: time.Now(), Version: "испытание"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(root); err != nil {
		t.Fatalf("корень обхода снесён: %v", err)
	}
}
