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
- tests/outbox_test.go (6 tests)
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs: SKIPPED per PIPELINE OVERRIDE (implementation + tests only)
Main Claims To Verify:
1. Order + outbox event persist atomically in one Tx (Commit/Rollback)
2. Naive dual-write loses events on broker failure (inconsistency)
3. Polling relay dispatches PENDING outbox msgs to broker, marks PROCESSED
4. Consumer is idempotent on duplicate delivery (at-least-once safe)
5. Rollback persists nothing; purge removes only PROCESSED
6. Concurrent writes are race-safe
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Relay publish-then-mark non-atomic (duplicate publish window)
- ConcurrentWrites test uses single shared order ID + zero assertions (weak proof)
- Relay Stop double-close / multi-Start lifecycle unsafe
