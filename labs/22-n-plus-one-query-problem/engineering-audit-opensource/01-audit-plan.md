# Engineering Audit Plan

Target Lab: `labs/22-n-plus-one-query-problem`
Implementation Files:
  - `internal/blog/models.go` (domain types: Author, Post, AuthorWithPosts)
  - `internal/blog/store.go` (in-memory mock DB + thread-safe query counter)
  - `internal/blog/repository.go` (naive N+1 fetch + eager-loading fetch)
Tests:
  - `internal/blog/repository_test.go` (query-count + data-equivalence tests)
Executable/Demo:
  - `cmd/demo/main.go`
Approved Research Inputs: `research/*` (out of scope for this implementation audit)
Main Claims To Verify:
  1. Naive author+posts fetch exhibits N+1: 1 query for authors + N queries for posts (= 4 for 3 authors).
  2. Eager-loading (batched) fetch reduces total queries to 2 (1 authors + 1 posts).
  3. Both approaches return identical result data.
  4. Store query counting is concurrency-safe (mutex guarded).
Commands To Run:
  - `go test -v ./...`
  - `go test -race ./...`
  - `go run ./cmd/demo`
Primary Risks:
  - Query counter not actually exercised / off-by-one in counting.
  - Eager result data divergence from N+1 result (ordering / missing posts).
  - Race condition in shared mutex-guarded query counter under concurrency.
  - Demo output fabricated vs. tests.
