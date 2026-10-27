package funciface

// GOAL
//
// Do the same job a second way, then decide which one you would ship.
//
//   - FilterFunc does what Filter does, but takes a plain function instead of
//     a Rule.
//   - RuleFunc is an adapter: a function type that has an Allow method, so any
//     function can be used wherever a Rule is expected. Fill in its Allow.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/a9-func-vs-interface/
//
// Then measure both versions and write your decision in THEORY.md:
//
//	go test -bench . -benchmem ./patterns/a9-func-vs-interface/

// RuleFunc lets an ordinary function be used as a Rule.
type RuleFunc func(n int) bool

// Allow calls the function itself.
func (f RuleFunc) Allow(n int) bool {
	// TODO
	return false
}

// FilterFunc returns the numbers f allows, in order.
func FilterFunc(nums []int, f func(n int) bool) []int {
	// TODO
	return nil
}
