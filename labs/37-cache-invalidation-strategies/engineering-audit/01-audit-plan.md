# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Implementation Files:
- internal/cache/store.go
- internal/cache/repo.go
- internal/cache/patterns.go
- internal/cache/stampede.go
- cmd/demo/main.go
Tests:
- tests/cache_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- research-revision/01-revision-plan.md
- research-revision/02-changes-made.md
- research-revision/03-revision-result.md
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
1. Cache-Aside pattern: reads load on miss, subsequent reads hit cache, updates write to DB then invalidate cache.
2. Write-Through pattern: updates write synchronously to DB and cache; subsequent reads hit cache directly.
3. Write-Behind pattern: updates write to cache immediately, asynchronous background worker flushes to DB.
4. SingleFlight stampede mitigation: coalesces concurrent reads on miss down to exactly 1 DB query.
5. XFetch algorithm: evaluates probabilistic early expiration via formula `-Δ · β · ln(U) > TTL_remaining` with correct negative sign.
6. Stale-While-Revalidate: serves stale cache entry immediately within `staleDelta` window while asynchronously triggering background DB revalidation.
7. TTL Jitter: adds positive random duration `[0, maxJitter)` to base TTL to desynchronize expirations.
8. Concurrency safety: all stores and services operate safely under Go race detector without data races or deadlocks.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent reads and cache updates.
- Inaccurate XFetch formula implementation (missing negative sign or bad edge values for random variable U).
- Goroutine or queue leaks in Write-Behind service or SWR revalidator.
- SWR revalidation thrashing (multiple goroutines spawning per expired key).
- Write-Behind silent write loss behavior under high load.
- Discrepancy between README claims and implemented code.
