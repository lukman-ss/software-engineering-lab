# Docs vs Code Audit

Target Lab: labs/21-outbox-pattern

## Documentation Review

File: `README.md`

### Claim 1: Architecture Components
- README lists:
  - `internal/outbox/db.go`: In-memory transactional database simulating `BeginTx`, `Commit`, and `Rollback` across `orders` and `outbox` records.
  - `internal/outbox/broker.go`: Thread-safe mock message broker simulating publish failures and event reception.
  - `internal/outbox/service.go`: Business logic comparing naive dual-write vs atomic outbox writes.
  - `internal/outbox/relay.go`: Asynchronous polling worker querying pending outbox records and dispatching them to the broker.
  - `internal/outbox/consumer.go`: Subscriber enforcing idempotency through event ID tracking.
- Code Reality: Exact match. All files exist at the specified paths with described responsibilities.
- Alignment: PASS.

### Claim 2: Running Tests
- README specifies:
  ```bash
  go test ./...
  go test -race ./...
  ```
- Code Reality: Both commands execute cleanly and pass with 0 race warnings.
- Alignment: PASS.

### Claim 3: Running Demo
- README specifies:
  ```bash
  go run ./cmd/demo
  ```
- Code Reality: Executes cleanly, demonstrating dual-write inconsistency, transactional outbox atomicity, and idempotent deduplication.
- Alignment: PASS.

## Summary of Discrepancies
- DOC_CODE_MISMATCH: 0
- TEST_CLAIM_MISMATCH: 0
- RESEARCH_IMPLEMENTATION_MISMATCH: 0
