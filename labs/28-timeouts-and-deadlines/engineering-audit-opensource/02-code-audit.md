# Code Audit

## Finding 1: Deadline Context Propagation

Location: `internal/deadline/deadline.go:13-28` (ExecuteWithBudget)

Claimed Behavior: Sub-context inherits parent deadline; worker receives budgeted context; parent timeout overrides local budget.

Observed Implementation:
```go
func ExecuteWithBudget(ctx context.Context, budget time.Duration, fn WorkerFunc) error {
    childCtx, cancel := context.WithTimeout(ctx, budget)
    defer cancel()

    done := make(chan error, 1)
    go func() {
        done <- fn(childCtx)
    }()

    select {
    case <-childCtx.Done():
        return childCtx.Err()
    case err := <-done:
        return err
    }
}
```

Assessment: PASS

Severity: LOW (no issue)

Notes:
- Uses `context.WithTimeout` correctly to bound execution.
- Worker receives `childCtx` so any internal deadline checks propagate correctly.
- Select case order ensures context cancellation wins over worker result if child context expires first.
- The worker may run longer than budget if it ignores `childCtx` but that's caller responsibility; function contracts to return context error when childCtx expires.

## Finding 2: Retry Jitter Implementation

Location: `internal/retry/retry.go:35-48` (CalculateBackoff) and `retry.go:50-76` (Do)

Claimed Behavior: Exponential backoff with full jitter (sleep = random[0, min(base*2^(attempt-1), max)]).

Observed Implementation:
```go
func (r *Retrier) CalculateBackoff(attempt int) time.Duration {
    multiplier := 1 << uint(attempt-1)
    temp := float64(r.cfg.BaseBackoff) * float64(multiplier)
    maxVal := float64(r.cfg.MaxBackoff)
    if temp > maxVal {
        temp = maxVal
    }
    // Full Jitter: random duration in [0, temp]
    sleep := rand.Float64() * temp
    return time.Duration(sleep)
}
```

Assessment: PASS

Severity: LOW

Notes:
- Matches full jitter formula.
- `rand.Float64()` returns [0,1), so sleep in [0, temp) inclusive-exclusive upper bound.
- Uses `math/rand/v2` which is Go 1.22+; acceptable per go.mod.
- Backoff calculation capped at MaxBackoff correctly.
- In `Do`, backoff applied between attempts (after failure, before next try), correct.

## Finding 3: Circuit Breaker Half-Open Single Trial

Location: `internal/circuit/circuit.go:64-89` (State, Allow, checkCooldown, RecordSuccess, RecordFailure)

Claimed Behavior: Circuit breaker allows one trial in HALF_OPEN state; if success, counts toward success threshold; if fail, trips open immediately.

Observed Implementation:
- `State()` calls `checkCooldown()` which moves OPEN→HALF_OPEN after cooldown.
- `Allow()` returns nil for CLOSED and HALF_OPEN, error for OPEN.
- `RecordSuccess()` increments successes only in HALF_OPEN; if successes >= threshold, transitions to CLOSED.
- `RecordFailure()` transitions HALF_OPEN→OPEN immediately on any failure.
- `Execute()` calls `Allow()` then `fn()` then records result.

Assessment: PASS

Severity: LOW

Notes:
- No explicit limit on concurrent trials in HALF_OPEN: multiple goroutines could pass `Allow()` before state updates.
  However, the mutex inside each method serializes state transitions; but consider:
    - G1 calls Allow -> sees HALF_OPEN -> returns nil
    - G2 calls Allow before G1 records result -> also sees HALF_OPEN -> returns nil
    Both proceed to fn() concurrently.
  This is a known limitation; typical circuit breaker implementations often allow a single trial.
  Claim does not explicitly state "single trial", only state machine. The lab's claim: "State machine (CLOSED, OPEN, HALF_OPEN) preventing cascading calls". Partial mitigation: success/failure thresholds require multiple events.
  For single failure in HALF_OPEN, it trips open immediately, limiting damage.
  Mark as WARNING.

## Finding 4: Idempotency Store TTL Deletion

Location: `internal/idempotency/idempotency.go:29-51` (Get, Set)

Claimed Behavior: In-memory deduplication store with TTL; entries expire after TTL to prevent unbounded growth.

Observed Implementation:
- `Get` checks `time.Since(rec.CreatedAt) > s.ttl`; if expired, returns not-found but does NOT delete entry.
- `Set` unconditionally overwrites.
- No background janitor; expired entries linger until overwritten.

Assessment: WARNING

Severity: MEDIUM

