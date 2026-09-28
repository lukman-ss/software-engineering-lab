## Deadline Propagation Hierarchy

```
Client Request (deadline: 500ms)
        │
        ▼
ExecuteWithBudget(budget: 300ms) ──▶ childCtx deadline = min(500ms, 300ms) = 300ms
        │
        ▼
Downstream Service Call
        │
        ▼
Jika childCtx.Done() terjadi karena:
  - parent timeout (500ms tercapai) → context.DeadlineExceeded
  - budget lokal timeout (300ms tercapai) → context.DeadlineExceeded
  - function selesai dengan error/nil → done channel
```

## Circuit Breaker State Transition Diagram

```
State Transitions:

[CLOSED] 
   │ FailureThreshold reached (consecutive failures)
   ▼
[OPEN] ──▶ Requests blocked immediately (ErrCircuitOpen)
   │ Cooldown elapsed
   ▼
[HALF_OPEN] 
   ├─ SuccessThreshold reached (consecutive successes) ───▶ [CLOSED]
   │
   └─ Any failure ───────────────────────────────────────▶ [OPEN]
```

## Retry with Jitter Distribution

```
Attempt 1: sleep ∈ [0, base] (0-10ms)
Attempt 2: sleep ∈ [0, base×2] (0-20ms)
Attempt 3: sleep ∈ [0, base×4] (0-40ms) — capped at MaxBackoff
Attempt 4+: same as attempt 3 if maxed

Uniform distribution prevents clients from clustering retry waktu yang sama.
```

## Idempotency Store Lifecycle

```
Key: "req-tx-99231"
    │
    ▼
SET: Response="Charged $100", CreatedAt=T0
    │
    ▼ (within TTL)
GET: Found, return "Charged $100"
    │
    ▼ (after TTL expires)
GET: Not found (key deleted on access), return ""
    │
    ▼
SET: new response
```