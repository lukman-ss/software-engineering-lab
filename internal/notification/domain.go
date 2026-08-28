package notification

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotificationFailed = errors.New("failed to send notification")
)

type Type string

const (
	TypeOrderCreated   Type = "order_created"
	TypePaymentSuccess Type = "payment_success"
	TypePaymentFailed  Type = "payment_failed"
)

type Notification struct {
	ID        string
	UserID    string
	Type      Type
	Payload   string
	SentAt    *time.Time
	CreatedAt time.Time
}

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	MarkAsSent(ctx context.Context, id string) error
}

// Service sends notifications (e.g. email, SMS).
// In the mini production system, this simulates an external call.
type Service interface {
	NotifyOrderCreated(ctx context.Context, userID string, orderID string) error
	NotifyPaymentResult(ctx context.Context, userID string, orderID string, success bool) error
}
