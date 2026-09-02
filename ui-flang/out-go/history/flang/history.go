// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
// Модуль flang: «History».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flanghistory/flangrt"
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

// VariantSpokoyno — вариант «Спокойно» суммы типов «Уровень».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantSpokoyno() rt.Value {
	return rt.Variant("Спокойно", nil)
}

// VariantVnimanie — вариант «Внимание» суммы типов «Уровень».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantVnimanie() rt.Value {
	return rt.Variant("Внимание", nil)
}

// VariantTrevoga — вариант «Тревога» суммы типов «Уровень».
//
// Дискриминант — имя варианта; проверяется через rt.VariantIs(значение, «Имя»).
// Приставка Variant в имени — это роль: функция flang с тем же именем даёт
// идентификатор без приставки, и одно объявление не спорит с другим.
func VariantTrevoga() rt.Value {
	return rt.Variant("Тревога", nil)
}

// HvostSpiska — функция flang «Хвост списка».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Хвостовой самовызов развёрнут в цикл: стек не растёт.
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр zamery — «замеры»: список: число.
// Параметр skolko — «сколько»: число.
// Результат — значение: список: число.
func HvostSpiska(ctx *rt.Ctx, zamery rt.Value, skolko rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Хвост списка"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	for {
		if rt.ChainEmpty(zamery) {
			return rt.List(nil), nil
		} else if rt.ChainCons(zamery) {
			// голова «голова»
			golova := rt.ChainHead(zamery)
			// хвост «хвост»
			hvost := rt.ChainTail(zamery)
			_ = golova
			// «длина»
			t1, e2 := rt.BLength(ctx, zamery)
			if e2 != nil {
				return rt.Value{}, e2
			}
			t3, e4 := rt.Gte(ctx, skolko, rt.Number(0.0))
			if e4 != nil {
				return rt.Value{}, e4
			}
			t5, e6 := rt.Cond(ctx, t3)
			if e6 != nil {
				return rt.Value{}, e6
			}
			var t7 rt.Value
			if t5 {
				t7 = skolko
			} else {
				t7 = rt.Number(0.0)
			}
			t8, e9 := rt.Lte(ctx, t1, t7)
			if e9 != nil {
				return rt.Value{}, e9
			}
			t10, e11 := rt.Cond(ctx, t8)
			if e11 != nil {
				return rt.Value{}, e11
			}
			if t10 {
				return zamery, nil
			} else {
				t12 := skolko
				zamery = hvost
				skolko = t12
				// виток цикла — тоже шаг вычисления: незавершающийся хвостовой
				// самовызов обязан упереться в лимит, а не крутиться вечно
				if e13 := ctx.Step("Хвост списка"); e13 != nil {
					return rt.Value{}, e13
				}
				continue
			}
		} else {
			return rt.Value{}, rt.MatchFail(ctx, zamery)
		}
	}
}

// DopisatZamer — функция flang «Дописать замер».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zamery — «замеры»: список: число.
// Параметр zamer — «замер»: число.
// Параметр predel — «предел»: число.
// Результат — значение: список: число.
func DopisatZamer(ctx *rt.Ctx, zamery rt.Value, zamer rt.Value, predel rt.Value) (rt.Value, error) {
	// «добавить»
	t14, e15 := rt.BAppend(ctx, zamer, zamery)
	if e15 != nil {
		return rt.Value{}, e15
	}
	t16, e17 := rt.Gt(ctx, predel, rt.Number(0.0))
	if e17 != nil {
		return rt.Value{}, e17
	}
	t18, e19 := rt.Cond(ctx, t16)
	if e19 != nil {
		return rt.Value{}, e19
	}
	var t20 rt.Value
	if t18 {
		t20 = predel
	} else {
		t20 = rt.Number(1.0)
	}
	t21, e22 := HvostSpiska(ctx, t14, t20)
	if e22 != nil {
		return rt.Value{}, e22
	}
	t23 := t21
	t24, e25 := rt.Gt(ctx, predel, rt.Number(0.0))
	if e25 != nil {
		return rt.Value{}, e25
	}
	t26, e27 := rt.Cond(ctx, t24)
	if e27 != nil {
		return rt.Value{}, e27
	}
	var t28 rt.Value
	if t26 {
		// «длина»
		t29, e30 := rt.BLength(ctx, t23)
		if e30 != nil {
			return rt.Value{}, e30
		}
		t31, e32 := rt.Lte(ctx, t29, predel)
		if e32 != nil {
			return rt.Value{}, e32
		}
		t28 = t31
	} else {
		t28 = rt.Flag(true)
	}
	// постусловие «история не длиннее предела»
	t33, e34 := rt.Post(ctx, t28, "история не длиннее предела", "Дописать замер")
	if e34 != nil {
		return rt.Value{}, e34
	}
	if !t33 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «история не длиннее предела» функции «Дописать замер»")
	}
	// «длина»
	t35, e36 := rt.BLength(ctx, t23)
	if e36 != nil {
		return rt.Value{}, e36
	}
	t37, e38 := rt.Gt(ctx, t35, rt.Number(0.0))
	if e38 != nil {
		return rt.Value{}, e38
	}
	// постусловие «история не пуста»
	t39, e40 := rt.Post(ctx, t37, "история не пуста", "Дописать замер")
	if e40 != nil {
		return rt.Value{}, e40
	}
	if !t39 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «история не пуста» функции «Дописать замер»")
	}
	return t23, nil
}

