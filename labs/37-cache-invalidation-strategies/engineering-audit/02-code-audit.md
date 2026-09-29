# Code Audit

Target Lab: labs/37-cache-invalidation-strategies

## Finding 1

Location: internal/cache/store.go:20-43 (`MemoryCache.Get`, `MemoryCache.Set`)
Claimed Behavior: Thread-safe in-memory cache with TTL expiration.
Observed Implementation: `sync.RWMutex` used correctly — `RLock` for reads, `Lock` for writes. Expiry check on `Get` returns `ErrCacheMiss` for expired items. TTL=0 treated as no expiration (non-expiring items).
Assessment: PASS
Severity: LOW
Notes: Correct mutex discipline throughout all cache operations.

## Finding 2

Location: internal/cache/store.go:46-51 (`GetRaw`)
Claimed Behavior: Returns item regardless of expiration, needed for SWR and XFetch.
Observed Implementation: Holds `RLock`, returns item and boolean presence flag without expiration check. Used correctly by XFetch and SWR services.
Assessment: PASS
Severity: LOW
Notes: Deliberate design decision, well-documented in implementation notes.

## Finding 3

Location: internal/cache/store.go:79-85 (`TTLWithJitter`)
Claimed Behavior: Adds random positive jitter `[0, maxJitter)` to base TTL.
Observed Implementation: `rand.Int63n(int64(maxJitter))` produces `[0, maxJitter)`. Guard for `maxJitter <= 0` returns base TTL unchanged. Uses `math/rand` standard package (not crypto/rand); marked with `ponytail:` comment acknowledging limitation.
Assessment: PASS
Severity: LOW
Notes: Non-cryptographic randomness is appropriate for TTL jitter. Ponytail comment correctly identifies upgrade path.

## Finding 4

Location: internal/cache/patterns.go:22-47 (`CacheAsideService.Get`, `CacheAsideService.Update`)
Claimed Behavior: Cache-Aside: read from cache → miss → DB query → populate cache. Update: DB write first, then cache invalidation.
Observed Implementation: Get: cache hit returns early, miss goes to DB, measures `ReadDelta`, sets cache with TTL. Update: DB write first, then `cache.Delete`. Correct "write-DB-then-invalidate" ordering, which avoids serving stale data after write.
Assessment: PASS
Severity: LOW
Notes: No race window between DB write and cache invalidation in sequential operation. Correct pattern implementation.

## Finding 5

Location: internal/cache/patterns.go:78-89 (`WriteThroughService.Update`)
Claimed Behavior: Synchronous write to DB and cache.
Observed Implementation: DB write first, then cache set with same value. Subsequent reads hit updated cache without DB query. `ReadDelta` is measured from DB write duration and stored — technically this is the write duration, not the read/compute duration for the item, but it is acceptable here since Write-Through does not use XFetch-style probabilistic logic.
Assessment: PASS
Severity: LOW
Notes: Minor semantic note: `ReadDelta` in `WriteThroughService` stores write latency, not fetch latency. Not a correctness issue since Write-Through doesn't use XFetch.

## Finding 6

Location: internal/cache/patterns.go:98-166 (`WriteBehindService`)
Claimed Behavior: Immediate cache update, async DB flush via buffered channel worker.
Observed Implementation: `Update` writes to cache immediately, sends to channel. `flushWorker` goroutine processes channel in background. `Close()` signals quit channel, then drains remaining items before returning. Buffer overflow handled with `select { default: }` (silently drops), documented with ponytail comment.
Assessment: PASS
Severity: LOW
Notes: Graceful shutdown drain loop (`for len(s.writeQueue) > 0`) is non-atomic with the channel close — in the window between channel close signal and drain loop, items could be added by concurrent Update calls in tests. However for test scenarios and demo, this is safe since Update is not called concurrently with Close.

## Finding 7

Location: internal/cache/stampede.go:16-41 (`NaiveStampedeService`)
Claimed Behavior: Intentionally naive — concurrent cache misses all query DB simultaneously.
Observed Implementation: No coordination — all goroutines that miss cache call `db.Query` independently. Correctly demonstrates the thundering herd problem.
Assessment: PASS
Severity: LOW
Notes: Correct — stampede behavior is intentional.

## Finding 8

Location: internal/cache/stampede.go:45-84 (`SingleFlightService.Get`)
Claimed Behavior: Coalesces concurrent miss requests to single DB query using `singleflight.Group`.
Observed Implementation: Uses `singleflight.Group.Do(key, ...)`. Double-check inside the flight function (re-checks cache before DB query). Correctly shares result across all concurrent callers. All 20 concurrent goroutines get the correct value from a single DB query.
Assessment: PASS
Severity: LOW
Notes: Double-check inside `Do` is a good defensive pattern to handle the case where a prior flight already populated the cache while the current flight was waiting.

## Finding 9

