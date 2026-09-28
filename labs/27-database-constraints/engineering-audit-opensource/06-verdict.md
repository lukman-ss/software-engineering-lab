# Engineering Audit Verdict

Target Lab: labs/27-database-constraints
Audit Date: Mon Sep 28 2026

## Summary

Code Files Reviewed:
- internal/dberr/errors.go
- internal/model/model.go
- internal/engine/engine.go
- internal/store/store.go
- internal/store/store_test.go
- cmd/demo/main.go

Tests Reviewed: 8 unit tests in internal/store/store_test.go
Commands Executed:
- go build ./... → PASS
- go test -v ./... → PASS (8/8)
- go test -race ./... → PASS (no data races)
- go run ./cmd/demo → PASS (see output below)

Failures: None (all commands succeeded).
Warnings: See Gap Analysis (05-gaps.md) and Code Audit (02-code-audit.md).

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: PASS (constraints prevent races, enforce invariants)
Documentation Accuracy: WARNING (see Findings in 04-docs-vs-code.md)

## Blocking Issues (HIGH/CRITICAL)

None.

## Non-Blocking Issues (MEDIUM/LOW)

### From Code Audit (02-code-audit.md):
- Finding 7: UnsafeStore genuinely exhibits the read-then-write race (intended) — WARNING, timing-dependent test but correctly shows the bug.
- Finding 8: Data-race safety — PASS.
- Finding 9: time.Sleep(1ms) in UnsafeStore is a timing-dependent flakiness source — WARNING, MEDIUM.
- Finding 10: SoftDeleteUser API signature awkward — WARNING, LOW.
- Finding 11: Cross-index (full-UNIQUE vs partial-UNIQUE) path inconsistent — WARNING, MEDIUM.
- Finding 12: No rollback / transaction semantics — PASS, LOW.
- Finding 13: No recovery / corruption handling — NOT_APPLICABLE.

### From Test Audit (03-test-audit.md):
- GAP 1: MISSING_TEST: MapToDomainError output never asserted — MEDIUM.

### From Docs vs Code (04-docs-vs-code.md):
- GAP 2: DOC_CODE_MISMATCH: engineering/01-design.md package tree does not match reality — MEDIUM.
- GAP 3: DOC_CODE_MISMATCH: engineering/03-execution-result.md omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition — MEDIUM.
- GAP 5: IMPLEMENTATION_OVERCLAIM: PRIMARY KEY constraint not demonstrated or tested — MEDIUM.

### From Gaps (05-gaps.md):
- MISSING_TEST: MapToDomainError output never asserted → MEDIUM
- DOC_CODE_MISMATCH: engineering/01-design.md package tree does not match reality → MEDIUM
- DOC_CODE_MISMATCH: engineering/03-execution-result.md omits TestConcurrentRegistration_Unsafe_SuffersRaceCondition → MEDIUM
- MISSING_EDGE_CASE: Mixed registration paths (full-UNIQUE vs partial-UNIQUE) not tested / documented → MEDIUM
- IMPLEMENTATION_OVERCLAIM: PRIMARY KEY constraint not demonstrated (PK overwrites on caller-supplied ID) → MEDIUM
- RACE_CONDITION: UnsafeStore race dependent on time.Sleep (flakiness risk if removed) → LOW

## Required Revisions

Based on the gaps above, the following revisions are recommended (non-blocking for approval):

1. Add a test asserting MapToDomainError output for each SQLState (fixes MISSING_TEST gap).
2. Update engineering/01-design.md to reflect the actual package names (internal/engine, model, store, dberr) or add a mapping note (fixes DOC_CODE_MISMATCH).
3. Amend engineering/03-execution-result.md to list all 8 tests including TestConcurrentRegistration_Unsafe_SuffersRaceCondition (fixes DOC_CODE_MISMATCH/TEST_CLAIM_MISMATCH).
4. Address the mixed-registration-path edge case: either document the namespaces are separate, unify the uniqueness check, or add a test demonstrating the behavior (fixes MISSING_EDGE_CASE).
5. Clarify PRIMARY KEY behavior: either add a check and test for caller-supplied PK collisions, or downgrade the claim in 01-design.md that PK enforces uniqueness (fixes IMPLEMENTATION_OVERCLAIM).
6. (Optional) Consider replacing the time.Sleep in UnsafeStore with a proper barrier if test purity is desired; otherwise accept as a deliberate demonstration crutch (low severity).

## Final Status

APPROVED_WITH_WARNINGS

**Justification**:  
- Code compiles.  
- All tests pass, including race detector.  
- Demo runs and matches claimed output (1 success / 49 rejections under concurrency, constraint violations mapped correctly).  
- No HIGH or CRITICAL blocking issues (no fake results, no broken core behavior, no data races).  
- The warnings are largely documentation or test coverage gaps that do not invalidate the claimed behavior: the lab correctly demonstrates that database-level constraints prevent race conditions and enforce invariants as specified.  
- A technical writer can trust the implementation matches its public claims after noting the minor doc/test improvements above.