# Code Audit

## Finding 1

Location: `internal/cache/store.go:79-85` (`TTLWithJitter`)
Claimed Behavior: Adds uniform random jitter in `[0, maxJitter)` to base TTL using `rand.Int63n`.
Observed Implementation: Uses package-level `math/rand` global source, which is not concurrency-safe before Go 1.20; however this code uses the global top-level functions which are safe (locked internally) since Go 1.0. Result is always `>= base` and `< base + maxJitter`.
Assessment: PASS
Severity: LOW
Notes: Using the global `math/rand` source is fine for jitter. No crypto requirement. Edge guard for `maxJitter <= 0` present. Correct range `[base, base+maxJitter)`.

---

## Finding 2

Location: `internal/cache/store.go:39-41` (`MemoryCache.Get`)
Claimed Behavior: Returns `ErrCacheMiss` for expired items.
Observed Implementation: Only checks expiry if `!item.ExpiresAt.IsZero()`. Items stored with TTL=0 skip expiry check and live forever. This is intentional by design (zero-TTL = no expiration).
Assessment: PASS
Severity: LOW
Notes: Correct sentinel handling.

---

## Finding 3

Location: `internal/cache/patterns.go:39-46` (`CacheAsideService.Update`)
Claimed Behavior: Write-then-invalidate: DB write first, then cache delete.
Observed Implementation: Correct order — DB write, then `cache.Delete(key)`. No rollback if cache.Delete fails (but Delete on a map cannot fail).
Assessment: PASS
Severity: LOW
Notes: Correct ordering avoids stale cache window that would result from delete-then-write.

---

## Finding 4

Location: `internal/cache/patterns.go:78-88` (`WriteThroughService.Update`)
Claimed Behavior: Synchronous write to DB and cache.
Observed Implementation: DB write first, then cache set. `delta` is measured as total DB write time, and the same delta is passed as `ReadDelta` when populating cache. This is a minor semantic mismatch (delta was defined as read/fetch duration in comments), but it does not affect correctness.
Assessment: PASS
Severity: LOW
Notes: Minor semantic inconsistency: `ReadDelta` field is populated with write latency in Write-Through path. Non-functional.

---

## Finding 5

Location: `internal/cache/patterns.go:120-134` (`WriteBehindService.flushWorker`)
Claimed Behavior: Background goroutine flushes writes to DB asynchronously; drains queue on graceful close.
Observed Implementation: `flushWorker` uses `select` with `quit` channel. On `quit`, drains `writeQueue` using a `len()` check loop. There is a subtle race: after `len(s.writeQueue) > 0` is evaluated, a concurrent `Update` could enqueue another item before the loop terminates, potentially causing that item to be drained in the next iteration or lost if the worker has already exited. However, since `Close()` is called by the test/user after all updates are done, this is an acceptable documentation-scoped limitation.
Assessment: WARNING
Severity: LOW
Notes: drain-on-close loop is not atomically aware of concurrent new enqueues during shutdown. Acceptable for demonstration purposes; not a correctness issue in the tested scenarios.

---

## Finding 6

Location: `internal/cache/patterns.go:152-161` (`WriteBehindService.Update`)
Claimed Behavior: Immediately writes to cache, enqueues async DB write.
Observed Implementation: Uses non-blocking `select` with `default` to drop writes when queue is full. This is documented with a comment. The drop behavior is acknowledged in the source.
Assessment: PASS
Severity: LOW
Notes: Overflow behavior is explicitly acknowledged as a demo limitation.

---

## Finding 7

Location: `internal/cache/stampede.go:62-84` (`SingleFlightService.Get`)
Claimed Behavior: Concurrent miss requests coalesce to a single DB query.
Observed Implementation: Uses `singleflight.Group.Do(key, ...)` with a double-check cache hit inside the flight function. All waiters receive the same result. Correct implementation.
Assessment: PASS
Severity: LOW
Notes: Double-check inside `Do` is good defensive practice.

---

## Finding 8

Location: `internal/cache/stampede.go:125-136` (`ShouldRecompute`)
Claimed Behavior: Implements `-Δ * β * ln(U) > TTL_remaining` formula correctly. Guard for degenerate `u <= 0` or `u >= 1`.
Observed Implementation: Guard `u <= 0 || u >= 1` returns `false` (no recompute) for boundary values. `math.Log(u)` is negative for `0 < u < 1`, so negating yields a positive value. Formula is mathematically correct.
Assessment: PASS
Severity: LOW
Notes: Edge case guard is correct. Boundary values `u=0` and `u=1` are safely handled.

---

## Finding 9

Location: `internal/cache/stampede.go:138-168` (`XFetchService.Get`)
Claimed Behavior: If no cached item or item expired, forces recompute. Otherwise evaluates probabilistic formula. On recompute failure, returns stale value as fallback if available.
Observed Implementation: Correctly distinguishes missing vs. expired vs. fresh. Fallback on DB error when raw item exists. Correct.
Assessment: PASS
Severity: LOW
Notes: Stale fallback on error is a good resilience practice.

---

## Finding 10

Location: `internal/cache/stampede.go:173-212` (`SWRService`)
Claimed Behavior: Returns stale cached data within stale window while triggering async revalidation. Uses deduplication guard to avoid concurrent revalidation goroutines for same key.
Observed Implementation: `triggerRevalidate` uses `sync.Mutex` to gate `revalidating[key]`. Background goroutine clears the key from `revalidating` map on completion. Atomic `revalCount` for inspection. Timeout of 10 seconds on revalidation context.
Assessment: PASS
Severity: LOW
Notes: Revalidation deduplication is correctly implemented.

---

## Finding 11

Location: `internal/cache/stampede.go:113-119` (`XFetchService.getRand`)
Claimed Behavior: Wraps `rand.Rand` with a mutex for concurrency safety.
Observed Implementation: `s.rand` is a non-global `*rand.Rand` (not concurrency-safe by itself). A `sync.Mutex` wraps it in `getRand`. Correct.
Assessment: PASS
Severity: LOW
Notes: Proper per-instance mutex protection.

---

## Finding 12

Location: `internal/cache/repo.go` (`MockDB.Write` latency path)
Claimed Behavior: `Write` uses `queryDelay` (the field name is shared for both read and write latency).
Observed Implementation: `Write` reuses the `queryDelay` field; there is no separate `writeDelay` field. Field name is mildly misleading but acceptable for a mock.
Assessment: PASS
Severity: LOW
Notes: Non-functional naming concern only.

---

## Summary

No HIGH or CRITICAL findings. No correctness failures. All core implementations match their claimed behavior. The one WARNING (Finding 5) is acknowledged in source comments and is a limitation specific to the demo's scoped use.
