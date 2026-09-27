# Gaps

## GAP-1: Store tests assert presence of error, not SQLSTATE code
Type: MISSING_TEST
Severity: MEDIUM
Location: internal/store/store_test.go:16-112 (TestNotNullConstraints, TestCheckConstraints, TestUniqueConstraint, TestForeignKeyConstraint)
Evidence: All constraint tests check `err == nil` only. Only TestErrorClassification checks codes, but it constructs errors directly via `dberr.New*` and never exercises the engine/store path. So `23502/23503/23505/23514` propagation through `SafeStore` is unasserted.
Impact: Core rejection behavior proven, code-level taxonomy mapping unproven at integration path.

## GAP-2: MapToDomainError severs error chain
Type: UNHANDLED_ERROR
Severity: MEDIUM
Location: internal/dberr/errors.go:83-100
Evidence: `MapToDomainError` uses `fmt.Errorf("...")` without `%w`. Returned domain error loses `*ConstraintError`, so `errors.As` / `IsConstraintViolation` fails on the value callers actually receive from `SafeStore`.
Impact: Human message correct; machine-readable code lost. No test covers this.

## GAP-3: context.Context accepted but never honored
Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/store/store.go:24,59,69,78
Evidence: `ctx` parameter unused in all four methods. No timeout/cancel test.
Impact: No functional effect for in-memory demo; API misrepresents cancellability.

## GAP-4: Design doc package map does not match code
Type: DOC_CODE_MISMATCH
Severity: LOW
Location: engineering/01-design.md:32-42 vs actual tree
Evidence: Doc cites `internal/db`, `internal/errors`, `internal/service`, `internal/domain`. Code has `internal/engine`, `internal/dberr`, `internal/model`, `internal/store`. No `internal/service` exists.
Impact: Docs-only; no behavioral claim tied to missing packages.

## GAP-5: Unsafe race demonstration is timing-dependent
Type: UNVERIFIED_RESULT
Severity: LOW
Location: internal/store/store.go:38 + store_test.go:210-239
Evidence: Race window enforced by `time.Sleep(1ms)`. Passed in audit runs (`go test -race -count=1` PASS) and demo shows safe-side only. Result is real but relies on artificial yield; without sleep the assertion `count > 1` could flake under different schedulers.
Impact: Illustrative only; SafeStore guarantee does not depend on it.

No FAKE_DEMO, FAKE_BENCHMARK, RACE_CONDITION, or BROKEN_IMPLEMENTATION found. No HIGH/CRITICAL gaps.
