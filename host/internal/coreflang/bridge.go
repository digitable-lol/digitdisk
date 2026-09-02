// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

//go:build flangcore

// Package coreflang binds the host to the decision layer written in flang and
// printed as a Go package into digitdisk/core/out-go.
//
// It is behind the "flangcore" build tag on purpose: the generated package is
// a separate Go module owned by another worker and regenerated from the flang
// source, so the host's default build must not depend on it being present.
//
// To build with the real decision layer, add these two lines to host/go.mod:
//
//	require flangprogram v0.0.0
//	replace flangprogram => ../core/out-go
//
// and build with: go build -tags flangcore -o digitdisk .
//
// The mapping is an identity on the wire: the variant names the flang module
// produces ("Файл", "Кэш", "МожноУбрать", …) are exactly the string values of
// the contract constants in package core, so nothing is translated, only
// checked.
package coreflang

import (
	"fmt"
	"sync"

	"digitdisk/internal/core"

	"flangprogram/flang"
	rt "flangprogram/flangrt"
)

// Bridge implements core.Decider on top of the flang module.
type Bridge struct {
	ctx *rt.Ctx

	// places is the справочник известных мест as flang values, built once
	// and handed to every call.  Building it per record would rebuild a
	// hundred records for every file on the disk.
	places rt.Value

	mu       sync.Mutex
	failures int
	firstErr error
}

// New returns a Bridge with the evaluation context the flang module asks for
// and an empty справочник — until UsePlaces is called it judges by приметы
// alone, which is exactly what contract version 0 did.
func New() *Bridge {
	return &Bridge{ctx: flang.NewContext(), places: rt.List(nil)}
}

// UsePlaces implements core.Placer.  The справочник is checked by the layer
// itself before it is kept: «Справочник ограничен» asserts that every цепь is
// bounded by slashes at both ends, which is the property that makes the
// comparison a comparison of components.  A справочник that fails it is
// refused whole — matching a chain without its slashes is the bug of
// 1 September all over again, and half a справочник is not better than none.
func (b *Bridge) UsePlaces(places []core.Place) error {
	values := make([]rt.Value, 0, len(places))
	for _, p := range places {
		class, ok := classVariant(p.Class)
		if !ok {
			return fmt.Errorf("справочник: разряд %q решающему слою не известен", p.Class)
		}
		anchor := flang.VariantOtKornya()
		switch p.Anchor {
		case core.AnchorRoot:
		case core.AnchorAnywhere:
			anchor = flang.VariantGdeUgodno()
		default:
			return fmt.Errorf("справочник: якорь %q решающему слою не известен", p.Anchor)
		}
		values = append(values, flang.SozdatMesto(class, anchor, rt.Text(p.Chain)))
	}
	list := rt.List(values)

	ok, err := flang.SpravochnikOgranichen(b.ctx, list)
	if err != nil {
		return fmt.Errorf("решающий слой не смог проверить справочник: %w", err)
	}
	if !ok.Flag {
		return fmt.Errorf("решающий слой отверг справочник: у какого-то места цепь не ограничена косыми с обеих сторон, " +
			"а без них сверка перестала бы быть сверкой по составляющим")
	}
	b.places = list
	return nil
}

// Name implements core.Decider.
func (b *Bridge) Name() string { return "flang «Опись диска» (core/out-go)" }

// Ready implements core.Decider.
func (b *Bridge) Ready() bool { return true }

// Failures reports how many records the decision layer refused, and the first
// reason.  A refusal is a fact to show, not something to hide behind a zero.
func (b *Bridge) Failures() (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures, b.firstErr
}

// classVariant maps разряд onto the sum type of the flang module.
func classVariant(c core.Class) (rt.Value, bool) {
	switch c {
	case core.ClassCache:
		return flang.VariantKesh(), true
	case core.ClassLog:
		return flang.VariantZhurnal(), true
	case core.ClassBuild:
		return flang.VariantSborka(), true
	case core.ClassDownload:
		return flang.VariantZagruzka(), true
	case core.ClassLarge:
		return flang.VariantKrupnoe(), true
	case core.ClassUnknown:
		return flang.VariantNeizvestnoe(), true
	}
	return rt.Nothing(), false
}

// Threshold implements core.Thresholder by asking «Порог разряда» — the same
// function the verdict rule itself calls.  The host therefore prints the
// number the decision was made with, not a copy of it.  «Крупное» and
// «Неизвестное» answer 0 there because they are not judged by age at all, and
// that is reported as "no threshold" rather than as a threshold of zero, which
// would read as "anything of any age".
func (b *Bridge) Threshold(c core.Class) (float64, bool) {
	variant, ok := classVariant(c)
	if !ok {
		return 0, false
	}
	v, err := flang.PorogRazryada(b.ctx, variant)
	if err != nil {
		b.note(err)
		return 0, false
	}
	if v.Num <= 0 {
		return 0, false
	}
	return v.Num, true
}

// kindVariant maps вид onto the sum type of the flang module.
func kindVariant(k core.Kind) rt.Value {
	switch k {
	case core.KindDir:
		return flang.VariantKatalog()
	case core.KindLink:
		return flang.VariantSsylka()
	default:
		return flang.VariantFayl()
	}
}

// Decide implements core.Decider by handing one Находка to «Решить находку».
func (b *Bridge) Decide(r core.Record) core.Decision {
	nahodka := flang.SozdatNahodka(
		rt.Text(r.Path),
		rt.Number(float64(r.Size)),
		rt.Number(r.AgeDays),
		kindVariant(r.Kind),
		rt.Flag(r.Accessible),
	)

	answer, err := flang.ReshitNahodku(b.ctx, nahodka, b.places)
	if err != nil {
		b.note(err)
		return core.Decision{Class: core.ClassUnknown, Verdict: core.VerdictKeep}
	}

	out := core.Decision{Class: core.ClassUnknown, Verdict: core.VerdictKeep}
	if v, err := rt.FieldGet(b.ctx, answer, "разряд"); err == nil {
		for _, c := range core.Classes {
			if rt.VariantIs(v, string(c)) {
				out.Class = c
				break
			}
		}
	} else {
		b.note(err)
	}
	if v, err := rt.FieldGet(b.ctx, answer, "приговор"); err == nil {
		for _, verdict := range core.Verdicts {
			if rt.VariantIs(v, string(verdict)) {
				out.Verdict = verdict
				break
			}
		}
	} else {
		b.note(err)
	}
	if v, err := rt.FieldGet(b.ctx, answer, "вес"); err == nil {
		out.Weight = v.Num
	} else {
		b.note(err)
	}
	return out
}

func (b *Bridge) note(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.firstErr == nil {
		b.firstErr = err
	}
}
