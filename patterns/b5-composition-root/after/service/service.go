// Package service works out order totals. Your work: see the GOAL in
// ../../after.go.
package service

import (
	"context"
	"fmt"

	"github.com/luckyparakh/goDP/patterns/b5-composition-root/after/store"
)

// Store is what Service needs.
type Store interface {
	Cents(ctx context.Context, id string) (int64, error)
}

// Service works out order totals.
type Service struct{ st Store }

// New returns a Service that reads totals from st.
func New(st Store) *Service {
	// TODO: keep st. Right now it is ignored, and Total builds its own store.
	return &Service{}
}

// Total returns the order total in cents.
func (s *Service) Total(ctx context.Context, id string) (int64, error) {
	// TODO: use s.st. This builds a store on every call, the way before/
	// builds one in New, because New is no longer allowed to fail.
	own, err := store.New("")
	if err != nil {
		return 0, fmt.Errorf("total: %w", err)
	}
	c, err := own.Cents(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("total: %w", err)
	}
	return c, nil
}
