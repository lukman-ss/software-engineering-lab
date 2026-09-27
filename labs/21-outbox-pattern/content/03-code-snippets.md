# Code Snippets

## Snippet 1 — Atomic Order Creation with Outbox

Source File: `internal/outbox/service.go:17-53`

Purpose: Demonstrating that both the `Order` and the `OutboxMessage` are written within the same transactional boundary (`BeginTx`, `Commit`, `Rollback`). Rollback is triggered if `json.Marshal` fails or `SaveOrder`/`SaveOutbox` returns an error.

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

---

## Snippet 2 — Naive Dual-Write (Failure-Prone)

Source File: `internal/outbox/service.go:55-90`

Purpose: Demonstrating the dual-write problem by committing the database transaction first and then publishing to the broker outside the transaction boundary. When `broker.Publish` fails, the order persists in the database but the corresponding event is never delivered — the lab proves this inconsistency in `TestDualWriteProblem_Failure`.

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
		// Failure here means DB has order, but message broker never received event!
		return fmt.Errorf("failed to publish to broker after DB commit: %w", err)
	}

	return nil
}
```

---

## Snippet 3 — Transactional Commit and Rollback

Source File: `internal/outbox/db.go:84-139`

Purpose: Showing how staged order and outbox changes are accumulated and then flushed atomically to the live maps on `Commit`. `Rollback` simply discards the staged data and marks the transaction closed.

```go
type Tx struct {
	mu           sync.Mutex
	db           *DB
	stagedOrders map[string]Order
	stagedOutbox map[string]OutboxMessage
	closed       bool
}

func (tx *Tx) SaveOrder(order Order) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.stagedOrders[order.ID] = order
	return nil
}

func (tx *Tx) SaveOutbox(msg OutboxMessage) error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.stagedOutbox[msg.ID] = msg
	return nil
}

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

---

## Snippet 4 — Polling Relay Dispatch

Source File: `internal/outbox/relay.go:50-67`

Purpose: Showing how the polling relay fetches all pending outbox records, attempts a publish to the broker for each one, and only marks a record as `PROCESSED` upon successful publish. If publish fails, the record remains `PENDING` until the next poll cycle.

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

---

## Snippet 5 — Idempotent Consumer

Source File: `internal/outbox/consumer.go:19-31`

Purpose: Demonstrating idempotency through event ID tracking. When `Handle` receives a message whose ID is already registered in `processedIDs`, it returns `false` (reject/ignore). Only the first delivery mutates state and returns `true`.

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

---

## Snippet 6 — Outbox Event Model

Source File: `internal/outbox/model.go:12-32`

Purpose: Defining the `Order` domain entity and the `OutboxMessage` event record, plus status enums `OrderStatus` and `MessageStatus`.

```go
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
```

---

## Snippet 7 — Demo Orchestration

Source File: `cmd/demo/main.go:10-63`

Purpose: Orchestrating the three end-to-end scenarios the lab verifies: the dual-write flaw, the atomic outbox success path, and the duplicate-delivery idempotent consumer.

```go
func main() {
	fmt.Println("=== Lab 21: Transactional Outbox Pattern Demo ===")

	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 50*time.Millisecond)
	service := outbox.NewOrderService(db)
	consumer := outbox.NewConsumer()

	// 1. Demonstrate Dual-Write Flaw
	fmt.Println("\n[Scenario 1: The Dual-Write Problem]")
	broker.SetFailNext(true) // broker failure simulation
	err := service.CreateOrderDualWriteNaive(broker, "order-dual-fail", "cust-1", 150.0)
	if err != nil {
		fmt.Printf("Direct write failed: %v\n", err)
	}
	_, dbFound := db.GetOrder("order-dual-fail")
	fmt.Printf("State Inconsistency: Order in DB = %v, Broker Message Count = %d\n", dbFound, len(broker.GetPublished()))

	// 2. Start Message Relay for Outbox
	fmt.Println("\n[Scenario 2: Transactional Outbox Solution]")
	relay.Start()
	defer relay.Stop()

	fmt.Println("Creating order with transactional outbox...")
	err = service.CreateOrderWithOutbox("order-outbox-success", "cust-2", 300.0)
	if err != nil {
		fmt.Printf("Failed to create order: %v\n", err)
		return
	}
	fmt.Println("Order and Outbox record atomically saved to DB.")

	// Let relay poll and dispatch
	time.Sleep(120 * time.Millisecond)

	published := broker.GetPublished()
	fmt.Printf("Broker received messages: %d\n", len(published))
	for _, m := range published {
		fmt.Printf(" - Event ID: %s, Type: %s, Payload: %s\n", m.ID, m.EventType, m.Payload)
		accepted := consumer.Handle(m)
		fmt.Printf(" - Consumer processing initial message: accepted=%v\n", accepted)
	}

	// 3. Demonstrate At-Least-Once & Idempotency
	fmt.Println("\n[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]")
	duplicateMsg := published[0]
	fmt.Println("Simulating duplicate delivery to consumer...")
	acceptedAgain := consumer.Handle(duplicateMsg)
	fmt.Printf("Consumer processing duplicate delivery: accepted=%v (Duplicate safely skipped!)\n", acceptedAgain)
	fmt.Printf("Total events processed by consumer: %d\n", consumer.GetReceivedCount())

	fmt.Println("\n=== Demo Complete ===")
}
```
