// Сгенерировано flang (бэкенд Go, flang/src/emit/go.mjs). Не редактировать руками.
// Модуль flang: «Опись диска».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flangprogram/flangrt"
)

// NewContext — контекст вычисления с настройками этой программы.
//
// Индексация строк, предел глубины вызовов и лимит шагов — это настройки
// программы, а не рантайма: печать могла идти с нулевой базой индексации,
// а пределы вызывающий вправе поменять прямо в возвращённом контексте.
func NewContext() *rt.Ctx {
	ctx := rt.NewCtx()
	ctx.IndexBase = 1
	ctx.MaxDepth = 10000
	ctx.MaxSteps = 1000000
	return ctx
}

// SozdatNahodka — запись FTS «Находка»: «путь», «размер», «возраст_дней», «вид», «доступен».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatNahodka(put rt.Value, razmer rt.Value, vozrastDney rt.Value, vid rt.Value, dostupen rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "путь", Value: put},
		{Name: "размер", Value: razmer},
		{Name: "возраст_дней", Value: vozrastDney},
		{Name: "вид", Value: vid},
		{Name: "доступен", Value: dostupen},
	})
}

// SozdatMesto — запись FTS «Место»: «разряд», «якорь», «цепь».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatMesto(razryad rt.Value, yakor rt.Value, cep rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "разряд", Value: razryad},
		{Name: "якорь", Value: yakor},
		{Name: "цепь", Value: cep},
	})
}

// SozdatReshenie — запись FTS «Решение»: «разряд», «приговор», «вес».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatReshenie(razryad rt.Value, prigovor rt.Value, ves rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "разряд", Value: razryad},
		{Name: "приговор", Value: prigovor},
		{Name: "вес", Value: ves},
	})
}

// SozdatStrokaOtchyota — запись FTS «Строка отчёта»: «путь», «решение».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatStrokaOtchyota(put rt.Value, reshenie rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "путь", Value: put},
		{Name: "решение", Value: reshenie},
	})
}

// SozdatSvod — запись FTS «Свод»: «кэш», «журнал», «сборка», «загрузка», «крупное», «неизвестное», «освободить».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatSvod(kesh rt.Value, zhurnal rt.Value, sborka rt.Value, zagruzka rt.Value, krupnoe rt.Value, neizvestnoe rt.Value, osvobodit rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "кэш", Value: kesh},
		{Name: "журнал", Value: zhurnal},
		{Name: "сборка", Value: sborka},
		{Name: "загрузка", Value: zagruzka},
		{Name: "крупное", Value: krupnoe},
		{Name: "неизвестное", Value: neizvestnoe},
		{Name: "освободить", Value: osvobodit},
	})
}

// VariantFayl — вариант «Файл» суммы типов «Вид».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantFayl() rt.Value {
	return rt.Variant("Файл", nil)
}

// VariantKatalog — вариант «Каталог» суммы типов «Вид».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantKatalog() rt.Value {
	return rt.Variant("Каталог", nil)
}

// VariantSsylka — вариант «Ссылка» суммы типов «Вид».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantSsylka() rt.Value {
	return rt.Variant("Ссылка", nil)
}

// VariantKesh — вариант «Кэш» суммы типов «Разряд».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantKesh() rt.Value {
	return rt.Variant("Кэш", nil)
}

// VariantZhurnal — вариант «Журнал» суммы типов «Разряд».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantZhurnal() rt.Value {
	return rt.Variant("Журнал", nil)
}

// VariantSborka — вариант «Сборка» суммы типов «Разряд».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantSborka() rt.Value {
	return rt.Variant("Сборка", nil)
}

// VariantZagruzka — вариант «Загрузка» суммы типов «Разряд».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantZagruzka() rt.Value {
	return rt.Variant("Загрузка", nil)
}

// VariantKrupnoe — вариант «Крупное» суммы типов «Разряд».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantKrupnoe() rt.Value {
	return rt.Variant("Крупное", nil)
}

// VariantNeizvestnoe — вариант «Неизвестное» суммы типов «Разряд».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantNeizvestnoe() rt.Value {
	return rt.Variant("Неизвестное", nil)
}

// VariantMozhnoUbrat — вариант «МожноУбрать» суммы типов «Приговор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantMozhnoUbrat() rt.Value {
	return rt.Variant("МожноУбрать", nil)
}

// VariantSprosit — вариант «Спросить» суммы типов «Приговор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantSprosit() rt.Value {
	return rt.Variant("Спросить", nil)
}

// VariantNeTrogat — вариант «НеТрогать» суммы типов «Приговор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantNeTrogat() rt.Value {
	return rt.Variant("НеТрогать", nil)
}

// VariantOtKornya — вариант «ОтКорня» суммы типов «Якорь».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantOtKornya() rt.Value {
	return rt.Variant("ОтКорня", nil)
}

// VariantGdeUgodno — вариант «ГдеУгодно» суммы типов «Якорь».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantGdeUgodno() rt.Value {
	return rt.Variant("ГдеУгодно", nil)
}

// PorogKrupnogo — функция flang «Порог крупного».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: число.
func PorogKrupnogo(ctx *rt.Ctx) (rt.Value, error) {
	return rt.Number(1073741824.0), nil
}

// PorogKesha — функция flang «Порог кэша».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: число.
func PorogKesha(ctx *rt.Ctx) (rt.Value, error) {
	return rt.Number(7.0), nil
}

// PorogZhurnala — функция flang «Порог журнала».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: число.
func PorogZhurnala(ctx *rt.Ctx) (rt.Value, error) {
	return rt.Number(30.0), nil
}

// PorogZagruzki — функция flang «Порог загрузки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: число.
func PorogZagruzki(ctx *rt.Ctx) (rt.Value, error) {
	return rt.Number(180.0), nil
}

// SostavlyayuschiePuti — функция flang «Составляющие пути».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение: список: строка.
func SostavlyayuschiePuti(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	// «разделить»
	t1, e2 := rt.BSplitProven(ctx, put, rt.Text("/"))
	if e2 != nil {
		return rt.Value{}, e2
	}
	return t1, nil
}

// ImyaVPuti — функция flang «Имя в пути».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение: строка.
func ImyaVPuti(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t3, e4 := SostavlyayuschiePuti(ctx, put)
	if e4 != nil {
		return rt.Value{}, e4
	}
	t5, e6 := rt.RequireList(ctx, t3, "свёртка")
	if e6 != nil {
		return rt.Value{}, e6
	}
	// «собрано»
	sobrano := rt.Text("")
	for t7 := range t5 {
		// «часть»
		chast := t5[t7]
		sobrano = chast
	}
	return sobrano, nil
}

// EstSostavlyayuschaya — функция flang «Есть составляющая».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Параметр imya — «имя»: строка.
// Результат — значение.
func EstSostavlyayuschaya(ctx *rt.Ctx, put rt.Value, imya rt.Value) (rt.Value, error) {
	t8, e9 := SostavlyayuschiePuti(ctx, put)
	if e9 != nil {
		return rt.Value{}, e9
	}
	// «содержит»
	t10, e11 := rt.BContains(ctx, t8, imya)
	if e11 != nil {
		return rt.Value{}, e11
	}
	return t10, nil
}

// OkanchivaetsyaNa — функция flang «Оканчивается на».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр hvost — «хвост»: строка.
// Результат — значение.
func OkanchivaetsyaNa(ctx *rt.Ctx, tekst rt.Value, hvost rt.Value) (rt.Value, error) {
	// «длина»
	t12, e13 := rt.BLength(ctx, tekst)
	if e13 != nil {
		return rt.Value{}, e13
	}
	// «длина»
	t14, e15 := rt.BLength(ctx, hvost)
	if e15 != nil {
		return rt.Value{}, e15
	}
	t16, e17 := rt.Lt(ctx, t12, t14)
	if e17 != nil {
		return rt.Value{}, e17
	}
	t18, e19 := rt.Cond(ctx, t16)
	if e19 != nil {
		return rt.Value{}, e19
	}
	if t18 {
		return rt.Flag(false), nil
	} else {
		// «длина»
		t20, e21 := rt.BLength(ctx, tekst)
		if e21 != nil {
			return rt.Value{}, e21
		}
		// «длина»
		t22, e23 := rt.BLength(ctx, hvost)
		if e23 != nil {
			return rt.Value{}, e23
		}
		t24, e25 := rt.Sub(ctx, t20, t22)
		if e25 != nil {
			return rt.Value{}, e25
		}
		t26, e27 := rt.Add(ctx, t24, rt.Number(1.0))
		if e27 != nil {
			return rt.Value{}, e27
		}
		// «длина»
		t28, e29 := rt.BLength(ctx, tekst)
		if e29 != nil {
			return rt.Value{}, e29
		}
		// «подстрока»
		t30, e31 := rt.BSubstring(ctx, tekst, t26, t28)
		if e31 != nil {
			return rt.Value{}, e31
		}
		return rt.Flag(rt.Equal(t30, hvost)), nil
	}
}

// SledPuti — функция flang «След пути».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение: строка.
func SledPuti(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t32, e33 := rt.Concat(ctx, put, rt.Text("/"))
	if e33 != nil {
		return rt.Value{}, e33
	}
	return t32, nil
}

// CepOgranichena — функция flang «Цепь ограничена».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр cep — «цепь»: строка.
// Результат — значение.
func CepOgranichena(ctx *rt.Ctx, cep rt.Value) (rt.Value, error) {
	// «начинается с»
	t34, e35 := rt.BStartsWith(ctx, cep, rt.Text("/"))
	if e35 != nil {
		return rt.Value{}, e35
	}
	t36, e37 := rt.Cond(ctx, t34)
	if e37 != nil {
		return rt.Value{}, e37
	}
	var t38 rt.Value
	if t36 {
		t39, e40 := OkanchivaetsyaNa(ctx, cep, rt.Text("/"))
		if e40 != nil {
			return rt.Value{}, e40
		}
		t38 = t39
	} else {
		t38 = rt.Flag(false)
	}
	t41, e42 := rt.Cond(ctx, t38)
	if e42 != nil {
		return rt.Value{}, e42
	}
	if t41 {
		// «длина»
		t43, e44 := rt.BLength(ctx, cep)
		if e44 != nil {
			return rt.Value{}, e44
		}
		t45, e46 := rt.Gt(ctx, t43, rt.Number(1.0))
		if e46 != nil {
			return rt.Value{}, e46
		}
		return t45, nil
	} else {
		return rt.Flag(false), nil
	}
}

// SpravochnikOgranichen — функция flang «Справочник ограничен».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение.
func SpravochnikOgranichen(ctx *rt.Ctx, spravochnik rt.Value) (rt.Value, error) {
	t47, e48 := rt.RequireList(ctx, spravochnik, "свёртка")
	if e48 != nil {
		return rt.Value{}, e48
	}
	// «акк»
	akk := rt.Flag(true)
	for t49 := range t47 {
		// «место»
		mesto := t47[t49]
		t50, e51 := rt.Cond(ctx, akk)
		if e51 != nil {
			return rt.Value{}, e51
		}
		var t52 rt.Value
		if t50 {
			t53, e54 := rt.FieldGet(ctx, mesto, "цепь")
			if e54 != nil {
				return rt.Value{}, e54
			}
			t55, e56 := CepOgranichena(ctx, t53)
			if e56 != nil {
				return rt.Value{}, e56
			}
			t52 = t55
		} else {
			t52 = rt.Flag(false)
		}
		akk = t52
	}
	return akk, nil
}

// RazryadMestaDopustim — функция flang «Разряд места допустим».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func RazryadMestaDopustim(ctx *rt.Ctx, razryad rt.Value) (rt.Value, error) {
	if rt.VariantIs(razryad, "Кэш") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Журнал") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Сборка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Загрузка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Крупное") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Неизвестное") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// MestoPodhodit — функция flang «Место подходит».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sled — «след»: строка.
// Параметр mesto — «место»: «Место».
// Результат — значение.
func MestoPodhodit(ctx *rt.Ctx, sled rt.Value, mesto rt.Value) (rt.Value, error) {
	t57, e58 := rt.FieldGet(ctx, mesto, "якорь")
	if e58 != nil {
		return rt.Value{}, e58
	}
	if rt.VariantIs(t57, "ОтКорня") {
		t59, e60 := rt.FieldGet(ctx, mesto, "цепь")
		if e60 != nil {
			return rt.Value{}, e60
		}
		// «начинается с»
		t61, e62 := rt.BStartsWith(ctx, sled, t59)
		if e62 != nil {
			return rt.Value{}, e62
		}
		return t61, nil
	} else if rt.VariantIs(t57, "ГдеУгодно") {
		t63, e64 := rt.FieldGet(ctx, mesto, "цепь")
		if e64 != nil {
			return rt.Value{}, e64
		}
		// «содержит»
		t65, e66 := rt.BContains(ctx, sled, t63)
		if e66 != nil {
			return rt.Value{}, e66
		}
		return t65, nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t57)
	}
}

