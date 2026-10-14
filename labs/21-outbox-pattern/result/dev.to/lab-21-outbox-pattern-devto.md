---
title: "The Transactional Outbox Pattern: Solving the Dual-Write Problem in Go"
published: false
description: "Eliminate silent event loss and state divergence in distributed systems with atomic outbox persistence, polling relays, and idempotent consumers."
tags: "go, architecture, microservices, distributed-systems"
---

Writing to a database and publishing an event to a message broker in separate operations is one of the most frequent sources of silent state corruption in distributed architectures. When an application commits a database transaction and subsequently attempts a network write to Kafka, RabbitMQ, or an event bus, any transient network partition, timeout, or broker restart breaks consistency. The entity is persisted, but downstream consumers never receive the notification.

In **Software Engineering Lab 21**, we examine and implement the core engineering mechanics of the **Transactional Outbox Pattern** in Go: atomic co-location of domain state and event records within a single ACID transaction, an asynchronous polling relay with retry capabilities, and deduplication tracking on downstream consumers to guarantee at-least-once delivery semantics without state drift.

```
Client Request ---> [ Order Service ]
                           |
            (Single DB Transaction: BEGIN ... COMMIT)
                           +---------------------------+
                           |  Domain Table: orders     |
                           |  Outbox Table: outbox     |
                           +---------------------------+
                                         ^
                                         | Poll (SELECT WHERE status='PENDING')
                                 [ Outbox Relay Worker ]
                                         |
                                         v Publish (retryable)
                                 [ Message Broker ]
                                         |
                                         v Deliver
                                [ Idempotent Consumer ]
```

---

## The Dual-Write Problem

Consider a standard e-commerce order creation flow. An application needs to update internal records and emit a notification event to trigger downstream billing and shipment pipelines:

1. `tx.SaveOrder(order)` inside a SQL transaction.
2. `tx.Commit()` succeeds.
3. `broker.Publish(event)` executes over TCP.

If the broker is unreachable, returns an acknowledgment timeout, or the application process crashes immediately following step 2, the order exists in the database, but no downstream event is ever published.

Inverting the sequence—publishing to the broker first and then committing the database transaction—introduces an even worse failure mode: if the database transaction encounters a lock timeout, serialization failure, or unique constraint violation, the database transaction rolls back, but the broker message has already been dispatched. Downstream services act on an order that does not exist.

Distributed two-phase commit (2PC / XA transactions) across heterogeneous databases and modern message brokers is impractical. Modern message brokers such as Apache Kafka and RabbitMQ do not participate in XA transactions, and 2PC introduces prohibitive latency, blocking locks, and fragile single points of failure.

---

## The Mental Model: Single Transaction Boundary

The Transactional Outbox pattern resolves this tension through a straightforward architectural rule:

> **Shift the second write into the same transaction boundary as the first.**

Instead of communicating across two independent network boundaries (database storage engine and messaging infrastructure), the service writes both the domain entity (`orders`) and an outbox record (`outbox_events`) inside the **exact same database transaction**.

Because both writes execute against the same database engine within a single ACID transaction, atomicity is guaranteed by the database's write-ahead log (WAL). Either both records persist, or neither does. A separate asynchronous worker—the message relay—scans the outbox table, dispatches pending records to the message broker, and marks them as processed upon delivery confirmation.

---

## Core Components and Implementation in Go

The reference implementation in `labs/21-outbox-pattern` defines five decoupled components:

### 1. Atomic Transaction Execution

In `internal/outbox/service.go`, `CreateOrderWithOutbox` guarantees that payload serialization, order creation, and outbox staging occur within one transaction:

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

If serialization fails or a database write constraint triggers an error, `tx.Rollback()` cleans up both staged records. No phantom event can escape to the broker.

### 2. The Naïve Dual-Write Anti-Pattern

For contrast, `service.go` implements the vulnerable naive pattern:

