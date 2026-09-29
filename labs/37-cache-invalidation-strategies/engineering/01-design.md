# Engineering Design

Target Lab: labs/37-cache-invalidation-strategies
Research Status: APPROVED

## Concept To Prove
Demonstrate and compare:
1. Cache invalidation patterns:
   - Cache-Aside (Lazy loading on miss; DB write first then cache invalidate).
   - Write-Through (Synchronous write to DB followed immediately by best-effort cache update).
   - Write-Behind / Write-Back (Immediate write to cache queue/dirty buffer, asynchronous flush to DB, demonstrating throughput and durability risk on un-flushed crashes).
2. Cache stampede (thundering herd) mitigations under high concurrency:
   - Naive (Broken) cache reload causing database stampede (`N` concurrent queries on expiry).
   - Single-Flight coalescing (`singleflight.Group`) ensuring exactly 1 DB query per expired key across concurrent requests.
   - Probabilistic early expiration (XFetch: `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`).
   - Stale-While-Revalidate (SWR: immediate return of stale value while triggering asynchronous background refresh).
   - Jittered TTL as an anti-synchronization mechanism across distinct keys.

## Expected Behavior
- **Cache-Aside**: Initial read misses and queries DB, populating cache. Subsequent reads hit cache. Updates write to DB, invalidate cache, forcing next read to refresh from DB.
- **Write-Through**: Updates write to DB and update cache in sequence. Subsequent reads immediately hit updated cache without DB query.
- **Write-Behind**: Updates write to in-memory/cache store and enqueue flush. Reads hit cache immediately. Background worker flushes to DB. If service terminates/crashes before flush, un-flushed items are lost.
- **Stampede / Single-Flight**: When a cached item expires and `N` goroutines query simultaneously, naive queries DB `N` times; single-flight coordinates readers so DB is queried exactly once.
- **XFetch**: Key refreshed proactively before official TTL based on measured computation time `Δ` and random draw `ln(rand())`.
- **Stale-While-Revalidate**: Expired key within stale window returns cached value immediately; asynchronous revalidation fetches latest data from DB and updates cache.

## Failure Scenario
- **Naive stampede**: On key expiration under `N` concurrent goroutines, DB is flooded with `N` identical queries (dog-piling/resource exhaustion).
- **Write-behind data loss**: Unflushed queue drops pending updates if memory is cleared before DB write.
- **Formula sign bug**: If XFetch is implemented as `Δ · β · ln(rand()) > TTL_remaining` without negative sign, `ln(rand())` is negative and early refresh never fires.

## Success Criteria
1. Full test suite passes including Go race detector (`go test -race ./...`).
2. Single-flight reduces concurrent miss DB queries from `N` to `1`.
3. XFetch triggers refresh before TTL when `-Δ · β · ln(rand()) > TTL_remaining`.
4. Stale-While-Revalidate serves stale data immediately while updating cache asynchronously in the background.
5. All 3 write patterns (Cache-Aside, Write-Through, Write-Behind) pass verified behavioral checks.
6. Demo executes all scenarios clearly in terminal.

## Architecture
- `internal/cache`:
  - `store.go`: Memory-backed cache with TTL, expiration metadata, and locking.
  - `repo.go`: Simulated SQL/database repository with query counters and simulated latency.
  - `patterns.go`: Cache-Aside, Write-Through, Write-Behind services.
  - `stampede.go`: Naive vs Single-Flight vs XFetch vs SWR implementations.
  - TTL jitter calculation is located in `store.go`.
- `cmd/demo/main.go`: End-to-end runnable comparison of patterns and stampede mitigations.
- `tests/`: Integration and concurrency unit tests verifying stampede reduction, XFetch mathematical trigger, SWR behavior, and write policies.

## Components
1. `MemoryCache`: In-memory thread-safe key-value store supporting TTL and expiration metadata inspection.
2. `MockDB`: Concurrent in-memory store simulating database latency and counting query invocations.
3. `CacheAside`: Checks cache -> DB fallback -> set cache. Write -> DB update -> delete cache.
4. `WriteThrough`: Write -> DB update -> update cache. Read -> check cache -> DB fallback.
5. `WriteBehind`: Write -> cache update -> queue -> async worker flushes to DB.
6. `SingleFlightService`: Uses Go `golang.org/x/sync/singleflight` to coalesce concurrent reads.
7. `XFetchService`: Computes probabilistic expiration condition `-delta * beta * math.Log(u) > remaining`.
8. `SWRService`: Evaluates `now < ttl` (fresh) vs `now < ttl + staleDelta` (serve stale + async revalidate) vs hard miss.

## Test Strategy
- Unit tests verifying exact query counts on single-flight under concurrent access.
- Tests verifying XFetch probabilistic formula logic with deterministic random draws.
- Tests verifying SWR serving stale data and asynchronously refreshing.
- Concurrency race detection (`-race`).

## Execution Plan
1. Setup Go module `labs/37-cache-invalidation-strategies`.
2. Implement `internal/cache` package with clean Go standard library and `golang.org/x/sync/singleflight`.
3. Write test suite in `tests/` and verify with race detector.
4. Build and run `cmd/demo/main.go`.
5. Capture execution results in `engineering/02-implementation-notes.md` and `engineering/03-execution-result.md`.

## Implementation Decisions
- In-memory thread-safe store and mock DB used instead of requiring external Redis/PostgreSQL instances so the lab is 100% reproducible and testable in standalone CLI environments.
- Correct mathematically verified XFetch formula `-Δ · β · ln(U)` used as instructed by Research Revision Findings.
