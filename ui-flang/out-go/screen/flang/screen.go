// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
// Модуль flang: «Screen».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flangscreen/flangrt"
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

// SozdatHod — запись FTS «Ход»: «сост», «ширина», «куски», «порог».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatHod(sost rt.Value, shirina rt.Value, kuski rt.Value, porog rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "сост", Value: sost},
		{Name: "ширина", Value: shirina},
		{Name: "куски", Value: kuski},
		{Name: "порог", Value: porog},
	})
}

// VariantZnaki — вариант «Знаки» суммы типов «Разбор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantZnaki() rt.Value {
	return rt.Variant("Знаки", nil)
}

// VariantZagolovok — вариант «Заголовок» суммы типов «Разбор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantZagolovok() rt.Value {
	return rt.Variant("Заголовок", nil)
}

// VariantVnutri — вариант «Внутри» суммы типов «Разбор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantVnutri() rt.Value {
	return rt.Variant("Внутри", nil)
}

// VariantOborvano — вариант «Оборвано» суммы типов «Разбор».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantOborvano() rt.Value {
	return rt.Variant("Оборвано", nil)
}

// DolyaVPredelah — функция flang «Доля в пределах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Результат — значение: число.
func DolyaVPredelah(ctx *rt.Ctx, dolya rt.Value) (rt.Value, error) {
	t1, e2 := rt.Lt(ctx, dolya, rt.Number(0.0))
	if e2 != nil {
		return rt.Value{}, e2
	}
	t3, e4 := rt.Cond(ctx, t1)
	if e4 != nil {
		return rt.Value{}, e4
	}
	var t5 rt.Value
	if t3 {
		t5 = rt.Number(0.0)
	} else {
		t6, e7 := rt.Gt(ctx, dolya, rt.Number(1.0))
		if e7 != nil {
			return rt.Value{}, e7
		}
		t8, e9 := rt.Cond(ctx, t6)
		if e9 != nil {
			return rt.Value{}, e9
		}
		var t10 rt.Value
		if t8 {
			t10 = rt.Number(1.0)
		} else {
			t10 = dolya
		}
		t5 = t10
	}
	t11 := t5
	t12, e13 := rt.Gte(ctx, t11, rt.Number(0.0))
	if e13 != nil {
		return rt.Value{}, e13
	}
	// постусловие «доля не меньше нуля»
	t14, e15 := rt.Post(ctx, t12, "доля не меньше нуля", "Доля в пределах")
	if e15 != nil {
		return rt.Value{}, e15
	}
	if !t14 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «доля не меньше нуля» функции «Доля в пределах»")
	}
	t16, e17 := rt.Lte(ctx, t11, rt.Number(1.0))
	if e17 != nil {
		return rt.Value{}, e17
	}
	// постусловие «доля не больше единицы»
	t18, e19 := rt.Post(ctx, t16, "доля не больше единицы", "Доля в пределах")
	if e19 != nil {
		return rt.Value{}, e19
	}
	if !t18 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «доля не больше единицы» функции «Доля в пределах»")
	}
	return t11, nil
}

// CeloeOt — функция flang «Целое от».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение: число.
func CeloeOt(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
	t20, e21 := rt.Mod(ctx, znachenie, rt.Number(1.0))
	if e21 != nil {
		return rt.Value{}, e21
	}
	t22, e23 := rt.Sub(ctx, znachenie, t20)
	if e23 != nil {
		return rt.Value{}, e23
	}
	return t22, nil
}

// KletokPolosy — функция flang «Клеток полосы».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Параметр shirina — «ширина»: «нат».
// Результат — значение: число.
func KletokPolosy(ctx *rt.Ctx, dolya rt.Value, shirina rt.Value) (rt.Value, error) {
	t24, e25 := DolyaVPredelah(ctx, dolya)
	if e25 != nil {
		return rt.Value{}, e25
	}
	t26, e27 := rt.Mul(ctx, t24, shirina)
	if e27 != nil {
		return rt.Value{}, e27
	}
	t28, e29 := rt.Add(ctx, t26, rt.Number(0.5))
	if e29 != nil {
		return rt.Value{}, e29
	}
	t30, e31 := CeloeOt(ctx, t28)
	if e31 != nil {
		return rt.Value{}, e31
	}
	t32 := t30
	t33, e34 := rt.Gte(ctx, t32, rt.Number(0.0))
	if e34 != nil {
		return rt.Value{}, e34
	}
	// постусловие «закрашенных клеток не меньше нуля»
	t35, e36 := rt.Post(ctx, t33, "закрашенных клеток не меньше нуля", "Клеток полосы")
	if e36 != nil {
		return rt.Value{}, e36
	}
	if !t35 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «закрашенных клеток не меньше нуля» функции «Клеток полосы»")
	}
	t37, e38 := rt.Lte(ctx, t32, shirina)
	if e38 != nil {
		return rt.Value{}, e38
	}
	// постусловие «закрашенных клеток не больше ширины»
	t39, e40 := rt.Post(ctx, t37, "закрашенных клеток не больше ширины", "Клеток полосы")
	if e40 != nil {
		return rt.Value{}, e40
	}
	if !t39 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «закрашенных клеток не больше ширины» функции «Клеток полосы»")
	}
	return t32, nil
}

// KletkiPolosy — функция flang «Клетки полосы».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр vsego — «всего»: «нат».
// Параметр zakrasheno — «закрашено»: число.
// Параметр polnyy — «полный»: строка.
// Параметр pustoy — «пустой»: строка.
// Результат — значение: список: строка.
func KletkiPolosy(ctx *rt.Ctx, vsego rt.Value, zakrasheno rt.Value, polnyy rt.Value, pustoy rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Клетки полосы"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t41, e42 := rt.Lte(ctx, vsego, rt.Number(0.0))
	if e42 != nil {
		return rt.Value{}, e42
	}
	t43, e44 := rt.Cond(ctx, t41)
	if e44 != nil {
		return rt.Value{}, e44
	}
	var t45 rt.Value
	if t43 {
		t45 = rt.List(nil)
	} else {
		t46, e47 := rt.Lte(ctx, vsego, zakrasheno)
		if e47 != nil {
			return rt.Value{}, e47
		}
		t48, e49 := rt.Cond(ctx, t46)
		if e49 != nil {
			return rt.Value{}, e49
		}
		var t50 rt.Value
		if t48 {
			t50 = polnyy
		} else {
			t50 = pustoy
		}
		t51, e52 := rt.Sub(ctx, vsego, rt.Number(1.0))
		if e52 != nil {
			return rt.Value{}, e52
		}
		t53, e54 := KletkiPolosy(ctx, t51, zakrasheno, polnyy, pustoy)
		if e54 != nil {
			return rt.Value{}, e54
		}
		// «добавить»
		t55, e56 := rt.BAppend(ctx, t50, t53)
		if e56 != nil {
			return rt.Value{}, e56
		}
		t45 = t55
	}
	t57 := t45
	// «длина»
	t58, e59 := rt.BLength(ctx, t57)
	if e59 != nil {
		return rt.Value{}, e59
	}
	// постусловие «клеток ровно столько, сколько просили»
	t60, e61 := rt.Post(ctx, rt.Flag(rt.Equal(t58, vsego)), "клеток ровно столько, сколько просили", "Клетки полосы")
	if e61 != nil {
		return rt.Value{}, e61
	}
	if !t60 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «клеток ровно столько, сколько просили» функции «Клетки полосы»")
	}
	return t57, nil
}

