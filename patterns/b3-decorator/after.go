package decorator

import (
	"io"
	"net/http"
	"time"
)

// GOAL
//
// Write each job once, as a wrapper, and put one chain of wrappers in front of
// every route.
//
//   - A Middleware takes a handler and returns a handler.
//   - Chain wraps h so the first middleware listed is the outermost: it sees
//     the request first and the response last.
//   - Logging writes one line per request: method, path and the status code
//     the inner handler wrote, such as "GET /orders 401".
//   - Auth lets a request through only if its X-Token header equals token.
//     Otherwise it answers 401, and the inner handler never runs.
//   - Timing calls now before and after the inner handler and passes the path
//     and the difference to record. Taking now as a parameter keeps the test
//     free of real clocks.
//   - Routes serves GET /orders, GET /report and GET /export. Each writes its
//     own name as the body. All three sit behind one chain: Logging, then Auth,
//     then Timing.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/b3-decorator/
//
// Then do "Read three" in THEORY.md.

// Middleware wraps a handler and returns a handler.
type Middleware func(http.Handler) http.Handler

// Chain wraps h in mws. The first one listed is the outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	// TODO
	return h
}

// Logging writes "METHOD PATH STATUS" to w for every request.
func Logging(w io.Writer) Middleware {
	// TODO
	return func(h http.Handler) http.Handler { return h }
}

// Auth answers 401 unless the X-Token header equals token.
func Auth(token string) Middleware {
	// TODO
	return func(h http.Handler) http.Handler { return h }
}

// Timing records how long the inner handler took, measured with now.
func Timing(now func() time.Time, record func(path string, d time.Duration)) Middleware {
	// TODO
	return func(h http.Handler) http.Handler { return h }
}

// Routes serves /orders, /report and /export behind Logging, Auth and Timing.
func Routes(token string, log io.Writer, now func() time.Time, record func(path string, d time.Duration)) http.Handler {
	// TODO
	return http.NotFoundHandler()
}
