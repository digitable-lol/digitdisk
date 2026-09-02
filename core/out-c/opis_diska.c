/*
 * Сгенерировано flang (бэкенд C, flang/src/emit/c.mjs). Не редактировать руками.
 * Модуль flang: «Опись диска».
 * Файл: реализация.
 * Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.
 */
#include "opis_diska.h"

#include <string.h>

/* Константы программы: имена полей и строковые литералы. */
static const char *const opis_diska_names_1[] = { "путь", "размер", "возраст_дней", "вид", "доступен" };
static const char *const opis_diska_names_2[] = { "разряд", "якорь", "цепь" };
static const char *const opis_diska_names_3[] = { "разряд", "приговор", "вес" };
static const char *const opis_diska_names_4[] = { "путь", "решение" };
static const char *const opis_diska_names_5[] = { "кэш", "журнал", "сборка", "загрузка", "крупное", "неизвестное", "освободить" };
static const fl_value opis_diska_text_6 = { FL_STRING, { .string = { "/", 1, 1 } } };
static const fl_value opis_diska_text_7 = { FL_STRING, { .string = { "", 0, 0 } } };
static const fl_value opis_diska_text_8 = { FL_STRING, { .string = { "0123456789abcdefABCDEF", 22, 22 } } };
static const fl_value opis_diska_text_9 = { FL_STRING, { .string = { "site-packages", 13, 13 } } };
static const fl_value opis_diska_text_10 = { FL_STRING, { .string = { "dist-packages", 13, 13 } } };
static const fl_value opis_diska_text_11 = { FL_STRING, { .string = { ".git", 4, 4 } } };
static const fl_value opis_diska_text_12 = { FL_STRING, { .string = { "/tmp/", 5, 5 } } };
static const fl_value opis_diska_text_13 = { FL_STRING, { .string = { "/var/tmp/", 9, 9 } } };
static const fl_value opis_diska_text_14 = { FL_STRING, { .string = { ".cache", 6, 6 } } };
static const fl_value opis_diska_text_15 = { FL_STRING, { .string = { "cache", 5, 5 } } };
static const fl_value opis_diska_text_16 = { FL_STRING, { .string = { "Caches", 6, 6 } } };
static const fl_value opis_diska_text_17 = { FL_STRING, { .string = { ".log", 4, 4 } } };
static const fl_value opis_diska_text_18 = { FL_STRING, { .string = { "log", 3, 3 } } };
static const fl_value opis_diska_text_19 = { FL_STRING, { .string = { "logs", 4, 4 } } };
static const fl_value opis_diska_text_20 = { FL_STRING, { .string = { "node_modules", 12, 12 } } };
static const fl_value opis_diska_text_21 = { FL_STRING, { .string = { "target", 6, 6 } } };
static const fl_value opis_diska_text_22 = { FL_STRING, { .string = { "build", 5, 5 } } };
static const fl_value opis_diska_text_23 = { FL_STRING, { .string = { "_build", 6, 6 } } };
static const fl_value opis_diska_text_24 = { FL_STRING, { .string = { ".gradle", 7, 7 } } };
static const fl_value opis_diska_text_25 = { FL_STRING, { .string = { "Downloads", 9, 9 } } };
static const fl_value opis_diska_text_26 = { FL_STRING, { .string = { "Загрузки", 16, 8 } } };


/* Фабрика записи FTS «Находка». */
fl_status opis_diska_sozdat_nahodka(fl_ctx *ctx, fl_value put, fl_value razmer, fl_value vozrast_dney, fl_value vid, fl_value dostupen, fl_value *out, fl_error *error) {
  fl_value values[5];
  values[0] = put; /* «путь» */
  values[1] = razmer; /* «размер» */
  values[2] = vozrast_dney; /* «возраст_дней» */
  values[3] = vid; /* «вид» */
  values[4] = dostupen; /* «доступен» */
  return fl_record_new(ctx, opis_diska_names_1, values, 5, out, error);
}

/* Фабрика записи FTS «Место». */
fl_status opis_diska_sozdat_mesto(fl_ctx *ctx, fl_value razryad, fl_value yakor, fl_value cep, fl_value *out, fl_error *error) {
  fl_value values[3];
  values[0] = razryad; /* «разряд» */
  values[1] = yakor; /* «якорь» */
  values[2] = cep; /* «цепь» */
  return fl_record_new(ctx, opis_diska_names_2, values, 3, out, error);
}

/* Фабрика записи FTS «Решение». */
fl_status opis_diska_sozdat_reshenie(fl_ctx *ctx, fl_value razryad, fl_value prigovor, fl_value ves, fl_value *out, fl_error *error) {
  fl_value values[3];
  values[0] = razryad; /* «разряд» */
  values[1] = prigovor; /* «приговор» */
  values[2] = ves; /* «вес» */
  return fl_record_new(ctx, opis_diska_names_3, values, 3, out, error);
}

/* Фабрика записи FTS «Строка отчёта». */
fl_status opis_diska_sozdat_stroka_otchyota(fl_ctx *ctx, fl_value put, fl_value reshenie, fl_value *out, fl_error *error) {
  fl_value values[2];
  values[0] = put; /* «путь» */
  values[1] = reshenie; /* «решение» */
  return fl_record_new(ctx, opis_diska_names_4, values, 2, out, error);
}

/* Фабрика записи FTS «Свод». */
fl_status opis_diska_sozdat_svod(fl_ctx *ctx, fl_value kesh, fl_value zhurnal, fl_value sborka, fl_value zagruzka, fl_value krupnoe, fl_value neizvestnoe, fl_value osvobodit, fl_value *out, fl_error *error) {
  fl_value values[7];
  values[0] = kesh; /* «кэш» */
  values[1] = zhurnal; /* «журнал» */
  values[2] = sborka; /* «сборка» */
  values[3] = zagruzka; /* «загрузка» */
  values[4] = krupnoe; /* «крупное» */
  values[5] = neizvestnoe; /* «неизвестное» */
  values[6] = osvobodit; /* «освободить» */
  return fl_record_new(ctx, opis_diska_names_5, values, 7, out, error);
}

/* Конструктор варианта «Файл» суммы «Вид». */
fl_status opis_diska_variant_fayl(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Файл", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Каталог» суммы «Вид». */
fl_status opis_diska_variant_katalog(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Каталог", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Ссылка» суммы «Вид». */
fl_status opis_diska_variant_ssylka(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Ссылка", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Кэш» суммы «Разряд». */
fl_status opis_diska_variant_kesh(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Кэш", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Журнал» суммы «Разряд». */
fl_status opis_diska_variant_zhurnal(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Журнал", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Сборка» суммы «Разряд». */
fl_status opis_diska_variant_sborka(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Сборка", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Загрузка» суммы «Разряд». */
fl_status opis_diska_variant_zagruzka(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Загрузка", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Крупное» суммы «Разряд». */
fl_status opis_diska_variant_krupnoe(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Крупное", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Неизвестное» суммы «Разряд». */
fl_status opis_diska_variant_neizvestnoe(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Неизвестное", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «МожноУбрать» суммы «Приговор». */
fl_status opis_diska_variant_mozhnoubrat(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "МожноУбрать", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «Спросить» суммы «Приговор». */
fl_status opis_diska_variant_sprosit(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "Спросить", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «НеТрогать» суммы «Приговор». */
fl_status opis_diska_variant_netrogat(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «ОтКорня» суммы «Якорь». */
fl_status opis_diska_variant_otkornya(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "ОтКорня", NULL, NULL, 0, out, error);
}

/* Конструктор варианта «ГдеУгодно» суммы «Якорь». */
fl_status opis_diska_variant_gdeugodno(fl_ctx *ctx, fl_value *out, fl_error *error) {
  return fl_variant_new(ctx, "ГдеУгодно", NULL, NULL, 0, out, error);
}

/*
 * Функция flang «Порог крупного».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_krupnogo(fl_ctx *ctx, fl_value *result, fl_error *error) {
  (void)ctx;
  (void)error;
  *result = fl_number(1073741824.0);
  return FL_OK;
}

/*
 * Функция flang «Порог кэша».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_kesha(fl_ctx *ctx, fl_value *result, fl_error *error) {
  (void)ctx;
  (void)error;
  *result = fl_number(7.0);
  return FL_OK;
}

/*
 * Функция flang «Порог журнала».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_zhurnala(fl_ctx *ctx, fl_value *result, fl_error *error) {
  (void)ctx;
  (void)error;
  *result = fl_number(30.0);
  return FL_OK;
}

/*
 * Функция flang «Порог загрузки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_zagruzki(fl_ctx *ctx, fl_value *result, fl_error *error) {
  (void)ctx;
  (void)error;
  *result = fl_number(180.0);
  return FL_OK;
}

/*
 * Функция flang «Составляющие пути».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение: список: строка
 */
fl_status opis_diska_sostavlyayuschie_puti(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t1 = fl_nothing(); /* «разделить» */
  FL_TRY(fl_b_razdelit_dokazano(ctx, put, opis_diska_text_6, &fl_t1, error));
  *result = fl_t1;
  return FL_OK;
}

/*
 * Функция flang «Имя в пути».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение: строка
 */
fl_status opis_diska_imya_v_puti(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t2 = fl_nothing();
  FL_TRY(opis_diska_sostavlyayuschie_puti(ctx, put, &fl_t2, error));
  fl_value fl_t3 = fl_nothing();
  FL_TRY(fl_require_list(ctx, fl_t2, "свёртка", &fl_t3, error));
  fl_value sobrano = opis_diska_text_7; /* «собрано» */
  const fl_mark fl_t5 = fl_region_open(ctx);
  for (size_t fl_t4 = 0; fl_t4 < fl_t3.as.list.count; fl_t4 += 1) {
    const fl_value chast = fl_t3.as.list.items[fl_t4]; /* «часть» */
    sobrano = chast;
    FL_TRY(fl_region_recycle(ctx, fl_t5, &sobrano, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t5, FL_OK, &sobrano, error));
  *result = sobrano;
  return FL_OK;
}

/*
 * Функция flang «Есть составляющая».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param imya — «имя»: строка
 * @return значение
 */
fl_status opis_diska_est_sostavlyayuschaya(fl_ctx *ctx, fl_value put, fl_value imya, fl_value *result, fl_error *error) {
  fl_value fl_t6 = fl_nothing();
  FL_TRY(opis_diska_sostavlyayuschie_puti(ctx, put, &fl_t6, error));
  fl_value fl_t7 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, fl_t6, imya, &fl_t7, error));
  *result = fl_t7;
  return FL_OK;
}

/*
 * Функция flang «Оканчивается на».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param tekst — «текст»: строка
 * @param hvost — «хвост»: строка
 * @return значение
 */
fl_status opis_diska_okanchivaetsya_na(fl_ctx *ctx, fl_value tekst, fl_value hvost, fl_value *result, fl_error *error) {
  fl_value fl_t8 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, tekst, &fl_t8, error));
  fl_value fl_t9 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, hvost, &fl_t9, error));
  if (fl_t8.tag != FL_NUMBER || fl_t9.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t8, fl_t9, error));
  bool fl_t10 = false;
  FL_TRY(fl_cond(ctx, fl_flag(fl_t8.as.number < fl_t9.as.number), &fl_t10, error));
  if (fl_t10) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    fl_value fl_t11 = fl_nothing(); /* «длина» */
    FL_TRY(fl_b_dlina(ctx, tekst, &fl_t11, error));
    fl_value fl_t12 = fl_nothing(); /* «длина» */
    FL_TRY(fl_b_dlina(ctx, hvost, &fl_t12, error));
    if (fl_t11.tag != FL_NUMBER || fl_t12.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "sub", fl_t11, fl_t12, error));
    fl_value fl_t13 = fl_nothing(); /* «длина» */
    FL_TRY(fl_b_dlina(ctx, tekst, &fl_t13, error));
    fl_value fl_t14 = fl_nothing(); /* «подстрока» */
    FL_TRY(fl_b_podstroka(ctx, tekst, fl_number((fl_t11.as.number - fl_t12.as.number) + 1.0), fl_t13, &fl_t14, error));
    *result = fl_flag(fl_equal(fl_t14, hvost));
    return FL_OK;
  }
}

/*
 * Функция flang «След пути».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение: строка
 */
fl_status opis_diska_sled_puti(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t15 = fl_nothing();
  FL_TRY(fl_concat(ctx, put, opis_diska_text_6, &fl_t15, error));
  *result = fl_t15;
  return FL_OK;
}

/*
 * Функция flang «Цепь ограничена».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param cep — «цепь»: строка
 * @return значение
 */
fl_status opis_diska_cep_ogranichena(fl_ctx *ctx, fl_value cep, fl_value *result, fl_error *error) {
  fl_value fl_t16 = fl_nothing(); /* «начинается с» */
  FL_TRY(fl_b_nachinaetsya_s(ctx, cep, opis_diska_text_6, &fl_t16, error));
  bool fl_t17 = false;
  FL_TRY(fl_cond(ctx, fl_t16, &fl_t17, error));
  fl_value fl_t18 = fl_nothing();
  if (fl_t17) {
    fl_value fl_t19 = fl_nothing();
    FL_TRY(opis_diska_okanchivaetsya_na(ctx, cep, opis_diska_text_6, &fl_t19, error));
    fl_t18 = fl_t19;
  } else {
    fl_t18 = fl_flag(false);
  }
  bool fl_t20 = false;
  FL_TRY(fl_cond(ctx, fl_t18, &fl_t20, error));
  if (fl_t20) {
    fl_value fl_t21 = fl_nothing(); /* «длина» */
    FL_TRY(fl_b_dlina(ctx, cep, &fl_t21, error));
    if (fl_t21.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t21, fl_number(1.0), error));
    *result = fl_flag(fl_t21.as.number > 1.0);
    return FL_OK;
  } else {
    *result = fl_flag(false);
    return FL_OK;
  }
}

/*
 * Функция flang «Справочник ограничен».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение
 */
fl_status opis_diska_spravochnik_ogranichen(fl_ctx *ctx, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t22 = fl_nothing();
  FL_TRY(fl_require_list(ctx, spravochnik, "свёртка", &fl_t22, error));
  fl_value akk = fl_flag(true); /* «акк» */
  const fl_mark fl_t24 = fl_region_open(ctx);
  for (size_t fl_t23 = 0; fl_t23 < fl_t22.as.list.count; fl_t23 += 1) {
    const fl_value mesto = fl_t22.as.list.items[fl_t23]; /* «место» */
    bool fl_t25 = false;
    FL_TRY(fl_cond(ctx, akk, &fl_t25, error));
    fl_value fl_t26 = fl_nothing();
    if (fl_t25) {
      fl_value fl_t27 = fl_nothing();
      FL_TRY(fl_field_get(ctx, mesto, "цепь", &fl_t27, error));
      fl_value fl_t28 = fl_nothing();
      FL_TRY(opis_diska_cep_ogranichena(ctx, fl_t27, &fl_t28, error));
      fl_t26 = fl_t28;
    } else {
      fl_t26 = fl_flag(false);
    }
    akk = fl_t26;
    FL_TRY(fl_region_recycle(ctx, fl_t24, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t24, FL_OK, &akk, error));
  *result = akk;
  return FL_OK;
}

/*
 * Функция flang «Разряд места допустим».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_mesta_dopustim(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error) {
  if (fl_variant_is(razryad, "Кэш")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Сборка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Загрузка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Крупное")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Неизвестное")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Место подходит».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param sled — «след»: строка
 * @param mesto — «место»: «Место»
 * @return значение
 */
fl_status opis_diska_mesto_podhodit(fl_ctx *ctx, fl_value sled, fl_value mesto, fl_value *result, fl_error *error) {
  fl_value fl_t29 = fl_nothing();
  FL_TRY(fl_field_get(ctx, mesto, "якорь", &fl_t29, error));
  if (fl_variant_is(fl_t29, "ОтКорня")) {
    fl_value fl_t30 = fl_nothing();
    FL_TRY(fl_field_get(ctx, mesto, "цепь", &fl_t30, error));
    fl_value fl_t31 = fl_nothing(); /* «начинается с» */
    FL_TRY(fl_b_nachinaetsya_s(ctx, sled, fl_t30, &fl_t31, error));
    *result = fl_t31;
    return FL_OK;
  } else if (fl_variant_is(fl_t29, "ГдеУгодно")) {
    fl_value fl_t32 = fl_nothing();
    FL_TRY(fl_field_get(ctx, mesto, "цепь", &fl_t32, error));
    fl_value fl_t33 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, sled, fl_t32, &fl_t33, error));
    *result = fl_t33;
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t29, error);
  }
}

