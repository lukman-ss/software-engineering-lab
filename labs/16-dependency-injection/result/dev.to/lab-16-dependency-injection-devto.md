---
title: "Dependency Injection & Inversion of Control: Separating Configuration from Use in Go"
published: false
description: "Explore Dependency Injection & Inversion of Control in Go: constructor injection, service locator anti-pattern, manual wiring, mocking, and test isolation."
tags: go, dependency-injection, testing, softwareengineering
---

# Dependency Injection & Inversion of Control: Separating Configuration from Use in Go

A payment processor embeds its gateway directly, making every test spin up real network connections. Lab 16 extracts that dependency so the business logic stays free of infrastructure assumptions.

> Constructor injection makes dependencies explicit; service locator hides them behind a container.

The lab uses a single `PaymentGateway` interface and a `Processor` that delegates to it. It does not implement real payments, connection pooling, or production-grade composition roots.

## Hard-Coded Dependencies Lock Business Logic to Infrastructure

When a class builds its own gateway, unit tests must reach outside the process. Changing providers forces edits across multiple modules instead of swapping one implementation.

```go
// Before: tight coupling
type Processor struct{}

func (p *Processor) ProcessPayment(amount int) error {
    gw := &RealGateway{}           // internal instantiation
    return gw.Charge(Money{...})
}
```

This couples validation, business rules, and network calls into one method. Mocking becomes impossible without refactoring production code.

## Why Manual Code Review Fails to Catch Hidden Dependencies

Reviewers focus on control flow and may miss an internal `new()` call buried in a helper. Manual review cannot verify that every external dependency was extracted to a constructor parameter.

On larger codebases, dependency chains spread across dozens of files. Relying on human attention alone leaves implicit coupling behind.

## Core Mental Model: Separate Configuration from Use

DI splits object creation into two phases. The composition root assembles concrete implementations; business modules consume interfaces without knowing their origins.

- **Configuration phase** (`main.go`): instantiate `RealGateway`, pass to `NewProcessor`.
- **Usage phase** (`Processor`): validate input, delegate to `gateway.Charge`, return errors.

Only along this path does the business module remain decoupled from provider selection.

## Manual DI Execution Pipeline

The pipeline starts at the composition root, flows through constructor injection, validates input, delegates to the gateway, and propagates errors upward.

```go
// main.go — composition root
realGateway := &di.RealGateway{}
processor := di.NewProcessor(realGateway)
err := processor.ProcessPayment(100)
```

Each step is explicit. No runtime container lookup occurs inside `ProcessPayment`.

## Constructor Injection: NewProcessor Accepts an Interface

The `Processor` struct stores a `PaymentGateway` interface. The constructor takes that interface; the struct never references `RealGateway` directly.

```go
// internal/di/processor.go
type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

Injected dependencies are immutable after construction. Swapping implementations requires only changing the composition root.

## Service Locator Hides Real Dependencies

`BadProcessor` receives a `Container` instead of `PaymentGateway`. The actual dependency becomes a runtime lookup rather than a compile-time contract.

```go
// internal/di/locator.go
type BadProcessor struct {
    container Container
}

func (p *BadProcessor) ProcessPayment(amount int) error {
    gw := p.container.GetPaymentGateway()  // hidden dependency
    // ...
}
```

Tests must configure the container, not inject the mock directly. This adds indirection and reduces readability.

## Constructor Injection vs Service Locator Comparison

| Aspect | Constructor Injection | Service Locator |
|--------|----------------------|-----------------|
| Dependencies | Explicit in constructor | Hidden in container |
| Mocking | Inject mock directly | Mock container + get() |
| Immutability | Field set once at construction | Can change if container changes |
| Complexity | Easy to understand | Extra indirection layer |
| Industry standard | Best practice | Anti-pattern |

The contrast is visible in test code: one passes a `MockGateway`, the other wires a `MockContainer`.

## Three Architectural Decisions in This Lab

The lab makes focused choices to keep the example tractable. These decisions define the structure without over-engineering.

1. **Choose constructor injection** as the default pattern for explicit, testable dependencies.
2. **Avoid service locator in domain code**; use it only to demonstrate the anti-pattern.
3. **Instantiate value objects directly** (`Money`) because they carry no infrastructure dependencies.

These choices preserve clarity while still showing the cost of the alternative.

## Six Test Cases Pass with Race Detector

The test suite isolates the processor from real gateways and verifies happy paths, error propagation, and invalid input.

```
=== RUN   TestProcessor_Success
--- PASS: TestProcessor_Success (0.00s)
=== RUN   TestProcessor_GatewayError
--- PASS: TestProcessor_GatewayError (0.00s)
=== RUN   TestProcessor_InvalidAmount
--- PASS: TestProcessor_InvalidAmount (0.00s)
=== RUN   TestBadProcessor_Success
--- PASS: TestBadProcessor_Success (0.00s)
=== RUN   TestBadProcessor_GatewayError
--- PASS: TestBadProcessor_GatewayError (0.00s)
=== RUN   TestBadProcessor_InvalidAmount
--- PASS: TestBadProcessor_InvalidAmount (0.00s)
PASS
ok      lab16/tests   (cached)
```

All six cases pass with `go test -race`. No shared mutable state exists between tests.

## Evidence Boundary: What This Lab Proves and What It Does Not

Lab 16 demonstrates isolation, error propagation, and input validation on six controlled scenarios. It does not prove production lifecycle management, performance overhead of manual wiring at scale, or empirical defect-rate reductions.

No quantitative measurements of test speed or system stability are included. Claims about DI benefits beyond the test suite require separate empirical validation.

## Trade-offs: Manual Wiring vs Framework DI

Manual wiring adds boilerplate in `main.go` but guarantees compile-time type safety with zero external dependencies. Framework-based DI reduces repetition but introduces reflection, configuration files, and opaque object graphs.

At small scales, manual wiring is often clearer. At very large object graphs, boilerplate can become burdensome. This lab accepts the boilerplate to keep the example self-contained.

## Race Safety in DI Implementation

Go's race detector checks concurrent memory accesses. This DI implementation uses no goroutines or shared state inside `Processor`. Tests construct instances sequentially, so no races exist to detect.

```go
func TestProcessor_Success(t *testing.T) {
    mock := &MockGateway{}
    proc := di.NewProcessor(mock)
    err := proc.ProcessPayment(50)
    if err != nil {
        t.Fatalf("expected nil error, got %v", err)
    }
}
```

Stateless components with injected dependencies naturally avoid race conditions.

## Before vs After: Hard-Coded to Constructor Injection

Before DI, the processor instantiated `RealGateway` internally. After DI, the processor receives `PaymentGateway` through its constructor. The test swaps the real gateway for a mock without touching production code.

| Before (Tight Coupling) | After (Manual DI) |
|-------------------------|-------------------|
| `Processor` creates `RealGateway` internally | `NewProcessor` accepts `PaymentGateway` |
| Tests cannot mock gateway | Tests inject `MockGateway` directly |
| Changing provider requires code edits | Changing provider requires one line in `main.go` |

DI does not eliminate all coupling; it moves the coupling point to the composition root.

## Interface PaymentGateway: Single Contract

The interface defines one method. All implementations share this contract.

```go
// internal/di/gateway.go
type Money struct {
    Amount   int
    Currency string
}

