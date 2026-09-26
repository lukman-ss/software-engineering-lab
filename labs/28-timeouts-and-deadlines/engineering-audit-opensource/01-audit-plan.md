# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Implementation Files:
- `internal/deadline/deadline.go`
- `internal/retry/retry.go`
- `internal/circuit/circuit.go`
- `internal/idempotency/idempotency.go`
- `cmd/demo/main.go`
- `go.mod`

Tests:
- `internal/deadline/deadline_test.go`
- `internal/retry/retry_test.go`
- `internal/circuit/circuit_test.go`
- `internal/idempotency/idempotency_test.go`
- `tests/integration_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- PIPELINE OVERRIDE active (audit implementation and tests only, do not audit research/content in this stage)
- Engineering specification: `engineering/01-design.md`, `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Compilation: package and demo compile with zero build errors under Go 1.22+.
2. Context deadline and budget propagation: `ExecuteWithBudget` enforces time bounds and inherits parent context cancellation.
3. Exponential backoff with full jitter: backoff increases exponentially up to max backoff, randomized uniformly across `[0, backoff]`.
4. Circuit breaker: transitions across `CLOSED` -> `OPEN` -> `HALF_OPEN` -> `CLOSED` based on configurable failure/success thresholds and cooldowns.
5. Idempotency store: thread-safe get/set with TTL expiration.
6. Race safety: passes `go test -race ./...` without data races.
7. Demo accuracy: `cmd/demo/main.go` execution output matches claimed behavior in `engineering/03-execution-result.md` and `README.md`.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Goroutine leakage in `ExecuteWithBudget` if worker does not monitor context cancellation.
- Race conditions or concurrency deadlocks in circuit breaker or idempotency store under contention.
- In-memory idempotency lack of TTL cleanup background loop (memory growth under long runs).
- Flaky tests caused by strict timing assertions in sleep/cooldown intervals.
