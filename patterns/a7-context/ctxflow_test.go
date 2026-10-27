package ctxflow

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// BEFORE: the caller gives up after 30ms, and the work carries on anyway.
func TestBefore_WorkContinuesAfterTheCallerGivesUp(t *testing.T) {
	var done atomic.Int64
	go FetchNoCtx(&done)

	time.Sleep(30 * time.Millisecond) // the caller gives up here
	time.Sleep(500 * time.Millisecond)

	if got := done.Load(); got != totalSteps {
		t.Fatalf("steps done = %d, want %d: the work should not have been stoppable", got, totalSteps)
	}
}

// AFTER: with plenty of time, all steps run and there is no error.
func TestAfter_FinishesWhenThereIsEnoughTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	n, err := Handle(ctx)
	if err != nil || n != totalSteps {
		t.Fatalf("Handle = %d, %v; want %d, nil", n, err, totalSteps)
	}
}

// AFTER: a timeout set at the top reaches the bottom and stops the work.
func TestAfter_TimeoutReachesTheBottom(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	n, err := Handle(ctx)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want it to wrap context.DeadlineExceeded", err)
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf("ctx.Err() = %v, want DeadlineExceeded", ctx.Err())
	}
	if n >= totalSteps {
		t.Fatalf("steps done = %d: the work was not stopped", n)
	}
}

// AFTER: a cancel at the top also reaches the bottom.
func TestAfter_CancelReachesTheBottom(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	n, err := Handle(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want it to wrap context.Canceled", err)
	}
	if n >= totalSteps {
		t.Fatalf("steps done = %d: the work was not stopped", n)
	}
}

// AFTER: a context that is already cancelled means no work at all.
func TestAfter_AlreadyCancelledDoesNoWork(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	n, err := Handle(ctx)

	if !errors.Is(err, context.Canceled) || n != 0 {
		t.Fatalf("Handle = %d, %v; want 0, Canceled", n, err)
	}
}
