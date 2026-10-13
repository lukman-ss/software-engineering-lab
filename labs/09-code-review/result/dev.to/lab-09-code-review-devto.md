---
title: "Code Review: When Checkout Success Still Breaks Inventory"
published: false
description: "Review a Go checkout through overselling, rollback, retry protection, ownership checks, and the failure windows that remain after improvement."
tags: go, backend, testing, codereview
---

A checkout can return success after reducing stock without creating an order. It can also accept two purchases when only one unit exists. Lab 09 makes those risks concrete by comparing `CheckoutNaive` with `CheckoutImproved` and inspecting the assertions around both implementations.

The useful review unit is a change in behavior: which requirement it serves, what can fail, who is affected, and what evidence protects the fix. A transaction or an idempotency interface is meaningful only when its actual boundary addresses that failure.

## Start with requirements and observable behavior

For this checkout, the important requirements are that stock stays nonnegative, stock reservation and order creation succeed together, a retry does not create another order, and a principal can only check out their own cart.

Review in this order: requirement → changed behavior → failure mode → impact → simplest sufficient fix → test evidence. Naming and extraction suggestions can wait when inventory integrity or authorization is at risk.

All technical observations below come from [Lab 09](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/09-code-review), inspected at commit `82af988b89439fa92b1200406c5afbd59ce72a56`. Go is unavailable in this content-generation environment; the tests were read, not executed. Described outcomes are assertions, not recorded run results. Go excerpts are copied from the source; surrounding code is omitted where indicated. The README's PHP example is illustrative pseudocode, not the executable implementation.

## Finding 1: success hides a partial checkout

### Code

The naive service updates stock inside its product loop and then ignores the result of order creation:

```go
// checkout_naive.go; excerpt, surrounding code omitted
_ = c.repo.CreateOrder(ctx, order)

_ = c.notify.SendOrderConfirmation(ctx, userID, orderID)

return &CheckoutResponse{Success: true, OrderID: orderID, Total: total}, nil
```

### Failure mode and impact

If order creation fails, previously applied stock deductions remain. The caller receives success and a notification is attempted for an order that may not exist. This is a correctness issue, not merely an error-message problem.

The actual naive code does propagate `GetCart` failure. It also skips an item when `GetProduct` fails; it does not dereference a nil product on that branch. Yet `toOrderItems` still includes all cart items, while the total only includes items that passed product lookup. That creates another inconsistency the abbreviated README example does not explain.

### Fix and test

The improved version reserves every item and creates the order inside one callback:

```go
// checkout_improved.go; excerpt
 txErr := c.txManager.WithinTransaction(ctx, func(tx CheckoutTx) error {
    for _, item := range cartItems {
        if reserveErr := tx.ReserveStock(ctx, item.ProductID, item.Quantity); reserveErr != nil {
            return reserveErr
        }
    }
    if createErr := tx.CreateOrder(ctx, order); createErr != nil {
        return createErr
    }
    return nil
})
```

`TestImprovedCheckout_IdempotencyFailureCanRetry` asserts that an injected order-creation failure leaves stock unchanged, creates no order, sends no notification, and permits a later retry. `TestImprovedCheckout_MultiProductFailureRollsBackEverything` checks that both product stocks retain their original values when order creation fails.

Those assertions describe the mock transaction contract. They do not validate a database adapter, because the lab contains no production database transaction implementation.

## Finding 2: synchronized memory can still oversell

### Code

```go
// checkout_naive.go; excerpt
product.Stock -= item.Quantity
_, _ = c.products.UpdateStock(ctx, item.ProductID, -item.Quantity)
```

In `MockProductRepository`, `GetProduct` returns a new `Product`. Changing its `Stock` field does not persist that value. `UpdateStock` applies a delta under a mutex, but it has no condition preventing a negative result. Its error is ignored by the naive service.

### Failure mode and impact

Two requests can each deduct one unit from an initial stock of one, leaving negative inventory. The business failure is overselling. This does not require an unsynchronized memory access: locks can protect each map access while still allowing an invalid business operation.

The barrier in `TestNaiveCheckout_DeterministicallyDemonstratesOverselling` makes both product reads occur before the requests proceed. Its assertions expect two logical successes and negative final stock. The naive path would also accept a sequential deduction larger than available stock; concurrency is not its only problem.

### Fix and test

The improved service calls `CheckoutTx.ReserveStock`. The actual mock transaction stages an accumulated deduction per product, rejects insufficient stock, and applies deductions on commit. Repeated entries for the same product share that accumulated reservation.

