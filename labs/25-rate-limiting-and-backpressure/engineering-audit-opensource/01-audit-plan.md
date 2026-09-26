# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Scope Override: implementation and tests only. Research/content excluded. No code modification.

Implementation Files:
- internal/ratelimit/bucket.go (TokenBucket, LeakyBucket)
- internal/ratelimit/registry.go (Registry per-tenant buckets)
- internal/backpressure/queue.go (BoundedQueue worker pool)
- internal/httputil/middleware.go (429 middleware)
- internal/retry/backoff.go (No/Full/Equal/Decorrelated jitter)
- cmd/demo/main.go (CLI demo)

Tests:
- internal/ratelimit/bucket_test.go (5 tests: burst, leaky, isolation, retry-after, race)
- internal/backpressure/queue_test.go (3 tests: rejection, concurrency, stop)
- internal/httputil/middleware_test.go (1 test: RFC6585 200->429)
- internal/retry/backoff_test.go (2 tests: bounds, decorrelated bounds)

Executable/Demo: cmd/demo via `go run ./cmd/demo`

Approved Research Inputs: EXCLUDED per pipeline override (not audited).

Main Claims To Verify:
1. TokenBucket burst B + refill R, thread-safe.
2. LeakyBucket constant drain R, burst rejection.
3. Per-tenant registry isolation (CGNAT/RFC6598 avoidance).
4. BoundedQueue non-blocking TrySubmit -> ErrQueueFull when full; Stop idempotent.
5. Backoff formulas match AWS spec (Full/Equal/No/Decorrelated).
6. Middleware returns RFC6585 429 + Retry-After + JSON on exhaustion.
7. Race detector clean; demo output real and reproducible.

Commands To Run:
- go build ./...
- go vet ./...
- go test ./... (record)
- go test -race ./... (record)
- go run ./cmd/demo (record, repeat 3x for flake check)
- go test -v -count=1 ./... (record)

Primary Risks:
- Wall-clock sleep tests flaky (design claims simulated time, code uses time.Sleep).
- BoundedQueue TrySubmit vs Stop race (send on closed channel panic).
- Registry unbounded growth (no eviction).
- Demo queue section timing-dependent (accepted/rejected counts may vary run to run).
- Thin HTTP negative coverage (no anonymous-fallback, body-shape, header-value tests).
- Concurrency tests assert no invariants (race-only).