Location: internal/cache/stampede.go:122-136 (`ShouldRecompute`)
Claimed Behavior: Implements XFetch formula `-Δ · β · ln(U) > TTL_remaining`. Correct sign: `math.Log(u)` is negative for `u ∈ (0,1)`, so negating produces a positive value.
Observed Implementation: `expiryCompute := -deltaSec * beta * math.Log(u)`. Guards `u <= 0 || u >= 1` to prevent `+Inf` or zero. Returns true when `expiryCompute > ttlRemainingSec`. This is mathematically correct.
Assessment: PASS
Severity: LOW
Notes: Formula matches the research-approved correction. Guard on u bounds is correct. Test `TestXFetchLogic` validates the formula numerically with concrete u values.

## Finding 10

Location: internal/cache/stampede.go:138-169 (`XFetchService.Get`)
Claimed Behavior: Probabilistic early expiration — triggers proactive DB fetch before TTL expiry when formula holds.
Observed Implementation: Calls `GetRaw`, checks if expired (forced recompute) or `ShouldRecompute`. If no recompute needed, returns cached value. If fetch fails but stale item exists, returns stale as fallback (resilient behavior). Injects `randFunc` for deterministic testing.
Assessment: PASS
Severity: LOW
Notes: Stale fallback on XFetch recompute failure is a valid, resilient behavior. `SetRandFunc` properly enables deterministic tests.

## Finding 11

Location: internal/cache/stampede.go:173-253 (`SWRService`)
Claimed Behavior: Fresh items returned directly. Expired items within stale window served stale + trigger async background revalidation. Items beyond stale window trigger synchronous fetch.
Observed Implementation: Three-tier logic (fresh / stale-window / hard-miss) correctly implemented. `triggerRevalidate` uses `revalidating` map under mutex to prevent concurrent background revalidations per key. Background goroutine acquires 10-second timeout context, queries DB, updates cache. `revalCount` atomic counter tracks revalidations.
Assessment: PASS
Severity: LOW
Notes: The `revalidating` map cleanup (`delete`) happens in a `defer` in the goroutine, ensuring the key is unlocked even if the DB query fails. Correct.

## Finding 12

Location: internal/cache/stampede.go:237-253 (SWR background goroutine)
Claimed Behavior: Async revalidation updates cache on success; no-op on error.
Observed Implementation: Errors from `db.Query` are silently ignored (no error channel, no log). For a lab/demo context this is acceptable; production would require monitoring.
Assessment: PASS
Severity: LOW
Notes: Acceptable simplification for lab scope; documented in implementation notes.

## Finding 13

Location: internal/cache/repo.go:54-68 (`MockDB.Write`)
Claimed Behavior: Write increments `writeCount`, respects context cancellation, simulates DB delay.
Observed Implementation: `writeCount.Add(1)` called before delay. This means even cancelled writes increment the counter. In practice, tests do not cancel mid-write, so no test impact.
Assessment: WARNING
Severity: LOW
Notes: If a context is cancelled mid-write before `db.data[key] = value`, the write count is incremented but the data is not actually written. This creates a minor discrepancy in write counter semantics. Does not affect any current test correctness.

## Finding 14

Location: engineering/01-design.md:46 (mentions `jitter.go`)
Claimed Behavior: Architecture mentions `jitter.go` as a separate file.
Observed Implementation: `TTLWithJitter` is implemented in `store.go`, not a separate `jitter.go`. No `jitter.go` file exists.
Assessment: WARNING
Severity: LOW
Notes: Design doc lists `jitter.go` as an architecture component that does not exist. Implementation consolidated TTL jitter into `store.go`. Minor doc inaccuracy; code is correct.

## Finding 15

Location: tests/cache_test.go (missing test coverage)
Claimed Behavior: Implementation notes assert SWR timing-sensitive test. Write-Behind drop behavior exists but overflow test only checks cache value, not drop count.
Observed Implementation: `TestWriteBehindService_QueueOverflow` verifies cache reflects latest update but does not assert how many writes were actually flushed to DB (i.e., overflow drop is not counted). `TestCachePatterns_FailurePaths` only covers Cache-Aside and Write-Through read errors; Write-Behind DB failure path is not tested.
Assessment: WARNING
Severity: MEDIUM
Notes: Queue overflow drop behavior is tested partially (cache correct), but the drop side-effect is not quantified/asserted.

## Finding 16

Location: cmd/demo/main.go (Write-Through demo, line 59)
Claimed Behavior: Demo shows `DB queries: 1` after Write-Through Update + GET (no extra query).
Observed Implementation: Write-Through `Update` calls `db.Write` (increments `writeCount` not `queryCount`). Post-update GET hits cache directly. `db.QueryCount()` stays at 1 (from initial miss). Output is correct and not misleading.
Assessment: PASS
Severity: LOW
Notes: Demo output correctly disambiguates query vs write counts.
