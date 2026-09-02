// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestЛокальЧитаетсяПоПорядкуPOSIX(t *testing.T) {
	env := func(pairs ...string) func(string) string {
		m := map[string]string{}
		for i := 0; i+1 < len(pairs); i += 2 {
			m[pairs[i]] = pairs[i+1]
		}
		return func(n string) string { return m[n] }
	}
	for _, tc := range []struct {
		имя   string
		env   func(string) string
		ждём  Lang
		знаем bool
	}{
		{"LC_ALL перевешивает всё", env("LC_ALL", "ru_RU.UTF-8", "LC_MESSAGES", "en_US", "LANG", "en_US"), RU, true},
		{"LC_MESSAGES перевешивает LANG", env("LC_MESSAGES", "en_GB", "LANG", "ru_RU.UTF-8"), EN, true},
		{"LANG, когда больше нечего", env("LANG", "ru_RU.UTF-8"), RU, true},
		{"C — это ответ, но не выбор", env("LANG", "C"), Default, false},
		{"POSIX — то же самое", env("LC_ALL", "POSIX"), Default, false},
		{"пусто", env(), Default, false},
		{"третий язык — не наш", env("LANG", "de_DE.UTF-8"), Default, false},
		{"третий язык не пропускает к LANG", env("LC_ALL", "de_DE.UTF-8", "LANG", "ru_RU.UTF-8"), Default, false},
		{"без кодировки", env("LANG", "ru"), RU, true},
		{"с модификатором", env("LANG", "en@quot"), EN, true},
	} {
		t.Run(tc.имя, func(t *testing.T) {
			l, ok := Locale(tc.env)
			if l != tc.ждём || ok != tc.знаем {
				t.Errorf("вышло %v/%v, ждали %v/%v", l, ok, tc.ждём, tc.знаем)
			}
		})
	}
	if _, ok := Locale(nil); ok && os.Getenv("LANG") == "" {
		t.Error("nil-окружение должно читаться из процесса")
	}
}

func TestРазборЯзыкаБерётТоЧтоЧеловекНаберёт(t *testing.T) {
	for _, s := range []string{"ru", "RU", " ru ", "рус", "русский", "1", "russian"} {
		if l, ok := Parse(s); !ok || l != RU {
			t.Errorf("%q не разобралось как русский", s)
		}
	}
	for _, s := range []string{"en", "EN", "English", "англ", "2", "eng"} {
		if l, ok := Parse(s); !ok || l != EN {
			t.Errorf("%q не разобралось как английский", s)
		}
	}
	for _, s := range []string{"", "de", "3", "клингонский", "ru_RU"} {
		if _, ok := Parse(s); ok {
			t.Errorf("%q разобралось, а не должно было", s)
		}
	}
}

// TestЧислаПишутсяПоЯзыку — десятичная запятая против точки, разделитель
// разрядов, порядок в дате. Число, читаемое по-разному, хуже числа на чужом
// языке: чужой язык виден, а «12,3» как английское — это двенадцать и три.
func TestЧислаПишутсяПоЯзыку(t *testing.T) {
	for _, tc := range []struct{ ru, en, что string }{
		{RU.Dec(12.34, 2), "12.34", "дробное"},
		{RU.Pct(99.5, 1), "99.5%", "проценты"},
		{RU.Num(1234567), "1,234,567", "разряды"},
		{RU.Bytes(1536), "1.5 KiB", "килобайты"},
		{RU.Bytes(512), "512 B", "байты"},
		{RU.Days(41), "41 d", "дни"},
	} {
		_ = tc
	}
	if got := RU.Dec(12.34, 2); got != "12,34" {
		t.Errorf("по-русски дробное вышло %q", got)
	}
	if got := EN.Dec(12.34, 2); got != "12.34" {
		t.Errorf("англ дробное вышло %q", got)
	}
	if got := RU.Num(1234567); got != "1 234 567" {
		t.Errorf("по-русски разряды вышли %q", got)
	}
	if got := EN.Num(1234567); got != "1,234,567" {
		t.Errorf("англ разряды вышли %q", got)
	}
	if got := RU.Bytes(1536); got != "1,5 КиБ" {
		t.Errorf("по-русски полтора кибибайта вышли %q", got)
	}
	if got := EN.Bytes(1536); got != "1.5 KiB" {
		t.Errorf("англ полтора кибибайта вышли %q", got)
	}
	if got := RU.Bytes(-2048); got != "-2,0 КиБ" {
		t.Errorf("отрицательное вышло %q", got)
	}
	when := time.Date(2026, 9, 2, 15, 4, 0, 0, time.UTC)
	if got := RU.DateTime(when); got != "02.09.2026 15:04" {
		t.Errorf("русская дата вышла %q", got)
	}
	if got := EN.DateTime(when); got != "2026-09-02 15:04" {
		t.Errorf("английская дата вышла %q", got)
	}
	if got := RU.Uptime(5*24*3600 + 3*3600 + 14*60); got != "5д 03:14" {
		t.Errorf("русское время работы вышло %q", got)
	}
	if got := EN.Uptime(5*24*3600 + 3*3600 + 14*60); got != "5d 03:14" {
		t.Errorf("английское время работы вышло %q", got)
	}
}

