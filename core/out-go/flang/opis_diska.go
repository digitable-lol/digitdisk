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

// PrimetaKesha — функция flang «Примета кэша».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaKesha(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	// «содержит»
	t1, e2 := rt.BContains(ctx, put, rt.Text("/.cache/"))
	if e2 != nil {
		return rt.Value{}, e2
	}
	t3, e4 := rt.Cond(ctx, t1)
	if e4 != nil {
		return rt.Value{}, e4
	}
	var t5 rt.Value
	if t3 {
		t5 = rt.Flag(true)
	} else {
		// «содержит»
		t6, e7 := rt.BContains(ctx, put, rt.Text("/cache/"))
		if e7 != nil {
			return rt.Value{}, e7
		}
		t5 = t6
	}
	t8, e9 := rt.Cond(ctx, t5)
	if e9 != nil {
		return rt.Value{}, e9
	}
	var t10 rt.Value
	if t8 {
		t10 = rt.Flag(true)
	} else {
		// «содержит»
		t11, e12 := rt.BContains(ctx, put, rt.Text("/Caches/"))
		if e12 != nil {
			return rt.Value{}, e12
		}
		t10 = t11
	}
	t13, e14 := rt.Cond(ctx, t10)
	if e14 != nil {
		return rt.Value{}, e14
	}
	if t13 {
		return rt.Flag(true), nil
	} else {
		// «содержит»
		t15, e16 := rt.BContains(ctx, put, rt.Text("/tmp/"))
		if e16 != nil {
			return rt.Value{}, e16
		}
		return t15, nil
	}
}

// PrimetaZhurnala — функция flang «Примета журнала».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaZhurnala(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	// «содержит»
	t17, e18 := rt.BContains(ctx, put, rt.Text("/log/"))
	if e18 != nil {
		return rt.Value{}, e18
	}
	t19, e20 := rt.Cond(ctx, t17)
	if e20 != nil {
		return rt.Value{}, e20
	}
	var t21 rt.Value
	if t19 {
		t21 = rt.Flag(true)
	} else {
		// «содержит»
		t22, e23 := rt.BContains(ctx, put, rt.Text("/logs/"))
		if e23 != nil {
			return rt.Value{}, e23
		}
		t21 = t22
	}
	t24, e25 := rt.Cond(ctx, t21)
	if e25 != nil {
		return rt.Value{}, e25
	}
	if t24 {
		return rt.Flag(true), nil
	} else {
		// «содержит»
		t26, e27 := rt.BContains(ctx, put, rt.Text(".log"))
		if e27 != nil {
			return rt.Value{}, e27
		}
		return t26, nil
	}
}

// PrimetaSborki — функция flang «Примета сборки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaSborki(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	// «содержит»
	t28, e29 := rt.BContains(ctx, put, rt.Text("/node_modules/"))
	if e29 != nil {
		return rt.Value{}, e29
	}
	t30, e31 := rt.Cond(ctx, t28)
	if e31 != nil {
		return rt.Value{}, e31
	}
	var t32 rt.Value
	if t30 {
		t32 = rt.Flag(true)
	} else {
		// «содержит»
		t33, e34 := rt.BContains(ctx, put, rt.Text("/target/"))
		if e34 != nil {
			return rt.Value{}, e34
		}
		t32 = t33
	}
	t35, e36 := rt.Cond(ctx, t32)
	if e36 != nil {
		return rt.Value{}, e36
	}
	var t37 rt.Value
	if t35 {
		t37 = rt.Flag(true)
	} else {
		// «содержит»
		t38, e39 := rt.BContains(ctx, put, rt.Text("/build/"))
		if e39 != nil {
			return rt.Value{}, e39
		}
		t37 = t38
	}
	t40, e41 := rt.Cond(ctx, t37)
	if e41 != nil {
		return rt.Value{}, e41
	}
	var t42 rt.Value
	if t40 {
		t42 = rt.Flag(true)
	} else {
		// «содержит»
		t43, e44 := rt.BContains(ctx, put, rt.Text("/_build/"))
		if e44 != nil {
			return rt.Value{}, e44
		}
		t42 = t43
	}
	t45, e46 := rt.Cond(ctx, t42)
	if e46 != nil {
		return rt.Value{}, e46
	}
	if t45 {
		return rt.Flag(true), nil
	} else {
		// «содержит»
		t47, e48 := rt.BContains(ctx, put, rt.Text("/.gradle/"))
		if e48 != nil {
			return rt.Value{}, e48
		}
		return t47, nil
	}
}

// PrimetaZagruzki — функция flang «Примета загрузки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaZagruzki(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	// «содержит»
	t49, e50 := rt.BContains(ctx, put, rt.Text("/Downloads/"))
	if e50 != nil {
		return rt.Value{}, e50
	}
	t51, e52 := rt.Cond(ctx, t49)
	if e52 != nil {
		return rt.Value{}, e52
	}
	if t51 {
		return rt.Flag(true), nil
	} else {
		// «содержит»
		t53, e54 := rt.BContains(ctx, put, rt.Text("/Загрузки/"))
		if e54 != nil {
			return rt.Value{}, e54
		}
		return t53, nil
	}
}

