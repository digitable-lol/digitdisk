// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
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

// VariantMusor — вариант «Мусор» суммы типов «Природа».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantMusor() rt.Value {
	return rt.Variant("Мусор", nil)
}

// VariantSvezhee — вариант «Свежее» суммы типов «Природа».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantSvezhee() rt.Value {
	return rt.Variant("Свежее", nil)
}

// VariantIshodniki — вариант «Исходники» суммы типов «Природа».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantIshodniki() rt.Value {
	return rt.Variant("Исходники", nil)
}

// VariantLichnoe — вариант «Личное» суммы типов «Природа».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantLichnoe() rt.Value {
	return rt.Variant("Личное", nil)
}

// VariantHranilische — вариант «Хранилище» суммы типов «Природа».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantHranilische() rt.Value {
	return rt.Variant("Хранилище", nil)
}

// VariantPodPrismotrom — вариант «ПодПрисмотром» суммы типов «Природа».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantPodPrismotrom() rt.Value {
	return rt.Variant("ПодПрисмотром", nil)
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

// RasshireniyaIshodnikov — функция flang «Расширения исходников».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: список: строка.
func RasshireniyaIshodnikov(ctx *rt.Ctx) (rt.Value, error) {
	t444 := make([]rt.Value, 30)
	t444[0] = rt.Text(".go")
	t444[1] = rt.Text(".c")
	t444[2] = rt.Text(".h")
	t444[3] = rt.Text(".cc")
	t444[4] = rt.Text(".cpp")
	t444[5] = rt.Text(".hpp")
	t444[6] = rt.Text(".rs")
	t444[7] = rt.Text(".py")
	t444[8] = rt.Text(".js")
	t444[9] = rt.Text(".ts")
	t444[10] = rt.Text(".java")
	t444[11] = rt.Text(".kt")
	t444[12] = rt.Text(".rb")
	t444[13] = rt.Text(".php")
	t444[14] = rt.Text(".cs")
	t444[15] = rt.Text(".swift")
	t444[16] = rt.Text(".ex")
	t444[17] = rt.Text(".exs")
	t444[18] = rt.Text(".erl")
	t444[19] = rt.Text(".lua")
	t444[20] = rt.Text(".pl")
	t444[21] = rt.Text(".sh")
	t444[22] = rt.Text(".sql")
	t444[23] = rt.Text(".flang")
	t444[24] = rt.Text(".md")
	t444[25] = rt.Text(".txt")
	t444[26] = rt.Text(".rst")
	t444[27] = rt.Text(".tex")
	t444[28] = rt.Text(".html")
	t444[29] = rt.Text(".css")
	return rt.List(t444), nil
}

// PrimetaIshodnika — функция flang «Примета исходника».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaIshodnika(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t445, e446 := ImyaVPuti(ctx, put)
	if e446 != nil {
		return rt.Value{}, e446
	}
	// пусть «имя»
	imya := t445
	t447, e448 := RasshireniyaIshodnikov(ctx)
	if e448 != nil {
		return rt.Value{}, e448
	}
	t449, e450 := rt.RequireList(ctx, t447, "отфильтровать")
	if e450 != nil {
		return rt.Value{}, e450
	}
	t451 := make([]rt.Value, 0, len(t449))
	for t452 := range t449 {
		// «хвост»
		hvost := t449[t452]
		t453, e454 := OkanchivaetsyaNa(ctx, imya, hvost)
		if e454 != nil {
			return rt.Value{}, e454
		}
		t455, e456 := rt.Keep(ctx, t453)
		if e456 != nil {
			return rt.Value{}, e456
		}
		if t455 {
			t451 = append(t451, hvost)
		}
	}
	// «длина»
	t457, e458 := rt.BLength(ctx, rt.List(t451))
	if e458 != nil {
		return rt.Value{}, e458
	}
	t459, e460 := rt.Gt(ctx, t457, rt.Number(0.0))
	if e460 != nil {
		return rt.Value{}, e460
	}
	return t459, nil
}

// PodPrismotromSistemyVersiy — функция flang «Под присмотром системы версий».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PodPrismotromSistemyVersiy(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t461, e462 := EstSostavlyayuschaya(ctx, put, rt.Text(".git"))
	if e462 != nil {
		return rt.Value{}, e462
	}
	t463, e464 := rt.Cond(ctx, t461)
	if e464 != nil {
		return rt.Value{}, e464
	}
	var t465 rt.Value
	if t463 {
		t465 = rt.Flag(true)
	} else {
		t466, e467 := EstSostavlyayuschaya(ctx, put, rt.Text(".hg"))
		if e467 != nil {
			return rt.Value{}, e467
		}
		t465 = t466
	}
	t468, e469 := rt.Cond(ctx, t465)
	if e469 != nil {
		return rt.Value{}, e469
	}
	if t468 {
		return rt.Flag(true), nil
	} else {
		return EstSostavlyayuschaya(ctx, put, rt.Text(".svn"))
	}
}

// MusornyyRazryad — функция flang «Мусорный разряд».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func MusornyyRazryad(ctx *rt.Ctx, razryad rt.Value) (rt.Value, error) {
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

