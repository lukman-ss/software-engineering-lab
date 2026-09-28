# Code Audit

Target Lab: labs/28-timeouts-and-deadlines

## Finding 1

Location: `internal/deadline/deadline.go:13-28`
Claimed Behavior: Executes worker function with budget; terminates on context cancellation or timeout.
Observed Implementation: Uses `context.WithTimeout(ctx, budget)`. Launches worker function in a goroutine and selects on `childCtx.Done()` vs `done` channel. Channel capacity is 1, preventing goroutine blockage if context times out.
Assessment: PASS
Severity: LOW
Notes: If `fn` ignores `childCtx.Done()`, the background goroutine continues until `fn` completes. This is standard Go behavior for asynchronous execution with budget, but callers must cooperate with context.

## Finding 2

Location: `internal/retry/retry.go:35-48`
Claimed Behavior: Exponential backoff with full jitter in range `[0, min(MaxBackoff, BaseBackoff * 2^(attempt-1))]`.
Observed Implementation: Uses `1 << uint(attempt-1)`, clips to `maxVal`, multiplies by `rand.Float64()`. Falls back to defaults in `NewRetrier`.
Assessment: PASS
Severity: LOW
Notes: Correctly implements the AWS full jitter algorithm. Zero configuration defaults properly initialized.

## Finding 3

Location: `internal/circuit/circuit.go:38-139`
Claimed Behavior: Thread-safe circuit breaker with `CLOSED`, `OPEN`, and `HALF_OPEN` state transitions based on failure/success thresholds and cooldown duration.
Observed Implementation: State machine guarded by `sync.RWMutex` (`mu.Lock()` used on both reads and state updates due to lazy cooldown transition `checkCooldown()`). State transitions follow specification.
Assessment: PASS
Severity: LOW
Notes: Mutex correctly prevents race conditions during concurrent `Allow()`, `RecordSuccess()`, and `RecordFailure()` calls.

## Finding 4

Location: `internal/idempotency/idempotency.go:13-52`
Claimed Behavior: Concurrent safe in-memory deduplication store with TTL lazy-eviction.
Observed Implementation: Protected by `sync.RWMutex` (`mu.Lock()` on both `Get` and `Set` to support in-place lazy deletion of expired keys).
Assessment: PASS
Severity: LOW
Notes: Clean and safe concurrent map implementation.
