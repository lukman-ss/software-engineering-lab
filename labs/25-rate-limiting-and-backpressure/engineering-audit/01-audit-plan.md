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
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

Main Claims To Verify:
1. Token bucket allows burst up to capacity $B$ and continuously refills at rate $R$.
2. Leaky bucket smooths flow to leak rate $R$, rejecting immediate bursts when water reaches capacity.
3. Multi-tenant registry isolates rate limit buckets by tenant key / API key to avoid CGNAT IP collisions (RFC 6598).
4. Bounded queue backpressure sheds excess jobs fast (`ErrQueueFull`) using non-blocking submission without unbounded latency or memory growth.
5. AWS retry backoff strategies (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter) adhere to AWS Architecture specifications.
6. HTTP middleware returns RFC 6585 status `429 Too Many Requests` with standard `Retry-After` header.
7. Concurrency safety across token buckets, registries, and bounded queue workers without data races.

Commands To Run:
- `go test -count=1 -v ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- Race conditions during concurrent token refill or queue submission.
- Flaky tests dependent on real-time sleep.
- Mismatch between README documentation and actual implementation behavior.
- Deviation from standard AWS retry formulas or RFC 6585 specifications.
