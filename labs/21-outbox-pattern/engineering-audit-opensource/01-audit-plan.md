# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern (Transactional Outbox Pattern, Go)

## Implementation Files
- labs/21-outbox-pattern/go.mod
- labs/21-outbox-pattern/README.md
- labs/21-outbox-pattern/cmd/demo/main.go
- labs/21-outbox-pattern/internal/outbox/model.go
- labs/21-outbox-pattern/internal/outbox/db.go
- labs/21-outbox-pattern/internal/outbox/broker.go
- labs/21-outbox-pattern/internal/outbox/service.go
- labs/21-outbox-pattern/internal/outbox/relay.go
- labs/21-outbox-pattern/internal/outbox/consumer.go

## Tests
- labs/21-outbox-pattern/tests/outbox_test.go (package `tests`)
  - TestTransactionalOutbox_HappyPath
  - TestTransactionalOutbox_Rollback
  - TestTransactionalOutbox_Idempotency_DuplicateDelivery
  - TestDualWriteProblem_Failure
  - TestTransactionalOutbox_ConcurrentWrites
  - TestTransactionalOutbox_PurgeProcessed

## Executable/Demo
- labs/21-outbox-pattern/cmd/demo/main.go (`go run ./cmd/demo`)

## Approved Research Inputs
- (Audit stage: implementation and tests only. Research/content not audited per pipeline override.)
- README claims: atomic persistence of domain entity + event log; decoupled polling relay dispatch; downstream consumer idempotency.

## Main Claims To Verify
1. Order + outbox message are persisted atomically within a single transaction (both committed or neither).
2. The relay polls the outbox table and dispatches events to the broker (decoupled dispatch).
3. Broker publish failure leaves the outbox message PENDING for retry (durable at-least-once delivery).
4. The consumer is idempotent on duplicate delivery (event-ID based dedup).
5. Rollback discards staged order and outbox writes.
6. Processed outbox records can be purged.
7. Concurrent writes to the in-memory DB are free of data races.
8. The demo runs and reproduces the claimed scenarios truthfully.

## Commands To Run
```bash
cd labs/21-outbox-pattern
go build ./...
go vet ./...
go test ./...
go test -race ./...
go run ./cmd/demo
```

## Primary Risks
- The relay's retry-on-broker-failure path is central to the outbox durability guarantee but may be untested.
- The concurrent-writes test reuses the same orderID across goroutines and may not assert correctness (only race-freedom).
- Service-level rollback path (e.g. marshal error) may be unreachable and therefore untested.
- `Relay.Stop()` closing `stopChan` twice would panic (latent, not triggered by current tests/demo).
