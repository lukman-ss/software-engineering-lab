# Database Connection Pooling Lab

This lab demonstrates database connection pool behavior, client-side resource sizing, and the failure modes associated with improper connection management.

## Components

- `internal/pool/mockdb.go`: Implements a Go `database/sql/driver` to simulate database server constraints (`max_connections`) and connection establishment latency.
- `internal/pool/service.go`: Provides services executing safe operations vs unsafe operations (holding a DB connection during external I/O).
- `tests/pool_test.go`: Automated test suite covering direct overhead, pool exhaustion, connection leaks, and concurrent operations.
- `cmd/demo/main.go`: Interactive CLI demo demonstrating connection overhead, rejection, and pool starvation.

## Running Tests

Execute the automated test suite:

```bash
go test -v ./...
```

Run race detection:

```bash
go test -race -v ./...
```

## Running the Demo

Execute the demo application:

```bash
go run ./cmd/demo
```
