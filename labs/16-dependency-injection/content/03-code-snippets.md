## Snippet 1 — Interface dan Value Object

Source File: `internal/di/gateway.go`
Purpose: Mendefinisikan kontrak abstraksi dan value object yang diinstansiasi langsung.

```go
// Money is a value object instantiated directly. (Finding 5)
type Money struct {
	Amount   int
	Currency string
}

// PaymentGateway abstracts the concrete implementation. (Finding 1)
type PaymentGateway interface {
	Charge(m Money) error
}

// RealGateway is the production implementation simulating a network call.
type RealGateway struct{}

func (g *RealGateway) Charge(m Money) error {
	fmt.Printf("RealGateway charging %d %s\n", m.Amount, m.Currency)
	return nil
}
```

Explanation:
`PaymentGateway` adalah kontrak tunggal yang memungkinkan swap implementasi tanpa mengubah consumer. `Money` adalah value object tanpa dependensi infrastruktur — dibuat langsung dengan struct literal. `RealGateway` mensimulasikan panggilan jaringan dengan print dan return nil.

## Snippet 2 — Constructor Injection

Source File: `internal/di/processor.go`
Purpose: Menunjukkan dependensi eksplisit melalui constructor.

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
`Processor` menerima `PaymentGateway` saat konstruksi, menjamin objek selalu valid dan dependensi bersifat eksplisit di signature. Validasi `amount <= 0` berjalan sebelum delegasi ke gateway. `Money` dibuat langsung tanpa DI.

## Snippet 3 — Service Locator Anti-Pattern

Source File: `internal/di/locator.go`
Purpose: Menunjukkan dependensi tersembunyi akibat injeksi container.

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
`BadProcessor` menerima `Container` bukan `PaymentGateway` — dependensi sebenarnya tersembunyi dan kelas terikat pada API container. Panggilan `p.container.GetPaymentGateway()` adalah indirection yang tidak ada pada Constructor Injection.

## Snippet 4 — Composition Root

Source File: `cmd/demo/main.go`
Purpose: Menunjukkan separation of configuration from use.

```go
type SimpleContainer struct {
	gateway di.PaymentGateway
}

func (c *SimpleContainer) GetPaymentGateway() di.PaymentGateway {
	return c.gateway
}

func main() {
	realGateway := &di.RealGateway{}

	processor := di.NewProcessor(realGateway)
	fmt.Println("--- Running Constructor Injection ---")
	err := processor.ProcessPayment(100)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	container := &SimpleContainer{gateway: realGateway}
	badProcessor := di.NewBadProcessor(container)
	fmt.Println("--- Running Service Locator ---")
	err = badProcessor.ProcessPayment(200)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
```

Explanation:
`main.go` adalah satu-satunya tempat yang membuat `RealGateway`. Constructor Injection menyuntikkannya langsung; Service Locator membungkusnya dulu ke `SimpleContainer`. Output demo terverifikasi: `RealGateway charging 100 USD` lalu `RealGateway charging 200 USD`.

## Snippet 5 — Isolated Testing dengan Mock

Source File: `tests/processor_test.go`
Purpose: Menunjukkan pengujian terisolasi tanpa infrastruktur nyata.

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

func TestProcessor_InvalidAmount(t *testing.T) {
	mock := &MockGateway{}
	proc := di.NewProcessor(mock)

	err := proc.ProcessPayment(-10)
	if err == nil {
		t.Fatalf("expected error for negative amount, got nil")
	}

	if mock.ChargedMoney.Amount != 0 {
		t.Errorf("gateway should not have been called, but got amount %d", mock.ChargedMoney.Amount)
	}
}
```

Explanation:
`MockGateway` mengimplementasikan `PaymentGateway` tanpa jaringan. Karena `Processor` bergantung pada interface, mock disuntikkan langsung via `NewProcessor(mock)`. `TestProcessor_Success` memverifikasi happy path; `TestProcessor_InvalidAmount` memverifikasi validasi mencegah pemanggilan gateway (`ChargedMoney.Amount` tetap 0).
