// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !darwin

package iokit

// Match has nothing to read where there is no IORegistry.  The package still
// compiles and its types still exist, so a caller can name them without a
// build tag of its own.
func Match(class string, parents int) ([]Entry, error) { return nil, notMac{} }
