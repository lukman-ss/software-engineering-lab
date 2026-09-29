## Finding 1

Location: internal/cache/store.go:31-43 (Get method)
Claimed Behavior: Get returns cached item if present and not expired; otherwise ErrCacheMiss.
Observed Implementation: Uses mu.RLock, checks ok, then if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) returns ErrCacheMiss.
Assessment: PASS
Severity: LOW
Notes: Correct for positive TTL. For ttl=0, ExpiresAt=zero time, IsZero() true, skips expiry check. Good.

## Finding 2

Location: internal/cache/store.go:45-51 (GetRaw method)
Claimed Behavior: GetRaw returns item regardless of expiration, needed for SWR and XFetch inspection.
Observed Implementation: Exactly returns (item, ok) without expiry check.
Assessment: PASS
Severity: LOW
Notes: Matches claim.

## Finding 3

Location: internal/cache/store.go:78-85 (TTLWithJitter)
Claimed Behavior: Adds random positive jitter [0, maxJitter) to base TTL.
Observed Implementation: base + time.Duration(rand.Int63n(int64(maxJitter))). For maxJitter>0, returns value in [base, base+maxJitter). For maxJitter<=0 returns base.
Assessment: PASS
Severity: LOW
Notes: Matches claim. ponytail comment notes std rand source; no crypto/rand needed for jitter.

## Finding 4

Location: internal/cache/repo.go:34-52 (Query method)
Claimed Behavior: Query increments queryCount, applies delay if set, returns data or ErrNotFound.
Observed Implementation: queryCount.Add(1) BEFORE delay/select. On delay>0, does time.After or ctx.Done select. Then mu.RLock, data lookup.
Assessment: PASS
Severity: LOW
Notes: queryCount increments even if ctx cancelled (minor). Otherwise correct.

## Finding 6

Location: internal/cache/patterns.go:22-37 (CacheAside Get)
Claimed Behavior: Miss→DB query→populate cache; hit→no query.
Observed Implementation: cache.Get; on hit return value; on miss db.Query(ctx), measure delta, cache.Set(key,val,ttl,delta), return val. DB errors propagated, nothing cached.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 7

Location: internal/cache/patterns.go:39-47 (CacheAside Update)
Claimed Behavior: DB write first, then cache invalidate.
Observed Implementation: db.Write first, return wrapped error on failure (cache untouched); on success cache.Delete.
Assessment: PASS
Severity: LOW
Notes: Correct ordering. Failure path preserves old cache entry; acceptable, documented pattern.

## Finding 8

Location: internal/cache/patterns.go:78-89 (WriteThrough Update)
Claimed Behavior: Synchronous DB write followed by synchronous cache update.
Observed Implementation: db.Write (wrapped error on failure, cache untouched); then cache.Set with measured delta.
Assessment: PASS
Severity: LOW
Notes: Correct. Read path (61-76) identical to CacheAside read. Good.

## Finding 9

Location: internal/cache/patterns.go:120-135 (flushWorker)
Claimed Behavior: Background worker flushes queue; Close drains remaining queue.
Observed Implementation: Loops on select queue/quit. On quit, drains len(queue) then returns. DB write errors ignored (`_ =`).
Assessment: WARNING
Severity: MEDIUM
Notes: Drain-on-close verified in code; flush errors silently swallowed so DB failure is invisible. Documented as demo scope in 02-implementation-notes; acceptable for lab but must stay a warning.

## Finding 10

Location: internal/cache/patterns.go:152-161 (WriteBehind Update)
Claimed Behavior: Immediate cache write + async enqueue; drops on full queue.
Observed Implementation: cache.Set immediately, then non-blocking send with empty default branch (silent drop).
Assessment: WARNING
Severity: MEDIUM
Notes: Matches documented decision + ponytail comment. Silent data loss on overflow is real; production would need persistent queue. No test covers overflow path (see 03-test-audit).

## Finding 11

Location: internal/cache/patterns.go:163-166 (WriteBehind Close)
Claimed Behavior: Graceful shutdown.
Observed Implementation: close(quit) + wg.Wait. Not idempotent (double Close panics); Update-after-Close enqueues with no consumer (write never flushed, no error returned).
Assessment: WARNING
Severity: LOW
Notes: Fine for demo/test lifecycle (defer Close once). No test for double-Close or post-Close Update.

## Finding 12

