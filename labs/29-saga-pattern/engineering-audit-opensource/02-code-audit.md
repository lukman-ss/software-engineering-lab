# Code Audit

## Finding 1
Location: internal/saga/orchestrator.go:92-112
Claimed Behavior: LIFO compensation rollback with error aggregation.
Observed Implementation: Compensation iterates over executed steps in reverse order (LIFO). Each compensation is called synchronously. Errors are logged and aggregated, with final return aggregating all compensation errors.
Assessment: PASS
Severity: LOW
Notes: Implementation matches specification; compensation is properly sequenced reverse of execution, ensuring earlier steps are compensated first.

## Finding 2
Location: internal/saga/orchestrator.go:49-90
Claimed Behavior: Sequential step execution with context cancellation support and compensation on failure.
Observed Implementation: Steps are copied from internal slice before execution to avoid mutation. Each step's Execute is called sequentially. Context cancellation triggers immediate compensation of all executed steps.
Assessment: PASS
Severity: LOW
Notes: Behavior matches documented happy path and failure handling; context cancellation is handled correctly, triggering compensation.

## Finding 3
Location: internal/services/services.go:29-40
Claimed Behavior: Semantic lock prevents concurrent sagas from modifying pending state.
Observed Implementation: OrderService maintains locks map[string]bool and checks before creating order. Lock is only removed on Approve or Cancel. However, if saga fails before completing, lock may remain, blocking future operations.
Assessment: WARNING
Severity: MEDIUM
Notes: Potential lock leak; if saga fails after CreateOrder but before Approve/Cancel, lock persists. Need compensating lock cleanup on compensation.

## Finding 4
Location: internal/services/services.go:82-97
Claimed Behavior: Idempotency prevents duplicate payments.
Observed Implementation: PaymentService uses processedID map keyed only by paymentID. If same paymentID is used for different orders, duplicate will be silently ignored, which may be acceptable for idempotency.
Assessment: PASS
Severity: LOW
Notes: Idempotency works; design decision to allow reuse across orders may be intentional.

## Finding 5
Location: internal/saga/choreography.go:44-52
Claimed Behavior: Event bus with subscription and publishing.
Observed Implementation: Publish reads handlers copy under RLock, then calls each handler synchronously. Handlers may themselves call Publish, potentially causing nested recursion but no deadlock.
Assessment: PASS
Severity: LOW
Notes: Event flow decoupled and thread-safe.

## Finding 6
Location: tests/saga_test.go:318-351
Claimed Behavior: Compensation error propagation.
Observed Implementation: When a compensation returns error, orchestrator logs StatusCompensateFailed and aggregates errors. However, orchestrator returns the compensation error aggregated with step error only if compensation fails; but step error is still returned first, which may obscure compensation error.
Assessment: WARNING
Severity: MEDIUM
Notes: Compensation error is logged but not emphasized; could cause silent loss of compensation failure.

## Finding 7
Location: internal/saga/orchestrator.go:63-67
Claimed Behavior: Context cancellation triggers compensation.
Observed Implementation: When ctx.Done(), logs failure and calls compensate. However, the compensate is called with `executed` which includes only steps executed so far. It is fine.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 8
Location: cmd/demo/main.go:70-110
Claimed Behavior: Demo shows happy path and rollback.
Observed Implementation: Two scenarios: Success with order creation, payment, inventory reservation, and approval; and failure at inventory reservation causing rollback of payment and order.
Assessment: PASS
Severity: LOW
Notes: Demo matches specification, showing LIFO rollback.

## Finding 9
Location: internal/services/services.go:114-151
Claimed Behavior: Inventory reservation with release.
Observed Implementation: Release checks reserved[item] < qty condition; could block rollback if release is called with insufficient reserved state (should not happen if orchestrated correctly).
Assessment: PASS
Severity: LOW
Notes: Simple but correct.

## Finding 10
Location: internal/services/services.go:69-112
Claimed Behavior: Payment refunds delete payment.
Observed Implementation: RefundPayment deletes payment entry, not just marking as refunded. That's fine for demo.
Assessment: PASS
Severity: LOW
Notes: Clear compensation.

## Finding 11
Location: tests/saga_test.go:172-171
Claimed Behavior: Semantic lock.
Observed Implementation: Second CreateOrder fails as expected.
Assessment: PASS
Severity: LOW
Notes: Test passes.

## Finding 12
Location: tests/saga_test.go:353-384
Claimed Behavior: Context cancellation triggers compensation.
Observed Implementation: Step1 cancels context, step2 never runs, step1 is compensated.
Assessment: PASS
Severity: LOW
Notes: Works.

## Finding 13
Location: internal/saga/orchestrator.go:92-112
Claimed Behavior: Compensation error aggregation.
Observed Implementation: compErrors slice collects compensation errors, but if there are multiple compensation errors, only the first is returned via fmt.Errorf("%v", compErrors). This is okay for demo.
Assessment: PASS
Severity: LOW
Notes: Not production-grade but acceptable for lab.