// PrirodaNahodki — функция flang «Природа находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение: «Природа».
func PrirodaNahodki(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value, prigovor rt.Value) (rt.Value, error) {
	t470, e471 := rt.FieldGet(ctx, nahodka, "путь")
	if e471 != nil {
		return rt.Value{}, e471
	}
	t472, e473 := PodPrismotromSistemyVersiy(ctx, t470)
	if e473 != nil {
		return rt.Value{}, e473
	}
	t474, e475 := rt.Cond(ctx, t472)
	if e475 != nil {
		return rt.Value{}, e475
	}
	var t476 rt.Value
	if t474 {
		t476 = rt.Variant("ПодПрисмотром", nil)
	} else {
		t477, e478 := rt.FieldGet(ctx, nahodka, "путь")
		if e478 != nil {
			return rt.Value{}, e478
		}
		t479, e480 := AdresuetsyaSoderzhimym(ctx, t477)
		if e480 != nil {
			return rt.Value{}, e480
		}
		t481, e482 := rt.Cond(ctx, t479)
		if e482 != nil {
			return rt.Value{}, e482
		}
		var t483 rt.Value
		if t481 {
			t483 = rt.Variant("Хранилище", nil)
		} else {
			t484, e485 := EtoMozhnoUbrat(ctx, prigovor)
			if e485 != nil {
				return rt.Value{}, e485
			}
			t486, e487 := rt.Cond(ctx, t484)
			if e487 != nil {
				return rt.Value{}, e487
			}
			var t488 rt.Value
			if t486 {
				t488 = rt.Variant("Мусор", nil)
			} else {
				t489, e490 := rt.FieldGet(ctx, nahodka, "путь")
				if e490 != nil {
					return rt.Value{}, e490
				}
				t491, e492 := PrimetaIshodnika(ctx, t489)
				if e492 != nil {
					return rt.Value{}, e492
				}
				t493, e494 := rt.Cond(ctx, t491)
				if e494 != nil {
					return rt.Value{}, e494
				}
				var t495 rt.Value
				if t493 {
					t495 = rt.Variant("Исходники", nil)
				} else {
					t496, e497 := MusornyyRazryad(ctx, razryad)
					if e497 != nil {
						return rt.Value{}, e497
					}
					t498, e499 := rt.Cond(ctx, t496)
					if e499 != nil {
						return rt.Value{}, e499
					}
					var t500 rt.Value
					if t498 {
						t500 = rt.Variant("Свежее", nil)
					} else {
						t500 = rt.Variant("Личное", nil)
					}
					t495 = t500
				}
				t488 = t495
			}
			t483 = t488
		}
		t476 = t483
	}
	t501 := t476
	t502, e503 := PrirodaObosnovana(ctx, nahodka, razryad, prigovor, t501)
	if e503 != nil {
		return rt.Value{}, e503
	}
	// постусловие «Природа обоснована»
	t504, e505 := rt.Post(ctx, t502, "Природа обоснована", "Природа находки")
	if e505 != nil {
		return rt.Value{}, e505
	}
	if !t504 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Природа обоснована» функции «Природа находки»")
	}
	return t501, nil
}

// PrirodaObosnovana — функция flang «Природа обоснована».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Параметр prigovor — «приговор»: «Приговор».
// Параметр priroda — «природа»: «Природа».
// Результат — значение.
func PrirodaObosnovana(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value, prigovor rt.Value, priroda rt.Value) (rt.Value, error) {
	if rt.VariantIs(priroda, "Мусор") {
		return EtoMozhnoUbrat(ctx, prigovor)
	} else if rt.VariantIs(priroda, "ПодПрисмотром") {
		t506, e507 := rt.FieldGet(ctx, nahodka, "путь")
		if e507 != nil {
			return rt.Value{}, e507
		}
		return PodPrismotromSistemyVersiy(ctx, t506)
	} else if rt.VariantIs(priroda, "Хранилище") {
		t508, e509 := rt.FieldGet(ctx, nahodka, "путь")
		if e509 != nil {
			return rt.Value{}, e509
		}
		t510, e511 := AdresuetsyaSoderzhimym(ctx, t508)
		if e511 != nil {
			return rt.Value{}, e511
		}
		t512, e513 := rt.Cond(ctx, t510)
		if e513 != nil {
			return rt.Value{}, e513
		}
		if t512 {
			t514, e515 := rt.FieldGet(ctx, nahodka, "путь")
			if e515 != nil {
				return rt.Value{}, e515
			}
			t516, e517 := PodPrismotromSistemyVersiy(ctx, t514)
			if e517 != nil {
				return rt.Value{}, e517
			}
			t518, e519 := rt.Cond(ctx, t516)
			if e519 != nil {
				return rt.Value{}, e519
			}
			if t518 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(priroda, "Исходники") {
		t520, e521 := rt.FieldGet(ctx, nahodka, "путь")
		if e521 != nil {
			return rt.Value{}, e521
		}
		t522, e523 := PrimetaIshodnika(ctx, t520)
		if e523 != nil {
			return rt.Value{}, e523
		}
		t524, e525 := rt.Cond(ctx, t522)
		if e525 != nil {
			return rt.Value{}, e525
		}
		if t524 {
			t526, e527 := EtoMozhnoUbrat(ctx, prigovor)
			if e527 != nil {
				return rt.Value{}, e527
			}
			t528, e529 := rt.Cond(ctx, t526)
			if e529 != nil {
				return rt.Value{}, e529
			}
			if t528 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(priroda, "Свежее") {
		t530, e531 := MusornyyRazryad(ctx, razryad)
		if e531 != nil {
			return rt.Value{}, e531
		}
		t532, e533 := rt.Cond(ctx, t530)
		if e533 != nil {
			return rt.Value{}, e533
		}
		if t532 {
			t534, e535 := EtoMozhnoUbrat(ctx, prigovor)
			if e535 != nil {
				return rt.Value{}, e535
			}
			t536, e537 := rt.Cond(ctx, t534)
			if e537 != nil {
				return rt.Value{}, e537
			}
			if t536 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(priroda, "Личное") {
		t538, e539 := EtoMozhnoUbrat(ctx, prigovor)
		if e539 != nil {
			return rt.Value{}, e539
		}
		t540, e541 := rt.Cond(ctx, t538)
		if e541 != nil {
			return rt.Value{}, e541
		}
		if t540 {
			return rt.Flag(false), nil
		} else {
			return rt.Flag(true), nil
		}
	} else {
		return rt.Value{}, rt.MatchFail(ctx, priroda)
	}
}

// PrirodaPoNahodke — функция flang «Природа по находке».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: «Природа».
func PrirodaPoNahodke(ctx *rt.Ctx, nahodka rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t542, e543 := rt.FieldGet(ctx, nahodka, "путь")
	if e543 != nil {
		return rt.Value{}, e543
	}
	t544, e545 := RazryadPoSpravochniku(ctx, t542, spravochnik)
	if e545 != nil {
		return rt.Value{}, e545
	}
	// пусть «место»
	mesto := t544
	t546, e547 := RazryadNahodki(ctx, nahodka, mesto)
	if e547 != nil {
		return rt.Value{}, e547
	}
	// пусть «разряд»
	razryad := t546
	t548, e549 := PrigovorNahodki(ctx, nahodka, razryad)
	if e549 != nil {
		return rt.Value{}, e549
	}
	// пусть «приговор»
	prigovor := t548
	return PrirodaNahodki(ctx, nahodka, razryad, prigovor)
}

// Strogost — функция flang «Строгость».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр priroda — «природа»: «Природа».
// Результат — значение: число.
func Strogost(ctx *rt.Ctx, priroda rt.Value) (rt.Value, error) {
	if rt.VariantIs(priroda, "Мусор") {
		return rt.Number(1.0), nil
	} else if rt.VariantIs(priroda, "Свежее") {
		return rt.Number(2.0), nil
	} else if rt.VariantIs(priroda, "Личное") {
		return rt.Number(2.0), nil
	} else if rt.VariantIs(priroda, "Исходники") {
		return rt.Number(2.0), nil
	} else if rt.VariantIs(priroda, "Хранилище") {
		return rt.Number(3.0), nil
	} else if rt.VariantIs(priroda, "ПодПрисмотром") {
		return rt.Number(3.0), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, priroda)
	}
}

// I1Derzhitsya — функция flang «И1 держится».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр reshenie — «решение»: «Решение».
// Результат — значение.
func I1Derzhitsya(ctx *rt.Ctx, reshenie rt.Value) (rt.Value, error) {
	t550, e551 := rt.FieldGet(ctx, reshenie, "приговор")
	if e551 != nil {
		return rt.Value{}, e551
	}
	if rt.VariantIs(t550, "МожноУбрать") {
		t552, e553 := rt.FieldGet(ctx, reshenie, "разряд")
		if e553 != nil {
			return rt.Value{}, e553
		}
		if rt.VariantIs(t552, "Кэш") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t552, "Журнал") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t552, "Сборка") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t552, "Загрузка") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t552, "Крупное") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t552, "Неизвестное") {
			return rt.Flag(false), nil
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t552)
		}
	} else if rt.VariantIs(t550, "Спросить") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t550, "НеТрогать") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t550)
	}
}

