//go:build ignore

package ocp

import "fmt"

// HOW IT WORKS
//
// The tiers live in a map from name to a small func. Quote looks the tier up
// and calls the func. Register adds an entry. Each file can call Register from
// its own init(), so a new tier is a new file and no existing file changes.
//
// The map is created where it is declared, so it exists before any init() in
// the package runs, whichever file that init() is in.
//
// WHY IT HELPS
//
// Adding "student" cannot break "regular" or "member", because their code is
// not touched. Nothing that already works has to be re-checked.
//
// WHEN NOT TO USE
//
// You build the lookup before you need it, and the rules are now spread over
// several files instead of one switch you can read top to bottom. If tiers
// almost never change, the switch was fine.

// Pricer turns a list price into a final price, both in cents.
type Pricer func(cents int64) int64

var registry = map[string]Pricer{}

// Register adds a tier. Call it from init(). The map is not safe for
// concurrent writes.
func Register(tier string, p Pricer) { registry[tier] = p }

func init() {
	Register("regular", func(cents int64) int64 { return cents })
	Register("member", func(cents int64) int64 { return cents * 90 / 100 })
}

// Quote applies the tier discount to a list price, both in cents.
func Quote(tier string, cents int64) (int64, error) {
	p, ok := registry[tier]
	if !ok {
		return 0, fmt.Errorf("quote %q: %w", tier, ErrUnknownTier)
	}
	return p(cents), nil
}
