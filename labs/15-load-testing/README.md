# Lab 15: Load Testing and Bottleneck Analysis

Demonstration of load testing principles, percentile analysis (P50, P95, P99), and resource exhaustion detection.

## Structure
- `cmd/demo`: Entry point running comparative Smoke vs. Stress test scenarios.
- `internal/server`: Mock service with constrained connection pool capacity to demonstrate saturation.
- `internal/loadtest`: Built-in concurrency runner and percentile calculator.
- `tests`: Automated integration tests validating metric accuracy and stress-induced latency growth.
- `engineering/`: Design, implementation notes, and execution results.

## Requirements
- Go 1.22+

## Running the Demo
```bash
go run ./cmd/demo
```

## Running Tests
```bash
go test -v ./...
go test -race ./...
```
