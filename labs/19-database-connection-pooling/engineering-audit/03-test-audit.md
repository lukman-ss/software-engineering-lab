# Test Audit

## Coverage Assessment

- **Happy Path:** Covered by `TestDirectConnectionOverhead` and `TestSafeProcessingConcurrently`.
- **Failure Path:** Covered by `TestOversizedPoolExhaustsServerConnections` (server rejection) and `TestConnectionStarvationDueToLeak` (context timeout).
- **Edge Cases:** Missing deadlock test (Research Finding 11: pool-locking deadlock).
- **Transitions:** Pool warmup transitions covered implicitly in overhead test.
- **Recovery:** Implicitly covered; context cancellations correctly unblock pool acquisition.
- **Concurrency:** Covered by `TestSafeProcessingConcurrently` and `TestOversizedPoolExhaustsServerConnections`.
- **Negative Cases:** Covered (timeout, rejection).

A passing test suite can still be weak. In this lab, tests successfully capture the business logic constraints (timeout and rejection) but expose a race condition in the underlying driver mock.

## Required Execution Results

### 1. `go test ./...`
```text
ok      github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/tests    0.185s
```

### 2. `go test -race ./...`
```text
==================
WARNING: DATA RACE
Write at 0x00c00032e00c by goroutine 42:
  sync/atomic.AddInt32()
...
  github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/internal/pool.(*mockConn).Close()
      /Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/19-database-connection-pooling/internal/pool/mockdb.go:73

Previous read at 0x00c00032e00c by goroutine 59:
  github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/internal/pool.(*MockDriver).Open()
      /Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/19-database-connection-pooling/internal/pool/mockdb.go:41
...
FAIL    github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/tests    0.178s
```
**Result: FAIL**

### 3. `go run ./cmd/demo`
```text
--- Database Connection Pooling Demo ---

1. Direct Connection Overhead Penalty
Unpooled (5 requests): 55.086041ms
Pooled (5 requests): 14.792µs

2. Oversized Pool Exhausting Server Connections
Client attempted: 30, Succeeded: 15, Server Rejected: 15

3. Connection Leak Starving Pool
Starting 2 unsafe orders (holding connection during slow external IO)
Attempting 3rd order with safe flow and short timeout...
Order 3 Failed: context deadline exceeded
```
**Result: PASS**
