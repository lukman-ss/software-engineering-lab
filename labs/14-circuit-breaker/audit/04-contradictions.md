# 04 Contradictions: Circuit Breaker Lab

## Contradiction 1

Statement A:
"Failure Modes: Ignoring error types: Counting 4xx (client errors) as circuit-breaking faults rather than restricting to 5xx/timeouts."
Location:
`README.md:136`, `research/09-failure-modes.md:6`

Statement B:
The `Execute(fn func() error) error` implementation unconditionally increments `failureCount` whenever `err != nil`. It treats all returned errors (including theoretical 4xx validation errors) equally.
Location:
`internal/circuitbreaker/circuit_breaker.go:118`

Type:
CODE_DOC_MISMATCH

Impact:
The lab documentation correctly warns against tripping the breaker on expected user validation failures. However, the reference Go implementation lacks an error classifier (e.g. `isCircuitFault(error) bool`) and trips on any generic error. Users copying the raw code might inadvertently trip production circuit breakers on `400 Bad Request` or `404 Not Found`.

Assessment:
MEDIUM. The codebase is an intentional simplification for an educational lab setting, emphasizing the core state machine transitions rather than edge cases. However, given that the documentation specifically warns against this exact behavior as a failure mode, the missing code-level error filter represents a discrepancy.

---

## Contradiction 2

Statement A:
"Key production metrics: circuit_state (gauge), circuit_open_count, rejected_call_count, failure_count, dependency_latency."
Location:
`README.md:78-83`, `research/08-observability.md:3-8`

Statement B:
The `CircuitBreaker` struct holds private fields for counters but does not export them, emit metrics, or expose hooks for a telemetry system (e.g., Prometheus / OpenTelemetry).
Location:
`internal/circuitbreaker/circuit_breaker.go`

Type:
CODE_DOC_MISMATCH

Impact:
The observability section provides theoretical guidance on how to monitor circuit breakers in a production environment, but the provided sample code does not actually implement this instrumentation.

Assessment:
WARNING (LOW). Educational scope. It is acceptable for a minimal lab implementation to omit full prometheus/metrics wiring to maintain readability, provided it's understood that the metrics section is architectural advice rather than an implemented feature.