// RazryadNahodki — функция flang «Разряд находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение: «Разряд».
func RazryadNahodki(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t55, e56 := rt.FieldGet(ctx, nahodka, "путь")
	if e56 != nil {
		return rt.Value{}, e56
	}
	t57, e58 := PrimetaKesha(ctx, t55)
	if e58 != nil {
		return rt.Value{}, e58
	}
	t59, e60 := rt.Cond(ctx, t57)
	if e60 != nil {
		return rt.Value{}, e60
	}
	var t61 rt.Value
	if t59 {
		t61 = rt.Variant("Кэш", nil)
	} else {
		t62, e63 := rt.FieldGet(ctx, nahodka, "путь")
		if e63 != nil {
			return rt.Value{}, e63
		}
		t64, e65 := PrimetaZhurnala(ctx, t62)
		if e65 != nil {
			return rt.Value{}, e65
		}
		t66, e67 := rt.Cond(ctx, t64)
		if e67 != nil {
			return rt.Value{}, e67
		}
		var t68 rt.Value
		if t66 {
			t68 = rt.Variant("Журнал", nil)
		} else {
			t69, e70 := rt.FieldGet(ctx, nahodka, "путь")
			if e70 != nil {
				return rt.Value{}, e70
			}
			t71, e72 := PrimetaSborki(ctx, t69)
			if e72 != nil {
				return rt.Value{}, e72
			}
			t73, e74 := rt.Cond(ctx, t71)
			if e74 != nil {
				return rt.Value{}, e74
			}
			var t75 rt.Value
			if t73 {
				t75 = rt.Variant("Сборка", nil)
			} else {
				t76, e77 := rt.FieldGet(ctx, nahodka, "путь")
				if e77 != nil {
					return rt.Value{}, e77
				}
				t78, e79 := PrimetaZagruzki(ctx, t76)
				if e79 != nil {
					return rt.Value{}, e79
				}
				t80, e81 := rt.Cond(ctx, t78)
				if e81 != nil {
					return rt.Value{}, e81
				}
				var t82 rt.Value
				if t80 {
					t82 = rt.Variant("Загрузка", nil)
				} else {
					t83, e84 := rt.FieldGet(ctx, nahodka, "размер")
					if e84 != nil {
						return rt.Value{}, e84
					}
					t85, e86 := PorogKrupnogo(ctx)
					if e86 != nil {
						return rt.Value{}, e86
					}
					t87, e88 := rt.Gte(ctx, t83, t85)
					if e88 != nil {
						return rt.Value{}, e88
					}
					t89, e90 := rt.Cond(ctx, t87)
					if e90 != nil {
						return rt.Value{}, e90
					}
					var t91 rt.Value
					if t89 {
						t91 = rt.Variant("Крупное", nil)
					} else {
						t91 = rt.Variant("Неизвестное", nil)
					}
					t82 = t91
				}
				t75 = t82
			}
			t68 = t75
		}
		t61 = t68
	}
	t92 := t61
	t93, e94 := RazryadObosnovan(ctx, nahodka, t92)
	if e94 != nil {
		return rt.Value{}, e94
	}
	// постусловие «Разряд обоснован приметой»
	t95, e96 := rt.Post(ctx, t93, "Разряд обоснован приметой", "Разряд находки")
	if e96 != nil {
		return rt.Value{}, e96
	}
	if !t95 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «Разряд обоснован приметой» ядра «Опись диска»")
	}
	return t92, nil
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
		t97, e98 := rt.FieldGet(ctx, nahodka, "размер")
		if e98 != nil {
			return rt.Value{}, e98
		}
		t99, e100 := PorogKrupnogo(ctx)
		if e100 != nil {
			return rt.Value{}, e100
		}
		t101, e102 := rt.Gte(ctx, t97, t99)
		if e102 != nil {
			return rt.Value{}, e102
		}
		return t101, nil
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
	t103, e104 := rt.FieldGet(ctx, nahodka, "вид")
	if e104 != nil {
		return rt.Value{}, e104
	}
	if rt.VariantIs(t103, "Каталог") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t103, "Файл") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t103, "Ссылка") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t103)
	}
}

// Ssylka — функция flang «Ссылка».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение.
func Ssylka(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t105, e106 := rt.FieldGet(ctx, nahodka, "вид")
	if e106 != nil {
		return rt.Value{}, e106
	}
	if rt.VariantIs(t105, "Ссылка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t105, "Файл") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t105, "Каталог") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t105)
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
	t107, e108 := Katalog(ctx, nahodka)
	if e108 != nil {
		return rt.Value{}, e108
	}
	t109, e110 := rt.Cond(ctx, t107)
	if e110 != nil {
		return rt.Value{}, e110
	}
	var t111 rt.Value
	if t109 {
		t111 = rt.Variant("Спросить", nil)
	} else {
		t112, e113 := rt.FieldGet(ctx, nahodka, "возраст_дней")
		if e113 != nil {
			return rt.Value{}, e113
		}
		t114, e115 := rt.Gte(ctx, t112, porog)
		if e115 != nil {
			return rt.Value{}, e115
		}
		t116, e117 := rt.Cond(ctx, t114)
		if e117 != nil {
			return rt.Value{}, e117
		}
		var t118 rt.Value
		if t116 {
			t118 = rt.Variant("МожноУбрать", nil)
		} else {
			t118 = rt.Variant("Спросить", nil)
		}
		t111 = t118
	}
	t119 := t111
	t120, e121 := Katalog(ctx, nahodka)
	if e121 != nil {
		return rt.Value{}, e121
	}
	t122, e123 := rt.Cond(ctx, t120)
	if e123 != nil {
		return rt.Value{}, e123
	}
	var t124 rt.Value
	if t122 {
		t125, e126 := EtoMozhnoUbrat(ctx, t119)
		if e126 != nil {
			return rt.Value{}, e126
		}
		t127, e128 := rt.Cond(ctx, t125)
		if e128 != nil {
			return rt.Value{}, e128
		}
		var t129 rt.Value
		if t127 {
			t129 = rt.Flag(false)
		} else {
			t129 = rt.Flag(true)
		}
		t124 = t129
	} else {
		t124 = rt.Flag(true)
	}
	// постусловие «Каталог не убирается»
	t130, e131 := rt.Post(ctx, t124, "Каталог не убирается", "Приговор мусора")
	if e131 != nil {
		return rt.Value{}, e131
	}
	if !t130 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «Каталог не убирается» ядра «Опись диска»")
	}
	return t119, nil
}

// PrigovorNahodki — функция flang «Приговор находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение: «Приговор».
func PrigovorNahodki(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value) (rt.Value, error) {
	t132, e133 := rt.FieldGet(ctx, nahodka, "доступен")
	if e133 != nil {
		return rt.Value{}, e133
	}
	t134, e135 := rt.Cond(ctx, t132)
	if e135 != nil {
		return rt.Value{}, e135
	}
	var t136 rt.Value
	if t134 {
		t136 = rt.Flag(false)
	} else {
		t136 = rt.Flag(true)
	}
	t137, e138 := rt.Cond(ctx, t136)
	if e138 != nil {
		return rt.Value{}, e138
	}
	var t139 rt.Value
	if t137 {
		t139 = rt.Variant("НеТрогать", nil)
	} else {
		t140, e141 := Ssylka(ctx, nahodka)
		if e141 != nil {
			return rt.Value{}, e141
		}
		t142, e143 := rt.Cond(ctx, t140)
		if e143 != nil {
			return rt.Value{}, e143
		}
		var t144 rt.Value
		if t142 {
			t144 = rt.Variant("НеТрогать", nil)
		} else {
			var t145 rt.Value
			if rt.VariantIs(razryad, "Кэш") {
				t146, e147 := PorogKesha(ctx)
				if e147 != nil {
					return rt.Value{}, e147
				}
				t148, e149 := PrigovorMusora(ctx, nahodka, t146)
				if e149 != nil {
					return rt.Value{}, e149
				}
				t145 = t148
			} else if rt.VariantIs(razryad, "Сборка") {
				t150, e151 := PorogKesha(ctx)
				if e151 != nil {
					return rt.Value{}, e151
				}
				t152, e153 := PrigovorMusora(ctx, nahodka, t150)
				if e153 != nil {
					return rt.Value{}, e153
				}
				t145 = t152
			} else if rt.VariantIs(razryad, "Журнал") {
				t154, e155 := PorogZhurnala(ctx)
				if e155 != nil {
					return rt.Value{}, e155
				}
				t156, e157 := PrigovorMusora(ctx, nahodka, t154)
				if e157 != nil {
					return rt.Value{}, e157
				}
				t145 = t156
			} else if rt.VariantIs(razryad, "Загрузка") {
				t158, e159 := rt.FieldGet(ctx, nahodka, "возраст_дней")
				if e159 != nil {
					return rt.Value{}, e159
				}
				t160, e161 := PorogZagruzki(ctx)
				if e161 != nil {
					return rt.Value{}, e161
				}
				t162, e163 := rt.Gte(ctx, t158, t160)
				if e163 != nil {
					return rt.Value{}, e163
				}
				t164, e165 := rt.Cond(ctx, t162)
				if e165 != nil {
					return rt.Value{}, e165
				}
				var t166 rt.Value
				if t164 {
					t166 = rt.Variant("Спросить", nil)
				} else {
					t166 = rt.Variant("НеТрогать", nil)
				}
				t145 = t166
			} else if rt.VariantIs(razryad, "Крупное") {
				t145 = rt.Variant("Спросить", nil)
			} else if rt.VariantIs(razryad, "Неизвестное") {
				t145 = rt.Variant("НеТрогать", nil)
			} else {
				return rt.Value{}, rt.MatchFail(ctx, razryad)
			}
			t144 = t145
		}
		t139 = t144
	}
	t167 := t139
	t168, e169 := PrigovorObosnovan(ctx, nahodka, razryad, t167)
	if e169 != nil {
		return rt.Value{}, e169
	}
	// постусловие «Приговор обоснован»
	t170, e171 := rt.Post(ctx, t168, "Приговор обоснован", "Приговор находки")
	if e171 != nil {
		return rt.Value{}, e171
	}
	if !t170 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «Приговор обоснован» ядра «Опись диска»")
	}
	return t167, nil
}