The README supplies this production-oriented SQL illustration:

```sql
UPDATE products
SET stock = stock - $1
WHERE id = $2
  AND stock >= $1;
```

Its explanation requires checking affected rows. This SQL is documentation, not a database adapter implemented or exercised by the lab. `BEGIN` and `COMMIT` alone do not express the stock condition.

`TestImprovedCheckout_ConcurrentStockReservationPreservesInvariant` asserts one success, one failure, one order, and final stock zero for two buyers competing for one unit. `TestImprovedCheckout_DuplicateProductInCartAccumulatesReservation` checks both sufficient and insufficient accumulated reservations.

### What the mock actually proves

`MockTransactionManager` holds one mutex per manager for the complete callback and commit. For requests sharing that manager, transactions are serialized. This is the global serialization boundary of the demonstration, not production database concurrency control.

`go test -race ./...` is relevant to unsafe concurrent memory access. It cannot establish the stock invariant by itself. Also, despite its name and comment, `TestImprovedCheckout_DeterministicNoRace` does not use the naive read hook on the improved path: improved calls `GetProducts`, which does not invoke that hook. Its outcome assertions remain useful, but it does not force the same read interleaving.

## Finding 3: retries need a claim, not just a key

### Code

```go
// checkout_improved.go; excerpt
scopedKey := fmt.Sprintf("checkout:%s:%s", principal.UserID, cmd.IdempotencyKey)
```

`hashRequest` marshals `CheckoutCommand` and hashes it with SHA-256. That command contains only `CartOwnerID` and `IdempotencyKey`; it does not contain quantities, product IDs, prices, or a cart version.

### Failure mode and impact

Without retry protection, repeating a checkout can deduct stock and create an order again. Hashing mutable cart contents would also make a completed retry depend on the cart's later state.

The improved service first checks for a completed record before loading the cart. If it finds a stored response, it returns that response immediately. Otherwise it loads and validates the cart, calculates the command hash, and calls `Claim`. Contrary to the README's broad wording, the atomic claim and hash comparison happen after cart loading; only the completed-record lookup happens before it.

### Fix and test

The mock repository protects `Claim` with a mutex. A new key becomes `PROCESSING`. The same hash while processing returns `ErrDuplicateRequest`; a different hash returns `ErrIdempotencyConflict`; a completed matching claim returns the stored response. `TestIdempotencyRepository_StateMachineDirect` asserts those repository behaviors with explicit test hashes.

At the service level, `TestImprovedCheckout_DuplicateRequestReturnsSameResult` expects the same order ID and one order after two requests. `TestImprovedCheckout_RetryAfterCartMutationReturnsOriginalResponse` empties the cart between attempts and asserts replay without another stock deduction. The same-key concurrency test asserts one order, one stock deduction, and one notification; it does not require both callers to receive success.

### The conflict boundary is narrower than it sounds

The early `Get` replay does not compare `RequestHash`. With the current command shape and ownership check, the same scoped key normally implies the same authorized command fields. Still, the service does not implement hash validation on every completed replay, and the hash cannot distinguish changed cart contents. If the command grows, the fast path needs review. The direct repository conflict test is not evidence that an expanded service payload would be checked correctly.

Failures before business commit generally release the claim so another attempt can run. Release can fail too. The transaction-error path joins that error with `errors.Join`; the missing-product-in-batch branch logs a release failure but returns only `ErrProductNotFound`. Error preservation is therefore not uniform across all branches.

## Finding 4: authentication context must constrain cart access

### Code

```go
// checkout_improved.go; excerpt
if principal.UserID != cmd.CartOwnerID {
    c.logger.Error(ctx, "authorization failed", "userID", principal.UserID, "cartOwnerID", cmd.CartOwnerID)
    return nil, ErrForbidden
}
```

### Failure mode, impact, fix, and test

Accepting a caller-selected cart owner without comparing it with the principal could permit checkout of another user's cart. Improved checks ownership before cart loading or replay. `TestImprovedCheckout_CannotCheckoutAnotherUsersCart` asserts forbidden access, unchanged stock, and no order.

Per-user key scoping solves a separate problem: two users submitting the same key string should not share an operation. `TestImprovedCheckout_IdempotencyKeyIsScopedPerUser` asserts two orders and the combined stock deduction.

