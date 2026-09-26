# Engineering Audit Plan

Target Lab: labs/22-n-plus-one-query-problem
Implementation Files:
- internal/blog/models.go
- internal/blog/store.go
- internal/blog/repository.go
- cmd/demo/main.go
Tests:
- internal/blog/repository_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/ (not audited in this stage)
Main Claims To Verify:
1. The N+1 problem is demonstrated correctly: 1 query for authors + N queries for posts.
2. The eager loading solution reduces queries to 2: 1 query for authors + 1 batched query for posts.
3. The results of both methods are equivalent (same data returned).
4. The store correctly tracks query counts with thread safety (using mutex).
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- The eager loading implementation might not actually batch correctly (if GetPostsByAuthorIDs is not used or if it still does per-author queries).
- The test might not cover edge cases (zero authors, duplicate author IDs, etc.).
- Concurrency safety: ensure mutex is used correctly in all store methods.