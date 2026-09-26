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
- Token Bucket vs Leaky Bucket rate limiting mechanisms
- RFC 6585 (`429 Too Many Requests` + `Retry-After`) & RFC 6598 (Tenant isolation over CGNAT IP limits)
- Bounded Queue Load Shedding under arrival rate exceeding service capacity
- AWS Retry Jitter strategies (Marc Brooker)

Main Claims To Verify:
1. Token bucket allows burst capacity $B$ and continuously refills at rate $R$.
2. Leaky bucket smooths flow to leak rate $R$, rejecting bursts exceeding capacity.
3. Bounded queue backpressure rejects excess jobs immediately (`TrySubmit`) with fast drop error when full.
4. AWS Jitter retry strategies satisfy mathematical interval bounds.
5. HTTP middleware returns RFC 6585 compliant HTTP 429 status code and `Retry-After` header when rate-limited.
6. Code compiles cleanly, `go test ./...` and `go test -race ./...` pass with 0 race warnings.
7. Interactive demo output (`go run ./cmd/demo`) matches claimed execution results.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or state mutation anomalies under heavy multi-goroutine access.
- Non-deterministic flakiness in bucket refill or retry jitter test assertions.
- Mismatches between README claims and actual internal struct/method behavior.