// Polosa — функция flang «Полоса».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Параметр shirina — «ширина»: «нат».
// Параметр polnyy — «полный»: строка.
// Параметр pustoy — «пустой»: строка.
// Результат — значение: строка.
func Polosa(ctx *rt.Ctx, dolya rt.Value, shirina rt.Value, polnyy rt.Value, pustoy rt.Value) (rt.Value, error) {
	t62, e63 := KletokPolosy(ctx, dolya, shirina)
	if e63 != nil {
		return rt.Value{}, e63
	}
	t64, e65 := KletkiPolosy(ctx, shirina, t62, polnyy, pustoy)
	if e65 != nil {
		return rt.Value{}, e65
	}
	// «соединить»
	t66, e67 := rt.BJoin(ctx, t64, rt.Text(""))
	if e67 != nil {
		return rt.Value{}, e67
	}
	t68 := t66
	// «длина»
	t69, e70 := rt.BLength(ctx, polnyy)
	if e70 != nil {
		return rt.Value{}, e70
	}
	t71, e72 := rt.Cond(ctx, rt.Flag(rt.Equal(t69, rt.Number(1.0))))
	if e72 != nil {
		return rt.Value{}, e72
	}
	var t73 rt.Value
	if t71 {
		// «длина»
		t74, e75 := rt.BLength(ctx, pustoy)
		if e75 != nil {
			return rt.Value{}, e75
		}
		t73 = rt.Flag(rt.Equal(t74, rt.Number(1.0)))
	} else {
		t73 = rt.Flag(false)
	}
	t76, e77 := rt.Cond(ctx, t73)
	if e77 != nil {
		return rt.Value{}, e77
	}
	var t78 rt.Value
	if t76 {
		// «длина»
		t79, e80 := rt.BLength(ctx, t68)
		if e80 != nil {
			return rt.Value{}, e80
		}
		t78 = rt.Flag(rt.Equal(t79, shirina))
	} else {
		t78 = rt.Flag(true)
	}
	// постусловие «полоса ровно в заказанную ширину»
	t81, e82 := rt.Post(ctx, t78, "полоса ровно в заказанную ширину", "Полоса")
	if e82 != nil {
		return rt.Value{}, e82
	}
	if !t81 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «полоса ровно в заказанную ширину» функции «Полоса»")
	}
	return t68, nil
}

// Esk — функция flang «Эск».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: строка.
func Esk(ctx *rt.Ctx) (rt.Value, error) {
	// «символ по коду»
	t83, e84 := rt.BCharFromCode(ctx, rt.Number(27.0))
	if e84 != nil {
		return rt.Value{}, e84
	}
	return t83, nil
}

// ZakryvaetPosledovatelnost — функция flang «Закрывает последовательность».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znak — «знак»: строка.
// Результат — значение.
func ZakryvaetPosledovatelnost(ctx *rt.Ctx, znak rt.Value) (rt.Value, error) {
	// «код символа»
	t85, e86 := rt.BCharCode(ctx, znak)
	if e86 != nil {
		return rt.Value{}, e86
	}
	// пусть «код»
	kod := t85
	t87, e88 := rt.Gte(ctx, kod, rt.Number(64.0))
	if e88 != nil {
		return rt.Value{}, e88
	}
	t89, e90 := rt.Cond(ctx, t87)
	if e90 != nil {
		return rt.Value{}, e90
	}
	var t91 rt.Value
	if t89 {
		t92, e93 := rt.Lte(ctx, kod, rt.Number(126.0))
		if e93 != nil {
			return rt.Value{}, e93
		}
		t91 = t92
	} else {
		t91 = rt.Flag(false)
	}
	t94, e95 := rt.Cond(ctx, t91)
	if e95 != nil {
		return rt.Value{}, e95
	}
	var t96 rt.Value
	if t94 {
		t97, e98 := rt.Cond(ctx, rt.Flag(rt.Equal(kod, rt.Number(91.0))))
		if e98 != nil {
			return rt.Value{}, e98
		}
		var t99 rt.Value
		if t97 {
			t99 = rt.Flag(false)
		} else {
			t99 = rt.Flag(true)
		}
		t96 = t99
	} else {
		t96 = rt.Flag(false)
	}
	t100, e101 := rt.Cond(ctx, t96)
	if e101 != nil {
		return rt.Value{}, e101
	}
	var t102 rt.Value
	if t100 {
		t103, e104 := rt.Cond(ctx, rt.Flag(rt.Equal(kod, rt.Number(59.0))))
		if e104 != nil {
			return rt.Value{}, e104
		}
		var t105 rt.Value
		if t103 {
			t105 = rt.Flag(false)
		} else {
			t105 = rt.Flag(true)
		}
		t102 = t105
	} else {
		t102 = rt.Flag(false)
	}
	t106, e107 := rt.Cond(ctx, t102)
	if e107 != nil {
		return rt.Value{}, e107
	}
	if t106 {
		t108, e109 := rt.Gte(ctx, kod, rt.Number(48.0))
		if e109 != nil {
			return rt.Value{}, e109
		}
		t110, e111 := rt.Cond(ctx, t108)
		if e111 != nil {
			return rt.Value{}, e111
		}
		var t112 rt.Value
		if t110 {
			t113, e114 := rt.Lte(ctx, kod, rt.Number(57.0))
			if e114 != nil {
				return rt.Value{}, e114
			}
			t112 = t113
		} else {
			t112 = rt.Flag(false)
		}
		t115, e116 := rt.Cond(ctx, t112)
		if e116 != nil {
			return rt.Value{}, e116
		}
		if t115 {
			return rt.Flag(false), nil
		} else {
			return rt.Flag(true), nil
		}
	} else {
		return rt.Flag(false), nil
	}
}

// HodNachala — функция flang «Ход начала».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр porog — «порог»: число.
// Результат — значение: «Ход».
func HodNachala(ctx *rt.Ctx, porog rt.Value) (rt.Value, error) {
	t117 := make([]rt.Field, 4)
	t117[0] = rt.Field{Name: "сост", Value: rt.Variant("Знаки", nil)}
	t117[1] = rt.Field{Name: "ширина", Value: rt.Number(0.0)}
	t117[2] = rt.Field{Name: "куски", Value: rt.List(nil)}
	t117[3] = rt.Field{Name: "порог", Value: porog}
	t118 := rt.Record(t117)
	t119, e120 := rt.FieldGet(ctx, t118, "ширина")
	if e120 != nil {
		return rt.Value{}, e120
	}
	// постусловие «начальная ширина нулевая»
	t121, e122 := rt.Post(ctx, rt.Flag(rt.Equal(t119, rt.Number(0.0))), "начальная ширина нулевая", "Ход начала")
	if e122 != nil {
		return rt.Value{}, e122
	}
	if !t121 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «начальная ширина нулевая» функции «Ход начала»")
	}
	return t118, nil
}

