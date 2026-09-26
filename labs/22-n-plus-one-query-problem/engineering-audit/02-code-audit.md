# Engineering Code Audit

Target Lab: `labs/22-n-plus-one-query-problem`

## Finding 1

Location: `internal/blog/store.go:6-80`
Claimed Behavior: Thread-safe in-memory store simulating database query execution and query counting.
Observed Implementation: All store access methods (`GetAllAuthors`, `GetPostsByAuthorID`, `GetPostsByAuthorIDs`, `GetQueryCount`, `ResetQueryCount`) synchronize access and mutation using `s.mu.Lock()` and `s.mu.Unlock()`. Query counter increments properly per method invocation.
Assessment: PASS
Severity: LOW
Notes: Straightforward, robust mutex locking.

## Finding 2

Location: `internal/blog/repository.go:13-28`
Claimed Behavior: Naive relationship loading triggering N+1 queries.
Observed Implementation: Executes `GetAllAuthors()` once, then iterates over slice of authors and executes `GetPostsByAuthorID(author.ID)` for each author. Correctly returns aggregated list of `AuthorWithPosts`.
Assessment: PASS
Severity: LOW
Notes: Accurately simulates the naive N+1 query pattern.

## Finding 3

Location: `internal/blog/repository.go:32-57`
Claimed Behavior: Eager loading (batching) relationship fetching triggering 2 queries.
Observed Implementation: Collects all author IDs into a slice, calls `GetPostsByAuthorIDs(authorIDs)` in one batched call, groups posts into a map in memory (`postsByAuthor`), and maps them back onto the authors list.
Assessment: PASS
Severity: LOW
Notes: Eliminates N round-trips while preserving ordering and author-post associations.

## Finding 4

Location: `internal/blog/repository.go:15-17, 34-36`
Claimed Behavior: Safe handling of empty author datasets.
Observed Implementation: Returns empty slice `[]AuthorWithPosts{}` immediately when `len(authors) == 0`.
Assessment: PASS
Severity: LOW
Notes: Correct guard clause.