// I1DerzhitsyaVsyudu — функция flang «И1 держится всюду».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр resheniya — «решения»: список: «Решение».
// Результат — значение.
func I1DerzhitsyaVsyudu(ctx *rt.Ctx, resheniya rt.Value) (rt.Value, error) {
	t554, e555 := rt.RequireList(ctx, resheniya, "свёртка")
	if e555 != nil {
		return rt.Value{}, e555
	}
	// «акк»
	akk := rt.Flag(true)
	for t556 := range t554 {
		// «решение»
		reshenie := t554[t556]
		t557, e558 := rt.Cond(ctx, akk)
		if e558 != nil {
			return rt.Value{}, e558
		}
		var t559 rt.Value
		if t557 {
			t560, e561 := I1Derzhitsya(ctx, reshenie)
			if e561 != nil {
				return rt.Value{}, e561
			}
			t559 = t560
		} else {
			t559 = rt.Flag(false)
		}
		akk = t559
	}
	return akk, nil
}

// PustoySvod — функция flang «Пустой свод».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: «Свод».
func PustoySvod(ctx *rt.Ctx) (rt.Value, error) {
	t562 := make([]rt.Field, 7)
	t562[0] = rt.Field{Name: "кэш", Value: rt.Number(0.0)}
	t562[1] = rt.Field{Name: "журнал", Value: rt.Number(0.0)}
	t562[2] = rt.Field{Name: "сборка", Value: rt.Number(0.0)}
	t562[3] = rt.Field{Name: "загрузка", Value: rt.Number(0.0)}
	t562[4] = rt.Field{Name: "крупное", Value: rt.Number(0.0)}
	t562[5] = rt.Field{Name: "неизвестное", Value: rt.Number(0.0)}
	t562[6] = rt.Field{Name: "освободить", Value: rt.Number(0.0)}
	return rt.Record(t562), nil
}