// OborvanoLi — функция flang «Оборвано ли».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sost — «сост»: «Разбор».
// Результат — значение.
func OborvanoLi(ctx *rt.Ctx, sost rt.Value) (rt.Value, error) {
	if rt.VariantIs(sost, "Знаки") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(sost, "Заголовок") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(sost, "Внутри") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(sost, "Оборвано") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, sost)
	}
}

// ShagShiriny — функция flang «Шаг ширины».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр hod — «ход»: «Ход».
// Параметр znak — «знак»: строка.
// Результат — значение: «Ход».
func ShagShiriny(ctx *rt.Ctx, hod rt.Value, znak rt.Value) (rt.Value, error) {
	var t123 rt.Value
	t124, e125 := rt.FieldGet(ctx, hod, "сост")
	if e125 != nil {
		return rt.Value{}, e125
	}
	if rt.VariantIs(t124, "Оборвано") {
		t123 = hod
	} else if rt.VariantIs(t124, "Знаки") {
		t126, e127 := Esk(ctx)
		if e127 != nil {
			return rt.Value{}, e127
		}
		t128, e129 := rt.Cond(ctx, rt.Flag(rt.Equal(znak, t126)))
		if e129 != nil {
			return rt.Value{}, e129
		}
		var t130 rt.Value
		if t128 {
			t131, e132 := rt.FieldGet(ctx, hod, "ширина")
			if e132 != nil {
				return rt.Value{}, e132
			}
			t133, e134 := rt.FieldGet(ctx, hod, "куски")
			if e134 != nil {
				return rt.Value{}, e134
			}
			t135, e136 := rt.FieldGet(ctx, hod, "порог")
			if e136 != nil {
				return rt.Value{}, e136
			}
			t137 := make([]rt.Field, 4)
			t137[0] = rt.Field{Name: "сост", Value: rt.Variant("Заголовок", nil)}
			t137[1] = rt.Field{Name: "ширина", Value: t131}
			t137[2] = rt.Field{Name: "куски", Value: t133}
			t137[3] = rt.Field{Name: "порог", Value: t135}
			t130 = rt.Record(t137)
		} else {
			t138, e139 := rt.FieldGet(ctx, hod, "ширина")
			if e139 != nil {
				return rt.Value{}, e139
			}
			t140, e141 := rt.Add(ctx, t138, rt.Number(1.0))
			if e141 != nil {
				return rt.Value{}, e141
			}
			t142, e143 := rt.FieldGet(ctx, hod, "куски")
			if e143 != nil {
				return rt.Value{}, e143
			}
			t144, e145 := rt.FieldGet(ctx, hod, "порог")
			if e145 != nil {
				return rt.Value{}, e145
			}
			t146 := make([]rt.Field, 4)
			t146[0] = rt.Field{Name: "сост", Value: rt.Variant("Знаки", nil)}
			t146[1] = rt.Field{Name: "ширина", Value: t140}
			t146[2] = rt.Field{Name: "куски", Value: t142}
			t146[3] = rt.Field{Name: "порог", Value: t144}
			t130 = rt.Record(t146)
		}
		t123 = t130
	} else if rt.VariantIs(t124, "Заголовок") {
		t147, e148 := rt.Cond(ctx, rt.Flag(rt.Equal(znak, rt.Text("["))))
		if e148 != nil {
			return rt.Value{}, e148
		}
		var t149 rt.Value
		if t147 {
			t149 = rt.Flag(true)
		} else {
			t149 = rt.Flag(rt.Equal(znak, rt.Text("]")))
		}
		t150, e151 := rt.Cond(ctx, t149)
		if e151 != nil {
			return rt.Value{}, e151
		}
		var t152 rt.Value
		if t150 {
			t152 = rt.Flag(true)
		} else {
			t152 = rt.Flag(rt.Equal(znak, rt.Text("?")))
		}
		t153, e154 := rt.Cond(ctx, t152)
		if e154 != nil {
			return rt.Value{}, e154
		}
		var t155 rt.Value
		if t153 {
			t155 = rt.Variant("Внутри", nil)
		} else {
			t155 = rt.Variant("Знаки", nil)
		}
		t156, e157 := rt.FieldGet(ctx, hod, "ширина")
		if e157 != nil {
			return rt.Value{}, e157
		}
		t158, e159 := rt.FieldGet(ctx, hod, "куски")
		if e159 != nil {
			return rt.Value{}, e159
		}
		t160, e161 := rt.FieldGet(ctx, hod, "порог")
		if e161 != nil {
			return rt.Value{}, e161
		}
		t162 := make([]rt.Field, 4)
		t162[0] = rt.Field{Name: "сост", Value: t155}
		t162[1] = rt.Field{Name: "ширина", Value: t156}
		t162[2] = rt.Field{Name: "куски", Value: t158}
		t162[3] = rt.Field{Name: "порог", Value: t160}
		t123 = rt.Record(t162)
	} else if rt.VariantIs(t124, "Внутри") {
		t163, e164 := ZakryvaetPosledovatelnost(ctx, znak)
		if e164 != nil {
			return rt.Value{}, e164
		}
		t165, e166 := rt.Cond(ctx, t163)
		if e166 != nil {
			return rt.Value{}, e166
		}
		var t167 rt.Value
		if t165 {
			t167 = rt.Variant("Знаки", nil)
		} else {
			t167 = rt.Variant("Внутри", nil)
		}
		t168, e169 := rt.FieldGet(ctx, hod, "ширина")
		if e169 != nil {
			return rt.Value{}, e169
		}
		t170, e171 := rt.FieldGet(ctx, hod, "куски")
		if e171 != nil {
			return rt.Value{}, e171
		}
		t172, e173 := rt.FieldGet(ctx, hod, "порог")
		if e173 != nil {
			return rt.Value{}, e173
		}
		t174 := make([]rt.Field, 4)
		t174[0] = rt.Field{Name: "сост", Value: t167}
		t174[1] = rt.Field{Name: "ширина", Value: t168}
		t174[2] = rt.Field{Name: "куски", Value: t170}
		t174[3] = rt.Field{Name: "порог", Value: t172}
		t123 = rt.Record(t174)
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t124)
	}
	t175 := t123
	t176, e177 := rt.FieldGet(ctx, t175, "ширина")
	if e177 != nil {
		return rt.Value{}, e177
	}
	t178, e179 := rt.FieldGet(ctx, hod, "ширина")
	if e179 != nil {
		return rt.Value{}, e179
	}
	t180, e181 := rt.Gte(ctx, t176, t178)
	if e181 != nil {
		return rt.Value{}, e181
	}
	// постусловие «ширина не убывает»
	t182, e183 := rt.Post(ctx, t180, "ширина не убывает", "Шаг ширины")
	if e183 != nil {
		return rt.Value{}, e183
	}
	if !t182 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина не убывает» функции «Шаг ширины»")
	}
	t184, e185 := rt.FieldGet(ctx, t175, "ширина")
	if e185 != nil {
		return rt.Value{}, e185
	}
	t186, e187 := rt.FieldGet(ctx, hod, "ширина")
	if e187 != nil {
		return rt.Value{}, e187
	}
	t188, e189 := rt.Add(ctx, t186, rt.Number(1.0))
	if e189 != nil {
		return rt.Value{}, e189
	}
	t190, e191 := rt.Lte(ctx, t184, t188)
	if e191 != nil {
		return rt.Value{}, e191
	}
	// постусловие «ширина растёт не более чем на одну клетку»
	t192, e193 := rt.Post(ctx, t190, "ширина растёт не более чем на одну клетку", "Шаг ширины")
	if e193 != nil {
		return rt.Value{}, e193
	}
	if !t192 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина растёт не более чем на одну клетку» функции «Шаг ширины»")
	}
	return t175, nil
}

