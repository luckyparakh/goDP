package lsp

// GOAL
//
// Write Contract so it checks the five promises documented on Store in
// before.go. Call newStore to get a store to test. Calling it more than once
// is allowed.
//
// It must return one text line per broken promise, and an empty slice when
// every promise is kept. It returns lines instead of calling t.Error so a test
// can say "this bad store must be caught" without shipping a failing test.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a3-lsp/

// Contract runs the Store promises against stores built by newStore and
// returns one line per broken promise.
func Contract(newStore func() Store) []string {
	// TODO
	return nil
}
