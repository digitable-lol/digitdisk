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
static const char *const opis_diska_names_2[] = { "разряд", "приговор", "вес" };
static const char *const opis_diska_names_3[] = { "путь", "решение" };
static const char *const opis_diska_names_4[] = { "кэш", "журнал", "сборка", "загрузка", "крупное", "неизвестное", "освободить" };
static const fl_value opis_diska_text_5 = { FL_STRING, { .string = { "/.cache/", 8, 8 } } };
static const fl_value opis_diska_text_6 = { FL_STRING, { .string = { "/cache/", 7, 7 } } };
static const fl_value opis_diska_text_7 = { FL_STRING, { .string = { "/Caches/", 8, 8 } } };
static const fl_value opis_diska_text_8 = { FL_STRING, { .string = { "/tmp/", 5, 5 } } };
static const fl_value opis_diska_text_9 = { FL_STRING, { .string = { "/log/", 5, 5 } } };
static const fl_value opis_diska_text_10 = { FL_STRING, { .string = { "/logs/", 6, 6 } } };
static const fl_value opis_diska_text_11 = { FL_STRING, { .string = { ".log", 4, 4 } } };
static const fl_value opis_diska_text_12 = { FL_STRING, { .string = { "/node_modules/", 14, 14 } } };
static const fl_value opis_diska_text_13 = { FL_STRING, { .string = { "/target/", 8, 8 } } };
static const fl_value opis_diska_text_14 = { FL_STRING, { .string = { "/build/", 7, 7 } } };
static const fl_value opis_diska_text_15 = { FL_STRING, { .string = { "/_build/", 8, 8 } } };
static const fl_value opis_diska_text_16 = { FL_STRING, { .string = { "/.gradle/", 9, 9 } } };
static const fl_value opis_diska_text_17 = { FL_STRING, { .string = { "/Downloads/", 11, 11 } } };
static const fl_value opis_diska_text_18 = { FL_STRING, { .string = { "/Загрузки/", 18, 10 } } };


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
 * Функция flang «Примета кэша».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_kesha(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error) {
  fl_value fl_t1 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_5, &fl_t1, error));
  bool fl_t2 = false;
  FL_TRY(fl_cond(ctx, fl_t1, &fl_t2, error));
  fl_value fl_t3 = fl_nothing();
  if (fl_t2) {
    fl_t3 = fl_flag(true);
  } else {
    fl_value fl_t4 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_6, &fl_t4, error));
    fl_t3 = fl_t4;
  }
  bool fl_t5 = false;
  FL_TRY(fl_cond(ctx, fl_t3, &fl_t5, error));
  fl_value fl_t6 = fl_nothing();
  if (fl_t5) {
    fl_t6 = fl_flag(true);
  } else {
    fl_value fl_t7 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_7, &fl_t7, error));
    fl_t6 = fl_t7;
  }
  bool fl_t8 = false;
  FL_TRY(fl_cond(ctx, fl_t6, &fl_t8, error));
  if (fl_t8) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t9 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_8, &fl_t9, error));
    *result = fl_t9;
    return FL_OK;
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
  fl_value fl_t10 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_9, &fl_t10, error));
  bool fl_t11 = false;
  FL_TRY(fl_cond(ctx, fl_t10, &fl_t11, error));
  fl_value fl_t12 = fl_nothing();
  if (fl_t11) {
    fl_t12 = fl_flag(true);
  } else {
    fl_value fl_t13 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_10, &fl_t13, error));
    fl_t12 = fl_t13;
  }
  bool fl_t14 = false;
  FL_TRY(fl_cond(ctx, fl_t12, &fl_t14, error));
  if (fl_t14) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t15 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_11, &fl_t15, error));
    *result = fl_t15;
    return FL_OK;
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
  fl_value fl_t16 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_12, &fl_t16, error));
  bool fl_t17 = false;
  FL_TRY(fl_cond(ctx, fl_t16, &fl_t17, error));
  fl_value fl_t18 = fl_nothing();
  if (fl_t17) {
    fl_t18 = fl_flag(true);
  } else {
    fl_value fl_t19 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_13, &fl_t19, error));
    fl_t18 = fl_t19;
  }
  bool fl_t20 = false;
  FL_TRY(fl_cond(ctx, fl_t18, &fl_t20, error));
  fl_value fl_t21 = fl_nothing();
  if (fl_t20) {
    fl_t21 = fl_flag(true);
  } else {
    fl_value fl_t22 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_14, &fl_t22, error));
    fl_t21 = fl_t22;
  }
  bool fl_t23 = false;
  FL_TRY(fl_cond(ctx, fl_t21, &fl_t23, error));
  fl_value fl_t24 = fl_nothing();
  if (fl_t23) {
    fl_t24 = fl_flag(true);
  } else {
    fl_value fl_t25 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_15, &fl_t25, error));
    fl_t24 = fl_t25;
  }
  bool fl_t26 = false;
  FL_TRY(fl_cond(ctx, fl_t24, &fl_t26, error));
  if (fl_t26) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t27 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_16, &fl_t27, error));
    *result = fl_t27;
    return FL_OK;
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
  fl_value fl_t28 = fl_nothing(); /* «содержит» */
  FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_17, &fl_t28, error));
  bool fl_t29 = false;
  FL_TRY(fl_cond(ctx, fl_t28, &fl_t29, error));
  if (fl_t29) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    fl_value fl_t30 = fl_nothing(); /* «содержит» */
    FL_TRY(fl_b_soderzhit(ctx, put, opis_diska_text_18, &fl_t30, error));
    *result = fl_t30;
    return FL_OK;
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
  fl_value fl_t31 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t31, error));
  fl_value fl_t32 = fl_nothing();
  FL_TRY(opis_diska_primeta_kesha(ctx, fl_t31, &fl_t32, error));
  bool fl_t33 = false;
  FL_TRY(fl_cond(ctx, fl_t32, &fl_t33, error));
  fl_value fl_t34 = fl_nothing();
  if (fl_t33) {
    fl_value fl_t35 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "Кэш", NULL, NULL, 0, &fl_t35, error));
    fl_t34 = fl_t35;
  } else {
    fl_value fl_t36 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t36, error));
    fl_value fl_t37 = fl_nothing();
    FL_TRY(opis_diska_primeta_zhurnala(ctx, fl_t36, &fl_t37, error));
    bool fl_t38 = false;
    FL_TRY(fl_cond(ctx, fl_t37, &fl_t38, error));
    fl_value fl_t39 = fl_nothing();
    if (fl_t38) {
      fl_value fl_t40 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "Журнал", NULL, NULL, 0, &fl_t40, error));
      fl_t39 = fl_t40;
    } else {
      fl_value fl_t41 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t41, error));
      fl_value fl_t42 = fl_nothing();
      FL_TRY(opis_diska_primeta_sborki(ctx, fl_t41, &fl_t42, error));
      bool fl_t43 = false;
      FL_TRY(fl_cond(ctx, fl_t42, &fl_t43, error));
      fl_value fl_t44 = fl_nothing();
      if (fl_t43) {
        fl_value fl_t45 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "Сборка", NULL, NULL, 0, &fl_t45, error));
        fl_t44 = fl_t45;
      } else {
        fl_value fl_t46 = fl_nothing();
        FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t46, error));
        fl_value fl_t47 = fl_nothing();
        FL_TRY(opis_diska_primeta_zagruzki(ctx, fl_t46, &fl_t47, error));
        bool fl_t48 = false;
        FL_TRY(fl_cond(ctx, fl_t47, &fl_t48, error));
        fl_value fl_t49 = fl_nothing();
        if (fl_t48) {
          fl_value fl_t50 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "Загрузка", NULL, NULL, 0, &fl_t50, error));
          fl_t49 = fl_t50;
        } else {
          fl_value fl_t51 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t51, error));
          fl_value fl_t52 = fl_nothing();
          FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t52, error));
          fl_value fl_t53 = fl_nothing();
          FL_TRY(fl_gte(ctx, fl_t51, fl_t52, &fl_t53, error));
          bool fl_t54 = false;
          FL_TRY(fl_cond(ctx, fl_t53, &fl_t54, error));
          fl_value fl_t55 = fl_nothing();
          if (fl_t54) {
            fl_value fl_t56 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Крупное", NULL, NULL, 0, &fl_t56, error));
            fl_t55 = fl_t56;
          } else {
            fl_value fl_t57 = fl_nothing();
            FL_TRY(fl_variant_new(ctx, "Неизвестное", NULL, NULL, 0, &fl_t57, error));
            fl_t55 = fl_t57;
          }
          fl_t49 = fl_t55;
        }
        fl_t44 = fl_t49;
      }
      fl_t39 = fl_t44;
    }
    fl_t34 = fl_t39;
  }
  const fl_value fl_t58 = fl_t34;
  fl_value fl_t59 = fl_nothing();
  FL_TRY(opis_diska_razryad_obosnovan(ctx, nahodka, fl_t58, &fl_t59, error));
  /* постусловие «Разряд обоснован приметой» */
  bool fl_t60 = false;
  FL_TRY(fl_post(ctx, fl_t59, "Разряд обоснован приметой", "Разряд находки", &fl_t60, error));
  if (!fl_t60) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «Разряд обоснован приметой» ядра «Опись диска»");
  }
  *result = fl_t58;
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
    fl_value fl_t61 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t61, error));
    fl_value fl_t62 = fl_nothing();
    FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t62, error));
    fl_value fl_t63 = fl_nothing();
    FL_TRY(fl_gte(ctx, fl_t61, fl_t62, &fl_t63, error));
    *result = fl_t63;
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
  fl_value fl_t64 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "вид", &fl_t64, error));
  if (fl_variant_is(fl_t64, "Каталог")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t64, "Файл")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(fl_t64, "Ссылка")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t64, error);
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
  fl_value fl_t65 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "вид", &fl_t65, error));
  if (fl_variant_is(fl_t65, "Ссылка")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t65, "Файл")) {
    *result = fl_flag(false);
    return FL_OK;
  } else if (fl_variant_is(fl_t65, "Каталог")) {
    *result = fl_flag(false);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t65, error);
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
  fl_value fl_t66 = fl_nothing();
  FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t66, error));
  bool fl_t67 = false;
  FL_TRY(fl_cond(ctx, fl_t66, &fl_t67, error));
  fl_value fl_t68 = fl_nothing();
  if (fl_t67) {
    fl_value fl_t69 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t69, error));
    fl_t68 = fl_t69;
  } else {
    fl_value fl_t70 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t70, error));
    fl_value fl_t71 = fl_nothing();
    FL_TRY(fl_gte(ctx, fl_t70, porog, &fl_t71, error));
    bool fl_t72 = false;
    FL_TRY(fl_cond(ctx, fl_t71, &fl_t72, error));
    fl_value fl_t73 = fl_nothing();
    if (fl_t72) {
      fl_value fl_t74 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "МожноУбрать", NULL, NULL, 0, &fl_t74, error));
      fl_t73 = fl_t74;
    } else {
      fl_value fl_t75 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t75, error));
      fl_t73 = fl_t75;
    }
    fl_t68 = fl_t73;
  }
  const fl_value fl_t76 = fl_t68;
  fl_value fl_t77 = fl_nothing();
  FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t77, error));
  bool fl_t78 = false;
  FL_TRY(fl_cond(ctx, fl_t77, &fl_t78, error));
  fl_value fl_t79 = fl_nothing();
  if (fl_t78) {
    fl_value fl_t80 = fl_nothing();
    FL_TRY(opis_diska_eto_mozhnoubrat(ctx, fl_t76, &fl_t80, error));
    bool fl_t81 = false;
    FL_TRY(fl_cond(ctx, fl_t80, &fl_t81, error));
    fl_value fl_t82 = fl_nothing();
    if (fl_t81) {
      fl_t82 = fl_flag(false);
    } else {
      fl_t82 = fl_flag(true);
    }
    fl_t79 = fl_t82;
  } else {
    fl_t79 = fl_flag(true);
  }
  /* постусловие «Каталог не убирается» */
  bool fl_t83 = false;
  FL_TRY(fl_post(ctx, fl_t79, "Каталог не убирается", "Приговор мусора", &fl_t83, error));
  if (!fl_t83) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «Каталог не убирается» ядра «Опись диска»");
  }
  *result = fl_t76;
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
  fl_value fl_t84 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "доступен", &fl_t84, error));
  bool fl_t85 = false;
  FL_TRY(fl_cond(ctx, fl_t84, &fl_t85, error));
  fl_value fl_t86 = fl_nothing();
  if (fl_t85) {
    fl_t86 = fl_flag(false);
  } else {
    fl_t86 = fl_flag(true);
  }
  bool fl_t87 = false;
  FL_TRY(fl_cond(ctx, fl_t86, &fl_t87, error));
  fl_value fl_t88 = fl_nothing();
  if (fl_t87) {
    fl_value fl_t89 = fl_nothing();
    FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t89, error));
    fl_t88 = fl_t89;
  } else {
    fl_value fl_t90 = fl_nothing();
    FL_TRY(opis_diska_ssylka(ctx, nahodka, &fl_t90, error));
    bool fl_t91 = false;
    FL_TRY(fl_cond(ctx, fl_t90, &fl_t91, error));
    fl_value fl_t92 = fl_nothing();
    if (fl_t91) {
      fl_value fl_t93 = fl_nothing();
      FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t93, error));
      fl_t92 = fl_t93;
    } else {
      fl_value fl_t94 = fl_nothing();
      if (fl_variant_is(razryad, "Кэш")) {
        fl_value fl_t95 = fl_nothing();
        FL_TRY(opis_diska_porog_kesha(ctx, &fl_t95, error));
        fl_value fl_t96 = fl_nothing();
        FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t95, &fl_t96, error));
        fl_t94 = fl_t96;
      } else if (fl_variant_is(razryad, "Сборка")) {
        fl_value fl_t97 = fl_nothing();
        FL_TRY(opis_diska_porog_kesha(ctx, &fl_t97, error));
        fl_value fl_t98 = fl_nothing();
        FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t97, &fl_t98, error));
        fl_t94 = fl_t98;
      } else if (fl_variant_is(razryad, "Журнал")) {
        fl_value fl_t99 = fl_nothing();
        FL_TRY(opis_diska_porog_zhurnala(ctx, &fl_t99, error));
        fl_value fl_t100 = fl_nothing();
        FL_TRY(opis_diska_prigovor_musora(ctx, nahodka, fl_t99, &fl_t100, error));
        fl_t94 = fl_t100;
      } else if (fl_variant_is(razryad, "Загрузка")) {
        fl_value fl_t101 = fl_nothing();
        FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t101, error));
        fl_value fl_t102 = fl_nothing();
        FL_TRY(opis_diska_porog_zagruzki(ctx, &fl_t102, error));
        fl_value fl_t103 = fl_nothing();
        FL_TRY(fl_gte(ctx, fl_t101, fl_t102, &fl_t103, error));
        bool fl_t104 = false;
        FL_TRY(fl_cond(ctx, fl_t103, &fl_t104, error));
        fl_value fl_t105 = fl_nothing();
        if (fl_t104) {
          fl_value fl_t106 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t106, error));
          fl_t105 = fl_t106;
        } else {
          fl_value fl_t107 = fl_nothing();
          FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t107, error));
          fl_t105 = fl_t107;
        }
        fl_t94 = fl_t105;
      } else if (fl_variant_is(razryad, "Крупное")) {
        fl_value fl_t108 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "Спросить", NULL, NULL, 0, &fl_t108, error));
        fl_t94 = fl_t108;
      } else if (fl_variant_is(razryad, "Неизвестное")) {
        fl_value fl_t109 = fl_nothing();
        FL_TRY(fl_variant_new(ctx, "НеТрогать", NULL, NULL, 0, &fl_t109, error));
        fl_t94 = fl_t109;
      } else {
        return fl_match_fail(ctx, razryad, error);
      }
      fl_t92 = fl_t94;
    }
    fl_t88 = fl_t92;
  }
  const fl_value fl_t110 = fl_t88;
  fl_value fl_t111 = fl_nothing();
  FL_TRY(opis_diska_prigovor_obosnovan(ctx, nahodka, razryad, fl_t110, &fl_t111, error));
  /* постусловие «Приговор обоснован» */
  bool fl_t112 = false;
  FL_TRY(fl_post(ctx, fl_t111, "Приговор обоснован", "Приговор находки", &fl_t112, error));
  if (!fl_t112) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «Приговор обоснован» ядра «Опись диска»");
  }
  *result = fl_t110;
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
  fl_value fl_t113 = fl_nothing();
  if (fl_variant_is(prigovor, "НеТрогать")) {
    fl_t113 = fl_number(0.0);
  } else if (fl_variant_is(prigovor, "МожноУбрать")) {
    fl_value fl_t114 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t114, error));
    fl_t113 = fl_t114;
  } else if (fl_variant_is(prigovor, "Спросить")) {
    fl_value fl_t115 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t115, error));
    fl_t113 = fl_t115;
  } else {
    return fl_match_fail(ctx, prigovor, error);
  }
  const fl_value fl_t116 = fl_t113;
  fl_value fl_t117 = fl_nothing();
  FL_TRY(opis_diska_ves_obosnovan(ctx, nahodka, prigovor, fl_t116, &fl_t117, error));
  /* постусловие «Вес обоснован» */
  bool fl_t118 = false;
  FL_TRY(fl_post(ctx, fl_t117, "Вес обоснован", "Вес находки", &fl_t118, error));
  if (!fl_t118) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «Вес обоснован» ядра «Опись диска»");
  }
  *result = fl_t116;
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
  fl_value fl_t119 = fl_nothing();
  FL_TRY(fl_gte(ctx, ves, fl_number(0.0), &fl_t119, error));
  bool fl_t120 = false;
  FL_TRY(fl_cond(ctx, fl_t119, &fl_t120, error));
  if (fl_t120) {
    fl_value fl_t121 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t121, error));
    fl_value fl_t122 = fl_nothing();
    FL_TRY(fl_lte(ctx, ves, fl_t121, &fl_t122, error));
    *result = fl_t122;
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
  fl_value fl_t123 = fl_nothing();
  FL_TRY(opis_diska_razryad_nahodki(ctx, nahodka, &fl_t123, error));
  const fl_value razryad = fl_t123; /* пусть «разряд» */
  fl_value fl_t124 = fl_nothing();
  FL_TRY(opis_diska_prigovor_nahodki(ctx, nahodka, razryad, &fl_t124, error));
  const fl_value prigovor = fl_t124; /* пусть «приговор» */
  fl_value fl_t125 = fl_nothing();
  FL_TRY(opis_diska_ves_nahodki(ctx, nahodka, prigovor, &fl_t125, error));
  fl_value fl_t127[3];
  fl_t127[0] = razryad; /* «разряд» */
  fl_t127[1] = prigovor; /* «приговор» */
  fl_t127[2] = fl_t125; /* «вес» */
  fl_value fl_t126 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_2, fl_t127, 3, &fl_t126, error));
  const fl_value fl_t128 = fl_t126;
  fl_value fl_t129 = fl_nothing();
  FL_TRY(opis_diska_i1_derzhitsya(ctx, fl_t128, &fl_t129, error));
  /* постусловие «И1: убрать можно только мусор» */
  bool fl_t130 = false;
  FL_TRY(fl_post(ctx, fl_t129, "И1: убрать можно только мусор", "Решить находку", &fl_t130, error));
  if (!fl_t130) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «И1: убрать можно только мусор» ядра «Опись диска»");
  }
  *result = fl_t128;
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
  fl_value fl_t131 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "отобразить", &fl_t131, error));
  fl_value *fl_t132 = NULL;
  size_t fl_t133 = 0;
  FL_TRY(fl_list_alloc(ctx, fl_t131.as.list.count, &fl_t132, error));
  for (size_t fl_t134 = 0; fl_t134 < fl_t131.as.list.count; fl_t134 += 1) {
    const fl_value nahodka = fl_t131.as.list.items[fl_t134]; /* «находка» */
    fl_value fl_t135 = fl_nothing();
    FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, &fl_t135, error));
    fl_t132[fl_t133] = fl_t135;
    fl_t133 += 1;
  }
  const fl_value fl_t136 = fl_list(fl_t132, fl_t133);
  fl_value fl_t137 = fl_nothing();
  FL_TRY(opis_diska_i1_derzhitsya_vsyudu(ctx, fl_t136, &fl_t137, error));
  /* постусловие «И1 всюду: убрать можно только мусор» */
  bool fl_t138 = false;
  FL_TRY(fl_post(ctx, fl_t137, "И1 всюду: убрать можно только мусор", "Решить всё", &fl_t138, error));
  if (!fl_t138) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «И1 всюду: убрать можно только мусор» ядра «Опись диска»");
  }
  *result = fl_t136;
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
  fl_value fl_t139 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "приговор", &fl_t139, error));
  if (fl_variant_is(fl_t139, "МожноУбрать")) {
    fl_value fl_t140 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "разряд", &fl_t140, error));
    if (fl_variant_is(fl_t140, "Кэш")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t140, "Журнал")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t140, "Сборка")) {
      *result = fl_flag(true);
      return FL_OK;
    } else if (fl_variant_is(fl_t140, "Загрузка")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(fl_t140, "Крупное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else if (fl_variant_is(fl_t140, "Неизвестное")) {
      *result = fl_flag(false);
      return FL_OK;
    } else {
      return fl_match_fail(ctx, fl_t140, error);
    }
  } else if (fl_variant_is(fl_t139, "Спросить")) {
    *result = fl_flag(true);
    return FL_OK;
  } else if (fl_variant_is(fl_t139, "НеТрогать")) {
    *result = fl_flag(true);
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t139, error);
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
  fl_value fl_t141 = fl_nothing();
  FL_TRY(fl_require_list(ctx, resheniya, "свёртка", &fl_t141, error));
  fl_value akk = fl_flag(true); /* «акк» */
  for (size_t fl_t142 = 0; fl_t142 < fl_t141.as.list.count; fl_t142 += 1) {
    const fl_value reshenie = fl_t141.as.list.items[fl_t142]; /* «решение» */
    bool fl_t143 = false;
    FL_TRY(fl_cond(ctx, akk, &fl_t143, error));
    fl_value fl_t144 = fl_nothing();
    if (fl_t143) {
      fl_value fl_t145 = fl_nothing();
      FL_TRY(opis_diska_i1_derzhitsya(ctx, reshenie, &fl_t145, error));
      fl_t144 = fl_t145;
    } else {
      fl_t144 = fl_flag(false);
    }
    akk = fl_t144;
  }
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
  fl_value fl_t147[7];
  fl_t147[0] = fl_number(0.0); /* «кэш» */
  fl_t147[1] = fl_number(0.0); /* «журнал» */
  fl_t147[2] = fl_number(0.0); /* «сборка» */
  fl_t147[3] = fl_number(0.0); /* «загрузка» */
  fl_t147[4] = fl_number(0.0); /* «крупное» */
  fl_t147[5] = fl_number(0.0); /* «неизвестное» */
  fl_t147[6] = fl_number(0.0); /* «освободить» */
  fl_value fl_t146 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t147, 7, &fl_t146, error));
  *result = fl_t146;
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
  fl_value fl_t148 = fl_nothing();
  fl_value fl_t149 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "приговор", &fl_t149, error));
  if (fl_variant_is(fl_t149, "МожноУбрать")) {
    fl_value fl_t150 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t150, error));
    fl_t148 = fl_t150;
  } else if (fl_variant_is(fl_t149, "Спросить")) {
    fl_t148 = fl_number(0.0);
  } else if (fl_variant_is(fl_t149, "НеТрогать")) {
    fl_t148 = fl_number(0.0);
  } else {
    return fl_match_fail(ctx, fl_t149, error);
  }
  const fl_value ubrat = fl_t148; /* пусть «убрать» */
  fl_value fl_t151 = fl_nothing();
  FL_TRY(fl_field_get(ctx, reshenie, "разряд", &fl_t151, error));
  if (fl_variant_is(fl_t151, "Кэш")) {
    fl_value fl_t152 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t152, error));
    fl_value fl_t153 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t153, error));
    fl_value fl_t154 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t152, fl_t153, &fl_t154, error));
    fl_value fl_t155 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t155, error));
    fl_value fl_t156 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t156, error));
    fl_value fl_t157 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t157, error));
    fl_value fl_t158 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t158, error));
    fl_value fl_t159 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t159, error));
    fl_value fl_t160 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t160, error));
    fl_value fl_t161 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t160, ubrat, &fl_t161, error));
    fl_value fl_t163[7];
    fl_t163[0] = fl_t154; /* «кэш» */
    fl_t163[1] = fl_t155; /* «журнал» */
    fl_t163[2] = fl_t156; /* «сборка» */
    fl_t163[3] = fl_t157; /* «загрузка» */
    fl_t163[4] = fl_t158; /* «крупное» */
    fl_t163[5] = fl_t159; /* «неизвестное» */
    fl_t163[6] = fl_t161; /* «освободить» */
    fl_value fl_t162 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t163, 7, &fl_t162, error));
    *result = fl_t162;
    return FL_OK;
  } else if (fl_variant_is(fl_t151, "Журнал")) {
    fl_value fl_t164 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t164, error));
    fl_value fl_t165 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t165, error));
    fl_value fl_t166 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t166, error));
    fl_value fl_t167 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t165, fl_t166, &fl_t167, error));
    fl_value fl_t168 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t168, error));
    fl_value fl_t169 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t169, error));
    fl_value fl_t170 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t170, error));
    fl_value fl_t171 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t171, error));
    fl_value fl_t172 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t172, error));
    fl_value fl_t173 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t172, ubrat, &fl_t173, error));
    fl_value fl_t175[7];
    fl_t175[0] = fl_t164; /* «кэш» */
    fl_t175[1] = fl_t167; /* «журнал» */
    fl_t175[2] = fl_t168; /* «сборка» */
    fl_t175[3] = fl_t169; /* «загрузка» */
    fl_t175[4] = fl_t170; /* «крупное» */
    fl_t175[5] = fl_t171; /* «неизвестное» */
    fl_t175[6] = fl_t173; /* «освободить» */
    fl_value fl_t174 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t175, 7, &fl_t174, error));
    *result = fl_t174;
    return FL_OK;
  } else if (fl_variant_is(fl_t151, "Сборка")) {
    fl_value fl_t176 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t176, error));
    fl_value fl_t177 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t177, error));
    fl_value fl_t178 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t178, error));
    fl_value fl_t179 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t179, error));
    fl_value fl_t180 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t178, fl_t179, &fl_t180, error));
    fl_value fl_t181 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t181, error));
    fl_value fl_t182 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t182, error));
    fl_value fl_t183 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t183, error));
    fl_value fl_t184 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t184, error));
    fl_value fl_t185 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t184, ubrat, &fl_t185, error));
    fl_value fl_t187[7];
    fl_t187[0] = fl_t176; /* «кэш» */
    fl_t187[1] = fl_t177; /* «журнал» */
    fl_t187[2] = fl_t180; /* «сборка» */
    fl_t187[3] = fl_t181; /* «загрузка» */
    fl_t187[4] = fl_t182; /* «крупное» */
    fl_t187[5] = fl_t183; /* «неизвестное» */
    fl_t187[6] = fl_t185; /* «освободить» */
    fl_value fl_t186 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t187, 7, &fl_t186, error));
    *result = fl_t186;
    return FL_OK;
  } else if (fl_variant_is(fl_t151, "Загрузка")) {
    fl_value fl_t188 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t188, error));
    fl_value fl_t189 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t189, error));
    fl_value fl_t190 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t190, error));
    fl_value fl_t191 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t191, error));
    fl_value fl_t192 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t192, error));
    fl_value fl_t193 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t191, fl_t192, &fl_t193, error));
    fl_value fl_t194 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t194, error));
    fl_value fl_t195 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t195, error));
    fl_value fl_t196 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t196, error));
    fl_value fl_t197 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t196, ubrat, &fl_t197, error));
    fl_value fl_t199[7];
    fl_t199[0] = fl_t188; /* «кэш» */
    fl_t199[1] = fl_t189; /* «журнал» */
    fl_t199[2] = fl_t190; /* «сборка» */
    fl_t199[3] = fl_t193; /* «загрузка» */
    fl_t199[4] = fl_t194; /* «крупное» */
    fl_t199[5] = fl_t195; /* «неизвестное» */
    fl_t199[6] = fl_t197; /* «освободить» */
    fl_value fl_t198 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t199, 7, &fl_t198, error));
    *result = fl_t198;
    return FL_OK;
  } else if (fl_variant_is(fl_t151, "Крупное")) {
    fl_value fl_t200 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t200, error));
    fl_value fl_t201 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t201, error));
    fl_value fl_t202 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t202, error));
    fl_value fl_t203 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t203, error));
    fl_value fl_t204 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t204, error));
    fl_value fl_t205 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t205, error));
    fl_value fl_t206 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t204, fl_t205, &fl_t206, error));
    fl_value fl_t207 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t207, error));
    fl_value fl_t208 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t208, error));
    fl_value fl_t209 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t208, ubrat, &fl_t209, error));
    fl_value fl_t211[7];
    fl_t211[0] = fl_t200; /* «кэш» */
    fl_t211[1] = fl_t201; /* «журнал» */
    fl_t211[2] = fl_t202; /* «сборка» */
    fl_t211[3] = fl_t203; /* «загрузка» */
    fl_t211[4] = fl_t206; /* «крупное» */
    fl_t211[5] = fl_t207; /* «неизвестное» */
    fl_t211[6] = fl_t209; /* «освободить» */
    fl_value fl_t210 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t211, 7, &fl_t210, error));
    *result = fl_t210;
    return FL_OK;
  } else if (fl_variant_is(fl_t151, "Неизвестное")) {
    fl_value fl_t212 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "кэш", &fl_t212, error));
    fl_value fl_t213 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "журнал", &fl_t213, error));
    fl_value fl_t214 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "сборка", &fl_t214, error));
    fl_value fl_t215 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "загрузка", &fl_t215, error));
    fl_value fl_t216 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "крупное", &fl_t216, error));
    fl_value fl_t217 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "неизвестное", &fl_t217, error));
    fl_value fl_t218 = fl_nothing();
    FL_TRY(fl_field_get(ctx, reshenie, "вес", &fl_t218, error));
    fl_value fl_t219 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t217, fl_t218, &fl_t219, error));
    fl_value fl_t220 = fl_nothing();
    FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t220, error));
    fl_value fl_t221 = fl_nothing();
    FL_TRY(fl_add(ctx, fl_t220, ubrat, &fl_t221, error));
    fl_value fl_t223[7];
    fl_t223[0] = fl_t212; /* «кэш» */
    fl_t223[1] = fl_t213; /* «журнал» */
    fl_t223[2] = fl_t214; /* «сборка» */
    fl_t223[3] = fl_t215; /* «загрузка» */
    fl_t223[4] = fl_t216; /* «крупное» */
    fl_t223[5] = fl_t219; /* «неизвестное» */
    fl_t223[6] = fl_t221; /* «освободить» */
    fl_value fl_t222 = fl_nothing();
    FL_TRY(fl_record_new(ctx, opis_diska_names_4, fl_t223, 7, &fl_t222, error));
    *result = fl_t222;
    return FL_OK;
  } else {
    return fl_match_fail(ctx, fl_t151, error);
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
  fl_value fl_t224 = fl_nothing();
  FL_TRY(opis_diska_reshit_vsyo(ctx, zapisi, &fl_t224, error));
  fl_value fl_t225 = fl_nothing();
  FL_TRY(fl_require_list(ctx, fl_t224, "свёртка", &fl_t225, error));
  fl_value fl_t226 = fl_nothing();
  FL_TRY(opis_diska_pustoy_svod(ctx, &fl_t226, error));
  fl_value svod = fl_t226; /* «свод» */
  for (size_t fl_t227 = 0; fl_t227 < fl_t225.as.list.count; fl_t227 += 1) {
    const fl_value reshenie = fl_t225.as.list.items[fl_t227]; /* «решение» */
    fl_value fl_t228 = fl_nothing();
    FL_TRY(opis_diska_pribavit_reshenie(ctx, svod, reshenie, &fl_t228, error));
    svod = fl_t228;
  }
  const fl_value fl_t229 = svod;
  fl_value fl_t230 = fl_nothing();
  FL_TRY(opis_diska_i2_derzhitsya(ctx, zapisi, fl_t229, &fl_t230, error));
  /* постусловие «И2: освобождаемое не больше убираемого» */
  bool fl_t231 = false;
  FL_TRY(fl_post(ctx, fl_t230, "И2: освобождаемое не больше убираемого", "Свести", &fl_t231, error));
  if (!fl_t231) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «И2: освобождаемое не больше убираемого» ядра «Опись диска»");
  }
  *result = fl_t229;
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
  fl_value fl_t232 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "свёртка", &fl_t232, error));
  fl_value akk = fl_number(0.0); /* «акк» */
  for (size_t fl_t233 = 0; fl_t233 < fl_t232.as.list.count; fl_t233 += 1) {
    const fl_value nahodka = fl_t232.as.list.items[fl_t233]; /* «находка» */
    fl_value fl_t234 = fl_nothing();
    fl_value fl_t235 = fl_nothing();
    FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, &fl_t235, error));
    fl_value fl_t236 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t235, "приговор", &fl_t236, error));
    if (fl_variant_is(fl_t236, "МожноУбрать")) {
      fl_value fl_t237 = fl_nothing();
      FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t237, error));
      fl_value fl_t238 = fl_nothing();
      FL_TRY(fl_add(ctx, akk, fl_t237, &fl_t238, error));
      fl_t234 = fl_t238;
    } else if (fl_variant_is(fl_t236, "Спросить")) {
      fl_t234 = akk;
    } else if (fl_variant_is(fl_t236, "НеТрогать")) {
      fl_t234 = akk;
    } else {
      return fl_match_fail(ctx, fl_t236, error);
    }
    akk = fl_t234;
  }
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
  fl_value fl_t239 = fl_nothing();
  FL_TRY(fl_field_get(ctx, svod, "освободить", &fl_t239, error));
  fl_value fl_t240 = fl_nothing();
  FL_TRY(opis_diska_summa_razmerov_ubiraemyh(ctx, zapisi, &fl_t240, error));
  fl_value fl_t241 = fl_nothing();
  FL_TRY(fl_lte(ctx, fl_t239, fl_t240, &fl_t241, error));
  *result = fl_t241;
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
  fl_value fl_t242 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t242, error));
  fl_value fl_t243 = fl_nothing();
  FL_TRY(opis_diska_reshit_nahodku(ctx, nahodka, &fl_t243, error));
  fl_value fl_t245[2];
  fl_t245[0] = fl_t242; /* «путь» */
  fl_t245[1] = fl_t243; /* «решение» */
  fl_value fl_t244 = fl_nothing();
  FL_TRY(fl_record_new(ctx, opis_diska_names_3, fl_t245, 2, &fl_t244, error));
  *result = fl_t244;
  return FL_OK;
}

