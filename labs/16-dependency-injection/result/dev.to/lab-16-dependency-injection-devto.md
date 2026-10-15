---
title: "Dependency Injection: Separating Configuration from Use in Go"
published: false
description: "Explore Dependency Injection in Go: constructor injection versus Service Locator anti-pattern, interface-based mocking, and what the test suite actually proves."
tags: go, testing, architecture, softwareengineering
---

# Dependency Injection: Separating Configuration from Use in Go

Hard-coded dependencies make isolated testing impossible and force full-service rebuilds when an external implementation changes. Lab 16 contrasts two designs for the same payment logic — `Processor` with constructor injection and `BadProcessor` using a Service Locator — to show how wiring location determines testability.

> DI separates configuration from use; it does not eliminate the need for validation or error handling.

The lab uses a minimal manual-wiring setup in pure Go, no DI framework. Six unit tests in `tests/processor_test.go` cover happy path, error propagation, and invalid input for both patterns, including a race-detector run.

## The problem: hard-coded dependencies

A component that constructs its own dependencies internally couples to a concrete type:

```go
// Before (not shown in lab, but illustrative of the anti-pattern)
type Processor struct{}
func (p *Processor) ProcessPayment(amount int) error {
    gateway := &RealGateway{} // hard-coded dependency
    return gateway.Charge(money)
}
```

`RealGateway` is locked inside `Processor`. Mocking it in a unit test requires either changing the source or accepting integration-level flakiness. Every change to the payment provider forces edits across every consumer class.

## Why manual refactoring leaves leaks

Rewriting constructors by hand is error-prone. A developer may update `NewProcessor` to accept an interface while leaving a stray instantiation at the old call site. A missing binding in a registry leaves nothing but a nil-pointer panic at runtime — Go's compile-time type system cannot guard against it when wiring is scattered.

## Mental model: DI as a dependency graph

Dependencies form a directed acyclic graph (DAG): nodes are interfaces or concrete types, edges are providers or injection paths. The composition root sits at the graph's entry point, resolving and wiring everything before any business logic runs.

The lab's graph is small by design:

```
main.go (composition root)
      │
      ├── RealGateway ──implements──▶ PaymentGateway
      │         │
      │         ▼ inject via NewProcessor
      │     Processor
      │         │
      │         └──▶ Money{Amount, Currency}  (value object, direct)
      │
      └── SimpleContainer ──implements──▶ Container
                │
                ▼ inject via NewBadProcessor
            BadProcessor
                │
                └──▶ Container.GetPaymentGateway() ──▶ PaymentGateway
```

## Pipeline: from wiring to execution

1. **Composition root** (`cmd/demo/main.go`) instantiates concrete types.
2. **Consumer constructors** receive interfaces, never concrete types.
3. **Consumer methods** validate input first, then delegate to injected dependencies.
4. **Tests** substitute `MockGateway` for `RealGateway` without touching `Processor`.

No reflection, no service locator lookup inside the domain layer.

## Code: the Provider / Interface contract

```go
type PaymentGateway interface {
    Charge(m Money) error
}

type Processor struct {
    gateway PaymentGateway
}

func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}
```

`NewProcessor` guarantees the `Processor` is created with a valid, non-nil dependency. The field is unexported, so the injected gateway cannot be swapped after construction — immutability at birth.

## Failure path: invalid input

`ProcessPayment` validates before delegating:

```go
func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

`TestProcessor_InvalidAmount` sends `-10`. It expects a non-nil error and `mock.ChargedMoney.Amount == 0`, confirming the gateway was never called. The same guard exists in `BadProcessor`, proving the validation belongs to the domain layer, not the DI infrastructure.

## DI vs No DI (testability & startup)

| Metric | DI (lab) | No DI (hard-coded) |
|---|---|---|
| Mockability | High — pass `MockGateway` directly | Low — `RealGateway` is internal |
| Test isolation | Full — zero network calls | None — integration test required |
| Startup wiring | Explicit in composition root | Scattered across consumers |
| Boilerplate | Small — one constructor per service | None per consumer, but rigid |

The table reflects what this lab demonstrates. It does not measure production startup latency or framework-generated wiring time.

## Three key decisions when adopting DI

1. **Bind per request or singleton?** — The lab creates objects fresh on each call. Lifecycle management (singleton, scoped, transient) is discussed in the README but not implemented.
2. **Prefer interface over concrete.** — `PaymentGateway` over `*RealGateway`.
3. **Inject only services that interact with external infrastructure.** — `Money` is a plain value object; instantiate it directly.

## Demo output

```
go run ./cmd/demo

