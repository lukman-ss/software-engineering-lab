# Test Audit

## Test Inventory

| Test | Path | Category |
|---|---|---|
| TestDirectConnectionOverhead | tests/pool_test.go:13 | Happy path / performance |
| TestOversizedPoolExhaustsServerConnections | tests/pool_test.go:51 | Failure path / exhaustion |
| TestConnectionStarvationDueToLeak | tests/pool_test.go:88 | Failure path / leak |
| TestSafeProcessingConcurrently | tests/pool_test.go:126 | Happy path / concurrency |
| TestPoolLockingDeadlock | tests/pool_test.go:154 | Edge case / deadlock |
| TestMockConnDoubleClose | tests/pool_test.go:176 | Edge case / idempotency |
| TestExternalCallErrorPropagation | tests/pool_test.go:199 | Failure path / error propagation |
| TestPreCancelledContextProcessOrderSafe | tests/pool_test.go:237 | Edge case / cancelled context |
| TestUnsafeLeakExecContextFailure | tests/pool_test.go:253 | Failure path / pool-exhausted exec |
| TestTotalCreatedPoolReuse | tests/pool_test.go:283 | Happy path / connection reuse count |

## Coverage Assessment

### Happy Path
- COVERED: Pooled connections reuse (TestDirectConnectionOverhead, TestTotalCreatedPoolReuse)
- COVERED: Safe concurrent processing (TestSafeProcessingConcurrently)

### Failure Path
- COVERED: Oversized pool exhausting server limits (TestOversizedPoolExhaustsServerConnections)
- COVERED: Connection starvation from leak (TestConnectionStarvationDueToLeak)
- COVERED: ExecContext failure when pool exhausted (TestUnsafeLeakExecContextFailure)
- COVERED: External call error propagation for both safe and unsafe paths (TestExternalCallErrorPropagation)

### Edge Cases
- COVERED: Pool deadlock / double acquisition (TestPoolLockingDeadlock)
- COVERED: Double-close idempotency on mock connection (TestMockConnDoubleClose)
- COVERED: Pre-cancelled context (TestPreCancelledContextProcessOrderSafe)

### Concurrency / Race
- COVERED: TestSafeProcessingConcurrently runs 20 goroutines against 5-slot pool
- COVERED: Race detector: `go test -race ./...` — all PASS

### Rollback
- NOT TESTED: No transaction rollback test. Lab does not claim rollback behavior — ACCEPTABLE.

### Negative Cases
- COVERED: Nil externalCall handled gracefully in both ProcessOrderSafe and ProcessOrderUnsafeLeak (guarded by `if externalCall != nil`)

## Execution Results (Actual, This Audit)

```
go test -v ./...

=== RUN   TestDirectConnectionOverhead
--- PASS: TestDirectConnectionOverhead (0.03s)
=== RUN   TestOversizedPoolExhaustsServerConnections
--- PASS: TestOversizedPoolExhaustsServerConnections (0.02s)
=== RUN   TestConnectionStarvationDueToLeak
--- PASS: TestConnectionStarvationDueToLeak (0.02s)
=== RUN   TestSafeProcessingConcurrently
--- PASS: TestSafeProcessingConcurrently (0.01s)
=== RUN   TestPoolLockingDeadlock
--- PASS: TestPoolLockingDeadlock (0.05s)
=== RUN   TestMockConnDoubleClose
--- PASS: TestMockConnDoubleClose (0.00s)
=== RUN   TestExternalCallErrorPropagation
    --- PASS: TestExternalCallErrorPropagation/ProcessOrderSafe (0.00s)
    --- PASS: TestExternalCallErrorPropagation/ProcessOrderUnsafeLeak (0.00s)
--- PASS: TestExternalCallErrorPropagation (0.00s)
=== RUN   TestPreCancelledContextProcessOrderSafe
--- PASS: TestPreCancelledContextProcessOrderSafe (0.00s)
=== RUN   TestUnsafeLeakExecContextFailure
--- PASS: TestUnsafeLeakExecContextFailure (0.02s)
=== RUN   TestTotalCreatedPoolReuse
--- PASS: TestTotalCreatedPoolReuse (0.00s)
PASS
ok  	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/tests
```

```
go test -race -v ./...

All 10 tests: PASS
Race detector: no data races detected
ok  	github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/tests	1.279s
```

## Weaknesses / Observations

1. `TestDirectConnectionOverhead` and `TestTotalCreatedPoolReuse` are timing/count-based. They are structurally sound, but on a machine with scheduler interference the timing test could theoretically flake. The 5ms/connection gap provides reasonable margin.

2. `TestOversizedPoolExhaustsServerConnections` asserts `errCount >= 10`. With 20 goroutines and server max=10, at least 10 must fail simultaneously. Given all goroutines hold connections for 20ms and are launched near-simultaneously, this is reliable in practice.

3. No benchmark tests (`func Benchmark...`). Not required by the design, but would strengthen performance claims.

4. No test for `ProcessOrderSafe` with a nil externalCall and a working DB — trivially covered by the starvation test (nil externalCall path), but not as an explicit dedicated test.
