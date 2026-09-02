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
	var t238 rt.Value
	if t236 {
		t238 = rt.Number(1024.0)
	} else {
		t239, e240 := rt.Gte(ctx, vitki, rt.Number(1024.0))
		if e240 != nil {
			return rt.Value{}, e240
		}
		t241, e242 := rt.Cond(ctx, t239)
		if e242 != nil {
			return rt.Value{}, e242
		}
		var t243 rt.Value
		if t241 {
			t244, e245 := DelenieNacelo(ctx, vitki, rt.Number(1024.0))
			if e245 != nil {
				return rt.Value{}, e245
			}
			t246, e247 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e247 != nil {
				return rt.Value{}, e247
			}
			t248, e249 := DelitelBayt(ctx, t244, t246)
			if e249 != nil {
				return rt.Value{}, e249
			}
			t250, e251 := rt.Mul(ctx, rt.Number(1024.0), t248)
			if e251 != nil {
				return rt.Value{}, e251
			}
			t243 = t250
		} else {
			t243 = rt.Number(1024.0)
		}
		t238 = t243
	}
	t252 := t238
	t253, e254 := rt.Gte(ctx, t252, rt.Number(1024.0))
	if e254 != nil {
		return rt.Value{}, e254
	}
	// постусловие «делитель не меньше тысячи двадцати четырёх»
	t255, e256 := rt.Post(ctx, t253, "делитель не меньше тысячи двадцати четырёх", "Делитель байт")
	if e256 != nil {
		return rt.Value{}, e256
	}
	if !t255 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «делитель не меньше тысячи двадцати четырёх» функции «Делитель байт»")
	}
	return t252, nil
}

