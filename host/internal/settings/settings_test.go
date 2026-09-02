// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"digitdisk/internal/lang"
)

// Здесь проверяется не язык, а ПОРЯДОК: кто решает, когда спрашивают и когда
// молчат, и что при этом появляется в домашнем каталоге. Инструмент, который
// молча заводит файлы у человека, — ровно то, за что ругают чистилки, поэтому
// «ничего не завелось» здесь такое же утверждение, как «завелось».

// env строит окружение из пар.
func env(pairs ...string) func(string) string {
	m := map[string]string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return func(name string) string { return m[name] }
}

// дом даёт пустой домашний каталог на один прогон.
func дом(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func файлНастроек(home string) string {
	return filepath.Join(home, FamilyDir, ToolDir, FileName)
}

func естьФайл(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestКлючПеревешиваетВсё(t *testing.T) {
	home := дом(t)
	o := Options{Home: home, Getenv: env("LANG", "ru_RU.UTF-8", "DIGITDISK_LANG", "ru")}
	if _, err := Save(o, Settings{Lang: lang.RU}); err != nil {
		t.Fatalf("настройки не записались: %v", err)
	}
	c := Decide(o, "en", Ask{})
	if c.Lang != lang.EN || c.Source != FromFlag {
		t.Fatalf("ключ не перевесил: %v из %v", c.Lang, c.Source)
	}
	if c.Saved != "" {
		t.Errorf("ключ на один запуск не должен ничего сохранять, а записал %s", c.Saved)
	}
}

func TestПеременнаяПеревешиваетФайл(t *testing.T) {
	home := дом(t)
	o := Options{Home: home, Getenv: env("DIGITDISK_LANG", "en", "LANG", "ru_RU.UTF-8")}
	if _, err := Save(o, Settings{Lang: lang.RU}); err != nil {
		t.Fatalf("настройки не записались: %v", err)
	}
	c := Decide(o, "", Ask{})
	if c.Lang != lang.EN || c.Source != FromEnv {
		t.Fatalf("переменная не перевесила: %v из %v", c.Lang, c.Source)
	}
}

func TestЗапомненныйВыборПеревешиваетЛокаль(t *testing.T) {
	home := дом(t)
	o := Options{Home: home, Getenv: env("LANG", "en_US.UTF-8")}
	if _, err := Save(o, Settings{Lang: lang.RU}); err != nil {
		t.Fatalf("настройки не записались: %v", err)
	}
	c := Decide(o, "", Ask{})
	if c.Lang != lang.RU || c.Source != FromFile {
		t.Fatalf("запомненный выбор не перевесил локаль: %v из %v", c.Lang, c.Source)
	}
}

func TestБезТерминалаНеСпрашиваютИНеПишут(t *testing.T) {
	// Труба, скрипт, CI, --json. Вопрос здесь повесил бы чужую работу
	// намертво, а файл в домашнем каталоге завёлся бы без спроса.
	for _, tc := range []struct {
		имя    string
		env    func(string) string
		ждём   lang.Lang
		откуда Source
	}{
		{"русская локаль", env("LANG", "ru_RU.UTF-8"), lang.RU, FromLocale},
		{"английская локаль", env("LC_ALL", "en_GB.UTF-8"), lang.EN, FromLocale},
		{"переносимая локаль", env("LANG", "C"), lang.Default, FromDefault},
		{"локали нет вовсе", env(), lang.Default, FromDefault},
		{"третий язык", env("LANG", "de_DE.UTF-8"), lang.Default, FromDefault},
	} {
		t.Run(tc.имя, func(t *testing.T) {
			home := дом(t)
			o := Options{Home: home, Getenv: tc.env}
			c := Decide(o, "", Ask{})
			if c.Lang != tc.ждём || c.Source != tc.откуда {
				t.Errorf("вышло %v из %v, ждали %v из %v", c.Lang, c.Source, tc.ждём, tc.откуда)
			}
			if естьФайл(файлНастроек(home)) {
				t.Error("завёлся файл настроек, а никто ни о чём не спрашивал")
			}
			if len(c.Notes) != 0 {
				t.Errorf("сказано лишнее: %v", c.Notes)
			}
		})
	}
}

func TestВопросЗадаётсяИОтветЗапоминается(t *testing.T) {
	home := дом(t)
	o := Options{Home: home, Getenv: env("LANG", "ru_RU.UTF-8")}
	var экран strings.Builder
	c := Decide(o, "", Ask{In: strings.NewReader("2\n"), Out: &экран, May: true})

	if c.Lang != lang.EN || c.Source != FromAsk {
		t.Fatalf("ответ не принят: %v из %v", c.Lang, c.Source)
	}
	// Вопрос задаётся на обоих языках сразу: на момент вопроса язык ещё
	// не выбран, и задать его на одном — значит выбрать за человека.
	спрошено := экран.String()
	for _, надо := range []string{"язык вывода", "output language", "русский", "English", "1)", "2)"} {
		if !strings.Contains(спрошено, надо) {
			t.Errorf("в вопросе нет %q:\n%s", надо, спрошено)
		}
	}
	// Локаль подсказывает умолчание и только: стрелка стоит у русского,
	// а выбран английский.
	if !strings.Contains(спрошено, "→ 1)") {
		t.Errorf("локаль не подсказала умолчание:\n%s", спрошено)
	}

	path := файлНастроек(home)
	if c.Saved != path {
		t.Errorf("сохранено в %q, ждали %q", c.Saved, path)
	}
	if !естьФайл(path) {
		t.Fatal("файл настроек не завёлся")
	}
	// Запись в домашний каталог — действие, и о нём говорят одной строкой.
	if len(c.Notes) != 1 || !strings.Contains(c.Notes[0].In(lang.EN), path) {
		t.Errorf("о заведённом файле не сказано: %v", c.Notes)
	}

	// И следующий запуск уже не спрашивает.
	again := Decide(o, "", Ask{In: strings.NewReader(""), Out: &экран, May: true})
	if again.Source != FromFile || again.Lang != lang.EN {
		t.Errorf("второй запуск снова спросил: %v из %v", again.Lang, again.Source)
	}
}

func TestПустойОтветБерётПодсказкуЛокали(t *testing.T) {
	home := дом(t)
	o := Options{Home: home, Getenv: env("LANG", "ru_RU.UTF-8")}
	var экран strings.Builder
	c := Decide(o, "", Ask{In: strings.NewReader("\n"), Out: &экран, May: true})
	if c.Lang != lang.RU || c.Source != FromAsk {
		t.Fatalf("Enter не взял подсказку: %v из %v", c.Lang, c.Source)
	}
}

func TestКонецВводаНичегоНеЗаводит(t *testing.T) {
	// Терминал закрылся, ответа нет. Прогон продолжается на языке машины, и
	// файл не заводится: никто ничего не выбрал.
	home := дом(t)
	o := Options{Home: home, Getenv: env("LANG", "ru_RU.UTF-8")}
	var экран strings.Builder
	c := Decide(o, "", Ask{In: strings.NewReader(""), Out: &экран, May: true})
	if c.Lang != lang.RU || c.Source != FromLocale {
		t.Fatalf("вышло %v из %v", c.Lang, c.Source)
	}
	if естьФайл(файлНастроек(home)) {
		t.Error("завёлся файл настроек, а ответа не было")
	}
}

func TestНедоступныйДомНеРонитПрогон(t *testing.T) {
	// Домашний каталог может быть недоступен для записи. Тогда инструмент
	// работает дальше на выбранном языке и внятно говорит, что не сохранил,
	// а не падает: маленькая беда не отвечается большой.
	if os.Geteuid() == 0 {
		t.Skip("под root права на запись не проверить")
	}
	home := дом(t)
	if err := os.Chmod(home, 0o500); err != nil {
		t.Fatalf("права не выставились: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(home, 0o700) })

	o := Options{Home: home, Getenv: env("LANG", "en_US.UTF-8")}
	var экран strings.Builder
	c := Decide(o, "", Ask{In: strings.NewReader("1\n"), Out: &экран, May: true})
	if c.Lang != lang.RU || c.Source != FromAsk {
		t.Fatalf("выбор не принят: %v из %v", c.Lang, c.Source)
	}
	if c.Saved != "" {
		t.Errorf("сохранение объявлено, а каталог недоступен: %s", c.Saved)
	}
	if len(c.Notes) != 1 || !strings.Contains(c.Notes[0].In(lang.RU), "не сохранён") {
		t.Errorf("не сказано, что не сохранил: %v", c.Notes)
	}
}

func TestНастройкиЧитаютсяНаДвухНаписаниях(t *testing.T) {
	// Файл читает и человек. Русский ключ в настройках англичанина — та же
	// невежливость, что и русский отчёт.
	for _, текст := range []string{"язык=ru\n", "lang=ru\n", "language = ru\n", "# так\nязык = ru\n"} {
		home := дом(t)
		dir := filepath.Join(home, FamilyDir, ToolDir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, FileName), []byte(текст), 0o644); err != nil {
			t.Fatal(err)
		}
		s, err := Load(Options{Home: home, Getenv: env()})
		if err != nil {
			t.Fatalf("%q: %v", текст, err)
		}
		if s.Lang != lang.RU {
			t.Errorf("%q прочиталось как %q", текст, s.Lang)
		}
	}
}

func TestИспорченныеНастройкиНеРонятПрогон(t *testing.T) {
	home := дом(t)
	dir := filepath.Join(home, FamilyDir, ToolDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, FileName), []byte("язык=клингонский\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	o := Options{Home: home, Getenv: env("LANG", "ru_RU.UTF-8")}
	c := Decide(o, "", Ask{})
	if c.Lang != lang.RU {
		t.Errorf("вышло %v, ждали язык машины", c.Lang)
	}
	if len(c.Notes) != 1 {
		t.Errorf("о непрочитанных настройках не сказано: %v", c.Notes)
	}
}

func TestНовыйДомПереднийСтарыйЧитается(t *testing.T) {
	home := дом(t)
	o := Options{Home: home, Getenv: env()}

	// Ничего нет.
	if _, _, ok := Find(o, "places.conf"); ok {
		t.Error("нашлось то, чего нет")
	}

	// Только старый дом — читается он, и это отмечается.
	old := filepath.Join(home, ".config", LegacyDir)
	if err := os.MkdirAll(old, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "places.conf"), []byte("# старый\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path, legacy, ok := Find(o, "places.conf")
	if !ok || !legacy || path != filepath.Join(old, "places.conf") {
		t.Fatalf("старый дом не прочитался: %q legacy=%v ok=%v", path, legacy, ok)
	}

	// Появился новый — берётся он.
	now := filepath.Join(home, FamilyDir, ToolDir)
	if err := os.MkdirAll(now, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(now, "places.conf"), []byte("# новый\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	path, legacy, ok = Find(o, "places.conf")
	if !ok || legacy || path != filepath.Join(now, "places.conf") {
		t.Fatalf("новый дом не перевесил: %q legacy=%v ok=%v", path, legacy, ok)
	}
}

func TestЗаписьНеТеряетЧужихКлючей(t *testing.T) {
	// Настройку новее этой версии старая версия обязана сохранить, а не
	// стереть, переписывая файл под себя.
	home := дом(t)
	dir := filepath.Join(home, FamilyDir, ToolDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, FileName), []byte("язык=ru\nпалитра=paper\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	o := Options{Home: home, Getenv: env()}
	s, err := Load(o)
	if err != nil {
		t.Fatal(err)
	}
	s.Lang = lang.EN
	if _, err := Save(o, s); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(dir, FileName))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "палитра=paper") {
		t.Errorf("чужой ключ стёрт:\n%s", body)
	}
	if !strings.Contains(string(body), "lang=en") {
		t.Errorf("новый язык не записан:\n%s", body)
	}
}
