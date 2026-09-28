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

Main Claims To Verify:
1. Token Bucket & Leaky Bucket algorithms enforce limits and allow burst accommodation vs traffic smoothing cleanly under concurrency.
2. Multi-tenant registry isolates tenant rate limits (RFC 6598 CGNAT context).
3. Bounded queue backpressure performs fast load shedding (`ErrQueueFull`) when buffer capacity is reached.
4. AWS retry backoff (Full, Equal, Decorrelated Jitter) adheres to Marc Brooker bounds and distribution equations.
5. HTTP middleware returns RFC 6585 `429 Too Many Requests` with standard `Retry-After` header.
6. Code compiles, tests pass cleanly under race detector (`go test -race ./...`), and demo runs without errors.
7. Documentation (`README.md`, `engineering/03-execution-result.md`) accurately reflects codebase execution and output.

Commands To Run:
- `go test -count=1 -v ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions or non-determinism in state transitions (refill/leak timing, registry access, queue stop).
- Mismatch between README / Execution notes and actual runtime output.
- Missing edge case handling (negative inputs, zero parameters, worker termination leaks).
