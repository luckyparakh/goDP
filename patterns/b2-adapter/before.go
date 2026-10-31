package adapter

// PROBLEM
//
// Look in before/orders. The orders service runs SQL on a *sql.DB itself and
// wraps whatever the database returns with %w. So the database's errors leak
// out to every caller:
//
//   - To find out an order does not exist, a caller must check sql.ErrNoRows,
//     which means importing database/sql. The HTTP layer now knows you use SQL.
//   - A caller can reach into the driver's own error type with errors.As, and
//     start depending on its codes.
//
// Swap the database, or the driver, and every one of those callers breaks.
//
// TRY IT
//
// Read the three TestBefore_* tests in adapter_test.go. Each one is a caller
// that has to know about database/sql or the driver to do its job.
//
// fakedb/ is a tiny stand-in driver so the tests need no real database. You do
// not need to read it.
//
// Do not edit before/. Your work goes in after/.
