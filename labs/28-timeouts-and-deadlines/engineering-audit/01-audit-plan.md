# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Implementation Files:
- `internal/deadline/deadline.go`
- `internal/retry/retry.go`
- `internal/circuit/circuit.go`
- `internal/idempotency/idempotency.go`
Tests:
- `internal/deadline/deadline_test.go`
- `internal/retry/retry_test.go`
- `internal/circuit/circuit_test.go`
- `internal/idempotency/idempotency_test.go`
- `tests/integration_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research/01-plan.md`
Main Claims To Verify:
1. Context deadline propagation enforces execution time budgets and inherits parent timeouts.
2. Retry mechanism uses exponential backoff with full jitter to avoid synchronization storms.
3. Circuit breaker transitions between CLOSED, OPEN, and HALF_OPEN correctly based on failure/success thresholds and cooldown timeouts.
4. Idempotency key store provides thread-safe response caching and lazy eviction on expiration.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent state transitions in Circuit Breaker or Idempotency Store.
- Goroutine leakage in `ExecuteWithBudget` when worker functions block indefinitely.
- Improper math or bounds handling in Full Jitter backoff calculations.
