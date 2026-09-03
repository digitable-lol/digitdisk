// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package iokit

import (
	"strings"

	"digitdisk/internal/gpuinfo"
)

// Source is what a card read here says about where its numbers came from.  It
// is the registry and nothing else: no outside program is run, and nothing is
// inferred from the model of the machine.
const Source = "IORegistry"

// accelClasses are the classes a Mac hangs graphics off, in the order they are
// tried.  IOAccelerator is the superclass every graphics driver registers
// under, on both architectures; IOPCIDevice is the fallback for a machine
// whose accelerator did not register — a card that is on the bus but has no
// driver is still a card, and saying its numbers is better than saying
// nothing.
var accelClasses = []string{"IOAccelerator", "IOPCIDevice"}

// Cards reads the video cards out of the IORegistry.
//
// Every field is filled only from a property that is actually there.  An
// integrated card publishes no memory of its own and gets no memory here; a
// card whose entry carries no model name is named by the class of its driver
// rather than by a guess.  That is the whole point of reading the registry
// instead of assuming: a Mac that says nothing about a field leaves it empty.
func Cards() []gpuinfo.Card {
	var out []gpuinfo.Card
	seen := map[string]bool{}
	for _, class := range accelClasses {
		entries, err := Match(class, 3)
		if err != nil {
			continue
		}
		for i := range entries {
			e := &entries[i]
			if class == "IOPCIDevice" && !isDisplay(e) {
				continue
			}
			c := card(e)
			key := c.VendorID + ":" + c.DeviceID + ":" + c.Name
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, c)
		}
	}
	return out
}

// isDisplay reads the PCI class code and answers whether this device is a
// display controller.  The code is four bytes little-endian, and its third
// byte is the base class; 0x03 is "display controller" in the PCI
// specification, which is also what /sys/bus/pci publishes on Linux.
func isDisplay(e *Entry) bool {
	v, ok := e.Props["class-code"]
	if !ok || v.Kind != KindData || len(v.Data) < 3 {
		return false
	}
	return v.Data[2] == 0x03
}

// card builds one card out of one registry entry and the nodes it hangs off.
func card(e *Entry) gpuinfo.Card {
	c := gpuinfo.Card{
		Node:   e.Name,
		Name:   text(e, "model"),
		Driver: text(e, "IOClass", "CFBundleIdentifier"),
		Source: Source,
	}
	if c.Name == "" {
		// No model name anywhere up the chain.  The name a person would
		// recognise is then the driver's own class — "AGXAcceleratorG16"
		// says which Apple GPU this is — and failing that the node name.
		c.Name = c.Driver
	}
	if c.Name == "" {
		c.Name = e.Name
	}
	c.VendorID = pciID(e, "vendor-id")
	c.DeviceID = pciID(e, "device-id")
	c.Vendor = vendorName(c.VendorID)
	if n, ok := vram(e); ok {
		c.MemoryTotalBytes = &n
	}
	return c
}

// text returns the first of several keys that is present and readable as text,
// looking up the chain of parents.
func text(e *Entry, keys ...string) string {
	for _, k := range keys {
		v, ok := e.Prop(k)
		if !ok {
			continue
		}
		switch v.Kind {
		case KindString:
			if s := strings.TrimSpace(v.Str); s != "" {
				return s
			}
		case KindData:
			if s, ok := asciiOf(v.Data); ok {
				if s = strings.TrimSpace(s); s != "" {
					return s
				}
			}
		}
	}
	return ""
}

// pciID renders vendor-id or device-id.  IOKit publishes both as four bytes,
// little-endian, of which the low two are the identifier — the same number
// /sys/bus/pci prints as 0x8086 on Linux, and it is printed the same way here
// so the two collectors do not describe one card two ways.
func pciID(e *Entry, key string) string {
	v, ok := e.Prop(key)
	if !ok {
		return ""
	}
	var n uint64
	switch {
	case v.Kind == KindData && len(v.Data) >= 2:
		n = uint64(v.Data[0]) | uint64(v.Data[1])<<8
	case v.Kind == KindNumber && v.Num > 0:
		n = uint64(v.Num) & 0xffff
	default:
		return ""
	}
	const hexDigits = "0123456789abcdef"
	b := []byte("0x0000")
	for i := 0; i < 4; i++ {
		b[5-i] = hexDigits[(n>>(4*i))&0xf]
	}
	return string(b)
}

// vendorName names the four vendors that ever shipped graphics in a Mac.  The
// Linux collector reads pci.ids for this; macOS ships no such file, and a
// table of four is a table, not a database.
func vendorName(id string) string {
	switch id {
	case "0x8086":
		return "Intel"
	case "0x1002", "0x1022":
		return "AMD"
	case "0x10de":
		return "NVIDIA"
	case "0x106b":
		return "Apple"
	}
	return ""
}

// vram reads how much memory the card has.  IOKit publishes it two ways and
// the drivers do not agree on which: "VRAM,totalMB" in megabytes, and
// "VRAM,totalsize" in bytes.  Either may be a number or four bytes; both are
// read, and neither is invented.
func vram(e *Entry) (uint64, bool) {
	if n, ok := amount(e, "VRAM,totalsize"); ok {
		return n, true
	}
	if n, ok := amount(e, "VRAM,totalMB"); ok {
		return n << 20, true
	}
	return 0, false
}

func amount(e *Entry, key string) (uint64, bool) {
	v, ok := e.Prop(key)
	if !ok {
		return 0, false
	}
	switch v.Kind {
	case KindNumber:
		if v.Num > 0 {
			return uint64(v.Num), true
		}
	case KindData:
		var n uint64
		for i := len(v.Data) - 1; i >= 0 && i < 8; i-- {
			n = n<<8 | uint64(v.Data[i])
		}
		if n > 0 {
			return n, true
		}
	}
	return 0, false
}
