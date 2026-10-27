package ocp

import (
	"errors"
	"testing"
)

var cases = []struct {
	tier string
	in   int64
	want int64
}{
	{"regular", 10000, 10000},
	{"member", 10000, 9000},
}

func TestBefore_Price(t *testing.T) {
	for _, c := range cases {
		got, err := Price(c.tier, c.in)
		if err != nil || got != c.want {
			t.Errorf("Price(%q, %d) = %d, %v; want %d, nil", c.tier, c.in, got, err, c.want)
		}
	}
}

func TestBefore_PriceUnknownTier(t *testing.T) {
	if _, err := Price("student", 10000); !errors.Is(err, ErrUnknownTier) {
		t.Fatalf("err = %v, want ErrUnknownTier", err)
	}
}

func TestAfter_Quote(t *testing.T) {
	for _, c := range cases {
		got, err := Quote(c.tier, c.in)
		if err != nil || got != c.want {
			t.Errorf("Quote(%q, %d) = %d, %v; want %d, nil", c.tier, c.in, got, err, c.want)
		}
	}
}

func TestAfter_QuoteUnknownTier(t *testing.T) {
	if _, err := Quote("nonexistent", 10000); !errors.Is(err, ErrUnknownTier) {
		t.Fatalf("err = %v, want ErrUnknownTier", err)
	}
}

// The extension exercise: this passes only after you add a new file that
// registers "student" (25% off) without touching any existing file.
func TestAfter_StudentTier(t *testing.T) {
	got, err := Quote("student", 10000)
	if err != nil || got != 7500 {
		t.Fatalf("Quote(student, 10000) = %d, %v; want 7500, nil", got, err)
	}
}
