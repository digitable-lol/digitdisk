// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build flangcore

package main

import (
	"digitdisk/internal/core"
	"digitdisk/internal/coreflang"
	"digitdisk/internal/lang"
)

// chosenDecider returns the flang decision layer from core/out-go.
func chosenDecider(lang.Lang) core.Decider { return coreflang.New() }
