# Code Audit

Target Lab: labs/22-n-plus-one-query-problem

## Finding 1

Location: internal/blog/store.go:30-80
Claimed Behavior: Thread-safe in-memory database simulation that accurately tracks query counts.
Observed Implementation: Store uses `sync.Mutex` across all data retrieval and query counter methods (`GetQueryCount`, `ResetQueryCount`, `GetAllAuthors`, `GetPostsByAuthorID`, `GetPostsByAuthorIDs`). Each query call increments `queryCount` by 1 under the lock.
Assessment: PASS
Severity: LOW
Notes: Concurrency-safe for concurrent readers/accessors. Returns internal slices directly; acceptable for deterministic mock store.

## Finding 2

Location: internal/blog/repository.go:13-25
Claimed Behavior: Naive relationship retrieval executes 1 query for parents and N queries for child records.
Observed Implementation: Calls `r.store.GetAllAuthors()`, iterates over the resulting slice, and calls `r.store.GetPostsByAuthorID(author.ID)` once per author.
Assessment: PASS
Severity: LOW
Notes: Perfectly models the iterative N+1 query antipattern.

## Finding 3

Location: internal/blog/repository.go:29-55
Claimed Behavior: Eager loading batching reduces query execution to 2 queries total.
Observed Implementation: Calls `r.store.GetAllAuthors()`, extracts author IDs, calls `r.store.GetPostsByAuthorIDs(authorIDs)`, and groups records in memory using a map before populating the result structs.
Assessment: PASS
Severity: LOW
Notes: Correctly reproduces the two-query eager loading pattern common in ORMs and DataLoaders.

## Finding 4

Location: internal/blog/repository.go:31-33 vs internal/blog/repository.go:16
Claimed Behavior: Symmetrical data retrieval return semantics between naive and eager methods.
Observed Implementation: When authors slice is empty, `GetAuthorsWithPostsEager` early-exits returning `nil`, whereas `GetAuthorsWithPostsNPlusOne` initializes `var result []AuthorWithPosts` and returns an empty non-nil slice (`[]AuthorWithPosts{}`).
Assessment: WARNING
Severity: LOW
Notes: Functional behavior is equivalent in Go slice handling, but creates minor return value inconsistency when empty.

## Finding 5

Location: internal/blog/repository.go
Claimed Behavior: Production N+1 performance bottleneck simulation.
Observed Implementation: The implementation is an in-memory simulation with zero I/O and zero latency overhead.
Assessment: PASS
Severity: LOW
Notes: The engineering notes explicitly declare this trade-off: in-memory deterministic query counting was chosen over real database engine setup to maintain zero external dependencies.
