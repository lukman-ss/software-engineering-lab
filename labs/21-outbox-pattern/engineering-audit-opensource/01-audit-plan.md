# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- labs/21-outbox-pattern/internal/outbox/db.go
- labs/21-outbox-pattern/internal/outbox/broker.go
- labs/21-outbox-pattern/internal/outbox/service.go
- labs/21-outbox-pattern/internal/outbox/relay.go
- labs/21-outbox-pattern/internal/outbox/consumer.go
- labs/21-outbox-pattern/internal/outbox/model.go
Tests: labs/21-outbox-pattern/tests/outbox_test.go
Executable/Demo: labs/21-outbox-pattern/cmd/demo/main.go
Approved Research Inputs: (Not audited in this stage per pipeline override)
Main Claims To Verify:
1. Atomic persistence between domain entities and event logs (outbox pattern)
2. Decoupled polling relay dispatch to a message broker
3. Downstream consumer idempotency
4. The dual-write problem leads to inconsistency
5. The outbox pattern prevents inconsistency even under failures
6. Concurrent writes are handled safely (no races)
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in the in-memory DB or broker
- Incorrect state transitions in outbox processing
- Missing error handling in relay or service
- Demo might not reflect actual test conditions