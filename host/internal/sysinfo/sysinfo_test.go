// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package sysinfo

import "testing"

// The tests that hold on every system.  The ones that build a fake /proc are
// in sysinfo_linux_test.go, next to the collector that reads it.

func TestHumanDuration(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, ""},
		{-5, ""},
		{61, "00:01"},
		{3600, "01:00"},
		{90061, "1д 01:01"},
		{2100448, "24д 07:27"},
	}
	for _, c := range cases {
		if got := HumanDuration(c.in); got != c.want {
			t.Errorf("HumanDuration(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestUnmeasuredNamesTheFactAndTheReason(t *testing.T) {
	st := Status{Missing: map[string]string{FactCPUBusy: "нужен Mach"}}
	why, ok := st.Unmeasured(FactCPUBusy)
	if !ok || why != "нужен Mach" {
		t.Errorf("Unmeasured(%q) = %q, %v", FactCPUBusy, why, ok)
	}
	if _, ok := st.Unmeasured(FactSensors); ok {
		t.Errorf("a fact nobody reported missing must not read as missing")
	}
	if _, ok := (Status{}).Unmeasured(FactCPUBusy); ok {
		t.Errorf("an empty snapshot reports no facts as unmeasured")
	}
}
