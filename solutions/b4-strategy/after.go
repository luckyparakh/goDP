//go:build ignore

package strategy

import (
	"fmt"
	"time"
)

// HOW IT WORKS
//
// A Backoff is a plain func from the retry number to a wait. Constant, Linear
// and Exponential are constructors: each returns a closure that has captured
// its own numbers. Retry holds whichever Backoff it was given and calls it.
// It has no switch and no list of names, so it cannot tell a built-in Backoff
// from one written in a test file.
//
// WHY IT HELPS
//
// A new policy is a new func, anywhere, and Retry does not change. A typo in
// a name is now a compile error (Exponental is undefined), not a silent wrong
// wait. And each policy can take its own numbers instead of sharing one
// hard-coded base.
//
// The choice has moved to whoever builds the call to Retry. In a service that
// is main (B5): if the policy comes from config, main turns the name into a
// Backoff once, and nothing below ever sees the string.
//
// WHEN NOT TO USE
//
// Two fixed options picked in one place read better as a switch you can see
// top to bottom. And this is not A2's registry: there, a name is looked up in
// a global table. Use that when the choice really must arrive as a name and
// new options must plug in without editing the code that does the lookup.

// Backoff says how long to wait before retry number attempt.
type Backoff func(attempt int) time.Duration

// Constant always waits d.
func Constant(d time.Duration) Backoff {
	return func(int) time.Duration { return d }
}

// Linear waits step times the retry number.
func Linear(step time.Duration) Backoff {
	return func(attempt int) time.Duration { return step * time.Duration(attempt) }
}

// Exponential waits base, then doubles each retry, but never more than limit.
func Exponential(base, limit time.Duration) Backoff {
	return func(attempt int) time.Duration {
		d := base
		for i := 1; i < attempt && d < limit; i++ {
			d *= 2
		}
		return min(d, limit)
	}
}

// Retry runs op up to attempts times, sleeping b(n) between attempts.
func Retry(attempts int, b Backoff, sleep func(time.Duration), op func() error) error {
	if attempts < 1 {
		return ErrNoAttempts
	}
	var err error
	for n := 1; n <= attempts; n++ {
		if err = op(); err == nil {
			return nil
		}
		if n < attempts {
			sleep(b(n))
		}
	}
	return fmt.Errorf("after %d attempts: %w", attempts, err)
}
