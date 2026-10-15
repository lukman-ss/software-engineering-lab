---
title: "Saga Pattern in Go: Distributed Transactions Without Two-Phase Commit"
published: false
description: "Explore Saga Pattern in Go: orchestration vs choreography, LIFO compensating transactions, semantic locks, and verifiable test invariants."
tags: go, microservices, architecture, testing
---

# Saga Pattern in Go: Distributed Transactions Without Two-Phase Commit

Modern microservices isolate relational databases, breaking traditional multi-table transactions. Lab 29 explores distributed coordination using independent local transactions, LIFO compensating rollbacks, and concurrency controls in Go.

> Saga pattern replaces 2PC with independent local transactions and compensating rollbacks.

## When Separate Databases Cannot Commit Together

Microservice architectures give each service its own autonomous datastore. Sequential business workflows spanning orders, payments, and inventory break single-database ACID boundaries. Two-phase commit (2PC) offers theoretical coordination, but locks remote resources synchronously. This blocking behavior degrades system throughput and increases vulnerability to network partitions.

## Why Naive Approaches Fail Cross-Service Transactions

Executing sequential cross-service calls without automated rollback mechanisms creates silent data inconsistencies. Calling services sequentially does not mean workflows recover automatically when downstream steps fail. If inventory reservation fails after payment succeeds, orders remain pending while funds stay deducted, forcing manual reconciliation.

```go
func (o *OrderService) CreateFullOrder(orderID string) error {
    o.CreateOrder(orderID)     // DB-commit #1
    p.ProcessPayment(orderID)  // DB-commit #2 — if this fails?
    i.ReserveInventory(...)    // DB-commit #3 — if this fails?
    o.ApproveOrder(orderID)    // DB-commit #4
}
```

## Mental Model: Transaction Relay with Backward Compensation

The Saga pattern models distributed execution as an estafet relay of local transactions. Each step executes against its local datastore and commits immediately. If step *n* fails, the orchestrator triggers compensations backward across completed steps, unwinding actions in reverse to achieve eventual consistency.

## Orchestration Pipeline: From Step Execution to Compensation

A centralized orchestrator coordinates forward steps and backward rollbacks. Recording execution logs does not automatically eliminate distributed state corruption if compensations stall, but when any step returns an error, forward execution halts and the orchestrator invokes compensating actions on completed steps in reverse order.

## Core Implementation: Orchestrator Step and LIFO Compensate

The `Step` struct binds an `Execute` and `Compensate` function pair. The orchestrator records completed steps and unwinds them in Last-In-First-Out (LIFO) order during failures, aggregating compensation errors to avoid masking partial recovery issues.

```go
type Step struct {
    Name       string
    Execute    func(ctx context.Context) error
    Compensate func(ctx context.Context) error
}

func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {
    for i := len(executed) - 1; i >= 0; i-- {
        step := executed[i]
        if step.Compensate != nil {
            _ = step.Compensate(ctx)
        }
    }
    return nil
}
```

## Failure Path: ReserveInventory Fails and Triggers LIFO Compensation

When inventory reservation fails due to exhausted stock, forward progress halts. Reversing execution order functions cleanly only on this verified path: `RefundPayment` runs first to return customer funds, followed by `CancelOrder` to mark the order `CANCELLED`, restoring the system state.

## Mechanism vs Manifestation: Saga Transaction Types

Saga steps classify into compensable, pivot, and retryable transactions. Designating a step as compensable does not mean the underlying database write is reversible without side effects. In this lab, `ApproveOrder` serves as the terminal pivot step without compensation.

| Transaction Type | Forward Execution | Failure Behavior | Compensation Strategy |
|---|---|---|---|
| Compensable | Local DB commit | Triggers rollback chain | Explicit undo (Refund, Cancel) |
| Pivot | Irreversible boundary | Aims for terminal state | None (point of no return) |
| Retryable | Post-pivot operation | Idempotent automated retry | Unneeded (expected to succeed) |

## Architectural Decisions: Orchestration vs Choreography

