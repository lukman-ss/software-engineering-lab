# Docs vs Code Audit

Comparing README, engineering design/notes, code, tests, and demo output.

## README ↔ Code

### Claim: Components — Database (Expand and Contract)
README: "In-memory database... Expand and Contract pattern (Parallel Change). Mock storage supports writing dual schema versions (legacy `Name` and modern `FirstName`/`LastName`), and transparent fallback logic when reading records created by differing versions."

Code (`internal/db/db.go`): ✅ matches — `SaveExpand`/`InsertLegacy`/`GetUser` implement legacy+modern fields and fallback read. No mismatch.

### Claim: Components — Server
README: "Exposes Liveness and Readiness probes. When a shutdown signal is sent, it executes a configurable `preStop` delay... graceful shutdown ensuring any in-flight requests complete before termination."

Code (`internal/server/server.go`): ✅ matches — probes + `Shutdown(ctx)` with `preStop` + `http.Server.Shutdown` + `wg.Wait`.

### Claim: Components — Worker
README: "Upon receiving a shutdown signal, it stops pulling new jobs but continues processing the current active job until completion."

Code (`internal/worker/worker.go`): ✅ matches — `Stop()` closes chan, in-flight jobs drain to completion (timeout escalation).

### Claim: Components — Demo
README: "A CLI orchestrator that wires these components together, simulates startup initialization, executes in-flight workloads, and sends a termination signal to demonstrate zero-downtime draining behavior."

Code (`cmd/demo/main.go`): ✅ matches — readiness gating, concurrent `/work?d=2s`, SIGTERM self-signal, graceful shutdown, worker drain.

## Design (`engineering/01-design.md`) ↔ Code

- Success criteria #1 readiness/traffic: ✅ `TestServerProbes`, `TestServerReadyUnreadyTransition`.
- Success criteria #2 in-flight completion: ✅ `TestServerGracefulShutdown`, `TestServerMultiRequestDrain`.
- Success criteria #3 preStop delay: ✅ `TestServerPreStopHook`.
- Success criteria #4 worker completes active job: ✅ `TestWorkerGracefulShutdown`.
- Success criteria #5 expand/contract compatibility: ✅ DB tests.
- Success criteria #6 race-free: ✅ `go test -race ./...` clean.

## Execution Result (`engineering/03-execution-result.md`) ↔ Demo

Documented result:
```
Starting Zero-Downtime Deployment Demo
Background worker started
Worker N starting job DemoJob-1
Server starting on 127.0.0.1:8080
Application is ready to receive traffic
Simulating SIGTERM ...
SIGTERM received ...
Server received shutdown request
Server marked unready ...
Executing preStop sleep for 1s ...
Worker N finished job DemoJob-1
Initiating graceful shutdown of HTTP listeners...
Client request completed with status: 200
All in-flight requests completed. Server stopped gracefully.
Worker receiving stop signal ...
Worker gracefully stopped
Demo finished cleanly. Zero downtime achieved.
```

Live run captured to `/tmp/demo-output.log` (17 lines). Every semantic line present and in order; only log timestamps differ (second-precision) — expected, not fabrication. Exit code 0. ✅ no FAKE_DEMO.

## Execution Result ↔ Tests

Documented: 18 tests (5 DB, 8 server, 5 worker), all PASS. Actual: ✅ confirmed — 18/18 PASS (`go test -v -count=1`), race-clean.

## Findings

| Type | Doc location | Issue |
|------|--------------|-------|
| DOC_CODE_MISMATCH | none | None found — README and engineering notes match code. |
| TEST_CLAIM_MISMATCH | none | None found — tests cover all documented claims. |
| RESEARCH_IMPLEMENTATION_MISMATCH | none (research excluded per pipeline) | N/A. |

Only note: `engineering-notes` reference "18 tests (5 DB, 8 server, 5 worker)" — matches actual. ✅
