package nullobject

import "log/slog"

// GOAL
//
// Service below starts as a copy of ServiceBefore, nil checks included, and
// New ignores its arguments. Make the nil checks unnecessary.
//
//   - New keeps what it is given. A nil logger becomes a logger that discards
//     everything: the standard library already has one, see
//     go doc log/slog.DiscardHandler. Nil metrics becomes a Metrics of your
//     own that does nothing.
//   - After New, no field is ever compared with nil. Charge and Refund call
//     s.log and s.metrics directly.
//   - Keep the field names log and metrics, and keep the behaviour of
//     ServiceBefore: the same errors, log messages and metric names.
//
// You are done when every test in this folder passes:
//
//	go test ./patterns/b6-null-object/
//
// Then measure how much the no-op costs, and paste the numbers in THEORY.md:
//
//	go test -bench . -benchmem ./patterns/b6-null-object/

// Service takes payments. Its logger and metrics are optional.
type Service struct {
	log     *slog.Logger
	metrics Metrics
}

// New returns a Service. log and m may be nil.
func New(log *slog.Logger, m Metrics) *Service {
	// TODO: use log and m. Right now they are ignored.
	return &Service{}
}

// Charge takes a payment, in cents.
func (s *Service) Charge(cents int64) error {
	if cents <= 0 {
		if s.metrics != nil {
			s.metrics.Inc("charge_rejected")
		}
		return ErrBadAmount
	}
	if s.log != nil {
		s.log.Info("charge", "cents", cents)
	}
	if s.metrics != nil {
		s.metrics.Inc("charges")
	}
	return nil
}

// Refund gives a payment back, in cents.
func (s *Service) Refund(cents int64) error {
	if cents <= 0 {
		return ErrBadAmount
	}
	if s.log != nil {
		s.log.Info("refund", "cents", cents)
	}
	s.metrics.Inc("refunds")
	return nil
}
