# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Implementation Files:
- `internal/cache/store.go`
- `internal/cache/repo.go`
- `internal/cache/patterns.go`
- `internal/cache/stampede.go`
- `cmd/demo/main.go`
- `go.mod`
- `go.sum`
Tests:
- `tests/cache_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-revision/01-revision-plan.md`
- `research-revision/02-changes-made.md`
- `research-revision/03-revision-result.md`
Main Claims To Verify:
1. Cache-Aside pattern: cache miss queries DB and sets cache; updates write DB and invalidate cache.
2. Write-Through pattern: updates write DB and cache synchronously; subsequent reads hit cache.
3. Write-Behind pattern: updates write cache immediately and enqueue async DB flush.
4. Stampede mitigation: Naive causes multiple DB queries on concurrent misses; SingleFlight coalesces `N` requests into 1 DB query.
5. XFetch early expiration: correct mathematical condition `-Δ · β · ln(U) > TTL_remaining` evaluated properly.
6. Stale-While-Revalidate: serves stale data immediately within window and executes async revalidation.
7. TTL Jitter: adds random positive jitter within `[0, maxJitter)` to prevent synchronized expiry.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go test -v -count=1 ./tests/...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions during concurrent cache reads and singleflight/SWR async executions.
- SWR test flakiness due to timing dependencies with `time.Sleep`.
- Overclaiming in docs regarding distributed caching or unbounded queue durability in Write-Behind.