type PaymentGateway interface {
    Charge(m Money) error
}
```

Consumers depend on the interface. Producers implement it. This separation enables substitution without modifying consumers.

## Money: Direct Instantiation of a Value Object

`Money` contains only data. It has no infrastructure dependencies, so it is created directly inside `ProcessPayment`.

```go
m := Money{Amount: amount, Currency: "USD"}
return p.gateway.Charge(m)
```

Injecting a value object through DI would add abstraction without benefit. Value objects belong in the domain layer, not the composition root.

## MockGateway: Test Double for Isolation

`MockGateway` records calls and can be configured to succeed or fail. Tests verify behavior without network access.

```go
// tests/processor_test.go
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

Mocks are structured types that satisfy interfaces. They replace infrastructure without changing production code.

## Composition Root: main.go as Wiring Point

`main.go` is the only file that knows about `RealGateway`. It assembles the object graph and passes dependencies inward.

```go
// cmd/demo/main.go
func main() {
    realGateway := &di.RealGateway{}
    processor := di.NewProcessor(realGateway)
    fmt.Println("--- Running Constructor Injection ---")
    err := processor.ProcessPayment(100)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

Business modules never import concrete implementations. This keeps the composition root centralized and the domain pure.

## Input Validation Before Gateway Delegation

`ProcessPayment` rejects invalid amounts before calling the gateway. This prevents unnecessary external calls and fails fast.

```go
if amount <= 0 {
    return errors.New("invalid amount")
}
```

Validation belongs to the domain layer. Gateways should not enforce business invariants; they execute them.

## Dependency Injection Implementation Checklist

- [ ] Dependencies declared as interfaces, not concrete types
- [ ] Constructor injection used as default for all services
- [ ] Composition root is the only place creating concrete implementations
- [ ] No `container.Get(...)` calls inside domain logic
- [ ] Value objects instantiated directly
- [ ] Each service testable with mocks, no infrastructure setup required
- [ ] Input validation performed before delegating to external dependencies
- [ ] Tests cover happy path, error propagation, and invalid input
- [ ] `go test -race ./...` passes without warnings

These criteria ensure dependencies remain explicit and testable.

## Common Anti-Patterns in Dependency Injection

| Anti-Pattern | Problem |
|--------------|---------|
| Service Locator in domain code | Hides dependencies, complicates testing |
| Abstracting value objects | Adds unnecessary DI indirection |
| Constructor bloat | Indicates too many responsibilities |
| DI overuse | Using DI for pure functions with no state |

Avoiding these pitfalls preserves the benefits of explicit dependency management.

## Independent Reproduction Using Go Commands

The implementation uses only the Go standard library. Run the following commands to verify behavior locally.

```bash
cd labs/16-dependency-injection
go test -v ./...       # run all test cases
go test -race ./...    # verify race safety
go run ./cmd/demo      # execute constructor injection demo
```

Zero external dependencies are required. The test suite validates isolation, error propagation, and input handling.

## Core Principle: Dependencies Must Be Explicit

DI architecture rests on three pillars:

1. **Separation of concerns** — configuration (composition root) is separate from use (business logic).
2. **Program to interfaces** — consumers depend on contracts, not implementations.
3. **Explicit over implicit** — all dependencies are declared in constructors; no magic lookups.

These principles make systems testable, understandable, and evolvable.

## Implementation Modules and Verification Status

Source files: `internal/di/gateway.go`, `internal/di/processor.go`, `internal/di/locator.go`, `cmd/demo/main.go`, `tests/processor_test.go`.

Commit `3d936ba` · Status: APPROVED  
Tests passing: 6 scenarios (success, error propagation, invalid amount)  
Race detector: Clean (`go test -race ./... PASS`)  
Dependencies: 100% Go standard library, zero external modules.

---

Source: [Software Engineering Lab 16](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), README.md, implementation, and tests. Commit snapshot `3d936ba`. Code snippets used within their original package context.
