package codereview

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEmptyCart           = errors.New("cart is empty")
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidQuantity     = errors.New("item quantity must be greater than zero")
	ErrInsufficientStock   = errors.New("insufficient stock for product")
	ErrIdempotencyConflict = errors.New("idempotency key conflict: request payload mismatch")
	ErrDuplicateRequest    = errors.New("request already processed or in progress")
)

type CartItem struct {
	ProductID string
	Quantity  int
	UnitPrice int64
}

type Product struct {
	ID        string
	Name      string
	Stock     int
	UnitPrice int64
}

type Order struct {
	ID             string
	UserID         string
	IdempotencyKey string
	TotalAmount    int64
	Items          []OrderItem
	CreatedAt      time.Time
}

type OrderItem struct {
	ProductID string
	Quantity  int
	UnitPrice int64
	Subtotal  int64
}

type CheckoutResponse struct {
	Success bool
	OrderID string
	Total   int64
}

type Logger interface {
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

type NotificationSender interface {
	SendOrderConfirmation(ctx context.Context, userID string, orderID string) error
}

type IdempotencyRepository interface {
	TryInsert(ctx context.Context, key string, hash string) (bool, error)
	MarkProcessed(ctx context.Context, key string, resp *CheckoutResponse) error
	Get(ctx context.Context, key string) (*IdempotencyRecord, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
}

type ProductRepository interface {
	GetProduct(ctx context.Context, productID string) (*Product, error)
	UpdateStock(ctx context.Context, productID string, delta int) (int, error)
}

type CartSource interface {
	GetCart(ctx context.Context, userID string) ([]CartItem, error)
}