# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- internal/outbox/model.go
- internal/outbox/db.go
- internal/outbox/broker.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
- internal/outbox/service.go
Tests:
- tests/outbox_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md (and research-revision/03-revision-result.md)
Main Claims To Verify:
1. Atomic persistence: Entity state change and outbox message are written atomically in a single local transaction.
2. Rollback safety: Aborted transactions persist neither entity changes nor outbox records.
3. Decoupled polling relay: Asynchronous relay polls pending outbox records and publishes to message broker.
4. Outbox state transition: Upon successful broker publish, outbox message status transitions from PENDING to PROCESSED.
5. Idempotent consumer / At-least-once delivery: Downstream consumer handles duplicate event deliveries idempotently via event ID deduplication.
6. Dual-write vulnerability comparison: Naive dual write leaves database and message broker in an inconsistent state when broker is unavailable.
7. Concurrency safety: Zero race conditions during concurrent writes and relay processing (`go test -race ./...`).

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or deadlocks in in-memory DB / Transaction mutex synchronization.
- Mock relay polling leaks or dangling goroutines upon test completion.
- Discrepancies between demo execution output and documented README claims.
