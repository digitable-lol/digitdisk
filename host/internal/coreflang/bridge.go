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
	"sync"

	"digitdisk/internal/core"

	"flangprogram/flang"
	rt "flangprogram/flangrt"
)

// Bridge implements core.Decider on top of the flang module.
type Bridge struct {
	ctx *rt.Ctx

	mu       sync.Mutex
	failures int
	firstErr error
}

// New returns a Bridge with the evaluation context the flang module asks for.
func New() *Bridge { return &Bridge{ctx: flang.NewContext()} }

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

	answer, err := flang.ReshitNahodku(b.ctx, nahodka)
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
