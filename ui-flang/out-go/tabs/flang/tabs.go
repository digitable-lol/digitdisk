// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
// Модуль flang: «Tabs».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flangtabs/flangrt"
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

// SozdatVkladki — запись FTS «Вкладки»: «открыта», «сколько».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatVkladki(otkryta rt.Value, skolko rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "открыта", Value: otkryta},
		{Name: "сколько", Value: skolko},
	})
}

// VariantSleduyuschaya — вариант «Следующая» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantSleduyuschaya() rt.Value {
	return rt.Variant("Следующая", nil)
}

// VariantPredyduschaya — вариант «Предыдущая» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantPredyduschaya() rt.Value {
	return rt.Variant("Предыдущая", nil)
}

// VariantNomerom — вариант «Номером» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantNomerom(nomer rt.Value) rt.Value {
	return rt.Variant("Номером", []rt.Field{
		{Name: "номер", Value: nomer},
	})
}

// VariantPervaya — вариант «Первая» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantPervaya() rt.Value {
	return rt.Variant("Первая", nil)
}

// VariantPoslednyaya — вариант «Последняя» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantPoslednyaya() rt.Value {
	return rt.Variant("Последняя", nil)
}

// VariantMimo — вариант «Мимо» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantMimo() rt.Value {
	return rt.Variant("Мимо", nil)
}

// SkolkoHotyaByOdna — функция flang «Сколько хотя бы одна».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Результат — значение: число.
func SkolkoHotyaByOdna(ctx *rt.Ctx, skolko rt.Value) (rt.Value, error) {
	t1, e2 := rt.Lt(ctx, skolko, rt.Number(1.0))
	if e2 != nil {
		return rt.Value{}, e2
	}
	t3, e4 := rt.Cond(ctx, t1)
	if e4 != nil {
		return rt.Value{}, e4
	}
	var t5 rt.Value
	if t3 {
		t5 = rt.Number(1.0)
	} else {
		t5 = skolko
	}
	t6 := t5
	t7, e8 := rt.Gte(ctx, t6, rt.Number(1.0))
	if e8 != nil {
		return rt.Value{}, e8
	}
	// постусловие «вкладок не меньше одной»
	t9, e10 := rt.Post(ctx, t7, "вкладок не меньше одной", "Сколько хотя бы одна")
	if e10 != nil {
		return rt.Value{}, e10
	}
	if !t9 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «вкладок не меньше одной» функции «Сколько хотя бы одна»")
	}
	return t6, nil
}

// VkladkaVPredelah — функция flang «Вкладка в пределах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр otkryta — «открыта»: число.
// Параметр skolko — «сколько»: число.
// Результат — значение: число.
func VkladkaVPredelah(ctx *rt.Ctx, otkryta rt.Value, skolko rt.Value) (rt.Value, error) {
	t11, e12 := SkolkoHotyaByOdna(ctx, skolko)
	if e12 != nil {
		return rt.Value{}, e12
	}
	// пусть «предел»
	predel := t11
	t13, e14 := rt.Lt(ctx, otkryta, rt.Number(0.0))
	if e14 != nil {
		return rt.Value{}, e14
	}
	t15, e16 := rt.Cond(ctx, t13)
	if e16 != nil {
		return rt.Value{}, e16
	}
	var t17 rt.Value
	if t15 {
		t17 = rt.Number(0.0)
	} else {
		t18, e19 := rt.Sub(ctx, predel, rt.Number(1.0))
		if e19 != nil {
			return rt.Value{}, e19
		}
		t20, e21 := rt.Gt(ctx, otkryta, t18)
		if e21 != nil {
			return rt.Value{}, e21
		}
		t22, e23 := rt.Cond(ctx, t20)
		if e23 != nil {
			return rt.Value{}, e23
		}
		var t24 rt.Value
		if t22 {
			t25, e26 := rt.Sub(ctx, predel, rt.Number(1.0))
			if e26 != nil {
				return rt.Value{}, e26
			}
			t24 = t25
		} else {
			t24 = otkryta
		}
		t17 = t24
	}
	t27 := t17
	t28, e29 := rt.Gte(ctx, t27, rt.Number(0.0))
	if e29 != nil {
		return rt.Value{}, e29
	}
	// постусловие «вкладка не отрицательна»
	t30, e31 := rt.Post(ctx, t28, "вкладка не отрицательна", "Вкладка в пределах")
	if e31 != nil {
		return rt.Value{}, e31
	}
	if !t30 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «вкладка не отрицательна» функции «Вкладка в пределах»")
	}
	t32, e33 := SkolkoHotyaByOdna(ctx, skolko)
	if e33 != nil {
		return rt.Value{}, e33
	}
	t34, e35 := rt.Lt(ctx, t27, t32)
	if e35 != nil {
		return rt.Value{}, e35
	}
	// постусловие «вкладка внутри набора»
	t36, e37 := rt.Post(ctx, t34, "вкладка внутри набора", "Вкладка в пределах")
	if e37 != nil {
		return rt.Value{}, e37
	}
	if !t36 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «вкладка внутри набора» функции «Вкладка в пределах»")
	}
	return t27, nil
}

