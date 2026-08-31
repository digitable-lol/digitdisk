module digitdisk

go 1.24

// Ядро на flang. Оно печатается в core/out-go и коммитится в дерево
// (AGENTS.md: «core/out-go печатается, а не пишется»), поэтому путь есть в
// любом чистом клоне и замена по пути ничего не ломает. Без сборочного признака
// `flangcore` хозяин пользуется заглушкой и говорит об этом в отчёте.
require flangprogram v0.0.0

replace flangprogram => ../core/out-go
