// Package fakedb is a tiny database/sql driver for tests, built only from the
// standard library. You do not need to read or edit it.
//
// It answers any query with one column, "cents", looked up by the first
// argument. A missing id gives no rows, so Scan returns sql.ErrNoRows, exactly
// as a real database would. Broken returns a database whose every query fails
// with a driver-specific *Error, the way a real driver reports a lost
// connection with its own error type and code.
package fakedb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
)

// Error is this driver's own error type, like the error types real drivers
// export. Code is a made-up status code.
type Error struct{ Code string }

func (e *Error) Error() string { return "fakedb: error " + e.Code }

// Open returns a database holding rows, a map from order id to cents.
func Open(rows map[string]int64) *sql.DB { return sql.OpenDB(connector{rows: rows}) }

// Broken returns a database whose every query fails with &Error{Code: code}.
func Broken(code string) *sql.DB { return sql.OpenDB(connector{fail: &Error{Code: code}}) }

type connector struct {
	rows map[string]int64
	fail error
}

func (c connector) Connect(context.Context) (driver.Conn, error) { return conn(c), nil }
func (c connector) Driver() driver.Driver                        { return drv(c) }

type drv connector

func (d drv) Open(string) (driver.Conn, error) { return conn(d), nil }

type conn connector

func (c conn) Prepare(string) (driver.Stmt, error) { return stmt(c), nil }
func (c conn) Close() error                        { return nil }
func (c conn) Begin() (driver.Tx, error)           { return nil, errors.New("fakedb: no transactions") }

type stmt connector

func (s stmt) Close() error  { return nil }
func (s stmt) NumInput() int { return -1 }

func (s stmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, errors.New("fakedb: read only")
}

func (s stmt) Query(args []driver.Value) (driver.Rows, error) {
	if s.fail != nil {
		return nil, s.fail
	}
	var id string
	if len(args) > 0 {
		id, _ = args[0].(string)
	}
	cents, ok := s.rows[id]
	return &rows{cents: cents, left: ok}, nil
}

type rows struct {
	cents int64
	left  bool
}

func (r *rows) Columns() []string { return []string{"cents"} }
func (r *rows) Close() error      { return nil }

func (r *rows) Next(dest []driver.Value) error {
	if !r.left {
		return io.EOF
	}
	dest[0], r.left = r.cents, false
	return nil
}
