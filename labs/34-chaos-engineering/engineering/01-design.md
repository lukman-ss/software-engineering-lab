# Engineering Design

Target Lab: labs/34-chaos-engineering
Research Status: APPROVED

## Concept To Prove
Chaos Engineering principles: defining steady state metrics, injecting controlled faults (latency, service errors), limiting blast radius, automatically aborting experiments on metric degradation, and verifying resilient system behavior (Circuit Breaker & Graceful Degradation).

## Expected Behavior
1. Under normal operation, steady state error rate remains 0% and response latency stays within bounds.
2. Injecting fault (e.g. downstream service latency or 500 internal errors) causes target service resilience mechanisms (Circuit Breaker) to open or fall back gracefully.
3. If steady state error rate or latency exceeds threshold, Chaos Experiment automatically triggers Abort Signal, neutralizing injected fault.

## Failure Scenario
Without Circuit Breaker / Graceful Degradation or Abort mechanism:
- Injected downstream fault causes cascading delay/failures and unmitigated steady state violation.
- Uncontrolled chaos fault runs continuously, violating blast radius boundaries.

## Success Criteria
- Automated injection of latency/errors into downstream dependency.
- Continuous steady-state health checks.
- Auto-abort when steady-state metrics breach defined thresholds.
- Circuit breaker state transition to Open under fault condition with fallback response.
- All unit and concurrency tests pass with Go race detector (`go test -race ./...`).

## Architecture
- `FaultInjector`: Intercepts calls to simulate latency or errors with configurable probabilities.
- `SteadyStateMonitor`: Periodically/continuously measures success rate and latency.
- `ExperimentEngine`: Manages experiment lifecycle (Hypothesis -> Inject -> Monitor -> Auto-Abort / Neutralize).
- `ResilientClient`: Client wrapping target service with Circuit Breaker and Fallback logic.

## Components
- `pkg/fault`: Fault injection primitives (Latency, Error).
- `pkg/monitor`: Steady-state metric calculation and abort triggers.
- `pkg/circuitbreaker`: Minimal Circuit Breaker state machine (Closed, Open, Half-Open).
- `pkg/experiment`: Chaos experiment control runner.
- `cmd/demo`: Executable demonstration of experiment run, failure injection, circuit breaking, and auto-abort.

## Test Strategy
- Unit tests for Fault Injector (probabilistic triggering and exact overrides).
- Unit tests for Circuit Breaker state transitions.
- Integration test for Chaos Experiment lifecycle (Start -> Inject Fault -> Detect Degradation -> Auto Abort -> Return to Steady State).
- Concurrency test with `go test -race` for thread-safe fault injection and metric tracking.

## Execution Plan
1. Implement `go.mod` (Go 1.22+ standard library).
2. Implement core components in `internal/` or `pkg/`.
3. Implement `cmd/demo/main.go`.
4. Implement tests in `tests/` or alongside packages.
5. Run tests, race detector, demo execution.
6. Record execution results and notes.

## Implementation Decisions
- **Standard Library Only**: Use Go standard library (`sync`, `sync/atomic`, `time`, `context`, `net/http`) to avoid external dependency overhead.
- **In-Memory Metrics**: Use thread-safe sliding window counters for steady state monitoring.
