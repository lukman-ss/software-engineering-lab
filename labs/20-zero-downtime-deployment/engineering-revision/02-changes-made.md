# Changes Made

Target Lab: labs/20-zero-downtime-deployment

---

## Revision 1

Audit Issue: GAP-01 / Timing Flakiness in `TestWorkerShutdownTimeout`
Severity: MEDIUM
Files Changed: `internal/worker/worker.go`, `tests/worker_test.go`
Action: In `internal/worker/worker.go`, replaced un-cancellable `time.Sleep` with `select` listening to `time.After(job.Duration)` and `w.ctx.Done()`. Also added a non-blocking `w.ctx.Done()` check immediately upon dequeuing a job from `jobChan` to avoid processing queued jobs if drain timeout has elapsed. In `tests/worker_test.go`, tuned job durations in `TestWorkerShutdownTimeout` (`job-slow`: 10ms, `job-dropped`: 100ms, drain timeout: 25ms) ensuring deterministic abandonment of timed-out queued jobs.
Verification: Ran `go test -count=10 -run TestWorkerShutdownTimeout ./...` (100% pass).
Status: RESOLVED

---

## Revision 2

Audit Issue: GAP-02 / Server Listener Bind Timing Race in `tests/server_test.go`
Severity: LOW
Files Changed: `tests/server_test.go`
Action: Replaced hardcoded `time.Sleep(50 * time.Millisecond)` server startup delays across all server test functions with a resilient `waitForServerReady(addr string)` polling helper that polls `/healthz/live` with a 2-second deadline and 5ms retry interval.
Verification: Ran `go test -count=10 -run TestServer ./...` without dial errors or connection resets.
Status: RESOLVED

---

## Revision 3

Audit Issue: GAP-03 / Documentation Drift in Test Results
Severity: LOW
Files Changed: `engineering/03-execution-result.md`
Action: Synchronized `engineering/03-execution-result.md` with accurate 18 test names, test counts (5 DB, 8 server, 5 worker), clean race detector status, and demo execution logs.
Verification: All counts and test function references match implementation code.
Status: RESOLVED