Choosing coordination patterns depends on workflow complexity. Adopting choreography does not automatically decouple business domains; it exchanges a coordinator for complex event wiring. Orchestration provides centralized control for linear flows, whereas choreography fits event-driven services that react independently.

## Execution Evidence: Deterministic Test Suite Passing

The Go test suite validates both orchestration and choreography across nine automated scenarios. Passing the test suite demonstrates logical correctness, but the scope is bounded by the tests: tests verify happy paths, LIFO rollbacks, payment idempotency, and semantic locks in 0.003s.

```sh
=== RUN   TestOrchestrator_HappyPath
    --- PASS: TestOrchestrator_HappyPath (0.00s)
=== RUN   TestOrchestrator_FailureCompensatesLIFO
    --- PASS: TestOrchestrator_FailureCompensatesLIFO (0.00s)
=== RUN   TestPayment_Idempotency
    --- PASS: TestPayment_Idempotency (0.00s)
=== RUN   TestSemanticLock
    --- PASS: TestSemanticLock (0.00s)
PASS: ok  labs/29-saga-pattern/tests  0.003s
```

## Evidence Boundaries: What Was Proved vs Unverified

This lab proves computational rollback logic and in-memory thread safety within a single Go process. However, this simulation is not equivalent to production distributed systems. The test suite does not yet prove behavior under physical network partitions, broker outages, or database crashes. Production still requires persistent outbox tables and dead-letter queues.

## Architectural Trade-offs: Availability vs Isolation

Sagas deliberately trade ACID isolation for high service availability. Accepting eventual consistency requires additional evidence from integration testing: because local steps commit immediately without holding distributed locks, intermediate states remain temporarily visible to concurrent transactions across services.

| Dimension | Monolithic ACID | Saga Pattern | Production Impact |
|---|---|---|---|
| Isolation | Full serializable | None cross-service | Dirty reads and race risks possible |
| Availability | Low (blocking 2PC) | High (local commits) | Services remain available during outages |
| Consistency | Strong immediate | Eventual consistency | Temporary inconsistencies during rollback |

## Concurrency Pipeline: Mutex Snapshot Pattern and Thread Safety

The orchestrator guards internal collections using `sync.Mutex`. Using mutex snapshots protects local struct fields but does not automatically guard against external service concurrency conflicts. `Execute` copies registered steps to a local slice before releasing the lock, avoiding contention during external RPC calls.

```go
func (o *Orchestrator) Execute(ctx context.Context) error {
    o.mu.Lock()
    steps := make([]Step, len(o.steps))
    copy(steps, o.steps)
    o.mu.Unlock()
    // Steps execute without holding the orchestrator lock
    return nil
}
```

## Before vs After: From Manual Chaining to Saga Orchestration

Refactoring from naive manual chaining to structured sagas replaces uncoordinated writes with traceable state transitions. Without sagas, partial failures leave orphaned orders and deducted payments. With sagas, the orchestrator logs transitions and executes LIFO compensations automatically.

## Service Layer Defenses: Semantic Locks, Idempotency Keys, and Atomic Mutations

Saga reliability relies heavily on defensive primitives within participating services. These service guards protect data integrity only on this verified path: `OrderService` uses semantic locks to reject duplicate orders, `PaymentService` checks idempotency keys, and `InventoryService` mutates stock atomically under local mutex guards.

## Choreography Pattern: EventBus Decoupling via Pub/Sub

The choreography pattern removes centralized coordinators in favor of asynchronous event pub/sub. Pub/sub event decoupling does not mean transactional boundaries disappear; it shifts coordination into subscriber handlers. In-memory `EventBus` pub/sub coordinates payment and inventory responses without central coordinators.

```go
type EventBus struct {
    mu       sync.RWMutex
    handlers map[EventType][]EventHandler
}

func (b *EventBus) Publish(ctx context.Context, e Event) {
    b.mu.RLock()
    handlers := append([]EventHandler(nil), b.handlers[e.Type]...)
    b.mu.RUnlock()
    for _, h := range handlers { h(ctx, e) }
}
```