// ZnakiGrafika — функция flang «Знаки графика».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: список: строка.
func ZnakiGrafika(ctx *rt.Ctx) (rt.Value, error) {
	t41 := make([]rt.Value, 8)
	t41[0] = rt.Text("▁")
	t41[1] = rt.Text("▂")
	t41[2] = rt.Text("▃")
	t41[3] = rt.Text("▄")
	t41[4] = rt.Text("▅")
	t41[5] = rt.Text("▆")
	t41[6] = rt.Text("▇")
	t41[7] = rt.Text("█")
	return rt.List(t41), nil
}

// ZnakPropuska — функция flang «Знак пропуска».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: строка.
func ZnakPropuska(ctx *rt.Ctx) (rt.Value, error) {
	return rt.Text("·"), nil
}

// CeloeKNulyu — функция flang «Целое к нулю».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение: число.
func CeloeKNulyu(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
	t42, e43 := rt.Mod(ctx, znachenie, rt.Number(1.0))
	if e43 != nil {
		return rt.Value{}, e43
	}
	t44, e45 := rt.Sub(ctx, znachenie, t42)
	if e45 != nil {
		return rt.Value{}, e45
	}
	return t44, nil
}

// ZnakZamera — функция flang «Знак замера».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Результат — значение: строка.
func ZnakZamera(ctx *rt.Ctx, dolya rt.Value) (rt.Value, error) {
	t46, e47 := rt.Lt(ctx, dolya, rt.Number(0.0))
	if e47 != nil {
		return rt.Value{}, e47
	}
	t48, e49 := rt.Cond(ctx, t46)
	if e49 != nil {
		return rt.Value{}, e49
	}
	var t50 rt.Value
	if t48 {
		t51, e52 := ZnakPropuska(ctx)
		if e52 != nil {
			return rt.Value{}, e52
		}
		t50 = t51
	} else {
		t53, e54 := rt.Gt(ctx, dolya, rt.Number(1.0))
		if e54 != nil {
			return rt.Value{}, e54
		}
		t55, e56 := rt.Cond(ctx, t53)
		if e56 != nil {
			return rt.Value{}, e56
		}
		var t57 rt.Value
		if t55 {
			t57 = rt.Number(1.0)
		} else {
			t57 = dolya
		}
		t58, e59 := rt.Mul(ctx, t57, rt.Number(7.0))
		if e59 != nil {
			return rt.Value{}, e59
		}
		t60, e61 := CeloeKNulyu(ctx, t58)
		if e61 != nil {
			return rt.Value{}, e61
		}
		t62, e63 := rt.Add(ctx, rt.Number(1.0), t60)
		if e63 != nil {
			return rt.Value{}, e63
		}
		t64, e65 := ZnakiGrafika(ctx)
		if e65 != nil {
			return rt.Value{}, e65
		}
		// «элемент»
		t66, e67 := rt.BElement(ctx, t62, t64)
		if e67 != nil {
			return rt.Value{}, e67
		}
		t50 = t66
	}
	t68 := t50
	// «длина»
	t69, e70 := rt.BLength(ctx, t68)
	if e70 != nil {
		return rt.Value{}, e70
	}
	// постусловие «знак — один символ»
	t71, e72 := rt.Post(ctx, rt.Flag(rt.Equal(t69, rt.Number(1.0))), "знак — один символ", "Знак замера")
	if e72 != nil {
		return rt.Value{}, e72
	}
	if !t71 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «знак — один символ» функции «Знак замера»")
	}
	return t68, nil
}

