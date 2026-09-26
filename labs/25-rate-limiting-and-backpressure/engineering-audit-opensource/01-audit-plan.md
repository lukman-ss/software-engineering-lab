# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Output: labs/25-rate-limiting-and-backpressure/engineering-audit-opensource/
Audit Date: 2026-09-26
Scope Override: implementation and tests only. No research/content audit. No code modification.

## Implementation Files

- internal/ratelimit/bucket.go (TokenBucket, LeakyBucket, RetryAfterSeconds)
- internal/ratelimit/registry.go (per-tenant Registry)
- internal/backpressure/queue.go (BoundedQueue, TrySubmit, Stop, Stats)
- internal/httputil/middleware.go (RateLimitMiddleware, 429 + Retry-After)
- internal/retry/backoff.go (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter)
- cmd/demo/main.go (4-stage CLI demo)

## Tests

- internal/ratelimit/bucket_test.go (5 tests: burst/refill, leak, isolation, retry-after, race smoke)
- internal/backpressure/queue_test.go (2 tests: rejection under load, concurrency safety)
- internal/httputil/middleware_test.go (1 test: RFC 6585 429 + Retry-After)
- internal/retry/backoff_test.go (2 tests: bounds, decorrelated bounds)

## Executable/Demo

- cmd/demo/main.go via `go run ./cmd/demo`

## Main Claims To Verify

1. TokenBucket burst to capacity B then refill at rate R (fractional, thread-safe).
2. LeakyBucket constant drain R, burst rejection at capacity.
3. Per-tenant registry isolation (CGNAT avoidance via X-API-Key).
4. BoundedQueue non-blocking TrySubmit with fast ErrQueueFull; no OOM/unbounded growth.
5. AWS jitter formulas (Brooker): Full/Equal/No/Decorrelated within specified bounds.
6. Middleware returns RFC 6585 429 + Retry-After header + JSON body.
7. Race-clean under concurrent load; demo output real.

## Commands To Run

- go build ./...
- go vet ./...
- go test -v -count=1 ./...
- go test -race -count=1 ./...
- go run ./cmd/demo

## Primary Risks

- Stop()/TrySubmit close-channel race (panic / job drop) — shutdown path.
- RetryAfterSeconds with refillRate=0 (float div-by-zero → bad int).
- Weak tests: bounds-only backoff checks, assertion-free race smoke test, single middleware test.
- Docs overclaim: "deterministic simulated time helpers" (tests use real sleeps); demo omits DecorrelatedJitter.
