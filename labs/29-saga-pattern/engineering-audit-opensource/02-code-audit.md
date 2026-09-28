## Finding 1

Location: internal/saga/orchestrator.go lines 49-90 (Execute method)
Claimed Behavior: Orchestrator manages forward steps and LIFO compensation on failure.
Observed Implementation: Steps are executed sequentially; on failure, compensations run in reverse LIFO order using copied slice of executed steps. Context cancellation triggers compensation with background context.
Assessment: PASS
Severity: LOW
Notes: Context cancellation uses background context for compensations, which may be intentional but could be noted as a design decision. No actual deviation from claimed behavior.

## Finding 2

Location: internal/services/services.go lines 29-39 (CreateOrder) and lines 42-51 (ApproveOrder)
Claimed Behavior: Semantic lock prevents concurrent creation of same order.
Observed Implementation: CreateOrder sets locks[orderID]=true; ApproveOrder releases lock only on success. If ApproveOrder fails because order is already cancelled, lock is not released, potentially causing a permanent lock leak for that orderID.
Assessment: WARNING
Severity: MEDIUM
Notes: In the current saga usage, cancellation is performed via CancelOrder (which releases lock), so the leak may not occur. However, the logic is fragile and could cause a denial-of-service if approval fails for other reasons. Recommend releasing lock in all ApproveOrder paths or using defer.

## Finding 3

Location: internal/saga/orchestrator.go lines 92-111 (compensate method)
Claimed Behavior: Compensating transactions execute in LIFO (reverse) order on failure.
Observed Implementation: Compensation loop iterates from len(executed)-1 down to 0, invoking each step's Compensate function. Errors are aggregated and returned.
Assessment: PASS
Severity: LOW
Notes: Compensation runs with context.Background(), discarding original context. This may be intentional to ensure compensations run to completion regardless of original timeout/cancellation.

## Finding 4

Location: internal/services/services.go lines 82-96 (ProcessPayment)
Claimed Behavior: Idempotency prevents duplicate step execution.
Observed Implementation: ProcessPayment checks processedID map; if already present, returns nil (success). If shouldFail is true, returns error without marking processed, allowing retries.
Assessment: PASS
Severity: LOW
Notes: This matches typical idempotency semantics (avoid duplicate successful side effects). Failed attempts are not considered processed, which is correct.

## Finding 5

Location: internal/saga/choreography.go lines 44-52 (Publish method)
Claimed Behavior: Decoupled event bus for choreographing saga steps across independent services.
Observed Implementation: Publish uses read lock, copies handler slice, releases lock, then invokes each handler. Allows concurrent subscription while publishing.
Assessment: PASS
Severity: LOW
Notes: Safe and correct implementation.

## Finding 6

Location: tests/saga_test.go
Claimed Behavior: Test suite covers happy path, failure path, rollback, idempotency, semantic lock, concurrency.
Observed Implementation: Tests include TestOrchestrator_HappyPath, TestOrchestrator_FailureCompensatesLIFO, TestPayment_Idempotency, TestSemanticLock, TestOrchestrator_Concurrency, TestChoreography_Flow, TestChoreography_FailureCompensates, TestOrchestrator_CompensationErrorPropagated, TestOrchestrator_ContextCancellation.
Assessment: PASS
Severity: LOW
Notes: Test suite appears comprehensive and validates the claimed behaviors.