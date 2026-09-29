# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Implementation Files:
- internal/cache/*.go
- cmd/demo/main.go
Tests:
- tests/cache_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs: (not audited now)
Main Claims To Verify:
1. Cache-Aside read/write/invalidation behavior.
2. Write-Through synchronous consistency.
3. Write-Behind async flush and queue handling.
4. Naive stampede causes multiple DB queries.
5. SingleFlight reduces to one query.
6. XFetch triggers early recompute per formula.
7. Stale-While-Revalidate serves stale data then refreshes.
8. TTL jitter adds random offset.
Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency bugs, queue overflow silent drops, timing‑sensitive tests, untested error paths.
