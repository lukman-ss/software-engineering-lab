# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Date: 2026-09-27
Audit Output Dir: labs/25-rate-limiting-and-backpressure/engineering-audit-opensource/

## Implementation Files

- internal/ratelimit/bucket.go (TokenBucket, LeakyBucket)
- internal/ratelimit/registry.go (per-tenant Registry)
- internal/backpressure/queue.go (BoundedQueue)
- internal/httputil/middleware.go (RFC 6585 429 middleware)
- internal/retry/backoff.go (AWS jitter backoff)
- cmd/demo/main.go (CLI demo)

## Tests

- internal/ratelimit/bucket_test.go (5 tests: BurstAndRefill, LeakRate, TenantIsolation, RetryAfterSeconds, ConcurrencyRace)
- internal/backpressure/queue_test.go (3 tests: RejectionUnderLoad, ConcurrencySafety, SubmitAfterStop)
- internal/httputil/middleware_test.go (1 test: RFC6585)
- internal/retry/backoff_test.go (2 tests: Bounds, DecorrelatedJitter_Bounds)

Total: 11 tests across 4 packages.

## Executable/Demo

- cmd/demo/main.go via `go run ./cmd/demo`

## Approved Research Inputs

Pipeline override: research/content NOT audited in this stage. Audit scoped to implementation + tests only. Engineering design inputs consulted for claims only: engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md, README.md.

## Main Claims To Verify

1. TokenBucket allows bursts up to capacity B, refills at rate R, thread-safe.
2. LeakyBucket smooths to leak rate R, rejects bursts at capacity, thread-safe.
3. Registry isolates per-tenant buckets (CGNAT avoidance).
4. BoundedQueue.TrySubmit rejects fast with ErrQueueFull when full; safe under concurrency; safe shutdown.
5. Retry backoff matches AWS formulas (Full/Equal/No/Decorrelated jitter) within bounds.
6. Middleware returns RFC 6585 429 + Retry-After header on exhaustion.
7. `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` all pass; demo output real.
8. README matches code.

## Commands To Run

- [x] go build ./... -> SUCCESS (exit 0)
- [x] go test -count=1 -v ./... -> all 11 PASS (exit 0)
- [x] go test -race -count=1 ./... -> all ok, no data races (exit 0)
- [x] go run ./cmd/demo -> exit 0, real output
- [x] go vet ./... -> clean
- [x] gofmt -l . -> 3 files flagged (queue.go, bucket.go, registry.go)
- [x] throwaway concurrent Stop+Submit repro (removed afterwards) -> PANIC confirmed

## Primary Risks

- Channel close vs send race in BoundedQueue.Stop/TrySubmit (no data race, so -race is blind to it).
- Demo backpressure section is timing-dependent; recorded output in 03-execution-result.md not reproducible run-to-run.
- 03-execution-result.md test list is stale (9 tests listed, 11 exist).
- Middleware test asserts header presence only, not value correctness or body shape.
