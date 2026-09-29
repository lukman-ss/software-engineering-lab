# Code Audit

Target Lab: labs/37-cache-invalidation-strategies

## Finding 1: In-Memory Cache Store Expiration and Thread Safety

Location: `internal/cache/store.go:20-76`
Claimed Behavior: Thread-safe in-memory key-value store with TTL expiration, raw inspection for metadata, and thread-safe write/delete.
Observed Implementation:
- Uses `sync.RWMutex` protecting `items map[string]Item`.
- `Get` checks `time.Now().After(item.ExpiresAt)` under `RLock`.
- `GetRaw` allows reading raw expiration/delta metadata for SWR and XFetch calculations.
- `Set` and `Delete` acquire write lock `c.mu.Lock()`.
Assessment: PASS
Severity: LOW
Notes: Correct locking semantics. Clean and lightweight.

## Finding 2: Correctness of XFetch Probabilistic Early Expiration Formula

Location: `internal/cache/stampede.go:122-136`
Claimed Behavior: Probabilistic early refresh according to formula `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`.
Observed Implementation:
- Guards `u <= 0 || u >= 1` returning false.
- Calculates `expiryCompute := -deltaSec * beta * math.Log(u)`.
- Compares `expiryCompute > ttlRemainingSec`.
- Avoids negative log inversion bug (since `ln(u) < 0` for `u in (0, 1)`, `-delta * beta * ln(u) > 0`).
Assessment: PASS
Severity: LOW
Notes: Accurately reflects approved research correction.

## Finding 3: SingleFlight Coalescing Implementation

Location: `internal/cache/stampede.go:45-84`
Claimed Behavior: Coalesces concurrent cache miss requests on the same key into a single DB query.
Observed Implementation:
- Uses `golang.org/x/sync/singleflight.Group`.
- Implements double-checking inside `Do()` callback (`item, err := s.cache.Get(key)`) to avoid redundant work if another flight already populated the cache.
Assessment: PASS
Severity: LOW
Notes: Correctly handles concurrency and errors.

## Finding 4: Stale-While-Revalidate Deduplication and Background Execution

Location: `internal/cache/stampede.go:173-254`
Claimed Behavior: Serves stale cached values immediately while triggering single background revalidation goroutine per key.
Observed Implementation:
- Checks `now.Before(staleUntil)` to serve stale cached value immediately.
- `triggerRevalidate` guards `revalidating[key]` with `s.mu.Lock()` to prevent multiple simultaneous background goroutines for the same key.
- Revalidation runs in detached goroutine with `context.WithTimeout(context.Background(), 10*time.Second)`.
Assessment: PASS
Severity: LOW
Notes: Thread-safe deduplication prevents stampede during revalidation window.

## Finding 5: Write-Behind Queue Overflow Behavior

Location: `internal/cache/patterns.go:156-161`
Claimed Behavior: Writes immediately to cache and asynchronously to DB via buffer.
Observed Implementation:
- Uses non-blocking channel send `select { case s.writeQueue <- WriteRequest{...}: default: }`.
- When queue capacity is reached, writes to DB are dropped silently while cache remains updated.
Assessment: WARNING
Severity: LOW
Notes: Marked with a `ponytail:` comment in code and explicitly documented in `02-implementation-notes.md` as an educational design choice rather than production persistence. Documented accurately.

## Finding 6: Write-Behind Graceful Shutdown

Location: `internal/cache/patterns.go:120-135`, `internal/cache/patterns.go:163-166`
Claimed Behavior: Drain remaining queue items on shutdown.
Observed Implementation:
- `Close()` closes `quit` channel and waits on `s.wg.Wait()`.
- Worker select loop drains remaining `s.writeQueue` items when `quit` receives.
Assessment: PASS
Severity: LOW
Notes: Ensures unflushed queue items are flushed to `MockDB` during normal termination.
