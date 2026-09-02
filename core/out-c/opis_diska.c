/*
 * Сгенерировано flang (бэкенд C, flang/self/emit-c.flang). Не редактировать руками.
 * Модуль flang: «Опись диска».
 * Файл: реализация.
 * Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.
 */
#include "opis_diska.h"

#include <string.h>

/* Константы программы: имена полей и строковые литералы. */
static const char *const opis_diska_names_1[] = { "путь", "размер", "возраст_дней", "вид", "доступен" };
static const char *const opis_diska_names_2[] = { "разряд", "приговор", "вес" };
static const char *const opis_diska_names_3[] = { "путь", "решение" };
static const char *const opis_diska_names_4[] = { "кэш", "журнал", "сборка", "загрузка", "крупное", "неизвестное", "освободить" };
static const fl_value opis_diska_text_5 = { FL_STRING, { .string = { "/", 1, 1 } } };
static const fl_value opis_diska_text_6 = { FL_STRING, { .string = { "", 0, 0 } } };
static const fl_value opis_diska_text_7 = { FL_STRING, { .string = { "0123456789abcdefABCDEF", 22, 22 } } };
static const fl_value opis_diska_text_8 = { FL_STRING, { .string = { "site-packages", 13, 13 } } };
static const fl_value opis_diska_text_9 = { FL_STRING, { .string = { "dist-packages", 13, 13 } } };
static const fl_value opis_diska_text_10 = { FL_STRING, { .string = { ".git", 4, 4 } } };
static const fl_value opis_diska_text_11 = { FL_STRING, { .string = { "/tmp/", 5, 5 } } };
static const fl_value opis_diska_text_12 = { FL_STRING, { .string = { "/var/tmp/", 9, 9 } } };
static const fl_value opis_diska_text_13 = { FL_STRING, { .string = { ".cache", 6, 6 } } };
static const fl_value opis_diska_text_14 = { FL_STRING, { .string = { "cache", 5, 5 } } };
static const fl_value opis_diska_text_15 = { FL_STRING, { .string = { "Caches", 6, 6 } } };
static const fl_value opis_diska_text_16 = { FL_STRING, { .string = { ".log", 4, 4 } } };
static const fl_value opis_diska_text_17 = { FL_STRING, { .string = { "log", 3, 3 } } };
static const fl_value opis_diska_text_18 = { FL_STRING, { .string = { "logs", 4, 4 } } };
static const fl_value opis_diska_text_19 = { FL_STRING, { .string = { "node_modules", 12, 12 } } };
static const fl_value opis_diska_text_20 = { FL_STRING, { .string = { "target", 6, 6 } } };
static const fl_value opis_diska_text_21 = { FL_STRING, { .string = { "build", 5, 5 } } };
static const fl_value opis_diska_text_22 = { FL_STRING, { .string = { "_build", 6, 6 } } };
static const fl_value opis_diska_text_23 = { FL_STRING, { .string = { ".gradle", 7, 7 } } };
static const fl_value opis_diska_text_24 = { FL_STRING, { .string = { "Downloads", 9, 9 } } };
static const fl_value opis_diska_text_25 = { FL_STRING, { .string = { "Загрузки", 16, 8 } } };


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

/* Фабрика записи FTS «Решение». */
fl_status opis_diska_sozdat_reshenie(fl_ctx *ctx, fl_value razryad, fl_value prigovor, fl_value ves, fl_value *out, fl_error *error) {
  fl_value values[3];
  values[0] = razryad; /* «разряд» */
  values[1] = prigovor; /* «приговор» */
  values[2] = ves; /* «вес» */
  return fl_record_new(ctx, opis_diska_names_2, values, 3, out, error);
}

/* Фабрика записи FTS «Строка отчёта». */
fl_status opis_diska_sozdat_stroka_otchyota(fl_ctx *ctx, fl_value put, fl_value reshenie, fl_value *out, fl_error *error) {
  fl_value values[2];
  values[0] = put; /* «путь» */
  values[1] = reshenie; /* «решение» */
  return fl_record_new(ctx, opis_diska_names_3, values, 2, out, error);
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
  return fl_record_new(ctx, opis_diska_names_4, values, 7, out, error);
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
  FL_TRY(fl_b_razdelit_dokazano(ctx, put, opis_diska_text_5, &fl_t1, error));
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
  fl_value sobrano = opis_diska_text_6; /* «собрано» */
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
 * Функция flang «Шестнадцатеричный знак».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param znak — «знак»: строка
 * @return значение
 */
fl_status opis_diska_shestnadcaterichnyy_znak(fl_ctx *ctx, fl_value znak, fl_value *result, fl_error *error) {
  fl_value fl_t15 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, opis_diska_text_7, znak, &fl_t15, error));
  *result = fl_t15;
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
  fl_value fl_t16 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, chast, &fl_t16, error));
  if (fl_t16.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t16, fl_number(32.0), error));
  bool fl_t17 = false;
  FL_TRY(fl_cond(ctx, fl_flag(fl_t16.as.number < 32.0), &fl_t17, error));
  if (fl_t17) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    fl_value fl_t18 = fl_nothing(); /* «символы» */
    FL_TRY(fl_b_simvoly(ctx, chast, &fl_t18, error));
    fl_value fl_t19 = fl_nothing();
    FL_TRY(fl_require_list(ctx, fl_t18, "свёртка", &fl_t19, error));
    fl_value sobrano = fl_flag(true); /* «собрано» */
    const fl_mark fl_t21 = fl_region_open(ctx);
    for (size_t fl_t20 = 0; fl_t20 < fl_t19.as.list.count; fl_t20 += 1) {
      const fl_value znak = fl_t19.as.list.items[fl_t20]; /* «знак» */
      bool fl_t22 = false;
      FL_TRY(fl_cond(ctx, sobrano, &fl_t22, error));
      fl_value fl_t23 = fl_nothing();
      if (fl_t22) {
        fl_value fl_t24 = fl_nothing();
        FL_TRY(opis_diska_shestnadcaterichnyy_znak(ctx, znak, &fl_t24, error));
        fl_t23 = fl_t24;
      } else {
        fl_t23 = fl_flag(false);
      }
      sobrano = fl_t23;
      FL_TRY(fl_region_recycle(ctx, fl_t21, &sobrano, error));
    }
    FL_TRY(fl_region_close(ctx, fl_t21, FL_OK, &sobrano, error));
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
  fl_value fl_t25 = fl_nothing();
  FL_TRY(opis_diska_sostavlyayuschie_puti(ctx, put, &fl_t25, error));
  fl_value fl_t26 = fl_nothing();
  FL_TRY(fl_require_list(ctx, fl_t25, "отфильтровать", &fl_t26, error));
  fl_value *fl_t27 = NULL;
  size_t fl_t28 = 0;
  FL_TRY(fl_list_alloc(ctx, fl_t26.as.list.count, &fl_t27, error));
  for (size_t fl_t29 = 0; fl_t29 < fl_t26.as.list.count; fl_t29 += 1) {
    const fl_value chast = fl_t26.as.list.items[fl_t29]; /* «часть» */
    fl_value fl_t30 = fl_nothing();
    FL_TRY(opis_diska_pohozhe_na_otpechatok(ctx, chast, &fl_t30, error));
    bool fl_t31 = false;
    FL_TRY(fl_keep(ctx, fl_t30, &fl_t31, error));
    if (fl_t31) {
      fl_t27[fl_t28] = chast;
      fl_t28 += 1;
    }
  }
  fl_value fl_t32 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, fl_list(fl_t27, fl_t28), &fl_t32, error));
  if (fl_t32.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t32, fl_number(0.0), error));
  bool fl_t33 = false;
  FL_TRY(fl_cond(ctx, fl_flag(fl_t32.as.number > 0.0), &fl_t33, error));
  fl_value fl_t34 = fl_nothing();
  if (fl_t33) {
    fl_t34 = fl_flag(true);
  } else {
    fl_value fl_t35 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_8, &fl_t35, error));
    fl_t34 = fl_t35;
  }
  bool fl_t36 = false;
  FL_TRY(fl_cond(ctx, fl_t34, &fl_t36, error));
  fl_value fl_t37 = fl_nothing();
  if (fl_t36) {
    fl_t37 = fl_flag(true);
  } else {
    fl_value fl_t38 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_9, &fl_t38, error));
    fl_t37 = fl_t38;
  }
  bool fl_t39 = false;
  FL_TRY(fl_cond(ctx, fl_t37, &fl_t39, error));
  if (fl_t39) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_10, result, error);
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
  fl_value fl_t40 = fl_nothing(); /* «начинается с» */
  FL_TRY(fl_b_nachinaetsya_s(ctx, put, opis_diska_text_11, &fl_t40, error));
  bool fl_t41 = false;
  FL_TRY(fl_cond(ctx, fl_t40, &fl_t41, error));
  if (fl_t41) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t42 = fl_nothing(); /* «начинается с» */
    FL_TRY(fl_b_nachinaetsya_s(ctx, put, opis_diska_text_12, &fl_t42, error));
    *result = fl_t42;
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
  fl_value fl_t43 = fl_nothing();
  FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_13, &fl_t43, error));
  bool fl_t44 = false;
  FL_TRY(fl_cond(ctx, fl_t43, &fl_t44, error));
  fl_value fl_t45 = fl_nothing();
  if (fl_t44) {
    fl_t45 = fl_flag(true);
  } else {
    fl_value fl_t46 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_14, &fl_t46, error));
    fl_t45 = fl_t46;
  }
  bool fl_t47 = false;
  FL_TRY(fl_cond(ctx, fl_t45, &fl_t47, error));
  fl_value fl_t48 = fl_nothing();
  if (fl_t47) {
    fl_t48 = fl_flag(true);
  } else {
    fl_value fl_t49 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_15, &fl_t49, error));
    fl_t48 = fl_t49;
  }
  bool fl_t50 = false;
  FL_TRY(fl_cond(ctx, fl_t48, &fl_t50, error));
  if (fl_t50) {
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
  fl_value fl_t51 = fl_nothing();
  FL_TRY(opis_diska_imya_v_puti(ctx, put, &fl_t51, error));
  fl_value fl_t52 = fl_nothing();
  FL_TRY(opis_diska_okanchivaetsya_na(ctx, fl_t51, opis_diska_text_16, &fl_t52, error));
  bool fl_t53 = false;
  FL_TRY(fl_cond(ctx, fl_t52, &fl_t53, error));
  fl_value fl_t54 = fl_nothing();
  if (fl_t53) {
    fl_t54 = fl_flag(true);
  } else {
    fl_value fl_t55 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_17, &fl_t55, error));
    fl_t54 = fl_t55;
  }
  bool fl_t56 = false;
  FL_TRY(fl_cond(ctx, fl_t54, &fl_t56, error));
  if (fl_t56) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_18, result, error);
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
  fl_value fl_t57 = fl_nothing();
  FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_19, &fl_t57, error));
  bool fl_t58 = false;
  FL_TRY(fl_cond(ctx, fl_t57, &fl_t58, error));
  fl_value fl_t59 = fl_nothing();
  if (fl_t58) {
    fl_t59 = fl_flag(true);
  } else {
    fl_value fl_t60 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_20, &fl_t60, error));
    fl_t59 = fl_t60;
  }
  bool fl_t61 = false;
  FL_TRY(fl_cond(ctx, fl_t59, &fl_t61, error));
  fl_value fl_t62 = fl_nothing();
  if (fl_t61) {
    fl_t62 = fl_flag(true);
  } else {
    fl_value fl_t63 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_21, &fl_t63, error));
    fl_t62 = fl_t63;
  }
  bool fl_t64 = false;
  FL_TRY(fl_cond(ctx, fl_t62, &fl_t64, error));
  fl_value fl_t65 = fl_nothing();
  if (fl_t64) {
    fl_t65 = fl_flag(true);
  } else {
    fl_value fl_t66 = fl_nothing();
    FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_22, &fl_t66, error));
    fl_t65 = fl_t66;
  }
  bool fl_t67 = false;
  FL_TRY(fl_cond(ctx, fl_t65, &fl_t67, error));
  if (fl_t67) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_23, result, error);
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
  fl_value fl_t68 = fl_nothing();
  FL_TRY(opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_24, &fl_t68, error));
  bool fl_t69 = false;
  FL_TRY(fl_cond(ctx, fl_t68, &fl_t69, error));
  if (fl_t69) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return opis_diska_est_sostavlyayuschaya(ctx, put, opis_diska_text_25, result, error);
  }
}

