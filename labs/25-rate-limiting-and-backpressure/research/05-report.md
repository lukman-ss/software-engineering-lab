# Technical Report

## Title: Rate Limiting & Backpressure Engineering

---

## 1. Overview

Rate limiting and backpressure are complementary mechanisms for controlling traffic load on distributed systems. Rate limiting restricts the rate of incoming requests; backpressure applies downstream capacity constraints to upstream callers. Together, they prevent resource exhaustion and cascading failures.

---

## 2. HTTP Rate Limiting

### 2.1 Standard Response: HTTP 429

The IETF standard for HTTP rate limiting responses is defined in RFC 6585 Section 4:

- Status Code: `429 Too Many Requests`
- Optional header: `Retry-After` (HTTP-date or delay-seconds)

```
HTTP/1.1 429 Too Many Requests
Retry-After: 120
Content-Type: application/json

{
  "error": "rate_limit_exceeded",
  "retry_after": 120
}
```

**Sources**: RFC 6585, RFC 9110

### 2.2 Rate Limiting Algorithms

#### Token Bucket

- Tokens refill at constant rate $R$ per second
- Bucket capacity $B$ allows burst of up to $B$ requests
- State: $B_t = \min(B, B_{t-1} + R \cdot \Delta t)$

```python
class TokenBucket:
    def __init__(self, capacity: int, refill_rate: float):
        self.capacity = capacity
        self.tokens = capacity
        self.refill_rate = refill_rate  # tokens per second
        self.last_refill = time.monotonic()

    def allow(self) -> bool:
        now = time.monotonic()
        elapsed = now - self.last_refill
        self.tokens = min(self.capacity, self.tokens + self.refill_rate * elapsed)
        self.last_refill = now

        if self.tokens >= 1:
            self.tokens -= 1
            return True
        return False
```

**Characteristics**: Allows bursts up to capacity. Commonly used in API gateways (e.g., AWS API Gateway, Stripe API).

#### Leaky Bucket

- Drains queue at constant rate $R$
- Smooths output to strict constant rate
- No burst accommodation

```python
class LeakyBucket:
    def __init__(self, capacity: int, leak_rate: float):
        self.capacity = capacity
        self.water = 0.0
        self.leak_rate = leak_rate
        self.last_leak = time.monotonic()

    def allow(self, now: float = None) -> bool:
        now = now or time.monotonic()
        elapsed = now - self.last_leak
        self.water = max(0, self.water - self.leak_rate * elapsed)
        self.last_leak = now

        if self.water < self.capacity:
            self.water += 1
            return True
        return False
```

**Characteristics**: Strict rate limiting, no burst tolerance. Suitable for traffic shaping (e.g., network buffers).

#### Fixed Window / Sliding Window

- **Fixed Window**: Counter resets every fixed interval (e.g., 100 requests per minute)
- **Sliding Window**: Time window slides, counting requests in the past N seconds
- **Rolling Window**: Combination for approximate fairness

**Sources**: RFC 2697 (token bucket/srTCM), RFC 6585 (HTTP 429)

---

## 3. Queueing Theory

### 3.1 Little's Law (Correct Formulation)

**Source**: J.D.C. Little, "A Proof for the Queuing Formula: L = λW", Operations Research, 1961.

**Formula**: $L = \lambda W$

Where:
- $L$ = average number of items in a stationary, stable system
- $\lambda$ = effective arrival rate (items/time, must be < service rate for stability)
- $W$ = average time an item spends in the system (residence time)

**Requirements**: System must be stable (average arrival rate < average service rate) and operating in steady state.

**Application**: Capacity planning — if target latency $W$ must be ≤ 50ms and throughput $\lambda$ = 1000 req/s, then max average queue length $L$ = 50 items.

### 3.2 Deterministic Queue Buildup (Fluid Model — NOT Little's Law)

**Formula**: $\Delta Q = (r_{in} - r_{out}) \cdot \Delta t$

Where:
- $Q$ = queue length
- $r_{in}$ = instantaneous arrival rate
- $r_{out}$ = instantaneous service rate
- $\Delta t$ = elapsed time

**Important**: This is a transient, deterministic approximation, not Little's Law. It describes rate differential accumulation, not steady-state behavior.

**Misattribution Note**: Presenting $\Delta Q = (r_{in} - r_{out}) \cdot \Delta t$ as Little's Law is a technical error. Little's Law relates time-averages; the fluid model describes instantaneous rate differentials.

### 3.3 Unbounded Queue Risks

**Source**: Google SRE Book, Chapter 22 — Addressing Cascading Failures

An unlimited queue:
- **Does not prevent** system failure
- **Shifts** failure point from immediate rejection → memory exhaustion (OOM)
- **Causes** queue age degradation (latency accumulation)
- **Triggers** timeout cascade (downstream deadlines exceeded)
- **Creates** head-of-line blocking (slow consumers block all items)

> "Queued requests consume memory and increase latency... If there is insufficient capacity to handle all the requests at steady state, the server will saturate its queues." — Google SRE Book, §22.10

**Recommendation**: Use bounded queues with explicit rejection (backpressure) or dead-letter queue policies.

---

## 4. Retry Strategies & Backoff

### 4.1 Exponential Backoff

- Each retry multiplies backoff by a constant factor
- Backoff is capped at a maximum value to bound total retry time
- Without jitter, retries synchronize ("thundering herd")

### 4.2 Jitter (AWS Standard)

**Source**: Marc Brooker, "Exponential Backoff And Jitter", AWS Architecture Blog, 2015.

