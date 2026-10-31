package adapter

import (
	"context"
	"database/sql"
	"errors"
	"os/exec"
	"strings"
	"testing"

	afterorders "github.com/luckyparakh/goDP/patterns/b2-adapter/after/orders"
	"github.com/luckyparakh/goDP/patterns/b2-adapter/after/sqlstore"
	beforeorders "github.com/luckyparakh/goDP/patterns/b2-adapter/before/orders"
	"github.com/luckyparakh/goDP/patterns/b2-adapter/fakedb"
)

var seed = map[string]int64{"42": 4200}

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

// BEFORE: the only way a caller can spot a missing order is to know it is SQL.
func TestBefore_CallerMustImportDatabaseSQL(t *testing.T) {
	_, err := beforeorders.New(fakedb.Open(seed)).Total(context.Background(), "nope")

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows to leak through; err = %v", err)
	}
}

// BEFORE: a caller can reach the driver's own error type and its code.
func TestBefore_DriverErrorLeaks(t *testing.T) {
	_, err := beforeorders.New(fakedb.Broken("08006")).Total(context.Background(), "42")

	var fe *fakedb.Error
	if !errors.As(err, &fe) || fe.Code != "08006" {
		t.Fatalf("expected the driver error to leak through; err = %v", err)
	}
}

// BEFORE: the service itself depends on database/sql.
func TestBefore_ServiceImportsDatabaseSQL(t *testing.T) {
	if !dependsOn(t, "./before/orders", "database/sql") {
		t.Fatal("expected before/orders to depend on database/sql")
	}
}

// AFTER: the service compiles with no database/sql anywhere in its imports.
func TestAfter_ServiceImportsNoDatabaseSQL(t *testing.T) {
	for _, bad := range []string{"database/sql", "/fakedb"} {
		if dependsOn(t, "./after/orders", bad) {
			t.Errorf("after/orders still depends on %s", bad)
		}
	}
}

// AFTER: a missing order is the service's own error, and the SQL one is gone.
func TestAfter_NotFoundIsADomainError(t *testing.T) {
	svc := afterorders.New(sqlstore.New(fakedb.Open(seed)))

	got, err := svc.Total(context.Background(), "42")
	if err != nil || got != 4200 {
		t.Fatalf("Total(42) = %d, %v; want 4200, nil", got, err)
	}

	_, err = svc.Total(context.Background(), "nope")
	if !errors.Is(err, afterorders.ErrNotFound) {
		t.Fatalf("errors.Is(err, orders.ErrNotFound) = false; err = %v", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		t.Fatal("sql.ErrNoRows still reaches the caller")
	}
}

// AFTER: an unknown driver failure becomes ErrUnavailable. The driver's text
// stays in the message for the logs, but its identity does not reach callers.
func TestAfter_DriverErrorIsHidden(t *testing.T) {
	_, err := afterorders.New(sqlstore.New(fakedb.Broken("08006"))).Total(context.Background(), "42")

	if !errors.Is(err, afterorders.ErrUnavailable) {
		t.Fatalf("errors.Is(err, orders.ErrUnavailable) = false; err = %v", err)
	}
	var fe *fakedb.Error
	if errors.As(err, &fe) {
		t.Fatal("callers can still reach the driver's error type with errors.As")
	}
	if err == nil || !strings.Contains(err.Error(), "08006") {
		t.Fatalf("the driver's code should stay in the message for the logs; err = %v", err)
	}
}
