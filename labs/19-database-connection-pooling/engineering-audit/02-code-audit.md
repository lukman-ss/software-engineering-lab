# Code Audit

## Finding 1

Location: `internal/pool/mockdb.go:41`
Claimed Behavior: Safe concurrent tracking of active connections.
Observed Implementation: In `MockDriver.Open()`, `d.activeConns` is read directly without atomic operation: `if d.maxConnections > 0 && d.activeConns >= d.maxConnections`. Concurrently, `mockConn.Close()` modifies `d.activeConns` via `atomic.AddInt32(&c.driver.activeConns, -1)` without acquiring `d.mu`. This causes a real data race caught by Go's race detector during concurrent connection release and creation.
Assessment: FAIL
Severity: HIGH
Notes: Violates Go memory safety and concurrency guarantees. Must use `atomic.LoadInt32(&d.activeConns)` in `Open()`.

## Finding 2

Location: `internal/pool/mockdb.go:41-43`
Claimed Behavior: Simulates server overload rejection (`FATAL: sorry, too many clients already`).
Observed Implementation: Returns `ErrServerOverloaded` when active connections reach `maxConnections`.
Assessment: PASS
Severity: LOW
Notes: Properly models PostgreSQL connection slot exhaustion behavior under load.

## Finding 3

Location: `internal/pool/service.go:21-50`
Claimed Behavior: Safe flow releases/bounds DB connection during external I/O; unsafe flow holds connection across network call.
Observed Implementation: `ProcessOrderSafe` executes external call outside of connection acquisition scope. `ProcessOrderUnsafeLeak` holds `*sql.Conn` open across `externalCall()`.
Assessment: PASS
Severity: LOW
Notes: Correctly illustrates the anti-pattern identified in Research Finding 9.

## Finding 4

Location: `internal/pool/mockdb.go`
Claimed Behavior: Sizing beyond hardware limits causes throughput degradation.
Observed Implementation: Mock driver only models hard slot rejection (`ErrServerOverloaded`), not throughput degradation or latency curve (the performance "knee" identified in Research Finding 4).
Assessment: WARNING
Severity: MEDIUM
Notes: Engineering design doc states "Test validates throughput degradation or max connection enforcement". Implementation chose max connection enforcement only.
