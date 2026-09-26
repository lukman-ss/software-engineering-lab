# Source Map

## State Machine and Core Concepts

Research:
- research/03-core-concepts.md (Three States, Key Configuration, Failure Detection, Thread Safety, What Circuit Breaker Does NOT Do)
- research/05-circuit-states.md (State transitions diagram)
- research/09-failure-modes.md (Premature Opening, Thundering Herd, Infinite Open, 4xx False Positives)

Implementation:
- internal/circuitbreaker/circuit_breaker.go (State type, Config, CircuitBreaker struct, Execute method, checkStateTransitionLocked)

Tests:
- internal/circuitbreaker/circuit_breaker_test.go (Tests 1-16 covering all states and transitions)

Master Draft Sections:
- Problem, Why This Matters, Mental Model, Core Concept

---

## Fail-Fast Behavior

Research:
- research/03-evidence.md (Evidence 3 — Circuit Breaker prevents operations likely to fail)
- research/03-core-concepts.md (Fail Fast purpose)

Implementation:
- internal/circuitbreaker/circuit_breaker.go:78-81 (OPEN state returns ErrCircuitOpen immediately)

Tests:
- internal/circuitbreaker/circuit_breaker_test.go Test 5 (OPEN calls fail fast)
- internal/circuitbreaker/circuit_breaker_test.go Test 6 (DOWNSTREAM function not executed when OPEN)

Demo:
- cmd/demo/main.go Scenario 2 (fail-fast in nanoseconds)
- engineering/03-execution-result.md (demo output)

Master Draft Section:
- Core Concept, Recovery / Rollback

---

## Recovery and HALF_OPEN

Research:
- research/03-core-concepts.md (HALF_OPEN state description)
- research/03-evidence.md (Evidence 4 — Half-Open limits traffic)
- research/05-circuit-states.md (State transitions)

Implementation:
- internal/circuitbreaker/circuit_breaker.go:64-71 (checkStateTransitionLocked — OPEN to HALF_OPEN)
- internal/circuitbreaker/circuit_breaker.go:82-121 (HALF_OPEN handling with HalfOpenMaxCalls)

Tests:
- internal/circuitbreaker/circuit_breaker_test.go Test 7 (cooldown moves toward HALF_OPEN)
- internal/circuitbreaker/circuit_breaker_test.go Test 8 (successful HALF_OPEN probe closes)
- internal/circuitbreaker/circuit_breaker_test.go Test 9 (failed HALF_OPEN probe opens again)
- internal/circuitbreaker/circuit_breaker_test.go Test 12 (HalfOpenMaxCalls > 1 limits probes)
- tests/integration_test.go (cooldown and recovery integration test)

Demo:
- cmd/demo/main.go Scenario 3 (Recovery) and Scenario 4 (Failed Recovery)
- engineering/03-execution-result.md (demo output)

Master Draft Sections:
- Recovery / Rollback

---

## Concurrency Safety

Research:
- research/03-core-concepts.md (Thread Safety — Mutex, Atomic, Channel patterns)

Implementation:
- internal/circuitbreaker/circuit_breaker.go:28 (sync.Mutex on CircuitBreaker struct)
- internal/circuitbreaker/circuit_breaker.go (all mutations guarded by cb.mu.Lock())

Tests:
- internal/circuitbreaker/circuit_breaker_test.go Test 11 (50 concurrent goroutines)
- engineering-audit/03-test-audit.md (Race Detector PASS — zero data races)

Engineering Audit:
- engineering-audit/02-code-audit.md Finding 2 (Safe concurrent access)

Master Draft Sections:
- Implementation

---

## Timeout vs Retry vs Circuit Breaker

Research:
- research/06-timeout-retry-backoff.md (Timeout bounds execution, Retry risks, Synergy)
- research/03-evidence.md (Evidence 5 — Retry storms from unbounded retries)

Master Draft Sections:
- Common Mistakes, Production Considerations

---

## Cascade Failure

Research:
- research/03-evidence.md (Evidence 1 — blocked requests hold critical resources)
- research/10-final-research.md (Circuit Breakers isolate network blast radius)

README:
- README.md (Cascade Failure section)

Master Draft Sections:
- Problem, Why This Matters

---

## Demo Scenarios (4 Case Studies)

Demo:
- cmd/demo/main.go (Scenarios 1-4)

Execution Results:
- engineering/03-execution-result.md (full demo output)
- engineering-audit/03-test-audit.md (demo verification)

Tests:
- tests/integration_test.go (Integration tests)
- internal/circuitbreaker/circuit_breaker_test.go (Unit tests 1-16)

Master Draft Sections:
- Case Study, What the Tests Prove

---

## Production Considerations and Warnings

Research:
- research/03-core-concepts.md (Key Configuration Parameters — illustrative ranges)
- research-revision/03-revision-result.md (Remaining Risks on illustrative values)

Engineering:
- engineering/02-implementation-notes.md (Known Limitations, Implementation-Specific Choices)
- engineering-revision/03-revision-result.md (No remaining risks after revision)

Engineering Audit:
- engineering-audit-opensource/06-verdict.md (LOW warning: README illustrative timing values)

README:
- README.md (Note on Lab Timeouts, Observability section, Failure Modes)

Master Draft Sections:
- Production Considerations, Common Mistakes

---

## What Circuit Breaker Does NOT Do

Research:
- research/03-core-concepts.md (Does not heal failing service, does not retry, does not replace timeouts)

Master Draft Sections:
- Key Takeaways

---

## Fallback

Research:
- research/03-core-concepts.md (Graceful Degradation purpose)

README:
- README.md (Fallback section)

Master Draft Sections:
- Key Takeaways

---

## Bulkhead

Research:
- research/02-sources.md (Source 9 — Bulkhead Pattern)

README:
- README.md (Bulkhead section)

---

## Key Takeaways Summary

Master Draft:
- content/05-key-takeaways.md (derived from research, implementation, tests, and demo)

---

## Full Source Inventory

### Research Files:
- research/01-research-plan.md
- research/02-sources.md (13 sources)
- research/03-evidence.md
- research/03-core-concepts.md
- research/04-contradictions.md
- research/05-circuit-states.md
- research/06-timeout-retry-backoff.md
- research/06-open-questions.md
- research/09-failure-modes.md
- research/10-final-research.md

### Research Audit:
- research-audit/01-audit-plan.md
- research-audit/02-source-audit.md
- research-audit/03-claim-audit.md
- research-audit/04-contradictions.md
- research-audit/05-code-audit.md
- research-audit/06-gaps.md
- research-audit/07-verdict.md (APPROVED)

### Engineering:
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md

### Engineering Audit:
- engineering-audit/01-audit-plan.md
- engineering-audit/02-code-audit.md
- engineering-audit/03-test-audit.md
- engineering-audit/04-docs-vs-code.md
- engineering-audit/05-gaps.md
- engineering-audit/06-verdict.md (APPROVED)

### Implementation Files:
- internal/circuitbreaker/circuit_breaker.go
- internal/circuitbreaker/circuit_breaker_test.go
- internal/payment/client.go
- internal/payment/fake_server.go
- internal/checkout/service.go
- cmd/demo/main.go

### Tests:
- tests/integration_test.go

### Revisions:
- research-revision/01-revision-plan.md
- research-revision/02-changes-made.md
- research-revision/03-revision-result.md
- engineering-revision/01-revision-plan.md
- engineering-revision/02-changes-made.md
- engineering-revision/03-revision-result.md