// EtoNeizvestnoe — функция flang «Это неизвестное».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func EtoNeizvestnoe(ctx *rt.Ctx, razryad rt.Value) (rt.Value, error) {
	if rt.VariantIs(razryad, "Неизвестное") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Кэш") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Журнал") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Сборка") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Загрузка") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Крупное") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// NomerRazryada — функция flang «Номер разряда».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение: число.
func NomerRazryada(ctx *rt.Ctx, razryad rt.Value) (rt.Value, error) {
	if rt.VariantIs(razryad, "Кэш") {
		return rt.Number(1.0), nil
	} else if rt.VariantIs(razryad, "Журнал") {
		return rt.Number(2.0), nil
	} else if rt.VariantIs(razryad, "Сборка") {
		return rt.Number(3.0), nil
	} else if rt.VariantIs(razryad, "Загрузка") {
		return rt.Number(4.0), nil
	} else if rt.VariantIs(razryad, "Крупное") {
		return rt.Number(5.0), nil
	} else if rt.VariantIs(razryad, "Неизвестное") {
		return rt.Number(6.0), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// TotZheRazryad — функция flang «Тот же разряд».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр pervyy — «первый»: «Разряд».
// Параметр vtoroy — «второй»: «Разряд».
// Результат — значение.
func TotZheRazryad(ctx *rt.Ctx, pervyy rt.Value, vtoroy rt.Value) (rt.Value, error) {
	t67, e68 := NomerRazryada(ctx, pervyy)
	if e68 != nil {
		return rt.Value{}, e68
	}
	t69, e70 := NomerRazryada(ctx, vtoroy)
	if e70 != nil {
		return rt.Value{}, e70
	}
	return rt.Flag(rt.Equal(t67, t69)), nil
}

// MestoObosnovano — функция flang «Место обосновано».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Параметр spravochnik — «справочник»: список: «Место».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func MestoObosnovano(ctx *rt.Ctx, put rt.Value, spravochnik rt.Value, razryad rt.Value) (rt.Value, error) {
	t71, e72 := EtoNeizvestnoe(ctx, razryad)
	if e72 != nil {
		return rt.Value{}, e72
	}
	t73, e74 := rt.Cond(ctx, t71)
	if e74 != nil {
		return rt.Value{}, e74
	}
	if t73 {
		return rt.Flag(true), nil
	} else {
		t75, e76 := RazryadMestaDopustim(ctx, razryad)
		if e76 != nil {
			return rt.Value{}, e76
		}
		t77, e78 := rt.Cond(ctx, t75)
		if e78 != nil {
			return rt.Value{}, e78
		}
		if t77 {
			t79, e80 := rt.RequireList(ctx, spravochnik, "отфильтровать")
			if e80 != nil {
				return rt.Value{}, e80
			}
			t81 := make([]rt.Value, 0, len(t79))
			for t82 := range t79 {
				// «место»
				mesto := t79[t82]
				t83, e84 := SledPuti(ctx, put)
				if e84 != nil {
					return rt.Value{}, e84
				}
				t85, e86 := MestoPodhodit(ctx, t83, mesto)
				if e86 != nil {
					return rt.Value{}, e86
				}
				t87, e88 := rt.Cond(ctx, t85)
				if e88 != nil {
					return rt.Value{}, e88
				}
				var t89 rt.Value
				if t87 {
					t90, e91 := rt.FieldGet(ctx, mesto, "разряд")
					if e91 != nil {
						return rt.Value{}, e91
					}
					t92, e93 := TotZheRazryad(ctx, t90, razryad)
					if e93 != nil {
						return rt.Value{}, e93
					}
					t89 = t92
				} else {
					t89 = rt.Flag(false)
				}
				t94, e95 := rt.Keep(ctx, t89)
				if e95 != nil {
					return rt.Value{}, e95
				}
				if t94 {
					t81 = append(t81, mesto)
				}
			}
			// «длина»
			t96, e97 := rt.BLength(ctx, rt.List(t81))
			if e97 != nil {
				return rt.Value{}, e97
			}
			t98, e99 := rt.Gt(ctx, t96, rt.Number(0.0))
			if e99 != nil {
				return rt.Value{}, e99
			}
			return t98, nil
		} else {
			return rt.Flag(false), nil
		}
	}
}

// RazryadPoSpravochniku — функция flang «Разряд по справочнику».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: «Разряд».
func RazryadPoSpravochniku(ctx *rt.Ctx, put rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t100, e101 := SledPuti(ctx, put)
	if e101 != nil {
		return rt.Value{}, e101
	}
	// пусть «след»
	sled := t100
	t102, e103 := rt.RequireList(ctx, spravochnik, "свёртка")
	if e103 != nil {
		return rt.Value{}, e103
	}
	// «найдено»
	naydeno := rt.Variant("Неизвестное", nil)
	for t104 := range t102 {
		// «место»
		mesto := t102[t104]
		t105, e106 := EtoNeizvestnoe(ctx, naydeno)
		if e106 != nil {
			return rt.Value{}, e106
		}
		t107, e108 := rt.Cond(ctx, t105)
		if e108 != nil {
			return rt.Value{}, e108
		}
		var t109 rt.Value
		if t107 {
			t110, e111 := MestoPodhodit(ctx, sled, mesto)
			if e111 != nil {
				return rt.Value{}, e111
			}
			t109 = t110
		} else {
			t109 = rt.Flag(false)
		}
		t112, e113 := rt.Cond(ctx, t109)
		if e113 != nil {
			return rt.Value{}, e113
		}
		var t114 rt.Value
		if t112 {
			t115, e116 := rt.FieldGet(ctx, mesto, "разряд")
			if e116 != nil {
				return rt.Value{}, e116
			}
			t114 = t115
		} else {
			t114 = naydeno
		}
		naydeno = t114
	}
	t117 := naydeno
	t118, e119 := MestoObosnovano(ctx, put, spravochnik, t117)
	if e119 != nil {
		return rt.Value{}, e119
	}
	// постусловие «Место обосновано записью справочника»
	t120, e121 := rt.Post(ctx, t118, "Место обосновано записью справочника", "Разряд по справочнику")
	if e121 != nil {
		return rt.Value{}, e121
	}
	if !t120 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Место обосновано записью справочника» функции «Разряд по справочнику»")
	}
	return t117, nil
}

// ShestnadcaterichnyyZnak — функция flang «Шестнадцатеричный знак».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znak — «знак»: строка.
// Результат — значение.
func ShestnadcaterichnyyZnak(ctx *rt.Ctx, znak rt.Value) (rt.Value, error) {
	// «содержит»
	t122, e123 := rt.BContains(ctx, rt.Text("0123456789abcdefABCDEF"), znak)
	if e123 != nil {
		return rt.Value{}, e123
	}
	return t122, nil
}

// PohozheNaOtpechatok — функция flang «Похоже на отпечаток».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр chast — «часть»: строка.
// Результат — значение.
func PohozheNaOtpechatok(ctx *rt.Ctx, chast rt.Value) (rt.Value, error) {
	// «длина»
	t124, e125 := rt.BLength(ctx, chast)
	if e125 != nil {
		return rt.Value{}, e125
	}
	t126, e127 := rt.Lt(ctx, t124, rt.Number(32.0))
	if e127 != nil {
		return rt.Value{}, e127
	}
	t128, e129 := rt.Cond(ctx, t126)
	if e129 != nil {
		return rt.Value{}, e129
	}
	if t128 {
		return rt.Flag(false), nil
	} else {
		// «символы»
		t130, e131 := rt.BCharacters(ctx, chast)
		if e131 != nil {
			return rt.Value{}, e131
		}
		t132, e133 := rt.RequireList(ctx, t130, "свёртка")
		if e133 != nil {
			return rt.Value{}, e133
		}
		// «собрано»
		sobrano := rt.Flag(true)
		for t134 := range t132 {
			// «знак»
			znak := t132[t134]
			t135, e136 := rt.Cond(ctx, sobrano)
			if e136 != nil {
				return rt.Value{}, e136
			}
			var t137 rt.Value
			if t135 {
				t138, e139 := ShestnadcaterichnyyZnak(ctx, znak)
				if e139 != nil {
					return rt.Value{}, e139
				}
				t137 = t138
			} else {
				t137 = rt.Flag(false)
			}
			sobrano = t137
		}
		return sobrano, nil
	}
}

// AdresuetsyaSoderzhimym — функция flang «Адресуется содержимым».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func AdresuetsyaSoderzhimym(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t140, e141 := SostavlyayuschiePuti(ctx, put)
	if e141 != nil {
		return rt.Value{}, e141
	}
	t142, e143 := rt.RequireList(ctx, t140, "отфильтровать")
	if e143 != nil {
		return rt.Value{}, e143
	}
	t144 := make([]rt.Value, 0, len(t142))
	for t145 := range t142 {
		// «часть»
		chast := t142[t145]
		t146, e147 := PohozheNaOtpechatok(ctx, chast)
		if e147 != nil {
			return rt.Value{}, e147
		}
		t148, e149 := rt.Keep(ctx, t146)
		if e149 != nil {
			return rt.Value{}, e149
		}
		if t148 {
			t144 = append(t144, chast)
		}
	}
	// «длина»
	t150, e151 := rt.BLength(ctx, rt.List(t144))
	if e151 != nil {
		return rt.Value{}, e151
	}
	t152, e153 := rt.Gt(ctx, t150, rt.Number(0.0))
	if e153 != nil {
		return rt.Value{}, e153
	}
	t154, e155 := rt.Cond(ctx, t152)
	if e155 != nil {
		return rt.Value{}, e155
	}
	var t156 rt.Value
	if t154 {
		t156 = rt.Flag(true)
	} else {
		t157, e158 := EstSostavlyayuschaya(ctx, put, rt.Text("site-packages"))
		if e158 != nil {
			return rt.Value{}, e158
		}
		t156 = t157
	}
	t159, e160 := rt.Cond(ctx, t156)
	if e160 != nil {
		return rt.Value{}, e160
	}
	var t161 rt.Value
	if t159 {
		t161 = rt.Flag(true)
	} else {
		t162, e163 := EstSostavlyayuschaya(ctx, put, rt.Text("dist-packages"))
		if e163 != nil {
			return rt.Value{}, e163
		}
		t161 = t162
	}
	t164, e165 := rt.Cond(ctx, t161)
	if e165 != nil {
		return rt.Value{}, e165
	}
	if t164 {
		return rt.Flag(true), nil
	} else {
		return EstSostavlyayuschaya(ctx, put, rt.Text(".git"))
	}
}

// PodSistemnymVremennym — функция flang «Под системным временным».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PodSistemnymVremennym(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	// «начинается с»
	t166, e167 := rt.BStartsWith(ctx, put, rt.Text("/tmp/"))
	if e167 != nil {
		return rt.Value{}, e167
	}
	t168, e169 := rt.Cond(ctx, t166)
	if e169 != nil {
		return rt.Value{}, e169
	}
	if t168 {
		return rt.Flag(true), nil
	} else {
		// «начинается с»
		t170, e171 := rt.BStartsWith(ctx, put, rt.Text("/var/tmp/"))
		if e171 != nil {
			return rt.Value{}, e171
		}
		return t170, nil
	}
}

// PrimetaKesha — функция flang «Примета кэша».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaKesha(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t172, e173 := EstSostavlyayuschaya(ctx, put, rt.Text(".cache"))
	if e173 != nil {
		return rt.Value{}, e173
	}
	t174, e175 := rt.Cond(ctx, t172)
	if e175 != nil {
		return rt.Value{}, e175
	}
	var t176 rt.Value
	if t174 {
		t176 = rt.Flag(true)
	} else {
		t177, e178 := EstSostavlyayuschaya(ctx, put, rt.Text("cache"))
		if e178 != nil {
			return rt.Value{}, e178
		}
		t176 = t177
	}
	t179, e180 := rt.Cond(ctx, t176)
	if e180 != nil {
		return rt.Value{}, e180
	}
	var t181 rt.Value
	if t179 {
		t181 = rt.Flag(true)
	} else {
		t182, e183 := EstSostavlyayuschaya(ctx, put, rt.Text("Caches"))
		if e183 != nil {
			return rt.Value{}, e183
		}
		t181 = t182
	}
	t184, e185 := rt.Cond(ctx, t181)
	if e185 != nil {
		return rt.Value{}, e185
	}
	if t184 {
		return rt.Flag(true), nil
	} else {
		return PodSistemnymVremennym(ctx, put)
	}
}

