# 03 Claim Audit: Circuit Breaker Lab

## Claim 1

Claim:
Remote calls across networks fail or hang. Without protection, callers block waiting for timeouts, holding threads, sockets, and memory until system resources deplete, triggering cascading failures.

Location:
`README.md` (Problem, Cascade Failure); `research/04-cascade-failure.md`

Evidence Provided:
Detailed description of worker thread starvation, HTTP connection pool exhaustion, health check timeouts, and upstream propagation. Demonstrated in Demo Scenario 1 where 3 requests block for 100ms each.

Source:
Martin Fowler (*Circuit Breaker*); Microsoft Learn (*Circuit Breaker Pattern*)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Both Martin Fowler and Microsoft Learn explicitly cite cascading resource exhaustion (threads, memory, sockets) as the primary rationale for introducing circuit breakers.

---

## Claim 2

Claim:
A Circuit Breaker monitors downstream failures, trips OPEN when failure threshold is reached, and fails fast in microsecond/nanosecond timeframes without executing downstream network calls.

Location:
`README.md` (Circuit Breaker Mental Model, OPEN, With Circuit Breaker); `research/03-core-concepts.md`

Evidence Provided:
Demo Scenario 2 execution logs showing requests 4, 5, and 6 failing fast in under 1µs with `ErrCircuitOpen`, while downstream call count stays fixed at 3. Tested in `TestCircuitBreaker/5. OPEN calls fail fast` and `TestCircuitBreaker/6. OPEN calls do not execute downstream function`.

Source:
Martin Fowler (*Circuit Breaker*); Microsoft Learn (*Circuit Breaker Pattern*)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core mechanism supported by canonical definitions and verified by lab execution.

---

## Claim 3

Claim:
In HALF-OPEN state, after cooldown timer (`OpenTimeout`) elapses, the circuit allows limited canary probe requests (`HalfOpenMaxCalls`). If the probe succeeds, the circuit resets to CLOSED; if it fails, the circuit returns to OPEN and restarts the cooldown timer.

Location:
`README.md` (HALF-OPEN, Recovery); `research/05-circuit-states.md`

Evidence Provided:
Demo Scenarios 3 and 4; unit tests 7, 8, 9, and 10 in `circuit_breaker_test.go`.

Source:
Martin Fowler (*Circuit Breaker*); Microsoft Learn (*Circuit Breaker Pattern*)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fully verified by state machine implementation and test assertions.

---

## Claim 4

Claim:
Retries without circuit breakers can cause retry storms and amplify downstream overload. Circuit breakers silence retries by failing fast when failures persist.

Location:
`README.md` (Timeout vs Retry vs Circuit Breaker); `research/06-timeout-retry-backoff.md`

Evidence Provided:
Analysis of retry amplification and exponential backoff synergy.

Source:
AWS Builders Library (*Timeouts, retries, and backoff with jitter*); Microsoft Learn (*Circuit Breaker Pattern - Note on Retry Pattern*)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
AWS Builders Library documents retry storms and backoff; Microsoft Learn explicitly states that retry logic should detect circuit breaker exceptions and cease attempts if the circuit is OPEN.

---

## Claim 5

Claim:
Key production metrics for circuit breakers include `circuit_state` (gauge 0, 1, 2), `circuit_open_count`, `rejected_call_count`, `failure_count`, and `dependency_latency`.

Location:
`README.md` (Observability); `research/08-observability.md`

Evidence Provided:
Textual description of metric semantics and observability recommendations.

Source:
Microsoft Learn (*Circuit Breaker Pattern - Monitoring*)

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
The metrics listed follow standard observability practices (e.g., Prometheus / OpenTelemetry conventions). Note that these are documented as guidance and are not instrumented inside the Go reference code.

---

## Claim 6

Claim:
Counting 4xx client errors as circuit-breaking faults rather than restricting to 5xx server errors and network timeouts is a failure mode.

Location:
`README.md` (Failure Modes); `research/09-failure-modes.md`

Evidence Provided:
Textual explanation of failure mode.

Source:
Microsoft Learn (*Circuit Breaker Pattern - Types of exceptions*)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
MEDIUM

Notes:
Microsoft Learn highlights that circuit breakers should evaluate exception types and not trip on normal client-level validation errors. While the research correctly identifies this, the reference Go implementation in `circuitbreaker.go` lacks an error classifier and counts any `err != nil` towards the failure threshold.

---

## Claim 7

Claim:
Illustrative lab timeouts are 100ms HTTP timeout, 300ms cooldown, and failure threshold of 3; production timeouts must be tuned to SLA, P99, and recovery profiles.

Location:
`README.md` (Note on Lab Timeouts); `cmd/demo/main.go`

Evidence Provided:
Working demo configuration runs in under 1.5 seconds. Explicit disclaimer present in `README.md:63`.

Source:
Implementation configuration and documentation.

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
The documentation transparently acknowledges that the 100ms/300ms values are chosen for interactive testing and educational demonstrations, not universal production defaults.

---

## Claim 8

Claim:
In real-world architectures (e.g., CMMS invoice generation with WhatsApp notifications, or PPOB payment gateway with Pulsa provider), asynchronous non-critical downstream failures must not abort core transaction persistence.

Location:
`README.md` (CMMS Example, PPOB Example); `research/10-final-research.md`

Evidence Provided:
Architectural breakdown distinguishing synchronous critical dependencies from asynchronous queue-backed non-critical dependencies.

Source:
Engineering architectural design examples (uncited domain case studies).

Source Actually Supports Claim:
PARTIAL

Classification:
EXAMPLE

Severity:
LOW

Notes:
Sound domain architecture modeling. These examples serve as contextual illustrations rather than empirical claims requiring academic citation.

---

## Claim 9

Claim:
The Go in-memory Circuit Breaker implementation is thread-safe and free of data races under concurrent execution.

Location:
`internal/circuitbreaker/circuit_breaker.go`; `internal/circuitbreaker/circuit_breaker_test.go`

Evidence Provided:
Use of `sync.Mutex` on all state transitions and counter access. Test 11 runs 50 concurrent goroutines. Verified clean with `go test -race -count=1 ./...`.

Source:
Go standard library (`sync.Mutex`); test execution output.

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Execution verified with race detector passing without warnings.
