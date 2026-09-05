// Сгенерировано flang (бэкенд Go, flang/self/emit-go.flang). Не редактировать руками.
// Модуль flang: «Format».
// Файл: реализация: функции, конструкторы значений, вызов по имени.
// Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.

// Контракт вызова: функция возвращает значение и nil либо нулевое значение
// и ошибку *flangrt.Error с кодом и текстом, дословно совпадающими с
// интерпретатором flang. Паник здесь нет: диагностика — это значение.
package flang

import (
	rt "flangformat/flangrt"
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

// SozdatHodRazryadov — запись FTS «Ход разрядов»: «собрано», «осталось», «разделитель».
//
// Запись flang тотальна: пропущенное поле — это «ничто», а не дырка.
func SozdatHodRazryadov(sobrano rt.Value, ostalos rt.Value, razdelitel rt.Value) rt.Value {
	return rt.Record([]rt.Field{
		{Name: "собрано", Value: sobrano},
		{Name: "осталось", Value: ostalos},
		{Name: "разделитель", Value: razdelitel},
	})
}

// CeloeKNulyu — функция flang «Целое к нулю».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение: число.
func CeloeKNulyu(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
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

// CeloeVniz — функция flang «Целое вниз».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение: число.
func CeloeVniz(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
	t5, e6 := CeloeKNulyu(ctx, znachenie)
	if e6 != nil {
		return rt.Value{}, e6
	}
	// пусть «к нулю»
	kNulyu := t5
	t7, e8 := rt.Gt(ctx, kNulyu, znachenie)
	if e8 != nil {
		return rt.Value{}, e8
	}
	t9, e10 := rt.Cond(ctx, t7)
	if e10 != nil {
		return rt.Value{}, e10
	}
	var t11 rt.Value
	if t9 {
		t12, e13 := rt.Sub(ctx, kNulyu, rt.Number(1.0))
		if e13 != nil {
			return rt.Value{}, e13
		}
		t11 = t12
	} else {
		t11 = kNulyu
	}
	t14 := t11
	t15, e16 := rt.Lte(ctx, t14, znachenie)
	if e16 != nil {
		return rt.Value{}, e16
	}
	// постусловие «целое вниз не больше самого числа»
	t17, e18 := rt.Post(ctx, t15, "целое вниз не больше самого числа", "Целое вниз")
	if e18 != nil {
		return rt.Value{}, e18
	}
	if !t17 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «целое вниз не больше самого числа» функции «Целое вниз»")
	}
	t19, e20 := rt.Sub(ctx, znachenie, t14)
	if e20 != nil {
		return rt.Value{}, e20
	}
	t21, e22 := rt.Lt(ctx, t19, rt.Number(1.0))
	if e22 != nil {
		return rt.Value{}, e22
	}
	// постусловие «целое вниз ближе единицы»
	t23, e24 := rt.Post(ctx, t21, "целое вниз ближе единицы", "Целое вниз")
	if e24 != nil {
		return rt.Value{}, e24
	}
	if !t23 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «целое вниз ближе единицы» функции «Целое вниз»")
	}
	return t14, nil
}

// DelenieNacelo — функция flang «Деление нацело».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр delimoe — «делимое»: число.
// Параметр delitel — «делитель»: число.
// Результат — значение: число.
func DelenieNacelo(ctx *rt.Ctx, delimoe rt.Value, delitel rt.Value) (rt.Value, error) {
	t25, e26 := rt.Cond(ctx, rt.Flag(rt.Equal(delitel, rt.Number(0.0))))
	if e26 != nil {
		return rt.Value{}, e26
	}
	if t25 {
		return rt.Number(0.0), nil
	} else {
		t27, e28 := rt.Div(ctx, delimoe, delitel)
		if e28 != nil {
			return rt.Value{}, e28
		}
		return CeloeKNulyu(ctx, t27)
	}
}

// DesyatVStepeni — функция flang «Десять в степени».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр stepen — «степень»: «нат».
// Результат — значение: число.
func DesyatVStepeni(ctx *rt.Ctx, stepen rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Десять в степени"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t29, e30 := rt.Lte(ctx, stepen, rt.Number(0.0))
	if e30 != nil {
		return rt.Value{}, e30
	}
	t31, e32 := rt.Cond(ctx, t29)
	if e32 != nil {
		return rt.Value{}, e32
	}
	var t33 rt.Value
	if t31 {
		t33 = rt.Number(1.0)
	} else {
		t34, e35 := rt.Sub(ctx, stepen, rt.Number(1.0))
		if e35 != nil {
			return rt.Value{}, e35
		}
		t36, e37 := DesyatVStepeni(ctx, t34)
		if e37 != nil {
			return rt.Value{}, e37
		}
		t38, e39 := rt.Mul(ctx, rt.Number(10.0), t36)
		if e39 != nil {
			return rt.Value{}, e39
		}
		t33 = t38
	}
	t40 := t33
	t41, e42 := rt.Gt(ctx, t40, rt.Number(0.0))
	if e42 != nil {
		return rt.Value{}, e42
	}
	// постусловие «степень десяти положительна»
	t43, e44 := rt.Post(ctx, t41, "степень десяти положительна", "Десять в степени")
	if e44 != nil {
		return rt.Value{}, e44
	}
	if !t43 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «степень десяти положительна» функции «Десять в степени»")
	}
	return t40, nil
}

// Povtorit — функция flang «Повторить».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр znak — «знак»: строка.
// Параметр skolko — «сколько»: число.
// Результат — значение: строка.
func Povtorit(ctx *rt.Ctx, znak rt.Value, skolko rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Повторить"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t45, e46 := rt.Lte(ctx, skolko, rt.Number(0.0))
	if e46 != nil {
		return rt.Value{}, e46
	}
	t47, e48 := rt.Cond(ctx, t45)
	if e48 != nil {
		return rt.Value{}, e48
	}
	var t49 rt.Value
	if t47 {
		t49 = rt.Text("")
	} else {
		t50, e51 := rt.Sub(ctx, skolko, rt.Number(1.0))
		if e51 != nil {
			return rt.Value{}, e51
		}
		// пусть «шаг»
		shag := t50
		t52, e53 := rt.Lt(ctx, shag, skolko)
		if e53 != nil {
			return rt.Value{}, e53
		}
		t54, e55 := rt.Cond(ctx, t52)
		if e55 != nil {
			return rt.Value{}, e55
		}
		var t56 rt.Value
		if t54 {
			t56 = shag
		} else {
			t57, e58 := MeraUbyvaet(ctx, shag, skolko)
			if e58 != nil {
				return rt.Value{}, e58
			}
			t56 = t57
		}
		t59, e60 := Povtorit(ctx, znak, t56)
		if e60 != nil {
			return rt.Value{}, e60
		}
		t61, e62 := rt.Concat(ctx, znak, t59)
		if e62 != nil {
			return rt.Value{}, e62
		}
		t49 = t61
	}
	t63 := t49
	// «длина»
	t64, e65 := rt.BLength(ctx, znak)
	if e65 != nil {
		return rt.Value{}, e65
	}
	t66, e67 := rt.Cond(ctx, rt.Flag(rt.Equal(t64, rt.Number(0.0))))
	if e67 != nil {
		return rt.Value{}, e67
	}
	var t68 rt.Value
	if t66 {
		// «длина»
		t69, e70 := rt.BLength(ctx, t63)
		if e70 != nil {
			return rt.Value{}, e70
		}
		t68 = rt.Flag(rt.Equal(t69, rt.Number(0.0)))
	} else {
		t68 = rt.Flag(true)
	}
	// постусловие «повтор пустого знака пуст»
	t71, e72 := rt.Post(ctx, t68, "повтор пустого знака пуст", "Повторить")
	if e72 != nil {
		return rt.Value{}, e72
	}
	if !t71 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «повтор пустого знака пуст» функции «Повторить»")
	}
	return t63, nil
}

// SlevaNulyami — функция flang «Слева нулями».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр shirina — «ширина»: «нат».
// Результат — значение: строка.
func SlevaNulyami(ctx *rt.Ctx, tekst rt.Value, shirina rt.Value) (rt.Value, error) {
	// «длина»
	t73, e74 := rt.BLength(ctx, tekst)
	if e74 != nil {
		return rt.Value{}, e74
	}
	t75, e76 := rt.Lte(ctx, shirina, t73)
	if e76 != nil {
		return rt.Value{}, e76
	}
	t77, e78 := rt.Cond(ctx, t75)
	if e78 != nil {
		return rt.Value{}, e78
	}
	var t79 rt.Value
	if t77 {
		t79 = tekst
	} else {
		// «длина»
		t80, e81 := rt.BLength(ctx, tekst)
		if e81 != nil {
			return rt.Value{}, e81
		}
		t82, e83 := rt.Sub(ctx, shirina, t80)
		if e83 != nil {
			return rt.Value{}, e83
		}
		t84, e85 := Povtorit(ctx, rt.Text("0"), t82)
		if e85 != nil {
			return rt.Value{}, e85
		}
		t86, e87 := rt.Concat(ctx, t84, tekst)
		if e87 != nil {
			return rt.Value{}, e87
		}
		t79 = t86
	}
	t88 := t79
	// «длина»
	t89, e90 := rt.BLength(ctx, t88)
	if e90 != nil {
		return rt.Value{}, e90
	}
	// «длина»
	t91, e92 := rt.BLength(ctx, tekst)
	if e92 != nil {
		return rt.Value{}, e92
	}
	t93, e94 := rt.Gte(ctx, t89, t91)
	if e94 != nil {
		return rt.Value{}, e94
	}
	// постусловие «короче заказанного не бывает»
	t95, e96 := rt.Post(ctx, t93, "короче заказанного не бывает", "Слева нулями")
	if e96 != nil {
		return rt.Value{}, e96
	}
	if !t95 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «короче заказанного не бывает» функции «Слева нулями»")
	}
	return t88, nil
}

// VerhnyayaPolovina — функция flang «Верхняя половина».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение: число.
func VerhnyayaPolovina(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
	t97, e98 := rt.Mul(ctx, rt.Number(134217729.0), znachenie)
	if e98 != nil {
		return rt.Value{}, e98
	}
	// пусть «сдвинутое»
	sdvinutoe := t97
	t99, e100 := rt.Sub(ctx, sdvinutoe, znachenie)
	if e100 != nil {
		return rt.Value{}, e100
	}
	t101, e102 := rt.Sub(ctx, sdvinutoe, t99)
	if e102 != nil {
		return rt.Value{}, e102
	}
	return t101, nil
}

// PogreshnostProizvedeniya — функция flang «Погрешность произведения».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр pervoe — «первое»: число.
// Параметр vtoroe — «второе»: число.
// Результат — значение: число.
func PogreshnostProizvedeniya(ctx *rt.Ctx, pervoe rt.Value, vtoroe rt.Value) (rt.Value, error) {
	t103, e104 := rt.Mul(ctx, pervoe, vtoroe)
	if e104 != nil {
		return rt.Value{}, e104
	}
	// пусть «произведение»
	proizvedenie := t103
	t105, e106 := VerhnyayaPolovina(ctx, pervoe)
	if e106 != nil {
		return rt.Value{}, e106
	}
	// пусть «верх первого»
	verhPervogo := t105
	t107, e108 := rt.Sub(ctx, pervoe, verhPervogo)
	if e108 != nil {
		return rt.Value{}, e108
	}
	// пусть «низ первого»
	nizPervogo := t107
	t109, e110 := VerhnyayaPolovina(ctx, vtoroe)
	if e110 != nil {
		return rt.Value{}, e110
	}
	// пусть «верх второго»
	verhVtorogo := t109
	t111, e112 := rt.Sub(ctx, vtoroe, verhVtorogo)
	if e112 != nil {
		return rt.Value{}, e112
	}
	// пусть «низ второго»
	nizVtorogo := t111
	t113, e114 := rt.Mul(ctx, verhPervogo, verhVtorogo)
	if e114 != nil {
		return rt.Value{}, e114
	}
	t115, e116 := rt.Sub(ctx, t113, proizvedenie)
	if e116 != nil {
		return rt.Value{}, e116
	}
	t117, e118 := rt.Mul(ctx, verhPervogo, nizVtorogo)
	if e118 != nil {
		return rt.Value{}, e118
	}
	t119, e120 := rt.Add(ctx, t115, t117)
	if e120 != nil {
		return rt.Value{}, e120
	}
	t121, e122 := rt.Mul(ctx, nizPervogo, verhVtorogo)
	if e122 != nil {
		return rt.Value{}, e122
	}
	t123, e124 := rt.Add(ctx, t119, t121)
	if e124 != nil {
		return rt.Value{}, e124
	}
	t125, e126 := rt.Mul(ctx, nizPervogo, nizVtorogo)
	if e126 != nil {
		return rt.Value{}, e126
	}
	t127, e128 := rt.Add(ctx, t123, t125)
	if e128 != nil {
		return rt.Value{}, e128
	}
	return t127, nil
}