// Nabor — функция flang «Набор».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр otkryta — «открыта»: число.
// Параметр skolko — «сколько»: число.
// Результат — значение: «Вкладки».
func Nabor(ctx *rt.Ctx, otkryta rt.Value, skolko rt.Value) (rt.Value, error) {
	t38, e39 := VkladkaVPredelah(ctx, otkryta, skolko)
	if e39 != nil {
		return rt.Value{}, e39
	}
	t40, e41 := SkolkoHotyaByOdna(ctx, skolko)
	if e41 != nil {
		return rt.Value{}, e41
	}
	t42 := make([]rt.Field, 2)
	t42[0] = rt.Field{Name: "открыта", Value: t38}
	t42[1] = rt.Field{Name: "сколько", Value: t40}
	t43 := rt.Record(t42)
	t44, e45 := rt.FieldGet(ctx, t43, "открыта")
	if e45 != nil {
		return rt.Value{}, e45
	}
	t46, e47 := rt.FieldGet(ctx, t43, "сколько")
	if e47 != nil {
		return rt.Value{}, e47
	}
	t48, e49 := rt.Lt(ctx, t44, t46)
	if e49 != nil {
		return rt.Value{}, e49
	}
	// постусловие «открытая вкладка внутри набора»
	t50, e51 := rt.Post(ctx, t48, "открытая вкладка внутри набора", "Набор")
	if e51 != nil {
		return rt.Value{}, e51
	}
	if !t50 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «открытая вкладка внутри набора» функции «Набор»")
	}
	t52, e53 := rt.FieldGet(ctx, t43, "открыта")
	if e53 != nil {
		return rt.Value{}, e53
	}
	t54, e55 := rt.Gte(ctx, t52, rt.Number(0.0))
	if e55 != nil {
		return rt.Value{}, e55
	}
	// постусловие «открытая вкладка не отрицательна»
	t56, e57 := rt.Post(ctx, t54, "открытая вкладка не отрицательна", "Набор")
	if e57 != nil {
		return rt.Value{}, e57
	}
	if !t56 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «открытая вкладка не отрицательна» функции «Набор»")
	}
	return t43, nil
}

