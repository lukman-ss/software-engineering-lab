# Gaps Analysis

Audited live (no code changes): `go test ./...` ok, `go test -race ./internal/... ./cmd/...` ok, `go run ./cmd/demo` ok.

## Gap 1
Type: UNHANDLED_ERROR / MISSING_EDGE_CASE
Location: internal/ratelimit/bucket.go:55-79 (TokenBucket.RetryAfterSeconds), :16-23 (NewTokenBucket)
Description: RetryAfterSeconds divides needed/tb.refillRate; NewTokenBucket accepts refillRate<=0 unchecked -> division by zero -> int(+Inf) runtime panic or hang.
Test Coverage: None.
Claim vs Reality: README implies robustness; not exercised. Severity MEDIUM. Mitigation: document non-zero refillRate; add zero-guard.

## Gap 2
Type: RACE_CONDITION (latent)
Location: internal/backpressure/queue.go:91-103 (Stop), :67-85 (TrySubmit)
Description: Stop() holds stopMu.Lock before close(queue); TrySubmit only holds RLock while enqueueing. Submitters past the stopped check and between RUnlock and enqueue can send on a closed channel -> panic. TestBoundedQueue_ConcurrentStopAndSubmit passed but does not target that window; race detector cannot prove absence of a narrow race.
Test Coverage: TestBoundedQueue_ConcurrentStopAndSubmit.
Claim vs Reality: Tests pass under -race; claim of "race-clean" empirically satisfied this run, but residual latent path exists. Severity MEDIUM.

## Gap 3
Type: MISSING_TEST / UNHANDLED_ERROR (resilience)
Location: internal/backpressure/queue.go:59 (`_ = job(bq.ctx)`)
Description: workerLoop discards job errors and has no recover; a panicking job kills the worker silently, shrinking pool to zero, while Stats().processed may under-count. Job errors never propagated.
Test Coverage: None submit failing/panicking jobs.
Claim vs Reality: README does not explicitly claim error-isolation; design doc says "failure handling" in audit plan scope. Severity MEDIUM.

## Gap 4
Type: MISSING_EDGE_CASE
Location: internal/backpressure/queue.go:59
Description: Jobs execute against queue-lifetime ctx only; no per-job timeout/deadline honored inside job body (job sleeps 50ms unconditionally in demo, no cancellation check inside heavy work).
Test Coverage: None. Severity LOW.

## Gap 5
Type: UNVERIFIED / IMPLEMENTATION_OVERCLAIM
Location: internal/retry/backoff.go:50-58 (DecorrelatedJitter)
Description: Test only checks [base,cap] bounds over 5 chained iterations; does not verify the decorrelated *chain* property sleep_i depends on sleep_{i-1}, nor monotonic-ish growth, nor math correctness vs Brooker formula beyond bounds.
Test Coverage: TestDecoratedJitter_Bounds.
Claim vs Reality: README "matching AWS Architecture specifications" — formula matches but correctness of sequence not proven. Severity LOW.

## Gap 6
Type: DOC_CODE_MISMATCH (low)
Location: engineering/01-design.md:93 (Decisions) vs bucket_test.go
Description: Design claims "Deterministic simulated-time helpers in tests to avoid flaky wall-clock sleeps"; tests use time.Sleep real-time waits, not simulated clock.
Test Coverage: N/A (doc claim).
Claim vs Reality: Mild mismatch; tests still pass reliably at human tolerance. Severity LOW.

## Gap 7
Type: UNVERIFIED_RESULT / FAKE_DEMO risk
Location: engineering/03-execution-result.md section 3/4
Description: Section 3 queue Stats (Processed=1) differs from live rerun (Processed=2) due to timing; section 4 jitter samples differ numerically across runs by design (disclosed). Document correctly disclaims. No fabrication. Severity LOW / informational.

## Gap 8
Type: MISSING_TEST
Location: README features vs tests
Description: README mentions "Exponential Backoff with jitter" and "429 response with Retry-After" — the Retry-After header value consistency with JSON body field not asserted; anonymous fallback Content-Type asserted but not body field equality.
Test Coverage: middleware_test.go asserts header present, Content-Type set, but not retry_after body == header value.
Claim vs Reality: Partial. Severity LOW.

## Summary

8 gaps. Severities: MEDIUM x3 (2,3,1), LOW x5. No CRITICAL (no fabricated results; demo verified live; race under -race clean this run). No FAKE_BENCHMARK present. No RESEARCH_MISMATCH (research not audited). No core-behavior failure.