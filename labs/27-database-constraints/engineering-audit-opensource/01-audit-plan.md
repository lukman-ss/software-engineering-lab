# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- internal/engine/engine.go
- internal/store/store.go
- internal/model/model.go
- internal/dberr/errors.go
- cmd/demo/main.go
Tests:
- internal/store/store_test.go (8 tests)
Executable/Demo:
- cmd/demo/main.go (`go run ./cmd/demo`)
Approved Research Inputs: SKIPPED per pipeline override (implementation+tests only)
Main Claims To Verify:
1. NOT NULL (23502) rejects empty email/username/user_id
2. CHECK (23514) rejects age<18, bad status, total_cents<=0
3. UNIQUE (23505) rejects duplicate email, exactly 1 win under concurrency
4. FOREIGN KEY (23503) rejects orphan order, accepts valid ref
5. PARTIAL UNIQUE INDEX allows email reuse after soft-delete, blocks 2nd active
6. UnsafeStore suffers duplicates, SafeStore enforces uniqueness
7. SQLSTATE taxonomy + domain error mapping
8. Demo output real, README matches code
Commands To Run:
- go vet ./...
- go test -v ./... (clean cache)
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- In-memory simulator vs real PG semantics (empty-string NULL proxy, coarse lock)
- Unsafe race test sleep-dependent flakiness
- Mapped domain errors lose SQLSTATE code
- Full vs partial index modes isolated, never co-enforced
- Stale design-doc paths/names, stale execution-result log
