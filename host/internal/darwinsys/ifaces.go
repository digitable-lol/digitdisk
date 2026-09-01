// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

// Constants of the routing-socket interface list, from <net/route.h> and
// <net/if_var.h>.
//
// WHY NOT getifaddrs(3).  getifaddrs is the documented way to reach an
// interface's struct if_data, and it is a libc function: calling it needs cgo,
// which this tree will not take on for a fact it can read directly.  getifaddrs
// itself reads the routing socket, and so do we: sysctl
// {CTL_NET, AF_ROUTE, 0, 0, NET_RT_IFLIST2, 0} answers with a stream of
// messages, and the RTM_IFINFO2 ones carry struct if_msghdr2 with the 64-bit
// struct if_data64 inside.  The "2" matters: the older NET_RT_IFLIST carries
// struct if_data, whose byte counters are 32 bits and wrap every 4 GiB — on a
// machine that has been up a week those are not the numbers anybody wants.
const (
	// RTM_VERSION, the routing message version this decoder understands.
	rtmVersion = 5
	// RTM_IFINFO2, the message type NET_RT_IFLIST2 adds.
	rtmIfInfo2 = 0x12

	// sizeof(struct if_msghdr2): 32 bytes of header, then if_data64.
	ifMsghdr2Size = 160
	// sizeof(struct if_data64).
	ifData64Size = 128

	// Offsets inside if_msghdr2.
	offMsglen = 0  // u_short ifm_msglen
	offVer    = 2  // u_char  ifm_version
	offType   = 3  // u_char  ifm_type
	offIndex  = 12 // u_short ifm_index
	offData   = 32 // struct if_data64 ifm_data

	// Offsets inside if_data64, relative to offData.
	offMTU        = offData + 8   // u_int32_t ifi_mtu
	offIPackets   = offData + 24  // u_int64_t ifi_ipackets
	offIErrors    = offData + 32  // u_int64_t ifi_ierrors
	offOPackets   = offData + 40  // u_int64_t ifi_opackets
	offOErrors    = offData + 48  // u_int64_t ifi_oerrors
	offIBytes     = offData + 64  // u_int64_t ifi_ibytes
	offOBytes     = offData + 72  // u_int64_t ifi_obytes
	offIQDrops    = offData + 96  // u_int64_t ifi_iqdrops
	offLastChange = offData + 120 // struct timeval32 ifi_lastchange
)

// IfCounters is one interface's counters out of struct if_data64.
//
// There is no outbound-drop counter: if_data64 ends at ifi_lastchange and has
// none, so the host reports transmit drops as unmeasured instead of zero.
type IfCounters struct {
	// Index is ifm_index, which is also the index net.Interface carries,
	// and is how the two halves of an interface are joined.
	Index int
	// MTU is ifi_mtu.  It is decoded for one reason: the caller knows the
	// MTU of every interface from the standard library already, so a
	// disagreement here proves the layout above is wrong.
	MTU       uint32
	RxBytes   uint64
	RxPackets uint64
	RxErrors  uint64
	RxDropped uint64
	TxBytes   uint64
	TxPackets uint64
	TxErrors  uint64
}

// ParseIfList2 walks the NET_RT_IFLIST2 answer and returns the counters of
// every interface in it, keyed by interface index.
//
// The walk is driven by each message's own ifm_msglen, so message types this
// decoder does not know (addresses, multicast addresses) are stepped over
// rather than misread.  A message that claims a length shorter than a header
// or longer than what is left stops the walk: past that point nothing in the
// buffer can be trusted.
func ParseIfList2(b []byte) map[int]IfCounters {
	out := map[int]IfCounters{}
	for off := 0; off+4 <= len(b); {
		msglen := int(u16(b, off+offMsglen))
		if msglen < 4 || off+msglen > len(b) {
			return out
		}
		msg := b[off : off+msglen]
		if b[off+offVer] == rtmVersion && b[off+offType] == rtmIfInfo2 && msglen >= ifMsghdr2Size {
			idx := int(u16(msg, offIndex))
			out[idx] = IfCounters{
				Index:     idx,
				MTU:       u32(msg, offMTU),
				RxBytes:   u64(msg, offIBytes),
				RxPackets: u64(msg, offIPackets),
				RxErrors:  u64(msg, offIErrors),
				RxDropped: u64(msg, offIQDrops),
				TxBytes:   u64(msg, offOBytes),
				TxPackets: u64(msg, offOPackets),
				TxErrors:  u64(msg, offOErrors),
			}
		}
		off += msglen
	}
	return out
}

// VerifyIfList reports whether the decoded counters describe the interfaces
// the caller already knows: for every interface both sides carry, the MTU
// decoded at the offset above must equal the MTU the standard library reports
// for the same index, and at least one interface must be checked that way.
//
// This is the same bargain as Verify for processes.  A shifted if_data64 puts
// some other field where ifi_mtu belongs, and an interface whose MTU is not
// the one the system itself names is a decoder that must publish nothing.
//
// Only the interfaces both sides carry are compared, on purpose: an interface
// that appears or disappears between the two reads is an ordinary race, not
// evidence about the layout, and must not throw away every counter on the
// machine.
func VerifyIfList(counters map[int]IfCounters, mtus map[int]int) bool {
	checked := 0
	for idx, c := range counters {
		mtu, ok := mtus[idx]
		if !ok {
			continue
		}
		if int(c.MTU) != mtu {
			return false
		}
		checked++
	}
	return checked > 0
}
