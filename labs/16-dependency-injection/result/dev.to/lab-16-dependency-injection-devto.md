---
title: "Dependency Injection: Decoupling Construction from Business Logic"
published: false
description: "Explore Constructor Injection vs Service Locator in Go. Decouple construction from use, isolate unit tests with mocks, and avoid container leaks."
tags: go, architecture, testing, softwareengineering
---

# Dependency Injection: Decoupling Construction from Business Logic

When a business component instantiates its own dependencies directly—such as instantiating an external payment gateway inside a transaction processor—infrastructure concerns bleed directly into domain execution. Swapping payment providers, altering transport protocols, or isolating business logic in tests requires modifying production business code.

> Externalize object creation. Separate the construction of an object graph from its runtime use.

Lab 16 investigates this boundary in Go. It contrasts clean **Constructor Injection** against the **Service Locator anti-pattern**, proves how mocks enable deterministic unit testing without network dependencies, and isolates value objects from superfluous dependency injection machinery.

```
CONSTRUCTOR INJECTION:
main() (Assembler) ──── instantiates ────> RealGateway
      │
      └─ passes RealGateway via NewProcessor ──> Processor(gateway PaymentGateway)
                                                     │
                                                     └─ calls Charge(Money)

SERVICE LOCATOR (ANTI-PATTERN):
main() ─── wraps RealGateway ───> SimpleContainer
                                       ▲
                                       │ container.GetPaymentGateway()
BadProcessor(container Container) ─────┘
```

The lab uses manual dependency wiring in `cmd/demo/main.go` and mock doubles in `tests/processor_test.go`. It does not rely on reflection-based DI containers or code generation frameworks.

---

## The Implemented Boundary

The core system consists of five distinct components in `internal/di`:

1. `PaymentGateway`: An abstraction contract (`gateway.go:12-14`) declaring `Charge(m Money) error`.
2. `RealGateway`: Concrete implementation (`gateway.go:17-22`) simulating network calls.
3. `Money`: A value object (`gateway.go:6-9`) representing currency and amount.
4. `Processor`: A business component using constructor injection (`processor.go:6-21`).
5. `BadProcessor`: A business component holding a service locator container (`locator.go:11-25`).

```go
// internal/di/gateway.go:11-14
type PaymentGateway interface {
	Charge(m Money) error
}

// internal/di/gateway.go:6-9
type Money struct {
	Amount   int
	Currency string
}
```

`Money` carries no external network dependencies, database handles, or runtime side effects. It is instantiated directly at the call site:

```go
// internal/di/processor.go:19
m := Money{Amount: amount, Currency: "USD"}
return p.gateway.Charge(m)
```

Direct instantiation of value objects bypasses DI containers entirely. Injecting simple data bags adds ceremony without architectural benefit.

---

## Constructor Injection vs. Service Locator

The structural contrast between clean injection and the service locator anti-pattern is explicit:

### 1. Constructor Injection (`Processor`)

```go
// internal/di/processor.go:6-13
type Processor struct {
	gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
	return &Processor{gateway: g}
}
```

- Dependencies are explicit in the function signature.
- `Processor` cannot be initialized in an incomplete state.
- The struct requires only `PaymentGateway`, having zero knowledge of any DI container or resolution mechanics.

### 2. Service Locator Anti-Pattern (`BadProcessor`)

```go
// internal/di/locator.go:6-17
type Container interface {
	GetPaymentGateway() PaymentGateway
}

type BadProcessor struct {
	container Container
}

func NewBadProcessor(c Container) *BadProcessor {
	return &BadProcessor{container: c}
}
```

In `BadProcessor`, the true dependency (`PaymentGateway`) is hidden behind the `Container` interface. The constructor accepts any object matching `Container`, but the consumer code queries the container during method execution:

```go
// internal/di/locator.go:19-25
func (p *BadProcessor) ProcessPayment(amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	m := Money{Amount: amount, Currency: "USD"}
	return p.container.GetPaymentGateway().Charge(m)
}
```

This violates the principle of least knowledge:
- `BadProcessor` couples to the container API.
- Static type analysis cannot determine which services `BadProcessor` actually requires without inspecting its method bodies.
- Unit testing requires synthesizing a mock container rather than providing the target interface directly.

---

## Execution Flow in the Composition Root

Assembly occurs exclusively in `cmd/demo/main.go`, separating runtime construction from domain logic:

