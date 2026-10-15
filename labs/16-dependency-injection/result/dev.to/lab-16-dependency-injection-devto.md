---
title: "Dependency Injection in Go: Constructor Injection vs Service Locator Trade-offs"
published: false
description: "Analyze Constructor Injection versus Service Locator in Go: object construction boundaries, mock isolation, nil safety, and test assertions."
tags: go, testing, architecture, softwareengineering
---

# Dependency Injection in Go: Constructor Injection vs Service Locator Trade-offs

A software component that directly instantiates its own dependencies couples business domain rules tightly to volatile infrastructure details. When an application service constructs third-party clients or network drivers within its internal methods, swapping out external vendors requires widespread modifications across the core codebase. Furthermore, unit testing becomes slow and fragile because tests cannot run without invoking real external systems or relying on complicated monkey patching mechanisms.

Lab 16 examines this architectural tension through a concrete Go payment processing implementation. By contrasting constructor injection against the service locator pattern, the lab demonstrates how inverting control over dependency construction simplifies test isolation while highlighting the operational boundaries that engineers must manage.

> When a component builds its own dependencies, replacing them requires editing the component itself.

## Coupling Locks the Gateway Inside Domain Logic

When the `Processor` component directly constructs `&RealGateway{}` internally, it tightly binds transaction execution to production network infrastructure. Swapping out a payment vendor or switching to an alternative provider forces engineers to alter domain source files rather than modifying a single configuration entrypoint. Under this design, automated unit testing becomes impractical because every test execution inherently attempts to establish live network connections.

## Constructors Without Interfaces Prevent Substitution

A naive constructor that takes zero arguments hides its external requirements from callers and test harnesses alike. Refactoring the constructor to accept a `PaymentGateway` interface parameter requires zero modifications to the core arithmetic inside `ProcessPayment`. However, this small structural change immediately shifts the instantiation decision upward to the composition root, allowing test suites to supply fast in-memory doubles without altering business logic.

## Inversion of Control Represents Control Direction

Inversion of Control (IoC) is a broad architectural paradigm where custom application logic receives its execution flow and dependencies from an external supervisor rather than driving the program lifecycle itself. Dependency Injection represents a specific realization of IoC focused strictly on dependency assembly. In this lab, the entrypoint at `cmd/demo/main.go` acts as the composition root, taking full ownership of concrete struct initialization before supplying those references to domain components that depend solely on abstract behavioral contracts.

## Request Execution Flow from Assembly to Settlement

The execution lifecycle of a payment request proceeds through four distinct operational phases: composition, invocation, validation, and delegation. During application boot, `main.go` creates an instance of `RealGateway` and injects it into `NewProcessor`. At runtime, a client calls `ProcessPayment(100)`. The processor first validates that the payment amount is strictly positive, builds an immutable `Money` value object directly on the stack, and finally delegates settlement to `gateway.Charge(m)`, returning any resultant errors to the caller.

## Constructors Make Dependencies Explicit in Signatures

Public constructor functions make external infrastructure needs explicit right within their type signatures. By accepting `PaymentGateway` as a required parameter and storing it in an unexported struct field, `Processor` communicates its operational requirements transparently to every consumer. This contract prevents callers from instantiating an unconfigured processor while ensuring internal methods remain decoupled from specific vendor implementations.

```go
func NewProcessor(g PaymentGateway) *Processor {
    return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
    if amount <= 0 { return errors.New("invalid amount") }
    m := Money{Amount: amount, Currency: "USD"}
    return p.gateway.Charge(m)
}
```

## Service Locators Conceal True Dependencies

In contrast to explicit injection, `BadProcessor` accepts a generic `Container` interface rather than an explicit `PaymentGateway`. During method execution, `ProcessPayment` retrieves the gateway dynamically by invoking `p.container.GetPaymentGateway()`. While both approaches produce identical terminal output in standard execution paths, the service locator pattern hides the component's true requirements and forces testing harnesses to construct artificial container structures rather than passing direct mocks.

## The Hidden Indirection Tax of Service Locators

While service locators initially appear convenient by keeping constructor parameter signatures short, they introduce an ongoing maintenance tax. Domain components become tightly coupled to container access interfaces, transforming clear compile-time dependency requirements into runtime lookup conventions that obscure the architectural graph from static analysis tools and code reviewers.

| Architectural Aspect | Constructor Injection | Service Locator Anti-Pattern |
|---|---|---|
| Test Isolation | Direct mock injection via interface parameter | Requires mocking container plus nested gateway |
| Signature Transparency | Explicit contracts visible in public signatures | Dependencies hidden behind generic container API |
| Vendor Swapping | Single line adjustment in composition root | Single line adjustment in composition root |
| Nil Pointer Exposure | Runtime panic if constructor lacks explicit nil guard | Runtime panic if container lookup returns nil |

## Isolated Test Harness Eliminates Network Calls

The test harness located in `tests/processor_test.go` exercises business logic without opening network sockets or calling third-party APIs. The test suite utilizes `MockGateway`, an in-memory double that records received transactions into a struct field and provides a `ShouldFail` boolean toggle. This setup allows unit tests to deterministically verify happy path settlement, boundary validation, and error propagation paths in complete isolation.

## Six Test Cases Pass with Clean Race Detection

Executing the full automated test suite confirms that all six test scenarios pass with zero failures. Both `Processor` and `BadProcessor` are tested across successful processing, simulated downstream provider failure, and non-positive input amounts. The Go race detector verifies that sequential test execution completes cleanly without detecting data contention issues across internal pointers.

```sh
$ go test -race -v ./tests/
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
PASS: ok lab16/tests (cached)
```

## Local Execution Speeds Are Not Production SLOs

