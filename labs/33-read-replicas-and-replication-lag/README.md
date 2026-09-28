# Lab 33: Read Replicas and Replication Lag

This lab implements and tests patterns for managing database replication lag when using read replicas for scaling, ensuring session guarantees (such as read-your-own-writes) while maximizing read offloading.

## Concepts Implemented
- **Asynchronous vs Synchronous Replication (`remote_apply`)**
- **Stale Read Detection & Mitigation**
- **Time-Based Sticky Routing** (routing reads to primary for $T_{sticky}$ duration after write)
- **Causal Token / Minimum LSN Routing** (blocking wait for replica catch-up before serving reads)
- **Lag-Aware Dynamic Routing & Primary Fallback**

## Structure
- `cmd/demo/`: Runnable demonstration CLI showcasing replication lag anomalies and mitigation strategies.
- `internal/cluster/`: Primary-replica database cluster simulation with WAL streaming and LSN tracking.
- `internal/router/`: Read/write splitting and smart routing logic with session state tracking.
- `tests/`: Automated unit, integration, and race-detector concurrency tests.
- `engineering/`: Design specifications, implementation notes, and execution results.

## Running the Lab

### Run Tests
```bash
go test -v ./...
```

### Run Race Detector
```bash
go test -race ./...
```

### Run Demo
```bash
go run ./cmd/demo
```