```go
// cmd/demo/main.go:16-35
func main() {
	// Separation of Configuration from Use
	realGateway := &di.RealGateway{}

	// 1. Constructor Injection
	processor := di.NewProcessor(realGateway)
	fmt.Println("--- Running Constructor Injection ---")
	_ = processor.ProcessPayment(100)

	// 2. Service Locator (Anti-Pattern)
	container := &SimpleContainer{gateway: realGateway}
	badProcessor := di.NewBadProcessor(container)
	fmt.Println("--- Running Service Locator ---")
	_ = badProcessor.ProcessPayment(200)
}
```

Running the executable confirms that both paths execute the same downstream call:

```text
$ go run ./cmd/demo
--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD
```

Both patterns achieve runtime dispatch, but Constructor Injection surfaces the contract at compile time, eliminating runtime container lookup.

---

## What the Tests Actually Prove

The test suite in `tests/processor_test.go` exercises both processors under three specific conditions across 6 assertions:

```go
// tests/processor_test.go:10-21
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
```

### 1. Isolated Unit Execution Without Network
`MockGateway` satisfies `di.PaymentGateway`. Tests run without opening ports, reading configuration files, or invoking external network APIs. The entire test suite executes in less than 2 milliseconds.

### 2. State Recording on Success
`TestProcessor_Success` passes `amount = 50`. The test verifies:
- Error is `nil`.
- `mock.ChargedMoney.Amount == 50`.
- `mock.ChargedMoney.Currency == "USD"`.

### 3. Error Propagation
When `mock.ShouldFail = true`, `TestProcessor_GatewayError` confirms that `ProcessPayment` propagates `"gateway unavailable"` verbatim without swallowing or mutating the error.

### 4. Input Validation Pre-check Guards Infrastructure
When `amount <= 0` (`TestProcessor_InvalidAmount` with `-10`), `ProcessPayment` returns `"invalid amount"`. The assertion:

```go
// tests/processor_test.go:56-58
if mock.ChargedMoney.Amount != 0 {
	t.Errorf("gateway should not have been called, but got amount %d", mock.ChargedMoney.Amount)
}
```

proves that domain validation halts execution before calling the injected infrastructure gateway.

---

## Verified Matrix vs. Unverified Claims

| Property | Implemented & Tested | Lab Limitation / Omission |
|---|---|---|
| Dependency Mocking | Yes (`MockGateway` replaces `RealGateway`) | No automated mocking framework used |
| Explicit Dependencies | Yes (`NewProcessor(g PaymentGateway)`) | Manual wiring only; no reflection or codegen container |
| Service Locator Contrast | Yes (`NewBadProcessor(c Container)`) | Minimal single-method container stub |
| Input Validation Guard | Yes (`amount <= 0` rejects before gateway) | Single scalar validation rule |
| Concurrency Safety | Yes (`go test -race ./...` passes) | Structs hold stateless dependencies |
| Container Lifecycle | No | Singleton, Scoped, and Transient lifecycles not implemented |
| Framework Performance | No | Zero memory or CPU allocation benchmarks recorded |
| Defect Reduction Rate | No | No empirical defect metric exists in source code |

---

## Production Trade-offs

1. **Manual Assembly vs. Auto-Wiring Frameworks:**
   Manual construction in `main.go` provides complete compile-time type safety and explicit call hierarchies without reflection magic. For systems with hundreds of interdependent components, manual composition roots require significant boilerplate. Code generation tools (e.g., `google/wire`) offer compile-time safety with automated assembly, while runtime reflection containers (e.g., `uber/dig`) introduce startup panics for missing bindings.

2. **Constructor Bloat as an Architectural Metric:**
   When a constructor grows past 4–6 dependencies, it indicates a violation of the Single Responsibility Principle. The class is orchestrating too many distinct concerns. Introducing a Service Locator to condense these into a single `Container` argument hides the design smell without fixing the underlying cohesion issue.

3. **Interface Segregation:**
   Keep injected interfaces minimal and focused on caller needs. `PaymentGateway` exposes a single method `Charge(m Money) error`. Large monolithic interfaces force consumers and mock doubles to satisfy unused methods.

---

## Repository Artifacts

- Implementation: `internal/di/gateway.go`, `internal/di/processor.go`, `internal/di/locator.go`
- Composition Root: `cmd/demo/main.go`
- Test Suite: `tests/processor_test.go`
- Module: `go.mod` (Go 1.22)
- Commit Reference: `8afe2a824667528d16d6d9b2ac434a18a6ac2db1`
