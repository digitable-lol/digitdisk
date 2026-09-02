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

// ShestnadcaterichnyyZnak — функция flang «Шестнадцатеричный знак».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр znak — «знак»: строка.
// Результат — значение.
func ShestnadcaterichnyyZnak(ctx *rt.Ctx, znak rt.Value) (rt.Value, error) {
	// «содержит»
	t32, e33 := rt.BContains(ctx, rt.Text("0123456789abcdefABCDEF"), znak)
	if e33 != nil {
		return rt.Value{}, e33
	}
	return t32, nil
}

// PohozheNaOtpechatok — функция flang «Похоже на отпечаток».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр chast — «часть»: строка.
// Результат — значение.
func PohozheNaOtpechatok(ctx *rt.Ctx, chast rt.Value) (rt.Value, error) {
	// «длина»
	t34, e35 := rt.BLength(ctx, chast)
	if e35 != nil {
		return rt.Value{}, e35
	}
	t36, e37 := rt.Lt(ctx, t34, rt.Number(32.0))
	if e37 != nil {
		return rt.Value{}, e37
	}
	t38, e39 := rt.Cond(ctx, t36)
	if e39 != nil {
		return rt.Value{}, e39
	}
	if t38 {
		return rt.Flag(false), nil
	} else {
		// «символы»
		t40, e41 := rt.BCharacters(ctx, chast)
		if e41 != nil {
			return rt.Value{}, e41
		}
		t42, e43 := rt.RequireList(ctx, t40, "свёртка")
		if e43 != nil {
			return rt.Value{}, e43
		}
		// «собрано»
		sobrano := rt.Flag(true)
		for t44 := range t42 {
			// «знак»
			znak := t42[t44]
			t45, e46 := rt.Cond(ctx, sobrano)
			if e46 != nil {
				return rt.Value{}, e46
			}
			var t47 rt.Value
			if t45 {
				t48, e49 := ShestnadcaterichnyyZnak(ctx, znak)
				if e49 != nil {
					return rt.Value{}, e49
				}
				t47 = t48
			} else {
				t47 = rt.Flag(false)
			}
			sobrano = t47
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
	t50, e51 := SostavlyayuschiePuti(ctx, put)
	if e51 != nil {
		return rt.Value{}, e51
	}
	t52, e53 := rt.RequireList(ctx, t50, "отфильтровать")
	if e53 != nil {
		return rt.Value{}, e53
	}
	t54 := make([]rt.Value, 0, len(t52))
	for t55 := range t52 {
		// «часть»
		chast := t52[t55]
		t56, e57 := PohozheNaOtpechatok(ctx, chast)
		if e57 != nil {
			return rt.Value{}, e57
		}
		t58, e59 := rt.Keep(ctx, t56)
		if e59 != nil {
			return rt.Value{}, e59
		}
		if t58 {
			t54 = append(t54, chast)
		}
	}
	// «длина»
	t60, e61 := rt.BLength(ctx, rt.List(t54))
	if e61 != nil {
		return rt.Value{}, e61
	}
	t62, e63 := rt.Gt(ctx, t60, rt.Number(0.0))
	if e63 != nil {
		return rt.Value{}, e63
	}
	t64, e65 := rt.Cond(ctx, t62)
	if e65 != nil {
		return rt.Value{}, e65
	}
	var t66 rt.Value
	if t64 {
		t66 = rt.Flag(true)
	} else {
		t67, e68 := EstSostavlyayuschaya(ctx, put, rt.Text("site-packages"))
		if e68 != nil {
			return rt.Value{}, e68
		}
		t66 = t67
	}
	t69, e70 := rt.Cond(ctx, t66)
	if e70 != nil {
		return rt.Value{}, e70
	}
	var t71 rt.Value
	if t69 {
		t71 = rt.Flag(true)
	} else {
		t72, e73 := EstSostavlyayuschaya(ctx, put, rt.Text("dist-packages"))
		if e73 != nil {
			return rt.Value{}, e73
		}
		t71 = t72
	}
	t74, e75 := rt.Cond(ctx, t71)
	if e75 != nil {
		return rt.Value{}, e75
	}
	if t74 {
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
	t76, e77 := rt.BStartsWith(ctx, put, rt.Text("/tmp/"))
	if e77 != nil {
		return rt.Value{}, e77
	}
	t78, e79 := rt.Cond(ctx, t76)
	if e79 != nil {
		return rt.Value{}, e79
	}
	if t78 {
		return rt.Flag(true), nil
	} else {
		// «начинается с»
		t80, e81 := rt.BStartsWith(ctx, put, rt.Text("/var/tmp/"))
		if e81 != nil {
			return rt.Value{}, e81
		}
		return t80, nil
	}
}

// PrimetaKesha — функция flang «Примета кэша».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр put — «путь»: строка.
// Результат — значение.
func PrimetaKesha(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t82, e83 := EstSostavlyayuschaya(ctx, put, rt.Text(".cache"))
	if e83 != nil {
		return rt.Value{}, e83
	}
	t84, e85 := rt.Cond(ctx, t82)
	if e85 != nil {
		return rt.Value{}, e85
	}
	var t86 rt.Value
	if t84 {
		t86 = rt.Flag(true)
	} else {
		t87, e88 := EstSostavlyayuschaya(ctx, put, rt.Text("cache"))
		if e88 != nil {
			return rt.Value{}, e88
		}
		t86 = t87
	}
	t89, e90 := rt.Cond(ctx, t86)
	if e90 != nil {
		return rt.Value{}, e90
	}
	var t91 rt.Value
	if t89 {
		t91 = rt.Flag(true)
	} else {
		t92, e93 := EstSostavlyayuschaya(ctx, put, rt.Text("Caches"))
		if e93 != nil {
			return rt.Value{}, e93
		}
		t91 = t92
	}
	t94, e95 := rt.Cond(ctx, t91)
	if e95 != nil {
		return rt.Value{}, e95
	}
	if t94 {
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
	t96, e97 := ImyaVPuti(ctx, put)
	if e97 != nil {
		return rt.Value{}, e97
	}
	t98, e99 := OkanchivaetsyaNa(ctx, t96, rt.Text(".log"))
	if e99 != nil {
		return rt.Value{}, e99
	}
	t100, e101 := rt.Cond(ctx, t98)
	if e101 != nil {
		return rt.Value{}, e101
	}
	var t102 rt.Value
	if t100 {
		t102 = rt.Flag(true)
	} else {
		t103, e104 := EstSostavlyayuschaya(ctx, put, rt.Text("log"))
		if e104 != nil {
			return rt.Value{}, e104
		}
		t102 = t103
	}
	t105, e106 := rt.Cond(ctx, t102)
	if e106 != nil {
		return rt.Value{}, e106
	}
	if t105 {
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
	t107, e108 := EstSostavlyayuschaya(ctx, put, rt.Text("node_modules"))
	if e108 != nil {
		return rt.Value{}, e108
	}
	t109, e110 := rt.Cond(ctx, t107)
	if e110 != nil {
		return rt.Value{}, e110
	}
	var t111 rt.Value
	if t109 {
		t111 = rt.Flag(true)
	} else {
		t112, e113 := EstSostavlyayuschaya(ctx, put, rt.Text("target"))
		if e113 != nil {
			return rt.Value{}, e113
		}
		t111 = t112
	}
	t114, e115 := rt.Cond(ctx, t111)
	if e115 != nil {
		return rt.Value{}, e115
	}
	var t116 rt.Value
	if t114 {
		t116 = rt.Flag(true)
	} else {
		t117, e118 := EstSostavlyayuschaya(ctx, put, rt.Text("build"))
		if e118 != nil {
			return rt.Value{}, e118
		}
		t116 = t117
	}
	t119, e120 := rt.Cond(ctx, t116)
	if e120 != nil {
		return rt.Value{}, e120
	}
	var t121 rt.Value
	if t119 {
		t121 = rt.Flag(true)
	} else {
		t122, e123 := EstSostavlyayuschaya(ctx, put, rt.Text("_build"))
		if e123 != nil {
			return rt.Value{}, e123
		}
		t121 = t122
	}
	t124, e125 := rt.Cond(ctx, t121)
	if e125 != nil {
		return rt.Value{}, e125
	}
	if t124 {
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
	t126, e127 := EstSostavlyayuschaya(ctx, put, rt.Text("Downloads"))
	if e127 != nil {
		return rt.Value{}, e127
	}
	t128, e129 := rt.Cond(ctx, t126)
	if e129 != nil {
		return rt.Value{}, e129
	}
	if t128 {
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
// Результат — значение.
func EstPrimeta(ctx *rt.Ctx, put rt.Value) (rt.Value, error) {
	t130, e131 := PrimetaKesha(ctx, put)
	if e131 != nil {
		return rt.Value{}, e131
	}
	t132, e133 := rt.Cond(ctx, t130)
	if e133 != nil {
		return rt.Value{}, e133
	}
	var t134 rt.Value
	if t132 {
		t134 = rt.Flag(true)
	} else {
		t135, e136 := PrimetaZhurnala(ctx, put)
		if e136 != nil {
			return rt.Value{}, e136
		}
		t134 = t135
	}
	t137, e138 := rt.Cond(ctx, t134)
	if e138 != nil {
		return rt.Value{}, e138
	}
	var t139 rt.Value
	if t137 {
		t139 = rt.Flag(true)
	} else {
		t140, e141 := PrimetaSborki(ctx, put)
		if e141 != nil {
			return rt.Value{}, e141
		}
		t139 = t140
	}
	t142, e143 := rt.Cond(ctx, t139)
	if e143 != nil {
		return rt.Value{}, e143
	}
	if t142 {
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
// Результат — значение: «Разряд».
func RazryadNahodki(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t144, e145 := rt.FieldGet(ctx, nahodka, "путь")
	if e145 != nil {
		return rt.Value{}, e145
	}
	t146, e147 := PrimetaKesha(ctx, t144)
	if e147 != nil {
		return rt.Value{}, e147
	}
	t148, e149 := rt.Cond(ctx, t146)
	if e149 != nil {
		return rt.Value{}, e149
	}
	var t150 rt.Value
	if t148 {
		t150 = rt.Variant("Кэш", nil)
	} else {
		t151, e152 := rt.FieldGet(ctx, nahodka, "путь")
		if e152 != nil {
			return rt.Value{}, e152
		}
		t153, e154 := PrimetaZhurnala(ctx, t151)
		if e154 != nil {
			return rt.Value{}, e154
		}
		t155, e156 := rt.Cond(ctx, t153)
		if e156 != nil {
			return rt.Value{}, e156
		}
		var t157 rt.Value
		if t155 {
			t157 = rt.Variant("Журнал", nil)
		} else {
			t158, e159 := rt.FieldGet(ctx, nahodka, "путь")
			if e159 != nil {
				return rt.Value{}, e159
			}
			t160, e161 := PrimetaSborki(ctx, t158)
			if e161 != nil {
				return rt.Value{}, e161
			}
			t162, e163 := rt.Cond(ctx, t160)
			if e163 != nil {
				return rt.Value{}, e163
			}
			var t164 rt.Value
			if t162 {
				t164 = rt.Variant("Сборка", nil)
			} else {
				t165, e166 := rt.FieldGet(ctx, nahodka, "путь")
				if e166 != nil {
					return rt.Value{}, e166
				}
				t167, e168 := PrimetaZagruzki(ctx, t165)
				if e168 != nil {
					return rt.Value{}, e168
				}
				t169, e170 := rt.Cond(ctx, t167)
				if e170 != nil {
					return rt.Value{}, e170
				}
				var t171 rt.Value
				if t169 {
					t171 = rt.Variant("Загрузка", nil)
				} else {
					t172, e173 := rt.FieldGet(ctx, nahodka, "размер")
					if e173 != nil {
						return rt.Value{}, e173
					}
					t174, e175 := PorogKrupnogo(ctx)
					if e175 != nil {
						return rt.Value{}, e175
					}
					t176, e177 := rt.Gte(ctx, t172, t174)
					if e177 != nil {
						return rt.Value{}, e177
					}
					t178, e179 := rt.Cond(ctx, t176)
					if e179 != nil {
						return rt.Value{}, e179
					}
					var t180 rt.Value
					if t178 {
						t180 = rt.Variant("Крупное", nil)
					} else {
						t180 = rt.Variant("Неизвестное", nil)
					}
					t171 = t180
				}
				t164 = t171
			}
			t157 = t164
		}
		t150 = t157
	}
	t181 := t150
	t182, e183 := RazryadObosnovan(ctx, nahodka, t181)
	if e183 != nil {
		return rt.Value{}, e183
	}
	// постусловие «Разряд обоснован приметой»
	t184, e185 := rt.Post(ctx, t182, "Разряд обоснован приметой", "Разряд находки")
	if e185 != nil {
		return rt.Value{}, e185
	}
	if !t184 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Разряд обоснован приметой» функции «Разряд находки»")
	}
	t186, e187 := rt.FieldGet(ctx, nahodka, "путь")
	if e187 != nil {
		return rt.Value{}, e187
	}
	t188, e189 := EstPrimeta(ctx, t186)
	if e189 != nil {
		return rt.Value{}, e189
	}
	t190, e191 := rt.Cond(ctx, t188)
	if e191 != nil {
		return rt.Value{}, e191
	}
	var t192 rt.Value
	if t190 {
		t192 = rt.Flag(false)
	} else {
		t192 = rt.Flag(true)
	}
	t193, e194 := rt.Cond(ctx, t192)
	if e194 != nil {
		return rt.Value{}, e194
	}
	var t195 rt.Value
	if t193 {
		t196, e197 := RazryadReshyonRazmerom(ctx, t181)
		if e197 != nil {
			return rt.Value{}, e197
		}
		t195 = t196
	} else {
		t195 = rt.Flag(true)
	}
	// постусловие «И4: без приметы-составляющей разряд решает размер»
	t198, e199 := rt.Post(ctx, t195, "И4: без приметы-составляющей разряд решает размер", "Разряд находки")
	if e199 != nil {
		return rt.Value{}, e199
	}
	if !t198 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И4: без приметы-составляющей разряд решает размер» функции «Разряд находки»")
	}
	return t181, nil
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
		t200, e201 := rt.FieldGet(ctx, nahodka, "размер")
		if e201 != nil {
			return rt.Value{}, e201
		}
		t202, e203 := PorogKrupnogo(ctx)
		if e203 != nil {
			return rt.Value{}, e203
		}
		t204, e205 := rt.Gte(ctx, t200, t202)
		if e205 != nil {
			return rt.Value{}, e205
		}
		return t204, nil
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
	t206, e207 := rt.FieldGet(ctx, nahodka, "вид")
	if e207 != nil {
		return rt.Value{}, e207
	}
	if rt.VariantIs(t206, "Каталог") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t206, "Файл") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t206, "Ссылка") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t206)
	}
}

// Ssylka — функция flang «Ссылка».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение.
func Ssylka(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t208, e209 := rt.FieldGet(ctx, nahodka, "вид")
	if e209 != nil {
		return rt.Value{}, e209
	}
	if rt.VariantIs(t208, "Ссылка") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t208, "Файл") {
		return rt.Flag(false), nil
	} else if rt.VariantIs(t208, "Каталог") {
		return rt.Flag(false), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t208)
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
	t210, e211 := Katalog(ctx, nahodka)
	if e211 != nil {
		return rt.Value{}, e211
	}
	t212, e213 := rt.Cond(ctx, t210)
	if e213 != nil {
		return rt.Value{}, e213
	}
	var t214 rt.Value
	if t212 {
		t214 = rt.Variant("Спросить", nil)
	} else {
		t215, e216 := rt.FieldGet(ctx, nahodka, "возраст_дней")
		if e216 != nil {
			return rt.Value{}, e216
		}
		t217, e218 := rt.Gte(ctx, t215, porog)
		if e218 != nil {
			return rt.Value{}, e218
		}
		t219, e220 := rt.Cond(ctx, t217)
		if e220 != nil {
			return rt.Value{}, e220
		}
		var t221 rt.Value
		if t219 {
			t221 = rt.Variant("МожноУбрать", nil)
		} else {
			t221 = rt.Variant("Спросить", nil)
		}
		t214 = t221
	}
	t222 := t214
	t223, e224 := Katalog(ctx, nahodka)
	if e224 != nil {
		return rt.Value{}, e224
	}
	t225, e226 := rt.Cond(ctx, t223)
	if e226 != nil {
		return rt.Value{}, e226
	}
	var t227 rt.Value
	if t225 {
		t228, e229 := EtoMozhnoUbrat(ctx, t222)
		if e229 != nil {
			return rt.Value{}, e229
		}
		t230, e231 := rt.Cond(ctx, t228)
		if e231 != nil {
			return rt.Value{}, e231
		}
		var t232 rt.Value
		if t230 {
			t232 = rt.Flag(false)
		} else {
			t232 = rt.Flag(true)
		}
		t227 = t232
	} else {
		t227 = rt.Flag(true)
	}
	// постусловие «Каталог не убирается»
	t233, e234 := rt.Post(ctx, t227, "Каталог не убирается", "Приговор мусора")
	if e234 != nil {
		return rt.Value{}, e234
	}
	if !t233 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Каталог не убирается» функции «Приговор мусора»")
	}
	return t222, nil
}

// PrigovorNahodki — функция flang «Приговор находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр razryad — «разряд»: «Разряд».
// Результат — значение: «Приговор».
func PrigovorNahodki(ctx *rt.Ctx, nahodka rt.Value, razryad rt.Value) (rt.Value, error) {
	t235, e236 := rt.FieldGet(ctx, nahodka, "доступен")
	if e236 != nil {
		return rt.Value{}, e236
	}
	t237, e238 := rt.Cond(ctx, t235)
	if e238 != nil {
		return rt.Value{}, e238
	}
	var t239 rt.Value
	if t237 {
		t239 = rt.Flag(false)
	} else {
		t239 = rt.Flag(true)
	}
	t240, e241 := rt.Cond(ctx, t239)
	if e241 != nil {
		return rt.Value{}, e241
	}
	var t242 rt.Value
	if t240 {
		t242 = rt.Variant("НеТрогать", nil)
	} else {
		t243, e244 := Ssylka(ctx, nahodka)
		if e244 != nil {
			return rt.Value{}, e244
		}
		t245, e246 := rt.Cond(ctx, t243)
		if e246 != nil {
			return rt.Value{}, e246
		}
		var t247 rt.Value
		if t245 {
			t247 = rt.Variant("НеТрогать", nil)
		} else {
			t248, e249 := rt.FieldGet(ctx, nahodka, "путь")
			if e249 != nil {
				return rt.Value{}, e249
			}
			t250, e251 := AdresuetsyaSoderzhimym(ctx, t248)
			if e251 != nil {
				return rt.Value{}, e251
			}
			t252, e253 := rt.Cond(ctx, t250)
			if e253 != nil {
				return rt.Value{}, e253
			}
			var t254 rt.Value
			if t252 {
				t254 = rt.Variant("НеТрогать", nil)
			} else {
				var t255 rt.Value
				if rt.VariantIs(razryad, "Кэш") {
					t256, e257 := PorogKesha(ctx)
					if e257 != nil {
						return rt.Value{}, e257
					}
					t258, e259 := PrigovorMusora(ctx, nahodka, t256)
					if e259 != nil {
						return rt.Value{}, e259
					}
					t255 = t258
				} else if rt.VariantIs(razryad, "Сборка") {
					t260, e261 := PorogKesha(ctx)
					if e261 != nil {
						return rt.Value{}, e261
					}
					t262, e263 := PrigovorMusora(ctx, nahodka, t260)
					if e263 != nil {
						return rt.Value{}, e263
					}
					t255 = t262
				} else if rt.VariantIs(razryad, "Журнал") {
					t264, e265 := PorogZhurnala(ctx)
					if e265 != nil {
						return rt.Value{}, e265
					}
					t266, e267 := PrigovorMusora(ctx, nahodka, t264)
					if e267 != nil {
						return rt.Value{}, e267
					}
					t255 = t266
				} else if rt.VariantIs(razryad, "Загрузка") {
					t268, e269 := rt.FieldGet(ctx, nahodka, "возраст_дней")
					if e269 != nil {
						return rt.Value{}, e269
					}
					t270, e271 := PorogZagruzki(ctx)
					if e271 != nil {
						return rt.Value{}, e271
					}
					t272, e273 := rt.Gte(ctx, t268, t270)
					if e273 != nil {
						return rt.Value{}, e273
					}
					t274, e275 := rt.Cond(ctx, t272)
					if e275 != nil {
						return rt.Value{}, e275
					}
					var t276 rt.Value
					if t274 {
						t276 = rt.Variant("Спросить", nil)
					} else {
						t276 = rt.Variant("НеТрогать", nil)
					}
					t255 = t276
				} else if rt.VariantIs(razryad, "Крупное") {
					t255 = rt.Variant("Спросить", nil)
				} else if rt.VariantIs(razryad, "Неизвестное") {
					t255 = rt.Variant("НеТрогать", nil)
				} else {
					return rt.Value{}, rt.MatchFail(ctx, razryad)
				}
				t254 = t255
			}
			t247 = t254
		}
		t242 = t247
	}
	t277 := t242
	t278, e279 := PrigovorObosnovan(ctx, nahodka, razryad, t277)
	if e279 != nil {
		return rt.Value{}, e279
	}
	// постусловие «Приговор обоснован»
	t280, e281 := rt.Post(ctx, t278, "Приговор обоснован", "Приговор находки")
	if e281 != nil {
		return rt.Value{}, e281
	}
	if !t280 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Приговор обоснован» функции «Приговор находки»")
	}
	t282, e283 := rt.FieldGet(ctx, nahodka, "путь")
	if e283 != nil {
		return rt.Value{}, e283
	}
	t284, e285 := AdresuetsyaSoderzhimym(ctx, t282)
	if e285 != nil {
		return rt.Value{}, e285
	}
	t286, e287 := rt.Cond(ctx, t284)
	if e287 != nil {
		return rt.Value{}, e287
	}
	var t288 rt.Value
	if t286 {
		t289, e290 := EtoNeTrogat(ctx, t277)
		if e290 != nil {
			return rt.Value{}, e290
		}
		t288 = t289
	} else {
		t288 = rt.Flag(true)
	}
	// постусловие «И3: адресуемое содержимым не убирается»
	t291, e292 := rt.Post(ctx, t288, "И3: адресуемое содержимым не убирается", "Приговор находки")
	if e292 != nil {
		return rt.Value{}, e292
	}
	if !t291 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И3: адресуемое содержимым не убирается» функции «Приговор находки»")
	}
	return t277, nil
}

