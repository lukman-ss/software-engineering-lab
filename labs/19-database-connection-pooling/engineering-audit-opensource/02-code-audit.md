## Finding 1

Location: internal/pool/mockdb.go:41-48
Claimed Behavior: Mock driver enforces server max_connections limit
Observed Implementation: Open() checks maxConnections under mutex and returns ErrServerOverloaded when exceeded
Assessment: PASS
Severity: N/A
Notes: Correct use of mutex for check-then-increment pattern. Atomic operations used for counters.

## Finding 2

Location: internal/pool/mockdb.go:69-76
Claimed Behavior: Connection close decrements active connection count
Observed Implementation: Close() uses atomic decrement when connection not already closed
Assessment: PASS
Severity: N/A
Notes: Thread-safe and idempotent close operation.

## Finding 3

Location: internal/pool/service.go:21-39
Claimed Behavior: ProcessOrderSafe performs external I/O before DB operation to avoid connection leaks
Observed Implementation: External call happens before acquiring DB connection via db.Conn()
Assessment: PASS
Severity: N/A
Notes: Correct ordering prevents holding DB connection during external I/O.

## Finding 4

Location: internal/pool/service.go:41-62
Claimed Behavior: ProcessOrderUnsafeLeak holds DB connection during external I/O (demonstrates leak)
Observed Implementation: External call happens after DB operation but before connection close
Assessment: PASS
Severity: N/A
Notes: Intentionally demonstrates the anti-pattern; connection properly closed via defer.

## Finding 5

Location: tests/pool_test.go
Claimed Behavior: Tests verify connection overhead, pool exhaustion, leaks, and concurrent safety
Observed Implementation: Four tests covering all claimed scenarios with appropriate assertions
Assessment: PASS
Severity: N/A
Notes: Tests are deterministic, use context timeouts appropriately, and validate expected behaviors.

## Finding 6

Location: cmd/demo/main.go
Claimed Behavior: Demo visually demonstrates overhead, exhaustion, and starvation
Observed Implementation: Three demo functions matching test scenarios with timing output
Assessment: PASS
Severity: N/A
Notes: Demo output matches expected patterns; ignores Exec errors for brevity (acceptable in demo).

## Finding 7

Location: internal/pool/mockdb.go:33-49 (Open method)
Claimed Behavior: Connection delay simulates network handshake penalty
Observed Implementation: time.Sleep(d.connectDelay) before attempting connection
Assessment: PASS
Severity: N/A
Notes: Delay is applied per connection attempt; affects unpooled but not pooled scenarios as expected.