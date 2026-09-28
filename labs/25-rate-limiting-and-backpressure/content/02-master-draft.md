# Rate Limiting & Backpressure: Principles, Algorithms, and Implementation

## Problem

Distributed systems face a fundamental tension: incoming requests arrive faster than downstream components can process them. Without protective mechanisms, this imbalance produces unbounded queue buildup, memory exhaustion, timeout cascades, and ultimately total service outage. Rate limiting and backpressure are complementary strategies that address this tension at different layers: rate limiting throttles traffic at system boundaries, while backpressure propagates pressure signals backward through the internal pipeline to prevent overload propagation.

## Why This Matters

Modern distributed systems—API gateways, microservice meshes, message queues—must serve diverse tenants with varying workload patterns on shared infrastructure. A single noisy tenant can degrade or deny service to all others. The absence of rate limiting invites cascading failures: when one component slows, upstream buffers fill, timeouts multiply, and retry storms amplify the original fault. Research confirms that multi-layer approaches combining rate limiters with load shedders are most effective for production systems (Stripe's four-tier limiter pattern). Understanding these mechanisms is essential for building resilient, predictable systems.

## Mental Model

The core mental model rests on three pillars:

1. **Token Bucket** — A container holding tokens that replenish at a fixed rate. Requests consume tokens. When tokens exist, requests pass; when empty, they are rejected. This allows controlled bursts up to bucket capacity while enforcing a long-term average rate.

2. **Bounded Queue Backpressure** — A fixed-capacity buffer between producers and consumers. When the buffer is full, new jobs are rejected immediately rather than queued indefinitely. This prevents memory exhaustion and signals to callers that the system is overloaded.

3. **Exponential Backoff with Jitter** — When a request fails, the client waits before retrying. The wait grows exponentially with each attempt. Adding randomized jitter prevents synchronized retries (thundering herd), which would otherwise cause periodic load spikes that amplify downstream degradation.

## Core Concept: Token Bucket vs. Leaky Bucket

The **Token Bucket** algorithm allows burst traffic up to bucket capacity while maintaining a long-term average rate limit. The bucket holds tokens that accumulate at a refill rate. Each request consumes one token. If tokens are available, the request is allowed; otherwise it is rejected. The burstiness of allowed traffic is determined by the bucket depth. This makes token bucket ideal for scenarios where occasional bursts are acceptable but sustained rates must be controlled.

The **Leaky Bucket** algorithm smooths traffic output to a constant rate. Water (requests) flows into the bucket; leaks drain at a constant rate. If the bucket fills (water reaches capacity), new arrivals are rejected. Unlike token bucket, leaky bucket enforces strict output pacing regardless of arrival patterns. NGINX implements rate limiting using leaky bucket as a meter via `limit_req_module`.

Mathematically, the leaky bucket as a meter is equivalent to a mirror image of the token bucket algorithm. The practical difference is perspective: token bucket controls input rate with burst tolerance, leaky bucket controls output rate with smoothing.

## Failure Scenario

Consider a system receiving 5,000,000 items at a processing rate of 2,000 items per second. Without rate limiting or backpressure, the queue grows continuously. Using Little's Law (L = λW), the time to drain is 5,000,000 / 2,000 = 2,500 seconds ≈ 41 minutes 40 seconds. During this period, memory accumulates unboundedly, latency degrades, and timeouts cascade upstream.

A synchronized retry storm compounds the failure. Without jitter, 100 contending clients retry simultaneously at fixed intervals, creating periodic load spikes that more than double the effective request count. The system never recovers because each recovery attempt triggers the next wave of failures.

Naive IP-based rate limiting introduces a different failure: carrier-grade NAT (RFC 6598) routes multiple tenants through a single IP address. A single noisy neighbor throttles all users behind that IP. Per-tenant key-based limiting (using API keys or tenant IDs) avoids this degradation.

## How It Works: The Architecture

The system architecture proceeds through four layers:

```
[Incoming Request]
        │
        ▼
[HTTP Middleware / Tenant Extractor] ── (RFC 6598 Tenant Key via X-API-Key)
        │
        ▼
[Token Bucket / Rate Limiter] ── (Exceeded? → HTTP 429 + Retry-After)
        │ (Passed)
        ▼
[Bounded Queue / Backpressure Channel] ── (Queue Full? → 503 Overloaded)
        │ (Enqueued)
        ▼
[Worker Pool / Consumer] ── (Little's Law L = λW steady state)
        │
        ▼
[Client with Full Jitter Retry] ── (AWS Jitter Backoff)
```

The HTTP middleware extracts the tenant identity from the `X-API-Key` header (falling back to `anonymous`), consults the per-tenant token bucket registry, and either allows the request to proceed or returns HTTP 429 with a `Retry-After` header. The bounded queue accepts jobs up to capacity, rejecting excess immediately. The worker pool processes enqueued jobs at a constant rate. The client implements exponential backoff with jitter to avoid synchronization.

## Implementation

The Go implementation uses only the standard library (`sync`, `time`, `net/http`) to maximize clarity and testability. Key design decisions include: monotonic time tracking via `time.Now()` to guard against system clock adjustments; floating-point token arithmetic for smooth fractional refill precision; non-blocking channel submission (`select-default`) for zero-allocation fast rejection; mutex synchronization chosen over lock-free CAS for simplicity and robust multi-field updates.

The `TokenBucket` struct maintains capacity, current tokens, refill rate, and last-refill timestamp. `Allow()` replenishes tokens based on elapsed time, caps at capacity, and consumes one token if available. `RetryAfterSeconds()` calculates wait time for a requested number of tokens. The `LeakyBucket` struct mirrors this with water level tracking and constant leak rate.

The `Registry` manages per-tenant token buckets keyed by API key/tenant ID, preventing CGNAT-based throttling collisions. `BoundedQueue` uses a fixed-capacity Go channel with `TrySubmit` employing a `select-default` pattern: if the channel send succeeds, the job is accepted; otherwise, the job is immediately rejected with `ErrQueueFull`. The worker loop processes jobs from the channel until context cancellation.

`ComputeBackoff` implements AWS Marc Brooker's jitter formulas. For Full Jitter: `sleep = random(0, min(cap, base × 2^attempt))`. For Equal Jitter: `sleep = min(cap, base × 2^attempt)/2 + random(0, half)`. For Decorrelated Jitter: `sleep = min(cap, random(base, prevSleep × 3))`. No Jitter returns the raw exponential value.

## Code Walkthrough

The demo (`cmd/demo/main.go`) exercises all four behaviors in sequence. It creates a token bucket with capacity 3 and refill rate 5/s, demonstrating that the first three requests are allowed (consuming tokens) and the fourth and fifth are rejected (bucket empty). After a 300ms pause, a token has replenished, allowing another request. The leaky bucket with capacity 3 and leak rate 10/s fills quickly: the first three requests are allowed but water reaches capacity, rejecting subsequent requests until drain occurs.

The bounded queue with capacity 3 and one worker accepts the first three jobs, then rejects jobs 4 and 6 when the channel buffer is full. Statistics confirm accepted=4, rejected=2, processed=1. The retry demonstration shows how Full Jitter produces lower latency on early attempts compared to No Jitter, with intervals spread uniformly across the exponential range.

The HTTP middleware (`internal/httputil/middleware.go`) enforces tenant rate limits, returning RFC 6585 standard `429 Too Many Requests` responses with `Retry-After` headers and JSON error payloads containing `error: "rate_limit_exceeded"` and the computed retry-after seconds.

## What the Tests Prove

The test suite validates correctness across all components:

- **TokenBucket tests** (`internal/ratelimit/bucket_test.go`): Verify burst capacity allows exactly capacity requests, refill timing restores tokens proportionally, token depletion blocks further requests, and concurrent access passes the Go race detector with 50 simultaneous goroutines.
- **LeakyBucket tests**: Confirm water level tracking, rejection when bucket reaches capacity, and concurrency safety.
- **Registry tests** (`internal/ratelimit/registry.go`): Prove per-tenant isolation—different tenants receive independent buckets, and concurrent access to the same key produces a single shared bucket.
- **BoundedQueue tests** (`internal/backpressure/queue_test.go`): Verify rejection under load returns `ErrQueueFull` without blocking the caller, concurrency safety under stress, and proper cleanup on `Stop()`.
- **Retry tests** (`internal/retry/backoff_test.go`): Validate that Full Jitter intervals fall within `[0, min(cap, base × 2^attempt)]`, Equal Jitter intervals fall within `[half, cap]`, and Decorrelated Jitter respects the multiplicative range.
- **HTTP middleware tests** (`internal/httputil/middleware_test.go`): Confirm RFC 6585 status code 429, `Retry-After` header presence, JSON response body structure, and anonymous fallback behavior.

All tests pass. Race detector runs are clean. The demo produces expected four-section output.

## Recovery / Rollback

The bounded queue's `Stop()` method provides graceful shutdown: it sets the stopped flag, cancels the context, closes the channel, and waits for all workers to finish via `sync.WaitGroup`. In-flight jobs complete processing; queued jobs drain. The registry requires no special teardown since token buckets are garbage-collected when no longer referenced.

For rate limiting recovery, the token bucket naturally self-heals as tokens replenish over time. The `RetryAfterSeconds()` method gives callers precise information about when they can retry, enabling intelligent backoff rather than blind repeated attempts.

## Production Considerations

The current implementation is in-memory only. Production distributed systems require external state storage—Redis is well-suited due to atomic operations (`INCR`, `EXPIRE`) and data structures (hashes, sorted sets, strings) that cover fixed-window counters, sliding windows, and token bucket algorithms. Redis provides sub-millisecond latency for the synchronous request path.

Default AWS SDK values (50ms base delay for transient errors, 1000ms for throttling, 20s maximum cap) are implementation-specific defaults calibrated to AWS service SLAs. These must be calibrated to downstream service latency and timeout profiles rather than adopted universally. The Full Jitter formula is one variant among three (Full, Equal, Decorrelated); AWS standardized on Full Jitter but the choice depends on specific contention patterns.

HTTP 429 versus 503 requires careful consideration: 429 indicates client quota exceeded (RFC 6585), while 503 may be more appropriate for server-side overload. NGINX defaults to 503 but makes it configurable. Both should include `Retry-After` headers when possible. Responses with 429 MUST NOT be cached by intermediaries.

Monitoring should track queue age rather than raw queue depth, as age provides earlier warning of degradation. Backlog growth follows the formula: `(arrival_rate - processing_rate) × time`.

## Common Mistakes

1. **IP-based rate limiting without tenant awareness**: Carrier-grade NAT (RFC 6598) means multiple users share one IP, causing collateral throttling. Use API keys or tenant IDs instead.
2. **No jitter in retries**: Synchronized retries create thundering herd problems, doubling or tripling effective load during recovery.
3. **Unbounded queues**: Allowing unlimited queue growth trades latency for memory exhaustion. Fixed-capacity queues with fast rejection prevent this.
4. **Confusing 429 and 503**: 429 is client quota exceeded; 503 is server unavailable. Using the wrong code misinforms clients about retry strategy.
5. **Ignoring retry budgets**: Google SRE recommends per-request budgets (max 3 attempts) and per-client budgets (max 10% retry ratio) to prevent retry storms during cascading overload.

## Case Study: Stripe's Four-Layer Approach

Stripe's engineering blog documents four production limiters. The **Request Rate Limiter** is the most important, enforcing per-user request rates. The **Concurrent Requests Limiter** caps in-flight requests per user. The **Fleet Usage Load Shedder** reserves capacity for critical traffic by shedding lower-priority requests. The **Worker Utilization Load Shedder** sheds test-mode traffic, then GET requests, then POST requests, preserving critical operations under pressure.

This layered approach enables graceful degradation: each layer catches a different failure mode, and lower layers activate only when higher layers are insufficient. Criticality levels (CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE) ensure that rejection cascades only after exhausting all lower-criticality options.

Google SRE extends this model with per-customer CPU quotas (Gmail 4000 CPU-seconds/sec, Calendar 4000, Android 3000, Google+ 2000, others 500) and utilization signals (executor load average with exponential decay smoothing, CPU, memory) to reject requests based on criticality thresholds. Client-side throttling prevents overwhelmed clients from amplifying failures: each client tracks requests and accepts, continuing only until requests are K times accepts.

## Checklist

- [ ] Identify rate limiting boundaries (API gateway, service mesh, individual service)
- [ ] Choose algorithm: Token Bucket (burst tolerance) or Leaky Bucket (strict pacing)
- [ ] Implement per-tenant key-based limiting (avoid IP-based)
- [ ] Configure bounded queues with capacity tuned to downstream processing rate
- [ ] Add HTTP 429 responses with Retry-After headers (RFC 6585)
- [ ] Implement retry with Full Jitter: `delay = random(0,1) × min(cap, base × 2^attempt)`
- [ ] Set retry budgets (per-request max 3, per-client max 10%)
- [ ] Monitor queue age and backlog growth
- [ ] Define criticality levels for graceful load shedding
- [ ] Run race detector tests (`go test -race`)
- [ ] Calibrate delay constants to downstream SLA/latency/timeout profiles

## Key Takeaways

1. Rate limiting and backpressure are complementary: boundary control plus internal pressure propagation.
2. Token Bucket permits bursts up to capacity; Leaky Bucket enforces strict output pacing.
3. HTTP 429 (RFC 6585) is the standard rate-limiting status code and MUST NOT be cached.
4. Little's Law (L = λW) provides the mathematical foundation for queue analysis and backlog prediction.
5. Exponential Backoff with Full Jitter (`random(0, min(cap, base × 2^attempt))`) prevents thundering herd synchronization.
6. Per-tenant key-based limiting avoids CGNAT collateral throttling.
7. Bounded queues with fast rejection prevent memory exhaustion and signal overload to callers.
8. Multi-layer approaches (Stripe's four limiters, Google SRE criticality levels) provide graceful degradation.
9. Retry budgets (per-request max 3, per-client max 10%) prevent retry storms during cascading failures.
10. In-memory implementations lack distributed coordination—Redis or similar state storage is required for production multi-instance deployments.

## Sources

- Token bucket — Wikipedia (https://en.wikipedia.org/wiki/Token_bucket)
- Leaky bucket — Wikipedia (https://en.wikipedia.org/wiki/Leaky_bucket)
- Rate limiting — Wikipedia (https://en.wikipedia.org/wiki/Rate_limiting)
- Exponential backoff — Wikipedia (https://en.wikipedia.org/wiki/Exponential_backoff)
- Little's law — Wikipedia (https://en.wikipedia.org/wiki/Little%27s_law)
- RFC 6585: Additional HTTP Status Codes — IETF (https://www.rfc-editor.org/rfc/rfc6585#section-4)
- Module ngx_http_limit_req_module — NGINX (https://nginx.org/en/docs/http/ngx_http_limit_req_module.html)
- Scaling your API with rate limiters — Stripe Engineering Blog (https://stripe.com/blog/rate-limiters)
- Exponential Backoff And Jitter — AWS Architecture Blog (https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)
- Retry behavior — AWS SDK Documentation (https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html)
- Handling Overload — Google SRE Workbook (https://landing.google.com/sre/sre-book/chapters/handling-overload/)
- Redis rate limiter documentation — Redis (https://redis.io/docs/latest/develop/use-cases/rate-limiter/)
- Consumer Prefetch & Queue Flow Control — RabbitMQ (https://www.rabbitmq.com/docs/consumer-prefetch)

Implementation files: `internal/ratelimit/bucket.go`, `internal/ratelimit/registry.go`, `internal/backpressure/queue.go`, `internal/retry/backoff.go`, `internal/httputil/middleware.go`, `cmd/demo/main.go`
Test files: `internal/ratelimit/bucket_test.go`, `internal/backpressure/queue_test.go`, `internal/retry/backoff_test.go`, `internal/httputil/middleware_test.go`
