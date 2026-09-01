// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package sysinfo

import "net"

// interfaceAddresses maps interface name -> assigned addresses.
//
// This is the same code on every system, and deliberately so: Linux does not
// publish IPv4 addresses as a readable table under /proc, macOS has no /proc
// at all, and the standard library already asks each kernel in its own way.  A
// failure yields an empty map rather than an error, and the address list then
// shows up empty instead of invented.
func interfaceAddresses() map[string][]string {
	out := map[string][]string{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return out
	}
	for _, in := range ifaces {
		addrs, err := in.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			out[in.Name] = append(out[in.Name], a.String())
		}
	}
	return out
}