/*
 * Функция flang «Есть примета».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_est_primeta(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t70 = fl_nothing();
  FL_TRY(opis_diska_primeta_kesha(ctx, put, &fl_t70, error));
  bool fl_t71 = false;
  FL_TRY(fl_cond(ctx, fl_t70, &fl_t71, error));
  fl_value fl_t72 = fl_nothing();
  if (fl_t71) {
    fl_t72 = fl_flag(true);
  } else {
    fl_value fl_t73 = fl_nothing();
    FL_TRY(opis_diska_primeta_zhurnala(ctx, put, &fl_t73, error));
    fl_t72 = fl_t73;
  }
  bool fl_t74 = false;
  FL_TRY(fl_cond(ctx, fl_t72, &fl_t74, error));
  fl_value fl_t75 = fl_nothing();
  if (fl_t74) {
    fl_t75 = fl_flag(true);
  } else {
    fl_value fl_t76 = fl_nothing();
    FL_TRY(opis_diska_primeta_sborki(ctx, put, &fl_t76, error));
    fl_t75 = fl_t76;
  }
  bool fl_t77 = false;
  FL_TRY(fl_cond(ctx, fl_t75, &fl_t77, error));
  if (fl_t77) {
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
 * @return значение: «Разряд»
 */
fl_status opis_diska_razryad_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error) {
  fl_value fl_t78 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t78, error));
  fl_value fl_t79 = fl_nothing();
  FL_TRY(opis_diska_primeta_kesha(ctx, fl_t78, &fl_t79, error));
  bool fl_t80 = false;
  FL_TRY(fl_cond(ctx, fl_t79, &fl_t80, error));
  fl_value fl_t81 = fl_nothing();
  if (fl_t80) {
    fl_value fl_t82 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "Кэш", NULL, NULL, 0, &fl_t82, error));
    fl_t81 = fl_t82;
  } else {
    fl_value fl_t83 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t83, error));
    fl_value fl_t84 = fl_nothing();
    FL_TRY(opis_diska_primeta_zhurnala(ctx, fl_t83, &fl_t84, error));
    bool fl_t85 = false;
    FL_TRY(fl_cond(ctx, fl_t84, &fl_t85, error));
    fl_value fl_t86 = fl_nothing();
    if (fl_t85) {
      fl_value fl_t87 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "Журнал", NULL, NULL, 0, &fl_t87, error));
      fl_t86 = fl_t87;
    } else {
      fl_value fl_t88 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t88, error));
      fl_value fl_t89 = fl_nothing();
      FL_TRY(opis_diska_primeta_sborki(ctx, fl_t88, &fl_t89, error));
      bool fl_t90 = false;
      FL_TRY(fl_cond(ctx, fl_t89, &fl_t90, error));
      fl_value fl_t91 = fl_nothing();
      if (fl_t90) {
        fl_value fl_t92 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "Сборка", NULL, NULL, 0, &fl_t92, error));
        fl_t91 = fl_t92;
      } else {
        fl_value fl_t93 = fl_nothing();
        FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t93, error));
        fl_value fl_t94 = fl_nothing();
        FL_TRY(opis_diska_primeta_zagruzki(ctx, fl_t93, &fl_t94, error));
        bool fl_t95 = false;
        FL_TRY(fl_cond(ctx, fl_t94, &fl_t95, error));
        fl_value fl_t96 = fl_nothing();
        if (fl_t95) {
          fl_value fl_t97 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "Загрузка", NULL, NULL, 0, &fl_t97, error));
          fl_t96 = fl_t97;
        } else {
          fl_value fl_t98 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t98, error));
          fl_value fl_t99 = fl_nothing();
          FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t99, error));
          if (fl_t98.tag != FL_NUMBER || fl_t99.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t98, fl_t99, error));
          bool fl_t100 = false;
          FL_TRY(fl_cond(ctx, fl_flag(fl_t98.as.number >= fl_t99.as.number), &fl_t100, error));
          fl_value fl_t101 = fl_nothing();
          if (fl_t100) {
            fl_value fl_t102 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Крупное", NULL, NULL, 0, &fl_t102, error));
            fl_t101 = fl_t102;
          } else {
            fl_value fl_t103 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Неизвестное", NULL, NULL, 0, &fl_t103, error));
            fl_t101 = fl_t103;
          }
          fl_t96 = fl_t101;
        }
        fl_t91 = fl_t96;
      }
      fl_t86 = fl_t91;
    }
    fl_t81 = fl_t86;
  }
  const fl_value fl_t104 = fl_t81;
  fl_value fl_t105 = fl_nothing();
  FL_TRY(opis_diska_razryad_obosnovan(ctx, nahodka, fl_t104, &fl_t105, error));
  /* постусловие «Разряд обоснован приметой» */
  bool fl_t106 = false;
  FL_TRY(fl_post(ctx, fl_t105, "Разряд обоснован приметой", "Разряд находки", &fl_t106, error));
  if (!fl_t106) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Разряд обоснован приметой» функции «Разряд находки»");
  }
  fl_value fl_t107 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t107, error));
  fl_value fl_t108 = fl_nothing();
  FL_TRY(opis_diska_est_primeta(ctx, fl_t107, &fl_t108, error));
  bool fl_t109 = false;
  FL_TRY(fl_cond(ctx, fl_t108, &fl_t109, error));
  fl_value fl_t110 = fl_nothing();
  if (fl_t109) {
    fl_t110 = fl_flag(false);
  } else {
    fl_t110 = fl_flag(true);
  }
  bool fl_t111 = false;
  FL_TRY(fl_cond(ctx, fl_t110, &fl_t111, error));
  fl_value fl_t112 = fl_nothing();
  if (fl_t111) {
    fl_value fl_t113 = fl_nothing();
    FL_TRY(opis_diska_razryad_reshyon_razmerom(ctx, fl_t104, &fl_t113, error));
    fl_t112 = fl_t113;
  } else {
    fl_t112 = fl_flag(true);
  }
  /* постусловие «И4: без приметы-составляющей разряд решает размер» */
  bool fl_t114 = false;
  FL_TRY(fl_post(ctx, fl_t112, "И4: без приметы-составляющей разряд решает размер", "Разряд находки", &fl_t114, error));
  if (!fl_t114) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И4: без приметы-составляющей разряд решает размер» функции «Разряд находки»");
  }
  *result = fl_t104;
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
    fl_value fl_t115 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t115, error));
    fl_value fl_t116 = fl_nothing();
    FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t116, error));
    if (fl_t115.tag != FL_NUMBER || fl_t116.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t115, fl_t116, error));
    *result = fl_flag(fl_t115.as.number >= fl_t116.as.number);
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
  fl_value fl_t117 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "вид", &fl_t117, error));
  if (fl_variant_is(fl_t117, "Каталог")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t117, "Файл")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(fl_t117, "Ссылка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t117, error);
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
  fl_value fl_t118 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "вид", &fl_t118, error));
  if (fl_variant_is(fl_t118, "Ссылка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t118, "Файл")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(fl_t118, "Каталог")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t118, error);
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
  fl_value fl_t119 = fl_nothing();
  FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t119, error));
  bool fl_t120 = false;
  FL_TRY(fl_cond(ctx, fl_t119, &fl_t120, error));
  fl_value fl_t121 = fl_nothing();
  if (fl_t120) {
    fl_value fl_t122 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t122, error));
    fl_t121 = fl_t122;
  } else {
    fl_value fl_t123 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t123, error));
    if (fl_t123.tag != FL_NUMBER || porog.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t123, porog, error));
    bool fl_t124 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_t123.as.number >= porog.as.number), &fl_t124, error));
    fl_value fl_t125 = fl_nothing();
    if (fl_t124) {
      fl_value fl_t126 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "МожноУбрать", NULL, NULL, 0, &fl_t126, error));
      fl_t125 = fl_t126;
    } else {
      fl_value fl_t127 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t127, error));
      fl_t125 = fl_t127;
    }
    fl_t121 = fl_t125;
  }
  const fl_value fl_t128 = fl_t121;
  fl_value fl_t129 = fl_nothing();
  FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t129, error));
  bool fl_t130 = false;
  FL_TRY(fl_cond(ctx, fl_t129, &fl_t130, error));
  fl_value fl_t131 = fl_nothing();
  if (fl_t130) {
    fl_value fl_t132 = fl_nothing();
    FL_TRY(opis_diska_eto_mozhnoubrat(ctx, fl_t128, &fl_t132, error));
    bool fl_t133 = false;
    FL_TRY(fl_cond(ctx, fl_t132, &fl_t133, error));
    fl_value fl_t134 = fl_nothing();
    if (fl_t133) {
      fl_t134 = fl_flag(false);
    } else {
      fl_t134 = fl_flag(true);
    }
    fl_t131 = fl_t134;
  } else {
    fl_t131 = fl_flag(true);
  }
  /* постусловие «Каталог не убирается» */
  bool fl_t135 = false;
  FL_TRY(fl_post(ctx, fl_t131, "Каталог не убирается", "Приговор мусора", &fl_t135, error));
  if (!fl_t135) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Каталог не убирается» функции «Приговор мусора»");
  }
  *result = fl_t128;
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
  fl_value fl_t136 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "доступен", &fl_t136, error));
  bool fl_t137 = false;
  FL_TRY(fl_cond(ctx, fl_t136, &fl_t137, error));
  fl_value fl_t138 = fl_nothing();
  if (fl_t137) {
    fl_t138 = fl_flag(false);
  } else {
    fl_t138 = fl_flag(true);
  }
  bool fl_t139 = false;
  FL_TRY(fl_cond(ctx, fl_t138, &fl_t139, error));
  fl_value fl_t140 = fl_nothing();
  if (fl_t139) {
    fl_value fl_t141 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t141, error));
    fl_t140 = fl_t141;
  } else {
    fl_value fl_t142 = fl_nothing();
    FL_TRY(opis_diska_ssylka(ctx, nahodka, &fl_t142, error));
    bool fl_t143 = false;
    FL_TRY(fl_cond(ctx, fl_t142, &fl_t143, error));
    fl_value fl_t144 = fl_nothing();
    if (fl_t143) {
      fl_value fl_t145 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t145, error));
      fl_t144 = fl_t145;
    } else {
      fl_value fl_t146 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t146, error));
      fl_value fl_t147 = fl_nothing();
      FL_TRY(opis_diska_adresuetsya_soderzhimym(ctx, fl_t146, &fl_t147, error));
      bool fl_t148 = false;
      FL_TRY(fl_cond(ctx, fl_t147, &fl_t148, error));
      fl_value fl_t149 = fl_nothing();
      if (fl_t148) {
        fl_value fl_t150 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t150, error));
        fl_t149 = fl_t150;
      } else {
        fl_value fl_t151 = fl_nothing();
        if (fl_variant_is(razryad, "Кэш")) {
          fl_value fl_t152 = fl_nothing();
          FL_TRY(opis_diska_porog_kesha(ctx, &fl_t152, error));
          fl_value fl_t153 = fl_nothing();
          FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t152, &fl_t153, error));
          fl_t151 = fl_t153;
        } else if (fl_variant_is(razryad, "Сборка")) {
          fl_value fl_t154 = fl_nothing();
          FL_TRY(opis_diska_porog_kesha(ctx, &fl_t154, error));
          fl_value fl_t155 = fl_nothing();
          FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t154, &fl_t155, error));
          fl_t151 = fl_t155;
        } else if (fl_variant_is(razryad, "Журнал")) {
          fl_value fl_t156 = fl_nothing();
          FL_TRY(opis_diska_porog_zhurnala(ctx, &fl_t156, error));
          fl_value fl_t157 = fl_nothing();
          FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t156, &fl_t157, error));
          fl_t151 = fl_t157;
        } else if (fl_variant_is(razryad, "Загрузка")) {
          fl_value fl_t158 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t158, error));
          fl_value fl_t159 = fl_nothing();
          FL_TRY(opis_diska_porog_zagruzki(ctx, &fl_t159, error));
          if (fl_t158.tag != FL_NUMBER || fl_t159.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t158, fl_t159, error));
          bool fl_t160 = false;
          FL_TRY(fl_cond(ctx, fl_flag(fl_t158.as.number >= fl_t159.as.number), &fl_t160, error));
          fl_value fl_t161 = fl_nothing();
          if (fl_t160) {
            fl_value fl_t162 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t162, error));
            fl_t161 = fl_t162;
          } else {
            fl_value fl_t163 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t163, error));
            fl_t161 = fl_t163;
          }
          fl_t151 = fl_t161;
        } else if (fl_variant_is(razryad, "Крупное")) {
          fl_value fl_t164 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t164, error));
          fl_t151 = fl_t164;
        } else if (fl_variant_is(razryad, "Неизвестное")) {
          fl_value fl_t165 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t165, error));
          fl_t151 = fl_t165;
        } else {
          return fl_match_fail(ctx, razryad, error);
        }
        fl_t149 = fl_t151;
      }
      fl_t144 = fl_t149;
    }
    fl_t140 = fl_t144;
  }
  const fl_value fl_t166 = fl_t140;
  fl_value fl_t167 = fl_nothing();
  FL_TRY(opis_diska_prigovor_obosnovan(ctx, nahodka, razryad, fl_t166, &fl_t167, error));
  /* постусловие «Приговор обоснован» */
  bool fl_t168 = false;
  FL_TRY(fl_post(ctx, fl_t167, "Приговор обоснован", "Приговор находки", &fl_t168, error));
  if (!fl_t168) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Приговор обоснован» функции «Приговор находки»");
  }
  fl_value fl_t169 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t169, error));
  fl_value fl_t170 = fl_nothing();
  FL_TRY(opis_diska_adresuetsya_soderzhimym(ctx, fl_t169, &fl_t170, error));
  bool fl_t171 = false;
  FL_TRY(fl_cond(ctx, fl_t170, &fl_t171, error));
  fl_value fl_t172 = fl_nothing();
  if (fl_t171) {
    fl_value fl_t173 = fl_nothing();
    FL_TRY(opis_diska_eto_netrogat(ctx, fl_t166, &fl_t173, error));
    fl_t172 = fl_t173;
  } else {
    fl_t172 = fl_flag(true);
  }
  /* постусловие «И3: адресуемое содержимым не убирается» */
  bool fl_t174 = false;
  FL_TRY(fl_post(ctx, fl_t172, "И3: адресуемое содержимым не убирается", "Приговор находки", &fl_t174, error));
  if (!fl_t174) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И3: адресуемое содержимым не убирается» функции «Приговор находки»");
  }
  *result = fl_t166;
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
  fl_value fl_t175 = fl_nothing();
  if (fl_variant_is(prigovor, "НеТрогать")) {
    fl_t175 = fl_number(0.0);
  } else if (fl_variant_is(prigovor, "МожноУбрать")) {
    fl_value fl_t176 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t176, error));
    fl_t175 = fl_t176;
  } else if (fl_variant_is(prigovor, "Спросить")) {
    fl_value fl_t177 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t177, error));
    fl_t175 = fl_t177;
  } else {
    return fl_match_fail(ctx, prigovor, error);
  }
  const fl_value fl_t178 = fl_t175;
  fl_value fl_t179 = fl_nothing();
  FL_TRY(opis_diska_ves_obosnovan(ctx, nahodka, prigovor, fl_t178, &fl_t179, error));
  /* постусловие «Вес обоснован» */
  bool fl_t180 = false;
  FL_TRY(fl_post(ctx, fl_t179, "Вес обоснован", "Вес находки", &fl_t180, error));
  if (!fl_t180) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Вес обоснован» функции «Вес находки»");
  }
  *result = fl_t178;
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
  bool fl_t181 = false;
  FL_TRY(fl_cond(ctx, fl_flag(ves.as.number >= 0.0), &fl_t181, error));
  if (fl_t181) {
    fl_value fl_t182 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t182, error));
    if (ves.tag != FL_NUMBER || fl_t182.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, ves, fl_t182, error));
    *result = fl_flag(ves.as.number <= fl_t182.as.number);
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
 * @return значение: «Решение»
 */