// OkruglitKChyotnomu — функция flang «Округлить к чётному».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр masshtab — «масштаб»: число.
// Параметр pogreshnost — «погрешность»: число.
// Результат — значение: число.
func OkruglitKChyotnomu(ctx *rt.Ctx, masshtab rt.Value, pogreshnost rt.Value) (rt.Value, error) {
	t129, e130 := CeloeVniz(ctx, masshtab)
	if e130 != nil {
		return rt.Value{}, e130
	}
	// пусть «низ»
	niz := t129
	t131, e132 := rt.Sub(ctx, masshtab, niz)
	if e132 != nil {
		return rt.Value{}, e132
	}
	// пусть «дробь»
	drob := t131
	t133, e134 := rt.Gt(ctx, drob, rt.Number(0.5))
	if e134 != nil {
		return rt.Value{}, e134
	}
	t135, e136 := rt.Cond(ctx, t133)
	if e136 != nil {
		return rt.Value{}, e136
	}
	var t137 rt.Value
	if t135 {
		t138, e139 := rt.Add(ctx, niz, rt.Number(1.0))
		if e139 != nil {
			return rt.Value{}, e139
		}
		t137 = t138
	} else {
		t140, e141 := rt.Lt(ctx, drob, rt.Number(0.5))
		if e141 != nil {
			return rt.Value{}, e141
		}
		t142, e143 := rt.Cond(ctx, t140)
		if e143 != nil {
			return rt.Value{}, e143
		}
		var t144 rt.Value
		if t142 {
			t144 = niz
		} else {
			t145, e146 := rt.Gt(ctx, pogreshnost, rt.Number(0.0))
			if e146 != nil {
				return rt.Value{}, e146
			}
			t147, e148 := rt.Cond(ctx, t145)
			if e148 != nil {
				return rt.Value{}, e148
			}
			var t149 rt.Value
			if t147 {
				t150, e151 := rt.Add(ctx, niz, rt.Number(1.0))
				if e151 != nil {
					return rt.Value{}, e151
				}
				t149 = t150
			} else {
				t152, e153 := rt.Lt(ctx, pogreshnost, rt.Number(0.0))
				if e153 != nil {
					return rt.Value{}, e153
				}
				t154, e155 := rt.Cond(ctx, t152)
				if e155 != nil {
					return rt.Value{}, e155
				}
				var t156 rt.Value
				if t154 {
					t156 = niz
				} else {
					t157, e158 := rt.Mod(ctx, niz, rt.Number(2.0))
					if e158 != nil {
						return rt.Value{}, e158
					}
					t159, e160 := rt.Cond(ctx, rt.Flag(rt.Equal(t157, rt.Number(0.0))))
					if e160 != nil {
						return rt.Value{}, e160
					}
					var t161 rt.Value
					if t159 {
						t161 = niz
					} else {
						t162, e163 := rt.Add(ctx, niz, rt.Number(1.0))
						if e163 != nil {
							return rt.Value{}, e163
						}
						t161 = t162
					}
					t156 = t161
				}
				t149 = t156
			}
			t144 = t149
		}
		t137 = t144
	}
	t164 := t137
	t165, e166 := rt.Gt(ctx, masshtab, t164)
	if e166 != nil {
		return rt.Value{}, e166
	}
	t167, e168 := rt.Cond(ctx, t165)
	if e168 != nil {
		return rt.Value{}, e168
	}
	var t169 rt.Value
	if t167 {
		t170, e171 := rt.Sub(ctx, masshtab, t164)
		if e171 != nil {
			return rt.Value{}, e171
		}
		t169 = t170
	} else {
		t172, e173 := rt.Sub(ctx, t164, masshtab)
		if e173 != nil {
			return rt.Value{}, e173
		}
		t169 = t172
	}
	t174, e175 := rt.Lte(ctx, t169, rt.Number(1.0))
	if e175 != nil {
		return rt.Value{}, e175
	}
	// постусловие «округление ближе единицы»
	t176, e177 := rt.Post(ctx, t174, "округление ближе единицы", "Округлить к чётному")
	if e177 != nil {
		return rt.Value{}, e177
	}
	if !t176 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «округление ближе единицы» функции «Округлить к чётному»")
	}
	return t164, nil
}

// ZnakMinus — функция flang «Знак минус».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Результат — значение.
func ZnakMinus(ctx *rt.Ctx, znachenie rt.Value) (rt.Value, error) {
	t178, e179 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e179 != nil {
		return rt.Value{}, e179
	}
	t180, e181 := rt.Cond(ctx, t178)
	if e181 != nil {
		return rt.Value{}, e181
	}
	if t180 {
		return rt.Flag(true), nil
	} else {
		t182, e183 := rt.Div(ctx, rt.Number(1.0), znachenie)
		if e183 != nil {
			return rt.Value{}, e183
		}
		t184, e185 := rt.Lt(ctx, t182, rt.Number(0.0))
		if e185 != nil {
			return rt.Value{}, e185
		}
		return t184, nil
	}
}

// Drobyu — функция flang «Дробью».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр znakov — «знаков»: «нат».
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func Drobyu(ctx *rt.Ctx, znachenie rt.Value, znakov rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t186, e187 := ZnakMinus(ctx, znachenie)
	if e187 != nil {
		return rt.Value{}, e187
	}
	t188, e189 := rt.Cond(ctx, t186)
	if e189 != nil {
		return rt.Value{}, e189
	}
	if t188 {
		t190, e191 := rt.Sub(ctx, rt.Number(0.0), znachenie)
		if e191 != nil {
			return rt.Value{}, e191
		}
		t192, e193 := DrobyuNeotricatelnogo(ctx, t190, znakov, razdelitel)
		if e193 != nil {
			return rt.Value{}, e193
		}
		t194, e195 := rt.Concat(ctx, rt.Text("-"), t192)
		if e195 != nil {
			return rt.Value{}, e195
		}
		return t194, nil
	} else {
		return DrobyuNeotricatelnogo(ctx, znachenie, znakov, razdelitel)
	}
}

// DrobyuNeotricatelnogo — функция flang «Дробью неотрицательного».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр znakov — «знаков»: «нат».
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func DrobyuNeotricatelnogo(ctx *rt.Ctx, znachenie rt.Value, znakov rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t196, e197 := DesyatVStepeni(ctx, znakov)
	if e197 != nil {
		return rt.Value{}, e197
	}
	// пусть «степень»
	stepen := t196
	t198, e199 := rt.Mul(ctx, znachenie, stepen)
	if e199 != nil {
		return rt.Value{}, e199
	}
	t200, e201 := PogreshnostProizvedeniya(ctx, znachenie, stepen)
	if e201 != nil {
		return rt.Value{}, e201
	}
	t202, e203 := OkruglitKChyotnomu(ctx, t198, t200)
	if e203 != nil {
		return rt.Value{}, e203
	}
	// пусть «целое»
	celoe := t202
	t204, e205 := DelenieNacelo(ctx, celoe, stepen)
	if e205 != nil {
		return rt.Value{}, e205
	}
	// пусть «целых»
	celyh := t204
	t206, e207 := rt.Mul(ctx, celyh, stepen)
	if e207 != nil {
		return rt.Value{}, e207
	}
	t208, e209 := rt.Sub(ctx, celoe, t206)
	if e209 != nil {
		return rt.Value{}, e209
	}
	// пусть «остаток»
	ostatok := t208
	t210, e211 := rt.Lte(ctx, znakov, rt.Number(0.0))
	if e211 != nil {
		return rt.Value{}, e211
	}
	t212, e213 := rt.Cond(ctx, t210)
	if e213 != nil {
		return rt.Value{}, e213
	}
	if t212 {
		// «к строке»
		t214, e215 := rt.BToString(ctx, celyh)
		if e215 != nil {
			return rt.Value{}, e215
		}
		return t214, nil
	} else {
		// «к строке»
		t216, e217 := rt.BToString(ctx, celyh)
		if e217 != nil {
			return rt.Value{}, e217
		}
		t218, e219 := rt.Concat(ctx, t216, razdelitel)
		if e219 != nil {
			return rt.Value{}, e219
		}
		// «к строке»
		t220, e221 := rt.BToString(ctx, ostatok)
		if e221 != nil {
			return rt.Value{}, e221
		}
		t222, e223 := SlevaNulyami(ctx, t220, znakov)
		if e223 != nil {
			return rt.Value{}, e223
		}
		t224, e225 := rt.Concat(ctx, t218, t222)
		if e225 != nil {
			return rt.Value{}, e225
		}
		return t224, nil
	}
}

// StupenBayt — функция flang «Ступень байт».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр vitki — «витки»: число.
// Параметр zapas — «запас»: «нат».
// Результат — значение: число.
func StupenBayt(ctx *rt.Ctx, vitki rt.Value, zapas rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Ступень байт"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t226, e227 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e227 != nil {
		return rt.Value{}, e227
	}
	t228, e229 := rt.Cond(ctx, t226)
	if e229 != nil {
		return rt.Value{}, e229
	}
	if t228 {
		return rt.Number(0.0), nil
	} else {
		t230, e231 := rt.Gte(ctx, vitki, rt.Number(1024.0))
		if e231 != nil {
			return rt.Value{}, e231
		}
		t232, e233 := rt.Cond(ctx, t230)
		if e233 != nil {
			return rt.Value{}, e233
		}
		if t232 {
			t234, e235 := DelenieNacelo(ctx, vitki, rt.Number(1024.0))
			if e235 != nil {
				return rt.Value{}, e235
			}
			t236, e237 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e237 != nil {
				return rt.Value{}, e237
			}
			t238, e239 := StupenBayt(ctx, t234, t236)
			if e239 != nil {
				return rt.Value{}, e239
			}
			t240, e241 := rt.Add(ctx, rt.Number(1.0), t238)
			if e241 != nil {
				return rt.Value{}, e241
			}
			return t240, nil
		} else {
			return rt.Number(0.0), nil
		}
	}
}

// DelitelBayt — функция flang «Делитель байт».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр vitki — «витки»: число.
// Параметр zapas — «запас»: «нат».
// Результат — значение: число.
func DelitelBayt(ctx *rt.Ctx, vitki rt.Value, zapas rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Делитель байт"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t242, e243 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e243 != nil {
		return rt.Value{}, e243
	}
	t244, e245 := rt.Cond(ctx, t242)
	if e245 != nil {
		return rt.Value{}, e245
	}
	if t244 {
		return rt.Number(1024.0), nil
	} else {
		t246, e247 := rt.Gte(ctx, vitki, rt.Number(1024.0))
		if e247 != nil {
			return rt.Value{}, e247
		}
		t248, e249 := rt.Cond(ctx, t246)
		if e249 != nil {
			return rt.Value{}, e249
		}
		if t248 {
			t250, e251 := DelenieNacelo(ctx, vitki, rt.Number(1024.0))
			if e251 != nil {
				return rt.Value{}, e251
			}
			t252, e253 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e253 != nil {
				return rt.Value{}, e253
			}
			t254, e255 := DelitelBayt(ctx, t250, t252)
			if e255 != nil {
				return rt.Value{}, e255
			}
			t256, e257 := rt.Mul(ctx, rt.Number(1024.0), t254)
			if e257 != nil {
				return rt.Value{}, e257
			}
			return t256, nil
		} else {
			return rt.Number(1024.0), nil
		}
	}
}

// EdinicaBayt — функция flang «Единица байт».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stupen — «ступень»: число.
// Результат — значение: строка.
func EdinicaBayt(ctx *rt.Ctx, stupen rt.Value) (rt.Value, error) {
	t258, e259 := rt.Lte(ctx, stupen, rt.Number(0.0))
	if e259 != nil {
		return rt.Value{}, e259
	}
	t260, e261 := rt.Cond(ctx, t258)
	if e261 != nil {
		return rt.Value{}, e261
	}
	if t260 {
		return rt.Text("КиБ"), nil
	} else {
		t262, e263 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e263 != nil {
			return rt.Value{}, e263
		}
		t264, e265 := rt.Cond(ctx, t262)
		if e265 != nil {
			return rt.Value{}, e265
		}
		if t264 {
			return rt.Text("ПиБ"), nil
		} else {
			t266, e267 := rt.Add(ctx, stupen, rt.Number(1.0))
			if e267 != nil {
				return rt.Value{}, e267
			}
			t268 := make([]rt.Value, 5)
			t268[0] = rt.Text("КиБ")
			t268[1] = rt.Text("МиБ")
			t268[2] = rt.Text("ГиБ")
			t268[3] = rt.Text("ТиБ")
			t268[4] = rt.Text("ПиБ")
			// «элемент»
			t269, e270 := rt.BElement(ctx, t266, rt.List(t268))
			if e270 != nil {
				return rt.Value{}, e270
			}
			return t269, nil
		}
	}
}

// BaytyZnakom — функция flang «Байты знаком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func BaytyZnakom(ctx *rt.Ctx, skolko rt.Value, razdelitel rt.Value) (rt.Value, error) {
	return BaytyYazykom(ctx, skolko, razdelitel, rt.Text("ru"))
}

