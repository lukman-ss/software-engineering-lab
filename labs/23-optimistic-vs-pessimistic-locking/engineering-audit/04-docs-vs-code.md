# Documentation vs Code Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Comparison Summary

| Item | README / Engineering Docs | Implementation / Code | Assessment |
|------|---------------------------|-----------------------|------------|
| Concurrency Strategies | Pessimistic, Optimistic (direct & retry), Atomic single-statement | Implemented in `internal/inventory/service.go` and `store.go` | PASS |
| Directory Structure | Stated structure in `README.md` | Matches exact layout on disk | PASS |
| Execution Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All 3 commands executed cleanly with expected outputs | PASS |
| Error Handling | `ErrNotFound`, `ErrInsufficientStock`, `ErrOptimisticLock`, `ErrInvalidQuantity` | Declared in `model.go` and returned in `store.go` | PASS |
| Demo Output | 5 sections matching strategies | Executed demo produces output exactly matching documentation claims | PASS |

## Discrepancies Found
- None. Documentation accurately reflects the codebase without overclaims or outdated references.
