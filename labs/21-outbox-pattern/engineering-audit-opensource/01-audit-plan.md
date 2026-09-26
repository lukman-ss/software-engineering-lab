# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern (Transactional Outbox Pattern)
Implementation Files: 
- internal/outbox/db.go
- internal/outbox/broker.go
- internal/outbox/service.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
- internal/outbox/model.go
Tests: tests/outbox_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: Not auditing research per pipeline override (implementation and tests only)
Main Claims To Verify:
1. Atomic persistence of Order and OutboxMessage in same transaction (CreateOrderWithOutbox)
2. Dual-write flaw demonstrated (CreateOrderDualWriteNaive fails after DB commit but before broker publish)
3. Asynchronous polling relay dispatching pending messages
4. Consumer idempotency via message ID tracking
5. Concurrency safety under race detector
6. Recovery from broker publish failures (message remains pending for retry)
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Weak concurrency test (TestTransactionalOutbox_ConcurrentWrites has no assertions)
- Missing test for broker failure during relay retry scenario
- Latent panic if Relay.Stop() called twice (not exercised in tests)