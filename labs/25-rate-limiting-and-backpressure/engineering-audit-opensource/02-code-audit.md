# Code Audit

## Finding 1: TokenBucket Implementation

Location: internal/ratelimit/bucket.go (lines 8-52)
Claimed Behavior: Thread-safe token bucket with fractional token replenishment, allowing bursts up to capacity and refilling at continuous rate.
Observed Implementation: 
- Uses `sync.Mutex` for thread safety
- `NewTokenBucket(capacity, refillRate)` initializes with full tokens
- `AllowN(n)` correctly refills tokens based on elapsed time, caps at capacity, and checks if tokens >= n
- `RetryAfterSeconds(n)` calculates wait time for n tokens with proper rounding
Assessment: PASS
Severity: N/A
Notes: Implementation matches specification. Concurrency safe with mutex guarding all shared state.

## Finding 2: LeakyBucket Implementation

Location: internal/ratelimit/bucket.go (lines 81-121)
Claimed Behavior: Thread-safe leaky bucket with constant leak rate, smoothing flow to rate R.
Observed Implementation: 
- Uses `sync.Mutex` for thread safety
- `NewLeakyBucket(capacity, leakRate)` initializes with 0 water
- `Allow()` drains water based on elapsed time, adds 1.0 if water+1 <= capacity
- `Water()` thread-safe read
Assessment: PASS
Severity: N/A
Notes: Implementation matches specification. Concurrency safe with mutex guarding all shared state.

## Finding 3: Registry Implementation

Location: internal/ratelimit/registry.go (lines 1-37)
Claimed Behavior: Per-tenant rate limiter registry to avoid CGNAT IP collisions (RFC 6598).
Observed Implementation: 
- Uses `sync.RWMutex` for thread safety
- `Get(tenantKey)` implements double-checked locking pattern for efficient concurrent access
- Creates new TokenBucket per tenant key on first access
Assessment: PASS
Severity: N/A
Notes: Implementation matches specification. Thread-safe with correct double-checked locking.

## Finding 4: BoundedQueue Implementation

Location: internal/backpressure/queue.go (lines 1-80)
Claimed Behavior: Fixed-capacity worker pool rejecting jobs fast when queue capacity is reached (non-blocking backpressure).
Observed Implementation: 
- Uses `sync/atomic.Int64` for counters and `sync.WaitGroup` for worker management
- `NewBoundedQueue(capacity, workers)` creates buffered channel and spawns workers
- `TrySubmit(job)` uses select-default pattern for non-blocking enqueue
- Returns `ErrQueueFull` when queue capacity exceeded
- `Stats()` returns accepted, rejected, processed counts and queue length
- `Stop()` cancels context, closes channel, waits for workers
Assessment: WARNING
Severity: MEDIUM
Notes: 
- Implementation is correct for normal usage patterns
- However, `Stop()` closes the channel, and subsequent calls to `TrySubmit()` will panic with "send on closed channel"
- This is an API contract issue rather than a bug in intended usage
- Demo and tests use `defer bq.Stop()` after all operations, so not triggered
- Recommendation: Document that `TrySubmit()` must not be called after `Stop()`, or improve error handling

## Finding 5: HTTP Middleware Implementation

Location: internal/httputil/middleware.go (lines 1-40)
Claimed Behavior: HTTP middleware enforcing tenant-based rate limits, returning RFC 6585 429 Too Many Requests with Retry-After header.
Observed Implementation: 
- Extracts tenant key from X-API-Key header (defaults to "anonymous")
- Gets bucket from registry via `Get(tenantKey)`
- If `!bucket.Allow()`, returns HTTP 429 with Retry-After header and JSON error body
- Otherwise, calls next handler
Assessment: PASS
Severity: N/A
Notes: Implementation matches specification. Correctly implements RFC 6585 semantics.

## Finding 6: Retry/Backoff Implementation

Location: internal/retry/backoff.go (lines 1-63)
Claimed Behavior: Exponential backoff implementation supporting AWS Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter.
Observed Implementation: 
- Defines BackoffStrategy constants matching AWS specifications
- `ComputeBackoff(strategy, attempt, cfg, prevSleep)` implements all four jitter algorithms:
  - NoJitter: `min(cap, base * 2^attempt)` ✓
  - FullJitter: `random(0, min(cap, base * 2^attempt))` ✓
  - EqualJitter: `min/2 + random(0, min/2)` ✓
  - DecorrelatedJitter: `random_between(base, prevSleep * 3)` clamped to cap ✓
- Uses `math/rand` which is thread-safe in Go
Assessment: PASS
Severity: N/A
Notes: Implementation matches AWS specifications exactly. Mathematical bounds are correct per verification in tests.

## Summary of Code Quality

The implementation demonstrates:
- Strong thread safety practices with appropriate synchronization primitives
- Correct algorithmic implementations matching published specifications
- Clear separation of concerns across packages
- Proper error handling with meaningful error types (e.g., `ErrQueueFull`)
- Clean, readable code with appropriate comments
- Effective use of Go idioms (select-default for non-blocking channels, context for cancellation)