// EdinicaBayt — функция flang «Единица байт».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stupen — «ступень»: число.
// Результат — значение: строка.
func EdinicaBayt(ctx *rt.Ctx, stupen rt.Value) (rt.Value, error) {
	t257, e258 := rt.Lte(ctx, stupen, rt.Number(0.0))
	if e258 != nil {
		return rt.Value{}, e258
	}
	t259, e260 := rt.Cond(ctx, t257)
	if e260 != nil {
		return rt.Value{}, e260
	}
	if t259 {
		return rt.Text("КиБ"), nil
	} else {
		t261, e262 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e262 != nil {
			return rt.Value{}, e262
		}
		t263, e264 := rt.Cond(ctx, t261)
		if e264 != nil {
			return rt.Value{}, e264
		}
		if t263 {
			return rt.Text("ПиБ"), nil
		} else {
			t265, e266 := rt.Add(ctx, stupen, rt.Number(1.0))
			if e266 != nil {
				return rt.Value{}, e266
			}
			t267 := make([]rt.Value, 5)
			t267[0] = rt.Text("КиБ")
			t267[1] = rt.Text("МиБ")
			t267[2] = rt.Text("ГиБ")
			t267[3] = rt.Text("ТиБ")
			t267[4] = rt.Text("ПиБ")
			// «элемент»
			t268, e269 := rt.BElement(ctx, t265, rt.List(t267))
			if e269 != nil {
				return rt.Value{}, e269
			}
			return t268, nil
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
	t270, e271 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e271 != nil {
		return rt.Value{}, e271
	}
	t272, e273 := rt.Cond(ctx, t270)
	if e273 != nil {
		return rt.Value{}, e273
	}
	var t274 rt.Value
	if t272 {
		t274 = rt.Text("-")
	} else {
		t274 = rt.Text("")
	}
	// пусть «знак»
	znak := t274
	t275, e276 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e276 != nil {
		return rt.Value{}, e276
	}
	t277, e278 := rt.Cond(ctx, t275)
	if e278 != nil {
		return rt.Value{}, e278
	}
	var t279 rt.Value
	if t277 {
		t280, e281 := rt.Sub(ctx, rt.Number(0.0), skolko)
		if e281 != nil {
			return rt.Value{}, e281
		}
		t279 = t280
	} else {
		t279 = skolko
	}
	t282, e283 := CeloeKNulyu(ctx, t279)
	if e283 != nil {
		return rt.Value{}, e283
	}
	// пусть «величина»
	velichina := t282
	t284, e285 := rt.Lt(ctx, velichina, rt.Number(1024.0))
	if e285 != nil {
		return rt.Value{}, e285
	}
	t286, e287 := rt.Cond(ctx, t284)
	if e287 != nil {
		return rt.Value{}, e287
	}
	if t286 {
		// «к строке»
		t288, e289 := rt.BToString(ctx, velichina)
		if e289 != nil {
			return rt.Value{}, e289
		}
		t290, e291 := rt.Concat(ctx, znak, t288)
		if e291 != nil {
			return rt.Value{}, e291
		}
		t292, e293 := rt.Concat(ctx, t290, rt.Text(" Б"))
		if e293 != nil {
			return rt.Value{}, e293
		}
		return t292, nil
	} else {
		t294, e295 := DelenieNacelo(ctx, velichina, rt.Number(1024.0))
		if e295 != nil {
			return rt.Value{}, e295
		}
		// пусть «витки»
		vitki := t294
		t296, e297 := DelitelBayt(ctx, vitki, rt.Number(4.0))
		if e297 != nil {
			return rt.Value{}, e297
		}
		t298, e299 := rt.Div(ctx, velichina, t296)
		if e299 != nil {
			return rt.Value{}, e299
		}
		t300, e301 := Drobyu(ctx, t298, rt.Number(1.0), razdelitel)
		if e301 != nil {
			return rt.Value{}, e301
		}
		t302, e303 := rt.Concat(ctx, znak, t300)
		if e303 != nil {
			return rt.Value{}, e303
		}
		t304, e305 := rt.Concat(ctx, t302, rt.Text(" "))
		if e305 != nil {
			return rt.Value{}, e305
		}
		t306, e307 := StupenBayt(ctx, vitki, rt.Number(4.0))
		if e307 != nil {
			return rt.Value{}, e307
		}
		t308, e309 := EdinicaBayt(ctx, t306)
		if e309 != nil {
			return rt.Value{}, e309
		}
		t310, e311 := rt.Concat(ctx, t304, t308)
		if e311 != nil {
			return rt.Value{}, e311
		}
		return t310, nil
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
	t312, e313 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e313 != nil {
		return rt.Value{}, e313
	}
	t314, e315 := rt.Cond(ctx, t312)
	if e315 != nil {
		return rt.Value{}, e315
	}
	if t314 {
		return rt.Number(0.0), nil
	} else {
		t316, e317 := rt.Gte(ctx, vitki, rt.Number(1000.0))
		if e317 != nil {
			return rt.Value{}, e317
		}
		t318, e319 := rt.Cond(ctx, t316)
		if e319 != nil {
			return rt.Value{}, e319
		}
		if t318 {
			t320, e321 := DelenieNacelo(ctx, vitki, rt.Number(1000.0))
			if e321 != nil {
				return rt.Value{}, e321
			}
			t322, e323 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e323 != nil {
				return rt.Value{}, e323
			}
			t324, e325 := StupenTysyach(ctx, t320, t322)
			if e325 != nil {
				return rt.Value{}, e325
			}
			t326, e327 := rt.Add(ctx, rt.Number(1.0), t324)
			if e327 != nil {
				return rt.Value{}, e327
			}
			return t326, nil
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
	t328, e329 := rt.Lte(ctx, zapas, rt.Number(0.0))
	if e329 != nil {
		return rt.Value{}, e329
	}
	t330, e331 := rt.Cond(ctx, t328)
	if e331 != nil {
		return rt.Value{}, e331
	}
	var t332 rt.Value
	if t330 {
		t332 = rt.Number(1000.0)
	} else {
		t333, e334 := rt.Gte(ctx, vitki, rt.Number(1000.0))
		if e334 != nil {
			return rt.Value{}, e334
		}
		t335, e336 := rt.Cond(ctx, t333)
		if e336 != nil {
			return rt.Value{}, e336
		}
		var t337 rt.Value
		if t335 {
			t338, e339 := DelenieNacelo(ctx, vitki, rt.Number(1000.0))
			if e339 != nil {
				return rt.Value{}, e339
			}
			t340, e341 := rt.Sub(ctx, zapas, rt.Number(1.0))
			if e341 != nil {
				return rt.Value{}, e341
			}
			t342, e343 := DelitelTysyach(ctx, t338, t340)
			if e343 != nil {
				return rt.Value{}, e343
			}
			t344, e345 := rt.Mul(ctx, rt.Number(1000.0), t342)
			if e345 != nil {
				return rt.Value{}, e345
			}
			t337 = t344
		} else {
			t337 = rt.Number(1000.0)
		}
		t332 = t337
	}
	t346 := t332
	t347, e348 := rt.Gte(ctx, t346, rt.Number(1000.0))
	if e348 != nil {
		return rt.Value{}, e348
	}
	// постусловие «делитель не меньше тысячи»
	t349, e350 := rt.Post(ctx, t347, "делитель не меньше тысячи", "Делитель тысяч")
	if e350 != nil {
		return rt.Value{}, e350
	}
	if !t349 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «делитель не меньше тысячи» функции «Делитель тысяч»")
	}
	return t346, nil
}

// EdinicaTysyach — функция flang «Единица тысяч».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр stupen — «ступень»: число.
// Результат — значение: строка.
func EdinicaTysyach(ctx *rt.Ctx, stupen rt.Value) (rt.Value, error) {
	t351, e352 := rt.Lte(ctx, stupen, rt.Number(0.0))
	if e352 != nil {
		return rt.Value{}, e352
	}
	t353, e354 := rt.Cond(ctx, t351)
	if e354 != nil {
		return rt.Value{}, e354
	}
	if t353 {
		return rt.Text("КБ"), nil
	} else {
		t355, e356 := rt.Gte(ctx, stupen, rt.Number(4.0))
		if e356 != nil {
			return rt.Value{}, e356
		}
		t357, e358 := rt.Cond(ctx, t355)
		if e358 != nil {
			return rt.Value{}, e358
		}
		if t357 {
			return rt.Text("ПБ"), nil
		} else {
			t359, e360 := rt.Add(ctx, stupen, rt.Number(1.0))
			if e360 != nil {
				return rt.Value{}, e360
			}
			t361 := make([]rt.Value, 5)
			t361[0] = rt.Text("КБ")
			t361[1] = rt.Text("МБ")
			t361[2] = rt.Text("ГБ")
			t361[3] = rt.Text("ТБ")
			t361[4] = rt.Text("ПБ")
			// «элемент»
			t362, e363 := rt.BElement(ctx, t359, rt.List(t361))
			if e363 != nil {
				return rt.Value{}, e363
			}
			return t362, nil
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
	t364, e365 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e365 != nil {
		return rt.Value{}, e365
	}
	t366, e367 := rt.Cond(ctx, t364)
	if e367 != nil {
		return rt.Value{}, e367
	}
	var t368 rt.Value
	if t366 {
		t368 = rt.Text("-")
	} else {
		t368 = rt.Text("")
	}
	// пусть «знак»
	znak := t368
	t369, e370 := rt.Lt(ctx, skolko, rt.Number(0.0))
	if e370 != nil {
		return rt.Value{}, e370
	}
	t371, e372 := rt.Cond(ctx, t369)
	if e372 != nil {
		return rt.Value{}, e372
	}
	var t373 rt.Value
	if t371 {
		t374, e375 := rt.Sub(ctx, rt.Number(0.0), skolko)
		if e375 != nil {
			return rt.Value{}, e375
		}
		t373 = t374
	} else {
		t373 = skolko
	}
	t376, e377 := CeloeKNulyu(ctx, t373)
	if e377 != nil {
		return rt.Value{}, e377
	}
	// пусть «величина»
	velichina := t376
	t378, e379 := rt.Lt(ctx, velichina, rt.Number(1000.0))
	if e379 != nil {
		return rt.Value{}, e379
	}
	t380, e381 := rt.Cond(ctx, t378)
	if e381 != nil {
		return rt.Value{}, e381
	}
	if t380 {
		// «к строке»
		t382, e383 := rt.BToString(ctx, velichina)
		if e383 != nil {
			return rt.Value{}, e383
		}
		t384, e385 := rt.Concat(ctx, znak, t382)
		if e385 != nil {
			return rt.Value{}, e385
		}
		t386, e387 := rt.Concat(ctx, t384, rt.Text(" Б"))
		if e387 != nil {
			return rt.Value{}, e387
		}
		return t386, nil
	} else {
		t388, e389 := DelenieNacelo(ctx, velichina, rt.Number(1000.0))
		if e389 != nil {
			return rt.Value{}, e389
		}
		// пусть «витки»
		vitki := t388
		t390, e391 := DelitelTysyach(ctx, vitki, rt.Number(4.0))
		if e391 != nil {
			return rt.Value{}, e391
		}
		t392, e393 := rt.Div(ctx, velichina, t390)
		if e393 != nil {
			return rt.Value{}, e393
		}
		t394, e395 := Drobyu(ctx, t392, rt.Number(1.0), razdelitel)
		if e395 != nil {
			return rt.Value{}, e395
		}
		t396, e397 := rt.Concat(ctx, znak, t394)
		if e397 != nil {
			return rt.Value{}, e397
		}
		t398, e399 := rt.Concat(ctx, t396, rt.Text(" "))
		if e399 != nil {
			return rt.Value{}, e399
		}
		t400, e401 := StupenTysyach(ctx, vitki, rt.Number(4.0))
		if e401 != nil {
			return rt.Value{}, e401
		}
		t402, e403 := EdinicaTysyach(ctx, t400)
		if e403 != nil {
			return rt.Value{}, e403
		}
		t404, e405 := rt.Concat(ctx, t398, t402)
		if e405 != nil {
			return rt.Value{}, e405
		}
		return t404, nil
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
	t406, e407 := rt.Cond(ctx, rt.Flag(rt.Equal(celoe, rt.Number(0.0))))
	if e407 != nil {
		return rt.Value{}, e407
	}
	if t406 {
		return rt.Number(0.0), nil
	} else {
		t408, e409 := rt.Div(ctx, chast, celoe)
		if e409 != nil {
			return rt.Value{}, e409
		}
		return t408, nil
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
	t410, e411 := rt.Cond(ctx, rt.Flag(rt.Equal(celoe, rt.Number(0.0))))
	if e411 != nil {
		return rt.Value{}, e411
	}
	if t410 {
		return rt.Number(0.0), nil
	} else {
		t412, e413 := rt.Mul(ctx, rt.Number(100.0), chast)
		if e413 != nil {
			return rt.Value{}, e413
		}
		t414, e415 := rt.Div(ctx, t412, celoe)
		if e415 != nil {
			return rt.Value{}, e415
		}
		return t414, nil
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
	t416, e417 := rt.Mul(ctx, dolya, rt.Number(100.0))
	if e417 != nil {
		return rt.Value{}, e417
	}
	t418, e419 := Drobyu(ctx, t416, rt.Number(1.0), razdelitel)
	if e419 != nil {
		return rt.Value{}, e419
	}
	t420, e421 := rt.Concat(ctx, t418, rt.Text("%"))
	if e421 != nil {
		return rt.Value{}, e421
	}
	return t420, nil
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
	t422, e423 := rt.Lte(ctx, sekundy, rt.Number(0.0))
	if e423 != nil {
		return rt.Value{}, e423
	}
	t424, e425 := rt.Cond(ctx, t422)
	if e425 != nil {
		return rt.Value{}, e425
	}
	if t424 {
		return rt.Text(""), nil
	} else {
		t426, e427 := CeloeKNulyu(ctx, sekundy)
		if e427 != nil {
			return rt.Value{}, e427
		}
		// пусть «всего»
		vsego := t426
		t428, e429 := DelenieNacelo(ctx, vsego, rt.Number(3600.0))
		if e429 != nil {
			return rt.Value{}, e429
		}
		// пусть «часов всего»
		chasovVsego := t428
		t430, e431 := DelenieNacelo(ctx, chasovVsego, rt.Number(24.0))
		if e431 != nil {
			return rt.Value{}, e431
		}
		// пусть «дней»
		dney := t430
		t432, e433 := rt.Mod(ctx, chasovVsego, rt.Number(24.0))
		if e433 != nil {
			return rt.Value{}, e433
		}
		// пусть «часов»
		chasov := t432
		t434, e435 := DelenieNacelo(ctx, vsego, rt.Number(60.0))
		if e435 != nil {
			return rt.Value{}, e435
		}
		t436, e437 := rt.Mod(ctx, t434, rt.Number(60.0))
		if e437 != nil {
			return rt.Value{}, e437
		}
		// пусть «минут»
		minut := t436
		// «к строке»
		t438, e439 := rt.BToString(ctx, chasov)
		if e439 != nil {
			return rt.Value{}, e439
		}
		t440, e441 := SlevaNulyami(ctx, t438, rt.Number(2.0))
		if e441 != nil {
			return rt.Value{}, e441
		}
		t442, e443 := rt.Concat(ctx, t440, rt.Text(":"))
		if e443 != nil {
			return rt.Value{}, e443
		}
		// «к строке»
		t444, e445 := rt.BToString(ctx, minut)
		if e445 != nil {
			return rt.Value{}, e445
		}
		t446, e447 := SlevaNulyami(ctx, t444, rt.Number(2.0))
		if e447 != nil {
			return rt.Value{}, e447
		}
		t448, e449 := rt.Concat(ctx, t442, t446)
		if e449 != nil {
			return rt.Value{}, e449
		}
		// пусть «время»
		vremya := t448
		t450, e451 := rt.Gt(ctx, dney, rt.Number(0.0))
		if e451 != nil {
			return rt.Value{}, e451
		}
		t452, e453 := rt.Cond(ctx, t450)
		if e453 != nil {
			return rt.Value{}, e453
		}
		if t452 {
			// «к строке»
			t454, e455 := rt.BToString(ctx, dney)
			if e455 != nil {
				return rt.Value{}, e455
			}
			t456, e457 := rt.Concat(ctx, t454, rt.Text("д"))
			if e457 != nil {
				return rt.Value{}, e457
			}
			t458, e459 := rt.Concat(ctx, t456, rt.Text(" "))
			if e459 != nil {
				return rt.Value{}, e459
			}
			t460, e461 := rt.Concat(ctx, t458, vremya)
			if e461 != nil {
				return rt.Value{}, e461
			}
			return t460, nil
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
	t462, e463 := rt.Lt(ctx, sekundy, rt.Number(3600.0))
	if e463 != nil {
		return rt.Value{}, e463
	}
	t464, e465 := rt.Cond(ctx, t462)
	if e465 != nil {
		return rt.Value{}, e465
	}
	if t464 {
		t466, e467 := DelenieNacelo(ctx, sekundy, rt.Number(60.0))
		if e467 != nil {
			return rt.Value{}, e467
		}
		// «к строке»
		t468, e469 := rt.BToString(ctx, t466)
		if e469 != nil {
			return rt.Value{}, e469
		}
		t470, e471 := rt.Concat(ctx, t468, rt.Text(" "))
		if e471 != nil {
			return rt.Value{}, e471
		}
		t472, e473 := rt.Concat(ctx, t470, rt.Text("мин"))
		if e473 != nil {
			return rt.Value{}, e473
		}
		return t472, nil
	} else {
		t474, e475 := rt.Lt(ctx, sekundy, rt.Number(172800.0))
		if e475 != nil {
			return rt.Value{}, e475
		}
		t476, e477 := rt.Cond(ctx, t474)
		if e477 != nil {
			return rt.Value{}, e477
		}
		if t476 {
			t478, e479 := DelenieNacelo(ctx, sekundy, rt.Number(3600.0))
			if e479 != nil {
				return rt.Value{}, e479
			}
			// «к строке»
			t480, e481 := rt.BToString(ctx, t478)
			if e481 != nil {
				return rt.Value{}, e481
			}
			t482, e483 := rt.Concat(ctx, t480, rt.Text(" "))
			if e483 != nil {
				return rt.Value{}, e483
			}
			t484, e485 := rt.Concat(ctx, t482, rt.Text("ч"))
			if e485 != nil {
				return rt.Value{}, e485
			}
			return t484, nil
		} else {
			t486, e487 := rt.Div(ctx, sekundy, rt.Number(3600.0))
			if e487 != nil {
				return rt.Value{}, e487
			}
			t488, e489 := rt.Div(ctx, t486, rt.Number(24.0))
			if e489 != nil {
				return rt.Value{}, e489
			}
			t490, e491 := CeloeKNulyu(ctx, t488)
			if e491 != nil {
				return rt.Value{}, e491
			}
			// «к строке»
			t492, e493 := rt.BToString(ctx, t490)
			if e493 != nil {
				return rt.Value{}, e493
			}
			t494, e495 := rt.Concat(ctx, t492, rt.Text(" "))
			if e495 != nil {
				return rt.Value{}, e495
			}
			t496, e497 := rt.Concat(ctx, t494, rt.Text("дн"))
			if e497 != nil {
				return rt.Value{}, e497
			}
			return t496, nil
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
	t498, e499 := rt.Lt(ctx, millisekundy, rt.Number(1000.0))
	if e499 != nil {
		return rt.Value{}, e499
	}
	t500, e501 := rt.Cond(ctx, t498)
	if e501 != nil {
		return rt.Value{}, e501
	}
	if t500 {
		t502, e503 := CeloeKNulyu(ctx, millisekundy)
		if e503 != nil {
			return rt.Value{}, e503
		}
		// «к строке»
		t504, e505 := rt.BToString(ctx, t502)
		if e505 != nil {
			return rt.Value{}, e505
		}
		t506, e507 := rt.Concat(ctx, t504, rt.Text(" мс"))
		if e507 != nil {
			return rt.Value{}, e507
		}
		return t506, nil
	} else {
		t508, e509 := rt.Div(ctx, millisekundy, rt.Number(1000.0))
		if e509 != nil {
			return rt.Value{}, e509
		}
		t510, e511 := Drobyu(ctx, t508, rt.Number(1.0), razdelitel)
		if e511 != nil {
			return rt.Value{}, e511
		}
		t512, e513 := rt.Concat(ctx, t510, rt.Text(" с"))
		if e513 != nil {
			return rt.Value{}, e513
		}
		return t512, nil
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
	t514, e515 := rt.FieldGet(ctx, hod, "осталось")
	if e515 != nil {
		return rt.Value{}, e515
	}
	t516, e517 := rt.Mod(ctx, t514, rt.Number(3.0))
	if e517 != nil {
		return rt.Value{}, e517
	}
	t518, e519 := rt.Cond(ctx, rt.Flag(rt.Equal(t516, rt.Number(0.0))))
	if e519 != nil {
		return rt.Value{}, e519
	}
	var t520 rt.Value
	if t518 {
		t521, e522 := rt.FieldGet(ctx, hod, "собрано")
		if e522 != nil {
			return rt.Value{}, e522
		}
		// «длина»
		t523, e524 := rt.BLength(ctx, t521)
		if e524 != nil {
			return rt.Value{}, e524
		}
		t525, e526 := rt.Gt(ctx, t523, rt.Number(0.0))
		if e526 != nil {
			return rt.Value{}, e526
		}
		t520 = t525
	} else {
		t520 = rt.Flag(false)
	}
	t527, e528 := rt.Cond(ctx, t520)
	if e528 != nil {
		return rt.Value{}, e528
	}
	var t529 rt.Value
	if t527 {
		t530, e531 := rt.FieldGet(ctx, hod, "собрано")
		if e531 != nil {
			return rt.Value{}, e531
		}
		t532, e533 := rt.FieldGet(ctx, hod, "разделитель")
		if e533 != nil {
			return rt.Value{}, e533
		}
		t534, e535 := rt.Concat(ctx, t530, t532)
		if e535 != nil {
			return rt.Value{}, e535
		}
		t536, e537 := rt.Concat(ctx, t534, cifra)
		if e537 != nil {
			return rt.Value{}, e537
		}
		t538, e539 := rt.FieldGet(ctx, hod, "осталось")
		if e539 != nil {
			return rt.Value{}, e539
		}
		t540, e541 := rt.Sub(ctx, t538, rt.Number(1.0))
		if e541 != nil {
			return rt.Value{}, e541
		}
		t542, e543 := rt.FieldGet(ctx, hod, "разделитель")
		if e543 != nil {
			return rt.Value{}, e543
		}
		t544 := make([]rt.Field, 3)
		t544[0] = rt.Field{Name: "собрано", Value: t536}
		t544[1] = rt.Field{Name: "осталось", Value: t540}
		t544[2] = rt.Field{Name: "разделитель", Value: t542}
		t529 = rt.Record(t544)
	} else {
		t545, e546 := rt.FieldGet(ctx, hod, "собрано")
		if e546 != nil {
			return rt.Value{}, e546
		}
		t547, e548 := rt.Concat(ctx, t545, cifra)
		if e548 != nil {
			return rt.Value{}, e548
		}
		t549, e550 := rt.FieldGet(ctx, hod, "осталось")
		if e550 != nil {
			return rt.Value{}, e550
		}
		t551, e552 := rt.Sub(ctx, t549, rt.Number(1.0))
		if e552 != nil {
			return rt.Value{}, e552
		}
		t553, e554 := rt.FieldGet(ctx, hod, "разделитель")
		if e554 != nil {
			return rt.Value{}, e554
		}
		t555 := make([]rt.Field, 3)
		t555[0] = rt.Field{Name: "собрано", Value: t547}
		t555[1] = rt.Field{Name: "осталось", Value: t551}
		t555[2] = rt.Field{Name: "разделитель", Value: t553}
		t529 = rt.Record(t555)
	}
	t556 := t529
	t557, e558 := rt.FieldGet(ctx, t556, "осталось")
	if e558 != nil {
		return rt.Value{}, e558
	}
	t559, e560 := rt.FieldGet(ctx, hod, "осталось")
	if e560 != nil {
		return rt.Value{}, e560
	}
	t561, e562 := rt.Lt(ctx, t557, t559)
	if e562 != nil {
		return rt.Value{}, e562
	}
	// постусловие «осталось убывает»
	t563, e564 := rt.Post(ctx, t561, "осталось убывает", "Шаг разряда")
	if e564 != nil {
		return rt.Value{}, e564
	}
	if !t563 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «осталось убывает» функции «Шаг разряда»")
	}
	t565, e566 := rt.FieldGet(ctx, t556, "разделитель")
	if e566 != nil {
		return rt.Value{}, e566
	}
	t567, e568 := rt.FieldGet(ctx, hod, "разделитель")
	if e568 != nil {
		return rt.Value{}, e568
	}
	// постусловие «разделитель не меняется»
	t569, e570 := rt.Post(ctx, rt.Flag(rt.Equal(t565, t567)), "разделитель не меняется", "Шаг разряда")
	if e570 != nil {
		return rt.Value{}, e570
	}
	if !t569 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «разделитель не меняется» функции «Шаг разряда»")
	}
	return t556, nil
}

// Razryadami — функция flang «Разрядами».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znachenie — «значение»: число.
// Параметр razdelitel — «разделитель»: строка.
// Результат — значение: строка.
func Razryadami(ctx *rt.Ctx, znachenie rt.Value, razdelitel rt.Value) (rt.Value, error) {
	t571, e572 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e572 != nil {
		return rt.Value{}, e572
	}
	t573, e574 := rt.Cond(ctx, t571)
	if e574 != nil {
		return rt.Value{}, e574
	}
	var t575 rt.Value
	if t573 {
		t575 = rt.Text("-")
	} else {
		t575 = rt.Text("")
	}
	// пусть «знак»
	znak := t575
	t576, e577 := rt.Lt(ctx, znachenie, rt.Number(0.0))
	if e577 != nil {
		return rt.Value{}, e577
	}
	t578, e579 := rt.Cond(ctx, t576)
	if e579 != nil {
		return rt.Value{}, e579
	}
	var t580 rt.Value
	if t578 {
		t581, e582 := rt.Sub(ctx, rt.Number(0.0), znachenie)
		if e582 != nil {
			return rt.Value{}, e582
		}
		t580 = t581
	} else {
		t580 = znachenie
	}
	t583, e584 := CeloeKNulyu(ctx, t580)
	if e584 != nil {
		return rt.Value{}, e584
	}
	// «к строке»
	t585, e586 := rt.BToString(ctx, t583)
	if e586 != nil {
		return rt.Value{}, e586
	}
	// пусть «цифры»
	cifry := t585
	// «символы»
	t587, e588 := rt.BCharacters(ctx, cifry)
	if e588 != nil {
		return rt.Value{}, e588
	}
	t589, e590 := rt.RequireList(ctx, t587, "свёртка")
	if e590 != nil {
		return rt.Value{}, e590
	}
	// «длина»
	t591, e592 := rt.BLength(ctx, cifry)
	if e592 != nil {
		return rt.Value{}, e592
	}
	t593 := make([]rt.Field, 3)
	t593[0] = rt.Field{Name: "собрано", Value: rt.Text("")}
	t593[1] = rt.Field{Name: "осталось", Value: t591}
	t593[2] = rt.Field{Name: "разделитель", Value: razdelitel}
	// «ход»
	hod := rt.Record(t593)
	for t594 := range t589 {
		// «цифра»
		cifra := t589[t594]
		t595, e596 := ShagRazryada(ctx, hod, cifra)
		if e596 != nil {
			return rt.Value{}, e596
		}
		hod = t595
	}
	t597, e598 := rt.FieldGet(ctx, hod, "собрано")
	if e598 != nil {
		return rt.Value{}, e598
	}
	t599, e600 := rt.Concat(ctx, znak, t597)
	if e600 != nil {
		return rt.Value{}, e600
	}
	return t599, nil
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
		t601, e602 := rt.BCharCodeProven(ctx, golova)
		if e602 != nil {
			return rt.Value{}, e602
		}
		t603, e604 := rt.Cond(ctx, rt.Flag(rt.Equal(t601, rt.Number(32.0))))
		if e604 != nil {
			return rt.Value{}, e604
		}
		var t605 rt.Value
		if t603 {
			t605 = rt.Flag(true)
		} else {
			// «код символа»
			t606, e607 := rt.BCharCodeProven(ctx, golova)
			if e607 != nil {
				return rt.Value{}, e607
			}
			t605 = rt.Flag(rt.Equal(t606, rt.Number(9.0)))
		}
		t608, e609 := rt.Cond(ctx, t605)
		if e609 != nil {
			return rt.Value{}, e609
		}
		var t610 rt.Value
		if t608 {
			t610 = rt.Flag(true)
		} else {
			// «код символа»
			t611, e612 := rt.BCharCodeProven(ctx, golova)
			if e612 != nil {
				return rt.Value{}, e612
			}
			t610 = rt.Flag(rt.Equal(t611, rt.Number(10.0)))
		}
		t613, e614 := rt.Cond(ctx, t610)
		if e614 != nil {
			return rt.Value{}, e614
		}
		var t615 rt.Value
		if t613 {
			t615 = rt.Flag(true)
		} else {
			// «код символа»
			t616, e617 := rt.BCharCodeProven(ctx, golova)
			if e617 != nil {
				return rt.Value{}, e617
			}
			t615 = rt.Flag(rt.Equal(t616, rt.Number(11.0)))
		}
		t618, e619 := rt.Cond(ctx, t615)
		if e619 != nil {
			return rt.Value{}, e619
		}
		var t620 rt.Value
		if t618 {
			t620 = rt.Flag(true)
		} else {
			// «код символа»
			t621, e622 := rt.BCharCodeProven(ctx, golova)
			if e622 != nil {
				return rt.Value{}, e622
			}
			t620 = rt.Flag(rt.Equal(t621, rt.Number(12.0)))
		}
		t623, e624 := rt.Cond(ctx, t620)
		if e624 != nil {
			return rt.Value{}, e624
		}
		var t625 rt.Value
		if t623 {
			t625 = rt.Flag(true)
		} else {
			// «код символа»
			t626, e627 := rt.BCharCodeProven(ctx, golova)
			if e627 != nil {
				return rt.Value{}, e627
			}
			t625 = rt.Flag(rt.Equal(t626, rt.Number(13.0)))
		}
		t628, e629 := rt.Cond(ctx, t625)
		if e629 != nil {
			return rt.Value{}, e629
		}
		var t630 rt.Value
		if t628 {
			t630 = rt.Flag(true)
		} else {
			// «код символа»
			t631, e632 := rt.BCharCodeProven(ctx, golova)
			if e632 != nil {
				return rt.Value{}, e632
			}
			t630 = rt.Flag(rt.Equal(t631, rt.Number(133.0)))
		}
		t633, e634 := rt.Cond(ctx, t630)
		if e634 != nil {
			return rt.Value{}, e634
		}
		if t633 {
			return rt.Flag(true), nil
		} else {
			// «код символа»
			t635, e636 := rt.BCharCodeProven(ctx, golova)
			if e636 != nil {
				return rt.Value{}, e636
			}
			return rt.Flag(rt.Equal(t635, rt.Number(160.0))), nil
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
	t637, e638 := rt.BCharacters(ctx, tekst)
	if e638 != nil {
		return rt.Value{}, e638
	}
	t639, e640 := rt.RequireList(ctx, t637, "отфильтровать")
	if e640 != nil {
		return rt.Value{}, e640
	}
	t641 := make([]rt.Value, 0, len(t639))
	for t642 := range t639 {
		// «знак»
		znak := t639[t642]
		t643, e644 := Probelnyy(ctx, znak)
		if e644 != nil {
			return rt.Value{}, e644
		}
		t645, e646 := rt.Cond(ctx, t643)
		if e646 != nil {
			return rt.Value{}, e646
		}
		var t647 rt.Value
		if t645 {
			t647 = rt.Flag(false)
		} else {
			t647 = rt.Flag(true)
		}
		t648, e649 := rt.Keep(ctx, t647)
		if e649 != nil {
			return rt.Value{}, e649
		}
		if t648 {
			t641 = append(t641, znak)
		}
	}
	// «длина»
	t650, e651 := rt.BLength(ctx, rt.List(t641))
	if e651 != nil {
		return rt.Value{}, e651
	}
	t652, e653 := rt.Cond(ctx, rt.Flag(rt.Equal(t650, rt.Number(0.0))))
	if e653 != nil {
		return rt.Value{}, e653
	}
	if t652 {
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
	t654, e655 := rt.BLength(ctx, tekst)
	if e655 != nil {
		return rt.Value{}, e655
	}
	t656, e657 := rt.Lte(ctx, t654, skolko)
	if e657 != nil {
		return rt.Value{}, e657
	}
	t658, e659 := rt.Cond(ctx, t656)
	if e659 != nil {
		return rt.Value{}, e659
	}
	var t660 rt.Value
	if t658 {
		t660 = tekst
	} else {
		t661, e662 := rt.Lte(ctx, skolko, rt.Number(0.0))
		if e662 != nil {
			return rt.Value{}, e662
		}
		t663, e664 := rt.Cond(ctx, t661)
		if e664 != nil {
			return rt.Value{}, e664
		}
		var t665 rt.Value
		if t663 {
			t665 = rt.Text("")
		} else {
			t666, e667 := rt.Lte(ctx, skolko, rt.Number(1.0))
			if e667 != nil {
				return rt.Value{}, e667
			}
			t668, e669 := rt.Cond(ctx, t666)
			if e669 != nil {
				return rt.Value{}, e669
			}
			var t670 rt.Value
			if t668 {
				// «подстрока»
				t671, e672 := rt.BSubstring(ctx, tekst, rt.Number(1.0), rt.Number(1.0))
				if e672 != nil {
					return rt.Value{}, e672
				}
				t670 = t671
			} else {
				t673, e674 := rt.Sub(ctx, skolko, rt.Number(1.0))
				if e674 != nil {
					return rt.Value{}, e674
				}
				// «подстрока»
				t675, e676 := rt.BSubstring(ctx, tekst, rt.Number(1.0), t673)
				if e676 != nil {
					return rt.Value{}, e676
				}
				t677, e678 := rt.Concat(ctx, t675, rt.Text("…"))
				if e678 != nil {
					return rt.Value{}, e678
				}
				t670 = t677
			}
			t665 = t670
		}
		t660 = t665
	}
	t679 := t660
	t680, e681 := rt.Gte(ctx, skolko, rt.Number(0.0))
	if e681 != nil {
		return rt.Value{}, e681
	}
	t682, e683 := rt.Cond(ctx, t680)
	if e683 != nil {
		return rt.Value{}, e683
	}
	var t684 rt.Value
	if t682 {
		// «длина»
		t685, e686 := rt.BLength(ctx, t679)
		if e686 != nil {
			return rt.Value{}, e686
		}
		// «длина»
		t687, e688 := rt.BLength(ctx, tekst)
		if e688 != nil {
			return rt.Value{}, e688
		}
		t689, e690 := rt.Lte(ctx, t687, skolko)
		if e690 != nil {
			return rt.Value{}, e690
		}
		t691, e692 := rt.Cond(ctx, t689)
		if e692 != nil {
			return rt.Value{}, e692
		}
		var t693 rt.Value
		if t691 {
			// «длина»
			t694, e695 := rt.BLength(ctx, tekst)
			if e695 != nil {
				return rt.Value{}, e695
			}
			t693 = t694
		} else {
			t693 = skolko
		}
		t696, e697 := rt.Lte(ctx, t685, t693)
		if e697 != nil {
			return rt.Value{}, e697
		}
		t684 = t696
	} else {
		t684 = rt.Flag(true)
	}
	// постусловие «обрезанное не длиннее заказа»
	t698, e699 := rt.Post(ctx, t684, "обрезанное не длиннее заказа", "Обрезать")
	if e699 != nil {
		return rt.Value{}, e699
	}
	if !t698 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «обрезанное не длиннее заказа» функции «Обрезать»")
	}
	return t679, nil
}

// MeraUbyvaet — функция flang «мера убывает».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр shag — «шаг»: число.
// Параметр mera — «мера»: число.
// Результат — значение: число.
func MeraUbyvaet(ctx *rt.Ctx, shag rt.Value, mera rt.Value) (rt.Value, error) {
	t700 := shag
	t701, e702 := rt.Lt(ctx, t700, mera)
	if e702 != nil {
		return rt.Value{}, e702
	}
	// постусловие «мера убывает»
	t703, e704 := rt.Post(ctx, t701, "мера убывает", "мера убывает")
	if e704 != nil {
		return rt.Value{}, e704
	}
	if !t703 {
		return rt.Value{}, rt.Fail("FLANG_MEASURE", "%s", "тотальная функция «Повторить»: мера не убыла — аргумент 2 вызова «Повторить» не стал меньше параметра «сколько». Завершение доказано убыванием этой меры, а числа flang — IEEE-754 double: при большом |«сколько»| постоянный шаг не меняет значение, и спуск не идёт. Отказ здесь честнее зацикливания")
	}
	return t700, nil
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
