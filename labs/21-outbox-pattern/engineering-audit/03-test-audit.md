# Test Audit

## Overview
Test suite located in `tests/outbox_test.go` contains 5 tests covering happy path, transaction rollback, consumer idempotency on duplicate deliveries, dual-write vulnerability, and concurrent writes under the race detector.

## Test Inventory & Assessment

### 1. `TestTransactionalOutbox_HappyPath`
- **Focus**: End-to-end transactional write, polling dispatch, status update, downstream consumer handling.
- **Coverage**:
  - Atomic DB write of `Order` and `OutboxMessage`.
  - Relay polling picks up message and publishes to `MockBroker`.
  - DB outbox record marked `PROCESSED`.
  - Consumer receives and processes message.
- **Assessment**: PASS

### 2. `TestTransactionalOutbox_Rollback`
- **Focus**: Transaction abort/rollback.
- **Coverage**:
  - Simulates failed transaction via `tx.Rollback()`.
  - Asserts neither `Order` nor `OutboxMessage` persisted in `db`.
  - Asserts relay sends zero messages to broker.
- **Assessment**: PASS

### 3. `TestTransactionalOutbox_Idempotency_DuplicateDelivery`
- **Focus**: Consumer-side idempotency against at-least-once delivery retries.
- **Coverage**:
  - Initial message processed successfully (`accepted=true`).
  - Identical duplicate message rejected (`accepted=false`).
  - Total processed count remains 1.
- **Assessment**: PASS

### 4. `TestDualWriteProblem_Failure`
- **Focus**: Demonstrating dual-write state inconsistency.
- **Coverage**:
  - Simulates broker down via `broker.SetFailNext(true)`.
  - DB commit succeeds, broker publish fails.
  - Asserts DB has persisted order while broker received 0 messages.
- **Assessment**: PASS

### 5. `TestTransactionalOutbox_ConcurrentWrites`
- **Focus**: Concurrency safety under high concurrent load.
- **Coverage**:
  - 10 concurrent goroutines writing 10 transactions each (100 operations total) alongside background polling relay worker.
  - Verified with `-race` flag. Zero data races detected.
- **Assessment**: PASS

## Execution Results

```text
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
PASS
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	1.439s
```

Race detector output:
- `go test -count=1 -race ./...` PASSED with 0 warnings or race warnings.
