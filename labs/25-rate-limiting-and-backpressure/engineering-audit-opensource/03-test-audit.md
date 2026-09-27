## Finding 1

Location: internal/ratelimit/bucket_test.go:9-30
Claimed Behavior: Burst capacity + refill timing test.
Observed Implementation: Creates bucket(3,10/sec), consumes 3 tokens, expects 4th fail, sleeps 200ms (~2 tokens), expects 1 token available.
Assessment: PASS
Severity: LOW
Notes: Wall-clock sleep — potential flakiness under heavy load, but CI passes; acceptable.

## Finding 2

Location: internal/ratelimit/bucket_test.go:32-49
Claimed Behavior: LeakyBucket leak rate.
Observed Implementation: capacity=2, leak=5/sec; admits first 2; third fails; sleep 250ms (>0.5*2 leaked), admits fourth.
Assessment: PASS
Severity: LOW
Notes: Uses deterministic sleeps; fine.

## Finding 3

Location: internal/ratelimit/bucket_test.go:51-67
Claimed Behavior: Per-tenant isolation via registry.
Observed Implementation: Registry(1,10); Get("tenant-a") and ("tenant-b"); tenant-a exhausts 1 token; tenant-b still admits.
Assessment: PASS
Severity: LOW
Notes: No concurrency test for registry — but unit test passes; race test in bucket covers mutex.

## Finding 4

Location: internal/ratelimit/bucket_test.go:69-81
Claimed Behavior: RetryAfterSeconds boundary behavior.
Observed Implementation: Bucket(2,2); consumes 2 tokens; RetryAfterSeconds(1) >0; sleep 600ms (>0.5s*2), expects 0.
Assessment: PASS
Severity: LOW
Notes: Good edge test.

## Finding 5

Location: internal/ratelimit/bucket_test.go:83-98
Claimed Behavior: Concurrency safety (no race).
Observed Implementation: 50 goroutines x 10 Allow() calls; no assertion on final token count — only exercises mutex.
Assessment: WARNING
Severity: MEDIUM
Notes: Test lacks post-condition check; could pass even if mutex broken (no data race but logic race). Recommend adding final Tokens() assertion.

## Finding 6

Location: internal/backpressure/queue_test.go:10-45
Claimed Behavior: RejectionUnderLoad fills queue then fast-rejects.
Observed Implementation: Submits one blocking job; fills buffer with capacity tries; next submit must ErrQueueFull.
Assessment: PASS
Severity: LOW
Notes: Correct use of channels; no flaky sleeps.

## Finding 7

Location: internal/backpressure/queue_test.go:47-68
Claimed Behavior: Concurrency safety stats.
Observed Implementation: 30 goroutines TrySubmit (no sleep); after Wait, checks accepted+rejected == 30.
Assessment: PASS
Severity: LOW
Notes: No processing waited — but processed not required. Good.

## Finding 8

Location: internal/backpressure/queue_test.go:70-81
Claimed Behavior: SubmitAfterStop returns ErrQueueStopped.
Observed Implementation: Stop(); TrySubmit -> ErrQueueStopped; double-stop safe.
Assessment: PASS
Severity: LOW
Notes: Idempotent stop correct.

## Finding 9

Location: internal/httputil/middleware_test.go:11-43
Claimed Behavior: RFC 6585 compliance.
Observed Implementation: Registry(1,10); first request 200 OK; second same key -> 429 + Retry-After header.
Assessment: PASS
Severity: LOW
Notes: No body validation; but header suffices for claim.

## Finding 10

Location: internal/retry/backoff_test.go:8-33
ClaimedBehavior: Full/Equal/No Jitter bounds.
ObservedImplementation: Loop attempts 0..9; checks sleep ∈ [0, cap] for Full/Equal; NoJitter ∈ [Base, Cap].
Assessment: PASS
Severity: LOW
Notes: Bounds correct; prevSleep fixed at 0 for first three; does not test Decorrelated with varying prev.

## Finding 11

Location: internal/retry/backoff_test.go:35-48
ClaimedBehavior: DecorrelatedJitter bounds.
ObservedImplementation: cfg.Base=50ms, Cap=500ms; loop 5 iterations; sets prev=cfg.Base; sleep=DecorrelatedJitter(i,cfg,prev); checks ∈ [Base, Cap]; updates prev=sleep.
Assessment: PASS
Severity: LOW
Notes: Verifies recursive bound property.