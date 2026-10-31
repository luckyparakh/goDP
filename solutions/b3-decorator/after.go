//go:build ignore

package decorator

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HOW IT WORKS
//
// Each job is a func that takes a handler and returns a new handler that does
// the job around it. Chain applies them from the last to the first, so the
// first one listed ends up outermost: its code runs first on the way in and
// last on the way out. Routes wraps the whole mux once, so every route, and
// every route added later, sits behind the same chain.
//
// Logging needs the status code, but the inner handler is the one that writes
// it. So Logging wraps the ResponseWriter too: statusWriter records the code
// as it passes through. That is a second decorator, on the response side. It
// has an Unwrap method so http.ResponseController can still reach the real
// writer for things like Flush.
//
// WHY IT HELPS
//
// Each job is written once. A route cannot forget auth, because routes no
// longer do auth. The order is one line in Routes, and the tests prove it:
// Logging sits outside Auth, so rejected requests are logged.
//
// WHEN NOT TO USE
//
// If a job differs per route, a shared wrapper hides a one-off behind
// indirection; put it in that handler. And keep chains short: a stack of ten
// wrappers around three lines of work is hard to follow in a stack trace.

// Middleware wraps a handler and returns a handler.
type Middleware func(http.Handler) http.Handler

// Chain wraps h in mws. The first one listed is the outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// statusWriter remembers the status code the inner handler wrote.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Unwrap lets http.ResponseController find the original writer.
func (s *statusWriter) Unwrap() http.ResponseWriter { return s.ResponseWriter }

// Logging writes "METHOD PATH STATUS" to w for every request.
func Logging(w io.Writer) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: rw, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			fmt.Fprintf(w, "%s %s %d\n", r.Method, r.URL.Path, sw.status)
		})
	}
}

// Auth answers 401 unless the X-Token header equals token.
func Auth(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Token") != token {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Timing records how long the inner handler took, measured with now.
func Timing(now func() time.Time, record func(path string, d time.Duration)) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := now()
			next.ServeHTTP(w, r)
			record(r.URL.Path, now().Sub(start))
		})
	}
}

// Routes serves /orders, /report and /export behind Logging, Auth and Timing.
func Routes(token string, log io.Writer, now func() time.Time, record func(path string, d time.Duration)) http.Handler {
	mux := http.NewServeMux()
	for _, name := range []string{"orders", "report", "export"} {
		mux.Handle("GET /"+name, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, name)
		}))
	}
	return Chain(mux, Logging(log), Auth(token), Timing(now, record))
}