// ShirinaBezPosledovatelnostey — функция flang «Ширина без последовательностей».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Результат — значение: число.
func ShirinaBezPosledovatelnostey(ctx *rt.Ctx, tekst rt.Value) (rt.Value, error) {
	// «символы»
	t194, e195 := rt.BCharacters(ctx, tekst)
	if e195 != nil {
		return rt.Value{}, e195
	}
	t196, e197 := rt.RequireList(ctx, t194, "свёртка")
	if e197 != nil {
		return rt.Value{}, e197
	}
	t198, e199 := HodNachala(ctx, rt.Number(0.0))
	if e199 != nil {
		return rt.Value{}, e199
	}
	// «ход»
	hod := t198
	for t200 := range t196 {
		// «знак»
		znak := t196[t200]
		t201, e202 := ShagShiriny(ctx, hod, znak)
		if e202 != nil {
			return rt.Value{}, e202
		}
		hod = t201
	}
	t203, e204 := rt.FieldGet(ctx, hod, "ширина")
	if e204 != nil {
		return rt.Value{}, e204
	}
	t205 := t203
	// «длина»
	t206, e207 := rt.BLength(ctx, tekst)
	if e207 != nil {
		return rt.Value{}, e207
	}
	t208, e209 := rt.Lte(ctx, t205, t206)
	if e209 != nil {
		return rt.Value{}, e209
	}
	// постусловие «ширина не больше числа знаков»
	t210, e211 := rt.Post(ctx, t208, "ширина не больше числа знаков", "Ширина без последовательностей")
	if e211 != nil {
		return rt.Value{}, e211
	}
	if !t210 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина не больше числа знаков» функции «Ширина без последовательностей»")
	}
	t212, e213 := rt.Gte(ctx, t205, rt.Number(0.0))
	if e213 != nil {
		return rt.Value{}, e213
	}
	// постусловие «ширина не меньше нуля»
	t214, e215 := rt.Post(ctx, t212, "ширина не меньше нуля", "Ширина без последовательностей")
	if e215 != nil {
		return rt.Value{}, e215
	}
	if !t214 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина не меньше нуля» функции «Ширина без последовательностей»")
	}
	return t205, nil
}

// PosledovatelnostiCely — функция flang «Последовательности целы».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Результат — значение.
func PosledovatelnostiCely(ctx *rt.Ctx, tekst rt.Value) (rt.Value, error) {
	// «символы»
	t216, e217 := rt.BCharacters(ctx, tekst)
	if e217 != nil {
		return rt.Value{}, e217
	}
	t218, e219 := rt.RequireList(ctx, t216, "свёртка")
	if e219 != nil {
		return rt.Value{}, e219
	}
	t220, e221 := HodNachala(ctx, rt.Number(0.0))
	if e221 != nil {
		return rt.Value{}, e221
	}
	// «ход»
	hod := t220
	for t222 := range t218 {
		// «знак»
		znak := t218[t222]
		t223, e224 := ShagShiriny(ctx, hod, znak)
		if e224 != nil {
			return rt.Value{}, e224
		}
		hod = t223
	}
	t225, e226 := rt.FieldGet(ctx, hod, "сост")
	if e226 != nil {
		return rt.Value{}, e226
	}
	if rt.VariantIs(t225, "Знаки") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t225, "Заголовок") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t225, "Внутри") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t225, "Оборвано") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t225)
	}
}

