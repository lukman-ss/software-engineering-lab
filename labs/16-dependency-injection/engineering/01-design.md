# Engineering Design

Target Lab: labs/16-dependency-injection
Research Status: APPROVED

## Concept To Prove
Demonstrate Dependency Injection, specifically comparing Constructor Injection (best practice) against Service Locator (anti-pattern). Verify that unit testing is trivial via mock injection and value objects are instantiated directly.

## Expected Behavior
- Constructor injection mandates all dependencies explicitly via interface.
- Service locator hides dependencies by passing a container instead.
- Simple domain/value objects (Money) do not use DI.

## Failure Scenario
- Passing invalid state (negative amount) rejects the payment without calling external services.
- Mocking a gateway error propagates back cleanly.

## Success Criteria
- Mock gateway records calls successfully in tests.
- Demo prints expected messages representing the real infrastructure calls.

## Architecture
- `di.PaymentGateway`: interface abstracting the external API.
- `di.RealGateway`: simulated concrete implementation.
- `di.Processor`: service using constructor injection.
- `di.BadProcessor`: service using the Service locator anti-pattern.
- `di.Money`: value object decoupled from infrastructure.

## Components
- `internal/di/gateway.go`
- `internal/di/processor.go`
- `internal/di/locator.go`
- `cmd/demo/main.go`

## Test Strategy
- Unit tests pass a `MockGateway` into `Processor` to assert call properties and test failures without network calls.

## Execution Plan
1. Implement the DI components and Service Locator.
2. Build demo CLI wiring real components.
3. Write unit tests utilizing mock implementations.
4. Run `go test` and `go run`.

## Implementation Decisions
- Kept the "container" extremely simple (just an interface returning the gateway) to isolate the anti-pattern behavior. No full-fledged DI framework (like `dig` or `wire`) added since YAGNI dictates standard Go is sufficient to prove the concepts.
