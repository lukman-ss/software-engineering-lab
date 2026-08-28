package payment

import (
	"context"
	"time"

	apperrors "github.com/lukman/software-engineer-lab/pkg/errors"
)

var (
	ErrPaymentNotFound = apperrors.ErrPaymentNotFound
	ErrIdempotencyKey  = apperrors.ErrIdempotencyConflict
	ErrPaymentFailed   = apperrors.ErrPaymentFailed
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
)

type Payment struct {
	ID             string
	OrderID        string
	Amount         int64
	Status         Status
	IdempotencyKey string
	PaymentMethod  string
	ExternalID     string
	PaidAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Repository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
}

type Service interface {
	ProcessPayment(ctx context.Context, idempotencyKey string, orderID string, amount int64, method string) (*Payment, error)
	GetPaymentStatus(ctx context.Context, id string) (*Payment, error)
}
