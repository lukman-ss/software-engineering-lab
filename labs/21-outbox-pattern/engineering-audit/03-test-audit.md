# Test Audit

Target Lab: labs/21-outbox-pattern

## Test Suite Overview

File: `tests/outbox_test.go`
Number of Tests: 7 tests

### 1. `TestTransactionalOutbox_HappyPath`
- Scope: End-to-end atomic order + outbox creation, relay polling, broker receipt, DB status transition (`PENDING` -> `PROCESSED`), and consumer delivery.
- Assessment: PASS. Verifies full lifecycle.

### 2. `TestTransactionalOutbox_Rollback`
- Scope: Negative / rollback case. Ensures uncommitted transactions discard both domain and outbox entities, preventing broker dispatch.
- Assessment: PASS. Verifies rollback atomicity.

### 3. `TestTransactionalOutbox_Idempotency_DuplicateDelivery`
- Scope: At-least-once delivery duplicate simulation.
- Assessment: PASS. Proves consumer deduplication on repeated event ID.

### 4. `TestDualWriteProblem_Failure`
- Scope: Flaw verification. Broker failure following DB commit.
- Assessment: PASS. Asserts DB contains record while broker has 0 messages.

### 5. `TestTransactionalOutbox_ConcurrentWrites`
- Scope: Concurrency stress testing with 10 workers submitting 10 orders each concurrently while relay worker polls.
- Assessment: PASS. Run with `go test -race` passes with zero race conditions.

### 6. `TestTransactionalOutbox_PurgeProcessed`
- Scope: Outbox table cleanup / retention purging of processed records.
- Assessment: PASS. Asserts processed records removed while pending records retained.

### 7. `TestTransactionalOutbox_RelayRetryAfterBrokerFailure`
- Scope: Fault recovery and relay start/stop idempotency. Injects broker failure on initial poll, verifies subsequent poll succeeds and marks record processed.
- Assessment: PASS. Proves recovery after transient failure.

## Test Execution Results

```text
go test -v ./...
=== RUN   TestTransactionalOutbox_HappyPath
--- PASS: TestTransactionalOutbox_HappyPath (0.05s)
=== RUN   TestTransactionalOutbox_Rollback
--- PASS: TestTransactionalOutbox_Rollback (0.03s)
=== RUN   TestTransactionalOutbox_Idempotency_DuplicateDelivery
--- PASS: TestTransactionalOutbox_Idempotency_DuplicateDelivery (0.00s)
=== RUN   TestDualWriteProblem_Failure
--- PASS: TestDualWriteProblem_Failure (0.00s)
=== RUN   TestTransactionalOutbox_ConcurrentWrites
--- PASS: TestTransactionalOutbox_ConcurrentWrites (0.05s)
=== RUN   TestTransactionalOutbox_PurgeProcessed
--- PASS: TestTransactionalOutbox_PurgeProcessed (0.00s)
=== RUN   TestTransactionalOutbox_RelayRetryAfterBrokerFailure
--- PASS: TestTransactionalOutbox_RelayRetryAfterBrokerFailure (0.08s)
PASS
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	0.224s
```

Race detector:
```text
go test -race ./...
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	0.750s
```

Demo execution:
```text
go run ./cmd/demo
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

## Summary
All required test dimensions (happy path, failure path, transitions, rollback, retry, concurrency, cleanup) are fully tested and proven by real executions.
