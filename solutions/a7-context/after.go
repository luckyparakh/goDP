//go:build ignore

package ctxflow

import (
	"context"
	"fmt"
	"time"
)

// HOW IT WORKS
//
// The context is the first parameter of every level and is passed down
// unchanged. Only the bottom level waits, so only Repo listens to it: each
// step waits on a timer and on ctx.Done() at the same time, and whichever
// happens first wins. When the context ends, Repo returns the steps it has
// finished and ctx.Err(). The levels above add context with %w, so
// errors.Is(err, context.DeadlineExceeded) still works at the top.
//
// WHY IT HELPS
//
// One timeout at the top bounds the whole call chain, and one cancel stops all
// of it. Nothing keeps working for a caller who has left, so there is nothing
// wasted and nothing left running.
//
// WHEN NOT TO USE
//
// Do not store a context in a struct, and do not pass one to a function that
// does no waiting and cannot fail. Do not use context.Background() in the
// middle of a call chain, because that cuts the cancel signal off. Context
// values are for request-scoped data, not for passing ordinary parameters.

// Handle is the top level. It calls Service.
func Handle(ctx context.Context) (int, error) {
	n, err := Service(ctx)
	if err != nil {
		return n, fmt.Errorf("handle: %w", err)
	}
	return n, nil
}

// Service is the middle level. It calls Repo.
func Service(ctx context.Context) (int, error) {
	n, err := Repo(ctx)
	if err != nil {
		return n, fmt.Errorf("service: %w", err)
	}
	return n, nil
}

// Repo is the bottom level. It does the slow work.
func Repo(ctx context.Context) (int, error) {
	done := 0
	for i := 0; i < totalSteps; i++ {
		timer := time.NewTimer(stepDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return done, ctx.Err()
		case <-timer.C:
			done++
		}
	}
	return done, nil
}
