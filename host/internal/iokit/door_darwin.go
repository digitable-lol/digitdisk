// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package iokit

import (
	"runtime"
	"strconv"
	"syscall"
	"unsafe"
)

// The symbols this package calls.  All of them are documented C functions, and
// none of them is private.
//
// From IOKit — the registry itself:
//
//	IOServiceMatching               a matching dictionary for one class name
//	IOServiceGetMatchingServices    the entries that dictionary matches
//	IOIteratorNext                  the next entry, 0 when there are no more
//	IOObjectRelease                 giving an entry or an iterator back
//	IOObjectGetClass                the class an entry is an instance of
//	IORegistryEntryGetName          the entry's own name in the tree
//	IORegistryEntryGetParentEntry   one step up a named plane
//	IORegistryEntryCreateCFProperties   the entry's whole property dictionary
//
// From CoreFoundation — reading what came back:
//
//	CFRelease, CFGetTypeID
//	CFStringGetTypeID, CFNumberGetTypeID, CFDataGetTypeID, CFBooleanGetTypeID
//	CFDictionaryGetCount, CFDictionaryGetKeysAndValues
//	CFStringGetCString, CFStringGetLength
//	CFNumberGetValue, CFBooleanGetValue
//	CFDataGetLength, CFDataGetBytePtr
//
// From libSystem — one function, for the same reason internal/libsystem needs
// it: copying out of memory that is not ours without making its address a Go
// pointer.
//
//	memcpy(3)
//
// The local names are prefixed iok_ rather than libc_ on purpose: a dynamic
// import is a link-time symbol in one namespace for the whole program, and
// internal/libsystem already imports memcpy as libc_memcpy.

//go:cgo_import_dynamic iok_IOServiceMatching IOServiceMatching "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IOServiceGetMatchingServices IOServiceGetMatchingServices "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IOIteratorNext IOIteratorNext "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IOObjectRelease IOObjectRelease "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IOObjectGetClass IOObjectGetClass "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IORegistryEntryGetName IORegistryEntryGetName "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IORegistryEntryGetParentEntry IORegistryEntryGetParentEntry "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"
//go:cgo_import_dynamic iok_IORegistryEntryCreateCFProperties IORegistryEntryCreateCFProperties "/System/Library/Frameworks/IOKit.framework/Versions/A/IOKit"

//go:cgo_import_dynamic iok_CFRelease CFRelease "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFGetTypeID CFGetTypeID "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFStringGetTypeID CFStringGetTypeID "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFNumberGetTypeID CFNumberGetTypeID "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFDataGetTypeID CFDataGetTypeID "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFBooleanGetTypeID CFBooleanGetTypeID "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFDictionaryGetCount CFDictionaryGetCount "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFDictionaryGetKeysAndValues CFDictionaryGetKeysAndValues "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFStringGetCString CFStringGetCString "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFStringGetLength CFStringGetLength "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFNumberGetValue CFNumberGetValue "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFBooleanGetValue CFBooleanGetValue "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFDataGetLength CFDataGetLength "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"
//go:cgo_import_dynamic iok_CFDataGetBytePtr CFDataGetBytePtr "/System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation"

//go:cgo_import_dynamic iok_memcpy memcpy "/usr/lib/libSystem.B.dylib"

// The addresses of the assembly stubs, filled in by the .s files next to this
// one — one per architecture, because a stub is two instructions and those
// differ.
var (
	iok_IOServiceMatching_trampoline_addr                 uintptr
	iok_IOServiceGetMatchingServices_trampoline_addr      uintptr
	iok_IOIteratorNext_trampoline_addr                    uintptr
	iok_IOObjectRelease_trampoline_addr                   uintptr
	iok_IOObjectGetClass_trampoline_addr                  uintptr
	iok_IORegistryEntryGetName_trampoline_addr            uintptr
	iok_IORegistryEntryGetParentEntry_trampoline_addr     uintptr
	iok_IORegistryEntryCreateCFProperties_trampoline_addr uintptr

	iok_CFRelease_trampoline_addr                    uintptr
	iok_CFGetTypeID_trampoline_addr                  uintptr
	iok_CFStringGetTypeID_trampoline_addr            uintptr
	iok_CFNumberGetTypeID_trampoline_addr            uintptr
	iok_CFDataGetTypeID_trampoline_addr              uintptr
	iok_CFBooleanGetTypeID_trampoline_addr           uintptr
	iok_CFDictionaryGetCount_trampoline_addr         uintptr
	iok_CFDictionaryGetKeysAndValues_trampoline_addr uintptr
	iok_CFStringGetCString_trampoline_addr           uintptr
	iok_CFStringGetLength_trampoline_addr            uintptr
	iok_CFNumberGetValue_trampoline_addr             uintptr
	iok_CFBooleanGetValue_trampoline_addr            uintptr
	iok_CFDataGetLength_trampoline_addr              uintptr
	iok_CFDataGetBytePtr_trampoline_addr             uintptr

	iok_memcpy_trampoline_addr uintptr
)

