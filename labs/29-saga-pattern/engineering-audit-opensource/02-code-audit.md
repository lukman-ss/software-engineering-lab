## Code Audit Findings

### Finding 1: Compensation error handling ignored
Location: internal/saga/orchestrator.go:79
Claimed Behavior: Compensation actions are invoked for rollback and errors are logged or propagated.
Observed Implementation: The return value of step.Compensate(ctx) is discarded with `_ =`.
Assessment: WARNING
Severity: MEDIUM
Notes: If a compensation fails, the system may be left in an inconsistent state, and the error is not surfaced to the caller. The orchestrator should at least log compensation errors or aggregate them.

### Finding 2: Missing context cancellation check during step execution
Location: internal/saga/orchestrator.go:56-70
Claimed Behavior: The orchestrator respects context cancellation.
Observed Implementation: The step.Execute(ctx) is called, but the orchestrator does not check ctx.Done() before invoking the next step or during iteration.
Assessment: WARNING
Severity: LOW
Notes: If the context is cancelled while a step is executing, the orchestrator will continue to the next step after the current step returns (or fails). This may lead to unnecessary work after cancellation.

### Finding 3: EventBus handler modification during publish
Location: internal/saga/choreography.go:44-51
Claimed Behavior: EventBus is safe for concurrent subscription and publishing.
Observed Implementation: Publish takes a read lock, copies the handler slice, then releases the lock before invoking handlers. Handlers can subscribe/unsubscribe (which take a write lock) concurrently.
Assessment: PASS
Severity: LOW
Notes: The pattern is safe because the handler slice is a copy; modifications to the underlying map do not affect the slice. However, if a handler blocks for a long time, new subscriptions will wait for the write lock.

### Finding 4: Semantic lock levers
Location: internal/services/services.go:33-38
Claimed Behavior: Semantic lock prevents concurrent creation of the same order.
Observed Implementation: Lock is checked and set within a mutex. Released on approval or cancellation.
Assessment: PASS
Severity: LOW
Notes: Correct implementation. The lock is cleared only on terminal states (approved/cancelled). If an order creation succeeds but the saga fails later, the compensating CancelOrder will release the lock.

### Finding 5: Idempotency key verification
Location: internal/services/services.go:86-88
Claimed Behavior: Payment processing is idempotent using processedID map.
Observed Implementation: If processedID[paymentID] is true, returns nil without verifying that the amount matches the original.
Assessment: PASS
Severity: LOW
Notes: This is acceptable for demonstration; in a real system, the amount should be checked to prevent mismatch.