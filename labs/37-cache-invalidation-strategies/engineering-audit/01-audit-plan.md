# Engineering Audit Plan

Target Lab: `labs/37-cache-invalidation-strategies`
Implementation Files:
- `internal/cache/store.go` (MemoryCache, Item, TTLWithJitter)
- `internal/cache/repo.go` (MockDB with atomic query/write counters and mock latencies)
- `internal/cache/patterns.go` (CacheAsideService, WriteThroughService, WriteBehindService)
- `internal/cache/stampede.go` (NaiveStampedeService, SingleFlightService, XFetchService, ShouldRecompute, SWRService)

Tests:
- `tests/cache_test.go` (TestCachePatterns, TestStampedeMitigation, TestXFetchLogic, TestStaleWhileRevalidate, TestJitter)

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-revision/02-changes-made.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Cache-Aside: Reads check cache, populate cache on miss. Updates write to DB and delete cache key.
2. Write-Through: Updates write synchronously to DB and cache. Subsequent reads hit cache.
3. Write-Behind: Updates write to cache immediately, enqueueing asynchronous background write to DB.
4. SingleFlight Stampede Mitigation: Concurrently incoming requests for same cold/expired key coalesce into a single DB query.
5. XFetch Probabilistic Early Expiration: Correctly implements the optimal formula `-Δ * β * ln(U) > TTL_remaining` without the common missing negative sign bug.
6. Stale-While-Revalidate (SWR): Returns stale data immediately within stale window while revalidating asynchronously in background.
7. TTL Jitter: Applies positive uniform jitter `[0, maxJitter)` to prevent synchronized cache expiration stampedes.
8. Concurrency safety: Zero data races under `go test -race ./...`.
9. Documentation accuracy: README claims match code APIs and demo outputs.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -v -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Mathematical accuracy of XFetch logarithmic formula under edge cases (`u <= 0` or `u >= 1`).
- Concurrency races on in-memory maps or counters during concurrent stampede simulation.
- Background goroutine leakage in `WriteBehindService` and `SWRService`.
- SWR revalidation duplicate invocation under concurrent requests.
