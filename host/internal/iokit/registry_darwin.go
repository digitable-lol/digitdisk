// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package iokit

// Match returns every entry of one IOKit class, together with up to parents
// levels of the nodes they hang off.
//
// Ownership is the whole difficulty of this call, and it is settled here so no
// caller has to think about it:
//
//   - the matching dictionary is consumed by IOServiceGetMatchingServices and
//     is therefore never released here;
//   - every entry IOIteratorNext hands out carries a reference, and every one
//     of them is released before this function returns;
//   - every property dictionary is created for us and released the same way.
//
// What comes back is Go values only.  No CFTypeRef outlives this call, which
// is what makes the rest of the tree able to use the answer without knowing
// that Core Foundation exists.
func Match(class string, parents int) ([]Entry, error) {
	iter, err := matchingServices(class)
	if err != nil {
		return nil, err
	}
	if iter == 0 {
		return nil, notFound(class)
	}
	defer objectRelease(iter)

	var out []Entry
	for obj := iteratorNext(iter); obj != 0; obj = iteratorNext(iter) {
		e := readEntry(obj, parents)
		objectRelease(obj)
		if e != nil {
			out = append(out, *e)
		}
	}
	if len(out) == 0 {
		return nil, notFound(class)
	}
	return out, nil
}

// readEntry reads one node and, if asked, walks up the IOService plane.  The
// walk stops at the root, where IORegistryEntryGetParentEntry answers with a
// failure rather than with a zero — that failure is the end of the tree and
// not an error to report.
func readEntry(obj ref, parents int) *Entry {
	e := &Entry{
		Name:  entryName(iok_IORegistryEntryGetName_trampoline_addr, obj),
		Class: entryName(iok_IOObjectGetClass_trampoline_addr, obj),
		Props: properties(obj),
	}
	if parents <= 0 {
		return e
	}
	parent, err := parentEntry(obj)
	if err != nil || parent == 0 {
		return e
	}
	e.Parent = readEntry(parent, parents-1)
	objectRelease(parent)
	return e
}

// properties copies one entry's property dictionary into a Go map.
func properties(obj ref) map[string]Value {
	dict, err := entryProperties(obj)
	if err != nil {
		return nil
	}
	defer cfRelease(dict)

	keys, values := cfDictionaryPairs(dict)
	out := make(map[string]Value, len(keys))
	for i := range keys {
		name, ok := cfString(keys[i])
		if !ok || name == "" {
			continue
		}
		out[name] = decode(values[i])
	}
	return out
}

// typeIDs are asked for once: CFStringGetTypeID and its kin are functions
// whose answer is fixed for the life of the process, and asking per property
// would be a call per property for a number that never changes.
var typeIDs = struct {
	str, num, boolean, data uintptr
	read                    bool
}{}

func kindOf(v ref) Kind {
	if !typeIDs.read {
		typeIDs.str = cfTypeIDOf(iok_CFStringGetTypeID_trampoline_addr)
		typeIDs.num = cfTypeIDOf(iok_CFNumberGetTypeID_trampoline_addr)
		typeIDs.boolean = cfTypeIDOf(iok_CFBooleanGetTypeID_trampoline_addr)
		typeIDs.data = cfTypeIDOf(iok_CFDataGetTypeID_trampoline_addr)
		typeIDs.read = true
	}
	switch cfTypeID(v) {
	case typeIDs.str:
		return KindString
	case typeIDs.num:
		return KindNumber
	case typeIDs.boolean:
		return KindBool
	case typeIDs.data:
		return KindData
	}
	return KindOther
}

// decode turns one Core Foundation value into a Go one.  A type this package
// does not read comes back as KindOther and empty: the point of the exercise
// is numbers that are true, and a half-read CFArray is not one.
func decode(v ref) Value {
	if v == 0 {
		return Value{}
	}
	switch kindOf(v) {
	case KindString:
		s, ok := cfString(v)
		if !ok {
			return Value{}
		}
		return Value{Kind: KindString, Str: s}
	case KindNumber:
		n, ok := cfNumber(v)
		if !ok {
			return Value{}
		}
		return Value{Kind: KindNumber, Num: n}
	case KindBool:
		return Value{Kind: KindBool, Bool: cfBool(v)}
	case KindData:
		return Value{Kind: KindData, Data: cfData(v)}
	}
	return Value{}
}
