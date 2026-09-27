# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- internal/model/model.go
- internal/dberr/errors.go
- internal/engine/engine.go
- internal/store/store.go
- cmd/demo/main.go
Tests:
- internal/store/store_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- research/01-plan.md
- research-audit/07-verdict.md (Status: APPROVED)
Main Claims To Verify:
1. Application-only validation without database constraints fails under concurrent requests (read-then-write race condition creates duplicates).
2. Database-level UNIQUE constraints atomically prevent duplicates under concurrency, returning SQLSTATE 23505 (unique_violation).
3. NOT NULL constraints reject missing mandatory attributes (SQLSTATE 23502).
4. CHECK constraints enforce row-scoped boolean invariants (SQLSTATE 23514).
5. FOREIGN KEY constraints prevent orphan records referring to non-existent parents (SQLSTATE 23503).
6. Partial Unique Indexes (`WHERE deleted_at IS NULL`) allow reuse of unique keys after soft-deletion while maintaining uniqueness among active records.
7. Error mapping translates low-level SQLSTATE codes into structured domain-layer errors.
Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions or deadlocks inside the simulated engine mutex hierarchy.
- In-memory simulation deviating from declared SQLSTATE codes or relational constraint semantics.
- Weak test assertions failing to actually reproduce concurrency collisions in `UnsafeStore`.
