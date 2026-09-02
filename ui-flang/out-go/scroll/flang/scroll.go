// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
// Модуль flang: «Scroll».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flangscroll/flangrt"
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

// VariantStrokoyVniz — вариант «Строкой вниз» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantStrokoyVniz() rt.Value {
	return rt.Variant("Строкой вниз", nil)
}

// VariantStrokoyVverh — вариант «Строкой вверх» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantStrokoyVverh() rt.Value {
	return rt.Variant("Строкой вверх", nil)
}

// VariantStraniceyVniz — вариант «Страницей вниз» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantStraniceyVniz() rt.Value {
	return rt.Variant("Страницей вниз", nil)
}

// VariantStraniceyVverh — вариант «Страницей вверх» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantStraniceyVverh() rt.Value {
	return rt.Variant("Страницей вверх", nil)
}

// VariantVNachalo — вариант «В начало» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantVNachalo() rt.Value {
	return rt.Variant("В начало", nil)
}

// VariantVKonec — вариант «В конец» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantVKonec() rt.Value {
	return rt.Variant("В конец", nil)
}

// VariantMimo — вариант «Мимо» суммы типов «Нажатие».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantMimo() rt.Value {
	return rt.Variant("Мимо", nil)
}

// VysotaTela — функция flang «Высота тела».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр strok — «строк»: число.
// Параметр obvyazka — «обвязка»: число.
// Результат — значение: число.
func VysotaTela(ctx *rt.Ctx, strok rt.Value, obvyazka rt.Value) (rt.Value, error) {
	t1, e2 := rt.Sub(ctx, strok, obvyazka)
	if e2 != nil {
		return rt.Value{}, e2
	}
	t3, e4 := rt.Lt(ctx, t1, rt.Number(1.0))
	if e4 != nil {
		return rt.Value{}, e4
	}
	t5, e6 := rt.Cond(ctx, t3)
	if e6 != nil {
		return rt.Value{}, e6
	}
	var t7 rt.Value
	if t5 {
		t7 = rt.Number(1.0)
	} else {
		t8, e9 := rt.Sub(ctx, strok, obvyazka)
		if e9 != nil {
			return rt.Value{}, e9
		}
		t7 = t8
	}
	t10 := t7
	t11, e12 := rt.Gte(ctx, t10, rt.Number(1.0))
	if e12 != nil {
		return rt.Value{}, e12
	}
	// постусловие «высота тела не меньше одной строки»
	t13, e14 := rt.Post(ctx, t11, "высота тела не меньше одной строки", "Высота тела")
	if e14 != nil {
		return rt.Value{}, e14
	}
	if !t13 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «высота тела не меньше одной строки» функции «Высота тела»")
	}
	return t10, nil
}

// ProkrutkaVPredelah — функция flang «Прокрутка в пределах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prokrutka — «прокрутка»: число.
// Параметр dlina — «длина»: число.
// Результат — значение: число.
func ProkrutkaVPredelah(ctx *rt.Ctx, prokrutka rt.Value, dlina rt.Value) (rt.Value, error) {
	t15, e16 := rt.Sub(ctx, dlina, rt.Number(1.0))
	if e16 != nil {
		return rt.Value{}, e16
	}
	t17, e18 := rt.Gt(ctx, prokrutka, t15)
	if e18 != nil {
		return rt.Value{}, e18
	}
	t19, e20 := rt.Cond(ctx, t17)
	if e20 != nil {
		return rt.Value{}, e20
	}
	var t21 rt.Value
	if t19 {
		t22, e23 := rt.Sub(ctx, dlina, rt.Number(1.0))
		if e23 != nil {
			return rt.Value{}, e23
		}
		t21 = t22
	} else {
		t21 = prokrutka
	}
	// пусть «дно»
	dno := t21
	t24, e25 := rt.Lt(ctx, dno, rt.Number(0.0))
	if e25 != nil {
		return rt.Value{}, e25
	}
	t26, e27 := rt.Cond(ctx, t24)
	if e27 != nil {
		return rt.Value{}, e27
	}
	var t28 rt.Value
	if t26 {
		t28 = rt.Number(0.0)
	} else {
		t28 = dno
	}
	t29 := t28
	t30, e31 := rt.Gte(ctx, t29, rt.Number(0.0))
	if e31 != nil {
		return rt.Value{}, e31
	}
	// постусловие «прокрутка не отрицательна»
	t32, e33 := rt.Post(ctx, t30, "прокрутка не отрицательна", "Прокрутка в пределах")
	if e33 != nil {
		return rt.Value{}, e33
	}
	if !t32 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «прокрутка не отрицательна» функции «Прокрутка в пределах»")
	}
	t34, e35 := rt.Lte(ctx, dlina, rt.Number(0.0))
	if e35 != nil {
		return rt.Value{}, e35
	}
	t36, e37 := rt.Cond(ctx, t34)
	if e37 != nil {
		return rt.Value{}, e37
	}
	var t38 rt.Value
	if t36 {
		t38 = rt.Number(0.0)
	} else {
		t39, e40 := rt.Sub(ctx, dlina, rt.Number(1.0))
		if e40 != nil {
			return rt.Value{}, e40
		}
		t38 = t39
	}
	t41, e42 := rt.Lte(ctx, t29, t38)
	if e42 != nil {
		return rt.Value{}, e42
	}
	// постусловие «прокрутка не заходит за последнюю строку»
	t43, e44 := rt.Post(ctx, t41, "прокрутка не заходит за последнюю строку", "Прокрутка в пределах")
	if e44 != nil {
		return rt.Value{}, e44
	}
	if !t43 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «прокрутка не заходит за последнюю строку» функции «Прокрутка в пределах»")
	}
	return t29, nil
}