// PrimetaZhurnala — функция flang «Примета журнала».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaZhurnala(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t186, e187 := ImyaVPuti(ctx, put)
	if e187 != nil {
		return rt.Value{}, e187
	}
	t188, e189 := OkanchivaetsyaNa(ctx, t186, rt.Text(".log"))
	if e189 != nil {
		return rt.Value{}, e189
	}
	t190, e191 := rt.Cond(ctx, t188)
	if e191 != nil {
		return rt.Value{}, e191
	}
	var t192 rt.Value
	if t190 {
		t192 = rt.Flag(true)
	} else {
		t193, e194 := EstSostavlyayuschaya(ctx, put, rt.Text("log"))
		if e194 != nil {
			return rt.Value{}, e194
		}
		t192 = t193
	}
	t195, e196 := rt.Cond(ctx, t192)
	if e196 != nil {
		return rt.Value{}, e196
	}
	if t195 {
		return rt.Flag(true), nil
	} else {
		return EstSostavlyayuschaya(ctx, put, rt.Text("logs"))
	}
}

// PrimetaSborki — функция flang «Примета сборки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaSborki(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t197, e198 := EstSostavlyayuschaya(ctx, put, rt.Text("node_modules"))
	if e198 != nil {
		return rt.Value{}, e198
	}
	t199, e200 := rt.Cond(ctx, t197)
	if e200 != nil {
		return rt.Value{}, e200
	}
	var t201 rt.Value
	if t199 {
		t201 = rt.Flag(true)
	} else {
		t202, e203 := EstSostavlyayuschaya(ctx, put, rt.Text("target"))
		if e203 != nil {
			return rt.Value{}, e203
		}
		t201 = t202
	}
	t204, e205 := rt.Cond(ctx, t201)
	if e205 != nil {
		return rt.Value{}, e205
	}
	var t206 rt.Value
	if t204 {
		t206 = rt.Flag(true)
	} else {
		t207, e208 := EstSostavlyayuschaya(ctx, put, rt.Text("build"))
		if e208 != nil {
			return rt.Value{}, e208
		}
		t206 = t207
	}
	t209, e210 := rt.Cond(ctx, t206)
	if e210 != nil {
		return rt.Value{}, e210
	}
	var t211 rt.Value
	if t209 {
		t211 = rt.Flag(true)
	} else {
		t212, e213 := EstSostavlyayuschaya(ctx, put, rt.Text("_build"))
		if e213 != nil {
			return rt.Value{}, e213
		}
		t211 = t212
	}
	t214, e215 := rt.Cond(ctx, t211)
	if e215 != nil {
		return rt.Value{}, e215
	}
	if t214 {
		return rt.Flag(true), nil
	} else {
		return EstSostavlyayuschaya(ctx, put, rt.Text(".gradle"))
	}
}

// PrimetaZagruzki — функция flang «Примета загрузки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaZagruzki(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t216, e217 := EstSostavlyayuschaya(ctx, put, rt.Text("Downloads"))
	if e217 != nil {
		return rt.Value{}, e217
	}
	t218, e219 := rt.Cond(ctx, t216)
	if e219 != nil {
		return rt.Value{}, e219
	}
	if t218 {
		return rt.Flag(true), nil
	} else {
		return EstSostavlyayuschaya(ctx, put, rt.Text("Загрузки"))
	}
}

// EstPrimeta — функция flang «Есть примета».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Параметр mesto — «место»: «Разряд».
// Результат — значение.
func EstPrimeta(ctx *rt.Ctx, put rt.Value, mesto rt.Value) (rt.Value, error) {
	t220, e221 := EtoNeizvestnoe(ctx, mesto)
	if e221 != nil {
		return rt.Value{}, e221
	}
	t222, e223 := rt.Cond(ctx, t220)
	if e223 != nil {
		return rt.Value{}, e223
	}
	var t224 rt.Value
	if t222 {
		t224 = rt.Flag(false)
	} else {
		t224 = rt.Flag(true)
	}
	t225, e226 := rt.Cond(ctx, t224)
	if e226 != nil {
		return rt.Value{}, e226
	}
	var t227 rt.Value
	if t225 {
		t227 = rt.Flag(true)
	} else {
		t228, e229 := PrimetaKesha(ctx, put)
		if e229 != nil {
			return rt.Value{}, e229
		}
		t227 = t228
	}
	t230, e231 := rt.Cond(ctx, t227)
	if e231 != nil {
		return rt.Value{}, e231
	}
	var t232 rt.Value
	if t230 {
		t232 = rt.Flag(true)
	} else {
		t233, e234 := PrimetaZhurnala(ctx, put)
		if e234 != nil {
			return rt.Value{}, e234
		}
		t232 = t233
	}
	t235, e236 := rt.Cond(ctx, t232)
	if e236 != nil {
		return rt.Value{}, e236
	}
	var t237 rt.Value
	if t235 {
		t237 = rt.Flag(true)
	} else {
		t238, e239 := PrimetaSborki(ctx, put)
		if e239 != nil {
			return rt.Value{}, e239
		}
		t237 = t238
	}
	t240, e241 := rt.Cond(ctx, t237)
	if e241 != nil {
		return rt.Value{}, e241
	}
	if t240 {
		return rt.Flag(true), nil
	} else {
		return PrimetaZagruzki(ctx, put)
	}
}

// RazryadReshyonRazmerom — функция flang «Разряд решён размером».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func RazryadReshyonRazmerom(ctx *rt.Ctx, razryad rt.Value) (rt.Value, error) {
	if rt.VariantIs(razryad, "Крупное") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Неизвестное") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Кэш") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Журнал") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Сборка") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(razryad, "Загрузка") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// RazryadNahodki — функция flang «Разряд находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр mesto — «место»: «Разряд».
// Результат — значение: «Разряд».
func RazryadNahodki(ctx *rt.Ctx, nahodka rt.Value, mesto rt.Value) (rt.Value, error) {
	t242, e243 := EtoNeizvestnoe(ctx, mesto)
	if e243 != nil {
		return rt.Value{}, e243
	}
	t244, e245 := rt.Cond(ctx, t242)
	if e245 != nil {
		return rt.Value{}, e245
	}
	var t246 rt.Value
	if t244 {
		t246 = rt.Flag(false)
	} else {
		t246 = rt.Flag(true)
	}
	t247, e248 := rt.Cond(ctx, t246)
	if e248 != nil {
		return rt.Value{}, e248
	}
	var t249 rt.Value
	if t247 {
		t249 = mesto
	} else {
		t250, e251 := rt.FieldGet(ctx, nahodka, "путь")
		if e251 != nil {
			return rt.Value{}, e251
		}
		t252, e253 := PrimetaKesha(ctx, t250)
		if e253 != nil {
			return rt.Value{}, e253
		}
		t254, e255 := rt.Cond(ctx, t252)
		if e255 != nil {
			return rt.Value{}, e255
		}
		var t256 rt.Value
		if t254 {
			t256 = rt.Variant("Кэш", nil)
		} else {
			t257, e258 := rt.FieldGet(ctx, nahodka, "путь")
			if e258 != nil {
				return rt.Value{}, e258
			}
			t259, e260 := PrimetaZhurnala(ctx, t257)
			if e260 != nil {
				return rt.Value{}, e260
			}
			t261, e262 := rt.Cond(ctx, t259)
			if e262 != nil {
				return rt.Value{}, e262
			}
			var t263 rt.Value
			if t261 {
				t263 = rt.Variant("Журнал", nil)
			} else {
				t264, e265 := rt.FieldGet(ctx, nahodka, "путь")
				if e265 != nil {
					return rt.Value{}, e265
				}
				t266, e267 := PrimetaSborki(ctx, t264)
				if e267 != nil {
					return rt.Value{}, e267
				}
				t268, e269 := rt.Cond(ctx, t266)
				if e269 != nil {
					return rt.Value{}, e269
				}
				var t270 rt.Value
				if t268 {
					t270 = rt.Variant("Сборка", nil)
				} else {
					t271, e272 := rt.FieldGet(ctx, nahodka, "путь")
					if e272 != nil {
						return rt.Value{}, e272
					}
					t273, e274 := PrimetaZagruzki(ctx, t271)
					if e274 != nil {
						return rt.Value{}, e274
					}
					t275, e276 := rt.Cond(ctx, t273)
					if e276 != nil {
						return rt.Value{}, e276
					}
					var t277 rt.Value
					if t275 {
						t277 = rt.Variant("Загрузка", nil)
					} else {
						t278, e279 := rt.FieldGet(ctx, nahodka, "размер")
						if e279 != nil {
							return rt.Value{}, e279
						}
						t280, e281 := PorogKrupnogo(ctx)
						if e281 != nil {
							return rt.Value{}, e281
						}
						t282, e283 := rt.Gte(ctx, t278, t280)
						if e283 != nil {
							return rt.Value{}, e283
						}
						t284, e285 := rt.Cond(ctx, t282)
						if e285 != nil {
							return rt.Value{}, e285
						}
						var t286 rt.Value
						if t284 {
							t286 = rt.Variant("Крупное", nil)
						} else {
							t286 = rt.Variant("Неизвестное", nil)
						}
						t277 = t286
					}
					t270 = t277
				}
				t263 = t270
			}
			t256 = t263
		}
		t249 = t256
	}
	t287 := t249
	t288, e289 := RazryadObosnovan(ctx, nahodka, mesto, t287)
	if e289 != nil {
		return rt.Value{}, e289
	}
	// постусловие «Разряд обоснован приметой или местом»
	t290, e291 := rt.Post(ctx, t288, "Разряд обоснован приметой или местом", "Разряд находки")
	if e291 != nil {
		return rt.Value{}, e291
	}
	if !t290 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Разряд обоснован приметой или местом» функции «Разряд находки»")
	}
	t292, e293 := rt.FieldGet(ctx, nahodka, "путь")
	if e293 != nil {
		return rt.Value{}, e293
	}
	t294, e295 := EstPrimeta(ctx, t292, mesto)
	if e295 != nil {
		return rt.Value{}, e295
	}
	t296, e297 := rt.Cond(ctx, t294)
	if e297 != nil {
		return rt.Value{}, e297
	}
	var t298 rt.Value
	if t296 {
		t298 = rt.Flag(false)
	} else {
		t298 = rt.Flag(true)
	}
	t299, e300 := rt.Cond(ctx, t298)
	if e300 != nil {
		return rt.Value{}, e300
	}
	var t301 rt.Value
	if t299 {
		t302, e303 := RazryadReshyonRazmerom(ctx, t287)
		if e303 != nil {
			return rt.Value{}, e303
		}
		t301 = t302
	} else {
		t301 = rt.Flag(true)
	}
	// постусловие «И4: без приметы-составляющей и без места разряд решает размер»
	t304, e305 := rt.Post(ctx, t301, "И4: без приметы-составляющей и без места разряд решает размер", "Разряд находки")
	if e305 != nil {
		return rt.Value{}, e305
	}
	if !t304 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И4: без приметы-составляющей и без места разряд решает размер» функции «Разряд находки»")
	}
	return t287, nil
}

// KrupnoeNeMelchePoroga — функция flang «Крупное не мельче порога».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func KrupnoeNeMelchePoroga(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value) (rt.Value, error) {
	if rt.VariantIs(razryad, "Крупное") {
		t306, e307 := rt.FieldGet(ctx, nahodka, "размер")
		if e307 != nil {
			return rt.Value{}, e307
		}
		t308, e309 := PorogKrupnogo(ctx)
		if e309 != nil {
			return rt.Value{}, e309
		}
		t310, e311 := rt.Gte(ctx, t306, t308)
		if e311 != nil {
			return rt.Value{}, e311
		}
		return t310, nil
	} else if rt.VariantIs(razryad, "Кэш") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Журнал") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Сборка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Загрузка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(razryad, "Неизвестное") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// Katalog — функция flang «Каталог».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение.
