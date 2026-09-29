# Code Audit

## Finding 1

Location: internal/cache/store.go
Claimed Behavior: Memory cache TTL expiration and jitter.
Observed Implementation: `Get` returns ErrCacheMiss on expiry; `Set` stores ExpresAt. TTLWithJitter adds [0,maxJitter) to base.
Assessment: PASS
Severity: LOW

## Finding 2

Location: internal/cache/patterns.go CacheAsideService
Claimed Behavior: Read misses -> DB, populate cache; Update -> DB write + invalidate.
Observed Implementation: Update writes DB then `Delete` cache before returning.
Assessment: WARNING
Severity: LOW
Notes: Update ordering (write DB, then invalidate) is intentional Cache-Aside. Acceptable; brief stale inconsistency inherent.

## Finding 3

Location: internal/cache/patterns.go WriteThroughService
Claimed Behavior: Update writes DB + cache synchronously.
Observed Implementation: Writes DB then Set cache.
Assessment: PASS

## Finding 4

Location: internal/cache/patterns.go WriteBehindService
Claimed Behavior: Cache updated immediately; DB flush async by worker.
Observed Implementation: Update enqueues to buffered channel; flushWorker drains; Close closes quit and drains remaining queue.
Assessment: WARNING
Severity: MEDIUM
Notes: Queue overflow silently drops writes (`default` branch). Flush worker discards DB write errors (`_ =`). Acceptable for demo but unpropagated errors noted.

## Finding 5

Location: internal/cache/stampede.go SingleFlightService
Claimed Behavior: Coalesces concurrent misses via singleflight.
Observed Implementation: `flight.Do` with double-check inside closure.
Assessment: PASS

## Finding 6

Location: internal/cache/stampede.go ShouldRecompute / XFetchService
Claimed Behavior: -Δ β ln(U) > TTL_remaining triggers early refresh.
Observed Implementation: math.Log(u), correct sign; boundary checks u<=0 || u>=1 returns false. Falls back to stale item on recompute failure.
Assessment: PASS

## Finding 7

Location: internal/cache/stampede.go SWRService
Claimed Behavior: Returns stale data within stale window, triggers async revalidate.
Observed Implementation: `revalidating` map guards duplicate triggers; background goroutine revalidates with timeout. No mutex on `GetRaw` item access after expiry, only triggers background.
Assessment: WARNING
Severity: MEDIUM
Notes: Revalidation sets cache concurrently; `triggerRevalidate` guard ensures single trigger. Stale window logic correct.

## Finding 8

Location: internal/cache/repo.go MockDB
Claimed Behavior: Atomic query counters.
Observed Implementation: atomic.Int64 counters; RWMutex on data. Context cancellation honored.
Assessment: PASS