// ShagObrezki — функция flang «Шаг обрезки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр hod — «ход»: «Ход».
// Параметр znak — «знак»: строка.
// Результат — значение: «Ход».
func ShagObrezki(ctx *rt.Ctx, hod rt.Value, znak rt.Value) (rt.Value, error) {
	var t227 rt.Value
	t228, e229 := rt.FieldGet(ctx, hod, "сост")
	if e229 != nil {
		return rt.Value{}, e229
	}
	if rt.VariantIs(t228, "Оборвано") {
		t227 = hod
	} else if rt.VariantIs(t228, "Знаки") {
		t230, e231 := Esk(ctx)
		if e231 != nil {
			return rt.Value{}, e231
		}
		t232, e233 := rt.Cond(ctx, rt.Flag(rt.Equal(znak, t230)))
		if e233 != nil {
			return rt.Value{}, e233
		}
		var t234 rt.Value
		if t232 {
			t235, e236 := rt.FieldGet(ctx, hod, "ширина")
			if e236 != nil {
				return rt.Value{}, e236
			}
			t237, e238 := rt.FieldGet(ctx, hod, "куски")
			if e238 != nil {
				return rt.Value{}, e238
			}
			// «добавить»
			t239, e240 := rt.BAppend(ctx, znak, t237)
			if e240 != nil {
				return rt.Value{}, e240
			}
			t241, e242 := rt.FieldGet(ctx, hod, "порог")
			if e242 != nil {
				return rt.Value{}, e242
			}
			t243 := make([]rt.Field, 4)
			t243[0] = rt.Field{Name: "сост", Value: rt.Variant("Заголовок", nil)}
			t243[1] = rt.Field{Name: "ширина", Value: t235}
			t243[2] = rt.Field{Name: "куски", Value: t239}
			t243[3] = rt.Field{Name: "порог", Value: t241}
			t234 = rt.Record(t243)
		} else {
			t244, e245 := rt.FieldGet(ctx, hod, "ширина")
			if e245 != nil {
				return rt.Value{}, e245
			}
			t246, e247 := rt.FieldGet(ctx, hod, "порог")
			if e247 != nil {
				return rt.Value{}, e247
			}
			t248, e249 := rt.Gte(ctx, t244, t246)
			if e249 != nil {
				return rt.Value{}, e249
			}
			t250, e251 := rt.Cond(ctx, t248)
			if e251 != nil {
				return rt.Value{}, e251
			}
			var t252 rt.Value
			if t250 {
				t253, e254 := rt.FieldGet(ctx, hod, "ширина")
				if e254 != nil {
					return rt.Value{}, e254
				}
				t255, e256 := rt.FieldGet(ctx, hod, "куски")
				if e256 != nil {
					return rt.Value{}, e256
				}
				t257, e258 := rt.FieldGet(ctx, hod, "порог")
				if e258 != nil {
					return rt.Value{}, e258
				}
				t259 := make([]rt.Field, 4)
				t259[0] = rt.Field{Name: "сост", Value: rt.Variant("Оборвано", nil)}
				t259[1] = rt.Field{Name: "ширина", Value: t253}
				t259[2] = rt.Field{Name: "куски", Value: t255}
				t259[3] = rt.Field{Name: "порог", Value: t257}
				t252 = rt.Record(t259)
			} else {
				t260, e261 := rt.FieldGet(ctx, hod, "ширина")
				if e261 != nil {
					return rt.Value{}, e261
				}
				t262, e263 := rt.Add(ctx, t260, rt.Number(1.0))
				if e263 != nil {
					return rt.Value{}, e263
				}
				t264, e265 := rt.FieldGet(ctx, hod, "куски")
				if e265 != nil {
					return rt.Value{}, e265
				}
				// «добавить»
				t266, e267 := rt.BAppend(ctx, znak, t264)
				if e267 != nil {
					return rt.Value{}, e267
				}
				t268, e269 := rt.FieldGet(ctx, hod, "порог")
				if e269 != nil {
					return rt.Value{}, e269
				}
				t270 := make([]rt.Field, 4)
				t270[0] = rt.Field{Name: "сост", Value: rt.Variant("Знаки", nil)}
				t270[1] = rt.Field{Name: "ширина", Value: t262}
				t270[2] = rt.Field{Name: "куски", Value: t266}
				t270[3] = rt.Field{Name: "порог", Value: t268}
				t252 = rt.Record(t270)
			}
			t234 = t252
		}
		t227 = t234
	} else if rt.VariantIs(t228, "Заголовок") {
		t271, e272 := rt.Cond(ctx, rt.Flag(rt.Equal(znak, rt.Text("["))))
		if e272 != nil {
			return rt.Value{}, e272
		}
		var t273 rt.Value
		if t271 {
			t273 = rt.Flag(true)
		} else {
			t273 = rt.Flag(rt.Equal(znak, rt.Text("]")))
		}
		t274, e275 := rt.Cond(ctx, t273)
		if e275 != nil {
			return rt.Value{}, e275
		}
		var t276 rt.Value
		if t274 {
			t276 = rt.Flag(true)
		} else {
			t276 = rt.Flag(rt.Equal(znak, rt.Text("?")))
		}
		t277, e278 := rt.Cond(ctx, t276)
		if e278 != nil {
			return rt.Value{}, e278
		}
		var t279 rt.Value
		if t277 {
			t279 = rt.Variant("Внутри", nil)
		} else {
			t279 = rt.Variant("Знаки", nil)
		}
		t280, e281 := rt.FieldGet(ctx, hod, "ширина")
		if e281 != nil {
			return rt.Value{}, e281
		}
		t282, e283 := rt.FieldGet(ctx, hod, "куски")
		if e283 != nil {
			return rt.Value{}, e283
		}
		// «добавить»
		t284, e285 := rt.BAppend(ctx, znak, t282)
		if e285 != nil {
			return rt.Value{}, e285
		}
		t286, e287 := rt.FieldGet(ctx, hod, "порог")
		if e287 != nil {
			return rt.Value{}, e287
		}
		t288 := make([]rt.Field, 4)
		t288[0] = rt.Field{Name: "сост", Value: t279}
		t288[1] = rt.Field{Name: "ширина", Value: t280}
		t288[2] = rt.Field{Name: "куски", Value: t284}
		t288[3] = rt.Field{Name: "порог", Value: t286}
		t227 = rt.Record(t288)
	} else if rt.VariantIs(t228, "Внутри") {
		t289, e290 := ZakryvaetPosledovatelnost(ctx, znak)
		if e290 != nil {
			return rt.Value{}, e290
		}
		t291, e292 := rt.Cond(ctx, t289)
		if e292 != nil {
			return rt.Value{}, e292
		}
		var t293 rt.Value
		if t291 {
			t293 = rt.Variant("Знаки", nil)
		} else {
			t293 = rt.Variant("Внутри", nil)
		}
		t294, e295 := rt.FieldGet(ctx, hod, "ширина")
		if e295 != nil {
			return rt.Value{}, e295
		}
		t296, e297 := rt.FieldGet(ctx, hod, "куски")
		if e297 != nil {
			return rt.Value{}, e297
		}
		// «добавить»
		t298, e299 := rt.BAppend(ctx, znak, t296)
		if e299 != nil {
			return rt.Value{}, e299
		}
		t300, e301 := rt.FieldGet(ctx, hod, "порог")
		if e301 != nil {
			return rt.Value{}, e301
		}
		t302 := make([]rt.Field, 4)
		t302[0] = rt.Field{Name: "сост", Value: t293}
		t302[1] = rt.Field{Name: "ширина", Value: t294}
		t302[2] = rt.Field{Name: "куски", Value: t298}
		t302[3] = rt.Field{Name: "порог", Value: t300}
		t227 = rt.Record(t302)
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t228)
	}
	t303 := t227
	t304, e305 := rt.FieldGet(ctx, t303, "ширина")
	if e305 != nil {
		return rt.Value{}, e305
	}
	t306, e307 := rt.FieldGet(ctx, hod, "ширина")
	if e307 != nil {
		return rt.Value{}, e307
	}
	t308, e309 := rt.Gte(ctx, t304, t306)
	if e309 != nil {
		return rt.Value{}, e309
	}
	// постусловие «ширина не убывает»
	t310, e311 := rt.Post(ctx, t308, "ширина не убывает", "Шаг обрезки")
	if e311 != nil {
		return rt.Value{}, e311
	}
	if !t310 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина не убывает» функции «Шаг обрезки»")
	}
	t312, e313 := rt.FieldGet(ctx, hod, "сост")
	if e313 != nil {
		return rt.Value{}, e313
	}
	t314, e315 := OborvanoLi(ctx, t312)
	if e315 != nil {
		return rt.Value{}, e315
	}
	t316, e317 := rt.Cond(ctx, t314)
	if e317 != nil {
		return rt.Value{}, e317
	}
	var t318 rt.Value
	if t316 {
		t318 = rt.Flag(rt.Equal(t303, hod))
	} else {
		t318 = rt.Flag(true)
	}
	// постусловие «оборванное больше не растёт»
	t319, e320 := rt.Post(ctx, t318, "оборванное больше не растёт", "Шаг обрезки")
	if e320 != nil {
		return rt.Value{}, e320
	}
	if !t319 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «оборванное больше не растёт» функции «Шаг обрезки»")
	}
	t321, e322 := rt.FieldGet(ctx, t303, "порог")
	if e322 != nil {
		return rt.Value{}, e322
	}
	t323, e324 := rt.FieldGet(ctx, hod, "порог")
	if e324 != nil {
		return rt.Value{}, e324
	}
	// постусловие «порог не меняется»
	t325, e326 := rt.Post(ctx, rt.Flag(rt.Equal(t321, t323)), "порог не меняется", "Шаг обрезки")
	if e326 != nil {
		return rt.Value{}, e326
	}
	if !t325 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «порог не меняется» функции «Шаг обрезки»")
	}
	t327, e328 := rt.FieldGet(ctx, t303, "куски")
	if e328 != nil {
		return rt.Value{}, e328
	}
	// «длина»
	t329, e330 := rt.BLength(ctx, t327)
	if e330 != nil {
		return rt.Value{}, e330
	}
	t331, e332 := rt.FieldGet(ctx, hod, "куски")
	if e332 != nil {
		return rt.Value{}, e332
	}
	// «длина»
	t333, e334 := rt.BLength(ctx, t331)
	if e334 != nil {
		return rt.Value{}, e334
	}
	t335, e336 := rt.Add(ctx, t333, rt.Number(1.0))
	if e336 != nil {
		return rt.Value{}, e336
	}
	t337, e338 := rt.Lte(ctx, t329, t335)
	if e338 != nil {
		return rt.Value{}, e338
	}
	// постусловие «кусков прибавляется не больше одного»
	t339, e340 := rt.Post(ctx, t337, "кусков прибавляется не больше одного", "Шаг обрезки")
	if e340 != nil {
		return rt.Value{}, e340
	}
	if !t339 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «кусков прибавляется не больше одного» функции «Шаг обрезки»")
	}
	return t303, nil
}