func Katalog(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t312, e313 := rt.FieldGet(ctx, nahodka, "вид")
	if e313 != nil {
		return rt.Value{}, e313
	}
	if rt.VariantIs(t312, "Каталог") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t312, "Файл") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t312, "Ссылка") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t312)
	}
}

// Ssylka — функция flang «Ссылка».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение.
func Ssylka(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t314, e315 := rt.FieldGet(ctx, nahodka, "вид")
	if e315 != nil {
		return rt.Value{}, e315
	}
	if rt.VariantIs(t314, "Ссылка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t314, "Файл") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t314, "Каталог") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t314)
	}
}

// PrigovorMusora — функция flang «Приговор мусора».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр porog — «порог»: число.
// Результат — значение: «Приговор».
func PrigovorMusora(ctx *rt.Ctx, nahodka rt.Value, porog rt.Value) (rt.Value, error) {
	t316, e317 := Katalog(ctx, nahodka)
	if e317 != nil {
		return rt.Value{}, e317
	}
	t318, e319 := rt.Cond(ctx, t316)
	if e319 != nil {
		return rt.Value{}, e319
	}
	var t320 rt.Value
	if t318 {
		t320 = rt.Variant("Спросить", nil)
	} else {
		t321, e322 := rt.FieldGet(ctx, nahodka, "возраст_дней")
		if e322 != nil {
			return rt.Value{}, e322
		}
		t323, e324 := rt.Gte(ctx, t321, porog)
		if e324 != nil {
			return rt.Value{}, e324
		}
		t325, e326 := rt.Cond(ctx, t323)
		if e326 != nil {
			return rt.Value{}, e326
		}
		var t327 rt.Value
		if t325 {
			t327 = rt.Variant("МожноУбрать", nil)
		} else {
			t327 = rt.Variant("Спросить", nil)
		}
		t320 = t327
	}
	t328 := t320
	t329, e330 := Katalog(ctx, nahodka)
	if e330 != nil {
		return rt.Value{}, e330
	}
	t331, e332 := rt.Cond(ctx, t329)
	if e332 != nil {
		return rt.Value{}, e332
	}
	var t333 rt.Value
	if t331 {
		t334, e335 := EtoMozhnoUbrat(ctx, t328)
		if e335 != nil {
			return rt.Value{}, e335
		}
		t336, e337 := rt.Cond(ctx, t334)
		if e337 != nil {
			return rt.Value{}, e337
		}
		var t338 rt.Value
		if t336 {
			t338 = rt.Flag(false)
		} else {
			t338 = rt.Flag(true)
		}
		t333 = t338
	} else {
		t333 = rt.Flag(true)
	}
	// постусловие «Каталог не убирается»
	t339, e340 := rt.Post(ctx, t333, "Каталог не убирается", "Приговор мусора")
	if e340 != nil {
		return rt.Value{}, e340
	}
	if !t339 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Каталог не убирается» функции «Приговор мусора»")
	}
	return t328, nil
}

// PrigovorNahodki — функция flang «Приговор находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение: «Приговор».
func PrigovorNahodki(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value) (rt.Value, error) {
	t341, e342 := rt.FieldGet(ctx, nahodka, "доступен")
	if e342 != nil {
		return rt.Value{}, e342
	}
	t343, e344 := rt.Cond(ctx, t341)
	if e344 != nil {
		return rt.Value{}, e344
	}
	var t345 rt.Value
	if t343 {
		t345 = rt.Flag(false)
	} else {
		t345 = rt.Flag(true)
	}
	t346, e347 := rt.Cond(ctx, t345)
	if e347 != nil {
		return rt.Value{}, e347
	}
	var t348 rt.Value
	if t346 {
		t348 = rt.Variant("НеТрогать", nil)
	} else {
		t349, e350 := Ssylka(ctx, nahodka)
		if e350 != nil {
			return rt.Value{}, e350
		}
		t351, e352 := rt.Cond(ctx, t349)
		if e352 != nil {
			return rt.Value{}, e352
		}
		var t353 rt.Value
		if t351 {
			t353 = rt.Variant("НеТрогать", nil)
		} else {
			t354, e355 := rt.FieldGet(ctx, nahodka, "путь")
			if e355 != nil {
				return rt.Value{}, e355
			}
			t356, e357 := AdresuetsyaSoderzhimym(ctx, t354)
			if e357 != nil {
				return rt.Value{}, e357
			}
			t358, e359 := rt.Cond(ctx, t356)
			if e359 != nil {
				return rt.Value{}, e359
			}
			var t360 rt.Value
			if t358 {
				t360 = rt.Variant("НеТрогать", nil)
			} else {
				var t361 rt.Value
				if rt.VariantIs(razryad, "Кэш") {
					t362, e363 := PorogKesha(ctx)
					if e363 != nil {
						return rt.Value{}, e363
					}
					t364, e365 := PrigovorMusora(ctx, nahodka, t362)
					if e365 != nil {
						return rt.Value{}, e365
					}
					t361 = t364
				} else if rt.VariantIs(razryad, "Сборка") {
					t366, e367 := PorogKesha(ctx)
					if e367 != nil {
						return rt.Value{}, e367
					}
					t368, e369 := PrigovorMusora(ctx, nahodka, t366)
					if e369 != nil {
						return rt.Value{}, e369
					}
					t361 = t368
				} else if rt.VariantIs(razryad, "Журнал") {
					t370, e371 := PorogZhurnala(ctx)
					if e371 != nil {
						return rt.Value{}, e371
					}
					t372, e373 := PrigovorMusora(ctx, nahodka, t370)
					if e373 != nil {
						return rt.Value{}, e373
					}
					t361 = t372
				} else if rt.VariantIs(razryad, "Загрузка") {
					t374, e375 := rt.FieldGet(ctx, nahodka, "возраст_дней")
					if e375 != nil {
						return rt.Value{}, e375
					}
					t376, e377 := PorogZagruzki(ctx)
					if e377 != nil {
						return rt.Value{}, e377
					}
					t378, e379 := rt.Gte(ctx, t374, t376)
					if e379 != nil {
						return rt.Value{}, e379
					}
					t380, e381 := rt.Cond(ctx, t378)
					if e381 != nil {
						return rt.Value{}, e381
					}
					var t382 rt.Value
					if t380 {
						t382 = rt.Variant("Спросить", nil)
					} else {
						t382 = rt.Variant("НеТрогать", nil)
					}
					t361 = t382
				} else if rt.VariantIs(razryad, "Крупное") {
					t361 = rt.Variant("Спросить", nil)
				} else if rt.VariantIs(razryad, "Неизвестное") {
					t361 = rt.Variant("НеТрогать", nil)
				} else {
					return rt.Value{}, rt.MatchFail(ctx, razryad)
				}
				t360 = t361
			}
			t353 = t360
		}
		t348 = t353
	}
	t383 := t348
	t384, e385 := PrigovorObosnovan(ctx, nahodka, razryad, t383)
	if e385 != nil {
		return rt.Value{}, e385
	}
	// постусловие «Приговор обоснован»
	t386, e387 := rt.Post(ctx, t384, "Приговор обоснован", "Приговор находки")
	if e387 != nil {
		return rt.Value{}, e387
	}
	if !t386 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Приговор обоснован» функции «Приговор находки»")
	}
	t388, e389 := rt.FieldGet(ctx, nahodka, "путь")
	if e389 != nil {
		return rt.Value{}, e389
	}
	t390, e391 := AdresuetsyaSoderzhimym(ctx, t388)
	if e391 != nil {
		return rt.Value{}, e391
	}
	t392, e393 := rt.Cond(ctx, t390)
	if e393 != nil {
		return rt.Value{}, e393
	}
	var t394 rt.Value
	if t392 {
		t395, e396 := EtoNeTrogat(ctx, t383)
		if e396 != nil {
			return rt.Value{}, e396
		}
		t394 = t395
	} else {
		t394 = rt.Flag(true)
	}
	// постусловие «И3: адресуемое содержимым не убирается»
	t397, e398 := rt.Post(ctx, t394, "И3: адресуемое содержимым не убирается", "Приговор находки")
	if e398 != nil {
		return rt.Value{}, e398
	}
	if !t397 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И3: адресуемое содержимым не убирается» функции «Приговор находки»")
	}
	return t383, nil
}

// VesNahodki — функция flang «Вес находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение: число.
func VesNahodki(ctx *rt.Ctx, nahodka rt.Value, prigovor rt.Value) (rt.Value, error) {
	var t399 rt.Value
	if rt.VariantIs(prigovor, "НеТрогать") {
		t399 = rt.Number(0.0)
	} else if rt.VariantIs(prigovor, "МожноУбрать") {
		t400, e401 := rt.FieldGet(ctx, nahodka, "размер")
		if e401 != nil {
			return rt.Value{}, e401
		}
		t399 = t400
	} else if rt.VariantIs(prigovor, "Спросить") {
		t402, e403 := rt.FieldGet(ctx, nahodka, "размер")
		if e403 != nil {
			return rt.Value{}, e403
		}
		t399 = t402
	} else {
		return rt.Value{}, rt.MatchFail(ctx, prigovor)
	}
	t404 := t399
	t405, e406 := VesObosnovan(ctx, nahodka, prigovor, t404)
	if e406 != nil {
		return rt.Value{}, e406
	}
	// постусловие «Вес обоснован»
	t407, e408 := rt.Post(ctx, t405, "Вес обоснован", "Вес находки")
	if e408 != nil {
		return rt.Value{}, e408
	}
	if !t407 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Вес обоснован» функции «Вес находки»")
	}
	return t404, nil
}

// VesVGranicah — функция flang «Вес в границах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр ves — «вес»: число.
// Результат — значение.
func VesVGranicah(ctx *rt.Ctx, nahodka rt.Value, ves rt.Value) (rt.Value, error) {
	t409, e410 := rt.Gte(ctx, ves, rt.Number(0.0))
	if e410 != nil {
		return rt.Value{}, e410
	}
	t411, e412 := rt.Cond(ctx, t409)
	if e412 != nil {
		return rt.Value{}, e412
	}
	if t411 {
		t413, e414 := rt.FieldGet(ctx, nahodka, "размер")
		if e414 != nil {
			return rt.Value{}, e414
		}
		t415, e416 := rt.Lte(ctx, ves, t413)
		if e416 != nil {
			return rt.Value{}, e416
		}
		return t415, nil
	} else {
		return rt.Flag(false), nil
	}
}

// ReshitNahodku — функция flang «Решить находку».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: «Решение».
func ReshitNahodku(ctx *rt.Ctx, nahodka rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t417, e418 := rt.FieldGet(ctx, nahodka, "путь")
	if e418 != nil {
		return rt.Value{}, e418
	}
	t419, e420 := RazryadPoSpravochniku(ctx, t417, spravochnik)
	if e420 != nil {
		return rt.Value{}, e420
	}
	// пусть «место»
	mesto := t419
	t421, e422 := RazryadNahodki(ctx, nahodka, mesto)
	if e422 != nil {
		return rt.Value{}, e422
	}
	// пусть «разряд»
	razryad := t421
	t423, e424 := PrigovorNahodki(ctx, nahodka, razryad)
	if e424 != nil {
		return rt.Value{}, e424
	}
	// пусть «приговор»
	prigovor := t423
	t425, e426 := VesNahodki(ctx, nahodka, prigovor)
	if e426 != nil {
		return rt.Value{}, e426
	}
	t427 := make([]rt.Field, 3)
	t427[0] = rt.Field{Name: "разряд", Value: razryad}
	t427[1] = rt.Field{Name: "приговор", Value: prigovor}
	t427[2] = rt.Field{Name: "вес", Value: t425}
	t428 := rt.Record(t427)
	t429, e430 := I1Derzhitsya(ctx, t428)
	if e430 != nil {
		return rt.Value{}, e430
	}
	// постусловие «И1: убрать можно только мусор»
	t431, e432 := rt.Post(ctx, t429, "И1: убрать можно только мусор", "Решить находку")
	if e432 != nil {
		return rt.Value{}, e432
	}
	if !t431 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И1: убрать можно только мусор» функции «Решить находку»")
	}
	return t428, nil
}

