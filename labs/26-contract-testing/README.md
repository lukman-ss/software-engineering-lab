# Lab 26: Contract Testing (Consumer-Driven Contracts)

Practical implementation of Consumer-Driven Contract (CDC) testing and CI/CD deployment gates in Go.

## Overview

Contract testing ensures services can communicate compatibly without maintaining brittle end-to-end environments. This lab demonstrates:
1. **Consumer-Driven Contracts**: Consumer declares minimal required schema/interactions.
2. **Provider CI Gate Verification**: Provider verifies implementation against consumer contracts before deployment.
3. **Breaking Change Detection**: Detects enum casing changes, field renames, and primitive type mutations.
4. **Safe API Evolution**: Preserving V1 contract compatibility while exposing V2 schemas.

## Project Structure

```text
labs/26-contract-testing/
├── cmd/
│   └── demo/main.go            # Executable multi-stage demo
├── internal/
│   ├── consumer/client.go      # Mobile consumer client & contract builder
│   ├── contract/verifier.go    # CDC engine and verification runner
│   ├── model/order.go          # Domain models and DTOs (V1, Breaking, V2)
│   └── provider/server.go      # HTTP handlers for V1, Breaking, and Dual providers
├── tests/
│   └── contract_test.go        # Unit, integration, and race detector tests
├── engineering/
│   ├── 01-design.md
│   ├── 02-implementation-notes.md
│   └── 03-execution-result.md
├── go.mod
└── README.md
```

## Running the Lab

### Run Tests and Race Detector

```bash
cd labs/26-contract-testing
go test -v ./...
go test -race ./...
```

### Run Demo

```bash
cd labs/26-contract-testing
go run ./cmd/demo
```
