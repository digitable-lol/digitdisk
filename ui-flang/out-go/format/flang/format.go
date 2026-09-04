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

// Drobyu — функция flang «Дробью».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр znakov — «знаков»: «нат».
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func Drobyu(ctx *rt.Ctx, znachenie rt.Value, znakov rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t178, e179 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e179 != nil {
		return rt.Value{}, e179
	}
	t180, e181 := rt.Cond(ctx, t178)
	if e181 != nil {
		return rt.Value{}, e181
	}
	if t180 {
		t182, e183 := rt.Sub(ctx, rt.Number(0.0), znachenie)
		if e183 != nil {
			return rt.Value{}, e183
		}
		t184, e185 := DrobyuNeotricatelnogo(ctx, t182, znakov, razdelitel)
		if e185 != nil {
			return rt.Value{}, e185
		}
		t186, e187 := rt.Concat(ctx, rt.Text("-"), t184)
		if e187 != nil {
			return rt.Value{}, e187
		}
		return t186, nil
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
	t188, e189 := DesyatVStepeni(ctx, znakov)
	if e189 != nil {
		return rt.Value{}, e189
	}
	// пусть «степень»
	stepen := t188
	t190, e191 := rt.Mul(ctx, znachenie, stepen)
	if e191 != nil {
		return rt.Value{}, e191
	}
	t192, e193 := PogreshnostProizvedeniya(ctx, znachenie, stepen)
	if e193 != nil {
		return rt.Value{}, e193
	}
	t194, e195 := OkruglitKChyotnomu(ctx, t190, t192)
	if e195 != nil {
		return rt.Value{}, e195
	}
	// пусть «целое»
	celoe := t194
	t196, e197 := DelenieNacelo(ctx, celoe, stepen)
	if e197 != nil {
		return rt.Value{}, e197
	}
	// пусть «целых»
	celyh := t196
	t198, e199 := rt.Mul(ctx, celyh, stepen)
	if e199 != nil {
		return rt.Value{}, e199
	}
	t200, e201 := rt.Sub(ctx, celoe, t198)
	if e201 != nil {
		return rt.Value{}, e201
	}
	// пусть «остаток»
	ostatok := t200
	t202, e203 := rt.Lte(ctx, znakov, rt.Number(0.0))
	if e203 != nil {
		return rt.Value{}, e203
	}
	t204, e205 := rt.Cond(ctx, t202)
	if e205 != nil {
		return rt.Value{}, e205
	}
	if t204 {
		// «к строке»
		t206, e207 := rt.BToString(ctx, celyh)
		if e207 != nil {
			return rt.Value{}, e207
		}
		return t206, nil
	} else {
		// «к строке»
		t208, e209 := rt.BToString(ctx, celyh)
		if e209 != nil {
			return rt.Value{}, e209
		}
		t210, e211 := rt.Concat(ctx, t208, razdelitel)
		if e211 != nil {
			return rt.Value{}, e211
		}
		// «к строке»
		t212, e213 := rt.BToString(ctx, ostatok)
		if e213 != nil {
			return rt.Value{}, e213
		}
		t214, e215 := SlevaNulyami(ctx, t212, znakov)
		if e215 != nil {
			return rt.Value{}, e215
		}
		t216, e217 := rt.Concat(ctx, t210, t214)
		if e217 != nil {
			return rt.Value{}, e217
		}
		return t216, nil
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
	t218, e219 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e219 != nil {
		return rt.Value{}, e219
	}
	t220, e221 := rt.Cond(ctx, t218)
	if e221 != nil {
		return rt.Value{}, e221
	}
	if t220 {
		return rt.Number(0.0), nil
	} else {
		t222, e223 := rt.Gte(ctx, vitki, rt.Number(1024.0))
		if e223 != nil {
			return rt.Value{}, e223
		}
		t224, e225 := rt.Cond(ctx, t222)
		if e225 != nil {
			return rt.Value{}, e225
		}
		if t224 {
			t226, e227 := DelenieNacelo(ctx, vitki, rt.Number(1024.0))
			if e227 != nil {
				return rt.Value{}, e227
			}
			t228, e229 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e229 != nil {
				return rt.Value{}, e229
			}
			t230, e231 := StupenBayt(ctx, t226, t228)
			if e231 != nil {
				return rt.Value{}, e231
			}
			t232, e233 := rt.Add(ctx, rt.Number(1.0), t230)
			if e233 != nil {
				return rt.Value{}, e233
			}
			return t232, nil
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
	t234, e235 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e235 != nil {
		return rt.Value{}, e235
	}
	t236, e237 := rt.Cond(ctx, t234)
	if e237 != nil {
		return rt.Value{}, e237
	}
	if t236 {
		return rt.Number(1024.0), nil
	} else {
		t238, e239 := rt.Gte(ctx, vitki, rt.Number(1024.0))
		if e239 != nil {
			return rt.Value{}, e239
		}
		t240, e241 := rt.Cond(ctx, t238)
		if e241 != nil {
			return rt.Value{}, e241
		}
		if t240 {
			t242, e243 := DelenieNacelo(ctx, vitki, rt.Number(1024.0))
			if e243 != nil {
				return rt.Value{}, e243
			}
			t244, e245 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e245 != nil {
				return rt.Value{}, e245
			}
			t246, e247 := DelitelBayt(ctx, t242, t244)
			if e247 != nil {
				return rt.Value{}, e247
			}
			t248, e249 := rt.Mul(ctx, rt.Number(1024.0), t246)
			if e249 != nil {
				return rt.Value{}, e249
			}
			return t248, nil
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
	t250, e251 := rt.Lte(ctx, stupen, rt.Number(0.0))
	if e251 != nil {
		return rt.Value{}, e251
	}
	t252, e253 := rt.Cond(ctx, t250)
	if e253 != nil {
		return rt.Value{}, e253
	}
	if t252 {
		return rt.Text("КиБ"), nil
	} else {
		t254, e255 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e255 != nil {
			return rt.Value{}, e255
		}
		t256, e257 := rt.Cond(ctx, t254)
		if e257 != nil {
			return rt.Value{}, e257
		}
		if t256 {
			return rt.Text("ПиБ"), nil
		} else {
			t258, e259 := rt.Add(ctx, stupen, rt.Number(1.0))
			if e259 != nil {
				return rt.Value{}, e259
			}
			t260 := make([]rt.Value, 5)
			t260[0] = rt.Text("КиБ")
			t260[1] = rt.Text("МиБ")
			t260[2] = rt.Text("ГиБ")
			t260[3] = rt.Text("ТиБ")
			t260[4] = rt.Text("ПиБ")
			// «элемент»
			t261, e262 := rt.BElement(ctx, t258, rt.List(t260))
			if e262 != nil {
				return rt.Value{}, e262
			}
			return t261, nil
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
	t263, e264 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e264 != nil {
		return rt.Value{}, e264
	}
	t265, e266 := rt.Cond(ctx, t263)
	if e266 != nil {
		return rt.Value{}, e266
	}
	var t267 rt.Value
	if t265 {
		t267 = rt.Text("-")
	} else {
		t267 = rt.Text("")
	}
	// пусть «знак»
	znak := t267
	t268, e269 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e269 != nil {
		return rt.Value{}, e269
	}
	t270, e271 := rt.Cond(ctx, t268)
	if e271 != nil {
		return rt.Value{}, e271
	}
	var t272 rt.Value
	if t270 {
		t273, e274 := rt.Sub(ctx, rt.Number(0.0), skolko)
		if e274 != nil {
			return rt.Value{}, e274
		}
		t272 = t273
	} else {
		t272 = skolko
	}
	t275, e276 := CeloeKNulyu(ctx, t272)
	if e276 != nil {
		return rt.Value{}, e276
	}
	// пусть «величина»
	velichina := t275
	t277, e278 := rt.Lt(ctx, velichina, rt.Number(1024.0))
	if e278 != nil {
		return rt.Value{}, e278
	}
	t279, e280 := rt.Cond(ctx, t277)
	if e280 != nil {
		return rt.Value{}, e280
	}
	if t279 {
		// «к строке»
		t281, e282 := rt.BToString(ctx, velichina)
		if e282 != nil {
			return rt.Value{}, e282
		}
		t283, e284 := rt.Concat(ctx, znak, t281)
		if e284 != nil {
			return rt.Value{}, e284
		}
		t285, e286 := rt.Concat(ctx, t283, rt.Text(" Б"))
		if e286 != nil {
			return rt.Value{}, e286
		}
		return t285, nil
	} else {
		t287, e288 := DelenieNacelo(ctx, velichina, rt.Number(1024.0))
		if e288 != nil {
			return rt.Value{}, e288
		}
		// пусть «витки»
		vitki := t287
		t289, e290 := DelitelBayt(ctx, vitki, rt.Number(4.0))
		if e290 != nil {
			return rt.Value{}, e290
		}
		t291, e292 := rt.Div(ctx, velichina, t289)
		if e292 != nil {
			return rt.Value{}, e292
		}
		t293, e294 := Drobyu(ctx, t291, rt.Number(1.0), razdelitel)
		if e294 != nil {
			return rt.Value{}, e294
		}
		t295, e296 := rt.Concat(ctx, znak, t293)
		if e296 != nil {
			return rt.Value{}, e296
		}
		t297, e298 := rt.Concat(ctx, t295, rt.Text(" "))
		if e298 != nil {
			return rt.Value{}, e298
		}
		t299, e300 := StupenBayt(ctx, vitki, rt.Number(4.0))
		if e300 != nil {
			return rt.Value{}, e300
		}
		t301, e302 := EdinicaBayt(ctx, t299)
		if e302 != nil {
			return rt.Value{}, e302
		}
		t303, e304 := rt.Concat(ctx, t297, t301)
		if e304 != nil {
			return rt.Value{}, e304
		}
		return t303, nil
	}
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
	t305, e306 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e306 != nil {
		return rt.Value{}, e306
	}
	t307, e308 := rt.Cond(ctx, t305)
	if e308 != nil {
		return rt.Value{}, e308
	}
	if t307 {
		return rt.Number(0.0), nil
	} else {
		t309, e310 := rt.Gte(ctx, vitki, rt.Number(1000.0))
		if e310 != nil {
			return rt.Value{}, e310
		}
		t311, e312 := rt.Cond(ctx, t309)
		if e312 != nil {
			return rt.Value{}, e312
		}
		if t311 {
			t313, e314 := DelenieNacelo(ctx, vitki, rt.Number(1000.0))
			if e314 != nil {
				return rt.Value{}, e314
			}
			t315, e316 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e316 != nil {
				return rt.Value{}, e316
			}
			t317, e318 := StupenTysyach(ctx, t313, t315)
			if e318 != nil {
				return rt.Value{}, e318
			}
			t319, e320 := rt.Add(ctx, rt.Number(1.0), t317)
			if e320 != nil {
				return rt.Value{}, e320
			}
			return t319, nil
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
	t321, e322 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e322 != nil {
		return rt.Value{}, e322
	}
	t323, e324 := rt.Cond(ctx, t321)
	if e324 != nil {
		return rt.Value{}, e324
	}
	if t323 {
		return rt.Number(1000.0), nil
	} else {
		t325, e326 := rt.Gte(ctx, vitki, rt.Number(1000.0))
		if e326 != nil {
			return rt.Value{}, e326
		}
		t327, e328 := rt.Cond(ctx, t325)
		if e328 != nil {
			return rt.Value{}, e328
		}
		if t327 {
			t329, e330 := DelenieNacelo(ctx, vitki, rt.Number(1000.0))
			if e330 != nil {
				return rt.Value{}, e330
			}
			t331, e332 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e332 != nil {
				return rt.Value{}, e332
			}
			t333, e334 := DelitelTysyach(ctx, t329, t331)
			if e334 != nil {
				return rt.Value{}, e334
			}
			t335, e336 := rt.Mul(ctx, rt.Number(1000.0), t333)
			if e336 != nil {
				return rt.Value{}, e336
			}
			return t335, nil
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
	t337, e338 := rt.Lte(ctx, stupen, rt.Number(0.0))
	if e338 != nil {
		return rt.Value{}, e338
	}
	t339, e340 := rt.Cond(ctx, t337)
	if e340 != nil {
		return rt.Value{}, e340
	}
	if t339 {
		return rt.Text("КБ"), nil
	} else {
		t341, e342 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e342 != nil {
			return rt.Value{}, e342
		}
		t343, e344 := rt.Cond(ctx, t341)
		if e344 != nil {
			return rt.Value{}, e344
		}
		if t343 {
			return rt.Text("ПБ"), nil
		} else {
			t345, e346 := rt.Add(ctx, stupen, rt.Number(1.0))
			if e346 != nil {
				return rt.Value{}, e346
			}
			t347 := make([]rt.Value, 5)
			t347[0] = rt.Text("КБ")
			t347[1] = rt.Text("МБ")
			t347[2] = rt.Text("ГБ")
			t347[3] = rt.Text("ТБ")
			t347[4] = rt.Text("ПБ")
			// «элемент»
			t348, e349 := rt.BElement(ctx, t345, rt.List(t347))
			if e349 != nil {
				return rt.Value{}, e349
			}
			return t348, nil
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
	t350, e351 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e351 != nil {
		return rt.Value{}, e351
	}
	t352, e353 := rt.Cond(ctx, t350)
	if e353 != nil {
		return rt.Value{}, e353
	}
	var t354 rt.Value
	if t352 {
		t354 = rt.Text("-")
	} else {
		t354 = rt.Text("")
	}
	// пусть «знак»
	znak := t354
	t355, e356 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e356 != nil {
		return rt.Value{}, e356
	}
	t357, e358 := rt.Cond(ctx, t355)
	if e358 != nil {
		return rt.Value{}, e358
	}
	var t359 rt.Value
	if t357 {
		t360, e361 := rt.Sub(ctx, rt.Number(0.0), skolko)
		if e361 != nil {
			return rt.Value{}, e361
		}
		t359 = t360
	} else {
		t359 = skolko
	}
	t362, e363 := CeloeKNulyu(ctx, t359)
	if e363 != nil {
		return rt.Value{}, e363
	}
	// пусть «величина»
	velichina := t362
	t364, e365 := rt.Lt(ctx, velichina, rt.Number(1000.0))
	if e365 != nil {
		return rt.Value{}, e365
	}
	t366, e367 := rt.Cond(ctx, t364)
	if e367 != nil {
		return rt.Value{}, e367
	}
	if t366 {
		// «к строке»
		t368, e369 := rt.BToString(ctx, velichina)
		if e369 != nil {
			return rt.Value{}, e369
		}
		t370, e371 := rt.Concat(ctx, znak, t368)
		if e371 != nil {
			return rt.Value{}, e371
		}
		t372, e373 := rt.Concat(ctx, t370, rt.Text(" Б"))
		if e373 != nil {
			return rt.Value{}, e373
		}
		return t372, nil
	} else {
		t374, e375 := DelenieNacelo(ctx, velichina, rt.Number(1000.0))
		if e375 != nil {
			return rt.Value{}, e375
		}
		// пусть «витки»
		vitki := t374
		t376, e377 := DelitelTysyach(ctx, vitki, rt.Number(4.0))
		if e377 != nil {
			return rt.Value{}, e377
		}
		t378, e379 := rt.Div(ctx, velichina, t376)
		if e379 != nil {
			return rt.Value{}, e379
		}
		t380, e381 := Drobyu(ctx, t378, rt.Number(1.0), razdelitel)
		if e381 != nil {
			return rt.Value{}, e381
		}
		t382, e383 := rt.Concat(ctx, znak, t380)
		if e383 != nil {
			return rt.Value{}, e383
		}
		t384, e385 := rt.Concat(ctx, t382, rt.Text(" "))
		if e385 != nil {
			return rt.Value{}, e385
		}
		t386, e387 := StupenTysyach(ctx, vitki, rt.Number(4.0))
		if e387 != nil {
			return rt.Value{}, e387
		}
		t388, e389 := EdinicaTysyach(ctx, t386)
		if e389 != nil {
			return rt.Value{}, e389
		}
		t390, e391 := rt.Concat(ctx, t384, t388)
		if e391 != nil {
			return rt.Value{}, e391
		}
		return t390, nil
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
	t392, e393 := rt.Cond(ctx, rt.Flag(rt.Equal(celoe, rt.Number(0.0))))
	if e393 != nil {
		return rt.Value{}, e393
	}
	if t392 {
		return rt.Number(0.0), nil
	} else {
		t394, e395 := rt.Div(ctx, chast, celoe)
		if e395 != nil {
			return rt.Value{}, e395
		}
		return t394, nil
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
	t396, e397 := rt.Cond(ctx, rt.Flag(rt.Equal(celoe, rt.Number(0.0))))
	if e397 != nil {
		return rt.Value{}, e397
	}
	if t396 {
		return rt.Number(0.0), nil
	} else {
		t398, e399 := rt.Mul(ctx, rt.Number(100.0), chast)
		if e399 != nil {
			return rt.Value{}, e399
		}
		t400, e401 := rt.Div(ctx, t398, celoe)
		if e401 != nil {
			return rt.Value{}, e401
		}
		return t400, nil
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
	t402, e403 := rt.Mul(ctx, dolya, rt.Number(100.0))
	if e403 != nil {
		return rt.Value{}, e403
	}
	t404, e405 := Drobyu(ctx, t402, rt.Number(1.0), razdelitel)
	if e405 != nil {
		return rt.Value{}, e405
	}
	t406, e407 := rt.Concat(ctx, t404, rt.Text("%"))
	if e407 != nil {
		return rt.Value{}, e407
	}
	return t406, nil
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
	t408, e409 := rt.Lte(ctx, sekundy, rt.Number(0.0))
	if e409 != nil {
		return rt.Value{}, e409
	}
	t410, e411 := rt.Cond(ctx, t408)
	if e411 != nil {
		return rt.Value{}, e411
	}
	if t410 {
		return rt.Text(""), nil
	} else {
		t412, e413 := CeloeKNulyu(ctx, sekundy)
		if e413 != nil {
			return rt.Value{}, e413
		}
		// пусть «всего»
		vsego := t412
		t414, e415 := DelenieNacelo(ctx, vsego, rt.Number(3600.0))
		if e415 != nil {
			return rt.Value{}, e415
		}
		// пусть «часов всего»
		chasovVsego := t414
		t416, e417 := DelenieNacelo(ctx, chasovVsego, rt.Number(24.0))
		if e417 != nil {
			return rt.Value{}, e417
		}
		// пусть «дней»
		dney := t416
		t418, e419 := rt.Mod(ctx, chasovVsego, rt.Number(24.0))
		if e419 != nil {
			return rt.Value{}, e419
		}
		// пусть «часов»
		chasov := t418
		t420, e421 := DelenieNacelo(ctx, vsego, rt.Number(60.0))
		if e421 != nil {
			return rt.Value{}, e421
		}
		t422, e423 := rt.Mod(ctx, t420, rt.Number(60.0))
		if e423 != nil {
			return rt.Value{}, e423
		}
		// пусть «минут»
		minut := t422
		// «к строке»
		t424, e425 := rt.BToString(ctx, chasov)
		if e425 != nil {
			return rt.Value{}, e425
		}
		t426, e427 := SlevaNulyami(ctx, t424, rt.Number(2.0))
		if e427 != nil {
			return rt.Value{}, e427
		}
		t428, e429 := rt.Concat(ctx, t426, rt.Text(":"))
		if e429 != nil {
			return rt.Value{}, e429
		}
		// «к строке»
		t430, e431 := rt.BToString(ctx, minut)
		if e431 != nil {
			return rt.Value{}, e431
		}
		t432, e433 := SlevaNulyami(ctx, t430, rt.Number(2.0))
		if e433 != nil {
			return rt.Value{}, e433
		}
		t434, e435 := rt.Concat(ctx, t428, t432)
		if e435 != nil {
			return rt.Value{}, e435
		}
		// пусть «время»
		vremya := t434
		t436, e437 := rt.Gt(ctx, dney, rt.Number(0.0))
		if e437 != nil {
			return rt.Value{}, e437
		}
		t438, e439 := rt.Cond(ctx, t436)
		if e439 != nil {
			return rt.Value{}, e439
		}
		if t438 {
			// «к строке»
			t440, e441 := rt.BToString(ctx, dney)
			if e441 != nil {
				return rt.Value{}, e441
			}
			t442, e443 := rt.Concat(ctx, t440, rt.Text("д"))
			if e443 != nil {
				return rt.Value{}, e443
			}
			t444, e445 := rt.Concat(ctx, t442, rt.Text(" "))
			if e445 != nil {
				return rt.Value{}, e445
			}
			t446, e447 := rt.Concat(ctx, t444, vremya)
			if e447 != nil {
				return rt.Value{}, e447
			}
			return t446, nil
		} else {
			return vremya, nil
		}
	}
}

// Proshlo — функция flang «Прошло».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр sekundy — «секунды»: число.
// Результат — значение: строка.
func Proshlo(ctx *rt.Ctx, sekundy rt.Value) (rt.Value, error) {
	t448, e449 := rt.Lt(ctx, sekundy, rt.Number(3600.0))
	if e449 != nil {
		return rt.Value{}, e449
	}
	t450, e451 := rt.Cond(ctx, t448)
	if e451 != nil {
		return rt.Value{}, e451
	}
	if t450 {
		t452, e453 := DelenieNacelo(ctx, sekundy, rt.Number(60.0))
		if e453 != nil {
			return rt.Value{}, e453
		}
		// «к строке»
		t454, e455 := rt.BToString(ctx, t452)
		if e455 != nil {
			return rt.Value{}, e455
		}
		t456, e457 := rt.Concat(ctx, t454, rt.Text(" "))
		if e457 != nil {
			return rt.Value{}, e457
		}
		t458, e459 := rt.Concat(ctx, t456, rt.Text("мин"))
		if e459 != nil {
			return rt.Value{}, e459
		}
		return t458, nil
	} else {
		t460, e461 := rt.Lt(ctx, sekundy, rt.Number(172800.0))
		if e461 != nil {
			return rt.Value{}, e461
		}
		t462, e463 := rt.Cond(ctx, t460)
		if e463 != nil {
			return rt.Value{}, e463
		}
		if t462 {
			t464, e465 := DelenieNacelo(ctx, sekundy, rt.Number(3600.0))
			if e465 != nil {
				return rt.Value{}, e465
			}
			// «к строке»
			t466, e467 := rt.BToString(ctx, t464)
			if e467 != nil {
				return rt.Value{}, e467
			}
			t468, e469 := rt.Concat(ctx, t466, rt.Text(" "))
			if e469 != nil {
				return rt.Value{}, e469
			}
			t470, e471 := rt.Concat(ctx, t468, rt.Text("ч"))
			if e471 != nil {
				return rt.Value{}, e471
			}
			return t470, nil
		} else {
			t472, e473 := rt.Div(ctx, sekundy, rt.Number(3600.0))
			if e473 != nil {
				return rt.Value{}, e473
			}
			t474, e475 := rt.Div(ctx, t472, rt.Number(24.0))
			if e475 != nil {
				return rt.Value{}, e475
			}
			t476, e477 := CeloeKNulyu(ctx, t474)
			if e477 != nil {
				return rt.Value{}, e477
			}
			// «к строке»
			t478, e479 := rt.BToString(ctx, t476)
			if e479 != nil {
				return rt.Value{}, e479
			}
			t480, e481 := rt.Concat(ctx, t478, rt.Text(" "))
			if e481 != nil {
				return rt.Value{}, e481
			}
			t482, e483 := rt.Concat(ctx, t480, rt.Text("дн"))
			if e483 != nil {
				return rt.Value{}, e483
			}
			return t482, nil
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
	t484, e485 := rt.Lt(ctx, millisekundy, rt.Number(1000.0))
	if e485 != nil {
		return rt.Value{}, e485
	}
	t486, e487 := rt.Cond(ctx, t484)
	if e487 != nil {
		return rt.Value{}, e487
	}
	if t486 {
		t488, e489 := CeloeKNulyu(ctx, millisekundy)
		if e489 != nil {
			return rt.Value{}, e489
		}
		// «к строке»
		t490, e491 := rt.BToString(ctx, t488)
		if e491 != nil {
			return rt.Value{}, e491
		}
		t492, e493 := rt.Concat(ctx, t490, rt.Text(" мс"))
		if e493 != nil {
			return rt.Value{}, e493
		}
		return t492, nil
	} else {
		t494, e495 := rt.Div(ctx, millisekundy, rt.Number(1000.0))
		if e495 != nil {
			return rt.Value{}, e495
		}
		t496, e497 := Drobyu(ctx, t494, rt.Number(1.0), razdelitel)
		if e497 != nil {
			return rt.Value{}, e497
		}
		t498, e499 := rt.Concat(ctx, t496, rt.Text(" с"))
		if e499 != nil {
			return rt.Value{}, e499
		}
		return t498, nil
	}
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
	t500, e501 := rt.FieldGet(ctx, hod, "осталось")
	if e501 != nil {
		return rt.Value{}, e501
	}
	t502, e503 := rt.Mod(ctx, t500, rt.Number(3.0))
	if e503 != nil {
		return rt.Value{}, e503
	}
	t504, e505 := rt.Cond(ctx, rt.Flag(rt.Equal(t502, rt.Number(0.0))))
	if e505 != nil {
		return rt.Value{}, e505
	}
	var t506 rt.Value
	if t504 {
		t507, e508 := rt.FieldGet(ctx, hod, "собрано")
		if e508 != nil {
			return rt.Value{}, e508
		}
		// «длина»
		t509, e510 := rt.BLength(ctx, t507)
		if e510 != nil {
			return rt.Value{}, e510
		}
		t511, e512 := rt.Gt(ctx, t509, rt.Number(0.0))
		if e512 != nil {
			return rt.Value{}, e512
		}
		t506 = t511
	} else {
		t506 = rt.Flag(false)
	}
	t513, e514 := rt.Cond(ctx, t506)
	if e514 != nil {
		return rt.Value{}, e514
	}
	var t515 rt.Value
	if t513 {
		t516, e517 := rt.FieldGet(ctx, hod, "собрано")
		if e517 != nil {
			return rt.Value{}, e517
		}
		t518, e519 := rt.FieldGet(ctx, hod, "разделитель")
		if e519 != nil {
			return rt.Value{}, e519
		}
		t520, e521 := rt.Concat(ctx, t516, t518)
		if e521 != nil {
			return rt.Value{}, e521
		}
		t522, e523 := rt.Concat(ctx, t520, cifra)
		if e523 != nil {
			return rt.Value{}, e523
		}
		t524, e525 := rt.FieldGet(ctx, hod, "осталось")
		if e525 != nil {
			return rt.Value{}, e525
		}
		t526, e527 := rt.Sub(ctx, t524, rt.Number(1.0))
		if e527 != nil {
			return rt.Value{}, e527
		}
		t528, e529 := rt.FieldGet(ctx, hod, "разделитель")
		if e529 != nil {
			return rt.Value{}, e529
		}
		t530 := make([]rt.Field, 3)
		t530[0] = rt.Field{Name: "собрано", Value: t522}
		t530[1] = rt.Field{Name: "осталось", Value: t526}
		t530[2] = rt.Field{Name: "разделитель", Value: t528}
		t515 = rt.Record(t530)
	} else {
		t531, e532 := rt.FieldGet(ctx, hod, "собрано")
		if e532 != nil {
			return rt.Value{}, e532
		}
		t533, e534 := rt.Concat(ctx, t531, cifra)
		if e534 != nil {
			return rt.Value{}, e534
		}
		t535, e536 := rt.FieldGet(ctx, hod, "осталось")
		if e536 != nil {
			return rt.Value{}, e536
		}
		t537, e538 := rt.Sub(ctx, t535, rt.Number(1.0))
		if e538 != nil {
			return rt.Value{}, e538
		}
		t539, e540 := rt.FieldGet(ctx, hod, "разделитель")
		if e540 != nil {
			return rt.Value{}, e540
		}
		t541 := make([]rt.Field, 3)
		t541[0] = rt.Field{Name: "собрано", Value: t533}
		t541[1] = rt.Field{Name: "осталось", Value: t537}
		t541[2] = rt.Field{Name: "разделитель", Value: t539}
		t515 = rt.Record(t541)
	}
	t542 := t515
	t543, e544 := rt.FieldGet(ctx, t542, "осталось")
	if e544 != nil {
		return rt.Value{}, e544
	}
	t545, e546 := rt.FieldGet(ctx, hod, "осталось")
	if e546 != nil {
		return rt.Value{}, e546
	}
	t547, e548 := rt.Lt(ctx, t543, t545)
	if e548 != nil {
		return rt.Value{}, e548
	}
	// постусловие «осталось убывает»
	t549, e550 := rt.Post(ctx, t547, "осталось убывает", "Шаг разряда")
	if e550 != nil {
		return rt.Value{}, e550
	}
	if !t549 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «осталось убывает» функции «Шаг разряда»")
	}
	t551, e552 := rt.FieldGet(ctx, t542, "разделитель")
	if e552 != nil {
		return rt.Value{}, e552
	}
	t553, e554 := rt.FieldGet(ctx, hod, "разделитель")
	if e554 != nil {
		return rt.Value{}, e554
	}
	// постусловие «разделитель не меняется»
	t555, e556 := rt.Post(ctx, rt.Flag(rt.Equal(t551, t553)), "разделитель не меняется", "Шаг разряда")
	if e556 != nil {
		return rt.Value{}, e556
	}
	if !t555 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «разделитель не меняется» функции «Шаг разряда»")
	}
	return t542, nil
}

// Razryadami — функция flang «Разрядами».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func Razryadami(ctx *rt.Ctx, znachenie rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t557, e558 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e558 != nil {
		return rt.Value{}, e558
	}
	t559, e560 := rt.Cond(ctx, t557)
	if e560 != nil {
		return rt.Value{}, e560
	}
	var t561 rt.Value
	if t559 {
		t561 = rt.Text("-")
	} else {
		t561 = rt.Text("")
	}
	// пусть «знак»
	znak := t561
	t562, e563 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e563 != nil {
		return rt.Value{}, e563
	}
	t564, e565 := rt.Cond(ctx, t562)
	if e565 != nil {
		return rt.Value{}, e565
	}
	var t566 rt.Value
	if t564 {
		t567, e568 := rt.Sub(ctx, rt.Number(0.0), znachenie)
		if e568 != nil {
			return rt.Value{}, e568
		}
		t566 = t567
	} else {
		t566 = znachenie
	}
	t569, e570 := CeloeKNulyu(ctx, t566)
	if e570 != nil {
		return rt.Value{}, e570
	}
	// «к строке»
	t571, e572 := rt.BToString(ctx, t569)
	if e572 != nil {
		return rt.Value{}, e572
	}
	// пусть «цифры»
	cifry := t571
	// «символы»
	t573, e574 := rt.BCharacters(ctx, cifry)
	if e574 != nil {
		return rt.Value{}, e574
	}
	t575, e576 := rt.RequireList(ctx, t573, "свёртка")
	if e576 != nil {
		return rt.Value{}, e576
	}
	// «длина»
	t577, e578 := rt.BLength(ctx, cifry)
	if e578 != nil {
		return rt.Value{}, e578
	}
	t579 := make([]rt.Field, 3)
	t579[0] = rt.Field{Name: "собрано", Value: rt.Text("")}
	t579[1] = rt.Field{Name: "осталось", Value: t577}
	t579[2] = rt.Field{Name: "разделитель", Value: razdelitel}
	// «ход»
	hod := rt.Record(t579)
	for t580 := range t575 {
		// «цифра»
		cifra := t575[t580]
		t581, e582 := ShagRazryada(ctx, hod, cifra)
		if e582 != nil {
			return rt.Value{}, e582
		}
		hod = t581
	}
	t583, e584 := rt.FieldGet(ctx, hod, "собрано")
	if e584 != nil {
		return rt.Value{}, e584
	}
	t585, e586 := rt.Concat(ctx, znak, t583)
	if e586 != nil {
		return rt.Value{}, e586
	}
	return t585, nil
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
		t587, e588 := rt.BCharCodeProven(ctx, golova)
		if e588 != nil {
			return rt.Value{}, e588
		}
		t589, e590 := rt.Cond(ctx, rt.Flag(rt.Equal(t587, rt.Number(32.0))))
		if e590 != nil {
			return rt.Value{}, e590
		}
		var t591 rt.Value
		if t589 {
			t591 = rt.Flag(true)
		} else {
			// «код символа»
			t592, e593 := rt.BCharCodeProven(ctx, golova)
			if e593 != nil {
				return rt.Value{}, e593
			}
			t591 = rt.Flag(rt.Equal(t592, rt.Number(9.0)))
		}
		t594, e595 := rt.Cond(ctx, t591)
		if e595 != nil {
			return rt.Value{}, e595
		}
		var t596 rt.Value
		if t594 {
			t596 = rt.Flag(true)
		} else {
			// «код символа»
			t597, e598 := rt.BCharCodeProven(ctx, golova)
			if e598 != nil {
				return rt.Value{}, e598
			}
			t596 = rt.Flag(rt.Equal(t597, rt.Number(10.0)))
		}
		t599, e600 := rt.Cond(ctx, t596)
		if e600 != nil {
			return rt.Value{}, e600
		}
		var t601 rt.Value
		if t599 {
			t601 = rt.Flag(true)
		} else {
			// «код символа»
			t602, e603 := rt.BCharCodeProven(ctx, golova)
			if e603 != nil {
				return rt.Value{}, e603
			}
			t601 = rt.Flag(rt.Equal(t602, rt.Number(11.0)))
		}
		t604, e605 := rt.Cond(ctx, t601)
		if e605 != nil {
			return rt.Value{}, e605
		}
		var t606 rt.Value
		if t604 {
			t606 = rt.Flag(true)
		} else {
			// «код символа»
			t607, e608 := rt.BCharCodeProven(ctx, golova)
			if e608 != nil {
				return rt.Value{}, e608
			}
			t606 = rt.Flag(rt.Equal(t607, rt.Number(12.0)))
		}
		t609, e610 := rt.Cond(ctx, t606)
		if e610 != nil {
			return rt.Value{}, e610
		}
		var t611 rt.Value
		if t609 {
			t611 = rt.Flag(true)
		} else {
			// «код символа»
			t612, e613 := rt.BCharCodeProven(ctx, golova)
			if e613 != nil {
				return rt.Value{}, e613
			}
			t611 = rt.Flag(rt.Equal(t612, rt.Number(13.0)))
		}
		t614, e615 := rt.Cond(ctx, t611)
		if e615 != nil {
			return rt.Value{}, e615
		}
		var t616 rt.Value
		if t614 {
			t616 = rt.Flag(true)
		} else {
			// «код символа»
			t617, e618 := rt.BCharCodeProven(ctx, golova)
			if e618 != nil {
				return rt.Value{}, e618
			}
			t616 = rt.Flag(rt.Equal(t617, rt.Number(133.0)))
		}
		t619, e620 := rt.Cond(ctx, t616)
		if e620 != nil {
			return rt.Value{}, e620
		}
		if t619 {
			return rt.Flag(true), nil
		} else {
			// «код символа»
			t621, e622 := rt.BCharCodeProven(ctx, golova)
			if e622 != nil {
				return rt.Value{}, e622
			}
			return rt.Flag(rt.Equal(t621, rt.Number(160.0))), nil
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
	t623, e624 := rt.BCharacters(ctx, tekst)
	if e624 != nil {
		return rt.Value{}, e624
	}
	t625, e626 := rt.RequireList(ctx, t623, "отфильтровать")
	if e626 != nil {
		return rt.Value{}, e626
	}
	t627 := make([]rt.Value, 0, len(t625))
	for t628 := range t625 {
		// «знак»
		znak := t625[t628]
		t629, e630 := Probelnyy(ctx, znak)
		if e630 != nil {
			return rt.Value{}, e630
		}
		t631, e632 := rt.Cond(ctx, t629)
		if e632 != nil {
			return rt.Value{}, e632
		}
		var t633 rt.Value
		if t631 {
			t633 = rt.Flag(false)
		} else {
			t633 = rt.Flag(true)
		}
		t634, e635 := rt.Keep(ctx, t633)
		if e635 != nil {
			return rt.Value{}, e635
		}
		if t634 {
			t627 = append(t627, znak)
		}
	}
	// «длина»
	t636, e637 := rt.BLength(ctx, rt.List(t627))
	if e637 != nil {
		return rt.Value{}, e637
	}
	t638, e639 := rt.Cond(ctx, rt.Flag(rt.Equal(t636, rt.Number(0.0))))
	if e639 != nil {
		return rt.Value{}, e639
	}
	if t638 {
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
	t640, e641 := rt.BLength(ctx, tekst)
	if e641 != nil {
		return rt.Value{}, e641
	}
	t642, e643 := rt.Lte(ctx, t640, skolko)
	if e643 != nil {
		return rt.Value{}, e643
	}
	t644, e645 := rt.Cond(ctx, t642)
	if e645 != nil {
		return rt.Value{}, e645
	}
	var t646 rt.Value
	if t644 {
		t646 = tekst
	} else {
		t647, e648 := rt.Lte(ctx, skolko, rt.Number(0.0))
		if e648 != nil {
			return rt.Value{}, e648
		}
		t649, e650 := rt.Cond(ctx, t647)
		if e650 != nil {
			return rt.Value{}, e650
		}
		var t651 rt.Value
		if t649 {
			t651 = rt.Text("")
		} else {
			t652, e653 := rt.Lte(ctx, skolko, rt.Number(1.0))
			if e653 != nil {
				return rt.Value{}, e653
			}
			t654, e655 := rt.Cond(ctx, t652)
			if e655 != nil {
				return rt.Value{}, e655
			}
			var t656 rt.Value
			if t654 {
				// «подстрока»
				t657, e658 := rt.BSubstring(ctx, tekst, rt.Number(1.0), rt.Number(1.0))
				if e658 != nil {
					return rt.Value{}, e658
				}
				t656 = t657
			} else {
				t659, e660 := rt.Sub(ctx, skolko, rt.Number(1.0))
				if e660 != nil {
					return rt.Value{}, e660
				}
				// «подстрока»
				t661, e662 := rt.BSubstring(ctx, tekst, rt.Number(1.0), t659)
				if e662 != nil {
					return rt.Value{}, e662
				}
				t663, e664 := rt.Concat(ctx, t661, rt.Text("…"))
				if e664 != nil {
					return rt.Value{}, e664
				}
				t656 = t663
			}
			t651 = t656
		}
		t646 = t651
	}
	t665 := t646
	t666, e667 := rt.Gte(ctx, skolko, rt.Number(0.0))
	if e667 != nil {
		return rt.Value{}, e667
	}
	t668, e669 := rt.Cond(ctx, t666)
	if e669 != nil {
		return rt.Value{}, e669
	}
	var t670 rt.Value
	if t668 {
		// «длина»
		t671, e672 := rt.BLength(ctx, t665)
		if e672 != nil {
			return rt.Value{}, e672
		}
		// «длина»
		t673, e674 := rt.BLength(ctx, tekst)
		if e674 != nil {
			return rt.Value{}, e674
		}
		t675, e676 := rt.Lte(ctx, t673, skolko)
		if e676 != nil {
			return rt.Value{}, e676
		}
		t677, e678 := rt.Cond(ctx, t675)
		if e678 != nil {
			return rt.Value{}, e678
		}
		var t679 rt.Value
		if t677 {
			// «длина»
			t680, e681 := rt.BLength(ctx, tekst)
			if e681 != nil {
				return rt.Value{}, e681
			}
			t679 = t680
		} else {
			t679 = skolko
		}
		t682, e683 := rt.Lte(ctx, t671, t679)
		if e683 != nil {
			return rt.Value{}, e683
		}
		t670 = t682
	} else {
		t670 = rt.Flag(true)
	}
	// постусловие «обрезанное не длиннее заказа»
	t684, e685 := rt.Post(ctx, t670, "обрезанное не длиннее заказа", "Обрезать")
	if e685 != nil {
		return rt.Value{}, e685
	}
	if !t684 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «обрезанное не длиннее заказа» функции «Обрезать»")
	}
	return t665, nil
}

// MeraUbyvaet — функция flang «мера убывает».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр shag — «шаг»: число.
// Параметр mera — «мера»: число.
// Результат — значение: число.
func MeraUbyvaet(ctx *rt.Ctx, shag rt.Value, mera rt.Value) (rt.Value, error) {
	t686 := shag
	t687, e688 := rt.Lt(ctx, t686, mera)
	if e688 != nil {
		return rt.Value{}, e688
	}
	// постусловие «мера убывает»
	t689, e690 := rt.Post(ctx, t687, "мера убывает", "мера убывает")
	if e690 != nil {
		return rt.Value{}, e690
	}
	if !t689 {
		return rt.Value{}, rt.Fail("FLANG_MEASURE", "%s", "тотальная функция «Повторить»: мера не убыла — аргумент 2 вызова «Повторить» не стал меньше параметра «сколько». Завершение доказано убыванием этой меры, а числа flang — IEEE-754 double: при большом |«сколько»| постоянный шаг не меняет значение, и спуск не идёт. Отказ здесь честнее зацикливания")
	}
	return t686, nil
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
