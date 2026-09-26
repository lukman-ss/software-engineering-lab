# Test Audit

Target Lab: labs/21-outbox-pattern

## Coverage Overview

The test suite in `tests/outbox_test.go` exercises the following scenarios:
1. `TestTransactionalOutbox_HappyPath`: Verifies atomic write, relay dispatch, outbox status transition to `PROCESSED`, and consumer receipt.
2. `TestTransactionalOutbox_Rollback`: Verifies staged order and outbox records are aborted when `tx.Rollback()` is called. No broker publication occurs.
3. `TestTransactionalOutbox_Idempotency_DuplicateDelivery`: Verifies consumer deduplicates messages with the same ID, preventing double processing.
4. `TestDualWriteProblem_Failure`: Verifies state divergence when direct broker publication fails after commit.
5. `TestTransactionalOutbox_ConcurrentWrites`: Executes 10 concurrent goroutines writing 10 orders each while background relay is actively polling.

## Execution Results

### 1. `go test ./...`
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	0.231s
```
Status: PASS

### 2. `go test -race ./...`
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	1.397s
```
Status: PASS (No data races detected)

### 3. `go run ./cmd/demo`
```text
=== Lab 21: Transactional Outbox Pattern Demo ===

[Scenario 1: The Dual-Write Problem]
Direct write failed: failed to publish to broker after DB commit: broker unavailable
State Inconsistency: Order in DB = true, Broker Message Count = 0

[Scenario 2: Transactional Outbox Solution]
Creating order with transactional outbox...
Order and Outbox record atomically saved to DB.
Broker received messages: 1
 - Event ID: evt-order-outbox-success, Type: OrderCreated, Payload: {"ID":"order-outbox-success","CustomerID":"cust-2","Amount":300,"Status":"CREATED"}
 - Consumer processing initial message: accepted=true

[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]
Simulating duplicate delivery to consumer...
Consumer processing duplicate delivery: accepted=false (Duplicate safely skipped!)
Total events processed by consumer: 1

=== Demo Complete ===
```
Status: PASS (Real executable output matches expected behavior)

## Assessment
The tests prove all claimed functional, concurrency, and reliability behaviors.
