# Engineering Code Audit

## Summary

Six source files reviewed. Build succeeds, `go vet` clean, race detector clean for the existing test suite. A latent send-on-closed-channel panic exists under concurrent shutdown that no current test exercises (the race detector is blind because there is no memory-safety race, only a panic).

## Finding 1

Location: internal/backpressure/queue.go:66-81 (TrySubmit) and queue.go:87-94 (Stop)
Claimed Behavior: BoundedQueue is concurrency-safe including shutdown (design doc success criteria #5: "Concurrency test passes with Go race detector clean"). TrySubmit rejects fast via select-default; Stop drains workers.
Observed Implementation: Stop() does bq.stopped.CompareAndSwap -> bq.cancel() -> close(bq.queue) -> wg.Wait(). TrySubmit checks stopped.Load() then runs a select with cases `<-bq.ctx.Done()`, `bq.queue <- job`, and `default`. Between the stopped.Load()==false check and the select, Stop() can cancel context AND close the channel. The select then has two ready cases: ctx.Done (returns ErrQueueStopped) and the send on a closed channel (panics). Go selects randomly, so ~50% of attempts hitting that window panic with "send on closed channel".
Assessment: FAIL
Severity: HIGH
Notes: Reproduced with a throwaway concurrent Stop+Submit test (50/50 panic rate, consistent across 3 runs); file removed afterwards to respect no-modification rule. Existing tests do not cover this path: TestBoundedQueue_ConcurrencySafety never calls Stop concurrently, and SubmitAfterStop submits strictly after Stop completes. The race detector does NOT catch this (no data race; it is a logic/ordering panic). This undermines the documented concurrency-safety claim. Fix direction: re-check a draining flag inside the select or guard the send against the closed channel, e.g. perform the stopped/ctx check as the first select case only after ctx is done, or use a mutex around submit+close. Also note: `default` in the select means a send to a closed-but-not-necessarily-full channel is chosen and panics rather than falling through.

## Finding 2

Location: internal/backpressure/queue.go:48-62 (workerLoop) and queue.go:87-94 (Stop)
Claimed Behavior: Graceful shutdown drains pending jobs.
Observed Implementation: Stop() calls close(bq.queue) immediately after cancel(). workers blocked on `<bq.queue` will finish the in-flight job then drain remaining buffered jobs until the channel is empty and returns (ok=false). This drains buffered jobs, OK.
Assessment: PASS (for buffer drain)
Severity: LOW
Notes: Combined with Finding 1, shutdown drains remaining buffered jobs; the only defect is the concurrent-submit path.

## Finding 3

Location: internal/ratelimit/bucket.go (TokenBucket and LeakyBucket)
Claimed Behavior: Thread-safe token/leaky buckets; TokenBucket allows bursts up to capacity with continuous refill; LeakyBucket smooths at leak rate.
Observed Implementation: Both use sync.Mutex around state. TokenBucket.AllowN refills proportionally (elapsed * refillRate), caps at capacity, subtracts n if tokens >= n. LeakyBucket.Allow drains proportionally (subtract elapsed*leakRate), floors at 0, accepts if water+1 <= capacity. Both correct; monotonic time via time.Now() used (standard, not a monotonic-clock API, but consistent within a lock).
Assessment: PASS
Severity: LOW
Notes: `time.Now()` is wall-clock, not `time.Since`/monotonic; design decision #2 claims monotonic guards against clock skew, but Go's `time.Now()` already returns a time with monotonic component when subtracted via Sub/Elapsed. Elapsed is computed via now.Sub(lastRefill) which DOES use the monotonic component, so the skew claim is effectively honored. Minor wording drift only.

## Finding 4

Location: internal/ratelimit/bucket.go:54-79 (RetryAfterSeconds)
Claimed Behavior: Returns seconds a caller should wait for n tokens.
Observed Implementation: Recomputes current token level, returns 0 if >= n; else needed/refillRate, ceiled via int cast + up-round, floored to 1. Correctly computed.
Assessment: PASS
Severity: LOW

## Finding 5

Location: internal/retry/backoff.go (ComputeBackoff)
Claimed Behavior: AWS jitter formulas: NoJitter=cap-or-exp, FullJitter=random(0..min(cap,exp)), EqualJitter=half+random(0..half), DecorrelatedJitter=min(cap, random between base and prev*3).
Observed Implementation: Matches spec. expBackoff = base*2^attempt; temp=min(cap, expBackoff). FullJitter = rand*temp. EqualJitter = half + rand*half (so in [half, temp]). Decorrelated: prevFloat=max(prev,base); rangeMax=prevFloat*3; sleep=base + rand*(rangeMax-base), capped at cap. Correct per AWS/Marc Brooker.
Assessment: PASS
Severity: LOW
Notes: Uses math/rand global; fine for this use. EqualJitter lower bound [half,temp] correct.

## Finding 6

Location: internal/httputil/middleware.go
Claimed Behavior: RFC 6585 429 with Retry-After header and JSON body identifying tenant via X-API-Key (fallback anonymous).
Observed Implementation: tenantKey from X-API-Key or "anonymous"; bucket.Allow() else writes Content-Type+json, Retry-After header (strconv.Itoa(retryAfter)), 429, JSON body {error,RetryAfter}. Correct.
Assessment: PASS
Severity: LOW

## Finding 7

Location: cmd/demo/main.go
Claimed Behavior: Demonstrates token bucket burst, leaky bucket smoothing, bounded queue shedding, jitter spread.
Observed Implementation: All four sections implemented; output non-deterministic (random jitter; timing-dependent queue accept/reject). See test-audit #3 and docs-vs-code #2.
Assessment: PASS (functionally) / WARNING (reproducibility)
Severity: MEDIUM
Notes: The demo itself is real and runs cleanly; issue is over-claimed deterministic output in engineering record.
