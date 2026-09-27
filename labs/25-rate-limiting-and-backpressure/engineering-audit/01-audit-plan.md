# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
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
- Token Bucket & Leaky Bucket algorithm specifications
- RFC 6585 (429 Too Many Requests, Retry-After header)
- RFC 6598 (Carrier-Grade NAT / tenant key isolation)
- AWS Architecture Exponential Backoff and Jitter strategies (Marc Brooker)
- Bounded Queue Backpressure and load shedding mechanics

Main Claims To Verify:
1. `TokenBucket` burst allowance up to capacity and fractional refill rate tracking.
2. `LeakyBucket` traffic smoothing and leak rate enforcement.
3. Multi-tenant key registry isolating tenant rate quotas (RFC 6598).
4. `BoundedQueue` non-blocking `TrySubmit` shedding excess load with `ErrQueueFull`.
5. Bounded queue lifecycle stop and shutdown behavior.
6. Exponential backoff jitter algorithms (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) honoring bounds.
7. HTTP 429 response formatting with RFC 6585 compliant `Retry-After` header.
8. Thread safety across rate limiters and queues under concurrent execution with Go race detector clean.
9. Demo executable runs cleanly and outputs real execution results matching claims.

Commands To Run:
```bash
cd labs/25-rate-limiting-and-backpressure
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions or data races on shared state (token updates, channel submission, map access).
- Clock skew or non-monotonic time handling during token refill calculations.
- Blocking or goroutine leaks on queue shutdown / worker lifecycle.
- Unbounded memory growth in registry or queue buffers.
- Discrepancy between code behavior and README documentation.
