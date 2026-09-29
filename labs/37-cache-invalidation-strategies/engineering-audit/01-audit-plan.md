# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Implementation Files:
- `internal/cache/store.go`
- `internal/cache/repo.go`
- `internal/cache/patterns.go`
- `internal/cache/stampede.go`
Tests:
- `tests/cache_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research-revision/03-revision-result.md`
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`

Main Claims To Verify:
1. Cache-Aside pattern: cache miss queries DB, cache hit skips DB, write invalidates cache forcing subsequent DB query.
2. Write-Through pattern: write synchronously updates DB and cache; subsequent reads hit cache without DB query.
3. Write-Behind pattern: write updates cache immediately, background worker flushes asynchronously to DB.
4. Naive stampede: concurrent cache miss causes multiple DB queries (`N > 1`).
5. SingleFlight stampede mitigation: concurrent cache miss coalesces into exactly 1 DB query across `N` concurrent goroutines.
6. XFetch probabilistic early expiration: evaluates formula `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`; early refresh fires when formula evaluates to true.
7. Stale-While-Revalidate (SWR): returns stale cached data immediately within stale window while asynchronously triggering background DB revalidation.
8. TTL Jitter: adds positive offset within `[0, maxJitter)` to prevent synchronized cache expiration.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent cache reads/writes and background revalidation.
- Flaky tests dependent on real-time sleeps in asynchronous workflows.
- Mathematical sign errors or edge conditions in XFetch `ShouldRecompute`.
- Unbounded background goroutine spawning in SWR triggering stampedes on background revalidation.
- Discrepancies between README documentation claims and actual code behavior.
