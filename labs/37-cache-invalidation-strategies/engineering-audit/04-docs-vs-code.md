# Documentation vs Code Audit

Target Lab: labs/37-cache-invalidation-strategies

## Comparison Matrix

| Claim / Topic | README.md Claim | Code Implementation | Status |
|---|---|---|---|
| Cache-Aside | Checks cache, loads from DB on miss, populates cache. Invalidates cache on update. | `CacheAsideService.Get` checks cache, queries `MockDB`, populates cache. `Update` writes to DB then deletes key from cache (`internal/cache/patterns.go:22-47`). | MATCH |
| Write-Through | Updates DB and cache synchronously on write; subsequent reads hit cache. | `WriteThroughService.Update` writes to DB then writes to cache synchronously (`internal/cache/patterns.go:78-89`). | MATCH |
| Write-Behind | Updates cache immediately; flushes to DB asynchronously via background worker queue. | `WriteBehindService.Update` writes to cache and enqueues to `writeQueue`. Background worker flushes (`internal/cache/patterns.go:120-161`). | MATCH |
| SingleFlight | Coalesces concurrent cache miss requests on a hot key down to a single DB query using `golang.org/x/sync/singleflight`. | `SingleFlightService.Get` wraps DB query in `s.flight.Do` (`internal/cache/stampede.go:56-84`). | MATCH |
| XFetch Formula | Early refresh using formula `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`. | Implemented in `ShouldRecompute` as `-deltaSec * beta * math.Log(u) > ttlRemainingSec` (`internal/cache/stampede.go:125-136`). | MATCH |
| Stale-While-Revalidate | Returns stale cached data immediately while asynchronously triggering background DB revalidation. | `SWRService.Get` checks `now.Before(staleUntil)`, calls `triggerRevalidate`, and returns stale value (`internal/cache/stampede.go:208-212`). | MATCH |
| TTL Jitter | Adds randomized offset `[0, maxJitter)` to base TTL to prevent synchronized key expiry. | `TTLWithJitter` calculates `base + time.Duration(rand.Int63n(int64(maxJitter)))` (`internal/cache/store.go:79-86`). | MATCH |
| CLI / Demo Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All commands run without modification and match output. | MATCH |

## Findings

1. **DOC_CODE_MISMATCH**: None detected.
2. **TEST_CLAIM_MISMATCH**: None detected.
3. **RESEARCH_IMPLEMENTATION_MISMATCH**: None detected. Research findings required negative sign in XFetch formula, SWR deduplication, and distinct write patterns; all are accurately implemented.
