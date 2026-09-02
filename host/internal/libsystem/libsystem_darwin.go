// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build darwin

package libsystem

import (
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

// The symbols this package calls, all of them documented C functions that
// /usr/lib/libSystem.B.dylib exports:
//
//	sysctl(3)            — a sysctl node named by its whole MIB
//	host_statistics      — machine-wide statistics, 32-bit flavors
//	host_statistics64    — machine-wide statistics, 64-bit flavors
//	host_processor_info  — the same statistics, one processor at a time
//	mach_host_self       — the host port those three are asked through
//	proc_pidinfo(3)      — what one process costs
//	task_self_trap       — this task's own port (<mach/mach_traps.h>)
//	vm_deallocate        — giving back what host_processor_info allocated
//	memcpy(3)            — copying out of it without making its address a pointer
//
// The umbrella library re-exports libsystem_kernel and libsystem_c, where
// these actually live, so one path resolves all nine.

//go:cgo_import_dynamic libc_sysctl sysctl "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_host_statistics host_statistics "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_host_statistics64 host_statistics64 "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_host_processor_info host_processor_info "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_mach_host_self mach_host_self "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_proc_pidinfo proc_pidinfo "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_task_self_trap task_self_trap "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_vm_deallocate vm_deallocate "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_memcpy memcpy "/usr/lib/libSystem.B.dylib"

// The addresses of the assembly stubs, filled in by the .s files next to this
// one — one per architecture, because a stub is two instructions and those
// differ.
var (
	libc_sysctl_trampoline_addr              uintptr
	libc_host_statistics_trampoline_addr     uintptr
	libc_host_statistics64_trampoline_addr   uintptr
	libc_host_processor_info_trampoline_addr uintptr
	libc_mach_host_self_trampoline_addr      uintptr
	libc_proc_pidinfo_trampoline_addr        uintptr
	libc_task_self_trap_trampoline_addr      uintptr
	libc_vm_deallocate_trampoline_addr       uintptr
	libc_memcpy_trampoline_addr              uintptr
)

// syscall.syscall6 and syscall.rawSyscall call a C function by address on
// macOS, where the standard library reaches the kernel through libSystem
// rather than through traps of its own.  Both are marked in the standard
// library as callable this way; this is the same door x/sys/unix uses.

//go:linkname syscallSix syscall.syscall6
func syscallSix(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname rawSyscallThree syscall.rawSyscall
func rawSyscallThree(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

// hostPort is the send right mach_host_self hands out.  It is taken once: each
// call adds a reference, and a command-line tool that asked for one per
// reading would leak a port per reading.
var hostPort = sync.OnceValue(func() uintptr {
	port, _, _ := rawSyscallThree(libc_mach_host_self_trampoline_addr, 0, 0, 0)
	return port
})

// HostStatistics asks the host port for one flavor of machine-wide statistics
// and returns exactly the bytes the kernel says it wrote.
//
// The element count is both an argument and a result: it goes in as the size
// of the buffer in four-byte words and comes back as the number of words
// actually filled.  Returning the answer trimmed to that length is what lets
// the decoders insist on an exact size — a kernel whose structure is not the
// one they describe hands back a different count and is refused rather than
// misread.  macOS 26 is exactly that case for the memory flavor: it grew the
// structure by a field, and answers a request for the older length with the
// older length.
func HostStatistics(flavor, count int) ([]byte, error) {
	return hostStatistics(libc_host_statistics_trampoline_addr, flavor, count)
}

// HostStatistics64 is HostStatistics for the flavors whose structures carry
// 64-bit counters.
func HostStatistics64(flavor, count int) ([]byte, error) {
	return hostStatistics(libc_host_statistics64_trampoline_addr, flavor, count)
}

func hostStatistics(fn uintptr, flavor, count int) ([]byte, error) {
	if count <= 0 {
		return nil, syscall.EINVAL
	}
	buf := make([]byte, 4*count)
	n := uint32(count)
	// kern_return_t host_statistics(host_t, host_flavor_t,
	//                               host_info_t, mach_msg_type_number_t *)
	rc, _, _ := syscallSix(fn, hostPort(), uintptr(uint32(flavor)),
		uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)), 0, 0)
	// A Mach call answers with a kern_return_t, not with the -1 and errno of
	// a system call, so the result is read here and not left to the wrapper.
	if kr := int32(rc); kr != 0 {
		return nil, kernError(kr)
	}
	if int(n) > count {
		return nil, syscall.EINVAL
	}
	return buf[:4*int(n)], nil
}

// taskPort is this task's own port, the one vm_deallocate is addressed to.
// Like hostPort it is taken once.  task_self_trap is the kernel trap the
// libSystem macro mach_task_self() is built on, and it is what a caller
// without cgo can reach: the macro reads a variable, and a variable is not a
// symbol this door opens.
var taskPort = sync.OnceValue(func() uintptr {
	port, _, _ := rawSyscallThree(libc_task_self_trap_trampoline_addr, 0, 0, 0)
	return port
})

// ProcessorInfo asks the host for one flavor of per-processor statistics and
// returns the number of processors together with a copy of the array.
//
// Unlike host_statistics, this call does not fill a buffer of ours: the
// kernel allocates the array in this task's address space and hands over the
// address, and the caller owns it from then on.  So the bytes are copied out
// and the allocation is given straight back — a screen that takes a snapshot
// every two seconds would otherwise lose a page of memory per reading, which
// on a machine with many processors adds up to a leak a person would notice
// by the end of the day.
//
// A task port of zero means the trap did not answer.  Nothing is published
// then: reading the array we could not free is not worth the leak.
func ProcessorInfo(flavor int) (cpus int, data []byte, err error) {
	if taskPort() == 0 {
		return 0, nil, syscall.EINVAL
	}
	var count uint32 // natural_t *out_processor_count
	var addr uintptr // processor_info_array_t *out_processor_info
	var words uint32 // mach_msg_type_number_t *out_processor_infoCnt
	// kern_return_t host_processor_info(host_t, processor_flavor_t,
	//                                   natural_t *, processor_info_array_t *,
	//                                   mach_msg_type_number_t *)
	rc, _, _ := syscallSix(libc_host_processor_info_trampoline_addr, hostPort(),
		uintptr(uint32(flavor)), uintptr(unsafe.Pointer(&count)),
		uintptr(unsafe.Pointer(&addr)), uintptr(unsafe.Pointer(&words)), 0)
	if kr := int32(rc); kr != 0 {
		return 0, nil, kernError(kr)
	}
	if addr == 0 || words == 0 || count == 0 {
		return 0, nil, syscall.EINVAL
	}
	size := 4 * uintptr(words)
	out := make([]byte, size)
	// The copy is made by memcpy(3) rather than by Go.  The address the
	// kernel handed over is not a Go pointer and must not become one: a
	// number turned into a pointer is exactly what `go vet` refuses, and it
	// refuses it for a good reason — the garbage collector would then be
	// looking at memory it does not own.  Handing the number to a C function
	// that expects a number keeps it a number the whole way.
	// void *memcpy(void *dst, const void *src, size_t n)
	syscallSix(libc_memcpy_trampoline_addr,
		uintptr(unsafe.Pointer(&out[0])), addr, size, 0, 0, 0)
	// kern_return_t vm_deallocate(vm_map_t, vm_address_t, vm_size_t)
	syscallSix(libc_vm_deallocate_trampoline_addr, taskPort(), addr, size, 0, 0, 0)
	return int(count), out, nil
}

// ProcTaskInfo asks proc_pidinfo for what one process costs, and returns the
// bytes it wrote.
//
// The wrapper in libSystem answers 0 for every failure and leaves the reason
// in errno, so a short answer is a failure and not a short structure.  The
// common failure is not a bug: PROC_PIDTASKINFO is refused for a process
// belonging to another user unless the caller is the administrator.
func ProcTaskInfo(pid, flavor, size int) ([]byte, error) {
	if size <= 0 {
		return nil, syscall.EINVAL
	}
	buf := make([]byte, size)
	// int proc_pidinfo(int pid, int flavor, uint64_t arg,
	//                  void *buffer, int buffersize)
	n, _, err := syscallSix(libc_proc_pidinfo_trampoline_addr,
		uintptr(uint32(pid)), uintptr(uint32(flavor)), 0,
		uintptr(unsafe.Pointer(&buf[0])), uintptr(uint32(size)), 0)
	if int32(n) <= 0 {
		if err != 0 {
			return nil, err
		}
		return nil, syscall.EINVAL
	}
	if int(n) > size {
		return nil, syscall.EINVAL
	}
	return buf[:int(n)], nil
}

// SysctlRaw reads a sysctl node named by its whole MIB rather than by a
// string, which is the only way to reach a node that takes a pid — the command
// line of one process, KERN_PROCARGS2.
//
// The buffer is allocated by the caller because that is how the node works:
// it is capped by kern.argmax and answers with as much as fits.
func SysctlRaw(mib []int32, size int) ([]byte, error) {
	if len(mib) == 0 || size <= 0 {
		return nil, syscall.EINVAL
	}
	buf := make([]byte, size)
	n := uintptr(size)
	// int sysctl(int *name, u_int namelen, void *oldp, size_t *oldlenp,
	//            void *newp, size_t newlen)
	_, _, err := syscallSix(libc_sysctl_trampoline_addr,
		uintptr(unsafe.Pointer(&mib[0])), uintptr(len(mib)),
		uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)), 0, 0)
	if err != 0 {
		return nil, err
	}
	if int(n) > size {
		return nil, syscall.EINVAL
	}
	return buf[:n], nil
}

// kernError renders a kern_return_t.  Mach numbers its failures in a space of
// its own that has nothing to do with errno, so translating one into the other
// would invent a meaning; the number is carried as it is.
type kernError int32

func (e kernError) Error() string {
	return "ядро ответило кодом " + strconv.Itoa(int(e))
}
