# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Audit Output Dir: labs/37-cache-invalidation-strategies/engineering-audit-opensource/
Audit Date: 2026-09-29

## Implementation Files
- internal/cache/store.go (MemoryCache, TTLWithJitter)
- internal/cache/repo.go (MockDB)
- internal/cache/patterns.go (CacheAside, WriteThrough, WriteBehind)
- internal/cache/stampede.go (Naive, SingleFlight, XFetch, SWR)
- cmd/demo/main.go (runnable demo)
- go.mod (module github.com/lukman/labs/37-cache-invalidation-strategies, go 1.22, x/sync v0.7.0)

## Tests
- tests/cache_test.go (TestCachePatterns, TestStampedeMitigation, TestXFetchLogic, TestStaleWhileRevalidate, TestJitter)

## Executable/Demo
- go run ./cmd/demo

## Approved Research Inputs
- OUT OF SCOPE per pipeline override. Implementation + tests only. Engineering design docs (engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md) and README.md used only as claimed-behavior reference vs code.

## Main Claims To Verify
1. Cache-Aside: miss→DB→populate; hit→no query; update→DB write + invalidate.
2. Write-Through: update writes DB+cache synchronously; next read hits cache.
3. Write-Behind: update visible in cache immediately; DB flushed async; Close drains.
4. Naive stampede: N concurrent misses → N DB queries.
5. SingleFlight: N concurrent misses → exactly 1 DB query, all correct.
6. XFetch formula: -Δ·β·ln(U) > TTL_remaining, correct sign, guarded u range.
7. SWR: stale served immediately + async background revalidation completes.
8. TTL jitter: result in [base, base+maxJitter).
9. Race-safe under concurrency; demo output real; README matches code.

## Commands To Run
- go build ./...
- go test -v ./...
- go test -race ./... (fresh, uncached)
- go run ./cmd/demo
- go vet ./... (supplementary)

## Primary Risks
- Timing-sensitive tests (SWR 20ms/300ms + Sleep; WriteBehind 50ms Sleep) flaky under load.
- XFetch service Get() logic vs pure-function test gap (formula tested, service path demo-only).
- WriteBehind silent drop on queue overflow + ignored flush errors.
- Failure paths (DB errors, ctx cancel) untested.
