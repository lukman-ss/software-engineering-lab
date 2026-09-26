# Zero-Downtime Deployment Lab

This lab demonstrates the architectural patterns necessary to achieve zero-downtime deployment by coordinating traffic routing, graceful termination, backward-compatible database schema changes, and cooperative worker shutdown.

## Components

The lab is divided into a background queue worker, an HTTP server, and an in-memory database:

- **Database (`internal/db`)**: Demonstrates the "Expand and Contract" pattern (Parallel Change). The mock storage supports writing dual schema versions (legacy `Name` and modern `FirstName`/`LastName`), and transparent fallback logic when reading records created by differing versions of the application.
- **Server (`internal/server`)**: Exposes Liveness and Readiness probes. When a shutdown signal is sent, it executes a configurable `preStop` delay to simulate load balancer detachment latency, then performs a graceful shutdown ensuring any in-flight requests complete before termination.
- **Worker (`internal/worker`)**: A background daemon that pulls jobs from a queue. Upon receiving a shutdown signal, it stops pulling new jobs but continues processing the current active job until completion.
- **Demo (`cmd/demo`)**: A CLI orchestrator that wires these components together, simulates startup initialization, executes in-flight workloads, and sends a termination signal to demonstrate zero-downtime draining behavior.

## Running the Demo

```bash
cd labs/20-zero-downtime-deployment
go run ./cmd/demo
```

## Running the Tests

```bash
cd labs/20-zero-downtime-deployment
go test -v ./...
go test -race ./...
```
