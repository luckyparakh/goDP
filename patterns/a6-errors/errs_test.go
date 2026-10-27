package errs

import (
	"errors"
	"strings"
	"testing"
)

// BEFORE: the message says "not found", but the error itself is gone.
func TestBefore_VerbVLosesTheSentinel(t *testing.T) {
	_, err := LoadBefore("nope")

	if errors.Is(err, ErrNotFound) {
		t.Fatal("did not expect errors.Is to find ErrNotFound: the error was wrapped without the w verb")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("the text should still say it: %q", err)
	}
}

// BEFORE: the message names the field, but there is no typed error to ask.
func TestBefore_VerbVLosesTheType(t *testing.T) {
	_, err := LoadBefore("")

	var ve *ValidationError
	if errors.As(err, &ve) {
		t.Fatal("did not expect errors.As to find a *ValidationError: the error was wrapped without the w verb")
	}
	if !strings.Contains(err.Error(), "invalid id") {
		t.Fatalf("the text should still say it: %q", err)
	}
}

// chainLength counts the errors from err down to the innermost one.
func chainLength(err error) int {
	n := 0
	for e := err; e != nil; e = errors.Unwrap(e) {
		n++
	}
	return n
}

// AFTER: the happy path is unchanged.
func TestAfter_HappyPath(t *testing.T) {
	got, err := Handle("42")
	if err != nil || got != "order 42" {
		t.Fatalf("Handle(42) = %q, %v; want order 42, nil", got, err)
	}
}

// AFTER: the sentinel is still found after two wraps, and each layer left
// its own context in the message.
func TestAfter_SentinelSurvivesTwoWraps(t *testing.T) {
	_, err := Handle("nope")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false; err = %v", err)
	}
	if n := chainLength(err); n < 3 {
		t.Fatalf("chain has %d errors, want at least 3 (handle, load, sentinel)", n)
	}
	for _, want := range []string{"handle", "load"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("message %q is missing context %q", err, want)
		}
	}
}

// AFTER: the typed error is still found after two wraps, with its data.
func TestAfter_TypedErrorSurvivesTwoWraps(t *testing.T) {
	_, err := Handle("")

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("errors.As found no *ValidationError; err = %v", err)
	}
	if ve.Field != "id" {
		t.Fatalf("ve.Field = %q, want id", ve.Field)
	}
	if errors.Is(err, ErrNotFound) {
		t.Fatal("an empty id is a validation problem, not 'not found'")
	}
}
