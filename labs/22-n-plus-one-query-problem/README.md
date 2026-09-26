# Lab 22: N+1 Query Problem

This lab demonstrates the N+1 query problem and its mitigation using eager loading (batching).

## Structure

- `cmd/demo/main.go`: Executable demonstrating query execution counts for naive vs. eager loading approaches.
- `internal/blog/models.go`: Domain models for Authors and Posts.
- `internal/blog/store.go`: In-memory mock data store with thread-safe query counting.
- `internal/blog/repository.go`: Repository providing naive and eager relationship fetching.
- `internal/blog/repository_test.go`: Tests validating query counts.

## Running the Demo

```bash
go run ./cmd/demo
```

## Running the Tests

```bash
go test -v ./...
go test -race ./...
```
