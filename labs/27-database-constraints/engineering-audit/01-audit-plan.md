# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- `internal/dberr/errors.go`
- `internal/model/model.go`
- `internal/engine/engine.go`
- `internal/store/store.go`
- `cmd/demo/main.go`
Tests:
- `internal/store/store_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`
Main Claims To Verify:
1. NOT NULL constraint enforcement (`SQLSTATE 23502`) on mandatory fields (`email`, `username`, `user_id`).
2. CHECK constraint validation (`SQLSTATE 23514`) on boundary invariants (`age >= 18`, status enum, `total_cents > 0`).
3. UNIQUE constraint enforcement (`SQLSTATE 23505`) preventing duplicates and race conditions.
4. FOREIGN KEY referential integrity (`SQLSTATE 23503`) blocking orphan records.
5. Partial unique index (`WHERE deleted_at IS NULL`) allowing re-registration after soft-delete while enforcing uniqueness across active records.
6. Concurrency safety: SafeStore enforces exactly 1 success under concurrent inserts; UnsafeStore demonstrates read-then-write race condition duplicate leaks.
7. Error classification and domain error mapping via SQLSTATE codes.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race condition or deadlocks in simulated storage engine under `-race`.
- Mismatch between README documentation and actual implemented constraint rules / error mappers.
- Fragile test assertions or unverified demo output.