/*
 * Функция flang «Это неизвестное».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_eto_neizvestnoe(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error) {
  if (fl_variant_is(razryad, "Неизвестное")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Кэш")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Сборка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Загрузка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Крупное")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Номер разряда».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение: число
 */
fl_status opis_diska_nomer_razryada(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error) {
  if (fl_variant_is(razryad, "Кэш")) {
    *result = fl_number(1.0);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    *result = fl_number(2.0);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Сборка")) {
    *result = fl_number(3.0);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Загрузка")) {
    *result = fl_number(4.0);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Крупное")) {
    *result = fl_number(5.0);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Неизвестное")) {
    *result = fl_number(6.0);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Тот же разряд».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param pervyy — «первый»: «Разряд»
 * @param vtoroy — «второй»: «Разряд»
 * @return значение
 */
fl_status opis_diska_tot_zhe_razryad(fl_ctx *ctx, fl_value pervyy, fl_value vtoroy, fl_value *result, fl_error *error) {
  fl_value fl_t34 = fl_nothing();
  FL_TRY(opis_diska_nomer_razryada(ctx, pervyy, &fl_t34, error));
  fl_value fl_t35 = fl_nothing();
  FL_TRY(opis_diska_nomer_razryada(ctx, vtoroy, &fl_t35, error));
  *result = fl_flag(fl_equal(fl_t34, fl_t35));
  return FL_OK;
}

/*
 * Функция flang «Место обосновано».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param spravochnik — «справочник»: список: «Место»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_mesto_obosnovano(fl_ctx *ctx, fl_value put, fl_value spravochnik, fl_value razryad, fl_value *result, fl_error *error) {
  fl_value fl_t36 = fl_nothing();
  FL_TRY(opis_diska_eto_neizvestnoe(ctx, razryad, &fl_t36, error));
  bool fl_t37 = false;
  FL_TRY(fl_cond(ctx, fl_t36, &fl_t37, error));
  if (fl_t37) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t38 = fl_nothing();
    FL_TRY(opis_diska_razryad_mesta_dopustim(ctx, razryad, &fl_t38, error));
    bool fl_t39 = false;
    FL_TRY(fl_cond(ctx, fl_t38, &fl_t39, error));
    if (fl_t39) {
      fl_value fl_t40 = fl_nothing();
      FL_TRY(fl_require_list(ctx, spravochnik, "отфильтровать", &fl_t40, error));
      fl_value *fl_t41 = NULL;
      size_t fl_t42 = 0;
      FL_TRY(fl_list_alloc(ctx, fl_t40.as.list.count, &fl_t41, error));
      for (size_t fl_t43 = 0; fl_t43 < fl_t40.as.list.count; fl_t43 += 1) {
        const fl_value mesto = fl_t40.as.list.items[fl_t43]; /* «место» */
        fl_value fl_t44 = fl_nothing();
        FL_TRY(opis_diska_sled_puti(ctx, put, &fl_t44, error));
        fl_value fl_t45 = fl_nothing();
        FL_TRY(opis_diska_mesto_podhodit(ctx, fl_t44, mesto, &fl_t45, error));
        bool fl_t46 = false;
        FL_TRY(fl_cond(ctx, fl_t45, &fl_t46, error));
        fl_value fl_t47 = fl_nothing();
        if (fl_t46) {
          fl_value fl_t48 = fl_nothing();
          FL_TRY(fl_field_get(ctx, mesto, "разряд", &fl_t48, error));
          fl_value fl_t49 = fl_nothing();
          FL_TRY(opis_diska_tot_zhe_razryad(ctx, fl_t48, razryad, &fl_t49, error));
          fl_t47 = fl_t49;
        } else {
          fl_t47 = fl_flag(false);
        }
        bool fl_t50 = false;
        FL_TRY(fl_keep(ctx, fl_t47, &fl_t50, error));
        if (fl_t50) {
          fl_t41[fl_t42] = mesto;
          fl_t42 += 1;
        }
      }
      fl_value fl_t51 = fl_nothing(); /* «длина» */
      FL_TRY(fl_b_dlina(ctx, fl_list(fl_t41, fl_t42), &fl_t51, error));
      if (fl_t51.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t51, fl_number(0.0), error));
      *result = fl_flag(fl_t51.as.number > 0.0);
      return FL_OK;
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  }
}

/*
 * Функция flang «Разряд по справочнику».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Разряд»
 */
fl_status opis_diska_razryad_po_spravochniku(fl_ctx *ctx, fl_value put, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t52 = fl_nothing();
  FL_TRY(opis_diska_sled_puti(ctx, put, &fl_t52, error));
  const fl_value sled = fl_t52; /* пусть «след» */
  fl_value fl_t53 = fl_nothing();
  FL_TRY(fl_require_list(ctx, spravochnik, "свёртка", &fl_t53, error));
  fl_value fl_t54 = fl_nothing();
  FL_TRY(fl_variant_new(ctx, "Неизвестное", NULL, NULL, 0, &fl_t54, error));
  fl_value naydeno = fl_t54; /* «найдено» */
  const fl_mark fl_t56 = fl_region_open(ctx);
  for (size_t fl_t55 = 0; fl_t55 < fl_t53.as.list.count; fl_t55 += 1) {
    const fl_value mesto = fl_t53.as.list.items[fl_t55]; /* «место» */
    fl_value fl_t57 = fl_nothing();
    FL_TRY(opis_diska_eto_neizvestnoe(ctx, naydeno, &fl_t57, error));
    bool fl_t58 = false;
    FL_TRY(fl_cond(ctx, fl_t57, &fl_t58, error));
    fl_value fl_t59 = fl_nothing();
    if (fl_t58) {
      fl_value fl_t60 = fl_nothing();
      FL_TRY(opis_diska_mesto_podhodit(ctx, sled, mesto, &fl_t60, error));
      fl_t59 = fl_t60;
    } else {
      fl_t59 = fl_flag(false);
    }
    bool fl_t61 = false;
    FL_TRY(fl_cond(ctx, fl_t59, &fl_t61, error));
    fl_value fl_t62 = fl_nothing();
    if (fl_t61) {
      fl_value fl_t63 = fl_nothing();
      FL_TRY(fl_field_get(ctx, mesto, "разряд", &fl_t63, error));
      fl_t62 = fl_t63;
    } else {
      fl_t62 = naydeno;
    }
    naydeno = fl_t62;
    FL_TRY(fl_region_recycle(ctx, fl_t56, &naydeno, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t56, FL_OK, &naydeno, error));
  const fl_value fl_t64 = naydeno;
  fl_value fl_t65 = fl_nothing();
  FL_TRY(opis_diska_mesto_obosnovano(ctx, put, spravochnik, fl_t64, &fl_t65, error));
  /* постусловие «Место обосновано записью справочника» */
  bool fl_t66 = false;
  FL_TRY(fl_post(ctx, fl_t65, "Место обосновано записью справочника", "Разряд по справочнику", &fl_t66, error));
  if (!fl_t66) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Место обосновано записью справочника» функции «Разряд по справочнику»");
  }
  *result = fl_t64;
  return FL_OK;
}

/*
 * Функция flang «Шестнадцатеричный знак».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param znak — «знак»: строка
 * @return значение
 */
fl_status opis_diska_shestnadcaterichnyy_znak(fl_ctx *ctx, fl_value znak, fl_value *result, fl_error *error) {
  fl_value fl_t67 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, opis_diska_text_8, znak, &fl_t67, error));
  *result = fl_t67;
  return FL_OK;
}

/*
 * Функция flang «Похоже на отпечаток».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param chast — «часть»: строка
 * @return значение
 */
fl_status opis_diska_pohozhe_na_otpechatok(fl_ctx *ctx, fl_value chast, fl_value *result, fl_error *error) {
  fl_value fl_t68 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, chast, &fl_t68, error));
  if (fl_t68.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t68, fl_number(32.0), error));
  bool fl_t69 = false;
  FL_TRY(fl_cond(ctx, fl_flag(fl_t68.as.number < 32.0), &fl_t69, error));
  if (fl_t69) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    fl_value fl_t70 = fl_nothing(); /* «символы» */
    FL_TRY(fl_b_simvoly(ctx, chast, &fl_t70, error));
    fl_value fl_t71 = fl_nothing();
    FL_TRY(fl_require_list(ctx, fl_t70, "свёртка", &fl_t71, error));
    fl_value sobrano = fl_flag(true); /* «собрано» */
    const fl_mark fl_t73 = fl_region_open(ctx);
    for (size_t fl_t72 = 0; fl_t72 < fl_t71.as.list.count; fl_t72 += 1) {
      const fl_value znak = fl_t71.as.list.items[fl_t72]; /* «знак» */
      bool fl_t74 = false;
      FL_TRY(fl_cond(ctx, sobrano, &fl_t74, error));
      fl_value fl_t75 = fl_nothing();
      if (fl_t74) {
        fl_value fl_t76 = fl_nothing();
        FL_TRY(opis_diska_shestnadcaterichnyy_znak(ctx, znak, &fl_t76, error));
        fl_t75 = fl_t76;
      } else {
        fl_t75 = fl_flag(false);
      }
      sobrano = fl_t75;
      FL_TRY(fl_region_recycle(ctx, fl_t73, &sobrano, error));
    }
    FL_TRY(fl_region_close(ctx, fl_t73, FL_OK, &sobrano, error));
    *result = sobrano;
    return FL_OK;
  }
}

/*
 * Функция flang «Адресуется содержимым».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_adresuetsya_soderzhimym(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t77 = fl_nothing();
  FL_TRY(opis_diska_sostavlyayuschie_puti(ctx, put, &fl_t77, error));
  fl_value fl_t78 = fl_nothing();
  FL_TRY(fl_require_list(ctx, fl_t77, "отфильтровать", &fl_t78, error));
  fl_value *fl_t79 = NULL;
  size_t fl_t80 = 0;
  FL_TRY(fl_list_alloc(ctx, fl_t78.as.list.count, &fl_t79, error));
  for (size_t fl_t81 = 0; fl_t81 < fl_t78.as.list.count; fl_t81 += 1) {
    const fl_value chast = fl_t78.as.list.items[fl_t81]; /* «часть» */
    fl_value fl_t82 = fl_nothing();
    FL_TRY(opis_diska_pohozhe_na_otpechatok(ctx, chast, &fl_t82, error));
    bool fl_t83 = false;
    FL_TRY(fl_keep(ctx, fl_t82, &fl_t83, error));
    if (fl_t83) {
      fl_t79[fl_t80] = chast;
      fl_t80 += 1;
    }
  }
  fl_value fl_t84 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, fl_list(fl_t79, fl_t80), &fl_t84, error));
  if (fl_t84.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t84, fl_number(0.0), error));
  bool fl_t85 = false;
  FL_TRY(fl_cond(ctx, fl_flag(fl_t84.as.number > 0.0), &fl_t85, error));
  fl_value fl_t86 = fl_nothing();
  if (fl_t85) {
    fl_t86 = fl_flag(true);
  } else {
    fl_value fl_t87 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_9, &fl_t87, error));
    fl_t86 = fl_t87;
  }
  bool fl_t88 = false;
  FL_TRY(fl_cond(ctx, fl_t86, &fl_t88, error));
  fl_value fl_t89 = fl_nothing();
  if (fl_t88) {
    fl_t89 = fl_flag(true);
  } else {
    fl_value fl_t90 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_10, &fl_t90, error));
    fl_t89 = fl_t90;
  }
  bool fl_t91 = false;
  FL_TRY(fl_cond(ctx, fl_t89, &fl_t91, error));
  if (fl_t91) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_11, result, error);
  }
}

