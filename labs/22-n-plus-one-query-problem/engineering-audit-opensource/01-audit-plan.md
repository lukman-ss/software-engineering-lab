# Engineering Audit Plan

Target Lab: labs/22-n-plus-one-query-problem (module: github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem)

Implementation Files:
- internal/blog/models.go        — Domain models: Author, Post, AuthorWithPosts
- internal/blog/store.go          — In-memory mock data store with thread-safe query counting
- internal/blog/repository.go     — Repository methods: GetAuthorsWithPostsNPlusOne (N+1), GetAuthorsWithPostsEager (batched)
- cmd/demo/main.go                — Executable demonstrating query execution counts for naive vs. eager loading approaches

Tests:
- internal/blog/repository_test.go
  - TestGetAuthorsWithPostsNPlusOne
  - TestGetAuthorsWithPostsEager
  - TestEmptyStore

Executable/Demo:
- cmd/demo/main.go (run via `go run ./cmd/demo`)

Approved Research Inputs:
- Audit implementation and tests only. Research/content are out of scope for this stage per pipeline override.

Main Claims To Verify:
1. The N+1 implementation executes N+1 queries (1 + N).
2. The eager-loading implementation executes exactly 2 queries (1 + 1) regardless of author/post counts.
3. Both implementations return equivalent author/post result sets.
4. The demo prints the correct query counts for both approaches.
5. The store query-counting is thread-safe.

Commands To Run:
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Query counting is simulated (in-memory), not a real DB; must verify counts reflect actual method-call semantics.
- N+1 count assumption (N+1 = 4 for 3 authors) must align with code.
- No negative/edge-case coverage for query-count behavior under eager loading.
