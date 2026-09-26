# Implementation Notes

## Files Added
- `internal/blog/models.go`: Domain models (`Author`, `Post`, `AuthorWithPosts`).
- `internal/blog/store.go`: Thread-safe mock data store with internal query counter.
- `internal/blog/repository.go`: Data access layer with N+1 and Eager loading implementations.
- `internal/blog/repository_test.go`: Unit tests proving query count reduction.
- `cmd/demo/main.go`: CLI application.

## Core Design Decisions
- Simulated database using in-memory structs and a `queryCount` integer incremented upon each data retrieval method invocation.

## Implementation-Specific Choices
- Used a generic loop-based resolution for eager loading to mimic what an ORM or DataLoader natively does under the hood.

## Known Limitations
- Does not demonstrate network-level latency costs, as the data store is purely in-memory.
- Memory bloat issue from unrestrained eager loading (as mentioned in the research) is not explicitly proven via OOM errors here, as the dataset is deliberately small.

## Trade-offs
- Sacrificed real I/O benchmarking for deterministic query counting.

## What Is Demonstrated
- The exact explosion of queries (N+1) when fetching relationships iteratively.
- The mitigation of the query explosion (reducing N queries to 1 batched query) using eager loading.

## What Is Not Demonstrated
- Memory bloat (Out of Memory) conditions caused by over-fetching.
- ORM specific mapping anomalies (e.g. JPA `FetchType.EAGER` mapping level issues).
