// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// One stub per imported symbol, and the address of that stub as data.
//
// A stub is a single jump, exactly as in internal/libsystem: the Go linker
// resolves the dynamic import by name at link time and binds it at start-up,
// while the caller needs an ordinary address to hand to syscall.syscall6.  The
// two files, this one and its twin for the other architecture, are identical
// line for line — the assembler needs the architecture in the file name, not
// in the code.

#include "textflag.h"

TEXT iok_IOServiceMatching_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IOServiceMatching(SB)
GLOBL	·iok_IOServiceMatching_trampoline_addr(SB), RODATA, $8
DATA	·iok_IOServiceMatching_trampoline_addr(SB)/8, $iok_IOServiceMatching_trampoline<>(SB)

TEXT iok_IOServiceGetMatchingServices_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IOServiceGetMatchingServices(SB)
GLOBL	·iok_IOServiceGetMatchingServices_trampoline_addr(SB), RODATA, $8
DATA	·iok_IOServiceGetMatchingServices_trampoline_addr(SB)/8, $iok_IOServiceGetMatchingServices_trampoline<>(SB)

TEXT iok_IOIteratorNext_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IOIteratorNext(SB)
GLOBL	·iok_IOIteratorNext_trampoline_addr(SB), RODATA, $8
DATA	·iok_IOIteratorNext_trampoline_addr(SB)/8, $iok_IOIteratorNext_trampoline<>(SB)

TEXT iok_IOObjectRelease_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IOObjectRelease(SB)
GLOBL	·iok_IOObjectRelease_trampoline_addr(SB), RODATA, $8
DATA	·iok_IOObjectRelease_trampoline_addr(SB)/8, $iok_IOObjectRelease_trampoline<>(SB)

TEXT iok_IOObjectGetClass_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IOObjectGetClass(SB)
GLOBL	·iok_IOObjectGetClass_trampoline_addr(SB), RODATA, $8
DATA	·iok_IOObjectGetClass_trampoline_addr(SB)/8, $iok_IOObjectGetClass_trampoline<>(SB)

TEXT iok_IORegistryEntryGetName_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IORegistryEntryGetName(SB)
GLOBL	·iok_IORegistryEntryGetName_trampoline_addr(SB), RODATA, $8
DATA	·iok_IORegistryEntryGetName_trampoline_addr(SB)/8, $iok_IORegistryEntryGetName_trampoline<>(SB)

TEXT iok_IORegistryEntryGetParentEntry_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IORegistryEntryGetParentEntry(SB)
GLOBL	·iok_IORegistryEntryGetParentEntry_trampoline_addr(SB), RODATA, $8
DATA	·iok_IORegistryEntryGetParentEntry_trampoline_addr(SB)/8, $iok_IORegistryEntryGetParentEntry_trampoline<>(SB)

TEXT iok_IORegistryEntryCreateCFProperties_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_IORegistryEntryCreateCFProperties(SB)
GLOBL	·iok_IORegistryEntryCreateCFProperties_trampoline_addr(SB), RODATA, $8
DATA	·iok_IORegistryEntryCreateCFProperties_trampoline_addr(SB)/8, $iok_IORegistryEntryCreateCFProperties_trampoline<>(SB)

TEXT iok_CFRelease_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFRelease(SB)
GLOBL	·iok_CFRelease_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFRelease_trampoline_addr(SB)/8, $iok_CFRelease_trampoline<>(SB)

TEXT iok_CFGetTypeID_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFGetTypeID(SB)
GLOBL	·iok_CFGetTypeID_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFGetTypeID_trampoline_addr(SB)/8, $iok_CFGetTypeID_trampoline<>(SB)

TEXT iok_CFStringGetTypeID_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFStringGetTypeID(SB)
GLOBL	·iok_CFStringGetTypeID_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFStringGetTypeID_trampoline_addr(SB)/8, $iok_CFStringGetTypeID_trampoline<>(SB)

TEXT iok_CFNumberGetTypeID_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFNumberGetTypeID(SB)
GLOBL	·iok_CFNumberGetTypeID_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFNumberGetTypeID_trampoline_addr(SB)/8, $iok_CFNumberGetTypeID_trampoline<>(SB)

TEXT iok_CFDataGetTypeID_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFDataGetTypeID(SB)
GLOBL	·iok_CFDataGetTypeID_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFDataGetTypeID_trampoline_addr(SB)/8, $iok_CFDataGetTypeID_trampoline<>(SB)

TEXT iok_CFBooleanGetTypeID_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFBooleanGetTypeID(SB)
GLOBL	·iok_CFBooleanGetTypeID_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFBooleanGetTypeID_trampoline_addr(SB)/8, $iok_CFBooleanGetTypeID_trampoline<>(SB)

TEXT iok_CFDictionaryGetCount_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFDictionaryGetCount(SB)
GLOBL	·iok_CFDictionaryGetCount_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFDictionaryGetCount_trampoline_addr(SB)/8, $iok_CFDictionaryGetCount_trampoline<>(SB)

TEXT iok_CFDictionaryGetKeysAndValues_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFDictionaryGetKeysAndValues(SB)
GLOBL	·iok_CFDictionaryGetKeysAndValues_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFDictionaryGetKeysAndValues_trampoline_addr(SB)/8, $iok_CFDictionaryGetKeysAndValues_trampoline<>(SB)

TEXT iok_CFStringGetCString_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFStringGetCString(SB)
GLOBL	·iok_CFStringGetCString_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFStringGetCString_trampoline_addr(SB)/8, $iok_CFStringGetCString_trampoline<>(SB)

TEXT iok_CFStringGetLength_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFStringGetLength(SB)
GLOBL	·iok_CFStringGetLength_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFStringGetLength_trampoline_addr(SB)/8, $iok_CFStringGetLength_trampoline<>(SB)

TEXT iok_CFNumberGetValue_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFNumberGetValue(SB)
GLOBL	·iok_CFNumberGetValue_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFNumberGetValue_trampoline_addr(SB)/8, $iok_CFNumberGetValue_trampoline<>(SB)

TEXT iok_CFBooleanGetValue_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFBooleanGetValue(SB)
GLOBL	·iok_CFBooleanGetValue_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFBooleanGetValue_trampoline_addr(SB)/8, $iok_CFBooleanGetValue_trampoline<>(SB)

TEXT iok_CFDataGetLength_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFDataGetLength(SB)
GLOBL	·iok_CFDataGetLength_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFDataGetLength_trampoline_addr(SB)/8, $iok_CFDataGetLength_trampoline<>(SB)

TEXT iok_CFDataGetBytePtr_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_CFDataGetBytePtr(SB)
GLOBL	·iok_CFDataGetBytePtr_trampoline_addr(SB), RODATA, $8
DATA	·iok_CFDataGetBytePtr_trampoline_addr(SB)/8, $iok_CFDataGetBytePtr_trampoline<>(SB)

TEXT iok_memcpy_trampoline<>(SB),NOSPLIT,$0-0
	JMP	iok_memcpy(SB)
GLOBL	·iok_memcpy_trampoline_addr(SB), RODATA, $8
DATA	·iok_memcpy_trampoline_addr(SB)/8, $iok_memcpy_trampoline<>(SB)
