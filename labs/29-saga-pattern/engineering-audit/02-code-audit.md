# Code Audit Findings

## Finding 1

Location: `internal/saga/orchestrator.go:49-90`
Claimed Behavior: Forward steps execute sequentially; upon any step failure or context cancellation, executed steps are rolled back in strict LIFO order.
Observed Implementation: `Execute()` captures executed steps in a slice. If `step.Execute(ctx)` returns an error or `ctx.Done()` fires, `o.compensate(context.Background(), executed)` is invoked which iterates in reverse order `for i := len(executed) - 1; i >= 0; i--`.
Assessment: PASS
Severity: LOW
Notes: Compensation runs on `context.Background()` to ensure rollback actions complete even if the parent saga context was cancelled.

## Finding 2

Location: `internal/saga/orchestrator.go:92-112`
Claimed Behavior: Compensation failure is captured, tracked in step logs as `StatusCompensateFailed`, and returned in the aggregated error.
Observed Implementation: When `step.Compensate(ctx)` returns an error, `o.logs` appends `StatusCompensateFailed`, and errors are accumulated into `compErrors` slice returned to caller.
Assessment: PASS
Severity: LOW
Notes: Robust compensation failure reporting aligned with research finding on compensation non-guarantees.

## Finding 3

Location: `internal/saga/choreography.go:27-52`
Claimed Behavior: Decoupled event publish/subscribe bus enables choreographed distributed transaction steps.
Observed Implementation: `EventBus` stores slice of handlers per `EventType`, protected with `sync.RWMutex`. Safe handler slice copying under read lock prevents race conditions during publication.
Assessment: PASS
Severity: LOW
Notes: Implementation cleanly captures the choreography pattern without unnecessary framework bloat.

## Finding 4

Location: `internal/services/services.go:16-68`
Claimed Behavior: Semantic locking countermeasure on `OrderService` prevents dirty writes/concurrent modifications on pending orders.
Observed Implementation: `CreateOrder()` checks `s.locks[orderID]` and returns error if already locked. Locks are released upon `ApproveOrder` or `CancelOrder`. State mutations protected by `sync.Mutex`.
Assessment: PASS
Severity: LOW
Notes: Valid implementation of semantic lock isolation countermeasure from research finding 7.

## Finding 5

Location: `internal/services/services.go:69-113`
Claimed Behavior: `PaymentService` enforces idempotency based on payment ID.
Observed Implementation: `processedID` map tracks already processed payment identifiers; duplicate calls return `nil` without duplicate balance deduction.
Assessment: PASS
Severity: LOW
Notes: Correctly fulfills research finding 8 regarding idempotent retry handling.

## Finding 6

Location: `cmd/demo/main.go:1-113`
Claimed Behavior: Runnable demonstration of orchestrator happy path and failure rollback scenarios.
Observed Implementation: Demo runs Scenario 1 (full success flow) and Scenario 2 (stock depletion rollback), producing exact matching output.
Assessment: PASS
Severity: LOW
Notes: No mock/fake outputs; runs live code paths.
