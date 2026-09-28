# Docs vs Code

Comparison of README, engineering notes, research claims, code, tests, and demo.

## README Claims vs Code
1. "Steady-state health metric monitoring" -> internal/monitor implements atomic counters and error-rate threshold. PASS
2. "Injected downstream service latency and errors" -> internal/fault implements latency + forced error injection. PASS
3. "Resilience validation via Circuit Breakers and Fallback mechanisms" -> internal/circuitbreaker implements state machine and fallback. PASS
4. "Blast radius control with automated experiment abort upon metric degradation" -> internal/experiment auto-aborts and calls injector.Clear(). PASS

All four README claims are implemented and proven by demo (see demo output in 06-verdict).

## Engineering Design vs Code
- Design doc references `pkg/fault`, `pkg/monitor`, `pkg/circuitbreaker`, `pkg/experiment`.
- Actual code lives under `internal/...` not `pkg/...`.
- Design doc Test Strategy mentions "probabilistic triggering and exact overrides" for Fault Injector.
- Fault injector has no probabilistic errorRate path (field declared, unused). Only deterministic forced error.
- Design doc claims "sliding window counters" in Architecture; engineering notes document that cumulative counters are used intentionally.

Assessment: DOC_CODE_MISMATCH (package path wording) — LOW.
Assessment: IMPMENTATION_OVERCLAIM (probabilistic fault) — LOW. This is an over-statement of the design plan, not of the implementation claims in README.

## Research vs Code
Research report findings (steady state, blast radius, circuit breaker, graceful degradation) all reflected in implementation. No RESEARCH_MISMATCH.

## Test vs Claim
README claims concurrency tests; `go test -race ./...` passes. Tests cover state transitions, fallback, auto-abort, concurrency. PASS. No TEST_CLAIM_MISMATCH.

## Demo vs Execution Result
Recorded `03-execution-result.md` output matches actual `go run ./cmd/demo` output exactly (experiment COMPLETED, then ABORTED on 33.33% error rate, injector neutralized, recovery to CLOSED). FAKE_DEMO not present.
