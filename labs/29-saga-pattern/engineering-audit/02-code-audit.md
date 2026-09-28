# Code Audit

## Finding 1

Location: `internal/saga/orchestrator.go:75-85`
Claimed Behavior: LIFO compensation execution upon step failure.
Observed Implementation: Iterates `executed` slice backwards from `len(executed)-1` down to 0, invoking `Compensate(ctx)` if non-nil, logging `StatusCompensated`.
Assessment: PASS
Severity: LOW
Notes: Properly ignores errors during compensation execution as intentionally simplified (annotated with ponytail comment).

## Finding 2

Location: `internal/saga/orchestrator.go:48-73`
Claimed Behavior: Thread-safe forward execution and state logging.
Observed Implementation: Copies `o.steps` under `o.mu` lock before iteration. Mutex locks protecting log updates during step execution and failure handling.
Assessment: PASS
Severity: LOW
Notes: Thread safe and prevents concurrent slice mutation.

## Finding 3

Location: `internal/saga/choreography.go:45-51`
Claimed Behavior: Thread-safe pub/sub event bus for choreography sagas.
Observed Implementation: Copies handler slice under RLock, then executes handlers outside lock to prevent deadlock.
Assessment: PASS
Severity: LOW
Notes: Synchronous handler execution allows deterministic choreography test assertions.

## Finding 4

Location: `internal/services/services.go:29-40`, `42-52`, `54-61`
Claimed Behavior: Semantic locking countermeasure on orders during saga execution.
Observed Implementation: `CreateOrder` sets `s.locks[orderID] = true`. Subsequent `CreateOrder` calls return error. `ApproveOrder` and `CancelOrder` release lock via `delete(s.locks, orderID)`.
Assessment: PASS
Severity: LOW
Notes: Simple and effective in-memory semantic locking implementation.

## Finding 5

Location: `internal/services/services.go:82-97`
Claimed Behavior: Idempotent payment processing.
Observed Implementation: Checks `s.processedID[paymentID]`. If true, returns nil immediately without charging again.
Assessment: PASS
Severity: LOW
Notes: Meets idempotency requirements.

## Finding 6

Location: `internal/services/services.go:127-151`
Claimed Behavior: Inventory reservation and releasing (compensation).
Observed Implementation: `Reserve` decrements stock and increments reserved. `Release` decrements reserved and restores stock, returning error if releasing more than reserved.
Assessment: PASS
Severity: LOW
Notes: Domain logic holds invariants correctly.
