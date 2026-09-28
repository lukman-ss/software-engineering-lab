# Lab 34: Chaos Engineering & Fault Injection

Runnable implementation of Chaos Engineering principles in Go 1.22+.

## Overview

This lab proves key Chaos Engineering principles:
1. Steady-state health metric monitoring.
2. Injected downstream service latency and errors.
3. Resilience validation via Circuit Breakers and Fallback mechanisms.
4. Blast radius control with automated experiment abort upon metric degradation.

## Requirements

- Go 1.22+

## Running Tests

Run all unit and concurrency tests:

```bash
go test ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

## Running the Demo

Execute the interactive demonstration:

```bash
go run ./cmd/demo
```

## Structure

- `internal/fault`: Latency and error injection primitives.
- `internal/circuitbreaker`: Circuit Breaker state machine (Closed, Open, Half-Open).
- `internal/monitor`: Steady-state health tracking.
- `internal/experiment`: Chaos experiment orchestration with auto-abort mechanism.
- `cmd/demo`: Executable demonstration script.
- `tests`: Unit and race condition tests.
- `engineering/`: Design, implementation notes, and execution results.
