# Engineering Audit Plan

Target Lab: labs/14-circuit-breaker
Implementation Files:
- internal/circuitbreaker/circuit_breaker.go
- internal/payment/client.go
- internal/payment/fake_server.go
- internal/checkout/service.go
- cmd/demo/main.go
Tests:
- internal/circuitbreaker/circuit_breaker_test.go (16 tests)
- tests/integration_test.go (2 tests, 3 subtests)
Executable/Demo: cmd/demo (4 scenarios)
Approved Research Inputs: SKIPPED per PIPELINE OVERRIDE (implementation and tests only)
Main Claims To Verify:
1. CLOSED -> OPEN on failures >= FailureThreshold
2. OPEN fails fast with ErrCircuitOpen, zero downstream calls
3. OPEN -> HALF_OPEN after OpenTimeout
4. HALF_OPEN probe success -> CLOSED; probe failure -> OPEN
5. Success in CLOSED resets failure count
6. HalfOpenMaxCalls throttles excess probes
7. Concurrency safety (mutex + generation counter)
8. Slow-dependency timeout trips breaker; fail-fast <10ms
Commands To Run:
- go build ./...
- go test -count=1 -v ./...
- go test -race -count=1 ./...
- go run ./cmd/demo
Primary Risks:
- Stale in-flight requests corrupting new state (generation handling)
- Half-open probe flood under concurrency
- Panic in guarded fn leaking lock or corrupting state
- Docs claiming behavior tests do not prove
