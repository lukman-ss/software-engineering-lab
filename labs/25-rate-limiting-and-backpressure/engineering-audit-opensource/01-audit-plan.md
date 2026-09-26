# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Scope: implementation and tests only (research/content excluded per pipeline override)
Audit Date: 2026-09-26

## Implementation Files

- internal/ratelimit/bucket.go (TokenBucket, LeakyBucket)
- internal/ratelimit/registry.go (per-tenant Registry)
- internal/backpressure/queue.go (BoundedQueue)
- internal/httputil/middleware.go (RFC 6585 429 middleware)
- internal/retry/backoff.go (No/Full/Equal/Decorrelated jitter)
- cmd/demo/main.go (CLI demo)

## Tests

- internal/ratelimit/bucket_test.go (5 tests)
- internal/backpressure/queue_test.go (2 tests)
- internal/httputil/middleware_test.go (1 test)
- internal/retry/backoff_test.go (2 tests)

## Executable/Demo

- cmd/demo (go run ./cmd/demo)

## Main Claims To Verify

1. TokenBucket allows bursts to capacity B, refills at rate R.
2. LeakyBucket smooths to leak rate R, rejects bursts at capacity.
3. Per-tenant registry isolates tenants (CGNAT/RFC 6598 rationale).
4. BoundedQueue TrySubmit fails fast with ErrQueueFull at capacity.
5. Backoff matches AWS jitter formulas with correct bounds.
6. Middleware returns RFC 6585 429 + Retry-After header + JSON body.
7. Concurrency-safe (race detector clean).
8. Demo output real, README matches code.

## Commands To Run

- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo

## Primary Risks

- Timing-sensitive tests (real time.Sleep) flaky on loaded CI.
- BoundedQueue Stop/close lifecycle vs concurrent submit.
- Zero/negative rate constructor inputs unvalidated.
- Test suite asserts bounds only, not formula exactness.