// ReshitVsyo — функция flang «Решить всё».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: список: «Решение».
func ReshitVsyo(ctx *rt.Ctx, zapisi rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t433, e434 := rt.RequireList(ctx, zapisi, "отобразить")
	if e434 != nil {
		return rt.Value{}, e434
	}
	t435 := make([]rt.Value, 0, len(t433))
	for t436 := range t433 {
		// «находка»
		nahodka := t433[t436]
		t437, e438 := ReshitNahodku(ctx, nahodka, spravochnik)
		if e438 != nil {
			return rt.Value{}, e438
		}
		t435 = append(t435, t437)
	}
	t439 := rt.List(t435)
	t440, e441 := I1DerzhitsyaVsyudu(ctx, t439)
	if e441 != nil {
		return rt.Value{}, e441
	}
	// постусловие «И1 всюду: убрать можно только мусор»
	t442, e443 := rt.Post(ctx, t440, "И1 всюду: убрать можно только мусор", "Решить всё")
	if e443 != nil {
		return rt.Value{}, e443
	}
	if !t442 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И1 всюду: убрать можно только мусор» функции «Решить всё»")
	}
	return t439, nil
}

// I1Derzhitsya — функция flang «И1 держится».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр reshenie — «решение»: «Решение».
// Результат — значение.
func I1Derzhitsya(ctx *rt.Ctx, reshenie rt.Value) (rt.Value, error) {
	t444, e445 := rt.FieldGet(ctx, reshenie, "приговор")
	if e445 != nil {
		return rt.Value{}, e445
	}
	if rt.VariantIs(t444, "МожноУбрать") {
		t446, e447 := rt.FieldGet(ctx, reshenie, "разряд")
		if e447 != nil {
			return rt.Value{}, e447
		}
		if rt.VariantIs(t446, "Кэш") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t446, "Журнал") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t446, "Сборка") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t446, "Загрузка") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t446, "Крупное") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t446, "Неизвестное") {
			return rt.Flag(false), nil
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t446)
		}
	} else if rt.VariantIs(t444, "Спросить") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t444, "НеТрогать") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t444)
	}
}

// I1DerzhitsyaVsyudu — функция flang «И1 держится всюду».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр resheniya — «решения»: список: «Решение».
// Результат — значение.
func I1DerzhitsyaVsyudu(ctx *rt.Ctx, resheniya rt.Value) (rt.Value, error) {
	t448, e449 := rt.RequireList(ctx, resheniya, "свёртка")
	if e449 != nil {
		return rt.Value{}, e449
	}
	// «акк»
	akk := rt.Flag(true)
	for t450 := range t448 {
		// «решение»
		reshenie := t448[t450]
		t451, e452 := rt.Cond(ctx, akk)
		if e452 != nil {
			return rt.Value{}, e452
		}
		var t453 rt.Value
		if t451 {
			t454, e455 := I1Derzhitsya(ctx, reshenie)
			if e455 != nil {
				return rt.Value{}, e455
			}
			t453 = t454
		} else {
			t453 = rt.Flag(false)
		}
		akk = t453
	}
	return akk, nil
}

// PustoySvod — функция flang «Пустой свод».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: «Свод».
func PustoySvod(ctx *rt.Ctx) (rt.Value, error) {
	t456 := make([]rt.Field, 7)
	t456[0] = rt.Field{Name: "кэш", Value: rt.Number(0.0)}
	t456[1] = rt.Field{Name: "журнал", Value: rt.Number(0.0)}
	t456[2] = rt.Field{Name: "сборка", Value: rt.Number(0.0)}
	t456[3] = rt.Field{Name: "загрузка", Value: rt.Number(0.0)}
	t456[4] = rt.Field{Name: "крупное", Value: rt.Number(0.0)}
	t456[5] = rt.Field{Name: "неизвестное", Value: rt.Number(0.0)}
	t456[6] = rt.Field{Name: "освободить", Value: rt.Number(0.0)}
	return rt.Record(t456), nil
}

// PribavitReshenie — функция flang «Прибавить решение».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр svod — «свод»: «Свод».
// Параметр reshenie — «решение»: «Решение».
// Результат — значение: «Свод».
func PribavitReshenie(ctx *rt.Ctx, svod rt.Value, reshenie rt.Value) (rt.Value, error) {
	var t457 rt.Value
	t458, e459 := rt.FieldGet(ctx, reshenie, "приговор")
	if e459 != nil {
		return rt.Value{}, e459
	}
	if rt.VariantIs(t458, "МожноУбрать") {
		t460, e461 := rt.FieldGet(ctx, reshenie, "вес")
		if e461 != nil {
			return rt.Value{}, e461
		}
		t457 = t460
	} else if rt.VariantIs(t458, "Спросить") {
		t457 = rt.Number(0.0)
	} else if rt.VariantIs(t458, "НеТрогать") {
		t457 = rt.Number(0.0)
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t458)
	}
	// пусть «убрать»
	ubrat := t457
	t462, e463 := rt.FieldGet(ctx, reshenie, "разряд")
	if e463 != nil {
		return rt.Value{}, e463
	}
	if rt.VariantIs(t462, "Кэш") {
		t464, e465 := rt.FieldGet(ctx, svod, "кэш")
		if e465 != nil {
			return rt.Value{}, e465
		}
		t466, e467 := rt.FieldGet(ctx, reshenie, "вес")
		if e467 != nil {
			return rt.Value{}, e467
		}
		t468, e469 := rt.Add(ctx, t464, t466)
		if e469 != nil {
			return rt.Value{}, e469
		}
		t470, e471 := rt.FieldGet(ctx, svod, "журнал")
		if e471 != nil {
			return rt.Value{}, e471
		}
		t472, e473 := rt.FieldGet(ctx, svod, "сборка")
		if e473 != nil {
			return rt.Value{}, e473
		}
		t474, e475 := rt.FieldGet(ctx, svod, "загрузка")
		if e475 != nil {
			return rt.Value{}, e475
		}
		t476, e477 := rt.FieldGet(ctx, svod, "крупное")
		if e477 != nil {
			return rt.Value{}, e477
		}
		t478, e479 := rt.FieldGet(ctx, svod, "неизвестное")
		if e479 != nil {
			return rt.Value{}, e479
		}
		t480, e481 := rt.FieldGet(ctx, svod, "освободить")
		if e481 != nil {
			return rt.Value{}, e481
		}
		t482, e483 := rt.Add(ctx, t480, ubrat)
		if e483 != nil {
			return rt.Value{}, e483
		}
		t484 := make([]rt.Field, 7)
		t484[0] = rt.Field{Name: "кэш", Value: t468}
		t484[1] = rt.Field{Name: "журнал", Value: t470}
		t484[2] = rt.Field{Name: "сборка", Value: t472}
		t484[3] = rt.Field{Name: "загрузка", Value: t474}
		t484[4] = rt.Field{Name: "крупное", Value: t476}
		t484[5] = rt.Field{Name: "неизвестное", Value: t478}
		t484[6] = rt.Field{Name: "освободить", Value: t482}
		return rt.Record(t484), nil
	} else if rt.VariantIs(t462, "Журнал") {
		t485, e486 := rt.FieldGet(ctx, svod, "кэш")
		if e486 != nil {
			return rt.Value{}, e486
		}
		t487, e488 := rt.FieldGet(ctx, svod, "журнал")
		if e488 != nil {
			return rt.Value{}, e488
		}
		t489, e490 := rt.FieldGet(ctx, reshenie, "вес")
		if e490 != nil {
			return rt.Value{}, e490
		}
		t491, e492 := rt.Add(ctx, t487, t489)
		if e492 != nil {
			return rt.Value{}, e492
		}
		t493, e494 := rt.FieldGet(ctx, svod, "сборка")
		if e494 != nil {
			return rt.Value{}, e494
		}
		t495, e496 := rt.FieldGet(ctx, svod, "загрузка")
		if e496 != nil {
			return rt.Value{}, e496
		}
		t497, e498 := rt.FieldGet(ctx, svod, "крупное")
		if e498 != nil {
			return rt.Value{}, e498
		}
		t499, e500 := rt.FieldGet(ctx, svod, "неизвестное")
		if e500 != nil {
			return rt.Value{}, e500
		}
		t501, e502 := rt.FieldGet(ctx, svod, "освободить")
		if e502 != nil {
			return rt.Value{}, e502
		}
		t503, e504 := rt.Add(ctx, t501, ubrat)
		if e504 != nil {
			return rt.Value{}, e504
		}
		t505 := make([]rt.Field, 7)
		t505[0] = rt.Field{Name: "кэш", Value: t485}
		t505[1] = rt.Field{Name: "журнал", Value: t491}
		t505[2] = rt.Field{Name: "сборка", Value: t493}
		t505[3] = rt.Field{Name: "загрузка", Value: t495}
		t505[4] = rt.Field{Name: "крупное", Value: t497}
		t505[5] = rt.Field{Name: "неизвестное", Value: t499}
		t505[6] = rt.Field{Name: "освободить", Value: t503}
		return rt.Record(t505), nil
	} else if rt.VariantIs(t462, "Сборка") {
		t506, e507 := rt.FieldGet(ctx, svod, "кэш")
		if e507 != nil {
			return rt.Value{}, e507
		}
		t508, e509 := rt.FieldGet(ctx, svod, "журнал")
		if e509 != nil {
			return rt.Value{}, e509
		}
		t510, e511 := rt.FieldGet(ctx, svod, "сборка")
		if e511 != nil {
			return rt.Value{}, e511
		}
		t512, e513 := rt.FieldGet(ctx, reshenie, "вес")
		if e513 != nil {
			return rt.Value{}, e513
		}
		t514, e515 := rt.Add(ctx, t510, t512)
		if e515 != nil {
			return rt.Value{}, e515
		}
		t516, e517 := rt.FieldGet(ctx, svod, "загрузка")
		if e517 != nil {
			return rt.Value{}, e517
		}
		t518, e519 := rt.FieldGet(ctx, svod, "крупное")
		if e519 != nil {
			return rt.Value{}, e519
		}
		t520, e521 := rt.FieldGet(ctx, svod, "неизвестное")
		if e521 != nil {
			return rt.Value{}, e521
		}
		t522, e523 := rt.FieldGet(ctx, svod, "освободить")
		if e523 != nil {
			return rt.Value{}, e523
		}
		t524, e525 := rt.Add(ctx, t522, ubrat)
		if e525 != nil {
			return rt.Value{}, e525
		}
		t526 := make([]rt.Field, 7)
		t526[0] = rt.Field{Name: "кэш", Value: t506}
		t526[1] = rt.Field{Name: "журнал", Value: t508}
		t526[2] = rt.Field{Name: "сборка", Value: t514}
		t526[3] = rt.Field{Name: "загрузка", Value: t516}
		t526[4] = rt.Field{Name: "крупное", Value: t518}
		t526[5] = rt.Field{Name: "неизвестное", Value: t520}
		t526[6] = rt.Field{Name: "освободить", Value: t524}
		return rt.Record(t526), nil
	} else if rt.VariantIs(t462, "Загрузка") {
		t527, e528 := rt.FieldGet(ctx, svod, "кэш")
		if e528 != nil {
			return rt.Value{}, e528
		}
		t529, e530 := rt.FieldGet(ctx, svod, "журнал")
		if e530 != nil {
			return rt.Value{}, e530
		}
		t531, e532 := rt.FieldGet(ctx, svod, "сборка")
		if e532 != nil {
			return rt.Value{}, e532
		}
		t533, e534 := rt.FieldGet(ctx, svod, "загрузка")
		if e534 != nil {
			return rt.Value{}, e534
		}
		t535, e536 := rt.FieldGet(ctx, reshenie, "вес")
		if e536 != nil {
			return rt.Value{}, e536
		}
		t537, e538 := rt.Add(ctx, t533, t535)
		if e538 != nil {
			return rt.Value{}, e538
		}
		t539, e540 := rt.FieldGet(ctx, svod, "крупное")
		if e540 != nil {
			return rt.Value{}, e540
		}
		t541, e542 := rt.FieldGet(ctx, svod, "неизвестное")
		if e542 != nil {
			return rt.Value{}, e542
		}
		t543, e544 := rt.FieldGet(ctx, svod, "освободить")
		if e544 != nil {
			return rt.Value{}, e544
		}
		t545, e546 := rt.Add(ctx, t543, ubrat)
		if e546 != nil {
			return rt.Value{}, e546
		}
		t547 := make([]rt.Field, 7)
		t547[0] = rt.Field{Name: "кэш", Value: t527}
		t547[1] = rt.Field{Name: "журнал", Value: t529}
		t547[2] = rt.Field{Name: "сборка", Value: t531}
		t547[3] = rt.Field{Name: "загрузка", Value: t537}
		t547[4] = rt.Field{Name: "крупное", Value: t539}
		t547[5] = rt.Field{Name: "неизвестное", Value: t541}
		t547[6] = rt.Field{Name: "освободить", Value: t545}
		return rt.Record(t547), nil
	} else if rt.VariantIs(t462, "Крупное") {
		t548, e549 := rt.FieldGet(ctx, svod, "кэш")
		if e549 != nil {
			return rt.Value{}, e549
		}
		t550, e551 := rt.FieldGet(ctx, svod, "журнал")
		if e551 != nil {
			return rt.Value{}, e551
		}
		t552, e553 := rt.FieldGet(ctx, svod, "сборка")
		if e553 != nil {
			return rt.Value{}, e553
		}
		t554, e555 := rt.FieldGet(ctx, svod, "загрузка")
		if e555 != nil {
			return rt.Value{}, e555
		}
		t556, e557 := rt.FieldGet(ctx, svod, "крупное")
		if e557 != nil {
			return rt.Value{}, e557
		}
		t558, e559 := rt.FieldGet(ctx, reshenie, "вес")
		if e559 != nil {
			return rt.Value{}, e559
		}
		t560, e561 := rt.Add(ctx, t556, t558)
		if e561 != nil {
			return rt.Value{}, e561
		}
		t562, e563 := rt.FieldGet(ctx, svod, "неизвестное")
		if e563 != nil {
			return rt.Value{}, e563
		}
		t564, e565 := rt.FieldGet(ctx, svod, "освободить")
		if e565 != nil {
			return rt.Value{}, e565
		}
		t566, e567 := rt.Add(ctx, t564, ubrat)
		if e567 != nil {
			return rt.Value{}, e567
		}
		t568 := make([]rt.Field, 7)
		t568[0] = rt.Field{Name: "кэш", Value: t548}
		t568[1] = rt.Field{Name: "журнал", Value: t550}
		t568[2] = rt.Field{Name: "сборка", Value: t552}
		t568[3] = rt.Field{Name: "загрузка", Value: t554}
		t568[4] = rt.Field{Name: "крупное", Value: t560}
		t568[5] = rt.Field{Name: "неизвестное", Value: t562}
		t568[6] = rt.Field{Name: "освободить", Value: t566}
		return rt.Record(t568), nil
	} else if rt.VariantIs(t462, "Неизвестное") {
		t569, e570 := rt.FieldGet(ctx, svod, "кэш")
		if e570 != nil {
			return rt.Value{}, e570
		}
		t571, e572 := rt.FieldGet(ctx, svod, "журнал")
		if e572 != nil {
			return rt.Value{}, e572
		}
		t573, e574 := rt.FieldGet(ctx, svod, "сборка")
		if e574 != nil {
			return rt.Value{}, e574
		}
		t575, e576 := rt.FieldGet(ctx, svod, "загрузка")
		if e576 != nil {
			return rt.Value{}, e576
		}
		t577, e578 := rt.FieldGet(ctx, svod, "крупное")
		if e578 != nil {
			return rt.Value{}, e578
		}
		t579, e580 := rt.FieldGet(ctx, svod, "неизвестное")
		if e580 != nil {
			return rt.Value{}, e580
		}
		t581, e582 := rt.FieldGet(ctx, reshenie, "вес")
		if e582 != nil {
			return rt.Value{}, e582
		}
		t583, e584 := rt.Add(ctx, t579, t581)
		if e584 != nil {
			return rt.Value{}, e584
		}
		t585, e586 := rt.FieldGet(ctx, svod, "освободить")
		if e586 != nil {
			return rt.Value{}, e586
		}
		t587, e588 := rt.Add(ctx, t585, ubrat)
		if e588 != nil {
			return rt.Value{}, e588
		}
		t589 := make([]rt.Field, 7)
		t589[0] = rt.Field{Name: "кэш", Value: t569}
		t589[1] = rt.Field{Name: "журнал", Value: t571}
		t589[2] = rt.Field{Name: "сборка", Value: t573}
		t589[3] = rt.Field{Name: "загрузка", Value: t575}
		t589[4] = rt.Field{Name: "крупное", Value: t577}
		t589[5] = rt.Field{Name: "неизвестное", Value: t583}
		t589[6] = rt.Field{Name: "освободить", Value: t587}
		return rt.Record(t589), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t462)
	}
}

