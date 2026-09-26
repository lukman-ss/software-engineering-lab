# Code Audit

Target Lab: labs/18-deadlock

## Finding 1

Location: `internal/bank/account.go:21-28`
Claimed Behavior: Simulates resource locking with timeout abort (mimicking database deadlock monitor/detection).
Observed Implementation: Uses channel-based locking with `select` listening on `a.ch` and `ctx.Done()`. Upon timeout, returns `bank.ErrDeadlock`.
Assessment: PASS
Severity: LOW
Notes: Correctly avoids Go runtime unrecoverable deadlock panics while accurately emulating database victim abort behavior.

## Finding 2

Location: `internal/transfer/transfer.go:10-27` (`TransferNaive`)
Claimed Behavior: Simulates deadlock condition through circular wait when executed concurrently in opposing directions.
Observed Implementation: Locks `from`, waits for `delay`, then attempts to lock `to`.
Assessment: PASS
Severity: LOW
Notes: Faithfully establishes the Coffman hold-and-wait and circular-wait conditions.

## Finding 3

Location: `internal/transfer/transfer.go:30-49` (`TransferOrdered`)
Claimed Behavior: Prevents circular wait by enforcing consistent alphabetical lock ordering.
Observed Implementation: Compares `acc1.ID > acc2.ID` and always acquires locks in deterministic order.
Assessment: PASS
Severity: LOW
Notes: Mathematical prevention of circular wait condition properly implemented.

## Finding 4

Location: `internal/transfer/transfer.go:52-71` (`TransferWithRetry`)
Claimed Behavior: Recovers from deadlock aborts via application-level retries.
Observed Implementation: Loops up to `maxRetries` attempting `TransferNaive` with per-attempt timeouts, retrying upon `ErrDeadlock` or `DeadlineExceeded` with backoff.
Assessment: PASS
Severity: LOW
Notes: Effectively recovers transactions. Uses constant backoff (2ms) rather than randomized exponential backoff with jitter, but this limitation is properly documented in engineering notes as a deliberate simplification for readability.
