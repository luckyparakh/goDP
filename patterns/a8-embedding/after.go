package embedding

// GOAL
//
// Get the behaviour you wanted from inheritance, the Go way: pass the thing
// that varies in, and use embedding only to forward.
//
//   - Greet returns "hello from " plus the name of the Namer it is given. It
//     works with any Namer: Base, Derived or Shout.
//   - Shout embeds a Namer, but its own Name method returns the inner name in
//     upper case. All other methods of the inner value are forwarded for free.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a8-embedding/

// Greet builds a greeting from the name of n.
func Greet(n Namer) string {
	// TODO
	return ""
}

// Shout wraps a Namer and shouts its name. It currently forwards Name
// unchanged. Add a Name method to change that.
type Shout struct{ Namer }