// Bayty — функция flang «Байты».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Результат — значение: строка.
func Bayty(ctx *rt.Ctx, skolko rt.Value) (rt.Value, error) {
	return BaytyZnakom(ctx, skolko, rt.Text("."))
}

// StupenTysyach — функция flang «Ступень тысяч».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр vitki — «витки»: число.
// Параметр zapas — «запас»: «нат».
// Результат — значение: число.
func StupenTysyach(ctx *rt.Ctx, vitki rt.Value, zapas rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Ступень тысяч"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t271, e272 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e272 != nil {
		return rt.Value{}, e272
	}
	t273, e274 := rt.Cond(ctx, t271)
	if e274 != nil {
		return rt.Value{}, e274
	}
	if t273 {
		return rt.Number(0.0), nil
	} else {
		t275, e276 := rt.Gte(ctx, vitki, rt.Number(1000.0))
		if e276 != nil {
			return rt.Value{}, e276
		}
		t277, e278 := rt.Cond(ctx, t275)
		if e278 != nil {
			return rt.Value{}, e278
		}
		if t277 {
			t279, e280 := DelenieNacelo(ctx, vitki, rt.Number(1000.0))
			if e280 != nil {
				return rt.Value{}, e280
			}
			t281, e282 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e282 != nil {
				return rt.Value{}, e282
			}
			t283, e284 := StupenTysyach(ctx, t279, t281)
			if e284 != nil {
				return rt.Value{}, e284
			}
			t285, e286 := rt.Add(ctx, rt.Number(1.0), t283)
			if e286 != nil {
				return rt.Value{}, e286
			}
			return t285, nil
		} else {
			return rt.Number(0.0), nil
		}
	}
}

// DelitelTysyach — функция flang «Делитель тысяч».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр vitki — «витки»: число.
// Параметр zapas — «запас»: «нат».
// Результат — значение: число.
func DelitelTysyach(ctx *rt.Ctx, vitki rt.Value, zapas rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Делитель тысяч"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	t287, e288 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e288 != nil {
		return rt.Value{}, e288
	}
	t289, e290 := rt.Cond(ctx, t287)
	if e290 != nil {
		return rt.Value{}, e290
	}
	if t289 {
		return rt.Number(1000.0), nil
	} else {
		t291, e292 := rt.Gte(ctx, vitki, rt.Number(1000.0))
		if e292 != nil {
			return rt.Value{}, e292
		}
		t293, e294 := rt.Cond(ctx, t291)
		if e294 != nil {
			return rt.Value{}, e294
		}
		if t293 {
			t295, e296 := DelenieNacelo(ctx, vitki, rt.Number(1000.0))
			if e296 != nil {
				return rt.Value{}, e296
			}
			t297, e298 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e298 != nil {
				return rt.Value{}, e298
			}
			t299, e300 := DelitelTysyach(ctx, t295, t297)
			if e300 != nil {
				return rt.Value{}, e300
			}
			t301, e302 := rt.Mul(ctx, rt.Number(1000.0), t299)
			if e302 != nil {
				return rt.Value{}, e302
			}
			return t301, nil
		} else {
			return rt.Number(1000.0), nil
		}
	}
}

// EdinicaTysyach — функция flang «Единица тысяч».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stupen — «ступень»: число.
// Результат — значение: строка.
func EdinicaTysyach(ctx *rt.Ctx, stupen rt.Value) (rt.Value, error) {
	t303, e304 := rt.Lte(ctx, stupen, rt.Number(0.0))
	if e304 != nil {
		return rt.Value{}, e304
	}
	t305, e306 := rt.Cond(ctx, t303)
	if e306 != nil {
		return rt.Value{}, e306
	}
	if t305 {
		return rt.Text("КБ"), nil
	} else {
		t307, e308 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e308 != nil {
			return rt.Value{}, e308
		}
		t309, e310 := rt.Cond(ctx, t307)
		if e310 != nil {
			return rt.Value{}, e310
		}
		if t309 {
			return rt.Text("ПБ"), nil
		} else {
			t311, e312 := rt.Add(ctx, stupen, rt.Number(1.0))
			if e312 != nil {
				return rt.Value{}, e312
			}
			t313 := make([]rt.Value, 5)
			t313[0] = rt.Text("КБ")
			t313[1] = rt.Text("МБ")
			t313[2] = rt.Text("ГБ")
			t313[3] = rt.Text("ТБ")
			t313[4] = rt.Text("ПБ")
			// «элемент»
			t314, e315 := rt.BElement(ctx, t311, rt.List(t313))
			if e315 != nil {
				return rt.Value{}, e315
			}
			return t314, nil
		}
	}
}

// BaytyDesyatichno — функция flang «Байты десятично».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func BaytyDesyatichno(ctx *rt.Ctx, skolko rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t316, e317 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e317 != nil {
		return rt.Value{}, e317
	}
	t318, e319 := rt.Cond(ctx, t316)
	if e319 != nil {
		return rt.Value{}, e319
	}
	var t320 rt.Value
	if t318 {
		t320 = rt.Text("-")
	} else {
		t320 = rt.Text("")
	}
	// пусть «знак»
	znak := t320
	t321, e322 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e322 != nil {
		return rt.Value{}, e322
	}
	t323, e324 := rt.Cond(ctx, t321)
	if e324 != nil {
		return rt.Value{}, e324
	}
	var t325 rt.Value
	if t323 {
		t326, e327 := rt.Sub(ctx, rt.Number(0.0), skolko)
		if e327 != nil {
			return rt.Value{}, e327
		}
		t325 = t326
	} else {
		t325 = skolko
	}
	t328, e329 := CeloeKNulyu(ctx, t325)
	if e329 != nil {
		return rt.Value{}, e329
	}
	// пусть «величина»
	velichina := t328
	t330, e331 := rt.Lt(ctx, velichina, rt.Number(1000.0))
	if e331 != nil {
		return rt.Value{}, e331
	}
	t332, e333 := rt.Cond(ctx, t330)
	if e333 != nil {
		return rt.Value{}, e333
	}
	if t332 {
		// «к строке»
		t334, e335 := rt.BToString(ctx, velichina)
		if e335 != nil {
			return rt.Value{}, e335
		}
		t336, e337 := rt.Concat(ctx, znak, t334)
		if e337 != nil {
			return rt.Value{}, e337
		}
		t338, e339 := rt.Concat(ctx, t336, rt.Text(" Б"))
		if e339 != nil {
			return rt.Value{}, e339
		}
		return t338, nil
	} else {
		t340, e341 := DelenieNacelo(ctx, velichina, rt.Number(1000.0))
		if e341 != nil {
			return rt.Value{}, e341
		}
		// пусть «витки»
		vitki := t340
		t342, e343 := DelitelTysyach(ctx, vitki, rt.Number(4.0))
		if e343 != nil {
			return rt.Value{}, e343
		}
		t344, e345 := rt.Div(ctx, velichina, t342)
		if e345 != nil {
			return rt.Value{}, e345
		}
		t346, e347 := Drobyu(ctx, t344, rt.Number(1.0), razdelitel)
		if e347 != nil {
			return rt.Value{}, e347
		}
		t348, e349 := rt.Concat(ctx, znak, t346)
		if e349 != nil {
			return rt.Value{}, e349
		}
		t350, e351 := rt.Concat(ctx, t348, rt.Text(" "))
		if e351 != nil {
			return rt.Value{}, e351
		}
		t352, e353 := StupenTysyach(ctx, vitki, rt.Number(4.0))
		if e353 != nil {
			return rt.Value{}, e353
		}
		t354, e355 := EdinicaTysyach(ctx, t352)
		if e355 != nil {
			return rt.Value{}, e355
		}
		t356, e357 := rt.Concat(ctx, t350, t354)
		if e357 != nil {
			return rt.Value{}, e357
		}
		return t356, nil
	}
}

// DolyaOt — функция flang «Доля от».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр chast — «часть»: число.
// Параметр celoe — «целое»: число.
// Результат — значение: число.
func DolyaOt(ctx *rt.Ctx, chast rt.Value, celoe rt.Value) (rt.Value, error) {
	t358, e359 := rt.Cond(ctx, rt.Flag(rt.Equal(celoe, rt.Number(0.0))))
	if e359 != nil {
		return rt.Value{}, e359
	}
	if t358 {
		return rt.Number(0.0), nil
	} else {
		t360, e361 := rt.Div(ctx, chast, celoe)
		if e361 != nil {
			return rt.Value{}, e361
		}
		return t360, nil
	}
}

// ProcentCelogo — функция flang «Процент целого».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр chast — «часть»: число.
// Параметр celoe — «целое»: число.
// Результат — значение: число.
func ProcentCelogo(ctx *rt.Ctx, chast rt.Value, celoe rt.Value) (rt.Value, error) {
	t362, e363 := rt.Cond(ctx, rt.Flag(rt.Equal(celoe, rt.Number(0.0))))
	if e363 != nil {
		return rt.Value{}, e363
	}
	if t362 {
		return rt.Number(0.0), nil
	} else {
		t364, e365 := rt.Mul(ctx, rt.Number(100.0), chast)
		if e365 != nil {
			return rt.Value{}, e365
		}
		t366, e367 := rt.Div(ctx, t364, celoe)
		if e367 != nil {
			return rt.Value{}, e367
		}
		return t366, nil
	}
}

// ProcentyZnakom — функция flang «Проценты знаком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func ProcentyZnakom(ctx *rt.Ctx, dolya rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t368, e369 := rt.Mul(ctx, dolya, rt.Number(100.0))
	if e369 != nil {
		return rt.Value{}, e369
	}
	t370, e371 := Drobyu(ctx, t368, rt.Number(1.0), razdelitel)
	if e371 != nil {
		return rt.Value{}, e371
	}
	t372, e373 := rt.Concat(ctx, t370, rt.Text("%"))
	if e373 != nil {
		return rt.Value{}, e373
	}
	return t372, nil
}

// Procenty — функция flang «Проценты».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр dolya — «доля»: число.
// Результат — значение: строка.
func Procenty(ctx *rt.Ctx, dolya rt.Value) (rt.Value, error) {
	return ProcentyZnakom(ctx, dolya, rt.Text("."))
}

// Dlitelnost — функция flang «Длительность».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sekundy — «секунды»: число.
// Результат — значение: строка.
func Dlitelnost(ctx *rt.Ctx, sekundy rt.Value) (rt.Value, error) {
	return DlitelnostYazykom(ctx, sekundy, rt.Text("ru"))
}

// Proshlo — функция flang «Прошло».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sekundy — «секунды»: число.
// Результат — значение: строка.
func Proshlo(ctx *rt.Ctx, sekundy rt.Value) (rt.Value, error) {
	t374, e375 := rt.Lt(ctx, sekundy, rt.Number(3600.0))
	if e375 != nil {
		return rt.Value{}, e375
	}
	t376, e377 := rt.Cond(ctx, t374)
	if e377 != nil {
		return rt.Value{}, e377
	}
	if t376 {
		t378, e379 := DelenieNacelo(ctx, sekundy, rt.Number(60.0))
		if e379 != nil {
			return rt.Value{}, e379
		}
		// «к строке»
		t380, e381 := rt.BToString(ctx, t378)
		if e381 != nil {
			return rt.Value{}, e381
		}
		t382, e383 := rt.Concat(ctx, t380, rt.Text(" "))
		if e383 != nil {
			return rt.Value{}, e383
		}
		t384, e385 := rt.Concat(ctx, t382, rt.Text("мин"))
		if e385 != nil {
			return rt.Value{}, e385
		}
		return t384, nil
	} else {
		t386, e387 := rt.Lt(ctx, sekundy, rt.Number(172800.0))
		if e387 != nil {
			return rt.Value{}, e387
		}
		t388, e389 := rt.Cond(ctx, t386)
		if e389 != nil {
			return rt.Value{}, e389
		}
		if t388 {
			t390, e391 := DelenieNacelo(ctx, sekundy, rt.Number(3600.0))
			if e391 != nil {
				return rt.Value{}, e391
			}
			// «к строке»
			t392, e393 := rt.BToString(ctx, t390)
			if e393 != nil {
				return rt.Value{}, e393
			}
			t394, e395 := rt.Concat(ctx, t392, rt.Text(" "))
			if e395 != nil {
				return rt.Value{}, e395
			}
			t396, e397 := rt.Concat(ctx, t394, rt.Text("ч"))
			if e397 != nil {
				return rt.Value{}, e397
			}
			return t396, nil
		} else {
			t398, e399 := rt.Div(ctx, sekundy, rt.Number(3600.0))
			if e399 != nil {
				return rt.Value{}, e399
			}
			t400, e401 := rt.Div(ctx, t398, rt.Number(24.0))
			if e401 != nil {
				return rt.Value{}, e401
			}
			t402, e403 := CeloeKNulyu(ctx, t400)
			if e403 != nil {
				return rt.Value{}, e403
			}
			// «к строке»
			t404, e405 := rt.BToString(ctx, t402)
			if e405 != nil {
				return rt.Value{}, e405
			}
			t406, e407 := rt.Concat(ctx, t404, rt.Text(" "))
			if e407 != nil {
				return rt.Value{}, e407
			}
			t408, e409 := rt.Concat(ctx, t406, rt.Text("дн"))
			if e409 != nil {
				return rt.Value{}, e409
			}
			return t408, nil
		}
	}
}