// VesNahodki — функция flang «Вес находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение: число.
func VesNahodki(ctx *rt.Ctx, nahodka rt.Value, prigovor rt.Value) (rt.Value, error) {
	var t172 rt.Value
	if rt.VariantIs(prigovor, "НеТрогать") {
		t172 = rt.Number(0.0)
	} else if rt.VariantIs(prigovor, "МожноУбрать") {
		t173, e174 := rt.FieldGet(ctx, nahodka, "размер")
		if e174 != nil {
			return rt.Value{}, e174
		}
		t172 = t173
	} else if rt.VariantIs(prigovor, "Спросить") {
		t175, e176 := rt.FieldGet(ctx, nahodka, "размер")
		if e176 != nil {
			return rt.Value{}, e176
		}
		t172 = t175
	} else {
		return rt.Value{}, rt.MatchFail(ctx, prigovor)
	}
	t177 := t172
	t178, e179 := VesObosnovan(ctx, nahodka, prigovor, t177)
	if e179 != nil {
		return rt.Value{}, e179
	}
	// постусловие «Вес обоснован»
	t180, e181 := rt.Post(ctx, t178, "Вес обоснован", "Вес находки")
	if e181 != nil {
		return rt.Value{}, e181
	}
	if !t180 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «Вес обоснован» ядра «Опись диска»")
	}
	return t177, nil
}

// VesVGranicah — функция flang «Вес в границах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр ves — «вес»: число.
// Результат — значение.
func VesVGranicah(ctx *rt.Ctx, nahodka rt.Value, ves rt.Value) (rt.Value, error) {
	t182, e183 := rt.Gte(ctx, ves, rt.Number(0.0))
	if e183 != nil {
		return rt.Value{}, e183
	}
	t184, e185 := rt.Cond(ctx, t182)
	if e185 != nil {
		return rt.Value{}, e185
	}
	if t184 {
		t186, e187 := rt.FieldGet(ctx, nahodka, "размер")
		if e187 != nil {
			return rt.Value{}, e187
		}
		t188, e189 := rt.Lte(ctx, ves, t186)
		if e189 != nil {
			return rt.Value{}, e189
		}
		return t188, nil
	} else {
		return rt.Flag(false), nil
	}
}

// ReshitNahodku — функция flang «Решить находку».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение: «Решение».
func ReshitNahodku(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t190, e191 := RazryadNahodki(ctx, nahodka)
	if e191 != nil {
		return rt.Value{}, e191
	}
	// пусть «разряд»
	razryad := t190
	t192, e193 := PrigovorNahodki(ctx, nahodka, razryad)
	if e193 != nil {
		return rt.Value{}, e193
	}
	// пусть «приговор»
	prigovor := t192
	t194, e195 := VesNahodki(ctx, nahodka, prigovor)
	if e195 != nil {
		return rt.Value{}, e195
	}
	t196 := make([]rt.Field, 3)
	t196[0] = rt.Field{Name: "разряд", Value: razryad}
	t196[1] = rt.Field{Name: "приговор", Value: prigovor}
	t196[2] = rt.Field{Name: "вес", Value: t194}
	t197 := rt.Record(t196)
	t198, e199 := I1Derzhitsya(ctx, t197)
	if e199 != nil {
		return rt.Value{}, e199
	}
	// постусловие «И1: убрать можно только мусор»
	t200, e201 := rt.Post(ctx, t198, "И1: убрать можно только мусор", "Решить находку")
	if e201 != nil {
		return rt.Value{}, e201
	}
	if !t200 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «И1: убрать можно только мусор» ядра «Опись диска»")
	}
	return t197, nil
}

// ReshitVsyo — функция flang «Решить всё».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: список: «Решение».
func ReshitVsyo(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t202, e203 := rt.RequireList(ctx, zapisi, "отобразить")
	if e203 != nil {
		return rt.Value{}, e203
	}
	t204 := make([]rt.Value, 0, len(t202))
	for t205 := range t202 {
		// «находка»
		nahodka := t202[t205]
		t206, e207 := ReshitNahodku(ctx, nahodka)
		if e207 != nil {
			return rt.Value{}, e207
		}
		t204 = append(t204, t206)
	}
	t208 := rt.List(t204)
	t209, e210 := I1DerzhitsyaVsyudu(ctx, t208)
	if e210 != nil {
		return rt.Value{}, e210
	}
	// постусловие «И1 всюду: убрать можно только мусор»
	t211, e212 := rt.Post(ctx, t209, "И1 всюду: убрать можно только мусор", "Решить всё")
	if e212 != nil {
		return rt.Value{}, e212
	}
	if !t211 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «И1 всюду: убрать можно только мусор» ядра «Опись диска»")
	}
	return t208, nil
}

// I1Derzhitsya — функция flang «И1 держится».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр reshenie — «решение»: «Решение».
// Результат — значение.
func I1Derzhitsya(ctx *rt.Ctx, reshenie rt.Value) (rt.Value, error) {
	t213, e214 := rt.FieldGet(ctx, reshenie, "приговор")
	if e214 != nil {
		return rt.Value{}, e214
	}
	if rt.VariantIs(t213, "МожноУбрать") {
		t215, e216 := rt.FieldGet(ctx, reshenie, "разряд")
		if e216 != nil {
			return rt.Value{}, e216
		}
		if rt.VariantIs(t215, "Кэш") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t215, "Журнал") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t215, "Сборка") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t215, "Загрузка") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t215, "Крупное") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t215, "Неизвестное") {
			return rt.Flag(false), nil
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t215)
		}
	} else if rt.VariantIs(t213, "Спросить") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t213, "НеТрогать") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t213)
	}
}

