# Gap Analysis

## GAP-1

Type: MISSING_TEST
Severity: MEDIUM
Location: tests/locking_test.go (all strategies)
Description: No test exercises ErrInvalidQuantity (qty <= 0). All four Deduct paths guard `qty <= 0` in store.go, but zero/negative quantity is never passed in any test.
Impact: Negative-case input validation at the trust boundary is unproven by tests (code inspection shows guards exist and are correct).

## GAP-2

Type: MISSING_TEST
Severity: MEDIUM
Location: tests/locking_test.go (all strategies)
Description: No test exercises ErrNotFound (deduct on unseeded product id). Get, PessimisticDeduct, OptimisticDeduct, and AtomicDeduct all have not-found branches that are never executed.
Impact: Failure-path error propagation for missing rows is unproven by tests.

## GAP-3

Type: MISSING_TEST
Severity: MEDIUM
Location: tests/locking_test.go (naive, optimistic direct, optimistic retry, atomic)
Description: Insufficient-stock failure path is tested only for pessimistic locking (TestPessimisticLockingInsufficientStock). Naive, OptimisticDeduct, DeductOptimisticWithRetry, and AtomicDeduct never get an insufficient-stock test.
Impact: Core claimed guard (`WHERE stock >= qty` semantics) is proven for the happy path via invariants but its rejection branch is untested in 3 of 4 strategies.

## GAP-4

Type: MISSING_TEST
Severity: MEDIUM
Location: tests/locking_test.go (retry logic, service.go:28-45)
Description: Retry exhaustion is never tested. DeductOptimisticWithRetry with maxRetries=0 (or contended retry that must return ErrOptimisticLock) has no test; the only retry test uses maxRetries=10 and ignores returned errors (`_ =`).
Impact: The bounded-retry terminal behavior (return ErrOptimisticLock after attempts exhausted) is unproven.

## GAP-5

Type: MISSING_EDGE_CASE
Severity: LOW
Location: tests/locking_test.go:69-83 (TestPessimisticLockingInsufficientStock)
Description: After the failed second deduct, the test does not read back stock to verify it is unchanged (still 0). Failure atomicity (no partial write on rejection) is asserted nowhere.
Impact: Minor; code inspection shows the check happens before mutation under lock, so correctness is evident but unasserted.

## GAP-6

Type: MISSING_EDGE_CASE
Severity: LOW
Location: tests/locking_test.go:117-119 (TestOptimisticLockingConflict, `conflictCount == 0` assertion)
Description: The assertion that at least one conflict occurred depends on goroutine interleaving timing (50us sleep makes it virtually certain with 20 goroutines, but it is not guaranteed by synchronization). Observed run: 1 success / 19 conflicts. A slow or single-CPU scheduler could in principle serialize all 20 and fail the test spuriously.
Impact: Test flakiness risk only; the stock invariant assertion on line 114 already proves the safety property deterministically. No implementation defect.

## GAP-7

Type: MISSING_TEST
Severity: LOW
Location: tests/locking_test.go:123-145 (TestOptimisticLockingWithRetry)
Description: Returned errors from DeductOptimisticWithRetry are discarded (`_ =`); the test never asserts all 20 goroutines succeeded. It asserts only the stock-vs-counter consistency invariant.
Impact: Minor; with maxRetries=10 and 20 contenders convergence is reliable (observed 20/20), but silent failures would go unnoticed.

## Explicitly checked, no gap

- RACE_CONDITION: none. `go test -race -count=1 ./...` passes (ok, 1.263s). Direct reads of atomic counters in test/demo occur only after `wg.Wait()`, so no concurrent access.
- FAKE_DEMO / FAKE_BENCHMARK / UNVERIFIED_RESULT: none. Demo re-executed live (exit 0); output matches engineering/03-execution-result.md modulo expected nondeterminism (conflict counts 47 vs 61, elapsed ms). No benchmarks claimed.
- BROKEN_IMPLEMENTATION: none. All four strategies behave as claimed.
- DOC_CODE_MISMATCH / RESEARCH_MISMATCH / IMPLEMENTATION_OVERCLAIM: none. README and engineering notes match code, tests, and demo output. Research audit out of scope per pipeline override.
