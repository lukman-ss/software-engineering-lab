# Changes Made

## Revision 1

Audit Issue: GAP-004 — TestConnectionStarvationDueToLeak uses 10ms unconditional sleep for goroutine synchronization
Severity: MEDIUM
Files Changed: tests/pool_test.go
Action: Replaced `time.Sleep(10ms)` with a buffered channel (`acquired`) closed by the leak goroutine's externalCall at the moment the connection is held. The safe call only proceeds after the channel is closed, guaranteeing deterministic ordering regardless of scheduler load. Added a 2s watchdog timeout to prevent hangs.
Verification: `go test -race -count=1 ./...` — PASS
Status: RESOLVED

---

## Revision 2

Audit Issue: GAP-005 — TestOversizedPoolExhaustsServerConnections assertion too weak (errCount > 0)
Severity: LOW
Files Changed: tests/pool_test.go
Action: Strengthened assertion from `errCount == 0` (expect at least 1 failure) to `errCount < 10` (expect at least 10 failures). With server max=10 and 20 concurrent goroutines each holding a connection for 20ms, exactly 10 goroutines must fail. This correctly validates the exhaustion claim.
Verification: `go test -race -count=1 ./...` — PASS
Status: RESOLVED

---

## Revision 3

Audit Issue: GAP-001 — externalCall error path not tested in ProcessOrderSafe or ProcessOrderUnsafeLeak
Severity: LOW
Files Changed: tests/pool_test.go
Action: Added `TestExternalCallErrorPropagation` with two sub-tests (ProcessOrderSafe, ProcessOrderUnsafeLeak). Each sub-test uses an isolated MockDriver and db instance. Verifies that the callErr is returned by the service function, and that after `db.Close()` no connections remain active (proving defer conn.Close() ran).
Verification: `go test -race -count=1 ./...` — PASS
Status: RESOLVED

---

## Revision 4

Audit Issue: GAP-002 — Pre-cancelled context path not tested in ProcessOrderSafe
Severity: LOW
Files Changed: tests/pool_test.go
Action: Added `TestPreCancelledContextProcessOrderSafe`. Creates a context, immediately cancels it, passes it to ProcessOrderSafe. Verifies that a non-nil error is returned. Standard database/sql behavior rejects Conn() on a cancelled context; this test proves the lab's service layer propagates it correctly.
Verification: `go test -race -count=1 ./...` — PASS
Status: RESOLVED