// I1DerzhitsyaVsyudu — функция flang «И1 держится всюду».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр resheniya — «решения»: список: «Решение».
// Результат — значение.
func I1DerzhitsyaVsyudu(ctx *rt.Ctx, resheniya rt.Value) (rt.Value, error) {
	t217, e218 := rt.RequireList(ctx, resheniya, "свёртка")
	if e218 != nil {
		return rt.Value{}, e218
	}
	// «акк»
	akk := rt.Flag(true)
	for t219 := range t217 {
		// «решение»
		reshenie := t217[t219]
		t220, e221 := rt.Cond(ctx, akk)
		if e221 != nil {
			return rt.Value{}, e221
		}
		var t222 rt.Value
		if t220 {
			t223, e224 := I1Derzhitsya(ctx, reshenie)
			if e224 != nil {
				return rt.Value{}, e224
			}
			t222 = t223
		} else {
			t222 = rt.Flag(false)
		}
		akk = t222
	}
	return akk, nil
}

// PustoySvod — функция flang «Пустой свод».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: «Свод».
func PustoySvod(ctx *rt.Ctx) (rt.Value, error) {
	t225 := make([]rt.Field, 7)
	t225[0] = rt.Field{Name: "кэш", Value: rt.Number(0.0)}
	t225[1] = rt.Field{Name: "журнал", Value: rt.Number(0.0)}
	t225[2] = rt.Field{Name: "сборка", Value: rt.Number(0.0)}
	t225[3] = rt.Field{Name: "загрузка", Value: rt.Number(0.0)}
	t225[4] = rt.Field{Name: "крупное", Value: rt.Number(0.0)}
	t225[5] = rt.Field{Name: "неизвестное", Value: rt.Number(0.0)}
	t225[6] = rt.Field{Name: "освободить", Value: rt.Number(0.0)}
	return rt.Record(t225), nil
}

