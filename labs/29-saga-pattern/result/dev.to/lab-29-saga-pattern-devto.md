---
title: "Saga Pattern: Distributed Transactions Without 2PC in Go"
published: false
description: "Master distributed transaction management in Go with the Saga pattern: Orchestration, Choreography, LIFO compensation, and semantic locks."
tags: go, microservices, architecture, systemdesign
---

# Saga Pattern: Distributed Transactions Without 2PC in Go

## 1. Introduction: The Database-per-Service Reality

In a monolithic architecture backed by a single relational database, maintaining transactional consistency is trivial. A business flow spanning order placement, customer credit deduction, and warehouse stock reservation executes inside a single database transaction: `BEGIN TRANSACTION`, run the mutations, then issue `COMMIT`. If any operation fails, the database engine rolls back the entire sequence using write-ahead logs and undo segments. Every table reverts to its pre-transaction state automatically.

Transforming a monolithic application into a microservices architecture breaks this model. The standard architectural baseline is Database-per-Service: each service owns its datastore exclusively. The Order Service writes to an orders database, the Payment Service manages payment ledgers, and the Inventory Service controls warehouse records. This isolation guarantees independent deployability, domain boundaries, and granular horizontal scaling.

However, Database-per-Service eliminates global ACID transactions. No single database engine can coordinate row locks or commit markers across independent physical databases separated by a network. A multi-step transaction becomes a series of disjoint local transactions. When the third service fails after the first two have committed their local changes, the application faces partial execution.

Distributed systems replace strict ACID guarantees with the BASE model:
- **Basically Available:** The system prioritizes service availability over instantaneous consistency, accepting transient intermediate states.
- **Soft State:** Data status changes over time as asynchronous compensation messages propagate across the network.
- **Eventual Consistency:** System state converges to a consistent final state once all local transactions or corrective compensations finish executing.

Lab 29 implements the Saga pattern in idiomatic Go without third-party frameworks. It provides a production-grade blueprint for managing multi-service workflows with sequential orchestration, event-driven choreography, LIFO compensation, semantic locking, and idempotency key enforcement.

## 2. Why Two-Phase Commit (2PC) Fails at Scale

Before adopting the Saga pattern, distributed architectures historically relied on Two-Phase Commit (2PC), standardized under the XA protocol. 2PC guarantees atomic commits across heterogeneous databases via two phases: a voting phase (Prepare) and an execution phase (Commit/Abort).

While 2PC provides strict serializability, it introduces structural failure modes that make it impractical for high-throughput cloud environments:

1. **Blocking Protocol:** During the prepare phase, participating databases lock their rows and commit logs while waiting for the coordinator's global commit signal. If the coordinator crashes or experiences a network partition while participants wait, those resources remain locked indefinitely.
2. **Elevated Latency:** 2PC requires multiple synchronous network round-trips between the coordinator and all participants. Total transaction duration is bound to the slowest responding node across the network.
3. **Single Point of Failure (SPOF):** If the centralized coordinator fails after participants vote to commit but before issuing the global commit command, participants enter an in-doubt state. They cannot abort unilaterally (other nodes may have committed) and cannot commit unilaterally (other nodes may have aborted).
4. **Severe Lock Contention:** Holding row locks across network round-trips exhausts database connection pools and blocks concurrent reads and writes, inducing cascading timeouts across upstream services.
5. **Network Partition Vulnerability:** Under the CAP theorem, 2PC chooses consistency over availability. In an asynchronous network where partitions occur, 2PC stalls completely, violating availability requirements.

The Saga pattern eliminates distributed locking entirely. Each local transaction commits immediately within its local database, holding physical row locks for milliseconds rather than seconds.

## 3. Architecture and Core Mechanics

A Saga is a sequence of local transactions: `T1, T2, ..., Tn`. Each transaction updates data within a single service boundary. If a step fails at index `k`, the Saga engine executes a series of compensating transactions `Ck-1, ..., C2, C1` in reverse order (Last-In, First-Out / LIFO).

```
=== SAGA STEP TAXONOMY ===

Step 1: Compensable      Step 2: Compensable      Step 3: Pivot            Step 4: Retryable
[CreateOrder (PENDING)] -> [ProcessPayment]     -> [ReserveInventory]    -> [ApproveOrder]
         |                       |                        |                        |
  Rollback: CancelOrder    Rollback: RefundPayment   Point of No Return       Idempotent Retry
```

### Transaction Classification

Following the formal taxonomy from distributed systems literature and the Microsoft Azure Architecture Center:

