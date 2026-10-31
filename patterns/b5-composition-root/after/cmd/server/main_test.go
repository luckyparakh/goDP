package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// AFTER: the composition root builds the real thing from config alone.
func TestAfter_BuildWiresTheRealThing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "orders.json")
	if err := os.WriteFile(path, []byte(`{"42": 4200}`), 0o644); err != nil {
		t.Fatal(err)
	}
	env := map[string]string{"ORDERS_FILE": path}

	h, err := build(func(key string) string { return env[key] })
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/orders/42", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "4200" {
		t.Fatalf("GET /orders/42 = %d %q, want 200 %q", rec.Code, rec.Body.String(), "4200")
	}
}
