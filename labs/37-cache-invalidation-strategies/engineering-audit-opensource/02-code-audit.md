# Code Audit

Target Lab: labs/37-cache-invalidation-strategies

## Finding 1: In-Memory Cache Thread Safety and TTL Handling

Location: `internal/cache/store.go:20-76`
Claimed Behavior: Thread-safe cache operations supporting TTL checks, raw extraction for SWR/XFetch, and key deletions.
Observed Implementation: `MemoryCache` protects `items` map using `sync.RWMutex`. `Get` uses `RLock`, correctly checks `time.Now().After(item.ExpiresAt)`. `GetRaw` uses `RLock` to retrieve unexpired or expired item metadata without filtering by TTL. `Set` and `Delete` use write `Lock`.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safety design using standard library primitives.

## Finding 2: TTL Jitter Range and Uniform Randomness

Location: `internal/cache/store.go:78-86`
Claimed Behavior: Prevents synchronized key expiration by adding a non-negative random jitter up to `maxJitter`.
Observed Implementation: `TTLWithJitter` checks `if maxJitter <= 0` and returns `base`. For positive jitter, computes `rand.Int63n(int64(maxJitter))` and returns `base + jitter`.
Assessment: PASS
Severity: LOW
Notes: Correctly guarded against non-positive jitter values. Annotated with `ponytail:` comment noting stdlib `math/rand` selection.

## Finding 3: Cache Invalidation Patterns (Aside, Write-Through, Write-Behind)

Location: `internal/cache/patterns.go:10-166`
Claimed Behavior:
- Cache-Aside: Reads check cache -> fallback to DB -> populate cache. Updates write DB -> delete key from cache.
- Write-Through: Updates write DB -> synchronously update cache.
- Write-Behind: Updates write cache -> enqueue to background buffer -> flush to DB asynchronously. Graceful shutdown drains queue.
Observed Implementation:
- `CacheAsideService.Get` and `Update` match pattern precisely; errors in DB write halt before cache deletion.
- `WriteThroughService.Update` writes to DB first; returns error immediately if DB write fails, preserving DB integrity before cache update.
- `WriteBehindService` initializes a buffered channel worker with `sync.WaitGroup` and `quit` channel. Graceful `Close()` signals quit and drains queued writes. Buffer overflow drops write via non-blocking select, appropriately documented.
Assessment: PASS
Severity: LOW
Notes: Behavior conforms to classic write lifecycle models.

## Finding 4: SingleFlight Stampede Coalescing

Location: `internal/cache/stampede.go:43-85`
Claimed Behavior: Coalesces concurrent read misses on a hot key to a single DB query using `golang.org/x/sync/singleflight`.
Observed Implementation: `SingleFlightService.Get` first checks `cache.Get(key)`. On miss, delegates to `s.flight.Do(key, func() ...)` with double-check inside closure. Query results are populated back into cache and returned to all waiting callers.
Assessment: PASS
Severity: LOW
Notes: Double-check locking inside the flight execution prevents redundant fetches if cache was populated between caller check and execution.

## Finding 5: XFetch Probabilistic Early Expiration

Location: `internal/cache/stampede.go:86-170`
Claimed Behavior: Evaluates `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`.
Observed Implementation:
- `ShouldRecompute` explicitly guards `u <= 0 || u >= 1`.
- Correctly computes `expiryCompute := -deltaSec * beta * math.Log(u)`.
- If recompute succeeds, updates cache with fresh value and new duration.
- On query error during proactive recompute, gracefully returns existing stale value if available (`if ok { return item.Value, nil }`).
Assessment: PASS
Severity: LOW
Notes: Fully aligns with research correction removing false positive formula negation bugs.

## Finding 6: Stale-While-Revalidate (SWR) Concurrency and In-Flight Deduplication

Location: `internal/cache/stampede.go:171-254`
Claimed Behavior: Serves stale cached values immediately while asynchronously revalidating in the background, without launching unbounded duplicate revalidation goroutines for the same key.
Observed Implementation:
- `SWRService.Get` checks `item.ExpiresAt`. If stale but within `staleUntil = item.ExpiresAt.Add(s.staleDelta)`, returns `item.Value` immediately and calls `triggerRevalidate(key)`.
- `triggerRevalidate` guards execution using `revalidating map[string]bool` under mutex lock. Duplicate triggers return immediately.
- Background goroutine executes with `context.WithTimeout(context.Background(), 10*time.Second)` and cleans up the key in `revalidating` upon exit.
Assessment: PASS
Severity: LOW
Notes: Effective deduplication prevents background goroutine leaks during repeated reads of stale keys.
