//go:build ignore

package embedding

import "strings"

// HOW IT WORKS
//
// Greet does not rely on Base's Hello. It takes a Namer and asks that Namer
// for the name, so the thing that varies is passed in explicitly. Shout
// embeds a Namer interface: every method of the inner value is forwarded
// automatically, and Shout writes only the one method it wants to change.
//
// WHY IT HELPS
//
// There is no hidden "which Name does Hello call?" question. The behaviour you
// get is the behaviour you pass in, which is easy to read and easy to test.
// Embedding still removes the boilerplate of forwarding many methods by hand.
//
// WHEN NOT TO USE
//
// Do not embed to get an "is-a" relationship. Derived is not a Base and cannot
// be passed where a Base is expected, because the compiler rejects it. Embedding
// also promotes every exported method of the inner type into your API, so
// embed only when you want all of them to be visible.

// Greet builds a greeting from the name of n.
func Greet(n Namer) string {
	return "hello from " + n.Name()
}

// Shout wraps a Namer and shouts its name. Every other method is forwarded.
type Shout struct{ Namer }

// Name returns the inner name in upper case.
func (s Shout) Name() string {
	return strings.ToUpper(s.Namer.Name())
}
