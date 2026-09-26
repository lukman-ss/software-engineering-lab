# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
(Success - no output)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
=== RUN   TestGetAuthorsWithPostsNPlusOne
--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)
=== RUN   TestGetAuthorsWithPostsEager
--- PASS: TestGetAuthorsWithPostsEager (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	0.576s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
?   	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/cmd/demo	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/22-n-plus-one-query-problem/internal/blog	1.398s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
--- 1. Simulating N+1 Query Problem ---
Loaded 3 authors with their posts.
Total queries executed: 4 (1 query for authors + 3 queries for posts)

--- 2. Simulating Eager Loading (Batching) ---
Loaded 3 authors with their posts.
Total queries executed: 2 (1 query for authors + 1 batched query for posts)
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
