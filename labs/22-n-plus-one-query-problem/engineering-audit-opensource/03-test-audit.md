# Test Audit

Test File: internal/blog/repository_test.go

## Coverage Matrix

| Scenario | Covered | Test Name | Evidence |
|---|---|---|---|
| Happy path (N+1 query count) | Yes | TestGetAuthorsWithPostsNPlusOne | Asserts queryCount == 4 |
| Happy path (eager query count) | Yes | TestGetAuthorsWithPostsEager | Asserts queryCount == 2 |
| Result equivalence (eager == N+1) | Yes | TestGetAuthorsWithPostsEager | reflect.DeepEqual comparison |
| Empty store (no authors/posts) | Yes | TestEmptyStore | Asserts len == 0 for both |
| Author with zero posts | No | — | Mock data ensures every author has posts; not tested |
| Single author (N=1) | No | — | Only 3-author case tested |
| Failure path (store returns error) | No | — | Store cannot fail in mock |
| Concurrency (parallel calls) | No | — | No concurrent access tests; -race passes because no concurrency used |
| Negative cases (post belongs to nonexistent author) | No | — | Mock data consistent |

## Test Quality Assessment

### TestGetAuthorsWithPostsNPlusOne (PASS)
- Creates store, repo, resets count, invokes N+1.
- Asserts len(result) == 3 (correct: 3 authors).
- Asserts queryCount == 4 (correct: 1 + N = 1 + 3 = 4).
- Comment is correct: "3 authors + 1 initial query = 4".

### TestGetAuthorsWithPostsEager (PASS)
- Creates store, repo, resets count, invokes eager loading.
- Asserts len(result) == 3 (correct).
- Asserts queryCount == 2 (correct: 1 + 1 = 2).
- Deep-compares result with N+1 output via reflect.DeepEqual.
  - This is good: it proves functional equivalence, not just count.

### TestEmptyStore (PASS)
- Creates an empty store & repo.
- Asserts both methods return empty results.
- Confirms N+1 and eager are equal for empty input.

## Strengths
1. Query counts are asserted precisely (not just "improved").
2. Functional equivalence between naive and optimized approaches is asserted.
3. Empty input edge case is covered.
4. All tests pass consistently.

## Weaknesses / Missing Coverage
1. **Author with zero posts**: Not tested. N+1 would issue a query even for authors with no posts; eager would batch (still 1 query). Test should confirm this.
2. **Single author (N=1)**: N+1 count should be 2; eager should be 2. Only N=3 case tested.
3. **Concurrency safety**: Although store is mutex-guarded, tests do not call methods concurrently to validate this assumption. Race detector passes vacuously because no goroutines are spawned.
4. **No negative case**: Posts with an AuthorID not belonging to any author are not tested (though mock data prevents this).

## Test Execution Results (Actual)

```
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
=== RUN   TestEmptyStore
--- PASS: TestEmptyStore (0.00s)
PASS
ok  internal/blog  (cached)
```

## Race Detector Result

```
ok  internal/blog (cached)
```

No data races detected (tests do not use goroutines; race detector does not trigger).

## Overall Assessment: MEDIUM

Tests pass and cover core claims. Query-count assertions are accurate. However, edge-case coverage is limited (author with no posts, single-author, concurrency stress). Race detector runs clean but is not meaningfully exercised.

These gaps are acceptable for a demonstration lab but do not constitute "proven" concurrency safety under parallel load.
