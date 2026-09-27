# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Implementation Files:
- go.mod
- internal/ratelimit/bucket.go
- internal/ratelimit/registry.go
- internal/backpressure/queue.go
- internal/retry/backoff.go
- internal/httputil/middleware.go
- cmd/demo/main.go
Tests:
- internal/ratelimit/bucket_test.go
- internal/backpressure/queue_test.go
- internal/retry/backoff_test.go
- internal/httputil/middleware_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs:
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
1. Token Bucket allows bursts up to capacity then enforces refill rate.
2. Leaky Bucket smooths flow to constant leak rate, rejecting bursts when capacity exceeded.
3. Bounded Queue rejects excess requests immediately when full (non-blocking backpressure).
4. Retry backoff produces jittered intervals matching AWS Full Jitter, Equal Jitter, etc.
5. HTTP middleware returns RFC 6585 429 with Retry-After header.
6. Thread-safe under concurrent access (race detector clean).
7. Demo shows burst allowance, leaky smoothing, backpressure shedding, and jitter spread.
Commands To Run:
- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Incorrect token/leaky bucket arithmetic leading to wrong rate limits.
- Queue blocking instead of fast rejection.
- Race conditions in mutable state.
- Non-compliant backoff distributions.
- Missing or incorrect HTTP headers.