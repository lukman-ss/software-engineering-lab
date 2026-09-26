# Engineering Audit Plan

Target Lab:
labs/21-outbox-pattern

Implementation Files:
- labs/21-outbox-pattern/internal/outbox/model.go
- labs/21-outbox-pattern/internal/outbox/db.go
- labs/21-outbox-pattern/internal/outbox/service.go
- labs/21-outbox-pattern/internal/outbox/broker.go
- labs/21-outbox-pattern/internal/outbox/relay.go
- labs/21-outbox-pattern/internal/outbox/consumer.go
- labs/21-outbox-pattern/cmd/demo/main.go

Tests:
- labs/21-outbox-pattern/tests/outbox_test.go

Executable/Demo:
- labs/21-outbox-pattern/cmd/demo/main.go

Approved Research Inputs:
- Research status approved per engineering/01-design.md
- Research docs in labs/21-outbox-pattern/research/ and labs/21-outbox-pattern/research-revision/

Main Claims To Verify:
1. Atomic persistence: order record + outbox event in single database transaction
2. Transaction rollback: neither order nor outbox persisted on rollback
3. Asynchronous relay: polls pending outbox, dispatches to broker, marks processed
4. Idempotent consumer: tracks event IDs to prevent duplicate processing
5. At-least-once delivery: message re-published on broker/relay failures
6. Dual-write flaw demonstrated: DB commit succeeds but broker fails → inconsistency
7. Race-free concurrent execution: no data races under load

Commands To Run:
- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Cleanup worker for processed outbox records not implemented (claimed in design)
- In-memory map DB used instead of SQLite (design intent vs implementation)
- Concurrent writes test lacks assertions and uses duplicate keys
- No explicit test for broker failure recovery during relay polling
- Relay Stop() double-close could panic (minor fragility)