fl_status opis_diska_reshit_nahodku(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error) {
  fl_value fl_t183 = fl_nothing();
  FL_TRY(opis_diska_razryad_nahodki(ctx, nahodka, &fl_t183, error));
  const fl_value razryad = fl_t183; /* пусть «разряд» */
  fl_value fl_t184 = fl_nothing();
  FL_TRY(opis_diska_prigovor_nahodki(ctx, nahodka, razryad, &fl_t184, error));
  const fl_value prigovor = fl_t184; /* пусть «приговор» */
  fl_value fl_t185 = fl_nothing();
  FL_TRY(opis_diska_ves_nahodki(ctx, nahodka, prigovor, &fl_t185, error));
  fl_value fl_t187[3];
  fl_t187[0] = razryad; /* «разряд» */
  fl_t187[1] = prigovor; /* «приговор» */
  fl_t187[2] = fl_t185; /* «вес» */
  fl_value fl_t186 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_2, fl_t187, 3, &fl_t186, error));
  const fl_value fl_t188 = fl_t186;
  fl_value fl_t189 = fl_nothing();
  FL_TRY(opis_diska_i1_derzhitsya(ctx, fl_t188, &fl_t189, error));
  /* постусловие «И1: убрать можно только мусор» */
  bool fl_t190 = false;
  FL_TRY(fl_post(ctx, fl_t189, "И1: убрать можно только мусор", "Решить находку", &fl_t190, error));
  if (!fl_t190) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И1: убрать можно только мусор» функции «Решить находку»");
  }
  *result = fl_t188;
  return FL_OK;
}

