# Observability

## Metrics Inventory for Circuit Breaker

### Core State Metrics
| Metric | Type | Description | Source |
|--------|------|-------------|--------|
| `circuit_state` | Gauge | Current state: 0=CLOSED, 1=OPEN, 2=HALF_OPEN | Architectural recommendation derived from Azure "Monitoring: A circuit breaker should provide clear observability into both failed and successful requests" — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker |
| `circuit_open_count` | Counter | Number of times circuit transitioned to OPEN | Martin Fowler: "Any change in breaker state should be logged and breakers should reveal details of their state for deeper monitoring." — https://martinfowler.com/bliki/CircuitBreaker.html |
| `rejected_call_count` | Counter | Requests rejected by fail-fast while OPEN | Derived from Azure's fail-fast mechanism description |

### Failure Metrics
| Metric | Type | Description | Source |
|--------|------|-------------|--------|
| `failure_count` | Counter/Gauge | Current consecutive or windowed failure count | Martin Fowler: "the breaker stores the block... failure count" |
| `failure_rate` | Gauge | Percentage of failed calls in sliding window | Resilience4j: "failure rate is equal or greater than a configurable threshold" — https://resilience4j.readme.io/docs/circuitbreaker |
| `timeout_count` | Counter | Failures specifically due to timeout | Resilience4j distinguishes exception types (recordExceptions vs ignoreExceptions) |
| `slow_call_rate` | Gauge | Percentage of calls exceeding slow threshold | Resilience4j: `slowCallRateThreshold` / `slowCallDurationThreshold` |

### Performance Metrics
| Metric | Type | Description | Source |
|--------|------|-------------|--------|
| `dependency_latency` | Histogram | P50, P95, P99 latencies for downstream calls | Azure: "dependency_latency" referenced in topic spec; latency histogram is standard practice |
| `fallback_count` | Counter | Number of fallback invocations when OPEN | Topic spec; derived from Azure "gracefully degrade by returning default or cached responses" |

### State Transition Events
**Claim**: Every state transition should be logged as an event with timestamp, from-state, to-state, and trigger reason.

**Evidence**:
- Martin Fowler: "Circuit breakers are a valuable place for monitoring. Any change in breaker state should be logged and breakers should reveal details of their state for deeper monitoring. Breaker behavior is often a good source of warnings about deeper troubles in the environment." — https://martinfowler.com/bliki/CircuitBreaker.html
- Azure: "If the circuit breaker raises an event each time it changes state, this information can help monitor the health of the protected system component or alert an administrator when a circuit breaker switches to the Open state." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

**Confidence**: HIGH.

### Alerting Recommendations
**Claim**: Operations staff should be alerted when a breaker trips to OPEN.

**Evidence**:
- Martin Fowler: "Usually you'll also want some kind of monitor alert if the circuit breaker trips." — https://martinfowler.com/bliki/CircuitBreaker.html
- Martin Fowler: "Operations staff should be able to trip or reset breakers."

**Confidence**: HIGH.

---

## NOT VERIFIED
- Specific numeric alert thresholds (e.g., "alert if failure_rate > 50%") — depends on business SLA, no universal production recommendation.
- Whether `failure_rate` should be computed over a count-based or time-based window — Resilience4j supports both; choice depends on traffic volume and detection latency requirements.
- Benchmarks of metric overhead (circuit breaker metrics causing performance impact) — no quantitative source found.

---

## Metrics Not Invented
This research does NOT include fabricated metric names, thresholds, or latency numbers. The metric inventory above is derived from:
1. Explicit patterns in Resilience4j (failureRateThreshold, slowCallRateThreshold, slowCallDurationThreshold)
2. Explicit monitoring recommendations in Martin Fowler and Azure documentation
3. Standard time-series metric naming conventions (counter/gauge/histogram)

All numeric values are labeled as architectural recommendations or explicitly marked NOT VERIFIED where no authoritative source exists.