// ObrezatPoYacheykam — функция flang «Обрезать по ячейкам».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр predel — «предел»: число.
// Параметр mnogotochie — «многоточие»: строка.
// Параметр sbros — «сброс»: строка.
// Результат — значение: строка.
func ObrezatPoYacheykam(ctx *rt.Ctx, tekst rt.Value, predel rt.Value, mnogotochie rt.Value, sbros rt.Value) (rt.Value, error) {
	t341, e342 := rt.Lte(ctx, predel, rt.Number(0.0))
	if e342 != nil {
		return rt.Value{}, e342
	}
	t343, e344 := rt.Cond(ctx, t341)
	if e344 != nil {
		return rt.Value{}, e344
	}
	var t345 rt.Value
	if t343 {
		t345 = rt.Text("")
	} else {
		t346, e347 := ShirinaBezPosledovatelnostey(ctx, tekst)
		if e347 != nil {
			return rt.Value{}, e347
		}
		t348, e349 := rt.Lte(ctx, t346, predel)
		if e349 != nil {
			return rt.Value{}, e349
		}
		t350, e351 := rt.Cond(ctx, t348)
		if e351 != nil {
			return rt.Value{}, e351
		}
		var t352 rt.Value
		if t350 {
			t352 = tekst
		} else {
			// «символы»
			t353, e354 := rt.BCharacters(ctx, tekst)
			if e354 != nil {
				return rt.Value{}, e354
			}
			t355, e356 := rt.RequireList(ctx, t353, "свёртка")
			if e356 != nil {
				return rt.Value{}, e356
			}
			t357, e358 := rt.Sub(ctx, predel, rt.Number(1.0))
			if e358 != nil {
				return rt.Value{}, e358
			}
			t359, e360 := HodNachala(ctx, t357)
			if e360 != nil {
				return rt.Value{}, e360
			}
			// «ход»
			hod := t359
			for t361 := range t355 {
				// «знак»
				znak := t355[t361]
				t362, e363 := ShagObrezki(ctx, hod, znak)
				if e363 != nil {
					return rt.Value{}, e363
				}
				hod = t362
			}
			t364, e365 := rt.FieldGet(ctx, hod, "куски")
			if e365 != nil {
				return rt.Value{}, e365
			}
			// «соединить»
			t366, e367 := rt.BJoin(ctx, t364, rt.Text(""))
			if e367 != nil {
				return rt.Value{}, e367
			}
			t368, e369 := rt.Concat(ctx, t366, mnogotochie)
			if e369 != nil {
				return rt.Value{}, e369
			}
			t370, e371 := rt.Concat(ctx, t368, sbros)
			if e371 != nil {
				return rt.Value{}, e371
			}
			t352 = t370
		}
		t345 = t352
	}
	t372 := t345
	t373, e374 := ShirinaBezPosledovatelnostey(ctx, mnogotochie)
	if e374 != nil {
		return rt.Value{}, e374
	}
	t375, e376 := rt.Cond(ctx, rt.Flag(rt.Equal(t373, rt.Number(1.0))))
	if e376 != nil {
		return rt.Value{}, e376
	}
	var t377 rt.Value
	if t375 {
		t378, e379 := ShirinaBezPosledovatelnostey(ctx, sbros)
		if e379 != nil {
			return rt.Value{}, e379
		}
		t377 = rt.Flag(rt.Equal(t378, rt.Number(0.0)))
	} else {
		t377 = rt.Flag(false)
	}
	t380, e381 := rt.Cond(ctx, t377)
	if e381 != nil {
		return rt.Value{}, e381
	}
	var t382 rt.Value
	if t380 {
		t383, e384 := ShirinaBezPosledovatelnostey(ctx, t372)
		if e384 != nil {
			return rt.Value{}, e384
		}
		t385, e386 := rt.Gte(ctx, predel, rt.Number(0.0))
		if e386 != nil {
			return rt.Value{}, e386
		}
		t387, e388 := rt.Cond(ctx, t385)
		if e388 != nil {
			return rt.Value{}, e388
		}
		var t389 rt.Value
		if t387 {
			t389 = predel
		} else {
			t389 = rt.Number(0.0)
		}
		t390, e391 := rt.Lte(ctx, t383, t389)
		if e391 != nil {
			return rt.Value{}, e391
		}
		t382 = t390
	} else {
		t382 = rt.Flag(true)
	}
	// постусловие «обрезанное не шире предела»
	t392, e393 := rt.Post(ctx, t382, "обрезанное не шире предела", "Обрезать по ячейкам")
	if e393 != nil {
		return rt.Value{}, e393
	}
	if !t392 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «обрезанное не шире предела» функции «Обрезать по ячейкам»")
	}
	t394, e395 := PosledovatelnostiCely(ctx, tekst)
	if e395 != nil {
		return rt.Value{}, e395
	}
	t396, e397 := rt.Cond(ctx, t394)
	if e397 != nil {
		return rt.Value{}, e397
	}
	var t398 rt.Value
	if t396 {
		t399, e400 := PosledovatelnostiCely(ctx, mnogotochie)
		if e400 != nil {
			return rt.Value{}, e400
		}
		t398 = t399
	} else {
		t398 = rt.Flag(false)
	}
	t401, e402 := rt.Cond(ctx, t398)
	if e402 != nil {
		return rt.Value{}, e402
	}
	var t403 rt.Value
	if t401 {
		t404, e405 := PosledovatelnostiCely(ctx, sbros)
		if e405 != nil {
			return rt.Value{}, e405
		}
		t403 = t404
	} else {
		t403 = rt.Flag(false)
	}
	t406, e407 := rt.Cond(ctx, t403)
	if e407 != nil {
		return rt.Value{}, e407
	}
	var t408 rt.Value
	if t406 {
		t409, e410 := PosledovatelnostiCely(ctx, t372)
		if e410 != nil {
			return rt.Value{}, e410
		}
		t408 = t409
	} else {
		t408 = rt.Flag(true)
	}
	// постусловие «обрезка не рвёт последовательности»
	t411, e412 := rt.Post(ctx, t408, "обрезка не рвёт последовательности", "Обрезать по ячейкам")
	if e412 != nil {
		return rt.Value{}, e412
	}
	if !t411 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «обрезка не рвёт последовательности» функции «Обрезать по ячейкам»")
	}
	return t372, nil
}

