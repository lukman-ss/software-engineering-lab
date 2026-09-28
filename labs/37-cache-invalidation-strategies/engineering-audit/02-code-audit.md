# Code Audit Findings

## Finding 1

Location: internal/cache/store.go:20-76
Claimed Behavior: Thread-safe in-memory cache supporting TTL, entry creation timestamps, compute delta tracking, and raw retrieval for stale evaluation.
Observed Implementation: `MemoryCache` wraps an `items map[string]Item` protected with `sync.RWMutex`. `Get` enforces TTL expiration check; `GetRaw` allows reading expired entries without deletion; `Set` and `Delete` acquire write lock.
Assessment: PASS
Severity: LOW
Notes: Clean, standard Go synchronization pattern.

## Finding 2

Location: internal/cache/repo.go:13-80
Claimed Behavior: Mock backing store simulating query latency and recording atomic write/query metrics.
Observed Implementation: `MockDB` uses `sync.RWMutex` for data map, `sync/atomic.Int64` for counters, and respects context cancellation when sleeping on simulated latency.
Assessment: PASS
Severity: LOW
Notes: Context-aware query delays avoid goroutine hangs.

## Finding 3

Location: internal/cache/patterns.go:12-47 (CacheAsideService)
Claimed Behavior: Cache-Aside reads populate cache on miss; writes update DB first then invalidate cache key.
Observed Implementation: `Get` checks cache, queries DB on miss, measures fetch duration `delta`, and saves item with measured delta. `Update` calls `db.Write` first, then deletes key from cache.
Assessment: PASS
Severity: LOW
Notes: Matches canonical Cache-Aside pattern.

## Finding 4

Location: internal/cache/patterns.go:51-89 (WriteThroughService)
Claimed Behavior: Write-Through writes to DB and cache synchronously; subsequent reads hit cache.
Observed Implementation: `Update` writes to DB, measures elapsed delta, and immediately sets the key in `MemoryCache`. `Get` reads from cache first, falls back to DB on miss.
Assessment: PASS
Severity: LOW
Notes: Correct synchronous dual-write sequencing.

## Finding 5

Location: internal/cache/patterns.go:98-166 (WriteBehindService)
Claimed Behavior: Write-Behind updates cache immediately and flushes asynchronously to DB via background worker.
Observed Implementation: `Update` updates `MemoryCache` and sends request to buffered `writeQueue` channel using non-blocking select (`default` drops overflow). `Close()` closes quit channel and drains remaining queue before terminating.
Assessment: PASS
Severity: LOW
Notes: Non-blocking channel drop on overflow is explicitly documented in implementation notes as a simplified educational design choice.

## Finding 6

Location: internal/cache/stampede.go:45-84 (SingleFlightService)
Claimed Behavior: Uses SingleFlight to collapse concurrent cache misses for the same key into a single DB query.
Observed Implementation: Uses `singleflight.Group.Do(key, fn)`. Inside the flight function, performs a double-check on `s.cache.Get(key)` before falling back to `s.db.Query`, then writes the result to cache.
Assessment: PASS
Severity: LOW
Notes: Double-check pattern prevents redundant queries if cache was populated just before flight execution.

## Finding 7

Location: internal/cache/stampede.go:88-169 (XFetchService)
Claimed Behavior: Probabilistic early expiration implementation using `-Δ · β · ln(U) > TTL_remaining`.
Observed Implementation: `ShouldRecompute` correctly implements `-deltaSec * beta * math.Log(u) > ttlRemainingSec` with bounds checks `u <= 0 || u >= 1`. Fallback to stale value provided if DB query fails on background recompute.
Assessment: PASS
Severity: LOW
Notes: Direct verification against research revision: negative sign properly included.

## Finding 8

Location: internal/cache/stampede.go:173-254 (SWRService)
Claimed Behavior: Stale-While-Revalidate returns cached stale data immediately while asynchronously triggering background DB revalidation within `staleDelta` window.
Observed Implementation: `Get` checks if item is fresh or within `item.ExpiresAt.Add(s.staleDelta)`. If stale, calls `triggerRevalidate` which uses a mutex-guarded map `revalidating` to prevent duplicate concurrent background revalidation goroutines for the same key. Background goroutine queries DB and updates cache.
Assessment: PASS
Severity: LOW
Notes: Single-flight revalidation guard prevents goroutine explosion during stale window.