1. **Compensable Transactions:** Steps that can be semantically undone by executing an opposing compensating transaction. Creating an order in `PENDING` state or charging a card are compensable steps.
2. **Pivot Transaction:** The decisive boundary in the Saga flow. If the pivot transaction succeeds, the Saga is guaranteed to proceed to completion. If the pivot transaction fails, all previously completed compensable steps must be rolled back. In Lab 29, `ReserveInventory` acts as the pivot transaction.
3. **Retryable Transactions:** Steps occurring after the pivot transaction. These transactions are designed to never fail permanently and must be retried until successful. They require strict idempotency.

### Coordination Models: Orchestration vs. Choreography

Lab 29 implements both coordination archetypes:

```
=== ORCHESTRATION TOPOLOGY ===
             +--------------------+
             |  SagaOrchestrator  |
             +---+-------+------+--+
                 |       |      |
      1. Execute |       |      | 3. Execute
                 v       |      v
          +----------+   |   +-------------+
          | OrderSvc |   |   | InventorySvc|
          +----------+   |   +-------------+
                         v
                  +------------+
                  | PaymentSvc |
                  +------------+
                   2. Execute

=== CHOREOGRAPHY TOPOLOGY ===
+----------+      OrderCreated      +------------+
| OrderSvc | ---------------------> | PaymentSvc |
+----------+                        +-----+------+
     ^                                    |
     |                                    | PaymentCompleted
     | InventoryFailed                    v
+----+------------------------------+-------------+
|               EventBus            | InventorySvc|
+-----------------------------------+-------------+
```

- **Orchestration:** A centralized engine (`Orchestrator`) holds the workflow definition, invokes service steps sequentially, monitors execution logs, and triggers LIFO rollback upon failure. Services remain decoupled from each other, depending only on the orchestrator contract.
- **Choreography:** Decentralized services publish and consume domain events over an `EventBus`. Each service handles its own local transaction upon receiving an event and emits the next event. If a step fails, a failure event (e.g., `InventoryFailed`) triggers compensations in listening services.

## 4. Implementation Walkthrough in Go

The core engine is located in `internal/saga/orchestrator.go`. It requires no external dependencies beyond the Go standard library.

### 4.1 Orchestrator Engine and Step Definition

A step encapsulates an execution closure and a compensating closure:

```go
package saga

import (
	"context"
	"fmt"
	"sync"
)

type StepStatus string

const (
	StatusPending          StepStatus = "PENDING"
	StatusExecuted         StepStatus = "EXECUTED"
	StatusFailed           StepStatus = "FAILED"
	StatusCompensated      StepStatus = "COMPENSATED"
	StatusCompensateFailed StepStatus = "COMPENSATE_FAILED"
)

type Step struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

type StepLog struct {
	Name   string
	Status StepStatus
}

type Orchestrator struct {
	mu    sync.Mutex
	steps []Step
	logs  []StepLog
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		steps: make([]Step, 0),
		logs:  make([]StepLog, 0),
	}
}

func (o *Orchestrator) AddStep(step Step) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.steps = append(o.steps, step)
}
```

### 4.2 Sequential Execution and LIFO Rollback

The `Execute` method processes steps in sequence. It monitors context cancellation and initiates reverse-order compensation if an execution error occurs:

```go
func (o *Orchestrator) Execute(ctx context.Context) error {
	o.mu.Lock()
	steps := make([]Step, len(o.steps))
	copy(steps, o.steps)
	o.mu.Unlock()

	executed := make([]Step, 0)

	for _, step := range steps {
		select {
		case <-ctx.Done():
			o.mu.Lock()
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})
			o.mu.Unlock()
			compErr := o.compensate(context.Background(), executed)
			if compErr != nil {
				return fmt.Errorf("saga cancelled: %w; compensation errors: %v", ctx.Err(), compErr)
			}
			return fmt.Errorf("saga cancelled: %w", ctx.Err())
		default:
		}

		err := step.Execute(ctx)
		o.mu.Lock()
		if err != nil {
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})
			o.mu.Unlock()

			compErr := o.compensate(context.Background(), executed)
			if compErr != nil {
				return fmt.Errorf("step %s failed: %w; compensation errors: %v", step.Name, err, compErr)
			}
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}
		o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusExecuted})
		o.mu.Unlock()
		executed = append(executed, step)
	}

	return nil
}

func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {
	var compErrors []error
	for i := len(executed) - 1; i >= 0; i-- {
		step := executed[i]
		if step.Compensate != nil {
			err := step.Compensate(ctx)
			o.mu.Lock()
			if err != nil {
				o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensateFailed})
				compErrors = append(compErrors, fmt.Errorf("compensation %s: %w", step.Name, err))
			} else {
				o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensated})
			}
			o.mu.Unlock()
		}
	}
	if len(compErrors) > 0 {
		return fmt.Errorf("%v", compErrors)
	}
	return nil
}
```

