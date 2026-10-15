---
title: "Dependency Injection: Separating Configuration from Use"
published: false
description: "Explore Dependency Injection in Go: constructor injection, test doubles, the Service Locator anti-pattern, and pure value object boundaries."
tags: go, architecture, testing, softwareengineering
---

# Dependency Injection: Separating Configuration from Use

A business service that instantiates its own infrastructure couples domain logic to external vendors. Lab 16 isolates execution from object construction using Go interfaces and constructor injection.

> Separate configuration from use: assemble dependencies at the boundary, execute pure logic inside.

The lab verifies in-memory isolation, error propagation, and defensive validation using the Go standard library. It does not measure production DI framework performance or distributed container lifecycle management.

## Hard-Coded Infrastructure Breaks Domain Isolation

Instantiating concrete dependencies inside methods locks business logic to single vendor implementations. Unit tests cannot run without opening network connections or configuring fragile mock servers.

```go
// Tightly coupled anti-pattern
func (p *BadProcessor) ChargeDirect(amount int) error {
    gw := &RealGateway{} // Hard-coded external dependency
    return gw.Charge(Money{Amount: amount, Currency: "USD"})
}
```

Modifying vendor implementations causes ripple effects across consumer modules. Testing edge cases becomes slow, brittle, and dependent on network availability.

## Manual Reviews Miss Hidden Dependencies

Passing a generic service locator or container hides the actual contracts required by a struct. Method signatures appear clean, but dependencies remain invisible until inspected line by line.

```go
// locator.go: Container interface hides method requirements
type Container interface {
    GetPaymentGateway() PaymentGateway
}
```

Reviewers assume modularity exists because an abstraction parameter is present. In practice, missing registrations compile cleanly but trigger runtime panics in production.

## Separation of Configuration from Use

Martin Fowler defines dependency injection as separating object assembly from object execution. Composition roots assemble instances at startup, leaving business services to process pure domain workflows.

```go
// Configuration Phase (Composition Root):
realGateway := &di.RealGateway{}
processor := di.NewProcessor(realGateway)

// Use Phase (Domain Execution):
err := processor.ProcessPayment(100)
```

The processor accepts a `PaymentGateway` interface without knowing the concrete struct behind it. Domain logic executes independently of construction mechanics.

## Pipeline of Dependency Injection Flow

Execution flows linearly from the application entry point down to abstracted infrastructure interfaces.

```
cmd/demo/main.go               -> Composition Root initializes RealGateway
di.NewProcessor(gateway)       -> Constructor receives PaymentGateway interface
processor.ProcessPayment(100)  -> Validates input (amount <= 0) and builds Money
gateway.Charge(Money)          -> Delegates execution to concrete implementation
```

The composition root owns construction, the processor enforces business invariants, and the gateway handles outbound network transport.

## Constructor Injection: Valid State at Birth

Unexported struct fields paired with explicit constructor parameters ensure instances are never partially initialized.

```go
type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}
```

Callers cannot bypass required dependencies. The constructor signature documents all structural requirements directly at compile time.

## Early Validation Cuts Downstream I/O Calls

Validating business parameters at domain boundaries stops invalid transactions before invoking infrastructure clients.

```go
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

An input of `-10` aborts execution immediately. Downstream payment gateways are never called with zero or negative amounts, preventing wasted network operations.

## Dependency Mechanism Trade-Off Matrix

Structural mechanisms balance contract transparency, compiler verification, and test setup complexity.

| Criteria | Hard-Coded New | Service Locator | Constructor Injection |
|---|---|---|---|
| Contract Visibility | Hidden inside methods | Hidden in container | Explicit in constructor |
| Compiler Safety | Strong (concrete type) | Weak (runtime lookup) | Strong (interface verified) |
| Test Double Setup | Very high (requires real I/O) | Medium (requires mock container) | Minimal (direct mock injection) |
| Architectural Standard | Monolith anti-pattern | PSR-11: SHOULD NOT pass | Fowler default recommendation |

Constructor injection exposes complete dependency graphs at compile time without runtime indirection overhead.

## Three Structural Design Decisions

Designing clean boundaries requires three sequential architectural rules:

1. **Program to Interfaces:** Abstract external I/O boundaries behind clean interfaces like `PaymentGateway`.
2. **Require Constructor Injection:** Reject ambient service locators; pass dependencies explicitly into `NewProcessor`.
3. **Exempt Pure Value Objects:** Instantiate data structures like `Money{Amount, Currency}` directly without DI wrappers.

Abstractions protect system boundaries while pure data remains lightweight and transparent.

## Terminal Output: Demonstration Run

Executing both patterns demonstrates functional equivalence alongside structural differences.

```sh
$ go run ./cmd/demo
--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD
```

Both approaches execute payment operations, but constructor injection decouples callers from container APIs entirely.

## Boundary Analysis: What the Lab Verifies

Controlled test suites establish verified behavior without overclaiming production scale guarantees.

* **Verified in Lab:** `MockGateway` isolates `Processor` with zero network overhead; validation pre-checks prevent gateway calls; errors propagate intact; race detector reports zero data races.
* **Not Verified in Lab:** Runtime performance overhead of reflection-based containers; defect reduction metrics across enterprise engineering teams; lifecycle scopes (Singleton, Scoped, Transient); automated DI frameworks (Wire, Dig) at scale.

Lab tests verify isolation and contract adherence, not ecosystem-scale framework performance.

## Architectural Trade-Offs

Decoupled architectures trade local construction simplicity for modularity and testability.

* **Benefits Gained:** Total test isolation with microsecond unit tests; vendor-agnostic business logic; immutable dependency references post-construction.
* **Costs Incurred:** Manual wiring boilerplate in composition roots; code navigation across interface boundaries; potential constructor bloat if Single Responsibility is violated.

Constructor overhead serves as an early architectural warning against structs handling too many distinct responsibilities.

## Validation Pipeline and Concurrency

Stateless services holding immutable dependency references remain safe across concurrent goroutine invocations.

```
Goroutine Request   -> Concurrent calls to ProcessPayment(amount)
Pure Pre-Check      -> amount <= 0 validated without shared state mutation
Value Object Build  -> Money{Amount, Currency} allocated per request
Thread-Safe Charge  -> Downstream gateway executed safely (Race clean)
```

Because `Processor` holds no mutable state, concurrent executions proceed without mutex contention or synchronization primitives.

## Before and After Dependency Injection

Comparing hard-coded instantiation against explicit injection highlights operational and architectural improvements.

* **Before (Tightly Coupled):** `Processor` instantiates `RealGateway` directly; tests require live network connections; vendor swaps demand business logic edits; network failures cannot be tested deterministically.
* **After (Constructor Injection):** `Processor` receives `PaymentGateway` via `NewProcessor`; tests inject `MockGateway` locally; vendor swaps change one line in `main.go`; failure scenarios run deterministically in memory.

Decoupled interfaces replace brittle end-to-end integration setups with fast, reliable unit test suites.

## Test Suite 01: Happy Path Isolation

Unit tests verify clean delegation without network traffic or external infrastructure dependencies.

```go
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

