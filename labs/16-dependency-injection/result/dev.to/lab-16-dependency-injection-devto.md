---
title: "Dependency Injection & Inversion of Control: Isolated Testing in Go"
published: false
description: "Separate object configuration from business logic using constructor injection versus service locator in Go, with mock-based unit tests and concurrency verification."
tags: go, testing, dependency-injection, softwarearchitecture, softwareengineering
---

# Dependency Injection & Inversion of Control: Isolated Testing in Go

Business logic that instantiates its own dependencies creates code that resists testing. Lab 16 demonstrates constructor injection versus service locator in pure Go with six deterministic unit tests and a zero-dependency demo.

> Constructor injection makes every dependency explicit in the signature; service locator hides them behind runtime lookup.

## Direct Instantiation Locks Domain to Infrastructure

When a class calls `new ConcreteService` inside a business method, domain logic becomes dependent on external infrastructure. The code may run without errors during early development; this does not mean the architecture can evolve cleanly when a third-party vendor must be replaced. The entire domain layer would require refactoring for what amounts to a single dependency swap.

## Service Locator Looks Clean but Conceals Coupling

Many teams adopt a Service Locator believing they eliminated coupling because direct `new` calls disappear. Injecting a generic container does not automatically produce loosely coupled code when the underlying dependencies remain hidden behind runtime lookup methods. The constructor signature `NewBadProcessor(Container)` tells no one which services are actually required by the processor.

## Separation of Configuration from Use

Fowler's foundational principle establishes two phases: object graph assembly at the composition root and business execution within domain logic. Only along this path do business components remain completely insulated from infrastructure details. Pure value objects like `Money` are instantiated directly and never passed through any container or injector.

## Payment Transaction Execution Pipeline

Execution initiates with concrete assembly at the composition root (`main.go`), followed by local input guards before delegating to the gateway interface. `Processor.ProcessPayment` checks `amount <= 0` and returns an error immediately without touching the gateway. Money structs are assembled directly on the stack. This architecture cleanly isolates logic, though it has not yet proven resiliency under high network latency.

## Constructor Injection: Explicit Dependencies

