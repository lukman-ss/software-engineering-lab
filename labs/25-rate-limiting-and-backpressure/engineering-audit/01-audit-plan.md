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
- Token Bucket & Leaky Bucket algorithms (burst capacity vs traffic smoothing)
- Bounded Queue Backpressure & Load Shedding (rejection on saturation)
- RFC 6585 HTTP 429 Too Many Requests with `Retry-After` header
- AWS Full Jitter, Equal Jitter, No Jitter, Decorrelated Jitter (Marc Brooker)
- RFC 6598 CGNAT tenant isolation via API key registry

Main Claims To Verify:
1. Token bucket allows burst up to capacity $B$ and replenishes at rate $R$.
2. Leaky bucket smooths flow to rate $R$ and drops requests exceeding capacity.
3. Multi-tenant registry isolates rate limit buckets per tenant key.
4. Bounded queue immediately rejects excess jobs under saturation (`ErrQueueFull`).
5. Retry backoff calculates bounds according to AWS jitter strategies.
6. HTTP middleware returns HTTP 429 with `Retry-After` header when limit exceeded.
7. Concurrency safety under race detector (`go test -race ./...`).
8. Demo runs successfully with real, uncorrupted output matching documented results.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during token refill or channel queue operations.
- Flaky tests dependent on real-time wall clocks.
- Inaccurate jitter formulas or mathematical boundary violations.
- Discrepancies between execution result logs and actual command outputs.
