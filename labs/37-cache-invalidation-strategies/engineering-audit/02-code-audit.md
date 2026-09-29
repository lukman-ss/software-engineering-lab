# Code Audit

## Finding 1

Location: `internal/cache/store.go:31-43` — `MemoryCache.Get`
Claimed Behavior: Returns item if present and not expired; returns `ErrCacheMiss` otherwise.
Observed Implementation: Acquires `RLock`, checks existence and expiry. Returns `ErrCacheMiss` for absent or expired items. Items with zero `ExpiresAt` are treated as non-expiring (TTL 0).
Assessment: PASS
Severity: LOW
Notes: Expired item is returned as `ErrCacheMiss` but NOT deleted on read. Stale entry persists in map until overwritten. Non-issue for correctness but causes unbounded memory growth over time if no eviction. No TTL eviction background loop — intentional per implementation notes (in-memory demo scope).

---

## Finding 2

Location: `internal/cache/store.go:79-85` — `TTLWithJitter`
Claimed Behavior: Adds random jitter in `[0, maxJitter)` to base TTL; prevents synchronized expiry.
Observed Implementation: Uses `rand.Int63n(int64(maxJitter))`. Returns `base + jitter`. Guards `maxJitter <= 0` (returns base unchanged).
Assessment: PASS
Severity: LOW
Notes: Uses package-level `math/rand` (not seeded per-call, uses default global source). Thread-safe since Go 1.20 (global source auto-seeded). Does not use cryptographic randomness — correct for this use case.

---

## Finding 3

Location: `internal/cache/patterns.go:39-47` — `CacheAsideService.Update`
Claimed Behavior: Write to DB first, then invalidate cache.
Observed Implementation: `db.Write` then `cache.Delete`. Correct write-then-invalidate sequence. No rollback if DB write succeeds but cache delete fails (impossible for in-memory delete, but conceptually noted).
Assessment: PASS
Severity: LOW
Notes: Classic Cache-Aside update pattern correctly implemented. The race window between DB write and cache delete is inherent to Cache-Aside; not a bug.

---

## Finding 4

Location: `internal/cache/patterns.go:78-89` — `WriteThroughService.Update`
Claimed Behavior: Write to DB and update cache synchronously.
Observed Implementation: Calls `db.Write`, then `cache.Set`. If `db.Write` fails, returns early without updating cache — correct. Delta measured across DB write for ReadDelta stored in cache; this is the write time, not the read latency, which is semantically inconsistent with XFetch intent (ReadDelta for XFetch should be DB query time). However, `WriteThroughService` is not used with XFetch.
Assessment: PASS
Severity: LOW
Notes: ReadDelta for write-through is the DB write latency — semantically wrong for XFetch but not a concern here because WriteThroughService is independent.

---

## Finding 5

Location: `internal/cache/patterns.go:120-135` — `WriteBehindService.flushWorker`
Claimed Behavior: Background goroutine flushes pending writes; drains on close.
Observed Implementation: `select` on `writeQueue` or `quit`. On quit signal, drains remaining items in `writeQueue` using `for len(s.writeQueue) > 0`. There is a race condition: items could be enqueued between `close(quit)` and the drain loop. However, this is a demonstration-scope implementation correctly documented with `ponytail:` comment.
Assessment: WARNING
Severity: MEDIUM
Notes: The drain loop `for len(s.writeQueue) > 0` is not fully safe — a producer calling `Update` after `Close()` initiates will panic (send on closed channel). However, in test usage `Close()` is deferred and no concurrent updates occur after the test assertion. Practically safe in test and demo, but the implementation has an acknowledged durability risk documented in notes. Not a test correctness issue.

---

## Finding 6

Location: `internal/cache/patterns.go:152-161` — `WriteBehindService.Update`
Claimed Behavior: Write to cache immediately, enqueue DB flush.
Observed Implementation: `cache.Set` with hardcoded `1*time.Millisecond` ReadDelta, then non-blocking `select` to queue. Drops silently on full queue with a comment.
Assessment: WARNING
Severity: MEDIUM
Notes: Hardcoded `1ms` ReadDelta is a placeholder, not measured. Acceptable for demo. Silent drop on queue overflow is undocumented at call site (documented in implementation notes). Production risk noted correctly.

---

## Finding 7

