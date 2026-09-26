# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- `internal/outbox/model.go`: Domain and outbox data models and status enums.
- `internal/outbox/db.go`: In-memory transactional database engine supporting `BeginTx`, `SaveOrder`, `SaveOutbox`, `Commit`, `Rollback`, and outbox querying.
- `internal/outbox/broker.go`: Thread-safe mock message broker simulating event publishing and failure injection.
- `internal/outbox/service.go`: Business logic contrasting dual-write and transactional outbox approaches.
- `internal/outbox/relay.go`: Asynchronous background polling worker querying pending outbox records and publishing to broker.
- `internal/outbox/consumer.go`: Message consumer implementing idempotent deduplication based on event IDs.

Tests:
- `tests/outbox_test.go`: End-to-end integration and concurrency unit tests.

Executable/Demo:
- `cmd/demo/main.go`: Runnable demonstration showing dual-write failure, outbox atomicity, and consumer idempotency.

Approved Research Inputs:
- `research-audit/07-verdict.md` (APPROVED)
- `research/2026-09-26-outbox-pattern/05-report.md`

Main Claims To Verify:
1. Atomicity of domain state changes and outbox event persistence in a single transaction.
2. Rollback safety: Outbox events must not persist if the transaction aborts.
3. Decoupled polling relay worker dispatches events asynchronously to a message broker.
4. Downstream consumer handles at-least-once duplicate delivery idempotently.
5. Dual-write failure path correctly demonstrates state inconsistency when direct publish fails.
6. Thread-safety under concurrent writes and race condition absence.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions between transaction commits, relay polling, and consumer deduplication.
- Sleep-based test timing flakiness.
- Discrepancies between documentation and implementation mechanics.
