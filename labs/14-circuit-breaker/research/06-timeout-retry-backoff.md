# 06 Timeout, Retry, and Exponential Backoff

## Timeout
- **Definition**: Absolute upper bound on outbound call duration.
- **Why mandatory**: Without timeout, slow downstream hangs caller indefinitely → resource exhaustion → cascade.
- **Go implementation**: `http.Client{Timeout: ...}` — standard library, no external deps.

**Sources**: Azure Circuit Breaker — "Inappropriate time-outs on external services... If the time-out is too long, a thread that runs a circuit breaker might be blocked for an extended period"; Google SRE Workbook — "Processes hold resources for all in-flight requests... By default, this deadline is a very large number... causes clients... to experience higher latency... at risk of running out of resources".

## Retry
- **Purpose**: Handle transient failures (network blip, GC pause).
- **Strategies (Azure Retry Pattern)**: Cancel, Immediate Retry, Retry After Delay (exponential increment).
- **Idempotency prerequisite**: Must guarantee safe re-execution (GET usually OK; POST needs idempotency key).
- **Bounded**: Max attempts (typically 2-3); total retry budget.

**Sources**: Azure Retry Pattern (2024-07-18) — "retry logic should be sensitive to any exceptions that the circuit breaker returns and stop retry attempts if the circuit breaker indicates that a fault isn't transient"; AWS Exponential Backoff and Jitter (2015-03-04).

## Exponential Backoff + Jitter
- Plain exponential backoff clusters retries → synchronized spikes → thundering herd.
- Jitter decorrelates retry times → ~constant arrival rate.
- **Full Jitter** (AWS recommended): `sleep = random(0, min(cap, base * 2^attempt))`
- **Equal Jitter**: `sleep = (base * 2^attempt) / 2 + random(0, (base * 2^attempt) / 2)`
- **Decorrelated Jitter**: `sleep = random(0, min(cap, last_sleep * multiplier))`

**Evidence**: AWS Architecture Blog (Marc Brooker, 2015) — "no-jitter exponential backoff is the clear loser... Full Jitter uses less work... Decorrelated slightly more work but slightly less time"; simulations at 100 contending clients show >50% call-count reduction with jitter.

## Interaction: Timeout → Retry → Circuit Breaker
```
Attempt 1
  └─ timeout bound
Attempt 2 (backoff + jitter)
  └─ if still failing
Attempt 3
  └─ if still failing → circuit threshold met → CB OPEN
     └─ subsequent calls fail fast without downstream + retry
```
- Retry must respect `ErrCircuitOpen` (stop retrying immediately).
- CB prevents retry storm; timeout bounds single attempt.
- Together: layered defense.

**Sources**: Azure Circuit Breaker — "retry logic should be sensitive to any exceptions that the circuit breaker returns and stop retry attempts if the circuit breaker indicates that a fault isn't transient"; Google SRE Handling Overload — per-request retry budget (max 3) + per-client retry ratio cap (10%).

## Fallback (vs Retry vs CB)
| Concern | Pattern |
|---|---|
| "Dependency slow" | Timeout |
| "Transient blip" | Retry (bounded, jittered) |
| "Dependency persistently failing" | Circuit Breaker (fail fast) |
| "User-visible graceful degradation" | Fallback (cached, degraded, queued) |

**Note**: Do not use fallback for financial writes (PPOB); queue async instead.