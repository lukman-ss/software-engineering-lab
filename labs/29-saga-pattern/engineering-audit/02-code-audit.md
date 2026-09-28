# Code Audit

## Finding 1

Location: internal/saga/orchestrator.go:78-83
Claimed Behavior: Step failure triggers LIFO compensating transactions and aggregates compensation errors.
Observed Implementation: `Execute` copies steps safely under mutex, iterates forward, logs failure status on error, and executes `compensate` stack in reverse index order. If compensation returns errors, it formats and propagates both the step error and compensation errors.
Assessment: PASS
Severity: LOW
Notes: LIFO compensation execution is correctly implemented and context cancellation is handled with a fallback `context.Background()` for compensation execution.

## Finding 2

Location: internal/saga/orchestrator.go:58-68
Claimed Behavior: Orchestrator checks context cancellation before executing steps and triggers compensation if context is done.
Observed Implementation: Evaluates `ctx.Done()` in a select block prior to executing each step. If cancelled, logs status `StatusFailed` for the unexecuted step, runs compensation on executed steps using `context.Background()`, and returns context error.
Assessment: PASS
Severity: LOW
Notes: Ensures context cancellation does not skip rollback of already executed steps.

## Finding 3

Location: internal/services/services.go:33-35, 18-20
Claimed Behavior: Semantic locking countermeasure prevents concurrent sagas from modifying pending order state.
Observed Implementation: `OrderService` maintains `locks map[string]bool`. `CreateOrder` returns error if `s.locks[orderID]` is true, and sets lock to true. `ApproveOrder` and `CancelOrder` release lock via `delete(s.locks, orderID)`.
Assessment: PASS
Severity: LOW
Notes: Implements semantic locking countermeasure accurately for in-memory model.

## Finding 4

Location: internal/services/services.go:86-88
Claimed Behavior: Idempotency keys prevent duplicate payments on retries.
Observed Implementation: `PaymentService` checks `s.processedID[paymentID]`. If true, returns `nil` without re-processing amount.
Assessment: PASS
Severity: LOW
Notes: Correctly handles idempotent retries for payment processing.

## Finding 5

Location: internal/saga/choreography.go:44-51
Claimed Behavior: Event bus provides decoupled choreography event dispatching.
Observed Implementation: `EventBus` protects handlers slice with `sync.RWMutex`, copies handlers under `RLock()`, and dispatches event handlers synchronously in loop under `Publish`.
Assessment: PASS
Severity: LOW
Notes: Thread-safe in-memory event bus implementation.
