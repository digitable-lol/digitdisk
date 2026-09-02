// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
// Модуль flang: «Colour».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flangcolour/flangrt"
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

// CeloeOt — функция flang «Целое от».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение: число.
func CeloeOt(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
	t1, e2 := rt.Mod(ctx, znachenie, rt.Number(1.0))
	if e2 != nil {
		return rt.Value{}, e2
	}
	t3, e4 := rt.Sub(ctx, znachenie, t1)
	if e4 != nil {
		return rt.Value{}, e4
	}
	return t3, nil
}

// DelenieNacelo — функция flang «Деление нацело».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр delimoe — «делимое»: число.
// Параметр delitel — «делитель»: число.
// Результат — значение: число.
func DelenieNacelo(ctx *rt.Ctx, delimoe rt.Value, delitel rt.Value) (rt.Value, error) {
	t5, e6 := rt.Cond(ctx, rt.Flag(rt.Equal(delitel, rt.Number(0.0))))
	if e6 != nil {
		return rt.Value{}, e6
	}
	if t5 {
		return rt.Number(0.0), nil
	} else {
		t7, e8 := rt.Div(ctx, delimoe, delitel)
		if e8 != nil {
			return rt.Value{}, e8
		}
		return CeloeOt(ctx, t7)
	}
}

// RaznicaPoModulyu — функция flang «Разница по модулю».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр pervoe — «первое»: число.
// Параметр vtoroe — «второе»: число.
// Результат — значение: число.
func RaznicaPoModulyu(ctx *rt.Ctx, pervoe rt.Value, vtoroe rt.Value) (rt.Value, error) {
	t9, e10 := rt.Gt(ctx, pervoe, vtoroe)
	if e10 != nil {
		return rt.Value{}, e10
	}
	t11, e12 := rt.Cond(ctx, t9)
	if e12 != nil {
		return rt.Value{}, e12
	}
	var t13 rt.Value
	if t11 {
		t14, e15 := rt.Sub(ctx, pervoe, vtoroe)
		if e15 != nil {
			return rt.Value{}, e15
		}
		t13 = t14
	} else {
		t16, e17 := rt.Sub(ctx, vtoroe, pervoe)
		if e17 != nil {
			return rt.Value{}, e17
		}
		t13 = t16
	}
	t18 := t13
	t19, e20 := rt.Gte(ctx, t18, rt.Number(0.0))
	if e20 != nil {
		return rt.Value{}, e20
	}
	// постусловие «разница неотрицательна»
	t21, e22 := rt.Post(ctx, t19, "разница неотрицательна", "Разница по модулю")
	if e22 != nil {
		return rt.Value{}, e22
	}
	if !t21 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «разница неотрицательна» функции «Разница по модулю»")
	}
	return t18, nil
}

// StupenKanala — функция flang «Ступень канала».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр kanal — «канал»: число.
// Результат — значение: число.
func StupenKanala(ctx *rt.Ctx, kanal rt.Value) (rt.Value, error) {
	t23, e24 := rt.Mul(ctx, kanal, rt.Number(5.0))
	if e24 != nil {
		return rt.Value{}, e24
	}
	t25, e26 := rt.Add(ctx, t23, rt.Number(127.0))
	if e26 != nil {
		return rt.Value{}, e26
	}
	return DelenieNacelo(ctx, t25, rt.Number(255.0))
}

