# Code Audit

Environment:
  - go version go1.26.7 darwin/arm64
  - module github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling

## Finding 1

Location: internal/pool/mockdb.go (MockDriver.Open)
Claimed Behavior: MockDriver emulates a database server that enforces `max_connections`;
  once the server's active connection count reaches the cap, Open returns ErrServerOverloaded.
Observed Implementation: Open checks `d.activeConns >= d.maxConnections` while holding
  `d.mu`, but performs the increment (`atomic.AddInt32(&d.activeConns, 1)`) *after*
  releasing `d.mu`. Close decrements via `atomic.AddInt32` without taking `d.mu`.
Assessment: WARNING
Severity: LOW
Notes: This is a check-then-act (TOCTOU) gap, not a data race — the race detector passes
  clean because atomics are used throughout. Two Open goroutines can both observe a
  sub-cap value, release the lock, and both increment, allowing the cap to be exceeded by
  a small margin under high concurrency. For a mock that exists only to *illustrate* server
  overload, the enforcement is approximate but functionally adequate: the tests assert
  "at least some errors occur", not an exact saturation count, and the demo reproduces
  the rejection behavior. Fix would be to move the atomic.AddInt32 *inside* d.mu so the
  cap is evaluated and enforced atomically.

## Finding 2

Location: internal/pool/service.go:10 (var ErrAcquireTimeout)
Claimed Behavior: (Declared) ErrAcquireTimeout is defined as the sentinel returned on
  acquire timeout.
Observed Implementation: ErrAcquireTimeout is never referenced — not assigned, not
  returned, not compared, not tested. ProcessOrderSafe/ProcessOrderUnsafeLeak rely on
  the *sql.DB.Conn context deadline (which surfaces a context.DeadlineExceeded /
  context.Canceled error) rather than returning ErrAcquireTimeout.
Assessment: WARNING
Severity: LOW
Notes: Dead code / unexported-but-unused sentinel. The timeout behavior itself is correct
  (verified: tests/pool_test.go::TestConnectionStarvationDueToLeak expects the failure and
  passes). The sentinel should either be returned on timeout or removed. Does not affect
  correctness of demonstrated behavior.

## Finding 3

Location: internal/pool/mockdb.go (mockRows.Next) and mockStmt
Claimed Behavior: (Support) A `database/sql/driver`-level mock implementing the driver
  contract for database/sql to use via sql.OpenDB.
Observed Implementation: mockConn implements Prepare/Close/Begin; mockStmt implements
  Close/NumInput/Exec/Query; mockTx implements Commit/Rollback; mockRows implements
  Columns/Close/Next. MockConnector wraps MockDriver so OpenDB uses the instance directly
  (no global driver registration). This matches database/sql/driver v1 interface for go1.26.
Assessment: PASS
Severity: N/A
Notes: No compile/vet errors. sql.DB interacts with the mock as expected during the demo
  and tests.

## Finding 4

Location: internal/pool/service.go (OrderService)
Claimed Behavior: ProcessOrderSafe performs external I/O *outside* the connection
  lifetime; ProcessOrderUnsafeLeak holds the connection open during external I/O.
Observed Implementation: Safe flow calls externalCall() before acquiring
  s.db.Conn(ctx), then runs a short UPDATE and `defer conn.Close()`. Unsafe flow
  acquires Conn, updates status, THEN calls externalCall() while Conn is still held,
  and only closes via defer after.
Assessment: PASS
Severity: N/A
Notes: The unsafe path is intentionally unsafe (documented in doc comments). Resource
  cleanup is correctly deferred, so even the "leaky" path releases the connection on
  return — the leak is logical (connection held longer than necessary) not a file/descriptor leak.

## Finding 5

Location: cmd/demo/main.go
Claimed Behavior: Three demos — direct overhead, oversized pool exhausting server
  connections, and leak-induced pool starvation with a timeout.
Observed Implementation: demoDirectOverhead (unpooled 0 idle / pooled 5 idle), 
  demoOversizedPool (server max=15, client pool=50, 30 concurrent, 10ms hold),
  demoConnectionLeak (pool=2, 2 unsafe holders for 500ms, 3rd safe request 100ms timeout).
Assessment: PASS
Severity: N/A
Notes: Concurrency uses sync.WaitGroup + atomic counters. Timeouts via
  context.WithTimeout. Matches README and engineering/03-execution-result.md claims.

## Finding 6

Location: tests/pool_test.go
Claimed Behavior: Covers direct overhead, pool exhaustion, leak starvation, concurrent safe path.
Observed Implementation: 4 tests exactly covering those scenarios. Timing-based assertions
  use connect delays / sleeps large enough that pooled is reliably faster.
Assessment: PASS (with caveat)
Severity: LOW
Notes: TestDirectConnectionOverhead asserts `pooledDuration < unpooledDuration` with a 5ms
  connect delay per connection; on an extremely slow CI box the unpooled loop could be
  dominated by scheduler noise, but 5ms*5=25ms floor makes failure unlikely. Acceptable for
  a pedagogical lab.