// VesNahodki — функция flang «Вес находки».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр prigovor — «приговор»: «Приговор».
// Результат — значение: число.
func VesNahodki(ctx *rt.Ctx, nahodka rt.Value, prigovor rt.Value) (rt.Value, error) {
	var t293 rt.Value
	if rt.VariantIs(prigovor, "НеТрогать") {
		t293 = rt.Number(0.0)
	} else if rt.VariantIs(prigovor, "МожноУбрать") {
		t294, e295 := rt.FieldGet(ctx, nahodka, "размер")
		if e295 != nil {
			return rt.Value{}, e295
		}
		t293 = t294
	} else if rt.VariantIs(prigovor, "Спросить") {
		t296, e297 := rt.FieldGet(ctx, nahodka, "размер")
		if e297 != nil {
			return rt.Value{}, e297
		}
		t293 = t296
	} else {
		return rt.Value{}, rt.MatchFail(ctx, prigovor)
	}
	t298 := t293
	t299, e300 := VesObosnovan(ctx, nahodka, prigovor, t298)
	if e300 != nil {
		return rt.Value{}, e300
	}
	// постусловие «Вес обоснован»
	t301, e302 := rt.Post(ctx, t299, "Вес обоснован", "Вес находки")
	if e302 != nil {
		return rt.Value{}, e302
	}
	if !t301 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Вес обоснован» функции «Вес находки»")
	}
	return t298, nil
}

