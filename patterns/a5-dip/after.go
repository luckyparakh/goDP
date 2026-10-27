package dip

// GOAL
//
// Greeter does the same job as FileService, but this file must not import
// "os" or "path/filepath". The interface below belongs to Greeter, the code
// that needs the answer, and not to whoever supplies users. Fill in its one
// method, keep the error wrapping as in before.go, and make NewGreeter use its
// argument.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a5-dip/
//
// Check the rule about imports with:
//
//	grep -nE '"os"|"path/filepath"' patterns/a5-dip/after.go
//
// (In PowerShell: Select-String -Path patterns/a5-dip/after.go -Pattern '"os"')

// Users is what Greeter needs from the outside world.
type Users interface {
	// TODO
}

// Greeter greets users fetched through Users.
type Greeter struct{ users Users }

// NewGreeter wires a Greeter to a source of users.
func NewGreeter(u Users) *Greeter {
	// TODO
	return &Greeter{}
}

// Greeting returns "Hello, <name>!" for the user with the given id.
func (g *Greeter) Greeting(id string) (string, error) {
	// TODO
	return "", nil
}
