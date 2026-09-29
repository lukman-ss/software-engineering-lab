# Documentation vs. Code Audit

Target Lab: `labs/37-cache-invalidation-strategies`

## Comparison Matrix

| Claim / Topic | Documentation (README / Engineering Notes) | Code Implementation | Status |
|---|---|---|---|
| Cache-Aside | Reads DB on miss, writes DB then deletes cache | `CacheAsideService` in `patterns.go` | MATCH |
| Write-Through | Writes DB and cache synchronously on update | `WriteThroughService` in `patterns.go` | MATCH |
| Write-Behind | Writes cache immediately, queues async DB flush | `WriteBehindService` in `patterns.go` | MATCH |
| SingleFlight | Coalesces concurrent cache misses to 1 query | `SingleFlightService` in `stampede.go` | MATCH |
| XFetch Formula | `-Δ · β · ln(U) > TTL_remaining` (`U ~ Uniform(0,1)`) | `ShouldRecompute` in `stampede.go` | MATCH |
| Stale-While-Revalidate | Serves stale value immediately + async revalidate | `SWRService` in `stampede.go` | MATCH |
| TTL Jitter | `[base, base + maxJitter)` positive offset | `TTLWithJitter` in `store.go` | MATCH |
| Run Commands | `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` | Tested and fully operational | MATCH |

## Detailed Checks

1. **DOC_CODE_MISMATCH**: None detected. All features documented in `README.md` and `engineering/01-design.md` match the Go codebase.
2. **TEST_CLAIM_MISMATCH**: None detected. The test suite directly asserts the exact behaviors outlined in design and implementation docs.
3. **RESEARCH_IMPLEMENTATION_MISMATCH**: None detected. Approved research revisions (notably the `-Δ · β · ln(U)` negative sign correction in XFetch and request-triggered SWR revalidation guard) are faithfully implemented.
4. **BENCHMARK / DEMO ACCURACY**: Demo prints live metrics from in-memory counters and time-tracked operations. No hardcoded or fabricated mock results present.