// Kub256 — функция flang «Куб 256».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр krasnyy — «красный»: число.
// Параметр zelyonyy — «зелёный»: число.
// Параметр siniy — «синий»: число.
// Результат — значение: число.
func Kub256(ctx *rt.Ctx, krasnyy rt.Value, zelyonyy rt.Value, siniy rt.Value) (rt.Value, error) {
	t27, e28 := RaznicaPoModulyu(ctx, krasnyy, zelyonyy)
	if e28 != nil {
		return rt.Value{}, e28
	}
	t29, e30 := rt.Lt(ctx, t27, rt.Number(12.0))
	if e30 != nil {
		return rt.Value{}, e30
	}
	t31, e32 := rt.Cond(ctx, t29)
	if e32 != nil {
		return rt.Value{}, e32
	}
	var t33 rt.Value
	if t31 {
		t34, e35 := RaznicaPoModulyu(ctx, zelyonyy, siniy)
		if e35 != nil {
			return rt.Value{}, e35
		}
		t36, e37 := rt.Lt(ctx, t34, rt.Number(12.0))
		if e37 != nil {
			return rt.Value{}, e37
		}
		t33 = t36
	} else {
		t33 = rt.Flag(false)
	}
	t38, e39 := rt.Cond(ctx, t33)
	if e39 != nil {
		return rt.Value{}, e39
	}
	var t40 rt.Value
	if t38 {
		t41, e42 := RaznicaPoModulyu(ctx, krasnyy, siniy)
		if e42 != nil {
			return rt.Value{}, e42
		}
		t43, e44 := rt.Lt(ctx, t41, rt.Number(12.0))
		if e44 != nil {
			return rt.Value{}, e44
		}
		t40 = t43
	} else {
		t40 = rt.Flag(false)
	}
	t45, e46 := rt.Cond(ctx, t40)
	if e46 != nil {
		return rt.Value{}, e46
	}
	var t47 rt.Value
	if t45 {
		t48, e49 := rt.Add(ctx, krasnyy, zelyonyy)
		if e49 != nil {
			return rt.Value{}, e49
		}
		t50, e51 := rt.Add(ctx, t48, siniy)
		if e51 != nil {
			return rt.Value{}, e51
		}
		t52, e53 := DelenieNacelo(ctx, t50, rt.Number(3.0))
		if e53 != nil {
			return rt.Value{}, e53
		}
		// пусть «серый»
		seryy := t52
		t54, e55 := rt.Lt(ctx, seryy, rt.Number(8.0))
		if e55 != nil {
			return rt.Value{}, e55
		}
		t56, e57 := rt.Cond(ctx, t54)
		if e57 != nil {
			return rt.Value{}, e57
		}
		var t58 rt.Value
		if t56 {
			t58 = rt.Number(16.0)
		} else {
			t59, e60 := rt.Gt(ctx, seryy, rt.Number(248.0))
			if e60 != nil {
				return rt.Value{}, e60
			}
			t61, e62 := rt.Cond(ctx, t59)
			if e62 != nil {
				return rt.Value{}, e62
			}
			var t63 rt.Value
			if t61 {
				t63 = rt.Number(231.0)
			} else {
				t64, e65 := rt.Sub(ctx, seryy, rt.Number(8.0))
				if e65 != nil {
					return rt.Value{}, e65
				}
				t66, e67 := rt.Mul(ctx, t64, rt.Number(23.0))
				if e67 != nil {
					return rt.Value{}, e67
				}
				t68, e69 := DelenieNacelo(ctx, t66, rt.Number(240.0))
				if e69 != nil {
					return rt.Value{}, e69
				}
				t70, e71 := rt.Add(ctx, rt.Number(232.0), t68)
				if e71 != nil {
					return rt.Value{}, e71
				}
				t63 = t70
			}
			t58 = t63
		}
		t47 = t58
	} else {
		t72, e73 := StupenKanala(ctx, krasnyy)
		if e73 != nil {
			return rt.Value{}, e73
		}
		t74, e75 := rt.Mul(ctx, rt.Number(36.0), t72)
		if e75 != nil {
			return rt.Value{}, e75
		}
		t76, e77 := StupenKanala(ctx, zelyonyy)
		if e77 != nil {
			return rt.Value{}, e77
		}
		t78, e79 := rt.Mul(ctx, rt.Number(6.0), t76)
		if e79 != nil {
			return rt.Value{}, e79
		}
		t80, e81 := StupenKanala(ctx, siniy)
		if e81 != nil {
			return rt.Value{}, e81
		}
		t82, e83 := rt.Add(ctx, t78, t80)
		if e83 != nil {
			return rt.Value{}, e83
		}
		t84, e85 := rt.Add(ctx, t74, t82)
		if e85 != nil {
			return rt.Value{}, e85
		}
		t86, e87 := rt.Add(ctx, rt.Number(16.0), t84)
		if e87 != nil {
			return rt.Value{}, e87
		}
		t47 = t86
	}
	t88 := t47
	t89, e90 := rt.Gte(ctx, t88, rt.Number(16.0))
	if e90 != nil {
		return rt.Value{}, e90
	}
	t91, e92 := rt.Cond(ctx, t89)
	if e92 != nil {
		return rt.Value{}, e92
	}
	var t93 rt.Value
	if t91 {
		t94, e95 := rt.Lte(ctx, t88, rt.Number(255.0))
		if e95 != nil {
			return rt.Value{}, e95
		}
		t93 = t94
	} else {
		t93 = rt.Flag(false)
	}
	// постусловие «код лежит в палитре xterm-256»
	t96, e97 := rt.Post(ctx, t93, "код лежит в палитре xterm-256", "Куб 256")
	if e97 != nil {
		return rt.Value{}, e97
	}
	if !t96 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «код лежит в палитре xterm-256» функции «Куб 256»")
	}
	return t88, nil
}

