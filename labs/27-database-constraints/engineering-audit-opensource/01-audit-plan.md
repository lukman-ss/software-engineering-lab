# Engineering Audit Plan

Target Lab: labs/27-database-constraints

Implementation Files:
- internal/engine/engine.go (constraint engine, 161 lines)
- internal/store/store.go (UnsafeStore / SafeStore, 84 lines)
- internal/dberr/errors.go (SQLSTATE taxonomy + mapping, 100 lines)
- internal/model/model.go (User, Order, 20 lines)
- cmd/demo/main.go (demo, 99 lines)

Tests:
- internal/store/store_test.go (8 tests, 261 lines)

Executable/Demo: cmd/demo (go run ./cmd/demo)

Approved Research Inputs: SKIPPED per pipeline override (implementation + tests only).

Main Claims To Verify:
1. NOT NULL rejects missing email/username/user_id (23502)
2. CHECK rejects age<18, bad status, total_cents<=0 (23514)
3. UNIQUE rejects duplicate email, exactly-1-wins under 20-goroutine race (23505)
4. FOREIGN KEY rejects orphan order (23503)
5. PARTIAL UNIQUE allows email reuse after soft-delete, blocks second active duplicate
6. UnsafeStore (app-level check only) suffers duplicates under concurrency
7. SQLSTATE taxonomy + domain error mapping (23502/23503/23505/23514)

Commands To Run:
- go vet ./...
- go test -v ./... (fresh, -count=1)
- go test -race -count=1 -v ./...
- go run ./cmd/demo

Primary Risks:
- Tests assert err != nil but never assert SQLSTATE codes on engine/store path
- MapToDomainError uses fmt.Errorf without %w — machine-readable code lost
- ctx accepted but never honored (no timeout/cancel path)
- Unsafe-race test relies on timing (sleep 1ms) for >1 duplicates
- Design-doc package paths (internal/db, internal/errors) differ from code