// VesVGranicah — функция flang «Вес в границах».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Параметр ves — «вес»: число.
// Результат — значение.
func VesVGranicah(ctx *rt.Ctx, nahodka rt.Value, ves rt.Value) (rt.Value, error) {
	t303, e304 := rt.Gte(ctx, ves, rt.Number(0.0))
	if e304 != nil {
		return rt.Value{}, e304
	}
	t305, e306 := rt.Cond(ctx, t303)
	if e306 != nil {
		return rt.Value{}, e306
	}
	if t305 {
		t307, e308 := rt.FieldGet(ctx, nahodka, "размер")
		if e308 != nil {
			return rt.Value{}, e308
		}
		t309, e310 := rt.Lte(ctx, ves, t307)
		if e310 != nil {
			return rt.Value{}, e310
		}
		return t309, nil
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
	t311, e312 := RazryadNahodki(ctx, nahodka)
	if e312 != nil {
		return rt.Value{}, e312
	}
	// пусть «разряд»
	razryad := t311
	t313, e314 := PrigovorNahodki(ctx, nahodka, razryad)
	if e314 != nil {
		return rt.Value{}, e314
	}
	// пусть «приговор»
	prigovor := t313
	t315, e316 := VesNahodki(ctx, nahodka, prigovor)
	if e316 != nil {
		return rt.Value{}, e316
	}
	t317 := make([]rt.Field, 3)
	t317[0] = rt.Field{Name: "разряд", Value: razryad}
	t317[1] = rt.Field{Name: "приговор", Value: prigovor}
	t317[2] = rt.Field{Name: "вес", Value: t315}
	t318 := rt.Record(t317)
	t319, e320 := I1Derzhitsya(ctx, t318)
	if e320 != nil {
		return rt.Value{}, e320
	}
	// постусловие «И1: убрать можно только мусор»
	t321, e322 := rt.Post(ctx, t319, "И1: убрать можно только мусор", "Решить находку")
	if e322 != nil {
		return rt.Value{}, e322
	}
	if !t321 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И1: убрать можно только мусор» функции «Решить находку»")
	}
	return t318, nil
}