// Svesti — функция flang «Свести».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: «Свод».
func Svesti(ctx *rt.Ctx, zapisi rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t590, e591 := ReshitVsyo(ctx, zapisi, spravochnik)
	if e591 != nil {
		return rt.Value{}, e591
	}
	t592, e593 := rt.RequireList(ctx, t590, "свёртка")
	if e593 != nil {
		return rt.Value{}, e593
	}
	t594, e595 := PustoySvod(ctx)
	if e595 != nil {
		return rt.Value{}, e595
	}
	// «свод»
	svod := t594
	for t596 := range t592 {
		// «решение»
		reshenie := t592[t596]
		t597, e598 := PribavitReshenie(ctx, svod, reshenie)
		if e598 != nil {
			return rt.Value{}, e598
		}
		svod = t597
	}
	t599 := svod
	t600, e601 := I2Derzhitsya(ctx, zapisi, spravochnik, t599)
	if e601 != nil {
		return rt.Value{}, e601
	}
	// постусловие «И2: освобождаемое не больше убираемого»
	t602, e603 := rt.Post(ctx, t600, "И2: освобождаемое не больше убираемого", "Свести")
	if e603 != nil {
		return rt.Value{}, e603
	}
	if !t602 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И2: освобождаемое не больше убираемого» функции «Свести»")
	}
	return t599, nil
}

// SummaRazmerovUbiraemyh — функция flang «Сумма размеров убираемых».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: число.
func SummaRazmerovUbiraemyh(ctx *rt.Ctx, zapisi rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t604, e605 := rt.RequireList(ctx, zapisi, "свёртка")
	if e605 != nil {
		return rt.Value{}, e605
	}
	// «акк»
	akk := rt.Number(0.0)
	for t606 := range t604 {
		// «находка»
		nahodka := t604[t606]
		var t607 rt.Value
		t608, e609 := ReshitNahodku(ctx, nahodka, spravochnik)
		if e609 != nil {
			return rt.Value{}, e609
		}
		t610, e611 := rt.FieldGet(ctx, t608, "приговор")
		if e611 != nil {
			return rt.Value{}, e611
		}
		if rt.VariantIs(t610, "МожноУбрать") {
			t612, e613 := rt.FieldGet(ctx, nahodka, "размер")
			if e613 != nil {
				return rt.Value{}, e613
			}
			t614, e615 := rt.Add(ctx, akk, t612)
			if e615 != nil {
				return rt.Value{}, e615
			}
			t607 = t614
		} else if rt.VariantIs(t610, "Спросить") {
			t607 = akk
		} else if rt.VariantIs(t610, "НеТрогать") {
			t607 = akk
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t610)
		}
		akk = t607
	}
	return akk, nil
}

// I2Derzhitsya — функция flang «И2 держится».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Параметр svod — «свод»: «Свод».
// Результат — значение.
func I2Derzhitsya(ctx *rt.Ctx, zapisi rt.Value, spravochnik rt.Value, svod rt.Value) (rt.Value, error) {
	t616, e617 := rt.FieldGet(ctx, svod, "освободить")
	if e617 != nil {
		return rt.Value{}, e617
	}
	t618, e619 := SummaRazmerovUbiraemyh(ctx, zapisi, spravochnik)
	if e619 != nil {
		return rt.Value{}, e619
	}
	t620, e621 := rt.Lte(ctx, t616, t618)
	if e621 != nil {
		return rt.Value{}, e621
	}
	return t620, nil
}

// StrokuOtchyota — функция flang «Строку отчёта».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: «Строка отчёта».
func StrokuOtchyota(ctx *rt.Ctx, nahodka rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t622, e623 := rt.FieldGet(ctx, nahodka, "путь")
	if e623 != nil {
		return rt.Value{}, e623
	}
	t624, e625 := ReshitNahodku(ctx, nahodka, spravochnik)
	if e625 != nil {
		return rt.Value{}, e625
	}
	t626 := make([]rt.Field, 2)
	t626[0] = rt.Field{Name: "путь", Value: t622}
	t626[1] = rt.Field{Name: "решение", Value: t624}
	return rt.Record(t626), nil
}

// VstavitPoVesu — функция flang «Вставить по весу».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр stroka — «строка»: «Строка отчёта».
// Параметр stroki — «строки»: список: «Строка отчёта».
// Результат — значение: список: «Строка отчёта».
func VstavitPoVesu(ctx *rt.Ctx, stroka rt.Value, stroki rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Вставить по весу"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	if rt.ChainEmpty(stroki) {
		t627 := make([]rt.Value, 1)
		t627[0] = stroka
		return rt.List(t627), nil
	} else if rt.ChainCons(stroki) {
		// голова «голова»
		golova := rt.ChainHead(stroki)
		// хвост «хвост»
		hvost := rt.ChainTail(stroki)
		t628, e629 := rt.FieldGet(ctx, stroka, "решение")
		if e629 != nil {
			return rt.Value{}, e629
		}
		t630, e631 := rt.FieldGet(ctx, t628, "вес")
		if e631 != nil {
			return rt.Value{}, e631
		}
		t632, e633 := rt.FieldGet(ctx, golova, "решение")
		if e633 != nil {
			return rt.Value{}, e633
		}
		t634, e635 := rt.FieldGet(ctx, t632, "вес")
		if e635 != nil {
			return rt.Value{}, e635
		}
		t636, e637 := rt.Gte(ctx, t630, t634)
		if e637 != nil {
			return rt.Value{}, e637
		}
		t638, e639 := rt.Cond(ctx, t636)
		if e639 != nil {
			return rt.Value{}, e639
		}
		if t638 {
			return PripisatStrokuOtchyota(ctx, stroka, stroki)
		} else {
			t640, e641 := VstavitPoVesu(ctx, stroka, hvost)
			if e641 != nil {
				return rt.Value{}, e641
			}
			return PripisatStrokuOtchyota(ctx, golova, t640)
		}
	} else {
		return rt.Value{}, rt.MatchFail(ctx, stroki)
	}
}

// PripisatStrokuOtchyota — функция flang «Приписать строку отчёта».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр pervaya — «первая»: «Строка отчёта».
// Параметр stroki — «строки»: список: «Строка отчёта».
// Результат — значение: список: «Строка отчёта».
func PripisatStrokuOtchyota(ctx *rt.Ctx, pervaya rt.Value, stroki rt.Value) (rt.Value, error) {
	t642, e643 := rt.RequireList(ctx, stroki, "свёртка")
	if e643 != nil {
		return rt.Value{}, e643
	}
	t644 := make([]rt.Value, 1)
	t644[0] = pervaya
	// «акк»
	akk := rt.List(t644)
	for t645 := range t642 {
		// «эл»
		el := t642[t645]
		// «добавить»
		t646, e647 := rt.BAppend(ctx, el, akk)
		if e647 != nil {
			return rt.Value{}, e647
		}
		akk = t646
	}
	return akk, nil
}

// Otchyot — функция flang «Отчёт».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: список: «Строка отчёта».
func Otchyot(ctx *rt.Ctx, zapisi rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t648, e649 := rt.RequireList(ctx, zapisi, "свёртка")
	if e649 != nil {
		return rt.Value{}, e649
	}
	// «акк»
	akk := rt.List(nil)
	for t650 := range t648 {
		// «находка»
		nahodka := t648[t650]
		t651, e652 := StrokuOtchyota(ctx, nahodka, spravochnik)
		if e652 != nil {
			return rt.Value{}, e652
		}
		t653, e654 := VstavitPoVesu(ctx, t651, akk)
		if e654 != nil {
			return rt.Value{}, e654
		}
		akk = t653
	}
	t655 := akk
	t656, e657 := OtchyotToyZheDliny(ctx, zapisi, t655)
	if e657 != nil {
		return rt.Value{}, e657
	}
	// постусловие «Отчёт той же длины»
	t658, e659 := rt.Post(ctx, t656, "Отчёт той же длины", "Отчёт")
	if e659 != nil {
		return rt.Value{}, e659
	}
	if !t658 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Отчёт той же длины» функции «Отчёт»")
	}
	return t655, nil
}

