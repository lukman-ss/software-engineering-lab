# Code Audit

Target Lab: labs/29-saga-pattern

## Finding 1

Location: `internal/saga/orchestrator.go:49-90`
Claimed Behavior: Step execution with sequential progression and LIFO compensation rollback on error or context cancellation.
Observed Implementation: Steps execute in order; if any step returns an error or context is done, `o.compensate` is called with executed steps in reverse order (LIFO), tracking logs safely with mutex locking.
Assessment: PASS
Severity: LOW
Notes: Compensation handles errors by collecting and returning aggregated compensation errors if any compensation fails.

## Finding 2

Location: `internal/saga/orchestrator.go:92-112`
Claimed Behavior: Compensating transactions executed in LIFO order; failures during compensation tracked with `StatusCompensateFailed`.
Observed Implementation: Iteration index runs from `len(executed)-1` down to `0`. Failures append to `compErrors` and mark status `StatusCompensateFailed`.
Assessment: PASS
Severity: LOW
Notes: Implementation handles missing compensate callbacks gracefully via nil check.

## Finding 3

Location: `internal/saga/choreography.go:27-52`
Claimed Behavior: Decoupled event bus supporting publish-subscribe for choreography sagas.
Observed Implementation: Handlers registered per `EventType` under `sync.RWMutex`. `Publish` copies slice under read lock to allow safe sequential dispatch.
Assessment: PASS
Severity: LOW
Notes: Synchronous dispatch guarantees deterministic execution within in-memory choreography scenarios.

## Finding 4

Location: `internal/services/services.go:16-67`
Claimed Behavior: Order service with state transitions and semantic locking.
Observed Implementation: `CreateOrder` sets `OrderPending` and locks order ID. `ApproveOrder` and `CancelOrder` release lock and transition status. Mutex protects all map operations.
Assessment: PASS
Severity: LOW
Notes: Double create on locked ID returns error as designed.

## Finding 5

Location: `internal/services/services.go:69-113`
Claimed Behavior: Idempotent payment processing.
Observed Implementation: `processedID` map stores processed transaction IDs. Duplicate calls return `nil` without double payment processing.
Assessment: PASS
Severity: LOW
Notes: Thread-safe via `sync.Mutex`.

## Finding 6

Location: `internal/services/services.go:114-157`
Claimed Behavior: Inventory reservation and compensating release.
Observed Implementation: `Reserve` decrements stock and increments reserved count. `Release` decrements reserved and restores stock with bounds check. Mutex protection on all access.
Assessment: PASS
Severity: LOW
Notes: Clean failure when stock is insufficient.
