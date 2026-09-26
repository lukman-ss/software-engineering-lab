# Gap Analysis

Allowed gap types applied:

## MISSING_TEST
- Invalid quantity handling for each strategy (Test<Strategy>InvalidQuantity)
- Missing product ID handling for each strategy (Test<Strategy>NotFound)
- Concurrent oversell boundary for optimistic path: multiple goroutines draining stock from 1→0 via OptimisticDeduct (expect ErrInsufficientStock or ErrOptimisticLock, never negative stock)
- Concurrent drain-to-zero for optimistic retry path: starting stock N, M>N goroutines each decrement 1 with retry, ensure final stock == 0 and no ErrInsufficientStock lost
- Concurrent Get/Seed race: one goroutine Seeding while another calls Get (should return ErrNotFound or seeded product consistently)

## MISSING_EDGE_CASE
(Same as MISSING_TEST above — edge cases not exercised)

## RACE_CONDITION
None: `go test -race ./...` reports zero races.

## DOC_CODE_MISMATCH
None.

## BROKEN_IMPLEMENTATION
None.

## UNHANDLED_ERROR
None: all error paths propagate correctly; mutexes always unlocked.

## IMPLEMENTATION_OVERCLAIM
None: README/engineering/docs align with observed behavior.

## FAKE_DEMO / FAKE_BENCHMARK / UNVERIFIED_RESULT
None: demo output varies per run but reflects true concurrency effects (lost updates, conflict counts, retry convergence); no fabricated/invariant-violating numbers observed.

## Summary
Primary gaps are missing negative/error-path tests and missing concurrent edge-case tests for optimistic strategies at stock=0 boundary. Core safety properties verified.