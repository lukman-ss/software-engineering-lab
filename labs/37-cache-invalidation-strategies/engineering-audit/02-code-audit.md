# Code Audit

## Finding 1

Location: `internal/cache/store.go:31-43`
Claimed Behavior: Thread-safe cache read with expiration check.
Observed Implementation: Uses `sync.RWMutex.RLock()`, checks map existence and `time.Now().After(item.ExpiresAt)`.
Assessment: PASS
Severity: LOW
Notes: Correctly handles non-expiring (`ExpiresAt.IsZero()`) and expired keys.

## Finding 2

Location: `internal/cache/store.go:78-86`
Claimed Behavior: TTL jitter generation to desynchronize expiration times.
Observed Implementation: Generates random duration using `rand.Int63n(int64(maxJitter))` and adds to `base`.
Assessment: PASS
Severity: LOW
Notes: Correctly guards `maxJitter <= 0` and produces TTL within `[base, base+maxJitter)`. Uses `math/rand` pseudo-random generator, annotated with ponytail comment.

## Finding 3

Location: `internal/cache/patterns.go:39-47`
Claimed Behavior: Cache-Aside invalidates cache after DB write.
Observed Implementation: Performs synchronous DB write first, then deletes key from `MemoryCache`.
Assessment: PASS
Severity: LOW
Notes: Order of operations guarantees invalidation on DB success.

## Finding 4

Location: `internal/cache/patterns.go:78-89`
Claimed Behavior: Write-Through synchronously writes to DB and cache.
Observed Implementation: Writes to DB, measures elapsed delta, then updates `MemoryCache` with new value.
Assessment: PASS
Severity: LOW
Notes: Read-after-write will hit cache immediately.

## Finding 5

Location: `internal/cache/patterns.go:120-135`
Claimed Behavior: Write-Behind background worker flushes queued writes to DB and handles graceful shutdown.
Observed Implementation: Worker selects between `writeQueue` channel and `quit` signal; on `quit`, drains remaining buffered items in `writeQueue`.
Assessment: PASS
Severity: LOW
Notes: `Close()` closes quit channel and waits for `Wait()` group completion. Queue overflow defaults to drop (demonstration decision).

## Finding 6

Location: `internal/cache/stampede.go:56-84`
Claimed Behavior: `SingleFlightService` coalesces concurrent cache misses to 1 DB query.
Observed Implementation: Uses `golang.org/x/sync/singleflight.Group`. Calls `s.flight.Do(key, ...)` which double-checks cache inside execution function before querying DB and calling `s.cache.Set(...)`.
Assessment: PASS
Severity: LOW
Notes: Thread-safe, double-check pattern prevents redundant DB queries for concurrent goroutines.

## Finding 7

Location: `internal/cache/stampede.go:125-136`
Claimed Behavior: Probabilistic early expiration (XFetch) evaluation.
Observed Implementation: `ShouldRecompute` evaluates `expiryCompute := -deltaSec * beta * math.Log(u)` and checks `expiryCompute > ttlRemainingSec`. Guards `u <= 0 || u >= 1`.
Assessment: PASS
Severity: LOW
Notes: Formula correctly incorporates negative sign (`-Δ · β · ln(U)`) per research approved revision.

## Finding 8

Location: `internal/cache/stampede.go:226-254`
Claimed Behavior: Stale-While-Revalidate triggers async background revalidation while guarding concurrent duplicate revalidations for the same key.
Observed Implementation: `triggerRevalidate` tracks in-flight revalidations in `revalidating map[string]bool` under mutex. Spawns single background goroutine with `context.WithTimeout(10s)` and cleans up map entry on completion.
Assessment: PASS
Severity: LOW
Notes: Prevents goroutine stampede on stale cache hits. Correct atomic increment on revalidation count.
