package outbox

import "time"

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID         string
	CustomerID string
	Amount     float64
	Status     OrderStatus
}

type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "PENDING"
	MessageStatusProcessed MessageStatus = "PROCESSED"
)

type OutboxMessage struct {
	ID        string
	EventType string
	Payload   string
	Status    MessageStatus
	CreatedAt time.Time
}
