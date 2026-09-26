# Code Audit

Target Lab: labs/28-timeouts-and-deadlines

## Finding 1: Context Deadline & Budget Propagation Correctness

Location: `internal/deadline/deadline.go:13-28`
Claimed Behavior: Execute function within specified budget duration while propagating parent context cancellation.
Observed Implementation:
- Creates subcontext `context.WithTimeout(ctx, budget)`.
- Dispatches worker goroutine pushing result to buffered channel `chan error` of size 1 (preventing goroutine block on channel send).
- Uses `select` on `childCtx.Done()` vs `<-done`.
Assessment: PASS
Severity: LOW
Notes: Worker function receives `childCtx` and is expected to honor cancellation. If worker function ignores `childCtx`, the goroutine continues running in background until completion, which is standard Go context semantics and correctly documented.

## Finding 2: Full Jitter Exponential Backoff Calculation

Location: `internal/retry/retry.go:35-48`
Claimed Behavior: Exponential backoff with Full Jitter algorithm `sleep = rand.Float64() * min(MaxBackoff, BaseBackoff * 2^(attempt-1))`.
Observed Implementation:
- Validates attempt bounds (`attempt <= 0` returns 0).
- Calculates `multiplier := 1 << uint(attempt-1)` with float multiplication.
- Caps value against `r.cfg.MaxBackoff`.
- Uses Go 1.22+ `math/rand/v2.Float64()` for uniform distribution across `[0, temp]`.
- Loop in `Do` checks `ctx.Done()` before and after backoff sleep using `time.After`.
Assessment: PASS
Severity: LOW
Notes: Correct standard Full Jitter implementation per AWS Architecture / Google SRE patterns.

## Finding 3: Circuit Breaker State Machine & Concurrency Safety

Location: `internal/circuit/circuit.go:38-139`
Claimed Behavior: Thread-safe 3-state circuit breaker (`CLOSED`, `OPEN`, `HALF_OPEN`) tracking failure and success thresholds with cooldown period.
Observed Implementation:
- Protected by `sync.RWMutex` (uses Lock on state changes and State/Allow evaluations).
- `checkCooldown()` correctly transitions `OPEN` to `HALF_OPEN` when `time.Since(lastStateChg) >= cfg.Cooldown`.
- `Allow()` rejects execution with `ErrCircuitOpen` when in `OPEN` state.
- `RecordFailure()` trips breaker to `OPEN` immediately upon failure in `HALF_OPEN`, or when failures reach `FailureThreshold` in `CLOSED`.
- `RecordSuccess()` closes breaker when successes reach `SuccessThreshold` in `HALF_OPEN`, and resets failure count in `CLOSED`.
- Default values assigned when configuration properties <= 0.
Assessment: PASS
Severity: LOW
Notes: Clean, deterministic, thread-safe implementation.

## Finding 4: In-Memory Idempotency Store

Location: `internal/idempotency/idempotency.go:13-51`
Claimed Behavior: In-memory thread-safe store for deduplicating requests with TTL expiration.
Observed Implementation:
- `Store` holds `map[string]Record` protected by `sync.RWMutex`.
- `Get()` acquires `s.mu.RLock()`, checks existence and TTL expiration (`time.Since(rec.CreatedAt) > s.ttl`).
- `Set()` acquires `s.mu.Lock()`, writes record with current timestamp.
- Default TTL fallback to 1 minute if <= 0.
Assessment: PASS
Severity: LOW
Notes: As noted in `engineering/02-implementation-notes.md`, this is an in-memory reference implementation intended for single-node demonstration and unit verification without background purging overhead.