/*
 * Функция flang «Под системным временным».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_pod_sistemnym_vremennym(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t92 = fl_nothing(); /* «начинается с» */
  FL_TRY(fl_b_nachinaetsya_s(ctx, put, opis_diska_text_12, &fl_t92, error));
  bool fl_t93 = false;
  FL_TRY(fl_cond(ctx, fl_t92, &fl_t93, error));
  if (fl_t93) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t94 = fl_nothing(); /* «начинается с» */
    FL_TRY(fl_b_nachinaetsya_s(ctx, put, opis_diska_text_13, &fl_t94, error));
    *result = fl_t94;
    return FL_OK;
  }
}

/*
 * Функция flang «Примета кэша».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_kesha(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t95 = fl_nothing();
  FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_14, &fl_t95, error));
  bool fl_t96 = false;
  FL_TRY(fl_cond(ctx, fl_t95, &fl_t96, error));
  fl_value fl_t97 = fl_nothing();
  if (fl_t96) {
    fl_t97 = fl_flag(true);
  } else {
    fl_value fl_t98 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_15, &fl_t98, error));
    fl_t97 = fl_t98;
  }
  bool fl_t99 = false;
  FL_TRY(fl_cond(ctx, fl_t97, &fl_t99, error));
  fl_value fl_t100 = fl_nothing();
  if (fl_t99) {
    fl_t100 = fl_flag(true);
  } else {
    fl_value fl_t101 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_16, &fl_t101, error));
    fl_t100 = fl_t101;
  }
  bool fl_t102 = false;
  FL_TRY(fl_cond(ctx, fl_t100, &fl_t102, error));
  if (fl_t102) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_pod_sistemnym_vremennym(ctx, put, result, error);
  }
}

/*
 * Функция flang «Примета журнала».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_zhurnala(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t103 = fl_nothing();
  FL_TRY(opis_diska_imya_v_puti(ctx, put, &fl_t103, error));
  fl_value fl_t104 = fl_nothing();
  FL_TRY(opis_diska_okanchivaetsya_na(ctx, fl_t103, opis_diska_text_17, &fl_t104, error));
  bool fl_t105 = false;
  FL_TRY(fl_cond(ctx, fl_t104, &fl_t105, error));
  fl_value fl_t106 = fl_nothing();
  if (fl_t105) {
    fl_t106 = fl_flag(true);
  } else {
    fl_value fl_t107 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_18, &fl_t107, error));
    fl_t106 = fl_t107;
  }
  bool fl_t108 = false;
  FL_TRY(fl_cond(ctx, fl_t106, &fl_t108, error));
  if (fl_t108) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_19, result, error);
  }
}

/*
 * Функция flang «Примета сборки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_sborki(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t109 = fl_nothing();
  FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_20, &fl_t109, error));
  bool fl_t110 = false;
  FL_TRY(fl_cond(ctx, fl_t109, &fl_t110, error));
  fl_value fl_t111 = fl_nothing();
  if (fl_t110) {
    fl_t111 = fl_flag(true);
  } else {
    fl_value fl_t112 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_21, &fl_t112, error));
    fl_t111 = fl_t112;
  }
  bool fl_t113 = false;
  FL_TRY(fl_cond(ctx, fl_t111, &fl_t113, error));
  fl_value fl_t114 = fl_nothing();
  if (fl_t113) {
    fl_t114 = fl_flag(true);
  } else {
    fl_value fl_t115 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_22, &fl_t115, error));
    fl_t114 = fl_t115;
  }
  bool fl_t116 = false;
  FL_TRY(fl_cond(ctx, fl_t114, &fl_t116, error));
  fl_value fl_t117 = fl_nothing();
  if (fl_t116) {
    fl_t117 = fl_flag(true);
  } else {
    fl_value fl_t118 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_23, &fl_t118, error));
    fl_t117 = fl_t118;
  }
  bool fl_t119 = false;
  FL_TRY(fl_cond(ctx, fl_t117, &fl_t119, error));
  if (fl_t119) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_24, result, error);
  }
}

/*
 * Функция flang «Примета загрузки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_zagruzki(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t120 = fl_nothing();
  FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_25, &fl_t120, error));
  bool fl_t121 = false;
  FL_TRY(fl_cond(ctx, fl_t120, &fl_t121, error));
  if (fl_t121) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_26, result, error);
  }
}

/*
 * Функция flang «Есть примета».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param mesto — «место»: «Разряд»
 * @return значение
 */
fl_status opis_diska_est_primeta(fl_ctx *ctx, fl_value put, fl_value mesto, fl_value *result, fl_error *error) {
  fl_value fl_t122 = fl_nothing();
  FL_TRY(opis_diska_eto_neizvestnoe(ctx, mesto, &fl_t122, error));
  bool fl_t123 = false;
  FL_TRY(fl_cond(ctx, fl_t122, &fl_t123, error));
  fl_value fl_t124 = fl_nothing();
  if (fl_t123) {
    fl_t124 = fl_flag(false);
  } else {
    fl_t124 = fl_flag(true);
  }
  bool fl_t125 = false;
  FL_TRY(fl_cond(ctx, fl_t124, &fl_t125, error));
  fl_value fl_t126 = fl_nothing();
  if (fl_t125) {
    fl_t126 = fl_flag(true);
  } else {
    fl_value fl_t127 = fl_nothing();
    FL_TRY(opis_diska_primeta_kesha(ctx, put, &fl_t127, error));
    fl_t126 = fl_t127;
  }
  bool fl_t128 = false;
  FL_TRY(fl_cond(ctx, fl_t126, &fl_t128, error));
  fl_value fl_t129 = fl_nothing();
  if (fl_t128) {
    fl_t129 = fl_flag(true);
  } else {
    fl_value fl_t130 = fl_nothing();
    FL_TRY(opis_diska_primeta_zhurnala(ctx, put, &fl_t130, error));
    fl_t129 = fl_t130;
  }
  bool fl_t131 = false;
  FL_TRY(fl_cond(ctx, fl_t129, &fl_t131, error));
  fl_value fl_t132 = fl_nothing();
  if (fl_t131) {
    fl_t132 = fl_flag(true);
  } else {
    fl_value fl_t133 = fl_nothing();
    FL_TRY(opis_diska_primeta_sborki(ctx, put, &fl_t133, error));
    fl_t132 = fl_t133;
  }
  bool fl_t134 = false;
  FL_TRY(fl_cond(ctx, fl_t132, &fl_t134, error));
  if (fl_t134) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_primeta_zagruzki(ctx, put, result, error);
  }
}

/*
 * Функция flang «Разряд решён размером».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_reshyon_razmerom(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error) {
  if (fl_variant_is(razryad, "Крупное")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Неизвестное")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Кэш")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Сборка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Загрузка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Разряд находки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param mesto — «место»: «Разряд»
 * @return значение: «Разряд»
 */
fl_status opis_diska_razryad_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value mesto, fl_value *result, fl_error *error) {
  fl_value fl_t135 = fl_nothing();
  FL_TRY(opis_diska_eto_neizvestnoe(ctx, mesto, &fl_t135, error));
  bool fl_t136 = false;
  FL_TRY(fl_cond(ctx, fl_t135, &fl_t136, error));
  fl_value fl_t137 = fl_nothing();
  if (fl_t136) {
    fl_t137 = fl_flag(false);
  } else {
    fl_t137 = fl_flag(true);
  }
  bool fl_t138 = false;
  FL_TRY(fl_cond(ctx, fl_t137, &fl_t138, error));
  fl_value fl_t139 = fl_nothing();
  if (fl_t138) {
    fl_t139 = mesto;
  } else {
    fl_value fl_t140 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t140, error));
    fl_value fl_t141 = fl_nothing();
    FL_TRY(opis_diska_primeta_kesha(ctx, fl_t140, &fl_t141, error));
    bool fl_t142 = false;
    FL_TRY(fl_cond(ctx, fl_t141, &fl_t142, error));
    fl_value fl_t143 = fl_nothing();
    if (fl_t142) {
      fl_value fl_t144 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "Кэш", NULL, NULL, 0, &fl_t144, error));
      fl_t143 = fl_t144;
    } else {
      fl_value fl_t145 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t145, error));
      fl_value fl_t146 = fl_nothing();
      FL_TRY(opis_diska_primeta_zhurnala(ctx, fl_t145, &fl_t146, error));
      bool fl_t147 = false;
      FL_TRY(fl_cond(ctx, fl_t146, &fl_t147, error));
      fl_value fl_t148 = fl_nothing();
      if (fl_t147) {
        fl_value fl_t149 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "Журнал", NULL, NULL, 0, &fl_t149, error));
        fl_t148 = fl_t149;
      } else {
        fl_value fl_t150 = fl_nothing();
        FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t150, error));
        fl_value fl_t151 = fl_nothing();
        FL_TRY(opis_diska_primeta_sborki(ctx, fl_t150, &fl_t151, error));
        bool fl_t152 = false;
        FL_TRY(fl_cond(ctx, fl_t151, &fl_t152, error));
        fl_value fl_t153 = fl_nothing();
        if (fl_t152) {
          fl_value fl_t154 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "Сборка", NULL, NULL, 0, &fl_t154, error));
          fl_t153 = fl_t154;
        } else {
          fl_value fl_t155 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t155, error));
          fl_value fl_t156 = fl_nothing();
          FL_TRY(opis_diska_primeta_zagruzki(ctx, fl_t155, &fl_t156, error));
          bool fl_t157 = false;
          FL_TRY(fl_cond(ctx, fl_t156, &fl_t157, error));
          fl_value fl_t158 = fl_nothing();
          if (fl_t157) {
            fl_value fl_t159 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Загрузка", NULL, NULL, 0, &fl_t159, error));
            fl_t158 = fl_t159;
          } else {
            fl_value fl_t160 = fl_nothing();
            FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t160, error));
            fl_value fl_t161 = fl_nothing();
            FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t161, error));
            if (fl_t160.tag != FL_NUMBER || fl_t161.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t160, fl_t161, error));
            bool fl_t162 = false;
            FL_TRY(fl_cond(ctx, fl_flag(fl_t160.as.number >= fl_t161.as.number), &fl_t162, error));
            fl_value fl_t163 = fl_nothing();
            if (fl_t162) {
              fl_value fl_t164 = fl_nothing();
              FL_TRY(fl_variant_new(ctx, "Крупное", NULL, NULL, 0, &fl_t164, error));
              fl_t163 = fl_t164;
            } else {
              fl_value fl_t165 = fl_nothing();
              FL_TRY(fl_variant_new(ctx, "Неизвестное", NULL, NULL, 0, &fl_t165, error));
              fl_t163 = fl_t165;
            }
            fl_t158 = fl_t163;
          }
          fl_t153 = fl_t158;
        }
        fl_t148 = fl_t153;
      }
      fl_t143 = fl_t148;
    }
    fl_t139 = fl_t143;
  }
  const fl_value fl_t166 = fl_t139;
  fl_value fl_t167 = fl_nothing();
  FL_TRY(opis_diska_razryad_obosnovan(ctx, nahodka, mesto, fl_t166, &fl_t167, error));
  /* постусловие «Разряд обоснован приметой или местом» */
  bool fl_t168 = false;
  FL_TRY(fl_post(ctx, fl_t167, "Разряд обоснован приметой или местом", "Разряд находки", &fl_t168, error));
  if (!fl_t168) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Разряд обоснован приметой или местом» функции «Разряд находки»");
  }
  fl_value fl_t169 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t169, error));
  fl_value fl_t170 = fl_nothing();
  FL_TRY(opis_diska_est_primeta(ctx, fl_t169, mesto, &fl_t170, error));
  bool fl_t171 = false;
  FL_TRY(fl_cond(ctx, fl_t170, &fl_t171, error));
  fl_value fl_t172 = fl_nothing();
  if (fl_t171) {
    fl_t172 = fl_flag(false);
  } else {
    fl_t172 = fl_flag(true);
  }
  bool fl_t173 = false;
  FL_TRY(fl_cond(ctx, fl_t172, &fl_t173, error));
  fl_value fl_t174 = fl_nothing();
  if (fl_t173) {
    fl_value fl_t175 = fl_nothing();
    FL_TRY(opis_diska_razryad_reshyon_razmerom(ctx, fl_t166, &fl_t175, error));
    fl_t174 = fl_t175;
  } else {
    fl_t174 = fl_flag(true);
  }
  /* постусловие «И4: без приметы-составляющей и без места разряд решает размер» */
  bool fl_t176 = false;
  FL_TRY(fl_post(ctx, fl_t174, "И4: без приметы-составляющей и без места разряд решает размер", "Разряд находки", &fl_t176, error));
  if (!fl_t176) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И4: без приметы-составляющей и без места разряд решает размер» функции «Разряд находки»");
  }
  *result = fl_t166;
  return FL_OK;
}

/*
 * Функция flang «Крупное не мельче порога».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_krupnoe_ne_melche_poroga(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error) {
  if (fl_variant_is(razryad, "Крупное")) {
    fl_value fl_t177 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t177, error));
    fl_value fl_t178 = fl_nothing();
    FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t178, error));
    if (fl_t177.tag != FL_NUMBER || fl_t178.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t177, fl_t178, error));
    *result = fl_flag(fl_t177.as.number >= fl_t178.as.number);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Кэш")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Сборка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Загрузка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Неизвестное")) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Каталог».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @return значение
 */
