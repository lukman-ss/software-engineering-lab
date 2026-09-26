# Docs vs Code

Target Lab: labs/20-zero-downtime-deployment

---

## README vs Code

### Claim: Database demonstrates Expand and Contract / Parallel Change pattern

Code: `internal/db/db.go` — `InsertLegacy` (Name only), `SaveExpand` (FirstName+LastName+combined Name), `GetUser` (backward-compatible read with dual-field fallback).
Status: MATCH

### Claim: Server exposes Liveness (/healthz/live) and Readiness (/healthz/ready) probes

Code: Both routes present in `server.go` lines 28-41.
Status: MATCH

### Claim: Server executes configurable preStop delay to simulate load balancer detachment

Code: `preStop time.Duration` field, used in Shutdown at line 91-98.
Status: MATCH

### Claim: Server allows in-flight requests to complete before termination

Code: `http.Server.Shutdown(ctx)` guarantees handler completion. Additionally tracked via `wg`.
Status: MATCH

### Claim: Worker stops pulling new jobs but completes current active job on shutdown

Code: `close(w.jobChan)` stops new dequeues; `wg.Wait()` waits for running goroutines; timeout path cancels context after deadline.
Status: MATCH (with noted caveat that `time.Sleep` inside handler is not preemptible — in-flight always completes)

### Claim: Demo wires all components, simulates startup, in-flight workloads, SIGTERM, graceful drain

Code: `cmd/demo/main.go` does exactly this. Actual execution output confirmed correct behavior.
Status: MATCH

### Claim: `go test -race ./...` passes

Verified by execution: PASS
Status: MATCH

---

## Engineering Design vs Code

### Design claim: Shutdown(ctx, preStopDelay)

Design doc (01-design.md line 47) describes signature as `Shutdown(ctx, preStopDelay)`. Actual signature is `Shutdown(ctx context.Context) error` — preStopDelay is a constructor parameter (`NewServer(addr, preStopDelay)`), not a Shutdown argument.
Status: DOC_CODE_MISMATCH (minor — design doc describes interface intent, not final signature; behavior is equivalent)
Severity: LOW

### Design claim: Worker.Stop() signals worker to finish current job and exit gracefully

Code: `Stop(timeout time.Duration)` — timeout parameter not mentioned in design. Actual implementation adds cooperative timeout/cancel behavior not specified in the original design call signature. This is an enhancement, not a contradiction.
Status: MATCH (enhancement)

---

## Research vs Implementation

Research (both runs) identifies four core patterns: health probes, graceful shutdown with preStop, backward-compatible DB migrations (Expand/Contract), and background worker cooperative drain. All four are implemented and tested.

No fabricated behavior found. No overclaim detected.

---

## Summary

| Area | Status |
|---|---|
| README vs code | MATCH |
| Engineering design vs code | MINOR MISMATCH (Shutdown signature description only) |
| Research claims vs implementation | MATCH |
| Demo output vs claimed behavior | MATCH (verified by execution) |
| Test claims vs actual test logic | MATCH |
