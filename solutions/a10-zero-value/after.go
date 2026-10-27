//go:build ignore

package zerovalue

import "sync"

// HOW IT WORKS
//
// Every field of Tally is already usable as its zero value. A sync.Mutex is an
// unlocked lock when it is zero, and reading a nil map is allowed and returns
// 0. The one thing that is not allowed, writing to a nil map, is handled where
// the write happens: Inc makes the map the first time it needs it, inside the
// lock, so two goroutines cannot both create it.
//
// WHY IT HELPS
//
// There is no constructor to forget. Tally works as a package variable, as a
// field in another struct, and as a local, with no setup. It is safe for many
// goroutines, so callers do not have to add their own lock.
//
// WHEN NOT TO USE
//
// Some types cannot have a useful zero value, for example one that needs a
// connection string or a required setting. Then use a constructor that checks
// its input, so an invalid value can never be built. Also, a struct holding a
// mutex must not be copied after first use. go vet reports this.

// Tally counts how often each key was seen.
type Tally struct {
	mu     sync.Mutex
	counts map[string]int
}

// Inc adds one to the count for key.
func (t *Tally) Inc(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.counts == nil {
		t.counts = make(map[string]int)
	}
	t.counts[key]++
}

// Count returns how many times key was seen. An unknown key counts as 0.
func (t *Tally) Count(key string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counts[key]
}
