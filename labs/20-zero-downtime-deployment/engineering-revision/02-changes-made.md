# Changes Made

Target Lab: labs/20-zero-downtime-deployment

---

## Revision 1

Audit Issue: GAP-01 — TOCTOU race between `Enqueue` and `Stop`
Severity: MEDIUM
Files Changed: `internal/worker/worker.go`
Action: Added `enqueueMu sync.Mutex` field to `Worker`. In `Enqueue`: lock `enqueueMu` before checking `stopped` and before sending on `jobChan`. In `Stop`: lock `enqueueMu` while setting `stopped=true` and calling `close(jobChan)`, then unlock before waiting. This eliminates the window between the `stopped.Load()` check and the channel send where `Stop` could close the channel.
Verification: `go test -race ./...` passes. `TestWorkerConcurrentEnqueueStop` (50 iterations, 10 concurrent Enqueue goroutines + Stop) exercises the concurrent path without panic.
Status: RESOLVED

---

## Revision 2

Audit Issue: GAP-03 — No test asserts Enqueue-after-Stop silent rejection
Severity: LOW
Files Changed: `tests/worker_test.go`
Action: Added `TestWorkerEnqueueAfterStop` — calls `Stop`, then `Enqueue`, then asserts `GetCompletedJobs()` returns 0.
Verification: PASS
Status: RESOLVED

---

## Revision 3

Audit Issue: GAP-01 (concurrent path) — No test for concurrent Enqueue + Stop
Severity: MEDIUM
Files Changed: `tests/worker_test.go`
Action: Added `TestWorkerConcurrentEnqueueStop` — 50 iterations each launching 10 goroutines calling `Enqueue` concurrently while `Stop` is called on the same worker. No panic. Race detector clean.
Verification: PASS, race-free
Status: RESOLVED

---

## Revision 4

Audit Issue: GAP-04 — No test for legacy record overwritten with SaveExpand
Severity: LOW
Files Changed: `tests/db_test.go`
Action: Added `TestDBLegacyOverwriteWithExpand` — inserts legacy record "John Doe", reads and verifies FirstName/LastName split, then overwrites with `SaveExpand("1", "Jonathan", "Doe")`, reads again and asserts all three fields.
Verification: PASS
Status: RESOLVED

---

## Revision 5

Audit Issue: GAP-05 — No test for multi-request concurrent drain
Severity: LOW
Files Changed: `tests/server_test.go`
Action: Added `TestServerMultiRequestDrain` on port 8088 — launches 3 concurrent `/work?d=100ms` requests, calls `Shutdown` after 20ms, waits for all goroutines, asserts all 3 returned `WORK COMPLETED`.
Verification: PASS
Status: RESOLVED

---

## Revision 6

Audit Issue: GAP-02 — `engineering/03-execution-result.md` recorded 5 tests; actual count was 14 (pre-revision) then 19 (post-revision)
Severity: LOW
Files Changed: `engineering/03-execution-result.md`
Action: Replaced stale 5-test result block with accurate 19-test run output. Updated race detector timing. Updated Final Engineering Status to `READY_FOR_ENGINEERING_REAUDIT` with revision notes.
Verification: Document now matches actual test suite.
Status: RESOLVED
