//go:build ignore

package dip

import "fmt"

// HOW IT WORKS
//
// Greeter does not know where names come from. It asks one small question
// through the Users interface, which Greeter itself declares. A file reader, a
// database or a plain map can answer it. The arrow now points up: the supplier
// fits Greeter's need, not the other way round.
//
// WHY IT HELPS
//
// A test hands in a map, with no folder to create and no cleanup. A failure
// is one line: a stand-in that returns the error you choose. Before, that test
// was only possible by breaking a real file.
//
// WHEN NOT TO USE
//
// If the file reader is the only source forever and you never test the failure
// path, the interface is extra machinery. Add it when you need a second source
// or a failure you cannot trigger for real.

// Users is what Greeter needs from the outside world.
type Users interface {
	Name(id string) (string, error)
}

// Greeter greets users fetched through Users.
type Greeter struct{ users Users }

// NewGreeter wires a Greeter to a source of users.
func NewGreeter(u Users) *Greeter { return &Greeter{users: u} }

// Greeting returns "Hello, <name>!" for the user with the given id.
func (g *Greeter) Greeting(id string) (string, error) {
	name, err := g.users.Name(id)
	if err != nil {
		return "", fmt.Errorf("greeting %q: %w", id, err)
	}
	return "Hello, " + name + "!", nil
}