// ReshitVsyo — функция flang «Решить всё».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: список: «Решение».
func ReshitVsyo(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t323, e324 := rt.RequireList(ctx, zapisi, "отобразить")
	if e324 != nil {
		return rt.Value{}, e324
	}
	t325 := make([]rt.Value, 0, len(t323))
	for t326 := range t323 {
		// «находка»
		nahodka := t323[t326]
		t327, e328 := ReshitNahodku(ctx, nahodka)
		if e328 != nil {
			return rt.Value{}, e328
		}
		t325 = append(t325, t327)
	}
	t329 := rt.List(t325)
	t330, e331 := I1DerzhitsyaVsyudu(ctx, t329)
	if e331 != nil {
		return rt.Value{}, e331
	}
	// постусловие «И1 всюду: убрать можно только мусор»
	t332, e333 := rt.Post(ctx, t330, "И1 всюду: убрать можно только мусор", "Решить всё")
	if e333 != nil {
		return rt.Value{}, e333
	}
	if !t332 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И1 всюду: убрать можно только мусор» функции «Решить всё»")
	}
	return t329, nil
}

// I1Derzhitsya — функция flang «И1 держится».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр reshenie — «решение»: «Решение».
// Результат — значение.
func I1Derzhitsya(ctx *rt.Ctx, reshenie rt.Value) (rt.Value, error) {
	t334, e335 := rt.FieldGet(ctx, reshenie, "приговор")
	if e335 != nil {
		return rt.Value{}, e335
	}
	if rt.VariantIs(t334, "МожноУбрать") {
		t336, e337 := rt.FieldGet(ctx, reshenie, "разряд")
		if e337 != nil {
			return rt.Value{}, e337
		}
		if rt.VariantIs(t336, "Кэш") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t336, "Журнал") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t336, "Сборка") {
			return rt.Flag(true), nil
		} else if rt.VariantIs(t336, "Загрузка") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t336, "Крупное") {
			return rt.Flag(false), nil
		} else if rt.VariantIs(t336, "Неизвестное") {
			return rt.Flag(false), nil
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t336)
		}
	} else if rt.VariantIs(t334, "Спросить") {
		return rt.Flag(true), nil
	} else if rt.VariantIs(t334, "НеТрогать") {
		return rt.Flag(true), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t334)
	}
}

// I1DerzhitsyaVsyudu — функция flang «И1 держится всюду».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр resheniya — «решения»: список: «Решение».
// Результат — значение.
func I1DerzhitsyaVsyudu(ctx *rt.Ctx, resheniya rt.Value) (rt.Value, error) {
	t338, e339 := rt.RequireList(ctx, resheniya, "свёртка")
	if e339 != nil {
		return rt.Value{}, e339
	}
	// «акк»
	akk := rt.Flag(true)
	for t340 := range t338 {
		// «решение»
		reshenie := t338[t340]
		t341, e342 := rt.Cond(ctx, akk)
		if e342 != nil {
			return rt.Value{}, e342
		}
		var t343 rt.Value
		if t341 {
			t344, e345 := I1Derzhitsya(ctx, reshenie)
			if e345 != nil {
				return rt.Value{}, e345
			}
			t343 = t344
		} else {
			t343 = rt.Flag(false)
		}
		akk = t343
	}
	return akk, nil
}

// PustoySvod — функция flang «Пустой свод».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Результат — значение: «Свод».
func PustoySvod(ctx *rt.Ctx) (rt.Value, error) {
	t346 := make([]rt.Field, 7)
	t346[0] = rt.Field{Name: "кэш", Value: rt.Number(0.0)}
	t346[1] = rt.Field{Name: "журнал", Value: rt.Number(0.0)}
	t346[2] = rt.Field{Name: "сборка", Value: rt.Number(0.0)}
	t346[3] = rt.Field{Name: "загрузка", Value: rt.Number(0.0)}
	t346[4] = rt.Field{Name: "крупное", Value: rt.Number(0.0)}
	t346[5] = rt.Field{Name: "неизвестное", Value: rt.Number(0.0)}
	t346[6] = rt.Field{Name: "освободить", Value: rt.Number(0.0)}
	return rt.Record(t346), nil
}