| Algorithm | Formula |
|-----------|---------|
| **No Jitter** | $sleep = \min(cap, base \cdot 2^{attempt})$ |
| **Equal Jitter** | $sleep = \min(cap, \frac{base \cdot 2^{attempt}}{2}) + random(\frac{base \cdot 2^{attempt}}{2})$ |
| **Full Jitter** | $sleep = random(\min(cap, base \cdot 2^{attempt}))$ |

**Full Jitter** recommended: $sleep = random(0, \min(cap, base \cdot 2^{attempt}))$

### 4.3 Retry Storm / Amplification

**Source**: Google SRE Book, Chapter 22 — §22.8 "Retries"

Retry amplification occurs when:
1. Multiple retry layers multiply attempts: 3 layers × 4 attempts = **64 attempts** on the backend (4³)
2. Retries without jitter synchronize, creating periodic load spikes
3. Retries on overloaded backends compound the failure

**Mitigations**:
- Use randomized exponential backoff (Full Jitter)
- Implement server-wide retry budget (e.g., 60 retries/minute per process)
- Avoid retrying at multiple layers simultaneously
- Return specific status codes (429, 503) to signal overload without retry

---

## 5. IP-Based Rate Limiting Limitations

**Source**: RFC 6598 — Shared Address Space (100.64.0.0/10)

CGNAT (Carrier-Grade NAT):
- Multiple subscribers share a single public IP address
- Rate-limiting by IP address penalizes innocent users behind the same NAT
- RFC 6598 formally defines the shared address space for CGNAT deployments

**Recommendation**: For API services, use tenant identifier (API key, OAuth bearer token) for per-tenant rate limiting.

---

## 6. Monitoring & Observability

### 6.1 Queue Age Metrics

**Queue age** (time spent in queue) is a superior early warning indicator for backpressure compared to queue depth alone:

| Broker | Metric | Purpose |
|--------|--------|---------|
| AWS SQS | `ApproximateAgeOfOldestMessage` | Time elapsed since oldest message entered queue |
| AWS SQS | `ApproximateNumberOfMessagesVisible` | Count of messages available for processing |
| Kafka | Consumer Lag (offset difference) | Messages behind latest offset |
| Kafka | Record queue time | Time from record production to consumption |

**Finding**: Queue age directly measures SLA violation risk. A queue may have low depth but high age if messages are large or consumers are slow. Both metrics should be monitored together.

### 6.2 Queue Depth Metrics

Queue depth (item count) provides operational signal but:
- Does not account for message size variability
- Does not capture consumer processing time variance
- May appear healthy while individual messages age out

**Source**: AWS SQS CloudWatch Documentation

---

## 7. Autoscaling & Downstream Bottlenecks

**Source**: Google SRE Book, Chapter 22 — §22.2 "Server Overload"

**Finding**: Autoscaling upstream consumers without scaling downstream dependencies (database, cache, external APIs) causes cascade failure:

1. Upstream workers scale from 10 → 100
2. All 100 workers hit downstream service simultaneously
3. Downstream service saturates → errors increase
4. Failed requests trigger retries → amplification
5. Downstream crashes → redistributes load to remaining instances
6. Cascade propagates across failure domains

**Recommendation**: 
- Scale all layers together (frontend, backend, database)
- Implement circuit breakers on downstream connections
- Use bulkhead pattern (per-tenant connection pools)
- Monitor downstream capacity saturation before scaling upstream

---

## 8. Fair Queuing & Per-Tenant Isolation

Fair queuing (FQ), Weighted Fair Queuing (WFQ), and WF2Q+ are established scheduling algorithms for equitable bandwidth distribution:

- **FQ**: Equal resource allocation per flow
- **WFQ**: Weighted allocation based on flow priority
- **WF2Q+**: Worst-case fair service with guaranteed upper bounds

**Application**: Rate limiting should use per-tenant identifiers as queue keys (API key, user ID) rather than shared IP addresses.

**Note**: While the scheduling theory is well-established (Nagle 1987, Jiang et al. 2005), specific weight assignment is service-dependent.

---

## 9. Summary

| Concept | Key Finding | Source |
|---------|-------------|--------|
| HTTP 429 | Standardized in RFC 6585 §4 | RFC 6585, RFC 9110 |
| Token Bucket | Bursts allowed up to capacity B | RFC 2697 |
| Little's Law | $L = \lambda W$ — stochastic, steady-state | Little (1961) |
| Queue Buildup | $\Delta Q = (r_{in} - r_{out}) \Delta t$ — deterministic, transient | Fluid model |
| Unbounded Queue | Shifts failure: rejection → OOM/age | Google SRE §22.10 |
| Exponential Backoff | Base × 2^attempt, capped | AWS Architecture Blog |
| Jitter | Full Jitter = random(0, min(cap, base·2^attempt)) | AWS Architecture Blog |
| Retry Budget | e.g., 60 retries/minute per process | Google SRE §22.8 |
| IP Rate Limiting | Problematic under CGNAT | RFC 6598 |
| Queue Age | `ApproximateAgeOfOldestMessage` (SQS), Consumer Lag (Kafka) | AWS/Kafka docs |
| Autoscaling | Must include downstream dependencies | Google SRE §22.2 |

---

## 10. Remaining Open Questions

1. **Empirically-validated thresholds**: What are optimal timeout/retry values per service type? (No universal best practice — service-specific)
2. **Cross-cloud benchmarking**: Comparative analysis of rate limiting implementations (AWS, GCP, Azure, Cloudflare)
3. **Real-world case studies**: Documented cascade failures and recovery patterns in production systems

---

## 11. Classification Legend

| Label | Meaning |
|-------|---------|
| FACT | Directly supported by authoritative source |
| INTERPRETATION | Supported by theory but implementation varies |
| NOT VERIFIED | Not yet supported by cited source |
| EXAMPLE | Illustrative, not generalizable