// syscall.syscall6 calls a C function by address on macOS, where the standard
// library reaches the kernel through libSystem rather than through traps of its
// own.  It is the same door internal/libsystem uses, and the same one
// x/sys/unix uses.

//go:linkname syscallSix syscall.syscall6
func syscallSix(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// Constants of the frameworks that are literals in the headers rather than
// exported symbols, and are therefore written out here.
const (
	// kIOMainPortDefault (kIOMasterPortDefault before macOS 12) is
	// MACH_PORT_NULL, and IOKit documents NULL as "the default port".
	mainPortDefault = 0
	// kCFAllocatorDefault is NULL, which every CF creator reads as "the
	// default allocator".
	allocatorDefault = 0
	// kNilOptions.
	nilOptions = 0
	// kCFStringEncodingUTF8.
	encodingUTF8 = 0x08000100
	// kCFNumberSInt64Type.
	numberSInt64 = 4
	// kIOServicePlane, the plane a device's parents are looked up in.
	planeService = "IOService\x00"
	// io_name_t is char[128]; every IOKit call that fills one expects a
	// buffer of exactly that size.
	nameLen = 128
)

// A ref is a CFTypeRef or an io_object_t: a number that names something the
// frameworks own.  It is deliberately not a Go pointer — the garbage collector
// must never look at memory it does not own, and `go vet` refuses the
// conversion that would make it do so.
type ref = uintptr

// kernError renders a kern_return_t.  Mach numbers its failures in a space of
// its own that has nothing to do with errno, so translating one into the other
// would invent a meaning; the number is carried as it is.
type kernError int32

func (e kernError) Error() string {
	return "ядро ответило кодом " + strconv.Itoa(int(e))
}

// notFound is the answer when a class name matches nothing.  It is not a
// failure: a Mac with no discrete card has no IOPCIDevice for one.
type notFound string

func (e notFound) Error() string {
	return "в реестре нет ни одной записи класса " + string(e)
}

// serviceMatching builds the matching dictionary for one class name.  The
// dictionary is handed to IOServiceGetMatchingServices, which CONSUMES a
// reference on it — so it is never released here, and releasing it would be a
// use-after-free rather than tidiness.
func serviceMatching(class string) ref {
	name := append([]byte(class), 0)
	d, _, _ := syscallSix(iok_IOServiceMatching_trampoline_addr,
		uintptr(unsafe.Pointer(&name[0])), 0, 0, 0, 0, 0)
	keepAlive(&name[0])
	return d
}

// matchingServices returns an iterator over every entry of one class.
func matchingServices(class string) (ref, error) {
	dict := serviceMatching(class)
	if dict == 0 {
		return 0, notFound(class)
	}
	var iter uint32
	rc, _, _ := syscallSix(iok_IOServiceGetMatchingServices_trampoline_addr,
		mainPortDefault, dict, uintptr(unsafe.Pointer(&iter)), 0, 0, 0)
	keepAlive(&iter)
	if kr := int32(rc); kr != 0 {
		return 0, kernError(kr)
	}
	return ref(iter), nil
}

func iteratorNext(iter ref) ref {
	r, _, _ := syscallSix(iok_IOIteratorNext_trampoline_addr, iter, 0, 0, 0, 0, 0)
	return ref(uint32(r))
}

func objectRelease(obj ref) {
	syscallSix(iok_IOObjectRelease_trampoline_addr, obj, 0, 0, 0, 0, 0)
}

// entryName reads one of the two io_name_t answers — the entry's own name in
// the tree, or the class it is an instance of.
func entryName(fn uintptr, obj ref) string {
	var buf [nameLen]byte
	rc, _, _ := syscallSix(fn, obj, uintptr(unsafe.Pointer(&buf[0])), 0, 0, 0, 0)
	keepAlive(&buf[0])
	if int32(rc) != 0 {
		return ""
	}
	return cstring(buf[:])
}

// parentEntry steps one level up the IOService plane.  A card's numbers are
// often not on the accelerator itself but on the device it hangs off.
func parentEntry(obj ref) (ref, error) {
	plane := []byte(planeService)
	var parent uint32
	rc, _, _ := syscallSix(iok_IORegistryEntryGetParentEntry_trampoline_addr,
		obj, uintptr(unsafe.Pointer(&plane[0])),
		uintptr(unsafe.Pointer(&parent)), 0, 0, 0)
	keepAlive(&plane[0])
	keepAlive(&parent)
	if kr := int32(rc); kr != 0 {
		return 0, kernError(kr)
	}
	return ref(parent), nil
}

// entryProperties copies an entry's whole property dictionary out of the
// kernel.  The caller owns the dictionary and must release it.
func entryProperties(obj ref) (ref, error) {
	var dict uintptr
	rc, _, _ := syscallSix(iok_IORegistryEntryCreateCFProperties_trampoline_addr,
		obj, uintptr(unsafe.Pointer(&dict)), allocatorDefault, nilOptions, 0, 0)
	keepAlive(&dict)
	if kr := int32(rc); kr != 0 {
		return 0, kernError(kr)
	}
	if dict == 0 {
		return 0, kernError(0)
	}
	return dict, nil
}

func cfRelease(obj ref) {
	if obj != 0 {
		syscallSix(iok_CFRelease_trampoline_addr, obj, 0, 0, 0, 0, 0)
	}
}

func cfTypeID(obj ref) uintptr {
	id, _, _ := syscallSix(iok_CFGetTypeID_trampoline_addr, obj, 0, 0, 0, 0, 0)
	return id
}

func cfTypeIDOf(fn uintptr) uintptr {
	id, _, _ := syscallSix(fn, 0, 0, 0, 0, 0, 0)
	return id
}

func cfDictionaryCount(dict ref) int {
	n, _, _ := syscallSix(iok_CFDictionaryGetCount_trampoline_addr, dict, 0, 0, 0, 0, 0)
	return int(int64(n))
}

// cfDictionaryPairs copies the keys and values of a dictionary into two slices
// of numbers.  They are numbers and not pointers on purpose: CF owns what they
// point at, and the collector must not be told otherwise.
func cfDictionaryPairs(dict ref) (keys, values []ref) {
	n := cfDictionaryCount(dict)
	if n <= 0 {
		return nil, nil
	}
	keys = make([]ref, n)
	values = make([]ref, n)
	syscallSix(iok_CFDictionaryGetKeysAndValues_trampoline_addr, dict,
		uintptr(unsafe.Pointer(&keys[0])), uintptr(unsafe.Pointer(&values[0])), 0, 0, 0)
	keepAlive(&keys[0])
	keepAlive(&values[0])
	return keys, values
}

// cfString renders a CFStringRef.  The buffer is sized from the character
// count, four bytes per character being the widest UTF-8 can be, plus the
// terminator CFStringGetCString insists on.
func cfString(s ref) (string, bool) {
	n, _, _ := syscallSix(iok_CFStringGetLength_trampoline_addr, s, 0, 0, 0, 0, 0)
	size := 4*int(int64(n)) + 1
	if size < 2 {
		size = 2
	}
	buf := make([]byte, size)
	ok, _, _ := syscallSix(iok_CFStringGetCString_trampoline_addr, s,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(size), encodingUTF8, 0, 0)
	keepAlive(&buf[0])
	if uint8(ok) == 0 {
		return "", false
	}
	return cstring(buf), true
}

func cfNumber(v ref) (int64, bool) {
	var out int64
	ok, _, _ := syscallSix(iok_CFNumberGetValue_trampoline_addr, v,
		numberSInt64, uintptr(unsafe.Pointer(&out)), 0, 0, 0)
	keepAlive(&out)
	return out, uint8(ok) != 0
}

func cfBool(v ref) bool {
	ok, _, _ := syscallSix(iok_CFBooleanGetValue_trampoline_addr, v, 0, 0, 0, 0, 0)
	return uint8(ok) != 0
}

// cfData copies a CFDataRef out.  The copy is made by memcpy(3) and not by Go
// for the same reason internal/libsystem gives: the address CF handed over is
// not a Go pointer and must not become one.
func cfData(v ref) []byte {
	n, _, _ := syscallSix(iok_CFDataGetLength_trampoline_addr, v, 0, 0, 0, 0, 0)
	size := int(int64(n))
	if size <= 0 {
		return nil
	}
	p, _, _ := syscallSix(iok_CFDataGetBytePtr_trampoline_addr, v, 0, 0, 0, 0, 0)
	if p == 0 {
		return nil
	}
	out := make([]byte, size)
	syscallSix(iok_memcpy_trampoline_addr,
		uintptr(unsafe.Pointer(&out[0])), p, uintptr(size), 0, 0, 0)
	keepAlive(&out[0])
	return out
}

// keepAlive holds a Go object alive until the framework function that was
// handed its address has returned.  The conversion to uintptr hides the
// reference from the compiler, and a buffer the collector believes is dead is
// a buffer CoreFoundation may be writing into.
func keepAlive[T any](p *T) { runtime.KeepAlive(p) }

// cstring cuts a C string out of a fixed buffer.
func cstring(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
