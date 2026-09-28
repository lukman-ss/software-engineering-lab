# Docs vs Code Audit

## Sources Compared
- README.md
- engineering/01-design.md (approved research input)
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
- Source code (internal/*, cmd/demo/main.go)
- Tests (internal/*_test.go, tests/integration_test.go)

## Findings

### DOC_CODE_MISMATCH: None detected

All major advertised behaviors are present and align with documentation:
1. **Context deadline propagation & timeout budgets** – `internal/deadline` implements `ExecuteWithBudget` returning context error on timeout.
2. **Exponential backoff with full jitter** – `internal/retry` computes `sleep = rand.Float64() * min(MaxBackoff, BaseBackoff * 2^(attempt-1))`.
3. **Circuit breaker (CLOSED/OPEN/HALF_OPEN)** – `internal/circuit` state machine transitions per thresholds and cooldown.
4. **Idempotency deduplication** – `internal/idempotency` store returns cached responses within TTL.

### TEST_CLAIM_MISMATCH: None detected

Test suite validates:
- Deadline propagation success, timeout, and parent-context inheritance.
- Retry success, max-attempt exhaustion, context cancellation.
- Circuit breaker state transitions and half-open behavior.
- Idempotent store get/set, expiration, concurrent access.
- Integration test combining retry + circuit breaker and idempotent retry.

### RESEARCH_IMPLEMENTATION_MISMATCH: None detected

Design doc (engineering/01-design.md) specifies:
- **Concept To Prove**: same four patterns.
- **Expected Behavior**: items 1-4 match implementation.
- **Architecture**: maps packages 1-4.
- **Implementation Decisions**: stdlib only, in-memory store limitation noted.
- **What Is Demonstrated**: propagation, jitter backoff, circuit breaker, deduplication — all shown in demo and tests.

No exaggeration or false claims found.

### Minor discrepancies (not mismatches)

- **README line 7**: "timeout budgets" – code uses `budget` param but does not track aggregate spend across nested calls; only per-call limit. Interpretation acceptable.
- **engineering/02-implementation-notes.md line 17**: cites AWS/Google for full jitter – implementation matches described formula.
- **engineering/02-implementation-notes.md line 18**: "context timeout overrides local function timeout budget" – `ExecuteWithBudget` uses min(parent ctx, budget) via nested `WithTimeout`, correct.

## Conclusion
Documentation accurately reflects implementation. No DOC_CODE_MISMATCH, TEST_CLAIM_MISMATCH, or RESEARCH_IMPLEMENTATION_MISMATCH found.