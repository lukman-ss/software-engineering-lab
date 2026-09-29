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
1. Cache-Aside, Write-Through, and Write-Behind write/invalidation/flush strategies operate as designed.
2. SingleFlight coalesces N concurrent cache miss queries into exactly 1 database query.
3. XFetch implements probabilistic early expiration using formula `-Δ · β · ln(U) > TTL_remaining`.
4. Stale-While-Revalidate (SWR) serves stale cached values immediately while asynchronously revalidating in the background.
5. TTL Jitter adds random positive offsets `[0, maxJitter)` to base TTLs.
6. Code compiles, tests pass, race detector passes, demo executes cleanly, and README matches implementation.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent cache reads/writes, singleflight coalescing, or background SWR / Write-Behind goroutines.
- Unhandled errors during DB read/write failures.
- Discrepancy between README claims and underlying code logic.
