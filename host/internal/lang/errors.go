// SPDX-FileCopyrightText: 2026 Marat Zimnurov <zimtir@mail.ru>
// SPDX-License-Identifier: BSD-2-Clause

package lang

import "errors"

// An Error is a refusal that has to be readable in either language.
//
// «нужен ровно один путь для обхода, получено 2» is printed by main and by
// nothing else, but it is built four levels down, where nobody knows what
// language the reader has.  So it travels as a Phrase and is rendered at the
// place it is printed — the same trick Phrase plays for the JSON, for the same
// reason: the text is made in one place and read in another.
//
// It is an error and behaves like one: errors.Is and errors.As see through it
// to whatever it wraps, so the checks the tool already makes on os.IsNotExist
// and friends go on working.
type Error struct {
	P    Phrase
	Wrap error
}

// Errorf builds an Error from a wording of the dictionary.  An argument that
// is itself an error is remembered as the wrapped one, so `errors.Is` keeps
// working through a translated message.
func Errorf(ru string, a ...any) error {
	e := &Error{P: Say(ru, a...)}
	for _, arg := range a {
		if err, ok := arg.(error); ok {
			e.Wrap = err
			break
		}
	}
	return e
}

// Error is the Russian rendering: this is what a log, a journal or a test that
// compares error text has always seen.
func (e *Error) Error() string { return e.P.String() }

// Unwrap gives errors.Is and errors.As the error this one was built around.
func (e *Error) Unwrap() error { return e.Wrap }

// In renders this error in a language.
func (e *Error) In(l Lang) string { return e.P.In(l) }

// InLang renders err for a reader of l: the translated wording when the error
// is one of ours, the system's own text when it is not.  A message from the
// operating system is left alone on purpose — an English sentence no strerror
// ever produced would be harder to search for than the real one.
func InLang(err error, l Lang) string {
	if err == nil {
		return ""
	}
	var ours *Error
	if errors.As(err, &ours) {
		return ours.In(l)
	}
	return err.Error()
}

// verbs is used by the completeness check and by nothing else: it is here
// rather than in the test so that the rule about %-verbs is stated where the
// formatting happens.
func verbs(format string) []string {
	var out []string
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		if i+1 < len(format) && format[i+1] == '%' {
			i++
			continue
		}
		j := i + 1
		for j < len(format) && !isVerbLetter(format[j]) {
			j++
		}
		if j < len(format) {
			out = append(out, format[i:j+1])
			i = j
		}
	}
	return out
}

func isVerbLetter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

// Verbs is verbs, for the check in complete_test.go.
func Verbs(format string) []string { return verbs(format) }
