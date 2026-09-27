# Engineering Audit Plan

Target Lab: `labs/25-rate-limiting-and-backpressure`
Implementation Files:
- `internal/ratelimit/bucket.go` (TokenBucket & LeakyBucket)
- `internal/ratelimit/registry.go` (Tenant Isolation Registry)
- `internal/backpressure/queue.go` (BoundedQueue with fast load shedding)
- `internal/httputil/middleware.go` (HTTP 429 RFC 6585 middleware)
- `internal/retry/backoff.go` (AWS Architecture jitter retry backoff)

Tests:
- `internal/ratelimit/bucket_test.go`
- `internal/backpressure/queue_test.go`
- `internal/httputil/middleware_test.go`
- `internal/retry/backoff_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Token Bucket allows burst up to capacity $B$ and continuously refills at rate $R$.
2. Leaky Bucket smooths traffic to leak rate $R$, rejecting bursts exceeding capacity.
3. Bounded queue performs fast load shedding (`ErrQueueFull`) without unbounded memory growth or queue age degradation.
4. HTTP middleware produces compliant RFC 6585 HTTP 429 responses with `Retry-After` headers.
5. AWS Retry jitter algorithms (Full, Equal, Decorrelated, NoJitter) mathematically match Marc Brooker's formulas within configured bounds $[0, \text{Cap}]$.
6. Concurrency safety across all shared state primitives under `-race`.
7. Zero fabricated benchmark results or non-executable demo claims.

Commands To Run:
- `go test -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Floating-point time calculations causing token refill drift or negative tokens.
- Race conditions during worker shutdown and job submission in BoundedQueue.
- Memory leak in Registry map over unbounded tenant keys.
- Pseudo-random number generator concurrency contention or non-deterministic bounds.
