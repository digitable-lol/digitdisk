// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package iokit

import "digitdisk/internal/gpuinfo"

// Cards reads the video cards out of the IORegistry.
//
// Every field is filled only from a property that is actually there.  A card
// that publishes no memory of its own gets no memory here, and a card whose
// entry carries no model name is named by the class of its driver rather than
// by a guess: that is the whole point of reading the registry instead of
// assuming what a Mac of this model usually has.
func Cards() []gpuinfo.Card {
	var out []gpuinfo.Card
	seen := map[string]bool{}
	for _, class := range AccelClasses {
		entries, err := Match(class, 3)
		if err != nil {
			continue
		}
		out = cardsFrom(class, entries, seen, out)
	}
	return out
}
