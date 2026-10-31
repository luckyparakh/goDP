package strategy

import (
	"errors"
	"slices"
	"testing"
	"time"
)

const ms = time.Millisecond

var errFail = errors.New("fail")

// alwaysFails is an op that never succeeds.
func alwaysFails() error { return errFail }

// sleeps runs retry with a sleep that records instead of sleeping.
func sleeps(retry func(sleep func(time.Duration)) error) ([]time.Duration, error) {
	var got []time.Duration
	err := retry(func(d time.Duration) { got = append(got, d) })
	return got, err
}

// BEFORE: the three named policies work.
func TestBefore_DelaysByName(t *testing.T) {
	cases := []struct {
		policy string
		want   []time.Duration // for attempts 1 to 5
	}{
		{"constant", []time.Duration{100 * ms, 100 * ms, 100 * ms, 100 * ms, 100 * ms}},
		{"linear", []time.Duration{100 * ms, 200 * ms, 300 * ms, 400 * ms, 500 * ms}},
		{"exponential", []time.Duration{100 * ms, 200 * ms, 400 * ms, 800 * ms, time.Second}},
	}
	for _, c := range cases {
		for i, want := range c.want {
			if got := Delay(c.policy, i+1); got != want {
				t.Errorf("Delay(%q, %d) = %v, want %v", c.policy, i+1, got, want)
			}
		}
	}
}

// BEFORE: a typo compiles, returns no error, and waits the wrong time.
func TestBefore_TypoFallsBackSilently(t *testing.T) {
	if got := Delay("exponental", 4); got != 100*ms {
		t.Fatalf("Delay(\"exponental\", 4) = %v; expected the silent fallback to 100ms", got)
	}
}

// AFTER: the three Backoffs wait exactly as the named policies did.
func TestAfter_SameDelays(t *testing.T) {
	cases := []struct {
		policy string
		b      Backoff
	}{
		{"constant", Constant(100 * ms)},
		{"linear", Linear(100 * ms)},
		{"exponential", Exponential(100*ms, time.Second)},
	}
	for _, c := range cases {
		want, _ := sleeps(func(s func(time.Duration)) error { return RetryBefore(6, c.policy, s, alwaysFails) })
		got, _ := sleeps(func(s func(time.Duration)) error { return Retry(6, c.b, s, alwaysFails) })
		if !slices.Equal(got, want) {
			t.Errorf("%s: slept %v, want %v", c.policy, got, want)
		}
	}
}

// AFTER: Retry stops at the first success.
func TestAfter_StopsOnSuccess(t *testing.T) {
	calls := 0
	op := func() error {
		calls++
		if calls < 2 {
			return errFail
		}
		return nil
	}
	got, err := sleeps(func(s func(time.Duration)) error { return Retry(5, Constant(100*ms), s, op) })

	if err != nil || calls != 2 || !slices.Equal(got, []time.Duration{100 * ms}) {
		t.Fatalf("err %v, calls %d, slept %v; want nil, 2, [100ms]", err, calls, got)
	}
}

// AFTER: the proof. This Backoff exists only in this test file, and Retry
// uses it without being edited.
func TestAfter_FourthStrategyNeedsNoEdit(t *testing.T) {
	steps := []time.Duration{7 * ms, 11 * ms, 13 * ms}
	fromList := Backoff(func(attempt int) time.Duration { return steps[attempt-1] })

	got, _ := sleeps(func(s func(time.Duration)) error { return Retry(4, fromList, s, alwaysFails) })

	if !slices.Equal(got, steps) {
		t.Fatalf("slept %v, want %v", got, steps)
	}
}

// AFTER: when every attempt fails, the error wraps the last one, and fewer
// than one attempt is refused.
func TestAfter_ReturnsLastError(t *testing.T) {
	errs := []error{errors.New("first"), errors.New("second"), errors.New("third")}
	calls := 0
	op := func() error { calls++; return errs[calls-1] }

	err := Retry(3, Constant(0), func(time.Duration) {}, op)
	if !errors.Is(err, errs[2]) || errors.Is(err, errs[0]) {
		t.Fatalf("err = %v, want it to wrap the last error only", err)
	}
	if err := Retry(0, Constant(0), func(time.Duration) {}, alwaysFails); !errors.Is(err, ErrNoAttempts) {
		t.Fatalf("Retry(0, ...) = %v, want ErrNoAttempts", err)
	}
}
