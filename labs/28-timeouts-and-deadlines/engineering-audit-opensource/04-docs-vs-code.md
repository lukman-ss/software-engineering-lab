# Docs vs Code

Target Lab: labs/28-timeouts-and-deadlines

## README.md (`labs/28-timeouts-and-deadlines/README.md`)
Sections:
- High-level: context deadline propagation, timeout budgets, exponential backoff with full jitter, circuit breaking, idempotency deduplication. ✓ matches components.
- Components Implemented: 4 bullets. ✓ matches implementation files.
- Execution: `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`. ✓ executable.

### Findings
- README omits known limitations and implementation-specific choices (these live in `engineering/02-implementation-notes.md`). Acceptable for a lab README; not required.
- README claims no research/content audit here (per override); engineering notes describe full jitter/backoff. README does not give formula, engineering notes do. ✓ consistent.

## engineering/01-design.md
- Concept To Prove, Expected Behavior, Architecture, Test Strategy, Execution Plan present. ✓.
- "Backoff calculation (`base * 2^attempt + jitter`)" slightly misstates code (`base * 2^(attempt-1)` then jitter). Low-severity doc imprecision; formula described in 02-implementation-notes.md is correct.

## engineering/02-implementation-notes.md
- Full jitter formula stated correctly. ✓.
- "Context deadline propagation: Sub-calls inherit parent context bounds, context timeout overrides local function timeout budget." Code: `context.WithTimeout(ctx, budget)` — parent bound is stricter, child gets min(parent,budget). ✓ correct semantics.
- "Circuit Breaker Integration: transitions into Open when reaching consecutive failures threshold, immediately rejecting subsequent retries." ✓ matches code.
- "Idempotency Deduplication: Memory store with TTL mapping request tokens to responses." ✓ matches code; note: TTL does not auto-evict (see code audit Finding 4).
- Known Limitations: non-persistence stated. TTL non-eviction NOT mentioned. ⚠ gap.

## engineering/03-execution-result.md
- Claims build SUCCESS, race PASS, demo output. Verified by execution: identical output. ✓.
- No fabricated result observed.

## research (out of scope)
- Not audited per pipeline override. Referenced only for alignment sanity check:
  - README and engineering notes align with implementation.
  - No claim in README/engineering contradicts code.

## Mismatches

### DOC_CODE_MISMATCH
1. engineering/01-design.md: "Backoff calculation (base * 2^attempt + jitter)" — code uses `base * 2^(attempt-1)`. Minor; corrected in 02. Severity LOW.

### TEST_CLAIM_MISMATCH
2. Success Criteria #2: "Race detector passes cleanly under concurrent operations." Passes — but concurrent operations tested only for idempotency (assert-less) and circuit (none). Tests pass, but concurrency correctness is not fully asserted. Severity MEDIUM.

### RESEARCH_IMPLEMENTATION_MISMATCH
3. None observed within audit scope (research excluded).

## Summary
Documentation largely matches code and demo. One low-severity formula typo; one medium concern that race test exists and passes but concurrency correctness is weakly asserted.
