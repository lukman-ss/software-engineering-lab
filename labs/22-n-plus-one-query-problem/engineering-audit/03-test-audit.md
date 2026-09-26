# Test Audit

Target Lab: labs/22-n-plus-one-query-problem

## Test Execution Summary

Command:
```bash
go test -v ./...
```
Output:
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	0.576s
```

Race Detector Command:
```bash
go test -race -v ./...
```
Output:
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	1.171s
```

## Coverage Assessment

1. Happy Path Coverage:
   - PASS: Verifies N+1 query count (`1 + 3 = 4`) for naive implementation.
   - PASS: Verifies batch query count (`1 + 1 = 2`) for eager implementation.

2. Content Equivalence Verification:
   - WARNING: Tests assert `len(result) == 3`, but do not assert that the post contents/IDs match between `GetAuthorsWithPostsNPlusOne` and `GetAuthorsWithPostsEager`.

3. Edge Cases:
   - WARNING: No test cases verify behavior when the database has 0 authors.
   - WARNING: No test cases verify behavior when an author has 0 posts.

4. Concurrency & Race Safety:
   - PASS: `go test -race ./...` passed with zero warnings. Mutex implementation in `Store` is sound.

## Test Suite Strengths
- Fast, zero external dependencies, deterministic verification of query counting formulas.

## Test Suite Weaknesses
- Asserts only query counts and author lengths; does not deeply assert relational integrity of returned posts.