// StrokaIliPusto — функция flang «Строка или пусто».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stroki — «строки»: список: строка.
// Параметр mesto — «место»: число.
// Результат — значение: строка.
func StrokaIliPusto(ctx *rt.Ctx, stroki rt.Value, mesto rt.Value) (rt.Value, error) {
	t413, e414 := rt.Gte(ctx, mesto, rt.Number(1.0))
	if e414 != nil {
		return rt.Value{}, e414
	}
	t415, e416 := rt.Cond(ctx, t413)
	if e416 != nil {
		return rt.Value{}, e416
	}
	var t417 rt.Value
	if t415 {
		// «длина»
		t418, e419 := rt.BLength(ctx, stroki)
		if e419 != nil {
			return rt.Value{}, e419
		}
		t420, e421 := rt.Lte(ctx, mesto, t418)
		if e421 != nil {
			return rt.Value{}, e421
		}
		t417 = t420
	} else {
		t417 = rt.Flag(false)
	}
	t422, e423 := rt.Cond(ctx, t417)
	if e423 != nil {
		return rt.Value{}, e423
	}
	if t422 {
		// «элемент»
		t424, e425 := rt.BElement(ctx, mesto, stroki)
		if e425 != nil {
			return rt.Value{}, e425
		}
		return t424, nil
	} else {
		return rt.Text(""), nil
	}
}

// OknoStrok — функция flang «Окно строк».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр stroki — «строки»: список: строка.
// Параметр nachalo — «начало»: число.
// Параметр skolko — «сколько»: «нат».
// Результат — значение: список: строка.
func OknoStrok(ctx *rt.Ctx, stroki rt.Value, nachalo rt.Value, skolko rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Окно строк"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t426, e427 := rt.Lte(ctx, skolko, rt.Number(0.0))
	if e427 != nil {
		return rt.Value{}, e427
	}
	t428, e429 := rt.Cond(ctx, t426)
	if e429 != nil {
		return rt.Value{}, e429
	}
	var t430 rt.Value
	if t428 {
		t430 = rt.List(nil)
	} else {
		t431, e432 := rt.Add(ctx, nachalo, skolko)
		if e432 != nil {
			return rt.Value{}, e432
		}
		t433, e434 := rt.Sub(ctx, t431, rt.Number(1.0))
		if e434 != nil {
			return rt.Value{}, e434
		}
		t435, e436 := StrokaIliPusto(ctx, stroki, t433)
		if e436 != nil {
			return rt.Value{}, e436
		}
		t437, e438 := rt.Sub(ctx, skolko, rt.Number(1.0))
		if e438 != nil {
			return rt.Value{}, e438
		}
		t439, e440 := OknoStrok(ctx, stroki, nachalo, t437)
		if e440 != nil {
			return rt.Value{}, e440
		}
		// «добавить»
		t441, e442 := rt.BAppend(ctx, t435, t439)
		if e442 != nil {
			return rt.Value{}, e442
		}
		t430 = t441
	}
	t443 := t430
	// «длина»
	t444, e445 := rt.BLength(ctx, t443)
	if e445 != nil {
		return rt.Value{}, e445
	}
	// постусловие «в окне ровно столько строк, сколько просили»
	t446, e447 := rt.Post(ctx, rt.Flag(rt.Equal(t444, skolko)), "в окне ровно столько строк, сколько просили", "Окно строк")
	if e447 != nil {
		return rt.Value{}, e447
	}
	if !t446 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «в окне ровно столько строк, сколько просили» функции «Окно строк»")
	}
	return t443, nil
}

// ProkrutkaVPredelah — функция flang «Прокрутка в пределах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stroki — «строки»: список: строка.
// Параметр prokrutka — «прокрутка»: число.
// Результат — значение: число.
func ProkrutkaVPredelah(ctx *rt.Ctx, stroki rt.Value, prokrutka rt.Value) (rt.Value, error) {
	// «длина»
	t448, e449 := rt.BLength(ctx, stroki)
	if e449 != nil {
		return rt.Value{}, e449
	}
	t450, e451 := rt.Sub(ctx, t448, rt.Number(1.0))
	if e451 != nil {
		return rt.Value{}, e451
	}
	t452, e453 := rt.Gt(ctx, prokrutka, t450)
	if e453 != nil {
		return rt.Value{}, e453
	}
	t454, e455 := rt.Cond(ctx, t452)
	if e455 != nil {
		return rt.Value{}, e455
	}
	var t456 rt.Value
	if t454 {
		// «длина»
		t457, e458 := rt.BLength(ctx, stroki)
		if e458 != nil {
			return rt.Value{}, e458
		}
		t459, e460 := rt.Sub(ctx, t457, rt.Number(1.0))
		if e460 != nil {
			return rt.Value{}, e460
		}
		t456 = t459
	} else {
		t456 = prokrutka
	}
	// пусть «дно»
	dno := t456
	t461, e462 := rt.Lt(ctx, dno, rt.Number(0.0))
	if e462 != nil {
		return rt.Value{}, e462
	}
	t463, e464 := rt.Cond(ctx, t461)
	if e464 != nil {
		return rt.Value{}, e464
	}
	var t465 rt.Value
	if t463 {
		t465 = rt.Number(0.0)
	} else {
		t465 = dno
	}
	t466 := t465
	t467, e468 := rt.Gte(ctx, t466, rt.Number(0.0))
	if e468 != nil {
		return rt.Value{}, e468
	}
	// постусловие «прокрутка не отрицательна»
	t469, e470 := rt.Post(ctx, t467, "прокрутка не отрицательна", "Прокрутка в пределах")
	if e470 != nil {
		return rt.Value{}, e470
	}
	if !t469 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «прокрутка не отрицательна» функции «Прокрутка в пределах»")
	}
	// «длина»
	t471, e472 := rt.BLength(ctx, stroki)
	if e472 != nil {
		return rt.Value{}, e472
	}
	t473, e474 := rt.Lte(ctx, t471, rt.Number(0.0))
	if e474 != nil {
		return rt.Value{}, e474
	}
	t475, e476 := rt.Cond(ctx, t473)
	if e476 != nil {
		return rt.Value{}, e476
	}
	var t477 rt.Value
	if t475 {
		t477 = rt.Number(0.0)
	} else {
		// «длина»
		t478, e479 := rt.BLength(ctx, stroki)
		if e479 != nil {
			return rt.Value{}, e479
		}
		t480, e481 := rt.Sub(ctx, t478, rt.Number(1.0))
		if e481 != nil {
			return rt.Value{}, e481
		}
		t477 = t480
	}
	t482, e483 := rt.Lte(ctx, t466, t477)
	if e483 != nil {
		return rt.Value{}, e483
	}
	// постусловие «прокрутка не заходит за последнюю строку»
	t484, e485 := rt.Post(ctx, t482, "прокрутка не заходит за последнюю строку", "Прокрутка в пределах")
	if e485 != nil {
		return rt.Value{}, e485
	}
	if !t484 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «прокрутка не заходит за последнюю строку» функции «Прокрутка в пределах»")
	}
	return t466, nil
}

