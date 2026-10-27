package ocp

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

// Quote applies the tier discount to a list price, both in cents.
func Quote(tier string, cents int64) (int64, error) {
	// TODO
	return 0, nil
}
