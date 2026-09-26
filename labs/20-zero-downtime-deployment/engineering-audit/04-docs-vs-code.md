# Docs vs Code

Target Lab: labs/20-zero-downtime-deployment

## Analysis

### Server
- Documented: "Exposes Liveness and Readiness probes. When a shutdown signal is sent, it executes a configurable preStop delay... then performs a graceful shutdown"
- Code: Implemented correctly. Liveness (`/healthz/live`) returns 200 OK. Readiness (`/healthz/ready`) respects atomic boolean flag. Shutdown triggers `SetReady(false)`, sleep, then `http.Server.Shutdown()`.
- Result: PASS

### Worker
- Documented: "Upon receiving a shutdown signal, it stops pulling new jobs but continues processing the current active job until completion."
- Code: Partially implemented. `Stop()` closes channel and fires `cancel()`. The worker loop `select` listens to `ctx.Done()` which causes it to return immediately when finished with the active job, abandoning buffered jobs already accepted into `w.jobChan`.
- Result: TEST_CLAIM_MISMATCH / BROKEN_IMPLEMENTATION. The claim says "stops pulling new jobs", but technically it also drops already-pulled/buffered jobs, which violates standard graceful drain principles.

### Database
- Documented: "Demonstrates the 'Expand and Contract' pattern... transparent fallback logic"
- Code: Fully implemented in `GetUser`.
- Result: PASS

### Demo
- Documented: "A CLI orchestrator that wires these components together... sends a termination signal to demonstrate zero-downtime draining behavior."
- Code: Runs correctly. Emulates orchestrator SIGTERM, waits for drain.
- Result: PASS
