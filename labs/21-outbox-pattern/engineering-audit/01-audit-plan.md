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
- research/2026-09-26-outbox-pattern/05-report.md
Main Claims To Verify:
1. Transactional write of domain data (Order) and Outbox log event in single DB transaction.
2. Naive dual-write failure leaves database updated but message omitted.
3. Polling relay dispatches pending outbox entries asynchronously.
4. Downstream consumer handles duplicate delivery idempotently based on message ID.
5. Thread-safe execution under high concurrency.
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions during concurrent mutations of in-memory DB or relay polling.
- Weak assertions or missing edge-case testing in unit test suite.
- Discrepancies between documentation (README.md) and actual implementation details.