// TestPhraseЕдетВJSONПоРусски — сердце обещания «--json не переводится».
func TestPhraseЕдетВJSONПоРусски(t *testing.T) {
	p := Say("путь вне указанного корня")
	got, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `"путь вне указанного корня"` {
		t.Errorf("в JSON поехало %s", got)
	}
	if p.In(RU) != "путь вне указанного корня" {
		t.Errorf("по-русски вышло %q", p.In(RU))
	}
	if p.In(EN) == p.In(RU) {
		t.Errorf("англ не перевелось: %q", p.In(EN))
	}

	// С доводами — те же байты, что дал бы Sprintf.
	q := Say("нужна ровно одна корзина, получено %d", 3)
	if q.String() != "нужна ровно одна корзина, получено 3" {
		t.Errorf("русский рендер вышел %q", q.String())
	}

	// Текст системы не переводится ни на каком языке.
	r := FromError(errors.New("permission denied"))
	if r.In(RU) != "permission denied" || r.In(EN) != "permission denied" {
		t.Errorf("сообщение системы тронуто: %q / %q", r.In(RU), r.In(EN))
	}
	if FromError(nil).Empty() != true {
		t.Error("nil-ошибка не пуста")
	}
}

func TestPhraseЧитаетсяИзЖурналаОбратно(t *testing.T) {
	// Слово словаря возвращается словом и переводится снова; всё прочее —
	// записью, какой её сделали.
	var p Phrase
	if err := json.Unmarshal([]byte(`"путь вне указанного корня"`), &p); err != nil {
		t.Fatal(err)
	}
	if p.Wording() == "" {
		t.Error("известное слово вернулось сырым текстом")
	}
	var q Phrase
	if err := json.Unmarshal([]byte(`"нужна ровно одна корзина, получено 3"`), &q); err != nil {
		t.Fatal(err)
	}
	if q.Wording() != "" {
		t.Error("подставленное предложение вернулось как слово словаря")
	}
	if q.In(EN) != "нужна ровно одна корзина, получено 3" {
		t.Errorf("запись журнала переписана: %q", q.In(EN))
	}
}

func TestОшибкаПечатаетсяНаЯзыкеЧитателяИЛовитсяerrorsIs(t *testing.T) {
	корень := errors.New("файла нет")
	err := Errorf("справочник %s: %s", "встроенный", корень)
	if err.Error() != "справочник встроенный: файла нет" {
		t.Errorf("русский текст ошибки вышел %q", err.Error())
	}
	if !errors.Is(err, корень) {
		t.Error("errors.Is не видит обёрнутую ошибку")
	}
	if InLang(err, EN) == err.Error() {
		t.Errorf("англ не перевелось: %q", InLang(err, EN))
	}
	if InLang(корень, EN) != "файла нет" {
		t.Error("чужая ошибка тронута")
	}
	if InLang(nil, EN) != "" {
		t.Error("nil дал не пустую строку")
	}
}

func TestСловоЯдраПереводитсяТолькоНаЭкране(t *testing.T) {
	// Значение не меняется — меняется слово, которым его называют.
	const разряд = "Кэш"
	if RU.Word(разряд) != разряд {
		t.Errorf("по-русски разряд стал %q", RU.Word(разряд))
	}
	if EN.Word(разряд) == разряд {
		t.Errorf("англ разряд не назван словом: %q", EN.Word(разряд))
	}
	// Слово, которого словарь не знает, показывается как есть: справочник
	// человек пишет сам, и отказаться его показать было бы хуже.
	if EN.Word("Своё") != "Своё" {
		t.Errorf("незнакомое слово подменено: %q", EN.Word("Своё"))
	}
}

func TestДругойЯзыкЭтоВторойИзДвух(t *testing.T) {
	if RU.Other() != EN || EN.Other() != RU {
		t.Error("переключение языка не круговое")
	}
	if !RU.Valid() || !EN.Valid() || Lang("de").Valid() {
		t.Error("Valid пропускает третий язык")
	}
}

// TestВложеннаяФразаПереводитсяТоже — предложение, построенное вокруг другого
// предложения, обязано перевестись целиком. Наполовину переведённая строка —
// ровно та беда, ради которой всё это писалось.
func TestВложеннаяФразаПереводитсяТоже(t *testing.T) {
	внутри := Errorf("путь вне указанного корня")
	снаружи := Errorf("справочник %s: %s", "/etc/digitdisk.conf", внутри)

	if снаружи.Error() != "справочник /etc/digitdisk.conf: путь вне указанного корня" {
		t.Errorf("русский рендер вышел %q", снаружи.Error())
	}
	англ := InLang(снаружи, EN)
	if strings.ContainsAny(англ, "абвгдеёжзийклмнопрстуфхцчшщъыьэюя") {
		t.Errorf("внутренность осталась по-русски: %q", англ)
	}
}
