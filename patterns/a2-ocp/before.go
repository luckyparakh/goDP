package ocp

import "errors"

// PROBLEM
//
// Price is one switch with a case per tier. Every new tier means editing this
// function, which already works and is already tested for the other tiers.
// That is the risk: you can break "regular" while adding "student".
//
// TRY IT
//
// Add a "student" tier (25% off) as a new case below. You just edited code
// that was fine, so you must re-run the tests for every other tier too. Then
// undo it:
//
//	git checkout -- patterns/a2-ocp/before.go
//
// Do not keep that edit. Your work goes in after.go.

// ErrUnknownTier is returned for a tier nobody has defined.
var ErrUnknownTier = errors.New("unknown tier")

// Price applies the tier discount to a list price, both in cents.
func Price(tier string, cents int64) (int64, error) {
	switch tier {
	case "regular":
		return cents, nil
	case "member":
		return cents * 90 / 100, nil // 10% off
	default:
		return 0, ErrUnknownTier
	}
}
