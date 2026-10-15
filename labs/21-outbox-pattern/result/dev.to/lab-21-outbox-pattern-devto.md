---
title: "Transactional Outbox Pattern: Atomic Database Writes Meet Eventual Message Delivery"
published: false
description: "Explore the Transactional Outbox Pattern: solves dual-write problem between database and message broker through atomic writes, asynchronous relay polling, and idempotent consumers."
tags: go, architecture, messaging, testing, softwareengineering
---

# Transactional Outbox Pattern: Atomic Database Writes Meet Eventual Message Delivery

A microservice writes an order to the database, then publishes `OrderCreated` to Kafka. If Kafka fails after DB commit, the order exists but the event never arrives. Lab 21 proves this with a Go simulation.

> The pattern decouples atomic writes from asynchronous delivery through a polling relay.

## Dual-Write Problem: Silent Event Loss

Database commit succeeds, storing the order. Broker publish fails due to network timeout. The order persists, but `OrderCreated` never reaches downstream systems. No retry exists for naive dual-write—this is silent event loss.

```go
tx := db.BeginTx()
tx.SaveOrder(order)
tx.Commit()          // DB succeeded
broker.Publish(msg)  // Broker failed → inconsistent
```

## Why Naive Dual-Write Fails

Naive code looks simple but offers no cross-system atomicity. The broker call happens after commit—outside any transaction. If it fails, the order is permanently saved with no message to replay. This couples business logic to broker availability.

```go
func CreateOrderDualWriteNaive(broker, orderID, amount) error {
    tx := db.BeginTx()
    tx.SaveOrder(order)
    tx.Commit()          // DB succeeded
    broker.Publish(msg)  // Broker failed → inconsistent
}
```

## Mental Model: Transactional Outbox

Write the business entity AND the event to an `outbox_events` table within the same database transaction. After commit, a background relay polls for `PENDING` records and publishes them. If the relay fails, the message stays `PENDING` for retry.

```go
type OutboxMessage struct {
    ID        string
    EventType string
    Payload   string
    Status    MessageStatus // PENDING, PROCESSED
    CreatedAt time.Time
}
```

## Execution Pipeline: From Transaction to Consumer

Client → OrderService.CreateOrderWithOutbox → DB.BeginTx() → SaveOrder + SaveOutbox → DB.Commit() → Relay polls PENDING → Broker.Publish() → Mark PROCESSED → Consumer.Handle() with idempotency.

## Atomic Outbox Write Implementation

```go
func (s *OrderService) CreateOrderWithOutbox(orderID, customerID string, amount float64) error {
    tx := s.db.BeginTx()
    order := Order{ID: orderID, CustomerID: customerID, Amount: amount, Status: OrderStatusCreated}
    payload, _ := json.Marshal(order)
    msg := OutboxMessage{ID: fmt.Sprintf("evt-%s", orderID), EventType: "OrderCreated", Payload: string(payload), Status: MessageStatusPending, CreatedAt: time.Now()}
    tx.SaveOrder(order)
    tx.SaveOutbox(msg)
    return tx.Commit()
}
```

The transaction stages both writes—commit only succeeds if both persist together.

## Failure Path: Dangling Orders from Naive Write

When broker publish follows DB commit, failure leaves orders without events. Tests prove: order exists in DB, but broker message count is zero. No trace remains for retry.

```
[Scenario 1: Dual-Write Problem]
Direct write failed: failed to publish to broker after DB commit: broker unavailable
State Inconsistency: Order in DB = true, Broker Message Count = 0
```

## Mechanism Comparison: Dual-Write vs Outbox

| Aspect | Dual-Write Naive | Transactional Outbox |
|---|---|---|
| Write Boundary | Separate (DB commit, then broker) | Same (order + outbox in 1 Tx) |
| Atomicity | Not cross-system atomic | Atomic for first write |
| Broker Failure Impact | Order persisted, event lost | Order + event persisted, relay retries |
| Idempotency Requirement | Not needed (consistent) | Required (at-least-once delivery) |
| Operational Complexity | Low, but high risk | Medium, with relay and cleanup |