// Propuski — функция flang «Пропуски».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр skolko — «сколько»: число.
// Результат — значение: строка.
func Propuski(ctx *rt.Ctx, skolko rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Пропуски"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t73, e74 := rt.Lte(ctx, skolko, rt.Number(0.0))
	if e74 != nil {
		return rt.Value{}, e74
	}
	t75, e76 := rt.Cond(ctx, t73)
	if e76 != nil {
		return rt.Value{}, e76
	}
	var t77 rt.Value
	if t75 {
		t77 = rt.Text("")
	} else {
		t78, e79 := ZnakPropuska(ctx)
		if e79 != nil {
			return rt.Value{}, e79
		}
		t80, e81 := rt.Sub(ctx, skolko, rt.Number(1.0))
		if e81 != nil {
			return rt.Value{}, e81
		}
		// пусть «шаг»
		shag := t80
		t82, e83 := rt.Lt(ctx, shag, skolko)
		if e83 != nil {
			return rt.Value{}, e83
		}
		t84, e85 := rt.Cond(ctx, t82)
		if e85 != nil {
			return rt.Value{}, e85
		}
		var t86 rt.Value
		if t84 {
			t86 = shag
		} else {
			t87, e88 := MeraUbyvaet(ctx, shag, skolko)
			if e88 != nil {
				return rt.Value{}, e88
			}
			t86 = t87
		}
		t89, e90 := Propuski(ctx, t86)
		if e90 != nil {
			return rt.Value{}, e90
		}
		t91, e92 := rt.Concat(ctx, t78, t89)
		if e92 != nil {
			return rt.Value{}, e92
		}
		t77 = t91
	}
	t93 := t77
	// «длина»
	t94, e95 := rt.BLength(ctx, t93)
	if e95 != nil {
		return rt.Value{}, e95
	}
	t96, e97 := rt.Gte(ctx, skolko, rt.Number(0.0))
	if e97 != nil {
		return rt.Value{}, e97
	}
	t98, e99 := rt.Cond(ctx, t96)
	if e99 != nil {
		return rt.Value{}, e99
	}
	var t100 rt.Value
	if t98 {
		t100 = skolko
	} else {
		t100 = rt.Number(0.0)
	}
	// постусловие «пропусков ровно столько, сколько просили»
	t101, e102 := rt.Post(ctx, rt.Flag(rt.Equal(t94, t100)), "пропусков ровно столько, сколько просили", "Пропуски")
	if e102 != nil {
		return rt.Value{}, e102
	}
	if !t101 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «пропусков ровно столько, сколько просили» функции «Пропуски»")
	}
	return t93, nil
}

// Grafik — функция flang «График».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zamery — «замеры»: список: число.
// Параметр shirina — «ширина»: число.
// Результат — значение: строка.
func Grafik(ctx *rt.Ctx, zamery rt.Value, shirina rt.Value) (rt.Value, error) {
	t103, e104 := rt.Lte(ctx, shirina, rt.Number(0.0))
	if e104 != nil {
		return rt.Value{}, e104
	}
	t105, e106 := rt.Cond(ctx, t103)
	if e106 != nil {
		return rt.Value{}, e106
	}
	var t107 rt.Value
	if t105 {
		t107 = rt.Text("")
	} else {
		t108, e109 := HvostSpiska(ctx, zamery, shirina)
		if e109 != nil {
			return rt.Value{}, e109
		}
		// пусть «окно»
		okno := t108
		// «длина»
		t110, e111 := rt.BLength(ctx, okno)
		if e111 != nil {
			return rt.Value{}, e111
		}
		t112, e113 := rt.Sub(ctx, shirina, t110)
		if e113 != nil {
			return rt.Value{}, e113
		}
		t114, e115 := Propuski(ctx, t112)
		if e115 != nil {
			return rt.Value{}, e115
		}
		t116, e117 := rt.RequireList(ctx, okno, "отобразить")
		if e117 != nil {
			return rt.Value{}, e117
		}
		t118 := make([]rt.Value, 0, len(t116))
		for t119 := range t116 {
			// «замер»
			zamer := t116[t119]
			t120, e121 := ZnakZamera(ctx, zamer)
			if e121 != nil {
				return rt.Value{}, e121
			}
			t118 = append(t118, t120)
		}
		// «соединить»
		t122, e123 := rt.BJoin(ctx, rt.List(t118), rt.Text(""))
		if e123 != nil {
			return rt.Value{}, e123
		}
		t124, e125 := rt.Concat(ctx, t114, t122)
		if e125 != nil {
			return rt.Value{}, e125
		}
		t107 = t124
	}
	t126 := t107
	// «длина»
	t127, e128 := rt.BLength(ctx, t126)
	if e128 != nil {
		return rt.Value{}, e128
	}
	t129, e130 := rt.Gte(ctx, shirina, rt.Number(0.0))
	if e130 != nil {
		return rt.Value{}, e130
	}
	t131, e132 := rt.Cond(ctx, t129)
	if e132 != nil {
		return rt.Value{}, e132
	}
	var t133 rt.Value
	if t131 {
		t133 = shirina
	} else {
		t133 = rt.Number(0.0)
	}
	// постусловие «график ровно в заказанную ширину»
	t134, e135 := rt.Post(ctx, rt.Flag(rt.Equal(t127, t133)), "график ровно в заказанную ширину", "График")
	if e135 != nil {
		return rt.Value{}, e135
	}
	if !t134 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «график ровно в заказанную ширину» функции «График»")
	}
	return t126, nil
}

