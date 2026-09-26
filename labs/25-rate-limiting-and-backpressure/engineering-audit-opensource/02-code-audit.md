# Engineering Code Audit

## Finding 1: TokenBucket Burst and Refill Logic

Location: internal/ratelimit/bucket.go:29-46

Claimed Behavior: Token bucket allows bursts up to token capacity B, then enforces refill rate R. Refills continuously at rate R.

Observed Implementation:
The `AllowN` method acquires a mutex, computes elapsed time since last refill, adds `elapsed * refillRate` tokens (capped at capacity), then checks if sufficient tokens remain for the requested amount. The refill is lazy (computed on each access using `time.Now().Sub()`).

Assessment: PASS
Severity: LOW
Notes: Implementation is correct. Mutex guards all state access. Refill calculation uses elapsed time properly. Token cap at capacity prevents over-refill.

---

## Finding 2: TokenBucket RetryAfterSeconds Calculation

Location: internal/ratelimit/bucket.go:55-79

Claimed Behavior: Calculates how long caller should wait for n tokens.

Observed Implementation:
Computes current token count (including refill since lastRefill), checks if already sufficient (returns 0), otherwise computes `needed / refillRate` and rounds up to nearest integer.

Assessment: PASS
Severity: LOW
Notes: Logic is sound. Returns 0 when tokens available. Rounds up to ensure sufficient tokens. Minimum return of 1 second when needed > 0 but rounded value is 0.

---

## Finding 3: LeakyBucket Drain Logic

Location: internal/ratelimit/bucket.go:81-115

Claimed Behavior: Enforces constant drain rate R, rejecting bursts when water reaches capacity.

Observed Implementation:
The `Allow` method drains `elapsed * leakRate` from water (floored at 0), then checks if `water + 1.0 <= capacity`. Adds 1.0 to water when allowed.

Assessment: PASS
Severity: LOW
Notes: Correct leaky bucket implementation. Water level is subtracted (drained) over time. New request adds 1.0 unit. Rejects when `water + 1.0 > capacity`. Mutex guards state.

---

## Finding 4: Registry Tenant Isolation

Location: internal/ratelimit/registry.go

Claimed Behavior: Per-tenant rate limiters to avoid CGNAT IP collisions (RFC 6598).

Observed Implementation:
Uses `sync.RWMutex`. Read lock checks if bucket exists; if so returns. Write lock (with double-check) creates new bucket if not exists.

Assessment: PASS
Severity: LOW
Notes: Double-checked locking pattern is correctly implemented. Tenant keys are used as map keys. Memory leak risk: buckets are never cleaned up, but this is acceptable for in-memory demo scope.

---

## Finding 5: BoundedQueue Fast Rejection

Location: internal/backpressure/queue.go:60-70

Claimed Behavior: TrySubmit enqueues if capacity allows, otherwise drops immediately with ErrQueueFull.

Observed Implementation:
Uses Go select-default channel pattern: attempts `bq.queue <- job` and returns nil on success, or falls to `default` case immediately returning ErrQueueFull.

Assessment: PASS
Severity: LOW
Notes: Non-blocking submission is correct. The select-default pattern is the idiomatic Go way to do non-blocking channel sends. `accepted` and `rejected` atomic counters are updated appropriately.

---

## Finding 6: BoundedQueue Worker Lifecycle

Location: internal/backpressure/queue.go:26-58

Claimed Behavior: Worker pool processes jobs from bounded channel.

Observed Implementation:
Creates `workers` goroutines, each running `workerLoop`. Worker loop selects on `ctx.Done()` (returns on shutdown) or job from queue. Processes job with `bq.ctx`, increments `processed` counter.

Assessment: WARNING
Severity: MEDIUM
Notes: Potential issue: when `Stop()` is called, `cancel()` is invoked followed by `close(bq.queue)`. Workers in `workerLoop` may still be processing a job (blocked on `job(bq.ctx)`) when the queue is closed. The close of the channel will not affect active job execution. After job completes, `processed.Add(1)` executes on a possibly-stopped queue, but this is just an atomic increment which is safe. However, `Stop()` waits for workers via `bq.wg.Wait()`, so any in-flight job will complete before Stop returns. This is acceptable behavior.

However, there is a subtle issue: after `Stop()` is called, if someone calls `TrySubmit`, it will attempt to send on a closed channel, causing a panic. The implementation does not guard against submissions after Stop.

---