/*
 * Функция flang «Решить всё».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @return значение: список: «Решение»
 */
fl_status opis_diska_reshit_vsyo(fl_ctx *ctx, fl_value zapisi, fl_value *result, fl_error *error) {
  fl_value fl_t191 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "отобразить", &fl_t191, error));
  fl_value *fl_t192 = NULL;
  size_t fl_t193 = 0;
  FL_TRY(fl_list_alloc(ctx, fl_t191.as.list.count, &fl_t192, error));
  for (size_t fl_t194 = 0; fl_t194 < fl_t191.as.list.count; fl_t194 += 1) {
    const fl_value nahodka = fl_t191.as.list.items[fl_t194]; /* «находка» */
    fl_value fl_t195 = fl_nothing();
    FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, &fl_t195, error));
    fl_t192[fl_t193] = fl_t195;
    fl_t193 += 1;
  }
  const fl_value fl_t196 = fl_list(fl_t192, fl_t193);
  fl_value fl_t197 = fl_nothing();
  FL_TRY(opis_diska_i1_derzhitsya_vsyudu(ctx, fl_t196, &fl_t197, error));
  /* постусловие «И1 всюду: убрать можно только мусор» */
  bool fl_t198 = false;
  FL_TRY(fl_post(ctx, fl_t197, "И1 всюду: убрать можно только мусор", "Решить всё", &fl_t198, error));
  if (!fl_t198) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И1 всюду: убрать можно только мусор» функции «Решить всё»");
  }
  *result = fl_t196;
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
  fl_value fl_t199 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "приговор", &fl_t199, error));
  if (fl_variant_is(fl_t199, "МожноУбрать")) {
    fl_value fl_t200 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "разряд", &fl_t200, error));
    if (fl_variant_is(fl_t200, "Кэш")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t200, "Журнал")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t200, "Сборка")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t200, "Загрузка")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(fl_t200, "Крупное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(fl_t200, "Неизвестное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else {
      return fl_match_fail(ctx, fl_t200, error);
    }
  } else if (fl_variant_is(fl_t199, "Спросить")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t199, "НеТрогать")) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t199, error);
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
  fl_value fl_t201 = fl_nothing();
  FL_TRY(fl_require_list(ctx, resheniya, "свёртка", &fl_t201, error));
  fl_value akk = fl_flag(true); /* «акк» */
  const fl_mark fl_t203 = fl_region_open(ctx);
  for (size_t fl_t202 = 0; fl_t202 < fl_t201.as.list.count; fl_t202 += 1) {
    const fl_value reshenie = fl_t201.as.list.items[fl_t202]; /* «решение» */
    bool fl_t204 = false;
    FL_TRY(fl_cond(ctx, akk, &fl_t204, error));
    fl_value fl_t205 = fl_nothing();
    if (fl_t204) {
      fl_value fl_t206 = fl_nothing();
      FL_TRY(opis_diska_i1_derzhitsya(ctx, reshenie, &fl_t206, error));
      fl_t205 = fl_t206;
    } else {
      fl_t205 = fl_flag(false);
    }
    akk = fl_t205;
    FL_TRY(fl_region_recycle(ctx, fl_t203, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t203, FL_OK, &akk, error));
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
  fl_value fl_t208[7];
  fl_t208[0] = fl_number(0.0); /* «кэш» */
  fl_t208[1] = fl_number(0.0); /* «журнал» */
  fl_t208[2] = fl_number(0.0); /* «сборка» */
  fl_t208[3] = fl_number(0.0); /* «загрузка» */
  fl_t208[4] = fl_number(0.0); /* «крупное» */
  fl_t208[5] = fl_number(0.0); /* «неизвестное» */
  fl_t208[6] = fl_number(0.0); /* «освободить» */
  fl_value fl_t207 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t208, 7, &fl_t207, error));
  *result = fl_t207;
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
  fl_value fl_t209 = fl_nothing();
  fl_value fl_t210 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "приговор", &fl_t210, error));
  if (fl_variant_is(fl_t210, "МожноУбрать")) {
    fl_value fl_t211 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t211, error));
    fl_t209 = fl_t211;
  } else if (fl_variant_is(fl_t210, "Спросить")) {
    fl_t209 = fl_number(0.0);
  } else if (fl_variant_is(fl_t210, "НеТрогать")) {
    fl_t209 = fl_number(0.0);
  } else {
    return fl_match_fail(ctx, fl_t210, error);
  }
  const fl_value ubrat = fl_t209; /* пусть «убрать» */
  fl_value fl_t212 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "разряд", &fl_t212, error));
  if (fl_variant_is(fl_t212, "Кэш")) {
    fl_value fl_t213 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t213, error));
    fl_value fl_t214 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t214, error));
    if (fl_t213.tag != FL_NUMBER || fl_t214.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t213, fl_t214, error));
    fl_value fl_t215 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t215, error));
    fl_value fl_t216 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t216, error));
    fl_value fl_t217 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t217, error));
    fl_value fl_t218 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t218, error));
    fl_value fl_t219 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t219, error));
    fl_value fl_t220 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t220, error));
    if (fl_t220.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t220, ubrat, error));
    fl_value fl_t222[7];
    fl_t222[0] = fl_number(fl_t213.as.number + fl_t214.as.number); /* «кэш» */
    fl_t222[1] = fl_t215; /* «журнал» */
    fl_t222[2] = fl_t216; /* «сборка» */
    fl_t222[3] = fl_t217; /* «загрузка» */
    fl_t222[4] = fl_t218; /* «крупное» */
    fl_t222[5] = fl_t219; /* «неизвестное» */
    fl_t222[6] = fl_number(fl_t220.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t221 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t222, 7, &fl_t221, error));
    *result = fl_t221;
    return FL_OK;
  } else if (fl_variant_is(fl_t212, "Журнал")) {
    fl_value fl_t223 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t223, error));
    fl_value fl_t224 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t224, error));
    fl_value fl_t225 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t225, error));
    if (fl_t224.tag != FL_NUMBER || fl_t225.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t224, fl_t225, error));
    fl_value fl_t226 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t226, error));
    fl_value fl_t227 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t227, error));
    fl_value fl_t228 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t228, error));
    fl_value fl_t229 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t229, error));
    fl_value fl_t230 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t230, error));
    if (fl_t230.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t230, ubrat, error));
    fl_value fl_t232[7];
    fl_t232[0] = fl_t223; /* «кэш» */
    fl_t232[1] = fl_number(fl_t224.as.number + fl_t225.as.number); /* «журнал» */
    fl_t232[2] = fl_t226; /* «сборка» */
    fl_t232[3] = fl_t227; /* «загрузка» */
    fl_t232[4] = fl_t228; /* «крупное» */
    fl_t232[5] = fl_t229; /* «неизвестное» */
    fl_t232[6] = fl_number(fl_t230.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t231 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t232, 7, &fl_t231, error));
    *result = fl_t231;
    return FL_OK;
  } else if (fl_variant_is(fl_t212, "Сборка")) {
    fl_value fl_t233 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t233, error));
    fl_value fl_t234 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t234, error));
    fl_value fl_t235 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t235, error));
    fl_value fl_t236 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t236, error));
    if (fl_t235.tag != FL_NUMBER || fl_t236.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t235, fl_t236, error));
    fl_value fl_t237 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t237, error));
    fl_value fl_t238 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t238, error));
    fl_value fl_t239 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t239, error));
    fl_value fl_t240 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t240, error));
    if (fl_t240.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t240, ubrat, error));
    fl_value fl_t242[7];
    fl_t242[0] = fl_t233; /* «кэш» */
    fl_t242[1] = fl_t234; /* «журнал» */
    fl_t242[2] = fl_number(fl_t235.as.number + fl_t236.as.number); /* «сборка» */
    fl_t242[3] = fl_t237; /* «загрузка» */
    fl_t242[4] = fl_t238; /* «крупное» */
    fl_t242[5] = fl_t239; /* «неизвестное» */
    fl_t242[6] = fl_number(fl_t240.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t241 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t242, 7, &fl_t241, error));
    *result = fl_t241;
    return FL_OK;
  } else if (fl_variant_is(fl_t212, "Загрузка")) {
    fl_value fl_t243 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t243, error));
    fl_value fl_t244 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t244, error));
    fl_value fl_t245 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t245, error));
    fl_value fl_t246 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t246, error));
    fl_value fl_t247 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t247, error));
    if (fl_t246.tag != FL_NUMBER || fl_t247.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t246, fl_t247, error));
    fl_value fl_t248 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t248, error));
    fl_value fl_t249 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t249, error));
    fl_value fl_t250 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t250, error));
    if (fl_t250.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t250, ubrat, error));
    fl_value fl_t252[7];
    fl_t252[0] = fl_t243; /* «кэш» */
    fl_t252[1] = fl_t244; /* «журнал» */
    fl_t252[2] = fl_t245; /* «сборка» */
    fl_t252[3] = fl_number(fl_t246.as.number + fl_t247.as.number); /* «загрузка» */
    fl_t252[4] = fl_t248; /* «крупное» */
    fl_t252[5] = fl_t249; /* «неизвестное» */
    fl_t252[6] = fl_number(fl_t250.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t251 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t252, 7, &fl_t251, error));
    *result = fl_t251;
    return FL_OK;
  } else if (fl_variant_is(fl_t212, "Крупное")) {
    fl_value fl_t253 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t253, error));
    fl_value fl_t254 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t254, error));
    fl_value fl_t255 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t255, error));
    fl_value fl_t256 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t256, error));
    fl_value fl_t257 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t257, error));
    fl_value fl_t258 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t258, error));
    if (fl_t257.tag != FL_NUMBER || fl_t258.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t257, fl_t258, error));
    fl_value fl_t259 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t259, error));
    fl_value fl_t260 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t260, error));
    if (fl_t260.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t260, ubrat, error));
    fl_value fl_t262[7];
    fl_t262[0] = fl_t253; /* «кэш» */
    fl_t262[1] = fl_t254; /* «журнал» */
    fl_t262[2] = fl_t255; /* «сборка» */
    fl_t262[3] = fl_t256; /* «загрузка» */
    fl_t262[4] = fl_number(fl_t257.as.number + fl_t258.as.number); /* «крупное» */
    fl_t262[5] = fl_t259; /* «неизвестное» */
    fl_t262[6] = fl_number(fl_t260.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t261 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t262, 7, &fl_t261, error));
    *result = fl_t261;
    return FL_OK;
  } else if (fl_variant_is(fl_t212, "Неизвестное")) {
    fl_value fl_t263 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t263, error));
    fl_value fl_t264 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t264, error));
    fl_value fl_t265 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t265, error));
    fl_value fl_t266 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t266, error));
    fl_value fl_t267 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t267, error));
    fl_value fl_t268 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t268, error));
    fl_value fl_t269 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t269, error));
    if (fl_t268.tag != FL_NUMBER || fl_t269.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t268, fl_t269, error));
    fl_value fl_t270 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t270, error));
    if (fl_t270.tag != FL_NUMBER || ubrat.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", fl_t270, ubrat, error));
    fl_value fl_t272[7];
    fl_t272[0] = fl_t263; /* «кэш» */
    fl_t272[1] = fl_t264; /* «журнал» */
    fl_t272[2] = fl_t265; /* «сборка» */
    fl_t272[3] = fl_t266; /* «загрузка» */
    fl_t272[4] = fl_t267; /* «крупное» */
    fl_t272[5] = fl_number(fl_t268.as.number + fl_t269.as.number); /* «неизвестное» */
    fl_t272[6] = fl_number(fl_t270.as.number + ubrat.as.number); /* «освободить» */
    fl_value fl_t271 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t272, 7, &fl_t271, error));
    *result = fl_t271;
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t212, error);
  }
}

