// Package orders works out order totals. Your work: see the GOAL in
// ../../after.go. This file starts as a copy of before/orders, plus the two
// domain errors, which nothing returns yet.
package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound means there is no order with that id.
var ErrNotFound = errors.New("order not found")

// ErrUnavailable means orders could not be read right now. Trying again later
// may work.
var ErrUnavailable = errors.New("orders unavailable")

// Service reads order totals from a database.
type Service struct{ db *sql.DB }

// New returns a Service that queries db.
func New(db *sql.DB) *Service { return &Service{db: db} }

// Total returns the order total in cents.
func (s *Service) Total(ctx context.Context, id string) (int64, error) {
	var cents int64
	err := s.db.QueryRowContext(ctx, "SELECT cents FROM orders WHERE id = ?", id).Scan(&cents)
	if err != nil {
		return 0, fmt.Errorf("total %q: %w", id, err)
	}
	return cents, nil
}
