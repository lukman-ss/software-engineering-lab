# Docs vs Code

Scope: README, engineering notes, code, tests, demo. Research/content skipped per PIPELINE OVERRIDE.

## README vs Code: PASS
- CLOSED/OPEN/HALF_OPEN semantics match circuit_breaker.go:91-157.
- `failures >= FailureThreshold` matches code `>=`.
- `ErrCircuitOpen` fail-fast, zero downstream calls: proven by TestOpenDoesNotCallDownstream + integration count==2 + demo downstream_calls=3.
- OpenTimeout 300ms / HTTP 100ms in demo match README note on illustrative timeouts.
- Scenario 1 (no breaker, no state field), scenario 2 (trip on 3rd, fail-fast ns), scenario 3 (HALF_OPEN->CLOSED), scenario 4 (HALF_OPEN->OPEN): all reproduced in audit run, output structure matches README Expected Behavior.
- Observability section explicitly marked "omitted in this minimal lab implementation": honest, no overclaim.
- No benchmarks claimed, none present. No FAKE_BENCHMARK.

## Engineering notes vs Code: WARNING
- DOC_CODE_MISMATCH (LOW): 02-implementation-notes.md references `checkStateTransitionLocked`; code function is `advanceLocked` (circuit_breaker.go:83). Behavior description correct, name stale.
- DOC_CODE_MISMATCH (LOW): 03-execution-result.md scenario 1 shows `state=CLOSED` on without-breaker requests. Code CheckoutWithoutBreaker sets HasBreaker=false, omits state. Audit demo run confirms no state field. README is correct; execution-result stale.
- DOC_CODE_MISMATCH (LOW): 03-execution-result.md lists 11 unit tests + 1 integration test. Current suite has 16 unit tests + 2 integration tests (added: TestSuccessInClosedResetsFailures, TestHalfOpenThrottlesExcessCalls, TestPanicInHalfOpenCleansUpState, TestTrailingInFlightRequestDoesNotCorruptNewState, TestInterleavedConcurrentTransitions, TestCircuitBreakerSlowDependencyTimeoutTrips). Stale list; actual run proves superset passes.

## Tests vs Claims: PASS
- All README transition claims have direct unit + integration coverage (see 03-test-audit.md).
- No TEST_CLAIM_MISMATCH. No IMPLEMENTATION_OVERCLAIM in README.
- Known limitations (consecutive counting, no sliding window, mutex contention) honestly documented in 02-implementation-notes.md.

## Demo: PASS (real, not fake)
- `go run ./cmd/demo` executed in audit, 4 scenarios complete, timings (~100ms timeouts, ns fail-fast) plausible, downstream_counts correct.
