// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package run

import "syscall"

// maxRSSBytes is ru_maxrss in bytes.
//
// The unit is not portable and the difference is a factor of a thousand:
// Darwin counts ru_maxrss in BYTES, Linux and the BSDs in kilobytes.  Getting
// it wrong shows a two-gigabyte build as two megabytes, which looks like a
// small number rather than like a bug — so the conversion lives in a file per
// platform, where it cannot be shared by accident.
func maxRSSBytes(ru *syscall.Rusage) uint64 {
	if ru == nil || ru.Maxrss < 0 {
		return 0
	}
	return uint64(ru.Maxrss)
}
