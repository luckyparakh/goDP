package nullobject

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"strings"
	"testing"
)

// recorder is a Metrics that remembers what it counted.
type recorder struct{ counts map[string]int }

func (r *recorder) Inc(name string) { r.counts[name]++ }

// panics runs f and reports whether it panicked. Tests use it so a panic
// fails one test instead of stopping the whole run.
func panics(f func()) (did bool) {
	defer func() {
		if recover() != nil {
			did = true
		}
	}()
	f()
	return false
}

// fieldNilChecks counts the comparisons of a field with nil in file, such as
// s.metrics != nil. A plain name compared with nil, like err or m, does not
// count.
func fieldNilChecks(t *testing.T, file string) int {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), file, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	ast.Inspect(f, func(node ast.Node) bool {
		b, ok := node.(*ast.BinaryExpr)
		if !ok || (b.Op != token.EQL && b.Op != token.NEQ) {
			return true
		}
		for _, pair := range [][2]ast.Expr{{b.X, b.Y}, {b.Y, b.X}} {
			_, isField := pair[0].(*ast.SelectorExpr)
			id, isIdent := pair[1].(*ast.Ident)
			if isField && isIdent && id.Name == "nil" {
				n++
			}
		}
		return true
	})
	return n
}

// BEFORE: the one forgotten check panics as soon as metrics are left out.
func TestBefore_ForgottenNilCheckPanics(t *testing.T) {
	if !panics(func() { _ = NewServiceBefore(nil, nil).Refund(500) }) {
		t.Fatal("expected Refund to panic when metrics is nil")
	}
}

// BEFORE: the guards are everywhere.
func TestBefore_GuardsAreEverywhere(t *testing.T) {
	if n := fieldNilChecks(t, "before.go"); n < 4 {
		t.Fatalf("before.go has %d field nil checks, expected at least 4", n)
	}
}

// AFTER: nothing inside the service is ever compared with nil.
func TestAfter_NothingIsNilInside(t *testing.T) {
	if n := fieldNilChecks(t, "after.go"); n != 0 {
		t.Fatalf("after.go still compares a field with nil %d times, want 0", n)
	}
}

// AFTER: a service built with nothing works, and never panics.
func TestAfter_NilDependenciesAreSafe(t *testing.T) {
	svc := New(nil, nil)
	var charge, refund, bad error
	if panics(func() {
		charge = svc.Charge(500)
		refund = svc.Refund(500)
		bad = svc.Charge(0)
	}) {
		t.Fatal("a Service built with nil logger and metrics panicked")
	}
	if charge != nil || refund != nil || !errors.Is(bad, ErrBadAmount) {
		t.Fatalf("Charge %v, Refund %v, Charge(0) %v; want nil, nil, ErrBadAmount", charge, refund, bad)
	}
}

// AFTER: real dependencies still get every call.
func TestAfter_RealDependenciesStillWork(t *testing.T) {
	var logs bytes.Buffer
	rec := &recorder{counts: map[string]int{}}
	svc := New(slog.New(slog.NewTextHandler(&logs, nil)), rec)

	if panics(func() {
		_ = svc.Charge(500)
		_ = svc.Refund(200)
		_ = svc.Charge(0)
	}) {
		t.Fatal("a Service built with a real logger and metrics panicked")
	}
	for _, name := range []string{"charges", "refunds", "charge_rejected"} {
		if rec.counts[name] != 1 {
			t.Errorf("metric %q counted %d times, want 1", name, rec.counts[name])
		}
	}
	for _, msg := range []string{"msg=charge", "msg=refund"} {
		if !strings.Contains(logs.String(), msg) {
			t.Errorf("log is missing %q:\n%s", msg, logs.String())
		}
	}
}

// The two shapes, side by side. Compare their ns/op and allocs/op.
func BenchmarkBefore_NilChecks(b *testing.B) {
	svc := NewServiceBefore(nil, nil)
	for b.Loop() {
		_ = svc.Charge(500)
	}
}

func BenchmarkAfter_NoOps(b *testing.B) {
	svc := New(nil, nil)
	for b.Loop() {
		_ = svc.Charge(500)
	}
}
