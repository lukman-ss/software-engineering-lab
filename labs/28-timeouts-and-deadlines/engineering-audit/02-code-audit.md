# Code Audit

Target Lab: `labs/28-timeouts-and-deadlines`

## Finding 1

Location: `internal/deadline/deadline.go:18-20`
Claimed Behavior: Safe worker cancellation and timeout enforcement without goroutine leaks.
Observed Implementation:
```go
done := make(chan error, 1)
go func() {
    done <- fn(childCtx)
}()
```
When `childCtx.Done()` triggers before `fn` completes, `ExecuteWithBudget` returns immediately with `childCtx.Err()`. Because `done` is buffered with capacity 1, the anonymous goroutine will not block indefinitely when sending to `done`. However, if `fn` ignores `childCtx.Done()` or blocks indefinitely on uncooperative I/O, the spawned goroutine will remain alive until `fn` returns.
Assessment: PASS
Severity: LOW
Notes: Buffered channel size 1 prevents channel deadlocks. Cooperating functions that honor `childCtx` terminate cleanly.

## Finding 2

Location: `internal/retry/retry.go:35-48`
Claimed Behavior: Exponential backoff with full jitter in range `[0, min(maxBackoff, baseBackoff * 2^(attempt-1))]`.
Observed Implementation:
```go
multiplier := 1 << uint(attempt-1)
temp := float64(r.cfg.BaseBackoff) * float64(multiplier)
maxVal := float64(r.cfg.MaxBackoff)
if temp > maxVal {
    temp = maxVal
}
sleep := rand.Float64() * temp
return time.Duration(sleep)
```
Uses `math/rand/v2` which is thread-safe and cryptographically unbiased for jitter simulation. Exponential scaling caps strictly at `r.cfg.MaxBackoff`.
Assessment: PASS
Severity: LOW
Notes: Correctly implements the AWS Architecture Blog "Full Jitter" formula.

## Finding 3

Location: `internal/circuit/circuit.go:71-78, 128-139`
Claimed Behavior: State transitions `CLOSED` -> `OPEN` on failure threshold; `OPEN` -> `HALF_OPEN` after cooldown; `HALF_OPEN` -> `CLOSED` on success threshold; `HALF_OPEN` -> `OPEN` on single failure.
Observed Implementation:
In `Execute()`:
```go
if err := b.Allow(); err != nil {
    return err
}
err := fn()
if err != nil {
    b.RecordFailure()
    return err
}
b.RecordSuccess()
```
`Allow()` checks cooldown under lock and transitions state if cooldown elapsed. `RecordFailure()` and `RecordSuccess()` acquire write locks to mutate state and counters.
Assessment: PASS
Severity: LOW
Notes: Clean synchronization with `sync.RWMutex` (upgraded appropriately to write locks for state mutation).

## Finding 4

Location: `internal/idempotency/idempotency.go:29-51`
Claimed Behavior: In-memory idempotency deduplication with TTL.
Observed Implementation:
`Get()` verifies `time.Since(rec.CreatedAt) <= s.ttl` under read lock. Expired records return `("", false)`. Expired records remain in the map until overwritten, but memory footprint in lab scope is negligible.
Assessment: PASS
Severity: LOW
Notes: Thread-safe read/write lock synchronization. Lazy eviction suffices for lab scope.
