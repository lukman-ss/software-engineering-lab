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
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Token Bucket allows burst up to capacity $B$ and continuously refills at rate $R$.
2. Leaky Bucket smooths traffic to leak rate $R$, rejecting immediate bursts exceeding capacity.
3. Multi-tenant key registry isolates quotas per tenant to avoid CGNAT IP collisions (RFC 6598).
4. Bounded Queue provides immediate backpressure rejection (`ErrQueueFull`) when capacity is full without blocking or memory leaks.
5. AWS retry backoff strategies (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) accurately adhere to Marc Brooker / AWS Architecture formulas.
6. HTTP middleware implements RFC 6585 compliance (`429 Too Many Requests` status, `Retry-After` header, structured JSON response).
7. Thread safety across all concurrent operations with Go race detector verification.

Commands To Run:
- `go test ./...`
- `go test -v -count=1 ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Floating point inaccuracies or negative tokens/water levels in bucket algorithms.
- Deadlocks, goroutine leaks, or race conditions during bounded queue teardown (`Stop()`).
- Flaky unit tests relying on real wall-clock sleeps.
- Discrepancy between README claims, design notes, and actual code implementation.