// PribavitReshenie — функция flang «Прибавить решение».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр svod — «свод»: «Свод».
// Параметр reshenie — «решение»: «Решение».
// Результат — значение: «Свод».
func PribavitReshenie(ctx *rt.Ctx, svod rt.Value, reshenie rt.Value) (rt.Value, error) {
	var t347 rt.Value
	t348, e349 := rt.FieldGet(ctx, reshenie, "приговор")
	if e349 != nil {
		return rt.Value{}, e349
	}
	if rt.VariantIs(t348, "МожноУбрать") {
		t350, e351 := rt.FieldGet(ctx, reshenie, "вес")
		if e351 != nil {
			return rt.Value{}, e351
		}
		t347 = t350
	} else if rt.VariantIs(t348, "Спросить") {
		t347 = rt.Number(0.0)
	} else if rt.VariantIs(t348, "НеТрогать") {
		t347 = rt.Number(0.0)
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t348)
	}
	// пусть «убрать»
	ubrat := t347
	t352, e353 := rt.FieldGet(ctx, reshenie, "разряд")
	if e353 != nil {
		return rt.Value{}, e353
	}
	if rt.VariantIs(t352, "Кэш") {
		t354, e355 := rt.FieldGet(ctx, svod, "кэш")
		if e355 != nil {
			return rt.Value{}, e355
		}
		t356, e357 := rt.FieldGet(ctx, reshenie, "вес")
		if e357 != nil {
			return rt.Value{}, e357
		}
		t358, e359 := rt.Add(ctx, t354, t356)
		if e359 != nil {
			return rt.Value{}, e359
		}
		t360, e361 := rt.FieldGet(ctx, svod, "журнал")
		if e361 != nil {
			return rt.Value{}, e361
		}
		t362, e363 := rt.FieldGet(ctx, svod, "сборка")
		if e363 != nil {
			return rt.Value{}, e363
		}
		t364, e365 := rt.FieldGet(ctx, svod, "загрузка")
		if e365 != nil {
			return rt.Value{}, e365
		}
		t366, e367 := rt.FieldGet(ctx, svod, "крупное")
		if e367 != nil {
			return rt.Value{}, e367
		}
		t368, e369 := rt.FieldGet(ctx, svod, "неизвестное")
		if e369 != nil {
			return rt.Value{}, e369
		}
		t370, e371 := rt.FieldGet(ctx, svod, "освободить")
		if e371 != nil {
			return rt.Value{}, e371
		}
		t372, e373 := rt.Add(ctx, t370, ubrat)
		if e373 != nil {
			return rt.Value{}, e373
		}
		t374 := make([]rt.Field, 7)
		t374[0] = rt.Field{Name: "кэш", Value: t358}
		t374[1] = rt.Field{Name: "журнал", Value: t360}
		t374[2] = rt.Field{Name: "сборка", Value: t362}
		t374[3] = rt.Field{Name: "загрузка", Value: t364}
		t374[4] = rt.Field{Name: "крупное", Value: t366}
		t374[5] = rt.Field{Name: "неизвестное", Value: t368}
		t374[6] = rt.Field{Name: "освободить", Value: t372}
		return rt.Record(t374), nil
	} else if rt.VariantIs(t352, "Журнал") {
		t375, e376 := rt.FieldGet(ctx, svod, "кэш")
		if e376 != nil {
			return rt.Value{}, e376
		}
		t377, e378 := rt.FieldGet(ctx, svod, "журнал")
		if e378 != nil {
			return rt.Value{}, e378
		}
		t379, e380 := rt.FieldGet(ctx, reshenie, "вес")
		if e380 != nil {
			return rt.Value{}, e380
		}
		t381, e382 := rt.Add(ctx, t377, t379)
		if e382 != nil {
			return rt.Value{}, e382
		}
		t383, e384 := rt.FieldGet(ctx, svod, "сборка")
		if e384 != nil {
			return rt.Value{}, e384
		}
		t385, e386 := rt.FieldGet(ctx, svod, "загрузка")
		if e386 != nil {
			return rt.Value{}, e386
		}
		t387, e388 := rt.FieldGet(ctx, svod, "крупное")
		if e388 != nil {
			return rt.Value{}, e388
		}
		t389, e390 := rt.FieldGet(ctx, svod, "неизвестное")
		if e390 != nil {
			return rt.Value{}, e390
		}
		t391, e392 := rt.FieldGet(ctx, svod, "освободить")
		if e392 != nil {
			return rt.Value{}, e392
		}
		t393, e394 := rt.Add(ctx, t391, ubrat)
		if e394 != nil {
			return rt.Value{}, e394
		}
		t395 := make([]rt.Field, 7)
		t395[0] = rt.Field{Name: "кэш", Value: t375}
		t395[1] = rt.Field{Name: "журнал", Value: t381}
		t395[2] = rt.Field{Name: "сборка", Value: t383}
		t395[3] = rt.Field{Name: "загрузка", Value: t385}
		t395[4] = rt.Field{Name: "крупное", Value: t387}
		t395[5] = rt.Field{Name: "неизвестное", Value: t389}
		t395[6] = rt.Field{Name: "освободить", Value: t393}
		return rt.Record(t395), nil
	} else if rt.VariantIs(t352, "Сборка") {
		t396, e397 := rt.FieldGet(ctx, svod, "кэш")
		if e397 != nil {
			return rt.Value{}, e397
		}
		t398, e399 := rt.FieldGet(ctx, svod, "журнал")
		if e399 != nil {
			return rt.Value{}, e399
		}
		t400, e401 := rt.FieldGet(ctx, svod, "сборка")
		if e401 != nil {
			return rt.Value{}, e401
		}
		t402, e403 := rt.FieldGet(ctx, reshenie, "вес")
		if e403 != nil {
			return rt.Value{}, e403
		}
		t404, e405 := rt.Add(ctx, t400, t402)
		if e405 != nil {
			return rt.Value{}, e405
		}
		t406, e407 := rt.FieldGet(ctx, svod, "загрузка")
		if e407 != nil {
			return rt.Value{}, e407
		}
		t408, e409 := rt.FieldGet(ctx, svod, "крупное")
		if e409 != nil {
			return rt.Value{}, e409
		}
		t410, e411 := rt.FieldGet(ctx, svod, "неизвестное")
		if e411 != nil {
			return rt.Value{}, e411
		}
		t412, e413 := rt.FieldGet(ctx, svod, "освободить")
		if e413 != nil {
			return rt.Value{}, e413
		}
		t414, e415 := rt.Add(ctx, t412, ubrat)
		if e415 != nil {
			return rt.Value{}, e415
		}
		t416 := make([]rt.Field, 7)
		t416[0] = rt.Field{Name: "кэш", Value: t396}
		t416[1] = rt.Field{Name: "журнал", Value: t398}
		t416[2] = rt.Field{Name: "сборка", Value: t404}
		t416[3] = rt.Field{Name: "загрузка", Value: t406}
		t416[4] = rt.Field{Name: "крупное", Value: t408}
		t416[5] = rt.Field{Name: "неизвестное", Value: t410}
		t416[6] = rt.Field{Name: "освободить", Value: t414}
		return rt.Record(t416), nil
	} else if rt.VariantIs(t352, "Загрузка") {
		t417, e418 := rt.FieldGet(ctx, svod, "кэш")
		if e418 != nil {
			return rt.Value{}, e418
		}
		t419, e420 := rt.FieldGet(ctx, svod, "журнал")
		if e420 != nil {
			return rt.Value{}, e420
		}
		t421, e422 := rt.FieldGet(ctx, svod, "сборка")
		if e422 != nil {
			return rt.Value{}, e422
		}
		t423, e424 := rt.FieldGet(ctx, svod, "загрузка")
		if e424 != nil {
			return rt.Value{}, e424
		}
		t425, e426 := rt.FieldGet(ctx, reshenie, "вес")
		if e426 != nil {
			return rt.Value{}, e426
		}
		t427, e428 := rt.Add(ctx, t423, t425)
		if e428 != nil {
			return rt.Value{}, e428
		}
		t429, e430 := rt.FieldGet(ctx, svod, "крупное")
		if e430 != nil {
			return rt.Value{}, e430
		}
		t431, e432 := rt.FieldGet(ctx, svod, "неизвестное")
		if e432 != nil {
			return rt.Value{}, e432
		}
		t433, e434 := rt.FieldGet(ctx, svod, "освободить")
		if e434 != nil {
			return rt.Value{}, e434
		}
		t435, e436 := rt.Add(ctx, t433, ubrat)
		if e436 != nil {
			return rt.Value{}, e436
		}
		t437 := make([]rt.Field, 7)
		t437[0] = rt.Field{Name: "кэш", Value: t417}
		t437[1] = rt.Field{Name: "журнал", Value: t419}
		t437[2] = rt.Field{Name: "сборка", Value: t421}
		t437[3] = rt.Field{Name: "загрузка", Value: t427}
		t437[4] = rt.Field{Name: "крупное", Value: t429}
		t437[5] = rt.Field{Name: "неизвестное", Value: t431}
		t437[6] = rt.Field{Name: "освободить", Value: t435}
		return rt.Record(t437), nil
	} else if rt.VariantIs(t352, "Крупное") {
		t438, e439 := rt.FieldGet(ctx, svod, "кэш")
		if e439 != nil {
			return rt.Value{}, e439
		}
		t440, e441 := rt.FieldGet(ctx, svod, "журнал")
		if e441 != nil {
			return rt.Value{}, e441
		}
		t442, e443 := rt.FieldGet(ctx, svod, "сборка")
		if e443 != nil {
			return rt.Value{}, e443
		}
		t444, e445 := rt.FieldGet(ctx, svod, "загрузка")
		if e445 != nil {
			return rt.Value{}, e445
		}
		t446, e447 := rt.FieldGet(ctx, svod, "крупное")
		if e447 != nil {
			return rt.Value{}, e447
		}
		t448, e449 := rt.FieldGet(ctx, reshenie, "вес")
		if e449 != nil {
			return rt.Value{}, e449
		}
		t450, e451 := rt.Add(ctx, t446, t448)
		if e451 != nil {
			return rt.Value{}, e451
		}
		t452, e453 := rt.FieldGet(ctx, svod, "неизвестное")
		if e453 != nil {
			return rt.Value{}, e453
		}
		t454, e455 := rt.FieldGet(ctx, svod, "освободить")
		if e455 != nil {
			return rt.Value{}, e455
		}
		t456, e457 := rt.Add(ctx, t454, ubrat)
		if e457 != nil {
			return rt.Value{}, e457
		}
		t458 := make([]rt.Field, 7)
		t458[0] = rt.Field{Name: "кэш", Value: t438}
		t458[1] = rt.Field{Name: "журнал", Value: t440}
		t458[2] = rt.Field{Name: "сборка", Value: t442}
		t458[3] = rt.Field{Name: "загрузка", Value: t444}
		t458[4] = rt.Field{Name: "крупное", Value: t450}
		t458[5] = rt.Field{Name: "неизвестное", Value: t452}
		t458[6] = rt.Field{Name: "освободить", Value: t456}
		return rt.Record(t458), nil
	} else if rt.VariantIs(t352, "Неизвестное") {
		t459, e460 := rt.FieldGet(ctx, svod, "кэш")
		if e460 != nil {
			return rt.Value{}, e460
		}
		t461, e462 := rt.FieldGet(ctx, svod, "журнал")
		if e462 != nil {
			return rt.Value{}, e462
		}
		t463, e464 := rt.FieldGet(ctx, svod, "сборка")
		if e464 != nil {
			return rt.Value{}, e464
		}
		t465, e466 := rt.FieldGet(ctx, svod, "загрузка")
		if e466 != nil {
			return rt.Value{}, e466
		}
		t467, e468 := rt.FieldGet(ctx, svod, "крупное")
		if e468 != nil {
			return rt.Value{}, e468
		}
		t469, e470 := rt.FieldGet(ctx, svod, "неизвестное")
		if e470 != nil {
			return rt.Value{}, e470
		}
		t471, e472 := rt.FieldGet(ctx, reshenie, "вес")
		if e472 != nil {
			return rt.Value{}, e472
		}
		t473, e474 := rt.Add(ctx, t469, t471)
		if e474 != nil {
			return rt.Value{}, e474
		}
		t475, e476 := rt.FieldGet(ctx, svod, "освободить")
		if e476 != nil {
			return rt.Value{}, e476
		}
		t477, e478 := rt.Add(ctx, t475, ubrat)
		if e478 != nil {
			return rt.Value{}, e478
		}
		t479 := make([]rt.Field, 7)
		t479[0] = rt.Field{Name: "кэш", Value: t459}
		t479[1] = rt.Field{Name: "журнал", Value: t461}
		t479[2] = rt.Field{Name: "сборка", Value: t463}
		t479[3] = rt.Field{Name: "загрузка", Value: t465}
		t479[4] = rt.Field{Name: "крупное", Value: t467}
		t479[5] = rt.Field{Name: "неизвестное", Value: t473}
		t479[6] = rt.Field{Name: "освободить", Value: t477}
		return rt.Record(t479), nil
	} else {
		return rt.Value{}, rt.MatchFail(ctx, t352)
	}
}