// UlozhitRazdel — функция flang «Уложить раздел».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stroki — «строки»: список: строка.
// Параметр vysota — «высота»: «нат».
// Параметр prokrutka — «прокрутка»: число.
// Результат — значение: список: строка.
func UlozhitRazdel(ctx *rt.Ctx, stroki rt.Value, vysota rt.Value, prokrutka rt.Value) (rt.Value, error) {
	t486, e487 := ProkrutkaVPredelah(ctx, stroki, prokrutka)
	if e487 != nil {
		return rt.Value{}, e487
	}
	t488, e489 := rt.Add(ctx, t486, rt.Number(1.0))
	if e489 != nil {
		return rt.Value{}, e489
	}
	t490, e491 := OknoStrok(ctx, stroki, t488, vysota)
	if e491 != nil {
		return rt.Value{}, e491
	}
	t492 := t490
	// «длина»
	t493, e494 := rt.BLength(ctx, t492)
	if e494 != nil {
		return rt.Value{}, e494
	}
	// постусловие «раздел занимает ровно свою высоту»
	t495, e496 := rt.Post(ctx, rt.Flag(rt.Equal(t493, vysota)), "раздел занимает ровно свою высоту", "Уложить раздел")
	if e496 != nil {
		return rt.Value{}, e496
	}
	if !t495 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «раздел занимает ровно свою высоту» функции «Уложить раздел»")
	}
	return t492, nil
}

// Kadr — функция flang «Кадр».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stroki — «строки»: список: строка.
// Параметр shirina — «ширина»: число.
// Параметр vysota — «высота»: «нат».
// Параметр prokrutka — «прокрутка»: число.
// Параметр mnogotochie — «многоточие»: строка.
// Параметр sbros — «сброс»: строка.
// Результат — значение: список: строка.
func Kadr(ctx *rt.Ctx, stroki rt.Value, shirina rt.Value, vysota rt.Value, prokrutka rt.Value, mnogotochie rt.Value, sbros rt.Value) (rt.Value, error) {
	t497, e498 := UlozhitRazdel(ctx, stroki, vysota, prokrutka)
	if e498 != nil {
		return rt.Value{}, e498
	}
	t499, e500 := rt.RequireList(ctx, t497, "отобразить")
	if e500 != nil {
		return rt.Value{}, e500
	}
	t501 := make([]rt.Value, 0, len(t499))
	for t502 := range t499 {
		// «строчка»
		strochka := t499[t502]
		t503, e504 := ObrezatPoYacheykam(ctx, strochka, shirina, mnogotochie, sbros)
		if e504 != nil {
			return rt.Value{}, e504
		}
		t501 = append(t501, t503)
	}
	t505 := rt.List(t501)
	// «длина»
	t506, e507 := rt.BLength(ctx, t505)
	if e507 != nil {
		return rt.Value{}, e507
	}
	// постусловие «в кадре ровно столько строк, сколько высота»
	t508, e509 := rt.Post(ctx, rt.Flag(rt.Equal(t506, vysota)), "в кадре ровно столько строк, сколько высота", "Кадр")
	if e509 != nil {
		return rt.Value{}, e509
	}
	if !t508 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «в кадре ровно столько строк, сколько высота» функции «Кадр»")
	}
	return t505, nil
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Доля в пределах":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Доля в пределах", 1, len(args))
		}
		return DolyaVPredelah(ctx, args[0])
	case "Целое от":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Целое от", 1, len(args))
		}
		return CeloeOt(ctx, args[0])
	case "Клеток полосы":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Клеток полосы", 2, len(args))
		}
		return KletokPolosy(ctx, args[0], args[1])
	case "Клетки полосы":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Клетки полосы", 4, len(args))
		}
		return KletkiPolosy(ctx, args[0], args[1], args[2], args[3])
	case "Полоса":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Полоса", 4, len(args))
		}
		return Polosa(ctx, args[0], args[1], args[2], args[3])
	case "Эск":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Эск", 0, len(args))
		}
		return Esk(ctx)
	case "Закрывает последовательность":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Закрывает последовательность", 1, len(args))
		}
		return ZakryvaetPosledovatelnost(ctx, args[0])
	case "Ход начала":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Ход начала", 1, len(args))
		}
		return HodNachala(ctx, args[0])
	case "Оборвано ли":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Оборвано ли", 1, len(args))
		}
		return OborvanoLi(ctx, args[0])
	case "Шаг ширины":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Шаг ширины", 2, len(args))
		}
		return ShagShiriny(ctx, args[0], args[1])
	case "Ширина без последовательностей":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Ширина без последовательностей", 1, len(args))
		}
		return ShirinaBezPosledovatelnostey(ctx, args[0])
	case "Последовательности целы":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Последовательности целы", 1, len(args))
		}
		return PosledovatelnostiCely(ctx, args[0])
	case "Шаг обрезки":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Шаг обрезки", 2, len(args))
		}
		return ShagObrezki(ctx, args[0], args[1])
	case "Обрезать по ячейкам":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Обрезать по ячейкам", 4, len(args))
		}
		return ObrezatPoYacheykam(ctx, args[0], args[1], args[2], args[3])
	case "Строка или пусто":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Строка или пусто", 2, len(args))
		}
		return StrokaIliPusto(ctx, args[0], args[1])
	case "Окно строк":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Окно строк", 3, len(args))
		}
		return OknoStrok(ctx, args[0], args[1], args[2])
	case "Прокрутка в пределах":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прокрутка в пределах", 2, len(args))
		}
		return ProkrutkaVPredelah(ctx, args[0], args[1])
	case "Уложить раздел":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Уложить раздел", 3, len(args))
		}
		return UlozhitRazdel(ctx, args[0], args[1], args[2])
	case "Кадр":
		if len(args) != 6 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Кадр", 6, len(args))
		}
		return Kadr(ctx, args[0], args[1], args[2], args[3], args[4], args[5])
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
