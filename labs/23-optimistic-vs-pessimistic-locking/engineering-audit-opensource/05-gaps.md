# Gap Analysis

## Gap 1: Missing Edge Case Tests — Zero/Negative Quantity

Gap Type: MISSING_EDGE_CASE
Location: tests/locking_test.go — no test for `qty <= 0`
Severity: MEDIUM
Description: All five deduction methods (`NaiveDeduct`, `PessimisticDeduct`, `OptimisticDeduct`, `AtomicDeduct`, and Service wrappers) check `if qty <= 0` and return `ErrInvalidQuantity`. No test verifies this guard for any strategy.
Why it matters: Without testing, a regression could remove the guard and allow zero/negative deductions that corrupt inventory.
When to add: Add a test case for `DeductXxx(id, 0)` and `DeductXxx(id, -1)` asserting `ErrInvalidQuantity`.

## Gap 2: Missing Edge Case Tests — Non-Existent Product

Gap Type: MISSING_EDGE_CASE
Location: tests/locking_test.go — no test for product ID not in store
Severity: MEDIUM
Description: `Get` returns `ErrNotFound`, `PessimisticDeduct` returns `ErrNotFound`, `OptimisticDeduct` returns `ErrNotFound`, `AtomicDeduct` returns `ErrNotFound`. None of these paths are tested.
Why it matters: Error propagation for missing resources is untested. A regression could panic (e.g., nil pointer dereference) instead of returning `ErrNotFound`.
When to add: Add a test that calls each deduction strategy with a non-existent product ID and asserts `ErrNotFound`.

## Gap 3: OptimisticLockingWithRetry Does Not Test Retry Exhaustion

Gap Type: MISSING_EDGE_CASE
Location: tests/locking_test.go:123-145 — `TestOptimisticLockingWithRetry`
Severity: LOW
Description: This test uses 20 goroutines with `maxRetries=10` and stock=100, so all requests always succeed. The retry exhaustion path (where all retries fail and `ErrOptimisticLock` is returned) is never exercised.
Why it matters: The backoff/retry logic's failure path is untested. A bug in the retry loop (e.g., infinite loop, wrong error type returned) would not be caught.
When to add: Add a test with extreme contention (e.g., 100 goroutines, stock=10, maxRetries=2) where some requests must fail after exhausting retries.

## Gap 4: No Multi-Product Concurrency Test (Lock Isolation)

Gap Type: MISSING_EDGE_CASE
Location: tests/locking_test.go — all tests use a single product per store
Severity: LOW
Description: No test verifies that operations on different products can proceed concurrently without interfering. For pessimistic locking (per-row locks), this should allow parallelism across products. For atomic/naive (global engine lock), it would serialize.
Why it matters: The per-row locking design (a key feature of pessimistic locking) is not validated. A regression that accidentally uses a global lock instead of per-row locks would not be detected.
When to add: Add a test with multiple products and concurrent goroutines operating on different product IDs.

## Gap 5: No Test for Insufficient Stock via Atomic Strategy Under Concurrency

Gap Type: MISSING_EDGE_CASE
Location: tests/locking_test.go:147-171 — `TestAtomicConditionalUpdate`
Severity: LOW
Description: The atomic test only validates the happy path (50 deducts of 1 from 100). It does not test the insufficient-stock boundary (e.g., 55 goroutines each deducting 2, stock 100 — 50 should succeed, 5 should fail).
Why it matters: The atomic conditional update (`WHERE stock >= qty`) is supposed to safely reject over-deduction. This rejection path under concurrency is untested for the atomic strategy.
When to add: Add a test with concurrent atomic deductions where some must fail due to insufficient stock, verifying the exact success/failure counts.

## Gap 6: No Test for Insufficient Stock via Optimistic Strategy Under Concurrency

Gap Type: MISSING_EDGE_CASE
Location: tests/locking_test.go — no such test
Severity: LOW
Description: The optimistic tests (`TestOptimisticLockingConflict`, `TestOptimisticLockingWithRetry`) do not test the insufficient-stock path under concurrency. The stock check in `OptimisticDeduct` is done on a stale read, and if stock runs out, subsequent retries should fail with `ErrInsufficientStock`.
Why it matters: The interaction between version conflicts and stock exhaustion in the retry loop is untested.
When to add: Add a test where stock is barely sufficient (e.g., stock=20, 50 goroutines each deducting 1) — some will succeed, some will fail on version conflict and retry, some will ultimately fail with insufficient stock.

