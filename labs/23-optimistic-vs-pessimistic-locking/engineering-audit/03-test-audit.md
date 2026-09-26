# Test Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Coverage Breakdown

| Test Name | File | Scenarios Covered | Concurrency Level | Assertion Validation | Result |
|---|---|---|---|---|---|
| `TestNaiveLostUpdate` | `tests/locking_test.go:10` | Concurrency anomaly, lost update | 50 goroutines | Asserts `stock != 50` | PASS |
| `TestPessimisticLocking` | `tests/locking_test.go:40` | Pessimistic serialization | 50 goroutines | Asserts `stock == 50` | PASS |
| `TestPessimisticLockingInsufficientStock` | `tests/locking_test.go:69` | Negative branch, bounds check | 1 caller (sequential) | Asserts `ErrInsufficientStock` | PASS |
| `TestOptimisticLockingConflict` | `tests/locking_test.go:85` | Conflict detection, zero retry | 20 goroutines | Asserts `success + stock == 100` and `conflicts > 0` | PASS |
| `TestOptimisticLockingWithRetry` | `tests/locking_test.go:123` | Retry convergence with backoff | 20 goroutines | Asserts `stock == 100 - Optimistically` | PASS |
| `TestAtomicConditionalUpdate` | `tests/locking_test.go:147` | Single-statement atomic update | 50 goroutines | Asserts `stock == 50` | PASS |

## Test Execution Verification

- Unit Tests: PASS (`go test -v -count=1 ./...`)
- Race Detector: PASS (`go test -v -race -count=1 ./...` with 0 race warnings)
- Determinism: Invariant checks assert exact conservation of stock across concurrent worker sets.