```go
func (s *OrderService) CreateOrderDualWriteNaive(broker Broker, orderID string, customerID string, amount float64) error {
	tx := s.db.BeginTx()
	order := Order{ID: orderID, CustomerID: customerID, Amount: amount, Status: OrderStatusCreated}

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

When `broker.Publish` fails, the calling function returns an error, but the order remains committed in the database. The system enters an inconsistent state.

### 3. Asynchronous Polling Relay Worker

The message relay runs out-of-band relative to user requests. In `internal/outbox/relay.go`, the relay queries pending records, dispatches them to the broker interface, and transitions their status:

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

If the broker experiences transient failure, `Publish` returns an error, and the record remains in `PENDING` status. On subsequent polling intervals, the relay retries delivery until successful.

### 4. Downstream Idempotent Consumer

Because network timeouts can occur after the broker receives a message but before the relay updates the outbox table to `PROCESSED`, the outbox pattern provides **at-least-once delivery**, not exactly-once delivery. Downstream consumers must implement idempotency.

In `internal/outbox/consumer.go`, message identifiers are tracked to filter out duplicates:

```go
func (c *Consumer) Handle(msg OutboxMessage) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.processedIDs[msg.ID] {
		// Duplicate detected, safely ignored
		return false
	}

	c.processedIDs[msg.ID] = true
	c.received = append(c.received, msg)
	return true
}
```

---

## Automated Verification Suite

The implementation is verified through 8 automated test cases in `tests/outbox_test.go`:

| Test Case | Scenario Verified | Assertions |
|---|---|---|
| `TestTransactionalOutbox_HappyPath` | Atomic order + outbox creation | Order exists in DB, outbox status transitions from `PENDING` to `PROCESSED`, consumer accepts message |
| `TestTransactionalOutbox_Rollback` | Transaction rollback on failure | 0 orders in DB, 0 outbox records staged, 0 messages sent to broker |
| `TestTransactionalOutbox_Idempotency_DuplicateDelivery` | Duplicate event delivery | First delivery returns `accepted=true`; second delivery with same ID returns `accepted=false`; consumer message count remains 1 |
| `TestDualWriteProblem_Failure` | Naive dual-write behavior | Order persisted in DB (`found=true`), but broker received 0 messages when broker is down |
| `TestTransactionalOutbox_ConcurrentWrites` | 10 concurrent goroutines writing 10 orders each | 100 orders in DB, 100 outbox records processed, 0 pending, 0 race conditions (`-race` clean) |
| `TestTransactionalOutbox_PurgeProcessed` | Retention cleanup routine | `PROCESSED` records purged; `PENDING` records preserved |
| `TestTransactionalOutbox_RelayRetryAfterBrokerFailure` | Transient broker disconnection | First dispatch fails and stays `PENDING`; second dispatch succeeds after broker recovery and transitions to `PROCESSED` |
| `TestTransactionalOutbox_ConcurrentConsumers` | 10 concurrent consumers processing duplicate dispatches | Exactly 1 consumer accepts the message; 9 duplicate attempts rejected |

Run the test suite with data race detection:

```bash
go test -race -v ./...
```

All 8 tests pass without race conditions.

---

## Production Trade-offs and Engineering Boundaries

When adopting the Transactional Outbox pattern in production systems, account for the following architectural constraints:

1. **Relay Architecture: Polling vs CDC**: This lab implements a Polling Publisher. Polling introduces query overhead on the database and latency bounded by the polling ticker. In high-throughput systems, Change Data Capture (CDC) via Debezium reading Postgres Write-Ahead Logs (WAL) or MySQL binlogs dispatches events without polling queries.
2. **Table Purging and Growth**: An outbox table without a retention policy grows indefinitely. In high-velocity platforms emitting millions of events per day, implement table partitioning by date and automate periodic deletion of records in `PROCESSED` status.
3. **Operational Monitoring**: Instrument two critical indicators:
   - `outbox_unprocessed_count`: Total backlog of uncommitted messages.
   - `outbox_oldest_pending_age_seconds`: Elapsed time since the oldest pending message was written. A surge in oldest pending age indicates relay failure, poison pill records, or broker network partitions.
4. **Thin Events vs Fat Events**: Avoid storing full relational graphs in outbox payloads. Large payloads exhaust database buffer pools and saturate network I/O. Store identifiers and essential state mutations (thin events), allowing downstream consumers to query entity caches if full context is needed.
5. **Lab Simplification Boundaries**: The implementation in this lab utilizes an in-memory transactional storage engine (`internal/outbox/db.go`) to isolate architectural concepts without requiring external PostgreSQL or Kafka infrastructure. Production deployments require persistent disk storage, proper index design (`WHERE status = 'PENDING'`), and poison-pill dead-letter queues.

---

## Verification Summary

- **Repository**: [github.com/lukman-ss/software-engineering-lab/tree/main/labs/21-outbox-pattern](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/21-outbox-pattern)
- **Commit**: `ffb1d0c42b78d71a36d1e3b9769d01bd8f458f67`
- **Suite**: 8/8 automated test cases passing under Go race detector.
