# Code Snippets

## Snippet 1 — Constructor Injection

Source File: `internal/di/processor.go`
Purpose: Demonstrates explicit dependency injection through a constructor.

```go
// Processor uses Constructor Injection (Finding 3)
type Processor struct {
	gateway PaymentGateway
}

// NewProcessor ensures valid state and immutability.
func NewProcessor(g PaymentGateway) *Processor {
	return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	m := Money{Amount: amount, Currency: "USD"} // Value object direct instantiation (Finding 5)
	return p.gateway.Charge(m)
}
```

Explanation:
The `Processor` requires `PaymentGateway` at instantiation time, guaranteeing valid initialization and explicitly communicating dependencies. The `Money` struct is instantiated directly without DI.

## Snippet 2 — Service Locator Anti-Pattern

Source File: `internal/di/locator.go`
Purpose: Demonstrates the hidden dependencies and coupling introduced by the Service Locator pattern.

```go
// Container acts as a Service Locator.
type Container interface {
	GetPaymentGateway() PaymentGateway
}

// BadProcessor injects the container directly (Anti-Pattern) (Finding 4)
type BadProcessor struct {
	container Container
}

func NewBadProcessor(c Container) *BadProcessor {
	return &BadProcessor{container: c}
}

func (p *BadProcessor) ProcessPayment(amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	m := Money{Amount: amount, Currency: "USD"}
	return p.container.GetPaymentGateway().Charge(m)
}
```

Explanation:
`BadProcessor` injects `Container` instead of `PaymentGateway`. This obscures what the class actually depends on and couples the class directly to the container interface.

## Snippet 3 — Isolated Unit Testing with Mocks

Source File: `tests/processor_test.go`
Purpose: Demonstrates using a mock implementation to verify behavior in isolation without external network calls.

```go
type MockGateway struct {
	ChargedMoney di.Money
	ShouldFail   bool
}

func (m *MockGateway) Charge(money di.Money) error {
	if m.ShouldFail {
		return errors.New("gateway unavailable")
	}
	m.ChargedMoney = money
	return nil
}

func TestProcessor_Success(t *testing.T) {
	mock := &MockGateway{}
	proc := di.NewProcessor(mock)

	err := proc.ProcessPayment(50)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if mock.ChargedMoney.Amount != 50 || mock.ChargedMoney.Currency != "USD" {
		t.Errorf("expected 50 USD, got %+v", mock.ChargedMoney)
	}
}
```

Explanation:
Because `Processor` relies on an interface injected via constructor, the test harness can easily substitute `RealGateway` with `MockGateway` to verify interaction and state without hitting external systems.
