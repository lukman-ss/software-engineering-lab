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
- Research Report (`research/05-report.md`)
- Research Audit Verdict (`research-audit/07-verdict.md`)

Main Claims To Verify:
1. Token Bucket algorithm allows bursts up to capacity and refills continuously at specified rate.
2. Leaky Bucket algorithm enforces constant output leak rate and rejects immediate bursts exceeding water capacity.
3. Multi-tenant Registry isolates rate limiters per API key / tenant key, guarding against RFC 6598 CGNAT IP collision.
4. Bounded Queue Backpressure performs zero-allocation fast load shedding (`ErrQueueFull`) when queue buffer capacity is reached.
5. Exponential Backoff with AWS Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter adheres to Marc Brooker (AWS Architecture) standards.
6. HTTP Middleware returns standard RFC 6585 `429 Too Many Requests` with exact `Retry-After` header and JSON error body.
7. Concurrency safety verified under Go race detector (`go test -race ./...`).

Commands To Run:
- `go test -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or state corruption under concurrent rate limit checks.
- Non-deterministic flakiness in wall-clock time-based unit tests.
- Incorrect RFC 6585 header formatting or status code in HTTP middleware.
- Unbounded queue age or deadlock on queue termination.
