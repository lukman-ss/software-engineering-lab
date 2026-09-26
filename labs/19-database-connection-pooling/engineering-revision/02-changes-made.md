# Changes Made

## Revision 1

Audit Issue: GAP-01 — `engineering/03-execution-result.md` listed 6 tests; actual suite had 8 (TestExternalCallErrorPropagation, TestPreCancelledContextProcessOrderSafe added in prior revision).
Severity: LOW
Files Changed: `engineering/03-execution-result.md`
Action: Updated test list in both `go test -v` and `go test -race` output blocks to reflect all 10 tests (8 prior + 2 new from this revision).
Verification: File updated; matches actual `go test -v ./...` output.
Status: RESOLVED

---

## Revision 2

Audit Issue: GAP-02 — `ErrAcquireTimeout` exported but never returned by any function and never tested (dead code).
Severity: LOW
Files Changed: `internal/pool/service.go`
Action: Removed `ErrAcquireTimeout` variable declaration and the `"errors"` import. The actual timeout error is `context.DeadlineExceeded` from `db.Conn(ctx)`, which is already surfaced correctly.
Verification: `go build ./...` passes. No callers reference `ErrAcquireTimeout` anywhere in the codebase.
Status: RESOLVED

---

## Revision 3

Audit Issue: GAP-03 — No test for `ProcessOrderUnsafeLeak` when the initial `db.Conn` call fails (pool exhausted / context timeout before connection is acquired).
Severity: LOW
Files Changed: `tests/pool_test.go`
Action: Added `TestUnsafeLeakExecContextFailure`. Test saturates a pool-of-1 with one held connection, then calls `ProcessOrderUnsafeLeak` with a 20ms context — verifies error is returned and all connections are released after `db.Close()`.
Verification: `go test -v -run TestUnsafeLeakExecContextFailure ./...` → PASS.
Status: RESOLVED

---

## Revision 4

Audit Issue: GAP-05 — `TotalCreated()` counter is untested as proof that pooling results in fewer connections than unpooled use.
Severity: LOW
Files Changed: `tests/pool_test.go`
Action: Added `TestTotalCreatedPoolReuse`. Runs 5 queries unpooled (SetMaxIdleConns(0)) and 5 queries pooled (SetMaxIdleConns(5)), then asserts pooled TotalCreated < unpooled TotalCreated using the same MockDriver instance.
Verification: `go test -v -run TestTotalCreatedPoolReuse ./...` → PASS.
Status: RESOLVED

---

## Revision 5 (skipped — GAP-04)

Audit Issue: GAP-04 — No test for context cancellation during slow `connectDelay` (MockDriver.Open uses `time.Sleep`, not context-aware).
Severity: LOW
Action: SKIPPED. The audit itself noted this is acceptable for a mock driver. The `sql` package's pool layer handles context cancellation at a higher level. Adding a test would require changes to the mock's `Open` signature and the `driver.Connector` interface which is out of scope for a pedagogical mock.
Status: UNRESOLVED (accepted — acceptable for mock)
