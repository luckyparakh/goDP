//go:build ignore

// Package service works out order totals from whatever Store it is given.
package service

import (
	"context"
	"fmt"
)

// Store is what Service needs.
type Store interface {
	Cents(ctx context.Context, id string) (int64, error)
}

// Service works out order totals.
type Service struct{ st Store }

// New returns a Service that reads totals from st.
func New(st Store) *Service { return &Service{st: st} }

// Total returns the order total in cents.
func (s *Service) Total(ctx context.Context, id string) (int64, error) {
	c, err := s.st.Cents(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("total: %w", err)
	}
	return c, nil
}
