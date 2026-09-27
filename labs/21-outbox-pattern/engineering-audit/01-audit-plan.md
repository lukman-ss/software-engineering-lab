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
- `research/05-report.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. Dual-write pattern causes state inconsistency when broker write fails after DB commit.
2. Transactional outbox guarantees atomic persistence of domain entity (Order) and outbox record in DB transaction.
3. Polling relay asynchronously reads pending outbox messages and dispatches them to message broker.
4. Outbox records transition status from PENDING to PROCESSED upon successful broker delivery.
5. Downstream consumers achieve idempotency via event deduplication tracking.
6. System supports concurrent outbox creation, polling, and idempotent consumption safely without race conditions.

Commands To Run:
- `cd labs/21-outbox-pattern && go test ./...`
- `cd labs/21-outbox-pattern && go test -race ./...`
- `cd labs/21-outbox-pattern && go run ./cmd/demo`

Primary Risks:
- Memory leaks or race conditions during concurrent DB staging and commit in mock DB.
- Lock contention or deadlock during outbox polling relay execution.
- Discrepancies between README setup/architecture descriptions and actual implementation components.
