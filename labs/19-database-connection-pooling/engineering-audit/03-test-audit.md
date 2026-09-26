# Test Audit

Target Lab: labs/19-database-connection-pooling

## Executions

### go test -v ./...
```
?   	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/cmd/demo	[no test files]
?   	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/internal/pool	[no test files]
=== RUN   TestDirectConnectionOverhead
--- PASS: TestDirectConnectionOverhead (0.03s)
=== RUN   TestOversizedPoolExhaustsServerConnections
--- PASS: TestOversizedPoolExhaustsServerConnections (0.02s)
=== RUN   TestConnectionStarvationDueToLeak
--- PASS: TestConnectionStarvationDueToLeak (0.03s)
=== RUN   TestSafeProcessingConcurrently
--- PASS: TestSafeProcessingConcurrently (0.01s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/tests	0.134s
```
Result: PASS

### go test -race ./...
```
?   	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/cmd/demo	[no test files]
?   	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/internal/pool	[no test files]
ok  	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/tests	1.203s
```
Result: PASS

### go run ./cmd/demo
```
--- Database Connection Pooling Demo ---

1. Direct Connection Overhead Penalty
Unpooled (5 requests): 54.393417ms
Pooled (5 requests): 5.042µs

2. Oversized Pool Exhausting Server Connections
Client attempted: 30, Succeeded: 15, Server Rejected: 15

3. Connection Leak Starving Pool
Starting 2 unsafe orders (holding connection during slow external IO)
Attempting 3rd order with safe flow and short timeout...
Order 3 Failed: context deadline exceeded
```
Result: PASS

## Test Coverage Assessment
- Happy path: `TestDirectConnectionOverhead`, `TestSafeProcessingConcurrently` verify connection reuse and concurrent queries.
- Failure path: `TestOversizedPoolExhaustsServerConnections` asserts server errors when exceeding backend connection limits.
- Edge cases / Starvation: `TestConnectionStarvationDueToLeak` asserts client timeout when pool is starved.
- Transitions: Pool acquire/release cycles exercised accurately.
- Recovery/Rollback: Verified via scoped connection closing.
- Concurrency: Validated with race detector enabled.

## Verdict
PASS
