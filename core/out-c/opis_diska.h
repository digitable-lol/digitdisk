/*
 * Сгенерировано flang (бэкенд C, flang/src/emit/c.mjs). Не редактировать руками.
 * Модуль flang: «Опись диска».
 * Файл: объявления: конструкторы значений и функции программы.
 * Правьте исходник на flang и печатайте заново: любая правка здесь потеряется.
 */
#ifndef OPIS_DISKA_H
#define OPIS_DISKA_H

#include "flang_runtime.h"

/*
 * Контракт вызова: функция кладёт результат в *result и возвращает FL_OK
 * либо НЕ трогает *result и возвращает FL_ERROR, заполнив *error (его можно
 * передать NULL). Результат живёт в арене контекста — до ближайшего
 * fl_arena_reset; чтобы сохранить его надолго, скопируйте в свою память.
 *
 *   fl_arena arena;
 *   fl_ctx ctx;
 *   fl_error error;
 *   fl_value result;
 *   fl_arena_init(&arena);
 *   fl_ctx_init(&ctx, &arena);
 *   if (…(&ctx, …, &result, &error) != FL_OK) { … error.code, error.message … }
 *   fl_arena_release(&arena);
 */

/* Запись FTS «Находка»: «путь», «размер», «возраст_дней», «вид», «доступен». */
/* Запись flang тотальна: пропущенное поле — это «ничто», а не дырка. */
fl_status opis_diska_sozdat_nahodka(fl_ctx *ctx, fl_value put, fl_value razmer, fl_value vozrast_dney, fl_value vid, fl_value dostupen, fl_value *out, fl_error *error);

/* Запись FTS «Место»: «разряд», «якорь», «цепь». */
/* Запись flang тотальна: пропущенное поле — это «ничто», а не дырка. */
fl_status opis_diska_sozdat_mesto(fl_ctx *ctx, fl_value razryad, fl_value yakor, fl_value cep, fl_value *out, fl_error *error);

/* Запись FTS «Решение»: «разряд», «приговор», «вес». */
/* Запись flang тотальна: пропущенное поле — это «ничто», а не дырка. */
fl_status opis_diska_sozdat_reshenie(fl_ctx *ctx, fl_value razryad, fl_value prigovor, fl_value ves, fl_value *out, fl_error *error);

/* Запись FTS «Строка отчёта»: «путь», «решение». */
/* Запись flang тотальна: пропущенное поле — это «ничто», а не дырка. */
fl_status opis_diska_sozdat_stroka_otchyota(fl_ctx *ctx, fl_value put, fl_value reshenie, fl_value *out, fl_error *error);

/* Запись FTS «Свод»: «кэш», «журнал», «сборка», «загрузка», «крупное», «неизвестное», «освободить». */
/* Запись flang тотальна: пропущенное поле — это «ничто», а не дырка. */
fl_status opis_diska_sozdat_svod(fl_ctx *ctx, fl_value kesh, fl_value zhurnal, fl_value sborka, fl_value zagruzka, fl_value krupnoe, fl_value neizvestnoe, fl_value osvobodit, fl_value *out, fl_error *error);

