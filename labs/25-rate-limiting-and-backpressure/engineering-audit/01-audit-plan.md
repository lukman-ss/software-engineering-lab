# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Implementation Files:
- internal/ratelimit/bucket.go
- internal/ratelimit/registry.go
- internal/backpressure/queue.go
- internal/retry/backoff.go
- internal/httputil/middleware.go
- go.mod

Tests:
- internal/ratelimit/bucket_test.go
- internal/backpressure/queue_test.go
- internal/retry/backoff_test.go
- internal/httputil/middleware_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- research-revision/03-revision-result.md
- research-audit/07-verdict.md

Main Claims To Verify:
1. Token Bucket allows burst up to capacity $B$ and refills at rate $R$.
2. Leaky Bucket smooths request rates and blocks burst when water level reaches capacity.
3. Multi-tenant Registry isolates limits per tenant key (API key / tenant ID) guarding against CGNAT IP collisions (RFC 6598).
4. BoundedQueue enforces backpressure via non-blocking rejection (`TrySubmit`) returning `ErrQueueFull` when buffer is full.
5. AWS retry backoff strategies (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) correctly calculate intervals adhering to AWS formulas.
6. HTTP middleware returns RFC 6585 status 429 with `Retry-After` header and JSON error payload.
7. Concurrency safety under race detector (`go test -race ./...`).
8. README and demo outputs represent actual working code without fabricated outputs.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Flaky tests due to wall-clock `time.Sleep` calls in rate refill or leaky bucket tests.
- Goroutine or channel leaks during BoundedQueue shutdown.
- Arithmetic edge cases in backoff calculation (e.g. division by zero, float precision, zero elapsed time).
