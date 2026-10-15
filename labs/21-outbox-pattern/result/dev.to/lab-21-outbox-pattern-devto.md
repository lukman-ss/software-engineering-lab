---
title: "Transactional Outbox Pattern: Atomic Database Writes Meet Eventual Message Delivery"
published: false
description: "Explore Transactional Outbox Pattern in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Transactional Outbox Pattern: Atomic Database Writes Meet Eventual Message Delivery

A production system must remain resilient under load. Lab 21 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Separate writes to independent systems cannot be atomic; decouple them with an outbox.

## Dual-Write Problem: Ketidakonsistenan Senyap

The dual-write problem occurs when an application attempts to write to a database and message broker separately within one business flow. Database transactions only control operations within the database itself; they cannot ROLLBACK external systems like Kafka or RabbitMQ. As a result, if database commit succeeds but broker publish fails, the system enters an inconsistent state: the order is persisted but the event is never sent. This condition is hard to detect in real-time and requires complex manual reconciliation mechanisms.

## Mengapa Pendekatan Naif Gagal?

The naive dual-write implementation seems intuitive: save the order to the database, then publish the event to the broker. However, because both operations are outside the same transactional boundary, a publish failure cannot be rolled back from the database side. The order remains persisted while the event is lost forever. This approach also introduces strong coupling between the ordering service and the broker, so any broker failure will abort the entire business flow.

```go
func CreateOrderDualWriteNaive(broker, orderID, amount) error {
    tx := db.BeginTx()
    tx.SaveOrder(order)
    tx.Commit()          // DB berhasil
    broker.Publish(msg)  // Broker gagal → inkonsisten
}
```

## Mental Model: Transactional Outbox

The core mental model of the Transactional Outbox pattern is to move the event write into the same transactional boundary as the business entity write. Instead of performing two separate writes, the application writes the order and outbox message simultaneously within one database transaction. After successful commit, both records are persisted. An asynchronous process called the message relay then periodically polls PENDING records and delivers them to the broker. If the relay fails, the record remains PENDING and will be retried in the next cycle.

```go
type OutboxMessage struct {
    ID        string
    EventType string
    Payload   string
    Status    MessageStatus  // PENDING, PROCESSED
    CreatedAt time.Time
}
```

## Alur Eksekusi: Dari Transaksi Hingga Consumer

The Transactional Outbox execution pipeline consists of five main stages. First, the service receives an order creation request and starts a database transaction. Second, the service stages both the order and outbox message within the transaction. Third, the service commits, which atomically flushes both staged records to the database. Fourth, a relay worker running periodically polls PENDING records and attempts to publish them to the broker. Fifth, the consumer receives the message and performs an idempotent check before further processing.

## Implementasi Atomic Outbox Write

The atomic outbox write implementation resides in the CreateOrderWithOutbox method in service.go. This method starts a new transaction, creates the order entity and outbox message, then stages both using SaveOrder and SaveOutbox. If JSON marshalling fails or any save operation errors, the method rolls back and returns an error. Only after both staged records are successfully saved does the method commit, which atomically flushes all changes to the database.

```go
func (s *OrderService) CreateOrderWithOutbox(orderID string, customerID string, amount float64) error {
    tx := s.db.BeginTx()
    order := Order{ID: orderID, CustomerID: customerID, Amount: amount, Status: OrderStatusCreated}
    payload, _ := json.Marshal(order)
    msg := OutboxMessage{ID: fmt.Sprintf("evt-%s", orderID), EventType: "OrderCreated", Payload: string(payload), Status: MessageStatusPending, CreatedAt: time.Now()}
    tx.SaveOrder(order)
    tx.SaveOutbox(msg)
    return tx.Commit()
}
```

## Failure Path: Dual-Write Naive Menghasilkan Dangling Order

This failure path proves how naive dual-write creates a dangling reference: an order persisted in the database without a corresponding event in the broker. In the demo, the broker is simulated to fail via SetFailNext(true). After successful database commit, the CreateOrderDualWriteNaive method attempts to publish to the failing broker, causing an error to be returned. However, the order is already permanently stored in the database. The TestDualWriteProblem_Failure test case asserts this condition: the order is found in the DB, but the broker received no messages.

