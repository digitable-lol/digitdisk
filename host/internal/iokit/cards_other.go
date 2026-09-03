// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !darwin

package iokit

import "digitdisk/internal/gpuinfo"

// Source names where a card read by this package came from.
const Source = "IORegistry"

// Cards has no registry to read outside macOS.
func Cards() []gpuinfo.Card { return nil }