fl_status opis_diska_katalog(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error) {
  fl_value fl_t179 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "вид", &fl_t179, error));
  if (fl_variant_is(fl_t179, "Каталог")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t179, "Файл")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(fl_t179, "Ссылка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t179, error);
  }
}

/*
 * Функция flang «Ссылка».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @return значение
 */
fl_status opis_diska_ssylka(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error) {
  fl_value fl_t180 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "вид", &fl_t180, error));
  if (fl_variant_is(fl_t180, "Ссылка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t180, "Файл")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(fl_t180, "Каталог")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t180, error);
  }
}

/*
 * Функция flang «Приговор мусора».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param porog — «порог»: число
 * @return значение: «Приговор»
 */
fl_status opis_diska_prigovor_musora(fl_ctx *ctx, fl_value nahodka, fl_value porog, fl_value *result, fl_error *error) {
  fl_value fl_t181 = fl_nothing();
  FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t181, error));
  bool fl_t182 = false;
  FL_TRY(fl_cond(ctx, fl_t181, &fl_t182, error));
  fl_value fl_t183 = fl_nothing();
  if (fl_t182) {
    fl_value fl_t184 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t184, error));
    fl_t183 = fl_t184;
  } else {
    fl_value fl_t185 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t185, error));
    if (fl_t185.tag != FL_NUMBER || porog.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t185, porog, error));
    bool fl_t186 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_t185.as.number >= porog.as.number), &fl_t186, error));
    fl_value fl_t187 = fl_nothing();
    if (fl_t186) {
      fl_value fl_t188 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "МожноУбрать", NULL, NULL, 0, &fl_t188, error));
      fl_t187 = fl_t188;
    } else {
      fl_value fl_t189 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t189, error));
      fl_t187 = fl_t189;
    }
    fl_t183 = fl_t187;
  }
  const fl_value fl_t190 = fl_t183;
  fl_value fl_t191 = fl_nothing();
  FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t191, error));
  bool fl_t192 = false;
  FL_TRY(fl_cond(ctx, fl_t191, &fl_t192, error));
  fl_value fl_t193 = fl_nothing();
  if (fl_t192) {
    fl_value fl_t194 = fl_nothing();
    FL_TRY(opis_diska_eto_mozhnoubrat(ctx, fl_t190, &fl_t194, error));
    bool fl_t195 = false;
    FL_TRY(fl_cond(ctx, fl_t194, &fl_t195, error));
    fl_value fl_t196 = fl_nothing();
    if (fl_t195) {
      fl_t196 = fl_flag(false);
    } else {
      fl_t196 = fl_flag(true);
    }
    fl_t193 = fl_t196;
  } else {
    fl_t193 = fl_flag(true);
  }
  /* постусловие «Каталог не убирается» */
  bool fl_t197 = false;
  FL_TRY(fl_post(ctx, fl_t193, "Каталог не убирается", "Приговор мусора", &fl_t197, error));
  if (!fl_t197) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Каталог не убирается» функции «Приговор мусора»");
  }
  *result = fl_t190;
  return FL_OK;
}

/*
 * Функция flang «Приговор находки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @return значение: «Приговор»
 */
fl_status opis_diska_prigovor_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error) {
  fl_value fl_t198 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "доступен", &fl_t198, error));
  bool fl_t199 = false;
  FL_TRY(fl_cond(ctx, fl_t198, &fl_t199, error));
  fl_value fl_t200 = fl_nothing();
  if (fl_t199) {
    fl_t200 = fl_flag(false);
  } else {
    fl_t200 = fl_flag(true);
  }
  bool fl_t201 = false;
  FL_TRY(fl_cond(ctx, fl_t200, &fl_t201, error));
  fl_value fl_t202 = fl_nothing();
  if (fl_t201) {
    fl_value fl_t203 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t203, error));
    fl_t202 = fl_t203;
  } else {
    fl_value fl_t204 = fl_nothing();
    FL_TRY(opis_diska_ssylka(ctx, nahodka, &fl_t204, error));
    bool fl_t205 = false;
    FL_TRY(fl_cond(ctx, fl_t204, &fl_t205, error));
    fl_value fl_t206 = fl_nothing();
    if (fl_t205) {
      fl_value fl_t207 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t207, error));
      fl_t206 = fl_t207;
    } else {
      fl_value fl_t208 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t208, error));
      fl_value fl_t209 = fl_nothing();
      FL_TRY(opis_diska_adresuetsya_soderzhimym(ctx, fl_t208, &fl_t209, error));
      bool fl_t210 = false;
      FL_TRY(fl_cond(ctx, fl_t209, &fl_t210, error));
      fl_value fl_t211 = fl_nothing();
      if (fl_t210) {
        fl_value fl_t212 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t212, error));
        fl_t211 = fl_t212;
      } else {
        fl_value fl_t213 = fl_nothing();
        if (fl_variant_is(razryad, "Кэш")) {
          fl_value fl_t214 = fl_nothing();
          FL_TRY(opis_diska_porog_kesha(ctx, &fl_t214, error));
          fl_value fl_t215 = fl_nothing();
          FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t214, &fl_t215, error));
          fl_t213 = fl_t215;
        } else if (fl_variant_is(razryad, "Сборка")) {
          fl_value fl_t216 = fl_nothing();
          FL_TRY(opis_diska_porog_kesha(ctx, &fl_t216, error));
          fl_value fl_t217 = fl_nothing();
          FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t216, &fl_t217, error));
          fl_t213 = fl_t217;
        } else if (fl_variant_is(razryad, "Журнал")) {
          fl_value fl_t218 = fl_nothing();
          FL_TRY(opis_diska_porog_zhurnala(ctx, &fl_t218, error));
          fl_value fl_t219 = fl_nothing();
          FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t218, &fl_t219, error));
          fl_t213 = fl_t219;
        } else if (fl_variant_is(razryad, "Загрузка")) {
          fl_value fl_t220 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t220, error));
          fl_value fl_t221 = fl_nothing();
          FL_TRY(opis_diska_porog_zagruzki(ctx, &fl_t221, error));
          if (fl_t220.tag != FL_NUMBER || fl_t221.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t220, fl_t221, error));
          bool fl_t222 = false;
          FL_TRY(fl_cond(ctx, fl_flag(fl_t220.as.number >= fl_t221.as.number), &fl_t222, error));
          fl_value fl_t223 = fl_nothing();
          if (fl_t222) {
            fl_value fl_t224 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t224, error));
            fl_t223 = fl_t224;
          } else {
            fl_value fl_t225 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t225, error));
            fl_t223 = fl_t225;
          }
          fl_t213 = fl_t223;
        } else if (fl_variant_is(razryad, "Крупное")) {
          fl_value fl_t226 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t226, error));
          fl_t213 = fl_t226;
        } else if (fl_variant_is(razryad, "Неизвестное")) {
          fl_value fl_t227 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t227, error));
          fl_t213 = fl_t227;
        } else {
          return fl_match_fail(ctx, razryad, error);
        }
        fl_t211 = fl_t213;
      }
      fl_t206 = fl_t211;
    }
    fl_t202 = fl_t206;
  }
  const fl_value fl_t228 = fl_t202;
  fl_value fl_t229 = fl_nothing();
  FL_TRY(opis_diska_prigovor_obosnovan(ctx, nahodka, razryad, fl_t228, &fl_t229, error));
  /* постусловие «Приговор обоснован» */
  bool fl_t230 = false;
  FL_TRY(fl_post(ctx, fl_t229, "Приговор обоснован", "Приговор находки", &fl_t230, error));
  if (!fl_t230) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Приговор обоснован» функции «Приговор находки»");
  }
  fl_value fl_t231 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t231, error));
  fl_value fl_t232 = fl_nothing();
  FL_TRY(opis_diska_adresuetsya_soderzhimym(ctx, fl_t231, &fl_t232, error));
  bool fl_t233 = false;
  FL_TRY(fl_cond(ctx, fl_t232, &fl_t233, error));
  fl_value fl_t234 = fl_nothing();
  if (fl_t233) {
    fl_value fl_t235 = fl_nothing();
    FL_TRY(opis_diska_eto_netrogat(ctx, fl_t228, &fl_t235, error));
    fl_t234 = fl_t235;
  } else {
    fl_t234 = fl_flag(true);
  }
  /* постусловие «И3: адресуемое содержимым не убирается» */
  bool fl_t236 = false;
  FL_TRY(fl_post(ctx, fl_t234, "И3: адресуемое содержимым не убирается", "Приговор находки", &fl_t236, error));
  if (!fl_t236) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И3: адресуемое содержимым не убирается» функции «Приговор находки»");
  }
  *result = fl_t228;
  return FL_OK;
}

/*
 * Функция flang «Вес находки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param prigovor — «приговор»: «Приговор»
 * @return значение: число
 */
fl_status opis_diska_ves_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value prigovor, fl_value *result, fl_error *error) {
  fl_value fl_t237 = fl_nothing();
  if (fl_variant_is(prigovor, "НеТрогать")) {
    fl_t237 = fl_number(0.0);
  } else if (fl_variant_is(prigovor, "МожноУбрать")) {
    fl_value fl_t238 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t238, error));
    fl_t237 = fl_t238;
  } else if (fl_variant_is(prigovor, "Спросить")) {
    fl_value fl_t239 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t239, error));
    fl_t237 = fl_t239;
  } else {
    return fl_match_fail(ctx, prigovor, error);
  }
  const fl_value fl_t240 = fl_t237;
  fl_value fl_t241 = fl_nothing();
  FL_TRY(opis_diska_ves_obosnovan(ctx, nahodka, prigovor, fl_t240, &fl_t241, error));
  /* постусловие «Вес обоснован» */
  bool fl_t242 = false;
  FL_TRY(fl_post(ctx, fl_t241, "Вес обоснован", "Вес находки", &fl_t242, error));
  if (!fl_t242) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Вес обоснован» функции «Вес находки»");
  }
  *result = fl_t240;
  return FL_OK;
}

/*
 * Функция flang «Вес в границах».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param ves — «вес»: число
 * @return значение
 */
fl_status opis_diska_ves_v_granicah(fl_ctx *ctx, fl_value nahodka, fl_value ves, fl_value *result, fl_error *error) {
  if (ves.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, ves, fl_number(0.0), error));
  bool fl_t243 = false;
  FL_TRY(fl_cond(ctx, fl_flag(ves.as.number >= 0.0), &fl_t243, error));
  if (fl_t243) {
    fl_value fl_t244 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t244, error));
    if (ves.tag != FL_NUMBER || fl_t244.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, ves, fl_t244, error));
    *result = fl_flag(ves.as.number <= fl_t244.as.number);
    return FL_OK;
  } else {
    *result = fl_flag(false);
    return FL_OK;
  }
}

/*
 * Функция flang «Решить находку».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Решение»
 */
fl_status opis_diska_reshit_nahodku(fl_ctx *ctx, fl_value nahodka, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t245 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t245, error));
  fl_value fl_t246 = fl_nothing();
  FL_TRY(opis_diska_razryad_po_spravochniku(ctx, fl_t245, spravochnik, &fl_t246, error));
  const fl_value mesto = fl_t246; /* пусть «место» */
  fl_value fl_t247 = fl_nothing();
  FL_TRY(opis_diska_razryad_nahodki(ctx, nahodka, mesto, &fl_t247, error));
  const fl_value razryad = fl_t247; /* пусть «разряд» */
  fl_value fl_t248 = fl_nothing();
  FL_TRY(opis_diska_prigovor_nahodki(ctx, nahodka, razryad, &fl_t248, error));
  const fl_value prigovor = fl_t248; /* пусть «приговор» */
  fl_value fl_t249 = fl_nothing();
  FL_TRY(opis_diska_ves_nahodki(ctx, nahodka, prigovor, &fl_t249, error));
  fl_value fl_t251[3];
  fl_t251[0] = razryad; /* «разряд» */
  fl_t251[1] = prigovor; /* «приговор» */
  fl_t251[2] = fl_t249; /* «вес» */
  fl_value fl_t250 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_3, fl_t251, 3, &fl_t250, error));
  const fl_value fl_t252 = fl_t250;
  fl_value fl_t253 = fl_nothing();
  FL_TRY(opis_diska_i1_derzhitsya(ctx, fl_t252, &fl_t253, error));
  /* постусловие «И1: убрать можно только мусор» */
  bool fl_t254 = false;
  FL_TRY(fl_post(ctx, fl_t253, "И1: убрать можно только мусор", "Решить находку", &fl_t254, error));
  if (!fl_t254) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И1: убрать можно только мусор» функции «Решить находку»");
  }
  *result = fl_t252;
  return FL_OK;
}

/*
 * Функция flang «Решить всё».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: список: «Решение»
 */
fl_status opis_diska_reshit_vsyo(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t255 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "отобразить", &fl_t255, error));
  fl_value *fl_t256 = NULL;
  size_t fl_t257 = 0;
  FL_TRY(fl_list_alloc(ctx, fl_t255.as.list.count, &fl_t256, error));
  for (size_t fl_t258 = 0; fl_t258 < fl_t255.as.list.count; fl_t258 += 1) {
    const fl_value nahodka = fl_t255.as.list.items[fl_t258]; /* «находка» */
    fl_value fl_t259 = fl_nothing();
    FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, spravochnik, &fl_t259, error));
    fl_t256[fl_t257] = fl_t259;
    fl_t257 += 1;
  }
  const fl_value fl_t260 = fl_list(fl_t256, fl_t257);
  fl_value fl_t261 = fl_nothing();
  FL_TRY(opis_diska_i1_derzhitsya_vsyudu(ctx, fl_t260, &fl_t261, error));
  /* постусловие «И1 всюду: убрать можно только мусор» */
  bool fl_t262 = false;
  FL_TRY(fl_post(ctx, fl_t261, "И1 всюду: убрать можно только мусор", "Решить всё", &fl_t262, error));
  if (!fl_t262) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И1 всюду: убрать можно только мусор» функции «Решить всё»");
  }
  *result = fl_t260;
  return FL_OK;
}