## Invariant Proofs: Payment Idempotency and Order Semantic Locks

Automated unit tests verify service-level defenses under edge-case conditions. Passing unit assertions has not proven resilience against network timeouts, but it validates core invariants: `TestPayment_Idempotency` verifies duplicate charges are rejected, and `TestSemanticLock` confirms active order creation collisions fail cleanly.

```go
func TestPayment_Idempotency(t *testing.T) {
    svc := services.NewPaymentService()
    _ = svc.ProcessPayment("tx-1", 100, false)
    if err := svc.ProcessPayment("tx-1", 100, false); err != nil {
        t.Fatalf("duplicate charge must succeed idempotently: %v", err)
    }
}
```

## LIFO Log Verification: Exact Sequence of Step Statuses

In failure scenarios, step logging must record the exact chronological order of execution and rollback. `TestOrchestrator_FailureCompensatesLIFO` confirms that inventory failure produces five ordered statuses: `EXECUTED`, `EXECUTED`, `FAILED`, followed by `COMPENSATED` for refund and `COMPENSATED` for cancellation.

```go
expectedOrder := []saga.StepStatus{
    saga.StatusExecuted,    // CreateOrder
    saga.StatusExecuted,    // ProcessPayment
    saga.StatusFailed,      // ReserveInventory
    saga.StatusCompensated, // ProcessPayment (Refund)
    saga.StatusCompensated, // CreateOrder (Cancel)
}
```

## Context Cancellation Safety: Detached Context Rollback

Distributed workflows must handle timeouts gracefully without abandoning state. Detaching context allows local compensation to complete, but production still requires distributed trace propagation: when context cancels, the orchestrator detects `ctx.Done()` and invokes `compensate()` with `context.Background()` to complete rollbacks.

```go
func (o *Orchestrator) Execute(ctx context.Context) error {
    select {
    case <-ctx.Done():
        _ = o.compensate(context.Background(), executed)
        return fmt.Errorf("saga cancelled: %w", ctx.Err())
    default:
    }
    return nil
}
```

## Implementation Checklist for Production Saga Patterns

Deploying saga patterns into production requires strict operational rigor across participating services:

- Every compensable step provides a deterministic, tested `Compensate` function.
- Forward executions and backward compensations implement idempotency keys.
- Semantic locks or version checks protect against concurrent dirty writes.
- Coordinators snapshot step definitions before execution for thread safety.
- Compensating rollbacks execute with detached, non-cancelled contexts.
- Step logs record `PENDING`, `EXECUTED`, `FAILED`, and `COMPENSATED` statuses.
- Automated race tests run cleanly under multi-goroutine workloads.
- Dead-letter queues (DLQ) handle unrecoverable compensation failures.

## Saga Anti-Patterns to Avoid in Production

Common architectural mistakes can cause data corruption in production sagas. Avoiding circular dependencies does not automatically resolve stalled worker threads; systems must enforce timeouts. Designing non-idempotent compensating functions causes duplicate refunds upon network retries, while omitting compensations leaves orphaned state.

## Practical Takeaways: Commands and Production Runbooks

Running local commands provides a baseline, but applying them to production requires additional evidence from chaos testing. Operating sagas in production requires monitoring failed steps, alerting on compensation errors, and maintaining operational runbooks for manual reconciliation.

```sh
# Run full automated test suite
go test -v ./...

# Verify multi-threaded execution safety
go test -race -v ./...

# Run interactive demonstration CLI
go run ./cmd/demo
```

## Long-Term Architectural Principles: Eventual Consistency as a Contract

The Saga pattern establishes that eventual consistency is a deliberate design contract, not an accidental compromise. Embracing local database autonomy and explicit compensation replaces fragile distributed locks with observable recovery paths across distributed boundaries.

## Source

- Repository: [software-engineering-lab/labs/29-saga-pattern](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern)
- Commit SHA: `ba4b6216`
- Source files inspected: `internal/saga/orchestrator.go`, `internal/saga/choreography.go`, `internal/services/services.go`, `tests/saga_test.go`, `cmd/demo/main.go`