The test suite reports zero-second execution times because mock operations operate entirely on local CPU memory without external I/O. While these metrics confirm that the dependency injection harness introduces negligible computational overhead in automated test environments, they should never be confused with real-world production service level objectives, which are dominated by remote network latency, TLS handshakes, and database contention.

| Performance Metric | Measured Lab Result | Engineering Interpretation | Operational Reality Boundary |
|---|---|---|---|
| Test Duration | 0.00s per test case | In-memory execution speed | Excludes network, DNS, and TLS latency |
| Data Race Warnings | 0 detected races | Sequential test paths clean | Does not validate concurrent goroutine access |
| Mock Recording | Value captured correctly | Validates parameter handoff | Does not test network timeouts or drops |
| Demo Execution | Stdout log verified | Confirms manual wiring | Manual run does not guarantee regression safety |

## Balancing Boilerplate Against Test Isolation

Manual dependency injection requires software engineers to write explicit wiring code within the composition root, manually instantiating dependencies and passing them down the object hierarchy. This slight increase in boilerplate pays significant dividends by eliminating reflection magic, preserving compile-time type verification, and enabling developers to mock any external dependency with minimal effort.

## Defensive Nil Guards Protect Against Runtime Panics

Although constructor injection guarantees that callers cannot omit required parameters, the Go compiler permits passing an uninitialized `nil` value for any interface type. If a caller mistakenly passes `nil` into `NewProcessor`, the constructor executes successfully, only to trigger a fatal nil pointer dereference when `ProcessPayment` attempts to invoke `gateway.Charge()`. Defensive production code should include explicit nil checks inside constructors.

```go
func NewProcessor(g PaymentGateway) *Processor {
    if g == nil {
        return nil // Or return a sentinel error / panic based on team standards
    }
    return &Processor{gateway: g}
}
```

## Before and After: Structural Decoupling

Comparing the tightly coupled model to constructor injection highlights that core domain business logic remains virtually untouched. The architectural refactoring exists entirely at the package boundary: substituting a concrete struct field for an abstract interface decouples implementation from consumption, transforming a rigid monolith into an adaptable, testable module.

## Boundaries and Deliberate Non-Goals of the Lab

This laboratory concentrates strictly on structural dependency inversion mechanics. It deliberately avoids exploring advanced lifecycle management concerns (such as singleton versus transient object scopes), reflection-based automated container frameworks (such as Wire, Dig, or Fx), or the concurrency synchronization primitives needed when sharing stateful gateways across multiple goroutines.

## Edge Cases: Unchecked Nil Dependencies

When a consumer supplies an uninitialized `nil` reference to `NewProcessor`, the compilation phase passes because `nil` satisfies any Go interface signature. The resulting processor struct represents an invalid state that silently awaits execution. Robust production systems must pair dependency injection with defensive assertions at the boundary to catch configuration errors before requests are accepted.

## Concrete Invariants Proven by Automated Tests

The test suite provides conclusive evidence for three foundational runtime behaviors: valid payment requests record precise money values in mock memory, downstream payment failures propagate immediately back to the caller without being swallowed, and non-positive transaction amounts abort instantly before external network calls can be initiated.

## Unproven Production Capabilities

A completely green unit test suite does not guarantee production reliability. The tests in this lab do not simulate network drops, context deadline cancellations, circuit breaker state transitions, or high-throughput concurrency contention. Engineers must supplement dependency injection with integration tests and load testing to establish true production readiness.

## Operational Risks of Unchecked Object Lifecycles

Implementing dependency injection without standardized lifecycle guidelines can introduce operational risks. Inadvertently instantiating new processor or gateway instances per incoming HTTP request can exhaust memory and connection pools. Additionally, stateful dependencies—such as gateways maintaining authentication tokens—require thread-safe synchronization to prevent data races during concurrent access.

## Implementation Checklist for Safe Dependency Injection

Engineering teams adopting dependency injection should adhere to clear structural conventions to avoid accidental overengineering:

- [ ] Declare infrastructure requirements as interfaces rather than concrete struct types.
- [ ] Utilize constructor injection as the default mechanism to keep all dependencies visible.
- [ ] Restrict concrete instantiation logic strictly to composition roots like `main.go`.
- [ ] Instantiate pure value objects (such as `Money`) directly without routing through DI containers.
- [ ] Execute input boundary validations prior to delegating calls to external interfaces.
- [ ] Execute `go test -race ./...` to guarantee that tested execution paths remain race-free.

## Anti-Pattern: Embedding Containers into Domain Logic

Injecting a container or service locator into core business domain entities violates the principle of encapsulation. This anti-pattern forces domain services to depend upon container retrieval conventions, obscures real dependencies from code reviewers, and unnecessarily complicates unit testing by requiring developers to construct mock containers rather than simple mock interfaces.

## Reproducing the Verification Pipeline Locally

Developers can reproduce all architectural assertions and execute the interactive demonstration using the standard Go compiler without installing any external libraries or tools:

```sh
cd labs/16-dependency-injection
go test -race -v ./tests/     # Execute all 6 isolated unit test assertions
go run ./cmd/demo            # Run both Constructor and Service Locator demos
go vet ./...                 # Verify that interfaces and structs adhere to conventions
```

## Constructors Make Dependencies Transparent

Dependency injection is fundamentally about control flow and interface decoupling, not complex third-party frameworks. By making dependencies explicit in constructor parameters, software engineers establish clear contracts, unlock fast and deterministic test isolation, and maintain maintainable architectures that adapt smoothly as infrastructure requirements evolve over time.

## Source

[Software Engineering Lab 16 - Dependency Injection in Go](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection), snapshot commit `d7fb619`. Verified files: `internal/di/processor.go`, `internal/di/locator.go`, `internal/di/gateway.go`, `cmd/demo/main.go`, and `tests/processor_test.go`.
