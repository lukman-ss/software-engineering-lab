# Documentation vs Code Audit

## Consistency Review

| Document Item | Code Implementation | Status |
|---|---|---|
| README: Cache-Aside checks cache, loads from DB on miss, invalidates on update | `CacheAsideService.Get` & `CacheAsideService.Update` in `internal/cache/patterns.go` | MATCH |
| README: Write-Through updates DB and cache synchronously | `WriteThroughService.Update` in `internal/cache/patterns.go` | MATCH |
| README: Write-Behind updates cache immediately; flushes to DB asynchronously | `WriteBehindService.Update` & `flushWorker` in `internal/cache/patterns.go` | MATCH |
| README: SingleFlight coalesces concurrent cache miss requests using `golang.org/x/sync/singleflight` | `SingleFlightService.Get` in `internal/cache/stampede.go` | MATCH |
| README: XFetch early refresh formula `-Δ · β · ln(U) > TTL_remaining` | `ShouldRecompute` in `internal/cache/stampede.go` | MATCH |
| README: SWR returns stale data immediately while asynchronously triggering background DB revalidation | `SWRService.Get` & `triggerRevalidate` in `internal/cache/stampede.go` | MATCH |
| README: TTL Jitter adds randomized offset `[0, maxJitter)` | `TTLWithJitter` in `internal/cache/store.go` | MATCH |
| Demo Output in execution result vs actual `go run ./cmd/demo` | Executed live: identical numbers, queries, and structure | MATCH |

## Identified Discrepancies

None found. The README accurately describes the file tree, execution commands, and behavioral semantics implemented in the codebase.
