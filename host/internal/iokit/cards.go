// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package iokit

import (
	"strconv"
	"strings"

	"digitdisk/internal/gpuinfo"
)

// Turning registry entries into cards.  Nothing here calls a framework: it
// works on the Entry values the darwin half reads, which is what lets it be
// tested on any machine — the same split internal/darwinsys makes between the
// call and the decoding of what came back.

// Source is what a card read here says about where its numbers came from.  It
// is the registry and nothing else: no outside program is run, and nothing is
// inferred from the model of the machine.
const Source = "IORegistry"

// AccelClasses are the classes a Mac hangs graphics off, in the order they are
// tried.  IOAccelerator is the superclass every graphics driver registers
// under, on both architectures; IOPCIDevice is the fallback for a card that is
// on the bus with no accelerator of its own, and its entries are filtered by
// the PCI class code so the other eleven devices of a Mac are not offered as
// video cards.
var AccelClasses = []string{"IOAccelerator", "IOPCIDevice"}

// cardsFrom builds the card list out of already-read entries, dropping the
// duplicate a card causes by being both an accelerator and a PCI device.
func cardsFrom(class string, entries []Entry, seen map[string]bool, out []gpuinfo.Card) []gpuinfo.Card {
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
//
// Which field is filled from where is decided by what the two Macs of the
// checks actually publish, and nothing is filled from anywhere else.  The name
// and the identifiers live on the PCI node, one or two steps above the
// accelerator; the driver lives on the accelerator itself, and asking for it
// up the chain would answer with the driver of the bus.
func card(e *Entry) gpuinfo.Card {
	c := gpuinfo.Card{
		Node:   e.Name,
		Name:   text(e, true, "model"),
		Driver: driver(e),
		Bus:    busAddress(e),
		Source: Source,
	}
	if c.Name == "" {
		// No model name anywhere up the chain — that is an Apple Silicon
		// Mac, where the GPU is not on a bus and has no PCI name.  The
		// name a person would recognise is then the driver's own class.
		c.Name = c.Driver
	}
	if c.Name == "" {
		c.Name = e.Name
	}
	c.VendorID = pciID(e, "vendor-id")
	c.DeviceID = pciID(e, "device-id")
	c.Vendor = gpuinfo.VendorName(c.VendorID)
	if n, ok := vram(e); ok {
		c.MemoryTotalBytes = &n
	}
	return c
}

// driver names what runs the card.  It is read from the entry itself and never
// from a parent: the node above an accelerator is the bus, and "AppleACPIPCI"
// is the name of a bus driver, not of a graphics one.  An entry that publishes
// no IOClass is named by the class it is an instance of, which IOKit answers
// for every entry there is.
func driver(e *Entry) string {
	for _, k := range []string{"IOClass", "CFBundleIdentifier"} {
		if v, ok := e.Props[k]; ok && v.Kind == KindString {
			if s := strings.TrimSpace(v.Str); s != "" {
				return s
			}
		}
	}
	return e.Class
}

// text returns the first of several keys that is present and readable as text.
// A CFData holding printable ASCII counts as text: IOKit publishes the model
// name of a PCI card exactly that way.
func text(e *Entry, up bool, keys ...string) string {
	for _, k := range keys {
		var v Value
		var ok bool
		if up {
			v, ok = e.Prop(k)
		} else {
			v, ok = e.Props[k]
		}
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

// busAddress renders where the card sits on the PCI bus in the notation Linux
// uses — 0000:00:1d.0 — so the two collectors do not describe one card two
// ways.  IOKit publishes the same three numbers in decimal under "pcidebug",
// and a card that is not on a bus has no such property and gets no address.
func busAddress(e *Entry) string {
	v, ok := e.Prop("pcidebug")
	if !ok {
		return ""
	}
	s := ""
	switch v.Kind {
	case KindString:
		s = v.Str
	case KindData:
		s, _ = asciiOf(v.Data)
	}
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) != 3 {
		return ""
	}
	// The function number sometimes carries a suffix in brackets.
	if i := strings.IndexByte(parts[2], '('); i >= 0 {
		parts[2] = parts[2][:i]
	}
	n := make([]int, 3)
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || v < 0 || v > 0xffff {
			return ""
		}
		n[i] = v
	}
	return "0000:" + hex2(n[0]) + ":" + hex2(n[1]) + "." + strconv.FormatInt(int64(n[2]), 16)
}

func hex2(n int) string {
	s := strconv.FormatInt(int64(n), 16)
	if len(s) < 2 {
		return "0" + s
	}
	return s
}

// pciID renders vendor-id or device-id.  IOKit publishes both as four bytes,
// little-endian, of which the low two are the identifier.  It is written here
// as four hexadecimal digits and nothing else, which is exactly what the Linux
// collector puts in the same field after stripping the 0x of
// /sys/bus/pci/devices/*/vendor: one card must not be described two ways
// depending on which system read it.
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
	b := []byte("0000")
	for i := 0; i < 4; i++ {
		b[3-i] = hexDigits[(n>>(4*i))&0xf]
	}
	return string(b)
}

// vram reads how much memory the card has.  IOKit publishes it two ways and
// the drivers do not agree on which: "VRAM,totalMB" in megabytes and
// "VRAM,totalsize" in bytes, either as a number or as four bytes.  Both are
// read; a card that publishes neither gets no memory, and that is the answer
// for every card that shares the machine's memory instead of having its own.
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
