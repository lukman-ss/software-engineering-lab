# Docs vs Code Audit

Target Lab: labs/28-timeouts-and-deadlines

## Comparison Matrix

| Component / Claim | README Claim | Code Implementation | Test Verification | Demo Execution | Alignment Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Deadline Budget** | `internal/deadline`: Context deadline propagation & budget execution | `deadline.ExecuteWithBudget` using `context.WithTimeout` | Verified in `deadline_test.go` | Demo 1 output | PASS |
| **Exponential Backoff + Jitter** | `internal/retry`: Exponential backoff with full jitter | `retry.Retrier` with `rand.Float64() * temp` | Verified in `retry_test.go` | Demo 2 output | PASS |
| **Circuit Breaker** | `internal/circuit`: State machine (CLOSED, OPEN, HALF_OPEN) | `circuit.Breaker` with mutex thread-safety | Verified in `circuit_test.go` | Demo 3 output | PASS |
| **Idempotency Store** | `internal/idempotency`: In-memory deduplication store | `idempotency.Store` with lazy TTL eviction | Verified in `idempotency_test.go` | Demo 4 output | PASS |

## Discrepancy Checks

- `DOC_CODE_MISMATCH`: None detected. README component list and commands match implementation exactly.
- `TEST_CLAIM_MISMATCH`: None detected. Tests cover unit and integration scenarios described.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None detected. Implementation follows approved research patterns.
