package consumeriface

import (
	"os/exec"
	"strings"
	"testing"

	afterbilling "github.com/luckyparakh/goDP/patterns/b1-consumer-iface/after/billing"
	afterstore "github.com/luckyparakh/goDP/patterns/b1-consumer-iface/after/store"
	beforebilling "github.com/luckyparakh/goDP/patterns/b1-consumer-iface/before/billing"
	beforestore "github.com/luckyparakh/goDP/patterns/b1-consumer-iface/before/store"
)

// dependsOn reports whether the package at pkg (a path relative to this
// folder) depends, directly or not, on a package whose path ends in suffix.
// It asks `go list -deps`, the same command you can run yourself.
func dependsOn(t *testing.T, pkg, suffix string) bool {
	t.Helper()
	gobin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("the go command is not on PATH")
	}
	out, err := exec.Command(gobin, "list", "-deps", pkg).CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps %s: %v\n%s", pkg, err, out)
	}
	for _, dep := range strings.Fields(string(out)) {
		if strings.HasSuffix(dep, suffix) {
			return true
		}
	}
	return false
}

// BEFORE: billing depends on store. The arrow points from the code that
// decides to the code that stores.
func TestBefore_BillingImportsStore(t *testing.T) {
	if !dependsOn(t, "./before/billing", "/before/store") {
		t.Fatal("expected before/billing to depend on before/store")
	}
	if dependsOn(t, "./before/store", "/before/billing") {
		t.Fatal("did not expect before/store to depend on before/billing")
	}
}

// BEFORE: it works. The problem is the shape, not the behaviour.
func TestBefore_Due(t *testing.T) {
	st := beforestore.NewMem()
	st.Add("inv-1", 10000, 2500)

	got, err := beforebilling.New(st).Due("inv-1")
	if err != nil || got != 7500 {
		t.Fatalf("Due = %d, %v; want 7500, nil", got, err)
	}
}

// AFTER: billing no longer depends on store at all.
func TestAfter_BillingDoesNotImportStore(t *testing.T) {
	if dependsOn(t, "./after/billing", "/after/store") {
		t.Fatal("after/billing still depends on after/store: billing should declare its own Invoice and interface")
	}
}

// AFTER: the arrow is reversed. store depends on billing, because it hands
// back billing's own Invoice type. And billing still works.
func TestAfter_StoreImportsBilling(t *testing.T) {
	if !dependsOn(t, "./after/store", "/after/billing") {
		t.Fatal("after/store does not depend on after/billing: store should return billing.Invoice")
	}

	st := afterstore.NewMem()
	st.Add("inv-1", 10000, 2500)

	got, err := afterbilling.New(st).Due("inv-1")
	if err != nil || got != 7500 {
		t.Fatalf("Due = %d, %v; want 7500, nil", got, err)
	}
}
