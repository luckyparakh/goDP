package zerovalue

// GOAL
//
// Write Tally, a counter whose zero value is ready to use. This must work with
// no constructor and no setup:
//
//	var t Tally
//	t.Inc("a")
//
// It must also be safe when many goroutines call Inc at the same time. Decide
// which fields you need. Do not add a NewTally.
//
// You are done when every test in this folder passes, including with the race
// detector:
//
//	go test -race ./patterns/a10-zero-value/

// Tally counts how often each key was seen.
type Tally struct {
	// TODO: add the fields you need
}

// Inc adds one to the count for key.
func (t *Tally) Inc(key string) {
	// TODO
}

// Count returns how many times key was seen. An unknown key counts as 0.
func (t *Tally) Count(key string) int {
	// TODO
	return 0
}
