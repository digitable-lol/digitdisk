// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

// Package core holds the data contract between the Go host (fact gathering,
// CLI) and the decision layer written in flang and emitted as a Go package.
//
// Contract version 0.  One record in, one verdict out:
//
//	in:  путь (строка), размер (число байт), возраст_дней (число),
//	     вид (Файл|Каталог|Ссылка), доступен (истина|ложь)
//	out: разряд (Кэш|Журнал|Сборка|Загрузка|Крупное|Неизвестное),
//	     приговор (МожноУбрать|Спросить|НеТрогать), вес (число)
//
// The enum values are part of the contract and are therefore carried as the
// literal Russian strings the decision layer produces.  Go identifiers stay in
// English.
//
// WIRING THE REAL DECIDER: when the flang layer is emitted into
// digitdisk/core/out-go, add one file to this package that implements Decider
// by delegating to that package and converting its record/verdict types to the
// ones below.  Nothing else in the host needs to change: the host only ever
// talks to the Decider interface.
package core

// ContractVersion is the version of the record/verdict shape above.
const ContractVersion = 0

// Kind is "вид" — what the path is on disk.
type Kind string

const (
	KindFile Kind = "Файл"
	KindDir  Kind = "Каталог"
	KindLink Kind = "Ссылка"
)

// Class is "разряд" — the bucket the decision layer sorts a record into.
type Class string

const (
	ClassCache    Class = "Кэш"
	ClassLog      Class = "Журнал"
	ClassBuild    Class = "Сборка"
	ClassDownload Class = "Загрузка"
	ClassLarge    Class = "Крупное"
	ClassUnknown  Class = "Неизвестное"
)

// Classes lists every разряд in report order.
var Classes = []Class{ClassCache, ClassLog, ClassBuild, ClassDownload, ClassLarge, ClassUnknown}

// Verdict is "приговор" — what the decision layer advises for a record.
type Verdict string

const (
	VerdictRemovable Verdict = "МожноУбрать"
	VerdictAsk       Verdict = "Спросить"
	VerdictKeep      Verdict = "НеТрогать"
)

// Verdicts lists every приговор in report order.
var Verdicts = []Verdict{VerdictRemovable, VerdictAsk, VerdictKeep}

// Record is the input side of the contract: one path observed by the host.
type Record struct {
	Path       string  `json:"путь"`
	Size       int64   `json:"размер"`
	AgeDays    float64 `json:"возраст_дней"`
	Kind       Kind    `json:"вид"`
	Accessible bool    `json:"доступен"`
}

// Decision is the output side of the contract: what the decision layer says
// about one record.
type Decision struct {
	Class   Class   `json:"разряд"`
	Verdict Verdict `json:"приговор"`
	Weight  float64 `json:"вес"`
}

// Decider is the whole surface the host requires of the decision layer.
type Decider interface {
	// Decide classifies one record.  It must be safe for repeated calls and
	// must not touch the filesystem: the host has already done the looking.
	Decide(Record) Decision
	// Name identifies the implementation in reports.
	Name() string
	// Ready reports whether a real decision layer is behind this Decider.
	// A false value tells the report to say so out loud instead of passing
	// off empty answers as analysis.
	Ready() bool
}
