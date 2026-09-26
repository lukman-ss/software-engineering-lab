# Test Audit

Test suite location: tests/pool_test.go
Test package: `tests` (external test package consuming the pool package via import).

Command results (executed against the real working directory):
  $ go test -v ./...
  → 4 PASS, 0 FAIL  (exit 0)
  $ go test -race -v ./...
  → 4 PASS, 0 FAIL, 0 DATA RACES  (exit 0)

Coverage of required test dimensions:

- Happy path: PASS
  TestSafeProcessingConcurrently runs 20 concurrent safe orders on a pool of 5 with a 2s
  timeout and expects all to succeed.

- Failure path: PASS
  TestConnectionStarvationDueToLeak asserts the safe order FAILS (context deadline
  exceeded) when the pool is starved by an unsafe holder.

- Edge cases: PARTIAL
  No test exercises the exact server-saturation boundary of MockDriver, no test asserts
  the value of ErrAcquireTimeout, no test exercises mockRows/mockTx rollback behavior.

- Transitions: PARTIAL
  No explicit test for idle↔open connection reuse accounting beyond the overhead test's
  timing comparison.

- Recovery / rollback: NOT COVERED
  No test that a failed operation rolls back or returns the connection to the pool for
  reuse. (processOrderUnsafeLeak uses defer conn.Close(), so recovery is implicit.)

- Concurrency: PASS
  Three tests exercise concurrency under the race detector; none report races.

- Negative cases: PASS
  TestOversizedPoolExhaustsServerConnections expects ≥1 failure from 20 concurrent
  acquirers against a server cap of 10; TestConnectionStarvationDueToLeak expects a error.

## Finding 1

Location: tests/pool_test.go
Claimed Behavior: Connection starvation is observable when an oversized pool exceeds
  server max_connections.
Observed Implementation: TestOversizedPoolExhaustsServerConnections uses server max=10,
  client pool=20, 20 concurrent acquirers, asserts errCount > 0.
Assessment: PASS
Severity: N/A
Notes: Robustly asserted (does not pin an exact count, accommodating the check-then-act
  gap documented in code audit Finding 1).

## Finding 2

Location: tests/pool_test.go
Claimed Behavior: Safe processing survives concurrency.
Observed Implementation: TestSafeProcessingConcurrently, 20 goroutines, pool of 5,
  2s timeout, expects no errors.
Assessment: PASS
Severity: N/A

## Finding 3

Location: tests/pool_test.go
Claimed Behavior: (Gap) ErrAcquireTimeout is never referenced in any test.
Observed Implementation: No test asserts that a pool acquire timeout yields ErrAcquireTimeout
  or any specific sentinel.
Assessment: WARNING
Severity: LOW
Notes: Related to code-audit Finding 2. The timeout *behavior* is tested (the order fails),
  but no test pins the error identity. Low impact for a pedagogical lab.

## Finding 4

Location: tests/pool_test.go
Claimed Behavior: (Gap) No negative test for "connection properly returned after error".
Observed Implementation: ProcessOrderUnsafeLeak updates status then calls externalCall;
  if externalCall errors the connection is still closed via defer. No test asserts that a
  subsequent acquirer can reuse that connection (i.e., no pool-recovery-after-failure test).
Assessment: WARNING
Severity: LOW
Notes: Implicit via defer; not a correctness bug, but a coverage gap relative to the
  "failure handling / recovery" audit dimensions.
