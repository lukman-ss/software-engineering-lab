# Code Audit — labs/28-timeouts-and-deadlines

## Finding 1

Location: internal/deadline/deadline.go:13-27
Claimed Behavior: Operations abort immediately once deadline expires, without leakage (engineering/01-design.md Success Criteria 1).
Observed Implementation: `ExecuteWithBudget` spawns `go func(){ done <- fn(childCtx) }()` with buffered `done chan error, 1`. On `childCtx.Done()` path returns immediately, worker goroutine continues until `fn` returns. No kill, no wait, no leak check.
Assessment: WARNING
Severity: MEDIUM
Notes: Standard Go limitation — cannot kill goroutine. Return is fast, but background work leaks until `fn` respects ctx. Worker ignoring ctx leaks permanently. `select` also races when both branches ready — may return DeadlineExceeded even if `fn` just succeeded. No panic recovery.

## Finding 2

Location: internal/deadline/deadline.go:9
Claimed Behavior: N/A (exported sentinel).
Observed Implementation: `var ErrDeadlineExceeded` defined, never returned. All timeout paths return `childCtx.Err()` (`context.DeadlineExceeded` / `context.Canceled`).
Assessment: WARNING
Severity: LOW
Notes: Dead code. Tests assert `context.DeadlineExceeded`, not sentinel. Remove or wrap.

## Finding 3

Location: internal/retry/retry.go:35-48, 50-77
Claimed Behavior: Exponential backoff with full jitter, avoids retry storms.
Observed Implementation: `sleep = rand.Float64() * min(MaxBackoff, Base*2^(attempt-1))`. Correct full jitter. `Do` checks ctx before each attempt and during backoff sleep via `select ctx.Done / time.After`. Returns `errors.Join(ErrMaxRetriesExceeded, lastErr)`.
Assessment: PASS
Severity: LOW
Notes: `1 << uint(attempt-1)` overflows only for absurd MaxAttempts. Retries all errors — no retryable/non-retryable classification. Acceptable, not claimed. `math/rand/v2` auto-seeded, no determinism issue.

## Finding 4

Location: internal/circuit/circuit.go:64-139
Claimed Behavior: 3-state CLOSED/OPEN/HALF_OPEN, mutex sync, prevents cascading calls.
Observed Implementation: `sync.RWMutex` (used as Mutex — `State()` takes write Lock). `checkCooldown` under lock, OPEN→HALF_OPEN after Cooldown. CLOSED success resets failures; CLOSED failures trip at threshold; HALF_OPEN success-count→CLOSED, any failure→OPEN. `Allow` rejects only OPEN. `Execute` = Allow→fn→Record.
Assessment: PASS
Severity: LOW
Notes: Correct single-threaded transitions, verified by tests. Two limits: (1) HALF_OPEN allows unlimited concurrent probes — no single-flight cap; concurrent successes/failures interleave but mutex keeps counters safe. (2) `Execute(fn func() error)` takes no ctx — cannot propagate deadline through breaker. `State()` using Lock not RLock is safe, just coarse.

## Finding 5

Location: internal/idempotency/idempotency.go:29-52
Claimed Behavior: Thread-safe dedup store with TTL expiration.
Observed Implementation: `sync.RWMutex` + map. `Get` takes write Lock (needed for lazy delete of expired key). `Set` overwrites + refreshes `CreatedAt`. No background sweeper — expired-but-never-read keys stay in map.
Assessment: PASS
Severity: LOW
Notes: Concurrency safe. TTL refresh-on-overwrite semantics undocumented but sane for demo. Unbounded growth for write-once-never-read keys — acceptable for lab, noted in 05-gaps.

## Finding 6

Location: cmd/demo/main.go:1-97
Claimed Behavior: Demo of deadline propagation, jitter retry, circuit transitions, idempotent retry.
Observed Implementation: Real executable, imports all four packages. Runs without flags, deterministic sleeps (50ms/80ms/60ms), prints all four sections.
Assessment: PASS
Severity: LOW
Notes: Output matches engineering/03-execution-result.md byte-for-byte modulo timing. No fabricated benchmark, no network, no hidden failure.
