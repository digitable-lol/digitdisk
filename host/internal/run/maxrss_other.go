// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !darwin

package run

import "syscall"

// maxRSSBytes is ru_maxrss in bytes.  Everywhere but Darwin the kernel counts
// it in kilobytes; see maxrss_darwin.go for why this is not one function.
func maxRSSBytes(ru *syscall.Rusage) uint64 {
	if ru == nil || ru.Maxrss < 0 {
		return 0
	}
	return uint64(ru.Maxrss) * 1024
}
