# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- `internal/outbox/model.go`
- `internal/outbox/db.go`
- `internal/outbox/broker.go`
- `internal/outbox/service.go`
- `internal/outbox/relay.go`
- `internal/outbox/consumer.go`
Tests:
- `tests/outbox_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/runs/2026-09-25-the-outbox-pattern/05-report.md`
- `research/2026-09-26-outbox-pattern/05-report.md`
- `engineering/01-design.md`
Main Claims To Verify:
1. Dual-write vulnerability without transactional outbox causes state inconsistency between DB and message broker when broker publish fails.
2. Transactional outbox guarantees atomic persistence of domain entity (`Order`) and event log (`OutboxMessage`) within a single database transaction.
3. Transaction rollback guarantees neither entity nor outbox record is persisted.
4. Background polling relay fetches pending outbox events, dispatches to message broker, and marks records as `PROCESSED`.
5. Downstream consumer achieves idempotency via event deduplication tracking despite at-least-once delivery.
6. DB cleanup / purge mechanism successfully deletes processed outbox entries.
7. Concurrency safety under race detector (`-race`).
Commands To Run:
```bash
cd labs/21-outbox-pattern
go build ./...
go test -v -count=1 ./...
go test -v -count=1 -race ./...
go run ./cmd/demo
```
Primary Risks:
- Thread-safety across mock DB transactions, relay polling, and consumer deduplication sets.
- Mismatch between design document mentioning SQLite and actual implementation using in-memory transactional mock engine.
- Weak or unverified test assertions for relay dispatch timing.