// ZanyaloZnakom — функция flang «Заняло знаком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр millisekundy — «миллисекунды»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func ZanyaloZnakom(ctx *rt.Ctx, millisekundy rt.Value, razdelitel rt.Value) (rt.Value, error) {
	return ZanyaloYazykom(ctx, millisekundy, razdelitel, rt.Text("ru"))
}

// Zanyalo — функция flang «Заняло».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр millisekundy — «миллисекунды»: число.
// Результат — значение: строка.
func Zanyalo(ctx *rt.Ctx, millisekundy rt.Value) (rt.Value, error) {
	return ZanyaloZnakom(ctx, millisekundy, rt.Text("."))
}

// ShagRazryada — функция flang «Шаг разряда».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр hod — «ход»: «Ход разрядов».
// Параметр cifra — «цифра»: строка.
// Результат — значение: «Ход разрядов».
func ShagRazryada(ctx *rt.Ctx, hod rt.Value, cifra rt.Value) (rt.Value, error) {
	t410, e411 := rt.FieldGet(ctx, hod, "осталось")
	if e411 != nil {
		return rt.Value{}, e411
	}
	t412, e413 := rt.Mod(ctx, t410, rt.Number(3.0))
	if e413 != nil {
		return rt.Value{}, e413
	}
	t414, e415 := rt.Cond(ctx, rt.Flag(rt.Equal(t412, rt.Number(0.0))))
	if e415 != nil {
		return rt.Value{}, e415
	}
	var t416 rt.Value
	if t414 {
		t417, e418 := rt.FieldGet(ctx, hod, "собрано")
		if e418 != nil {
			return rt.Value{}, e418
		}
		// «длина»
		t419, e420 := rt.BLength(ctx, t417)
		if e420 != nil {
			return rt.Value{}, e420
		}
		t421, e422 := rt.Gt(ctx, t419, rt.Number(0.0))
		if e422 != nil {
			return rt.Value{}, e422
		}
		t416 = t421
	} else {
		t416 = rt.Flag(false)
	}
	t423, e424 := rt.Cond(ctx, t416)
	if e424 != nil {
		return rt.Value{}, e424
	}
	var t425 rt.Value
	if t423 {
		t426, e427 := rt.FieldGet(ctx, hod, "собрано")
		if e427 != nil {
			return rt.Value{}, e427
		}
		t428, e429 := rt.FieldGet(ctx, hod, "разделитель")
		if e429 != nil {
			return rt.Value{}, e429
		}
		t430, e431 := rt.Concat(ctx, t426, t428)
		if e431 != nil {
			return rt.Value{}, e431
		}
		t432, e433 := rt.Concat(ctx, t430, cifra)
		if e433 != nil {
			return rt.Value{}, e433
		}
		t434, e435 := rt.FieldGet(ctx, hod, "осталось")
		if e435 != nil {
			return rt.Value{}, e435
		}
		t436, e437 := rt.Sub(ctx, t434, rt.Number(1.0))
		if e437 != nil {
			return rt.Value{}, e437
		}
		t438, e439 := rt.FieldGet(ctx, hod, "разделитель")
		if e439 != nil {
			return rt.Value{}, e439
		}
		t440 := make([]rt.Field, 3)
		t440[0] = rt.Field{Name: "собрано", Value: t432}
		t440[1] = rt.Field{Name: "осталось", Value: t436}
		t440[2] = rt.Field{Name: "разделитель", Value: t438}
		t425 = rt.Record(t440)
	} else {
		t441, e442 := rt.FieldGet(ctx, hod, "собрано")
		if e442 != nil {
			return rt.Value{}, e442
		}
		t443, e444 := rt.Concat(ctx, t441, cifra)
		if e444 != nil {
			return rt.Value{}, e444
		}
		t445, e446 := rt.FieldGet(ctx, hod, "осталось")
		if e446 != nil {
			return rt.Value{}, e446
		}
		t447, e448 := rt.Sub(ctx, t445, rt.Number(1.0))
		if e448 != nil {
			return rt.Value{}, e448
		}
		t449, e450 := rt.FieldGet(ctx, hod, "разделитель")
		if e450 != nil {
			return rt.Value{}, e450
		}
		t451 := make([]rt.Field, 3)
		t451[0] = rt.Field{Name: "собрано", Value: t443}
		t451[1] = rt.Field{Name: "осталось", Value: t447}
		t451[2] = rt.Field{Name: "разделитель", Value: t449}
		t425 = rt.Record(t451)
	}
	t452 := t425
	t453, e454 := rt.FieldGet(ctx, t452, "осталось")
	if e454 != nil {
		return rt.Value{}, e454
	}
	t455, e456 := rt.FieldGet(ctx, hod, "осталось")
	if e456 != nil {
		return rt.Value{}, e456
	}
	t457, e458 := rt.Lt(ctx, t453, t455)
	if e458 != nil {
		return rt.Value{}, e458
	}
	// постусловие «осталось убывает»
	t459, e460 := rt.Post(ctx, t457, "осталось убывает", "Шаг разряда")
	if e460 != nil {
		return rt.Value{}, e460
	}
	if !t459 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «осталось убывает» функции «Шаг разряда»")
	}
	t461, e462 := rt.FieldGet(ctx, t452, "разделитель")
	if e462 != nil {
		return rt.Value{}, e462
	}
	t463, e464 := rt.FieldGet(ctx, hod, "разделитель")
	if e464 != nil {
		return rt.Value{}, e464
	}
	// постусловие «разделитель не меняется»
	t465, e466 := rt.Post(ctx, rt.Flag(rt.Equal(t461, t463)), "разделитель не меняется", "Шаг разряда")
	if e466 != nil {
		return rt.Value{}, e466
	}
	if !t465 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «разделитель не меняется» функции «Шаг разряда»")
	}
	return t452, nil
}

// Razryadami — функция flang «Разрядами».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func Razryadami(ctx *rt.Ctx, znachenie rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t467, e468 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e468 != nil {
		return rt.Value{}, e468
	}
	t469, e470 := rt.Cond(ctx, t467)
	if e470 != nil {
		return rt.Value{}, e470
	}
	var t471 rt.Value
	if t469 {
		t471 = rt.Text("-")
	} else {
		t471 = rt.Text("")
	}
	// пусть «знак»
	znak := t471
	t472, e473 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e473 != nil {
		return rt.Value{}, e473
	}
	t474, e475 := rt.Cond(ctx, t472)
	if e475 != nil {
		return rt.Value{}, e475
	}
	var t476 rt.Value
	if t474 {
		t477, e478 := rt.Sub(ctx, rt.Number(0.0), znachenie)
		if e478 != nil {
			return rt.Value{}, e478
		}
		t476 = t477
	} else {
		t476 = znachenie
	}
	t479, e480 := CeloeKNulyu(ctx, t476)
	if e480 != nil {
		return rt.Value{}, e480
	}
	// «к строке»
	t481, e482 := rt.BToString(ctx, t479)
	if e482 != nil {
		return rt.Value{}, e482
	}
	// пусть «цифры»
	cifry := t481
	// «символы»
	t483, e484 := rt.BCharacters(ctx, cifry)
	if e484 != nil {
		return rt.Value{}, e484
	}
	t485, e486 := rt.RequireList(ctx, t483, "свёртка")
	if e486 != nil {
		return rt.Value{}, e486
	}
	// «длина»
	t487, e488 := rt.BLength(ctx, cifry)
	if e488 != nil {
		return rt.Value{}, e488
	}
	t489 := make([]rt.Field, 3)
	t489[0] = rt.Field{Name: "собрано", Value: rt.Text("")}
	t489[1] = rt.Field{Name: "осталось", Value: t487}
	t489[2] = rt.Field{Name: "разделитель", Value: razdelitel}
	// «ход»
	hod := rt.Record(t489)
	for t490 := range t485 {
		// «цифра»
		cifra := t485[t490]
		t491, e492 := ShagRazryada(ctx, hod, cifra)
		if e492 != nil {
			return rt.Value{}, e492
		}
		hod = t491
	}
	t493, e494 := rt.FieldGet(ctx, hod, "собрано")
	if e494 != nil {
		return rt.Value{}, e494
	}
	t495, e496 := rt.Concat(ctx, znak, t493)
	if e496 != nil {
		return rt.Value{}, e496
	}
	return t495, nil
}

// Probelnyy — функция flang «Пробельный».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znak — «знак»: строка.
// Результат — значение.
func Probelnyy(ctx *rt.Ctx, znak rt.Value) (rt.Value, error) {
	if rt.ChainEmpty(znak) {
		return rt.Flag(false), nil
	} else if rt.ChainCons(znak) {
		// голова «голова»
		golova := rt.ChainHead(znak)
		// хвост «хвост»
		hvost := rt.ChainTail(znak)
		_ = hvost
		// «код символа»
		t497, e498 := rt.BCharCodeProven(ctx, golova)
		if e498 != nil {
			return rt.Value{}, e498
		}
		t499, e500 := rt.Cond(ctx, rt.Flag(rt.Equal(t497, rt.Number(32.0))))
		if e500 != nil {
			return rt.Value{}, e500
		}
		var t501 rt.Value
		if t499 {
			t501 = rt.Flag(true)
		} else {
			// «код символа»
			t502, e503 := rt.BCharCodeProven(ctx, golova)
			if e503 != nil {
				return rt.Value{}, e503
			}
			t501 = rt.Flag(rt.Equal(t502, rt.Number(9.0)))
		}
		t504, e505 := rt.Cond(ctx, t501)
		if e505 != nil {
			return rt.Value{}, e505
		}
		var t506 rt.Value
		if t504 {
			t506 = rt.Flag(true)
		} else {
			// «код символа»
			t507, e508 := rt.BCharCodeProven(ctx, golova)
			if e508 != nil {
				return rt.Value{}, e508
			}
			t506 = rt.Flag(rt.Equal(t507, rt.Number(10.0)))
		}
		t509, e510 := rt.Cond(ctx, t506)
		if e510 != nil {
			return rt.Value{}, e510
		}
		var t511 rt.Value
		if t509 {
			t511 = rt.Flag(true)
		} else {
			// «код символа»
			t512, e513 := rt.BCharCodeProven(ctx, golova)
			if e513 != nil {
				return rt.Value{}, e513
			}
			t511 = rt.Flag(rt.Equal(t512, rt.Number(11.0)))
		}
		t514, e515 := rt.Cond(ctx, t511)
		if e515 != nil {
			return rt.Value{}, e515
		}
		var t516 rt.Value
		if t514 {
			t516 = rt.Flag(true)
		} else {
			// «код символа»
			t517, e518 := rt.BCharCodeProven(ctx, golova)
			if e518 != nil {
				return rt.Value{}, e518
			}
			t516 = rt.Flag(rt.Equal(t517, rt.Number(12.0)))
		}
		t519, e520 := rt.Cond(ctx, t516)
		if e520 != nil {
			return rt.Value{}, e520
		}
		var t521 rt.Value
		if t519 {
			t521 = rt.Flag(true)
		} else {
			// «код символа»
			t522, e523 := rt.BCharCodeProven(ctx, golova)
			if e523 != nil {
				return rt.Value{}, e523
			}
			t521 = rt.Flag(rt.Equal(t522, rt.Number(13.0)))
		}
		t524, e525 := rt.Cond(ctx, t521)
		if e525 != nil {
			return rt.Value{}, e525
		}
		var t526 rt.Value
		if t524 {
			t526 = rt.Flag(true)
		} else {
			// «код символа»
			t527, e528 := rt.BCharCodeProven(ctx, golova)
			if e528 != nil {
				return rt.Value{}, e528
			}
			t526 = rt.Flag(rt.Equal(t527, rt.Number(133.0)))
		}
		t529, e530 := rt.Cond(ctx, t526)
		if e530 != nil {
			return rt.Value{}, e530
		}
		if t529 {
			return rt.Flag(true), nil
		} else {
			// «код символа»
			t531, e532 := rt.BCharCodeProven(ctx, golova)
			if e532 != nil {
				return rt.Value{}, e532
			}
			return rt.Flag(rt.Equal(t531, rt.Number(160.0))), nil
		}
	} else {
		return rt.Value{}, rt.MatchFail(ctx, znak)
	}
}

