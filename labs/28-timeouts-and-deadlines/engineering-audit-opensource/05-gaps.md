# Gap Analysis

## Findings by Type

1. [WARNING] CIRCULAR_EXECUTE_RACE
   Location: internal/circuit/circuit.go:128-139 (`Execute`)
   Description: `Execute` is non-atomic. `Allow()` releases the lock, `fn` executes, then `RecordSuccess`/`RecordFailure` reacquire the lock. In HALF_OPEN state under concurrent calls, multiple probes can execute before any records a success/failure, allowing more than one concurrent probe beyond the intended single-probe pattern. No data race (mutex correct), but logical race.
   Test Coverage: No concurrent test exercises this path.
   Severity: MEDIUM
   Gap Type: RACE_CONDITION (logical)

2. [WARNING] IDEMPOTENCY_CHECK_THEN_ACT
   Location: internal/idempotency/idempotency.go (API design)
   Description: Store exposes `Get` and `Set` separately. Callers doing check-then-act (`if !store.Get(k) { store.Set(k, v); doWork() }`) risk double execution under concurrent same-key requests. The store is internally thread-safe (no data race) but does not provide atomic get-or-create semantics. Demo and integration test are single-goroutine so this is not exercised.
   Test Coverage: `TestStore_ConcurrentAccess` only checks value correctness under concurrency, not that work() executes exactly once.
   Severity: MEDIUM
   Gap Type: MISSING_EDGE_CASE (concurrency edge case)

3. [WARNING] GOROUTINE_LEAK_ON_CONTEXT_IGNORING_FN
   Location: internal/deadline/deadline.go:18-20
   Description: If `fn` passed to `ExecuteWithBudget` ignores `childCtx` and runs indefinitely, the goroutine leaks when the context deadline fires (only the channel send is buffered). Tests use context-aware `fn` so this is never triggered.
   Test Coverage: No test with a context-ignoring, long-running function.
   Severity: LOW
   Gap Type: MISSING_EDGE_CASE

4. [INFO] MINIMAL_STRESS_COVERAGE
   Description: No load/stress tests for retrier or deadline under high concurrency. Only idempotency has a concurrent test.
   Severity: LOW
   Gap Type: MISSING_TEST

5. [INFO] NO_FUZZ
   Description: Backoff calculation and circuit thresholds are not fuzz-tested.
   Severity: LOW
   Gap Type: MISSING_TEST
