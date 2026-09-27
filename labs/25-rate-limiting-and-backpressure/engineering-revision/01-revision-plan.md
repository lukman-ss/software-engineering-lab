# Engineering Revision Plan

Target Lab: labs/25-rate-limiting-and-backpressure
Previous Verdict: NEEDS_REVISION

## Blocking Issues

1. **RACE_CONDITION / Shutdown Panic in BoundedQueue**: Concurrent `Stop()` and `TrySubmit()` causes panic (`send on closed channel`) due to race window between `stopped.Load()` and `select { case bq.queue <- job: ... }`. (Severity: HIGH)
2. **IMPLEMENTATION_OVERCLAIM**: Design/docs claim full concurrency safety including shutdown, but concurrent submit during `Stop()` panics. (Severity: HIGH)

## Non-Blocking Issues

1. **TEST_CLAIM_MISMATCH**: `engineering/03-execution-result.md` transcript omits `TestBoundedQueue_SubmitAfterStop` and `TestTokenBucket_RetryAfterSeconds`. (Severity: MEDIUM)
2. **UNVERIFIED_RESULT**: `engineering/03-execution-result.md` presents non-deterministic demo output as canonical without timing caveats. (Severity: MEDIUM)
3. **MISSING_TESTS**: Missing coverage for Registry concurrent same-key `Get`, LeakyBucket concurrency, Middleware HTTP body/header details and anonymous fallback, and edge cases. (Severity: LOW)
4. **DOC_CODE_MISMATCH**: `README.md` omits `engineering-revision/` and audit directories from structure tree. (Severity: LOW)
5. **CODE_FORMATTING**: 3 files flagged by `gofmt`. (Severity: LOW)

## Files To Change

- `internal/backpressure/queue.go`: Fix race window between `Stop()` and `TrySubmit()` using a Mutex guard or state check to prevent send on closed channel.
- `internal/backpressure/queue_test.go`: Add `TestBoundedQueue_ConcurrentStopAndSubmit`.
- `internal/ratelimit/bucket_test.go`: Add `TestLeakyBucket_ConcurrencyRace`.
- `internal/ratelimit/registry_test.go`: Add `TestRegistry_ConcurrentSameKeyGet`.
- `internal/httputil/middleware_test.go`: Add assertions for JSON body shape, `Retry-After` numeric value, and anonymous tenant fallback.
- `internal/retry/backoff_test.go`: Add test for unknown strategy fallback.
- `README.md`: Update directory structure to include `engineering-revision/` and `engineering-audit-*` directories.
- `engineering/03-execution-result.md`: Update test list and add non-determinism note to demo output.

## Tests To Add/Modify

- `TestBoundedQueue_ConcurrentStopAndSubmit` in `internal/backpressure/queue_test.go`
- `TestLeakyBucket_ConcurrencyRace` in `internal/ratelimit/bucket_test.go`
- `TestRegistry_ConcurrentSameKeyGet` in `internal/ratelimit/registry_test.go`
- `TestRateLimitMiddleware_AnonymousFallbackAndBody` in `internal/httputil/middleware_test.go`
- `TestComputeBackoff_UnknownStrategy` in `internal/retry/backoff_test.go`

## Validation Commands

```bash
cd labs/25-rate-limiting-and-backpressure
gofmt -l .
go vet ./...
go test -v ./...
go test -race -count=1 ./...
go run ./cmd/demo
```
