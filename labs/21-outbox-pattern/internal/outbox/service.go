package outbox

import (
	"encoding/json"
	"fmt"
	"time"
)

type OrderService struct {
	db *DB
}

func NewOrderService(db *DB) *OrderService {
	return &OrderService{db: db}
}

// CreateOrderWithOutbox writes order and outbox event in the same transaction
func (s *OrderService) CreateOrderWithOutbox(orderID string, customerID string, amount float64) error {
	tx := s.db.BeginTx()

	order := Order{
		ID:         orderID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     OrderStatusCreated,
	}

	payload, err := json.Marshal(order)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to marshal order payload: %w", err)
	}

	msg := OutboxMessage{
		ID:        fmt.Sprintf("evt-%s", orderID),
		EventType: "OrderCreated",
		Payload:   string(payload),
		Status:    MessageStatusPending,
		CreatedAt: time.Now(),
	}

	if err := tx.SaveOrder(order); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.SaveOutbox(msg); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// CreateOrderDualWriteNaive attempts naive dual write (DB write then broker write outside tx).
// Demonstrates dual-write vulnerability when broker fails.
func (s *OrderService) CreateOrderDualWriteNaive(broker Broker, orderID string, customerID string, amount float64) error {
	tx := s.db.BeginTx()
	order := Order{
		ID:         orderID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     OrderStatusCreated,
	}

	if err := tx.SaveOrder(order); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Direct write to broker outside DB transaction
	msg := OutboxMessage{
		ID:        fmt.Sprintf("evt-%s", orderID),
		EventType: "OrderCreated",
		Payload:   fmt.Sprintf(`{"ID":"%s","Amount":%f}`, orderID, amount),
		Status:    MessageStatusPending,
		CreatedAt: time.Now(),
	}

	if err := broker.Publish(msg); err != nil {
		// Failure here means DB has order, but message broker never received event!
		return fmt.Errorf("failed to publish to broker after DB commit: %w", err)
	}

	return nil
}