Location: internal/cache/stampede.go:26-41 (Naive Get)
Claimed Behavior: Naive reload causes stampede (N concurrent misses → N DB queries).
Observed Implementation: cache check, else direct db.Query + Set with no coordination.
Assessment: PASS
Severity: LOW
Notes: Correctly broken by design; demo measured 20/20.

## Finding 13

Location: internal/cache/stampede.go:56-84 (SingleFlight Get)
Claimed Behavior: Coalesces concurrent misses on one key to exactly 1 DB query.
Observed Implementation: Fast-path cache check; flight.Do(key, double-check cache → db.Query → Set). Error propagated; res.(string) only on nil error.
Assessment: PASS
Severity: LOW
Notes: Double-check inside flight correct. Known singleflight scoping: in-process only (documented limitation); shared ctx of first caller (standard, minor).

## Finding 14

Location: internal/cache/stampede.go:125-136 (ShouldRecompute)
Claimed Behavior: -Δ·β·ln(U) > TTL_remaining with correct negative sign; guards u range.
Observed Implementation: Returns false for u<=0||u>=1; computes -deltaSec*beta*math.Log(u) vs ttlRemainingSec.
Assessment: PASS
Severity: LOW
Notes: Formula sign correct; guard prevents +Inf at u=0. Verified by TestXFetchLogic + fresh test run.

## Finding 15

Location: internal/cache/stampede.go:138-169 (XFetch Get)
Claimed Behavior: Cold/expired miss always recomputes; fresh key recomputes only when formula fires; recompute failure falls back to stale.
Observed Implementation: GetRaw + now; u drawn even on cold miss; !ok||expired → true; else remaining-based ShouldRecompute. On DB error with ok → returns stale value, nil error; without ok → "", err.
Assessment: PASS
Severity: LOW
Notes: Stale-fallback masks recompute errors (returns nil error with stale data); reasonable for cache lab, untested path. Rand draw consumed on cold miss is harmless.

## Finding 16

Location: internal/cache/stampede.go:197-224 (SWR Get)
Claimed Behavior: Fresh → serve; within stale window → serve stale + async revalidate; else sync fetch.
Observed Implementation: GetRaw; IsZero||now.Before(ExpiresAt) → fresh return; staleUntil=ExpiresAt+staleDelta, now.Before → triggerRevalidate + stale return; else sync db.Query + Set.
Assessment: PASS
Severity: LOW
Notes: Boundary now==ExpiresAt treated as stale (correct). Zero-TTL items (IsZero) always fresh; consistent with store.

## Finding 17

Location: internal/cache/stampede.go:226-254 (triggerRevalidate)
Claimed Behavior: Single background revalidation per key; 10s timeout; errors ignored.
Observed Implementation: Mutex-guarded revalidating map; revalCount++ per trigger; goroutine with context timeout; Set only on success; flag cleared via defer.
Assessment: PASS
Severity: LOW
Notes: Single-revalidation-in-flight correct. Background query errors silently keep stale; acceptable, noted. revalCount exposed but never asserted in tests (see Finding 24).

## Finding 18

Location: cross-cutting (concurrency safety)
Claimed Behavior: Race-safe under high concurrency.
Observed Implementation: MemoryCache RWMutex; MockDB RWMutex + atomic counters; XFetch mu guards rand.Rand; SWR mu guards map + atomic revalCount; singleflight.Group goroutine-safe.
Assessment: PASS
Severity: LOW
Notes: go test -race -count=1 ./... fresh run PASS (1.614s). go vet clean.

## Finding 19

Location: cross-cutting (timeouts / ctx propagation)
Claimed Behavior: DB delay respects cancellation; SWR background bounded.
Observed Implementation: MockDB Query/Write select on time.After vs ctx.Done; services pass caller ctx through; SWR background uses 10s timeout.
Assessment: PASS
Severity: LOW
Notes: No test exercises ctx cancellation (see 03-test-audit).

## Finding 20

Location: cross-cutting (complexity)
Claimed Behavior: Minimal demo implementation.
Observed Implementation: 4 small files, stdlib + singleflight only, no extra abstraction layers.
Assessment: PASS
Severity: LOW
Notes: Boring and appropriately scoped. ponytail comments mark deliberate simplifications.

## Finding 5

Location: internal/cache/repo.go:54-68 (Write method)
Claimed Behavior: Write increments writeCount, applies delay, stores key/value.
Observed Implementation: writeCount.Add(1) BEFORE delay/select. On delay>0, time.After or ctx.Done select. Then mu.Lock, update map.
Assessment: PASS
Severity: LOW
Notes: Same note as Query.