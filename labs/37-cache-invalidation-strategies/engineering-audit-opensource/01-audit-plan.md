# Engineering Audit Plan

Target Lab: labs/37-cache-invalidation-strategies
Implementation Files: internal/cache/*.go, cmd/demo/main.go
Tests: tests/cache_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: N/A (audit stage focuses on implementation)
Main Claims To Verify:
- Cache-Aside, Write-Through, Write-Behind behave as described.
- SingleFlight coalesces concurrent misses.
- XFetch early expiration logic correct.
- SWR serves stale data and triggers async revalidation.
- TTL jitter adds jitter within range.
Commands To Run:
```
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
Primary Risks:
- Concurrency safety in WriteBehind flush.
- Correctness of XFetch formula.
- SWR revalidation race conditions.
