# Engineering Revision Plan

Target Lab: labs/27-database-constraints
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. `engineering/01-design.md` specified `TestConcurrentRegistration_Unsafe_SuffersRaceCondition`, but only the safe concurrency enforcement test was implemented in `store_test.go`.
2. `UnsafeStore` in `internal/store/store.go` was unused scaffolding.

## Files To Change
- `labs/27-database-constraints/internal/store/store.go`
- `labs/27-database-constraints/internal/store/store_test.go`
- `labs/27-database-constraints/internal/engine/engine.go`

## Tests To Add/Modify
- Add `TestConcurrentRegistration_Unsafe_SuffersRaceCondition` in `internal/store/store_test.go`.

## Validation Commands
```bash
cd labs/27-database-constraints
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