// Svesti — функция flang «Свести».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: «Свод».
func Svesti(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t480, e481 := ReshitVsyo(ctx, zapisi)
	if e481 != nil {
		return rt.Value{}, e481
	}
	t482, e483 := rt.RequireList(ctx, t480, "свёртка")
	if e483 != nil {
		return rt.Value{}, e483
	}
	t484, e485 := PustoySvod(ctx)
	if e485 != nil {
		return rt.Value{}, e485
	}
	// «свод»
	svod := t484
	for t486 := range t482 {
		// «решение»
		reshenie := t482[t486]
		t487, e488 := PribavitReshenie(ctx, svod, reshenie)
		if e488 != nil {
			return rt.Value{}, e488
		}
		svod = t487
	}
	t489 := svod
	t490, e491 := I2Derzhitsya(ctx, zapisi, t489)
	if e491 != nil {
		return rt.Value{}, e491
	}
	// постусловие «И2: освобождаемое не больше убираемого»
	t492, e493 := rt.Post(ctx, t490, "И2: освобождаемое не больше убираемого", "Свести")
	if e493 != nil {
		return rt.Value{}, e493
	}
	if !t492 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «И2: освобождаемое не больше убираемого» функции «Свести»")
	}
	return t489, nil
}

// SummaRazmerovUbiraemyh — функция flang «Сумма размеров убираемых».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр zapisi — «записи»: список: «Находка».
// Результат — значение: число.
func SummaRazmerovUbiraemyh(ctx *rt.Ctx, zapisi rt.Value) (rt.Value, error) {
	t494, e495 := rt.RequireList(ctx, zapisi, "свёртка")
	if e495 != nil {
		return rt.Value{}, e495
	}
	// «акк»
	akk := rt.Number(0.0)
	for t496 := range t494 {
		// «находка»
		nahodka := t494[t496]
		var t497 rt.Value
		t498, e499 := ReshitNahodku(ctx, nahodka)
		if e499 != nil {
			return rt.Value{}, e499
		}
		t500, e501 := rt.FieldGet(ctx, t498, "приговор")
		if e501 != nil {
			return rt.Value{}, e501
		}
		if rt.VariantIs(t500, "МожноУбрать") {
			t502, e503 := rt.FieldGet(ctx, nahodka, "размер")
			if e503 != nil {
				return rt.Value{}, e503
			}
			t504, e505 := rt.Add(ctx, akk, t502)
			if e505 != nil {
				return rt.Value{}, e505
			}
			t497 = t504
		} else if rt.VariantIs(t500, "Спросить") {
			t497 = akk
		} else if rt.VariantIs(t500, "НеТрогать") {
			t497 = akk
		} else {
			return rt.Value{}, rt.MatchFail(ctx, t500)
		}
		akk = t497
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
	t506, e507 := rt.FieldGet(ctx, svod, "освободить")
	if e507 != nil {
		return rt.Value{}, e507
	}
	t508, e509 := SummaRazmerovUbiraemyh(ctx, zapisi)
	if e509 != nil {
		return rt.Value{}, e509
	}
	t510, e511 := rt.Lte(ctx, t506, t508)
	if e511 != nil {
		return rt.Value{}, e511
	}
	return t510, nil
}

// StrokuOtchyota — функция flang «Строку отчёта».
//
// Тотальная: завершение доказано анализом завершаемости (totality.mjs).
//
// Параметр nahodka — «находка»: «Находка».
// Результат — значение: «Строка отчёта».
func StrokuOtchyota(ctx *rt.Ctx, nahodka rt.Value) (rt.Value, error) {
	t512, e513 := rt.FieldGet(ctx, nahodka, "путь")
	if e513 != nil {
		return rt.Value{}, e513
	}
	t514, e515 := ReshitNahodku(ctx, nahodka)
	if e515 != nil {
		return rt.Value{}, e515
	}
	t516 := make([]rt.Field, 2)
	t516[0] = rt.Field{Name: "путь", Value: t512}
	t516[1] = rt.Field{Name: "решение", Value: t514}
	return rt.Record(t516), nil
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
		t517 := make([]rt.Value, 1)
		t517[0] = stroka
		return rt.List(t517), nil
	} else if rt.ChainCons(stroki) {
		// голова «голова»
		golova := rt.ChainHead(stroki)
		// хвост «хвост»
		hvost := rt.ChainTail(stroki)
		t518, e519 := rt.FieldGet(ctx, stroka, "решение")
		if e519 != nil {
			return rt.Value{}, e519
		}
		t520, e521 := rt.FieldGet(ctx, t518, "вес")
		if e521 != nil {
			return rt.Value{}, e521
		}
		t522, e523 := rt.FieldGet(ctx, golova, "решение")
		if e523 != nil {
			return rt.Value{}, e523
		}
		t524, e525 := rt.FieldGet(ctx, t522, "вес")
		if e525 != nil {
			return rt.Value{}, e525
		}
		t526, e527 := rt.Gte(ctx, t520, t524)
		if e527 != nil {
			return rt.Value{}, e527
		}
		t528, e529 := rt.Cond(ctx, t526)
		if e529 != nil {
			return rt.Value{}, e529
		}
		if t528 {
			return PripisatStrokuOtchyota(ctx, stroka, stroki)
		} else {
			t530, e531 := VstavitPoVesu(ctx, stroka, hvost)
			if e531 != nil {
				return rt.Value{}, e531
			}
			return PripisatStrokuOtchyota(ctx, golova, t530)
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
	t532, e533 := rt.RequireList(ctx, stroki, "свёртка")
	if e533 != nil {
		return rt.Value{}, e533
	}
	t534 := make([]rt.Value, 1)
	t534[0] = pervaya
	// «акк»
	akk := rt.List(t534)
	for t535 := range t532 {
		// «эл»
		el := t532[t535]
		// «добавить»
		t536, e537 := rt.BAppend(ctx, el, akk)
		if e537 != nil {
			return rt.Value{}, e537
		}
		akk = t536
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
	t538, e539 := rt.RequireList(ctx, zapisi, "свёртка")
	if e539 != nil {
		return rt.Value{}, e539
	}
	// «акк»
	akk := rt.List(nil)
	for t540 := range t538 {
		// «находка»
		nahodka := t538[t540]
		t541, e542 := StrokuOtchyota(ctx, nahodka)
		if e542 != nil {
			return rt.Value{}, e542
		}
		t543, e544 := VstavitPoVesu(ctx, t541, akk)
		if e544 != nil {
			return rt.Value{}, e544
		}
		akk = t543
	}
	t545 := akk
	t546, e547 := OtchyotToyZheDliny(ctx, zapisi, t545)
	if e547 != nil {
		return rt.Value{}, e547
	}
	// постусловие «Отчёт той же длины»
	t548, e549 := rt.Post(ctx, t546, "Отчёт той же длины", "Отчёт")
	if e549 != nil {
		return rt.Value{}, e549
	}
	if !t548 {
		return rt.Value{}, rt.Fail("FLANG_PROPERTY", "%s", "нарушено свойство «Отчёт той же длины» функции «Отчёт»")
	}
	return t545, nil
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
	t550, e551 := rt.BLength(ctx, stroki)
	if e551 != nil {
		return rt.Value{}, e551
	}
	// «длина»
	t552, e553 := rt.BLength(ctx, zapisi)
	if e553 != nil {
		return rt.Value{}, e553
	}
	return rt.Flag(rt.Equal(t550, t552)), nil
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
	t554, e555 := EtoMozhnoUbrat(ctx, prigovor)
	if e555 != nil {
		return rt.Value{}, e555
	}
	t556, e557 := rt.Cond(ctx, t554)
	if e557 != nil {
		return rt.Value{}, e557
	}
	var t558 rt.Value
	if t556 {
		t558 = rt.Flag(false)
	} else {
		t558 = rt.Flag(true)
	}
	t559, e560 := rt.Cond(ctx, t558)
	if e560 != nil {
		return rt.Value{}, e560
	}
	if t559 {
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
	t561, e562 := rt.FieldGet(ctx, nahodka, "путь")
	if e562 != nil {
		return rt.Value{}, e562
	}
	t563, e564 := PrimetaKesha(ctx, t561)
	if e564 != nil {
		return rt.Value{}, e564
	}
	// пусть «кэш»
	kesh := t563
	t565, e566 := rt.FieldGet(ctx, nahodka, "путь")
	if e566 != nil {
		return rt.Value{}, e566
	}
	t567, e568 := PrimetaZhurnala(ctx, t565)
	if e568 != nil {
		return rt.Value{}, e568
	}
	// пусть «журнал»
	zhurnal := t567
	t569, e570 := rt.FieldGet(ctx, nahodka, "путь")
	if e570 != nil {
		return rt.Value{}, e570
	}
	t571, e572 := PrimetaSborki(ctx, t569)
	if e572 != nil {
		return rt.Value{}, e572
	}
	// пусть «сборка»
	sborka := t571
	t573, e574 := rt.FieldGet(ctx, nahodka, "путь")
	if e574 != nil {
		return rt.Value{}, e574
	}
	t575, e576 := PrimetaZagruzki(ctx, t573)
	if e576 != nil {
		return rt.Value{}, e576
	}
	// пусть «загрузка»
	zagruzka := t575
	t577, e578 := rt.FieldGet(ctx, nahodka, "размер")
	if e578 != nil {
		return rt.Value{}, e578
	}
	t579, e580 := PorogKrupnogo(ctx)
	if e580 != nil {
		return rt.Value{}, e580
	}
	t581, e582 := rt.Gte(ctx, t577, t579)
	if e582 != nil {
		return rt.Value{}, e582
	}
	// пусть «крупное»
	krupnoe := t581
	if rt.VariantIs(razryad, "Кэш") {
		return kesh, nil
	} else if rt.VariantIs(razryad, "Журнал") {
		t583, e584 := rt.Cond(ctx, zhurnal)
		if e584 != nil {
			return rt.Value{}, e584
		}
		if t583 {
			t585, e586 := rt.Cond(ctx, kesh)
			if e586 != nil {
				return rt.Value{}, e586
			}
			if t585 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Сборка") {
		t587, e588 := rt.Cond(ctx, sborka)
		if e588 != nil {
			return rt.Value{}, e588
		}
		var t589 rt.Value
		if t587 {
			t590, e591 := rt.Cond(ctx, kesh)
			if e591 != nil {
				return rt.Value{}, e591
			}
			var t592 rt.Value
			if t590 {
				t592 = rt.Flag(false)
			} else {
				t592 = rt.Flag(true)
			}
			t589 = t592
		} else {
			t589 = rt.Flag(false)
		}
		t593, e594 := rt.Cond(ctx, t589)
		if e594 != nil {
			return rt.Value{}, e594
		}
		if t593 {
			t595, e596 := rt.Cond(ctx, zhurnal)
			if e596 != nil {
				return rt.Value{}, e596
			}
			if t595 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Загрузка") {
		t597, e598 := rt.Cond(ctx, zagruzka)
		if e598 != nil {
			return rt.Value{}, e598
		}
		var t599 rt.Value
		if t597 {
			t600, e601 := rt.Cond(ctx, kesh)
			if e601 != nil {
				return rt.Value{}, e601
			}
			var t602 rt.Value
			if t600 {
				t602 = rt.Flag(false)
			} else {
				t602 = rt.Flag(true)
			}
			t599 = t602
		} else {
			t599 = rt.Flag(false)
		}
		t603, e604 := rt.Cond(ctx, t599)
		if e604 != nil {
			return rt.Value{}, e604
		}
		var t605 rt.Value
		if t603 {
			t606, e607 := rt.Cond(ctx, zhurnal)
			if e607 != nil {
				return rt.Value{}, e607
			}
			var t608 rt.Value
			if t606 {
				t608 = rt.Flag(false)
			} else {
				t608 = rt.Flag(true)
			}
			t605 = t608
		} else {
			t605 = rt.Flag(false)
		}
		t609, e610 := rt.Cond(ctx, t605)
		if e610 != nil {
			return rt.Value{}, e610
		}
		if t609 {
			t611, e612 := rt.Cond(ctx, sborka)
			if e612 != nil {
				return rt.Value{}, e612
			}
			if t611 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Крупное") {
		t613, e614 := rt.Cond(ctx, krupnoe)
		if e614 != nil {
			return rt.Value{}, e614
		}
		var t615 rt.Value
		if t613 {
			t616, e617 := rt.Cond(ctx, kesh)
			if e617 != nil {
				return rt.Value{}, e617
			}
			var t618 rt.Value
			if t616 {
				t618 = rt.Flag(false)
			} else {
				t618 = rt.Flag(true)
			}
			t615 = t618
		} else {
			t615 = rt.Flag(false)
		}
		t619, e620 := rt.Cond(ctx, t615)
		if e620 != nil {
			return rt.Value{}, e620
		}
		var t621 rt.Value
		if t619 {
			t622, e623 := rt.Cond(ctx, zhurnal)
			if e623 != nil {
				return rt.Value{}, e623
			}
			var t624 rt.Value
			if t622 {
				t624 = rt.Flag(false)
			} else {
				t624 = rt.Flag(true)
			}
			t621 = t624
		} else {
			t621 = rt.Flag(false)
		}
		t625, e626 := rt.Cond(ctx, t621)
		if e626 != nil {
			return rt.Value{}, e626
		}
		var t627 rt.Value
		if t625 {
			t628, e629 := rt.Cond(ctx, sborka)
			if e629 != nil {
				return rt.Value{}, e629
			}
			var t630 rt.Value
			if t628 {
				t630 = rt.Flag(false)
			} else {
				t630 = rt.Flag(true)
			}
			t627 = t630
		} else {
			t627 = rt.Flag(false)
		}
		t631, e632 := rt.Cond(ctx, t627)
		if e632 != nil {
			return rt.Value{}, e632
		}
		if t631 {
			t633, e634 := rt.Cond(ctx, zagruzka)
			if e634 != nil {
				return rt.Value{}, e634
			}
			if t633 {
				return rt.Flag(false), nil
			} else {
				return rt.Flag(true), nil
			}
		} else {
			return rt.Flag(false), nil
		}
	} else if rt.VariantIs(razryad, "Неизвестное") {
		t635, e636 := rt.Cond(ctx, kesh)
		if e636 != nil {
			return rt.Value{}, e636
		}
		var t637 rt.Value
		if t635 {
			t637 = rt.Flag(false)
		} else {
			t637 = rt.Flag(true)
		}
		t638, e639 := rt.Cond(ctx, t637)
		if e639 != nil {
			return rt.Value{}, e639
		}
		var t640 rt.Value
		if t638 {
			t641, e642 := rt.Cond(ctx, zhurnal)
			if e642 != nil {
				return rt.Value{}, e642
			}
			var t643 rt.Value
			if t641 {
				t643 = rt.Flag(false)
			} else {
				t643 = rt.Flag(true)
			}
			t640 = t643
		} else {
			t640 = rt.Flag(false)
		}
		t644, e645 := rt.Cond(ctx, t640)
		if e645 != nil {
			return rt.Value{}, e645
		}
		var t646 rt.Value
		if t644 {
			t647, e648 := rt.Cond(ctx, sborka)
			if e648 != nil {
				return rt.Value{}, e648
			}
			var t649 rt.Value
			if t647 {
				t649 = rt.Flag(false)
			} else {
				t649 = rt.Flag(true)
			}
			t646 = t649
		} else {
			t646 = rt.Flag(false)
		}
		t650, e651 := rt.Cond(ctx, t646)
		if e651 != nil {
			return rt.Value{}, e651
		}
		var t652 rt.Value
		if t650 {
			t653, e654 := rt.Cond(ctx, zagruzka)
			if e654 != nil {
				return rt.Value{}, e654
			}
			var t655 rt.Value
			if t653 {
				t655 = rt.Flag(false)
			} else {
				t655 = rt.Flag(true)
			}
			t652 = t655
		} else {
			t652 = rt.Flag(false)
		}
		t656, e657 := rt.Cond(ctx, t652)
		if e657 != nil {
			return rt.Value{}, e657
		}
		if t656 {
			t658, e659 := rt.Cond(ctx, krupnoe)
			if e659 != nil {
				return rt.Value{}, e659
			}
			if t658 {
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
	t660, e661 := rt.FieldGet(ctx, nahodka, "доступен")
	if e661 != nil {
		return rt.Value{}, e661
	}
	t662, e663 := rt.Cond(ctx, t660)
	if e663 != nil {
		return rt.Value{}, e663
	}
	var t664 rt.Value
	if t662 {
		t664 = rt.Flag(false)
	} else {
		t664 = rt.Flag(true)
	}
	t665, e666 := rt.Cond(ctx, t664)
	if e666 != nil {
		return rt.Value{}, e666
	}
	if t665 {
		return EtoNeTrogat(ctx, prigovor)
	} else {
		t667, e668 := Ssylka(ctx, nahodka)
		if e668 != nil {
			return rt.Value{}, e668
		}
		t669, e670 := rt.Cond(ctx, t667)
		if e670 != nil {
			return rt.Value{}, e670
		}
		if t669 {
			return EtoNeTrogat(ctx, prigovor)
		} else {
			t671, e672 := EtoMozhnoUbrat(ctx, prigovor)
			if e672 != nil {
				return rt.Value{}, e672
			}
			t673, e674 := rt.Cond(ctx, t671)
			if e674 != nil {
				return rt.Value{}, e674
			}
			if t673 {
				t675, e676 := I1NaPare(ctx, razryad, prigovor)
				if e676 != nil {
					return rt.Value{}, e676
				}
				t677, e678 := rt.Cond(ctx, t675)
				if e678 != nil {
					return rt.Value{}, e678
				}
				var t679 rt.Value
				if t677 {
					t680, e681 := Katalog(ctx, nahodka)
					if e681 != nil {
						return rt.Value{}, e681
					}
					t682, e683 := rt.Cond(ctx, t680)
					if e683 != nil {
						return rt.Value{}, e683
					}
					var t684 rt.Value
					if t682 {
						t684 = rt.Flag(false)
					} else {
						t684 = rt.Flag(true)
					}
					t679 = t684
				} else {
					t679 = rt.Flag(false)
				}
				t685, e686 := rt.Cond(ctx, t679)
				if e686 != nil {
					return rt.Value{}, e686
				}
				if t685 {
					t687, e688 := rt.FieldGet(ctx, nahodka, "возраст_дней")
					if e688 != nil {
						return rt.Value{}, e688
					}
					t689, e690 := PorogRazryada(ctx, razryad)
					if e690 != nil {
						return rt.Value{}, e690
					}
					t691, e692 := rt.Gte(ctx, t687, t689)
					if e692 != nil {
						return rt.Value{}, e692
					}
					return t691, nil
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
	t693, e694 := EtoNeTrogat(ctx, prigovor)
	if e694 != nil {
		return rt.Value{}, e694
	}
	t695, e696 := rt.Cond(ctx, t693)
	if e696 != nil {
		return rt.Value{}, e696
	}
	if t695 {
		return rt.Flag(rt.Equal(ves, rt.Number(0.0))), nil
	} else {
		t697, e698 := rt.FieldGet(ctx, nahodka, "размер")
		if e698 != nil {
			return rt.Value{}, e698
		}
		t699, e700 := rt.Cond(ctx, rt.Flag(rt.Equal(ves, t697)))
		if e700 != nil {
			return rt.Value{}, e700
		}
		if t699 {
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
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Есть примета", 1, len(args))
		}
		return EstPrimeta(ctx, args[0])
	case "Разряд решён размером":
		if len(args) != 1 {
			return rt.Value{}, rt.Fail(rt.CodeType,
				"функция «%s» принимает %d аргум., получено %d",
				"Разряд решён размером", 1, len(args))
		}
		return RazryadReshyonRazmerom(ctx, args[0])
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