/*
 * Функция flang «Свести».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @return значение: «Свод»
 */
fl_status opis_diska_svesti(fl_ctx *ctx, fl_value zapisi, fl_value *result, fl_error *error) {
  fl_value fl_t273 = fl_nothing();
  FL_TRY(opis_diska_reshit_vsyo(ctx, zapisi, &fl_t273, error));
  fl_value fl_t274 = fl_nothing();
  FL_TRY(fl_require_list(ctx, fl_t273, "свёртка", &fl_t274, error));
  fl_value fl_t275 = fl_nothing();
  FL_TRY(opis_diska_pustoy_svod(ctx, &fl_t275, error));
  fl_value svod = fl_t275; /* «свод» */
  const fl_mark fl_t277 = fl_region_open(ctx);
  for (size_t fl_t276 = 0; fl_t276 < fl_t274.as.list.count; fl_t276 += 1) {
    const fl_value reshenie = fl_t274.as.list.items[fl_t276]; /* «решение» */
    fl_value fl_t278 = fl_nothing();
    FL_TRY(opis_diska_pribavit_reshenie(ctx, svod, reshenie, &fl_t278, error));
    svod = fl_t278;
    FL_TRY(fl_region_recycle(ctx, fl_t277, &svod, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t277, FL_OK, &svod, error));
  const fl_value fl_t279 = svod;
  fl_value fl_t280 = fl_nothing();
  FL_TRY(opis_diska_i2_derzhitsya(ctx, zapisi, fl_t279, &fl_t280, error));
  /* постусловие «И2: освобождаемое не больше убираемого» */
  bool fl_t281 = false;
  FL_TRY(fl_post(ctx, fl_t280, "И2: освобождаемое не больше убираемого", "Свести", &fl_t281, error));
  if (!fl_t281) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «И2: освобождаемое не больше убираемого» функции «Свести»");
  }
  *result = fl_t279;
  return FL_OK;
}

/*
 * Функция flang «Сумма размеров убираемых».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @return значение: число
 */