Location: `internal/cache/stampede.go:56-84` — `SingleFlightService.Get`
Claimed Behavior: Coalesces `N` concurrent misses to 1 DB query.
Observed Implementation: First checks cache outside flight (fast path). Inside flight function, double-checks cache (prevents redundant DB hit if another goroutine already populated). Uses `singleflight.Group.Do` per key. Correct pattern.
Assessment: PASS
Severity: LOW
Notes: Double-checked locking pattern is correct and idiomatic for singleflight. The shared `ctx` passed into the flight function is the caller's context — if the first caller's context is cancelled mid-flight, all waiting goroutines receive the error. Acceptable for this scope.

---

## Finding 8

Location: `internal/cache/stampede.go:122-136` — `ShouldRecompute`
Claimed Behavior: Implements XFetch formula `-Δ · β · ln(U) > TTL_remaining`.
Observed Implementation:
```go
expiryCompute := -deltaSec * beta * math.Log(u)
return expiryCompute > ttlRemainingSec
```
Guards `u <= 0 || u >= 1` returning false (prevents `+Inf` from `log(0)` and meaningless `log(1) = 0`).
Assessment: PASS
Severity: LOW
Notes: Formula is mathematically correct. Guard clause correctly handles degenerate inputs. Exported function allows direct unit testing of the formula — good design.

---

## Finding 9

Location: `internal/cache/stampede.go:138-169` — `XFetchService.Get`
Claimed Behavior: Returns cached value or proactively recomputes based on XFetch condition.
Observed Implementation: Gets `u` from randFunc/rand, then evaluates `ShouldRecompute`. Falls back to stale item if recompute fails (`ok` check). Correctly handles cache miss (item not present) vs. expired item.
Assessment: PASS
Severity: LOW
Notes: `u` is drawn before the recompute check — correct. One subtle point: `u` is drawn even on a hard miss (no item), which is wasteful but harmless. The stale-fallback on DB error (`if ok { return item.Value }`) is a sound resilience decision.

---

## Finding 10

Location: `internal/cache/stampede.go:197-223` — `SWRService.Get`
Claimed Behavior: Fresh → return; stale within window → return stale + async revalidate; hard miss → sync fetch.
Observed Implementation: Evaluates `item.ExpiresAt` and `staleUntil`. Triggers `triggerRevalidate` under the stale window. Falls through to synchronous fetch on full expiry or miss.
Assessment: PASS
Severity: LOW
Notes: SWR correctly implements 3-tier freshness model.

---

## Finding 11

Location: `internal/cache/stampede.go:226-253` — `SWRService.triggerRevalidate`
Claimed Behavior: Starts one background goroutine per key; prevents multiple concurrent revalidations.
Observed Implementation: Acquires mutex, checks `revalidating[key]`, returns early if already in-flight. Sets flag, releases mutex. Goroutine clears flag on completion. Uses `context.WithTimeout(10s)` for background fetch.
Assessment: PASS
Severity: LOW
Notes: Single-revalidation-per-key guard is correct. The goroutine is not tracked via a WaitGroup — SWR service has no `Close()` method, so background goroutines may outlive the service in tests. In test scope this is acceptable (goroutines complete within test duration). Production would benefit from a `Close()` with context cancellation.

---

## Finding 12

Location: `internal/cache/repo.go:34-52` — `MockDB.Query`
Claimed Behavior: Simulates DB latency, increments query count, returns value.
Observed Implementation: Atomically increments `queryCount` before the delay — counts all calls including context-cancelled ones. Correct for measuring stampede (counts all concurrent initiating queries). `RLock` used for concurrent reads.
Assessment: PASS
Severity: LOW
Notes: Atomic counter is race-safe. Context cancellation handled during delay. Correct.

---

## Finding 13

Location: `internal/cache/stampede.go:93-97` — `XFetchService` struct with `rand *rand.Rand` and `mu sync.Mutex`
Claimed Behavior: Thread-safe random number generation.
Observed Implementation: `rand.Rand` (not the global source) is used with a dedicated `mu sync.Mutex`. `getRand()` acquires the mutex. The `randFunc` injection bypasses the mutex entirely.
Assessment: WARNING
Severity: LOW
Notes: `randFunc` injection has no mutex protection. Injected `randFunc` in tests is deterministic (returns constant), so no race in tests. Concurrent real usage with injected `randFunc` would race. Minor since this is test-only injection.

---

## Finding 14

Location: `internal/cache/store.go` — Design: no eviction background loop
Claimed Behavior: Cache with TTL support.
Observed Implementation: TTL is enforced lazily (on read only). Expired entries accumulate in the map until overwritten. No background sweeper.
Assessment: WARNING
Severity: LOW
Notes: Acceptable for lab scope — explicitly noted in implementation notes as in-memory demo. Not a correctness issue for demonstrated patterns.
