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
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
- `research/05-report.md`

Main Claims To Verify:
1. Naive dual-write causes inconsistency when message broker publish fails after database commit.
2. Transactional outbox pattern atomically writes business entity (`Order`) and event message (`OutboxMessage`) within a single database transaction.
3. Transaction rollback discards both business entity and outbox message without publishing.
4. Polling relay worker asynchronously queries `PENDING` outbox records, publishes to broker, and marks records as `PROCESSED`.
5. Downstream consumer handles duplicate deliveries idempotently via event identifier tracking.
6. Processed outbox records can be safely purged without affecting pending messages.
7. Concurrency and race safety holds across multiple writers and relay polling (`go test -race ./...`).

Commands To Run:
- `go test -v ./...`
- `go test -race -v ./...`
- `go run ./cmd/demo`

Primary Risks:
- Thread-safety / race conditions during concurrent transactions and relay polling.
- In-memory database transaction simulation leaking state across rollbacks or commits.
- Hardcoded sleeps in tests causing flaky test assertions or false race clean passes.
- Discrepancy between architectural design notes (SQLite vs in-memory struct) and implementation.
