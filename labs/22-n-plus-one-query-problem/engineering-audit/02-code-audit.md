# Code Audit

Target Lab: `labs/22-n-plus-one-query-problem`

## Finding 1: Thread-Safe Query Tracking Store

Location: `internal/blog/store.go:6-80`
Claimed Behavior: In-memory store safely tracks query count across reads and resets.
Observed Implementation: `Store` uses `sync.Mutex` protecting all operations (`queryCount`, `authors`, `posts`). All methods lock and defer unlock properly.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safe design.

## Finding 2: Exact N+1 Query Execution Path

Location: `internal/blog/repository.go:13-28`
Claimed Behavior: `GetAuthorsWithPostsNPlusOne` executes 1 initial query for authors, then iterates over each author calling `GetPostsByAuthorID`, yielding N queries for posts (Total = 1 + N).
Observed Implementation: Fetches all authors (1 query), checks length, iterates through N authors calling `GetPostsByAuthorID` (N queries). Total = 1 + N.
Assessment: PASS
Severity: LOW
Notes: Accurately replicates ORM lazy loading traversal.

## Finding 3: Eager Loading / Batching Execution Path

Location: `internal/blog/repository.go:32-57`
Claimed Behavior: `GetAuthorsWithPostsEager` fetches authors (1 query) and all associated posts in a single batched query using IDs (1 query), resulting in 2 queries total, followed by in-memory stitching.
Observed Implementation: Fetches all authors (1 query), builds `authorIDs` slice, fetches all matching posts via `GetPostsByAuthorIDs` (1 query), constructs map `postsByAuthor`, and returns stitched `AuthorWithPosts` slice preserving author ordering.
Assessment: PASS
Severity: LOW
Notes: Accurately implements `IN (...)` batching pattern.

## Finding 4: Empty Result Edge Case

Location: `internal/blog/repository.go:15-17, 34-36`
Claimed Behavior: When no authors exist, returns empty slice immediately without issuing secondary queries.
Observed Implementation: Returns `[]AuthorWithPosts{}` immediately on `len(authors) == 0`.
Assessment: PASS
Severity: LOW
Notes: Correctly handles empty datasets with 1 total query.
