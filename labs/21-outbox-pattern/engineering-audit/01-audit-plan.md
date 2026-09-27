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
- research/05-report.md (or research/runs/2026-09-25-the-outbox-pattern/05-report.md)
Main Claims To Verify:
- Direct write outside of transaction causes data inconsistency when message broker fails (Dual-Write problem).
- Atomic persistence of domain entity (`Order`) and `OutboxMessage` within a single transaction boundary.
- Database rollback discards both domain mutation and outbox event log.
- Asynchronous polling relay reads pending outbox events, dispatches them to broker, and marks status as processed.
- Broker failures trigger relay retry without losing events.
- Consumer handles duplicate message delivery idempotently using event ID tracking.
- Concurrent order generation and relay dispatch operate cleanly with zero data races.
- Processed outbox purge/cleanup capability works correctly.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions between transaction commits and background polling relay.
- Deadlocks or concurrent map access in simulated in-memory DB / Broker / Consumer.
- Incomplete duplicate handling or state leakage in consumer deduplication.
- Unhandled errors during message serialization, broker publishing, or status updating.
