# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Audit Scope (pipeline override): implementation + tests only. No research/content audit. No code modifications.
Output Dir: labs/28-timeouts-and-deadlines/engineering-audit-opensource/

Implementation Files:
- internal/deadline/deadline.go (28 lines)
- internal/retry/retry.go (77 lines)
- internal/circuit/circuit.go (139 lines)
- internal/idempotency/idempotency.go (52 lines)
- cmd/demo/main.go (97 lines)

Tests:
- internal/deadline/deadline_test.go (3 tests)
- internal/retry/retry_test.go (5 tests)
- internal/circuit/circuit_test.go (3 tests)
- internal/idempotency/idempotency_test.go (3 tests)
- tests/integration_test.go (2 tests)

Executable/Demo: cmd/demo (4 demos: deadline propagation, backoff jitter, circuit states, idempotent retry)

Main Claims To Verify:
1. Context deadline propagation aborts downstream work on expiry
2. Exponential backoff with full jitter, ctx-aware retry loop
3. Circuit breaker CLOSED->OPEN->HALF_OPEN transitions, thread-safe
4. Idempotency store deduplicates retried operations, TTL expiry, thread-safe
5. Race detector clean; demo output real; README matches code

Commands To Run:
- go vet ./...
- go test -count=1 -v ./...
- go test -race -count=1 ./...
- go run ./cmd/demo

Primary Risks:
- Goroutine leak in ExecuteWithBudget if worker ignores ctx
- Circuit Execute non-atomic (concurrent half-open probes)
- Idempotency check-then-act race under concurrent same-key use
- Thin concurrency test coverage (only idempotency has one)