// OtchyotToyZheDliny — функция flang «Отчёт той же длины».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр stroki — «строки»: список: «Строка отчёта».
// Результат — значение.
func OtchyotToyZheDliny(ctx *rt.Ctx, zapisi rt.Value, stroki rt.Value) (rt.Value, error) {
	// «длина»
	t660, e661 := rt.BLength(ctx, stroki)
	if e661 != nil {
		return rt.Value{}, e661
	}
	// «длина»
	t662, e663 := rt.BLength(ctx, zapisi)
	if e663 != nil {
		return rt.Value{}, e663
	}
	return rt.Flag(rt.Equal(t660, t662)), nil
}

// EtoMozhnoUbrat — функция flang «Это МожноУбрать».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение.
func EtoMozhnoUbrat(ctx *rt.Ctx, prigovor rt.Value) (rt.Value, error) {
	if rt.VariantIs(prigovor, "МожноУбрать") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(prigovor, "Спросить") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(prigovor, "НеТрогать") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, prigovor)
	}
}

// EtoNeTrogat — функция flang «Это НеТрогать».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение.
func EtoNeTrogat(ctx *rt.Ctx, prigovor rt.Value) (rt.Value, error) {
	if rt.VariantIs(prigovor, "НеТрогать") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(prigovor, "Спросить") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(prigovor, "МожноУбрать") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, prigovor)
	}
}

// I1NaPare — функция flang «И1 на паре».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение.
func I1NaPare(ctx *rt.Ctx, razryad rt.Value, prigovor rt.Value) (rt.Value, error) {
	t664, e665 := EtoMozhnoUbrat(ctx, prigovor)
	if e665 != nil {
		return rt.Value{}, e665
	}
	t666, e667 := rt.Cond(ctx, t664)
	if e667 != nil {
		return rt.Value{}, e667
	}
	var t668 rt.Value
	if t666 {
		t668 = rt.Flag(false)
	} else {
		t668 = rt.Flag(true)
	}
	t669, e670 := rt.Cond(ctx, t668)
	if e670 != nil {
		return rt.Value{}, e670
	}
	if t669 {
		return rt.Flag(true), nil
	} else {
		if rt.VariantIs(razryad, "Кэш") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(razryad, "Журнал") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(razryad, "Сборка") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(razryad, "Загрузка") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(razryad, "Крупное") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(razryad, "Неизвестное") {
			return rt.Flag(false), nil
		} else {
			return rt.Value{}, rt.MatchFail(ctx, razryad)
		}
	}
}

