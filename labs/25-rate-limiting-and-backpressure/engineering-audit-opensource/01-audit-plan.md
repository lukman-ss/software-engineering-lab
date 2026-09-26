# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Implementation Files: 
- go.mod
- internal/ratelimit/bucket.go
- internal/ratelimit/bucket_test.go
- internal/ratelimit/registry.go
- internal/backpressure/queue.go
- internal/backpressure/queue_test.go
- internal/httputil/middleware.go
- internal/httputil/middleware_test.go
- internal/retry/backoff.go
- internal/retry/backoff_test.go
- cmd/demo/main.go
Tests: 
- All *_test.go files in internal/ packages
Executable/Demo: ./cmd/demo/main.go
Approved Research Inputs: research/05-report.md
Main Claims To Verify:
- Token bucket allows bursts up to capacity then enforces refill rate
- Leaky bucket smooths flow to leak rate, rejecting bursts when at capacity
- Bounded queue rejects excess requests immediately when buffer is full
- Retry backoff produces jittered intervals adhering to AWS formulas
- HTTP rate limiter returns status 429 with correct Retry-After header
- Concurrency safety under load
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in shared state (token bucket, registry, bounded queue)
- Incorrect timing calculations in rate limiting algorithms
- Missing error handling or edge cases
- Concurrency safety not fully tested
- Demo output not matching actual implementation behavior