// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package core

import "testing"

func TestStubDecidesNothing(t *testing.T) {
	var d Decider = Default()
	if d.Ready() {
		t.Errorf("the placeholder must not claim to be a working decision layer")
	}
	for _, r := range []Record{
		{Path: "/var/cache/apt", Size: 1 << 30, AgeDays: 400, Kind: KindDir, Accessible: true},
		{Path: "/var/log/syslog", Size: 1 << 20, AgeDays: 1, Kind: KindFile, Accessible: true},
		{Path: "/x", Kind: KindLink, Accessible: false},
	} {
		got := d.Decide(r)
		want := Decision{Class: ClassUnknown, Verdict: VerdictKeep, Weight: 0}
		if got != want {
			t.Errorf("Decide(%q) = %+v, want %+v — guessing here would both mislead the operator and pre-empt the flang layer", r.Path, got, want)
		}
	}
}

func TestContractEnumsAreComplete(t *testing.T) {
	if ContractVersion != 0 {
		t.Errorf("ContractVersion = %d, want 0", ContractVersion)
	}
	if len(Classes) != 6 {
		t.Errorf("разрядов %d, want 6", len(Classes))
	}
	if len(Verdicts) != 3 {
		t.Errorf("приговоров %d, want 3", len(Verdicts))
	}
	seen := map[string]bool{}
	for _, c := range Classes {
		if seen[string(c)] {
			t.Errorf("duplicate разряд %q", c)
		}
		seen[string(c)] = true
	}
	for _, want := range []Class{ClassCache, ClassLog, ClassBuild, ClassDownload, ClassLarge, ClassUnknown} {
		if !seen[string(want)] {
			t.Errorf("разряд %q missing from Classes", want)
		}
	}
}
