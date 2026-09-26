# Code Audit

Target: labs/25-rate-limiting-and-backpressure

## Finding 1

Location: internal/ratelimit/bucket.go:29-46 (TokenBucket.AllowN)
Claimed Behavior: Burst to capacity B, refill at rate R, thread-safe.
Observed Implementation: Lazy refill on each call under sync.Mutex, clamp to capacity. Correct.
Assessment: PASS
Severity: LOW
Notes: Fractional arithmetic correct. No validation of capacity/rate <= 0.

## Finding 2

Location: internal/ratelimit/bucket.go:55-79 (RetryAfterSeconds)
Claimed Behavior: Retry-After hint for 429 responses.
Observed Implementation: Read-only refill estimate, ceil to seconds, 0 when available. Correct. Does not mutate lastRefill.
Assessment: PASS
Severity: LOW
Notes: Div-by-zero if refillRate==0 (needed/0 = +Inf, float->int conversion undefined). Constructor allows it. Edge unhandled.

## Finding 3

Location: internal/ratelimit/bucket.go:48-52 (Tokens)
Claimed Behavior: Remaining tokens reporter (used in demo).
Observed Implementation: Returns stored field without applying pending refill. Stale low by elapsed*R.
Assessment: WARNING
Severity: LOW
Notes: Display-only. Allow() itself refills correctly, so limiting behavior unaffected.

## Finding 4

Location: internal/ratelimit/bucket.go:98-115 (LeakyBucket.Allow)
Claimed Behavior: Constant drain R, reject burst at capacity.
Observed Implementation: Lazy leak under mutex, clamp at 0, admit iff water+1 <= capacity. Correct.
Assessment: PASS
Severity: LOW
Notes: Same zero/negative constructor-input caveat as Finding 2.

## Finding 5

Location: internal/ratelimit/registry.go:21-37
Claimed Behavior: Per-tenant isolation, CGNAT-safe (RFC 6598 rationale).
Observed Implementation: Double-checked locking with RWMutex. Correct, race-clean.
Assessment: PASS
Severity: LOW
Notes: No eviction/TTL — unbounded map growth under unbounded tenant cardinality. Acknowledged in-memory scope; OOM ceiling unmentioned. LOW.

## Finding 6

Location: internal/backpressure/queue.go:61-70 (TrySubmit)
Claimed Behavior: Non-blocking shed with ErrQueueFull.
Observed Implementation: select send/default on buffered chan, atomic counters. Correct, race-clean.
Assessment: PASS
Severity: LOW
Notes: Core backpressure claim proven.

## Finding 7

Location: internal/backpressure/queue.go:76-80 (Stop)
Claimed Behavior: Shutdown.
Observed Implementation: cancel() + close(queue) + wg.Wait(). Not idempotent; concurrent TrySubmit vs Stop sends on closed channel -> panic. Remaining queued jobs dropped, not drained. Job errors swallowed (job(bq.ctx) return ignored by design).
Assessment: WARNING
Severity: MEDIUM
Notes: Safe in tests/demo (Stop after submits settle, deferred once). No lifecycle/concurrency test. Also NewBoundedQueue(capacity<=0/workers<=0) unvalidated: unbuffered chan or undrained queue.

## Finding 8

Location: internal/httputil/middleware.go:17-40
Claimed Behavior: RFC 6585 429 + Retry-After + JSON body, tenant key extraction.
Observed Implementation: X-API-Key fallback "anonymous", bucket.Allow gate, correct status/header/body. Correct.
Assessment: PASS
Severity: LOW
Notes: Anonymous fallback merges all unauthenticated callers into one bucket (intended, matches notes). RetryAfter div-zero caveat inherits Finding 2.

## Finding 9

Location: internal/retry/backoff.go:24-63
Claimed Behavior: AWS/No/Full/Equal/Decorrelated formulas.
Observed Implementation: temp=min(cap,base*2^attempt); Full=rand*temp; Equal=half+rand*half; Decorrelated=min(cap,rand(base,3*max(base,prev))). Matches spec. Global math/rand safe for concurrent use, auto-seeded on Go 1.22.
Assessment: PASS
Severity: LOW
Notes: Bounds-only tests; exact-distribution conformance not asserted (acceptable for lab).

## Finding 10

Location: cmd/demo/main.go
Claimed Behavior: Burst, smoothing, shedding, jitter spread demo.
Observed Implementation: Real execution against same packages. Ran green; deterministic sections byte-identical to recorded log, jitter values vary run-to-run as expected from randomness.
Assessment: PASS
Severity: LOW
Notes: No FAKE_DEMO. Recorded jitter numbers are one sample, not reproducible literals.
