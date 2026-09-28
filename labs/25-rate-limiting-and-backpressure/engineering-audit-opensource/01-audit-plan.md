# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Date: 2026-09-28
Output Dir: labs/25-rate-limiting-and-backpressure/engineering-audit-opensource/
Scope Override: implementation + tests only. No research/content audit. No code modification.

## Implementation Files

- internal/ratelimit/bucket.go (TokenBucket, LeakyBucket)
- internal/ratelimit/registry.go (per-tenant Registry)
- internal/backpressure/queue.go (BoundedQueue)
- internal/retry/backoff.go (No/Full/Equal/Decorrelated jitter)
- internal/httputil/middleware.go (RFC 6585 429 middleware)
- cmd/demo/main.go (4-section CLI demo)

## Tests

- internal/ratelimit/bucket_test.go (8 tests)
- internal/backpressure/queue_test.go (4 tests)
- internal/httputil/middleware_test.go (2 tests)
- internal/retry/backoff_test.go (3 tests)

## Executable/Demo

- cmd/demo/main.go via `go run ./cmd/demo`

## Approved Research Inputs

NOT AUDITED (pipeline override).

## Main Claims To Verify

1. TokenBucket burst B + refill R; LeakyBucket constant drain R.
2. Per-tenant isolation (CGNAT-safe, key-based not IP-based).
3. BoundedQueue fast ErrQueueFull rejection; ErrQueueStopped after Stop.
4. AWS Brooker jitter formulas + bounds.
5. RFC 6585 429 + Retry-After header + JSON body.
6. Race-clean concurrency.
7. Demo output real; README matches code; no fake benchmark.

## Commands To Run

- go test ./...
- go test -race ./... (GOPROXY=off; toolchain stalls on network)
- go run ./cmd/demo
- go vet ./... (attempted; env-timeout, recorded)

## Primary Risks

- Wall-clock timing tests flaky under load.
- Registry unbounded map growth (no eviction).
- RetryAfterSeconds div-by-zero at refillRate=0.
- Queue Stop vs in-flight submit; job panic kills worker.
- Design doc claims simulated-time helpers tests do not use.
