// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

// Английские половины разделов о железе: значок системы, ядра по отдельности,
// видеокарты.
//
// Ширина здесь — часть перевода, как и во всём выводе. Подписи СИСТЕМЫ живут в
// поле на 14 ячеек, подписи датчиков — на 12, строки печатного отчёта стоят в
// колонке шириной 16, а заголовки таблицы видеокарт разложены по %8s, %22s и
// двум %9s. Английская половина, которая длиннее русской, не переводит строку,
// а сдвигает колонку, поэтому она такой же длины или короче.
//
// Чего здесь нет: имена неизмеренного про железо («загрузка по ядрам»,
// «видеокарты», «оболочка», «рабочий стол» и прочие константы
// internal/sysinfo/hardware.go) приходят значениями и заведены вокабулой рядом
// с договором, в dict_core.go. Часть из них служит заодно подписью на экране —
// статья одна, и спрашивают её обоими способами.
func init() {
	add(map[string]string{
		// ── СИСТЕМА: подписи рядом со значком (поле на 14) ───────────
		"модель":         "model",
		"процессор":      "processor",
		"терминал":       "terminal",
		", %d-разрядная": ", %d-bit",
		"%s из %s (%s)":  "%s of %s (%s)",

		// ── ЗАГРУЗКА: ядра по отдельности ────────────────────────────
		"по ядрам":      "by core",
		"разброс":       "spread",
		"под нагрузкой": "loaded",
		"ядро %d":       "core %d",
		"замерено %d из %s  (окно %d мс)":            "measured %d of %s  (window %d ms)",
		"мин %s · медиана %s · макс %s (ядро %d)":    "min %s · median %s · max %s (core %d)",
		"%d ядер из %d заняты больше чем наполовину": "%d cores of %d are busy more than half the time",
		"ПО ЯДРАМ":      "BY CORE",
		"КАРТА ЯДЕР":    "CORE MAP",
		"САМЫЕ ЗАНЯТЫЕ": "BUSIEST",
		"ячейка — ядро: ▁ пусто, █ занято целиком": "a cell is a core: ▁ idle, █ fully busy",
		"  по %d ядер в ячейке":                    "  %d cores per cell",
		"почему — digitdisk status --why":          "why — digitdisk status --why",

		// Те же ядра в печатном отчёте: колонка шириной 16.
		"  процессор     %s":      "  processor     %s",
		"  процессор     %s × %d": "  processor     %s × %d",
		"  по ядрам      мин %s / медиана %s / макс %s (ядро %d); занято больше половины %d из %d": "  by core       min %s / median %s / max %s (core %d); more than half busy %d of %d",

		// ── ВИДЕОКАРТЫ ───────────────────────────────────────────────
		"ВИДЕОКАРТЫ":                    "VIDEO CARDS",
		"— (видеокарт не нашлось)":      "— (no video cards found)",
		"…и ещё %d — раздел ВИДЕОКАРТЫ": "…and %d more — the VIDEO CARDS section",
		"температура":                   "temperature",
		"питание":                       "power",
		"всего %s, занято %s":           "total %s, used %s",
		"%s Вт":                         "%s W",
		"%s МГц":                        "%s MHz",

		// Заголовки таблицы карт в печатном отчёте: %8s, %22s, %9s, %9s.
		"карта":    "card",
		"темп.":    "temp.",
		"мощность": "power",

		// Откуда взялись числа карты — строкой под ними.
		"шина %s":     "bus %s",
		"драйвер %s":  "driver %s",
		"числа из %s": "numbers from %s",
		"числа от чужой программы %s": "numbers from an outside program, %s",

		// ── Ключ --gpu-tool: справка и пояснение самого flag ─────────
		"спросить о видеокартах программу их драйвера (nvidia-smi)":                "ask the driver's own program (nvidia-smi) about the video cards",
		"  --gpu-tool        status, run: спросить о видеокартах чужую программу":  "  --gpu-tool        status, run: ask an outside program about the cards",
		"                    (nvidia-smi) — то, чего драйвер не публикует файлами": "                    (nvidia-smi) — what the driver publishes in no file",

		// ── Причины неизмеренного: печатает только --why ─────────────
		"ядро не публикует счётчики по каждому процессору":                                                    "the kernel publishes no per-processor counters",
		"за окно замера счётчики процессоров не сдвинулись":                                                   "the processor counters did not move over the sample window",
		"среднее по ядрам разошлось с общей загрузкой машины — список ядер не публикуем":                      "the per-core average disagreed with the load of the whole machine — the core list is not published",
		"окно замера нулевое — доля занятого времени каждого ядра не измерялась":                              "the sample window is zero — the busy share of each core was not measured",
		"переменная окружения SHELL пуста — оболочку назвать нечем":                                           "the SHELL environment variable is empty — there is nothing to name the shell by",
		"переменные окружения XDG_CURRENT_DESKTOP и DESKTOP_SESSION пусты — рабочего стола в этом сеансе нет": "the XDG_CURRENT_DESKTOP and DESKTOP_SESSION environment variables are empty — this session has no desktop",
		"прошивка не публикует имя машины в %s":                                                               "the firmware publishes no machine name in %s",
		"в %s нет строки с названием процессора":                                                              "%s has no line naming the processor",

		// То же на macOS.
		"macOS публикует сведения о видеокартах только через IOKit, объектами Core Foundation; без cgo мы их не читаем, а угадывать не станем": "macOS publishes what it knows about video cards through IOKit alone, as Core Foundation objects; without cgo we do not read them, and we will not guess",
		"узел machdep.cpu.brand_string ничего не ответил — процессор себя не назвал":                                                           "the machdep.cpu.brand_string node answered nothing — the processor did not name itself",
		"ядро не дало счётчики по каждому процессору отдельно":                                                                                 "the kernel gave no counters for each processor separately",
		"процессоров в ответе ядра не столько, сколько машина насчитала у себя, — по ядрам не публикуем":                                       "the kernel answered about a different number of processors than the machine counts for itself — per-core figures are not published",

		// Почему у найденных карт нет чисел.
		"в %s и на шине PCI видеокарт не нашлось":                                                "no video cards were found in %s or on the PCI bus",
		"у найденной карты нет драйвера, а без него ядро не публикует о ней ничего, кроме имени": "the card that was found has no driver, and without one the kernel publishes nothing about it but its name",
		"драйвер %s не публикует ни загрузки, ни памяти карты в файлах ядра":                     "the %s driver publishes neither the load nor the memory of the card in kernel files",
		"драйвер nvidia не публикует ни загрузки, ни памяти, ни температуры карты в файлах — их знает только его собственная программа nvidia-smi, и запускается она лишь по ключу --gpu-tool":    "the nvidia driver publishes neither the load, nor the memory, nor the temperature of the card in files — only its own program nvidia-smi knows them, and that is run by --gpu-tool alone",
		"драйвер отдаёт мощность карты не в микроваттах, как обещает документация hwmon: полученное число меньше полуватта, то есть единица измерения у него другая — такое число мы не печатаем": "the driver reports the power of the card in something other than the microwatts the hwmon documentation promises: the number that came back is under half a watt, so its unit is a different one — such a number we do not print",
	})
}
