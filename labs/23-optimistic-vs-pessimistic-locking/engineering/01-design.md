# Engineering Design

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Research Status: APPROVED

## Concept To Prove
Prove how concurrent transactions executing naive read-modify-write patterns result in silent lost updates under default isolation levels, and demonstrate the three approved remediation patterns:
1. Pessimistic Locking (`SELECT ... FOR UPDATE` row exclusivity)
2. Optimistic Locking (version guard in WHERE clause with retry / conflict detection)
3. Atomic Single-Statement Operations (`UPDATE ... SET stock = stock - N WHERE stock >= N`)

## Expected Behavior
- Naive read-modify-write demonstrates lost updates under concurrent goroutines.
- Pessimistic locking prevents concurrency conflicts entirely by blocking writers on a per-row lock.
- Optimistic locking detects version mismatch and safely rejects stale writes without data corruption; application retry loop eventually converges successfully.
- Atomic single-statement updates serialize safely at the database statement level without manual transaction lock blocks.

## Failure Scenario
- Concurrent workers overwrite each other's changes without noticing, leading to stock discrepancies and inventory corruption (e.g. 50 successful deduct calls yield final stock > 50).

## Success Criteria
- Automated test demonstrates lost update under unsynchronized concurrent read-modify-write.
- Automated test validates pessimistic locking maintains exact inventory invariants.
- Automated test validates optimistic conflict detection and retry convergence with backoff.
- Automated test validates atomic conditional decrement guarantees consistency.
- Race detector passes with zero race warnings.

## Architecture
- In-memory thread-safe datastore (`Store`) simulating SQL storage engine semantics (row-level locks, versioned rows, atomic updates).
- Service layer (`Service`) implementing naive, pessimistic, optimistic, and atomic decrement logic.
- Automated test suite validating behavior and invariants under concurrent load.
- CLI demo running side-by-side scenarios.

## Components
- `internal/inventory/model.go`: Domain models and error definitions.
- `internal/inventory/store.go`: Simulated engine supporting row-level mutexes, version incrementing, and atomic updates.
- `internal/inventory/service.go`: Business operations for deduction under different concurrency strategies.
- `cmd/demo/main.go`: Visual CLI runner showing all 4 scenarios.
- `tests/locking_test.go`: Concurrency and invariant tests.

## Test Strategy
- Unit & concurrency tests with 20-50 parallel goroutines per scenario.
- Invariant assertions: `InitialStock - Deductions == FinalStock` (or accounting for rejected conflicts).
- Go race detector validation (`go test -race ./...`).

## Execution Plan
1. Run `go test ./...`
2. Run `go test -race ./...`
3. Execute `go run ./cmd/demo`
4. Inspect outputs and record results.

## Implementation Decisions
- Implementation-specific decision: In-memory store uses granular per-row mutexes (`sync.Mutex`) to model `SELECT ... FOR UPDATE` behavior rather than external PostgreSQL / MySQL instance, keeping the lab dependency-free, cross-platform, and fast.
- Jittered exponential backoff implemented in optimistic retry loop to prevent thundering herd during lock contention.
