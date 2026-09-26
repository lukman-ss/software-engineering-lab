# Code Audit

## Finding 1

Location: `internal/blog/repository.go:13-28` (`GetAuthorsWithPostsNPlusOne`)
Claimed Behavior: Demonstrates N+1 query problem by making 1 query for all authors and N queries for posts.
Observed Implementation: Fetches all authors using `r.store.GetAllAuthors()` (1 query), then iterates over each author calling `r.store.GetPostsByAuthorID(author.ID)` (N queries).
Assessment: PASS
Severity: LOW
Notes: Accurately models naive lazy loading query explosion.

## Finding 2

Location: `internal/blog/repository.go:32-57` (`GetAuthorsWithPostsEager`)
Claimed Behavior: Demonstrates eager loading mitigation by reducing query count to 2.
Observed Implementation: Fetches all authors (1 query), extracts IDs into a slice, fetches posts with `r.store.GetPostsByAuthorIDs(authorIDs)` (1 query), maps posts in memory, and constructs `AuthorWithPosts`.
Assessment: PASS
Severity: LOW
Notes: Idiomatic batching implementation cleanly matching design.

## Finding 3

Location: `internal/blog/store.go:30-80` (`Store` mutex usage)
Claimed Behavior: Thread-safe data store with query counting.
Observed Implementation: All accessors (`GetQueryCount`, `ResetQueryCount`, `GetAllAuthors`, `GetPostsByAuthorID`, `GetPostsByAuthorIDs`) lock `s.mu.Lock()` and defer unlock before mutating or reading `queryCount` and underlying slice data.
Assessment: PASS
Severity: LOW
Notes: No race conditions detected.

## Finding 4

Location: `internal/blog/repository.go:45-48` (`GetAuthorsWithPostsEager`)
Claimed Behavior: Correctly maps posts by author ID.
Observed Implementation: Maps all returned posts into `postsByAuthor[post.AuthorID]` slice map. Authors without posts correctly map to a `nil` slice which yields an empty slice behavior in JSON/Go domain models without nil panics.
Assessment: PASS
Severity: LOW
Notes: Fully handles empty and multi-item author post slices correctly.
