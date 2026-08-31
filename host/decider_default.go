// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build !flangcore

package main

import "digitdisk/internal/core"

// chosenDecider returns the placeholder decision layer.  See coreflang/bridge.go
// for how to build against the real one.
func chosenDecider() core.Decider { return core.Default() }
