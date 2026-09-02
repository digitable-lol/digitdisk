module digitdisk

go 1.25

// Ядро на flang. Оно печатается в core/out-go и коммитится в дерево
// (AGENTS.md: «core/out-go печатается, а не пишется»), поэтому путь есть в
// любом чистом клоне и замена по пути ничего не ломает. Без сборочного признака
// `flangcore` хозяин пользуется заглушкой и говорит об этом в отчёте.
require flangprogram v0.0.0

replace flangprogram => ../core/out-go

// Раскладка экрана на flang. Библиотека flang-tui подключена подмодулем в
// ui-flang/flang-tui и напечатана в ui-flang/out-go; печать коммитится, как и
// core/out-go, поэтому путь есть в любом чистом клоне. Без сборочного признака
// `flangui` хозяин считает раскладку рукописным Go (host/internal/ui/
// layout_stock.go) и эти модули не собираются вовсе.
require (
	flangcolour v0.0.0
	flangformat v0.0.0
	flanghistory v0.0.0
	flangscreen v0.0.0
	flangscroll v0.0.0
	flangtabs v0.0.0
)

replace flangcolour => ../ui-flang/out-go/colour

replace flangformat => ../ui-flang/out-go/format

replace flanghistory => ../ui-flang/out-go/history

replace flangscreen => ../ui-flang/out-go/screen

replace flangscroll => ../ui-flang/out-go/scroll

replace flangtabs => ../ui-flang/out-go/tabs
