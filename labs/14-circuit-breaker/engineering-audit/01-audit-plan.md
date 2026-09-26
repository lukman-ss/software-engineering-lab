# Engineering Audit Plan

Target Lab: labs/14-circuit-breaker
Implementation Files:
- `internal/circuitbreaker/circuit_breaker.go`
- `internal/checkout/service.go`
- `internal/payment/client.go`
- `internal/payment/fake_server.go`

Tests:
- `internal/circuitbreaker/circuit_breaker_test.go`
- `tests/integration_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/10-final-research.md`
- `research/05-circuit-states.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Circuit Breaker transitions through CLOSED -> OPEN -> HALF-OPEN -> CLOSED / OPEN states.
2. In CLOSED state: successes reset failure count; failures increment counter until threshold.
3. In OPEN state: incoming calls fail fast immediately with `ErrCircuitOpen` without dispatching downstream I/O.
4. After cooldown (`OpenTimeout`): breaker lazily advances to HALF-OPEN.
5. In HALF-OPEN state: throttles concurrent probe requests (`HalfOpenMaxCalls`). Success closes breaker; failure reopens breaker and resets timer.
6. Concurrency safety: thread-safe state inspection and transitions; generation counters prevent stale in-flight results from corrupting future state transitions.
7. Panics inside probe execution cleanly transition state and avoid deadlock / stuck breaker.
8. Demo output matches README scenario outputs.

Commands To Run:
```bash
go test -count=1 -v ./...
go test -count=1 -race ./...
go run ./cmd/demo
```

Primary Risks:
- Race conditions during state transitions under high concurrency.
- Stale asynchronous responses corrupting generation state.
- Unhandled panic leaving mutex locked or state permanently throttled.
- Discrepancies between README output and actual demo execution.
