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
//
// Version 1 added the справочник известных мест: the decision layer is handed
// a list of Place values alongside each record.  A layer that ignores it
// decides exactly what version 0 decided — the справочник only ever adds a
// разряд where the general приметы found none.
const ContractVersion = 1

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

// Thresholder is an optional capability a decision layer may also offer: the
// порог in days it applied to a разряд.  The host asks for it only to print
// why a verdict came out the way it did — "возраст 41 дн ≥ порог 7 дн" — so
// that the number in the explanation comes from the layer that decided, not
// from a copy of it kept in the host.  A host that copied the thresholds would
// go on printing 7 after the rule changed to 5, and the explanation would be a
// lie about a decision that was still correct.
//
// A Decider need not implement it; a host that finds it missing says порог "—"
// rather than guessing.
type Thresholder interface {
	// Threshold returns the age threshold in days the layer uses for a
	// разряд, and whether the разряд has one at all.
	Threshold(Class) (float64, bool)
}

// Sizer is another optional capability, and it exists for the same reason
// Thresholder does: a number the host shows must be the layer's number.
//
// «Порог крупного» is the size at which the decision layer stops calling a
// file ordinary.  The host asks for it when it has to decide HOW HARD to ask
// before doing something irreversible — see the забой key of the walk screen,
// which confirms a small erasure with one key and a large one only by count.
// "Large" there means large by the layer's reckoning, not by a number typed
// into the screen: a screen holding its own gigabyte would go on calling
// 900 МиБ small after the rule said otherwise.
//
// A Decider need not implement it.  A host that finds it missing does not
// guess a size — it takes the strict road every time, which is the safe answer
// to "I do not know how large is large".
type Sizer interface {
	// LargeBytes returns the size in bytes at which the layer calls a file
	// «Крупное», and whether the layer has such a threshold at all.
	LargeBytes() (int64, bool)
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

// Anchor is «якорь» — how a place's chain must sit in a path.
type Anchor string

const (
	// AnchorRoot: the chain must be the beginning of the path.
	AnchorRoot Anchor = "ОтКорня"
	// AnchorAnywhere: the chain may appear at any depth.
	AnchorAnywhere Anchor = "ГдеУгодно"
)

// Place is one entry of the справочник известных мест as the decision layer
// takes it: a разряд, an anchor, and a chain — the place's path with a slash
// at both ends, `/home/u/.npm/_cacache/`.
//
// The slashes are not decoration.  They are what makes the layer's comparison
// a comparison of whole path components rather than of substrings: the trail
// `/home/u/x.npm/_cacache/` does not contain `/home/u/.npm/_cacache/`, because
// there is no slash in front of `.npm`.  A chain without them is refused by
// the layer («Справочник ограничен»), not matched loosely.
//
// Only Кэш, Журнал, Сборка and Загрузка may appear here.  Крупное is decided
// by size and Неизвестное means "no place matched"; a справочник asserting
// either would be lying about where the разряд came from, and the layer's
// постусловие «Место обосновано» rejects it.
type Place struct {
	Class  Class  `json:"разряд"`
	Anchor Anchor `json:"якорь"`
	Chain  string `json:"цепь"`
}

// Placer is an optional capability of a decision layer: it accepts a
// справочник известных мест, once, before any record is judged.  A Decider
// that does not implement it judges by приметы alone, which is what the layer
// did before the справочник existed — so a missing capability is a smaller
// answer, never a wrong one.
type Placer interface {
	// UsePlaces hands the layer the справочник.  It returns an error when
	// the layer refuses it — a chain without its bounding slashes, say —
	// and the host then reports the refusal instead of cleaning with a
	// справочник the layer would not check.
	UsePlaces([]Place) error
}
