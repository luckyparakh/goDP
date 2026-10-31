//go:build ignore

// Package orders works out order totals. It knows nothing about SQL: it
// declares what it needs from storage and the errors it understands.
package orders

import (
	"context"
	"errors"
	"fmt"
)

// ErrNotFound means there is no order with that id.
var ErrNotFound = errors.New("order not found")

// ErrUnavailable means orders could not be read right now. Trying again later
// may work.
var ErrUnavailable = errors.New("orders unavailable")

// Store is what orders needs from storage. Errors it returns must wrap
// ErrNotFound, ErrUnavailable or a context error.
type Store interface {
	Cents(ctx context.Context, id string) (int64, error)
}

// Service works out order totals.
type Service struct{ store Store }

// New returns a Service that reads from st.
func New(st Store) *Service { return &Service{store: st} }

// Total returns the order total in cents.
func (s *Service) Total(ctx context.Context, id string) (int64, error) {
	cents, err := s.store.Cents(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("total: %w", err)
	}
	return cents, nil
}
