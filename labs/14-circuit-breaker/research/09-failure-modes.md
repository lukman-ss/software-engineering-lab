# 09 Failure Modes

## 1. Premature Opening
**Definition**: Threshold too low or evaluation window too short; transient blips flip circuit OPEN unnecessarily.
**Effect**: False rejection of healthy traffic; reduced availability.
**Mitigation**: tune threshold to observed baseline error rate; use rolling-window error % rather than consecutive count if traffic spiky; min sample volume (Hystrix `requestVolumeThreshold`).
**Sources**: Azure Circuit Breaker Problems — "If the circuit breaker remains in the Open state for a long period, it can raise exceptions even if the reason for the failure is resolved" is *late* opening; early opening is inverse.

## 2. Thundering Herd on HALF_OPEN Probe
**Definition**: Too many concurrent probes overwhelm recovering dependency, causing probe failures → immediate re-trip.
**Effect**: prolongs outage; delays recovery.
**Mitigation**: limit `HalfOpenMaxCalls` (gobreaker default 1; Hystrix 1; Azure "limited number"); stagger probe start across instances (not implemented in-memory breaker; requires external coordination via external store or sidecar).

**Sources**: Azure Circuit Breaker — "The Half-Open state helps prevent a recovering service from suddenly being flooded with requests"; Google SRE Workbook Managing Load — Dressy case study where load shedding + load balancing misconfigured isolated caused imbalance; "add error handling to load balancer logic" principle.

## 3. Infinite Open (Stuck Open)
**Definition**: cooldown too long or dependency never passes health check; circuit never leaves OPEN.
**Effect**: permanent fail-fast until manual override; increased MTTR.
**Mitigation**: reasonable `OpenTimeout` (seconds-minutes not hours); manual override button; health-check endpoint (separate from circuit breaker) to trigger reset; adaptive timeout increase only after repeated failures.

**Sources**: Hystrix Wiki — `circuitBreakerSleepWindowInMilliseconds()`; Azure Circuit Breaker Problems — "provide a manual reset option... force a circuit breaker into the Open state".

## 4. 4xx False Positives (Client Errors Tripping Breaker)
**Definition**: Counting HTTP 4xx (e.g. 400 Bad Request, 404 Not Found, 422 Unprocessable Entity) as circuit-breaker failures.
**Effect**: breaker opens on user-input validation errors; legitimate retries of same input always fail; reduces availability incorrectly.
**Mitigation**: configure `IsSuccessful` / `IsExcluded` to exclude 4xx unless service validates that specific 4xx indicates systemic fault (rare). Standard practice: trip on 5xx, timeout, connection error; not on 4xx.

**Evidence**: Azure Circuit Breaker Problems — "4xx False Positives: Tripping the breaker on user-input validation errors (400, 404, 422) instead of system faults (500, 502, 503, 504, timeouts)" — exact match from research/04-contradictions.md.

## 5. Metric Skew (Roll-over / Window Misalignment)
**Definition**: sliding window or interval misalignment causes sudden jumps in reported error rate.
**Effect**: spurious trips or delayed trips.
**Mitigation**: use sufficiently large window; align sampling with traffic periodicity if known; accept small inaccuracy for stability.

**Source**: cep21/circuit rolling stats discussion — windowing strategy impacts metric smoothness.

## 6. Panic Propagation
**Definition**: protected function panics; if circuit breaker does not recover panic, brings down caller.
**Effect**: 100% failure rate for caller until breaker reset.
**Mitigation**: circuit breaker must `recover()` panic and treat as error (not let it propagate). Go: `defer func() { if r:=recover(); r!=nil { /* handle */ } }()`.

**Evidence**: cep21/circuit — "recoverable panic()" feature listed; gobreaker does not mention panic handling explicitly (assumed same pattern).

## 7. Resource Leak in Half-Open
**Definition**: probe allocates resource but fails to release it on error → leak over many cycles → eventual OOM.
**Mitigation**: ensure probe execution path releases resources on any exit (success, error, panic, context cancellation).

**Lab note**: Lab implementations use short-lived func() {} probes; no persistent allocations per probe.

## 8. Incorrect Concurrency Guard
**Definition**: missing mutex on state/counters → race on concurrent Execute → lost updates, spurious state, double-trips.
**Effect**: flapping state, incorrect counts, test flakiness under race detector.
**Fix**: `sync.Mutex` guarding all state transitions and counters (lab and gobreaker use this).

**Sources**: gobreaker source — `mux sync.Mutex` guards `state`, `counts`, `lastStateChangeTime`; cep21/circuit — internal locking not visible in surface API but stress-tested.

## Summary
Failure modes cluster into: tuning (thresholds, timeouts), concurrency (races, leaks), semantics (error classification), and operational (manual overrides, monitoring). Lab avoids: infinite open (reasonable cooldown), 4xx FP (only counts timeouts/connection errors as failures), panic (uses `recover`), races (mutex).