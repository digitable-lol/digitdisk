// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package iokit

import (
	"encoding/hex"
	"strconv"
	"strings"
)

// Kind says which Core Foundation type a property came back as.  A property
// whose type this package does not decode keeps KindOther and its bytes are
// not invented.
type Kind int

// The Core Foundation types worth decoding.  Everything else — arrays,
// dictionaries, the odd CFSet — is left as KindOther: the registry entries a
// video card fills carry strings, numbers, flags and raw bytes, and decoding a
// type nobody reads would be code with no reader.
const (
	KindOther Kind = iota
	KindString
	KindNumber
	KindBool
	KindData
)

// Value is one property of a registry entry.
type Value struct {
	Kind Kind
	Str  string
	Num  int64
	Bool bool
	Data []byte
}

// Text renders a value the way a person reads it, and says so when there is
// nothing to render.  A CFData that holds printable ASCII is shown as the text
// it is — IOKit publishes the model name of a PCI card exactly that way — and
// as hex otherwise.
func (v Value) Text() string {
	switch v.Kind {
	case KindString:
		return v.Str
	case KindNumber:
		return strconv.FormatInt(v.Num, 10)
	case KindBool:
		if v.Bool {
			return "true"
		}
		return "false"
	case KindData:
		if s, ok := asciiOf(v.Data); ok {
			return s
		}
		return hex.EncodeToString(v.Data)
	}
	return "?"
}

// asciiOf reads a CFData as the C string IOKit often puts there: printable
// bytes and at most one terminator.  Anything else is not text, and guessing
// that it is would print rubbish beside a card.
func asciiOf(b []byte) (string, bool) {
	if len(b) == 0 {
		return "", false
	}
	if b[len(b)-1] == 0 {
		b = b[:len(b)-1]
	}
	if len(b) == 0 {
		return "", false
	}
	for _, c := range b {
		if c < 0x20 || c > 0x7e {
			return "", false
		}
	}
	return string(b), true
}

// Entry is one node of the IORegistry: what it is called, what class it is an
// instance of, everything it publishes, and — when the caller asked for them —
// the same for the node it hangs off.
type Entry struct {
	Name   string
	Class  string
	Props  map[string]Value
	Parent *Entry
}

// Prop returns one property by key, from the entry or, failing that, from its
// parents.  A video card is described in two places at once on a Mac: the
// accelerator says which driver runs it, and the device it hangs off says what
// it is and how much memory it has.
func (e *Entry) Prop(key string) (Value, bool) {
	for p := e; p != nil; p = p.Parent {
		if v, ok := p.Props[key]; ok {
			return v, true
		}
	}
	return Value{}, false
}

// Dump renders an entry and its parents as lines, sorted by key, with long
// values cut.  It exists for the measurement: what a given Mac publishes is a
// fact to be read off a run, not a list to be remembered.
func (e *Entry) Dump(limit int) []string {
	var out []string
	depth := 0
	for p := e; p != nil; p = p.Parent {
		pad := strings.Repeat("  ", depth)
		out = append(out, pad+"["+p.Class+"] "+p.Name)
		for _, k := range sortedKeys(p.Props) {
			s := p.Props[k].Text()
			if limit > 0 && len(s) > limit {
				s = s[:limit] + "…"
			}
			out = append(out, pad+"  "+k+" = "+s)
		}
		depth++
	}
	return out
}

func sortedKeys(m map[string]Value) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

// notMac is what every call answers where there is no IORegistry to read.
type notMac struct{}

func (notMac) Error() string { return "реестр IOKit есть только на macOS" }
