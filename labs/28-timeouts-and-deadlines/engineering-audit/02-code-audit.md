# Code Audit

Target Lab: labs/28-timeouts-and-deadlines

## Finding 1

Location: `internal/deadline/deadline.go:13-28`
Claimed Behavior: Context deadline execution budget and timeout inheritance.
Observed Implementation: `ExecuteWithBudget` creates a child context with `context.WithTimeout(ctx, budget)`. Spawns worker in goroutine writing to buffered channel of size 1. Selects on `childCtx.Done()` vs `done`.
Assessment: PASS
Severity: LOW
Notes: If `fn(childCtx)` does not honor `childCtx.Done()`, the spawned goroutine will remain running until `fn` returns. This is standard Go context semantics; callers must monitor `ctx.Done()`. Channel buffer of 1 prevents goroutine leaks on completion after timeout.

## Finding 2

Location: `internal/retry/retry.go:35-48`
Claimed Behavior: Exponential backoff with Full Jitter: $sleep = \text{rand}(0, \min(M, B \cdot 2^{attempt-1}))$.
Observed Implementation: Uses `math/rand/v2`, calculates `temp = float64(r.cfg.BaseBackoff) * float64(1 << uint(attempt-1))`, clamps to `MaxBackoff`, and returns `rand.Float64() * temp`.
Assessment: PASS
Severity: LOW
Notes: Correctly checks context cancellation before attempts and during backoff wait.

## Finding 3

Location: `internal/circuit/circuit.go:64-126`
Claimed Behavior: Circuit breaker state transitions: CLOSED -> OPEN on failure threshold; OPEN -> HALF_OPEN after cooldown; HALF_OPEN -> CLOSED on success threshold; HALF_OPEN -> OPEN on single failure.
Observed Implementation: State transitions guarded by `sync.RWMutex` (`mu.Lock()` on all state mutating paths and checks). `checkCooldown` lazily transitions OPEN to HALF_OPEN when cooldown duration elapses.
Assessment: PASS
Severity: LOW
Notes: `State()` uses `mu.Lock()` to allow lazy state updates on cooldown expiry, avoiding race conditions.

## Finding 4

Location: `internal/idempotency/idempotency.go:29-52`
Claimed Behavior: Thread-safe idempotency response store with lazy TTL eviction.
Observed Implementation: Guarded by `sync.RWMutex`. `Get` takes `mu.Lock()` to safely delete expired records lazily. `Set` stores record with `time.Now()`.
Assessment: PASS
Severity: LOW
Notes: Eviction is lazy on key access. Memory usage scales with unique keys until accessed post-TTL. Acceptable for in-memory lab demo scope.
