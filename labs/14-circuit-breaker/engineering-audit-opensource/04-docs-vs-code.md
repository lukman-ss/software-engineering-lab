# Docs vs Code Audit

## README vs Code

### Claim: CLOSED accumulates failures; trips to OPEN at FailureThreshold; success resets count
- README (CLOSED section): "When `failures >= FailureThreshold`, breaker transitions to OPEN."
- Code: circuit_breaker.go:145 `if cb.failureCount >= cb.config.FailureThreshold` then trip to OPEN; line 152 `cb.failureCount = 0` on success.
- Match: PASS

### Claim: OPEN fail-fasts with ErrCircuitOpen, zero downstream calls
- README (OPEN): "All incoming calls fail immediately with ErrCircuitOpen."
- README (Expected Behavior): "downstream_calls=3 (downstream calls stopped once OPEN)"
- Code: circuit_breaker.go:78-80 returns ErrCircuitOpen without invoking fn. FakeServer.RequestCount proves zero calls (tested + demo).
- Match: PASS

### Claim: HALF_OPEN triggered by OpenTimeout cooldown; limited probes via HalfOpenMaxCalls; success->CLOSED, failure->OPEN
- README (HALF OPEN): "Triggered after OpenTimeout elapses. Allows limited probe calls (HalfOpenMaxCalls)."
- Code: checkStateTransitionLocked (lines 65-66) OPEN->HALF_OPEN; Execute HALF_OPEN branch (lines 82-121) enforces halfOpenCalls limit and transitions.
- Match: PASS

### Claim: Demo scenarios (4 scenarios) reproduce fail-fast and recovery
- README (Expected Behavior): 4 scenarios with states/durations.
- Code/demo main.go: identical 4 scenarios, identical output structure.
- Match: PASS for behavior/state/durations-as-timings. NOTE on timing: README durations (e.g. "499.25µs", "100.208µs") differ from actual run ("735.875µs", "182.167µs"). README explicitly marked illustrative via "Note on Lab Timeouts". This is a non-fabricated, timing-dependent variance, not a behavioral mismatch.
- Match: DOC_CODE_MISMATCH (severity LOW) — illustrative timing only.

## Research vs Code

### Research 05 (circuit-states.md): state diagram and semantics
- CLOSED: failures increment; successes reset. OPEN: blocked + ErrCircuitOpen. HALF_OPEN: canary probe; success->CLOSED, fail->OPEN.
- Code implements exactly this.
- Match: PASS

### Research 10-final + research-audit verdict: isolates blast radius; fails fast; async fallback via queues
- Code implements state machine matching the synthesis. Async-queue/fallback/bulkhead described as architectural patterns, not asserted as implemented.
- Match: PASS (no overclaim)

### Engineering design 01 / 02 (implementation scope)
- "simple consecutive failure counting rather than sliding time-window" — documented as known limitation. No sliding-window claim in README.
- "Concurrency during HALF_OPEN permits up to HalfOpenMaxCalls probes; excess requests fail fast." — implemented and tested (subtest 12).
- Match: PASS (scope correctly stated, no overclaim)

## Claim Mismatches Identified

1. DOC_CODE_MISMATCH — README "Expected Behavior" Scenario 2 timing values are illustrative and do not match actual run timing (timing-dependent). LOW severity. Not fabricated.

No TEST_CLAIM_MISMATCH, no RESEARCH_IMPLEMENTATION_MISMATCH, no FAKE_DEMO, no FAKE_BENCHMARK.
