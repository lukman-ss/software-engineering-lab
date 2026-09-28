# Engineering Audit Plan

Target Lab: labs/28-timeouts-and-deadlines
Audit Scope (pipeline override): implementation + tests only. No research/content audit. No code modification.

## Implementation Files
- `internal/deadline/deadline.go` (28 lines)
- `internal/retry/retry.go` (77 lines)
- `internal/circuit/circuit.go` (139 lines)
- `internal/idempotency/idempotency.go` (51 lines)
- `cmd/demo/main.go` (97 lines)
- `go.mod` (module `timeouts-and-deadlines`, go 1.22)

## Tests
- `internal/deadline/deadline_test.go` (3 tests)
- `internal/retry/retry_test.go` (4 tests)
- `internal/circuit/circuit_test.go` (1 test)
- `internal/idempotency/idempotency_test.go` (2 tests)
- `tests/integration_test.go` (2 tests)

## Executable/Demo
- `cmd/demo/main.go` via `go run ./cmd/demo`

## Approved Research Inputs
- Not in scope per pipeline override. Design reference only: `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`.

## Main Claims To Verify
1. Context deadline propagates; work aborts on expiry (deadline).
2. Exponential backoff with full jitter (retry).
3. Circuit breaker CLOSED→OPEN→HALF_OPEN→CLOSED transitions; rejects while OPEN (circuit).
4. Idempotency store deduplicates retried execution; TTL expiry (idempotency).
5. Integrated retry+circuit and retry+idempotency behave as claimed.
6. Race-safe under concurrency; demo output real; README matches code.

## Commands To Run
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
- `go build ./...` (supplemental)
- `go vet ./...` (supplemental)

## Primary Risks
- Goroutine leak in `ExecuteWithBudget` if worker ignores ctx.
- Half-open concurrent flood (no single-trial gate).
- Idempotency check-then-set race + no expired-entry eviction.
- Thin tests: 1 circuit test, no half-open-failure path, no backoff-bound assertion, concurrency tests prove absence of race only, not correctness.
