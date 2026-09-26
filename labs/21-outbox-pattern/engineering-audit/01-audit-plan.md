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
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/05-report.md
- research/06-open-questions.md
Main Claims To Verify:
1. Dual-write pattern creates inconsistency on broker failure.
2. Transactional Outbox atomically commits domain entity and outbox message in a single database transaction.
3. Database transaction rollback discards both domain entity and outbox message.
4. Outbox Relay worker polls pending outbox records, dispatches them to message broker, and marks outbox status as PROCESSED.
5. Downstream consumer achieves idempotency by tracking processed event IDs to filter duplicate deliveries.
6. Execution and concurrency safety (`go test -race ./...`).
7. Accuracy of design doc (`engineering/01-design.md`), execution results (`engineering/03-execution-result.md`), and README (`README.md`) against actual codebase.
Commands To Run:
- go test -v ./...
- go test -count=1 -race ./...
- go run ./cmd/demo
Primary Risks:
- In-memory mock DB implementation might not properly enforce transaction isolation or lock safety during relay polling.
- Missing cleanup/archive implementation claimed in success criteria of design doc.
- Race conditions during concurrent outbox polling or writes.