## Mekanisme vs Gejala: Dual-Write vs Outbox

This comparison table highlights the fundamental differences between naive dual-write and transactional outbox approaches. In naive dual-write, the write boundary is separated between database commit and broker publish, so there is no cross-system atomicity. When the broker fails after commit, the order persists while the event is lost. Conversely, in transactional outbox, both writes are performed within one database transaction, guaranteeing atomicity for the first stage. Broker failure only delays delivery, not deletes the event, because the relay will retry until success.

Aspek | Dual-Write Naive | Transactional Outbox
|---|---|---|
Write Boundary | Terpisah (DB commit dulu, lalu broker) | Sama (order + outbox dalam satu Tx)
Atomicity | Tidak atomik lintas sistem | Atomik untuk write pertama
Broker Failure Impact | Order persisten, event hilang | Order + event persisten, relay retry
Idempotency Requirement | Tidak diperlukan (karena konsisten) | Diperlukan (at-least-once delivery)
Operational Complexity | Rendah, tetapi risiko tinggi | Sedang, dengan relay dan cleanup

## 3 Keputusan Desain Kritis

These three critical design decisions form the foundation of a robust transactional outbox implementation. The first decision establishes an atomic transaction boundary for business and outbox writes. The second decision selects a reliable relay mechanism with retry logic. The third decision ensures the consumer can safely handle duplicate delivery through idempotency checks. These three decisions are interdependent and must be implemented together to achieve the desired eventual consistency.

## Bukti Eksekusi: Terminal Output Demo

The terminal output above is direct execution proof from the demo command running three scenarios. Scenario 1 demonstrates dual-write failure: the order exists in the DB but the broker message count equals zero. Scenario 2 demonstrates outbox success: one message was successfully published to the broker after relay polling. Scenario 3 demonstrates the idempotent consumer: duplicate delivery is rejected with accepted=false, ensuring the received count remains one. This output verifies that the outbox pattern successfully overcomes the dual-write problem.

## YANG TIDAK DIBUKTIKAN

This laboratory has several limitations that need to be noted. First, the database used is an in-memory mock, not a production RDBMS like PostgreSQL or MySQL. Second, only the Polling Publisher strategy is implemented; the Transaction Log Tailing approach using CDC tools like Debezium is not shown. Third, there is no dead-letter queue to handle payloads that fail repeatedly. Fourth, the OutboxMessage schema is simplified without the aggregate columns typically present in production designs. Fifth, the SLA values mentioned are illustrative and must be adjusted to each system's SLOs.

## Trade-off Arsitektur: Outbox Pattern

Implementing the Transactional Outbox pattern brings trade-offs that need to be considered. On the benefit side, the pattern eliminates silent event loss that often occurs in naive dual-write and guarantees at-least-once delivery through retry mechanisms. It also decouples the ordering service from specific broker implementations. However, there are significant drawbacks: additional latency because events do not reach the broker immediately but must wait for polling intervals. Additionally, there is operational overhead in monitoring oldest-event age and pending counts, as well as added complexity for relay processes, cleanup jobs, and idempotency tracking.

## Concurrency Safety: 10 Workers × 10 Orders = 100 Messages

Concurrency testing proves that the outbox pattern implementation is safe against multi-goroutine access. In TestTransactionalOutbox_ConcurrentWrites, 10 worker goroutines each generate 10 orders simultaneously, totaling 100 orders. Go's race detector reports no data races. After the relay finishes processing all pending messages, the broker receives exactly 100 messages and none remain in pending state. Another test, TestTransactionalOutbox_ConcurrentConsumers, verifies that 10 workers sending messages with identical IDs result in only 20 unique processed messages, proving idempotency works correctly under concurrency.

```go
func TestTransactionalOutbox_ConcurrentWrites(t *testing.T) {
    workers := 10
    ordersPerWorker := 10
    totalExpected := workers * ordersPerWorker  // 100
    // Goroutines menulis order + outbox secara konkuren
    // Relay polling dan publish
    // Assert: len(broker.GetPublished()) == 100
    // Assert: len(db.GetPendingOutbox()) == 0
}
```

