# Code Audit

## Finding 1

Location: `internal/saga/orchestrator.go:75-85`
Claimed Behavior: LIFO compensation executed on step failure.
Observed Implementation: Iterates `executed` slice backwards, executing non-nil `Compensate` callbacks and appending to `logs`. Ignores errors returned by compensation routines (`_ = step.Compensate(ctx)`).
Assessment: PASS
Severity: LOW
Notes: Standard simplification marked explicitly with a `ponytail:` comment noting retry/recovery ceiling.

## Finding 2

Location: `internal/services/services.go:33-40`
Claimed Behavior: Semantic locking prevents concurrent state mutation on pending orders.
Observed Implementation: `CreateOrder` checks `locks[orderID]`. Returns error if locked. `ApproveOrder` and `CancelOrder` release lock (`delete(s.locks, orderID)`).
Assessment: PASS
Severity: LOW
Notes: Mutex protection (`s.mu`) guarantees thread-safe map access.

## Finding 3

Location: `internal/services/services.go:86-88`
Claimed Behavior: Idempotent payment processing.
Observed Implementation: `PaymentService` checks `processedID[paymentID]`. Returns `nil` without re-processing if already executed.
Assessment: PASS
Severity: LOW
Notes: Correctly handles duplicate invocations safely.

## Finding 4

Location: `internal/saga/choreography.go:44-51`
Claimed Behavior: Event bus dispatching for event-driven choreography.
Observed Implementation: `Publish` fetches handlers under `RUnlock` copy and executes callbacks synchronously.
Assessment: PASS
Severity: LOW
Notes: Thread-safe, minimal event bus implementation for in-memory choreography demonstration.
