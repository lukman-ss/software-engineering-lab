# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- internal/outbox/model.go
- internal/outbox/db.go
- internal/outbox/broker.go
- internal/outbox/service.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
Tests:
- tests/outbox_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
Main Claims To Verify:
1. Dual-write vulnerability demonstrates state inconsistency (DB write succeeds, broker write fails).
2. Transactional outbox pattern guarantees atomic write of business state and outbox record in a single database transaction.
3. Outbox polling relay delivers messages from outbox table to broker and updates outbox status asynchronously.
4. Downstream consumer handles messages idempotently via unique event ID tracking.
5. Code compiles cleanly, tests pass with `-race`, demo executes as claimed.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Memory DB simulation concurrency/race issues.
- Mismatch between research/design (SQLite claimed in design vs In-Memory DB in code).
- Incomplete failure/recovery path in relay worker.