// PribavitReshenie — функция flang «Прибавить решение».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр svod — «свод»: «Свод».
// Параметр reshenie — «решение»: «Решение».
// Результат — значение: «Свод».
func PribavitReshenie(ctx *rt.Ctx, svod rt.Value, reshenie rt.Value) (rt.Value, error) {
	var t226 rt.Value
	t227, e228 := rt.FieldGet(ctx, reshenie, "приговор")
	if e228 != nil {
		return rt.Value{}, e228
	}
	if rt.VariantIs(t227, "МожноУбрать") {
		t229, e230 := rt.FieldGet(ctx, reshenie, "вес")
		if e230 != nil {
			return rt.Value{}, e230
		}
		t226 = t229
	} else if rt.VariantIs(t227, "Спросить") {
		t226 = rt.Number(0.0)
	} else if rt.VariantIs(t227, "НеТрогать") {
		t226 = rt.Number(0.0)
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t227)
	}
	// пусть «убрать»
	ubrat := t226
	t231, e232 := rt.FieldGet(ctx, reshenie, "разряд")
	if e232 != nil {
		return rt.Value{}, e232
	}
	if rt.VariantIs(t231, "Кэш") {
		t233, e234 := rt.FieldGet(ctx, svod, "кэш")
		if e234 != nil {
			return rt.Value{}, e234
		}
		t235, e236 := rt.FieldGet(ctx, reshenie, "вес")
		if e236 != nil {
			return rt.Value{}, e236
		}
		t237, e238 := rt.Add(ctx, t233, t235)
		if e238 != nil {
			return rt.Value{}, e238
		}
		t239, e240 := rt.FieldGet(ctx, svod, "журнал")
		if e240 != nil {
			return rt.Value{}, e240
		}
		t241, e242 := rt.FieldGet(ctx, svod, "сборка")
		if e242 != nil {
			return rt.Value{}, e242
		}
		t243, e244 := rt.FieldGet(ctx, svod, "загрузка")
		if e244 != nil {
			return rt.Value{}, e244
		}
		t245, e246 := rt.FieldGet(ctx, svod, "крупное")
		if e246 != nil {
			return rt.Value{}, e246
		}
		t247, e248 := rt.FieldGet(ctx, svod, "неизвестное")
		if e248 != nil {
			return rt.Value{}, e248
		}
		t249, e250 := rt.FieldGet(ctx, svod, "освободить")
		if e250 != nil {
			return rt.Value{}, e250
		}
		t251, e252 := rt.Add(ctx, t249, ubrat)
		if e252 != nil {
			return rt.Value{}, e252
		}
		t253 := make([]rt.Field, 7)
		t253[0] = rt.Field{Name: "кэш", Value: t237}
		t253[1] = rt.Field{Name: "журнал", Value: t239}
		t253[2] = rt.Field{Name: "сборка", Value: t241}
		t253[3] = rt.Field{Name: "загрузка", Value: t243}
		t253[4] = rt.Field{Name: "крупное", Value: t245}
		t253[5] = rt.Field{Name: "неизвестное", Value: t247}
		t253[6] = rt.Field{Name: "освободить", Value: t251}
		return rt.Record(t253), nil
	} else if rt.VariantIs(t231, "Журнал") {
		t254, e255 := rt.FieldGet(ctx, svod, "кэш")
		if e255 != nil {
			return rt.Value{}, e255
		}
		t256, e257 := rt.FieldGet(ctx, svod, "журнал")
		if e257 != nil {
			return rt.Value{}, e257
		}
		t258, e259 := rt.FieldGet(ctx, reshenie, "вес")
		if e259 != nil {
			return rt.Value{}, e259
		}
		t260, e261 := rt.Add(ctx, t256, t258)
		if e261 != nil {
			return rt.Value{}, e261
		}
		t262, e263 := rt.FieldGet(ctx, svod, "сборка")
		if e263 != nil {
			return rt.Value{}, e263
		}
		t264, e265 := rt.FieldGet(ctx, svod, "загрузка")
		if e265 != nil {
			return rt.Value{}, e265
		}
		t266, e267 := rt.FieldGet(ctx, svod, "крупное")
		if e267 != nil {
			return rt.Value{}, e267
		}
		t268, e269 := rt.FieldGet(ctx, svod, "неизвестное")
		if e269 != nil {
			return rt.Value{}, e269
		}
		t270, e271 := rt.FieldGet(ctx, svod, "освободить")
		if e271 != nil {
			return rt.Value{}, e271
		}
		t272, e273 := rt.Add(ctx, t270, ubrat)
		if e273 != nil {
			return rt.Value{}, e273
		}
		t274 := make([]rt.Field, 7)
		t274[0] = rt.Field{Name: "кэш", Value: t254}
		t274[1] = rt.Field{Name: "журнал", Value: t260}
		t274[2] = rt.Field{Name: "сборка", Value: t262}
		t274[3] = rt.Field{Name: "загрузка", Value: t264}
		t274[4] = rt.Field{Name: "крупное", Value: t266}
		t274[5] = rt.Field{Name: "неизвестное", Value: t268}
		t274[6] = rt.Field{Name: "освободить", Value: t272}
		return rt.Record(t274), nil
	} else if rt.VariantIs(t231, "Сборка") {
		t275, e276 := rt.FieldGet(ctx, svod, "кэш")
		if e276 != nil {
			return rt.Value{}, e276
		}
		t277, e278 := rt.FieldGet(ctx, svod, "журнал")
		if e278 != nil {
			return rt.Value{}, e278
		}
		t279, e280 := rt.FieldGet(ctx, svod, "сборка")
		if e280 != nil {
			return rt.Value{}, e280
		}
		t281, e282 := rt.FieldGet(ctx, reshenie, "вес")
		if e282 != nil {
			return rt.Value{}, e282
		}
		t283, e284 := rt.Add(ctx, t279, t281)
		if e284 != nil {
			return rt.Value{}, e284
		}
		t285, e286 := rt.FieldGet(ctx, svod, "загрузка")
		if e286 != nil {
			return rt.Value{}, e286
		}
		t287, e288 := rt.FieldGet(ctx, svod, "крупное")
		if e288 != nil {
			return rt.Value{}, e288
		}
		t289, e290 := rt.FieldGet(ctx, svod, "неизвестное")
		if e290 != nil {
			return rt.Value{}, e290
		}
		t291, e292 := rt.FieldGet(ctx, svod, "освободить")
		if e292 != nil {
			return rt.Value{}, e292
		}
		t293, e294 := rt.Add(ctx, t291, ubrat)
		if e294 != nil {
			return rt.Value{}, e294
		}
		t295 := make([]rt.Field, 7)
		t295[0] = rt.Field{Name: "кэш", Value: t275}
		t295[1] = rt.Field{Name: "журнал", Value: t277}
		t295[2] = rt.Field{Name: "сборка", Value: t283}
		t295[3] = rt.Field{Name: "загрузка", Value: t285}
		t295[4] = rt.Field{Name: "крупное", Value: t287}
		t295[5] = rt.Field{Name: "неизвестное", Value: t289}
		t295[6] = rt.Field{Name: "освободить", Value: t293}
		return rt.Record(t295), nil
	} else if rt.VariantIs(t231, "Загрузка") {
		t296, e297 := rt.FieldGet(ctx, svod, "кэш")
		if e297 != nil {
			return rt.Value{}, e297
		}
		t298, e299 := rt.FieldGet(ctx, svod, "журнал")
		if e299 != nil {
			return rt.Value{}, e299
		}
		t300, e301 := rt.FieldGet(ctx, svod, "сборка")
		if e301 != nil {
			return rt.Value{}, e301
		}
		t302, e303 := rt.FieldGet(ctx, svod, "загрузка")
		if e303 != nil {
			return rt.Value{}, e303
		}
		t304, e305 := rt.FieldGet(ctx, reshenie, "вес")
		if e305 != nil {
			return rt.Value{}, e305
		}
		t306, e307 := rt.Add(ctx, t302, t304)
		if e307 != nil {
			return rt.Value{}, e307
		}
		t308, e309 := rt.FieldGet(ctx, svod, "крупное")
		if e309 != nil {
			return rt.Value{}, e309
		}
		t310, e311 := rt.FieldGet(ctx, svod, "неизвестное")
		if e311 != nil {
			return rt.Value{}, e311
		}
		t312, e313 := rt.FieldGet(ctx, svod, "освободить")
		if e313 != nil {
			return rt.Value{}, e313
		}
		t314, e315 := rt.Add(ctx, t312, ubrat)
		if e315 != nil {
			return rt.Value{}, e315
		}
		t316 := make([]rt.Field, 7)
		t316[0] = rt.Field{Name: "кэш", Value: t296}
		t316[1] = rt.Field{Name: "журнал", Value: t298}
		t316[2] = rt.Field{Name: "сборка", Value: t300}
		t316[3] = rt.Field{Name: "загрузка", Value: t306}
		t316[4] = rt.Field{Name: "крупное", Value: t308}
		t316[5] = rt.Field{Name: "неизвестное", Value: t310}
		t316[6] = rt.Field{Name: "освободить", Value: t314}
		return rt.Record(t316), nil
	} else if rt.VariantIs(t231, "Крупное") {
		t317, e318 := rt.FieldGet(ctx, svod, "кэш")
		if e318 != nil {
			return rt.Value{}, e318
		}
		t319, e320 := rt.FieldGet(ctx, svod, "журнал")
		if e320 != nil {
			return rt.Value{}, e320
		}
		t321, e322 := rt.FieldGet(ctx, svod, "сборка")
		if e322 != nil {
			return rt.Value{}, e322
		}
		t323, e324 := rt.FieldGet(ctx, svod, "загрузка")
		if e324 != nil {
			return rt.Value{}, e324
		}
		t325, e326 := rt.FieldGet(ctx, svod, "крупное")
		if e326 != nil {
			return rt.Value{}, e326
		}
		t327, e328 := rt.FieldGet(ctx, reshenie, "вес")
		if e328 != nil {
			return rt.Value{}, e328
		}
		t329, e330 := rt.Add(ctx, t325, t327)
		if e330 != nil {
			return rt.Value{}, e330
		}
		t331, e332 := rt.FieldGet(ctx, svod, "неизвестное")
		if e332 != nil {
			return rt.Value{}, e332
		}
		t333, e334 := rt.FieldGet(ctx, svod, "освободить")
		if e334 != nil {
			return rt.Value{}, e334
		}
		t335, e336 := rt.Add(ctx, t333, ubrat)
		if e336 != nil {
			return rt.Value{}, e336
		}
		t337 := make([]rt.Field, 7)
		t337[0] = rt.Field{Name: "кэш", Value: t317}
		t337[1] = rt.Field{Name: "журнал", Value: t319}
		t337[2] = rt.Field{Name: "сборка", Value: t321}
		t337[3] = rt.Field{Name: "загрузка", Value: t323}
		t337[4] = rt.Field{Name: "крупное", Value: t329}
		t337[5] = rt.Field{Name: "неизвестное", Value: t331}
		t337[6] = rt.Field{Name: "освободить", Value: t335}
		return rt.Record(t337), nil
	} else if rt.VariantIs(t231, "Неизвестное") {
		t338, e339 := rt.FieldGet(ctx, svod, "кэш")
		if e339 != nil {
			return rt.Value{}, e339
		}
		t340, e341 := rt.FieldGet(ctx, svod, "журнал")
		if e341 != nil {
			return rt.Value{}, e341
		}
		t342, e343 := rt.FieldGet(ctx, svod, "сборка")
		if e343 != nil {
			return rt.Value{}, e343
		}
		t344, e345 := rt.FieldGet(ctx, svod, "загрузка")
		if e345 != nil {
			return rt.Value{}, e345
		}
		t346, e347 := rt.FieldGet(ctx, svod, "крупное")
		if e347 != nil {
			return rt.Value{}, e347
		}
		t348, e349 := rt.FieldGet(ctx, svod, "неизвестное")
		if e349 != nil {
			return rt.Value{}, e349
		}
		t350, e351 := rt.FieldGet(ctx, reshenie, "вес")
		if e351 != nil {
			return rt.Value{}, e351
		}
		t352, e353 := rt.Add(ctx, t348, t350)
		if e353 != nil {
			return rt.Value{}, e353
		}
		t354, e355 := rt.FieldGet(ctx, svod, "освободить")
		if e355 != nil {
			return rt.Value{}, e355
		}
		t356, e357 := rt.Add(ctx, t354, ubrat)
		if e357 != nil {
			return rt.Value{}, e357
		}
		t358 := make([]rt.Field, 7)
		t358[0] = rt.Field{Name: "кэш", Value: t338}
		t358[1] = rt.Field{Name: "журнал", Value: t340}
		t358[2] = rt.Field{Name: "сборка", Value: t342}
		t358[3] = rt.Field{Name: "загрузка", Value: t344}
		t358[4] = rt.Field{Name: "крупное", Value: t346}
		t358[5] = rt.Field{Name: "неизвестное", Value: t352}
		t358[6] = rt.Field{Name: "освободить", Value: t356}
		return rt.Record(t358), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t231)
	}
}