// Procherk — функция flang «Прочерк».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Результат — значение: строка.
func Procherk(ctx *rt.Ctx, tekst rt.Value) (rt.Value, error) {
	// «символы»
	t533, e534 := rt.BCharacters(ctx, tekst)
	if e534 != nil {
		return rt.Value{}, e534
	}
	t535, e536 := rt.RequireList(ctx, t533, "отфильтровать")
	if e536 != nil {
		return rt.Value{}, e536
	}
	t537 := make([]rt.Value, 0, len(t535))
	for t538 := range t535 {
		// «знак»
		znak := t535[t538]
		t539, e540 := Probelnyy(ctx, znak)
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
		t544, e545 := rt.Keep(ctx, t543)
		if e545 != nil {
			return rt.Value{}, e545
		}
		if t544 {
			t537 = append(t537, znak)
		}
	}
	// «длина»
	t546, e547 := rt.BLength(ctx, rt.List(t537))
	if e547 != nil {
		return rt.Value{}, e547
	}
	t548, e549 := rt.Cond(ctx, rt.Flag(rt.Equal(t546, rt.Number(0.0))))
	if e549 != nil {
		return rt.Value{}, e549
	}
	if t548 {
		return rt.Text("—"), nil
	} else {
		return tekst, nil
	}
}

// Obrezat — функция flang «Обрезать».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр skolko — «сколько»: число.
// Результат — значение: строка.
func Obrezat(ctx *rt.Ctx, tekst rt.Value, skolko rt.Value) (rt.Value, error) {
	// «длина»
	t550, e551 := rt.BLength(ctx, tekst)
	if e551 != nil {
		return rt.Value{}, e551
	}
	t552, e553 := rt.Lte(ctx, t550, skolko)
	if e553 != nil {
		return rt.Value{}, e553
	}
	t554, e555 := rt.Cond(ctx, t552)
	if e555 != nil {
		return rt.Value{}, e555
	}
	var t556 rt.Value
	if t554 {
		t556 = tekst
	} else {
		t557, e558 := rt.Lte(ctx, skolko, rt.Number(0.0))
		if e558 != nil {
			return rt.Value{}, e558
		}
		t559, e560 := rt.Cond(ctx, t557)
		if e560 != nil {
			return rt.Value{}, e560
		}
		var t561 rt.Value
		if t559 {
			t561 = rt.Text("")
		} else {
			t562, e563 := rt.Lte(ctx, skolko, rt.Number(1.0))
			if e563 != nil {
				return rt.Value{}, e563
			}
			t564, e565 := rt.Cond(ctx, t562)
			if e565 != nil {
				return rt.Value{}, e565
			}
			var t566 rt.Value
			if t564 {
				// «подстрока»
				t567, e568 := rt.BSubstring(ctx, tekst, rt.Number(1.0), rt.Number(1.0))
				if e568 != nil {
					return rt.Value{}, e568
				}
				t566 = t567
			} else {
				t569, e570 := rt.Sub(ctx, skolko, rt.Number(1.0))
				if e570 != nil {
					return rt.Value{}, e570
				}
				// «подстрока»
				t571, e572 := rt.BSubstring(ctx, tekst, rt.Number(1.0), t569)
				if e572 != nil {
					return rt.Value{}, e572
				}
				t573, e574 := rt.Concat(ctx, t571, rt.Text("…"))
				if e574 != nil {
					return rt.Value{}, e574
				}
				t566 = t573
			}
			t561 = t566
		}
		t556 = t561
	}
	t575 := t556
	t576, e577 := rt.Gte(ctx, skolko, rt.Number(0.0))
	if e577 != nil {
		return rt.Value{}, e577
	}
	t578, e579 := rt.Cond(ctx, t576)
	if e579 != nil {
		return rt.Value{}, e579
	}
	var t580 rt.Value
	if t578 {
		// «длина»
		t581, e582 := rt.BLength(ctx, t575)
		if e582 != nil {
			return rt.Value{}, e582
		}
		// «длина»
		t583, e584 := rt.BLength(ctx, tekst)
		if e584 != nil {
			return rt.Value{}, e584
		}
		t585, e586 := rt.Lte(ctx, t583, skolko)
		if e586 != nil {
			return rt.Value{}, e586
		}
		t587, e588 := rt.Cond(ctx, t585)
		if e588 != nil {
			return rt.Value{}, e588
		}
		var t589 rt.Value
		if t587 {
			// «длина»
			t590, e591 := rt.BLength(ctx, tekst)
			if e591 != nil {
				return rt.Value{}, e591
			}
			t589 = t590
		} else {
			t589 = skolko
		}
		t592, e593 := rt.Lte(ctx, t581, t589)
		if e593 != nil {
			return rt.Value{}, e593
		}
		t580 = t592
	} else {
		t580 = rt.Flag(true)
	}
	// постусловие «обрезанное не длиннее заказа»
	t594, e595 := rt.Post(ctx, t580, "обрезанное не длиннее заказа", "Обрезать")
	if e595 != nil {
		return rt.Value{}, e595
	}
	if !t594 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «обрезанное не длиннее заказа» функции «Обрезать»")
	}
	return t575, nil
}

// Umestit — функция flang «Уместить».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр shirina — «ширина»: число.
// Результат — значение: строка.
func Umestit(ctx *rt.Ctx, tekst rt.Value, shirina rt.Value) (rt.Value, error) {
	t596, e597 := rt.Lte(ctx, shirina, rt.Number(0.0))
	if e597 != nil {
		return rt.Value{}, e597
	}
	t598, e599 := rt.Cond(ctx, t596)
	if e599 != nil {
		return rt.Value{}, e599
	}
	var t600 rt.Value
	if t598 {
		t600 = rt.Text("")
	} else {
		// «длина»
		t601, e602 := rt.BLength(ctx, tekst)
		if e602 != nil {
			return rt.Value{}, e602
		}
		t603, e604 := rt.Cond(ctx, rt.Flag(rt.Equal(t601, shirina)))
		if e604 != nil {
			return rt.Value{}, e604
		}
		var t605 rt.Value
		if t603 {
			t605 = tekst
		} else {
			// «длина»
			t606, e607 := rt.BLength(ctx, tekst)
			if e607 != nil {
				return rt.Value{}, e607
			}
			t608, e609 := rt.Lt(ctx, t606, shirina)
			if e609 != nil {
				return rt.Value{}, e609
			}
			t610, e611 := rt.Cond(ctx, t608)
			if e611 != nil {
				return rt.Value{}, e611
			}
			var t612 rt.Value
			if t610 {
				// «длина»
				t613, e614 := rt.BLength(ctx, tekst)
				if e614 != nil {
					return rt.Value{}, e614
				}
				t615, e616 := rt.Sub(ctx, shirina, t613)
				if e616 != nil {
					return rt.Value{}, e616
				}
				t617, e618 := Povtorit(ctx, rt.Text(" "), t615)
				if e618 != nil {
					return rt.Value{}, e618
				}
				t619, e620 := rt.Concat(ctx, tekst, t617)
				if e620 != nil {
					return rt.Value{}, e620
				}
				t612 = t619
			} else {
				t621, e622 := rt.Cond(ctx, rt.Flag(rt.Equal(shirina, rt.Number(1.0))))
				if e622 != nil {
					return rt.Value{}, e622
				}
				var t623 rt.Value
				if t621 {
					t623 = rt.Text("…")
				} else {
					t624, e625 := rt.Sub(ctx, shirina, rt.Number(1.0))
					if e625 != nil {
						return rt.Value{}, e625
					}
					// «подстрока»
					t626, e627 := rt.BSubstring(ctx, tekst, rt.Number(1.0), t624)
					if e627 != nil {
						return rt.Value{}, e627
					}
					t628, e629 := rt.Concat(ctx, t626, rt.Text("…"))
					if e629 != nil {
						return rt.Value{}, e629
					}
					t623 = t628
				}
				t612 = t623
			}
			t605 = t612
		}
		t600 = t605
	}
	t630 := t600
	t631, e632 := rt.Gt(ctx, shirina, rt.Number(0.0))
	if e632 != nil {
		return rt.Value{}, e632
	}
	t633, e634 := rt.Cond(ctx, t631)
	if e634 != nil {
		return rt.Value{}, e634
	}
	var t635 rt.Value
	if t633 {
		// «длина»
		t636, e637 := rt.BLength(ctx, t630)
		if e637 != nil {
			return rt.Value{}, e637
		}
		t635 = rt.Flag(rt.Equal(t636, shirina))
	} else {
		// «длина»
		t638, e639 := rt.BLength(ctx, t630)
		if e639 != nil {
			return rt.Value{}, e639
		}
		t635 = rt.Flag(rt.Equal(t638, rt.Number(0.0)))
	}
	// постусловие «ширина ровно заказанная»
	t640, e641 := rt.Post(ctx, t635, "ширина ровно заказанная", "Уместить")
	if e641 != nil {
		return rt.Value{}, e641
	}
	if !t640 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина ровно заказанная» функции «Уместить»")
	}
	return t630, nil
}

// Sprava — функция flang «Справа».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр shirina — «ширина»: число.
// Результат — значение: строка.
func Sprava(ctx *rt.Ctx, tekst rt.Value, shirina rt.Value) (rt.Value, error) {
	// «длина»
	t642, e643 := rt.BLength(ctx, tekst)
	if e643 != nil {
		return rt.Value{}, e643
	}
	t644, e645 := rt.Sub(ctx, shirina, t642)
	if e645 != nil {
		return rt.Value{}, e645
	}
	t646, e647 := rt.Gt(ctx, t644, rt.Number(0.0))
	if e647 != nil {
		return rt.Value{}, e647
	}
	t648, e649 := rt.Cond(ctx, t646)
	if e649 != nil {
		return rt.Value{}, e649
	}
	var t650 rt.Value
	if t648 {
		// «длина»
		t651, e652 := rt.BLength(ctx, tekst)
		if e652 != nil {
			return rt.Value{}, e652
		}
		t653, e654 := rt.Sub(ctx, shirina, t651)
		if e654 != nil {
			return rt.Value{}, e654
		}
		t655, e656 := Povtorit(ctx, rt.Text(" "), t653)
		if e656 != nil {
			return rt.Value{}, e656
		}
		t657, e658 := rt.Concat(ctx, t655, tekst)
		if e658 != nil {
			return rt.Value{}, e658
		}
		t650 = t657
	} else {
		t659, e660 := Umestit(ctx, tekst, shirina)
		if e660 != nil {
			return rt.Value{}, e660
		}
		t650 = t659
	}
	t661 := t650
	t662, e663 := rt.Gt(ctx, shirina, rt.Number(0.0))
	if e663 != nil {
		return rt.Value{}, e663
	}
	t664, e665 := rt.Cond(ctx, t662)
	if e665 != nil {
		return rt.Value{}, e665
	}
	var t666 rt.Value
	if t664 {
		// «длина»
		t667, e668 := rt.BLength(ctx, t661)
		if e668 != nil {
			return rt.Value{}, e668
		}
		t666 = rt.Flag(rt.Equal(t667, shirina))
	} else {
		// «длина»
		t669, e670 := rt.BLength(ctx, t661)
		if e670 != nil {
			return rt.Value{}, e670
		}
		t666 = rt.Flag(rt.Equal(t669, rt.Number(0.0)))
	}
	// постусловие «ширина ровно заказанная»
	t671, e672 := rt.Post(ctx, t666, "ширина ровно заказанная", "Справа")
	if e672 != nil {
		return rt.Value{}, e672
	}
	if !t671 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «ширина ровно заказанная» функции «Справа»")
	}
	return t661, nil
}

// Angliyskiy — функция flang «Английский».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр yazyk — «язык»: строка.
// Результат — значение.
func Angliyskiy(ctx *rt.Ctx, yazyk rt.Value) (rt.Value, error) {
	return rt.Flag(rt.Equal(yazyk, rt.Text("en"))), nil
}

// DrobyuYazykom — функция flang «Дробью языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр znakov — «знаков»: «нат».
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func DrobyuYazykom(ctx *rt.Ctx, znachenie rt.Value, znakov rt.Value, yazyk rt.Value) (rt.Value, error) {
	t673, e674 := Angliyskiy(ctx, yazyk)
	if e674 != nil {
		return rt.Value{}, e674
	}
	t675, e676 := rt.Cond(ctx, t673)
	if e676 != nil {
		return rt.Value{}, e676
	}
	var t677 rt.Value
	if t675 {
		t677 = rt.Text(".")
	} else {
		t677 = rt.Text(",")
	}
	return Drobyu(ctx, znachenie, znakov, t677)
}

// Procentom — функция flang «Процентом».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр znakov — «знаков»: «нат».
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func Procentom(ctx *rt.Ctx, znachenie rt.Value, znakov rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t678, e679 := Drobyu(ctx, znachenie, znakov, razdelitel)
	if e679 != nil {
		return rt.Value{}, e679
	}
	t680, e681 := rt.Concat(ctx, t678, rt.Text("%"))
	if e681 != nil {
		return rt.Value{}, e681
	}
	return t680, nil
}