fl_status opis_diska_summa_razmerov_ubiraemyh(fl_ctx *ctx, fl_value zapisi, fl_value *result, fl_error *error) {
  fl_value fl_t282 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "свёртка", &fl_t282, error));
  fl_value akk = fl_number(0.0); /* «акк» */
  const fl_mark fl_t284 = fl_region_open(ctx);
  for (size_t fl_t283 = 0; fl_t283 < fl_t282.as.list.count; fl_t283 += 1) {
    const fl_value nahodka = fl_t282.as.list.items[fl_t283]; /* «находка» */
    fl_value fl_t285 = fl_nothing();
    fl_value fl_t286 = fl_nothing();
    FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, &fl_t286, error));
    fl_value fl_t287 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t286, "приговор", &fl_t287, error));
    if (fl_variant_is(fl_t287, "МожноУбрать")) {
      fl_value fl_t288 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t288, error));
      if (akk.tag != FL_NUMBER || fl_t288.tag != FL_NUMBER) FL_TRY(fl_not_numbers(ctx, "add", akk, fl_t288, error));
      fl_t285 = fl_number(akk.as.number + fl_t288.as.number);
    } else if (fl_variant_is(fl_t287, "Спросить")) {
      fl_t285 = akk;
    } else if (fl_variant_is(fl_t287, "НеТрогать")) {
      fl_t285 = akk;
    } else {
      return fl_match_fail(ctx, fl_t287, error);
    }
    akk = fl_t285;
    FL_TRY(fl_region_recycle(ctx, fl_t284, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t284, FL_OK, &akk, error));
  *result = akk;
  return FL_OK;
}

/*
 * Функция flang «И2 держится».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param svod — «свод»: «Свод»
 * @return значение
 */
fl_status opis_diska_i2_derzhitsya(fl_ctx *ctx, fl_value zapisi, fl_value svod, fl_value *result, fl_error *error) {
  fl_value fl_t289 = fl_nothing();
  FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t289, error));
  fl_value fl_t290 = fl_nothing();
  FL_TRY(opis_diska_summa_razmerov_ubiraemyh(ctx, zapisi, &fl_t290, error));
  if (fl_t289.tag != FL_NUMBER || fl_t290.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t289, fl_t290, error));
  *result = fl_flag(fl_t289.as.number <= fl_t290.as.number);
  return FL_OK;
}

/*
 * Функция flang «Строку отчёта».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @return значение: «Строка отчёта»
 */
fl_status opis_diska_stroku_otchyota(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error) {
  fl_value fl_t291 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t291, error));
  fl_value fl_t292 = fl_nothing();
  FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, &fl_t292, error));
  fl_value fl_t294[2];
  fl_t294[0] = fl_t291; /* «путь» */
  fl_t294[1] = fl_t292; /* «решение» */
  fl_value fl_t293 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_3, fl_t294, 2, &fl_t293, error));
  *result = fl_t293;
  return FL_OK;
}

