# Documentation vs Code Verification

## 1. README vs Code Comparison

| README Claim | Code / Demo Reality | Assessment |
| :--- | :--- | :--- |
| **NOT NULL (`23502`)**: Ensures required columns cannot accept null/empty values | Enforced in `engine.go:51-56` & `engine.go:128-130`. Validated in `TestNotNullConstraints` & Demo Step [1]. | PASS |
| **CHECK Constraints (`23514`)**: Evaluates boolean predicate logic on row data | Enforced in `engine.go:60-68` & `engine.go:133-135`. Validated in `TestCheckConstraints` & Demo Step [2]. | PASS |
| **UNIQUE Constraints (`23505`)**: Enforces single occurrence across rows and prevents concurrent read-then-write race conditions | Enforced in `engine.go:78-83`. Validated in `TestUniqueConstraint`, `TestConcurrentRegistration_Safe_EnforcesUniqueness`, & Demo Step [5]. | PASS |
| **FOREIGN KEY Constraints (`23503`)**: Enforces referential integrity preventing orphan rows | Enforced in `engine.go:138-140`. Validated in `TestForeignKeyConstraint` & Demo Step [3]. | PASS |
| **PARTIAL UNIQUE INDEX**: Demonstrates conditional uniqueness (`WHERE deleted_at IS NULL`) | Enforced in `engine.go:73-77`. Validated in `TestPartialUniqueIndex` & Demo Step [4]. | PASS |
| **Run Commands**: `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | Executed directly during audit. All commands compiled and executed cleanly without error. | PASS |

## 2. Research Claims vs Implementation Comparison

| Research Claim | Implementation Reality | Assessment |
| :--- | :--- | :--- |
| PostgreSQL SQLSTATE codes (`23502`, `23514`, `23505`, `23503`) should be defined and tested | Defined in `internal/dberr/errors.go` and verified in tests & demo outputs. | PASS |
| Application-level check-then-act fails under concurrency, DB constraints succeed | Demonstrated in `UnsafeStore` vs `SafeStore` in `store_test.go` and `cmd/demo/main.go`. | PASS |
| Soft delete + re-registration requires partial unique index `WHERE deleted_at IS NULL` | Simulated in `engine.go` active email map index management during insert and soft delete operations. | PASS |

## Discrepancies Found

None. `DOC_CODE_MISMATCH`, `TEST_CLAIM_MISMATCH`, and `RESEARCH_IMPLEMENTATION_MISMATCH` are all clean with zero occurrences.
