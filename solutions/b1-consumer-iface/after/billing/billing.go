//go:build ignore

// Package billing works out what customers still owe. It owns the Invoice
// type and declares the one thing it needs from storage.
package billing

import "fmt"

// Invoice is one bill, in cents.
type Invoice struct {
	ID        string
	Cents     int64
	PaidCents int64
}

// Invoices is what billing needs from whoever stores invoices.
type Invoices interface {
	Invoice(id string) (Invoice, error)
}

// Service works out amounts due.
type Service struct{ invoices Invoices }

// New returns a Service that reads invoices from inv.
func New(inv Invoices) *Service { return &Service{invoices: inv} }

// Due returns what is still owed on an invoice, in cents.
func (svc *Service) Due(id string) (int64, error) {
	inv, err := svc.invoices.Invoice(id)
	if err != nil {
		return 0, fmt.Errorf("due: %w", err)
	}
	return inv.Cents - inv.PaidCents, nil
}
