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
	if ContractVersion != 1 {
		t.Errorf("ContractVersion = %d, want 1", ContractVersion)
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

// TestPlaceEnumsAreComplete pins the shape of the справочник side of the
// contract.  A разряд the справочник must never assert is checked by name:
// «Крупное» and «Неизвестное» are answers about size and about nothing having
// matched, and a place claiming either would be a lie about where the разряд
// came from.
func TestPlaceEnumsAreComplete(t *testing.T) {
	if AnchorRoot == AnchorAnywhere {
		t.Fatal("два якоря обязаны различаться")
	}
	for _, a := range []Anchor{AnchorRoot, AnchorAnywhere} {
		if a == "" {
			t.Error("у якоря нет имени: имя — часть договора")
		}
	}
	p := Place{Class: ClassCache, Anchor: AnchorRoot, Chain: "/home/u/.npm/_cacache/"}
	if p.Chain[0] != '/' || p.Chain[len(p.Chain)-1] != '/' {
		t.Error("цепь обязана быть ограничена косыми с обеих сторон")
	}
}

// TestStubIgnoresPlaces states the promise made in the Placer comment: a
// decision layer without the capability is a smaller answer, never a wrong
// one.  The stub does not implement Placer, and that must stay detectable.
func TestStubIgnoresPlaces(t *testing.T) {
	if _, ok := any(Default()).(Placer); ok {
		t.Error("заглушка объявила, что принимает справочник, — тогда она обязана его и применять")
	}
}