The naive service only accepts a `userID`; it has no independent principal/owner comparison. There is no HTTP authentication middleware in this lab, so it does not establish how that identity is authenticated upstream. The equality check also does not validate that a principal ID is nonempty.

## Finding 5: repeated lookup is a measurable call pattern

### Code and failure mode

Naive calls `GetProduct` for each cart item. Improved collects unique product IDs and calls `GetProducts` once. Repeated product lookup is the N+1-style issue here; multiple stock updates are a separate write pattern.

### Impact, fix, and test

Batch loading reduces repeated repository lookup work. `TestNaiveCheckout_NPlusOneProductLookups` expects three single-product calls for three items. `TestImprovedCheckout_BatchLoadsProducts` expects one batch call and zero single calls. `TestImprovedCheckout_BatchLoadUsesUniqueProductIDs` checks that repeated product entries produce two unique requested IDs in its fixture.

These counters measure repository calls, not database round trips or query latency. There is no benchmark result to report.

The improved total uses `product.UnitPrice`, whereas naive uses the cart's price. This changes the price source, but quantity positivity is the explicit validation; negative prices, multiplication overflow, and accumulated-total overflow are not guarded. The mock fixes product prices at 1000, and the batch-loading test does not assert pricing behavior.

## Finding 6: commit is not the end of every operation

### Code

```go
// checkout_improved.go; excerpt
if markErr := c.idempotency.MarkCompleted(ctx, scopedKey, resp); markErr != nil {
    c.logger.Error(ctx, "failed to finalize idempotency record; business transaction was committed", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", markErr.Error())
    return nil, fmt.Errorf("%w: %v", ErrIdempotencyFinalize, markErr)
}
```

### Failure mode and impact

Business commit occurs before `MarkCompleted`. Finalization failure can therefore return an error after stock has been deducted and the order exists. The claim remains processing in the mock, and notification is not reached. This is an atomicity gap between two state transitions.

`ErrIdempotencyFinalize` is wrapped for `errors.Is`; the underlying `markErr` is formatted with `%v`, not wrapped as a second cause. This error must not be read as proof that checkout rolled back.

### Fix boundary and test

`TestImprovedCheckout_IdempotencyFinalizeFailureAfterCommit` injects a finalization failure and asserts the sentinel error, stock 8 from an initial 10, one order, no notification, and no completed record. The README discusses a shared transactional boundary, reconciliation, and a unique business constraint as production options. Those options are not implemented here, nor is an Outbox.

After successful finalization, notification runs synchronously outside the transaction. Its error is logged without changing checkout success. A crash before delivery can lose the message, and completed replay does not resend it. The supplied notification mock always succeeds, so the test suite does not exercise a notification-send failure.

Contextual logger calls include user, key, order, and error information where applicable. `MockLogger` only stores messages, however, so these tests do not verify preservation of the structured fields. `TestImprovedCheckout_PreservesRepositoryErrors` checks that a batch infrastructure error is not mislabeled as `ErrProductNotFound`; it does not assert `errors.Is(err, infraErr)`.

## Rank the risk before polishing the code

Inventory corruption and cross-user access deserve attention before naming. Partial state, swallowed errors, duplicate execution, and post-commit ambiguity change what a client can safely assume. Batch loading matters too, but its urgency depends on workload and impact; a call count does not establish severity on its own.

The README's general findings include hardcoded environment configuration and duplicate validation code. Those are discussion topics, not verified concrete defects in the executable checkout. Several README test names are illustrative rather than the names present in `checkout_test.go`; the references above use the actual test functions.

The improved service is a teaching implementation with in-memory repositories. Its `atomic.Int64` counter belongs to each service instance, so IDs are only unique within that instance's running counter; they can collide across instances or replicas and after restart. There is no implemented TTL/recovery mechanism for stale processing claims, database uniqueness constraint, durable notification delivery, or general exactly-once guarantee.

**Code review reduces risk by protecting behavior and business invariants.** A review is stronger when it connects a concrete branch to a failure, a bounded fix, and an assertion that would expose regression.

## Source files

- [README.md](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/README.md): review framework and documented SQL illustration.
- [checkout_naive.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/checkout_naive.go): intentionally broken checkout.
- [checkout_improved.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/checkout_improved.go): service orchestration and remaining gaps.
- [domain.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/domain.go): command, interfaces, and errors.
- [mock.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/mock.go): in-memory transaction and idempotency semantics.
- [checkout_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/checkout_test.go): behavior assertions.
- [go.mod](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/09-code-review/go.mod): module and Go 1.25 requirement.
