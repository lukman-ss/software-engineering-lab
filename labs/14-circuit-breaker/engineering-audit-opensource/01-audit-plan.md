# Engineering Audit Plan

Target Lab: `labs/14-circuit-breaker` — Go circuit breaker demonstrating fail-fast, OPEN/HALF-OPEN/CLOSED transitions, and recovery.

Implementation Files:
- `internal/circuitbreaker/circuit_breaker.go` — core state machine (`Execute`, `State`, `New`, `advanceLocked`).
- `internal/circuitbreaker/circuit_breaker_test.go` — unit + concurrency tests.
- `internal/payment/client.go` — HTTP payment client (downstream dependency stub).
- `internal/payment/fake_server.go` — `httptest` fake with HEALTHY/SLOW/DOWN modes + request counter.
- `internal/checkout/service.go` — Checkout Service composing payment client + circuit breaker.
- `cmd/demo/main.go` — runnable demo (4 scenarios).
- `tests/integration_test.go` — end-to-end integration test.
- `go.mod` — module `circuitbreaker`, Go 1.22.

Tests:
- `internal/circuitbreaker/circuit_breaker_test.go` (13 tests).
- `tests/integration_test.go` (1 sub-test group).

Executable/Demo:
- `go run ./cmd/demo` (4 scenarios: no-breaker slow, fail-fast, recovery, failed recovery).

Approved Research Inputs: (out of scope for this engineering-only audit — ignored).

Main Claims To Verify:
1. Initial state is CLOSED.
2. Successes reset failure count in CLOSED.
3. `FailureThreshold` consecutive failures trip to OPEN.
4. OPEN state returns `ErrCircuitOpen` without calling downstream (fail-fast).
5. OPEN→HALF_OPEN transition occurs after `OpenTimeout` cooldown.
6. Successful HALF_OPEN probe transitions BACK to CLOSED and resumes traffic.
7. Failed HALF_OPEN probe transitions back to OPEN; next request still fails fast.
8. Half-open probe concurrency throttling via `HalfOpenMaxCalls`.
9. Thread safety under concurrent `Execute`/`State` (race detector clean).
10. Recovery after dependency becomes healthy.
11. README "Expected Behavior" matches `go run ./cmd/demo` output.

Commands To Run:
- `go vet ./...`
- `go test -count=1 -race ./...`
- `go test -count=1 -race -v ./...`
- `go run ./cmd/demo`
- `gofmt -l .`

Primary Risks:
- Timing-based transitions (cooldown) introduce flakiness potential under load.
- README "Expected Behavior" snippet may not match actual error-message format.
- `Execute` releases lock between decision and downstream call — verify throttle correctness in HALF_OPEN under concurrency.
- `State()` mutates state via `advanceLocked` (lock-free reads of `openedAt`) — verify no TOCTOU in transition to OPEN.
