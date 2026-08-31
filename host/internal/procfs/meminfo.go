// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package procfs contains pure parsers for the text formats exposed by the
// Linux /proc and /sys pseudo-filesystems.  Every exported function here takes
// the file content as a string (or []byte) and returns parsed values, so the
// parsers can be exercised in tests with captured samples instead of the live
// kernel.
package procfs

import (
	"strconv"
	"strings"
)

// Memory holds the subset of /proc/meminfo digitdisk reports, in bytes.
// Zero values mean "the kernel did not expose this key"; callers distinguish
// that via the Present map.
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

	m := Memory{Raw: raw}
	m.Total = raw["MemTotal"]
	m.Free = raw["MemFree"]
	m.Available = raw["MemAvailable"]
	m.Buffers = raw["Buffers"]
	m.Cached = raw["Cached"]
	m.Shared = raw["Shmem"]
	m.BuffCache = raw["Buffers"] + raw["Cached"] + raw["SReclaimable"]
	if m.Total >= m.Available {
		m.Used = m.Total - m.Available
	}
	m.SwapTotal = raw["SwapTotal"]
	m.SwapFree = raw["SwapFree"]
	if m.SwapTotal >= m.SwapFree {
		m.SwapUsed = m.SwapTotal - m.SwapFree
	}
	return m
}