// ProcentomYazykom — функция flang «Процентом языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр znakov — «знаков»: «нат».
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func ProcentomYazykom(ctx *rt.Ctx, znachenie rt.Value, znakov rt.Value, yazyk rt.Value) (rt.Value, error) {
	t682, e683 := Angliyskiy(ctx, yazyk)
	if e683 != nil {
		return rt.Value{}, e683
	}
	t684, e685 := rt.Cond(ctx, t682)
	if e685 != nil {
		return rt.Value{}, e685
	}
	var t686 rt.Value
	if t684 {
		t686 = rt.Text(".")
	} else {
		t686 = rt.Text(",")
	}
	return Procentom(ctx, znachenie, znakov, t686)
}

// RazryadyYazykom — функция flang «Разряды языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func RazryadyYazykom(ctx *rt.Ctx, znachenie rt.Value, yazyk rt.Value) (rt.Value, error) {
	t687, e688 := Angliyskiy(ctx, yazyk)
	if e688 != nil {
		return rt.Value{}, e688
	}
	t689, e690 := rt.Cond(ctx, t687)
	if e690 != nil {
		return rt.Value{}, e690
	}
	var t691 rt.Value
	if t689 {
		t691 = rt.Text(",")
	} else {
		t691 = rt.Text(" ")
	}
	return Razryadami(ctx, znachenie, t691)
}

// EdinicaBaytYazykom — функция flang «Единица байт языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stupen — «ступень»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func EdinicaBaytYazykom(ctx *rt.Ctx, stupen rt.Value, yazyk rt.Value) (rt.Value, error) {
	t692, e693 := rt.Lte(ctx, stupen, rt.Number(-1.0))
	if e693 != nil {
		return rt.Value{}, e693
	}
	t694, e695 := rt.Cond(ctx, t692)
	if e695 != nil {
		return rt.Value{}, e695
	}
	if t694 {
		t696, e697 := Angliyskiy(ctx, yazyk)
		if e697 != nil {
			return rt.Value{}, e697
		}
		t698, e699 := rt.Cond(ctx, t696)
		if e699 != nil {
			return rt.Value{}, e699
		}
		if t698 {
			return rt.Text("B"), nil
		} else {
			return rt.Text("Б"), nil
		}
	} else {
		t700, e701 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e701 != nil {
			return rt.Value{}, e701
		}
		t702, e703 := rt.Cond(ctx, t700)
		if e703 != nil {
			return rt.Value{}, e703
		}
		if t702 {
			t704, e705 := Angliyskiy(ctx, yazyk)
			if e705 != nil {
				return rt.Value{}, e705
			}
			t706, e707 := rt.Cond(ctx, t704)
			if e707 != nil {
				return rt.Value{}, e707
			}
			if t706 {
				return rt.Text("PiB"), nil
			} else {
				return rt.Text("ПиБ"), nil
			}
		} else {
			t708, e709 := Angliyskiy(ctx, yazyk)
			if e709 != nil {
				return rt.Value{}, e709
			}
			t710, e711 := rt.Cond(ctx, t708)
			if e711 != nil {
				return rt.Value{}, e711
			}
			if t710 {
				t712, e713 := rt.Add(ctx, stupen, rt.Number(1.0))
				if e713 != nil {
					return rt.Value{}, e713
				}
				t714 := make([]rt.Value, 5)
				t714[0] = rt.Text("KiB")
				t714[1] = rt.Text("MiB")
				t714[2] = rt.Text("GiB")
				t714[3] = rt.Text("TiB")
				t714[4] = rt.Text("PiB")
				// «элемент»
				t715, e716 := rt.BElement(ctx, t712, rt.List(t714))
				if e716 != nil {
					return rt.Value{}, e716
				}
				return t715, nil
			} else {
				t717, e718 := rt.Add(ctx, stupen, rt.Number(1.0))
				if e718 != nil {
					return rt.Value{}, e718
				}
				t719 := make([]rt.Value, 5)
				t719[0] = rt.Text("КиБ")
				t719[1] = rt.Text("МиБ")
				t719[2] = rt.Text("ГиБ")
				t719[3] = rt.Text("ТиБ")
				t719[4] = rt.Text("ПиБ")
				// «элемент»
				t720, e721 := rt.BElement(ctx, t717, rt.List(t719))
				if e721 != nil {
					return rt.Value{}, e721
				}
				return t720, nil
			}
		}
	}
}

// BaytyYazykom — функция flang «Байты языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Параметр razdelitel — «разделитель»: строка.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func BaytyYazykom(ctx *rt.Ctx, skolko rt.Value, razdelitel rt.Value, yazyk rt.Value) (rt.Value, error) {
	t722, e723 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e723 != nil {
		return rt.Value{}, e723
	}
	t724, e725 := rt.Cond(ctx, t722)
	if e725 != nil {
		return rt.Value{}, e725
	}
	var t726 rt.Value
	if t724 {
		t726 = rt.Text("-")
	} else {
		t726 = rt.Text("")
	}
	// пусть «знак»
	znak := t726
	t727, e728 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e728 != nil {
		return rt.Value{}, e728
	}
	t729, e730 := rt.Cond(ctx, t727)
	if e730 != nil {
		return rt.Value{}, e730
	}
	var t731 rt.Value
	if t729 {
		t732, e733 := rt.Sub(ctx, rt.Number(0.0), skolko)
		if e733 != nil {
			return rt.Value{}, e733
		}
		t731 = t732
	} else {
		t731 = skolko
	}
	t734, e735 := CeloeKNulyu(ctx, t731)
	if e735 != nil {
		return rt.Value{}, e735
	}
	// пусть «величина»
	velichina := t734
	t736, e737 := rt.Lt(ctx, velichina, rt.Number(1024.0))
	if e737 != nil {
		return rt.Value{}, e737
	}
	t738, e739 := rt.Cond(ctx, t736)
	if e739 != nil {
		return rt.Value{}, e739
	}
	if t738 {
		// «к строке»
		t740, e741 := rt.BToString(ctx, velichina)
		if e741 != nil {
			return rt.Value{}, e741
		}
		t742, e743 := rt.Concat(ctx, znak, t740)
		if e743 != nil {
			return rt.Value{}, e743
		}
		t744, e745 := rt.Concat(ctx, t742, rt.Text(" "))
		if e745 != nil {
			return rt.Value{}, e745
		}
		t746, e747 := EdinicaBaytYazykom(ctx, rt.Number(-1.0), yazyk)
		if e747 != nil {
			return rt.Value{}, e747
		}
		t748, e749 := rt.Concat(ctx, t744, t746)
		if e749 != nil {
			return rt.Value{}, e749
		}
		return t748, nil
	} else {
		t750, e751 := DelenieNacelo(ctx, velichina, rt.Number(1024.0))
		if e751 != nil {
			return rt.Value{}, e751
		}
		// пусть «витки»
		vitki := t750
		t752, e753 := DelitelBayt(ctx, vitki, rt.Number(4.0))
		if e753 != nil {
			return rt.Value{}, e753
		}
		t754, e755 := rt.Div(ctx, velichina, t752)
		if e755 != nil {
			return rt.Value{}, e755
		}
		t756, e757 := Drobyu(ctx, t754, rt.Number(1.0), razdelitel)
		if e757 != nil {
			return rt.Value{}, e757
		}
		t758, e759 := rt.Concat(ctx, znak, t756)
		if e759 != nil {
			return rt.Value{}, e759
		}
		t760, e761 := rt.Concat(ctx, t758, rt.Text(" "))
		if e761 != nil {
			return rt.Value{}, e761
		}
		t762, e763 := StupenBayt(ctx, vitki, rt.Number(4.0))
		if e763 != nil {
			return rt.Value{}, e763
		}
		t764, e765 := EdinicaBaytYazykom(ctx, t762, yazyk)
		if e765 != nil {
			return rt.Value{}, e765
		}
		t766, e767 := rt.Concat(ctx, t760, t764)
		if e767 != nil {
			return rt.Value{}, e767
		}
		return t766, nil
	}
}

// BaytyKakEst — функция flang «Байты как есть».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func BaytyKakEst(ctx *rt.Ctx, skolko rt.Value, yazyk rt.Value) (rt.Value, error) {
	t768, e769 := RazryadyYazykom(ctx, skolko, yazyk)
	if e769 != nil {
		return rt.Value{}, e769
	}
	t770, e771 := rt.Concat(ctx, t768, rt.Text(" "))
	if e771 != nil {
		return rt.Value{}, e771
	}
	t772, e773 := EdinicaBaytYazykom(ctx, rt.Number(-1.0), yazyk)
	if e773 != nil {
		return rt.Value{}, e773
	}
	t774, e775 := rt.Concat(ctx, t770, t772)
	if e775 != nil {
		return rt.Value{}, e775
	}
	return t774, nil
}

// SlovoDney — функция flang «Слово дней».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func SlovoDney(ctx *rt.Ctx, yazyk rt.Value) (rt.Value, error) {
	t776, e777 := Angliyskiy(ctx, yazyk)
	if e777 != nil {
		return rt.Value{}, e777
	}
	t778, e779 := rt.Cond(ctx, t776)
	if e779 != nil {
		return rt.Value{}, e779
	}
	if t778 {
		return rt.Text("d"), nil
	} else {
		return rt.Text("дн"), nil
	}
}

// Dney — функция flang «Дней».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр skolko — «сколько»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func Dney(ctx *rt.Ctx, skolko rt.Value, yazyk rt.Value) (rt.Value, error) {
	t780, e781 := DrobyuYazykom(ctx, skolko, rt.Number(0.0), yazyk)
	if e781 != nil {
		return rt.Value{}, e781
	}
	t782, e783 := rt.Concat(ctx, t780, rt.Text(" "))
	if e783 != nil {
		return rt.Value{}, e783
	}
	t784, e785 := SlovoDney(ctx, yazyk)
	if e785 != nil {
		return rt.Value{}, e785
	}
	t786, e787 := rt.Concat(ctx, t782, t784)
	if e787 != nil {
		return rt.Value{}, e787
	}
	return t786, nil
}

// Nazad — функция flang «Назад».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sekundy — «секунды»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func Nazad(ctx *rt.Ctx, sekundy rt.Value, yazyk rt.Value) (rt.Value, error) {
	t788, e789 := rt.Lt(ctx, sekundy, rt.Number(60.0))
	if e789 != nil {
		return rt.Value{}, e789
	}
	t790, e791 := rt.Cond(ctx, t788)
	if e791 != nil {
		return rt.Value{}, e791
	}
	if t790 {
		t792, e793 := CeloeKNulyu(ctx, sekundy)
		if e793 != nil {
			return rt.Value{}, e793
		}
		// «к строке»
		t794, e795 := rt.BToString(ctx, t792)
		if e795 != nil {
			return rt.Value{}, e795
		}
		t796, e797 := rt.Concat(ctx, t794, rt.Text(" "))
		if e797 != nil {
			return rt.Value{}, e797
		}
		t798, e799 := Angliyskiy(ctx, yazyk)
		if e799 != nil {
			return rt.Value{}, e799
		}
		t800, e801 := rt.Cond(ctx, t798)
		if e801 != nil {
			return rt.Value{}, e801
		}
		var t802 rt.Value
		if t800 {
			t802 = rt.Text("s")
		} else {
			t802 = rt.Text("с")
		}
		t803, e804 := rt.Concat(ctx, t796, t802)
		if e804 != nil {
			return rt.Value{}, e804
		}
		return t803, nil
	} else {
		t805, e806 := rt.Lt(ctx, sekundy, rt.Number(3600.0))
		if e806 != nil {
			return rt.Value{}, e806
		}
		t807, e808 := rt.Cond(ctx, t805)
		if e808 != nil {
			return rt.Value{}, e808
		}
		if t807 {
			t809, e810 := rt.Div(ctx, sekundy, rt.Number(60.0))
			if e810 != nil {
				return rt.Value{}, e810
			}
			t811, e812 := CeloeKNulyu(ctx, t809)
			if e812 != nil {
				return rt.Value{}, e812
			}
			// «к строке»
			t813, e814 := rt.BToString(ctx, t811)
			if e814 != nil {
				return rt.Value{}, e814
			}
			t815, e816 := rt.Concat(ctx, t813, rt.Text(" "))
			if e816 != nil {
				return rt.Value{}, e816
			}
			t817, e818 := Angliyskiy(ctx, yazyk)
			if e818 != nil {
				return rt.Value{}, e818
			}
			t819, e820 := rt.Cond(ctx, t817)
			if e820 != nil {
				return rt.Value{}, e820
			}
			var t821 rt.Value
			if t819 {
				t821 = rt.Text("min")
			} else {
				t821 = rt.Text("мин")
			}
			t822, e823 := rt.Concat(ctx, t815, t821)
			if e823 != nil {
				return rt.Value{}, e823
			}
			return t822, nil
		} else {
			t824, e825 := rt.Lt(ctx, sekundy, rt.Number(172800.0))
			if e825 != nil {
				return rt.Value{}, e825
			}
			t826, e827 := rt.Cond(ctx, t824)
			if e827 != nil {
				return rt.Value{}, e827
			}
			if t826 {
				t828, e829 := rt.Div(ctx, sekundy, rt.Number(3600.0))
				if e829 != nil {
					return rt.Value{}, e829
				}
				t830, e831 := CeloeKNulyu(ctx, t828)
				if e831 != nil {
					return rt.Value{}, e831
				}
				// «к строке»
				t832, e833 := rt.BToString(ctx, t830)
				if e833 != nil {
					return rt.Value{}, e833
				}
				t834, e835 := rt.Concat(ctx, t832, rt.Text(" "))
				if e835 != nil {
					return rt.Value{}, e835
				}
				t836, e837 := Angliyskiy(ctx, yazyk)
				if e837 != nil {
					return rt.Value{}, e837
				}
				t838, e839 := rt.Cond(ctx, t836)
				if e839 != nil {
					return rt.Value{}, e839
				}
				var t840 rt.Value
				if t838 {
					t840 = rt.Text("h")
				} else {
					t840 = rt.Text("ч")
				}
				t841, e842 := rt.Concat(ctx, t834, t840)
				if e842 != nil {
					return rt.Value{}, e842
				}
				return t841, nil
			} else {
				t843, e844 := rt.Div(ctx, sekundy, rt.Number(86400.0))
				if e844 != nil {
					return rt.Value{}, e844
				}
				t845, e846 := CeloeKNulyu(ctx, t843)
				if e846 != nil {
					return rt.Value{}, e846
				}
				// «к строке»
				t847, e848 := rt.BToString(ctx, t845)
				if e848 != nil {
					return rt.Value{}, e848
				}
				t849, e850 := rt.Concat(ctx, t847, rt.Text(" "))
				if e850 != nil {
					return rt.Value{}, e850
				}
				t851, e852 := Angliyskiy(ctx, yazyk)
				if e852 != nil {
					return rt.Value{}, e852
				}
				t853, e854 := rt.Cond(ctx, t851)
				if e854 != nil {
					return rt.Value{}, e854
				}
				var t855 rt.Value
				if t853 {
					t855 = rt.Text("d")
				} else {
					t855 = rt.Text("дн")
				}
				t856, e857 := rt.Concat(ctx, t849, t855)
				if e857 != nil {
					return rt.Value{}, e857
				}
				return t856, nil
			}
		}
	}
}