// PribavitReshenie — функция flang «Прибавить решение».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр svod — «свод»: «Свод».
// Параметр reshenie — «решение»: «Решение».
// Результат — значение: «Свод».
func PribavitReshenie(ctx *rt.Ctx, svod rt.Value, reshenie rt.Value) (rt.Value, error) {
	var t563 rt.Value
	t564, e565 := rt.FieldGet(ctx, reshenie, "приговор")
	if e565 != nil {
		return rt.Value{}, e565
	}
	if rt.VariantIs(t564, "МожноУбрать") {
		t566, e567 := rt.FieldGet(ctx, reshenie, "вес")
		if e567 != nil {
			return rt.Value{}, e567
		}
		t563 = t566
	} else if rt.VariantIs(t564, "Спросить") {
		t563 = rt.Number(0.0)
	} else if rt.VariantIs(t564, "НеТрогать") {
		t563 = rt.Number(0.0)
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t564)
	}
	// пусть «убрать»
	ubrat := t563
	t568, e569 := rt.FieldGet(ctx, reshenie, "разряд")
	if e569 != nil {
		return rt.Value{}, e569
	}
	if rt.VariantIs(t568, "Кэш") {
		t570, e571 := rt.FieldGet(ctx, svod, "кэш")
		if e571 != nil {
			return rt.Value{}, e571
		}
		t572, e573 := rt.FieldGet(ctx, reshenie, "вес")
		if e573 != nil {
			return rt.Value{}, e573
		}
		t574, e575 := rt.Add(ctx, t570, t572)
		if e575 != nil {
			return rt.Value{}, e575
		}
		t576, e577 := rt.FieldGet(ctx, svod, "журнал")
		if e577 != nil {
			return rt.Value{}, e577
		}
		t578, e579 := rt.FieldGet(ctx, svod, "сборка")
		if e579 != nil {
			return rt.Value{}, e579
		}
		t580, e581 := rt.FieldGet(ctx, svod, "загрузка")
		if e581 != nil {
			return rt.Value{}, e581
		}
		t582, e583 := rt.FieldGet(ctx, svod, "крупное")
		if e583 != nil {
			return rt.Value{}, e583
		}
		t584, e585 := rt.FieldGet(ctx, svod, "неизвестное")
		if e585 != nil {
			return rt.Value{}, e585
		}
		t586, e587 := rt.FieldGet(ctx, svod, "освободить")
		if e587 != nil {
			return rt.Value{}, e587
		}
		t588, e589 := rt.Add(ctx, t586, ubrat)
		if e589 != nil {
			return rt.Value{}, e589
		}
		t590 := make([]rt.Field, 7)
		t590[0] = rt.Field{Name: "кэш", Value: t574}
		t590[1] = rt.Field{Name: "журнал", Value: t576}
		t590[2] = rt.Field{Name: "сборка", Value: t578}
		t590[3] = rt.Field{Name: "загрузка", Value: t580}
		t590[4] = rt.Field{Name: "крупное", Value: t582}
		t590[5] = rt.Field{Name: "неизвестное", Value: t584}
		t590[6] = rt.Field{Name: "освободить", Value: t588}
		return rt.Record(t590), nil
	} else if rt.VariantIs(t568, "Журнал") {
		t591, e592 := rt.FieldGet(ctx, svod, "кэш")
		if e592 != nil {
			return rt.Value{}, e592
		}
		t593, e594 := rt.FieldGet(ctx, svod, "журнал")
		if e594 != nil {
			return rt.Value{}, e594
		}
		t595, e596 := rt.FieldGet(ctx, reshenie, "вес")
		if e596 != nil {
			return rt.Value{}, e596
		}
		t597, e598 := rt.Add(ctx, t593, t595)
		if e598 != nil {
			return rt.Value{}, e598
		}
		t599, e600 := rt.FieldGet(ctx, svod, "сборка")
		if e600 != nil {
			return rt.Value{}, e600
		}
		t601, e602 := rt.FieldGet(ctx, svod, "загрузка")
		if e602 != nil {
			return rt.Value{}, e602
		}
		t603, e604 := rt.FieldGet(ctx, svod, "крупное")
		if e604 != nil {
			return rt.Value{}, e604
		}
		t605, e606 := rt.FieldGet(ctx, svod, "неизвестное")
		if e606 != nil {
			return rt.Value{}, e606
		}
		t607, e608 := rt.FieldGet(ctx, svod, "освободить")
		if e608 != nil {
			return rt.Value{}, e608
		}
		t609, e610 := rt.Add(ctx, t607, ubrat)
		if e610 != nil {
			return rt.Value{}, e610
		}
		t611 := make([]rt.Field, 7)
		t611[0] = rt.Field{Name: "кэш", Value: t591}
		t611[1] = rt.Field{Name: "журнал", Value: t597}
		t611[2] = rt.Field{Name: "сборка", Value: t599}
		t611[3] = rt.Field{Name: "загрузка", Value: t601}
		t611[4] = rt.Field{Name: "крупное", Value: t603}
		t611[5] = rt.Field{Name: "неизвестное", Value: t605}
		t611[6] = rt.Field{Name: "освободить", Value: t609}
		return rt.Record(t611), nil
	} else if rt.VariantIs(t568, "Сборка") {
		t612, e613 := rt.FieldGet(ctx, svod, "кэш")
		if e613 != nil {
			return rt.Value{}, e613
		}
		t614, e615 := rt.FieldGet(ctx, svod, "журнал")
		if e615 != nil {
			return rt.Value{}, e615
		}
		t616, e617 := rt.FieldGet(ctx, svod, "сборка")
		if e617 != nil {
			return rt.Value{}, e617
		}
		t618, e619 := rt.FieldGet(ctx, reshenie, "вес")
		if e619 != nil {
			return rt.Value{}, e619
		}
		t620, e621 := rt.Add(ctx, t616, t618)
		if e621 != nil {
			return rt.Value{}, e621
		}
		t622, e623 := rt.FieldGet(ctx, svod, "загрузка")
		if e623 != nil {
			return rt.Value{}, e623
		}
		t624, e625 := rt.FieldGet(ctx, svod, "крупное")
		if e625 != nil {
			return rt.Value{}, e625
		}
		t626, e627 := rt.FieldGet(ctx, svod, "неизвестное")
		if e627 != nil {
			return rt.Value{}, e627
		}
		t628, e629 := rt.FieldGet(ctx, svod, "освободить")
		if e629 != nil {
			return rt.Value{}, e629
		}
		t630, e631 := rt.Add(ctx, t628, ubrat)
		if e631 != nil {
			return rt.Value{}, e631
		}
		t632 := make([]rt.Field, 7)
		t632[0] = rt.Field{Name: "кэш", Value: t612}
		t632[1] = rt.Field{Name: "журнал", Value: t614}
		t632[2] = rt.Field{Name: "сборка", Value: t620}
		t632[3] = rt.Field{Name: "загрузка", Value: t622}
		t632[4] = rt.Field{Name: "крупное", Value: t624}
		t632[5] = rt.Field{Name: "неизвестное", Value: t626}
		t632[6] = rt.Field{Name: "освободить", Value: t630}
		return rt.Record(t632), nil
	} else if rt.VariantIs(t568, "Загрузка") {
		t633, e634 := rt.FieldGet(ctx, svod, "кэш")
		if e634 != nil {
			return rt.Value{}, e634
		}
		t635, e636 := rt.FieldGet(ctx, svod, "журнал")
		if e636 != nil {
			return rt.Value{}, e636
		}
		t637, e638 := rt.FieldGet(ctx, svod, "сборка")
		if e638 != nil {
			return rt.Value{}, e638
		}
		t639, e640 := rt.FieldGet(ctx, svod, "загрузка")
		if e640 != nil {
			return rt.Value{}, e640
		}
		t641, e642 := rt.FieldGet(ctx, reshenie, "вес")
		if e642 != nil {
			return rt.Value{}, e642
		}
		t643, e644 := rt.Add(ctx, t639, t641)
		if e644 != nil {
			return rt.Value{}, e644
		}
		t645, e646 := rt.FieldGet(ctx, svod, "крупное")
		if e646 != nil {
			return rt.Value{}, e646
		}
		t647, e648 := rt.FieldGet(ctx, svod, "неизвестное")
		if e648 != nil {
			return rt.Value{}, e648
		}
		t649, e650 := rt.FieldGet(ctx, svod, "освободить")
		if e650 != nil {
			return rt.Value{}, e650
		}
		t651, e652 := rt.Add(ctx, t649, ubrat)
		if e652 != nil {
			return rt.Value{}, e652
		}
		t653 := make([]rt.Field, 7)
		t653[0] = rt.Field{Name: "кэш", Value: t633}
		t653[1] = rt.Field{Name: "журнал", Value: t635}
		t653[2] = rt.Field{Name: "сборка", Value: t637}
		t653[3] = rt.Field{Name: "загрузка", Value: t643}
		t653[4] = rt.Field{Name: "крупное", Value: t645}
		t653[5] = rt.Field{Name: "неизвестное", Value: t647}
		t653[6] = rt.Field{Name: "освободить", Value: t651}
		return rt.Record(t653), nil
	} else if rt.VariantIs(t568, "Крупное") {
		t654, e655 := rt.FieldGet(ctx, svod, "кэш")
		if e655 != nil {
			return rt.Value{}, e655
		}
		t656, e657 := rt.FieldGet(ctx, svod, "журнал")
		if e657 != nil {
			return rt.Value{}, e657
		}
		t658, e659 := rt.FieldGet(ctx, svod, "сборка")
		if e659 != nil {
			return rt.Value{}, e659
		}
		t660, e661 := rt.FieldGet(ctx, svod, "загрузка")
		if e661 != nil {
			return rt.Value{}, e661
		}
		t662, e663 := rt.FieldGet(ctx, svod, "крупное")
		if e663 != nil {
			return rt.Value{}, e663
		}
		t664, e665 := rt.FieldGet(ctx, reshenie, "вес")
		if e665 != nil {
			return rt.Value{}, e665
		}
		t666, e667 := rt.Add(ctx, t662, t664)
		if e667 != nil {
			return rt.Value{}, e667
		}
		t668, e669 := rt.FieldGet(ctx, svod, "неизвестное")
		if e669 != nil {
			return rt.Value{}, e669
		}
		t670, e671 := rt.FieldGet(ctx, svod, "освободить")
		if e671 != nil {
			return rt.Value{}, e671
		}
		t672, e673 := rt.Add(ctx, t670, ubrat)
		if e673 != nil {
			return rt.Value{}, e673
		}
		t674 := make([]rt.Field, 7)
		t674[0] = rt.Field{Name: "кэш", Value: t654}
		t674[1] = rt.Field{Name: "журнал", Value: t656}
		t674[2] = rt.Field{Name: "сборка", Value: t658}
		t674[3] = rt.Field{Name: "загрузка", Value: t660}
		t674[4] = rt.Field{Name: "крупное", Value: t666}
		t674[5] = rt.Field{Name: "неизвестное", Value: t668}
		t674[6] = rt.Field{Name: "освободить", Value: t672}
		return rt.Record(t674), nil
	} else if rt.VariantIs(t568, "Неизвестное") {
		t675, e676 := rt.FieldGet(ctx, svod, "кэш")
		if e676 != nil {
			return rt.Value{}, e676
		}
		t677, e678 := rt.FieldGet(ctx, svod, "журнал")
		if e678 != nil {
			return rt.Value{}, e678
		}
		t679, e680 := rt.FieldGet(ctx, svod, "сборка")
		if e680 != nil {
			return rt.Value{}, e680
		}
		t681, e682 := rt.FieldGet(ctx, svod, "загрузка")
		if e682 != nil {
			return rt.Value{}, e682
		}
		t683, e684 := rt.FieldGet(ctx, svod, "крупное")
		if e684 != nil {
			return rt.Value{}, e684
		}
		t685, e686 := rt.FieldGet(ctx, svod, "неизвестное")
		if e686 != nil {
			return rt.Value{}, e686
		}
		t687, e688 := rt.FieldGet(ctx, reshenie, "вес")
		if e688 != nil {
			return rt.Value{}, e688
		}
		t689, e690 := rt.Add(ctx, t685, t687)
		if e690 != nil {
			return rt.Value{}, e690
		}
		t691, e692 := rt.FieldGet(ctx, svod, "освободить")
		if e692 != nil {
			return rt.Value{}, e692
		}
		t693, e694 := rt.Add(ctx, t691, ubrat)
		if e694 != nil {
			return rt.Value{}, e694
		}
		t695 := make([]rt.Field, 7)
		t695[0] = rt.Field{Name: "кэш", Value: t675}
		t695[1] = rt.Field{Name: "журнал", Value: t677}
		t695[2] = rt.Field{Name: "сборка", Value: t679}
		t695[3] = rt.Field{Name: "загрузка", Value: t681}
		t695[4] = rt.Field{Name: "крупное", Value: t683}
		t695[5] = rt.Field{Name: "неизвестное", Value: t689}
		t695[6] = rt.Field{Name: "освободить", Value: t693}
		return rt.Record(t695), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t568)
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
	t696, e697 := ReshitVsyo(ctx, zapisi, spravochnik)
	if e697 != nil {
		return rt.Value{}, e697
	}
	t698, e699 := rt.RequireList(ctx, t696, "свёртка")
	if e699 != nil {
		return rt.Value{}, e699
	}
	t700, e701 := PustoySvod(ctx)
	if e701 != nil {
		return rt.Value{}, e701
	}
	// «свод»
	svod := t700
	for t702 := range t698 {
		// «решение»
		reshenie := t698[t702]
		t703, e704 := PribavitReshenie(ctx, svod, reshenie)
		if e704 != nil {
			return rt.Value{}, e704
		}
		svod = t703
	}
	t705 := svod
	t706, e707 := I2Derzhitsya(ctx, zapisi, spravochnik, t705)
	if e707 != nil {
		return rt.Value{}, e707
	}
	// постусловие «И2: освобождаемое не больше убираемого»
	t708, e709 := rt.Post(ctx, t706, "И2: освобождаемое не больше убираемого", "Свести")
	if e709 != nil {
		return rt.Value{}, e709
	}
	if !t708 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И2: освобождаемое не больше убираемого» функции «Свести»")
	}
	return t705, nil
}

