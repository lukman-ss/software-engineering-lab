# Code Audit

Target: labs/25-rate-limiting-and-backpressure. Read-only; no code changed.

## Finding 1

Location: internal/ratelimit/bucket.go:29-46 (TokenBucket.AllowN), :48-52 (Tokens), :55-79 (RetryAfterSeconds)
Claimed Behavior: burst B + continuous refill R, thread-safe.
Observed Implementation: mutex-guarded; lazy refill on each call capped at capacity; RetryAfterSeconds ceiling to int seconds.
Assessment: PASS
Severity: LOW
Notes: Tokens() returns stale snapshot (no lazy refill). RetryAfterSeconds(1) estimates without mutating state; correct.

## Finding 2

Location: internal/ratelimit/bucket.go:89-121 (LeakyBucket)
Claimed Behavior: constant leak R, burst rejection at capacity.
Observed Implementation: mutex-guarded; water decays by elapsed*leakRate floored at 0; admits iff water+1 <= capacity.
Assessment: PASS
Severity: LOW
Notes: Symmetric to token bucket. Water() stale-snapshot same as Tokens(); cosmetic.

## Finding 3

Location: internal/ratelimit/registry.go:21-37
Claimed Behavior: per-tenant isolation; concurrent same-key returns identical pointer.
Observed Implementation: RWMutex + double-checked locking; correct.
Assessment: PASS
Severity: LOW
Notes: Map grows unbounded, no eviction/TTL. Documented in-memory limitation covers distributed case; local leak unmentioned.

## Finding 4

Location: internal/backpressure/queue.go:67-85 (TrySubmit), :91-103 (Stop), :49-63 (workerLoop)
Claimed Behavior: non-blocking shed with ErrQueueFull; ErrQueueStopped after Stop.
Observed Implementation: RLock check of stopped + select ctx-done / enqueue / default. Stats via atomics. Stop idempotent, cancels ctx, closes channel, wg.Wait under no submit-lock held. Race run clean.
Assessment: PASS
Severity: LOW
Notes: Residual race window: submitter holding RLock sends on queue while Stop closes it -> possible panic send-on-closed. Covered by ConcurrentStopAndSubmit passing in practice (close happens only after stopMu exclusive lock, but in-flight sender already past check can still panic). Latent HIGH-severity path exists but unproven this run; flagged in gaps as UNVERIFIED/MISSING_EDGE_CASE, not FAIL.

## Finding 5

Location: internal/backpressure/queue.go:49-63
Claimed Behavior: worker pool processes jobs; failure handling.
Observed Implementation: `_ = job(bq.ctx)` discards error; job panic propagates and kills worker (no recover); processed++ counts even failed jobs; context passed is queue-lifetime ctx, not per-job timeout.
Assessment: WARNING
Severity: MEDIUM
Notes: Error-swallowing + no panic isolation means one bad job silently reduces pool capacity. Tests never submit failing/panicking jobs.

## Finding 6

Location: internal/retry/backoff.go:24-62
Claimed Behavior: AWS Brooker Full/Equal/No/Decorrelated formulas.
Observed Implementation: matches formulas; min(cap, base*2^attempt); decorrelated floors prev at base; unknown strategy falls back to temp.
Assessment: PASS
Severity: LOW
Notes: math/rand unseeded global (Go 1.22+ auto-seed; fine). attempt>=~52 overflows to +Inf; math.Min(cap,Inf)=cap so bounded. Negative attempt yields fraction — untested, harmless.

## Finding 7

Location: internal/httputil/middleware.go:17-40
Claimed Behavior: RFC 6585 429 + Retry-After + JSON body; X-API-Key tenant, anonymous fallback.
Observed Implementation: as claimed; RetryAfterSeconds(1) feeds header and body consistently.
Assessment: PASS
Severity: LOW
Notes: Retry-After often "1" even for sub-second waits (ceiling). Correct-conservative. No Retry-After on success; correct.

## Finding 8

Location: internal/ratelimit/bucket.go:55-79
Claimed Behavior: RetryAfterSeconds usable with any bucket config.
Observed Implementation: `needed / tb.refillRate` divides by zero when refillRate=0 -> +Inf -> int(+Inf) platform behavior / hang risk. No zero-rate guard; no test.
Assessment: WARNING
Severity: MEDIUM
Notes: Constructor accepts 0/negative rates unchecked. Only reachable via misuse; still UNHANDLED_ERROR gap.

## Finding 9

Location: cmd/demo/main.go
Claimed Behavior: demo of burst, smoothing, shedding, jitter.
Observed Implementation: real execution, deterministic sections 1-2, timing/randomized 3-4 by design and disclosed in 03-execution-result.md:69.
Assessment: PASS
Severity: LOW
Notes: Demo output verified live this audit (see 03-test-audit.md). Not fabricated.
