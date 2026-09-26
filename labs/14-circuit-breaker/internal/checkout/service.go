package checkout

import (
	"circuitbreaker/internal/circuitbreaker"
	"circuitbreaker/internal/payment"
	"context"
	"fmt"
	"time"
)

type Result struct {
	Err          error
	Duration     time.Duration
	DownstreamOK bool
	State        circuitbreaker.State
}

type Service struct {
	payment *payment.Client
	breaker *circuitbreaker.Breaker
}

func New(payClient *payment.Client, breaker *circuitbreaker.Breaker) *Service {
	return &Service{payment: payClient, breaker: breaker}
}

// NewService is the integration test constructor (alias for New).
func NewService(payClient *payment.Client, breaker *circuitbreaker.Breaker) *Service {
	return New(payClient, breaker)
}

// CheckoutWithBreaker executes a payment call through the circuit breaker.
func (s *Service) CheckoutWithBreaker() Result {
	start := time.Now()
	if err := s.breaker.Execute(func() error {
		return s.payment.ProcessPayment(context.Background())
	}); err != nil {
		return Result{Err: err, Duration: time.Since(start), State: s.breaker.State()}
	}
	return Result{Duration: time.Since(start), State: s.breaker.State(), DownstreamOK: true}
}

// CheckoutWithoutBreaker executes a payment call directly (no circuit breaker).
func (s *Service) CheckoutWithoutBreaker() Result {
	start := time.Now()
	err := s.payment.ProcessPayment(context.Background())
	return Result{Err: err, Duration: time.Since(start)}
}

// Checkout is the integration test method - uses circuit breaker.
func (s *Service) Checkout(ctx context.Context) error {
	return s.breaker.Execute(func() error {
		return s.payment.ProcessPayment(ctx)
	})
}

func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("err=%s duration=%s state=%s", r.Err, r.Duration, r.State)
	}
	return fmt.Sprintf("err=nil duration=%s state=%s", r.Duration, r.State)
}
