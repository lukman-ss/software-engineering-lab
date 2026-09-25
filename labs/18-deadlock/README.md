# Lab 18: Deadlock Simulation and Prevention

This lab demonstrates deadlock dynamics in concurrent systems, illustrating why circular wait happens, how systems abort transactions (deadlock victims), how **Lock Ordering** prevents deadlocks, and how **Application Retries** recover aborted transactions.

## Structure

- `cmd/demo/main.go`: End-to-end runnable demo showcasing the three scenarios.
- `internal/bank/account.go`: Account model with channel-based lock acquisition and context timeout.
- `internal/transfer/transfer.go`:
  - `TransferNaive`: Acquires locks based on caller input sequence. Can deadlock.
  - `TransferOrdered`: Acquires locks based on account ID sort order. Completely prevents deadlocks.
  - `TransferWithRetry`: Retries `TransferNaive` upon encountering `ErrDeadlock`.
- `tests/transfer_test.go`: Automated tests for all core claims.

## How to Run

### Tests
```bash
go test -v ./...
```

### Race Detector
```bash
go test -race ./...
```

### Run Demo
```bash
go run ./cmd/demo
```
