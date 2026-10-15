---
title: "Dependency Injection & IoC: Decoupling Object Construction in Go"
published: false
description: "Explore Dependency Injection and Inversion of Control in Go, contrasting explicit constructor injection against the Service Locator anti-pattern."
tags: go, architecture, testing, cleanarchitecture
---

# Dependency Injection & IoC: Decoupling Object Construction in Go

When a business service creates its own infrastructure dependencies—such as invoking a concrete HTTP client or database driver inside domain methods—the system couples tightly to external implementations. Testing requires actual network boundaries, third-party outages break local test suites, and swapping external vendors forces modifications across domain code.

> Dependency injection externalizes object creation: modules receive dependencies from the outside rather than assembling infrastructure themselves.

## Hardcoded dependencies break isolated testing

Directly instantiating concrete adapters inside domain logic binds business rules to network transports. In Lab 16, embedding `&RealGateway{}` inside `Processor` forces any test execution to hit real payment networks or crash when offline. Unit tests degrade into flaky integration tests.

## Service locators hide component contracts

The intuitive fix—passing a container registry into every struct—merely shifts coupling to runtime. The constructor signature `NewBadProcessor(c Container)` hides the real dependency on `PaymentGateway`. If a required key is missing from the container, failures occur during runtime request execution rather than at compile time.

## IoC is an architectural direction, not a library

Inversion of Control dictates the direction of execution authority. Rather than having domain models reach out to instantiate database pools or payment gateways, an external composition root assembles the dependency graph. The domain remains passive: it defines what it needs, and the caller injects it.

## Request flow crosses abstraction boundaries

The client calls `ProcessPayment(amount)` on the domain service. The `Processor` applies validation rules, constructs value objects, and delegates transport to `PaymentGateway`. Because the domain references only the interface contract, alternate gateway implementations swap in without modifying domain code.

```go
type PaymentGateway interface {
    Charge(m Money) error
}
```

## Constructors expose required dependencies

Using explicit constructor parameters transforms hidden couplings into clear interface contracts. `NewProcessor(g PaymentGateway)` makes missing dependencies visible during compilation. Callers cannot instantiate the struct without satisfying the interface signature.

```go
type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}
```

## Nil interfaces cause runtime panics

Go interfaces can hold a nil value while satisfying the type signature. Passing `nil` into `NewProcessor(nil)` compiles cleanly, but calling `p.gateway.Charge(m)` triggers a runtime nil pointer dereference. Constructors in production systems require explicit nil guards.

## Manual wiring incurs composition boilerplate

Manual dependency injection requires explicit assembly in `cmd/demo/main.go`. While this increases boilerplate compared to magic reflection containers or code-generation tools like Google Wire, it eliminates hidden container lookups and retains clear compile-time traceability across the object graph.

```go
realGateway := &di.RealGateway{}
processor := di.NewProcessor(realGateway)
```

## Test harness isolates two architectural paths

The experiment configures six test cases in `tests/processor_test.go`: three targeting constructor injection in `Processor` and three targeting the `BadProcessor` service locator pattern. A `MockGateway` in-memory test double captures invocations without triggering external I/O.

## Demo output confirms identical execution paths

Both architectural patterns execute valid payment workflows successfully against the demo gateway. The distinction between constructor injection and service locator lies not in raw runtime throughput, but in testability, boundary isolation, and contract transparency.

```bash
$ go run ./cmd/demo
--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD
```

## Mock tests verify logic isolation, not production scale

Passing mock tests proves that domain logic routes calls correctly and isolates third-party side effects. It does not prove network resiliency, connection pooling stability, or timeout handling when external gateways degrade under heavy load.

## Explicit contracts balance against constructor bloat

Constructor injection guarantees that every dependency is visible at creation time. However, as services grow, constructors risk accumulating dozens of parameters. This friction is not a flaw in DI; it serves as a compile-time signal that a struct is violating the Single Responsibility Principle.

## Value objects bypass dependency injection

