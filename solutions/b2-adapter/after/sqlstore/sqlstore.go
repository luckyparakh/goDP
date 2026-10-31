//go:build ignore

// Package sqlstore adapts a *sql.DB to what orders needs, and translates
// database errors into orders' own errors.
package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/luckyparakh/goDP/patterns/b2-adapter/after/orders"
)

// Store reads orders from a SQL database. It satisfies orders.Store.
type Store struct{ db *sql.DB }

// New returns a Store that queries db.
func New(db *sql.DB) *Store { return &Store{db: db} }

// Cents returns the total for id, with errors in orders' terms.
func (s *Store) Cents(ctx context.Context, id string) (int64, error) {
	var cents int64
	err := s.db.QueryRowContext(ctx, "SELECT cents FROM orders WHERE id = ?", id).Scan(&cents)
	switch {
	case err == nil:
		return cents, nil
	case errors.Is(err, sql.ErrNoRows):
		return 0, fmt.Errorf("order %q: %w", id, orders.ErrNotFound)
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return 0, fmt.Errorf("order %q: %w", id, err)
	default:
		// %v on purpose: keep the driver's text, hide its identity.
		return 0, fmt.Errorf("order %q: %w: %v", id, orders.ErrUnavailable, err)
	}
}
