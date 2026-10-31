// Package billing works out what customers still owe. In this version it
// imports store to get both the Invoice type and the interface it calls.
package billing

import (
	"fmt"

	"github.com/luckyparakh/goDP/patterns/b1-consumer-iface/before/store"
)

// Service works out amounts due.
type Service struct{ invoices store.Store }

// New returns a Service that reads invoices from s.
func New(s store.Store) *Service { return &Service{invoices: s} }

// Due returns what is still owed on an invoice, in cents.
func (svc *Service) Due(id string) (int64, error) {
	inv, err := svc.invoices.Invoice(id)
	if err != nil {
		return 0, fmt.Errorf("due: %w", err)
	}
	return inv.Cents - inv.PaidCents, nil
}