Pure data structures without behavior tied to external infrastructure—such as `Money{Amount: 100, Currency: "USD"}`—should be instantiated directly with struct literals. Wrapping value objects in factories or injecting them through containers introduces gratuitous abstraction.

```go
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

## Comparing hardcoded coupling to constructor injection

Directly instantiating concrete gateways traps business logic inside specific vendor packages. Constructor injection moves instantiation decisions to the composition root, allowing test doubles and production adapters to swap seamlessly.

| Design Aspect | Hardcoded Gateway | Constructor Injection |
|---|---|---|
| Dependency Location | Created inside domain methods | Injected via constructor signature |
| Unit Test Isolation | Requires real network | Swaps with in-memory MockGateway |
| Extensibility | Modifies domain code | Swaps at composition root |

## Boundaries of the manual injection lab

This lab isolates the structural mechanics of dependency injection using pure Go standard library idioms. It intentionally omits automated reflection containers, lifecycle scopes (singleton vs transient), and dynamic circular dependency resolution algorithms.

## Domain validation guards infrastructure boundaries

Before delegating execution to the injected gateway, `ProcessPayment` validates business inputs. Testing with negative values (`amount = -10`) verifies that domain errors reject execution immediately, ensuring that `mock.ChargedMoney.Amount` remains zero and external gateways are never called.

```go
func TestProcessor_InvalidAmount(t *testing.T) {
    mock := &MockGateway{}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(-10)
    if err == nil {
        t.Fatalf("expected error for negative amount, got nil")
    }
}
```

## Verified claims from the test suite

Automated test execution confirms six core assertions across the suite:
1. `Processor` delegates valid amounts to `MockGateway` without network calls.
2. Gateway errors propagate transparently back to callers.
3. Negative amounts abort before touching infrastructure interfaces.
4. `BadProcessor` mirrors functional delegation through the container.
5. In-memory doubles run in under 0.01 seconds.
6. `go test -race ./...` executes cleanly without data race warnings.

## Evidence limits and unverified properties

The unit test suite does not demonstrate performance improvements or defect reduction rates across large engineering teams. Furthermore, clean execution under `go test -race` reflects only the tested sequential paths, not high-concurrency production lifecycles.

## Mitigating god objects and parameter explosion

When a constructor requires excessive interface arguments, developers often revert to service locators to tidy the signature. The robust solution is architectural refactoring: decompose the oversized service into focused domain handlers with single responsibilities.

## Architectural checklist for dependency injection

1. Define interface contracts at consumption boundaries, not provider packages.
2. Adopt constructor injection as the default instantiation pattern.
3. Validate non-nil dependencies explicitly within constructors.
4. Centralize component wiring inside a dedicated composition root (`main.go`).
5. Instantiate pure value objects directly without container indirection.

## Anti-pattern: passing containers into domain services

PSR-11 and software engineering standards specify that applications should not inject general containers into domain models. Doing so couples business logic to container lookup APIs, defeats compile-time safety, and turns unit testing into container mocking overhead.

```go
// Anti-Pattern: Domain depends on container lookup
type BadProcessor struct {
    container Container
}
func (p *BadProcessor) ProcessPayment(amount int) error {
    return p.container.GetPaymentGateway().Charge(m)
}
```

## Reproducing the verification suite locally

The lab contains zero external third-party dependencies and runs entirely via the Go toolchain.

```bash
cd labs/16-dependency-injection
go test -v ./...    # verify 6 isolated unit test assertions
go test -race ./... # verify race detector on test paths
go run ./cmd/demo   # execute manual wiring vs container demo
```

## Engineering principles for decoupled systems

Dependency injection makes architectural dependencies visible at compile time. Separating configuration from execution allows systems to evolve independently, enabling fast unit testing while keeping domain logic clean of infrastructure details.

## Source

[Software Engineering Lab 16 - Dependency Injection](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `76f0059`. Source: README, implementation (`internal/di/`), tests (`tests/processor_test.go`), and demo execution records.
