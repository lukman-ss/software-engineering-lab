# Documentation vs Code Audit

Target Lab: labs/21-outbox-pattern

## Itemized Comparison

### 1. Architectural Components
- `README.md` lists:
  - `internal/outbox/db.go`
  - `internal/outbox/broker.go`
  - `internal/outbox/service.go`
  - `internal/outbox/relay.go`
  - `internal/outbox/consumer.go`
- Codebase check: All 5 files exist and implement the exact roles described.
- Status: PASS

### 2. Running Instructions
- `README.md` provides commands:
  - `go test ./...`
  - `go test -race ./...`
  - `go run ./cmd/demo`
- Codebase check: All commands execute cleanly with zero errors or race conditions.
- Status: PASS

### 3. Claims vs Implementation
- Atomicity claim: Domain entity and outbox message committed together.
  - Implemented in `internal/outbox/service.go:52` via `tx.Commit()`. Tested in `TestTransactionalOutbox_HappyPath`.
- Decoupled Polling Relay claim: Background worker queries pending records and forwards to broker.
  - Implemented in `internal/outbox/relay.go`. Tested in `TestTransactionalOutbox_HappyPath`.
- Consumer Idempotency claim: Consumer skips duplicate deliveries.
  - Implemented in `internal/outbox/consumer.go:23-26`. Tested in `TestTransactionalOutbox_Idempotency_DuplicateDelivery`.
- Dual-Write Vulnerability claim: Naive dual write leaves DB inconsistent if broker fails.
  - Implemented in `internal/outbox/service.go:57-89`. Tested in `TestDualWriteProblem_Failure` and showcased in `cmd/demo/main.go`.

## Mismatches Found
- None. Documentation directly matches the implementation and execution.
