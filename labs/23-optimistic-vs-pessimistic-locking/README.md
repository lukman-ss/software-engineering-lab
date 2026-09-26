# Lab 23: Optimistic vs Pessimistic Locking & Atomic Updates

This lab demonstrates how concurrent processes cause the silent "lost update" anomaly when using naive read-modify-write patterns, and proves the three primary remediation strategies:
1. **Pessimistic Locking** (`SELECT ... FOR UPDATE`)
2. **Optimistic Locking** (Version check in `WHERE` clause with retry logic)
3. **Atomic Single-Statement Operations** (`UPDATE ... SET stock = stock - N WHERE stock >= N`)

---

## Directory Structure

```text
labs/23-optimistic-vs-pessimistic-locking/
├── cmd/
│   └── demo/
│       └── main.go           # Executable demo comparing all concurrency approaches
├── internal/
│   └── inventory/
│       ├── model.go          # Data structures and domain error definitions
│       ├── service.go        # Business logic implementing each locking strategy
│       └── store.go          # Simulated storage engine with row locks & version guards
├── tests/
│   └── locking_test.go       # Automated concurrency and invariant tests
├── engineering/
│   ├── 01-design.md
│   ├── 02-implementation-notes.md
│   └── 03-execution-result.md
├── go.mod
└── README.md
```

---

## How to Run

### Run Unit & Concurrency Tests
```bash
go test -v ./...
```

### Run Race Detector
```bash
go test -race ./...
```

### Run Demonstration
```bash
go run ./cmd/demo
```