/* Тело «Вставить по весу»; глубину считает обёртка ниже. */
static fl_status opis_diska_vstavit_po_vesu_body(fl_ctx *ctx, fl_value stroka, fl_value stroki, fl_value *result, fl_error *error) {
  if (fl_chain_empty(stroki)) {
    fl_value *fl_t295 = NULL;
    FL_TRY(fl_list_alloc(ctx, 1, &fl_t295, error));
    fl_t295[0] = stroka;
    *result = fl_list(fl_t295, 1);
    return FL_OK;
  } else if (fl_chain_cons(stroki)) {
    const fl_value golova = fl_chain_head(stroki); /* голова «голова» */
    const fl_value hvost = fl_chain_tail(stroki); /* хвост «хвост» */
    fl_value fl_t296 = fl_nothing();
    FL_TRY(fl_field_get(ctx, stroka, "решение", &fl_t296, error));
    fl_value fl_t297 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t296, "вес", &fl_t297, error));
    fl_value fl_t298 = fl_nothing();
    FL_TRY(fl_field_get(ctx, golova, "решение", &fl_t298, error));
    fl_value fl_t299 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t298, "вес", &fl_t299, error));
    if (fl_t297.tag != FL_NUMBER || fl_t299.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t297, fl_t299, error));
    bool fl_t300 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_t297.as.number >= fl_t299.as.number), &fl_t300, error));
    if (fl_t300) {
      return opis_diska_pripisat_stroku_otchyota(ctx, stroka, stroki, result, error);
    } else {
      fl_value fl_t301 = fl_nothing();
      FL_TRY(opis_diska_vstavit_po_vesu(ctx, stroka, hvost, &fl_t301, error));
      return opis_diska_pripisat_stroku_otchyota(ctx, golova, fl_t301, result, error);
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
  fl_value fl_t302 = fl_nothing();
  FL_TRY(fl_require_list(ctx, stroki, "свёртка", &fl_t302, error));
  fl_value *fl_t303 = NULL;
  FL_TRY(fl_list_alloc(ctx, 1, &fl_t303, error));
  fl_t303[0] = pervaya;
  fl_value akk = fl_list(fl_t303, 1); /* «акк» */
  const fl_mark fl_t305 = fl_region_open(ctx);
  for (size_t fl_t304 = 0; fl_t304 < fl_t302.as.list.count; fl_t304 += 1) {
    const fl_value el = fl_t302.as.list.items[fl_t304]; /* «эл» */
    fl_value fl_t306 = fl_nothing(); /* «добавить» */
    FL_TRY(fl_b_dobavit(ctx, el, akk, &fl_t306, error));
    akk = fl_t306;
    FL_TRY(fl_region_recycle(ctx, fl_t305, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t305, FL_OK, &akk, error));
  *result = akk;
  return FL_OK;
}

/*
 * Функция flang «Отчёт».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @return значение: список: «Строка отчёта»
 */
fl_status opis_diska_otchyot(fl_ctx *ctx, fl_value zapisi, fl_value *result, fl_error *error) {
  fl_value fl_t307 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "свёртка", &fl_t307, error));
  fl_value akk = fl_list(NULL, 0); /* «акк» */
  const fl_mark fl_t309 = fl_region_open(ctx);
  for (size_t fl_t308 = 0; fl_t308 < fl_t307.as.list.count; fl_t308 += 1) {
    const fl_value nahodka = fl_t307.as.list.items[fl_t308]; /* «находка» */
    fl_value fl_t310 = fl_nothing();
    FL_TRY(opis_diska_stroku_otchyota(ctx, nahodka, &fl_t310, error));
    fl_value fl_t311 = fl_nothing();
    FL_TRY(opis_diska_vstavit_po_vesu(ctx, fl_t310, akk, &fl_t311, error));
    akk = fl_t311;
    FL_TRY(fl_region_recycle(ctx, fl_t309, &akk, error));
  }
  FL_TRY(fl_region_close(ctx, fl_t309, FL_OK, &akk, error));
  const fl_value fl_t312 = akk;
  fl_value fl_t313 = fl_nothing();
  FL_TRY(opis_diska_otchyot_toy_zhe_dliny(ctx, zapisi, fl_t312, &fl_t313, error));
  /* постусловие «Отчёт той же длины» */
  bool fl_t314 = false;
  FL_TRY(fl_post(ctx, fl_t313, "Отчёт той же длины", "Отчёт", &fl_t314, error));
  if (!fl_t314) {
    return fl_fail(ctx, error, "FLANG_PROPERTY", "%s", "нарушено свойство «Отчёт той же длины» функции «Отчёт»");
  }
  *result = fl_t312;
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
  fl_value fl_t315 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, stroki, &fl_t315, error));
  fl_value fl_t316 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, zapisi, &fl_t316, error));
  *result = fl_flag(fl_equal(fl_t315, fl_t316));
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
  fl_value fl_t317 = fl_nothing();
  FL_TRY(opis_diska_eto_mozhnoubrat(ctx, prigovor, &fl_t317, error));
  bool fl_t318 = false;
  FL_TRY(fl_cond(ctx, fl_t317, &fl_t318, error));
  fl_value fl_t319 = fl_nothing();
  if (fl_t318) {
    fl_t319 = fl_flag(false);
  } else {
    fl_t319 = fl_flag(true);
  }
  bool fl_t320 = false;
  FL_TRY(fl_cond(ctx, fl_t319, &fl_t320, error));
  if (fl_t320) {
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
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error) {
  fl_value fl_t321 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t321, error));
  fl_value fl_t322 = fl_nothing();
  FL_TRY(opis_diska_primeta_kesha(ctx, fl_t321, &fl_t322, error));
  const fl_value kesh = fl_t322; /* пусть «кэш» */
  fl_value fl_t323 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t323, error));
  fl_value fl_t324 = fl_nothing();
  FL_TRY(opis_diska_primeta_zhurnala(ctx, fl_t323, &fl_t324, error));
  const fl_value zhurnal = fl_t324; /* пусть «журнал» */
  fl_value fl_t325 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t325, error));
  fl_value fl_t326 = fl_nothing();
  FL_TRY(opis_diska_primeta_sborki(ctx, fl_t325, &fl_t326, error));
  const fl_value sborka = fl_t326; /* пусть «сборка» */
  fl_value fl_t327 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t327, error));
  fl_value fl_t328 = fl_nothing();
  FL_TRY(opis_diska_primeta_zagruzki(ctx, fl_t327, &fl_t328, error));
  const fl_value zagruzka = fl_t328; /* пусть «загрузка» */
  fl_value fl_t329 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t329, error));
  fl_value fl_t330 = fl_nothing();
  FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t330, error));
  if (fl_t329.tag != FL_NUMBER || fl_t330.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t329, fl_t330, error));
  const fl_value krupnoe = fl_flag(fl_t329.as.number >= fl_t330.as.number); /* пусть «крупное» */
  if (fl_variant_is(razryad, "Кэш")) {
    *result = kesh;
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    bool fl_t331 = false;
    FL_TRY(fl_cond(ctx, zhurnal, &fl_t331, error));
    if (fl_t331) {
      bool fl_t332 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t332, error));
      if (fl_t332) {
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
    bool fl_t333 = false;
    FL_TRY(fl_cond(ctx, sborka, &fl_t333, error));
    fl_value fl_t334 = fl_nothing();
    if (fl_t333) {
      bool fl_t335 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t335, error));
      fl_value fl_t336 = fl_nothing();
      if (fl_t335) {
        fl_t336 = fl_flag(false);
      } else {
        fl_t336 = fl_flag(true);
      }
      fl_t334 = fl_t336;
    } else {
      fl_t334 = fl_flag(false);
    }
    bool fl_t337 = false;
    FL_TRY(fl_cond(ctx, fl_t334, &fl_t337, error));
    if (fl_t337) {
      bool fl_t338 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t338, error));
      if (fl_t338) {
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
    bool fl_t339 = false;
    FL_TRY(fl_cond(ctx, zagruzka, &fl_t339, error));
    fl_value fl_t340 = fl_nothing();
    if (fl_t339) {
      bool fl_t341 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t341, error));
      fl_value fl_t342 = fl_nothing();
      if (fl_t341) {
        fl_t342 = fl_flag(false);
      } else {
        fl_t342 = fl_flag(true);
      }
      fl_t340 = fl_t342;
    } else {
      fl_t340 = fl_flag(false);
    }
    bool fl_t343 = false;
    FL_TRY(fl_cond(ctx, fl_t340, &fl_t343, error));
    fl_value fl_t344 = fl_nothing();
    if (fl_t343) {
      bool fl_t345 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t345, error));
      fl_value fl_t346 = fl_nothing();
      if (fl_t345) {
        fl_t346 = fl_flag(false);
      } else {
        fl_t346 = fl_flag(true);
      }
      fl_t344 = fl_t346;
    } else {
      fl_t344 = fl_flag(false);
    }
    bool fl_t347 = false;
    FL_TRY(fl_cond(ctx, fl_t344, &fl_t347, error));
    if (fl_t347) {
      bool fl_t348 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t348, error));
      if (fl_t348) {
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
    bool fl_t349 = false;
    FL_TRY(fl_cond(ctx, krupnoe, &fl_t349, error));
    fl_value fl_t350 = fl_nothing();
    if (fl_t349) {
      bool fl_t351 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t351, error));
      fl_value fl_t352 = fl_nothing();
      if (fl_t351) {
        fl_t352 = fl_flag(false);
      } else {
        fl_t352 = fl_flag(true);
      }
      fl_t350 = fl_t352;
    } else {
      fl_t350 = fl_flag(false);
    }
    bool fl_t353 = false;
    FL_TRY(fl_cond(ctx, fl_t350, &fl_t353, error));
    fl_value fl_t354 = fl_nothing();
    if (fl_t353) {
      bool fl_t355 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t355, error));
      fl_value fl_t356 = fl_nothing();
      if (fl_t355) {
        fl_t356 = fl_flag(false);
      } else {
        fl_t356 = fl_flag(true);
      }
      fl_t354 = fl_t356;
    } else {
      fl_t354 = fl_flag(false);
    }
    bool fl_t357 = false;
    FL_TRY(fl_cond(ctx, fl_t354, &fl_t357, error));
    fl_value fl_t358 = fl_nothing();
    if (fl_t357) {
      bool fl_t359 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t359, error));
      fl_value fl_t360 = fl_nothing();
      if (fl_t359) {
        fl_t360 = fl_flag(false);
      } else {
        fl_t360 = fl_flag(true);
      }
      fl_t358 = fl_t360;
    } else {
      fl_t358 = fl_flag(false);
    }
    bool fl_t361 = false;
    FL_TRY(fl_cond(ctx, fl_t358, &fl_t361, error));
    if (fl_t361) {
      bool fl_t362 = false;
      FL_TRY(fl_cond(ctx, zagruzka, &fl_t362, error));
      if (fl_t362) {
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
    bool fl_t363 = false;
    FL_TRY(fl_cond(ctx, kesh, &fl_t363, error));
    fl_value fl_t364 = fl_nothing();
    if (fl_t363) {
      fl_t364 = fl_flag(false);
    } else {
      fl_t364 = fl_flag(true);
    }
    bool fl_t365 = false;
    FL_TRY(fl_cond(ctx, fl_t364, &fl_t365, error));
    fl_value fl_t366 = fl_nothing();
    if (fl_t365) {
      bool fl_t367 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t367, error));
      fl_value fl_t368 = fl_nothing();
      if (fl_t367) {
        fl_t368 = fl_flag(false);
      } else {
        fl_t368 = fl_flag(true);
      }
      fl_t366 = fl_t368;
    } else {
      fl_t366 = fl_flag(false);
    }
    bool fl_t369 = false;
    FL_TRY(fl_cond(ctx, fl_t366, &fl_t369, error));
    fl_value fl_t370 = fl_nothing();
    if (fl_t369) {
      bool fl_t371 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t371, error));
      fl_value fl_t372 = fl_nothing();
      if (fl_t371) {
        fl_t372 = fl_flag(false);
      } else {
        fl_t372 = fl_flag(true);
      }
      fl_t370 = fl_t372;
    } else {
      fl_t370 = fl_flag(false);
    }
    bool fl_t373 = false;
    FL_TRY(fl_cond(ctx, fl_t370, &fl_t373, error));
    fl_value fl_t374 = fl_nothing();
    if (fl_t373) {
      bool fl_t375 = false;
      FL_TRY(fl_cond(ctx, zagruzka, &fl_t375, error));
      fl_value fl_t376 = fl_nothing();
      if (fl_t375) {
        fl_t376 = fl_flag(false);
      } else {
        fl_t376 = fl_flag(true);
      }
      fl_t374 = fl_t376;
    } else {
      fl_t374 = fl_flag(false);
    }
    bool fl_t377 = false;
    FL_TRY(fl_cond(ctx, fl_t374, &fl_t377, error));
    if (fl_t377) {
      bool fl_t378 = false;
      FL_TRY(fl_cond(ctx, krupnoe, &fl_t378, error));
      if (fl_t378) {
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
  fl_value fl_t379 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "доступен", &fl_t379, error));
  bool fl_t380 = false;
  FL_TRY(fl_cond(ctx, fl_t379, &fl_t380, error));
  fl_value fl_t381 = fl_nothing();
  if (fl_t380) {
    fl_t381 = fl_flag(false);
  } else {
    fl_t381 = fl_flag(true);
  }
  bool fl_t382 = false;
  FL_TRY(fl_cond(ctx, fl_t381, &fl_t382, error));
  if (fl_t382) {
    return opis_diska_eto_netrogat(ctx, prigovor, result, error);
  } else {
    fl_value fl_t383 = fl_nothing();
    FL_TRY(opis_diska_ssylka(ctx, nahodka, &fl_t383, error));
    bool fl_t384 = false;
    FL_TRY(fl_cond(ctx, fl_t383, &fl_t384, error));
    if (fl_t384) {
      return opis_diska_eto_netrogat(ctx, prigovor, result, error);
    } else {
      fl_value fl_t385 = fl_nothing();
      FL_TRY(opis_diska_eto_mozhnoubrat(ctx, prigovor, &fl_t385, error));
      bool fl_t386 = false;
      FL_TRY(fl_cond(ctx, fl_t385, &fl_t386, error));
      if (fl_t386) {
        fl_value fl_t387 = fl_nothing();
        FL_TRY(opis_diska_i1_na_pare(ctx, razryad, prigovor, &fl_t387, error));
        bool fl_t388 = false;
        FL_TRY(fl_cond(ctx, fl_t387, &fl_t388, error));
        fl_value fl_t389 = fl_nothing();
        if (fl_t388) {
          fl_value fl_t390 = fl_nothing();
          FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t390, error));
          bool fl_t391 = false;
          FL_TRY(fl_cond(ctx, fl_t390, &fl_t391, error));
          fl_value fl_t392 = fl_nothing();
          if (fl_t391) {
            fl_t392 = fl_flag(false);
          } else {
            fl_t392 = fl_flag(true);
          }
          fl_t389 = fl_t392;
        } else {
          fl_t389 = fl_flag(false);
        }
        bool fl_t393 = false;
        FL_TRY(fl_cond(ctx, fl_t389, &fl_t393, error));
        if (fl_t393) {
          fl_value fl_t394 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t394, error));
          fl_value fl_t395 = fl_nothing();
          FL_TRY(opis_diska_porog_razryada(ctx, razryad, &fl_t395, error));
          if (fl_t394.tag != FL_NUMBER || fl_t395.tag != FL_NUMBER) FL_TRY(fl_not_order(ctx, fl_t394, fl_t395, error));
          *result = fl_flag(fl_t394.as.number >= fl_t395.as.number);
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
  fl_value fl_t396 = fl_nothing();
  FL_TRY(opis_diska_eto_netrogat(ctx, prigovor, &fl_t396, error));
  bool fl_t397 = false;
  FL_TRY(fl_cond(ctx, fl_t396, &fl_t397, error));
  if (fl_t397) {
    *result = fl_flag(fl_equal(ves, fl_number(0.0)));
    return FL_OK;
  } else {
    fl_value fl_t398 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t398, error));
    bool fl_t399 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_equal(ves, fl_t398)), &fl_t399, error));
    if (fl_t399) {
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
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Есть примета", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_est_primeta(ctx, args[0], result, error);
  }
  if (strcmp(name, "Разряд решён размером") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд решён размером", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_razryad_reshyon_razmerom(ctx, args[0], result, error);
  }
  if (strcmp(name, "Разряд находки") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд находки", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_razryad_nahodki(ctx, args[0], result, error);
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
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Решить находку", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_reshit_nahodku(ctx, args[0], result, error);
  }
  if (strcmp(name, "Решить всё") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Решить всё", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_reshit_vsyo(ctx, args[0], result, error);
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
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Свести", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_svesti(ctx, args[0], result, error);
  }
  if (strcmp(name, "Сумма размеров убираемых") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Сумма размеров убираемых", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_summa_razmerov_ubiraemyh(ctx, args[0], result, error);
  }
  if (strcmp(name, "И2 держится") == 0) {
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "И2 держится", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_i2_derzhitsya(ctx, args[0], args[1], result, error);
  }
  if (strcmp(name, "Строку отчёта") == 0) {
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Строку отчёта", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_stroku_otchyota(ctx, args[0], result, error);
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
    if (count != 1) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Отчёт", (unsigned long)1, (unsigned long)count);
    }
    return opis_diska_otchyot(ctx, args[0], result, error);
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
    if (count != 2) {
      return fl_fail(ctx, error, FL_CODE_TYPE, "функция «%s» принимает %lu аргум., получено %lu",
                     "Разряд обоснован", (unsigned long)2, (unsigned long)count);
    }
    return opis_diska_razryad_obosnovan(ctx, args[0], args[1], result, error);
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
  { "путь", 0 },
  { "размер", 3 },
  { "возраст_дней", 3 },
  { "вид", 4 },
  { "доступен", 5 },
  { "разряд", 1 },
  { "приговор", 6 },
  { "вес", 3 },
  { "кэш", 3 },
  { "журнал", 3 },
  { "сборка", 3 },
  { "загрузка", 3 },
  { "крупное", 3 },
  { "неизвестное", 3 },
  { "освободить", 3 },
  { "путь", 0 },
  { "решение", 8 },
};

