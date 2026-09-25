# 16 - Dependency Injection

This lab demonstrates Dependency Injection (DI) and Inversion of Control (IoC), contrasting proper constructor injection with the Service Locator anti-pattern.

## Core Concepts Proven

1. **Separation of Configuration from Use:** Object instantiation is externalized (`main.go`), decoupling logic from infrastructure.
2. **Fast, Isolated Unit Testing:** Dependencies are mocked without real network calls (`tests/processor_test.go`).
3. **Constructor Injection:** `NewProcessor` ensures components are fully initialized with explicit dependencies.
4. **Service Locator Anti-Pattern:** `NewBadProcessor` injects a `Container`, hiding real dependencies and coupling the object to framework APIs.
5. **Value Objects Bypass DI:** `Money` is directly instantiated as it lacks behavior tied to external infrastructure.

## Usage

### Run Tests
```bash
go test -race ./...
```

### Run Demo
```bash
go run ./cmd/demo
```