// ShagProkrutki — функция flang «Шаг прокрутки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prokrutka — «прокрутка»: число.
// Параметр vysota — «высота»: число.
// Параметр nazhatie — «нажатие»: «Нажатие».
// Параметр dlina — «длина»: число.
// Результат — значение: число.
func ShagProkrutki(ctx *rt.Ctx, prokrutka rt.Value, vysota rt.Value, nazhatie rt.Value, dlina rt.Value) (rt.Value, error) {
	if rt.VariantIs(nazhatie, "Строкой вниз") {
		t45, e46 := rt.Add(ctx, prokrutka, rt.Number(1.0))
		if e46 != nil {
			return rt.Value{}, e46
		}
		return t45, nil
	} else if rt.VariantIs(nazhatie, "Строкой вверх") {
		t47, e48 := rt.Sub(ctx, prokrutka, rt.Number(1.0))
		if e48 != nil {
			return rt.Value{}, e48
		}
		return t47, nil
	} else if rt.VariantIs(nazhatie, "Страницей вниз") {
		t49, e50 := rt.Add(ctx, prokrutka, vysota)
		if e50 != nil {
			return rt.Value{}, e50
		}
		return t49, nil
	} else if rt.VariantIs(nazhatie, "Страницей вверх") {
		t51, e52 := rt.Sub(ctx, prokrutka, vysota)
		if e52 != nil {
			return rt.Value{}, e52
		}
		return t51, nil
	} else if rt.VariantIs(nazhatie, "В начало") {
		return rt.Number(0.0), nil
	} else if rt.VariantIs(nazhatie, "В конец") {
		t53, e54 := rt.Lte(ctx, dlina, rt.Number(0.0))
		if e54 != nil {
			return rt.Value{}, e54
		}
		t55, e56 := rt.Cond(ctx, t53)
		if e56 != nil {
			return rt.Value{}, e56
		}
		if t55 {
			return rt.Number(0.0), nil
		} else {
			t57, e58 := rt.Sub(ctx, dlina, rt.Number(1.0))
			if e58 != nil {
				return rt.Value{}, e58
			}
			return t57, nil
		}
	} else if rt.VariantIs(nazhatie, "Мимо") {
		return prokrutka, nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, nazhatie)
	}
}

