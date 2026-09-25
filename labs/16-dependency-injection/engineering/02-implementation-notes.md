# Implementation Notes

## Files Added
- `cmd/demo/main.go`
- `internal/di/gateway.go`
- `internal/di/processor.go`
- `internal/di/locator.go`
- `tests/processor_test.go`
- `README.md`

## Core Design Decisions
- Interfaces (`PaymentGateway`) live close to the components consuming them.
- `Money` was explicitly kept as a value object without interfaces to prove Finding 5.

## Implementation-Specific Choices
- Opted for manual dependency injection in `main.go` over using a framework. Go's simplicity favors manual wiring (constructor injection) for most cases.
- The `Container` interface is a minimal proxy to highlight the Service Locator flaw without requiring heavy dependencies.

## Known Limitations
- Does not demonstrate reflection-based or code-gen based DI containers (e.g., `google/wire`), as the structural concepts apply regardless of automation tool.

## Trade-offs
- Manual dependency injection requires slightly more boilerplate in `main.go` compared to auto-wiring magic, but guarantees type-safety and compile-time validation.

## What Is Demonstrated
- Constructor Injection.
- Service Locator anti-pattern.
- Isolated Mock Testing.
- Value Object instantiation.

## What Is Not Demonstrated
- Advanced DI Container lifecycle management (Singleton vs Scoped vs Transient).
