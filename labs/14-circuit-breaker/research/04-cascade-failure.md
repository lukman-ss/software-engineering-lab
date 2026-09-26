# Cascade Failure Analysis

## Definition
A cascading failure occurs when the failure of one component causes a chain reaction that brings down other components and eventually the entire system.

**Source**: Google SRE Handbook (Source 5), AWS Well-Architected (Source 4)

## Mechanism (Without Circuit Breaker)

```
Downstream Service (slow/down)
↓
Caller requests block waiting for response
↓
Caller threads/connections remain occupied
↓
Caller cannot serve new requests
↓
Upstream callers also block
↓
System-wide degradation or outage
```

### Resource Exhaustion
1. **Thread/Connection Pool Exhaustion**: Worker threads block on I/O, pool drains
2. **Timeout Accumulation**: Requests queue up, each waiting full timeout duration
3. **Memory Pressure**: Request contexts accumulate in memory
4. **Backpressure Propagation**: Upstream services experience increased latency

**Source**: Google SRE Handbook (Source 5), Netflix Hystrix (Source 2)

## Example Timeline (Illustrative)

```
t=0s:   Payment service becomes slow (latency 30s)
t=0s:   Checkout receives 100 requests
t=0s:   All 100 workers spawn outbound calls to Payment
t=5s:   Workers still blocked, no new capacity
t=10s:  Upstream queue builds, timeouts start
t=30s:  First batch times out, but new requests already queued
t=30s+: System remains saturated, recovery impossible under load
```

## How Circuit Breaker Prevents Cascade
1. **Threshold**: After N failures, circuit opens
2. **Fail Fast**: New requests rejected in microseconds, not seconds
3. **Resource Release**: No threads/connections held waiting
4. **Recovery Window**: Downstream gets breathing room

**Source**: Microsoft Azure (Source 3), Martin Fowler (Source 1)

## Quantitative Difference (Illustrative, NOT Production Data)
- Without CB: Each request costs full timeout (e.g., 5s), 100 requests = 500s of worker time
- With CB: First N requests cost timeout, rest cost ~1ms (fail fast)
- Resource savings: ~99% reduction in blocked worker time after circuit opens

**Note**: Exact numbers depend on timeout configuration and traffic patterns. No universal benchmarks exist.

## Cascade vs Related Patterns
| Pattern | Prevents | How |
|---------|----------|-----|
| Circuit Breaker | Downstream failure propagation | Fail fast when threshold reached |
| Timeout | Individual request hanging | Bound maximum wait time |
| Bulkhead | Cross-dependency contamination | Isolate resources per dependency |
| Load Shedding | Overload from excess traffic | Drop low-priority requests |

**Source**: Microsoft Azure (Source 3, 9), Google SRE Workbook (Source 10)