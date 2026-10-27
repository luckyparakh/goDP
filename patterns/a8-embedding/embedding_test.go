package embedding

import "testing"

// BEFORE: calling Name directly on Derived does use the replacement.
func TestBefore_OverrideWorksOnDirectCalls(t *testing.T) {
	if got := (Derived{}).Name(); got != "derived" {
		t.Fatalf("Derived.Name() = %q, want derived", got)
	}
}

// BEFORE: the surprise. Hello is a method of Base, so it sees Base's Name.
func TestBefore_HelloIgnoresTheOverride(t *testing.T) {
	if got := (Derived{}).Hello(); got != "hello from base" {
		t.Fatalf("Derived.Hello() = %q, want hello from base", got)
	}
}

// AFTER: Greet uses whichever Namer it is handed.
func TestAfter_GreetUsesTheNamerItIsGiven(t *testing.T) {
	if got := Greet(Base{}); got != "hello from base" {
		t.Errorf("Greet(Base) = %q, want hello from base", got)
	}
	if got := Greet(Derived{}); got != "hello from derived" {
		t.Errorf("Greet(Derived) = %q, want hello from derived", got)
	}
}

// AFTER: Shout replaces one method and forwards the rest.
func TestAfter_ShoutOverridesName(t *testing.T) {
	if got := (Shout{Derived{}}).Name(); got != "DERIVED" {
		t.Errorf("Shout.Name() = %q, want DERIVED", got)
	}
	if got := Greet(Shout{Base{}}); got != "hello from BASE" {
		t.Errorf("Greet(Shout) = %q, want hello from BASE", got)
	}
}

// Shout still counts as a Namer, so it can be passed anywhere a Namer fits.
var _ Namer = Shout{}