The test double captures `ChargedMoney` in memory, proving accurate argument delegation in zero milliseconds.

## Test Suite 02: Gateway Error Propagation

Simulating infrastructure failures ensures business services forward downstream errors without swallowing context.

```go
func TestProcessor_GatewayError(t *testing.T) {
    mock := &MockGateway{ShouldFail: true}
    proc := di.NewProcessor(mock)

    err := proc.ProcessPayment(50)
    if err == nil {
        t.Fatalf("expected error from gateway, got nil")
    }
}
```

Setting `ShouldFail: true` simulates downstream outages deterministically, guaranteeing the caller receives actionable failure signals.

## Test Suite 03: Short-Circuit Invalid Input

Defensive domain checks abort execution before interacting with downstream infrastructure.

```go
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

Asserting `ChargedMoney.Amount == 0` confirms that invalid inputs never trigger downstream API charges.

## Design Rationale: Value Objects Bypass DI

Pure data structures containing no behavior or external side effects must not be managed by DI containers.

```go
// gateway.go: Pure Value Object
type Money struct {
    Amount   int
    Currency string
}
```

Passing `Money` through containers creates needless indirection. Services crossing I/O boundaries require injection; transient data models should be instantiated directly.

## Service Locator Coupling Risks: PSR-11 Warning

Injecting containers into domain logic replaces one tight coupling with an even broader dependency surface.

```go
// BadProcessor couples to Container API instead of generic interface
type BadProcessor struct {
    container Container
}

func (p *BadProcessor) ProcessPayment(amount int) error {
    return p.container.GetPaymentGateway().Charge(Money{Amount: amount, Currency: "USD"})
}
```

PHP-FIG PSR-11 explicitly warns: *"Users SHOULD NOT pass a container into an object"*. Tests require `MockContainer` wrappers, complicating test fixtures unnecessarily.

## Architectural Implementation Checklist

Enforce clear dependency boundaries during architecture reviews using actionable guidelines:

* Declare dependencies as minimal interfaces rather than concrete structs.
* Adopt constructor injection as the default standard; reject optional setter injection.
* Centralize concrete object creation exclusively inside composition roots (`main.go`).
* Prohibit passing container or locator instances into domain logic.
* Instantiate pure value objects and data transfer objects directly.
* Provide isolated test doubles for every business method test suite.

Following these practices prevents infrastructure leakage into domain business logic.

## Anti-Patterns and Common Pitfalls

Adopting dependency injection without discipline introduces three recurring design flaws:

1. **Constructor Bloat:** Constructors accepting excessive parameters signal Single Responsibility Principle violations. Split the class rather than grouping parameters into a generic container.
2. **Service Locator Disguise:** Passing parameter maps or service registries hides real dependencies behind ambient lookups.
3. **Over-Abstracting Everything:** Defining 1:1 interfaces for internal utilities without multiple implementations adds cognitive overhead without testability benefits.

Constructor bloat serves as design feedback that a class handles too many separate concerns.

## Local Test Suite Reproduction

Verify all invariants, concurrency guarantees, and demonstrations locally using standard Go tooling:

```sh
cd labs/16-dependency-injection
go test -v ./...       # Verifies 6 unit test scenarios
go test -race ./...    # Confirms race-free concurrent execution
go run ./cmd/demo      # Runs comparison between Ctor Injection and Locator
```

The implementation relies strictly on the Go standard library, running deterministically in milliseconds without external tools.

## Long-Term Principle: The Hollywood Principle

Dependency injection implements Inversion of Control via the Hollywood Principle: *"Don't call us, we'll call you."*

High-level domain policies must not depend on low-level infrastructure details. Both must depend on abstractions. Externalizing object configuration preserves business logic maintainability across changing technological stacks.
