# Code Audit

## Finding 1

Location: internal/ratelimit/registry.go:21-36 (Get method)
Claimed Behavior: per-tenant TokenBucket isolation with double-checked locking.
Observed Implementation: Double-checked map lookup with RLock then Lock, creates bucket on miss.
Assessment: PASS
Severity: -
Notes: Correct double-check prevents lost-create race; no redundant allocations. RWMutex used appropriately.

## Finding 2

Location: internal/ratelimit/bucket.go:54-79 (RetryAfterSeconds)
Claimed Behavior: Computes integer seconds to wait for n tokens.
Observed Implementation: Tokens+elapsed*rate, capped, then (n-tokens)/rate ceiling.
Edge: refillRate=0 → division by zero → panic on int conversion.
Assessment: WARNING
Severity: MEDIUM
Notes: refillRate is float64; zero makes denominator zero → needed/infinite → 0 → returns 1 due to `if secs <= 0: return 1`. Actually safe by coincidence but unclear contract; guard clause recommended.

## Finding 3

Location: internal/backpressure/queue.go:60-70 (TrySubmit)
Claimed Behavior: Non-blocking enqueue; fast drop on full.
Observed Implementation: Select default on channel send; accepted/rejected counters atomic.
Assessment: PASS
Severity: -
Notes: Exactly the Go idiom for non-blocking offer; counters accurate.

## Finding 4

Location: internal/backpressure/queue.go:76-80 (Stop)
Claimed Behavior: Cancel ctx, close queue, wait wg.
Observed Implementation: Cancel then close then Wait. window: worker may see closed channel after ctx.Done() but before consuming; job dropped silently (select default path not used in workerLoop).
Assessment: WARNING
Severity: LOW
Notes: workerLoop selects on ctx first, so ctx cancellation observed before queue close; safe but odd order (cancel before close). No panic risk.

## Finding 5

Location: internal/backpressure/queue.go:14-24 (struct)
Claimed Behavior: Fields for stats and coordination.
Observed Implementation: accepted/rejected/processed atomic.Int64; ctx cancel; wg.
Assessment: PASS
Severity: -
Notes: Correct sync primitives; no mutex needed for counters.

## Finding 6

Location: internal/ratelimit/bucket.go:25-46 (AllowN)
Claimed Behavior: Thread-safe token consumption with refill.
Observed Implementation: Mutex lock; refill based on elapsed; cap; consume; unlock.
Assessment: PASS
Severity: -
Notes: Standard token bucket; monotonic time.Now().

## Finding 7

Location: internal/retry/backoff.go:50-58 (DecorrelatedJitter)
Claimed Behavior: sleep = min(cap, random_between(base, prevSleep * 3)).
Observed Implementation: baseFloat + rand*(rangeMax-baseFloat) where rangeMax = prevFloat*3, clamped to base; matches formula.
Assessment: PASS
Severity: -
Notes: Correct DecorrelatedJitter per AWS/Brooker.

## Finding 8

Location: internal/httputil/middleware.go:17-40
Claimed Behavior: RFC 6585 429 + Retry-After + JSON body on rate limit.
Observed Implementation: Header set, StatusTooManyRequests, JSON encode.
Assessment: PASS
Severity: -
Notes: Compliant minimal middleware.

## Finding 9

Location: internal/ratelimit/bucket_test.go:83-98 (ConcurrencyRace)
Claimed Behavior: Detect races under 50*10 Allow calls.
Observed Implementation: Spin 50 goroutines each calling Allow 10x; no assertions on outcome.
Assessment: WARNING
Severity: MEDIUM
Notes: Race detector validates sync; but test lacks behavioral assertion (e.g., total allowed <= capacity + refill*duration). Weak coverage.

## Finding 10

Location: cmd/demo/main.go:13-65
Claimed Behavior: Demo all four subsystems.
Observed Implementation: Prints token bucket, leaky bucket, bounded queue, retry jitter.
Missing: DecorrelatedJitter not shown; tenant key not varied; HTTP 429 not exercised.
Assessment: WARNING
Severity: MEDIUM
Notes: Demo incomplete per claims; omits DecorrelatedJitter and middleware path.