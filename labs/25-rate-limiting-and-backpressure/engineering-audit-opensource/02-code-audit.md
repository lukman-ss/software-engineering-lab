# Code Audit

Target Lab: labs/25-rate-limiting-and-backpressure
Date: 2026-09-27
Files reviewed: bucket.go, registry.go, queue.go, backoff.go, middleware.go, cmd/demo/main.go

## Finding 1

Location: internal/ratelimit/bucket.go:29-46 (TokenBucket.AllowN)
Claimed Behavior: Burst up to capacity B, then enforce refill rate R.
Observed Implementation: Elapsed-time refill with cap clamp, then deduct. Mutex-guarded.
Assessment: PASS
Severity: LOW
Notes: Fractional arithmetic correct. Tokens() returns last-stored value without live refill; stale by microseconds. Demo display only.

## Finding 2

Location: internal/ratelimit/bucket.go:55-79 (RetryAfterSeconds)
Claimed Behavior: Caller wait hint for Retry-After header.
Observed Implementation: Read-only recompute, ceil((n-tokens)/rate), min 1 when positive deficit, 0 when satisfied.
Assessment: PASS
Severity: LOW
Notes: No state mutation, safe. Division by zero if refillRate=0 (needed/0 = +Inf, int conversion undefined). No zero-rate guard. Constructor accepts any float, no validation.

## Finding 3

Location: internal/ratelimit/bucket.go:98-115 (LeakyBucket.Allow)
Claimed Behavior: Constant leak rate R, reject when water reaches capacity.
Observed Implementation: water -= elapsed*leakRate floored at 0, admit if water+1 <= capacity.
Assessment: PASS
Severity: LOW
Notes: Boundary exact-capacity admit correct. Water() same staleness note as Tokens().

## Finding 4

Location: internal/ratelimit/registry.go:21-37 (Registry.Get)
Claimed Behavior: Per-tenant isolation avoiding CGNAT IP collision.
Observed Implementation: RWMutex double-checked lazy init, one TokenBucket per key.
Assessment: PASS
Severity: LOW
Notes: Correct. No eviction/TTL; unbounded distinct-tenant growth. Documented as in-memory; cardinality leak not documented.

## Finding 5

Location: internal/backpressure/queue.go:66-81 (TrySubmit)
Claimed Behavior: Non-blocking fast ErrQueueFull when full.
Observed Implementation: select send with default branch, atomic accepted/rejected counters.
Assessment: PASS
Severity: LOW
Notes: Genuinely non-blocking. Correct fast-shed semantics.

## Finding 6

Location: internal/backpressure/queue.go:87-94 (Stop) vs TrySubmit send
Claimed Behavior: Safe worker-pool shutdown.
Observed Implementation: Stop CAS stopped flag, cancel, close(queue), wg.Wait. In-flight TrySubmit that passed stopped check can select send on channel concurrently with close -> panic (send on closed channel).
Assessment: WARNING
Severity: MEDIUM
Notes: Tests only Stop-then-submit (safe early return). No concurrent Stop-during-Submit test. Safe under defer-Stop-after-submits pattern used in demo/tests. Fix: remove close(queue) and rely on cancel, or guard submit/close with mutex. RACE_CONDITION / MISSING_EDGE_CASE.

## Finding 7

Location: internal/backpressure/queue.go:48-62 (workerLoop)
Claimed Behavior: Worker pool consumer.
Observed Implementation: Selects ctx.Done vs queue receive. Job error discarded (_ = job).
Assessment: PASS
Severity: LOW
Notes: Buffered jobs abandoned on Stop (cancel wins over drain). Undocumented drop-on-shutdown. Job error propagation absent; no claim requires it.

## Finding 8

Location: internal/retry/backoff.go:24-62 (ComputeBackoff)
Claimed Behavior: AWS formulas: NoJitter=temp, FullJitter=U(0,temp), EqualJitter=temp/2+U(0,temp/2), Decorrelated=min(cap,U(base,prev*3)).
Observed Implementation: Matches formulas. temp=min(cap,base*2^attempt).
Assessment: PASS
Severity: LOW
Notes: Uses math/rand global (auto-seeded Go 1.20+, concurrency-safe via locked source). Unknown strategy falls back to temp. Negative attempt yields fractional temp; untested, harmless.

## Finding 9

Location: internal/httputil/middleware.go:17-40
Claimed Behavior: RFC 6585 429 with Retry-After on exhaustion.
Observed Implementation: X-API-Key tenant, anonymous fallback, bucket.Allow gate, JSON body + Retry-After header + 429.
Assessment: PASS
Severity: LOW
Notes: Retry-After always >=1 in normal configs. Content-Type application/json correct. No 503 mapping for queue-full; queue error never reaches HTTP layer.

## Finding 10

Location: cmd/demo/main.go
Claimed Behavior: Demo shows burst, smoothing, shedding, jitter spread.
Observed Implementation: Real calls into all four packages. Output verified live.
Assessment: WARNING
Severity: MEDIUM
Notes: Section 3 stats nondeterministic (observed Accepted=4/Rejected=2 vs recorded 3/3) due to worker consuming during submit loop. Section 4 jitter values random per run. Not fake, but recorded snapshot in 03-execution-result.md not exactly reproducible.
