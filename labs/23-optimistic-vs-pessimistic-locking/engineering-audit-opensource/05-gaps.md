# Gap Analysis

## Gap 1
Location: tests/locking_test.go:34
Type: TEST_CLAIM_MISMATCH
Severity: MEDIUM
Claim: Test proves lost-update anomaly.
Issue: Assertion `if p.Stock == 50` is too weak. Stock 99 satisfies the test but
the test never asserts that 50 calls returned success while only 1 decrement applied.
Recommendation: Assert `p.Stock > 50` (strictly, to confirm lost updates) and
optionally verify all call returns were nil.

## Gap 2
Location: tests/locking_test.go:141
Type: TEST_CLAIM_MISMATCH
Severity: MEDIUM
Claim: Test proves optimistic retry convergence.
Issue: Asserts against `store.Optimistically` (implementation counter) rather than
an independent invariant. A bug in counter increment would pass silently.
Recommendation: Use invariant `successCount + finalStock == initialStock`.

## Gap 3
Location: tests/locking_test.go
Type: MISSING_TEST
Severity: MEDIUM
Claim: Error handling for invalid inputs is tested.
Issue: No test covers `ErrInvalidQuantity` (qty <= 0).
Recommendation: Add test calling each deduct method with qty=0/-1.

## Gap 4
Location: tests/locking_test.go
Type: MISSING_TEST
Severity: MEDIUM
Claim: Non-existent product handling is tested.
Issue: No test covers `ErrNotFound`.
Recommendation: Add test deducting from a non-seeded product ID.

## Gap 5
Location: tests/locking_test.go
Type: MISSING_TEST
Severity: LOW
Issue: `AtomicDeduct` insufficient-stock path untested.
Recommendation: Seed stock < qty, expect `ErrInsufficientStock`.

## Gap 6
Location: engineering/03-execution-result.md
Type: UNVERIFIED_RESULT
Severity: LOW
Issue: Recorded demo output shows `OptimisticFails: 61` and `Elapsed: 49.859ms`;
actual re-run shows `45` conflicts and `23.075ms`. Non-deterministic concurrent
output recorded as if deterministic.
Recommendation: Clarify in engineering notes that conflict counts and timing are
non-deterministic; record a range or omit exact values.

## Gap 7
Location: tests/locking_test.go
Type: MISSING_EDGE_CASE
Severity: LOW
Issue: `TestAtomicConditionalUpdate` tests only happy path.
Issue: No concurrent test combines optimistic-with-retry against a scenario
where stock could hit zero (boundary exhaustion). Demo uses 100→80 (not exhausted).
Recommendation: Add test where stock is small relative to concurrent demand to
exercise `ErrOptimisticLock` after retry exhaustion.
