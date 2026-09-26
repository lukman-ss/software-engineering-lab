## Snippet 1 — Transactional Order Creation

Source File: `internal/outbox/service.go:18-53`
Purpose: Demonstrates atomic persistence of order and outbox message within a single transaction.

```go
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
```

Explanation: Initiates a transaction, stages order and outbox message, then commits atomically. If any error occurs before commit, rollback discards all staged changes.

## Snippet 2 — Transaction Commit

Source File: `internal/outbox/db.go:99-116`
Purpose: Atomic write of both orders and outbox messages to the database.

```go
func (tx *Tx) Commit() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.closed = true

	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()
	for k, v := range tx.stagedOrders {
		tx.db.orders[k] = v
	}
	for k, v := range tx.stagedOutbox {
		tx.db.outbox[k] = v
	}
	return nil
}
```

Explanation: Under database mutex, all staged orders and outbox messages are written atomically. Both maps are updated together or not at all.

## Snippet 3 — Transaction Rollback

Source File: `internal/outbox/db.go:118-127`
Purpose: Discards staged mutations without persisting any changes.

```go
func (tx *Tx) Rollback() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.closed = true
	// discard staged mutations
	return nil
}
```

Explanation: Marks transaction as closed and discards staged changes. No data is written to the database.

## Snippet 4 — Message Relay Polling

Source File: `internal/outbox/relay.go:43-59`
Purpose: Polls pending outbox messages, publishes to broker, and marks as processed.

```go
func (r *Relay) PollAndDispatch() int {
	pending := r.db.GetPendingOutbox()
	dispatched := 0
	for _, msg := range pending {
		err := r.broker.Publish(msg)
		if err == nil {
			err = r.db.MarkOutboxProcessed(msg.ID)
			if err != nil {
				log.Printf("failed to mark outbox msg %s as processed: %v\n", msg.ID, err)
			} else {
				dispatched++
			}
		} else {
			log.Printf("failed to publish outbox msg %s: %v\n", msg.ID, err)
		}
	}
	return dispatched
}
```

Explanation: Retrieves all pending messages, publishes each to broker, and updates status to PROCESSED on success. Failed publishes leave status as PENDING for retry.

## Snippet 5 — Idempotent Consumer

Source File: `internal/outbox/consumer.go:19-31`
Purpose: Handles messages idempotently by tracking processed event IDs.

```go
func (c *Consumer) Handle(msg OutboxMessage) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.processedIDs[msg.ID] {
		// Duplicate detected, ignore
		return false
	}

	c.processedIDs[msg.ID] = true
	c.received = append(c.received, msg)
	return true
}
```

Explanation: Returns `true` if message is new and processed, `false` if duplicate. This enables at-least-once delivery without side effects from duplicates.

## Snippet 6 — Rollback Test

Source File: `tests/outbox_test.go:61-90`
Purpose: Verifies that rollback discards both order and outbox records.

```go
func TestTransactionalOutbox_Rollback(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 10*time.Millisecond)

	relay.Start()
	defer relay.Stop()

	tx := db.BeginTx()
	_ = tx.SaveOrder(outbox.Order{ID: "o-2", Amount: 50})
	_ = tx.SaveOutbox(outbox.OutboxMessage{ID: "evt-o-2", Status: outbox.MessageStatusPending})

	_ = tx.Rollback()

	_, ok := db.GetOrder("o-2")
	if ok {
		t.Fatalf("expected order to not be saved")
	}
	_, ok = db.GetOutbox("evt-o-2")
	if ok {
		t.Fatalf("expected outbox to not be saved")
	}

	time.Sleep(30 * time.Millisecond)
	if len(broker.GetPublished()) > 0 {
		t.Fatalf("expected no messages sent to broker")
	}
}
```

Explanation: Creates a transaction, stages both order and outbox, then rolls back. Asserts that neither record exists in the database and no message was published.

## Snippet 7 — Duplicate Delivery Test

Source File: `tests/outbox_test.go:92-115`
Purpose: Verifies consumer rejects duplicate message deliveries.

```go
func TestTransactionalOutbox_Idempotency_DuplicateDelivery(t *testing.T) {
	consumer := outbox.NewConsumer()

	msg := outbox.OutboxMessage{
		ID:        "evt-duplicate",
		EventType: "OrderCreated",
		Payload:   "data",
	}

	p1 := consumer.Handle(msg)
	if !p1 {
		t.Fatalf("expected first delivery to be processed")
	}

	p2 := consumer.Handle(msg)
	if p2 {
		t.Fatalf("expected duplicate delivery to be rejected by consumer")
	}

	if consumer.GetReceivedCount() != 1 {
		t.Fatalf("expected exactly 1 message processed despite duplicate delivery")
	}
}
```

Explanation: First call processes message, second call rejects it as duplicate. Consumer count remains 1, proving deduplication works.

## Snippet 8 — Dual-Write Failure Demonstration

Source File: `tests/outbox_test.go:117-139`
Purpose: Demonstrates state inconsistency when dual-write fails.

```go
func TestDualWriteProblem_Failure(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	service := outbox.NewOrderService(db)

	broker.SetFailNext(true)

	err := service.CreateOrderDualWriteNaive(broker, "o-bug", "c-2", 200)
	if err == nil {
		t.Fatalf("expected error from dual write naive approach")
	}

	_, ok := db.GetOrder("o-bug")
	if !ok {
		t.Fatalf("expected order to be saved in DB")
	}

	if len(broker.GetPublished()) != 0 {
		t.Fatalf("expected no message published to broker")
	}
}
```

Explanation: Forces broker failure after DB commit. Result: order exists in DB, but no message was published to broker — proving the dual-write problem.

## Snippet 9 — Dual-Write Naive Implementation

Source File: `internal/outbox/service.go:55-90`
Purpose: Demonstrates the vulnerable dual-write approach (for comparison).

```go
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
		return fmt.Errorf("failed to publish to broker after DB commit: %w", err)
	}

	return nil
}
```

Explanation: Commits order to DB first, then attempts broker publish outside the transaction. If broker fails, DB has committed data but event is lost.

## Snippet 10 — Relay Startup with Goroutine

Source File: `internal/outbox/relay.go:24-37`
Purpose: Starts the polling relay as a background goroutine.

```go
func (r *Relay) Start() {
	go func() {
		ticker := time.NewTicker(r.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.PollAndDispatch()
			case <-r.stopChan:
				return
			}
		}
	}()
}
```

Explanation: Spawns a goroutine that polls at fixed interval and dispatches messages. Clean shutdown via stop channel.

## Snippet 11 — Outbox Message Model

Source File: `internal/outbox/model.go:26-32`
Purpose: Defines the outbox message data structure.

```go
type OutboxMessage struct {
	ID        string
	EventType string
	Payload   string
	Status    MessageStatus
	CreatedAt time.Time
}
```

Explanation: Contains unique ID, event type, JSON payload, status (PENDING/PROCESSED), and creation timestamp.

## Snippet 12 — Order Model

Source File: `internal/outbox/model.go:12-17`
Purpose: Defines the order domain model.

```go
type Order struct {
	ID         string
	CustomerID string
	Amount     float64
	Status     OrderStatus
}
```

Explanation: Simple order entity with ID, customer reference, amount, and status.