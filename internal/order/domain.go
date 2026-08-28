package order

import (
	"context"
	"time"

	apperrors "github.com/lukman/software-engineer-lab/pkg/errors"
)

var (
	ErrOrderNotFound = apperrors.ErrOrderNotFound
	ErrInvalidStatus = apperrors.ErrInvalidStatus
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPaid      Status = "paid"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
	StatusShipped   Status = "shipped"
)

type Order struct {
	ID          string
	UserID      string
	Status      Status
	TotalAmount int64
	Items       []OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice int64
	Subtotal  int64
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
}

type Service interface {
	CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	MarkAsPaid(ctx context.Context, id string) error
	CancelOrder(ctx context.Context, id string) error
}
