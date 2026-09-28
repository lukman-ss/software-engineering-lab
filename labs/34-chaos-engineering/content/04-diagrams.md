# Architectural & Behavioral Diagrams

## 1. Circuit Breaker State Transition Diagram

Diagram transisi state pada `internal/circuitbreaker/circuitbreaker.go`:

```
                 [Normal Operation]
                         │
                         ▼
                ┌──────────────────┐
        ┌──────>│      CLOSED      │<────────────────────────┐
        │       └──────────────────┘                         │
        │                 │                                  │
Success │                 │ failures >= threshold            │ Success
resets  │                 ▼                                  │ resets
failure │       ┌──────────────────┐                         │ failure
count   │       │       OPEN       │                         │ count
        │       └──────────────────┘                         │
        │                 │                                  │
        │                 │ cooldown elapsed                 │
        │                 ▼                                  │
        │       ┌──────────────────┐                         │
        └───────┤    HALF-OPEN     ├─────────────────────────┘
                └──────────────────┘
                          │
                          │ failure detected
                          ▼
                     (Back to OPEN)
```

---

## 2. Chaos Experiment Lifecycle & Auto-Abort Architecture

Diagram siklus hidup eksperimen dan penahanan blast radius pada `internal/experiment/runner.go`:

```
             ┌────────────────────────┐
             │    Experiment Init     │
             │     State: PENDING     │
             └───────────┬────────────┘
                         │
                         ▼
             ┌────────────────────────┐
             │    SetFault() Active   │
             │     State: RUNNING     │
             └───────────┬────────────┘
                         │
          ┌──────────────┴──────────────┐
          ▼                             ▼
    [Wait Duration]           [Monitor Ticker (Every T)]
          │                             │
          │ Duration                    │
          │ Elapsed                     ▼
          │                   IsHealthy() == false?
          │                   (ErrorRate > MaxThreshold)
          │                             │
          │                   ┌─────────┴─────────┐
          │               Yes │                No │
          │                   ▼                   ▼
          │          ┌─────────────────┐    (Keep Running)
          │          │  terminate()    │
          │          │  State: ABORTED │
          │          │  injector.Clear │
          │          └─────────────────┘
          ▼
   ┌─────────────────┐
   │   terminate()   │
   │ State: COMPLETED│
   │ injector.Clear  │
   └─────────────────┘
```

---

## 3. Resilient Request Flow with Graceful Degradation

Alur eksekusi request melalui client berketahanan (`cmd/demo/main.go`):

```
Client Call
    │
    ▼
CircuitBreaker.Execute()
    │
    ├─► State == OPEN? ──Yes─► Execute Fallback() ──► Return "Payment Queued (Fallback)"
    │                                                      │
    │ No                                                   ▼
    ▼                                              RecordSuccess() to Monitor
Call Downstream Service (with FaultInjector)
    │
    ├──► Downstream OK? ──Yes─► Return Success ──► RecordSuccess()
    │
    └──► Downstream Error? ──Yes─► Increment Failures
                                          │
                                          ├─► Failures >= Threshold? ──Yes─► State = OPEN
                                          │
                                          └─► Fallback Provided?
                                                    │
                                                    ├─► Yes ─► Execute Fallback() ──► RecordSuccess()
                                                    └─► No  ─► Return Error       ──► RecordFailure()
```