/*
 * Функция flang «И1 держится».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param reshenie — «решение»: «Решение»
 * @return значение
 */
fl_status opis_diska_i1_derzhitsya(fl_ctx *ctx, fl_value reshenie, fl_value *result, fl_error *error) {
  fl_value fl_t263 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "приговор", &fl_t263, error));
  if (fl_variant_is(fl_t263, "МожноУбрать")) {
    fl_value fl_t264 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "разряд", &fl_t264, error));
    if (fl_variant_is(fl_t264, "Кэш")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t264, "Журнал")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t264, "Сборка")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t264, "Загрузка")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(fl_t264, "Крупное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(fl_t264, "Неизвестное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else {
      return fl_match_fail(ctx, fl_t264, error);
    }
  } else if (fl_variant_is(fl_t263, "Спросить")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t263, "НеТрогать")) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t263, error);
  }
}

/*
 * Функция flang «И1 держится всюду».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param resheniya — «решения»: список: «Решение»
 * @return значение
 */
fl_status opis_diska_i1_derzhitsya_vsyudu(fl_ctx *ctx, fl_value resheniya, fl_value *result, fl_error *error) {
  fl_value fl_t265 = fl_nothing();
  FL_TRY(fl_require_list(ctx, resheniya, "свёртка", &fl_t265, error));
  fl_value akk = fl_flag(true); /* «акк» */
  const fl_mark fl_t267 = fl_region_open(ctx);
  for (size_t fl_t266 = 0; fl_t266 < fl_t265.as.list.count; fl_t266 += 1) {
    const fl_value reshenie = fl_t265.as.list.items[fl_t266]; /* «решение» */
    bool fl_t268 = false;
    FL_TRY(fl_cond(ctx, akk, &fl_t268, error));
    fl_value fl_t269 = fl_nothing();
    if (fl_t268) {
      fl_value fl_t270 = fl_nothing();
      FL_TRY(opis_diska_i1_derzhitsya(ctx, reshenie, &fl_t270, error));
      fl_t269 = fl_t270;
    } else {
      fl_t269 = fl_flag(false);
    }
    akk = fl_t269;
    FL_TRY(fl_region_recycle(ctx, fl_t267, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t267, FL_OK, &akk, error));
  *result = akk;
  return FL_OK;
}

/*
 * Функция flang «Пустой свод».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: «Свод»
 */
fl_status opis_diska_pustoy_svod(fl_ctx *ctx, fl_value *result, fl_error *error) {
  fl_value fl_t272[7];
  fl_t272[0] = fl_number(0.0); /* «кэш» */
  fl_t272[1] = fl_number(0.0); /* «журнал» */
  fl_t272[2] = fl_number(0.0); /* «сборка» */
  fl_t272[3] = fl_number(0.0); /* «загрузка» */
  fl_t272[4] = fl_number(0.0); /* «крупное» */
  fl_t272[5] = fl_number(0.0); /* «неизвестное» */
  fl_t272[6] = fl_number(0.0); /* «освободить» */
  fl_value fl_t271 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t272, 7, &fl_t271, error));
  *result = fl_t271;
  return FL_OK;
}

/*
 * Функция flang «Прибавить решение».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param svod — «свод»: «Свод»
 * @param reshenie — «решение»: «Решение»
 * @return значение: «Свод»
 */
fl_status opis_diska_pribavit_reshenie(fl_ctx *ctx, fl_value svod, fl_value reshenie, fl_value *result, fl_error *error) {
  fl_value fl_t273 = fl_nothing();
  fl_value fl_t274 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "приговор", &fl_t274, error));
  if (fl_variant_is(fl_t274, "МожноУбрать")) {
    fl_value fl_t275 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t275, error));
    fl_t273 = fl_t275;
  } else if (fl_variant_is(fl_t274, "Спросить")) {
    fl_t273 = fl_number(0.0);
  } else if (fl_variant_is(fl_t274, "НеТрогать")) {
    fl_t273 = fl_number(0.0);
  } else {
    return fl_match_fail(ctx, fl_t274, error);
  }
  const fl_value ubrat = fl_t273; /* пусть «убрать» */
  fl_value fl_t276 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "разряд", &fl_t276, error));
  if (fl_variant_is(fl_t276, "Кэш")) {
    fl_value fl_t277 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t277, error));
    fl_value fl_t278 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t278, error));
    if (fl_t277.tag != FL_NUMBER || fl_t278.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t277, fl_t278, error));
    fl_value fl_t279 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t279, error));
    fl_value fl_t280 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t280, error));
    fl_value fl_t281 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t281, error));
    fl_value fl_t282 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t282, error));
    fl_value fl_t283 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t283, error));
    fl_value fl_t284 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t284, error));
    if (fl_t284.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t284, ubrat, error));
    fl_value fl_t286[7];
    fl_t286[0] = fl_number(fl_t277.as.number + fl_t278.as.number); /* «кэш» */
    fl_t286[1] = fl_t279; /* «журнал» */
    fl_t286[2] = fl_t280; /* «сборка» */
    fl_t286[3] = fl_t281; /* «загрузка» */
    fl_t286[4] = fl_t282; /* «крупное» */
    fl_t286[5] = fl_t283; /* «неизвестное» */
    fl_t286[6] = fl_number(fl_t284.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t285 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t286, 7, &fl_t285, error));
    *result = fl_t285;
    return FL_OK;
  } else if (fl_variant_is(fl_t276, "Журнал")) {
    fl_value fl_t287 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t287, error));
    fl_value fl_t288 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t288, error));
    fl_value fl_t289 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t289, error));
    if (fl_t288.tag != FL_NUMBER || fl_t289.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t288, fl_t289, error));
    fl_value fl_t290 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t290, error));
    fl_value fl_t291 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t291, error));
    fl_value fl_t292 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t292, error));
    fl_value fl_t293 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t293, error));
    fl_value fl_t294 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t294, error));
    if (fl_t294.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t294, ubrat, error));
    fl_value fl_t296[7];
    fl_t296[0] = fl_t287; /* «кэш» */
    fl_t296[1] = fl_number(fl_t288.as.number + fl_t289.as.number); /* «журнал» */
    fl_t296[2] = fl_t290; /* «сборка» */
    fl_t296[3] = fl_t291; /* «загрузка» */
    fl_t296[4] = fl_t292; /* «крупное» */
    fl_t296[5] = fl_t293; /* «неизвестное» */
    fl_t296[6] = fl_number(fl_t294.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t295 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t296, 7, &fl_t295, error));
    *result = fl_t295;
    return FL_OK;
  } else if (fl_variant_is(fl_t276, "Сборка")) {
    fl_value fl_t297 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t297, error));
    fl_value fl_t298 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t298, error));
    fl_value fl_t299 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t299, error));
    fl_value fl_t300 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t300, error));
    if (fl_t299.tag != FL_NUMBER || fl_t300.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t299, fl_t300, error));
    fl_value fl_t301 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t301, error));
    fl_value fl_t302 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t302, error));
    fl_value fl_t303 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t303, error));
    fl_value fl_t304 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t304, error));
    if (fl_t304.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t304, ubrat, error));
    fl_value fl_t306[7];
    fl_t306[0] = fl_t297; /* «кэш» */
    fl_t306[1] = fl_t298; /* «журнал» */
    fl_t306[2] = fl_number(fl_t299.as.number + fl_t300.as.number); /* «сборка» */
    fl_t306[3] = fl_t301; /* «загрузка» */
    fl_t306[4] = fl_t302; /* «крупное» */
    fl_t306[5] = fl_t303; /* «неизвестное» */
    fl_t306[6] = fl_number(fl_t304.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t305 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t306, 7, &fl_t305, error));
    *result = fl_t305;
    return FL_OK;
  } else if (fl_variant_is(fl_t276, "Загрузка")) {
    fl_value fl_t307 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t307, error));
    fl_value fl_t308 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t308, error));
    fl_value fl_t309 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t309, error));
    fl_value fl_t310 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t310, error));
    fl_value fl_t311 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t311, error));
    if (fl_t310.tag != FL_NUMBER || fl_t311.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t310, fl_t311, error));
    fl_value fl_t312 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t312, error));
    fl_value fl_t313 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t313, error));
    fl_value fl_t314 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t314, error));
    if (fl_t314.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t314, ubrat, error));
    fl_value fl_t316[7];
    fl_t316[0] = fl_t307; /* «кэш» */
    fl_t316[1] = fl_t308; /* «журнал» */
    fl_t316[2] = fl_t309; /* «сборка» */
    fl_t316[3] = fl_number(fl_t310.as.number + fl_t311.as.number); /* «загрузка» */
    fl_t316[4] = fl_t312; /* «крупное» */
    fl_t316[5] = fl_t313; /* «неизвестное» */
    fl_t316[6] = fl_number(fl_t314.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t315 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t316, 7, &fl_t315, error));
    *result = fl_t315;
    return FL_OK;
  } else if (fl_variant_is(fl_t276, "Крупное")) {
    fl_value fl_t317 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t317, error));
    fl_value fl_t318 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t318, error));
    fl_value fl_t319 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t319, error));
    fl_value fl_t320 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t320, error));
    fl_value fl_t321 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t321, error));
    fl_value fl_t322 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t322, error));
    if (fl_t321.tag != FL_NUMBER || fl_t322.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t321, fl_t322, error));
    fl_value fl_t323 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t323, error));
    fl_value fl_t324 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t324, error));
    if (fl_t324.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t324, ubrat, error));
    fl_value fl_t326[7];
    fl_t326[0] = fl_t317; /* «кэш» */
    fl_t326[1] = fl_t318; /* «журнал» */
    fl_t326[2] = fl_t319; /* «сборка» */
    fl_t326[3] = fl_t320; /* «загрузка» */
    fl_t326[4] = fl_number(fl_t321.as.number + fl_t322.as.number); /* «крупное» */
    fl_t326[5] = fl_t323; /* «неизвестное» */
    fl_t326[6] = fl_number(fl_t324.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t325 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t326, 7, &fl_t325, error));
    *result = fl_t325;
    return FL_OK;
  } else if (fl_variant_is(fl_t276, "Неизвестное")) {
    fl_value fl_t327 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t327, error));
    fl_value fl_t328 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t328, error));
    fl_value fl_t329 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t329, error));
    fl_value fl_t330 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t330, error));
    fl_value fl_t331 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t331, error));
    fl_value fl_t332 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t332, error));
    fl_value fl_t333 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t333, error));
    if (fl_t332.tag != FL_NUMBER || fl_t333.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t332, fl_t333, error));
    fl_value fl_t334 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t334, error));
    if (fl_t334.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t334, ubrat, error));
    fl_value fl_t336[7];
    fl_t336[0] = fl_t327; /* «кэш» */
    fl_t336[1] = fl_t328; /* «журнал» */
    fl_t336[2] = fl_t329; /* «сборка» */
    fl_t336[3] = fl_t330; /* «загрузка» */
    fl_t336[4] = fl_t331; /* «крупное» */
    fl_t336[5] = fl_number(fl_t332.as.number + fl_t333.as.number); /* «неизвестное» */
    fl_t336[6] = fl_number(fl_t334.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t335 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_5, fl_t336, 7, &fl_t335, error));
    *result = fl_t335;
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t276, error);
  }
}

/*
 * Функция flang «Свести».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Свод»
 */
fl_status opis_diska_svesti(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t337 = fl_nothing();
  FL_TRY(opis_diska_reshit_vsyo(ctx, zapisi, spravochnik, &fl_t337, error));
  fl_value fl_t338 = fl_nothing();
  FL_TRY(fl_require_list(ctx, fl_t337, "свёртка", &fl_t338, error));
  fl_value fl_t339 = fl_nothing();
  FL_TRY(opis_diska_pustoy_svod(ctx, &fl_t339, error));
  fl_value svod = fl_t339; /* «свод» */
  const fl_mark fl_t341 = fl_region_open(ctx);
  for (size_t fl_t340 = 0; fl_t340 < fl_t338.as.list.count; fl_t340 += 1) {
    const fl_value reshenie = fl_t338.as.list.items[fl_t340]; /* «решение» */
    fl_value fl_t342 = fl_nothing();
    FL_TRY(opis_diska_pribavit_reshenie(ctx, svod, reshenie, &fl_t342, error));
    svod = fl_t342;
    FL_TRY(fl_region_recycle(ctx, fl_t341, &svod, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t341, FL_OK, &svod, error));
  const fl_value fl_t343 = svod;
  fl_value fl_t344 = fl_nothing();
  FL_TRY(opis_diska_i2_derzhitsya(ctx, zapisi, spravochnik, fl_t343, &fl_t344, error));
  /* постусловие «И2: освобождаемое не больше убираемого» */
  bool fl_t345 = false;
  FL_TRY(fl_post(ctx, fl_t344, "И2: освобождаемое не больше убираемого", "Свести", &fl_t345, error));
  if (!fl_t345) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И2: освобождаемое не больше убираемого» функции «Свести»");
  }
  *result = fl_t343;
  return FL_OK;
}

/*
 * Функция flang «Сумма размеров убираемых».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: число
 */
fl_status opis_diska_summa_razmerov_ubiraemyh(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t346 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "свёртка", &fl_t346, error));
  fl_value akk = fl_number(0.0); /* «акк» */
  const fl_mark fl_t348 = fl_region_open(ctx);
  for (size_t fl_t347 = 0; fl_t347 < fl_t346.as.list.count; fl_t347 += 1) {
    const fl_value nahodka = fl_t346.as.list.items[fl_t347]; /* «находка» */
    fl_value fl_t349 = fl_nothing();
    fl_value fl_t350 = fl_nothing();
    FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, spravochnik, &fl_t350, error));
    fl_value fl_t351 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t350, "приговор", &fl_t351, error));
    if (fl_variant_is(fl_t351, "МожноУбрать")) {
      fl_value fl_t352 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t352, error));
      if (akk.tag != FL_NUMBER || fl_t352.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", akk, fl_t352, error));
      fl_t349 = fl_number(akk.as.number + fl_t352.as.number);
    } else if (fl_variant_is(fl_t351, "Спросить")) {
      fl_t349 = akk;
    } else if (fl_variant_is(fl_t351, "НеТрогать")) {
      fl_t349 = akk;
    } else {
      return fl_match_fail(ctx, fl_t351, error);
    }
    akk = fl_t349;
    FL_TRY(fl_region_recycle(ctx, fl_t348, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t348, FL_OK, &akk, error));
  *result = akk;
  return FL_OK;
}

