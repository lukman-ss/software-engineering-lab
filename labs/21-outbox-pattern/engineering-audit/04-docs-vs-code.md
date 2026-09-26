# Docs vs Code Audit

## Comparisons

### 1. Architecture Claims in README.md vs Code
- **README Claim**: `internal/outbox/db.go` provides in-memory transactional database simulating `BeginTx`, `Commit`, and `Rollback` across `orders` and `outbox` records.
  - **Code Reality**: Exact match (`db.go`).
- **README Claim**: `internal/outbox/broker.go` provides thread-safe mock message broker simulating publish failures and event reception.
  - **Code Reality**: Exact match (`broker.go` with mutex protection and `failNext` flag).
- **README Claim**: `internal/outbox/service.go` provides business logic comparing naive dual-write vs atomic outbox writes.
  - **Code Reality**: Exact match (`service.go` with `CreateOrderWithOutbox` and `CreateOrderDualWriteNaive`).
- **README Claim**: `internal/outbox/relay.go` provides asynchronous polling worker querying pending outbox records and dispatching them to the broker.
  - **Code Reality**: Exact match (`relay.go` with ticker-based polling).
- **README Claim**: `internal/outbox/consumer.go` provides subscriber enforcing idempotency through event ID tracking.
  - **Code Reality**: Exact match (`consumer.go` with set of processed event IDs).

### 2. Commands Verification
- `go test ./...` in README: Verified and passes.
- `go test -race ./...` in README: Verified and passes cleanly with no race conditions detected.
- `go run ./cmd/demo` in README: Verified and runs successfully demonstrating dual-write problem, outbox resolution, and duplicate suppression.

## Mismatch Findings
- `DOC_CODE_MISMATCH`: None
- `TEST_CLAIM_MISMATCH`: None
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None
