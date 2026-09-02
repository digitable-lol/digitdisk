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
//	mach_host_self       — the host port those two are asked through
//	proc_pidinfo(3)      — what one process costs
//
// The umbrella library re-exports libsystem_kernel and libsystem_c, where
// these actually live, so one path resolves all five.

//go:cgo_import_dynamic libc_sysctl sysctl "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_host_statistics host_statistics "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_host_statistics64 host_statistics64 "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_mach_host_self mach_host_self "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic libc_proc_pidinfo proc_pidinfo "/usr/lib/libSystem.B.dylib"

// The addresses of the assembly stubs, filled in by the .s files next to this
// one — one per architecture, because a stub is two instructions and those
// differ.
var (
	libc_sysctl_trampoline_addr            uintptr
	libc_host_statistics_trampoline_addr   uintptr
	libc_host_statistics64_trampoline_addr uintptr
	libc_mach_host_self_trampoline_addr    uintptr
	libc_proc_pidinfo_trampoline_addr      uintptr
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
