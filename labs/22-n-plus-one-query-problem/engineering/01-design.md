# Engineering Design

Target Lab: labs/22-n-plus-one-query-problem
Research Status: APPROVED_WITH_WARNINGS

## Concept To Prove
Demonstrate the N+1 query problem and how eager loading (batching queries) reduces database round-trips from N+1 down to 2 queries.

## Expected Behavior
- Executing an un-batched query to fetch entities and their relationships results in 1 initial query plus N additional queries.
- Executing a batched/eager loading query reduces the total query count to 2 queries regardless of the number of parent records.

## Failure Scenario
The naive approach generates excessive database calls leading to aggregate latency overhead and potential connection pool exhaustion under load.

## Success Criteria
- Automated test verifies N+1 queries executed in naive approach (N=3, Total=4).
- Automated test verifies exactly 2 queries executed in eager loading approach (Total=2).
- Both return identical data models.

## Architecture
- In-memory mock database tracking query executions.
- Repository layer providing both naive (N+1) and eager (batched) relationship retrieval methods.
- Executable demo verifying behavior visually.

## Components
- `Store`: Mock in-memory database with query counter.
- `Repository`: Data access patterns exposing `GetAuthorsWithPostsNPlusOne` and `GetAuthorsWithPostsEager`.
- `Demo`: CLI runner executing and outputting query statistics.

## Test Strategy
- Unit tests verifying the query counter in both scenarios against expected thresholds.

## Execution Plan
1. Initialize repository with sample data.
2. Query parents and lazily query children inside a loop.
3. Query parents and batch-fetch children with a single `IN` query.
4. Assert query counts.

## Implementation Decisions
- Used an in-memory thread-safe mock datastore with explicit query counting rather than an external database engine to keep the lab self-contained and zero-dependency.
