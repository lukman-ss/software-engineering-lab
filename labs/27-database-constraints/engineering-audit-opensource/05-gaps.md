# Gap Analysis

## Identified Gaps

| Type | Location | Severity | Description |
|------|----------|----------|-------------|
| MISSING_TEST | internal/store/store.go (UnsafeStore), internal/store/store_test.go | MEDIUM | `UnsafeStore` is declared and described as demonstrating a vulnerable read-then-write pattern, but no test exercises it. The documentation claims "comprehensive unit and concurrency race condition tests", yet the vulnerable path has zero coverage. |
| MISSING_TEST | internal/store/store_test.go | MEDIUM | Tests verify `err == nil` / `err != nil` and constructor-level SQLSTATE, but **never** verify: <br> • The SQLSTATE of an error returned from a store method (`RegisterUser`, `CreateOrder`) <br> • That `MapToDomainError` produces the correct human-readable message for each SQLSTATE class <br> • End-to-end validation of the "Always check SQLSTATE, not error text" research principle. |
| UNHANDLED_ERROR | internal/dberr/errors.go:35-37, 83-100 | LOW | `MapToDomainError` drops the original `ConstraintError` type (and thus the inspectable `SQLSTATE`). After mapping, callers cannot call `IsConstraintViolation` on the returned error, making programmatic handling by SQLSTATE (e.g., retry on 40001) impossible. |
| IMPLEMENTATION_OVERCLAIM | internal/store/store.go:40-42 | MEDIUM | Comment: "Engine insert bypassing unique constraints (simulating unconstrained table)". In reality: `InsertUser(u, true)` still enforces the partial unique index for active users (`DeletedAt == nil`). There is no code path that actually bypasses constraint enforcement. The UnsafeStore does not simulate an unconstrained table. |
| RACE_CONDITION | internal/engine/engine.go:14 | LOW | The engine uses a coarse-grained `sync.RWMutex` (table-level lock) rather than the row-level B-tree page latches described in research/05-report.md:50-55. This is a documented trade-off (engineering/02-implementation-notes.md:25) and behaviorally equivalent for all verified cases. |
| FAKE_DEMO | *none* | LOW | The demo's "SQLSTATE Taxonomy Verification" (`cmd/demo/main.go:98`) constructs a synthetic error rather than asserting the taxonomy on an actual store-level duplicate error. While not fabricating a result, it does not prove the taxonomy works end-to-end through the store pipeline. The output is real, but the verification step is cosmetic. |
| MISSING_EDGE_CASE | internal/engine/engine.go:89-107 | LOW | No test for `SoftDeleteUser` with a non-existent user ID (returns plain `fmt.Errorf`). |
| DEAD_CODE | internal/dberr/errors.go:11, 16-17 | LOW | Unused SQLSTATE constants: `SQLStateRestrictViolation` (23001), `SQLStateExclusionViolation` (23P01), `SQLStateSerializationFail` (40001). The latter is mentioned in research as retry-trigger but never produced by the engine. |

## Summary of Allowed Gap Types Used

- `MISSING_TEST`: 2 instances (UnsafeStore untested, SQLSTATE verification on store returns untested)
- `UNHANDLED_ERROR`: 1 instance (`MapToDomainError` loses SQLSTATE inspectability)
- `IMPLEMENTATION_OVERCLAIM`: 1 instance (UnsafeStore falsely claims to bypass constraints)
- `RACE_CONDITION`: 1 instance (coarse-grained vs row-level lock, LOW severity)
- `FAKE_DEMO`: 1 instance (demo SQLSTATE line is cosmetic)
- `MISSING_EDGE_CASE`: 1 instance (SoftDeleteUser missing-user path untested)
- `DEAD_CODE`: 1 instance (unused SQLSTATE constants)

Total: 7 gaps across 6 types.