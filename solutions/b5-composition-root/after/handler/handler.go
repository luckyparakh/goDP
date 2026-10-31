//go:build ignore

// Package handler serves order totals over HTTP from whatever Service it is
// given.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Service is what Handler needs.
type Service interface {
	Total(ctx context.Context, id string) (int64, error)
}

// Handler answers GET /orders/{id}.
type Handler struct{ svc Service }

// New returns a Handler that asks svc for totals.
func New(svc Service) *Handler { return &Handler{svc: svc} }

// ServeHTTP writes the order total in cents.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	cents, err := h.svc.Total(r.Context(), id)
	if err != nil {
		// Every error is a 404 here, to keep this task about wiring.
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprint(w, cents)
}
