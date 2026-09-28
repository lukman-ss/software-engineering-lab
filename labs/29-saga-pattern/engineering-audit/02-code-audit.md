# Code Audit

## Finding 1

Location: `internal/saga/orchestrator.go:49-90`
Claimed Behavior: Forward step execution, LIFO compensation on step error or context cancellation, and thread-safe step/log management.
Observed Implementation:
- Uses `sync.Mutex` to safely read step slice and update logs.
- Evaluates `ctx.Done()` before executing steps and triggers LIFO compensation with `context.Background()`.
- On step execution error, immediately executes `o.compensate(context.Background(), executed)`.
Assessment: PASS
Severity: LOW
Notes: `context.Background()` is intentionally passed to `compensate` so compensation is guaranteed to complete even if the parent `ctx` was cancelled.

## Finding 2

Location: `internal/saga/orchestrator.go:92-112`
Claimed Behavior: LIFO compensation execution and error aggregation.
Observed Implementation:
- Iterates backwards over `executed` steps (`len(executed)-1` down to 0).
- Records `StatusCompensated` or `StatusCompensateFailed`.
- Aggregates compensation errors into a single error if any compensation step fails.
Assessment: PASS
Severity: LOW
Notes: Code clean and correct.

## Finding 3

Location: `internal/saga/choreography.go:27-52`
Claimed Behavior: Decoupled pub/sub event bus supporting choreography.
Observed Implementation:
- Thread-safe handlers map using `sync.RWMutex`.
- Copying handler slice under `RLock` before dispatching to prevent lock contention during handler invocation.
Assessment: PASS
Severity: LOW
Notes: No blocking issues observed.

## Finding 4

Location: `internal/services/services.go:29-40, 82-97, 127-152`
Claimed Behavior: Domain services with semantic locking, payment idempotency, and inventory reservation.
Observed Implementation:
- `OrderService` locks order state on creation (`locks[orderID] = true`), checked during `CreateOrder`.
- `PaymentService` checks idempotency map (`processedID[paymentID]`) to prevent double charges.
- `InventoryService` checks and updates stock/reserved balances safely under `sync.Mutex`.
Assessment: PASS
Severity: LOW
Notes: Clean thread-safe stdlib-only implementation.