```go
type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

Processor stores its gateway interface in an unexported field and accepts it strictly during initialization. Relying on interfaces does not imply zero-cost abstractions; it deliberately accepts interface pointer overhead to achieve complete domain isolation.

## Service Locator Failure: Hidden Container Coupling

```go
// locator.go — anti-pattern
func (p *BadProcessor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.container.GetPaymentGateway().Charge(m)
}
```

The Service Locator masks actual dependencies inside method bodies, risking nil-pointer panics when `GetPaymentGateway()` returns nil at runtime. Only along this path are unit tests forced to construct irrelevant container mocks instead of directly supplying mock services.

## Constructor Injection vs Service Locator Matrix

| Dimension | Constructor Injection (Processor) | Service Locator (BadProcessor) |
|---|---|---|
| Dependency Declaration | Explicit in `NewProcessor` params | Hidden behind Container |
| Type Safety | Compile-time check, valid at birth | Runtime resolution, nil-panic risk |
| Test Mocking Overhead | Minimal: inject MockGateway directly | Double: MockContainer + MockGateway |
| Domain Logic Coupling | Zero to framework/container | Tied to Container interface |
| Industry Recommendation | Modern best practice | Discouraged (PSR-11: SHOULD NOT) |

This matrix contrasts Constructor Injection against the Service Locator pattern. Constructor injection guarantees compile-time type safety and explicit dependencies, whereas service locators defer missing-dependency detection to runtime.

## Dependency Assembly Decisions

The architectural sequence establishes three priorities: first, use Constructor Injection as the default for all core services to ensure valid state at birth. Second, instantiate plain value objects like `Money` directly—do not manage DTOs through an injector. Third, use manual wiring in `main.go`; choosing manual wiring does not reflect anti-framework dogma, but rather a pragmatic commitment to Go simplicity and YAGNI.

## Demo Execution and Six Passing Tests

```
=== RUN   TestProcessor_Success
--- PASS: TestProcessor_Success (0.00s)
=== RUN   TestProcessor_GatewayError
--- PASS: TestProcessor_GatewayError (0.00s)
=== RUN   TestProcessor_InvalidAmount
--- PASS: TestProcessor_InvalidAmount (0.00s)
=== RUN   TestBadProcessor_Success / GatewayError / InvalidAmount
--- PASS: TestBadProcessor_* (0.00s)
PASS: ok lab16/tests 1.356s (race detector: clean)
--- Running Constructor Injection: RealGateway charging 100 USD
--- Running Service Locator: RealGateway charging 200 USD
```

Verification demonstrates all six unit test scenarios passing deterministically in sub-second execution. The demo run confirms both injection mechanisms process payments successfully without triggering data races.

## Lab Evidence Boundaries

This lab proves isolated unit testing, input guard execution prior to external infrastructure calls, and gateway error propagation. It does not prove automated lifecycle management (singleton/scoped/transient), industrial defect-rate reduction claims, or behavior under high network latency. Clearly delimit verified properties from unproven enterprise heuristics.

## DI Trade-offs: Indirection and Boilerplate

Adopting DI introduces cognitive indirection because concrete implementations are detached from invocation sites. Manual wiring adds composition lines at the root as the object graph grows. Introducing interfaces and assemblers does not automatically improve architecture if the domain does not require implementation swapping—the cost must be justified by actual flexibility needs.

## Concurrency Safety of Stateless Components

```go
// Processor is stateless after construction
type Processor struct {
    gateway PaymentGateway // read-only after init
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

The Processor struct remains stateless post-construction, rendering it inherently safe for concurrent execution across goroutines. While this simulation is not equivalent to distributed connection pools under real network load, it passes the Go race detector cleanly.

## Design Transformation: Hardcoded vs Constructor Injection

Before DI, business workflows directly instantiated concrete network gateways inside methods, making offline unit testing impossible. Refactoring to constructor injection does not mean adding unnecessary complexity—it adds approximately three lines but transforms flaky integration tests into instant unit tests. The domain layer becomes independently verifiable.

## Minimal Interface Contract and Value Object Exception

```go
type Money struct {
    Amount   int
    Currency string
}

type PaymentGateway interface {
    Charge(m Money) error
}
```

The PaymentGateway interface declares a single Charge method while Money is instantiated directly as a value object. Scoping DI strictly to services rather than value objects reflects design boundaries validated within this lab's test scenarios.

## Isolated Testing with MockGateway

```go
type MockGateway struct {
    ChargedMoney di.Money
    ShouldFail   bool
}

func (m *MockGateway) Charge(money di.Money) error {
    if m.ShouldFail { return errors.New("gateway unavailable") }
    m.ChargedMoney = money
    return nil
}
```

MockGateway simulates network outages and records transaction payloads without touching HTTP sockets. Passing unit tests with mocks does not automatically guarantee third-party integration integrity if production API schemas change.

## Pre-Check Guards Prevent Wasteful Calls

```go
func TestProcessor_InvalidAmount(t *testing.T) {
    mock := &MockGateway{}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(-10)
    if err == nil { t.Fatalf("expected error for negative amount") }
    if mock.ChargedMoney.Amount != 0 {
        t.Errorf("gateway should not have been called, got %d", mock.ChargedMoney.Amount)
    }
}
```

Tests confirm negative amounts are rejected at the domain boundary before contacting the gateway. Only along this pathway does the mock assert that Amount remains 0, confirming that wasteful network invocations were prevented at the domain layer.

## Transparent Error Propagation

```go
func TestProcessor_GatewayError(t *testing.T) {
    mock := &MockGateway{ShouldFail: true}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(50)
    if err == nil { t.Fatalf("expected error from gateway, got nil") }
}
```

The Processor bubbles gateway errors transparently back to the caller. Propagating raw errors does not imply the system lacks retry policies; it confirms that interface boundaries accurately convey downstream failure states without swallowing exceptions.

## Composition Root in cmd/demo/main.go

```go
func main() {
    realGateway := &di.RealGateway{}
    processor := di.NewProcessor(realGateway)
    _ = processor.ProcessPayment(100)

    container := &SimpleContainer{gateway: realGateway}
    badProcessor := di.NewBadProcessor(container)
    _ = badProcessor.ProcessPayment(200)
}
```

The main entrypoint functions as the sole composition root where concrete structs are instantiated and wired. This minimal structure avoids third-party DI reflection frameworks, with its scope strictly matching lab verification needs.

## DI Implementation Checklist

- Infrastructure dependencies declared as interfaces, never concrete structs, in domain code.
- Constructor injection used as default; all dependencies explicit in constructor parameters.
- Composition root (main.go) is the sole assembler of concrete implementations.
- No generic `container.Get()` calls inside business domain logic.
- Pure value objects (Money, DateTime) instantiated directly without any container.
- Every service has isolated unit tests using mocks without network access.

## Common DI Anti-Patterns

Service Locator leak: passing a Container into domain objects hides real dependencies behind runtime method lookups. Constructor bloat: a constructor accepting a dozen parameters signals a violation of the Single Responsibility Principle. Value object injection: managing simple data structs inside a DI container does not make code cleaner; it adds pointless ceremony without measurable benefit.

## Local Reproduction Commands

```sh
cd labs/16-dependency-injection
go test -v ./...    # validate isolated Processor and BadProcessor unit tests
go test -race ./... # verify stateless concurrency is race-free
go run ./cmd/demo   # run Constructor Injection vs Service Locator comparison
```

The entire test harness relies solely on the pure Go standard library with zero third-party dependencies. While this validates architectural isolation locally, production systems may still require code-generation tools like wire at enterprise scale.

## Core Architectural Principles

Dependency Injection fundamentally separates configuration concerns from operational business logic. Leveraging interfaces and constructors does not automatically fix architectural debt if domain boundaries are poorly defined initially. The three pillars are: separate configuration from use at the composition root, declare all dependencies explicitly through constructors, and enforce strict domain boundaries before attempting to apply DI.

## Source

[Software Engineering Lab 16 - Dependency Injection](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `5563838`. Source: internal/di/processor.go, internal/di/gateway.go, internal/di/locator.go, cmd/demo/main.go, tests/processor_test.go, engineering/03-execution-result.md, and go.mod. Code excerpts belong in their original package context.
