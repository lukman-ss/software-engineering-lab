# Engineering Revision Log

## Revision 1

Audit Issue: GAP 1 (MISSING_TEST) - `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` planned in `01-design.md` was omitted.
Severity: LOW
Files Changed:
- `labs/27-database-constraints/internal/store/store_test.go`
Action: Added `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` asserting race condition under concurrency in unsafe store implementation.
Verification: `go test -v -run TestConcurrentRegistration_Unsafe_SuffersRaceCondition ./internal/store` (PASS)
Status: RESOLVED

## Revision 2

Audit Issue: GAP 2 (UNUSED_CODE) - `UnsafeStore` was unused scaffolding and called indexing insert.
Severity: LOW
Files Changed:
- `labs/27-database-constraints/internal/engine/engine.go`
- `labs/27-database-constraints/internal/store/store.go`
Action: Added `InsertUserUnsafe` to engine to bypass index constraints. Connected `UnsafeStore.RegisterUser` to `InsertUserUnsafe` with artificial interleaving yield to demonstrate read-then-write race.
Verification: `go test -v ./...` & `go test -race ./...` (PASS)
Status: RESOLVED