Notes:
- Expired entries accumulate until same key reused (Set) or process restart.
  Under high-cardinality keys (e.g., per-request UUIDs), memory leak.
  Claim: "In-memory deduplication store preventing double execution during retries." Does not specify persistence or cleanup bounds.
  Lab's limitation noted in engineering/02-implementation-notes.md: "The in-memory idempotency store does not persist across application crashes." Does not mention leak.
  Fix: delete on expired Get or add periodic cleanup.
  Since lab is self-contained demo, risk low but still a resource leak.

## Finding 5: Idempotency Store Race on Get-Then-Set

Location: `internal/idempotency/idempotency.go:29-51`

Claimed Behavior: Store is safe for concurrent use.

Observed Implementation:
- Get uses RLock, Set uses Lock.
- However, typical pattern: check Get, if miss then Set; between Get and Set another goroutine could Set then first goroutine sets again (overwrite).
  This is benign for idempotency (same key, possibly different response) but could cause wasted work if response generation expensive.
  Claim: "Idempotent request handlers deduplicate retried operations using idempotency keys." Implies exactly-one storage per key; current allows multiple sets.
  However, the demo's usage wraps Get-then-Set in user function; lab does not provide atomic GetOrSet.
  Mark as WARNING.

Severity: MEDIUM

Notes:
- In `cmd/demo/main.go`, the payment function does:
    if res, ok := store.Get(key); ok { ... } else { process; store.Set }
  This is classic check-then-set race.
  Under concurrent retries for same key, multiple payment processes may occur before first sets.
  However, idempotency key is meant to be stable per logical request; if retries happen rapidly, window small.
  Still, a race exists.

## Finding 6: Demo Output Matches Implementation

Location: `cmd/demo/main.go`

Claimed Behavior: Demo shows:
  1. Deadline propagation: context deadline exceeded
  2. Retry: 3 attempts then success
  3. Circuit breaker: transitions as described
  4. Idempotency: second request deduplicated

Observed Implementation: Demo code matches described behavior.

Assessment: PASS

Severity: LOW

Notes:
- Ran demo; output identical to engineering/03-execution-result.md.
- No mocking; real execution.

## Finding 7: Test Coverage Gaps

Location: Various test files

Claimed Behavior: Unit tests verify happy path, failure path, edge cases.

Observed Implementation:
- deadline_test.go: success, timeout, parent timeout inheritance (good)
- retry_test.go: success first try, retry until success, exceed max attempts, context canceled (good)
- circuit_test.go: only one test covering state transitions Closed->Open->HalfOpen->Closed; missing:
    - Open state rejects immediately (covered)
    - Half-open failure -> Open (not tested)
    - Concurrent callers (not tested)
    - Success threshold >1 (test uses 2)
- idempotency_test.go: Get/Set/TLL, concurrent access (no assertions on concurrent result, just waits)
- integration_test.go: retry+circuit (all attempts fail -> open), idempotent retry (deduplication works)
  Missing: retry+idempotency+circuit interactions, success after circuit half-open, etc.

Assessment: WARNING

Severity: MEDIUM

Notes:
- Tests pass but do not exhaustively verify all claimed behaviors under edge conditions.
- Particularly circuit breaker half-open failure path missing.
- Idempotency test concurrent access does not validate correctness (no assertion on number of sets/gets results).

## Finding 8: Race Detector Clean

Location: All packages

Claimed Behavior: No data races.

Observed Implementation: `go test -race ./...` passes.

Assessment: PASS

Severity: LOW

Notes:
- No race reported.

## Finding 9: Context Cancellation Propagation in Retry

Location: `internal/retry/retry.go:50-76` (Do)

Claimed Behavior: Retry loop respects parent context cancellation.

Observed Implementation:
```go
func (r *Retrier) Do(ctx context.Context, fn func(ctx context.Context) error) error {
    var lastErr error
    for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        err := fn(ctx)
        // ...
        // backoff select also checks ctx.Done()
    }
    return errors.Join(ErrMaxRetriesExceeded, lastErr)
}
```

Assessment: PASS

Severity: LOW

Notes:
- Checks ctx.Done() at start of each iteration and before backoff.
- Passes ctx to fn so fn can also respect context.
- Correct.

## Finding 10: Timeout Budget Overrun If Worker Ignores Ctx

Location: `internal/deadline/deadline.go:13-28`

Claimed Behavior: Function returns when budget exceeded.

Observed Implementation: If worker ignores childCtx and does not check for cancellation, it may run past budget; function will still return ctx.Err() when childCtx.Done() fires (select case). However, worker goroutine continues; it is leaked until function returns or context done.

Assessment: WARNING

Severity: MEDIUM

Notes:
- The function cannot force-terminate worker; relies on worker cooperating.
- Claim: "operations abort immediately once context deadline expires." Only true if worker checks ctx.
- Mitigation: document that fn must respect childCtx.
- Lab's demo uses `select { case <-time.After: ... case <-childCtx.Done(): }` so cooperating.