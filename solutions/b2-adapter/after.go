//go:build ignore

package adapter

// HOW IT WORKS
//
// orders declares a one-method Store interface and its own two errors, and
// imports nothing about SQL. sqlstore is the adapter: it holds the *sql.DB,
// has the method orders asks for, and translates at the boundary.
//
//   - sql.ErrNoRows becomes orders.ErrNotFound. The SQL error carries nothing
//     a caller needs, so it is dropped.
//   - context.Canceled and context.DeadlineExceeded pass through with %w. They
//     are the caller's own signal (A7), not a storage detail.
//   - Anything else becomes orders.ErrUnavailable with %w, and the driver's
//     error is added with %v. Its text stays in the message for the logs, but
//     errors.Is and errors.As cannot reach it, so no caller can depend on it.
//
// sqlstore imports orders, as B1 taught: the provider uses the consumer's
// errors, so the arrow points from the adapter to the service.
//
// The solution files are after/orders/orders.go and after/sqlstore/sqlstore.go.
//
// WHY IT HELPS
//
// Callers check orders.ErrNotFound and never import database/sql. Swap the
// database, or the driver, and only sqlstore changes. The service can be tested
// with a map that returns orders.ErrNotFound, no driver at all.
//
// WHEN NOT TO USE
//
// A small tool over one database you will never swap, where nothing above
// checks the error kind, gains only a layer to read through. And do not hide
// errors a caller genuinely needs: translate them into your own words instead.
