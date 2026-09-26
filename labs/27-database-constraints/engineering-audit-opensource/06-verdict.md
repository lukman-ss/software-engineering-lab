# Engineering Audit Verdict

Target Lab:
labs/27-database-constraints

Audit Date:
2026-09-26

## Summary

Code Files Reviewed:
- internal/engine/engine.go
- internal/store/store.go
- internal/dberr/errors.go
- internal/model/model.go
- cmd/demo/main.go

Tests Reviewed:
- internal/store/store_test.go

Commands Executed:
- go test -v ./... → PASS (8/8)
- go test -race ./... → PASS (clean, no data race)
- go run ./cmd/demo → ran, output matches recorded claims; concurrency result (50 goroutines → 1 success, 49 errors, "Database Integrity Intact: true") reproduced

Failures:
- None. Build, tests, race detector, and demo all succeeded.

Warnings:
- UnsafeStore comment falsely states it bypasses constraints; it does not (MEDIUM).
- No test asserts SQLSTATE or domain messages on store-returned errors (MEDIUM).
- Domain error mapping discards the ConstraintError type, losing SQLSTATE inspectability (LOW).

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: WARNING
Documentation Accuracy: PASS

## Blocking Issues

1. None that would REJECT. All README claims (the five constraint types, partial unique index, concurrency race prevention, SQLSTATE taxonomy) are reproduced by execution.

## Non-Blocking Issues

1. UnsafeStore misleadingly claims to bypass constraints — it does NOT. The safe-vs-unsafe contrast is not actually demonstrated. (IMPLEMENTATION_OVERCLAIM)
2. No test asserts SQLSTATE or `*ConstraintError` type on store method returns; tests only check `err == nil`. (MISSING_TEST)
3. `MapToDomainError` discards the original `ConstraintError`, so `IsConstraintViolation` cannot be applied to store results; programmatic SQLSTATE handling (e.g., retry-on-40001 per research) is impossible on store return values. (UNHANDLED_ERROR)
4. Demo's "SQLSTATE Taxonomy Verification" uses a synthetic error instead of an actual store-level failure. (FAKE_DEMO, cosmetic)
5. Coarse-grained table lock vs research's row-level locking (documented, behaviorally equivalent). (RACE_CONDITION, LOW)
6. Unused SQLSTATE constants 23001/23P01/40001; no engine path emits 40001. (DEAD_CODE, LOW)
7. `Order.Status` field never validated/used. (DEAD_CODE, LOW)

## Required Revisions

Before publication the Technical Writer should be aware of:
1. The UnsafeStore "vulnerable" path is illustrative only — it is not actually vulnerable and is untested. Remove or genuinely model the unconstrained-table behavior so the contrast is real.
2. Add tests asserting the SQLSTATE (`dberr.IsConstraintViolation`) and domain message content for at least one store-returned error per constraint class (23502, 23503, 23505, 23514).
3. Either preserve the `ConstraintError` through the store layer (e.g., wrap rather than replace) so SQLSTATE remains inspectable, or document that classification is intentionally hidden from callers.
4. Replace the synthetic demo taxonomy line with an assertion that derives the SQLSTATE from an actual duplicate registration.

## Final Status

APPROVED_WITH_WARNINGS

Rationale: compilation passes, all required tests pass, race detector is clean, the demo runs and reproduces every README claim. The core behavior (constraints prevent concurrent duplicate registration → exactly one success) is proven. The issues are incomplete coverage and misleading comments rather than broken or fabricated core behavior. The verdict is raised to APPROVED_WITH_WARNINGS (not APPROVED) because the SQLSTATE-verification gap and the false UnsafeStore claim represent incomplete proof and a documented mismatch that the Technical Writer should resolve before the lab is presented as a reference.