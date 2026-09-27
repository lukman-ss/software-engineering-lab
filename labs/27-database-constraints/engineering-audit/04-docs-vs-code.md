# Docs vs Code Audit

Target Lab: labs/27-database-constraints

## Comparison Matrix

| Documented Claim / Element | Code / Test Implementation | Alignment Status | Notes |
|---|---|---|---|
| NOT NULL (`23502`) | `internal/engine/engine.go:51-56`, `TestNotNullConstraints` | MATCH | Implemented and verified |
| CHECK (`23514`) | `internal/engine/engine.go:58-69`, `TestCheckConstraints` | MATCH | Verified for age, status, order total |
| UNIQUE (`23505`) | `internal/engine/engine.go:79-83`, `TestUniqueConstraint`, `TestConcurrentRegistration_Safe_EnforcesUniqueness` | MATCH | Full full-table unique index |
| FOREIGN KEY (`23503`) | `internal/engine/engine.go:127-148`, `TestForeignKeyConstraint` | MATCH | Verified referential check against users table |
| PARTIAL UNIQUE INDEX (`WHERE deleted_at IS NULL`) | `internal/engine/engine.go:71-77`, `TestPartialUniqueIndex` | MATCH | Re-registration after soft deletion validated |
| Concurrency Stress Demo | `cmd/demo/main.go:59-95` | MATCH | Real 50-worker test in demo and 20-worker in tests |
| SQLSTATE Error Mapping | `internal/dberr/errors.go:83-99` | MATCH | Verified error code checks and domain translations |

## Discrepancy Findings
- `DOC_CODE_MISMATCH`: None observed.
- `TEST_CLAIM_MISMATCH`: None observed.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None observed.
