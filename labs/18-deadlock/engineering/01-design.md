# Engineering Design

Target Lab: labs/18-deadlock
Research Status: APPROVED

## Concept To Prove
1. Deadlock as circular wait condition.
2. Deadlock handling via abort (Deadlock Victim).
3. Deadlock prevention using Lock Ordering.
4. Transaction duration impacting deadlock probability.
5. Application-level retry mechanism to recover from deadlock aborts.

## Expected Behavior
- Naive transfers running concurrently in opposite directions (A->B and B->A) deadlock and time out (abort as victim).
- Ordered transfers running concurrently avoid deadlocks completely.
- Retry-wrapped transfers eventually succeed even under contention and deadlocks.
- Longer transactions reliably increase deadlock frequency.

## Failure Scenario
- Concurrent bidirectional naive transfers trigger deadlocks, resulting in `ErrDeadlock` (simulating database victim selection).

## Success Criteria
- Tests verify deadlocks occur under naive bidirectional transfer.
- Tests verify ordered transfers never deadlock.
- Tests verify retry mechanism handles deadlocks and completes successfully.
- Tests verify transaction delay increases deadlock frequency.
- Race detector passes without data races.
- Demo runs end-to-end and prints readable metrics.

## Architecture
- `internal/bank`: Models accounts with lockable resources (simulating row-level database locks with timeouts).
- `internal/transfer`: Contains `TransferNaive`, `TransferOrdered`, and `TransferWithRetry`.
- `cmd/demo`: Executable demonstration illustrating deadlock occurrence, ordering prevention, and retry recovery.

## Components
- `bank.Account`: Thread-safe entity with balanced state and timeout-capable lock channel.
- `transfer`: Transfer logic implementing different synchronization and mitigation strategies.

## Test Strategy
- Unit tests verifying balances and basic transfer functionality.
- Concurrency tests executing bidirectional transfers to intentionally cause deadlocks.
- Concurrency tests demonstrating lock ordering prevents circular waits.
- Concurrency tests demonstrating retry loops eventually succeed.

## Execution Plan
1. Implement bank account locking mechanics.
2. Implement naive, ordered, and retry-based transfer strategies.
3. Write test suite to validate all 5 findings.
4. Build and execute demo.
5. Record execution results.

## Implementation Decisions
- Used Go channels to simulate row locks with context timeouts (mimics database lock acquisition timeouts / deadlock detection).
- Simulates transaction duration via controlled sleep to reliably observe concurrency behavior.
