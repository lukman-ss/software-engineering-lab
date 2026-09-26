# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Audit Scope: implementation + tests only (no research/content audit, no code modification)
Output Dir: labs/27-database-constraints/engineering-audit-opensource/

Implementation Files:
- internal/engine/engine.go (161 lines)
- internal/store/store.go (84 lines)
- internal/model/model.go (20 lines)
- internal/dberr/errors.go (100 lines)
- cmd/demo/main.go (99 lines)
- go.mod (module: github.com/lukman/software-engineering-lab/labs/27-database-constraints)

Tests:
- internal/store/store_test.go (8 tests):
  1. TestNotNullConstraints
  2. TestCheckConstraints
  3. TestUniqueConstraint
  4. TestForeignKeyConstraint
  5. TestPartialUniqueIndex
  6. TestConcurrentRegistration_Safe_EnforcesUniqueness
  7. TestConcurrentRegistration_Unsafe_SuffersRaceCondition
  8. TestErrorClassification

Executable/Demo:
- cmd/demo/main.go (`go run ./cmd/demo`)

Approved Research Inputs: out of scope per pipeline override (not audited here)

Main Claims To Verify:
1. NOT NULL rejects missing email/username/user_id (23502)
2. CHECK rejects age<18, bad status, total_cents<=0 (23514)
3. UNIQUE rejects duplicate email, exactly 1 winner under concurrency (23505)
4. FOREIGN KEY rejects orphan orders (23503)
5. PARTIAL UNIQUE allows email reuse after soft-delete, blocks double-active
6. Error taxonomy maps 23502/23503/23505/23514 to domain errors
7. UnsafeStore demonstrates read-then-write duplicate race; SafeStore prevents it
8. README accurately describes all implemented constraints and run commands

Commands To Run:
- go build ./...
- go test -v -count=1 ./...
- go test -race -count=1 ./...
- go test -race -run TestConcurrentRegistration_Unsafe -count=10 ./internal/store/
- go run ./cmd/demo

Primary Risks:
- Design doc (engineering/01-design.md) references package names that do not match actual package layout
- Execution result log (engineering/03-execution-result.md) shows 7 tests but claims 8
- MapToDomainError loses ConstraintError type info, breaking IsConstraintViolation on mapped errors
- InsertUser allows caller-supplied non-zero ID to silently overwrite existing records
- Stale execution log may mislead Technical Writer about test coverage