/*
 * Функция flang «И2 держится».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @param svod — «свод»: «Свод»
 * @return значение
 */
fl_status opis_diska_i2_derzhitsya(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value svod, fl_value *result, fl_error *error) {
  fl_value fl_t353 = fl_nothing();
  FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t353, error));
  fl_value fl_t354 = fl_nothing();
  FL_TRY(opis_diska_summa_razmerov_ubiraemyh(ctx, zapisi, spravochnik, &fl_t354, error));
  if (fl_t353.tag != FL_NUMBER || fl_t354.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t353, fl_t354, error));
  *result = fl_flag(fl_t353.as.number <= fl_t354.as.number);
  return FL_OK;
}

/*
 * Функция flang «Строку отчёта».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Строка отчёта»
 */
fl_status opis_diska_stroku_otchyota(fl_ctx *ctx, fl_value nahodka, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t355 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t355, error));
  fl_value fl_t356 = fl_nothing();
  FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, spravochnik, &fl_t356, error));
  fl_value fl_t358[2];
  fl_t358[0] = fl_t355; /* «путь» */
  fl_t358[1] = fl_t356; /* «решение» */
  fl_value fl_t357 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t358, 2, &fl_t357, error));
  *result = fl_t357;
  return FL_OK;
}

/* Тело «Вставить по весу»; глубину считает обёртка ниже. */
static fl_status opis_diska_vstavit_po_vesu_body(fl_ctx *ctx, fl_value stroka, fl_value stroki, fl_value *result, fl_error *error) {
  if (fl_chain_empty(stroki)) {
    fl_value *fl_t359 = NULL;
    FL_TRY(fl_list_alloc(ctx, 1, &fl_t359, error));
    fl_t359[0] = stroka;
    *result = fl_list(fl_t359, 1);
    return FL_OK;
  } else if (fl_chain_cons(stroki)) {
    const fl_value golova = fl_chain_head(stroki); /* голова «голова» */
    const fl_value hvost = fl_chain_tail(stroki); /* хвост «хвост» */
    fl_value fl_t360 = fl_nothing();
    FL_TRY(fl_field_get(ctx, stroka, "решение", &fl_t360, error));
    fl_value fl_t361 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t360, "вес", &fl_t361, error));
    fl_value fl_t362 = fl_nothing();
    FL_TRY(fl_field_get(ctx, golova, "решение", &fl_t362, error));
    fl_value fl_t363 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t362, "вес", &fl_t363, error));
    if (fl_t361.tag != FL_NUMBER || fl_t363.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t361, fl_t363, error));
    bool fl_t364 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_t361.as.number >= fl_t363.as.number), &fl_t364, error));
    if (fl_t364) {
      return opis_diska_pripisat_stroku_otchyota(ctx, stroka, stroki, result, error);
    } else {
      fl_value fl_t365 = fl_nothing();
      FL_TRY(opis_diska_vstavit_po_vesu(ctx, stroka, hvost, &fl_t365, error));
      return opis_diska_pripisat_stroku_otchyota(ctx, golova, fl_t365, result, error);
    }
  } else {
    return fl_match_fail(ctx, stroki, error);
  }
}

/*
 * Функция flang «Вставить по весу».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 *
 * Рекурсивная: считает глубину, на превышении — FLANG_RECURSION_LIMIT.
 * @param stroka — «строка»: «Строка отчёта»
 * @param stroki — «строки»: список: «Строка отчёта»
 * @return значение: список: «Строка отчёта»
 */
fl_status opis_diska_vstavit_po_vesu(fl_ctx *ctx, fl_value stroka, fl_value stroki, fl_value *result, fl_error *error) {
  FL_TRY(fl_enter(ctx, "Вставить по весу", error));
  {
    const fl_mark region = fl_region_open(ctx);
    const fl_status status = opis_diska_vstavit_po_vesu_body(ctx, stroka, stroki, result, error);
    fl_leave(ctx);
    return fl_region_close(ctx, region, status, result, error);
  }
}

/*
 * Функция flang «Приписать строку отчёта».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param pervaya — «первая»: «Строка отчёта»
 * @param stroki — «строки»: список: «Строка отчёта»
 * @return значение: список: «Строка отчёта»
 */
fl_status opis_diska_pripisat_stroku_otchyota(fl_ctx *ctx, fl_value pervaya, fl_value stroki, fl_value *result, fl_error *error) {
  fl_value fl_t366 = fl_nothing();
  FL_TRY(fl_require_list(ctx, stroki, "свёртка", &fl_t366, error));
  fl_value *fl_t367 = NULL;
  FL_TRY(fl_list_alloc(ctx, 1, &fl_t367, error));
  fl_t367[0] = pervaya;
  fl_value akk = fl_list(fl_t367, 1); /* «акк» */
  const fl_mark fl_t369 = fl_region_open(ctx);
  for (size_t fl_t368 = 0; fl_t368 < fl_t366.as.list.count; fl_t368 += 1) {
    const fl_value el = fl_t366.as.list.items[fl_t368]; /* «эл» */
    fl_value fl_t370 = fl_nothing(); /* «добавить» */
    FL_TRY(fl_b_dobavit(ctx, el, akk, &fl_t370, error));
    akk = fl_t370;
    FL_TRY(fl_region_recycle(ctx, fl_t369, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t369, FL_OK, &akk, error));
  *result = akk;
  return FL_OK;
}

/*
 * Функция flang «Отчёт».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: список: «Строка отчёта»
 */
fl_status opis_diska_otchyot(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error) {
  fl_value fl_t371 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "свёртка", &fl_t371, error));
  fl_value akk = fl_list(NULL, 0); /* «акк» */
  const fl_mark fl_t373 = fl_region_open(ctx);
  for (size_t fl_t372 = 0; fl_t372 < fl_t371.as.list.count; fl_t372 += 1) {
    const fl_value nahodka = fl_t371.as.list.items[fl_t372]; /* «находка» */
    fl_value fl_t374 = fl_nothing();
    FL_TRY(opis_diska_stroku_otchyota(ctx, nahodka, spravochnik, &fl_t374, error));
    fl_value fl_t375 = fl_nothing();
    FL_TRY(opis_diska_vstavit_po_vesu(ctx, fl_t374, akk, &fl_t375, error));
    akk = fl_t375;
    FL_TRY(fl_region_recycle(ctx, fl_t373, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t373, FL_OK, &akk, error));
  const fl_value fl_t376 = akk;
  fl_value fl_t377 = fl_nothing();
  FL_TRY(opis_diska_otchyot_toy_zhe_dliny(ctx, zapisi, fl_t376, &fl_t377, error));
  /* постусловие «Отчёт той же длины» */
  bool fl_t378 = false;
  FL_TRY(fl_post(ctx, fl_t377, "Отчёт той же длины", "Отчёт", &fl_t378, error));
  if (!fl_t378) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Отчёт той же длины» функции «Отчёт»");
  }
  *result = fl_t376;
  return FL_OK;
}

/*
 * Функция flang «Отчёт той же длины».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param stroki — «строки»: список: «Строка отчёта»
 * @return значение
 */
fl_status opis_diska_otchyot_toy_zhe_dliny(fl_ctx *ctx, fl_value zapisi, fl_value stroki, fl_value *result, fl_error *error) {
  fl_value fl_t379 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, stroki, &fl_t379, error));
  fl_value fl_t380 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, zapisi, &fl_t380, error));
  *result = fl_flag(fl_equal(fl_t379, fl_t380));
  return FL_OK;
}

/*
 * Функция flang «Это МожноУбрать».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_eto_mozhnoubrat(fl_ctx *ctx, fl_value prigovor, fl_value *result, fl_error *error) {
  if (fl_variant_is(prigovor, "МожноУбрать")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(prigovor, "Спросить")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(prigovor, "НеТрогать")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, prigovor, error);
  }
}

/*
 * Функция flang «Это НеТрогать».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_eto_netrogat(fl_ctx *ctx, fl_value prigovor, fl_value *result, fl_error *error) {
  if (fl_variant_is(prigovor, "НеТрогать")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(prigovor, "Спросить")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(prigovor, "МожноУбрать")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, prigovor, error);
  }
}

/*
 * Функция flang «И1 на паре».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_i1_na_pare(fl_ctx *ctx, fl_value razryad, fl_value prigovor, fl_value *result, fl_error *error) {
  fl_value fl_t381 = fl_nothing();
  FL_TRY(opis_diska_eto_mozhnoubrat(ctx, prigovor, &fl_t381, error));
  bool fl_t382 = false;
  FL_TRY(fl_cond(ctx, fl_t381, &fl_t382, error));
  fl_value fl_t383 = fl_nothing();
  if (fl_t382) {
    fl_t383 = fl_flag(false);
  } else {
    fl_t383 = fl_flag(true);
  }
  bool fl_t384 = false;
  FL_TRY(fl_cond(ctx, fl_t383, &fl_t384, error));
  if (fl_t384) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    if (fl_variant_is(razryad, "Кэш")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(razryad, "Журнал")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(razryad, "Сборка")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(razryad, "Загрузка")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(razryad, "Крупное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(razryad, "Неизвестное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else {
      return fl_match_fail(ctx, razryad, error);
    }
  }
}

/*
 * Функция flang «Порог разряда».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение: число
 */
