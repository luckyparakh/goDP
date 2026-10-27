//go:build ignore

package lsp

import (
	"errors"
	"fmt"
)

// HOW IT WORKS
//
// Contract is the five promises written as code. It asks newStore for a store,
// does what a user would do, and records a line for each promise that does not
// hold. Any implementation can be passed in, so one checklist covers them all.
//
// WHY IT HELPS
//
// The compiler cannot tell the honest store from the lying one. A test run
// against every implementation can. It also stops a new implementation from
// quietly breaking what callers rely on.
//
// WHEN NOT TO USE
//
// A checklist only catches promises someone wrote down. A store can still break
// a sixth promise nobody listed. If there is only one implementation, a normal
// test is enough.

// Contract runs the Store promises against stores built by newStore and
// returns one line per broken promise.
func Contract(newStore func() Store) []string {
	var v []string

	// Promise 1: a missing key is ErrNotFound.
	if _, err := newStore().Get("missing"); !errors.Is(err, ErrNotFound) {
		v = append(v, fmt.Sprintf("1: Get(missing) error = %v, want ErrNotFound", err))
	}

	// Promise 2: Put succeeds.
	if err := newStore().Put("k", "v1"); err != nil {
		v = append(v, fmt.Sprintf("2: Put returned %v, want nil", err))
	}

	// Promise 3: Get returns what Put stored.
	s := newStore()
	if err := s.Put("k", "v1"); err == nil {
		if got, err := s.Get("k"); err != nil || got != "v1" {
			v = append(v, fmt.Sprintf("3: Get after Put = %q, %v, want v1, nil", got, err))
		}
	}

	// Promises 4 and 5: overwriting succeeds and the newest value wins.
	s = newStore()
	if err := s.Put("k", "v1"); err == nil {
		if err := s.Put("k", "v2"); err != nil {
			v = append(v, fmt.Sprintf("4: overwrite Put returned %v, want nil", err))
		} else if got, err := s.Get("k"); err != nil || got != "v2" {
			v = append(v, fmt.Sprintf("5: Get after overwrite = %q, %v, want v2, nil", got, err))
		}
	}

	return v
}