// Svesti — функция flang «Свести».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: «Свод».
func Svesti(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t359, e360 := ReshitVsyo(ctx, zapisi)
	if e360 != nil {
		return rt.Value{}, e360
	}
	t361, e362 := rt.RequireList(ctx, t359, "свёртка")
	if e362 != nil {
		return rt.Value{}, e362
	}
	t363, e364 := PustoySvod(ctx)
	if e364 != nil {
		return rt.Value{}, e364
	}
	// «свод»
	svod := t363
	for t365 := range t361 {
		// «решение»
		reshenie := t361[t365]
		t366, e367 := PribavitReshenie(ctx, svod, reshenie)
		if e367 != nil {
			return rt.Value{}, e367
		}
		svod = t366
	}
	t368 := svod
	t369, e370 := I2Derzhitsya(ctx, zapisi, t368)
	if e370 != nil {
		return rt.Value{}, e370
	}
	// постусловие «И2: освобождаемое не больше убираемого»
	t371, e372 := rt.Post(ctx, t369, "И2: освобождаемое не больше убираемого", "Свести")
	if e372 != nil {
		return rt.Value{}, e372
	}
	if !t371 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «И2: освобождаемое не больше убираемого» ядра «Опись диска»")
	}
	return t368, nil
}

// SummaRazmerovUbiraemyh — функция flang «Сумма размеров убираемых».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: число.
func SummaRazmerovUbiraemyh(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t373, e374 := rt.RequireList(ctx, zapisi, "свёртка")
	if e374 != nil {
		return rt.Value{}, e374
	}
	// «акк»
	akk := rt.Number(0.0)
	for t375 := range t373 {
		// «находка»
		nahodka := t373[t375]
		var t376 rt.Value
		t377, e378 := ReshitNahodku(ctx, nahodka)
		if e378 != nil {
			return rt.Value{}, e378
		}
		t379, e380 := rt.FieldGet(ctx, t377, "приговор")
		if e380 != nil {
			return rt.Value{}, e380
		}
		if rt.VariantIs(t379, "МожноУбрать") {
			t381, e382 := rt.FieldGet(ctx, nahodka, "размер")
			if e382 != nil {
				return rt.Value{}, e382
			}
			t383, e384 := rt.Add(ctx, akk, t381)
			if e384 != nil {
				return rt.Value{}, e384
			}
			t376 = t383
		} else if rt.VariantIs(t379, "Спросить") {
			t376 = akk
		} else if rt.VariantIs(t379, "НеТрогать") {
			t376 = akk
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t379)
		}
		akk = t376
	}
	return akk, nil
}

// I2Derzhitsya — функция flang «И2 держится».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Параметр svod — «свод»: «Свод».
// Результат — значение.
func I2Derzhitsya(ctx *rt.Ctx, zapisi rt.Value, svod rt.Value) (rt.Value, error) {
	t385, e386 := rt.FieldGet(ctx, svod, "освободить")
	if e386 != nil {
		return rt.Value{}, e386
	}
	t387, e388 := SummaRazmerovUbiraemyh(ctx, zapisi)
	if e388 != nil {
		return rt.Value{}, e388
	}
	t389, e390 := rt.Lte(ctx, t385, t387)
	if e390 != nil {
		return rt.Value{}, e390
	}
	return t389, nil
}

