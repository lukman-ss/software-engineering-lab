# Documentation vs Code Audit

## Comparison Matrix

| Claim in Documentation | Implementation File | Test / Verification | Status |
| :--- | :--- | :--- | :--- |
| NOT NULL (`23502`) enforcement on `email`, `username`, `user_id` | `internal/engine/engine.go:51-56, 128-130` | `TestNotNullConstraints` | MATCH |
| CHECK (`23514`) enforcement on `age >= 18`, `status`, `total_cents > 0` | `internal/engine/engine.go:59-68, 133-135` | `TestCheckConstraints` | MATCH |
| UNIQUE (`23505`) enforcement & race condition prevention | `internal/engine/engine.go:78-83` | `TestUniqueConstraint`, `TestConcurrentRegistration_Safe_EnforcesUniqueness` | MATCH |
| FOREIGN KEY (`23503`) integrity | `internal/engine/engine.go:138-140` | `TestForeignKeyConstraint` | MATCH |
| Partial Unique Index (`WHERE deleted_at IS NULL`) | `internal/engine/engine.go:71-77, 102-120` | `TestPartialUniqueIndex` | MATCH |
| Demo and test commands | `README.md:17, 24, 32` | Verified executable via bash | MATCH |

## Findings
- No discrepancies detected between README.md, engineering design documents, demo output, and underlying code implementation.
- All SQLSTATE error codes and descriptions match standard PostgreSQL semantics.
