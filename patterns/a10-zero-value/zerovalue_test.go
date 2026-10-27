package zerovalue

import (
	"sync"
	"testing"
)

// BEFORE: the natural two lines crash.
func TestBefore_ZeroValuePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic from writing to a nil map")
		}
	}()

	var c Counter
	c.Inc("a")
}

// BEFORE: it only works through the constructor.
func TestBefore_NewCounterWorks(t *testing.T) {
	c := NewCounter()
	c.Inc("a")
	c.Inc("a")

	if got := c.Count("a"); got != 2 {
		t.Fatalf("Count(a) = %d, want 2", got)
	}
}

// AFTER: no constructor, and it works the first time.
func TestAfter_ZeroValueIsReady(t *testing.T) {
	var tally Tally
	tally.Inc("a")
	tally.Inc("a")
	tally.Inc("b")

	if got := tally.Count("a"); got != 2 {
		t.Errorf("Count(a) = %d, want 2", got)
	}
	if got := tally.Count("b"); got != 1 {
		t.Errorf("Count(b) = %d, want 1", got)
	}
	if got := tally.Count("never-seen"); got != 0 {
		t.Errorf("Count(never-seen) = %d, want 0", got)
	}
}

// AFTER: many goroutines at once. Run this with -race.
func TestAfter_SafeForConcurrentUse(t *testing.T) {
	var tally Tally

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				tally.Inc("k")
			}
		}()
	}
	wg.Wait()

	if got := tally.Count("k"); got != 5000 {
		t.Fatalf("Count(k) = %d, want 5000", got)
	}
}
