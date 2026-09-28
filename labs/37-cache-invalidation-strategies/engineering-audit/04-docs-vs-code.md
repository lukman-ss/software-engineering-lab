# Docs vs Code Audit

## 1. README vs Code Alignment

| README Claim | Code Implementation | Status |
|---|---|---|
| Cache-Aside: checks cache, loads on miss, invalidates on update | `CacheAsideService.Get` and `Update` in `internal/cache/patterns.go` | MATCH |
| Write-Through: updates DB and cache synchronously | `WriteThroughService.Update` in `internal/cache/patterns.go` | MATCH |
| Write-Behind: updates cache immediately, flushes asynchronously | `WriteBehindService.Update` and `flushWorker` in `internal/cache/patterns.go` | MATCH |
| SingleFlight: coalesces concurrent misses to 1 DB query | `SingleFlightService.Get` using `singleflight.Group` in `internal/cache/stampede.go` | MATCH |
| XFetch: `-Δ · β · ln(U) > TTL_remaining` | `ShouldRecompute` and `XFetchService.Get` in `internal/cache/stampede.go` | MATCH |
| Stale-While-Revalidate: returns stale immediately, triggers async refresh | `SWRService.Get` in `internal/cache/stampede.go` | MATCH |
| TTL Jitter: adds randomized offset `[0, maxJitter)` | `TTLWithJitter` in `internal/cache/store.go` | MATCH |

## 2. Directory Structure Verification

README directory structure matches repository layout exactly:
- `cmd/demo/main.go`
- `internal/cache/{store.go, repo.go, patterns.go, stampede.go}`
- `tests/cache_test.go`
- `engineering/{01-design.md, 02-implementation-notes.md, 03-execution-result.md}`
- `go.mod` / `go.sum`

## 3. Demo Output Verification

README commands and demo output in `engineering/03-execution-result.md` were directly compared against actual execution of `go run ./cmd/demo`:
- Cache-Aside 1st GET: 1 DB query, 2nd GET: 1 DB query, update & GET: 2 DB queries.
- Write-Through 1st GET: 1 DB query, update & GET: 1 DB query (direct cache hit).
- Write-Behind: 0 writes immediately, 1 DB write after 50ms flush delay.
- Stampede: Naive 20 queries vs SingleFlight 1 query.
- XFetch: Low rand draw triggers early proactive refresh (2 queries), High rand draw serves cached data (1 query).
- SWR: Initial v1 -> Stale v1 -> Fresh v2.
- TTL Jitter: All sampled TTLs within `[5m0s, 5m30s)`.

Everything matches verified runtime behavior.

Assessment: PASS (no mismatches found)
