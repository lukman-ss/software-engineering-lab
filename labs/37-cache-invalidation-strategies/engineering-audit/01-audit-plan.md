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
1. Cache-Aside invalidates cache on update, forcing fresh DB query.
2. Write-Through synchronously updates cache & DB on update.
3. Write-Behind updates cache immediately and flushes asynchronously to DB.
4. Singleflight coalesces 20 concurrent cache-miss reads into 1 DB query.
5. XFetch early refresh formula triggers correctly when `-Δ · β · ln(U) > TTL_remaining`.
6. Stale-While-Revalidate serves stale data immediately while triggering async background revalidation.
7. TTL Jitter adds randomized offset `[0, maxJitter)` to base TTL.
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Data race conditions under concurrent cache/DB access.
- Inaccurate mathematical implementation of XFetch probabilistic formula.
- Unhandled channel blocking or goroutine leaks in Write-Behind / SWR background workers.
- Documentation/README mismatches with code reality.
