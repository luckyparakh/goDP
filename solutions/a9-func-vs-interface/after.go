//go:build ignore

package funciface

// HOW IT WORKS
//
// FilterFunc takes the decision as a plain function, so a one-line rule is one
// line and needs no type. RuleFunc is a function type with an Allow method that
// just calls the function. That one method is enough to make any function
// satisfy Rule, which is how net/http.HandlerFunc lets a function be a Handler.
//
// WHY IT HELPS
//
// You pick the weight to match the job. A function is the lightest form for
// one decision. The adapter means choosing the function form never locks
// callers out of an API that asks for the interface.
//
// WHEN NOT TO USE
//
// A function is the wrong shape when the thing has several methods, needs
// identity, or must be inspected, for example with a type switch. Then use an
// interface. Do not add RuleFunc if nobody calls the interface version.

// RuleFunc lets an ordinary function be used as a Rule.
type RuleFunc func(n int) bool

// Allow calls the function itself.
func (f RuleFunc) Allow(n int) bool { return f(n) }

// FilterFunc returns the numbers f allows, in order.
func FilterFunc(nums []int, f func(n int) bool) []int {
	var out []int
	for _, n := range nums {
		if f(n) {
			out = append(out, n)
		}
	}
	return out
}
