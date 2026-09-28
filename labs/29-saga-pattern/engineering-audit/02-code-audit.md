# Code Audit

## Finding 1

Location: `labs/29-saga-pattern/internal/saga/orchestrator.go:78`
Claimed Behavior: Step failure triggers compensating actions in reverse order (LIFO) for completed steps.
Observed Implementation: `Execute` creates a slice of executed steps and on step failure invokes `o.compensate(context.Background(), executed)`, which iterates backwards from `len(executed)-1` down to 0, invoking `Compensate`.
Assessment: PASS
Severity: LOW
Notes: Compensation receives a background context to ensure rollback isn't prematurely aborted by an expired parent context.

## Finding 2

Location: `labs/29-saga-pattern/internal/saga/orchestrator.go:58-69`
Claimed Behavior: Context cancellation stops step execution and triggers compensation.
Observed Implementation: Each step loop checks `select { case <-ctx.Done(): ... }`. If cancelled, it logs `StatusFailed` for the pending step and executes compensations for already executed steps.
Assessment: PASS
Severity: LOW
Notes: Correctly handles context deadline/cancellation without leaving completed steps uncompensated.

## Finding 3

Location: `labs/29-saga-pattern/internal/services/services.go:33-39, 42-51`
Claimed Behavior: Semantic lock countermeasure prevents concurrent sagas from modifying pending state.
Observed Implementation: `CreateOrder` sets `locks[orderID] = true` and `orders[orderID] = OrderPending`. Subsequent `CreateOrder` calls for the same locked `orderID` fail with `"order is semantically locked"`. `ApproveOrder` and `CancelOrder` release the lock by deleting `locks[orderID]`.
Assessment: PASS
Severity: LOW
Notes: Properly enforces semantic lock state transitions guarded by `sync.Mutex`.

## Finding 4

Location: `labs/29-saga-pattern/internal/services/services.go:86-88`
Claimed Behavior: Idempotency keys prevent double-charging/duplicate execution.
Observed Implementation: `PaymentService.ProcessPayment` tracks `processedID[paymentID]`. If `paymentID` is already marked true, it returns `nil` immediately without re-charging.
Assessment: PASS
Severity: LOW
Notes: Idempotency check occurs under mutex lock before state mutation.

## Finding 5

Location: `labs/29-saga-pattern/internal/saga/choreography.go:38-51`
Claimed Behavior: Choreography decoupled event bus handles event subscription and publishing concurrently safely.
Observed Implementation: `EventBus` uses `sync.RWMutex`. `Subscribe` acquires write lock, `Publish` acquires read lock, makes a copy of event handlers, releases read lock, and invokes handlers synchronously.
Assessment: PASS
Severity: LOW
Notes: Synchronous handler execution in `Publish` avoids race conditions on handler slices while keeping invocation simple for in-memory demonstration.