fl_status opis_diska_porog_razryada(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error) {
  if (fl_variant_is(razryad, "Кэш")) {
    return opis_diska_porog_kesha(ctx, result, error);
  } else if (fl_variant_is(razryad, "Сборка")) {
    return opis_diska_porog_kesha(ctx, result, error);
  } else if (fl_variant_is(razryad, "Журнал")) {
    return opis_diska_porog_zhurnala(ctx, result, error);
  } else if (fl_variant_is(razryad, "Загрузка")) {
    return opis_diska_porog_zagruzki(ctx, result, error);
  } else if (fl_variant_is(razryad, "Крупное")) {
    *result = fl_number(0.0);
    return FL_OK;
  } else if (fl_variant_is(razryad, "Неизвестное")) {
    *result = fl_number(0.0);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Разряд обоснован».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param mesto — «место»: «Разряд»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value mesto, fl_value razryad, fl_value *result, fl_error *error) {
  fl_value fl_t385 = fl_nothing();
  FL_TRY(opis_diska_eto_neizvestnoe(ctx, mesto, &fl_t385, error));
  bool fl_t386 = false;
  FL_TRY(fl_cond(ctx, fl_t385, &fl_t386, error));
  fl_value fl_t387 = fl_nothing();
  if (fl_t386) {
    fl_t387 = fl_flag(false);
  } else {
    fl_t387 = fl_flag(true);
  }
  bool fl_t388 = false;
  FL_TRY(fl_cond(ctx, fl_t387, &fl_t388, error));
  if (fl_t388) {
    return opis_diska_tot_zhe_razryad(ctx, razryad, mesto, result, error);
  } else {
    return opis_diska_razryad_obosnovan_primetoy(ctx, nahodka, razryad, result, error);
  }
}

/*
 * Функция flang «Разряд обоснован приметой».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_obosnovan_primetoy(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error) {
  fl_value fl_t389 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t389, error));
  fl_value fl_t390 = fl_nothing();
  FL_TRY(opis_diska_primeta_kesha(ctx, fl_t389, &fl_t390, error));
  const fl_value kesh = fl_t390; /* пусть «кэш» */
  fl_value fl_t391 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t391, error));
  fl_value fl_t392 = fl_nothing();
  FL_TRY(opis_diska_primeta_zhurnala(ctx, fl_t391, &fl_t392, error));
  const fl_value zhurnal = fl_t392; /* пусть «журнал» */
  fl_value fl_t393 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t393, error));
  fl_value fl_t394 = fl_nothing();
  FL_TRY(opis_diska_primeta_sborki(ctx, fl_t393, &fl_t394, error));
  const fl_value sborka = fl_t394; /* пусть «сборка» */
  fl_value fl_t395 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t395, error));
  fl_value fl_t396 = fl_nothing();
  FL_TRY(opis_diska_primeta_zagruzki(ctx, fl_t395, &fl_t396, error));
  const fl_value zagruzka = fl_t396; /* пусть «загрузка» */
  fl_value fl_t397 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t397, error));
  fl_value fl_t398 = fl_nothing();
  FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t398, error));
  if (fl_t397.tag != FL_NUMBER || fl_t398.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t397, fl_t398, error));
  const fl_value krupnoe = fl_flag(fl_t397.as.number >= fl_t398.as.number); /* пусть «крупное» */
  if (fl_variant_is(razryad, "Кэш")) {
    *result = kesh;
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    bool fl_t399 = false;
    FL_TRY(fl_cond(ctx, zhurnal, &fl_t399, error));
    if (fl_t399) {
      bool fl_t400 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t400, error));
      if (fl_t400) {
        *result = fl_flag(false);
        return FL_OK;
      } else {
        *result = fl_flag(true);
        return FL_OK;
      }
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  } else if (fl_variant_is(razryad, "Сборка")) {
    bool fl_t401 = false;
    FL_TRY(fl_cond(ctx, sborka, &fl_t401, error));
    fl_value fl_t402 = fl_nothing();
    if (fl_t401) {
      bool fl_t403 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t403, error));
      fl_value fl_t404 = fl_nothing();
      if (fl_t403) {
        fl_t404 = fl_flag(false);
      } else {
        fl_t404 = fl_flag(true);
      }
      fl_t402 = fl_t404;
    } else {
      fl_t402 = fl_flag(false);
    }
    bool fl_t405 = false;
    FL_TRY(fl_cond(ctx, fl_t402, &fl_t405, error));
    if (fl_t405) {
      bool fl_t406 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t406, error));
      if (fl_t406) {
        *result = fl_flag(false);
        return FL_OK;
      } else {
        *result = fl_flag(true);
        return FL_OK;
      }
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  } else if (fl_variant_is(razryad, "Загрузка")) {
    bool fl_t407 = false;
    FL_TRY(fl_cond(ctx, zagruzka, &fl_t407, error));
    fl_value fl_t408 = fl_nothing();
    if (fl_t407) {
      bool fl_t409 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t409, error));
      fl_value fl_t410 = fl_nothing();
      if (fl_t409) {
        fl_t410 = fl_flag(false);
      } else {
        fl_t410 = fl_flag(true);
      }
      fl_t408 = fl_t410;
    } else {
      fl_t408 = fl_flag(false);
    }
    bool fl_t411 = false;
    FL_TRY(fl_cond(ctx, fl_t408, &fl_t411, error));
    fl_value fl_t412 = fl_nothing();
    if (fl_t411) {
      bool fl_t413 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t413, error));
      fl_value fl_t414 = fl_nothing();
      if (fl_t413) {
        fl_t414 = fl_flag(false);
      } else {
        fl_t414 = fl_flag(true);
      }
      fl_t412 = fl_t414;
    } else {
      fl_t412 = fl_flag(false);
    }
    bool fl_t415 = false;
    FL_TRY(fl_cond(ctx, fl_t412, &fl_t415, error));
    if (fl_t415) {
      bool fl_t416 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t416, error));
      if (fl_t416) {
        *result = fl_flag(false);
        return FL_OK;
      } else {
        *result = fl_flag(true);
        return FL_OK;
      }
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  } else if (fl_variant_is(razryad, "Крупное")) {
    bool fl_t417 = false;
    FL_TRY(fl_cond(ctx, krupnoe, &fl_t417, error));
    fl_value fl_t418 = fl_nothing();
    if (fl_t417) {
      bool fl_t419 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t419, error));
      fl_value fl_t420 = fl_nothing();
      if (fl_t419) {
        fl_t420 = fl_flag(false);
      } else {
        fl_t420 = fl_flag(true);
      }
      fl_t418 = fl_t420;
    } else {
      fl_t418 = fl_flag(false);
    }
    bool fl_t421 = false;
    FL_TRY(fl_cond(ctx, fl_t418, &fl_t421, error));
    fl_value fl_t422 = fl_nothing();
    if (fl_t421) {
      bool fl_t423 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t423, error));
      fl_value fl_t424 = fl_nothing();
      if (fl_t423) {
        fl_t424 = fl_flag(false);
      } else {
        fl_t424 = fl_flag(true);
      }
      fl_t422 = fl_t424;
    } else {
      fl_t422 = fl_flag(false);
    }
    bool fl_t425 = false;
    FL_TRY(fl_cond(ctx, fl_t422, &fl_t425, error));
    fl_value fl_t426 = fl_nothing();
    if (fl_t425) {
      bool fl_t427 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t427, error));
      fl_value fl_t428 = fl_nothing();
      if (fl_t427) {
        fl_t428 = fl_flag(false);
      } else {
        fl_t428 = fl_flag(true);
      }
      fl_t426 = fl_t428;
    } else {
      fl_t426 = fl_flag(false);
    }
    bool fl_t429 = false;
    FL_TRY(fl_cond(ctx, fl_t426, &fl_t429, error));
    if (fl_t429) {
      bool fl_t430 = false;
      FL_TRY(fl_cond(ctx, zagruzka, &fl_t430, error));
      if (fl_t430) {
        *result = fl_flag(false);
        return FL_OK;
      } else {
        *result = fl_flag(true);
        return FL_OK;
      }
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  } else if (fl_variant_is(razryad, "Неизвестное")) {
    bool fl_t431 = false;
    FL_TRY(fl_cond(ctx, kesh, &fl_t431, error));
    fl_value fl_t432 = fl_nothing();
    if (fl_t431) {
      fl_t432 = fl_flag(false);
    } else {
      fl_t432 = fl_flag(true);
    }
    bool fl_t433 = false;
    FL_TRY(fl_cond(ctx, fl_t432, &fl_t433, error));
    fl_value fl_t434 = fl_nothing();
    if (fl_t433) {
      bool fl_t435 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t435, error));
      fl_value fl_t436 = fl_nothing();
      if (fl_t435) {
        fl_t436 = fl_flag(false);
      } else {
        fl_t436 = fl_flag(true);
      }
      fl_t434 = fl_t436;
    } else {
      fl_t434 = fl_flag(false);
    }
    bool fl_t437 = false;
    FL_TRY(fl_cond(ctx, fl_t434, &fl_t437, error));
    fl_value fl_t438 = fl_nothing();
    if (fl_t437) {
      bool fl_t439 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t439, error));
      fl_value fl_t440 = fl_nothing();
      if (fl_t439) {
        fl_t440 = fl_flag(false);
      } else {
        fl_t440 = fl_flag(true);
      }
      fl_t438 = fl_t440;
    } else {
      fl_t438 = fl_flag(false);
    }
    bool fl_t441 = false;
    FL_TRY(fl_cond(ctx, fl_t438, &fl_t441, error));
    fl_value fl_t442 = fl_nothing();
    if (fl_t441) {
      bool fl_t443 = false;
      FL_TRY(fl_cond(ctx, zagruzka, &fl_t443, error));
      fl_value fl_t444 = fl_nothing();
      if (fl_t443) {
        fl_t444 = fl_flag(false);
      } else {
        fl_t444 = fl_flag(true);
      }
      fl_t442 = fl_t444;
    } else {
      fl_t442 = fl_flag(false);
    }
    bool fl_t445 = false;
    FL_TRY(fl_cond(ctx, fl_t442, &fl_t445, error));
    if (fl_t445) {
      bool fl_t446 = false;
      FL_TRY(fl_cond(ctx, krupnoe, &fl_t446, error));
      if (fl_t446) {
        *result = fl_flag(false);
        return FL_OK;
      } else {
        *result = fl_flag(true);
        return FL_OK;
      }
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  } else {
    return fl_match_fail(ctx, razryad, error);
  }
}

/*
 * Функция flang «Приговор обоснован».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_prigovor_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value prigovor, fl_value *result, fl_error *error) {
  fl_value fl_t447 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "доступен", &fl_t447, error));
  bool fl_t448 = false;
  FL_TRY(fl_cond(ctx, fl_t447, &fl_t448, error));
  fl_value fl_t449 = fl_nothing();
  if (fl_t448) {
    fl_t449 = fl_flag(false);
  } else {
    fl_t449 = fl_flag(true);
  }
  bool fl_t450 = false;
  FL_TRY(fl_cond(ctx, fl_t449, &fl_t450, error));
  if (fl_t450) {
    return opis_diska_eto_netrogat(ctx, prigovor, result, error);
  } else {
    fl_value fl_t451 = fl_nothing();
    FL_TRY(opis_diska_ssylka(ctx, nahodka, &fl_t451, error));
    bool fl_t452 = false;
    FL_TRY(fl_cond(ctx, fl_t451, &fl_t452, error));
    if (fl_t452) {
      return opis_diska_eto_netrogat(ctx, prigovor, result, error);
    } else {
      fl_value fl_t453 = fl_nothing();
      FL_TRY(opis_diska_eto_mozhnoubrat(ctx, prigovor, &fl_t453, error));
      bool fl_t454 = false;
      FL_TRY(fl_cond(ctx, fl_t453, &fl_t454, error));
      if (fl_t454) {
        fl_value fl_t455 = fl_nothing();
        FL_TRY(opis_diska_i1_na_pare(ctx, razryad, prigovor, &fl_t455, error));
        bool fl_t456 = false;
        FL_TRY(fl_cond(ctx, fl_t455, &fl_t456, error));
        fl_value fl_t457 = fl_nothing();
        if (fl_t456) {
          fl_value fl_t458 = fl_nothing();
          FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t458, error));
          bool fl_t459 = false;
          FL_TRY(fl_cond(ctx, fl_t458, &fl_t459, error));
          fl_value fl_t460 = fl_nothing();
          if (fl_t459) {
            fl_t460 = fl_flag(false);
          } else {
            fl_t460 = fl_flag(true);
          }
          fl_t457 = fl_t460;
        } else {
          fl_t457 = fl_flag(false);
        }
        bool fl_t461 = false;
        FL_TRY(fl_cond(ctx, fl_t457, &fl_t461, error));
        if (fl_t461) {
          fl_value fl_t462 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t462, error));
          fl_value fl_t463 = fl_nothing();
          FL_TRY(opis_diska_porog_razryada(ctx, razryad, &fl_t463, error));
          if (fl_t462.tag != FL_NUMBER || fl_t463.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t462, fl_t463, error));
          *result = fl_flag(fl_t462.as.number >= fl_t463.as.number);
          return FL_OK;
        } else {
          *result = fl_flag(false);
          return FL_OK;
        }
      } else {
        *result = fl_flag(true);
        return FL_OK;
      }
    }
  }
}

/*
 * Функция flang «Вес обоснован».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param prigovor — «приговор»: «Приговор»
 * @param ves — «вес»: число
 * @return значение
 */
fl_status opis_diska_ves_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value prigovor, fl_value ves, fl_value *result, fl_error *error) {
  fl_value fl_t464 = fl_nothing();
  FL_TRY(opis_diska_eto_netrogat(ctx, prigovor, &fl_t464, error));
  bool fl_t465 = false;
  FL_TRY(fl_cond(ctx, fl_t464, &fl_t465, error));
  if (fl_t465) {
    *result = fl_flag(fl_equal(ves, fl_number(0.0)));
    return FL_OK;
  } else {
    fl_value fl_t466 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t466, error));
    bool fl_t467 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_equal(ves, fl_t466)), &fl_t467, error));
    if (fl_t467) {
      return opis_diska_ves_v_granicah(ctx, nahodka, ves, result, error);
    } else {
      *result = fl_flag(false);
      return FL_OK;
    }
  }
}

/*
 * Вызов по исходному имени flang. Коды и тексты — те же, что у
 * интерпретатора: «не найдена функция …» и «функция … принимает N аргум.».
 */
