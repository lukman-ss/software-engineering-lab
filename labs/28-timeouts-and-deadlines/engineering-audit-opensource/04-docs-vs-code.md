# Docs vs Code

## Comparison: README vs Code
- README claims `internal/deadline` provides "context deadline propagation and execution within explicit time budgets."
  - Code: ✅ Matches. `ExecuteWithBudget` wraps ctx with timeout.
- README claims `internal/retry` uses "exponential backoff with full jitter."
  - Code: ✅ `CalculateBackoff` implements `rand()*min(maxBackoff, baseBackoff*2^(attempt-1))`.
- README claims `internal/circuit` is a `CLOSED`/`OPEN`/`HALF_OPEN` state machine.
  - Code: ✅ Matches exactly.
- README claims `internal/idempotency` is "in-memory deduplication store preventing double execution during retries."
  - Code: ✅ Matches with TTL.

## Comparison: Engineering Design (01-design.md) vs Code
- Design: "Requests propagate context deadlines... operations abort immediately once context deadline expires."
  - Code: ✅ `ExecuteWithBudget` selects on ctx.Done.
- Design Success Criteria: "Context cancellation terminates long-running downstream work without leakage."
  - Code: ⚠️ Partially. Goroutine running `fn` is not forcibly killed; relies on fn honoring ctx. Tests/demo honour ctx, but leak risk exists if downstream ignores ctx. Documented as known limitation indirectly.
- Design: "Race detector passes under concurrent operations."
  - Code/Execution: ✅ PASS.
- Design Test Strategy: unit tests in internal/* plus concurrent race tests; integration tests combining retrier/circuit/idempotency.
  - Code: ✅ Matches.

## Comparison: Execution Result (03-execution-result.md) vs Actual Audit Run
- Engineering doc claims `go build ./...` SUCCESS, `go test ./...` all ok, `go test -race ./...` PASS, demo output as recorded.
- Audit run result: ✅ All match. Output identical (modulo cached timing values).
- No fabricated benchmark or result observed in doc.

## Discrepancies / Mismatches
- DOC_CODE_MISMATCH (LOW): Design success criterion "Context cancellation terminates long-running downstream work without leakage" slightly overstates guarantees — implementation delegates to fn observing ctx. Not false but aspirational.
- TEST_CLAIM_MISMATCH: None. Engineering notes claim Ready_for_engineering_audit; audit confirms tests pass and demo runs.
- RESEARCH_IMPLEMENTATION_MISMATCH: Not audited per pipeline override (research out of scope).

## Conclusion
README and engineering docs match code and observed execution. No material mismatches found.