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
1. Context deadline propagation & execution budgets (`internal/deadline`).
2. Exponential backoff with full jitter calculation (`internal/retry`).
3. Circuit breaker state transitions (`CLOSED`, `OPEN`, `HALF_OPEN`) with thread-safety (`internal/circuit`).
4. Idempotency store with TTL lazy-eviction and concurrent access safety (`internal/idempotency`).
5. All tests compile and pass without race conditions (`go test -race ./...`).
6. Demo binary runs cleanly to completion (`go run ./cmd/demo`).
7. README documentation accurately describes codebase structure and execution commands.

Commands To Run:
- `go test -count=1 ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Go goroutine leaks in deadline execution (`ExecuteWithBudget`).
- Race conditions during concurrent mutations of Breaker or Idempotency Store.
- Docs/Code discrepancies in public APIs or configuration defaults.
