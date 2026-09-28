# Engineering Gaps Analysis

Target Lab: `labs/28-timeouts-and-deadlines`

## Gap Summary

No critical, high, or medium severity gaps were identified.

## Minor Observations (Low Severity)

1. **Passive Idempotency Store Eviction**:
   - Gap Type: `MISSING_EDGE_CASE` (Low)
   - Description: Keys in `idempotency.Store` are only invalidated when queried via `Get()`. Unqueried expired keys remain in the map.
   - Impact: Acceptable for in-memory educational lab demo.

2. **Caller Goroutine Lifecycle Responsibility in `deadline.ExecuteWithBudget`**:
   - Gap Type: `MISSING_EDGE_CASE` (Low)
   - Description: `ExecuteWithBudget` unblocks immediately upon timeout via channel select, but if `fn` does not monitor `ctx.Done()`, `fn` continues executing in background.
   - Impact: Standard Go idiom. Calling code must observe context cancellation.
