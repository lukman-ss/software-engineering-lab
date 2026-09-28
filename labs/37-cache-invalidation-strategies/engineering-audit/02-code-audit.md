# Code Audit Findings

## Finding 1

Location: internal/cache/stampede.go:125-136
Claimed Behavior: Probabilistic early refresh follows `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`.
Observed Implementation: `ShouldRecompute` explicitly computes `expiryCompute := -deltaSec * beta * math.Log(u)` and checks `expiryCompute > ttlRemainingSec`, with input bounds checking `u <= 0 || u >= 1`.
Assessment: PASS
Severity: LOW
Notes: Mathematical formula accurately implements the research correction, avoiding negative ln evaluation without the minus sign.

## Finding 2

Location: internal/cache/stampede.go:62-79
Claimed Behavior: Singleflight coalesces simultaneous cache misses down to 1 DB query per key.
Observed Implementation: Uses `golang.org/x/sync/singleflight.Group.Do(key, func() (interface{}, error) {...})` with double-check inside flight execution.
Assessment: PASS
Severity: LOW
Notes: Safe, clean, thread-safe deduplication per key within process boundary.

## Finding 3

Location: internal/cache/patterns.go:120-135, 163-166
Claimed Behavior: Write-Behind background worker drains queue and shuts down safely.
Observed Implementation: Worker selects between `writeQueue` and `quit` channel. When `quit` fires, it drains `len(s.writeQueue) > 0` before returning. `Close()` closes `quit` and calls `s.wg.Wait()`.
Assessment: PASS
Severity: LOW
Notes: Safe shutdown sequence; no orphaned background goroutines.

## Finding 4

Location: internal/cache/patterns.go:156-160
Claimed Behavior: Non-blocking enqueue on full write queue.
Observed Implementation: Uses `select` with `default: // drop or handle overflow`.
Assessment: PASS
Severity: LOW
Notes: Explicitly documented as an intentional simplification in implementation notes. Appropriate for lab scope.

## Finding 5

Location: internal/cache/stampede.go:226-253
Claimed Behavior: SWR triggers single background revalidation per key without spawning redundant concurrent revalidations.
Observed Implementation: Protected by `s.mu.Lock()`, checks `revalidating[key]`. If already in flight, exits early. Spawns goroutine with 10s context timeout and cleans up key upon completion in defer.
Assessment: PASS
Severity: LOW
Notes: Properly avoids goroutine explosion under continuous stale read load.

## Finding 6

Location: internal/cache/store.go:20-76
Claimed Behavior: Thread-safe in-memory cache operations with TTL expiration.
Observed Implementation: `RWMutex` correctly guards read (`RLock`/`RUnlock`) and write (`Lock`/`Unlock`) paths across `Get`, `GetRaw`, `Set`, and `Delete`.
Assessment: PASS
Severity: LOW
Notes: No race conditions detected under race detector execution.