// PorogRazryada — функция flang «Порог разряда».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение: число.
func PorogRazryada(ctx *rt.Ctx, razryad rt.Value) (rt.Value, error) {
	if rt.VariantIs(razryad, "Кэш") {
		return PorogKesha(ctx)
	} else if rt.VariantIs(razryad, "Сборка") {
		return PorogKesha(ctx)
	} else if rt.VariantIs(razryad, "Журнал") {
		return PorogZhurnala(ctx)
	} else if rt.VariantIs(razryad, "Загрузка") {
		return PorogZagruzki(ctx)
	} else if rt.VariantIs(razryad, "Крупное") {
		return rt.Number(0.0), nil
	} else if rt.VariantIs(razryad, "Неизвестное") {
		return rt.Number(0.0), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// RazryadObosnovan — функция flang «Разряд обоснован».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр mesto — «место»: «Разряд».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func RazryadObosnovan(ctx *rt.Ctx, nahodka rt.Value, mesto rt.Value, razryad rt.Value) (rt.Value, error) {
	t671, e672 := EtoNeizvestnoe(ctx, mesto)
	if e672 != nil {
		return rt.Value{}, e672
	}
	t673, e674 := rt.Cond(ctx, t671)
	if e674 != nil {
		return rt.Value{}, e674
	}
	var t675 rt.Value
	if t673 {
		t675 = rt.Flag(false)
	} else {
		t675 = rt.Flag(true)
	}
	t676, e677 := rt.Cond(ctx, t675)
	if e677 != nil {
		return rt.Value{}, e677
	}
	if t676 {
		return TotZheRazryad(ctx, razryad, mesto)
	} else {
		return RazryadObosnovanPrimetoy(ctx, nahodka, razryad)
	}
}

// RazryadObosnovanPrimetoy — функция flang «Разряд обоснован приметой».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func RazryadObosnovanPrimetoy(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value) (rt.Value, error) {
	t678, e679 := rt.FieldGet(ctx, nahodka, "путь")
	if e679 != nil {
		return rt.Value{}, e679
	}
	t680, e681 := PrimetaKesha(ctx, t678)
	if e681 != nil {
		return rt.Value{}, e681
	}
	// пусть «кэш»
	kesh := t680
	t682, e683 := rt.FieldGet(ctx, nahodka, "путь")
	if e683 != nil {
		return rt.Value{}, e683
	}
	t684, e685 := PrimetaZhurnala(ctx, t682)
	if e685 != nil {
		return rt.Value{}, e685
	}
	// пусть «журнал»
	zhurnal := t684
	t686, e687 := rt.FieldGet(ctx, nahodka, "путь")
	if e687 != nil {
		return rt.Value{}, e687
	}
	t688, e689 := PrimetaSborki(ctx, t686)
	if e689 != nil {
		return rt.Value{}, e689
	}
	// пусть «сборка»
	sborka := t688
	t690, e691 := rt.FieldGet(ctx, nahodka, "путь")
	if e691 != nil {
		return rt.Value{}, e691
	}
	t692, e693 := PrimetaZagruzki(ctx, t690)
	if e693 != nil {
		return rt.Value{}, e693
	}
	// пусть «загрузка»
	zagruzka := t692
	t694, e695 := rt.FieldGet(ctx, nahodka, "размер")
	if e695 != nil {
		return rt.Value{}, e695
	}
	t696, e697 := PorogKrupnogo(ctx)
	if e697 != nil {
		return rt.Value{}, e697
	}
	t698, e699 := rt.Gte(ctx, t694, t696)
	if e699 != nil {
		return rt.Value{}, e699
	}
	// пусть «крупное»
	krupnoe := t698
	if rt.VariantIs(razryad, "Кэш") {
		return kesh, nil
	} else if rt.VariantIs(razryad, "Журнал") {
		t700, e701 := rt.Cond(ctx, zhurnal)
		if e701 != nil {
			return rt.Value{}, e701
		}
		if t700 {
			t702, e703 := rt.Cond(ctx, kesh)
			if e703 != nil {
				return rt.Value{}, e703
			}
			if t702 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Сборка") {
		t704, e705 := rt.Cond(ctx, sborka)
		if e705 != nil {
			return rt.Value{}, e705
		}
		var t706 rt.Value
		if t704 {
			t707, e708 := rt.Cond(ctx, kesh)
			if e708 != nil {
				return rt.Value{}, e708
			}
			var t709 rt.Value
			if t707 {
				t709 = rt.Flag(false)
			} else {
				t709 = rt.Flag(true)
			}
			t706 = t709
		} else {
			t706 = rt.Flag(false)
		}
		t710, e711 := rt.Cond(ctx, t706)
		if e711 != nil {
			return rt.Value{}, e711
		}
		if t710 {
			t712, e713 := rt.Cond(ctx, zhurnal)
			if e713 != nil {
				return rt.Value{}, e713
			}
			if t712 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Загрузка") {
		t714, e715 := rt.Cond(ctx, zagruzka)
		if e715 != nil {
			return rt.Value{}, e715
		}
		var t716 rt.Value
		if t714 {
			t717, e718 := rt.Cond(ctx, kesh)
			if e718 != nil {
				return rt.Value{}, e718
			}
			var t719 rt.Value
			if t717 {
				t719 = rt.Flag(false)
			} else {
				t719 = rt.Flag(true)
			}
			t716 = t719
		} else {
			t716 = rt.Flag(false)
		}
		t720, e721 := rt.Cond(ctx, t716)
		if e721 != nil {
			return rt.Value{}, e721
		}
		var t722 rt.Value
		if t720 {
			t723, e724 := rt.Cond(ctx, zhurnal)
			if e724 != nil {
				return rt.Value{}, e724
			}
			var t725 rt.Value
			if t723 {
				t725 = rt.Flag(false)
			} else {
				t725 = rt.Flag(true)
			}
			t722 = t725
		} else {
			t722 = rt.Flag(false)
		}
		t726, e727 := rt.Cond(ctx, t722)
		if e727 != nil {
			return rt.Value{}, e727
		}
		if t726 {
			t728, e729 := rt.Cond(ctx, sborka)
			if e729 != nil {
				return rt.Value{}, e729
			}
			if t728 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Крупное") {
		t730, e731 := rt.Cond(ctx, krupnoe)
		if e731 != nil {
			return rt.Value{}, e731
		}
		var t732 rt.Value
		if t730 {
			t733, e734 := rt.Cond(ctx, kesh)
			if e734 != nil {
				return rt.Value{}, e734
			}
			var t735 rt.Value
			if t733 {
				t735 = rt.Flag(false)
			} else {
				t735 = rt.Flag(true)
			}
			t732 = t735
		} else {
			t732 = rt.Flag(false)
		}
		t736, e737 := rt.Cond(ctx, t732)
		if e737 != nil {
			return rt.Value{}, e737
		}
		var t738 rt.Value
		if t736 {
			t739, e740 := rt.Cond(ctx, zhurnal)
			if e740 != nil {
				return rt.Value{}, e740
			}
			var t741 rt.Value
			if t739 {
				t741 = rt.Flag(false)
			} else {
				t741 = rt.Flag(true)
			}
			t738 = t741
		} else {
			t738 = rt.Flag(false)
		}
		t742, e743 := rt.Cond(ctx, t738)
		if e743 != nil {
			return rt.Value{}, e743
		}
		var t744 rt.Value
		if t742 {
			t745, e746 := rt.Cond(ctx, sborka)
			if e746 != nil {
				return rt.Value{}, e746
			}
			var t747 rt.Value
			if t745 {
				t747 = rt.Flag(false)
			} else {
				t747 = rt.Flag(true)
			}
			t744 = t747
		} else {
			t744 = rt.Flag(false)
		}
		t748, e749 := rt.Cond(ctx, t744)
		if e749 != nil {
			return rt.Value{}, e749
		}
		if t748 {
			t750, e751 := rt.Cond(ctx, zagruzka)
			if e751 != nil {
				return rt.Value{}, e751
			}
			if t750 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Неизвестное") {
		t752, e753 := rt.Cond(ctx, kesh)
		if e753 != nil {
			return rt.Value{}, e753
		}
		var t754 rt.Value
		if t752 {
			t754 = rt.Flag(false)
		} else {
			t754 = rt.Flag(true)
		}
		t755, e756 := rt.Cond(ctx, t754)
		if e756 != nil {
			return rt.Value{}, e756
		}
		var t757 rt.Value
		if t755 {
			t758, e759 := rt.Cond(ctx, zhurnal)
			if e759 != nil {
				return rt.Value{}, e759
			}
			var t760 rt.Value
			if t758 {
				t760 = rt.Flag(false)
			} else {
				t760 = rt.Flag(true)
			}
			t757 = t760
		} else {
			t757 = rt.Flag(false)
		}
		t761, e762 := rt.Cond(ctx, t757)
		if e762 != nil {
			return rt.Value{}, e762
		}
		var t763 rt.Value
		if t761 {
			t764, e765 := rt.Cond(ctx, sborka)
			if e765 != nil {
				return rt.Value{}, e765
			}
			var t766 rt.Value
			if t764 {
				t766 = rt.Flag(false)
			} else {
				t766 = rt.Flag(true)
			}
			t763 = t766
		} else {
			t763 = rt.Flag(false)
		}
		t767, e768 := rt.Cond(ctx, t763)
		if e768 != nil {
			return rt.Value{}, e768
		}
		var t769 rt.Value
		if t767 {
			t770, e771 := rt.Cond(ctx, zagruzka)
			if e771 != nil {
				return rt.Value{}, e771
			}
			var t772 rt.Value
			if t770 {
				t772 = rt.Flag(false)
			} else {
				t772 = rt.Flag(true)
			}
			t769 = t772
		} else {
			t769 = rt.Flag(false)
		}
		t773, e774 := rt.Cond(ctx, t769)
		if e774 != nil {
			return rt.Value{}, e774
		}
		if t773 {
			t775, e776 := rt.Cond(ctx, krupnoe)
			if e776 != nil {
				return rt.Value{}, e776
			}
			if t775 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else {
		return rt.Value{}, rt.MatchFail(ctx, razryad)
	}
}

// PrigovorObosnovan — функция flang «Приговор обоснован».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение.
func PrigovorObosnovan(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value, prigovor rt.Value) (rt.Value, error) {
	t777, e778 := rt.FieldGet(ctx, nahodka, "доступен")
	if e778 != nil {
		return rt.Value{}, e778
	}
	t779, e780 := rt.Cond(ctx, t777)
	if e780 != nil {
		return rt.Value{}, e780
	}
	var t781 rt.Value
	if t779 {
		t781 = rt.Flag(false)
	} else {
		t781 = rt.Flag(true)
	}
	t782, e783 := rt.Cond(ctx, t781)
	if e783 != nil {
		return rt.Value{}, e783
	}
	if t782 {
		return EtoNeTrogat(ctx, prigovor)
	} else {
		t784, e785 := Ssylka(ctx, nahodka)
		if e785 != nil {
			return rt.Value{}, e785
		}
		t786, e787 := rt.Cond(ctx, t784)
		if e787 != nil {
			return rt.Value{}, e787
		}
		if t786 {
			return EtoNeTrogat(ctx, prigovor)
		} else {
			t788, e789 := EtoMozhnoUbrat(ctx, prigovor)
			if e789 != nil {
				return rt.Value{}, e789
			}
			t790, e791 := rt.Cond(ctx, t788)
			if e791 != nil {
				return rt.Value{}, e791
			}
			if t790 {
				t792, e793 := I1NaPare(ctx, razryad, prigovor)
				if e793 != nil {
					return rt.Value{}, e793
				}
				t794, e795 := rt.Cond(ctx, t792)
				if e795 != nil {
					return rt.Value{}, e795
				}
				var t796 rt.Value
				if t794 {
					t797, e798 := Katalog(ctx, nahodka)
					if e798 != nil {
						return rt.Value{}, e798
					}
					t799, e800 := rt.Cond(ctx, t797)
					if e800 != nil {
						return rt.Value{}, e800
					}
					var t801 rt.Value
					if t799 {
						t801 = rt.Flag(false)
					} else {
						t801 = rt.Flag(true)
					}
					t796 = t801
				} else {
					t796 = rt.Flag(false)
				}
				t802, e803 := rt.Cond(ctx, t796)
				if e803 != nil {
					return rt.Value{}, e803
				}
				if t802 {
					t804, e805 := rt.FieldGet(ctx, nahodka, "возраст_дней")
					if e805 != nil {
						return rt.Value{}, e805
					}
					t806, e807 := PorogRazryada(ctx, razryad)
					if e807 != nil {
						return rt.Value{}, e807
					}
					t808, e809 := rt.Gte(ctx, t804, t806)
					if e809 != nil {
						return rt.Value{}, e809
					}
					return t808, nil
				} else {
					return rt.Flag(false), nil
				}
			} else {
				return rt.Flag(true), nil
			}
		}
	}
}

// VesObosnovan — функция flang «Вес обоснован».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр prigovor — «приговор»: «Приговор».
// Параметр ves — «вес»: число.
// Результат — значение.
func VesObosnovan(ctx *rt.Ctx, nahodka rt.Value, prigovor rt.Value, ves rt.Value) (rt.Value, error) {
	t810, e811 := EtoNeTrogat(ctx, prigovor)
	if e811 != nil {
		return rt.Value{}, e811
	}
	t812, e813 := rt.Cond(ctx, t810)
	if e813 != nil {
		return rt.Value{}, e813
	}
	if t812 {
		return rt.Flag(rt.Equal(ves, rt.Number(0.0))), nil
	} else {
		t814, e815 := rt.FieldGet(ctx, nahodka, "размер")
		if e815 != nil {
			return rt.Value{}, e815
		}
		t816, e817 := rt.Cond(ctx, rt.Flag(rt.Equal(ves, t814)))
		if e817 != nil {
			return rt.Value{}, e817
		}
		if t816 {
			return VesVGranicah(ctx, nahodka, ves)
		} else {
			return rt.Flag(false), nil
		}
	}
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Порог крупного":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Порог крупного", 0, len(args))
		}
		return PorogKrupnogo(ctx)
	case "Порог кэша":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Порог кэша", 0, len(args))
		}
		return PorogKesha(ctx)
	case "Порог журнала":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Порог журнала", 0, len(args))
		}
		return PorogZhurnala(ctx)
	case "Порог загрузки":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Порог загрузки", 0, len(args))
		}
		return PorogZagruzki(ctx)
	case "Составляющие пути":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Составляющие пути", 1, len(args))
		}
		return SostavlyayuschiePuti(ctx, args[0])
	case "Имя в пути":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Имя в пути", 1, len(args))
		}
		return ImyaVPuti(ctx, args[0])
	case "Есть составляющая":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Есть составляющая", 2, len(args))
		}
		return EstSostavlyayuschaya(ctx, args[0], args[1])
	case "Оканчивается на":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Оканчивается на", 2, len(args))
		}
		return OkanchivaetsyaNa(ctx, args[0], args[1])
	case "След пути":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"След пути", 1, len(args))
		}
		return SledPuti(ctx, args[0])
	case "Цепь ограничена":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Цепь ограничена", 1, len(args))
		}
		return CepOgranichena(ctx, args[0])
	case "Справочник ограничен":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Справочник ограничен", 1, len(args))
		}
		return SpravochnikOgranichen(ctx, args[0])
	case "Разряд места допустим":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд места допустим", 1, len(args))
		}
		return RazryadMestaDopustim(ctx, args[0])
	case "Место подходит":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Место подходит", 2, len(args))
		}
		return MestoPodhodit(ctx, args[0], args[1])
	case "Это неизвестное":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Это неизвестное", 1, len(args))
		}
		return EtoNeizvestnoe(ctx, args[0])
	case "Номер разряда":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Номер разряда", 1, len(args))
		}
		return NomerRazryada(ctx, args[0])
	case "Тот же разряд":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Тот же разряд", 2, len(args))
		}
		return TotZheRazryad(ctx, args[0], args[1])
	case "Место обосновано":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Место обосновано", 3, len(args))
		}
		return MestoObosnovano(ctx, args[0], args[1], args[2])
	case "Разряд по справочнику":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд по справочнику", 2, len(args))
		}
		return RazryadPoSpravochniku(ctx, args[0], args[1])
	case "Шестнадцатеричный знак":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Шестнадцатеричный знак", 1, len(args))
		}
		return ShestnadcaterichnyyZnak(ctx, args[0])
	case "Похоже на отпечаток":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Похоже на отпечаток", 1, len(args))
		}
		return PohozheNaOtpechatok(ctx, args[0])
	case "Адресуется содержимым":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Адресуется содержимым", 1, len(args))
		}
		return AdresuetsyaSoderzhimym(ctx, args[0])
	case "Под системным временным":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Под системным временным", 1, len(args))
		}
		return PodSistemnymVremennym(ctx, args[0])
	case "Примета кэша":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Примета кэша", 1, len(args))
		}
		return PrimetaKesha(ctx, args[0])
	case "Примета журнала":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Примета журнала", 1, len(args))
		}
		return PrimetaZhurnala(ctx, args[0])
	case "Примета сборки":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Примета сборки", 1, len(args))
		}
		return PrimetaSborki(ctx, args[0])
	case "Примета загрузки":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Примета загрузки", 1, len(args))
		}
		return PrimetaZagruzki(ctx, args[0])
	case "Есть примета":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Есть примета", 2, len(args))
		}
		return EstPrimeta(ctx, args[0], args[1])
	case "Разряд решён размером":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд решён размером", 1, len(args))
		}
		return RazryadReshyonRazmerom(ctx, args[0])
	case "Разряд находки":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд находки", 2, len(args))
		}
		return RazryadNahodki(ctx, args[0], args[1])
	case "Крупное не мельче порога":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Крупное не мельче порога", 2, len(args))
		}
		return KrupnoeNeMelchePoroga(ctx, args[0], args[1])
	case "Каталог":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Каталог", 1, len(args))
		}
		return Katalog(ctx, args[0])
	case "Ссылка":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Ссылка", 1, len(args))
		}
		return Ssylka(ctx, args[0])
	case "Приговор мусора":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Приговор мусора", 2, len(args))
		}
		return PrigovorMusora(ctx, args[0], args[1])
	case "Приговор находки":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Приговор находки", 2, len(args))
		}
		return PrigovorNahodki(ctx, args[0], args[1])
	case "Вес находки":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Вес находки", 2, len(args))
		}
		return VesNahodki(ctx, args[0], args[1])
	case "Вес в границах":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Вес в границах", 2, len(args))
		}
		return VesVGranicah(ctx, args[0], args[1])
	case "Решить находку":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Решить находку", 2, len(args))
		}
		return ReshitNahodku(ctx, args[0], args[1])
	case "Решить всё":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Решить всё", 2, len(args))
		}
		return ReshitVsyo(ctx, args[0], args[1])
	case "И1 держится":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"И1 держится", 1, len(args))
		}
		return I1Derzhitsya(ctx, args[0])
	case "И1 держится всюду":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"И1 держится всюду", 1, len(args))
		}
		return I1DerzhitsyaVsyudu(ctx, args[0])
	case "Пустой свод":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Пустой свод", 0, len(args))
		}
		return PustoySvod(ctx)
	case "Прибавить решение":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прибавить решение", 2, len(args))
		}
		return PribavitReshenie(ctx, args[0], args[1])
	case "Свести":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Свести", 2, len(args))
		}
		return Svesti(ctx, args[0], args[1])
	case "Сумма размеров убираемых":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Сумма размеров убираемых", 2, len(args))
		}
		return SummaRazmerovUbiraemyh(ctx, args[0], args[1])
	case "И2 держится":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"И2 держится", 3, len(args))
		}
		return I2Derzhitsya(ctx, args[0], args[1], args[2])
	case "Строку отчёта":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Строку отчёта", 2, len(args))
		}
		return StrokuOtchyota(ctx, args[0], args[1])
	case "Вставить по весу":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Вставить по весу", 2, len(args))
		}
		return VstavitPoVesu(ctx, args[0], args[1])
	case "Приписать строку отчёта":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Приписать строку отчёта", 2, len(args))
		}
		return PripisatStrokuOtchyota(ctx, args[0], args[1])
	case "Отчёт":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Отчёт", 2, len(args))
		}
		return Otchyot(ctx, args[0], args[1])
	case "Отчёт той же длины":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Отчёт той же длины", 2, len(args))
		}
		return OtchyotToyZheDliny(ctx, args[0], args[1])
	case "Это МожноУбрать":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Это МожноУбрать", 1, len(args))
		}
		return EtoMozhnoUbrat(ctx, args[0])
	case "Это НеТрогать":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Это НеТрогать", 1, len(args))
		}
		return EtoNeTrogat(ctx, args[0])
	case "И1 на паре":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"И1 на паре", 2, len(args))
		}
		return I1NaPare(ctx, args[0], args[1])
	case "Порог разряда":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Порог разряда", 1, len(args))
		}
		return PorogRazryada(ctx, args[0])
	case "Разряд обоснован":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд обоснован", 3, len(args))
		}
		return RazryadObosnovan(ctx, args[0], args[1], args[2])
	case "Разряд обоснован приметой":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд обоснован приметой", 2, len(args))
		}
		return RazryadObosnovanPrimetoy(ctx, args[0], args[1])
	case "Приговор обоснован":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Приговор обоснован", 3, len(args))
		}
		return PrigovorObosnovan(ctx, args[0], args[1], args[2])
	case "Вес обоснован":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Вес обоснован", 3, len(args))
		}
		return VesObosnovan(ctx, args[0], args[1], args[2])
	}
	return rt.Value{}, rt.Fail(rt.CodeUnknownName, "не найдена функция «%s»", name)
}

// Граница входа: объявленные типы параметров данными.
//
// Прогонщик сверяет по ним значения, пришедшие снаружи, ДО вызова
// (rt.CheckEntry). Виды rt.TypeUnknown (значение-функция, параметр
// полиморфизма, применение типа с аргументами) не сверяются — ровно как
// молчит о них проверка значений свидетеля.
var entryTypes = []rt.Type{}

var entryFields = []rt.TypeField{}

var entryVariants = []rt.TypeVariant{}

var entryParams = []rt.EntryParam{}

var entryTable = rt.EntryTable{Types: entryTypes, Fields: entryFields, Variants: entryVariants, Params: entryParams}

// Entry — объявленные типы параметров: по ним сверяется вход извне.
func Entry() *rt.EntryTable {
	return &entryTable
}
