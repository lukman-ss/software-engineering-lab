# Cascade Failure

## Definition
Cascade failure (cascading failure) occurs when the failure of a single service triggers a chain of failures in dependent services, ultimately bringing down the entire system.

Source: Microsoft Azure Circuit Breaker pattern — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker (2025-02-05)
> "If a service is busy, failure in one part of the system might lead to cascading failures."

Source: Google SRE Book, Handling Overload — https://sre.google/sre-book/handling-overload/ (2017)
> "Left unchecked, the failure in a subset of a system (such as an individual backend task) might trigger the failure of other system components, potentially causing the entire system (or a considerable subset) to fail."

## Mechanism (Without Circuit Breaker)
```
slow/down dependency
→ requests wait (blocked threads)
→ latency increases (P99 spikes)
→ resources remain occupied (threads, sockets, memory, DB connections)
→ caller thread pool exhausts
→ caller's own health checks fail
→ upstream caller degrades
→ entire platform outage
```

### Evidence Chain

**Claim 1: Blocked requests hold critical resources**
- Source: Microsoft Azure Circuit Breaker (Tier 1)
- URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Evidence: "These blocked requests might hold critical system resources, such as memory, threads, and database connections. This problem can exhaust resources, which might fail other unrelated parts of the system that need to use the same resources."
- Confidence: HIGH

**Claim 2: Many callers on unresponsive supplier exhaust resources → cascading failures**
- Source: Martin Fowler, Circuit Breaker (Tier 1)
- URL: https://martinfowler.com/bliki/CircuitBreaker.html (2014-03-06)
- Evidence: "What's worse if you have many callers on a unresponsive supplier, then you can run out of critical resources leading to cascading failures across multiple systems."
- Confidence: HIGH — corroborated by Azure and Google SRE

**Claim 3: Backend should continue serving at provisioned rate under overload (not crash)**
- Source: Google SRE Book, Handling Overload (Tier 1)
- URL: https://sre.google/sre-book/handling-overload/
- Evidence: "A backend task provisioned to serve a certain traffic rate should continue to serve traffic at that rate without any significant impact on latency, regardless of how much excess traffic is thrown at the task."
- Confidence: HIGH

## Why Timeout Alone Is Insufficient
- Timeout prevents infinite hang but still blocks for the full timeout duration (e.g., 100ms–30s)
- Under high concurrency, N requests × timeout = thread pool exhaustion before any timeout fires
- Azure notes: "To resolve this problem, set a shorter time-out. But ensure that the time-out is long enough for the operation to succeed most of the time." — tradeoff, not a solution to cascade

## How Circuit Breaker Prevents Cascade
- Fails fast in nanoseconds/microseconds once OPEN (no thread blocking)
- Releases caller resources immediately
- Gives downstream idle recovery time
- Google SRE: backend should "accept only the requests that it can process and reject the rest gracefully" — circuit breaker is client-side enforcement of this

## NOT VERIFIED
- Specific quantitative thresholds for when cascade triggers (e.g., "at X% thread pool utilization cascade begins") — no authoritative source provides a universal number; depends on architecture.

## Related Concepts
- **Load Shedding**: Server-side rejection of excess load (see Google SRE criticality levels)
- **Bulkhead**: Resource isolation so failure of dependency A does not consume resources needed for dependency B
- **Retry Storm**: Retries amplify load on already-failing service, accelerating cascade (see 06-timeout-retry-backoff.md)
