# Diagrams

## State Machine Diagram

```text
        ┌───────────┐
        │  CLOSED   │◄──────────────────────────┐
        └─────┬─────┘                           │
              │ failures >= FailureThreshold    │ probe success
              ▼                                 │
        ┌───────────┐                           │
        │   OPEN    │                           │
        └─────┬─────┘                           │
              │ cooldown elapsed                │
              ▼                                 │
        ┌───────────┐                           │
        │ HALF_OPEN ├─────────────────────────┘
        └─────┬─────┘
              │ probe fail
              ▼
        (to OPEN)
```

## Architecture Diagram

```text
Checkout Service
       │
       ▼
Circuit Breaker Proxy
  ├── [CLOSED]   ──► Call Payment Client ──► Fake Payment Server
  ├── [OPEN]     ──► Fail Fast (ErrCircuitOpen)
  └── [HALF_OPEN]──► Allow Canary Probe
```

## Without vs With Circuit Breaker — Request Flow

```text
WITHOUT CIRCUIT BREAKER (SLOW DEPENDENCY):
  Request 1  ──► [blocked until timeout (~100ms)] ──► Payment Server
  Request 2  ──► [blocked until timeout (~100ms)] ──► Payment Server
  Request 3  ──► [blocked until timeout (~100ms)] ──► Payment Server
  ↓
  downstream_calls=3
  total_latency = 3 * timeout

WITH CIRCUIT BREAKER (DOWN DEPENDENCY):
  Request 1  ──► [CALL]  Payment Server (error)  state=CLOSED   ~500µs
  Request 2  ──► [CALL]  Payment Server (error)  state=CLOSED   ~170µs
  Request 3  ──► [CALL]  Payment Server (error)  tripped          ~170µs
  Request 4  ──► [FAIL FAST]  ErrCircuitOpen      state=OPEN      ~700ns
  Request 5  ──► [FAIL FAST]  ErrCircuitOpen      state=OPEN      ~500ns
  Request 6  ──► [FAIL FAST]  ErrCircuitOpen      state=OPEN      ~375ns
  ↓
  downstream_calls=3 (stopped once OPEN)
```

## Recovery Flow

```text
  State = OPEN
     │
     │ wait OpenTimeout (300ms)
     │
     ▼
  State = HALF_OPEN
     │
     │ send probe request
     │
     ├─ probe success ──► State = CLOSED (resume normal traffic)
     └─ probe fail ────► State = OPEN (restart cooldown)
```

## Concurrency Safety

```text
  Goroutine 1        Goroutine 2        Goroutine N
       │                    │                  │
       ├─ cb.Execute(fn) ──►│                  │
       │  ├─ cb.mu.Lock() ──►│                  │
       │  ├─ state check    │                  │
       │  ├─ cb.mu.Unlock() │                  │
       │  └─ fn()           │                  │
       │                    ├─ cb.Execute() ──►│
       │                    │  ├─ cb.mu.Lock() │
       │                    │  └─ blocked...   │
       │                    │                  └─ ...
       └─ return            └─ return          └─ return
```

All state mutations and checks are guarded by `sync.Mutex`, verified by Go race detector.