// StrokuOtchyota — функция flang «Строку отчёта».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение: «Строка отчёта».
func StrokuOtchyota(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t391, e392 := rt.FieldGet(ctx, nahodka, "путь")
	if e392 != nil {
		return rt.Value{}, e392
	}
	t393, e394 := ReshitNahodku(ctx, nahodka)
	if e394 != nil {
		return rt.Value{}, e394
	}
	t395 := make([]rt.Field, 2)
	t395[0] = rt.Field{Name: "путь", Value: t391}
	t395[1] = rt.Field{Name: "решение", Value: t393}
	return rt.Record(t395), nil
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
		t396 := make([]rt.Value, 1)
		t396[0] = stroka
		return rt.List(t396), nil
	} else if rt.ChainCons(stroki) {
		// голова «голова»
		golova := rt.ChainHead(stroki)
		// хвост «хвост»
		hvost := rt.ChainTail(stroki)
		t397, e398 := rt.FieldGet(ctx, stroka, "решение")
		if e398 != nil {
			return rt.Value{}, e398
		}
		t399, e400 := rt.FieldGet(ctx, t397, "вес")
		if e400 != nil {
			return rt.Value{}, e400
		}
		t401, e402 := rt.FieldGet(ctx, golova, "решение")
		if e402 != nil {
			return rt.Value{}, e402
		}
		t403, e404 := rt.FieldGet(ctx, t401, "вес")
		if e404 != nil {
			return rt.Value{}, e404
		}
		t405, e406 := rt.Gte(ctx, t399, t403)
		if e406 != nil {
			return rt.Value{}, e406
		}
		t407, e408 := rt.Cond(ctx, t405)
		if e408 != nil {
			return rt.Value{}, e408
		}
		if t407 {
			return PripisatStrokuOtchyota(ctx, stroka, stroki)
		} else {
			t409, e410 := VstavitPoVesu(ctx, stroka, hvost)
			if e410 != nil {
				return rt.Value{}, e410
			}
			return PripisatStrokuOtchyota(ctx, golova, t409)
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
	t411, e412 := rt.RequireList(ctx, stroki, "свёртка")
	if e412 != nil {
		return rt.Value{}, e412
	}
	t413 := make([]rt.Value, 1)
	t413[0] = pervaya
	// «акк»
	akk := rt.List(t413)
	for t414 := range t411 {
		// «эл»
		el := t411[t414]
		// «добавить»
		t415, e416 := rt.BAppend(ctx, el, akk)
		if e416 != nil {
			return rt.Value{}, e416
		}
		akk = t415
	}
	return akk, nil
}

// Otchyot — функция flang «Отчёт».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: список: «Строка отчёта».
func Otchyot(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t417, e418 := rt.RequireList(ctx, zapisi, "свёртка")
	if e418 != nil {
		return rt.Value{}, e418
	}
	// «акк»
	akk := rt.List(nil)
	for t419 := range t417 {
		// «находка»
		nahodka := t417[t419]
		t420, e421 := StrokuOtchyota(ctx, nahodka)
		if e421 != nil {
			return rt.Value{}, e421
		}
		t422, e423 := VstavitPoVesu(ctx, t420, akk)
		if e423 != nil {
			return rt.Value{}, e423
		}
		akk = t422
	}
	t424 := akk
	t425, e426 := OtchyotToyZheDliny(ctx, zapisi, t424)
	if e426 != nil {
		return rt.Value{}, e426
	}
	// постусловие «Отчёт той же длины»
	t427, e428 := rt.Post(ctx, t425, "Отчёт той же длины", "Отчёт")
	if e428 != nil {
		return rt.Value{}, e428
	}
	if !t427 {
		return rt.Value{}, rt.Fail("DIGITDISK_RULE", "%s", "нарушено постусловие «Отчёт той же длины» ядра «Опись диска»")
	}
	return t424, nil
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
	t429, e430 := rt.BLength(ctx, stroki)
	if e430 != nil {
		return rt.Value{}, e430
	}
	// «длина»
	t431, e432 := rt.BLength(ctx, zapisi)
	if e432 != nil {
		return rt.Value{}, e432
	}
	return rt.Flag(rt.Equal(t429, t431)), nil
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
	t433, e434 := EtoMozhnoUbrat(ctx, prigovor)
	if e434 != nil {
		return rt.Value{}, e434
	}
	t435, e436 := rt.Cond(ctx, t433)
	if e436 != nil {
		return rt.Value{}, e436
	}
	var t437 rt.Value
	if t435 {
		t437 = rt.Flag(false)
	} else {
		t437 = rt.Flag(true)
	}
	t438, e439 := rt.Cond(ctx, t437)
	if e439 != nil {
		return rt.Value{}, e439
	}
	if t438 {
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
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение.
func RazryadObosnovan(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value) (rt.Value, error) {
	t440, e441 := rt.FieldGet(ctx, nahodka, "путь")
	if e441 != nil {
		return rt.Value{}, e441
	}
	t442, e443 := PrimetaKesha(ctx, t440)
	if e443 != nil {
		return rt.Value{}, e443
	}
	// пусть «кэш»
	kesh := t442
	t444, e445 := rt.FieldGet(ctx, nahodka, "путь")
	if e445 != nil {
		return rt.Value{}, e445
	}
	t446, e447 := PrimetaZhurnala(ctx, t444)
	if e447 != nil {
		return rt.Value{}, e447
	}
	// пусть «журнал»
	zhurnal := t446
	t448, e449 := rt.FieldGet(ctx, nahodka, "путь")
	if e449 != nil {
		return rt.Value{}, e449
	}
	t450, e451 := PrimetaSborki(ctx, t448)
	if e451 != nil {
		return rt.Value{}, e451
	}
	// пусть «сборка»
	sborka := t450
	t452, e453 := rt.FieldGet(ctx, nahodka, "путь")
	if e453 != nil {
		return rt.Value{}, e453
	}
	t454, e455 := PrimetaZagruzki(ctx, t452)
	if e455 != nil {
		return rt.Value{}, e455
	}
	// пусть «загрузка»
	zagruzka := t454
	t456, e457 := rt.FieldGet(ctx, nahodka, "размер")
	if e457 != nil {
		return rt.Value{}, e457
	}
	t458, e459 := PorogKrupnogo(ctx)
	if e459 != nil {
		return rt.Value{}, e459
	}
	t460, e461 := rt.Gte(ctx, t456, t458)
	if e461 != nil {
		return rt.Value{}, e461
	}
	// пусть «крупное»
	krupnoe := t460
	if rt.VariantIs(razryad, "Кэш") {
		return kesh, nil
	} else if rt.VariantIs(razryad, "Журнал") {
		t462, e463 := rt.Cond(ctx, zhurnal)
		if e463 != nil {
			return rt.Value{}, e463
		}
		if t462 {
			t464, e465 := rt.Cond(ctx, kesh)
			if e465 != nil {
				return rt.Value{}, e465
			}
			if t464 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Сборка") {
		t466, e467 := rt.Cond(ctx, sborka)
		if e467 != nil {
			return rt.Value{}, e467
		}
		var t468 rt.Value
		if t466 {
			t469, e470 := rt.Cond(ctx, kesh)
			if e470 != nil {
				return rt.Value{}, e470
			}
			var t471 rt.Value
			if t469 {
				t471 = rt.Flag(false)
			} else {
				t471 = rt.Flag(true)
			}
			t468 = t471
		} else {
			t468 = rt.Flag(false)
		}
		t472, e473 := rt.Cond(ctx, t468)
		if e473 != nil {
			return rt.Value{}, e473
		}
		if t472 {
			t474, e475 := rt.Cond(ctx, zhurnal)
			if e475 != nil {
				return rt.Value{}, e475
			}
			if t474 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Загрузка") {
		t476, e477 := rt.Cond(ctx, zagruzka)
		if e477 != nil {
			return rt.Value{}, e477
		}
		var t478 rt.Value
		if t476 {
			t479, e480 := rt.Cond(ctx, kesh)
			if e480 != nil {
				return rt.Value{}, e480
			}
			var t481 rt.Value
			if t479 {
				t481 = rt.Flag(false)
			} else {
				t481 = rt.Flag(true)
			}
			t478 = t481
		} else {
			t478 = rt.Flag(false)
		}
		t482, e483 := rt.Cond(ctx, t478)
		if e483 != nil {
			return rt.Value{}, e483
		}
		var t484 rt.Value
		if t482 {
			t485, e486 := rt.Cond(ctx, zhurnal)
			if e486 != nil {
				return rt.Value{}, e486
			}
			var t487 rt.Value
			if t485 {
				t487 = rt.Flag(false)
			} else {
				t487 = rt.Flag(true)
			}
			t484 = t487
		} else {
			t484 = rt.Flag(false)
		}
		t488, e489 := rt.Cond(ctx, t484)
		if e489 != nil {
			return rt.Value{}, e489
		}
		if t488 {
			t490, e491 := rt.Cond(ctx, sborka)
			if e491 != nil {
				return rt.Value{}, e491
			}
			if t490 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Крупное") {
		t492, e493 := rt.Cond(ctx, krupnoe)
		if e493 != nil {
			return rt.Value{}, e493
		}
		var t494 rt.Value
		if t492 {
			t495, e496 := rt.Cond(ctx, kesh)
			if e496 != nil {
				return rt.Value{}, e496
			}
			var t497 rt.Value
			if t495 {
				t497 = rt.Flag(false)
			} else {
				t497 = rt.Flag(true)
			}
			t494 = t497
		} else {
			t494 = rt.Flag(false)
		}
		t498, e499 := rt.Cond(ctx, t494)
		if e499 != nil {
			return rt.Value{}, e499
		}
		var t500 rt.Value
		if t498 {
			t501, e502 := rt.Cond(ctx, zhurnal)
			if e502 != nil {
				return rt.Value{}, e502
			}
			var t503 rt.Value
			if t501 {
				t503 = rt.Flag(false)
			} else {
				t503 = rt.Flag(true)
			}
			t500 = t503
		} else {
			t500 = rt.Flag(false)
		}
		t504, e505 := rt.Cond(ctx, t500)
		if e505 != nil {
			return rt.Value{}, e505
		}
		var t506 rt.Value
		if t504 {
			t507, e508 := rt.Cond(ctx, sborka)
			if e508 != nil {
				return rt.Value{}, e508
			}
			var t509 rt.Value
			if t507 {
				t509 = rt.Flag(false)
			} else {
				t509 = rt.Flag(true)
			}
			t506 = t509
		} else {
			t506 = rt.Flag(false)
		}
		t510, e511 := rt.Cond(ctx, t506)
		if e511 != nil {
			return rt.Value{}, e511
		}
		if t510 {
			t512, e513 := rt.Cond(ctx, zagruzka)
			if e513 != nil {
				return rt.Value{}, e513
			}
			if t512 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Неизвестное") {
		t514, e515 := rt.Cond(ctx, kesh)
		if e515 != nil {
			return rt.Value{}, e515
		}
		var t516 rt.Value
		if t514 {
			t516 = rt.Flag(false)
		} else {
			t516 = rt.Flag(true)
		}
		t517, e518 := rt.Cond(ctx, t516)
		if e518 != nil {
			return rt.Value{}, e518
		}
		var t519 rt.Value
		if t517 {
			t520, e521 := rt.Cond(ctx, zhurnal)
			if e521 != nil {
				return rt.Value{}, e521
			}
			var t522 rt.Value
			if t520 {
				t522 = rt.Flag(false)
			} else {
				t522 = rt.Flag(true)
			}
			t519 = t522
		} else {
			t519 = rt.Flag(false)
		}
		t523, e524 := rt.Cond(ctx, t519)
		if e524 != nil {
			return rt.Value{}, e524
		}
		var t525 rt.Value
		if t523 {
			t526, e527 := rt.Cond(ctx, sborka)
			if e527 != nil {
				return rt.Value{}, e527
			}
			var t528 rt.Value
			if t526 {
				t528 = rt.Flag(false)
			} else {
				t528 = rt.Flag(true)
			}
			t525 = t528
		} else {
			t525 = rt.Flag(false)
		}
		t529, e530 := rt.Cond(ctx, t525)
		if e530 != nil {
			return rt.Value{}, e530
		}
		var t531 rt.Value
		if t529 {
			t532, e533 := rt.Cond(ctx, zagruzka)
			if e533 != nil {
				return rt.Value{}, e533
			}
			var t534 rt.Value
			if t532 {
				t534 = rt.Flag(false)
			} else {
				t534 = rt.Flag(true)
			}
			t531 = t534
		} else {
			t531 = rt.Flag(false)
		}
		t535, e536 := rt.Cond(ctx, t531)
		if e536 != nil {
			return rt.Value{}, e536
		}
		if t535 {
			t537, e538 := rt.Cond(ctx, krupnoe)
			if e538 != nil {
				return rt.Value{}, e538
			}
			if t537 {
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
	t539, e540 := rt.FieldGet(ctx, nahodka, "доступен")
	if e540 != nil {
		return rt.Value{}, e540
	}
	t541, e542 := rt.Cond(ctx, t539)
	if e542 != nil {
		return rt.Value{}, e542
	}
	var t543 rt.Value
	if t541 {
		t543 = rt.Flag(false)
	} else {
		t543 = rt.Flag(true)
	}
	t544, e545 := rt.Cond(ctx, t543)
	if e545 != nil {
		return rt.Value{}, e545
	}
	if t544 {
		return EtoNeTrogat(ctx, prigovor)
	} else {
		t546, e547 := Ssylka(ctx, nahodka)
		if e547 != nil {
			return rt.Value{}, e547
		}
		t548, e549 := rt.Cond(ctx, t546)
		if e549 != nil {
			return rt.Value{}, e549
		}
		if t548 {
			return EtoNeTrogat(ctx, prigovor)
		} else {
			t550, e551 := EtoMozhnoUbrat(ctx, prigovor)
			if e551 != nil {
				return rt.Value{}, e551
			}
			t552, e553 := rt.Cond(ctx, t550)
			if e553 != nil {
				return rt.Value{}, e553
			}
			if t552 {
				t554, e555 := I1NaPare(ctx, razryad, prigovor)
				if e555 != nil {
					return rt.Value{}, e555
				}
				t556, e557 := rt.Cond(ctx, t554)
				if e557 != nil {
					return rt.Value{}, e557
				}
				var t558 rt.Value
				if t556 {
					t559, e560 := Katalog(ctx, nahodka)
					if e560 != nil {
						return rt.Value{}, e560
					}
					t561, e562 := rt.Cond(ctx, t559)
					if e562 != nil {
						return rt.Value{}, e562
					}
					var t563 rt.Value
					if t561 {
						t563 = rt.Flag(false)
					} else {
						t563 = rt.Flag(true)
					}
					t558 = t563
				} else {
					t558 = rt.Flag(false)
				}
				t564, e565 := rt.Cond(ctx, t558)
				if e565 != nil {
					return rt.Value{}, e565
				}
				if t564 {
					t566, e567 := rt.FieldGet(ctx, nahodka, "возраст_дней")
					if e567 != nil {
						return rt.Value{}, e567
					}
					t568, e569 := PorogRazryada(ctx, razryad)
					if e569 != nil {
						return rt.Value{}, e569
					}
					t570, e571 := rt.Gte(ctx, t566, t568)
					if e571 != nil {
						return rt.Value{}, e571
					}
					return t570, nil
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
	t572, e573 := EtoNeTrogat(ctx, prigovor)
	if e573 != nil {
		return rt.Value{}, e573
	}
	t574, e575 := rt.Cond(ctx, t572)
	if e575 != nil {
		return rt.Value{}, e575
	}
	if t574 {
		return rt.Flag(rt.Equal(ves, rt.Number(0.0))), nil
	} else {
		t576, e577 := rt.FieldGet(ctx, nahodka, "размер")
		if e577 != nil {
			return rt.Value{}, e577
		}
		t578, e579 := rt.Cond(ctx, rt.Flag(rt.Equal(ves, t576)))
		if e579 != nil {
			return rt.Value{}, e579
		}
		if t578 {
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
	case "Разряд находки":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд находки", 1, len(args))
		}
		return RazryadNahodki(ctx, args[0])
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
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Решить находку", 1, len(args))
		}
		return ReshitNahodku(ctx, args[0])
	case "Решить всё":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Решить всё", 1, len(args))
		}
		return ReshitVsyo(ctx, args[0])
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
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Свести", 1, len(args))
		}
		return Svesti(ctx, args[0])
	case "Сумма размеров убираемых":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Сумма размеров убираемых", 1, len(args))
		}
		return SummaRazmerovUbiraemyh(ctx, args[0])
	case "И2 держится":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"И2 держится", 2, len(args))
		}
		return I2Derzhitsya(ctx, args[0], args[1])
	case "Строку отчёта":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Строку отчёта", 1, len(args))
		}
		return StrokuOtchyota(ctx, args[0])
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
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Отчёт", 1, len(args))
		}
		return Otchyot(ctx, args[0])
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
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд обоснован", 2, len(args))
		}
		return RazryadObosnovan(ctx, args[0], args[1])
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