// SummaRazmerovUbiraemyh — функция flang «Сумма размеров убираемых».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: число.
func SummaRazmerovUbiraemyh(ctx *rt.Ctx, zapisi rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t710, e711 := rt.RequireList(ctx, zapisi, "свёртка")
	if e711 != nil {
		return rt.Value{}, e711
	}
	// «акк»
	akk := rt.Number(0.0)
	for t712 := range t710 {
		// «находка»
		nahodka := t710[t712]
		var t713 rt.Value
		t714, e715 := ReshitNahodku(ctx, nahodka, spravochnik)
		if e715 != nil {
			return rt.Value{}, e715
		}
		t716, e717 := rt.FieldGet(ctx, t714, "приговор")
		if e717 != nil {
			return rt.Value{}, e717
		}
		if rt.VariantIs(t716, "МожноУбрать") {
			t718, e719 := rt.FieldGet(ctx, nahodka, "размер")
			if e719 != nil {
				return rt.Value{}, e719
			}
			t720, e721 := rt.Add(ctx, akk, t718)
			if e721 != nil {
				return rt.Value{}, e721
			}
			t713 = t720
		} else if rt.VariantIs(t716, "Спросить") {
			t713 = akk
		} else if rt.VariantIs(t716, "НеТрогать") {
			t713 = akk
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t716)
		}
		akk = t713
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
	t722, e723 := rt.FieldGet(ctx, svod, "освободить")
	if e723 != nil {
		return rt.Value{}, e723
	}
	t724, e725 := SummaRazmerovUbiraemyh(ctx, zapisi, spravochnik)
	if e725 != nil {
		return rt.Value{}, e725
	}
	t726, e727 := rt.Lte(ctx, t722, t724)
	if e727 != nil {
		return rt.Value{}, e727
	}
	return t726, nil
}

