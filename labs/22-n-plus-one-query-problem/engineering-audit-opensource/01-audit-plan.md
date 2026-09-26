# Engineering Audit Plan

Target Lab: labs/22-n-plus-one-query-problem
Audit Output Directory: labs/22-n-plus-one-query-problem/engineering-audit-opensource/

## Implementation Files
- `internal/blog/models.go`: Domain models (`Author`, `Post`, `AuthorWithPosts`).
- `internal/blog/store.go`: In-memory mock data store with mutex-protected query counting.
- `internal/blog/repository.go`: Repository with `GetAuthorsWithPostsNPlusOne` and `GetAuthorsWithPostsEager`.

## Tests
- `internal/blog/repository_test.go`: Unit tests asserting query counts and data equivalence.

## Executable/Demo
- `cmd/demo/main.go`: CLI runner printing query counts for naive vs. eager loading.

## Approved Research Inputs
- Research directory (`research/`) exists with a report. Engineering design (`engineering/01-design.md`) marks "Research Status: APPROVED".
> Note: Per pipeline override, this audit covers implementation and tests only; research content is not audited in this stage.

## Main Claims To Verify
1. Naive fetch performs N+1 queries (expected total = 4 for N=3 authors).
2. Eager/batched fetch performs 2 queries total regardless of N.
3. Naive and eager results return identical data (deep equality).
4. Empty store returns empty results without panics.
5. Query counting is thread-safe (Store protected by mutex).

## Commands To Run
- `go build ./...`
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
- `go vet ./...`

## Primary Risks
- The in-memory mock may not reflect real DB I/O (by-design limitation; documented).
- Query counting is the sole correctness proxy; tests must assert exact counts.
- Concurrency safety must be verified since Store is accessed via methods guarded by a mutex.
