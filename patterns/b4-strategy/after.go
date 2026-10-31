package strategy

import (
	"errors"
	"time"
)

// GOAL
//
// Make the wait a value the caller is handed, not a name it looks up.
//
//   - A Backoff is a plain func: given the retry number (1 for the first
//     retry), it returns how long to wait.
//   - Constant, Linear and Exponential build the three Backoffs from before.go.
//     Each takes its numbers as arguments, so nothing is hard-coded inside.
//     Exponential doubles from base and never goes past limit.
//   - Retry behaves like RetryBefore, but takes a Backoff instead of a policy
//     name. It must never mention any particular Backoff.
//   - Keep the errors as in RetryBefore: ErrNoAttempts for fewer than one
//     attempt, and the last error from op, wrapped, when every attempt fails.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/b4-strategy/
//
// One test passes Retry a Backoff written only in the test file. That is the
// proof: a fourth strategy, and Retry is not edited.

// Backoff says how long to wait before retry number attempt.
type Backoff func(attempt int) time.Duration

var errTODO = errors.New("TODO")

// Constant always waits d.
func Constant(d time.Duration) Backoff {
	// TODO
	return func(int) time.Duration { return 0 }
}

// Linear waits step times the retry number.
func Linear(step time.Duration) Backoff {
	// TODO
	return func(int) time.Duration { return 0 }
}

// Exponential waits base, then doubles each retry, but never more than limit.
func Exponential(base, limit time.Duration) Backoff {
	// TODO
	return func(int) time.Duration { return 0 }
}

// Retry runs op up to attempts times, sleeping b(n) between attempts.
func Retry(attempts int, b Backoff, sleep func(time.Duration), op func() error) error {
	// TODO
	return errTODO
}
