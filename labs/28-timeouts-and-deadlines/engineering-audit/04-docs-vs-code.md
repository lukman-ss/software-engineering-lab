# Docs vs Code Audit

## Comparison Matrix

| Claim / Doc Reference | Code Implementation | Test / Demo Evidence | Status |
| :--- | :--- | :--- | :--- |
| `internal/deadline`: Context deadline propagation & execution within time budget | `ExecuteWithBudget(ctx, budget, fn)` in `internal/deadline/deadline.go` | `deadline_test.go`, Demo Section 1 | MATCH |
| `internal/retry`: Exponential backoff with full jitter | `Retrier.Do` & `CalculateBackoff` in `internal/retry/retry.go` | `retry_test.go`, Demo Section 2 | MATCH |
| `internal/circuit`: 3-state machine (CLOSED, OPEN, HALF_OPEN) | `Breaker` in `internal/circuit/circuit.go` | `circuit_test.go`, Demo Section 3 | MATCH |
| `internal/idempotency`: In-memory deduplication store | `Store` with TTL in `internal/idempotency/idempotency.go` | `idempotency_test.go`, Demo Section 4 | MATCH |
| `README.md` execution instructions (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo`) | `go.mod`, package structure in `internal/`, `cmd/demo/main.go` | All commands run verbatim without error | MATCH |

## Discrepancy Findings

- `DOC_CODE_MISMATCH`: None detected.
- `TEST_CLAIM_MISMATCH`: None detected.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None detected.
