# Code Audit

Scope: `internal/blog/` (`models.go`, `store.go`, `repository.go`) and `cmd/demo/main.go`.

## Finding 1
Location: `internal/blog/store.go:42-47` — `GetAllAuthors`
Claimed Behavior: Returns a snapshot of authors and increments the query counter (1 query).
Observed Implementation: Holds the mutex, increments `queryCount`, then returns the internal slice header directly (no copy).
Assessment: PASS for counting. The returned slice is a reference to internal state; mutation by a caller would break encapsulation, but no current caller mutates it.
Severity: LOW
Notes: Returning an internal slice is a minor encapsulation leak; acceptable for an in-memory demo but should return a copy if the type is ever reused in concurrent mutation scenarios.

## Finding 2
Location: `internal/blog/store.go:49-80` — `GetPostsByAuthorID`, `GetPostsByAuthorIDs`
Claimed Behavior: Each call = exactly 1 query, guarded by mutex.
Observed Implementation: Each method locks `s.mu`, increments `queryCount` by one, filters posts, and unlocks. Counting is correct.
Assessment: PASS
Severity: N/A
Notes: `GetPostsByAuthorIDs` builds an `idMap` (O(M)) then scans all posts (O(N)); fine for a mock store.

## Finding 3
Location: `internal/blog/store.go:30-40` — query-count accessors
Claimed Behavior: Thread-safe query counting.
Observed Implementation: `GetQueryCount` and `ResetQueryCount` both acquire `s.mu`; the three data-query methods also acquire `s.mu` before mutating `queryCount`.
Assessment: PASS (concurrency-safe by construction)
Severity: N/A
Notes: The mutex guards every read/write of `queryCount`. Race detector passes.

## Finding 4
Location: `internal/blog/repository.go:13-28` — `GetAuthorsWithPostsNPlusOne`
Claimed Behavior: 1 query for authors + N queries for posts (N = number of authors).
Observed Implementation: Calls `GetAllAuthors()` once, then calls `GetPostsByAuthorID(author.ID)` once per author. For the seeded store (3 authors) => 4 queries total.
Assessment: PASS
Severity: N/A
Notes: Empty-input fast path returns `[]AuthorWithPosts{}` without touching the store when `len(authors)==0`.

## Finding 5
Location: `internal/blog/repository.go:32-58` — `GetAuthorsWithPostsEager`
Claimed Behavior: Eager loading via batching => 1 authors query + 1 posts query = 2 queries.
Observed Implementation: Calls `GetAllAuthors()` once and `GetPostsByAuthorIDs(authorIDs)` once, then partitions posts via a `map[int][]Post`. Total 2 queries for the seeded store.
Assessment: PASS
Severity: N/A
Notes: Posts per author preserve store order (scan is in insertion order), matching the N+1 result ordering so a `DeepEqual` between the two approaches is valid.

## Finding 6
Location: `internal/blog/repository.go:15-17` / `:34-36` — empty-store handling
Claimed Behavior: Handles zero authors gracefully.
Observed Implementation: Both functions return `[]AuthorWithPosts{}` when `len(authors)==0`.
Assessment: PASS
Severity: N/A
Notes: No query performed in the empty path, so no spurious counter increments.

## Finding 7
Location: `internal/blog/models.go`
Claimed Behavior: Domain models for Author, Post, and a joined `AuthorWithPosts`.
Observed Implementation: Plain structs with value fields. `AuthorWithPosts.Posts` is a `[]Post` slice (nil-safe: nil maps and nil slices compare equal under `reflect.DeepEqual`).
Assessment: PASS
Severity: N/A
Notes: For authors that have posts, both paths yield a populated `[]Post`; for authors with no posts both yield `nil`. Consistent.

## Findings not applicable (simple in-memory, no I/O)
- Timeout behavior: N/A — no network/disk IO, no timeouts.
- Recovery/rollback: N/A — no mutations; reads only.
- Cleanup: N/A — no resources to release; no goroutines spawned.
- Error propagation: N/A — Store methods do not return errors (in-memory mock); no error values to propagate.

## Overall
No correctness defects. Counting, batching logic, empty-state handling, and mutex-based concurrency guarding are all implemented correctly. The only low-severity note is the encapsulation leak of returning the internal `authors` slice, which does not affect any current behavior.
