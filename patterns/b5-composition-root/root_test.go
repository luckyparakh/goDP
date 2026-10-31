package compositionroot

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	afterhandler "github.com/luckyparakh/goDP/patterns/b5-composition-root/after/handler"
	afterservice "github.com/luckyparakh/goDP/patterns/b5-composition-root/after/service"
	beforehandler "github.com/luckyparakh/goDP/patterns/b5-composition-root/before/handler"
)

// builds lists the calls in dir that build a dependency: a New function from
// one of this task's packages, or a read of the environment. It is the grep
// from after.go, done with the Go parser so comments cannot fool it.
func builds(t *testing.T, dir string) []string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var hits []string
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		imported := map[string]string{} // name used in this file -> import path
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			name := p[strings.LastIndex(p, "/")+1:]
			if imp.Name != nil {
				name = imp.Name.Name
			}
			imported[name] = p
		}
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			pkg, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}
			p, fn := imported[pkg.Name], sel.Sel.Name
			ours := strings.Contains(p, "/b5-composition-root/") && strings.HasPrefix(fn, "New")
			env := p == "os" && (fn == "Getenv" || fn == "LookupEnv")
			if ours || env {
				hits = append(hits, fmt.Sprintf("%s: %s.%s", fset.Position(call.Pos()), pkg.Name, fn))
			}
			return true
		})
	}
	return hits
}

// writeOrders writes a small orders file and returns its path.
func writeOrders(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "orders.json")
	if err := os.WriteFile(path, []byte(`{"42": 4200}`), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// get sends GET path to h.
func get(h http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

// fakeTotals stands in for the service, or the store: a map, no file.
type fakeTotals map[string]int64

func (f fakeTotals) lookup(id string) (int64, error) {
	c, ok := f[id]
	if !ok {
		return 0, errors.New("no such order")
	}
	return c, nil
}

type fakeService struct{ fakeTotals }

func (f fakeService) Total(_ context.Context, id string) (int64, error) { return f.lookup(id) }

type fakeStore struct{ fakeTotals }

func (f fakeStore) Cents(_ context.Context, id string) (int64, error) { return f.lookup(id) }

// BEFORE: a test of an HTTP handler has to set up an environment variable and
// a real file, because the handler builds the store it never mentions.
func TestBefore_HandlerNeedsEnvAndFile(t *testing.T) {
	t.Setenv("ORDERS_FILE", "")
	if _, err := beforehandler.New(); err == nil {
		t.Fatal("expected handler.New to fail with no ORDERS_FILE")
	}

	t.Setenv("ORDERS_FILE", writeOrders(t))
	h, err := beforehandler.New()
	if err != nil {
		t.Fatal(err)
	}
	if rec := get(h, "/orders/42"); rec.Code != http.StatusOK || rec.Body.String() != "4200" {
		t.Fatalf("GET /orders/42 = %d %q, want 200 %q", rec.Code, rec.Body.String(), "4200")
	}
}

// BEFORE: every package builds its own dependency.
func TestBefore_PackagesBuildTheirOwnDependencies(t *testing.T) {
	for _, dir := range []string{"before/handler", "before/service", "before/store"} {
		if len(builds(t, dir)) == 0 {
			t.Errorf("expected %s to build its own dependency", dir)
		}
	}
}

// AFTER: no package below main builds anything.
func TestAfter_NoPackageBuildsItsOwnDependencies(t *testing.T) {
	for _, dir := range []string{"after/handler", "after/service", "after/store"} {
		for _, hit := range builds(t, dir) {
			t.Errorf("only main may build dependencies or read the environment: %s", hit)
		}
	}
}

// AFTER: the handler works with a fake service. No environment, no file.
func TestAfter_HandlerTakesAFakeService(t *testing.T) {
	t.Setenv("ORDERS_FILE", "")
	h := afterhandler.New(fakeService{fakeTotals{"7": 700}})

	if rec := get(h, "/orders/7"); rec.Code != http.StatusOK || rec.Body.String() != "700" {
		t.Fatalf("GET /orders/7 = %d %q, want 200 %q", rec.Code, rec.Body.String(), "700")
	}
}

// AFTER: the service works with a fake store. No environment, no file.
func TestAfter_ServiceTakesAFakeStore(t *testing.T) {
	t.Setenv("ORDERS_FILE", "")
	got, err := afterservice.New(fakeStore{fakeTotals{"7": 700}}).Total(context.Background(), "7")

	if err != nil || got != 700 {
		t.Fatalf("Total(7) = %d, %v; want 700, nil", got, err)
	}
}

// AFTER: main is where the three are built.
func TestAfter_MainIsTheOnlyPlaceThatWires(t *testing.T) {
	hits := strings.Join(builds(t, "after/cmd/server"), "\n")
	for _, want := range []string{"store.New", "service.New", "handler.New"} {
		if !strings.Contains(hits, want) {
			t.Errorf("after/cmd/server never calls %s", want)
		}
	}
}
