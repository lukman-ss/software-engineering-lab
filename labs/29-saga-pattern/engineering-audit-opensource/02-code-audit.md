# Code Audit

## Finding 1

Location: `internal/saga/orchestrator.go:49-90`
Claimed Behavior: Forward steps execute sequentially; upon failure or context cancellation, executed steps roll back in LIFO order using `context.Background()` to ensure cleanup completes even if original context was cancelled.
Observed Implementation: `Execute` records executed steps and calls `compensate(context.Background(), executed)` in reverse order upon step failure or context cancellation.
Assessment: PASS
Severity: LOW
Notes: Properly isolates compensation context from parent cancellation so cleanup is not aborted midway.

## Finding 2

Location: `internal/saga/choreography.go:27-52`
Claimed Behavior: Decoupled event bus allows pub/sub event-driven saga choreography.
Observed Implementation: Thread-safe `EventBus` using `sync.RWMutex` dispatches published events to registered handlers synchronously.
Assessment: PASS
Severity: LOW
Notes: Clean in-memory pub/sub pattern implementation suitable for saga choreography demonstration.

## Finding 3

Location: `internal/services/services.go:16-61`
Claimed Behavior: Semantic locks prevent concurrent modification of active pending orders.
Observed Implementation: `OrderService` maintains an in-memory lock table `locks[orderID] = true` during pending state, rejecting duplicate order creation or conflicting actions until approved or cancelled.
Assessment: PASS
Severity: LOW
Notes: Accurately models the semantic locking countermeasure described in Saga research.

## Finding 4

Location: `internal/services/services.go:69-113`
Claimed Behavior: Payment processing is idempotent using transaction key tracking.
Observed Implementation: `PaymentService` maintains `processedID[paymentID]`, returning immediate success on duplicate execution requests.
Assessment: PASS
Severity: LOW
Notes: Prevents double-charging on network retry / saga step re-execution.
