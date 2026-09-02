// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

import (
	"encoding/json"
	"fmt"
)

// A Phrase is a sentence that has to be two things at once: a value in the
// JSON, where it must not move a byte, and a line on the screen, where it must
// be in the reader's language.
//
// The refusals of `clean`, the notes of `places`, the reasons of `--why` are
// all of that kind: they are printed AND they are fields — «отказ»,
// «замечание», «не_сделано», «missing» — that scripts already read.  A Phrase
// keeps the Russian wording and its arguments, renders Russian for the machine
// and the chosen language for the person, and so lets the output be translated
// without the JSON changing at all.
//
// Errors from the system are a Phrase too, made by Raw: nobody wrote them,
// this tool cannot translate them honestly, and pretending otherwise would
// produce an English sentence that no `strerror` ever said.
type Phrase struct {
	ru   string
	args []any
	raw  string
}

// Say builds a Phrase from a wording of the dictionary and its arguments.
func Say(ru string, a ...any) Phrase { return Phrase{ru: ru, args: a} }

// Raw builds a Phrase from a text this tool did not write — the message of an
// error from the operating system, the line of a file a person edited.  It is
// shown as it is, in either language.
func Raw(s string) Phrase { return Phrase{raw: s} }

// FromError is Raw for an error, and the empty Phrase for a nil one.
func FromError(err error) Phrase {
	if err == nil {
		return Phrase{}
	}
	return Raw(err.Error())
}

// Empty reports whether nothing was said.
func (p Phrase) Empty() bool { return p.ru == "" && p.raw == "" }

// String is the Russian rendering: the wording as it was written, filled in.
// This is what goes into the JSON and into the журнал, and it is what this
// tool has always written there.
func (p Phrase) String() string {
	if p.raw != "" {
		return p.raw
	}
	if p.ru == "" {
		return ""
	}
	if len(p.args) == 0 {
		return p.ru
	}
	return fmt.Sprintf(p.ru, p.args...)
}

// In is the rendering for a reader.
//
// An argument that is itself one of ours — a Phrase, or an error built by
// Errorf — is rendered in the same language before it is put in.  Without
// that, a sentence built around another sentence («справочник %s: %s») would
// come out half translated, which is the one outcome this whole package
// exists to prevent.  Everything else is left exactly as it is: a path, a
// number, a message from the system.
func (p Phrase) In(l Lang) string {
	if p.raw != "" {
		return p.raw
	}
	if p.ru == "" {
		return ""
	}
	if len(p.args) == 0 {
		return l.T(p.ru)
	}
	return l.F(p.ru, translated(l, p.args)...)
}

// translated renders the arguments that are ours in l and leaves the rest
// alone.
func translated(l Lang, args []any) []any {
	out := make([]any, len(args))
	for i, a := range args {
		switch v := a.(type) {
		case Phrase:
			out[i] = v.In(l)
		case *Phrase:
			out[i] = v.In(l)
		case error:
			out[i] = InLang(v, l)
		default:
			out[i] = a
		}
	}
	return out
}

// MarshalJSON writes the Russian rendering, and nothing else ever.  The whole
// point of this type is that `--json` cannot tell which language the person
// running the command reads.
func (p Phrase) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

// UnmarshalJSON reads a Phrase back out of a журнал written earlier.
//
// A text that is a wording of the dictionary comes back as that wording and is
// translated again; a text that is not — a filled-in sentence, an error from
// the system, a line written by a version that spelled it differently — comes
// back raw and is shown as it was recorded.  A журнал is a record of what
// happened, and rewriting the record to suit the reader's language would make
// it a worse record.
func (p *Phrase) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s != "" && Known(s) {
		*p = Phrase{ru: s}
		return nil
	}
	*p = Phrase{raw: s}
	return nil
}

// Wording is the Russian key of this Phrase, empty for a Raw one.  Tests use
// it to prove that a refusal a person can be shown is a refusal the dictionary
// has a pair for.
func (p Phrase) Wording() string { return p.ru }

// IsZero tells encoding/json that a Phrase nobody filled in is the zero value.
//
// This is the whole of the promise that translating a refusal moves no byte of
// the JSON.  The fields that hold a Phrase — «замечание», «беда»,
// «не_сделано» — were strings with `omitempty`, and an empty string was left
// out of the object entirely.  A struct is never "empty" to `omitempty`, so
// without this the same fields would start appearing as `""` in every record
// that has nothing to say, and every script that tests for their presence
// would start finding them.  With `omitzero` on the tag and this method, they
// go on being absent.
func (p Phrase) IsZero() bool { return p.Empty() }
