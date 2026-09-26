# 03 Core Concepts — Circuit Breaker Pattern

## Definition
Circuit breaker wraps a protected function call; monitors failures; trips at threshold; subsequent calls fail without invoking downstream.

**Evidence**: "You wrap a protected function call in a circuit breaker object, which monitors for failures. Once the failures reach a certain threshold, the circuit breaker trips, and all further calls to the circuit breaker return with an error, without the protected call being made at all." — Martin Fowler, Circuit Breaker (2014-03-06) — https://martinfowler.com/bliki/CircuitBreaker.html — Tier 2 — Confidence HIGH — Corroborated by Azure Circuit Breaker pattern (2025-02-05) and Hystrix Wiki (2017-07-03).

## Purpose
- Fail fast (microsecond rejection vs second-long timeout)
- Prevent resource exhaustion cascade
- Give downstream recovery window
- Enable fallback / graceful degradation
- Provide monitoring signal (state change alerts)

**Sources**: Azure Circuit Breaker (2025-02-05) — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker ; Hystrix Wiki — https://github.com/Netflix/Hystrix/wiki/How-it-Works — "Fail fast and rapidly recover. Fallback and gracefully degrade when possible. Enable near real-time monitoring, alerting, and operational control."

## Analogy
Electrical breaker: over-current trips physical switch to protect wiring; manual reset. Software breaker: over-failure trips logical switch; auto-reset via timed probe.

## States Overview
See 05-circuit-states.md for full transitions. Summary:
- CLOSED: pass-through, count failures
- OPEN: short-circuit, return ErrCircuitOpen, no downstream call
- HALF_OPEN: limited probes to test recovery

## Configuration Dimensions (Illustrative, Not Prescriptive)
No universal recommended values exist. Production thresholds are tuned per dependency SLO.

| Parameter | Azure / Hystrix / gobreaker Name | Role |
|---|---|---|
| FailureThreshold / RequestVolumeThreshold + ErrorThresholdPercentage | consec. failures (gobreaker default 5) vs rolling window % | when to trip |
| Timeout / SleepWindow | open → half-open delay (Hystrix default 5s, gobreaker 60s) | cooldown |
| MaxRequests (half-open) | probes allowed (gobreaker 1, Hystrix 1) | thundering-herd control |
| Interval | counts window reset (gobreaker) | prevent stale counts |

**Sources**: gobreaker Settings docs — https://github.com/sony/gobreaker — MaxRequests, Interval, Timeout, ReadyToTrip; Hystrix Wiki Circuit Breaker section — sleepWindowInMilliseconds; Azure — failure threshold within given time period.

## What Breaker Does NOT Do
- Does not heal downstream
- Does not replace timeout (needs timeout to detect slowness)
- Does not replace retry (complementary; retry should respect CB open)
- Does not isolate resource pools by itself (needs bulkhead)

**Source**: Azure Circuit Breaker — "The Circuit Breaker pattern serves a different purpose than the Retry pattern... An application can combine these two patterns"; Google SRE Handling Overload — retry budgets separate from overload protection.
