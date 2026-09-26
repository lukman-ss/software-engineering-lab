## Finding 1

Location: `internal/pool/mockdb.go`
Claimed Behavior: Simulate connection establishment latency and backend max_connections limits.
Observed Implementation: Uses a custom `database/sql/driver` with atomic counters and mutexes to enforce `maxConnections` and simulated `connectDelay`.
Assessment: PASS
Severity: LOW
Notes: Correctly mocks backend limitations. Concurrency safe using atomic operations.

## Finding 2

Location: `internal/pool/service.go`
Claimed Behavior: Safe connection management vs unsafe connection leakage during external I/O.
Observed Implementation: `ProcessOrderSafe` performs external I/O outside DB connection lifecycle. `ProcessOrderUnsafeLeak` holds DB connection open during external call.
Assessment: PASS
Severity: LOW
Notes: Clearly demonstrates the connection leak anti-pattern.
