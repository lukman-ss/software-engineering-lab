# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Implementation Files:
- internal/deadline/deadline.go
- internal/retry/retry.go
- internal/circuit/circuit.go
- internal/idempotency/idempotency.go

Tests:
- internal/deadline/deadline_test.go
- internal/retry/retry_test.go
- internal/circuit/circuit_test.go
- internal/idempotency/idempotency_test.go
- tests/integration_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- research/03-evidence.md

Main Claims To Verify:
1. Deadline propagation with explicit budget management using context.WithTimeout.
2. Exponential backoff with full jitter in retries.
3. Circuit breaker state machine transitions (CLOSED -> OPEN -> HALF_OPEN -> CLOSED/OPEN).
4. In-memory idempotency deduplication with TTL expiry.
5. Integration of retry with circuit breaker and retry with idempotency.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Worker goroutine leak in `ExecuteWithBudget` when fn blocks indefinitely after timeout.
- State check race conditions in circuit breaker `State()` method.
- Memory leak in idempotency store due to lack of passive/active cleanup background task.
- Unhandled non-retryable errors in retry logic.
