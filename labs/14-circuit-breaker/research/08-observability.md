# 08 Observability

## Required / Common Metrics

| Metric | Type | Notes |
|---|---|---|
| `circuit_state` | Enum/gauge (0=CLOSED, 1=OPEN, 2=HALF_OPEN) | dashboard color |
| `circuit_open_count` | Counter | alerts on transition to OPEN |
| `failure_count` | Gauge | rolling or consecutive; Hystrix `ErrFailures`, `ErrTimeouts` |
| `failure_rate` / `error_percentage` | Gauge | Hystrix tripping condition: error % over rolling window |
| `timeout_count` | Counter | distinguish timeout vs connection error |
| `fallback_count` | Counter | degraded responses emitted |
| `rejected_call_count` | Counter | fast-failed due to OPEN |
| `dependency_latency` | Histogram | per-attempt latency; overflows drive timeout tuning |

**Evidence**: Azure Circuit Breaker — "Circuit breakers should provide clear observability into both failed and successful requests...". Hystrix Wiki — "reports successes, failures, rejections, and timeouts to the circuit breaker, which maintains a rolling set of counters" + metrics for dashboards/streams. — Confidence HIGH.

## Reference: Go Implementations
- **gobreaker** (`sony/gobreaker`): `Counts {Requests, TotalSuccesses, TotalFailures, TotalExclusions, ConsecutiveSuccesses, ConsecutiveFailures}` + `OnStateChange` callback for state-transition metrics. — https://github.com/sony/gobreaker
- **cep21/circuit**: `CmdMetricCollector`, `FallbackMetricCollector`, `rolling.StatFactory`, `MetricEventStream` streaming, `expvar` publisher `h.Var()`, `responsetimeslo.Factory` for "X% requests faster than Y ms" SLO tracking. — https://github.com/cep21/circuit (Tier 1).

## Alphabetization Note
- Hystrix Dashboard stream: `/hystrix.stream` endpoint + Turbine aggregation across fleet.
- In-process counters without external aggregation lose fleet visibility (each replica's view differs).

## Alerting Philosophy (SRE)
- Alert on state-change events (circuit opened) rather than raw failure count.
- Correlate `dependency_latency` rise with `timeout_count` before tripping—latency often precedes hard errors.

**Source**: Azure Circuit Breaker Problems & Considerations — "Circuit breakers should reveal details of their state for deeper monitoring... Alert an administrator when a circuit breaker switches to the Open state".

## Lab Minimal Observability
This lab's `State()` + counters are sufficient for tests/demos; production adds: Prometheus histograms, structured log events per transition, distributed tracing spans annotated with `cb.state`.

**Lab note**: Lab counters are authoritative for "expected behavior" claims; no fabricated benchmarks.