# Implementation Notes

## Files Added
- `internal/bank/account.go`: Contains `Account` struct with channel-based resource locking to simulate timeouts.
- `internal/transfer/transfer.go`: Contains core logic demonstrating naive locks, lock ordering, and application retry.
- `tests/transfer_test.go`: Suite of concurrency tests proving assertions from the research.
- `cmd/demo/main.go`: Executable printing scenarios end-to-end.

## Core Design Decisions
- **Channel-based locks with context timeout**: Instead of using `sync.Mutex` which blocks indefinitely causing an unrecoverable `fatal error: all goroutines are asleep - deadlock!`, I implemented a custom `RowLock` mechanism utilizing Go channels and `select` with `ctx.Done()`. This accurately simulates a database's "deadlock monitor" which detects wait cycles and aborts one transaction as the "deadlock victim".
- **Artificial Delay**: Inserted `time.Sleep` between lock acquisitions to deterministically trigger race conditions and circular waits in the naive implementation, validating Finding 4 (long transactions increase deadlock frequency).

## Implementation-Specific Choices
- Wait time values (e.g. `20ms`, `50ms`) are arbitrary and specifically tuned for tests to run quickly while consistently reproducing the deadlock behavior.

## Known Limitations
- This is an in-memory simulation at the application level. An actual RDBMS would use a Wait-for Graph (WFG) or Lock Manager to instantaneously detect cycles rather than relying purely on timeouts.
- Does not demonstrate complex multi-row deadlocks across multiple tables, limiting scope to simple A-to-B scenarios.

## Trade-offs
- Used simple boolean `ID` string comparison for Lock Ordering. This implies all entities have universally comparable identifiers.
- Retry logic uses simple loop with basic sleep instead of full exponential backoff with jitter, prioritizing readability for the lab.

## What Is Demonstrated
- Deadlocks occur when resources are requested out of order.
- Lock Ordering fully prevents deadlocks.
- Application-level Retry patterns can recover transactions aborted by deadlock detection.
- Transaction duration correlates with deadlock probability.

## What Is Not Demonstrated
- RDBMS internal deadlock monitors (WFG analysis).
- Database isolation levels.
- Multi-resource (3+ locks) deadlocks.
