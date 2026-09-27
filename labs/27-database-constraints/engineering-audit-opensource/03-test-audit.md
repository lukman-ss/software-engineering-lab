# Test Audit — labs/27-database-constraints

Files reviewed: internal/store/store_test.go (261 lines, 8 tests). Commands executed live.

## Coverage Matrix

| Constraint            | Test                                  | PASS live | Edge / negative covered |
|-----------------------|---------------------------------------|-----------|------------------------|
| NOT NULL              | TestNotNullConstraints                | yes       | missing email, missing username, missing user_id |
| CHECK                 | TestCheckConstraints                  | yes       | underage, invalid status, total<=0 |
| UNIQUE                | TestUniqueConstraint                  | yes       | exact duplicate email  |
| FOREIGN KEY           | TestForeignKeyConstraint              | yes       | orphan user_id, valid insert |
| Partial Unique Index  | TestPartialUniqueIndex                | yes       | active dup blocked, re-use post-soft-delete, re-block, soft-deleted row insert allowed |
| Concurrency (safe)    | TestConcurrentRegistration_Safe       | yes       | exactly 1 success / 19 rejected, final count=1 |
| Concurrency (unsafe)  | TestConcurrentRegistration_Unsafe     | yes       | >1 rows inserted, race proven |
| Error classification  | TestErrorClassification               | yes       | all four SQLSTATE codes via IsConstraintViolation |

## Live Execution Results

Command: go test -v ./...   Result: PASS (8 tests, 0.00s) exit 0
Command: go test -race ./... Result: PASS (cached) exit 0

## Findings

### Finding T1  PASS — Happy / failure / edge covered
Each constraint has both a happy path (valid insert) and a failure path (violation). NOT NULL has 3 negative cases; CHECK has 3; FK has 2; partial index walks 5 transitions (insert, block, soft-delete, re-insert, re-block, soft-delete-bypass).

### Finding T2  WARNING — No timeout/lifecycle tests
Success criteria mentions race detector; no timeout, cancellation, or `context` cancellation tests. Acceptable gap: no I/O or blocking waits exist except the engine lock. Severity LOW.

### Finding T3  WARNING — Error-mapping not tested end-to-end
`MapToDomainError` is used by SafeStore but no test asserts on its output. Tests classify raw `dberr` constructor errors only. SQLSTATE taxonomy proven at `dberr` layer, not through store path. Severity MEDIUM.

### Finding T4  WARNING — No dedicated test listing
Test names do not exactly match engineering/01-design.md Test Strategy (design lists `TestNotNullConstraint`, `TestCheckConstraint`, `TestUniqueConstraint`, `TestForeignKeyConstraint`, `TestErrorMapping`; actual names pluralized/with `s` and `TestErrorClassification`). Test logic matches design intent; only naming differs. Severity LOW.

### Finding T5  WARNING — Tests reference nonexistent `tests/...` path
Design §4 "Write unit and concurrency tests (`internal/store/...` and `tests/...`)" mentions a `tests/` directory. It does not exist; tests live only in `internal/store/store_test.go`. Coverage is complete, location differs from doc. Severity LOW.

### Finding T6  PASS — Concurrency race proven
Unsafe test reliably produces >1 rows (proves race); safe test reliably produces exactly 1 (proves elimination). Both pass under race detector. Strongest evidence in the lab.