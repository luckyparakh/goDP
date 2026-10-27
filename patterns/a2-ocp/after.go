package ocp

import (
	"fmt"
	"sync"
)

// GOAL
//
// Part 1. Make Quote behave exactly like Price in before.go, but without a
// switch. You need:
//   - a plain func type for "how to price one tier" (not an interface)
//   - a map from tier name to that func
//   - a Register function that adds to the map
//   - an init() that registers "regular" and "member"
//
// An unknown tier must wrap ErrUnknownTier, so errors.Is(err, ErrUnknownTier)
// is true.
//
// Part 2. Commit part 1. Then add a "student" tier (25% off) in a brand-new
// file in this folder, without editing any existing file. To prove it:
//
//	git add -N patterns/a2-ocp
//	git diff --name-status
//
// You should see "A" for your new file and no "M" lines.
//
// One question to answer before you write init(): what guarantees the map
// exists before any init() in the package runs, in whichever file it lives?

// Quote applies the Pricer discount to a list price, both in cents.
type Pricer = func(int64) int64

var l sync.Mutex
var mapper map[string]Pricer = make(map[string]Pricer)

func register(tier string, tierFunc Pricer) {
	if _, e := mapper[tier]; e {
		panic("tier already exists")
	}

	// Although here its safe here as register is being called from init only.
	// And initS are called in seq by GO
	l.Lock()
	mapper[tier] = tierFunc
	l.Unlock()
}

func init() {
	register("regular", regular)
	register("member", member)
}

func Quote(tier string, cents int64) (int64, error) {
	f, exists := mapper[tier]
	if !exists {
		return 0, fmt.Errorf("no tier name %s, %w", tier, ErrUnknownTier)
	}

	return f(cents), nil
}

func regular(cents int64) int64 {
	return cents
}

func member(cents int64) int64 {
	return cents * 90 / 100
}
