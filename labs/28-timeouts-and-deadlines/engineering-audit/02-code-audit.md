# Code Audit

Target Lab: labs/28-timeouts-and-deadlines

## Finding 1

Location: internal/deadline/deadline.go:18-20
Claimed Behavior: Context cancellation terminates long-running downstream work without leakage.
Observed Implementation:
```go
done := make(chan error, 1)
go func() {
    done <- fn(childCtx)
}()
```
The channel is buffered with size 1 (`make(chan error, 1)`), preventing goroutine leak if `childCtx` cancels before `fn` completes, provided `fn` respects `childCtx` and terminates.
Assessment: PASS
Severity: LOW
Notes: `fn` must honor `childCtx.Done()`. The buffered channel prevents the launched goroutine from hanging on channel send upon cancellation.

## Finding 2

Location: internal/retry/retry.go:35-48
Claimed Behavior: Exponential backoff with full jitter in range `[0, min(MaxBackoff, BaseBackoff * 2^(attempt-1))]`.
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
Standard full jitter implementation using `math/rand/v2`.
Assessment: PASS
Severity: LOW
Notes: Correctly handles shift bound, max backoff capping, and zero attempt base cases.

## Finding 3

Location: internal/circuit/circuit.go:38-139
Claimed Behavior: Thread-safe 3-state circuit breaker (`CLOSED`, `OPEN`, `HALF_OPEN`) with cooldown and threshold resets.
Observed Implementation: All state transitions (`checkCooldown`, `Allow`, `RecordSuccess`, `RecordFailure`, `State`) are synchronized using `sync.RWMutex` (upgraded to exclusive lock `mu.Lock()` on mutating checks).
Assessment: PASS
Severity: LOW
Notes: Thread-safety verified under race detector. Transition rules accurately reflect 3-state circuit breaker pattern.

## Finding 4

Location: internal/idempotency/idempotency.go:13-51
Claimed Behavior: Thread-safe in-memory key-response deduplication store with TTL.
Observed Implementation:
- `Get` uses `s.mu.RLock()` / `s.mu.RUnlock()`.
- `Set` uses `s.mu.Lock()` / `s.mu.Unlock()`.
- Expiry evaluated via `time.Since(rec.CreatedAt) > s.ttl`.
Assessment: PASS
Severity: LOW
Notes: Clean, minimal thread-safe dictionary implementation.