--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD
```

Both paths produce the same observable output. The difference is structural, not behavioral.

## What the tests do NOT prove

- **Lifecycle management.** Singleton, scoped, and transient scopes are not implemented or tested.
- **Performance overhead.** No benchmark compares manual wiring against a DI framework.
- **Concurrency beyond the race detector.** `go test -race` passes, but the processor has no mutable shared state to contend with — concurrency safety here follows from immutability, not locking.
- **Production-scale wiring.** The demo wires two services manually; a large object graph with hundreds of bindings was not evaluated.

## Trade-offs: cleaner code vs added indirection

**Benefits:**
- Testable without infrastructure setup.
- Swappable implementations without touching domain code.
- Compile-time type safety on constructor signatures.

**Costs:**
- Composition-root boilerplate grows with service count.
- Behavior is separated from construction, making tracing harder.
- Framework adoption (wire, dig) adds a dependency with its own learning curve.

## Before vs After: Constructor Injection

```go
// Before
svc := &serviceImpl{}   // hard-coded, no mock path

// After
svc := di.NewProcessor(realGateway)  // interface, injectable
```

The consumer (`Processor`) no longer knows what concrete gateway it talks to.

## Test Coverage: the six cases

| Test | Pattern | Input | Asserts |
|---|---|---|---|
| `TestProcessor_Success` | Constructor | 50 USD | `ChargedMoney` = 50 USD |
| `TestProcessor_GatewayError` | Constructor | ShouldFail=true | non-nil error returned |
| `TestProcessor_InvalidAmount` | Constructor | -10 | error + gateway never called |
| `TestBadProcessor_Success` | Service Locator | 75 USD | `ChargedMoney` = 75 USD |
| `TestBadProcessor_GatewayError` | Service Locator | ShouldFail=true | non-nil error returned |
| `TestBadProcessor_InvalidAmount` | Service Locator | -5 | error + gateway never called |

All six pass with `go test -v ./...` and `go test -race ./...` reports no data races.

## Error propagation: identical in both patterns

`TestProcessor_GatewayError` and `TestBadProcessor_GatewayError` both set `ShouldFail: true` on `MockGateway`. Both tests expect `err != nil`, confirming that DI does not swallow or transform gateway errors — it preserves the caller's ability to handle them.

## Performance: no benchmark evidence

The README's comparative metrics are illustrative, not measured. This lab does not include a benchmark comparing constructor-injection latency against direct `new` or container-lookup latency. Do not cite the "15 ms vs 40 ms" table as empirical data.

## Reflection safety: not used

The lab avoids reflection-based DI entirely. Manual wiring provides full compile-time type checking — `NewProcessor(realGateway)` rejects a type mismatch at build time. No runtime type assertion or `reflect` usage appears in the implementation.

## Checklist: implementing DI correctly

- [ ] Dependencies declared as interfaces, not concrete types.
- [ ] Constructor injection used by default; all dependencies explicit in signature.
- [ ] Composition root is the only place that instantiates concrete implementations.
- [ ] No `container.Get(...)` calls inside domain logic.
- [ ] Value objects instantiated directly, not through the container.
- [ ] Every service testable with a mock, no infrastructure setup required.
- [ ] Input validation runs before delegating to external dependencies.
- [ ] Test suite covers happy path, error propagation, and invalid input.
- [ ] `go test -race ./...` passes with no warnings.

## Anti-pattern: Service Locator abuse

`BadProcessor` accepts a `Container` interface and calls `GetPaymentGateway()` inside `ProcessPayment`:

```go
func (p *BadProcessor) ProcessPayment(amount int) error {
    if amount <= 0 {
        return errors.New("invalid amount")
    }
    m := Money{Amount: amount, Currency: "USD"}
    return p.container.GetPaymentGateway().Charge(m)
}
```

The actual dependency (`PaymentGateway`) is hidden behind the container. The consumer depends on the container API, not the service contract. PSR-11 states users **SHOULD NOT** pass a container into an object — a strong recommendation, not an absolute prohibition.

Problems with Service Locator:
- Dependencies are invisible from the constructor signature.
- Each test must stub out the container, not the service.
- Refactoring the service name breaks the locator call silently only at runtime if types diverge.

## Practical command reference

```bash
go run ./cmd/demo          # run both patterns end-to-end
go test -v ./...           # run all six unit tests
go test -race ./...        # verify no data races
```

## Long-term principle

> Depend on abstractions, not concretions.

Constructor injection achieves this by making every dependency explicit at construction time. The service becomes a stable consumer of an interface; the implementation can change, be swapped, or be replaced with a test double without touching the consumer.

DI is a structural discipline, not a performance optimization. It makes systems cheaper to test and cheaper to evolve — but it does not fix bad validation, missing error handling, or architectural problems that predate the wiring.

## Repository and source map

- [Lab 16](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection)
- [README.md](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/16-dependency-injection/README.md)
- [internal/di/gateway.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/16-dependency-injection/internal/di/gateway.go) — `PaymentGateway` interface, `Money` value object, `RealGateway`
- [internal/di/processor.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/16-dependency-injection/internal/di/processor.go) — `Processor` with constructor injection
- [internal/di/locator.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/16-dependency-injection/internal/di/locator.go) — `Container` interface and `BadProcessor` Service Locator anti-pattern
- [tests/processor_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/16-dependency-injection/tests/processor_test.go) — six unit tests covering both patterns
- [cmd/demo/main.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/16-dependency-injection/cmd/demo/main.go) — composition root demo
