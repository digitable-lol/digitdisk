// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build flangcore

package main

import (
	"digitdisk/internal/core"
	"digitdisk/internal/coreflang"
)

// chosenDecider returns the flang decision layer from core/out-go.
func chosenDecider() core.Decider { return coreflang.New() }
