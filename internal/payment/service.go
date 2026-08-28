package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/lukman/software-engineer-lab/pkg/util"
)

type appService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &appService{repo: repo}
}

func (s *appService) ProcessPayment(ctx context.Context, idempotencyKey, orderID string, amount int64, method string) (*Payment, error) {
	// Check idempotency key first
	existing, err := s.repo.GetByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("check idempotency: %w", err)
	}
	if existing != nil {
		return existing, nil
	}

	// Create payment
	now := util.Now()
	payment := &Payment{
		ID:             util.NewPaymentID(),
		OrderID:        orderID,
		Amount:         amount,
		Status:         StatusPending,
		IdempotencyKey: idempotencyKey,
		PaymentMethod:  method,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	// Simulate payment processing (mock gateway call)
	// In production, this calls external payment gateway
	payment.Status = StatusSucceeded
	paidAt := time.Now()
	payment.PaidAt = &paidAt
	payment.ExternalID = "ext_" + payment.ID

	if err := s.repo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("update payment: %w", err)
	}

	return payment, nil
}

func (s *appService) GetPaymentStatus(ctx context.Context, id string) (*Payment, error) {
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
