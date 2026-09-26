# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Implementation Files:
- internal/deadline/deadline.go
- internal/retry/retry.go
- internal/circuit/circuit.go
- internal/idempotency/idempotency.go
- cmd/demo/main.go
- go.mod

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
- research-audit/07-verdict.md (Status: APPROVED)
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. Context deadline propagation respects parent context and overrides longer local budget.
2. Exponential backoff with full jitter calculates randomized delay within `[0, min(maxBackoff, baseBackoff * 2^(attempt-1))]`.
3. Circuit breaker transitions across CLOSED -> OPEN -> HALF_OPEN -> CLOSED with failure/success thresholds and cooldown window.
4. Idempotency store provides concurrent key-response deduplication with TTL expiration.
5. Integration between retries, circuit breaker, and idempotency store safely terminates without race conditions or retry storms.

Commands To Run:
- go test -v -count=1 ./...
- go test -race -v -count=1 ./...
- go run ./cmd/demo

Primary Risks:
- Goroutine leak in `ExecuteWithBudget` if worker function blocks indefinitely after timeout.
- Concurrency race conditions in state transitions (`Breaker`) or cache lookups (`Store`).
- Discrepancies between advertised claims in README/engineering docs vs observed runtime behavior.
