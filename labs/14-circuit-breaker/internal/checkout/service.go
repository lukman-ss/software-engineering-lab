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
	HasBreaker   bool
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
		return Result{Err: err, Duration: time.Since(start), State: s.breaker.State(), HasBreaker: true}
	}
	return Result{Duration: time.Since(start), State: s.breaker.State(), DownstreamOK: true, HasBreaker: true}
}

// CheckoutWithoutBreaker executes a payment call directly (no circuit breaker).
func (s *Service) CheckoutWithoutBreaker() Result {
	start := time.Now()
	err := s.payment.ProcessPayment(context.Background())
	return Result{Err: err, Duration: time.Since(start), HasBreaker: false}
}

// Checkout is the integration test method - uses circuit breaker.
func (s *Service) Checkout(ctx context.Context) error {
	return s.breaker.Execute(func() error {
		return s.payment.ProcessPayment(ctx)
	})
}

func (r Result) String() string {
	stateStr := ""
	if r.HasBreaker {
		stateStr = fmt.Sprintf(" state=%s", r.State)
	}
	if r.Err != nil {
		return fmt.Sprintf("err=%s duration=%s%s", r.Err, r.Duration, stateStr)
	}
	return fmt.Sprintf("err=nil duration=%s%s", r.Duration, stateStr)
}
