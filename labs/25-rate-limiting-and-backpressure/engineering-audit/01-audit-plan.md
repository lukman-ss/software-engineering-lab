# Engineering Audit Plan

Target Lab: `labs/25-rate-limiting-and-backpressure`
Implementation Files:
- `internal/ratelimit/bucket.go`
- `internal/ratelimit/registry.go`
- `internal/backpressure/queue.go`
- `internal/retry/backoff.go`
- `internal/httputil/middleware.go`

Tests:
- `internal/ratelimit/bucket_test.go`
- `internal/backpressure/queue_test.go`
- `internal/retry/backoff_test.go`
- `internal/httputil/middleware_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)
- `engineering/01-design.md`

Main Claims To Verify:
1. Token Bucket supports burst capacity up to limit and continuous refill rate.
2. Leaky Bucket smooths request rates and limits peak burst capacity.
3. Multi-tenant key registry isolates quotas between distinct tenants.
4. BoundedQueue implements fast drop / backpressure without blocking callers when capacity is reached.
5. AWS retry backoff jitter formulas (Full Jitter, Equal Jitter, No Jitter, Decorrelated Jitter) conform to Marc Brooker specifications.
6. HTTP middleware correctly returns RFC 6585 status 429 and `Retry-After` header.
7. Concurrency safety under multi-goroutine access verified with `-race`.

Commands To Run:
- `go test -count=1 -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent token consumption or queue operations.
- Flaky tests due to wall-clock timing in bucket refill tests.
- Divergence between demo console output and claimed execution logs.