## Three Critical Design Decisions

1. **Atomic Boundary**: Include business entity + outbox event write in one transaction.
2. **Reliable Relay**: Choose polling interval (1-5 seconds) or CDC log-tailing. Only mark `PROCESSED` after broker ACK.
3. **Idempotent Consumer**: Track processed event IDs; second delivery is ignored.

## Terminal Output: Demo Execution

```
=== Lab 21: Transactional Outbox Pattern Demo ===

[Scenario 1: The Dual-Write Problem]
Direct write failed: failed to publish to broker after DB commit: broker unavailable
State Inconsistency: Order in DB = true, Broker Message Count = 0

[Scenario 2: Transactional Outbox Solution]
Broker received messages: 1
 - Event ID: evt-order-outbox-success, Type: OrderCreated, Payload: {...}

[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]
Consumer processing duplicate delivery: accepted=false
```

## Evidence Limitations: What This Lab Does Not Prove

This lab uses in-memory mock DB, not PostgreSQL/MySQL. Only polling publisher is implemented; CDC/log-tailing is not shown. No dead-letter queue for poison payloads. Example SLAs (oldest-event ~2s normal, 47+ min abnormal) are illustrative only.

## Architectural Trade-offs

**Benefits**: Eliminates silent event loss, guarantees at-least-once delivery via retry, decouples service from broker implementation, idempotent consumer prevents double processing.

**Drawbacks**: Added latency (wait for polling interval), operational overhead (monitor oldest-event age, pending count), added complexity (relay process, cleanup, idempotency tracking), growing outbox table until purge.

## Concurrency Safety: 10 Workers × 10 Orders

```go
func TestTransactionalOutbox_ConcurrentWrites(t *testing.T) {
    workers := 10
    ordersPerWorker := 10
    totalExpected := workers * ordersPerWorker // 100
    // Goroutines write order + outbox concurrently
    // Relay polling and publish
    // Assert: len(broker.GetPublished()) == 100
    // Assert: len(db.GetPendingOutbox()) == 0
}
```

Result: Race detector PASS, total published: 100, pending after: 0, consumer rejects duplicates correctly.

## Before vs After: Architectural Transformation

**Before (Dual-Write Naive)**: Service → DB.Commit → Broker.Publish. Order persists, event may lose. Broker down = inconsistent state. Testing hard to assert cross-system consistency.

**After (Transactional Outbox)**: Service → DB.Tx(Order+Outbox) → Relay → Broker. Order + event persist atomically. Broker down = retry, nothing lost. Testing: unit + integration tests verify consistency.

## Test 1: Happy Path Verification

```go
func TestTransactionalOutbox_HappyPath(t *testing.T) {
    db := outbox.NewDB()
    broker := outbox.NewMockBroker()
    relay := outbox.NewRelay(db, broker, 10*time.Millisecond)
    service := outbox.NewOrderService(db)
    relay.Start()
    err := service.CreateOrderWithOutbox("o-1", "c-1", 100.50)
    assert.NoError(t, err)
    time.Sleep(50 * time.Millisecond)
    published := broker.GetPublished()
    assert.Len(t, published, 1)
    assert.Equal(t, "OrderCreated", published[0].EventType)
}
```

## Test 2: Rollback Verification

When transaction rolls back, neither order nor outbox message persists. No events reach the broker.

```go
func TestTransactionalOutbox_Rollback(t *testing.T) {
    tx := db.BeginTx()
    tx.SaveOrder(outbox.Order{ID: "o-2", Amount: 50})
    tx.SaveOutbox(outbox.OutboxMessage{ID: "evt-o-2", Status: outbox.MessageStatusPending})
    tx.Rollback()
    _, ok := db.GetOrder("o-2")
    assert.False(t, ok, "order should not exist")
    _, ok = db.GetOutbox("evt-o-2")
    assert.False(t, ok, "outbox should not exist")
}
```