static const fl_type_variant opis_diska_entry_variants[] = {
  { "Кэш", 0, 0 },
  { "Журнал", 0, 0 },
  { "Сборка", 0, 0 },
  { "Загрузка", 0, 0 },
  { "Крупное", 0, 0 },
  { "Неизвестное", 0, 0 },
  { "Файл", 0, 0 },
  { "Каталог", 0, 0 },
  { "Ссылка", 0, 0 },
  { "МожноУбрать", 5, 0 },
  { "Спросить", 5, 0 },
  { "НеТрогать", 5, 0 },
};

static const fl_type opis_diska_entry_types[] = {
  { FL_TYPE_STRING, "строка", "", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 0 },
  { FL_TYPE_SUM, "«Разряд»", "Разряд", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 6 },
  { FL_TYPE_RECORD, "«Находка»", "Находка", false, false, false, 0.0, 0.0, 0, 0, 5, 0, 0 },
  { FL_TYPE_NUMBER, "число", "", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 0 },
  { FL_TYPE_SUM, "«Вид»", "Вид", false, false, false, 0.0, 0.0, 0, 0, 0, 6, 3 },
  { FL_TYPE_FLAG, "признак", "", false, false, false, 0.0, 0.0, 0, 0, 0, 0, 0 },
  { FL_TYPE_SUM, "«Приговор»", "Приговор", false, false, false, 0.0, 0.0, 0, 0, 0, 9, 3 },
  { FL_TYPE_LIST, "список «Находка»", "", false, false, false, 0.0, 0.0, 2, 0, 0, 0, 0 },
  { FL_TYPE_RECORD, "«Решение»", "Решение", false, false, false, 0.0, 0.0, 0, 5, 3, 0, 0 },
  { FL_TYPE_LIST, "список «Решение»", "", false, false, false, 0.0, 0.0, 8, 0, 0, 0, 0 },
  { FL_TYPE_RECORD, "«Свод»", "Свод", false, false, false, 0.0, 0.0, 0, 8, 7, 0, 0 },
  { FL_TYPE_RECORD, "«Строка отчёта»", "Строка отчёта", false, false, false, 0.0, 0.0, 0, 15, 2, 0, 0 },
  { FL_TYPE_LIST, "список «Строка отчёта»", "", false, false, false, 0.0, 0.0, 11, 0, 0, 0, 0 },
};

static const fl_entry_param opis_diska_entry_params[] = {
  { "Составляющие пути", "путь", 0 },
  { "Имя в пути", "путь", 0 },
  { "Есть составляющая", "путь", 0 },
  { "Есть составляющая", "имя", 0 },
  { "Оканчивается на", "текст", 0 },
  { "Оканчивается на", "хвост", 0 },
  { "Шестнадцатеричный знак", "знак", 0 },
  { "Похоже на отпечаток", "часть", 0 },
  { "Адресуется содержимым", "путь", 0 },
  { "Под системным временным", "путь", 0 },
  { "Примета кэша", "путь", 0 },
  { "Примета журнала", "путь", 0 },
  { "Примета сборки", "путь", 0 },
  { "Примета загрузки", "путь", 0 },
  { "Есть примета", "путь", 0 },
  { "Разряд решён размером", "разряд", 1 },
  { "Разряд находки", "находка", 2 },
  { "Крупное не мельче порога", "находка", 2 },
  { "Крупное не мельче порога", "разряд", 1 },
  { "Каталог", "находка", 2 },
  { "Ссылка", "находка", 2 },
  { "Приговор мусора", "находка", 2 },
  { "Приговор мусора", "порог", 3 },
  { "Приговор находки", "находка", 2 },
  { "Приговор находки", "разряд", 1 },
  { "Вес находки", "находка", 2 },
  { "Вес находки", "приговор", 6 },
  { "Вес в границах", "находка", 2 },
  { "Вес в границах", "вес", 3 },
  { "Решить находку", "находка", 2 },
  { "Решить всё", "записи", 7 },
  { "И1 держится", "решение", 8 },
  { "И1 держится всюду", "решения", 9 },
  { "Прибавить решение", "свод", 10 },
  { "Прибавить решение", "решение", 8 },
  { "Свести", "записи", 7 },
  { "Сумма размеров убираемых", "записи", 7 },
  { "И2 держится", "записи", 7 },
  { "И2 держится", "свод", 10 },
  { "Строку отчёта", "находка", 2 },
  { "Вставить по весу", "строка", 11 },
  { "Вставить по весу", "строки", 12 },
  { "Приписать строку отчёта", "первая", 11 },
  { "Приписать строку отчёта", "строки", 12 },
  { "Отчёт", "записи", 7 },
  { "Отчёт той же длины", "записи", 7 },
  { "Отчёт той же длины", "строки", 12 },
  { "Это МожноУбрать", "приговор", 6 },
  { "Это НеТрогать", "приговор", 6 },
  { "И1 на паре", "разряд", 1 },
  { "И1 на паре", "приговор", 6 },
  { "Порог разряда", "разряд", 1 },
  { "Разряд обоснован", "находка", 2 },
  { "Разряд обоснован", "разряд", 1 },
  { "Приговор обоснован", "находка", 2 },
  { "Приговор обоснован", "разряд", 1 },
  { "Приговор обоснован", "приговор", 6 },
  { "Вес обоснован", "находка", 2 },
  { "Вес обоснован", "приговор", 6 },
  { "Вес обоснован", "вес", 3 },
};

static const fl_entry_table opis_diska_entry_table = {
  opis_diska_entry_types, 13,
  opis_diska_entry_fields, 17,
  opis_diska_entry_variants, 12,
  opis_diska_entry_params, 60
};

const fl_entry_table *opis_diska_entry(void) {
  return &opis_diska_entry_table;
}
