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
- research/05-report.md (XFetch mathematical formula fix: `-Δ · β · ln(U) > TTL_remaining`, Singleflight coalescing, SWR, Write patterns)
Main Claims To Verify:
1. Cache-Aside, Write-Through, Write-Behind read/write patterns operate correctly according to their invariants.
2. Singleflight coalesces 20 concurrent cache miss requests on a single hot key to exactly 1 database query.
3. XFetch probabilistic early expiration implementation correctly implements `-Δ · β · ln(U) > TTL_remaining` and triggers recomputation prior to TTL expiration under low random draw.
4. Stale-While-Revalidate (SWR) returns stale cached values immediately while asynchronously executing background database revalidation.
5. TTL Jitter produces positive random offsets within `[0, maxJitter)`.
6. Code compiles, passes test suite with race detector (`go test -race ./...`), and demo runs without errors.
7. README accurately describes code structure, test commands, and implementation features.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions in concurrent cache access, background flush worker (Write-Behind), or background revalidation worker (SWR).
- Formula sign errors in XFetch implementation (`math.Log(u)` vs `-math.Log(u)`).
- Silent queue drops in Write-Behind during queue overflow without test/documentation coverage.
- Discrepancies between README claims and underlying code implementation.
