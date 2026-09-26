# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files: internal/outbox/broker.go, internal/outbox/consumer.go, internal/outbox/db.go, internal/outbox/model.go, internal/outbox/relay.go, internal/outbox/service.go, cmd/demo/main.go
Tests: tests/outbox_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/ (including research/2026-09-26-outbox-pattern/ and research/runs/2026-09-25-the-outbox-pattern/)
Main Claims To Verify:
1. Atomic persistence between domain entities and event logs (outbox pattern ensures both are saved in same transaction)
2. Decoupled polling relay dispatch to a message broker (relay polls outbox and publishes to broker without blocking transaction)
3. Downstream consumer idempotency (consumer tracks processed event IDs to avoid duplicate processing)
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in relay polling and broker publish
- Improper transaction rollback on failure
- Idempotency key collision or missing tracking
- Deadlock in concurrent access to shared resources