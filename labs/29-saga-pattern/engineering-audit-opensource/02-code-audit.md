# Code Audit

## Finding 1

Location: internal/saga/orchestrator.go:49-90 (Execute)
Claimed Behavior: Steps execute sequentially; on failure compensates completed steps in reverse order (LIFO).
Observed Implementation: Execute copies steps under mutex, iterates sequentially, appends failed step log, calls compensate() for executed steps. compensate() iterates executed in reverse, appends per-step log entries.
Assessment: PASS
Severity: LOW
Notes: compErr handling has a redundant double-return branch (lines 65-68); if compensation fails during cancellation the earlier return with both ctx.Err() and compErr fires first, making lines 67-68 unreachable. Not a safety issue, but dead code.

## Finding 2

Location: internal/saga/orchestrator.go:92-112 (compensate)
Claimed Behavior: LIFO rollback; logs compensation outcomes.
Observed Implementation: Reverse loop over executed; calls step.Compensate; logs Compensated/CompensateFailed. Aggregates errors into a slice and returns fmt.Errorf("%v", compErrors).
Assessment: PASS
Severity: LOW
Notes: Known limitation (compensation assumed to succeed) documented in engineering/02-implementation-notes.md. Error aggregation wraps all compensation errors but discards order info beyond step name — acceptable for this scope.

## Finding 3

Location: internal/services/services.go:29-40 (OrderService.CreateOrder)
Claimed Behavior: Semantic lock prevents concurrent saga from modifying pending state.
Observed Implementation: CreateOrder locks, checks s.locks[orderID], sets PENDING and lock=true. CancelOrder and ApproveOrder delete the lock.
Assessment: WARNING
Severity: MEDIUM
Notes: Lock is never released if CreateOrder succeeds but the saga is never completed (e.g., process crashes before Approve/Cancel). Lock persists indefinitely for that orderID. No TTL or cleanup path. Not a concurrency bug, but a resource leak scenario.

## Finding 4

Location: internal/services/services.go:82-97 (PaymentService.ProcessPayment)
Claimed Behavior: Idempotency via processedID key; retries produce identical state.
Observed Implementation: On duplicate paymentID, returns nil without writing payments map. First call writes amount and marks processedID=true.
Assessment: PASS
Severity: LOW
Notes: Idempotency works for identical paymentID and amount; no validation that amount matches on retry (silent overwrite of amount on first call only). Sufficient for lab scope.

## Finding 5

Location: internal/saga/orchestrator.go:30-34 (Orchestrator struct)
Claimed Behavior: Thread-safe orchestrator for concurrent saga executions.
Observed Implementation: Orchestrator holds steps and logs under a single mutex; Execute copies steps slice then unlocks before execution loop. compensate() takes mutex per step for logging.
Assessment: WARNING
Severity: MEDIUM
Notes: A single Orchestrator instance is not safe for concurrent Execute calls (steps slice mutated by AddStep during execution). Tests create one orch per goroutine so no actual race. Doc/README does not warn against shared orchestrator reuse.

## Finding 6

Location: internal/saga/choreography.go:44-51 (EventBus.Publish)
Claimed Behavior: Decoupled event bus for choreography model.
Observed Implementation: RLock, copy handler slice, RUnlock, then call handlers sequentially.
Assessment: PASS
Severity: LOW
Notes: Handlers execute synchronously in Publish caller; slow handlers block subsequent ones and the caller. Acceptable for in-memory simulation.

## Finding 7

Location: cmd/demo/main.go:65-66 (Demo Scenario 1)
Claimed Behavior: Happy path marks OrderState=APPROVED, stock=0.
Observed Implementation: Matches output; ApproveOrder deletes semantic lock, sets APPROVED.
Assessment: PASS
Severity: LOW

## Finding 8

Location: cmd/demo/main.go (Scenario 2)
Claimed Behavior: ReserveInventory fails → compensates Payment then Order in LIFO.
Observed Implementation: Matches output; stock remains 0, payment removed, order cancelled.
Assessment: PASS
Severity: LOW

## Finding 9

Location: tests/saga_test.go:318-351 (TestOrchestrator_CompensationErrorPropagated)
Claimed Behavior: Compensation failure is surfaced.
Observed Implementation: Step2 fails, Step1 compensation returns error; test expects logs[2].Status == StatusCompensateFailed.
Assessment: PASS
Severity: LOW
Notes: Logs show CompensateFailed but compensate() continues iterating remaining steps — correct behavior.

## Finding 10

Location: tests/saga_test.go:353-385 (TestOrchestrator_ContextCancellation)
Claimed Behavior: Cancellation triggers compensation of executed steps.
Observed Implementation: Step1 Execute calls cancel(); Step2 Execute asserts.Fatal. Orchestrator sees ctx.Done(), logs StatusFailed, compensates Step1.
Assessment: PASS
Severity: LOW
Notes: Context cancellation only checked at top of loop iteration; a long-running Execute callback would not be interrupted mid-flight. Known stdlib limitation.