## Finding 7: AWS Jitter Backoff Formulas

Location: internal/retry/backoff.go:24-63

Claimed Behavior: Implements Full Jitter, Equal Jitter, and Decorrelated Jitter matching AWS Marc Brooker formulas.

Observed Implementation:
- NoJitter: `min(cap, base * 2^attempt)` — correct
- FullJitter: `random_between(0, min(cap, base * 2^attempt))` — matches AWS formula `random(0,1) * min(cap, base*2^attempt)`
- EqualJitter: `temp/2 + random_between(0, temp/2)` — matches AWS Equal Jitter
- DecorrelatedJitter: `min(cap, random_between(base, prevSleep * 3))` — matches AWS Decorrelated Jitter

Assessment: PASS
Severity: LOW
Notes: All four strategies correctly implement the AWS jitter formulas. Uses `math/rand` (not crypto/rand, which is appropriate for this use case). No thread-safety concern since functions are stateless (no shared mutable state).

---

## Finding 8: HTTP 429 Middleware

Location: internal/httputil/middleware.go:17-40

Claimed Behavior: Returns HTTP 429 Too Many Requests (RFC 6585) with Retry-After header when rate limit exceeded.

Observed Implementation:
Extracts tenant key from `X-API-Key` header (falls back to "anonymous"). Gets per-tenant bucket from registry. If `Allow()` returns false, sets retry-after, writes 429 status with JSON error body. Otherwise delegates to next handler.

Assessment: PASS
Severity: LOW
Notes: Correctly implements RFC 6585 429 response. Properly sets Content-Type, Retry-After header, and JSON body. Response is written before returning, preventing further writes.

---

## Finding 9: Concurrency Safety in TokenBucket

Location: internal/ratelimit/bucket.go:8-14

Claimed Behavior: Thread-safe under concurrent access.

Observed Implementation:
Uses `sync.Mutex` on all methods that access `tokens`, `lastRefill`, and `refillRate`. All state mutations occur under lock.

Assessment: PASS
Severity: LOW
Notes: Mutex properly guards all internal state. No race conditions possible within the TokenBucket.

---

## Finding 10: Concurrency Safety in LeakyBucket

Location: internal/ratelimit/bucket.go:81-87

Claimed Behavior: Thread-safe under concurrent access.

Observed Implementation:
Uses `sync.Mutex` on `Allow()` and `Water()` methods. All state access is guarded.

Assessment: PASS
Severity: LOW
Notes: Same pattern as TokenBucket. Mutex properly guards all internal state.

---

## Finding 11: Atomic Counters in BoundedQueue

Location: internal/backpressure/queue.go:20-24

Claimed Behavior: Thread-safe statistics tracking.

Observed Implementation:
Uses `atomic.Int64` for `accepted`, `rejected`, and `processed` counters. These are incremented under no lock but atomically.

Assessment: PASS
Severity: LOW
Notes: Atomic operations ensure safe concurrent access. `Stats()` method reads these values atomically.

---

## Finding 12: Registry Double-Checked Locking

Location: internal/ratelimit/registry.go:21-37

Claimed Behavior: Thread-safe tenant bucket lookup and creation.

Observed Implementation:
Uses `sync.RWMutex`. Read lock for fast path (existing bucket), write lock with double-check for new bucket creation. This is the correct double-checked locking pattern.

Assessment: PASS
Severity: LOW
Notes: Correct implementation of double-checked locking. No race condition or duplicate creation possible.

---

## Finding 13: Missing BoundedQueue Stop Guard

Location: internal/backpressure/queue.go:60-70

Claimed Behavior: TrySubmit is safe to call at any time.

Observed Implementation:
After `Stop()` closes the queue channel, any subsequent `TrySubmit` call will panic because sending on a closed channel causes a runtime panic.

Assessment: WARNING
Severity: MEDIUM
Notes: There is no guard preventing `TrySubmit` calls after `Stop()`. In production code this could be a serious issue. However, in the demo and tests, `TrySubmit` is only called before `Stop()`. This is a design gap but not an active bug in the current usage. The `default` branch of select prevents blocking, but does NOT prevent sending on a closed channel.

---

## Summary

Total Findings: 13
PASS: 10
WARNING: 3
FAIL: 0

All core algorithms are correctly implemented. The main concerns are edge case handling around the bounded queue lifecycle (submissions after Stop) and the fact that the queue lacks a guard against this scenario. These are design gaps rather than active bugs in the current code paths.