# Engineering Audit Plan

Target Lab: labs/14-circuit-breaker

Implementation Files:
- internal/circuitbreaker/circuit_breaker.go — core state machine (mutex-guarded)
- internal/payment/client.go — HTTP payment client with timeout
- internal/payment/fake_server.go — controllable fake downstream (HEALTHY / SLOW / DOWN)
- internal/checkout/service.go — checkout entrypoint wrapping call in circuit breaker
- cmd/demo/main.go — runnable demonstration across 4 scenarios

Tests:
- internal/circuitbreaker/circuit_breaker_test.go — 16 unit subtests (state transitions, panic, concurrency)
- tests/integration_test.go — end-to-end with real fake server

Executable/Demo:
- cmd/demo (run via `go run ./cmd/demo`)

Approved Research Inputs:
- research/01-plan.md
- research/02-sources.md
- research/03-core-concepts.md
- research/04-cascade-failure.md
- research/05-circuit-states.md
- research/06-timeout-retry-backoff.md
- research/07-fallback-bulkhead.md
- research/08-observability.md
- research/09-failure-modes.md
- research/10-final-research.md
- research-audit/07-verdict.md (research APPROVED)

Main Claims To Verify:
1. CLOSED: all calls route downstream; failures accumulate; success resets counter; trips to OPEN at FailureThreshold.
2. OPEN: fail-fast with ErrCircuitOpen, zero downstream network calls, cooldown (OpenTimeout) elapses.
3. HALF_OPEN: cooldown expiry transitions OPEN -> HALF_OPEN; limited probes (HalfOpenMaxCalls); success -> CLOSED; failure -> OPEN.
4. Concurrency safety under Go race detector.
5. Demo reproduces cascade-failure mitigation, fail-fast, recovery, and re-trip scenarios.
6. README "Expected Behavior" output reflects actual run.

Commands To Run:
- go build ./...
- go vet ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Mutex discipline under concurrent HALF_OPEN probes (probe-counting race / double-decrement).
- Panic recovery correctness (state restoration, double-unlock, panic propagation).
- README illustrative timing vs actual duration values.
- Observability metrics documented as "recommended / omitted" — must stay scoped, not claimed as implemented.
