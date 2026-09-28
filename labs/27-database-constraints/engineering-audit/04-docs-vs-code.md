# Documentation vs Code Verification

## Mismatch Inspection

1. **README Claims**:
   - Lists 5 implemented constraint categories: NOT NULL (`23502`), CHECK (`23514`), UNIQUE (`23505`), FOREIGN KEY (`23503`), PARTIAL UNIQUE INDEX (`WHERE deleted_at IS NULL`).
   - Code Verification: `internal/engine/engine.go` and `internal/dberr/errors.go` explicitly declare and implement these exact 5 constraints and SQLSTATE codes.
   - Status: MATCH

2. **Demo Output Claims**:
   - `cmd/demo/main.go` runs all 5 scenarios plus 50-goroutine stress test and SQLSTATE taxonomy check.
   - Execution Verification: `go run ./cmd/demo` executes without error and matches documented outputs verbatim.
   - Status: MATCH

3. **Research Claims vs Implementation**:
   - Approved research focuses on PostgreSQL error state mapping and concurrency safety.
   - Code matches research specifications without excess or unapproved scope.
   - Status: MATCH
