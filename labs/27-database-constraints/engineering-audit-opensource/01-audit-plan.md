# Engineering Audit Plan

Target Lab:
labs/27-database-constraints

Implementation Files:
- internal/engine/engine.go
- internal/store/store.go
- internal/dberr/errors.go
- internal/model/model.go
- cmd/demo/main.go

Tests:
- internal/store/store_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md (Database Constraints & Data Integrity)

Main Claims To Verify:
1. NOT NULL constraints (23502) prevent null values
2. CHECK constraints (23514) validate row-level predicates
3. UNIQUE constraints (23505) enforce uniqueness and prevent read-then-write race conditions
4. FOREIGN KEY constraints (23503) enforce referential integrity
5. PARTIAL UNIQUE INDEX enables soft-delete re-registration while maintaining active uniqueness
6. Concurrency safety: under load, exactly one write succeeds for conflicting unique keys
7. Error handling maps SQLSTATE codes to user-facing messages

Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
1. Misleading implementation comment: UnsafeStore claims to bypass constraints but does not
2. Error classification not verified in integration tests (only constructors tested)
3. Domain error mapping loses SQLSTATE information (no programmatic retry path for 40001)
4. Coarse-grained locking vs research's row-level lock mechanism (documented trade-off)
5. UnsafeStore/vulnerable path lacks test coverage