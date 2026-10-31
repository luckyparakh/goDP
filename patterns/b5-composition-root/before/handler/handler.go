// Package handler serves order totals over HTTP. In this version it builds its
// own service, which builds its own store.
package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/luckyparakh/goDP/patterns/b5-composition-root/before/service"
)

// Handler answers GET /orders/{id}.
type Handler struct{ svc *service.Service }

// New builds a Handler, and the service it needs.
func New() (*Handler, error) {
	svc, err := service.New()
	if err != nil {
		return nil, fmt.Errorf("handler: %w", err)
	}
	return &Handler{svc: svc}, nil
}

// ServeHTTP writes the order total in cents.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	cents, err := h.svc.Total(r.Context(), id)
	if err != nil {
		// Every error is a 404 here, to keep this task about wiring. Telling
		// errors apart is B2.
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprint(w, cents)
}
