# Implementation Notes

## Files Added

- `go.mod`: Module definition (`labs/25-rate-limiting-and-backpressure`).
- `internal/ratelimit/bucket.go`: Thread-safe `TokenBucket` and `LeakyBucket` implementations with `sync.Mutex` and monotonic time tracking.
- `internal/ratelimit/registry.go`: `Registry` managing per-tenant token buckets identified by API keys / tenant IDs to avoid CGNAT IP collision (RFC 6598).
- `internal/ratelimit/bucket_test.go`: Unit and concurrency race tests for rate limiters.
- `internal/backpressure/queue.go`: Bounded channel worker pool (`BoundedQueue`) rejecting jobs fast when queue capacity is reached.
- `internal/backpressure/queue_test.go`: Unit tests proving rejection under load and concurrency safety.
- `internal/retry/backoff.go`: Exponential backoff implementation supporting AWS Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter.
- `internal/retry/backoff_test.go`: Boundary verification tests for jitter distributions.
- `internal/httputil/middleware.go`: HTTP middleware enforcing tenant rate limits, returning RFC 6585 status `429 Too Many Requests` with `Retry-After` header.
- `internal/httputil/middleware_test.go`: Integration test verifying HTTP status 429 response format and headers.
- `cmd/demo/main.go`: Interactive CLI execution demonstrating token bucket bursts, leaky bucket smoothing, backpressure load shedding, and retry jitter distributions.

## Core Design Decisions

1. **Thread-Safe Mutex Guarding**: Synchronized internal state in token and leaky buckets using `sync.Mutex` to guarantee safety under concurrent callers.
2. **Monotonic Time Calculation**: Used `time.Now().Sub(...)` for continuous rate calculation, avoiding clock-skew vulnerabilities.
3. **Non-Blocking Channel Backpressure**: Used Go select-default channel pattern in `BoundedQueue.TrySubmit` for zero-allocation fast rejection when queue capacity is filled.
4. **AWS Jitter Specification Compliance**: Implemented Full Jitter ($sleep = random(0, \min(cap, base \cdot 2^{attempt}))$) matching AWS Architecture Blog standards by Marc Brooker.

## Implementation-Specific Choices

- Floating point token arithmetic for smooth fractional refill precision.
- Default HTTP header `X-API-Key` for tenant extraction, fallback to `anonymous`.

## Known Limitations

- In-memory rate limiting state without distributed synchronization across multiple process replicas (Redis/Memcached state storage deferred to future distributed extensions).

## Trade-offs

- Mutex synchronization vs Lock-free CAS atomic counters: Mutexes chosen for simplicity, code clarity, and robust multi-field updates (tokens + timestamps).

## What Is Demonstrated

- Token bucket burst tolerance vs Leaky bucket output smoothing.
- Load shedding via bounded queue backpressure preventing memory unbounded growth.
- AWS randomized exponential backoff reducing thundering herd risk.
- RFC 6585 standard HTTP 429 response header generation.

## What Is Not Demonstrated

- Distributed rate limiting across multi-region clusters.
- Dynamic queue capacity autoscaling.
