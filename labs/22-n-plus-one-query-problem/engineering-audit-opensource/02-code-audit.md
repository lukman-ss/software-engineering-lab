## Finding 1: Query Count Accuracy

Location: 
- internal/blog/store.go:45 (GetAllAuthors), 52 (GetPostsByAuthorID), 66 (GetPostsByAuthorIDs)
- internal/blog/repository.go:14,21,33,43

Claimed Behavior: 
- N+1 approach executes 1 query for authors + N queries for posts.
- Eager approach executes 1 query for authors + 1 batched query for posts.

Observed Implementation:
- GetAllAuthors increments queryCount by 1.
- GetPostsByAuthorID increments queryCount by 1 per call.
- GetPostsByAuthorIDs increments queryCount by 1 per call (regardless of number of IDs).

Assessment: PASS
Severity: N/A
Notes: 
Query counting accurately reflects number of store method invocations. 
N+1 test: 3 authors -> 1 (GetAllAuthors) + 3*GetPostsByAuthorID = 4 queries.
Eager test: 1 (GetAllAuthors) + 1 (GetPostsByAuthorIDs) = 2 queries.
Matches expectations.

## Finding 2: Concurrency Safety

Location: internal/blog/store.go (lines 7,30-34,36-40,42-47,49-52,52-62,63-80)

Claimed Behavior: Store is thread-safe via mutex.

Observed Implementation:
- All exported methods lock the mutex (mu) before accessing shared state.
- No missed locks on read/write paths.
- GetQueryCount and ResetQueryCount also lock.

Assessment: PASS
Severity: N/A
Notes: 
Mutex protects authors, posts, and queryCount.
No obvious race conditions in store methods.

## Finding 3: Error Handling Absence

Location: internal/blog/store.go:42-61

Claimed Behavior: None claimed (in-memory mock).

Observed Implementation:
- GetAllAuthors returns slice; no error possible.
- GetPostsByAuthorID/ByAuthorIDs return slice; no error possible.
- No error propagation.

Assessment: WARNING
Severity: LOW
Notes: 
In a real database, these methods could return errors. 
The mock does not simulate error conditions, but the lab is about query counting, not error handling. 
Acceptable for scope.

## Finding 4: State Cleanup/Reset

Location: internal/blog/store.go:36-40

Claimed Behavior: ResetQueryCount resets queryCount to 0.

Observed Implementation:
- Locks mutex, sets queryCount = 0, unlocks.

Assessment: PASS
Severity: N/A
Notes: 
Used correctly in tests and demo.

## Finding 5: Unnecessary Complexity

Location: internal/blog/repository.go:30-58 (GetAuthorsWithPostsEager)

Claimed Behavior: Eager loading via batching.

Observed Implementation:
- Collects authorIDs, makes one batched query, builds postsByAuthor map, then reconstructs result.
- Straightforward and clear.

Assessment: PASS
Severity: N/A
Notes: 
No unnecessary complexity; implementation is minimal and correct.

## Finding 6: Test Coverage for Edge Cases

Location: internal/blog/repository_test.go

Claimed Behavior: Tests cover happy path and empty store.

Observed Implementation:
- Tests for N+1 and eager loading with 3 authors.
- TestEmptyStore covers zero authors/posts.
- No test for single author, or mismatch between authorIDs and posts.

Assessment: WARNING
Severity: MEDIUM
Notes: 
Edge cases like single author (N=1 -> N+1=2, eager=2) are not explicitly tested but would pass.
No test for authors with zero posts; however, mock data ensures each author has posts.
Still, adding a test for author with zero posts would improve coverage.

## Finding 7: Consistency Between N+1 and Eager Results

Location: internal/blog/repository_test.go:45-49

Claimed Behavior: Eager result should deep equal N+1 result.

Observed Implementation:
- TestGetAuthorsWithPostsEager calls N+1 result and compares via reflect.DeepEqual.

Assessment: PASS
Severity: N/A
Notes: 
Test passes, confirming functional equivalence.

## Finding 8: Demo Output Matches Implementation

Location: cmd/demo/main.go

Claimed Behavior: Demo prints query counts for N+1 and eager.

Observed Implementation:
- Prints exactly as expected based on store query counting.

Assessment: PASS
Severity: N/A
Notes: 
Demo output matches test expectations.

## Finding 9: No Timeout or Recovery Behavior

Location: N/A

Claimed Behavior: None claimed.

Observed Implementation:
- No network or external dependencies; no timeouts or recovery needed.

Assessment: PASS
Severity: N/A
Notes: 
Not applicable to in-memory mock.

## Finding 10: Potential Issue: GetPostsByAuthorIDs Parameter Validation

Location: internal/blog/store.go:63-80

Claimed Behavior: Accepts slice of authorIDs.

Observed Implementation:
- No validation for nil or empty slice; if authorIDs is nil, panic on range.
- If authorIDs empty, returns empty posts (correct).

Assessment: WARNING
Severity: LOW
Notes: 
Nil slice would cause panic. However, callers (repository) always pass non-empty slice because they check len(authors)==0 early and return empty result. 
Still, defensive coding would guard against nil.