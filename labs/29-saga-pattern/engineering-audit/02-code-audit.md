# Code Audit

## Finding 1

Location: internal/saga/orchestrator.go:49-90
Claimed Behavior: Step failure or cancellation triggers compensating actions in reverse (LIFO) order for executed steps.
Observed Implementation: Steps execute sequentially; completed steps are tracked in `executed` slice; on error or context cancellation, `compensate` iterates backwards through `executed` slice executing `Compensate` callbacks. Errors during compensation are aggregated and returned.
Assessment: PASS
Severity: LOW
Notes: Compensation executes against `context.Background()` to ensure cleanup runs even if parent context timed out or was cancelled.

## Finding 2

Location: internal/saga/orchestrator.go:92-112
Claimed Behavior: Thread-safe compensation logging and status updates.
Observed Implementation: Mutex `o.mu` protects updates to `o.logs` slice during forward execution and compensation.
Assessment: PASS
Severity: LOW
Notes: No race conditions observed when logging compensation states.

## Finding 3

Location: internal/saga/choreography.go:27-52
Claimed Behavior: Decoupled event bus for choreography model with safe concurrent publication and subscription.
Observed Implementation: `EventBus` uses `sync.RWMutex`. `Publish` copies slice of handlers under read lock, preventing data races if subscriptions happen concurrently.
Assessment: PASS
Severity: LOW
Notes: Handlers are called synchronously within `Publish`, preserving event ordering in single-threaded tests.

## Finding 4

Location: internal/services/services.go:82-97
Claimed Behavior: Idempotent payment processing using idempotency keys.
Observed Implementation: `ProcessPayment` checks `processedID` map with mutex protection and returns nil if transaction ID was previously processed.
Assessment: PASS
Severity: LOW
Notes: Idempotency keys are kept per session.

## Finding 5

Location: internal/services/services.go:29-62
Claimed Behavior: Semantic locking countermeasure preventing concurrent order updates.
Observed Implementation: `OrderService` maintains `locks` map. `CreateOrder` checks lock existence before setting status to `OrderPending` and establishing lock. `ApproveOrder` and `CancelOrder` release lock via `delete(s.locks, orderID)`.
Assessment: PASS
Severity: LOW
Notes: Correct semantic lock enforcement preventing concurrent operation on locked order.
