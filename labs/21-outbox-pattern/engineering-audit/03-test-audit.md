# Test Audit

## Executed Commands & Actual Outputs

### Command 1: Unit & Integration Tests
```bash
go test -v ./...
```
Output:
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
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
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	0.324s
```

### Command 2: Race Detector
```bash
go test -count=1 -race ./...
```
Output:
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	1.273s
```

### Command 3: Executable Demo
```bash
go run ./cmd/demo
```
Output:
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

## Coverage Verification

1. **Happy Path**: `TestTransactionalOutbox_HappyPath` verifies atomic save, relay dispatch, outbox status transition to `PROCESSED`, and consumer processing.
2. **Failure Path**: `TestDualWriteProblem_Failure` reproduces partial write failure where DB write commits but broker publish fails.
3. **Rollback**: `TestTransactionalOutbox_Rollback` asserts staged entity and outbox message are completely discarded upon `tx.Rollback()`.
4. **Edge Cases & Deduplication**: `TestTransactionalOutbox_Idempotency_DuplicateDelivery` verifies consumer detects duplicate event IDs and skips double-execution.
5. **Concurrency**: `TestTransactionalOutbox_ConcurrentWrites` runs 10 goroutines issuing 100 concurrent writes alongside active relay polling with zero data races detected under `-race`.
