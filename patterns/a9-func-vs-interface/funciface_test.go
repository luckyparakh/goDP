package funciface

import (
	"slices"
	"testing"
)

var nums = []int{1, 2, 3, 4, 5, 6}

// BEFORE: the interface version works.
func TestBefore_InterfaceVersionWorks(t *testing.T) {
	if got := Filter(nums, AtLeast{Min: 4}); !slices.Equal(got, []int{4, 5, 6}) {
		t.Errorf("AtLeast{4} = %v, want [4 5 6]", got)
	}
	if got := Filter(nums, Even{}); !slices.Equal(got, []int{2, 4, 6}) {
		t.Errorf("Even = %v, want [2 4 6]", got)
	}
}

// AFTER: a plain function does the same job with no new type.
func TestAfter_FilterFuncMatchesFilter(t *testing.T) {
	got := FilterFunc(nums, func(n int) bool { return n >= 4 })
	want := Filter(nums, AtLeast{Min: 4})
	if !slices.Equal(got, want) {
		t.Fatalf("FilterFunc = %v, want %v", got, want)
	}
}

// AFTER: the adapter lets a function be passed where a Rule is expected.
func TestAfter_RuleFuncSatisfiesRule(t *testing.T) {
	byThree := RuleFunc(func(n int) bool { return n%3 == 0 })

	if got := Filter(nums, byThree); !slices.Equal(got, []int{3, 6}) {
		t.Fatalf("Filter(RuleFunc) = %v, want [3 6]", got)
	}
}

// AFTER: a closure can keep state without a new type.
func TestAfter_ClosureKeepsState(t *testing.T) {
	calls := 0
	FilterFunc(nums, func(n int) bool {
		calls++
		return true
	})
	if calls != len(nums) {
		t.Fatalf("closure was called %d times, want %d", calls, len(nums))
	}
}

var sink []int

var bench = func() []int {
	s := make([]int, 1000)
	for i := range s {
		s[i] = i
	}
	return s
}()

func BenchmarkFilterInterface(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = Filter(bench, AtLeast{Min: 500})
	}
}

func BenchmarkFilterFunc(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink = FilterFunc(bench, func(n int) bool { return n >= 500 })
	}
}
