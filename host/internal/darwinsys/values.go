// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "time"

// TimevalSize is sizeof(struct timeval) on 64-bit macOS: a 64-bit tv_sec, a
// 32-bit tv_usec, and four bytes of tail padding (<sys/_types/_timeval.h>).
const TimevalSize = 16

// ParseTimeval decodes a struct timeval, which is what kern.boottime answers
// with.  It refuses a time at or before the epoch and a microsecond field
// outside its range: both mean the bytes are not a timeval, and a boot time
// invented out of such bytes would silently become a wrong uptime.
func ParseTimeval(b []byte) (time.Time, bool) {
	buf, ok := Padded(b, TimevalSize)
	if !ok {
		return time.Time{}, false
	}
	sec := int64(u64(buf, 0))
	usec := int64(i32(buf, 8))
	if sec <= 0 || usec < 0 || usec >= 1_000_000 {
		return time.Time{}, false
	}
	return time.Unix(sec, usec*1000), true
}

// LoadAvgSize is sizeof(struct loadavg) on 64-bit macOS: fixpt_t ldavg[3] —
// three 32-bit fixed-point numbers — then four bytes of padding, then a
// 64-bit long fscale (<sys/resource.h>, exposed as vm.loadavg).
const LoadAvgSize = 24

// ParseLoadAvg decodes vm.loadavg into the three averages.
//
// The scale comes from the structure itself rather than from the usual
// constant 2048: the kernel ships the divisor precisely so a caller does not
// have to assume one.  A zero scale is a refusal, not a division.
func ParseLoadAvg(b []byte) (one, five, fifteen float64, ok bool) {
	buf, ok := Padded(b, LoadAvgSize)
	if !ok {
		return 0, 0, 0, false
	}
	scale := float64(u64(buf, 16))
	if scale <= 0 {
		return 0, 0, 0, false
	}
	return float64(u32(buf, 0)) / scale,
		float64(u32(buf, 4)) / scale,
		float64(u32(buf, 8)) / scale,
		true
}

// SwapUsageSize is sizeof(struct xsw_usage) (<sys/sysctl.h>): three 64-bit
// byte counts, a 32-bit page size and a 32-bit boolean.
const SwapUsageSize = 32

// SwapUsage is vm.swapusage: the swap file totals macOS reports.
type SwapUsage struct {
	Total     uint64
	Avail     uint64
	Used      uint64
	PageSize  uint32
	Encrypted bool
}

// ParseSwapUsage decodes vm.swapusage.
//
// The validation is arithmetic the kernel guarantees and a wrong layout would
// break at once: neither the used nor the available part can exceed the total.
// A machine with swap turned off reports three zeros, which is a measurement
// and passes.
func ParseSwapUsage(b []byte) (SwapUsage, bool) {
	buf, ok := Padded(b, SwapUsageSize)
	if !ok {
		return SwapUsage{}, false
	}
	s := SwapUsage{
		Total:     u64(buf, 0),
		Avail:     u64(buf, 8),
		Used:      u64(buf, 16),
		PageSize:  u32(buf, 24),
		Encrypted: u32(buf, 28) != 0,
	}
	if s.Used > s.Total || s.Avail > s.Total {
		return SwapUsage{}, false
	}
	return s, true
}

// ParseUint64 decodes a 64-bit sysctl node (hw.memsize and its kind).
func ParseUint64(b []byte) (uint64, bool) {
	buf, ok := Padded(b, 8)
	if !ok {
		return 0, false
	}
	return u64(buf, 0), true
}

// ParseUint32 decodes a 32-bit sysctl node.
func ParseUint32(b []byte) (uint32, bool) {
	buf, ok := Padded(b, 4)
	if !ok {
		return 0, false
	}
	return u32(buf, 0), true
}
