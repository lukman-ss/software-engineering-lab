# Engineering Audit Plan

Target Lab: `labs/22-n-plus-one-query-problem`
Implementation Files:
- `internal/blog/models.go`
- `internal/blog/store.go`
- `internal/blog/repository.go`
Tests:
- `internal/blog/repository_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md` (APPROVED)
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
Main Claims To Verify:
1. Iterative/naive relationship fetching produces N+1 total queries (1 parent query + N child queries).
2. Eager loading (batch fetching with `IN`) resolves all relationships in exactly 2 queries (1 parent + 1 batched child).
3. Both methods return functionally identical datasets.
4. Concurrency safety of store query tracking (`sync.Mutex`).
5. Zero dependency build and test execution without race conditions.
Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -count=1 ./...`
- `go run ./cmd/demo`
Primary Risks:
- Thread-safety / race conditions during concurrent access or query counting.
- Data divergence between naive and eager loading results.
- Mismatch between README instructions and codebase structure/behavior.
