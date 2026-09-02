// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package darwinsys

import "testing"

// ticksBytes writes one struct processor_cpu_load_info the way the kernel
// does: four little-endian natural_t.
func ticksBytes(user, system, idle, nice uint32) []byte {
	out := make([]byte, CPUTicksSize)
	for i, v := range []uint32{user, system, idle, nice} {
		out[4*i] = byte(v)
		out[4*i+1] = byte(v >> 8)
		out[4*i+2] = byte(v >> 16)
		out[4*i+3] = byte(v >> 24)
	}
	return out
}

func TestProcessorTicksSplitPerProcessor(t *testing.T) {
	b := append(ticksBytes(10, 1, 100, 0), ticksBytes(20, 2, 200, 0)...)
	got, ok := ParseProcessorTicks(b, 2)
	if !ok {
		t.Fatal("два процессора не разобрались")
	}
	if len(got) != 2 {
		t.Fatalf("процессоров %d", len(got))
	}
	if got[0][0] != 10 || got[1][0] != 20 {
		t.Errorf("счётчики перепутаны: %v", got)
	}
	// The share is a ratio of two differences, so it comes out the same for
	// one processor as for the machine.
	later, _ := ParseProcessorTicks(append(ticksBytes(20, 1, 100, 0), ticksBytes(20, 2, 400, 0)...), 2)
	if share, ok := got[0].BusyShare(later[0]); !ok || share != 100 {
		t.Errorf("доля первого ядра %v (ok=%v), ждали 100", share, ok)
	}
	if share, ok := got[1].BusyShare(later[1]); !ok || share != 0 {
		t.Errorf("доля второго ядра %v (ok=%v), ждали 0", share, ok)
	}
}

// An array that is not the length the processor count implies is refused: the
// alternative is reading every processor after the first at the wrong offset.
func TestProcessorTicksRefuseTheWrongLength(t *testing.T) {
	b := append(ticksBytes(1, 1, 1, 1), ticksBytes(2, 2, 2, 2)...)
	for _, cpus := range []int{0, 1, 3, -1} {
		if _, ok := ParseProcessorTicks(b, cpus); ok {
			t.Errorf("массив на два процессора принят за %d", cpus)
		}
	}
	if _, ok := ParseProcessorTicks(b[:len(b)-1], 2); ok {
		t.Error("укороченный массив принят")
	}
}