// ZanyaloYazykom — функция flang «Заняло языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр millisekundy — «миллисекунды»: число.
// Параметр razdelitel — «разделитель»: строка.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func ZanyaloYazykom(ctx *rt.Ctx, millisekundy rt.Value, razdelitel rt.Value, yazyk rt.Value) (rt.Value, error) {
	t858, e859 := rt.Lt(ctx, millisekundy, rt.Number(1000.0))
	if e859 != nil {
		return rt.Value{}, e859
	}
	t860, e861 := rt.Cond(ctx, t858)
	if e861 != nil {
		return rt.Value{}, e861
	}
	if t860 {
		t862, e863 := CeloeKNulyu(ctx, millisekundy)
		if e863 != nil {
			return rt.Value{}, e863
		}
		// «к строке»
		t864, e865 := rt.BToString(ctx, t862)
		if e865 != nil {
			return rt.Value{}, e865
		}
		t866, e867 := rt.Concat(ctx, t864, rt.Text(" "))
		if e867 != nil {
			return rt.Value{}, e867
		}
		t868, e869 := Angliyskiy(ctx, yazyk)
		if e869 != nil {
			return rt.Value{}, e869
		}
		t870, e871 := rt.Cond(ctx, t868)
		if e871 != nil {
			return rt.Value{}, e871
		}
		var t872 rt.Value
		if t870 {
			t872 = rt.Text("ms")
		} else {
			t872 = rt.Text("мс")
		}
		t873, e874 := rt.Concat(ctx, t866, t872)
		if e874 != nil {
			return rt.Value{}, e874
		}
		return t873, nil
	} else {
		t875, e876 := rt.Div(ctx, millisekundy, rt.Number(1000.0))
		if e876 != nil {
			return rt.Value{}, e876
		}
		t877, e878 := Drobyu(ctx, t875, rt.Number(1.0), razdelitel)
		if e878 != nil {
			return rt.Value{}, e878
		}
		t879, e880 := rt.Concat(ctx, t877, rt.Text(" "))
		if e880 != nil {
			return rt.Value{}, e880
		}
		t881, e882 := Angliyskiy(ctx, yazyk)
		if e882 != nil {
			return rt.Value{}, e882
		}
		t883, e884 := rt.Cond(ctx, t881)
		if e884 != nil {
			return rt.Value{}, e884
		}
		var t885 rt.Value
		if t883 {
			t885 = rt.Text("s")
		} else {
			t885 = rt.Text("с")
		}
		t886, e887 := rt.Concat(ctx, t879, t885)
		if e887 != nil {
			return rt.Value{}, e887
		}
		return t886, nil
	}
}

// BezHvostovyhNuley — функция flang «Без хвостовых нулей».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Хвостовой самовызов развёрнут в цикл: стек не растёт.
//
// Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
//
// Параметр tekst — «текст»: строка.
// Параметр zapas — «запас»: «нат».
// Результат — значение: строка.
func BezHvostovyhNuley(ctx *rt.Ctx, tekst rt.Value, zapas rt.Value) (rt.Value, error) {
	if err := ctx.Enter("Без хвостовых нулей"); err != nil {
		return rt.Value{}, err
	}
	defer ctx.Leave()
	for {
		t888, e889 := rt.Lte(ctx, zapas, rt.Number(0.0))
		if e889 != nil {
			return rt.Value{}, e889
		}
		t890, e891 := rt.Cond(ctx, t888)
		if e891 != nil {
			return rt.Value{}, e891
		}
		if t890 {
			return tekst, nil
		} else {
			// «длина»
			t892, e893 := rt.BLength(ctx, tekst)
			if e893 != nil {
				return rt.Value{}, e893
			}
			t894, e895 := rt.Lte(ctx, t892, rt.Number(0.0))
			if e895 != nil {
				return rt.Value{}, e895
			}
			t896, e897 := rt.Cond(ctx, t894)
			if e897 != nil {
				return rt.Value{}, e897
			}
			if t896 {
				return tekst, nil
			} else {
				// «длина»
				t898, e899 := rt.BLength(ctx, tekst)
				if e899 != nil {
					return rt.Value{}, e899
				}
				// «длина»
				t900, e901 := rt.BLength(ctx, tekst)
				if e901 != nil {
					return rt.Value{}, e901
				}
				// «подстрока»
				t902, e903 := rt.BSubstring(ctx, tekst, t898, t900)
				if e903 != nil {
					return rt.Value{}, e903
				}
				t904, e905 := rt.Cond(ctx, rt.Flag(rt.Equal(t902, rt.Text("0"))))
				if e905 != nil {
					return rt.Value{}, e905
				}
				if t904 {
					// «длина»
					t906, e907 := rt.BLength(ctx, tekst)
					if e907 != nil {
						return rt.Value{}, e907
					}
					t908, e909 := rt.Sub(ctx, t906, rt.Number(1.0))
					if e909 != nil {
						return rt.Value{}, e909
					}
					// «подстрока»
					t910, e911 := rt.BSubstring(ctx, tekst, rt.Number(1.0), t908)
					if e911 != nil {
						return rt.Value{}, e911
					}
					t912, e913 := rt.Sub(ctx, zapas, rt.Number(1.0))
					if e913 != nil {
						return rt.Value{}, e913
					}
					tekst = t910
					zapas = t912
					// виток цикла — тоже шаг вычисления: незавершающийся хвостовой
					// самовызов обязан упереться в лимит, а не крутиться вечно
					if e914 := ctx.Step("Без хвостовых нулей"); e914 != nil {
						return rt.Value{}, e914
					}
					continue
				} else {
					return tekst, nil
				}
			}
		}
	}
}

// BezHvostovogoRazdelitelya — функция flang «Без хвостового разделителя».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр tekst — «текст»: строка.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func BezHvostovogoRazdelitelya(ctx *rt.Ctx, tekst rt.Value, razdelitel rt.Value) (rt.Value, error) {
	// «длина»
	t915, e916 := rt.BLength(ctx, tekst)
	if e916 != nil {
		return rt.Value{}, e916
	}
	t917, e918 := rt.Lte(ctx, t915, rt.Number(0.0))
	if e918 != nil {
		return rt.Value{}, e918
	}
	t919, e920 := rt.Cond(ctx, t917)
	if e920 != nil {
		return rt.Value{}, e920
	}
	var t921 rt.Value
	if t919 {
		t921 = tekst
	} else {
		// «длина»
		t922, e923 := rt.BLength(ctx, tekst)
		if e923 != nil {
			return rt.Value{}, e923
		}
		// «длина»
		t924, e925 := rt.BLength(ctx, tekst)
		if e925 != nil {
			return rt.Value{}, e925
		}
		// «подстрока»
		t926, e927 := rt.BSubstring(ctx, tekst, t922, t924)
		if e927 != nil {
			return rt.Value{}, e927
		}
		t928, e929 := rt.Cond(ctx, rt.Flag(rt.Equal(t926, razdelitel)))
		if e929 != nil {
			return rt.Value{}, e929
		}
		var t930 rt.Value
		if t928 {
			// «длина»
			t931, e932 := rt.BLength(ctx, tekst)
			if e932 != nil {
				return rt.Value{}, e932
			}
			t933, e934 := rt.Sub(ctx, t931, rt.Number(1.0))
			if e934 != nil {
				return rt.Value{}, e934
			}
			// «подстрока»
			t935, e936 := rt.BSubstring(ctx, tekst, rt.Number(1.0), t933)
			if e936 != nil {
				return rt.Value{}, e936
			}
			t930 = t935
		} else {
			t930 = tekst
		}
		t921 = t930
	}
	t937 := t921
	// «длина»
	t938, e939 := rt.BLength(ctx, t937)
	if e939 != nil {
		return rt.Value{}, e939
	}
	// «длина»
	t940, e941 := rt.BLength(ctx, tekst)
	if e941 != nil {
		return rt.Value{}, e941
	}
	t942, e943 := rt.Lte(ctx, t938, t940)
	if e943 != nil {
		return rt.Value{}, e943
	}
	// постусловие «длиннее не становится»
	t944, e945 := rt.Post(ctx, t942, "длиннее не становится", "Без хвостового разделителя")
	if e945 != nil {
		return rt.Value{}, e945
	}
	if !t944 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «длиннее не становится» функции «Без хвостового разделителя»")
	}
	return t937, nil
}

// Kazhdye — функция flang «Каждые».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр millisekundy — «миллисекунды»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func Kazhdye(ctx *rt.Ctx, millisekundy rt.Value, yazyk rt.Value) (rt.Value, error) {
	t946, e947 := rt.Lt(ctx, millisekundy, rt.Number(1000.0))
	if e947 != nil {
		return rt.Value{}, e947
	}
	t948, e949 := rt.Cond(ctx, t946)
	if e949 != nil {
		return rt.Value{}, e949
	}
	if t948 {
		t950, e951 := CeloeKNulyu(ctx, millisekundy)
		if e951 != nil {
			return rt.Value{}, e951
		}
		// «к строке»
		t952, e953 := rt.BToString(ctx, t950)
		if e953 != nil {
			return rt.Value{}, e953
		}
		t954, e955 := rt.Concat(ctx, t952, rt.Text(" "))
		if e955 != nil {
			return rt.Value{}, e955
		}
		t956, e957 := Angliyskiy(ctx, yazyk)
		if e957 != nil {
			return rt.Value{}, e957
		}
		t958, e959 := rt.Cond(ctx, t956)
		if e959 != nil {
			return rt.Value{}, e959
		}
		var t960 rt.Value
		if t958 {
			t960 = rt.Text("ms")
		} else {
			t960 = rt.Text("мс")
		}
		t961, e962 := rt.Concat(ctx, t954, t960)
		if e962 != nil {
			return rt.Value{}, e962
		}
		return t961, nil
	} else {
		t963, e964 := Angliyskiy(ctx, yazyk)
		if e964 != nil {
			return rt.Value{}, e964
		}
		t965, e966 := rt.Cond(ctx, t963)
		if e966 != nil {
			return rt.Value{}, e966
		}
		var t967 rt.Value
		if t965 {
			t967 = rt.Text(".")
		} else {
			t967 = rt.Text(",")
		}
		// пусть «разделитель»
		razdelitel := t967
		t968, e969 := rt.Div(ctx, millisekundy, rt.Number(1000.0))
		if e969 != nil {
			return rt.Value{}, e969
		}
		t970, e971 := Drobyu(ctx, t968, rt.Number(3.0), razdelitel)
		if e971 != nil {
			return rt.Value{}, e971
		}
		// пусть «с тремя знаками»
		sTremyaZnakami := t970
		t972, e973 := BezHvostovyhNuley(ctx, sTremyaZnakami, rt.Number(3.0))
		if e973 != nil {
			return rt.Value{}, e973
		}
		// пусть «без нулей»
		bezNuley := t972
		t974, e975 := BezHvostovogoRazdelitelya(ctx, bezNuley, razdelitel)
		if e975 != nil {
			return rt.Value{}, e975
		}
		t976, e977 := rt.Concat(ctx, t974, rt.Text(" "))
		if e977 != nil {
			return rt.Value{}, e977
		}
		t978, e979 := Angliyskiy(ctx, yazyk)
		if e979 != nil {
			return rt.Value{}, e979
		}
		t980, e981 := rt.Cond(ctx, t978)
		if e981 != nil {
			return rt.Value{}, e981
		}
		var t982 rt.Value
		if t980 {
			t982 = rt.Text("s")
		} else {
			t982 = rt.Text("с")
		}
		t983, e984 := rt.Concat(ctx, t976, t982)
		if e984 != nil {
			return rt.Value{}, e984
		}
		return t983, nil
	}
}

