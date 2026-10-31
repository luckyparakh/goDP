// Package store keeps invoices in memory. In this version it also declares
// the Invoice type and the Store interface, so billing has to import it.
package store

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned for an invoice id nobody added.
var ErrNotFound = errors.New("invoice not found")

// Invoice is one bill, in cents.
type Invoice struct {
	ID        string
	Cents     int64
	PaidCents int64
}

// Store is what billing needs. It is declared here, next to the code that
// implements it.
type Store interface {
	Invoice(id string) (Invoice, error)
}

// Mem keeps invoices in a map.
type Mem struct{ m map[string]Invoice }

// NewMem returns an empty store.
func NewMem() *Mem { return &Mem{m: map[string]Invoice{}} }

// Add stores an invoice for cents, of which paidCents is already paid.
func (s *Mem) Add(id string, cents, paidCents int64) {
	s.m[id] = Invoice{ID: id, Cents: cents, PaidCents: paidCents}
}

// Invoice returns the invoice stored under id.
func (s *Mem) Invoice(id string) (Invoice, error) {
	inv, ok := s.m[id]
	if !ok {
		return Invoice{}, fmt.Errorf("invoice %q: %w", id, ErrNotFound)
	}
	return inv, nil
}
