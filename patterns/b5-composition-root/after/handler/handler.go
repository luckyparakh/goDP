// Package handler serves order totals over HTTP. Your work: see the GOAL in
// ../../after.go.
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/luckyparakh/goDP/patterns/b5-composition-root/after/service"
)

// Service is what Handler needs.
type Service interface {
	Total(ctx context.Context, id string) (int64, error)
}

// Handler answers GET /orders/{id}.
type Handler struct{ svc Service }

// New returns a Handler that asks svc for totals.
func New(svc Service) *Handler {
	// TODO: keep svc. Right now it is ignored, and ServeHTTP builds its own
	// service.
	return &Handler{}
}

// ServeHTTP writes the order total in cents.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: use h.svc instead of building a service here.
	own := service.New(nil)

	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	cents, err := own.Total(r.Context(), id)
	if err != nil {
		// Every error is a 404 here, to keep this task about wiring. Telling
		// errors apart is B2.
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprint(w, cents)
}