// UrovenDoli — функция flang «Уровень доли».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Результат — значение: «Уровень».
func UrovenDoli(ctx *rt.Ctx, dolya rt.Value) (rt.Value, error) {
	t136, e137 := rt.Gte(ctx, dolya, rt.Number(0.9))
	if e137 != nil {
		return rt.Value{}, e137
	}
	t138, e139 := rt.Cond(ctx, t136)
	if e139 != nil {
		return rt.Value{}, e139
	}
	if t138 {
		return rt.Variant("Тревога", nil), nil
	} else {
		t140, e141 := rt.Gte(ctx, dolya, rt.Number(0.75))
		if e141 != nil {
			return rt.Value{}, e141
		}
		t142, e143 := rt.Cond(ctx, t140)
		if e143 != nil {
			return rt.Value{}, e143
		}
		if t142 {
			return rt.Variant("Внимание", nil), nil
		} else {
			return rt.Variant("Спокойно", nil), nil
		}
	}
}

// MeraUbyvaet — функция flang «мера убывает».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр shag — «шаг»: число.
// Параметр mera — «мера»: число.
// Результат — значение: число.
func MeraUbyvaet(ctx *rt.Ctx, shag rt.Value, mera rt.Value) (rt.Value, error) {
	t144 := shag
	t145, e146 := rt.Lt(ctx, t144, mera)
	if e146 != nil {
		return rt.Value{}, e146
	}
	// постусловие «мера убывает»
	t147, e148 := rt.Post(ctx, t145, "мера убывает", "мера убывает")
	if e148 != nil {
		return rt.Value{}, e148
	}
	if !t147 {
		return rt.Value{}, rt.Fail("FLANG_MEASURE", "%s", "тотальная функция «Пропуски»: мера не убыла — аргумент 1 вызова «Пропуски» не стал меньше параметра «сколько». Завершение доказано убыванием этой меры, а числа flang — IEEE-754 double: при большом |«сколько»| постоянный шаг не меняет значение, и спуск не идёт. Отказ здесь честнее зацикливания")
	}
	return t144, nil
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Хвост списка":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Хвост списка", 2, len(args))
		}
		return HvostSpiska(ctx, args[0], args[1])
	case "Дописать замер":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Дописать замер", 3, len(args))
		}
		return DopisatZamer(ctx, args[0], args[1], args[2])
	case "Знаки графика":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Знаки графика", 0, len(args))
		}
		return ZnakiGrafika(ctx)
	case "Знак пропуска":
		if len(args) != 0 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Знак пропуска", 0, len(args))
		}
		return ZnakPropuska(ctx)
	case "Целое к нулю":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Целое к нулю", 1, len(args))
		}
		return CeloeKNulyu(ctx, args[0])
	case "Знак замера":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Знак замера", 1, len(args))
		}
		return ZnakZamera(ctx, args[0])
	case "Пропуски":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Пропуски", 1, len(args))
		}
		return Propuski(ctx, args[0])
	case "График":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"График", 2, len(args))
		}
		return Grafik(ctx, args[0], args[1])
	case "Уровень доли":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Уровень доли", 1, len(args))
		}
		return UrovenDoli(ctx, args[0])
	case "мера убывает":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"мера убывает", 2, len(args))
		}
		return MeraUbyvaet(ctx, args[0], args[1])
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
