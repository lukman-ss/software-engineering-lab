# Database Constraints & Data Integrity

This lab demonstrates database-level constraints as the foundation for data integrity, race condition prevention, and schema invariants.

## Implemented Constraints

1. **NOT NULL (`23502`)**: Ensures required columns (e.g. `email`, `username`, `user_id`) cannot accept null/empty values.
2. **CHECK Constraints (`23514`)**: Evaluates boolean predicate logic on row data (e.g. `age >= 18`, `status IN (...)`, `total_cents > 0`).
3. **UNIQUE Constraints (`23505`)**: Enforces single occurrence across rows and prevents concurrent read-then-write race conditions.
4. **FOREIGN KEY Constraints (`23503`)**: Enforces referential integrity preventing orphan rows.
5. **PARTIAL UNIQUE INDEX**: Demonstrates conditional uniqueness (`WHERE deleted_at IS NULL`) enabling soft delete re-registration while maintaining active uniqueness.

## Running Tests

Run the standard test suite:

```bash
go test -v ./...
```

Run tests with the Go race detector:

```bash
go test -race ./...
```

## Running the Demo

Execute the interactive demo:

```bash
go run ./cmd/demo
```
