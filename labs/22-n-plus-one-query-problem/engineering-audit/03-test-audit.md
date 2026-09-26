# Test Audit

## Coverage Analysis

1. `TestGetAuthorsWithPostsNPlusOne`:
   - Verifies author count is 3.
   - Asserts query count is exactly 4 (1 + 3).
   - Verifies N+1 query problem behavior.

2. `TestGetAuthorsWithPostsEager`:
   - Verifies author count is 3.
   - Asserts query count is exactly 2 (1 + 1).
   - Validates eager loading mitigation behavior.
   - Uses `reflect.DeepEqual` to verify that eager loading outputs identical data structure to N+1 loading.

3. `TestEmptyStore`:
   - Edge case with empty authors and posts.
   - Confirms no panics and empty result slices returned from both functions.
   - Deep equality check between both results on empty inputs.

## Execution Output

### `go test -v ./...`
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
=== RUN   TestEmptyStore
--- PASS: TestEmptyStore (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	0.108s
```

### `go test -race ./...`
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	1.131s
```

### `go run ./cmd/demo`
```text
--- 1. Simulating N+1 Query Problem ---
Loaded 3 authors with their posts.
Total queries executed: 4 (1 query for authors + 3 queries for posts)

--- 2. Simulating Eager Loading (Batching) ---
Loaded 3 authors with their posts.
Total queries executed: 2 (1 query for authors + 1 batched query for posts)
```

Assessment: PASS
All claimed properties are covered and verified by automated tests and runtime demo.
