// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import (
	"encoding/binary"
	"testing"
	"time"
)

// WHAT THESE TESTS PROVE, AND WHAT THEY DO NOT.
//
// Every buffer below is built here, from the field offsets as the XNU headers
// give them — written out again as literal numbers, not taken from the
// constants under test, so a constant that drifts breaks a test instead of
// agreeing with itself.  That proves the decoders read the layout they claim
// to read, and that they refuse what they cannot vouch for.
//
// It does not prove the layout is the one a Mac actually answers with: no Mac
// was available while this was written, and no captured buffer from one is in
// the tree.  That check happens at run time, on the machine itself, in Verify
// and VerifyIfList — which is why those exist.

// put64/put32 write a little-endian field at an explicit offset.
func put64(b []byte, off int, v uint64) { binary.LittleEndian.PutUint64(b[off:], v) }
func put32(b []byte, off int, v uint32) { binary.LittleEndian.PutUint32(b[off:], v) }
func put16(b []byte, off int, v uint16) { binary.LittleEndian.PutUint16(b[off:], v) }

func TestPaddedRestoresTheByteSysctlTrims(t *testing.T) {
	full := make([]byte, 8)
	put64(full, 0, 8*1024*1024*1024) // 8 GiB: the top byte is zero

	// syscall.Sysctl strips one trailing NUL, so this is what a caller sees.
	trimmed := full[:7]
	got, ok := Padded(trimmed, 8)
	if !ok {
		t.Fatal("a buffer one byte short must be restorable")
	}
	if v, _ := ParseUint64(got); v != 8*1024*1024*1024 {
		t.Errorf("restored value = %d, want %d", v, uint64(8*1024*1024*1024))
	}
	if _, ok := Padded(full[:6], 8); ok {
		t.Error("two bytes short is a different structure and must be refused")
	}
	if _, ok := Padded(nil, 8); ok {
		t.Error("an empty buffer must be refused")
	}
	long, ok := Padded(make([]byte, 12), 8)
	if !ok || len(long) != 8 {
		t.Errorf("a grown structure must be read by its prefix, got %d bytes, ok=%v", len(long), ok)
	}
}

func TestParseTimeval(t *testing.T) {
	// struct timeval: int64 tv_sec at 0, int32 tv_usec at 8, 4 bytes padding.
	b := make([]byte, 16)
	put64(b, 0, 1700000000)
	put32(b, 8, 250000)

	got, ok := ParseTimeval(b)
	if !ok {
		t.Fatal("a valid timeval was refused")
	}
	if want := time.Unix(1700000000, 250000000); !got.Equal(want) {
		t.Errorf("boot time = %v, want %v", got, want)
	}
	// the same value as sysctl hands it over, one byte short
	if got2, ok := ParseTimeval(b[:15]); !ok || !got2.Equal(got) {
		t.Errorf("trimmed timeval = %v (ok=%v), want %v", got2, ok, got)
	}

	zero := make([]byte, 16)
	if _, ok := ParseTimeval(zero); ok {
		t.Error("a zero timeval is not a boot time and must be refused")
	}
	bad := make([]byte, 16)
	put64(bad, 0, 1700000000)
	put32(bad, 8, 9_000_000) // microseconds out of range
	if _, ok := ParseTimeval(bad); ok {
		t.Error("microseconds outside their range must be refused")
	}
}

func TestParseLoadAvg(t *testing.T) {
	// struct loadavg: fixpt_t ldavg[3] at 0/4/8, long fscale at 16.
	b := make([]byte, 24)
	put32(b, 0, 2048)  // 1.00
	put32(b, 4, 1024)  // 0.50
	put32(b, 8, 5120)  // 2.50
	put64(b, 16, 2048) // fscale

	one, five, fifteen, ok := ParseLoadAvg(b)
	if !ok {
		t.Fatal("a valid loadavg was refused")
	}
	if one != 1 || five != 0.5 || fifteen != 2.5 {
		t.Errorf("averages = %v/%v/%v, want 1/0.5/2.5", one, five, fifteen)
	}

	noScale := make([]byte, 24)
	put32(noScale, 0, 2048)
	if _, _, _, ok := ParseLoadAvg(noScale); ok {
		t.Error("a zero scale must be refused, not divided by")
	}
	if _, _, _, ok := ParseLoadAvg(b[:4]); ok {
		t.Error("a short buffer must be refused")
	}
}

func TestParseSwapUsage(t *testing.T) {
	// struct xsw_usage: u_int64 total/avail/used at 0/8/16, u_int32 pagesize
	// at 24, boolean_t encrypted at 28.
	b := make([]byte, 32)
	put64(b, 0, 4<<30)
	put64(b, 8, 3<<30)
	put64(b, 16, 1<<30)
	put32(b, 24, 4096)
	put32(b, 28, 1)

	s, ok := ParseSwapUsage(b)
	if !ok {
		t.Fatal("a valid xsw_usage was refused")
	}
	if s.Total != 4<<30 || s.Avail != 3<<30 || s.Used != 1<<30 || s.PageSize != 4096 || !s.Encrypted {
		t.Errorf("swap = %+v", s)
	}

	// swap turned off is a measurement, not a failure
	if s, ok := ParseSwapUsage(make([]byte, 32)); !ok || s.Total != 0 {
		t.Errorf("a machine without swap must parse as zeros, got %+v ok=%v", s, ok)
	}

	impossible := make([]byte, 32)
	put64(impossible, 0, 1<<30)
	put64(impossible, 16, 4<<30) // used > total
	if _, ok := ParseSwapUsage(impossible); ok {
		t.Error("used swap above the total means a wrong layout and must be refused")
	}
}
