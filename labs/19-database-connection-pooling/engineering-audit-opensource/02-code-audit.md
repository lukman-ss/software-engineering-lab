# Code Audit

Target Lab: labs/19-database-connection-pooling
Files: internal/pool/mockdb.go, internal/pool/service.go, cmd/demo/main.go

## Finding 1

Location: internal/pool/mockdb.go:33-49 (MockDriver.Open)
Claimed Behavior: Simulate server max_connections + connect latency.
Observed Implementation: Sleep outside lock, check+increment under mu.Lock. Atomic ops inside mutex redundant but correct.
Assessment: PASS
Severity: LOW
Notes: No race. Race detector clean.

## Finding 2

Location: internal/pool/mockdb.go:145-147 (MockConnector.Connect)
Claimed Behavior: Open DB connections.
Observed Implementation: Ignores ctx, calls blocking Open with time.Sleep.
Assessment: WARNING
Severity: LOW
Notes: Cancel during connectDelay not honored by driver; sql package still enforces Conn(ctx) timeout. Limits realism only.

## Finding 3

Location: internal/pool/mockdb.go:69-77 (mockConn.Close)
Claimed Behavior: Connection release.
Observed Implementation: closed flag + mutex, idempotent, decrements once.
Assessment: PASS
Severity: LOW
Notes: TestMockConnDoubleClose proves it.

## Finding 4

Location: internal/pool/mockdb.go:79-113 (mockTx, mockStmt.Exec)
Claimed Behavior: Query/tx simulation.
Observed Implementation: Exec always success, Tx Commit/Rollback no-ops, never exercised by service (service uses Conn+ExecContext, no Begin).
Assessment: WARNING
Severity: LOW
Notes: Dead tx path. No failure injection for query errors. Unnecessary surface, harmless.

## Finding 5

Location: internal/pool/service.go:17-34 (ProcessOrderSafe)
Claimed Behavior: External I/O outside connection hold, short bounded DB op.
Observed Implementation: externalCall first, then db.Conn + ExecContext + defer Close. Error propagated.
Assessment: PASS
Severity: LOW
Notes: Correct ordering. Proven by starvation + concurrency tests.

## Finding 6

Location: internal/pool/service.go:37-57 (ProcessOrderUnsafeLeak)
Claimed Behavior: Hold connection during external I/O to starve pool.
Observed Implementation: db.Conn, Exec, then externalCall while holding, defer Close releases after.
Assessment: WARNING
Severity: LOW
Notes: Name says Leak but it releases via defer; it holds, not leaks. Behavior still proves starvation. externalCall takes no ctx so not cancellable; fine for demo.

## Finding 7

Location: internal/pool/service.go error paths
Claimed Behavior: Error propagation + cleanup.
Observed Implementation: Both paths return externalCall/Exec errors directly; defer conn.Close() always runs.
Assessment: PASS
Severity: LOW
Notes: Proven by TestExternalCallErrorPropagation checking ActiveConnections()==0.

## Finding 8

Location: concurrency safety (MockDriver.mu + atomics, mockConn.mu, service stateless)
Claimed Behavior: Safe concurrent use.
Observed Implementation: Shared counters guarded; service holds no mutable state.
Assessment: PASS
Severity: LOW
Notes: go test -race clean. TestSafeProcessingConcurrently 20 goroutines x pool 5 passes.
