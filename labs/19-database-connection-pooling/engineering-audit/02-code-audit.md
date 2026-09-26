## Finding 1

Location: internal/pool/mockdb.go
Claimed Behavior: Simulates backend connection limits and handshake latency.
Observed Implementation: Uses mutex and atomics to bound connections and track limits. Enforces connectDelay.
Assessment: PASS
Severity: LOW
Notes: `database/sql/driver` implementation correctly simulates connection rejection without network complexity. Concurrency is handled safely.

## Finding 2

Location: internal/pool/service.go
Claimed Behavior: Demonstrates safe vs unsafe connection holding during external I/O.
Observed Implementation: `ProcessOrderSafe` performs external I/O outside DB connection lifecycle. `ProcessOrderUnsafeLeak` holds DB lock during external I/O.
Assessment: PASS
Severity: LOW
Notes: Correctly illustrates the starvation mechanism.