Notice the critical resilience detail: when `ctx.Done()` triggers cancellation, compensation executes with `context.Background()`. This ensures rollbacks finish cleanly even if the caller's request context has timed out.

### 4.3 Semantic Locking and Idempotency in Domain Services

In `internal/services/services.go`, domain services implement semantic defense mechanisms:

```go
type OrderService struct {
	mu     sync.Mutex
	orders map[string]OrderState
	locks  map[string]bool // semantic lock
}

func (s *OrderService) CreateOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.locks[orderID] {
		return errors.New("order is semantically locked by another operation")
	}

	s.orders[orderID] = OrderPending
	s.locks[orderID] = true
	return nil
}

func (s *OrderService) ApproveOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.orders[orderID] == OrderCancelled {
		return errors.New("cannot approve cancelled order")
	}
	s.orders[orderID] = OrderApproved
	delete(s.locks, orderID)
	return nil
}

func (s *OrderService) CancelOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[orderID] = OrderCancelled
	delete(s.locks, orderID)
	return nil
}
```

The `locks` map acts as a semantic lock. Any concurrent attempt to recreate or reprocess an order ID while it sits in `PENDING` state is rejected at the domain boundary.

The `PaymentService` enforces idempotency via deduplication keys:

```go
type PaymentService struct {
	mu          sync.Mutex
	payments    map[string]float64
	processedID map[string]bool // idempotency key
}

func (s *PaymentService) ProcessPayment(paymentID string, amount float64, shouldFail bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.processedID[paymentID] {
		return nil // Idempotent success
	}

	if shouldFail {
		return errors.New("insufficient funds")
	}

	s.payments[paymentID] = amount
	s.processedID[paymentID] = true
	return nil
}
```

## 5. Verification and Test Execution

The codebase provides test coverage in `tests/saga_test.go` and an interactive execution flow in `cmd/demo/main.go`.

### 5.1 CLI Demonstration Output

Running `go run cmd/demo/main.go` exercises both the happy path and a failure scenario triggered by an out-of-stock condition:

```
=== Saga Pattern Demonstration ===

--- Scenario 1: Happy Path (Orchestrator) ---
 -> [OrderService] Creating order: ord-success
 -> [PaymentService] Processing payment for: ord-success
 -> [InventoryService] Reserving 1 laptop
 -> [OrderService] Finalizing order approval
Scenario 1 Result: error=<nil>, OrderState=APPROVED, Remaining Stock=0

--- Scenario 2: Rollback on Failure (Orchestrator) ---
 -> [OrderService] Creating order: ord-failed
 -> [PaymentService] Processing payment for: ord-failed
 -> [InventoryService] Reserving 1 laptop (current stock: 0 )
 <- [PaymentService] Compensating: Refunding payment for: ord-failed
 <- [OrderService] Compensating: Cancelling order: ord-failed
Scenario 2 Result: error=step ReserveInventory failed: out of stock, OrderState=CANCELLED, HasPayment=false, Stock=0

=== Demo Completed Successfully ===
```

In Scenario 2, `ReserveInventory` fails because remaining stock is 0. The orchestrator intercepts the error and executes compensations in LIFO order: `PaymentService.RefundPayment` runs first, followed by `OrderService.CancelOrder`.

### 5.2 Test Suite Execution with Go Race Detector

Executing the test suite with race detection confirms zero data races across concurrent workers:

```bash
go test -v -race ./...
```

Output:

```
=== RUN   TestOrchestrator_HappyPath
--- PASS: TestOrchestrator_HappyPath (0.00s)
=== RUN   TestOrchestrator_FailureCompensatesLIFO
--- PASS: TestOrchestrator_FailureCompensatesLIFO (0.00s)
=== RUN   TestPayment_Idempotency
--- PASS: TestPayment_Idempotency (0.00s)
=== RUN   TestSemanticLock
--- PASS: TestSemanticLock (0.00s)
=== RUN   TestOrchestrator_Concurrency
--- PASS: TestOrchestrator_Concurrency (0.00s)
=== RUN   TestChoreography_Flow
--- PASS: TestChoreography_Flow (0.00s)
=== RUN   TestChoreography_FailureCompensates
--- PASS: TestChoreography_FailureCompensates (0.00s)
=== RUN   TestOrchestrator_CompensationErrorPropagated
--- PASS: TestOrchestrator_CompensationErrorPropagated (0.00s)
=== RUN   TestOrchestrator_ContextCancellation
--- PASS: TestOrchestrator_ContextCancellation (0.00s)
PASS
ok  	labs/29-saga-pattern/tests	1.450s
```

