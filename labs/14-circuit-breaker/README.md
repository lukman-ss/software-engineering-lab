# Lab 14 — Circuit Breaker

## Problem
Remote calls across networks fail or hang. Without protection, callers block waiting for timeouts, holding threads, sockets, and memory until system resources deplete.

## Cascade Failure
1. Downstream payment service degrades or hangs.
2. Checkout service worker threads block on slow HTTP responses.
3. Thread pools and connection pools exhaust.
4. Upstream API gateway and health-check endpoints freeze.
5. Entire platform suffers total outage from single downstream failure.

## Circuit Breaker Mental Model
Monitors downstream failures. Trips open to block calls when threshold exceeded. Fails fast in nanoseconds. Safely probes downstream before restoring traffic.

## CLOSED
Normal operational state. All calls routed downstream. Successes reset failure count. Failures increment counter. When `failures >= FailureThreshold`, breaker transitions to `OPEN`.

## OPEN
Tripped state. All incoming calls fail immediately with `ErrCircuitOpen`. Zero network calls dispatched. Downstream dependency spared from traffic. Cooldown timer (`OpenTimeout`) begins.

## HALF-OPEN
Probe state. Triggered after `OpenTimeout` elapses. Allows limited probe calls (`HalfOpenMaxCalls`):
- Probe succeeds: resets failure count and transitions to `CLOSED`.
- Probe fails: transitions back to `OPEN` and restarts cooldown timer.

## Architecture

```text
Checkout Service
       │
       ▼
Circuit Breaker Proxy
  ├── [CLOSED]   ──► Call Payment Client ──► Fake Payment Server
  ├── [OPEN]     ──► Fail Fast (ErrCircuitOpen)
  └── [HALF-OPEN]──► Allow Canary Probe
```

## Without Circuit Breaker
Every checkout call hits downstream. When downstream hangs or fails:
- Requests block until HTTP client timeout expires (e.g. 100ms in demo).
- Downstream receives 100% of load despite being broken.
- Total latency equals `N * timeout`.

## With Circuit Breaker
- Initial requests hit downstream and record errors.
- Breaker trips to `OPEN` once `FailureThreshold` reached.
- Subsequent requests fail fast in microseconds without network I/O.
- Downstream call counter freezes.

## Recovery
- Breaker sleeps during `OpenTimeout` cooldown.
- Enters `HALF_OPEN`.
- Dispatches canary probe.
- If probe returns 200 OK, breaker transitions to `CLOSED` and resumes normal throughput.
- If probe fails, breaker flips back to `OPEN`.

## Timeout vs Retry vs Circuit Breaker
- **Timeout**: Enforces maximum wait time per request (prevents infinite hang).
- **Retry**: Re-issues failed request assuming transient glitch. Can trigger retry storms.
- **Circuit Breaker**: Halts all requests to downstream once failure rate spikes. Silences retries.

*Note on Lab Timeouts*: The timeouts used in this lab (100ms HTTP timeout, 300ms breaker cooldown) are illustrative for fast automated testing and demonstration. Production timeouts must be tuned to actual service SLA, P99 latency, and downstream recovery profiles.

## Fallback
Alternative response strategy when circuit is `OPEN`:
- Return cached data.
- Return default placeholder.
- Queue request for asynchronous background execution.
- Never use silent fallback for critical state-altering mutations (e.g., balance debit).

## Bulkhead
Isolates resources (thread pools, connection limits, CPU slices) per dependency.
- Circuit breaker stops traffic to failing dependency.
- Bulkhead ensures failure of dependency A does not exhaust threads needed for dependency B.

## Observability
Key production metrics:
- `circuit_state`: Metric gauge (0=CLOSED, 1=OPEN, 2=HALF_OPEN).
- `circuit_open_count`: Counter incremented on state trip.
- `rejected_call_count`: Requests dropped via fail-fast.
- `failure_count`: Recent consecutive or windowed failure counter.
- `dependency_latency`: P50, P95, P99 call durations.

## Running the Lab

Execute demo:
```bash
go run ./cmd/demo
```

## Running Tests

Run all unit and concurrency tests:
```bash
go test -v ./...
```

Run race detector:
```bash
go test -race ./...
```

## Expected Behavior
```text
=== SCENARIO 1: WITHOUT CIRCUIT BREAKER (SLOW DEPENDENCY) ===
request=1 result=timeout/error duration=100ms
request=2 result=timeout/error duration=101ms
request=3 result=timeout/error duration=101ms
downstream_calls=3 (all requests blocked and hit downstream)

=== SCENARIO 2: WITH CIRCUIT BREAKER (FAIL-FAST ON DOWN DEPENDENCY) ===
request=1 result=payment_error   state=CLOSED duration=~1ms
request=2 result=payment_error   state=CLOSED duration=~200µs
request=3 result=payment_error   state=OPEN   duration=~150µs
request=4 result=circuit_open    state=OPEN   duration=<1µs
request=5 result=circuit_open    state=OPEN   duration=<1µs
downstream_calls=3 (downstream calls stopped once OPEN)

=== SCENARIO 3: RECOVERY (HALF_OPEN -> CLOSED) ===
waiting for cooldown...
payment server recovered to HEALTHY. Current CB state=HALF_OPEN
sending probe request...
probe result: err=<nil>, state after probe=CLOSED
subsequent regular request: state=CLOSED

=== SCENARIO 4: FAILED RECOVERY (HALF_OPEN -> OPEN AGAIN) ===
dependency still DOWN. Current CB state=HALF_OPEN
probe result: err=true, state after failed probe=OPEN
next request: circuit_open, state=OPEN
```

## Failure Modes
- **Threshold too low**: Breaker trips on momentary blips, dropping valid traffic.
- **Probe flood**: Too many parallel probe calls during HALF_OPEN overwhelm recovering service.
- **Ignoring error types**: Counting 4xx (client errors) as circuit-breaking faults rather than restricting to 5xx/timeouts.

## CMMS Example
Flow: `Create Invoice` -> `Generate PDF` -> `Send WhatsApp`
- Design rule: WhatsApp down != Invoice creation down.
- Pattern:
  - Save invoice synchronously in primary DB.
  - Push notification job to background message queue (RabbitMQ/Kafka).
  - Worker calls WhatsApp API wrapped in Circuit Breaker + exponential backoff retry.
  - Breaker trips OPEN if WhatsApp API fails. Queue retains message; retries pause.

## PPOB Example
Flow: `User` -> `Order` -> `Payment Gateway` -> `Provider Pulsa` -> `WhatsApp Notification`
- Critical vs non-critical dependencies:
  - **Payment Gateway**: Synchronous critical dependency. If down, fail fast; do not deduct balance or complete order.
  - **Provider Pulsa**: Asynchronous retryable dependency with strict idempotency key.
  - **WhatsApp**: Asynchronous non-critical dependency. Failure must not abort order fulfillment.

## Key Takeaways
- Circuit Breaker does not heal a broken dependency.
- Circuit Breaker stops caller self-destruction by failing fast and preserving system threads/memory.
- Downstream service gets idle recovery time without request flooding.

## Sources
- Martin Fowler, *Circuit Breaker*: https://martinfowler.com/bliki/CircuitBreaker.html
- Microsoft Learn, *Circuit Breaker Pattern*: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- AWS Builders Library, *Timeouts, retries, and backoff with jitter*: https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/
