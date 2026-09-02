// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package procfs contains pure parsers for the text formats exposed by the
// Linux /proc and /sys pseudo-filesystems.  Every exported function here takes
// the file content as a string (or []byte) and returns parsed values, so the
// parsers can be exercised in tests with captured samples instead of the live
// kernel.  That is also why the package carries no build tag: a parser of
// captured text is worth testing on any machine, and only the collector that
// names /proc paths is behind //go:build linux.
//
// The types here — Memory, LoadAvg, NetCounters — are the shape of the
// snapshot rather than a shape of /proc: the macOS collector fills the same
// structs from sysctl, and package darwinsys decodes the binary structures it
// gets back the way this package decodes text.
package procfs

import (
	"strconv"
	"strings"
)

// Memory holds the subset of /proc/meminfo digitdisk reports, in bytes.
// Zero values mean "the source did not expose this key"; callers distinguish
// a measured zero from an absent measurement via the Present map (Has).
//
// The shape is shared by every host, not just the Linux one: macOS fills the
// same struct out of sysctl, where only some of the fields exist at all.  That
// is exactly why Present is part of it — a field nobody measured must not
// arrive as a zero that reads like a measurement.
type Memory struct {
	Total     uint64 `json:"total_bytes"`
	Free      uint64 `json:"free_bytes"`
	Available uint64 `json:"available_bytes"`
	Buffers   uint64 `json:"buffers_bytes"`
	Cached    uint64 `json:"cached_bytes"`
	// BuffCache is Buffers + Cached + SReclaimable, i.e. the "buff/cache"
	// column reported by procps-ng free(1).
	BuffCache uint64 `json:"buff_cache_bytes"`
	Shared    uint64 `json:"shared_bytes"`
	// Used is Total - Available, which is what procps-ng 4.x free(1) prints
	// in its "used" column.
	Used      uint64 `json:"used_bytes"`
	SwapTotal uint64 `json:"swap_total_bytes"`
	SwapFree  uint64 `json:"swap_free_bytes"`
	SwapUsed  uint64 `json:"swap_used_bytes"`

	// Raw keeps every key seen, in bytes, so nothing is silently invented.
	Raw map[string]uint64 `json:"-"`

	// Present says which of the fields above were actually measured, keyed
	// by the JSON name of the field.  A nil map means "everything here was
	// measured" — see Has.
	Present map[string]bool `json:"present,omitempty"`
}

// Field names used as keys of Present.  They are the JSON names of the fields,
// so a reader of the JSON output looks a flag up by the name it just read.
const (
	FieldTotal     = "total_bytes"
	FieldFree      = "free_bytes"
	FieldAvailable = "available_bytes"
	FieldBuffers   = "buffers_bytes"
	FieldCached    = "cached_bytes"
	FieldBuffCache = "buff_cache_bytes"
	FieldShared    = "shared_bytes"
	FieldUsed      = "used_bytes"
	FieldSwapTotal = "swap_total_bytes"
	FieldSwapFree  = "swap_free_bytes"
	FieldSwapUsed  = "swap_used_bytes"
)

// Keys the macOS collector adds to Raw, in bytes.  They have no counterpart in
// /proc/meminfo and no field of their own, because only one system publishes
// them — but wired and compressed memory are the two numbers a Mac owner reads
// first, so the report prints them when they are there.
const (
	RawWired       = "wired_bytes"
	RawCompressed  = "compressed_bytes"
	RawActive      = "active_bytes"
	RawInactive    = "inactive_bytes"
	RawSpeculative = "speculative_bytes"
	RawPurgeable   = "purgeable_bytes"
	RawAnonymous   = "anonymous_bytes"
)

// Has reports whether the named field carries a measurement.
//
// A nil Present map answers true: a Memory built by hand (a test, a caller
// that fills every field it has) is taken at face value.  Only a source that
// knows it left gaps — the macOS collector, or a /proc/meminfo without
// MemAvailable — fills Present, and then a false answer means "not measured",
// never "measured zero".
func (m Memory) Has(field string) bool {
	if m.Present == nil {
		return true
	}
	return m.Present[field]
}

// ParseMeminfo parses the "Key:  value kB" format of /proc/meminfo.
// Values are converted to bytes.  Unknown units are rejected (skipped) rather
// than guessed at.
func ParseMeminfo(text string) Memory {
	raw := make(map[string]uint64, 64)
	for _, line := range strings.Split(text, "\n") {
		key, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}
		switch {
		case len(fields) == 1:
			// unit-less (e.g. HugePages_Total): keep as-is
		case strings.EqualFold(fields[1], "kB"):
			n *= 1024
		default:
			continue
		}
		raw[strings.TrimSpace(key)] = n
	}

	m := Memory{Raw: raw, Present: make(map[string]bool, 11)}
	// seen reports whether every named key was in the text; anySeen,
	// whether any of them was.
	seen := func(keys ...string) bool {
		for _, k := range keys {
			if _, ok := raw[k]; !ok {
				return false
			}
		}
		return true
	}
	anySeen := func(keys ...string) bool {
		for _, k := range keys {
			if _, ok := raw[k]; ok {
				return true
			}
		}
		return false
	}
	m.Total = raw["MemTotal"]
	m.Free = raw["MemFree"]
	m.Available = raw["MemAvailable"]
	m.Buffers = raw["Buffers"]
	m.Cached = raw["Cached"]
	m.Shared = raw["Shmem"]
	m.BuffCache = raw["Buffers"] + raw["Cached"] + raw["SReclaimable"]
	// The derived numbers are computed only from keys that were actually
	// there.  Without MemAvailable, "total minus available" would come out
	// as the whole of memory and look like a machine with nothing free.
	if seen("MemTotal", "MemAvailable") && m.Total >= m.Available {
		m.Used = m.Total - m.Available
	}
	m.SwapTotal = raw["SwapTotal"]
	m.SwapFree = raw["SwapFree"]
	if seen("SwapTotal", "SwapFree") && m.SwapTotal >= m.SwapFree {
		m.SwapUsed = m.SwapTotal - m.SwapFree
	}
	m.Present[FieldTotal] = seen("MemTotal")
	m.Present[FieldFree] = seen("MemFree")
	m.Present[FieldAvailable] = seen("MemAvailable")
	m.Present[FieldBuffers] = seen("Buffers")
	m.Present[FieldCached] = seen("Cached")
	m.Present[FieldBuffCache] = anySeen("Buffers", "Cached", "SReclaimable")
	m.Present[FieldShared] = seen("Shmem")
	m.Present[FieldUsed] = seen("MemTotal", "MemAvailable")
	m.Present[FieldSwapTotal] = seen("SwapTotal")
	m.Present[FieldSwapFree] = seen("SwapFree")
	m.Present[FieldSwapUsed] = seen("SwapTotal", "SwapFree")
	return m
}
