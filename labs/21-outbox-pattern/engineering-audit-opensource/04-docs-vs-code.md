# Documentation vs Code Comparison

## Source Documents Reviewed
- README.md (labs/21-outbox-pattern/README.md)
- Implementation files (Go source)
- Test files (labs/21-outbox-pattern/tests/outbox_test.go)
- Demo (labs/21-outbox-pattern/cmd/demo/main.go)

## Findings

### DOC_CODE_MISMATCH: None

All core claims in README are matched by implementation and verified by tests/demo.

### Detailed Comparison

**README Claim**: "A demonstration of the Transactional Outbox pattern implemented in Go, proving atomic persistence between domain entities and event logs, decoupled polling relay dispatch to a message broker, and downstream consumer idempotency."

**Code Evidence**:
- Atomic persistence: `internal/outbox/db.go` Tx.Commit writes both order and outbox atomically under db.mu.Lock(); Tx.Rollback discards both; `service.go` CreateOrderWithOutbox uses Tx.
- Decoupled polling relay dispatch: `internal/outbox/relay.go` polls `db.GetPendingOutbox()` every `pollInterval`, publishes via broker, marks processed on success.
- Downstream consumer idempotency: `internal/outbox/consumer.go` uses `processedIDs` map to reject duplicate event IDs.

**README Architecture List**:
1. `internal/outbox/db.go` — exists ✓
2. `internal/outbox/broker.go` — exists ✓
3. `internal/outbox/service.go` — exists ✓
4. `internal/outbox/relay.go` — exists ✓
5. `internal/outbox/consumer.go` — exists ✓
6. `internal/outbox/model.go` — exists (defines Order, OutboxMessage, status constants) ✓ (not listed in README "consists of" but is an implementation detail; omission not a mismatch)

**README Running Tests**:
```bash
go test ./...
go test -race ./...
```
- Both commands run successfully (tests pass, race detector passes). ✓

**README Running Demo**:
```bash
go run ./cmd/demo
```
- Command runs successfully and prints expected output showing:
  - Scenario 1: Dual-write failure → order persisted, broker empty (inconsistency).
  - Scenario 2: Transactional outbox success → order+outbox saved, broker receives one message, consumer accepts it.
  - Scenario 3: Duplicate delivery → consumer rejects duplicate (idempotent).
- Output matches described behavior exactly. ✓

**Test Coverage Alignment**:
- Tests exercise:
  - Atomic write + relay dispatch + consumer accept (`TestTransactionalOutbox_HappyPath`).
  - Tx rollback discards staged data (`TestTransactionalOutbox_Rollback`).
  - Consumer idempotency (`TestTransactionalOutbox_Idempotency_DuplicateDelivery`).
  - Dual-write inconsistency (`TestDualWriteProblem_Failure`).
  - Processed message purge (`TestTransactionalOutbox_PurgeProcessed`).
  - Concurrency race-safety (`TestTransactionalOutbox_ConcurrentWrites` via `-race`).

**No DOC_CODE_MISMATCH found.**

## MISMATCH TYPES NOT OBSERVED
- TEST_CLAIM_MISMATCH: No test asserts false claims; test assertions match code behavior.
- RESEARCH_IMPLEMENTATION_MISMATCH: (Not audited per pipeline override.)

## Summary
Documentation accurately reflects implementation. No mismatches detected between README claims and actual code/tests/demo behavior.