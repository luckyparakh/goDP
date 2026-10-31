// Package service works out order totals. In this version it builds its own
// store.
package service

import (
	"context"
	"fmt"

	"github.com/luckyparakh/goDP/patterns/b5-composition-root/before/store"
)

// Service works out order totals.
type Service struct{ st *store.File }

// New builds a Service, and the store it needs.
func New() (*Service, error) {
	st, err := store.New()
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}
	return &Service{st: st}, nil
}

// Total returns the order total in cents.
func (s *Service) Total(ctx context.Context, id string) (int64, error) {
	c, err := s.st.Cents(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("total: %w", err)
	}
	return c, nil
}
