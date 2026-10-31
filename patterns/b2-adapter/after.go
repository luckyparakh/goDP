package adapter

// GOAL
//
// Split after/orders into a service that knows nothing about SQL, and an
// adapter in after/sqlstore that knows everything about it.
//
//   - orders declares the one-method interface it needs from storage, and
//     returns only its own errors: ErrNotFound and ErrUnavailable. It must not
//     import database/sql, even indirectly.
//   - sqlstore holds the *sql.DB, runs the query, and translates errors at the
//     boundary:
//     no rows becomes orders.ErrNotFound;
//     context.Canceled and context.DeadlineExceeded pass through unchanged;
//     anything else becomes orders.ErrUnavailable. Keep the driver's text in
//     the message for the logs, but callers must not be able to reach the
//     driver's error with errors.Is or errors.As.
//   - orders.New(sqlstore.New(db)) must keep compiling, so the tests do too.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/b2-adapter/...
//
// Check the import rule yourself:
//
//	go list -deps ./patterns/b2-adapter/after/orders | grep database/sql
//
// It should print nothing.
