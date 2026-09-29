# Code Audit

Target Lab: `labs/29-saga-pattern`

## Finding 1: Core Orchestrator Sequential Execution and Thread Safety

Location: `internal/saga/orchestrator.go:49-90`
Claimed Behavior: Orchestrator executes steps sequentially, records logs, and handles context cancellation/failures safely.
Observed Implementation:
- Uses `sync.Mutex` to safely read `steps` into a local slice before iteration.
- Checks `ctx.Done()` before executing each step.
- Logs step state (`StatusExecuted` or `StatusFailed`) under mutex protection.
- Collects executed steps and passes to `compensate()` upon failure or cancellation.
Assessment: PASS
Severity: LOW
Notes: Thread-safe, protects logs and step execution order.

---

## Finding 2: LIFO Compensation and Error Propagation

Location: `internal/saga/orchestrator.go:92-112`
Claimed Behavior: Compensations execute in reverse (LIFO) order of execution, and compensation errors are recorded and returned.
Observed Implementation:
- Iterates backwards `for i := len(executed) - 1; i >= 0; i--`.
- Executes `step.Compensate(ctx)`.
- Updates `o.logs` with `StatusCompensated` or `StatusCompensateFailed`.
- Returns combined errors if any compensation fails.
Assessment: PASS
Severity: LOW
Notes: Accurately implements LIFO rollback. As noted in implementation notes, compensations assume no automatic retry/exponential backoff.

---

## Finding 3: Semantic Lock Countermeasure

Location: `internal/services/services.go:29-40, 42-61`
Claimed Behavior: Prevent dirty reads / concurrent modifications during saga execution using a semantic lock on pending orders.
Observed Implementation:
- `CreateOrder` sets `s.locks[orderID] = true` and checks if already locked.
- `ApproveOrder` and `CancelOrder` release lock via `delete(s.locks, orderID)`.
- All access guarded by `sync.Mutex`.
Assessment: PASS
Severity: LOW
Notes: Cleanly demonstrates semantic locking countermeasure for ACD (Atomic, Consistent, Durable) saga properties without ACID isolation.

---

## Finding 4: Idempotency Key Enforcement

Location: `internal/services/services.go:82-97`
Claimed Behavior: Payment processing checks idempotency key to prevent double charging on retry.
Observed Implementation:
- Maintains `processedID map[string]bool`.
- Returns `nil` immediately if key already exists.
- Guarded by `sync.Mutex`.
Assessment: PASS
Severity: LOW
Notes: Idempotency key pattern correctly implemented.

---

## Finding 5: Choreography Event Bus

Location: `internal/saga/choreography.go:27-52`
Claimed Behavior: Thread-safe event bus supporting pub/sub event choreography across decoupled services.
Observed Implementation:
- Uses `sync.RWMutex` to protect subscriber map.
- `Publish` copies handler slice under `RLock` and invokes handlers sequentially outside lock to avoid deadlock.
Assessment: PASS
Severity: LOW
Notes: Idiomatic standard library Go pub/sub implementation.