// Esk — функция flang «Эск».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: строка.
func Esk(ctx *rt.Ctx) (rt.Value, error) {
	// «символ по коду»
	t98, e99 := rt.BCharFromCode(ctx, rt.Number(27.0))
	if e99 != nil {
		return rt.Value{}, e99
	}
	return t98, nil
}

// CvetIstinnyy — функция flang «Цвет истинный».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр fonom — «фоном».
// Параметр krasnyy — «красный»: число.
// Параметр zelyonyy — «зелёный»: число.
// Параметр siniy — «синий»: число.
// Результат — значение: строка.
func CvetIstinnyy(ctx *rt.Ctx, fonom rt.Value, krasnyy rt.Value, zelyonyy rt.Value, siniy rt.Value) (rt.Value, error) {
	t100 := make([]rt.Value, 10)
	t101, e102 := Esk(ctx)
	if e102 != nil {
		return rt.Value{}, e102
	}
	t100[0] = t101
	t100[1] = rt.Text("[")
	t103, e104 := rt.Cond(ctx, fonom)
	if e104 != nil {
		return rt.Value{}, e104
	}
	var t105 rt.Value
	if t103 {
		t105 = rt.Number(48.0)
	} else {
		t105 = rt.Number(38.0)
	}
	// «к строке»
	t106, e107 := rt.BToString(ctx, t105)
	if e107 != nil {
		return rt.Value{}, e107
	}
	t100[2] = t106
	t100[3] = rt.Text(";2;")
	// «к строке»
	t108, e109 := rt.BToString(ctx, krasnyy)
	if e109 != nil {
		return rt.Value{}, e109
	}
	t100[4] = t108
	t100[5] = rt.Text(";")
	// «к строке»
	t110, e111 := rt.BToString(ctx, zelyonyy)
	if e111 != nil {
		return rt.Value{}, e111
	}
	t100[6] = t110
	t100[7] = rt.Text(";")
	// «к строке»
	t112, e113 := rt.BToString(ctx, siniy)
	if e113 != nil {
		return rt.Value{}, e113
	}
	t100[8] = t112
	t100[9] = rt.Text("m")
	// «соединить»
	t114, e115 := rt.BJoin(ctx, rt.List(t100), rt.Text(""))
	if e115 != nil {
		return rt.Value{}, e115
	}
	t116 := t114
	// «длина»
	t117, e118 := rt.BLength(ctx, t116)
	if e118 != nil {
		return rt.Value{}, e118
	}
	// «длина»
	t119, e120 := rt.BLength(ctx, t116)
	if e120 != nil {
		return rt.Value{}, e120
	}
	// «подстрока»
	t121, e122 := rt.BSubstring(ctx, t116, t117, t119)
	if e122 != nil {
		return rt.Value{}, e122
	}
	// постусловие «последовательность кончается буквой m»
	t123, e124 := rt.Post(ctx, rt.Flag(rt.Equal(t121, rt.Text("m"))), "последовательность кончается буквой m", "Цвет истинный")
	if e124 != nil {
		return rt.Value{}, e124
	}
	if !t123 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «последовательность кончается буквой m» функции «Цвет истинный»")
	}
	return t116, nil
}

