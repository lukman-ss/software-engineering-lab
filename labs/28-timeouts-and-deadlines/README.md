# Timeouts & Deadlines Laboratory

This lab implements and demonstrates distributed system resilience patterns including context deadline propagation, timeout budgets, exponential backoff with full jitter, circuit breaking, and idempotency deduplication.

## Components Implemented

1. **`internal/deadline`**: Context deadline propagation and execution within explicit time budgets.
2. **`internal/retry`**: Exponential backoff with full jitter to avoid synchronized retry storms.
3. **`internal/circuit`**: State machine (`CLOSED`, `OPEN`, `HALF_OPEN`) preventing cascading calls to failing services.
4. **`internal/idempotency`**: In-memory deduplication store preventing double execution during retries.

## Execution

### Run Tests
```bash
go test ./...
go test -race ./...
```

### Run Demo
```bash
go run ./cmd/demo
```
