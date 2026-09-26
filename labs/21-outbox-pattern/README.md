# Lab 21: Transactional Outbox Pattern

A demonstration of the Transactional Outbox pattern implemented in Go, proving atomic persistence between domain entities and event logs, decoupled polling relay dispatch to a message broker, and downstream consumer idempotency.

## Architecture

The project consists of:
- `internal/outbox/db.go`: In-memory transactional database simulating `BeginTx`, `Commit`, and `Rollback` across `orders` and `outbox` records.
- `internal/outbox/broker.go`: Thread-safe mock message broker simulating publish failures and event reception.
- `internal/outbox/service.go`: Business logic comparing naive dual-write vs atomic outbox writes.
- `internal/outbox/relay.go`: Asynchronous polling worker querying pending outbox records and dispatching them to the broker.
- `internal/outbox/consumer.go`: Subscriber enforcing idempotency through event ID tracking.

## Running Tests

Run tests and race detector:

```bash
go test ./...
go test -race ./...
```

## Running Demo

Execute the end-to-end demonstration:

```bash
go run ./cmd/demo
```