## Before vs After: Transformasi Arsitektur

The transformation from naive dual-write to transactional outbox fundamentally changes the architecture. Before implementing this pattern, the service writes to the database and broker sequentially without atomicity guarantees. If the broker fails, the order persists but the event is lost, creating an inconsistent state that is hard to detect. After implementing the pattern, both writes are performed within one database transaction, so both order and event persist together. Broker failure only delays delivery, not deletes the event, because the relay will keep retrying until success.

## Test 1: Happy Path - Atomic Write + Relay Dispatch

This test proves the complete happy path of the transactional outbox. First, the service creates an order with ID o-1 atomically along with the outbox message. After transaction commit, the relay running at 10ms intervals will poll and find the PENDING message. The relay then publishes to the broker and marks the message as PROCESSED. The consumer receives the message and processes it. This test asserts that the broker receives exactly one message with EventType OrderCreated, and the outbox message status changes to PROCESSED. All assertions pass, proving the end-to-end flow works correctly.

```go
func TestTransactionalOutbox_HappyPath(t *testing.T) {
    db := outbox.NewDB()
    broker := outbox.NewMockBroker()
    relay := outbox.NewRelay(db, broker, 10*time.Millisecond)
    service := outbox.NewOrderService(db)
    consumer := outbox.NewConsumer()
    relay.Start()
    err := service.CreateOrderWithOutbox("o-1", "c-1", 100.50)
    assert.NoError(t, err)
    time.Sleep(50 * time.Millisecond)
    published := broker.GetPublished()
    assert.Len(t, published, 1)
    assert.Equal(t, "OrderCreated", published[0].EventType)
    msg, _ := db.GetOutbox("evt-o-1")
    assert.Equal(t, outbox.MessageStatusProcessed, msg.Status)
}
```

## Test 2: Rollback - Nothing Persists

The rollback test proves that staged mutations are completely discarded when a transaction is rolled back. The SaveOrder and SaveOutbox methods only place data in the stagedOrders and stagedOutbox maps, without writing to the main database. When Rollback is called, both staged maps are discarded without ever touching the database. The test asserts that the order and outbox message cannot be retrieved from the database after rollback. Additionally, after waiting for relay polling, the broker receives no messages, proving that rollback prevents event dispatch entirely.

```go
func TestTransactionalOutbox_Rollback(t *testing.T) {
    db := outbox.NewDB()
    tx := db.BeginTx()
    tx.SaveOrder(outbox.Order{ID: "o-2", Amount: 50})
    tx.SaveOutbox(outbox.OutboxMessage{ID: "evt-o-2", Status: outbox.MessageStatusPending})
    tx.Rollback()
    _, ok := db.GetOrder("o-2")
    assert.False(t, ok, "order should not exist")
    _, ok = db.GetOutbox("evt-o-2")
    assert.False(t, ok, "outbox should not exist")
    time.Sleep(30 * time.Millisecond)
    assert.Len(t, broker.GetPublished(), 0, "no messages to broker")
}
```

## Test 3: Idempotency - Duplicate Delivery Rejected

The idempotency test proves that the consumer can safely handle duplicate delivery. First, Handle is called with a unique message and returns true (accepted). Second, Handle is called again with the same message and returns false (rejected). The test asserts that the received count remains 1, proving that duplicates do not add to the consumer's state. This scenario simulates a case where the relay crashes after publish but before marking the message as PROCESSED, causing the same message to be republished by the relay in the next cycle.

```go
func TestTransactionalOutbox_Idempotency_DuplicateDelivery(t *testing.T) {
    consumer := outbox.NewConsumer()
    msg := outbox.OutboxMessage{ID: "evt-duplicate", EventType: "OrderCreated", Payload: "data"}
    p1 := consumer.Handle(msg)
    assert.True(t, p1, "first delivery should be accepted")
    p2 := consumer.Handle(msg)
    assert.False(t, p2, "duplicate delivery should be rejected")
    assert.Equal(t, 1, consumer.GetReceivedCount(), "count stays 1")
}
```

## Test 4: Relay Retry - Broker Failure Tolerated

