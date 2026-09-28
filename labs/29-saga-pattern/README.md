# Lab 29: Saga Pattern

Demonstration of distributed transaction management using the Saga Pattern in Go, covering both Orchestration and Choreography models, compensating transactions (LIFO rollback), idempotency keys, and semantic locking countermeasures.

## Components
- `internal/saga/orchestrator.go`: Centralized saga coordinator managing forward steps and LIFO compensation.
- `internal/saga/choreography.go`: Decoupled event bus for choreographing saga steps across independent services.
- `internal/services/services.go`: Domain services (Order, Payment, Inventory) with local state, idempotency, and semantic locks.
- `cmd/demo/main.go`: Console demo showing happy path and failure compensation.
- `tests/saga_test.go`: Suite of unit, rollback, idempotency, semantic lock, and concurrency tests.

## Running Tests

```bash
go test -v ./...
go test -race ./...
```

## Running Demo

```bash
go run ./cmd/demo
```
