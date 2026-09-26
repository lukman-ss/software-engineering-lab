# Docs vs Code

Target Lab: labs/20-zero-downtime-deployment

---

## README vs Implementation

### Claim: Database demonstrates "Expand and Contract" pattern with dual-schema support

Code: `internal/db/db.go` implements `InsertLegacy` (Name-only), `SaveExpand` (FirstName+LastName+Name), and `GetUser` with fallback read logic.
Assessment: **MATCH**

### Claim: Server exposes /healthz/live and /healthz/ready

Code: Both routes registered in `NewServer`. Live always 200. Ready depends on `atomic.Bool`.
Assessment: **MATCH**

### Claim: Server executes configurable preStop delay before shutdown

Code: `preStop time.Duration` field, configured via `NewServer(addr, preStopDelay)`, `select` in `Shutdown`.
Assessment: **MATCH**

### Claim: Worker stops pulling new jobs but completes current active job

Code: `Stop()` sets `stopped=true`, closes channel; goroutines exit between jobs via `ctx.Done()` or channel drain; `time.Sleep` within active job is not interrupted.
Assessment: **MATCH** (with caveat: the claim "completes current active job" is accurate for the normal stop path; on timeout, the currently-running job's sleep is also not interrupted — it still completes. Channel-buffered queued jobs are dropped on timeout.)

### Claim: Demo wires components together, sends SIGTERM, demonstrates zero-downtime draining

Code: `cmd/demo/main.go` starts worker, server, enqueues job, starts in-flight HTTP request, injects SIGTERM via channel, gracefully shuts down. Actual output verified.
Assessment: **MATCH**

### README Running Instructions

README says: `go run ./cmd/demo` and `go test -v ./...` / `go test -race ./...`
Code: Commands work, exit 0.
Assessment: **MATCH**

---

## Engineering Design vs Implementation

### Design Claim: Worker.Stop() signals worker to finish current job and exit gracefully

Implementation: Confirmed. Stop() closes channel (no new jobs fetched), context cancel only on timeout (forces exit between jobs, not within a job sleep).
Assessment: **MATCH**

### Design Claim: All unit and integration tests pass without race conditions

Implementation: 14 tests pass, race detector clean.
Assessment: **MATCH**

### Design Claim: PreStop hook waits for configured delay before closing listeners

Implementation: Confirmed. Shutdown method enforces this ordering.
Assessment: **MATCH**

---

## Mismatches Found

### DOC_CODE_MISMATCH: engineering/03-execution-result.md test count

The execution result recorded by the engineer shows 5 tests passing. The actual test suite at audit time contains 14 tests. The engineering revision phase added 9 tests (`TestDBNotFound`, `TestDBSingleNameLegacy`, `TestDBSaveExpandEmptyFields`, `TestServerPreStopContextCancellation`, `TestServerInvalidDurationFallback`, `TestServerReadyUnreadyTransition`, `TestServerWorkRequestCancellation`, `TestWorkerConcurrency` — confirmed via engineering-revision/02-changes-made.md directory existence). The execution result doc was not updated to reflect the revision additions.

Severity: LOW — the doc is stale, not incorrect in intent; all tests pass in the current state.

### No other mismatches found.