The relay retry test proves the failure tolerance mechanism against broker failures. First, the broker is set to fail on the next publish via SetFailNext(true). The service creates an order with outbox, then the relay starts. On the first poll, the publish fails because the broker is in fail mode, so the message remains PENDING. On the next poll (after the ticker interval), the broker returns to normal, the publish succeeds, and the message is marked PROCESSED. This test asserts that after retry, the broker receives the message and the outbox status changes to PROCESSED.

```go
func TestTransactionalOutbox_RelayRetryAfterBrokerFailure(t *testing.T) {
    db := outbox.NewDB()
    broker := outbox.NewMockBroker()
    relay := outbox.NewRelay(db, broker, 10*time.Millisecond)
    service := outbox.NewOrderService(db)
    broker.SetFailNext(true)  // first publish fails
    service.CreateOrderWithOutbox("o-retry", "c-retry", 55.0)
    relay.Start()
    time.Sleep(25 * time.Millisecond)  // first poll fails
    msg, _ := db.GetOutbox("evt-o-retry")
    assert.Equal(t, outbox.MessageStatusPending, msg.Status)
    time.Sleep(50 * time.Millisecond)  // second poll succeeds
    published := broker.GetPublished()
    assert.Len(t, published, 1)
    msg, _ = db.GetOutbox("evt-o-retry")
    assert.Equal(t, outbox.MessageStatusProcessed, msg.Status)
}
```

## Test 5: Purge Processed - Memory Leak Prevention

The purge processed test proves the outbox table cleanup mechanism to prevent memory leaks. The database initially contains two messages: m1 with PENDING status and m2 with PROCESSED status. The PurgeProcessedOutbox method will delete all records with PROCESSED status and return the count of deleted records. The test asserts that purged count equals 1, the PENDING message still exists, and the PROCESSED message has been deleted. Without such a cleanup mechanism, the outbox table would grow indefinitely because PROCESSED records are never deleted.

```go
func TestTransactionalOutbox_PurgeProcessed(t *testing.T) {
    db := outbox.NewDB()
    tx := db.BeginTx()
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

## Checklist Implementasi Outbox Pattern

This implementation checklist covers essential steps for applying the transactional outbox pattern in production. First, create an outbox_events table with an adequate schema including aggregate columns for Kafka partitioning. Second, ensure every business transaction includes an INSERT to outbox_events. Third, implement a polling relay with an interval balanced between latency and database load. Fourth, the relay only marks messages as PROCESSED after the broker provides ACK. Fifth, the consumer must have an idempotency check based on event ID. Sixth, add a scheduled job to purge processed records to prevent table bloat.

## Anti-Pattern: Apa yang Harus Dihindari

Anti-patterns in outbox pattern implementation can undermine all the consistency guarantees intended. Writing to the broker before database commit reverts to naive dual-write. Using non-unique event IDs causes deduplication to fail. Consumers without idempotency checks will process duplicate messages twice. No purge for processed records will cause the outbox table to grow indefinitely. A relay with too-long polling intervals will significantly increase delivery latency.

## Takeaway Praktis & Command Block

The practical takeaway from this lab is understanding how to reproduce and verify the outbox pattern implementation. The go test -v ./... command runs the entire test suite with verbose output. The go test -race ./... command runs the race detector to ensure thread safety. The go run ./cmd/demo command runs the interactive demo displaying three scenarios. Best practices include using stable event IDs based on aggregate IDs, monitoring oldest-event age as a degradation signal, maintaining pending counts below thresholds, and implementing dead-letter queues for payloads that fail repeatedly.

```sh
cd labs/21-outbox-pattern
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## Prinsip Arsitektural Jangka Panjang

These long-term architectural principles extend beyond the specific transactional outbox implementation. First, eventual consistency is often more practical than strong consistency that requires expensive distributed coordination. Second, separating the write boundary from the delivery boundary allows each part to be optimized independently. Third, at-least-once delivery is safer than at-most-once because it prevents data loss. Fourth, idempotency is not an optional feature but a necessity because retries are normal. Fifth, observable metrics like oldest-event age and pending counts must be implemented for early degradation detection.

## Source

[Software Engineering Lab 21 - Transactional Outbox Pattern](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/21-outbox-pattern), snapshot commit `b2d32740`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
