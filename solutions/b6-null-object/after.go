//go:build ignore

package nullobject

import "log/slog"

// HOW IT WORKS
//
// New is the one place that knows a dependency may be missing. A nil logger
// is replaced by slog.New(slog.DiscardHandler), the standard library's own
// null object. Nil metrics is replaced by nopMetrics, an empty struct whose
// Inc does nothing. After New, both fields always hold something usable, so
// Charge and Refund call them with no checks at all.
//
// WHY IT HELPS
//
// The checks drop from four, in every method, to two, in one constructor.
// A new method cannot forget a check, because there is none to forget. The
// cost is a call to an empty method, and for the logger not even that much
// work: DiscardHandler.Enabled returns false, so slog drops the call before
// it formats anything. The benchmark shows what that costs on your machine.
//
// WHEN NOT TO USE
//
// If the dependency is required, do not paper over its absence: missing
// metrics in production should fail at startup, not be silently dropped. A
// null object for something that returns data, such as a store that always
// says "not found", hides real bugs. And watch for a typed nil: New(p) where p
// is a nil *recorder is not == nil, so it gets past the check. Either refuse
// it, or make the method safe to call on a nil receiver.

type nopMetrics struct{}

func (nopMetrics) Inc(string) {}

// Service takes payments. Its logger and metrics are optional.
type Service struct {
	log     *slog.Logger
	metrics Metrics
}

// New returns a Service. log and m may be nil.
func New(log *slog.Logger, m Metrics) *Service {
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}
	if m == nil {
		m = nopMetrics{}
	}
	return &Service{log: log, metrics: m}
}

// Charge takes a payment, in cents.
func (s *Service) Charge(cents int64) error {
	if cents <= 0 {
		s.metrics.Inc("charge_rejected")
		return ErrBadAmount
	}
	s.log.Info("charge", "cents", cents)
	s.metrics.Inc("charges")
	return nil
}

// Refund gives a payment back, in cents.
func (s *Service) Refund(cents int64) error {
	if cents <= 0 {
		return ErrBadAmount
	}
	s.log.Info("refund", "cents", cents)
	s.metrics.Inc("refunds")
	return nil
}
