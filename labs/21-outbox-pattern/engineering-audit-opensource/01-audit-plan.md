# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- internal/outbox/db.go
- internal/outbox/broker.go
- internal/outbox/service.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
- internal/outbox/model.go
- cmd/demo/main.go
Tests: tests/outbox_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: research/ (specifically the research that approved the lab, but we are not auditing research per pipeline override)
Main Claims To Verify:
1. Atomic persistence of order and outbox event in a single transaction.
2. Rollback on transaction failure leaves no partial state.
3. Relay polls pending outbox events, dispatches to broker, and marks as processed.
4. Consumer enforces idempotency via event ID tracking.
5. Dual-write naive approach demonstrates inconsistency when broker fails.
6. Relay retries after broker failure and eventually succeeds.
7. Purge processed outbox records.
8. Concurrent writes and consumer processing are safe (no races).
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Relay may publish but fail to mark processed, leading to duplicate sends (mitigated by idempotent consumer).
- Broker failure during publish may cause temporary inconsistency until retry.
- Consumer idempotency relies on in-memory map; not persistent across restarts (but demo does not require persistence).
- The in-memory DB uses mutexes; potential for deadlock if not careful (but implementation avoids nested locks).