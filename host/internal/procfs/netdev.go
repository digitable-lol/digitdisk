// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package procfs

import (
	"strconv"
	"strings"
)

// NetCounters holds the receive/transmit counters of one interface as found in
// /proc/net/dev.
type NetCounters struct {
	Name      string `json:"name"`
	RxBytes   uint64 `json:"rx_bytes"`
	RxPackets uint64 `json:"rx_packets"`
	RxErrors  uint64 `json:"rx_errors"`
	RxDropped uint64 `json:"rx_dropped"`
	TxBytes   uint64 `json:"tx_bytes"`
	TxPackets uint64 `json:"tx_packets"`
	TxErrors  uint64 `json:"tx_errors"`
	TxDropped uint64 `json:"tx_dropped"`
}

// ParseNetDev parses /proc/net/dev.  The two header lines are skipped; the
// interface name is the text before the first ':' on the line, which may be
// glued to the first counter when the value is wide.
func ParseNetDev(text string) []NetCounters {
	var out []NetCounters
	for _, line := range strings.Split(text, "\n") {
		name, rest, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		if name == "" || strings.Contains(name, " ") {
			continue // header rows ("Inter-|   Receive ...")
		}
		f := strings.Fields(rest)
		if len(f) < 16 {
			continue
		}
		num := func(i int) uint64 {
			n, _ := strconv.ParseUint(f[i], 10, 64)
			return n
		}
		out = append(out, NetCounters{
			Name:      name,
			RxBytes:   num(0),
			RxPackets: num(1),
			RxErrors:  num(2),
			RxDropped: num(3),
			TxBytes:   num(8),
			TxPackets: num(9),
			TxErrors:  num(10),
			TxDropped: num(11),
		})
	}
	return out
}