## Test 3: Idempotency Verification

Consumer rejects duplicate deliveries; received count stays 1.

```go
func TestTransactionalOutbox_Idempotency_DuplicateDelivery(t *testing.T) {
    consumer := outbox.NewConsumer()
    msg := outbox.OutboxMessage{ID: "evt-duplicate", EventType: "OrderCreated", Payload: "data"}
    p1 := consumer.Handle(msg)
    assert.True(t, p1, "first delivery should be accepted")
    p2 := consumer.Handle(msg)
    assert.False(t, p2, "duplicate delivery should be rejected")
    assert.Equal(t, 1, consumer.GetReceivedCount())
}
```

## Test 4: Relay Retry Verification

When broker fails once, relay retries and succeeds on next poll. Status transitions from PENDING to PROCESSED.

```go
func TestTransactionalOutbox_RelayRetryAfterBrokerFailure(t *testing.T) {
    broker.SetFailNext(true)  // first publish fails
    service.CreateOrderWithOutbox("o-retry", "c-retry", 55.0)
    relay.Start()
    time.Sleep(25 * time.Millisecond)  // first poll fails
    msg, _ := db.GetOutbox("evt-o-retry")
    assert.Equal(t, outbox.MessageStatusPending, msg.Status)
    time.Sleep(50 * time.Millisecond)  // second poll succeeds
    published := broker.GetPublished()
    assert.Len(t, published, 1)
}
```

## Test 5: Processed Purge Verification

Cleanup mechanism prevents unbounded table growth.

```go
func TestTransactionalOutbox_PurgeProcessed(t *testing.T) {
    tx.SaveOutbox(outbox.OutboxMessage{ID: "m1", Status: outbox.MessageStatusPending})
    tx.SaveOutbox(outbox.OutboxMessage{ID: "m2", Status: outbox.MessageStatusProcessed})
    tx.Commit()
    purged := db.PurgeProcessedOutbox()
    assert.Equal(t, 1, purged)
    _, ok := db.GetOutbox("m1")
    assert.True(t, ok, "pending should remain")
    _, ok = db.GetOutbox("m2")
    assert.False(t, ok, "processed should be purged")
}
```

## Implementation Checklist

- [ ] Create `outbox_events` table: id, aggregate_type, aggregate_id, type, payload, status, created_at
- [ ] Business transaction must include INSERT to outbox_events
- [ ] Implement polling relay with appropriate interval (1-5 seconds)
- [ ] Relay marks PROCESSED only after broker ACK
- [ ] Consumer implements idempotent check by event ID
- [ ] Schedule job to purge processed records
- [ ] Monitor oldest-event age and pending count for degradation detection
- [ ] Add dead-letter queue for non-retryable poison payloads

## Anti-Patterns to Avoid

1. Writing to broker before DB commit reverts to naive dual-write.
2. Using non-unique event IDs causes deduplication failures.
3. Consumer without idempotency processes duplicates twice.
4. No purge for processed records causes table bloat.
5. Excessive polling interval (>1 minute) increases latency.

## Practical Takeaway & Commands

```sh
cd labs/21-outbox-pattern
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Best practices: stable event IDs, monitor oldest-event age, keep pending count low, implement dead-letter queue for repeated failures.

## Long-Term Architectural Principles

1. Eventual consistency often beats expensive strong consistency.
2. Decouple write boundary from delivery boundary.
3. At-least-once delivery is safer than at-most-once.
4. Idempotency is required, not optional.
5. Observable metrics (oldest-event age, pending count) must be implemented.

## References

- [Software Engineering Lab 21 - Transactional Outbox](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/21-outbox-pattern)
- All source files: `internal/outbox/`, `tests/outbox_test.go`, `cmd/demo/main.go`
- Commit: `b2d32740bf482b67be48a216c69e5f8e6b83fe8b`

The outbox pattern eliminates dual-write risk through atomic database writes and eventual message delivery. All tests pass with race detector clean.