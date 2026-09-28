# Code Audit

## Finding 1
Location: internal/fault/injector.go
Claimed Behavior: Inject configurable latency and forced errors, clear on reset, safe for concurrent access.
Observed Implementation: SetFault/Clear/IsEnabled guarded by RWMutex; Execute takes a snapshot under RLock before blocking. Latency uses select racing time.After vs ctx.Done.
Assessment: PASS
Severity: LOW
Notes: Error rate/probability field declared but unused; not a defect, just dead config reserved for future use.

## Finding 2
Location: internal/circuitbreaker/circuitbreaker.go
Claimed Behavior: Closed -> Open on threshold failures; OPEN fast-fails with ErrCircuitOpen or fallback; cooldown moves to HALF-OPEN; success returns to CLOSED.
Observed Implementation: Execute releases mutex before invoking user fn and fallback so callbacks never deadlock on reentrancy. failure counter increments and transitions OPEN. Success in HALF-OPEN resets state; success in CLOSED resets failures.
Assessment: PASS
Severity: LOW
Notes: No timer reset race; lastStateChg updated atomically under mutex.

## Finding 3
Location: internal/monitor/monitor.go
Claimed Behavior: Tracks total/failed counters and computes error rate; IsHealthy enforces threshold after minimum sample.
Observed Implementation: Counters are atomic; snapshot in Metrics() reads atomically. IsHealthy returns true before 5 samples and otherwise compares rate to maxErrorRate.
Assessment: PASS
Severity: LOW
Notes: Cumulative (not sliding-window) counters, explicitly documented as a design simplification; matches engineering notes.

## Finding 4
Location: internal/experiment/runner.go
Claimed Behavior: Inject fault at start, periodically check steady state, auto-abort + neutralize injector on breach, else complete on timeout.
Observed Implementation: Run sets RUNNING, calls injector.SetFault, then select loop on ctx/timeout/ticker. terminate() locks, clears injector, sets state. Abort path returns formatted error.
Assessment: PASS
Severity: MEDIUM
Notes: Run is not idempotent; calling twice concurrently is unprotected, but demo/test call it once per instance. Within documented single-call usage this is acceptable.

## Finding 5
Location: cmd/demo/main.go
Claimed Behavior: Baseline traffic, mitigated chaos (circuit breaker + fallback), unmitigated fault triggering auto-abort, then recovery.
Observed Implementation: Uses shared injector/monitor/circuitbreaker. Launches experiment in goroutine and drives client requests. Reports experiment state, abort reason, and injector neutralization.
Assessment: PASS
Severity: LOW
Notes: No explicit synchronization between exp.Run goroutine and final exp.State() read, but State() is lock-protected and Run completes before final reads in practice due to sleeps. Acceptable for demo.
