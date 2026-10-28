package lsp

import (
	"errors"
	"fmt"
)

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
	var errs []string
	store1, store2 := newStore(), newStore()

	// Promise 2
	if err := store1.Put("k", "v"); err != nil {
		errs = append(errs, fmt.Sprintf("2. put error %s", err.Error()))
	}

	// Promise 3
	got, getErr := store1.Get("k")
	if got != "v" {
		errs = append(errs, fmt.Sprintf("3. got: %s, expected v", got))
	}
	if getErr != nil {
		errs = append(errs, fmt.Sprintf("3. get error %s", getErr.Error()))
	}

	// Promise 1
	if _, err := newStore().Get("no"); !errors.Is(err, ErrNotFound) {
		errs = append(errs, "1. no key found error")
	}

	// Promise 4 and 5
	if err := store2.Put("k", "v1"); err != nil {
		errs = append(errs, fmt.Sprintf("4. put error %s", err.Error()))
	}
	if err := store2.Put("k", "v2"); err != nil {
		errs = append(errs, fmt.Sprintf("4. put error %s", err.Error()))
	}
	got, getErr = store2.Get("k")
	if got != "v2" {
		errs = append(errs, fmt.Sprintf("5. got: %s, expected v2", got))
	}
	if getErr != nil {
		errs = append(errs, fmt.Sprintf("5. get error %s", getErr.Error()))
	}

	return errs
}
