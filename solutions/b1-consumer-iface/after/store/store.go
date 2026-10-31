//go:build ignore

// Package store keeps invoices in memory. It imports billing, because it
// hands back billing's own Invoice type.
package store

import (
	"errors"
	"fmt"

	"github.com/luckyparakh/goDP/patterns/b1-consumer-iface/after/billing"
)

// ErrNotFound is returned for an invoice id nobody added.
var ErrNotFound = errors.New("invoice not found")

// Mem keeps invoices in a map. It satisfies billing.Invoices.
type Mem struct{ m map[string]billing.Invoice }

// NewMem returns an empty store.
func NewMem() *Mem { return &Mem{m: map[string]billing.Invoice{}} }

// Add stores an invoice for cents, of which paidCents is already paid.
func (s *Mem) Add(id string, cents, paidCents int64) {
	s.m[id] = billing.Invoice{ID: id, Cents: cents, PaidCents: paidCents}
}

// Invoice returns the invoice stored under id.
func (s *Mem) Invoice(id string) (billing.Invoice, error) {
	inv, ok := s.m[id]
	if !ok {
		return billing.Invoice{}, fmt.Errorf("invoice %q: %w", id, ErrNotFound)
	}
	return inv, nil
}
