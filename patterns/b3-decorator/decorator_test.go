package decorator

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const token = "s3cret"

// serve sends GET path to h, with tok as the X-Token header if it is not empty.
func serve(h http.Handler, path, tok string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if tok != "" {
		req.Header.Set("X-Token", tok)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// noRecord is a Timing sink for tests that do not care about durations.
func noRecord(string, time.Duration) {}

// BEFORE: one copy of the pasted code forgot the token check.
func TestBefore_ExportForgotAuth(t *testing.T) {
	if rec := serve(RoutesBefore(token, io.Discard), "/export", ""); rec.Code != http.StatusOK {
		t.Fatalf("expected the bug: /export with no token = %d, want 200", rec.Code)
	}
}

// BEFORE: the copies that did paste the check work.
func TestBefore_OrdersNeedsToken(t *testing.T) {
	if rec := serve(RoutesBefore(token, io.Discard), "/orders", ""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("/orders with no token = %d, want 401", rec.Code)
	}
}

// AFTER: the first middleware is the outermost.
func TestAfter_ChainOrder(t *testing.T) {
	var got []string
	trace := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = append(got, name+">")
				next.ServeHTTP(w, r)
				got = append(got, "<"+name)
			})
		}
	}
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { got = append(got, "h") })

	serve(Chain(h, trace("A"), trace("B"), trace("C")), "/", "")

	if want := "A> B> C> h <C <B <A"; strings.Join(got, " ") != want {
		t.Fatalf("order = %q, want %q", strings.Join(got, " "), want)
	}
}

// AFTER: a request with the wrong token never reaches the handler.
func TestAfter_AuthStopsTheRequest(t *testing.T) {
	calls := 0
	h := Auth(token)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))

	if rec := serve(h, "/", "wrong"); rec.Code != http.StatusUnauthorized || calls != 0 {
		t.Fatalf("wrong token: code %d, handler calls %d; want 401, 0", rec.Code, calls)
	}
	if rec := serve(h, "/", token); rec.Code != http.StatusOK || calls != 1 {
		t.Fatalf("right token: code %d, handler calls %d; want 200, 1", rec.Code, calls)
	}
}

// AFTER: Logging is outside Auth, so it records the requests Auth turns away.
func TestAfter_RejectedRequestIsLogged(t *testing.T) {
	var log bytes.Buffer
	serve(Routes(token, &log, time.Now, noRecord), "/orders", "")

	if !strings.Contains(log.String(), "GET /orders 401") {
		t.Fatalf("log = %q, want a line containing %q", log.String(), "GET /orders 401")
	}
}

// AFTER: one chain covers every route, including the one that was forgotten.
func TestAfter_EveryRouteIsProtected(t *testing.T) {
	h := Routes(token, io.Discard, time.Now, noRecord)
	for _, name := range []string{"orders", "report", "export"} {
		path := "/" + name
		if rec := serve(h, path, ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("%s with no token = %d, want 401", path, rec.Code)
		}
		rec := serve(h, path, token)
		if rec.Code != http.StatusOK || rec.Body.String() != name {
			t.Errorf("%s with token = %d %q, want 200 %q", path, rec.Code, rec.Body.String(), name)
		}
	}
}

// AFTER: Timing measures with the clock it is given.
func TestAfter_TimingRecordsTheDuration(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ticks := 0
	now := func() time.Time { // each call is 5ms later than the last
		ticks++
		return start.Add(time.Duration(ticks) * 5 * time.Millisecond)
	}
	var gotPath string
	var gotD time.Duration
	record := func(path string, d time.Duration) { gotPath, gotD = path, d }

	h := Timing(now, record)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	serve(h, "/orders", "")

	if gotPath != "/orders" || gotD != 5*time.Millisecond {
		t.Fatalf("recorded %q %v, want %q %v", gotPath, gotD, "/orders", 5*time.Millisecond)
	}
}
