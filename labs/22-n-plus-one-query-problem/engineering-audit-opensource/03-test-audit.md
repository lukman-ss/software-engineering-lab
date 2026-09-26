# Test Audit

File reviewed: `internal/blog/repository_test.go` (3 tests).

## Coverage Matrix

| Scenario                     | Covered? | Test |
|------------------------------|----------|------|
| Happy path (N+1 query count) | YES      | `TestGetAuthorsWithPostsNPlusOne` — asserts 3 authors, 4 queries |
| Happy path (eager query count)| YES     | `TestGetAuthorsWithPostsEager` — asserts 3 authors, 2 queries |
| Data equivalence (eager == N+1)| YES    | `TestGetAuthorsWithPostsEager` — `reflect.DeepEqual` |
| Empty input                  | YES      | `TestEmptyStore` — nil authors/posts |
| Concurrency / race           | NO       | `go test -race` runs but no tests spawn concurrent goroutines against the Store — race-safety is verified by the detector's build instrumentation only, not by a dedicated concurrency test |
| Failure path                 | NO (N/A) | Store never returns errors; in-memory mock has no failure surfaces to inject |
| Negative cases               | PARTIAL  | Empty store covered; missing/duplicate author IDs not covered (acceptable for in-memory mock) |

## Findings

### Finding 1 — Query-count assertions match implementation
`TestGetAuthorsWithPostsNPlusOne` expects 4 queries (1 + 3 authors). `TestGetAuthorsWithPostsEager` expects 2 (1 + 1). Both pass and match `repository.go`. **PASS**.

### Finding 2 — Data-equivalence assertion
`TestGetAuthorsWithPostsEager` deep-compares eager output with the N+1 output. Ordering is preserved because both scan the store in insertion order. Passes. **PASS — strong test**.

### Finding 3 — Empty-store edge case
`TestEmptyStore` constructs a zero-value `Store{authors:nil, posts:nil}` and confirms both methods return empty, equal results without spurious queries. **PASS — good edge case**.

### Finding 4 — Missing concurrency test
Thread-safety of the query counter is *implemented* (mutex on every access) and the race detector reports no data races, but no test actually exercises concurrent calls to `GetAllAuthors`/`GetPostsByAuthorID`/`GetPostsByAuthorIDs`. The mutex therefore remains *unproven under load*, only proven by inspection.
Assessment: WARNING
Severity: LOW
Notes: Add a test that calls the query methods from multiple goroutines (e.g., via `errgroup`/ `sync.WaitGroup`) and asserts the counter equals the number of invocations. Not blocking for this lab.

### Finding 5 — No assertion on post contents per author
Tests assert *counts and shape* (len, DeepEqual vs. N+1) but never assert *specific post titles/IDs per author* (e.g., that Alice has 2 posts). Because eager is cross-checked against N+1 via DeepEqual, a data bug would have to exist identically in both implementations to be missed. For a demo this is acceptable.
Assessment: LOW
Severity: LOW
Notes: A content-specific assertion would strengthen the suite but is not required.

## Overall
Test suite is correct and passes (incl. `-race`). The two main claims (N+1 => 4 queries; eager => 2 queries; equivalent data) are proven. The one soft gap is lack of a dedicated concurrency test; this is non-blocking for the lab.