All 9 test cases pass, verifying LIFO ordering, event-driven choreography transitions, payment idempotency, and concurrent stock allocation across 10 goroutines.

## 6. Architectural Trade-off Analysis

The table below contrasts Monolithic ACID, Two-Phase Commit (2PC), and the Saga Pattern:

| Dimension | Monolithic ACID | Two-Phase Commit (2PC) | Saga Pattern |
|---|---|---|---|
| **Consistency Guarantee** | Immediate Consistency | Immediate Consistency | Eventual Consistency (BASE) |
| **Data Isolation** | Full Isolation (Serializability / Read Committed) | Full Isolation across nodes | No ACID Isolation; intermediate states are visible |
| **Rollback Mechanism** | Automatic DB rollback via undo log | Global abort by coordinator | Explicit compensating business transactions |
| **System Availability** | High for single node, down if DB crashes | Low; blocked by the slowest or unreachable node | High; local transactions commit without holding global locks |
| **Lock Duration** | Milliseconds | Extended across multiple network RTTs | Milliseconds per local service step |
| **Horizontal Scalability**| Limited to single instance or read replicas | Very poor due to coordination overhead | High; services scale independently |
| **Implementation Complexity** | Low; handled natively by database engine | High at infrastructure layer (XA drivers) | High in application logic (custom compensations) |

## 7. Operational Trade-offs and Mitigation Strategies

Adopting Sagas introduces specific operational challenges that teams must design for:

### The Loss of Isolation: Mitigating Data Anomalies

Because local transactions commit immediately, other processes can inspect and modify data before the Saga finishes. This introduces three risks:
1. **Dirty Reads:** A process reads uncommitted business state (e.g., viewing an order as processed while payment is being refunded).
2. **Lost Updates:** Concurrent Sagas overwrite the same record without taking prior mutations into account.
3. **Non-Repeatable Reads:** Re-reading data within a workflow returns different results due to intervening commits.

To mitigate these anomalies, applications rely on:
- **Semantic Locks:** Setting status fields to `PENDING` or `BLOCKED` to warn external readers that a record is undergoing modification.
- **Commutative Updates:** Designing operations so execution order does not alter final outcomes (e.g., relative credit/debit balances rather than absolute overrides).
- **Pessimistic Views:** Filtering pending records from public-facing customer queries until the Saga reaches an approved terminal state.

### Compensation Failures and Dead Letter Queues (DLQ)

Compensating transactions can fail due to upstream outages, network partition, or external service errors. If a refund call fails, the orchestrator records `StatusCompensateFailed`.

In production, uncompensated steps require automated remediation:
- **Persistent Event Logs:** Sagas must persist state machines in durable storage (PostgreSQL/Kafka) before executing each step.
- **Dead Letter Queues (DLQ):** Failed compensation events are published to a DLQ, where background workers retry with exponential backoff.
- **Manual Reconciliation Dashboards:** Unresolvable compensation failures trigger alerts for operational and finance teams to execute manual account adjustments.

### In-Memory vs. Production Storage

Lab 29 uses in-memory maps guarded by `sync.Mutex` to demonstrate concepts clearly without external infrastructure dependencies. If the process terminates abruptly, in-flight transaction states are lost. Production implementations require durable state storage (e.g., Temporal, Cadence, or a database-backed state machine) to resume interrupted Sagas after process restarts.

## 8. Selection Criteria: When to Use Which Pattern

- **Choose Orchestration when:**
  - The workflow contains 4 or more sequential steps.
  - Centralized auditability and workflow visualization are operational requirements.
  - You want to prevent complex circular dependencies between event publishers and subscribers.
- **Choose Choreography when:**
  - The workflow involves 2 to 3 services in a straightforward linear progression.
  - Extreme throughput is required and asynchronous event publishing avoids coordinator overhead.
  - Services are managed by decentralized teams without a shared orchestration platform.
- **Stick with Local ACID when:**
  - The application is a monolith or modular monolith sharing a single relational database.
  - All mutated entities reside within the same service boundary.

## 9. Conclusion and Source Code

The Saga pattern provides a pragmatic approach to transaction management in microservices architectures. By replacing distributed locks with sequential local transactions and LIFO compensations, systems achieve resilience, high throughput, and horizontal scalability.

### Laboratory Resources

- **Source Code Repository:** [software-engineering-lab / labs / 29-saga-pattern](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern)
- **Verified Commit SHA:** `9c525944dc2de00e16513830a25ffa1c4834d742`
- **Further Reading:** Chris Richardson, *Microservices Patterns* (Manning Publications); Hector Garcia-Molina & Kenneth Salem, *Sagas* (ACM, 1987); Microsoft Azure Architecture Center, *Saga Distributed Transactions Pattern*.
