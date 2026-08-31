// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package core

// Stub is the placeholder decision layer used until the flang-emitted package
// exists.  It deliberately decides nothing: every record comes back as
// Неизвестное / НеТрогать with zero вес.  Guessing here would both mislead the
// operator and pre-empt the layer being written in flang.
type Stub struct{}

// Decide implements Decider.
func (Stub) Decide(r Record) Decision {
	return Decision{Class: ClassUnknown, Verdict: VerdictKeep, Weight: 0}
}

// Name implements Decider.
func (Stub) Name() string {
	return "заглушка (решающий слой ещё не подключён)"
}

// Ready implements Decider.
func (Stub) Ready() bool { return false }

// Default returns the decision layer the host will use.  When the flang layer
// is wired in, this is the single place that changes.
func Default() Decider { return Stub{} }
