package checkout

import (
	"context"
	"fmt"

	"circuitbreaker/internal/circuitbreaker"
	"circuitbreaker/internal/payment"
)

type Service struct {
	paymentClient *payment.Client
	cb            *circuitbreaker.CircuitBreaker
}

func NewService(paymentClient *payment.Client, cb *circuitbreaker.CircuitBreaker) *Service {
	return &Service{
		paymentClient: paymentClient,
		cb:            cb,
	}
}

func (s *Service) Checkout(ctx context.Context) error {
	if s.cb != nil {
		err := s.cb.Execute(func() error {
			return s.paymentClient.ProcessPayment(ctx)
		})
		if err != nil {
			return fmt.Errorf("checkout payment failed (with CB): %w", err)
		}
		return nil
	}

	// Without Circuit Breaker
	err := s.paymentClient.ProcessPayment(ctx)
	if err != nil {
		return fmt.Errorf("checkout payment failed (no CB): %w", err)
	}
	return nil
}

func (s *Service) CircuitState() string {
	if s.cb == nil {
		return "N/A"
	}
	return string(s.cb.State())
}