// Prokrutit — функция flang «Прокрутить».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prokrutka — «прокрутка»: число.
// Параметр vysota — «высота»: число.
// Параметр dlina — «длина»: число.
// Параметр nazhatie — «нажатие»: «Нажатие».
// Результат — значение: число.
func Prokrutit(ctx *rt.Ctx, prokrutka rt.Value, vysota rt.Value, dlina rt.Value, nazhatie rt.Value) (rt.Value, error) {
	t59, e60 := ShagProkrutki(ctx, prokrutka, vysota, nazhatie, dlina)
	if e60 != nil {
		return rt.Value{}, e60
	}
	t61, e62 := ProkrutkaVPredelah(ctx, t59, dlina)
	if e62 != nil {
		return rt.Value{}, e62
	}
	t63 := t61
	t64, e65 := rt.Gte(ctx, t63, rt.Number(0.0))
	if e65 != nil {
		return rt.Value{}, e65
	}
	// постусловие «за верх не выходит»
	t66, e67 := rt.Post(ctx, t64, "за верх не выходит", "Прокрутить")
	if e67 != nil {
		return rt.Value{}, e67
	}
	if !t66 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «за верх не выходит» функции «Прокрутить»")
	}
	t68, e69 := rt.Lte(ctx, dlina, rt.Number(0.0))
	if e69 != nil {
		return rt.Value{}, e69
	}
	t70, e71 := rt.Cond(ctx, t68)
	if e71 != nil {
		return rt.Value{}, e71
	}
	var t72 rt.Value
	if t70 {
		t72 = rt.Number(0.0)
	} else {
		t73, e74 := rt.Sub(ctx, dlina, rt.Number(1.0))
		if e74 != nil {
			return rt.Value{}, e74
		}
		t72 = t73
	}
	t75, e76 := rt.Lte(ctx, t63, t72)
	if e76 != nil {
		return rt.Value{}, e76
	}
	// постусловие «за низ не выходит»
	t77, e78 := rt.Post(ctx, t75, "за низ не выходит", "Прокрутить")
	if e78 != nil {
		return rt.Value{}, e78
	}
	if !t77 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «за низ не выходит» функции «Прокрутить»")
	}
	return t63, nil
}

// VidnoStrok — функция flang «Видно строк».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prokrutka — «прокрутка»: число.
// Параметр vysota — «высота»: число.
// Параметр dlina — «длина»: число.
// Результат — значение: число.
func VidnoStrok(ctx *rt.Ctx, prokrutka rt.Value, vysota rt.Value, dlina rt.Value) (rt.Value, error) {
	t79, e80 := ProkrutkaVPredelah(ctx, prokrutka, dlina)
	if e80 != nil {
		return rt.Value{}, e80
	}
	// пусть «внутри»
	vnutri := t79
	t81, e82 := rt.Lte(ctx, dlina, rt.Number(0.0))
	if e82 != nil {
		return rt.Value{}, e82
	}
	t83, e84 := rt.Cond(ctx, t81)
	if e84 != nil {
		return rt.Value{}, e84
	}
	var t85 rt.Value
	if t83 {
		t85 = rt.Number(0.0)
	} else {
		t86, e87 := rt.Sub(ctx, dlina, vnutri)
		if e87 != nil {
			return rt.Value{}, e87
		}
		t85 = t86
	}
	// пусть «осталось»
	ostalos := t85
	t88, e89 := rt.Gte(ctx, vysota, rt.Number(0.0))
	if e89 != nil {
		return rt.Value{}, e89
	}
	t90, e91 := rt.Cond(ctx, t88)
	if e91 != nil {
		return rt.Value{}, e91
	}
	var t92 rt.Value
	if t90 {
		t92 = vysota
	} else {
		t92 = rt.Number(0.0)
	}
	t93, e94 := rt.Gt(ctx, ostalos, t92)
	if e94 != nil {
		return rt.Value{}, e94
	}
	t95, e96 := rt.Cond(ctx, t93)
	if e96 != nil {
		return rt.Value{}, e96
	}
	var t97 rt.Value
	if t95 {
		t98, e99 := rt.Gte(ctx, vysota, rt.Number(0.0))
		if e99 != nil {
			return rt.Value{}, e99
		}
		t100, e101 := rt.Cond(ctx, t98)
		if e101 != nil {
			return rt.Value{}, e101
		}
		var t102 rt.Value
		if t100 {
			t102 = vysota
		} else {
			t102 = rt.Number(0.0)
		}
		t97 = t102
	} else {
		t97 = ostalos
	}
	t103 := t97
	t104, e105 := rt.Gte(ctx, vysota, rt.Number(0.0))
	if e105 != nil {
		return rt.Value{}, e105
	}
	t106, e107 := rt.Cond(ctx, t104)
	if e107 != nil {
		return rt.Value{}, e107
	}
	var t108 rt.Value
	if t106 {
		t108 = vysota
	} else {
		t108 = rt.Number(0.0)
	}
	t109, e110 := rt.Lte(ctx, t103, t108)
	if e110 != nil {
		return rt.Value{}, e110
	}
	// постусловие «видно не больше высоты»
	t111, e112 := rt.Post(ctx, t109, "видно не больше высоты", "Видно строк")
	if e112 != nil {
		return rt.Value{}, e112
	}
	if !t111 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «видно не больше высоты» функции «Видно строк»")
	}
	t113, e114 := rt.Gte(ctx, t103, rt.Number(0.0))
	if e114 != nil {
		return rt.Value{}, e114
	}
	// постусловие «видно не меньше нуля»
	t115, e116 := rt.Post(ctx, t113, "видно не меньше нуля", "Видно строк")
	if e116 != nil {
		return rt.Value{}, e116
	}
	if !t115 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «видно не меньше нуля» функции «Видно строк»")
	}
	return t103, nil
}

