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

Tests:
- internal/store/store_test.go (261 lines, 8 tests)

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

Commands To Run:
- go build ./...
- go test ./...
- go test -v ./...
- go test -race -count=1 ./...
- go test -run TestConcurrentRegistration_Unsafe -count=10 ./internal/store/
- go run ./cmd/demo

Primary Risks:
- Logical race demo flakiness (sleep-based, not data race)
- Mapped domain errors losing SQLSTATE type info
- PK duplicate-ID overwrite (no explicit-ID guard)
- Stale execution log omitting 1 test
