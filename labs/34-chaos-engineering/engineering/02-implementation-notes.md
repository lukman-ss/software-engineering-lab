# Implementation Notes

## Files Added
- `go.mod`: Module definition for Go 1.22+.
- `internal/fault/injector.go`: Fault injection primitive supporting latency injection and forced error return.
- `internal/circuitbreaker/circuitbreaker.go`: State machine for Closed, Open, Half-Open states with fallback support.
- `internal/monitor/monitor.go`: Steady-state metrics collector evaluating error rates against configured threshold.
- `internal/experiment/runner.go`: Chaos experiment manager with automated monitoring and auto-abort switch.
- `cmd/demo/main.go`: Interactive executable demonstrating baseline traffic, mitigated chaos fault (via Circuit Breaker + Fallback), unmitigated fault leading to auto-abort, and post-experiment recovery.
- `tests/chaos_test.go`: Unit tests for fault injection, circuit breaker state transitions, graceful degradation, experiment auto-abort, and concurrent execution under race detector.

## Core Design Decisions
- **Standard Library Implementation**: Built exclusively using standard library primitives (`sync`, `sync/atomic`, `context`, `time`) without external dependencies.
- **Thread Safety**: Used `sync.RWMutex` for fault configuration and atomic operations for steady-state counter updates to ensure concurrent test safety under `go test -race`.

## Implementation-Specific Choices
- **Immediate Neutralization on Abort**: When an experiment auto-aborts due to steady-state metric breach, `injector.Clear()` is executed synchronously inside `terminate()` to strictly restrict blast radius.
- **Circuit Breaker Fallback**: Circuit breaker `Execute` accepts an optional fallback function so that failures count as steady-state metric success when graceful degradation is active.

## Known Limitations
- Metrics use cumulative counting rather than a sliding time window (simplified for demonstration lab scope).
- Fault injection is in-memory rather than network layer/proxy level.

## Trade-offs
- Simple threshold-based error rate metric chosen over complex percentile (p99) latency distribution monitoring.

## What Is Demonstrated
- Hypothesis-driven fault injection.
- Circuit Breaker preventing cascading failure.
- Graceful degradation through fallback mechanisms.
- Automated abort when steady-state metrics breach bounds.
- Full recovery of system health after chaos experiment termination.

## What Is Not Demonstrated
- Dynamic production traffic routing/canary deployment injection.
- Distributed tracing (OpenTelemetry) integration.
