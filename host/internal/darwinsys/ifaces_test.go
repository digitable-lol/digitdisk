// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "testing"

// ifinfo2 builds one RTM_IFINFO2 message the way <net/if_var.h> lays out
// struct if_msghdr2 with struct if_data64 inside it.  Offsets are literal
// here — see the note in darwinsys_test.go.
func ifinfo2(index int, mtu uint32, ibytes, obytes uint64) []byte {
	b := make([]byte, 160)
	put16(b, 0, 160) // ifm_msglen
	b[2] = 5         // ifm_version = RTM_VERSION
	b[3] = 0x12      // ifm_type = RTM_IFINFO2
	put16(b, 12, uint16(index))
	put32(b, 40, mtu)     // ifm_data.ifi_mtu
	put64(b, 56, 100)     // ifi_ipackets
	put64(b, 64, 1)       // ifi_ierrors
	put64(b, 72, 200)     // ifi_opackets
	put64(b, 80, 2)       // ifi_oerrors
	put64(b, 96, ibytes)  // ifi_ibytes
	put64(b, 104, obytes) // ifi_obytes
	put64(b, 128, 3)      // ifi_iqdrops
	return b
}

// otherMessage is a message of a type this decoder does not read; it must be
// stepped over by its own length, not misread.
func otherMessage(length int) []byte {
	b := make([]byte, length)
	put16(b, 0, uint16(length))
	b[2] = 5
	b[3] = 0x0e // RTM_IFINFO, the 32-bit one
	return b
}

func TestParseIfList2(t *testing.T) {
	var buf []byte
	buf = append(buf, ifinfo2(1, 16384, 1<<40, 1<<39)...) // lo0, past 4 GiB
	buf = append(buf, otherMessage(96)...)
	buf = append(buf, ifinfo2(14, 1500, 500, 600)...)

	got := ParseIfList2(buf)
	if len(got) != 2 {
		t.Fatalf("interfaces = %d, want 2 (the RTM_IFINFO message is not one)", len(got))
	}
	lo := got[1]
	if lo.MTU != 16384 || lo.RxBytes != 1<<40 || lo.TxBytes != 1<<39 {
		t.Errorf("lo0 = %+v", lo)
	}
	if lo.RxPackets != 100 || lo.TxPackets != 200 || lo.RxErrors != 1 || lo.TxErrors != 2 || lo.RxDropped != 3 {
		t.Errorf("lo0 counters = %+v", lo)
	}
	if en := got[14]; en.Index != 14 || en.RxBytes != 500 || en.TxBytes != 600 {
		t.Errorf("en0 = %+v", en)
	}
}

func TestParseIfList2StopsOnANonsenseLength(t *testing.T) {
	buf := ifinfo2(1, 1500, 10, 20)
	buf = append(buf, otherMessage(96)[:20]...) // a message that claims 96 bytes and has 20
	got := ParseIfList2(buf)
	if len(got) != 1 {
		t.Errorf("the whole message before the broken one must survive, got %d", len(got))
	}

	zero := make([]byte, 8) // ifm_msglen = 0 would loop forever if not refused
	if got := ParseIfList2(zero); len(got) != 0 {
		t.Errorf("a zero-length message must stop the walk, got %+v", got)
	}
}

func TestVerifyIfList(t *testing.T) {
	got := ParseIfList2(ifinfo2(14, 1500, 1, 2))
	if !VerifyIfList(got, map[int]int{14: 1500}) {
		t.Error("the MTU the system reports must confirm the layout")
	}
	if VerifyIfList(got, map[int]int{14: 9000}) {
		t.Error("a disagreeing MTU means a wrong layout and must not verify")
	}
	if !VerifyIfList(got, map[int]int{14: 1500, 1: 16384}) {
		t.Error("an interface that came and went between the two reads must not throw away the rest")
	}
	if VerifyIfList(got, map[int]int{1: 16384}) {
		t.Error("with no interface in common there is nothing to confirm the layout")
	}
	if VerifyIfList(nil, map[int]int{14: 1500}) {
		t.Error("an empty answer verifies nothing")
	}
	if VerifyIfList(got, nil) {
		t.Error("with nothing to check against, nothing is verified")
	}
}
