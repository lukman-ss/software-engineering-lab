# Engineering Design

Target Lab: labs/25-rate-limiting-and-backpressure
Research Status: APPROVED

## Concept To Prove

1. Token Bucket & Leaky Bucket algorithms: rate limiting with burst accommodation (token bucket) vs strict rate smoothing (leaky bucket).
2. HTTP 429 semantics: RFC 6585 compliance with `Retry-After` headers and JSON error payloads.
3. Bounded Queue Backpressure: prevent queue latency buildup and OOM under arrival rate exceeding service capacity ($\Delta Q = (r_{in} - r_{out}) \cdot \Delta t$), rejecting fast when capacity is exceeded.
4. Retry strategies with jitter: Full Jitter, Equal Jitter, and Decorrelated Jitter (Marc Brooker / AWS Architecture Blog) preventing retry storms and synchronization.
5. Multi-tenant key-based rate limiting (API key / tenant ID) rather than naive IP-based limits (RFC 6598 CGNAT avoidance).

## Expected Behavior

- Token bucket allows bursts up to token capacity $B$, then enforces refill rate $R$.
- Leaky bucket smooths flow to leak rate $R$, rejecting immediate bursts when water reaches capacity.
- Under sustained load exceeding service capacity, bounded queue worker accepts up to queue limit, then returns fast backpressure rejection (503 / 429).
- Backoff with Full Jitter distributes retry intervals uniformly in $[0, \min(cap, base \cdot 2^{attempt})]$, avoiding herd synchronization.
- HTTP middleware inspects API keys / tokens and returns RFC 6585 standard HTTP 429 responses with `Retry-After` header when exhausted.

## Failure Scenario

- Fast incoming traffic without rate limiting causes unbounded queue buildup, memory starvation, and timeout cascades.
- Synchronized retries without jitter cause periodic load spikes (retry storms) amplifying downstream degradation.
- Naive IP-only rate limiting throttles multiple users behind carrier-grade NAT (RFC 6598).

## Success Criteria

- TokenBucket and LeakyBucket accurately track and reject requests exceeding capacity/rates under concurrent access.
- Bounded queue rejects excess requests immediately with backpressure when buffer is full.
- Retry backoff produces jittered intervals adhering to AWS formulas.
- HTTP rate limiter handler returns status 429 with correct `Retry-After` header.
- Concurrency test passes with Go race detector clean (`go test -race ./...`).
- CLI demo runs showing rate limiting, backpressure drops, and retry jitter distributions.

## Architecture

```
[Incoming Request]
        │
        ▼
[HTTP Middleware / Tenant Extractor] ── (RFC 6598 Tenant Key)
        │
        ▼
[Token Bucket / Rate Limiter] ── (Exceeded? -> HTTP 429 + Retry-After)
        │ (Passed)
        ▼
[Bounded Queue / Backpressure Channel] ── (Queue Full? -> 503 Overloaded)
        │ (Enqueued)
        ▼
[Worker Pool / Consumer] ── (Little's Law L = λW steady state)
        │
        ▼
[Client with Full Jitter Retry] ── (AWS Jitter Backoff)
```

## Components

1. `ratelimit`:
   - `TokenBucket`: Thread-safe token bucket with fractional token replenishment.
   - `LeakyBucket`: Thread-safe leaky bucket with constant leak rate.
   - `Limiter`: Multi-tenant key-based registry managing per-tenant buckets.
2. `backpressure`:
   - `BoundedQueue`: Fixed-capacity worker pool rejecting jobs when full (fast drop / shed).
3. `retry`:
   - `JitterType`: NoJitter, FullJitter, EqualJitter, DecorrelatedJitter implementations following AWS specifications.
4. `http`:
   - HTTP middleware returning RFC 6585 `429 Too Many Requests` + `Retry-After`.
5. `cmd/demo`:
   - Interactive CLI displaying burst allowance, leaky bucket smoothing, bounded queue shedding, and retry jitter spread.

## Test Strategy

- `ratelimit_test.go`: Unit tests for burst capacity, refill timing, token depletion, leaky bucket drain, and concurrency safety.
- `backpressure_test.go`: Verify bounded queue capacity drops excess requests without blocking caller indefinitely.
- `retry_test.go`: Verify jitter interval bounds for Full Jitter, Equal Jitter, and No Jitter.
- `http_test.go`: Verify HTTP 429 status code and `Retry-After` header parsing.
- Concurrency & Race Detector: 50 concurrent goroutines pounding token bucket and queue.

## Execution Plan

1. Initialize Go module `labs/25-rate-limiting-and-backpressure`.
2. Implement packages in `internal/ratelimit`, `internal/backpressure`, `internal/retry`, and `internal/httputil`.
3. Create `cmd/demo/main.go` demonstrating all core behaviors.
4. Execute `go test ./...` and `go test -race ./...`.
5. Run demo and capture output for `03-execution-result.md`.
6. Write `02-implementation-notes.md` and `README.md`.

## Implementation Decisions

- Decision 1: Pure Go standard library (`sync`, `time`, `net/http`) without third-party dependencies to maximize clarity and testability.
- Decision 2: Monotonic clock comparisons via Go's `time.Now()` for rate limit refills to guard against system clock adjustments.
- Decision 3: Deterministic simulated time helpers in tests to avoid flaky wall-clock sleeps while retaining real-time execution in demo.
