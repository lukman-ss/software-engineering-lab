# 05 Circuit States

```text
       ┌───────────┐
       │  CLOSED   │◄─────────────────────────┐
       └─────┬─────┘                          │
             │ failures >= FailureThreshold   │ probe success
             ▼                                │
       ┌───────────┐                          │
       │   OPEN    │                          │
       └─────┬─────┘                          │
             │ cooldown elapsed               │
             ▼                                │
       ┌───────────┐                          │
       │ HALF_OPEN ├──────────────────────────┘
       └─────┬─────┘
             │ probe fail
             ▼
          (to OPEN)
```

- **CLOSED**: Traffic passes through. Failures increment counter.
- **OPEN**: Traffic blocked immediately with `ErrCircuitOpen`. No network requests dispatched.
- **HALF_OPEN**: Cooldown expired. Canary probe(s) allowed. Success moves state to CLOSED; failure resets breaker to OPEN.
