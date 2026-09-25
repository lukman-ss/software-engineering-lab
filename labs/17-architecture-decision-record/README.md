# Lab 17: Architecture Decision Record (ADR)

This lab demonstrates an automated validation and linting mechanism for Architecture Decision Records (ADRs) to enforce structural invariants, monotonic numbering, and supersession lineage.

## Implementation Details

- `internal/adr`:
  - `models.go`: Defines the core `Record` structure and valid statuses (`Proposed`, `Accepted`, `Superseded`, `Deprecated`, `Rejected`).
  - `parser.go`: Parses markdown ADR text into structured data.
  - `linter.go`: Validates monotonic numbering and referential integrity of supersession links concurrently.
- `cmd/demo/main.go`: Demonstrates parsing and validating an ADR sequence representing the architectural progression from a Modular Monolith to Microservices.
- `tests/`: Contains automated tests verifying parser robustness and linter rule enforcement.

## Running Tests

Run the test suite:

```bash
go test -v ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

## Running Demo

Execute the demonstration:

```bash
go run ./cmd/demo
```