fl_status opis_diska_call(fl_ctx *ctx, const char *name, const fl_value *args, size_t count,
                    fl_value *result, fl_error *error) {
  if (strcmp(name, "Порог крупного") == 0) {
    if (count != 0) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Порог крупного", (unsigned long)0, (unsigned long)count);
    }
    return opis_diska_porog_krupnogo(ctx, result, error);
  }
  if (strcmp(name, "Порог кэша") == 0) {
    if (count != 0) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Порог кэша", (unsigned long)0, (unsigned long)count);
    }
    return opis_diska_porog_kesha(ctx, result, error);
  }
  if (strcmp(name, "Порог журнала") == 0) {
    if (count != 0) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Порог журнала", (unsigned long)0, (unsigned long)count);
    }
    return opis_diska_porog_zhurnala(ctx, result, error);
  }
  if (strcmp(name, "Порог загрузки") == 0) {
    if (count != 0) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Порог загрузки", (unsigned long)0, (unsigned long)count);
    }
    return opis_diska_porog_zagruzki(ctx, result, error);
  }
  if (strcmp(name, "Составляющие пути") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Составляющие пути", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_sostavlyayuschie_puti(ctx, args[0], result, error);
  }
  if (strcmp(name, "Имя в пути") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Имя в пути", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_imya_v_puti(ctx, args[0], result, error);
  }
  if (strcmp(name, "Есть составляющая") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Есть составляющая", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_est_sostavlyayuschaya(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Оканчивается на") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Оканчивается на", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_okanchivaetsya_na(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "След пути") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "След пути", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_sled_puti(ctx, args[0], result, error);
  }
  if (strcmp(name, "Цепь ограничена") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Цепь ограничена", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_cep_ogranichena(ctx, args[0], result, error);
  }
  if (strcmp(name, "Справочник ограничен") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Справочник ограничен", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_spravochnik_ogranichen(ctx, args[0], result, error);
  }
  if (strcmp(name, "Разряд места допустим") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд места допустим", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_razryad_mesta_dopustim(ctx, args[0], result, error);
  }
  if (strcmp(name, "Место подходит") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Место подходит", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_mesto_podhodit(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Это неизвестное") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Это неизвестное", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_eto_neizvestnoe(ctx, args[0], result, error);
  }
  if (strcmp(name, "Номер разряда") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Номер разряда", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_nomer_razryada(ctx, args[0], result, error);
  }
  if (strcmp(name, "Тот же разряд") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Тот же разряд", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_tot_zhe_razryad(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Место обосновано") == 0) {
    if (count != 3) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Место обосновано", (unsigned long)3, (unsigned long)count);
    }
    return opis_diska_mesto_obosnovano(ctx, args[0], args[1], args[2], result, error);
  }
  if (strcmp(name, "Разряд по справочнику") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд по справочнику", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_razryad_po_spravochniku(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Шестнадцатеричный знак") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Шестнадцатеричный знак", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_shestnadcaterichnyy_znak(ctx, args[0], result, error);
  }
  if (strcmp(name, "Похоже на отпечаток") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Похоже на отпечаток", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_pohozhe_na_otpechatok(ctx, args[0], result, error);
  }
  if (strcmp(name, "Адресуется содержимым") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Адресуется содержимым", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_adresuetsya_soderzhimym(ctx, args[0], result, error);
  }
  if (strcmp(name, "Под системным временным") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Под системным временным", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_pod_sistemnym_vremennym(ctx, args[0], result, error);
  }
  if (strcmp(name, "Примета кэша") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Примета кэша", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_primeta_kesha(ctx, args[0], result, error);
  }
  if (strcmp(name, "Примета журнала") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Примета журнала", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_primeta_zhurnala(ctx, args[0], result, error);
  }
  if (strcmp(name, "Примета сборки") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Примета сборки", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_primeta_sborki(ctx, args[0], result, error);
  }
  if (strcmp(name, "Примета загрузки") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Примета загрузки", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_primeta_zagruzki(ctx, args[0], result, error);
  }
  if (strcmp(name, "Есть примета") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Есть примета", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_est_primeta(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Разряд решён размером") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд решён размером", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_razryad_reshyon_razmerom(ctx, args[0], result, error);
  }
  if (strcmp(name, "Разряд находки") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд находки", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_razryad_nahodki(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Крупное не мельче порога") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Крупное не мельче порога", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_krupnoe_ne_melche_poroga(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Каталог") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Каталог", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_katalog(ctx, args[0], result, error);
  }
  if (strcmp(name, "Ссылка") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Ссылка", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_ssylka(ctx, args[0], result, error);
  }
  if (strcmp(name, "Приговор мусора") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Приговор мусора", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_prigovor_musora(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Приговор находки") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Приговор находки", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_prigovor_nahodki(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Вес находки") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Вес находки", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_ves_nahodki(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Вес в границах") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Вес в границах", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_ves_v_granicah(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Решить находку") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Решить находку", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_reshit_nahodku(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Решить всё") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Решить всё", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_reshit_vsyo(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "И1 держится") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "И1 держится", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_i1_derzhitsya(ctx, args[0], result, error);
  }
  if (strcmp(name, "И1 держится всюду") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "И1 держится всюду", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_i1_derzhitsya_vsyudu(ctx, args[0], result, error);
  }
  if (strcmp(name, "Пустой свод") == 0) {
    if (count != 0) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Пустой свод", (unsigned long)0, (unsigned long)count);
    }
    return opis_diska_pustoy_svod(ctx, result, error);
  }
  if (strcmp(name, "Прибавить решение") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Прибавить решение", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_pribavit_reshenie(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Свести") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Свести", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_svesti(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Сумма размеров убираемых") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Сумма размеров убираемых", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_summa_razmerov_ubiraemyh(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "И2 держится") == 0) {
    if (count != 3) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "И2 держится", (unsigned long)3, (unsigned long)count);
    }
    return opis_diska_i2_derzhitsya(ctx, args[0], args[1], args[2], result, error);
  }
  if (strcmp(name, "Строку отчёта") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Строку отчёта", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_stroku_otchyota(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Вставить по весу") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Вставить по весу", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_vstavit_po_vesu(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Приписать строку отчёта") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Приписать строку отчёта", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_pripisat_stroku_otchyota(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Отчёт") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Отчёт", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_otchyot(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Отчёт той же длины") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Отчёт той же длины", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_otchyot_toy_zhe_dliny(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Это МожноУбрать") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Это МожноУбрать", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_eto_mozhnoubrat(ctx, args[0], result, error);
  }
  if (strcmp(name, "Это НеТрогать") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Это НеТрогать", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_eto_netrogat(ctx, args[0], result, error);
  }
  if (strcmp(name, "И1 на паре") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "И1 на паре", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_i1_na_pare(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Порог разряда") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Порог разряда", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_porog_razryada(ctx, args[0], result, error);
  }
  if (strcmp(name, "Разряд обоснован") == 0) {
    if (count != 3) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд обоснован", (unsigned long)3, (unsigned long)count);
    }
    return opis_diska_razryad_obosnovan(ctx, args[0], args[1], args[2], result, error);
  }
  if (strcmp(name, "Разряд обоснован приметой") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд обоснован приметой", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_razryad_obosnovan_primetoy(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Приговор обоснован") == 0) {
    if (count != 3) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Приговор обоснован", (unsigned long)3, (unsigned long)count);
    }
    return opis_diska_prigovor_obosnovan(ctx, args[0], args[1], args[2], result, error);
  }
  if (strcmp(name, "Вес обоснован") == 0) {
    if (count != 3) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Вес обоснован", (unsigned long)3, (unsigned long)count);
    }
    return opis_diska_ves_obosnovan(ctx, args[0], args[1], args[2], result, error);
  }
  return fl_fail(ctx, error, FL_CODE_UNKNOWN_NAME, "не найдена функция «%s»", name);
}

/*
 * ТА ЖЕ ДВЕРЬ, НО С ГРАНИЦЕЙ ВХОДА: сначала объявленные типы параметров
 * (fl_check_entry по таблице внизу файла), потом вызов. Зовите ЭТУ, если
 * значения пришли снаружи — из JSON, из другого языка, от человека.
 *
 * Почему не сверяет сам `_call`. Он обязан отвечать значение в значение так
 * же, как `interpret` у свидетеля, а тот объявленных типов не сверяет тоже:
 * сверяет их `flang run`. Здесь ровно та же пара — `_call` вычислитель,
 * `_enter` дверь, — и разойдись они, у языка стало бы два ответа на вопрос
 * «подходит ли значение типу».
 */
fl_status opis_diska_enter(fl_ctx *ctx, const char *name, const fl_value *args, size_t count,
                    fl_value *result, fl_error *error) {
  FL_TRY(fl_check_entry(ctx, opis_diska_entry(), name, args, count, error));
  return opis_diska_call(ctx, name, args, count, result, error);
}

/*
 * Граница входа: объявленные типы параметров данными. Прогонщик сверяет по
 * ним значения, пришедшие снаружи, ДО вызова (fl_check_entry).
 *
 * Виды `неизвестно` (значение-функция, параметр полиморфизма, применение
 * типа с аргументами) не сверяются — ровно как молчит о них проверка
 * значений свидетеля.
 */
static const fl_type_field opis_diska_entry_fields[] = {
  { "разряд", 3 },
  { "якорь", 4 },
  { "цепь", 0 },
  { "путь", 0 },
  { "размер", 6 },
  { "возраст_дней", 6 },
  { "вид", 7 },
  { "доступен", 8 },
  { "разряд", 3 },
  { "приговор", 9 },
  { "вес", 6 },
  { "кэш", 6 },
  { "журнал", 6 },
  { "сборка", 6 },
  { "загрузка", 6 },
  { "крупное", 6 },
  { "неизвестное", 6 },
  { "освободить", 6 },
  { "путь", 0 },
  { "решение", 11 },
};

static const fl_type_variant opis_diska_entry_variants[] = {
  { "Кэш", 0, 0 },
  { "Журнал", 0, 0 },
  { "Сборка", 0, 0 },
  { "Загрузка", 0, 0 },
  { "Крупное", 0, 0 },
  { "Неизвестное", 0, 0 },
  { "ОтКорня", 0, 0 },
  { "ГдеУгодно", 0, 0 },
  { "Файл", 3, 0 },
  { "Каталог", 3, 0 },
  { "Ссылка", 3, 0 },
  { "МожноУбрать", 8, 0 },
  { "Спросить", 8, 0 },
  { "НеТрогать", 8, 0 },
};

static const fl_type opis_diska_entry_types[] = {
  { FL_TYPE_STRING, "строка", "", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 0 },
  { FL_TYPE_LIST, "список «Место»", "", false, false, false, 0.0, 0.0, 2, 0, 0, 0, 0 },
  { FL_TYPE_RECORD, "«Место»", "Место", false, false, false, 0.0, 0.0, 0, 0, 3, 0, 0 },
  { FL_TYPE_SUM, "«Разряд»", "Разряд", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 6 },
  { FL_TYPE_SUM, "«Якорь»", "Якорь", false, false, false, 0.0, 0.0, 0, 0, 0, 6, 2 },
  { FL_TYPE_RECORD, "«Находка»", "Находка", false, false, false, 0.0, 0.0, 0, 3, 5, 0, 0 },
  { FL_TYPE_NUMBER, "число", "", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 0 },
  { FL_TYPE_SUM, "«Вид»", "Вид", false, false, false, 0.0, 0.0, 0, 0, 0, 8, 3 },
  { FL_TYPE_FLAG, "признак", "", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 0 },
  { FL_TYPE_SUM, "«Приговор»", "Приговор", false, false, false, 0.0, 0.0, 0, 0, 0, 11, 3 },
  { FL_TYPE_LIST, "список «Находка»", "", false, false, false, 0.0, 0.0, 5, 0, 0, 0, 0 },
  { FL_TYPE_RECORD, "«Решение»", "Решение", false, false, false, 0.0, 0.0, 0, 8, 3, 0, 0 },
  { FL_TYPE_LIST, "список «Решение»", "", false, false, false, 0.0, 0.0, 11, 0, 0, 0, 0 },
  { FL_TYPE_RECORD, "«Свод»", "Свод", false, false, false, 0.0, 0.0, 0, 11, 7, 0, 0 },
  { FL_TYPE_RECORD, "«Строка отчёта»", "Строка отчёта", false, false, false, 0.0, 0.0, 0, 18, 2, 0, 0 },
  { FL_TYPE_LIST, "список «Строка отчёта»", "", false, false, false, 0.0, 0.0, 14, 0, 0, 0, 0 },
};

static const fl_entry_param opis_diska_entry_params[] = {
  { "Составляющие пути", "путь", 0 },
  { "Имя в пути", "путь", 0 },
  { "Есть составляющая", "путь", 0 },
  { "Есть составляющая", "имя", 0 },
  { "Оканчивается на", "текст", 0 },
  { "Оканчивается на", "хвост", 0 },
  { "След пути", "путь", 0 },
  { "Цепь ограничена", "цепь", 0 },
  { "Справочник ограничен", "справочник", 1 },
  { "Разряд места допустим", "разряд", 3 },
  { "Место подходит", "след", 0 },
  { "Место подходит", "место", 2 },
  { "Это неизвестное", "разряд", 3 },
  { "Номер разряда", "разряд", 3 },
  { "Тот же разряд", "первый", 3 },
  { "Тот же разряд", "второй", 3 },
  { "Место обосновано", "путь", 0 },
  { "Место обосновано", "справочник", 1 },
  { "Место обосновано", "разряд", 3 },
  { "Разряд по справочнику", "путь", 0 },
  { "Разряд по справочнику", "справочник", 1 },
  { "Шестнадцатеричный знак", "знак", 0 },
  { "Похоже на отпечаток", "часть", 0 },
  { "Адресуется содержимым", "путь", 0 },
  { "Под системным временным", "путь", 0 },
  { "Примета кэша", "путь", 0 },
  { "Примета журнала", "путь", 0 },
  { "Примета сборки", "путь", 0 },
  { "Примета загрузки", "путь", 0 },
  { "Есть примета", "путь", 0 },
  { "Есть примета", "место", 3 },
  { "Разряд решён размером", "разряд", 3 },
  { "Разряд находки", "находка", 5 },
  { "Разряд находки", "место", 3 },
  { "Крупное не мельче порога", "находка", 5 },
  { "Крупное не мельче порога", "разряд", 3 },
  { "Каталог", "находка", 5 },
  { "Ссылка", "находка", 5 },
  { "Приговор мусора", "находка", 5 },
  { "Приговор мусора", "порог", 6 },
  { "Приговор находки", "находка", 5 },
  { "Приговор находки", "разряд", 3 },
  { "Вес находки", "находка", 5 },
  { "Вес находки", "приговор", 9 },
  { "Вес в границах", "находка", 5 },
  { "Вес в границах", "вес", 6 },
  { "Решить находку", "находка", 5 },
  { "Решить находку", "справочник", 1 },
  { "Решить всё", "записи", 10 },
  { "Решить всё", "справочник", 1 },
  { "И1 держится", "решение", 11 },
  { "И1 держится всюду", "решения", 12 },
  { "Прибавить решение", "свод", 13 },
  { "Прибавить решение", "решение", 11 },
  { "Свести", "записи", 10 },
  { "Свести", "справочник", 1 },
  { "Сумма размеров убираемых", "записи", 10 },
  { "Сумма размеров убираемых", "справочник", 1 },
  { "И2 держится", "записи", 10 },
  { "И2 держится", "справочник", 1 },
  { "И2 держится", "свод", 13 },
  { "Строку отчёта", "находка", 5 },
  { "Строку отчёта", "справочник", 1 },
  { "Вставить по весу", "строка", 14 },
  { "Вставить по весу", "строки", 15 },
  { "Приписать строку отчёта", "первая", 14 },
  { "Приписать строку отчёта", "строки", 15 },
  { "Отчёт", "записи", 10 },
  { "Отчёт", "справочник", 1 },
  { "Отчёт той же длины", "записи", 10 },
  { "Отчёт той же длины", "строки", 15 },
  { "Это МожноУбрать", "приговор", 9 },
  { "Это НеТрогать", "приговор", 9 },
  { "И1 на паре", "разряд", 3 },
  { "И1 на паре", "приговор", 9 },
  { "Порог разряда", "разряд", 3 },
  { "Разряд обоснован", "находка", 5 },
  { "Разряд обоснован", "место", 3 },
  { "Разряд обоснован", "разряд", 3 },
  { "Разряд обоснован приметой", "находка", 5 },
  { "Разряд обоснован приметой", "разряд", 3 },
  { "Приговор обоснован", "находка", 5 },
  { "Приговор обоснован", "разряд", 3 },
  { "Приговор обоснован", "приговор", 9 },
  { "Вес обоснован", "находка", 5 },
  { "Вес обоснован", "приговор", 9 },
  { "Вес обоснован", "вес", 6 },
};

static const fl_entry_table opis_diska_entry_table = {
  opis_diska_entry_types, 16,
  opis_diska_entry_fields, 20,
  opis_diska_entry_variants, 14,
  opis_diska_entry_params, 87
};

const fl_entry_table *opis_diska_entry(void) {
  return &opis_diska_entry_table;
}
