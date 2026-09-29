# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Implementation Files:
- internal/cache/store.go
- internal/cache/repo.go
- internal/cache/patterns.go
- internal/cache/stampede.go
Tests:
- tests/cache_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research-revision/03-revision-result.md (Corrected XFetch formula `-Δ · β · ln(U) > TTL_remaining`, SWR single-revalidation guard, Cache-Aside/Write-Through/Write-Behind lifecycles)
- research/05-report.md
Main Claims To Verify:
1. Cache-Aside loads from DB on miss, caches result, and invalidates cache on write.
2. Write-Through updates DB and cache synchronously on write; subsequent reads hit cache without DB query.
3. Write-Behind updates cache immediately and flushes asynchronously to DB via background worker.
4. SingleFlight collapses concurrent requests on expired keys to exactly 1 DB query.
5. XFetch evaluates probabilistic early expiration formula `-Δ · β · ln(U) > TTL_remaining` with correct mathematical sign.
6. Stale-While-Revalidate returns stale cached data immediately while asynchronously triggering revalidation.
7. TTL Jitter introduces randomized positive offset `[0, maxJitter)` to prevent synchronized expiry.
8. Implementation contains no fake/hardcoded benchmarks or data fabrications.
9. Concurrency is race-free under `go test -race ./...`.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions in concurrent cache access, background flush, or SWR revalidation.
- Incorrect mathematical formulation for XFetch probabilistic trigger.
- Misalignment between documented claims and actual behavior under failure/edge cases.
- Data loss or unhandled queue overflow in Write-Behind service.