// DlitelnostYazykom — функция flang «Длительность языком».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sekundy — «секунды»: число.
// Параметр yazyk — «язык»: строка.
// Результат — значение: строка.
func DlitelnostYazykom(ctx *rt.Ctx, sekundy rt.Value, yazyk rt.Value) (rt.Value, error) {
	t985, e986 := rt.Lte(ctx, sekundy, rt.Number(0.0))
	if e986 != nil {
		return rt.Value{}, e986
	}
	t987, e988 := rt.Cond(ctx, t985)
	if e988 != nil {
		return rt.Value{}, e988
	}
	if t987 {
		return rt.Text(""), nil
	} else {
		t989, e990 := CeloeKNulyu(ctx, sekundy)
		if e990 != nil {
			return rt.Value{}, e990
		}
		// пусть «всего»
		vsego := t989
		t991, e992 := DelenieNacelo(ctx, vsego, rt.Number(3600.0))
		if e992 != nil {
			return rt.Value{}, e992
		}
		// пусть «часов всего»
		chasovVsego := t991
		t993, e994 := DelenieNacelo(ctx, chasovVsego, rt.Number(24.0))
		if e994 != nil {
			return rt.Value{}, e994
		}
		// пусть «дней»
		dney := t993
		t995, e996 := rt.Mod(ctx, chasovVsego, rt.Number(24.0))
		if e996 != nil {
			return rt.Value{}, e996
		}
		// пусть «часов»
		chasov := t995
		t997, e998 := DelenieNacelo(ctx, vsego, rt.Number(60.0))
		if e998 != nil {
			return rt.Value{}, e998
		}
		t999, e1000 := rt.Mod(ctx, t997, rt.Number(60.0))
		if e1000 != nil {
			return rt.Value{}, e1000
		}
		// пусть «минут»
		minut := t999
		// «к строке»
		t1001, e1002 := rt.BToString(ctx, chasov)
		if e1002 != nil {
			return rt.Value{}, e1002
		}
		t1003, e1004 := SlevaNulyami(ctx, t1001, rt.Number(2.0))
		if e1004 != nil {
			return rt.Value{}, e1004
		}
		t1005, e1006 := rt.Concat(ctx, t1003, rt.Text(":"))
		if e1006 != nil {
			return rt.Value{}, e1006
		}
		// «к строке»
		t1007, e1008 := rt.BToString(ctx, minut)
		if e1008 != nil {
			return rt.Value{}, e1008
		}
		t1009, e1010 := SlevaNulyami(ctx, t1007, rt.Number(2.0))
		if e1010 != nil {
			return rt.Value{}, e1010
		}
		t1011, e1012 := rt.Concat(ctx, t1005, t1009)
		if e1012 != nil {
			return rt.Value{}, e1012
		}
		// пусть «время»
		vremya := t1011
		t1013, e1014 := rt.Gt(ctx, dney, rt.Number(0.0))
		if e1014 != nil {
			return rt.Value{}, e1014
		}
		t1015, e1016 := rt.Cond(ctx, t1013)
		if e1016 != nil {
			return rt.Value{}, e1016
		}
		if t1015 {
			// «к строке»
			t1017, e1018 := rt.BToString(ctx, dney)
			if e1018 != nil {
				return rt.Value{}, e1018
			}
			t1019, e1020 := Angliyskiy(ctx, yazyk)
			if e1020 != nil {
				return rt.Value{}, e1020
			}
			t1021, e1022 := rt.Cond(ctx, t1019)
			if e1022 != nil {
				return rt.Value{}, e1022
			}
			var t1023 rt.Value
			if t1021 {
				t1023 = rt.Text("d")
			} else {
				t1023 = rt.Text("д")
			}
			t1024, e1025 := rt.Concat(ctx, t1017, t1023)
			if e1025 != nil {
				return rt.Value{}, e1025
			}
			t1026, e1027 := rt.Concat(ctx, t1024, rt.Text(" "))
			if e1027 != nil {
				return rt.Value{}, e1027
			}
			t1028, e1029 := rt.Concat(ctx, t1026, vremya)
			if e1029 != nil {
				return rt.Value{}, e1029
			}
			return t1028, nil
		} else {
			return vremya, nil
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
	t1030 := shag
	t1031, e1032 := rt.Lt(ctx, t1030, mera)
	if e1032 != nil {
		return rt.Value{}, e1032
	}
	// постусловие «мера убывает»
	t1033, e1034 := rt.Post(ctx, t1031, "мера убывает", "мера убывает")
	if e1034 != nil {
		return rt.Value{}, e1034
	}
	if !t1033 {
		return rt.Value{}, rt.Fail("FLANG_MEASURE", "%s", "тотальная функция «Повторить»: мера не убыла — аргумент 2 вызова «Повторить» не стал меньше параметра «сколько». Завершение доказано убыванием этой меры, а числа flang — IEEE-754 double: при большом |«сколько»| постоянный шаг не меняет значение, и спуск не идёт. Отказ здесь честнее зацикливания")
	}
	return t1030, nil
}

// Call — вызов функции по её исходному имени flang.
//
// Нужен прогонщику и всякому, кто связывает программу с внешним миром
// динамически (скрипт, тест, служба). Коды и тексты — те же, что у
// интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
func Call(ctx *rt.Ctx, name string, args []rt.Value) (rt.Value, error) {
	switch name {
	case "Целое к нулю":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Целое к нулю", 1, len(args))
		}
		return CeloeKNulyu(ctx, args[0])
	case "Целое вниз":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Целое вниз", 1, len(args))
		}
		return CeloeVniz(ctx, args[0])
	case "Деление нацело":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Деление нацело", 2, len(args))
		}
		return DelenieNacelo(ctx, args[0], args[1])
	case "Десять в степени":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Десять в степени", 1, len(args))
		}
		return DesyatVStepeni(ctx, args[0])
	case "Повторить":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Повторить", 2, len(args))
		}
		return Povtorit(ctx, args[0], args[1])
	case "Слева нулями":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Слева нулями", 2, len(args))
		}
		return SlevaNulyami(ctx, args[0], args[1])
	case "Верхняя половина":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Верхняя половина", 1, len(args))
		}
		return VerhnyayaPolovina(ctx, args[0])
	case "Погрешность произведения":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Погрешность произведения", 2, len(args))
		}
		return PogreshnostProizvedeniya(ctx, args[0], args[1])
	case "Округлить к чётному":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Округлить к чётному", 2, len(args))
		}
		return OkruglitKChyotnomu(ctx, args[0], args[1])
	case "Знак минус":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Знак минус", 1, len(args))
		}
		return ZnakMinus(ctx, args[0])
	case "Дробью":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Дробью", 3, len(args))
		}
		return Drobyu(ctx, args[0], args[1], args[2])
	case "Дробью неотрицательного":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Дробью неотрицательного", 3, len(args))
		}
		return DrobyuNeotricatelnogo(ctx, args[0], args[1], args[2])
	case "Ступень байт":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Ступень байт", 2, len(args))
		}
		return StupenBayt(ctx, args[0], args[1])
	case "Делитель байт":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Делитель байт", 2, len(args))
		}
		return DelitelBayt(ctx, args[0], args[1])
	case "Единица байт":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Единица байт", 1, len(args))
		}
		return EdinicaBayt(ctx, args[0])
	case "Байты знаком":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Байты знаком", 2, len(args))
		}
		return BaytyZnakom(ctx, args[0], args[1])
	case "Байты":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Байты", 1, len(args))
		}
		return Bayty(ctx, args[0])
	case "Ступень тысяч":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Ступень тысяч", 2, len(args))
		}
		return StupenTysyach(ctx, args[0], args[1])
	case "Делитель тысяч":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Делитель тысяч", 2, len(args))
		}
		return DelitelTysyach(ctx, args[0], args[1])
	case "Единица тысяч":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Единица тысяч", 1, len(args))
		}
		return EdinicaTysyach(ctx, args[0])
	case "Байты десятично":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Байты десятично", 2, len(args))
		}
		return BaytyDesyatichno(ctx, args[0], args[1])
	case "Доля от":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Доля от", 2, len(args))
		}
		return DolyaOt(ctx, args[0], args[1])
	case "Процент целого":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Процент целого", 2, len(args))
		}
		return ProcentCelogo(ctx, args[0], args[1])
	case "Проценты знаком":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Проценты знаком", 2, len(args))
		}
		return ProcentyZnakom(ctx, args[0], args[1])
	case "Проценты":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Проценты", 1, len(args))
		}
		return Procenty(ctx, args[0])
	case "Длительность":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Длительность", 1, len(args))
		}
		return Dlitelnost(ctx, args[0])
	case "Прошло":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прошло", 1, len(args))
		}
		return Proshlo(ctx, args[0])
	case "Заняло знаком":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Заняло знаком", 2, len(args))
		}
		return ZanyaloZnakom(ctx, args[0], args[1])
	case "Заняло":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Заняло", 1, len(args))
		}
		return Zanyalo(ctx, args[0])
	case "Шаг разряда":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Шаг разряда", 2, len(args))
		}
		return ShagRazryada(ctx, args[0], args[1])
	case "Разрядами":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разрядами", 2, len(args))
		}
		return Razryadami(ctx, args[0], args[1])
	case "Пробельный":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Пробельный", 1, len(args))
		}
		return Probelnyy(ctx, args[0])
	case "Прочерк":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Прочерк", 1, len(args))
		}
		return Procherk(ctx, args[0])
	case "Обрезать":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Обрезать", 2, len(args))
		}
		return Obrezat(ctx, args[0], args[1])
	case "Уместить":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Уместить", 2, len(args))
		}
		return Umestit(ctx, args[0], args[1])
	case "Справа":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Справа", 2, len(args))
		}
		return Sprava(ctx, args[0], args[1])
	case "Английский":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Английский", 1, len(args))
		}
		return Angliyskiy(ctx, args[0])
	case "Дробью языком":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Дробью языком", 3, len(args))
		}
		return DrobyuYazykom(ctx, args[0], args[1], args[2])
	case "Процентом":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Процентом", 3, len(args))
		}
		return Procentom(ctx, args[0], args[1], args[2])
	case "Процентом языком":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Процентом языком", 3, len(args))
		}
		return ProcentomYazykom(ctx, args[0], args[1], args[2])
	case "Разряды языком":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряды языком", 2, len(args))
		}
		return RazryadyYazykom(ctx, args[0], args[1])
	case "Единица байт языком":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Единица байт языком", 2, len(args))
		}
		return EdinicaBaytYazykom(ctx, args[0], args[1])
	case "Байты языком":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Байты языком", 3, len(args))
		}
		return BaytyYazykom(ctx, args[0], args[1], args[2])
	case "Байты как есть":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Байты как есть", 2, len(args))
		}
		return BaytyKakEst(ctx, args[0], args[1])
	case "Слово дней":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Слово дней", 1, len(args))
		}
		return SlovoDney(ctx, args[0])
	case "Дней":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Дней", 2, len(args))
		}
		return Dney(ctx, args[0], args[1])
	case "Назад":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Назад", 2, len(args))
		}
		return Nazad(ctx, args[0], args[1])
	case "Заняло языком":
		if len(args) != 3 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Заняло языком", 3, len(args))
		}
		return ZanyaloYazykom(ctx, args[0], args[1], args[2])
	case "Без хвостовых нулей":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Без хвостовых нулей", 2, len(args))
		}
		return BezHvostovyhNuley(ctx, args[0], args[1])
	case "Без хвостового разделителя":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Без хвостового разделителя", 2, len(args))
		}
		return BezHvostovogoRazdelitelya(ctx, args[0], args[1])
	case "Каждые":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Каждые", 2, len(args))
		}
		return Kazhdye(ctx, args[0], args[1])
	case "Длительность языком":
		if len(args) != 2 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Длительность языком", 2, len(args))
		}
		return DlitelnostYazykom(ctx, args[0], args[1])
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
