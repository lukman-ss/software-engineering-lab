# Code Audit

## Finding 1

Location: `internal/blog/store.go:30-80`
Claimed Behavior: Thread-safe data store tracking query executions accurately.
Observed Implementation: All accessor and mutator methods (`GetQueryCount`, `ResetQueryCount`, `GetAllAuthors`, `GetPostsByAuthorID`, `GetPostsByAuthorIDs`) acquire `s.mu.Lock()` with deferred unlock.
Assessment: PASS
Severity: LOW
Notes: Mutex protection is properly structured across all methods.

## Finding 2

Location: `internal/blog/repository.go:13-28`
Claimed Behavior: Naive relationship loading produces N+1 query pattern.
Observed Implementation: `GetAllAuthors()` fetches authors (1 query), then iterates through authors and calls `GetPostsByAuthorID(author.ID)` in each iteration (N queries). Total = N + 1.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates N+1 query execution behavior.

## Finding 3

Location: `internal/blog/repository.go:32-58`
Claimed Behavior: Eager loading batches child queries into 1 query, reducing total queries to 2.
Observed Implementation: Collects all author IDs into `authorIDs` slice, invokes `GetPostsByAuthorIDs(authorIDs)` once (1 query), maps posts by `AuthorID`, and constructs result slice preserving author ordering. Total = 2 queries.
Assessment: PASS
Severity: LOW
Notes: Correctly handles batching and in-memory aggregation. Handles empty author list by returning empty slice without calling post query.

## Finding 4

Location: `internal/blog/store.go:63-79`
Claimed Behavior: Mock batch query returns all matching posts for given author IDs.
Observed Implementation: Uses hash map lookup `idMap[p.AuthorID]` over internal posts.
Assessment: PASS
Severity: LOW
Notes: Clean implementation with linear complexity.
