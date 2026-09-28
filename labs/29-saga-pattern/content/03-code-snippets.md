# Code Snippets

## Snippet 1 — Saga Orchestrator Step and Engine Definitions

Source File: `internal/saga/orchestrator.go`
Purpose: Defines the `Step`, `StepLog`, and `Orchestrator` structures that manage sequential execution and LIFO compensation.

```go
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
```

Explanation: Encapsulates the core building blocks of orchestration-based saga execution, tracking execution status per step.

---

## Snippet 2 — Orchestrator Execution and LIFO Compensation Engine

Source File: `internal/saga/orchestrator.go`
Purpose: Implements sequential step execution and reverse-order (LIFO) compensation upon failure.

```go
func (o *Orchestrator) Execute(ctx context.Context) error {
	o.mu.Lock()
	steps := make([]Step, len(o.steps))
	copy(steps, o.steps)
	o.mu.Unlock()

	executed := make([]Step, 0)

	for _, step := range steps {
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

Explanation: Iterates through registered steps; if any step fails, loops backward through successfully executed steps to invoke their `Compensate` functions.

---

## Snippet 3 — Choreography Event Bus

Source File: `internal/saga/choreography.go`
Purpose: Implements a lightweight event-driven message bus for choreography-based sagas.

```go
type EventType string

type Event struct {
	Type    EventType
	SagaID  string
	Payload any
}

type EventHandler func(ctx context.Context, e Event)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[EventType][]EventHandler
}

func (b *EventBus) Subscribe(t EventType, h EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

func (b *EventBus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	handlers := append([]EventHandler(nil), b.handlers[e.Type]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(ctx, e)
	}
}
```

Explanation: Enables decoupled event publication and subscription across independent microservice handlers.

---

## Snippet 4 — Services with Semantic Locks and Idempotency Keys

Source File: `internal/services/services.go`
Purpose: Implements mock domain services with semantic locking (`OrderService`) and idempotency (`PaymentService`).

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

Explanation: Protects intermediate states against concurrent modifications using semantic locks and ensures duplicate payment requests do not re-charge customers.
