# Code Audit

## Finding 1

Location: internal/pool/mockdb.go:33-48 (MockDriver.Open)
Claimed Behavior: Enforce maxConnections limit and track active/total counts atomically.
Observed Implementation: Uses `d.mu.Lock()` (mutex) to guard the capacity check, then uses `atomic.AddInt32` inside the same lock. The atomic calls inside the mutex are redundant but not unsafe. `ActiveConnections()` and `TotalCreated()` read via `atomic.LoadInt32` without the mutex, which is safe given the values are only written under the mutex or atomically.
Assessment: PASS
Severity: LOW
Notes: Mixed sync strategy (mutex + atomic) is redundant but functionally correct. Reads outside the lock are safe because atomic operations provide memory ordering guarantees.

---

## Finding 2

Location: internal/pool/mockdb.go:69-77 (mockConn.Close)
Claimed Behavior: Double-close must be idempotent — no panic, no double-decrement.
Observed Implementation: Uses `c.mu.Lock()` and guards with `if !c.closed`. Correctly decrements `activeConns` exactly once. Test `TestMockConnDoubleClose` verifies this.
Assessment: PASS
Severity: LOW
Notes: Correct.

---

## Finding 3

Location: internal/pool/service.go:17-34 (ProcessOrderSafe)
Claimed Behavior: External call completes before acquiring DB connection; connection released immediately after query.
Observed Implementation: `externalCall()` executes first, then `db.Conn(ctx)` is called. `defer conn.Close()` ensures release. No connection held during external I/O.
Assessment: PASS
Severity: LOW
Notes: Exactly matches the safe pattern claim.

---

## Finding 4

Location: internal/pool/service.go:37-57 (ProcessOrderUnsafeLeak)
Claimed Behavior: Holds DB connection open during external I/O call — demonstrating the leak pattern.
Observed Implementation: Acquires `conn` via `db.Conn()`, runs `ExecContext`, then calls `externalCall()` while still holding the connection (deferred `conn.Close()` fires only after externalCall returns). This is the correct "unsafe" demonstration.
Assessment: PASS
Severity: LOW
Notes: Intentional unsafe pattern; correctly implemented for demonstration purposes.

---

## Finding 5

Location: tests/pool_test.go:62-63 (TestOversizedPoolExhaustsServerConnections)
Claimed Behavior: `errCount` incremented atomically.
Observed Implementation: `errCount` is declared as `int32` and incremented via `mu.Lock()` + `errCount++` (plain integer addition), NOT `atomic.AddInt32`. A separate `sync.Mutex` guards it. This is safe but inconsistent — not a race condition because the mutex serializes access.
Assessment: PASS
Severity: LOW
Notes: Functionally correct; mixing mutex-guarded plain int with non-atomic access is fine here since mutex provides the necessary memory barrier.

---

## Finding 6

Location: tests/pool_test.go:13-48 (TestDirectConnectionOverhead)
Claimed Behavior: Pooled DB is measurably faster than unpooled DB.
Observed Implementation: Both `unpooledDB` and `pooledDB` are created from the **same** `mockDriver` instance. With `connectDelay=5ms`, unpooled creates 5 connections (5×5ms = ~25ms). Pooled path: the first `Exec` creates a connection and returns it to the idle pool; subsequent 4 reuse it (no delay). The test assertion `pooledDuration < unpooledDuration` is valid.
Assessment: PASS
Severity: LOW
Notes: Test is timing-based. On heavily loaded CI machines it could theoretically flake, but the 5ms/connection gap provides sufficient margin.

---

## Finding 7

Location: tests/pool_test.go:88-124 (TestConnectionStarvationDueToLeak)
Claimed Behavior: Unsafe goroutine holds single-slot pool; second request times out.
Observed Implementation: Pool size = 1. Goroutine calls `ProcessOrderUnsafeLeak` which acquires connection then signals `acquired` channel, then sleeps 100ms holding connection. Main goroutine waits for `acquired`, then tries `ProcessOrderSafe` with 20ms timeout. 20ms < 100ms remaining hold = timeout guaranteed.
Assessment: PASS
Severity: LOW
Notes: Timing logic is sound. Channel synchronization ensures the goroutine has acquired the connection before the second request starts.

---

## Finding 8

Location: tests/pool_test.go:283-311 (TestTotalCreatedPoolReuse)
Claimed Behavior: Pooled execution creates fewer total connections than unpooled.
Observed Implementation: Both phases share the same `mockDriver`, so `totalCreated` is cumulative. `pooledTotal = mockDriver.TotalCreated() - unpooledTotal` correctly isolates the pooled phase count. Sequential (not concurrent), so no race. Unpooled=5 connections, pooled=1 connection expected.
Assessment: PASS
Severity: LOW
Notes: Correct arithmetic for isolation.

---

## Finding 9

Location: tests/pool_test.go:199-235 (TestExternalCallErrorPropagation)
Claimed Behavior: Error from externalCall propagates correctly; no active connections remain after db.Close().
Observed Implementation: For `ProcessOrderSafe`: externalCall fails before `db.Conn()` is called, so no connection is ever acquired — active count stays 0. For `ProcessOrderUnsafeLeak`: connection acquired, ExecContext succeeds, externalCall fails; `defer conn.Close()` still fires, returning connection to pool; `db.Close()` drains pool → active=0.
Assessment: PASS
Severity: LOW
Notes: Both error paths correctly verified.

---

## Finding 10

Location: cmd/demo/main.go:42-48 (demoDirectOverhead warm-up)
Claimed Behavior: Pool is warmed up before timing the pooled section.
Observed Implementation: `pooledDB.Exec("SELECT 1")` is called once outside the timed loop, priming the idle connection. Then the 5-iteration loop runs against warmed pool. This is correct and prevents first-connection overhead from skewing results.
Assessment: PASS
Severity: LOW
Notes: Demo output (54ms unpooled, 12µs pooled) is reproducible and matches documented results.

---

## Finding 11

Location: go.mod:3
Claimed Behavior: Standard Go module, no external dependencies.
Observed Implementation: `go 1.26.7` — this Go version does not exist as of the audit date (latest stable is 1.23.x). The `go` directive in go.mod specifies a minimum version requirement; Go toolchains handle forward compatibility for newer directives by using the installed toolchain. Build succeeded, so the installed toolchain (likely 1.21+) handled it correctly.
Assessment: WARNING
Severity: LOW
Notes: `go 1.26.7` is a non-existent version at the time of this audit. This is likely a speculative/future-dated version number. It does not break compilation but is technically inaccurate documentation.

---

## Finding 12

Location: internal/pool/mockdb.go — no transaction rollback testing
Claimed Behavior: N/A (no rollback claim in this lab)
Observed Implementation: `mockTx.Rollback()` is implemented (returns nil) but no test exercises rollback behavior. The lab does not claim to test transaction rollback.
Assessment: PASS
Severity: LOW
Notes: Rollback path exists structurally; no claim made about it so absence of test is acceptable.
