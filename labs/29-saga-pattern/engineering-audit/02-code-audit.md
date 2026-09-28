# Code Audit

## Finding 1

Location: `internal/saga/orchestrator.go:49-90`
Claimed Behavior: Forward steps executed sequentially, failing step halts execution and initiates rollback.
Observed Implementation: `Execute` iterates over copied step slice. Checked context cancellation per iteration. On error, marks status `FAILED` and calls `o.compensate(context.Background(), executed)`.
Assessment: PASS
Severity: LOW
Notes: Correctly passes un-cancelled `context.Background()` to `compensate` so rollback transactions are not aborted by the forward request's cancelled context.

## Finding 2

Location: `internal/saga/orchestrator.go:92-112`
Claimed Behavior: Compensations executed in strict reverse order (LIFO), errors logged and aggregated.
Observed Implementation: Loops `for i := len(executed) - 1; i >= 0; i--`. Appends logs as `COMPENSATED` or `COMPENSATE_FAILED`. Collects and returns aggregated `compErrors`.
Assessment: PASS
Severity: LOW
Notes: LIFO reversal is exact. Status logging is protected by mutex.

## Finding 3

Location: `internal/saga/choreography.go:27-52`
Claimed Behavior: Thread-safe pub/sub event bus supporting multiple decoupled event handlers.
Observed Implementation: Protected with `sync.RWMutex`. `Publish` copies handler slice under read-lock before invocation to avoid deadlock during handler executions.
Assessment: PASS
Severity: LOW
Notes: Concurrency safe.

## Finding 4

Location: `internal/services/services.go:16-67`
Claimed Behavior: Semantic locking mechanism on order creation until approved/cancelled.
Observed Implementation: Mutex guarded map `locks[orderID]`. Returns error if order is already locked. Unlocks on `ApproveOrder` and `CancelOrder`.
Assessment: PASS
Severity: LOW
Notes: Prevents dirty overwrites during saga execution.

## Finding 5

Location: `internal/services/services.go:69-112`
Claimed Behavior: Idempotent payment processing using idempotency keys.
Observed Implementation: `processedID[paymentID]` checked under lock; returns `nil` early if already processed.
Assessment: PASS
Severity: LOW
Notes: Safe against duplicate retries.
