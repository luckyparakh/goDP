// Package orders works out order totals. In this version it runs the SQL
// itself and passes database errors straight up to its callers.
package orders

import (
	"context"
	"database/sql"
	"fmt"
)

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
