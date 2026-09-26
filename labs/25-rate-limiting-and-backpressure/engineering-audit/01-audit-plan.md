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
- Research Plan & Findings: Token Bucket, Leaky Bucket, Bounded Queue Backpressure, RFC 6585 (429 + Retry-After), RFC 6598 (CGNAT collision avoidance), Marc Brooker / AWS Jitter Backoff.

Main Claims To Verify:
1. Token bucket accommodates initial bursts up to capacity and refills at constant rate.
2. Leaky bucket prevents bursts exceeding capacity and drains at leak rate.
3. Bounded queue immediately rejects excess jobs under load (`ErrQueueFull`) without unbounded memory growth.
4. HTTP 429 middleware sets standard RFC 6585 `Retry-After` header and JSON response.
5. Multi-tenant registry isolates tenant rate limits to prevent CGNAT IP starvation.
6. Retry backoff adheres to AWS Full Jitter, Equal Jitter, No Jitter bounds.
7. Concurrency safety verified with Go race detector without data races.
8. Demo execution produces legitimate live metrics.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions in rate limit token accumulation or queue status updates.
- Flaky tests caused by real-time wall-clock `time.Sleep`.
- Mismatch between README claims and actual exposed package APIs.