## Gap 7: Version Not Incremented by Non-Optimistic Write Paths

Gap Type: RESEARCH_MISMATCH (implementation fidelity)
Location: internal/inventory/store.go — `NaiveDeduct` (line 84), `PessimisticDeduct` (line 111), `AtomicDeduct` (line 170)
Severity: MEDIUM
Description: In real SQL optimistic locking, the version/timestamp column is bumped on every write via `UPDATE ... SET version = version + 1 WHERE ...`. In this simulation, only `OptimisticDeduct` (line 149: `curr.Version++`) increments the version. The other write paths modify `curr.Stock` without bumping version. This means if strategies were mixed on the same product, optimistic locking would not detect changes from other strategies.
Why it matters: This is a simulation fidelity issue. It does not affect the current tests (each uses a dedicated store with one strategy), but it means the optimistic locking simulation is not fully sound in a mixed-operation scenario.
When to add: Either increment version in all write paths (to match SQL semantics) or document the limitation more explicitly in the implementation notes.

## Gap 8: Demo Output Timing-Dependent Values Recorded as Static

Gap Type: DOC_CODE_MISMATCH
Location: engineering/03-execution-result.md:29, 80
Severity: LOW
Description: The recorded execution result shows `[4] Total Attempted Conflicts Retried: 61` and `Elapsed Time: 49.859125ms`. Actual execution produces varying values (e.g., 43 conflicts, 26ms on the audit run; 45 conflicts on an earlier run). The core invariants (stock=80, 20 successful deductions) are stable, but timing-dependent metrics are recorded as fixed values.
Why it matters: A reader reproducing the lab will see different numbers and may think something is wrong. The conflict count and elapsed time are inherently non-deterministic.
When to add: Add a note in the execution result that conflict counts and elapsed times are timing-dependent and will vary per run.

## Gap 9: Demo Not Verified by Automated Test

Gap Type: MISSING_TEST
Location: tests/locking_test.go — no demo output verification
Severity: LOW
Description: The demo (`cmd/demo/main.go`) prints specific claims about success/failure and "LOST UPDATE DETECTED!" but no test programmatically verifies the demo's output or invariants. The demo is only verified by manual execution recorded in `03-execution-result.md`.
Why it matters: If the demo's output logic is broken (e.g., swapped scenarios, wrong assertions), it would not be caught by CI.
When to add: Consider adding a test that runs the demo logic as a function and verifies the output values.

## Gap 10: Seed Does Not Validate Input

Gap Type: MISSING_EDGE_CASE
Location: internal/inventory/store.go:28-40 — `Seed`
Severity: LOW
Description: `Seed` accepts any `stock` value (including negative) and always sets `Version: 1`. No input validation is performed. If stock were negative, all subsequent deduction attempts would pass the `p.Stock < qty` check for small positive qty values.
Why it matters: While no test or caller currently passes negative stock, defensive validation would prevent misuse.
When to add: Add input validation to `Seed` (e.g., `if stock < 0 { panic(...) }`) or document the precondition.

## Gap Summary

| Gap ID | Type | Severity | Status |
|---|---|---|---|
| 1 | MISSING_EDGE_CASE | MEDIUM | Open |
| 2 | MISSING_EDGE_CASE | MEDIUM | Open |
| 3 | MISSING_EDGE_CASE | LOW | Open |
| 4 | MISSING_EDGE_CASE | LOW | Open |
| 5 | MISSING_EDGE_CASE | LOW | Open |
| 6 | MISSING_EDGE_CASE | LOW | Open |
| 7 | RESEARCH_MISMATCH | MEDIUM | Open |
| 8 | DOC_CODE_MISMATCH | LOW | Open |
| 9 | MISSING_TEST | LOW | Open |
| 10 | MISSING_EDGE_CASE | LOW | Open |

Total gaps: 10
- CRITICAL: 0
- HIGH: 0
- MEDIUM: 3 (Gaps 1, 2, 7)
- LOW: 7 (Gaps 3, 4, 5, 6, 8, 9, 10)

All gaps are non-blocking for the core claim verification (the three remediation strategies prevent lost updates). The MEDIUM gaps relate to incomplete edge case testing and simulation fidelity, not to incorrect core behavior.