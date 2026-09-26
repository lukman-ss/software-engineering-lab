# Test Audit

## Test Suite Execution Results

### 1. `go test ./...`
Status: PASS
Execution Output:
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	0.226s
```

### 2. `go test -race ./...`
Status: PASS
Execution Output:
```text
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/cmd/demo	[no test files]
?   	github.com/software-engineering-lab/labs/21-outbox-pattern/internal/outbox	[no test files]
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	1.240s
```

### 3. `go run ./cmd/demo`
Status: PASS
Execution Output:
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

## Coverage Assessment

1. `TestTransactionalOutbox_HappyPath`: Covers end-to-end atomic commit, relay polling/dispatch, outbox status update (`PROCESSED`), and idempotent consumer intake.
2. `TestTransactionalOutbox_Rollback`: Verifies transaction rollback prevents both order and outbox persistence; proves broker receives no events.
3. `TestTransactionalOutbox_Idempotency_DuplicateDelivery`: Validates consumer deduplication under duplicate delivery.
4. `TestDualWriteProblem_Failure`: Verifies the dual-write flaw where DB commit succeeds but broker failure results in partial commit state.
5. `TestTransactionalOutbox_ConcurrentWrites`: Executes 10 goroutines performing 100 total concurrent writes during background relay polling under race detector.
6. `TestTransactionalOutbox_PurgeProcessed`: Validates outbox housekeeping and table compaction for processed messages.

Assessment: PASS. All test assertions are rigorous, verifiable, and executable.