// Pereklyuchit — функция flang «Переключить».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр vkladki — «вкладки»: «Вкладки».
// Параметр nazhatie — «нажатие»: «Нажатие».
// Результат — значение: «Вкладки».
func Pereklyuchit(ctx *rt.Ctx, vkladki rt.Value, nazhatie rt.Value) (rt.Value, error) {
	t58, e59 := rt.FieldGet(ctx, vkladki, "сколько")
	if e59 != nil {
		return rt.Value{}, e59
	}
	t60, e61 := SkolkoHotyaByOdna(ctx, t58)
	if e61 != nil {
		return rt.Value{}, e61
	}
	// пусть «предел»
	predel := t60
	t62, e63 := rt.FieldGet(ctx, vkladki, "открыта")
	if e63 != nil {
		return rt.Value{}, e63
	}
	t64, e65 := VkladkaVPredelah(ctx, t62, predel)
	if e65 != nil {
		return rt.Value{}, e65
	}
	// пусть «теперь»
	teper := t64
	var t66 rt.Value
	if rt.VariantIs(nazhatie, "Следующая") {
		t67, e68 := rt.Add(ctx, teper, rt.Number(1.0))
		if e68 != nil {
			return rt.Value{}, e68
		}
		t69, e70 := rt.Mod(ctx, t67, predel)
		if e70 != nil {
			return rt.Value{}, e70
		}
		t71, e72 := Nabor(ctx, t69, predel)
		if e72 != nil {
			return rt.Value{}, e72
		}
		t66 = t71
	} else if rt.VariantIs(nazhatie, "Предыдущая") {
		t73, e74 := rt.Add(ctx, teper, predel)
		if e74 != nil {
			return rt.Value{}, e74
		}
		t75, e76 := rt.Sub(ctx, t73, rt.Number(1.0))
		if e76 != nil {
			return rt.Value{}, e76
		}
		t77, e78 := rt.Mod(ctx, t75, predel)
		if e78 != nil {
			return rt.Value{}, e78
		}
		t79, e80 := Nabor(ctx, t77, predel)
		if e80 != nil {
			return rt.Value{}, e80
		}
		t66 = t79
	} else if rt.VariantIs(nazhatie, "Номером") {
		// поле «номер»
		nomer, e81 := rt.VariantField(ctx, nazhatie, "номер")
		if e81 != nil {
			return rt.Value{}, e81
		}
		t82, e83 := rt.Gte(ctx, nomer, rt.Number(1.0))
		if e83 != nil {
			return rt.Value{}, e83
		}
		t84, e85 := rt.Cond(ctx, t82)
		if e85 != nil {
			return rt.Value{}, e85
		}
		var t86 rt.Value
		if t84 {
			t87, e88 := rt.Sub(ctx, nomer, rt.Number(1.0))
			if e88 != nil {
				return rt.Value{}, e88
			}
			t89, e90 := rt.Lt(ctx, t87, predel)
			if e90 != nil {
				return rt.Value{}, e90
			}
			t86 = t89
		} else {
			t86 = rt.Flag(false)
		}
		t91, e92 := rt.Cond(ctx, t86)
		if e92 != nil {
			return rt.Value{}, e92
		}
		var t93 rt.Value
		if t91 {
			t94, e95 := rt.Sub(ctx, nomer, rt.Number(1.0))
			if e95 != nil {
				return rt.Value{}, e95
			}
			t96, e97 := Nabor(ctx, t94, predel)
			if e97 != nil {
				return rt.Value{}, e97
			}
			t93 = t96
		} else {
			t98, e99 := Nabor(ctx, teper, predel)
			if e99 != nil {
				return rt.Value{}, e99
			}
			t93 = t98
		}
		t66 = t93
	} else if rt.VariantIs(nazhatie, "Первая") {
		t100, e101 := Nabor(ctx, rt.Number(0.0), predel)
		if e101 != nil {
			return rt.Value{}, e101
		}
		t66 = t100
	} else if rt.VariantIs(nazhatie, "Последняя") {
		t102, e103 := rt.Sub(ctx, predel, rt.Number(1.0))
		if e103 != nil {
			return rt.Value{}, e103
		}
		t104, e105 := Nabor(ctx, t102, predel)
		if e105 != nil {
			return rt.Value{}, e105
		}
		t66 = t104
	} else if rt.VariantIs(nazhatie, "Мимо") {
		t106, e107 := Nabor(ctx, teper, predel)
		if e107 != nil {
			return rt.Value{}, e107
		}
		t66 = t106
	} else {
		return rt.Value{}, rt.MatchFail(ctx, nazhatie)
	}
	t108 := t66
	t109, e110 := rt.FieldGet(ctx, t108, "сколько")
	if e110 != nil {
		return rt.Value{}, e110
	}
	t111, e112 := rt.FieldGet(ctx, vkladki, "сколько")
	if e112 != nil {
		return rt.Value{}, e112
	}
	t113, e114 := SkolkoHotyaByOdna(ctx, t111)
	if e114 != nil {
		return rt.Value{}, e114
	}
	// постусловие «вкладок столько же»
	t115, e116 := rt.Post(ctx, rt.Flag(rt.Equal(t109, t113)), "вкладок столько же", "Переключить")
	if e116 != nil {
		return rt.Value{}, e116
	}
	if !t115 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «вкладок столько же» функции «Переключить»")
	}
	t117, e118 := rt.FieldGet(ctx, t108, "открыта")
	if e118 != nil {
		return rt.Value{}, e118
	}
	t119, e120 := rt.FieldGet(ctx, t108, "сколько")
	if e120 != nil {
		return rt.Value{}, e120
	}
	t121, e122 := rt.Lt(ctx, t117, t119)
	if e122 != nil {
		return rt.Value{}, e122
	}
	// постусловие «открытая вкладка внутри набора»
	t123, e124 := rt.Post(ctx, t121, "открытая вкладка внутри набора", "Переключить")
	if e124 != nil {
		return rt.Value{}, e124
	}
	if !t123 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «открытая вкладка внутри набора» функции «Переключить»")
	}
	t125, e126 := rt.FieldGet(ctx, t108, "открыта")
	if e126 != nil {
		return rt.Value{}, e126
	}
	t127, e128 := rt.Gte(ctx, t125, rt.Number(0.0))
	if e128 != nil {
		return rt.Value{}, e128
	}
	// постусловие «открытая вкладка не отрицательна»
	t129, e130 := rt.Post(ctx, t127, "открытая вкладка не отрицательна", "Переключить")
	if e130 != nil {
		return rt.Value{}, e130
	}
	if !t129 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «открытая вкладка не отрицательна» функции «Переключить»")
	}
	return t108, nil
}