// StrokuOtchyota — функция flang «Строку отчёта».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр spravochnik — «справочник»: список: «Место».
// Результат — значение: «Строка отчёта».
func StrokuOtchyota(ctx *rt.Ctx, nahodka rt.Value, spravochnik rt.Value) (rt.Value, error) {
	t728, e729 := rt.FieldGet(ctx, nahodka, "путь")
	if e729 != nil {
		return rt.Value{}, e729
	}
	t730, e731 := ReshitNahodku(ctx, nahodka, spravochnik)
	if e731 != nil {
		return rt.Value{}, e731
	}
	t732 := make([]rt.Field, 2)
	t732[0] = rt.Field{Name: "путь", Value: t728}
	t732[1] = rt.Field{Name: "решение", Value: t730}
	return rt.Record(t732), nil
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
		t733 := make([]rt.Value, 1)
		t733[0] = stroka
		return rt.List(t733), nil
	} else if rt.ChainCons(stroki) {
		// голова «голова»
		golova := rt.ChainHead(stroki)
		// хвост «хвост»
		hvost := rt.ChainTail(stroki)
		t734, e735 := rt.FieldGet(ctx, stroka, "решение")
		if e735 != nil {
			return rt.Value{}, e735
		}
		t736, e737 := rt.FieldGet(ctx, t734, "вес")
		if e737 != nil {
			return rt.Value{}, e737
		}
		t738, e739 := rt.FieldGet(ctx, golova, "решение")
		if e739 != nil {
			return rt.Value{}, e739
		}
		t740, e741 := rt.FieldGet(ctx, t738, "вес")
		if e741 != nil {
			return rt.Value{}, e741
		}
		t742, e743 := rt.Gte(ctx, t736, t740)
		if e743 != nil {
			return rt.Value{}, e743
		}
		t744, e745 := rt.Cond(ctx, t742)
		if e745 != nil {
			return rt.Value{}, e745
		}
		if t744 {
			return PripisatStrokuOtchyota(ctx, stroka, stroki)
		} else {
			t746, e747 := VstavitPoVesu(ctx, stroka, hvost)
			if e747 != nil {
				return rt.Value{}, e747
			}
			return PripisatStrokuOtchyota(ctx, golova, t746)
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
	t748, e749 := rt.RequireList(ctx, stroki, "свёртка")
	if e749 != nil {
		return rt.Value{}, e749
	}
	t750 := make([]rt.Value, 1)
	t750[0] = pervaya
	// «акк»
	akk := rt.List(t750)
	for t751 := range t748 {
		// «эл»
		el := t748[t751]
		// «добавить»
		t752, e753 := rt.BAppend(ctx, el, akk)
		if e753 != nil {
			return rt.Value{}, e753
		}
		akk = t752
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
	t754, e755 := rt.RequireList(ctx, zapisi, "свёртка")
	if e755 != nil {
		return rt.Value{}, e755
	}
	// «акк»
	akk := rt.List(nil)
	for t756 := range t754 {
		// «находка»
		nahodka := t754[t756]
		t757, e758 := StrokuOtchyota(ctx, nahodka, spravochnik)
		if e758 != nil {
			return rt.Value{}, e758
		}
		t759, e760 := VstavitPoVesu(ctx, t757, akk)
		if e760 != nil {
			return rt.Value{}, e760
		}
		akk = t759
	}
	t761 := akk
	t762, e763 := OtchyotToyZheDliny(ctx, zapisi, t761)
	if e763 != nil {
		return rt.Value{}, e763
	}
	// постусловие «Отчёт той же длины»
	t764, e765 := rt.Post(ctx, t762, "Отчёт той же длины", "Отчёт")
	if e765 != nil {
		return rt.Value{}, e765
	}
	if !t764 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Отчёт той же длины» функции «Отчёт»")
	}
	return t761, nil
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
	t766, e767 := rt.BLength(ctx, stroki)
	if e767 != nil {
		return rt.Value{}, e767
	}
	// «длина»
	t768, e769 := rt.BLength(ctx, zapisi)
	if e769 != nil {
		return rt.Value{}, e769
	}
	return rt.Flag(rt.Equal(t766, t768)), nil
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
	t770, e771 := EtoMozhnoUbrat(ctx, prigovor)
	if e771 != nil {
		return rt.Value{}, e771
	}
	t772, e773 := rt.Cond(ctx, t770)
	if e773 != nil {
		return rt.Value{}, e773
	}
	var t774 rt.Value
	if t772 {
		t774 = rt.Flag(false)
	} else {
		t774 = rt.Flag(true)
	}
	t775, e776 := rt.Cond(ctx, t774)
	if e776 != nil {
		return rt.Value{}, e776
	}
	if t775 {
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
	t777, e778 := EtoNeizvestnoe(ctx, mesto)
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
	t784, e785 := rt.FieldGet(ctx, nahodka, "путь")
	if e785 != nil {
		return rt.Value{}, e785
	}
	t786, e787 := PrimetaKesha(ctx, t784)
	if e787 != nil {
		return rt.Value{}, e787
	}
	// пусть «кэш»
	kesh := t786
	t788, e789 := rt.FieldGet(ctx, nahodka, "путь")
	if e789 != nil {
		return rt.Value{}, e789
	}
	t790, e791 := PrimetaZhurnala(ctx, t788)
	if e791 != nil {
		return rt.Value{}, e791
	}
	// пусть «журнал»
	zhurnal := t790
	t792, e793 := rt.FieldGet(ctx, nahodka, "путь")
	if e793 != nil {
		return rt.Value{}, e793
	}
	t794, e795 := PrimetaSborki(ctx, t792)
	if e795 != nil {
		return rt.Value{}, e795
	}
	// пусть «сборка»
	sborka := t794
	t796, e797 := rt.FieldGet(ctx, nahodka, "путь")
	if e797 != nil {
		return rt.Value{}, e797
	}
	t798, e799 := PrimetaZagruzki(ctx, t796)
	if e799 != nil {
		return rt.Value{}, e799
	}
	// пусть «загрузка»
	zagruzka := t798
	t800, e801 := rt.FieldGet(ctx, nahodka, "размер")
	if e801 != nil {
		return rt.Value{}, e801
	}
	t802, e803 := PorogKrupnogo(ctx)
	if e803 != nil {
		return rt.Value{}, e803
	}
	t804, e805 := rt.Gte(ctx, t800, t802)
	if e805 != nil {
		return rt.Value{}, e805
	}
	// пусть «крупное»
	krupnoe := t804
	if rt.VariantIs(razryad, "Кэш") {
		return kesh, nil
	} else if rt.VariantIs(razryad, "Журнал") {
		t806, e807 := rt.Cond(ctx, zhurnal)
		if e807 != nil {
			return rt.Value{}, e807
		}
		if t806 {
			t808, e809 := rt.Cond(ctx, kesh)
			if e809 != nil {
				return rt.Value{}, e809
			}
			if t808 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Сборка") {
		t810, e811 := rt.Cond(ctx, sborka)
		if e811 != nil {
			return rt.Value{}, e811
		}
		var t812 rt.Value
		if t810 {
			t813, e814 := rt.Cond(ctx, kesh)
			if e814 != nil {
				return rt.Value{}, e814
			}
			var t815 rt.Value
			if t813 {
				t815 = rt.Flag(false)
			} else {
				t815 = rt.Flag(true)
			}
			t812 = t815
		} else {
			t812 = rt.Flag(false)
		}
		t816, e817 := rt.Cond(ctx, t812)
		if e817 != nil {
			return rt.Value{}, e817
		}
		if t816 {
			t818, e819 := rt.Cond(ctx, zhurnal)
			if e819 != nil {
				return rt.Value{}, e819
			}
			if t818 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Загрузка") {
		t820, e821 := rt.Cond(ctx, zagruzka)
		if e821 != nil {
			return rt.Value{}, e821
		}
		var t822 rt.Value
		if t820 {
			t823, e824 := rt.Cond(ctx, kesh)
			if e824 != nil {
				return rt.Value{}, e824
			}
			var t825 rt.Value
			if t823 {
				t825 = rt.Flag(false)
			} else {
				t825 = rt.Flag(true)
			}
			t822 = t825
		} else {
			t822 = rt.Flag(false)
		}
		t826, e827 := rt.Cond(ctx, t822)
		if e827 != nil {
			return rt.Value{}, e827
		}
		var t828 rt.Value
		if t826 {
			t829, e830 := rt.Cond(ctx, zhurnal)
			if e830 != nil {
				return rt.Value{}, e830
			}
			var t831 rt.Value
			if t829 {
				t831 = rt.Flag(false)
			} else {
				t831 = rt.Flag(true)
			}
			t828 = t831
		} else {
			t828 = rt.Flag(false)
		}
		t832, e833 := rt.Cond(ctx, t828)
		if e833 != nil {
			return rt.Value{}, e833
		}
		if t832 {
			t834, e835 := rt.Cond(ctx, sborka)
			if e835 != nil {
				return rt.Value{}, e835
			}
			if t834 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Крупное") {
		t836, e837 := rt.Cond(ctx, krupnoe)
		if e837 != nil {
			return rt.Value{}, e837
		}
		var t838 rt.Value
		if t836 {
			t839, e840 := rt.Cond(ctx, kesh)
			if e840 != nil {
				return rt.Value{}, e840
			}
			var t841 rt.Value
			if t839 {
				t841 = rt.Flag(false)
			} else {
				t841 = rt.Flag(true)
			}
			t838 = t841
		} else {
			t838 = rt.Flag(false)
		}
		t842, e843 := rt.Cond(ctx, t838)
		if e843 != nil {
			return rt.Value{}, e843
		}
		var t844 rt.Value
		if t842 {
			t845, e846 := rt.Cond(ctx, zhurnal)
			if e846 != nil {
				return rt.Value{}, e846
			}
			var t847 rt.Value
			if t845 {
				t847 = rt.Flag(false)
			} else {
				t847 = rt.Flag(true)
			}
			t844 = t847
		} else {
			t844 = rt.Flag(false)
		}
		t848, e849 := rt.Cond(ctx, t844)
		if e849 != nil {
			return rt.Value{}, e849
		}
		var t850 rt.Value
		if t848 {
			t851, e852 := rt.Cond(ctx, sborka)
			if e852 != nil {
				return rt.Value{}, e852
			}
			var t853 rt.Value
			if t851 {
				t853 = rt.Flag(false)
			} else {
				t853 = rt.Flag(true)
			}
			t850 = t853
		} else {
			t850 = rt.Flag(false)
		}
		t854, e855 := rt.Cond(ctx, t850)
		if e855 != nil {
			return rt.Value{}, e855
		}
		if t854 {
			t856, e857 := rt.Cond(ctx, zagruzka)
			if e857 != nil {
				return rt.Value{}, e857
			}
			if t856 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Неизвестное") {
		t858, e859 := rt.Cond(ctx, kesh)
		if e859 != nil {
			return rt.Value{}, e859
		}
		var t860 rt.Value
		if t858 {
			t860 = rt.Flag(false)
		} else {
			t860 = rt.Flag(true)
		}
		t861, e862 := rt.Cond(ctx, t860)
		if e862 != nil {
			return rt.Value{}, e862
		}
		var t863 rt.Value
		if t861 {
			t864, e865 := rt.Cond(ctx, zhurnal)
			if e865 != nil {
				return rt.Value{}, e865
			}
			var t866 rt.Value
			if t864 {
				t866 = rt.Flag(false)
			} else {
				t866 = rt.Flag(true)
			}
			t863 = t866
		} else {
			t863 = rt.Flag(false)
		}
		t867, e868 := rt.Cond(ctx, t863)
		if e868 != nil {
			return rt.Value{}, e868
		}
		var t869 rt.Value
		if t867 {
			t870, e871 := rt.Cond(ctx, sborka)
			if e871 != nil {
				return rt.Value{}, e871
			}
			var t872 rt.Value
			if t870 {
				t872 = rt.Flag(false)
			} else {
				t872 = rt.Flag(true)
			}
			t869 = t872
		} else {
			t869 = rt.Flag(false)
		}
		t873, e874 := rt.Cond(ctx, t869)
		if e874 != nil {
			return rt.Value{}, e874
		}
		var t875 rt.Value
		if t873 {
			t876, e877 := rt.Cond(ctx, zagruzka)
			if e877 != nil {
				return rt.Value{}, e877
			}
			var t878 rt.Value
			if t876 {
				t878 = rt.Flag(false)
			} else {
				t878 = rt.Flag(true)
			}
			t875 = t878
		} else {
			t875 = rt.Flag(false)
		}
		t879, e880 := rt.Cond(ctx, t875)
		if e880 != nil {
			return rt.Value{}, e880
		}
		if t879 {
			t881, e882 := rt.Cond(ctx, krupnoe)
			if e882 != nil {
				return rt.Value{}, e882
			}
			if t881 {
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
	t883, e884 := rt.FieldGet(ctx, nahodka, "доступен")
	if e884 != nil {
		return rt.Value{}, e884
	}
	t885, e886 := rt.Cond(ctx, t883)
	if e886 != nil {
		return rt.Value{}, e886
	}
	var t887 rt.Value
	if t885 {
		t887 = rt.Flag(false)
	} else {
		t887 = rt.Flag(true)
	}
	t888, e889 := rt.Cond(ctx, t887)
	if e889 != nil {
		return rt.Value{}, e889
	}
	if t888 {
		return EtoNeTrogat(ctx, prigovor)
	} else {
		t890, e891 := Ssylka(ctx, nahodka)
		if e891 != nil {
			return rt.Value{}, e891
		}
		t892, e893 := rt.Cond(ctx, t890)
		if e893 != nil {
			return rt.Value{}, e893
		}
		if t892 {
			return EtoNeTrogat(ctx, prigovor)
		} else {
			t894, e895 := EtoMozhnoUbrat(ctx, prigovor)
			if e895 != nil {
				return rt.Value{}, e895
			}
			t896, e897 := rt.Cond(ctx, t894)
			if e897 != nil {
				return rt.Value{}, e897
			}
			if t896 {
				t898, e899 := I1NaPare(ctx, razryad, prigovor)
				if e899 != nil {
					return rt.Value{}, e899
				}
				t900, e901 := rt.Cond(ctx, t898)
				if e901 != nil {
					return rt.Value{}, e901
				}
				var t902 rt.Value
				if t900 {
					t903, e904 := Katalog(ctx, nahodka)
					if e904 != nil {
						return rt.Value{}, e904
					}
					t905, e906 := rt.Cond(ctx, t903)
					if e906 != nil {
						return rt.Value{}, e906
					}
					var t907 rt.Value
					if t905 {
						t907 = rt.Flag(false)
					} else {
						t907 = rt.Flag(true)
					}
					t902 = t907
				} else {
					t902 = rt.Flag(false)
				}
				t908, e909 := rt.Cond(ctx, t902)
				if e909 != nil {
					return rt.Value{}, e909
				}
				if t908 {
					t910, e911 := rt.FieldGet(ctx, nahodka, "возраст_дней")
					if e911 != nil {
						return rt.Value{}, e911
					}
					t912, e913 := PorogRazryada(ctx, razryad)
					if e913 != nil {
						return rt.Value{}, e913
					}
					t914, e915 := rt.Gte(ctx, t910, t912)
					if e915 != nil {
						return rt.Value{}, e915
					}
					return t914, nil
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
	t916, e917 := EtoNeTrogat(ctx, prigovor)
	if e917 != nil {
		return rt.Value{}, e917
	}
	t918, e919 := rt.Cond(ctx, t916)
	if e919 != nil {
		return rt.Value{}, e919
	}
	if t918 {
		return rt.Flag(rt.Equal(ves, rt.Number(0.0))), nil
	} else {
		t920, e921 := rt.FieldGet(ctx, nahodka, "размер")
		if e921 != nil {
			return rt.Value{}, e921
		}
		t922, e923 := rt.Cond(ctx, rt.Flag(rt.Equal(ves, t920)))
		if e923 != nil {
			return rt.Value{}, e923
		}
		if t922 {
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
	case "Расширения исходников":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Расширения исходников", 0, len(args))
		}
		return RasshireniyaIshodnikov(ctx)
	case "Примета исходника":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Примета исходника", 1, len(args))
		}
		return PrimetaIshodnika(ctx, args[0])
	case "Под присмотром системы версий":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Под присмотром системы версий", 1, len(args))
		}
		return PodPrismotromSistemyVersiy(ctx, args[0])
	case "Мусорный разряд":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Мусорный разряд", 1, len(args))
		}
		return MusornyyRazryad(ctx, args[0])
	case "Природа находки":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Природа находки", 3, len(args))
		}
		return PrirodaNahodki(ctx, args[0], args[1], args[2])
	case "Природа обоснована":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Природа обоснована", 4, len(args))
		}
		return PrirodaObosnovana(ctx, args[0], args[1], args[2], args[3])
	case "Природа по находке":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Природа по находке", 2, len(args))
		}
		return PrirodaPoNahodke(ctx, args[0], args[1])
	case "Строгость":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Строгость", 1, len(args))
		}
		return Strogost(ctx, args[0])
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
