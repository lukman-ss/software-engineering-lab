# Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Finding 1

Location: `internal/ratelimit/bucket.go:8-46` (`TokenBucket`)
Claimed Behavior: Thread-safe token bucket allowing bursts up to capacity $B$ and continuously replenishing at rate $R$.
Observed Implementation: Uses `sync.Mutex` for mutual exclusion across all state modifications. Refill is dynamically computed using `time.Now().Sub(tb.lastRefill).Seconds() * tb.refillRate` and capped at `tb.capacity`.
Assessment: PASS
Severity: LOW
Notes: Implementation handles floating-point tokens correctly. When tokens are insufficient, state remains unchanged and `false` is returned immediately.

## Finding 2

Location: `internal/ratelimit/bucket.go:55-79` (`RetryAfterSeconds`)
Claimed Behavior: Calculates rounded integer seconds caller should wait until $n$ tokens are refilled.
Observed Implementation: Evaluates current token balance with refill progression up to current invocation, computes `needed / refillRate`, and ceils fractional seconds to integer.
Assessment: PASS
Severity: LOW
Notes: Correctly handles case where bucket already has sufficient tokens (returns 0).

## Finding 3

Location: `internal/ratelimit/bucket.go:81-121` (`LeakyBucket`)
Claimed Behavior: Thread-safe leaky bucket enforcing output smoothing at leak rate $R$, rejecting bursts exceeding water capacity.
Observed Implementation: Uses `sync.Mutex`. Drains accumulated water over elapsed time `now.Sub(lb.lastLeak).Seconds() * lb.leakRate`, floors water at 0. Rejects when `lb.water + 1.0 > lb.capacity`.
Assessment: PASS
Severity: LOW
Notes: State updates only when request is allowed (`lb.water += 1.0`), preventing unearned water accumulation on rejected requests.

## Finding 4

Location: `internal/ratelimit/registry.go:21-37` (`Registry.Get`)
Claimed Behavior: Multi-tenant token bucket registry providing isolated quotas per tenant key, guarding against CGNAT IP sharing (RFC 6598).
Observed Implementation: Implements double-checked locking with `sync.RWMutex` (`RLock` first, followed by write `Lock` if bucket missing). Returns dedicated `*TokenBucket` per key.
Assessment: PASS
Severity: LOW
Notes: Thread-safe against concurrent initialization of the same tenant key.

## Finding 5

Location: `internal/backpressure/queue.go:30-94` (`BoundedQueue`)
Claimed Behavior: Bounded worker queue with non-blocking load shedding (`ErrQueueFull`) and clean lifecycle management.
Observed Implementation: Channel buffer sized to `capacity`. `TrySubmit` uses `select` with `default:` returning `ErrQueueFull` instantly without blocking. `Stop()` uses `atomic.Bool.CompareAndSwap` ensuring single execution, cancels worker context, closes channel, and awaits worker goroutines via `sync.WaitGroup`.
Assessment: PASS
Severity: LOW
Notes: Workers check context cancellation and channel closure safely.

## Finding 6

Location: `internal/retry/backoff.go:24-63` (`ComputeBackoff`)
Claimed Behavior: Implements AWS Marc Brooker exponential backoff algorithms (NoJitter, FullJitter, EqualJitter, DecorrelatedJitter).
Observed Implementation: Correct mathematical formulas implemented:
- Exponential base: `min(cap, base * 2^attempt)`
- Full Jitter: `rand.Float64() * temp`
- Equal Jitter: `temp/2 + rand.Float64()*(temp/2)`
- Decorrelated Jitter: `min(cap, rand(base, prevSleep * 3))`
Assessment: PASS
Severity: LOW
Notes: Uses global `math/rand` which in Go 1.22+ is thread-safe and auto-seeded.

## Finding 7

Location: `internal/httputil/middleware.go:17-39` (`RateLimitMiddleware`)
Claimed Behavior: HTTP middleware enforcing RFC 6585 status 429, setting `Retry-After` header and returning JSON error body.
Observed Implementation: Reads `X-API-Key` header with fallback to `"anonymous"`. When rate-limited, sets `Retry-After` header with seconds string, writes HTTP 429 header, and encodes JSON `{ "error": "rate_limit_exceeded", "retry_after": ... }`.
Assessment: PASS
Severity: LOW
Notes: Complies with standard HTTP semantics and RFC 6585.
