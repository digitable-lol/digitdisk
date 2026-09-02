// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// One stub per imported symbol, and the address of that stub as data.
//
// A stub is a single jump.  It exists because the Go linker resolves the
// dynamic import by name at link time and binds it at start-up, while the
// caller needs an ordinary address to hand to syscall.syscall6; the jump is
// what turns one into the other.  The two files, this one and its twin for the
// other architecture, are identical line for line — the assembler needs the
// architecture in the file name, not in the code.

#include "textflag.h"

TEXT libc_sysctl_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_sysctl(SB)
GLOBL	·libc_sysctl_trampoline_addr(SB), RODATA, $8
DATA	·libc_sysctl_trampoline_addr(SB)/8, $libc_sysctl_trampoline<>(SB)

TEXT libc_host_statistics_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_host_statistics(SB)
GLOBL	·libc_host_statistics_trampoline_addr(SB), RODATA, $8
DATA	·libc_host_statistics_trampoline_addr(SB)/8, $libc_host_statistics_trampoline<>(SB)

TEXT libc_host_statistics64_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_host_statistics64(SB)
GLOBL	·libc_host_statistics64_trampoline_addr(SB), RODATA, $8
DATA	·libc_host_statistics64_trampoline_addr(SB)/8, $libc_host_statistics64_trampoline<>(SB)

TEXT libc_mach_host_self_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_mach_host_self(SB)
GLOBL	·libc_mach_host_self_trampoline_addr(SB), RODATA, $8
DATA	·libc_mach_host_self_trampoline_addr(SB)/8, $libc_mach_host_self_trampoline<>(SB)

TEXT libc_proc_pidinfo_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_proc_pidinfo(SB)
GLOBL	·libc_proc_pidinfo_trampoline_addr(SB), RODATA, $8
DATA	·libc_proc_pidinfo_trampoline_addr(SB)/8, $libc_proc_pidinfo_trampoline<>(SB)

TEXT libc_host_processor_info_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_host_processor_info(SB)
GLOBL	·libc_host_processor_info_trampoline_addr(SB), RODATA, $8
DATA	·libc_host_processor_info_trampoline_addr(SB)/8, $libc_host_processor_info_trampoline<>(SB)

TEXT libc_task_self_trap_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_task_self_trap(SB)
GLOBL	·libc_task_self_trap_trampoline_addr(SB), RODATA, $8
DATA	·libc_task_self_trap_trampoline_addr(SB)/8, $libc_task_self_trap_trampoline<>(SB)

TEXT libc_vm_deallocate_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_vm_deallocate(SB)
GLOBL	·libc_vm_deallocate_trampoline_addr(SB), RODATA, $8
DATA	·libc_vm_deallocate_trampoline_addr(SB)/8, $libc_vm_deallocate_trampoline<>(SB)

TEXT libc_memcpy_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_memcpy(SB)
GLOBL	·libc_memcpy_trampoline_addr(SB), RODATA, $8
DATA	·libc_memcpy_trampoline_addr(SB)/8, $libc_memcpy_trampoline<>(SB)