// NadpisPolozheniya — функция flang «Надпись положения».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр prokrutka — «прокрутка»: число.
// Параметр vysota — «высота»: число.
// Параметр dlina — «длина»: число.
// Результат — значение: строка.
func NadpisPolozheniya(ctx *rt.Ctx, prokrutka rt.Value, vysota rt.Value, dlina rt.Value) (rt.Value, error) {
	t117, e118 := rt.Lte(ctx, dlina, vysota)
	if e118 != nil {
		return rt.Value{}, e118
	}
	t119, e120 := rt.Cond(ctx, t117)
	if e120 != nil {
		return rt.Value{}, e120
	}
	if t119 {
		return rt.Text(""), nil
	} else {
		t121, e122 := ProkrutkaVPredelah(ctx, prokrutka, dlina)
		if e122 != nil {
			return rt.Value{}, e122
		}
		// пусть «внутри»
		vnutri := t121
		t123, e124 := rt.Add(ctx, vnutri, rt.Number(1.0))
		if e124 != nil {
			return rt.Value{}, e124
		}
		// «к строке»
		t125, e126 := rt.BToString(ctx, t123)
		if e126 != nil {
			return rt.Value{}, e126
		}
		t127, e128 := rt.Concat(ctx, rt.Text("  строки "), t125)
		if e128 != nil {
			return rt.Value{}, e128
		}
		t129, e130 := rt.Concat(ctx, t127, rt.Text("–"))
		if e130 != nil {
			return rt.Value{}, e130
		}
		t131, e132 := VidnoStrok(ctx, vnutri, vysota, dlina)
		if e132 != nil {
			return rt.Value{}, e132
		}
		t133, e134 := rt.Add(ctx, vnutri, t131)
		if e134 != nil {
			return rt.Value{}, e134
		}
		// «к строке»
		t135, e136 := rt.BToString(ctx, t133)
		if e136 != nil {
			return rt.Value{}, e136
		}
		t137, e138 := rt.Concat(ctx, t129, t135)
		if e138 != nil {
			return rt.Value{}, e138
		}
		t139, e140 := rt.Concat(ctx, t137, rt.Text(" из "))
		if e140 != nil {
			return rt.Value{}, e140
		}
		// «к строке»
		t141, e142 := rt.BToString(ctx, dlina)
		if e142 != nil {
			return rt.Value{}, e142
		}
		t143, e144 := rt.Concat(ctx, t139, t141)
		if e144 != nil {
			return rt.Value{}, e144
		}
		return t143, nil
	}
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Высота тела":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Высота тела", 2, len(args))
		}
		return VysotaTela(ctx, args[0], args[1])
	case "Прокрутка в пределах":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прокрутка в пределах", 2, len(args))
		}
		return ProkrutkaVPredelah(ctx, args[0], args[1])
	case "Шаг прокрутки":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Шаг прокрутки", 4, len(args))
		}
		return ShagProkrutki(ctx, args[0], args[1], args[2], args[3])
	case "Прокрутить":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прокрутить", 4, len(args))
		}
		return Prokrutit(ctx, args[0], args[1], args[2], args[3])
	case "Видно строк":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Видно строк", 3, len(args))
		}
		return VidnoStrok(ctx, args[0], args[1], args[2])
	case "Надпись положения":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Надпись положения", 3, len(args))
		}
		return NadpisPolozheniya(ctx, args[0], args[1], args[2])
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