/* Сумма типов FTS «Вид»: «Файл» | «Каталог» | «Ссылка». */
/* Дискриминант — имя варианта; проверяется через fl_variant_is(значение, "Имя"). */
fl_status opis_diska_variant_fayl(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_katalog(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_ssylka(fl_ctx *ctx, fl_value *out, fl_error *error);

/* Сумма типов FTS «Разряд»: «Кэш» | «Журнал» | «Сборка» | «Загрузка» | «Крупное» | «Неизвестное». */
/* Дискриминант — имя варианта; проверяется через fl_variant_is(значение, "Имя"). */
fl_status opis_diska_variant_kesh(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_zhurnal(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_sborka(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_zagruzka(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_krupnoe(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_neizvestnoe(fl_ctx *ctx, fl_value *out, fl_error *error);

/* Сумма типов FTS «Приговор»: «МожноУбрать» | «Спросить» | «НеТрогать». */
/* Дискриминант — имя варианта; проверяется через fl_variant_is(значение, "Имя"). */
fl_status opis_diska_variant_mozhnoubrat(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_sprosit(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_netrogat(fl_ctx *ctx, fl_value *out, fl_error *error);

/* Сумма типов FTS «Якорь»: «ОтКорня» | «ГдеУгодно». */
/* Дискриминант — имя варианта; проверяется через fl_variant_is(значение, "Имя"). */
fl_status opis_diska_variant_otkornya(fl_ctx *ctx, fl_value *out, fl_error *error);
fl_status opis_diska_variant_gdeugodno(fl_ctx *ctx, fl_value *out, fl_error *error);

/*
 * Функция flang «Порог крупного».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_krupnogo(fl_ctx *ctx, fl_value *result, fl_error *error);

/*
 * Функция flang «Порог кэша».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_kesha(fl_ctx *ctx, fl_value *result, fl_error *error);

/*
 * Функция flang «Порог журнала».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_zhurnala(fl_ctx *ctx, fl_value *result, fl_error *error);

/*
 * Функция flang «Порог загрузки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: число
 */
fl_status opis_diska_porog_zagruzki(fl_ctx *ctx, fl_value *result, fl_error *error);

/*
 * Функция flang «Составляющие пути».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение: список: строка
 */
fl_status opis_diska_sostavlyayuschie_puti(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Имя в пути».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение: строка
 */
fl_status opis_diska_imya_v_puti(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Есть составляющая».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param imya — «имя»: строка
 * @return значение
 */
fl_status opis_diska_est_sostavlyayuschaya(fl_ctx *ctx, fl_value put, fl_value imya, fl_value *result, fl_error *error);

/*
 * Функция flang «Оканчивается на».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param tekst — «текст»: строка
 * @param hvost — «хвост»: строка
 * @return значение
 */
fl_status opis_diska_okanchivaetsya_na(fl_ctx *ctx, fl_value tekst, fl_value hvost, fl_value *result, fl_error *error);

/*
 * Функция flang «След пути».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение: строка
 */
fl_status opis_diska_sled_puti(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Цепь ограничена».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param cep — «цепь»: строка
 * @return значение
 */
fl_status opis_diska_cep_ogranichena(fl_ctx *ctx, fl_value cep, fl_value *result, fl_error *error);

/*
 * Функция flang «Справочник ограничен».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение
 */
fl_status opis_diska_spravochnik_ogranichen(fl_ctx *ctx, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «Разряд места допустим».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_mesta_dopustim(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Место подходит».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param sled — «след»: строка
 * @param mesto — «место»: «Место»
 * @return значение
 */
fl_status opis_diska_mesto_podhodit(fl_ctx *ctx, fl_value sled, fl_value mesto, fl_value *result, fl_error *error);

/*
 * Функция flang «Это неизвестное».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_eto_neizvestnoe(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Номер разряда».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение: число
 */
fl_status opis_diska_nomer_razryada(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Тот же разряд».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param pervyy — «первый»: «Разряд»
 * @param vtoroy — «второй»: «Разряд»
 * @return значение
 */
fl_status opis_diska_tot_zhe_razryad(fl_ctx *ctx, fl_value pervyy, fl_value vtoroy, fl_value *result, fl_error *error);

/*
 * Функция flang «Место обосновано».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param spravochnik — «справочник»: список: «Место»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_mesto_obosnovano(fl_ctx *ctx, fl_value put, fl_value spravochnik, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Разряд по справочнику».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Разряд»
 */
fl_status opis_diska_razryad_po_spravochniku(fl_ctx *ctx, fl_value put, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «Шестнадцатеричный знак».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param znak — «знак»: строка
 * @return значение
 */
fl_status opis_diska_shestnadcaterichnyy_znak(fl_ctx *ctx, fl_value znak, fl_value *result, fl_error *error);

/*
 * Функция flang «Похоже на отпечаток».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param chast — «часть»: строка
 * @return значение
 */
fl_status opis_diska_pohozhe_na_otpechatok(fl_ctx *ctx, fl_value chast, fl_value *result, fl_error *error);

/*
 * Функция flang «Адресуется содержимым».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_adresuetsya_soderzhimym(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Под системным временным».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_pod_sistemnym_vremennym(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Примета кэша».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_kesha(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Примета журнала».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_zhurnala(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Примета сборки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_sborki(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Примета загрузки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @return значение
 */
fl_status opis_diska_primeta_zagruzki(fl_ctx *ctx, fl_value put, fl_value *result, fl_error *error);

/*
 * Функция flang «Есть примета».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param put — «путь»: строка
 * @param mesto — «место»: «Разряд»
 * @return значение
 */
fl_status opis_diska_est_primeta(fl_ctx *ctx, fl_value put, fl_value mesto, fl_value *result, fl_error *error);

/*
 * Функция flang «Разряд решён размером».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_reshyon_razmerom(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Разряд находки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param mesto — «место»: «Разряд»
 * @return значение: «Разряд»
 */
fl_status opis_diska_razryad_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value mesto, fl_value *result, fl_error *error);

/*
 * Функция flang «Крупное не мельче порога».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_krupnoe_ne_melche_poroga(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Каталог».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @return значение
 */
fl_status opis_diska_katalog(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error);

/*
 * Функция flang «Ссылка».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @return значение
 */
fl_status opis_diska_ssylka(fl_ctx *ctx, fl_value nahodka, fl_value *result, fl_error *error);

/*
 * Функция flang «Приговор мусора».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param porog — «порог»: число
 * @return значение: «Приговор»
 */
fl_status opis_diska_prigovor_musora(fl_ctx *ctx, fl_value nahodka, fl_value porog, fl_value *result, fl_error *error);

/*
 * Функция flang «Приговор находки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @return значение: «Приговор»
 */
fl_status opis_diska_prigovor_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Вес находки».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param prigovor — «приговор»: «Приговор»
 * @return значение: число
 */
fl_status opis_diska_ves_nahodki(fl_ctx *ctx, fl_value nahodka, fl_value prigovor, fl_value *result, fl_error *error);

/*
 * Функция flang «Вес в границах».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param ves — «вес»: число
 * @return значение
 */
fl_status opis_diska_ves_v_granicah(fl_ctx *ctx, fl_value nahodka, fl_value ves, fl_value *result, fl_error *error);

/*
 * Функция flang «Решить находку».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Решение»
 */
fl_status opis_diska_reshit_nahodku(fl_ctx *ctx, fl_value nahodka, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «Решить всё».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: список: «Решение»
 */
fl_status opis_diska_reshit_vsyo(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «И1 держится».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param reshenie — «решение»: «Решение»
 * @return значение
 */
fl_status opis_diska_i1_derzhitsya(fl_ctx *ctx, fl_value reshenie, fl_value *result, fl_error *error);

/*
 * Функция flang «И1 держится всюду».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param resheniya — «решения»: список: «Решение»
 * @return значение
 */
fl_status opis_diska_i1_derzhitsya_vsyudu(fl_ctx *ctx, fl_value resheniya, fl_value *result, fl_error *error);

/*
 * Функция flang «Пустой свод».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @return значение: «Свод»
 */
fl_status opis_diska_pustoy_svod(fl_ctx *ctx, fl_value *result, fl_error *error);

/*
 * Функция flang «Прибавить решение».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param svod — «свод»: «Свод»
 * @param reshenie — «решение»: «Решение»
 * @return значение: «Свод»
 */
fl_status opis_diska_pribavit_reshenie(fl_ctx *ctx, fl_value svod, fl_value reshenie, fl_value *result, fl_error *error);

/*
 * Функция flang «Свести».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Свод»
 */
fl_status opis_diska_svesti(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «Сумма размеров убираемых».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: число
 */
fl_status opis_diska_summa_razmerov_ubiraemyh(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «И2 держится».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @param svod — «свод»: «Свод»
 * @return значение
 */
fl_status opis_diska_i2_derzhitsya(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value svod, fl_value *result, fl_error *error);

/*
 * Функция flang «Строку отчёта».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: «Строка отчёта»
 */
fl_status opis_diska_stroku_otchyota(fl_ctx *ctx, fl_value nahodka, fl_value spravochnik, fl_value *result, fl_error *error);

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
fl_status opis_diska_vstavit_po_vesu(fl_ctx *ctx, fl_value stroka, fl_value stroki, fl_value *result, fl_error *error);

/*
 * Функция flang «Приписать строку отчёта».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param pervaya — «первая»: «Строка отчёта»
 * @param stroki — «строки»: список: «Строка отчёта»
 * @return значение: список: «Строка отчёта»
 */
fl_status opis_diska_pripisat_stroku_otchyota(fl_ctx *ctx, fl_value pervaya, fl_value stroki, fl_value *result, fl_error *error);

/*
 * Функция flang «Отчёт».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param spravochnik — «справочник»: список: «Место»
 * @return значение: список: «Строка отчёта»
 */
fl_status opis_diska_otchyot(fl_ctx *ctx, fl_value zapisi, fl_value spravochnik, fl_value *result, fl_error *error);

/*
 * Функция flang «Отчёт той же длины».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param zapisi — «записи»: список: «Находка»
 * @param stroki — «строки»: список: «Строка отчёта»
 * @return значение
 */
fl_status opis_diska_otchyot_toy_zhe_dliny(fl_ctx *ctx, fl_value zapisi, fl_value stroki, fl_value *result, fl_error *error);

/*
 * Функция flang «Это МожноУбрать».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_eto_mozhnoubrat(fl_ctx *ctx, fl_value prigovor, fl_value *result, fl_error *error);

/*
 * Функция flang «Это НеТрогать».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_eto_netrogat(fl_ctx *ctx, fl_value prigovor, fl_value *result, fl_error *error);

/*
 * Функция flang «И1 на паре».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_i1_na_pare(fl_ctx *ctx, fl_value razryad, fl_value prigovor, fl_value *result, fl_error *error);

/*
 * Функция flang «Порог разряда».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param razryad — «разряд»: «Разряд»
 * @return значение: число
 */
fl_status opis_diska_porog_razryada(fl_ctx *ctx, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Разряд обоснован».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param mesto — «место»: «Разряд»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value mesto, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Разряд обоснован приметой».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @return значение
 */
fl_status opis_diska_razryad_obosnovan_primetoy(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value *result, fl_error *error);

/*
 * Функция flang «Приговор обоснован».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param razryad — «разряд»: «Разряд»
 * @param prigovor — «приговор»: «Приговор»
 * @return значение
 */
fl_status opis_diska_prigovor_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value razryad, fl_value prigovor, fl_value *result, fl_error *error);

/*
 * Функция flang «Вес обоснован».
 *
 * Тотальная: завершение доказано анализом завершаемости (totality.mjs).
 * @param nahodka — «находка»: «Находка»
 * @param prigovor — «приговор»: «Приговор»
 * @param ves — «вес»: число
 * @return значение
 */
fl_status opis_diska_ves_obosnovan(fl_ctx *ctx, fl_value nahodka, fl_value prigovor, fl_value ves, fl_value *result, fl_error *error);

/*
 * Вызов функции по её исходному имени flang. Нужен прогонщику и всякому,
 * кто связывает программу с внешним миром динамически (скрипт, FFI, тест).
 */
fl_status opis_diska_call(fl_ctx *ctx, const char *name, const fl_value *args, size_t count,
                    fl_value *result, fl_error *error);

/*
 * ТО ЖЕ, НО ЧЕРЕЗ ГРАНИЦУ ВХОДА: объявленные типы сверяются ДО вызова.
 * Значения, пришедшие снаружи — из JSON, из другого языка, от человека, —
 * обязаны заходить ЗДЕСЬ: `_call` их не сверяет, а доказательство завершения
 * `тотальной` стоит НА ТИПЕ и вместе с типом теряется.
 */
fl_status opis_diska_enter(fl_ctx *ctx, const char *name, const fl_value *args, size_t count,
                    fl_value *result, fl_error *error);

/*
 * Объявленные типы параметров — данными. Прогонщик сверяет по ним значения,
 * пришедшие снаружи, ДО вызова: доказательство завершения `тотальной` стоит
 * на типе, и значение вне типа выносит вместе с типом и доказательство.
 */
const fl_entry_table *opis_diska_entry(void);

#endif /* OPIS_DISKA_H */