/* Тело «Вставить по весу»; глубину считает обёртка ниже. */
static fl_status opis_diska_vstavit_po_vesu_body(fl_ctx *ctx, fl_value stroka, fl_value stroki, fl_value *result, fl_error *error) {
  if (fl_chain_empty(stroki)) {
    fl_value *fl_t246 = NULL;
    FL_TRY(fl_list_alloc(ctx, 1, &fl_t246, error));
    fl_t246[0] = stroka;
    *result = fl_list(fl_t246, 1);
    return FL_OK;
  } else if (fl_chain_cons(stroki)) {
    const fl_value golova = fl_chain_head(stroki); /* голова «голова» */
    const fl_value hvost = fl_chain_tail(stroki); /* хвост «хвост» */
    fl_value fl_t247 = fl_nothing();
    FL_TRY(fl_field_get(ctx, stroka, "решение", &fl_t247, error));
    fl_value fl_t248 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t247, "вес", &fl_t248, error));
    fl_value fl_t249 = fl_nothing();
    FL_TRY(fl_field_get(ctx, golova, "решение", &fl_t249, error));
    fl_value fl_t250 = fl_nothing();
    FL_TRY(fl_field_get(ctx, fl_t249, "вес", &fl_t250, error));
    fl_value fl_t251 = fl_nothing();
    FL_TRY(fl_gte(ctx, fl_t248, fl_t250, &fl_t251, error));
    bool fl_t252 = false;
    FL_TRY(fl_cond(ctx, fl_t251, &fl_t252, error));
    if (fl_t252) {
      return opis_diska_pripisat_stroku_otchyota(ctx, stroka, stroki, result, error);
    } else {
      fl_value fl_t253 = fl_nothing();
      FL_TRY(opis_diska_vstavit_po_vesu(ctx, stroka, hvost, &fl_t253, error));
      return opis_diska_pripisat_stroku_otchyota(ctx, golova, fl_t253, result, error);
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
    const fl_status status = opis_diska_vstavit_po_vesu_body(ctx, stroka, stroki, result, error);
    fl_leave(ctx);
    return status;
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
  fl_value fl_t254 = fl_nothing();
  FL_TRY(fl_require_list(ctx, stroki, "свёртка", &fl_t254, error));
  fl_value *fl_t255 = NULL;
  FL_TRY(fl_list_alloc(ctx, 1, &fl_t255, error));
  fl_t255[0] = pervaya;
  fl_value akk = fl_list(fl_t255, 1); /* «акк» */
  for (size_t fl_t256 = 0; fl_t256 < fl_t254.as.list.count; fl_t256 += 1) {
    const fl_value el = fl_t254.as.list.items[fl_t256]; /* «эл» */
    fl_value fl_t257 = fl_nothing(); /* «добавить» */
    FL_TRY(fl_b_dobavit(ctx, el, akk, &fl_t257, error));
    akk = fl_t257;
  }
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
  fl_value fl_t258 = fl_nothing();
  FL_TRY(fl_require_list(ctx, zapisi, "свёртка", &fl_t258, error));
  fl_value akk = fl_list(NULL, 0); /* «акк» */
  for (size_t fl_t259 = 0; fl_t259 < fl_t258.as.list.count; fl_t259 += 1) {
    const fl_value nahodka = fl_t258.as.list.items[fl_t259]; /* «находка» */
    fl_value fl_t260 = fl_nothing();
    FL_TRY(opis_diska_stroku_otchyota(ctx, nahodka, &fl_t260, error));
    fl_value fl_t261 = fl_nothing();
    FL_TRY(opis_diska_vstavit_po_vesu(ctx, fl_t260, akk, &fl_t261, error));
    akk = fl_t261;
  }
  const fl_value fl_t262 = akk;
  fl_value fl_t263 = fl_nothing();
  FL_TRY(opis_diska_otchyot_toy_zhe_dliny(ctx, zapisi, fl_t262, &fl_t263, error));
  /* постусловие «Отчёт той же длины» */
  bool fl_t264 = false;
  FL_TRY(fl_post(ctx, fl_t263, "Отчёт той же длины", "Отчёт", &fl_t264, error));
  if (!fl_t264) {
    return fl_fail(ctx, error, "DIGITDISK_RULE", "%s", "нарушено постусловие «Отчёт той же длины» ядра «Опись диска»");
  }
  *result = fl_t262;
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
  fl_value fl_t265 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, stroki, &fl_t265, error));
  fl_value fl_t266 = fl_nothing(); /* «длина» */
  FL_TRY(fl_b_dlina(ctx, zapisi, &fl_t266, error));
  *result = fl_flag(fl_equal(fl_t265, fl_t266));
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
  fl_value fl_t267 = fl_nothing();
  FL_TRY(opis_diska_eto_mozhnoubrat(ctx, prigovor, &fl_t267, error));
  bool fl_t268 = false;
  FL_TRY(fl_cond(ctx, fl_t267, &fl_t268, error));
  fl_value fl_t269 = fl_nothing();
  if (fl_t268) {
    fl_t269 = fl_flag(false);
  } else {
    fl_t269 = fl_flag(true);
  }
  bool fl_t270 = false;
  FL_TRY(fl_cond(ctx, fl_t269, &fl_t270, error));
  if (fl_t270) {
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
  fl_value fl_t271 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t271, error));
  fl_value fl_t272 = fl_nothing();
  FL_TRY(opis_diska_primeta_kesha(ctx, fl_t271, &fl_t272, error));
  const fl_value kesh = fl_t272; /* пусть «кэш» */
  fl_value fl_t273 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t273, error));
  fl_value fl_t274 = fl_nothing();
  FL_TRY(opis_diska_primeta_zhurnala(ctx, fl_t273, &fl_t274, error));
  const fl_value zhurnal = fl_t274; /* пусть «журнал» */
  fl_value fl_t275 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t275, error));
  fl_value fl_t276 = fl_nothing();
  FL_TRY(opis_diska_primeta_sborki(ctx, fl_t275, &fl_t276, error));
  const fl_value sborka = fl_t276; /* пусть «сборка» */
  fl_value fl_t277 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "путь", &fl_t277, error));
  fl_value fl_t278 = fl_nothing();
  FL_TRY(opis_diska_primeta_zagruzki(ctx, fl_t277, &fl_t278, error));
  const fl_value zagruzka = fl_t278; /* пусть «загрузка» */
  fl_value fl_t279 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t279, error));
  fl_value fl_t280 = fl_nothing();
  FL_TRY(opis_diska_porog_krupnogo(ctx, &fl_t280, error));
  fl_value fl_t281 = fl_nothing();
  FL_TRY(fl_gte(ctx, fl_t279, fl_t280, &fl_t281, error));
  const fl_value krupnoe = fl_t281; /* пусть «крупное» */
  if (fl_variant_is(razryad, "Кэш")) {
    *result = kesh;
    return FL_OK;
  } else if (fl_variant_is(razryad, "Журнал")) {
    bool fl_t282 = false;
    FL_TRY(fl_cond(ctx, zhurnal, &fl_t282, error));
    if (fl_t282) {
      bool fl_t283 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t283, error));
      if (fl_t283) {
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
    bool fl_t284 = false;
    FL_TRY(fl_cond(ctx, sborka, &fl_t284, error));
    fl_value fl_t285 = fl_nothing();
    if (fl_t284) {
      bool fl_t286 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t286, error));
      fl_value fl_t287 = fl_nothing();
      if (fl_t286) {
        fl_t287 = fl_flag(false);
      } else {
        fl_t287 = fl_flag(true);
      }
      fl_t285 = fl_t287;
    } else {
      fl_t285 = fl_flag(false);
    }
    bool fl_t288 = false;
    FL_TRY(fl_cond(ctx, fl_t285, &fl_t288, error));
    if (fl_t288) {
      bool fl_t289 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t289, error));
      if (fl_t289) {
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
    bool fl_t290 = false;
    FL_TRY(fl_cond(ctx, zagruzka, &fl_t290, error));
    fl_value fl_t291 = fl_nothing();
    if (fl_t290) {
      bool fl_t292 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t292, error));
      fl_value fl_t293 = fl_nothing();
      if (fl_t292) {
        fl_t293 = fl_flag(false);
      } else {
        fl_t293 = fl_flag(true);
      }
      fl_t291 = fl_t293;
    } else {
      fl_t291 = fl_flag(false);
    }
    bool fl_t294 = false;
    FL_TRY(fl_cond(ctx, fl_t291, &fl_t294, error));
    fl_value fl_t295 = fl_nothing();
    if (fl_t294) {
      bool fl_t296 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t296, error));
      fl_value fl_t297 = fl_nothing();
      if (fl_t296) {
        fl_t297 = fl_flag(false);
      } else {
        fl_t297 = fl_flag(true);
      }
      fl_t295 = fl_t297;
    } else {
      fl_t295 = fl_flag(false);
    }
    bool fl_t298 = false;
    FL_TRY(fl_cond(ctx, fl_t295, &fl_t298, error));
    if (fl_t298) {
      bool fl_t299 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t299, error));
      if (fl_t299) {
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
    bool fl_t300 = false;
    FL_TRY(fl_cond(ctx, krupnoe, &fl_t300, error));
    fl_value fl_t301 = fl_nothing();
    if (fl_t300) {
      bool fl_t302 = false;
      FL_TRY(fl_cond(ctx, kesh, &fl_t302, error));
      fl_value fl_t303 = fl_nothing();
      if (fl_t302) {
        fl_t303 = fl_flag(false);
      } else {
        fl_t303 = fl_flag(true);
      }
      fl_t301 = fl_t303;
    } else {
      fl_t301 = fl_flag(false);
    }
    bool fl_t304 = false;
    FL_TRY(fl_cond(ctx, fl_t301, &fl_t304, error));
    fl_value fl_t305 = fl_nothing();
    if (fl_t304) {
      bool fl_t306 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t306, error));
      fl_value fl_t307 = fl_nothing();
      if (fl_t306) {
        fl_t307 = fl_flag(false);
      } else {
        fl_t307 = fl_flag(true);
      }
      fl_t305 = fl_t307;
    } else {
      fl_t305 = fl_flag(false);
    }
    bool fl_t308 = false;
    FL_TRY(fl_cond(ctx, fl_t305, &fl_t308, error));
    fl_value fl_t309 = fl_nothing();
    if (fl_t308) {
      bool fl_t310 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t310, error));
      fl_value fl_t311 = fl_nothing();
      if (fl_t310) {
        fl_t311 = fl_flag(false);
      } else {
        fl_t311 = fl_flag(true);
      }
      fl_t309 = fl_t311;
    } else {
      fl_t309 = fl_flag(false);
    }
    bool fl_t312 = false;
    FL_TRY(fl_cond(ctx, fl_t309, &fl_t312, error));
    if (fl_t312) {
      bool fl_t313 = false;
      FL_TRY(fl_cond(ctx, zagruzka, &fl_t313, error));
      if (fl_t313) {
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
    bool fl_t314 = false;
    FL_TRY(fl_cond(ctx, kesh, &fl_t314, error));
    fl_value fl_t315 = fl_nothing();
    if (fl_t314) {
      fl_t315 = fl_flag(false);
    } else {
      fl_t315 = fl_flag(true);
    }
    bool fl_t316 = false;
    FL_TRY(fl_cond(ctx, fl_t315, &fl_t316, error));
    fl_value fl_t317 = fl_nothing();
    if (fl_t316) {
      bool fl_t318 = false;
      FL_TRY(fl_cond(ctx, zhurnal, &fl_t318, error));
      fl_value fl_t319 = fl_nothing();
      if (fl_t318) {
        fl_t319 = fl_flag(false);
      } else {
        fl_t319 = fl_flag(true);
      }
      fl_t317 = fl_t319;
    } else {
      fl_t317 = fl_flag(false);
    }
    bool fl_t320 = false;
    FL_TRY(fl_cond(ctx, fl_t317, &fl_t320, error));
    fl_value fl_t321 = fl_nothing();
    if (fl_t320) {
      bool fl_t322 = false;
      FL_TRY(fl_cond(ctx, sborka, &fl_t322, error));
      fl_value fl_t323 = fl_nothing();
      if (fl_t322) {
        fl_t323 = fl_flag(false);
      } else {
        fl_t323 = fl_flag(true);
      }
      fl_t321 = fl_t323;
    } else {
      fl_t321 = fl_flag(false);
    }
    bool fl_t324 = false;
    FL_TRY(fl_cond(ctx, fl_t321, &fl_t324, error));
    fl_value fl_t325 = fl_nothing();
    if (fl_t324) {
      bool fl_t326 = false;
      FL_TRY(fl_cond(ctx, zagruzka, &fl_t326, error));
      fl_value fl_t327 = fl_nothing();
      if (fl_t326) {
        fl_t327 = fl_flag(false);
      } else {
        fl_t327 = fl_flag(true);
      }
      fl_t325 = fl_t327;
    } else {
      fl_t325 = fl_flag(false);
    }
    bool fl_t328 = false;
    FL_TRY(fl_cond(ctx, fl_t325, &fl_t328, error));
    if (fl_t328) {
      bool fl_t329 = false;
      FL_TRY(fl_cond(ctx, krupnoe, &fl_t329, error));
      if (fl_t329) {
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
  fl_value fl_t330 = fl_nothing();
  FL_TRY(fl_field_get(ctx, nahodka, "доступен", &fl_t330, error));
  bool fl_t331 = false;
  FL_TRY(fl_cond(ctx, fl_t330, &fl_t331, error));
  fl_value fl_t332 = fl_nothing();
  if (fl_t331) {
    fl_t332 = fl_flag(false);
  } else {
    fl_t332 = fl_flag(true);
  }
  bool fl_t333 = false;
  FL_TRY(fl_cond(ctx, fl_t332, &fl_t333, error));
  if (fl_t333) {
    return opis_diska_eto_netrogat(ctx, prigovor, result, error);
  } else {
    fl_value fl_t334 = fl_nothing();
    FL_TRY(opis_diska_ssylka(ctx, nahodka, &fl_t334, error));
    bool fl_t335 = false;
    FL_TRY(fl_cond(ctx, fl_t334, &fl_t335, error));
    if (fl_t335) {
      return opis_diska_eto_netrogat(ctx, prigovor, result, error);
    } else {
      fl_value fl_t336 = fl_nothing();
      FL_TRY(opis_diska_eto_mozhnoubrat(ctx, prigovor, &fl_t336, error));
      bool fl_t337 = false;
      FL_TRY(fl_cond(ctx, fl_t336, &fl_t337, error));
      if (fl_t337) {
        fl_value fl_t338 = fl_nothing();
        FL_TRY(opis_diska_i1_na_pare(ctx, razryad, prigovor, &fl_t338, error));
        bool fl_t339 = false;
        FL_TRY(fl_cond(ctx, fl_t338, &fl_t339, error));
        fl_value fl_t340 = fl_nothing();
        if (fl_t339) {
          fl_value fl_t341 = fl_nothing();
          FL_TRY(opis_diska_katalog(ctx, nahodka, &fl_t341, error));
          bool fl_t342 = false;
          FL_TRY(fl_cond(ctx, fl_t341, &fl_t342, error));
          fl_value fl_t343 = fl_nothing();
          if (fl_t342) {
            fl_t343 = fl_flag(false);
          } else {
            fl_t343 = fl_flag(true);
          }
          fl_t340 = fl_t343;
        } else {
          fl_t340 = fl_flag(false);
        }
        bool fl_t344 = false;
        FL_TRY(fl_cond(ctx, fl_t340, &fl_t344, error));
        if (fl_t344) {
          fl_value fl_t345 = fl_nothing();
          FL_TRY(fl_field_get(ctx, nahodka, "возраст_дней", &fl_t345, error));
          fl_value fl_t346 = fl_nothing();
          FL_TRY(opis_diska_porog_razryada(ctx, razryad, &fl_t346, error));
          fl_value fl_t347 = fl_nothing();
          FL_TRY(fl_gte(ctx, fl_t345, fl_t346, &fl_t347, error));
          *result = fl_t347;
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
  fl_value fl_t348 = fl_nothing();
  FL_TRY(opis_diska_eto_netrogat(ctx, prigovor, &fl_t348, error));
  bool fl_t349 = false;
  FL_TRY(fl_cond(ctx, fl_t348, &fl_t349, error));
  if (fl_t349) {
    *result = fl_flag(fl_equal(ves, fl_number(0.0)));
    return FL_OK;
  } else {
    fl_value fl_t350 = fl_nothing();
    FL_TRY(fl_field_get(ctx, nahodka, "размер", &fl_t350, error));
    bool fl_t351 = false;
    FL_TRY(fl_cond(ctx, fl_flag(fl_equal(ves, fl_t350)), &fl_t351, error));
    if (fl_t351) {
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