// CvetIzPalitry — функция flang «Цвет из палитры».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр fonom — «фоном».
// Параметр krasnyy — «красный»: число.
// Параметр zelyonyy — «зелёный»: число.
// Параметр siniy — «синий»: число.
// Результат — значение: строка.
func CvetIzPalitry(ctx *rt.Ctx, fonom rt.Value, krasnyy rt.Value, zelyonyy rt.Value, siniy rt.Value) (rt.Value, error) {
	t125 := make([]rt.Value, 6)
	t126, e127 := Esk(ctx)
	if e127 != nil {
		return rt.Value{}, e127
	}
	t125[0] = t126
	t125[1] = rt.Text("[")
	t128, e129 := rt.Cond(ctx, fonom)
	if e129 != nil {
		return rt.Value{}, e129
	}
	var t130 rt.Value
	if t128 {
		t130 = rt.Number(48.0)
	} else {
		t130 = rt.Number(38.0)
	}
	// «к строке»
	t131, e132 := rt.BToString(ctx, t130)
	if e132 != nil {
		return rt.Value{}, e132
	}
	t125[2] = t131
	t125[3] = rt.Text(";5;")
	t133, e134 := Kub256(ctx, krasnyy, zelyonyy, siniy)
	if e134 != nil {
		return rt.Value{}, e134
	}
	// «к строке»
	t135, e136 := rt.BToString(ctx, t133)
	if e136 != nil {
		return rt.Value{}, e136
	}
	t125[4] = t135
	t125[5] = rt.Text("m")
	// «соединить»
	t137, e138 := rt.BJoin(ctx, rt.List(t125), rt.Text(""))
	if e138 != nil {
		return rt.Value{}, e138
	}
	t139 := t137
	// «длина»
	t140, e141 := rt.BLength(ctx, t139)
	if e141 != nil {
		return rt.Value{}, e141
	}
	// «длина»
	t142, e143 := rt.BLength(ctx, t139)
	if e143 != nil {
		return rt.Value{}, e143
	}
	// «подстрока»
	t144, e145 := rt.BSubstring(ctx, t139, t140, t142)
	if e145 != nil {
		return rt.Value{}, e145
	}
	// постусловие «последовательность кончается буквой m»
	t146, e147 := rt.Post(ctx, rt.Flag(rt.Equal(t144, rt.Text("m"))), "последовательность кончается буквой m", "Цвет из палитры")
	if e147 != nil {
		return rt.Value{}, e147
	}
	if !t146 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «последовательность кончается буквой m» функции «Цвет из палитры»")
	}
	return t139, nil
}

// CvetIzShestnadcati — функция flang «Цвет из шестнадцати».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр fonom — «фоном».
// Параметр ansiKod — «ансиКод»: число.
// Результат — значение: строка.
func CvetIzShestnadcati(ctx *rt.Ctx, fonom rt.Value, ansiKod rt.Value) (rt.Value, error) {
	t148 := make([]rt.Value, 4)
	t149, e150 := Esk(ctx)
	if e150 != nil {
		return rt.Value{}, e150
	}
	t148[0] = t149
	t148[1] = rt.Text("[")
	t151, e152 := rt.Cond(ctx, fonom)
	if e152 != nil {
		return rt.Value{}, e152
	}
	var t153 rt.Value
	if t151 {
		t154, e155 := rt.Add(ctx, ansiKod, rt.Number(10.0))
		if e155 != nil {
			return rt.Value{}, e155
		}
		t153 = t154
	} else {
		t153 = ansiKod
	}
	// «к строке»
	t156, e157 := rt.BToString(ctx, t153)
	if e157 != nil {
		return rt.Value{}, e157
	}
	t148[2] = t156
	t148[3] = rt.Text("m")
	// «соединить»
	t158, e159 := rt.BJoin(ctx, rt.List(t148), rt.Text(""))
	if e159 != nil {
		return rt.Value{}, e159
	}
	return t158, nil
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Целое от":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Целое от", 1, len(args))
		}
		return CeloeOt(ctx, args[0])
	case "Деление нацело":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Деление нацело", 2, len(args))
		}
		return DelenieNacelo(ctx, args[0], args[1])
	case "Разница по модулю":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разница по модулю", 2, len(args))
		}
		return RaznicaPoModulyu(ctx, args[0], args[1])
	case "Ступень канала":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Ступень канала", 1, len(args))
		}
		return StupenKanala(ctx, args[0])
	case "Куб 256":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Куб 256", 3, len(args))
		}
		return Kub256(ctx, args[0], args[1], args[2])
	case "Эск":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Эск", 0, len(args))
		}
		return Esk(ctx)
	case "Цвет истинный":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Цвет истинный", 4, len(args))
		}
		return CvetIstinnyy(ctx, args[0], args[1], args[2], args[3])
	case "Цвет из палитры":
		if len(args) != 4 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Цвет из палитры", 4, len(args))
		}
		return CvetIzPalitry(ctx, args[0], args[1], args[2], args[3])
	case "Цвет из шестнадцати":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Цвет из шестнадцати", 2, len(args))
		}
		return CvetIzShestnadcati(ctx, args[0], args[1])
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
