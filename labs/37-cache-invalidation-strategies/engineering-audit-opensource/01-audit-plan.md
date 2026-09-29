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
- research/05-report.md
- research-revision/03-revision-result.md

Main Claims To Verify:
1. Cache-Aside pattern: cache miss queries DB and populates cache; update invalidates cache.
2. Write-Through pattern: update synchronously writes DB and updates cache.
3. Write-Behind pattern: update writes cache immediately and queues background DB write; queue overflow drops writes when channel buffer is exceeded.
4. SingleFlight stampede mitigation: coalesces concurrent requests on cache miss to 1 DB query using `golang.org/x/sync/singleflight`.
5. XFetch probabilistic early expiration: triggers proactive recomputation using `-Δ · β · ln(U) > TTL_remaining`.
6. Stale-While-Revalidate (SWR): serves stale data immediately within stale window while asynchronously triggering deduplicated background revalidation.
7. TTL Jitter: adds randomized offset `[0, maxJitter)` to base TTL to prevent synchronized key expiration.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Data race conditions in concurrent cache or background goroutines.
- Incorrect mathematical formulation in probabilistic calculations (e.g. XFetch formula sign errors).
- Flaky tests dependent on real wall-clock sleeps (`time.Sleep`).
- Incomplete error propagation or missing failure path handling.