// ProkrutkuSbrosit — функция flang «Прокрутку сбросить».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nazhatie — «нажатие»: «Нажатие».
// Результат — значение.
func ProkrutkuSbrosit(ctx *rt.Ctx, nazhatie rt.Value) (rt.Value, error) {
	if rt.VariantIs(nazhatie, "Мимо") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(nazhatie, "Следующая") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(nazhatie, "Предыдущая") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(nazhatie, "Номером") {
		// поле «номер»
		nomer, e131 := rt.VariantField(ctx, nazhatie, "номер")
		if e131 != nil {
			return rt.Value{}, e131
		}
		_ = nomer
		return rt.Flag(true), nil
	} else if rt.VariantIs(nazhatie, "Первая") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(nazhatie, "Последняя") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, nazhatie)
	}
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Сколько хотя бы одна":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Сколько хотя бы одна", 1, len(args))
		}
		return SkolkoHotyaByOdna(ctx, args[0])
	case "Вкладка в пределах":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Вкладка в пределах", 2, len(args))
		}
		return VkladkaVPredelah(ctx, args[0], args[1])
	case "Набор":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Набор", 2, len(args))
		}
		return Nabor(ctx, args[0], args[1])
	case "Переключить":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Переключить", 2, len(args))
		}
		return Pereklyuchit(ctx, args[0], args[1])
	case "Прокрутку сбросить":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прокрутку сбросить", 1, len(args))
		}
		return ProkrutkuSbrosit(ctx, args[0])
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
