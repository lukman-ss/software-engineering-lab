# Test Audit

Target Lab: `labs/22-n-plus-one-query-problem`

## Test Coverage Summary

1. `TestGetAuthorsWithPostsNPlusOne` (`internal/blog/repository_test.go:8-25`):
   - Validates parent count equals 3.
   - Verifies query count equals 4 (1 + 3).
   - Confirms N+1 query explosion behavior.

2. `TestGetAuthorsWithPostsEager` (`internal/blog/repository_test.go:27-50`):
   - Validates parent count equals 3.
   - Verifies query count equals 2 (1 author query + 1 batch post query).
   - Validates deep structural equality (`reflect.DeepEqual`) between eager result and N+1 result to prove data equivalence.

3. `TestEmptyStore` (`internal/blog/repository_test.go:52-65`):
   - Tests nil/empty datasets for both methods.
   - Asserts non-nil empty slice equivalence.

## Execution Verification

### Command: `go test -v ./...`
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
=== RUN   TestEmptyStore
--- PASS: TestEmptyStore (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	(cached)
```

### Command: `go test -race -v ./...`
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
=== RUN   TestEmptyStore
--- PASS: TestEmptyStore (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	1.321s
```

### Assessment
- Happy path covered: PASS
- Edge cases covered: PASS
- Query count assertions: PASS
- Concurrency & Race detector: PASS
- Equivalence assertions: PASS
