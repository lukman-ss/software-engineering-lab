# software-engineer-lab

> **A hands-on laboratory for software engineering beyond CRUD.**

## What This Repository Is

This is a **mini production system** used to explore real-world failure modes in distributed systems. It is not a framework, not a tutorial collection for basic Go syntax, and not a collection of examples to copy.

It is an **engineering lab** where you:

1. Reproduce a failure scenario
2. Understand the root cause
3. Implement a correct solution
4. Verify with tests and benchmarks

## Why Beyond CRUD

CRUD applications assume sequential single-user access with perfect networks. Production systems must handle:

- **Concurrent users** → race conditions, lost updates
- **Unreliable networks** → timeouts, partial failures, retries
- **Partial outages** → circuit breakers, rate limiting
- **Data consistency** → transactions, distributed coordination

## Learning Philosophy

### Satu Lab, Satu Mental Model

Setiap lab berfokus pada **satu failure mode spesifik** dan **satu mental model utama**. 

> A lab may mention related failure modes, but should not fully implement topics owned by later labs.

Batasan kepemilikan konsep antar lab:
- **Lab 01 Idempotency** → berfokus pada *repeated logical operation* & *retry safety* (tidak mengajarkan locking mendalam).
- **Lab 05 Race Condition** → berfokus pada memahami bagaimana concurrent access dapat merusak business invariant.
- **Lab 11 Pessimistic Locking** → menjaga critical read-modify-write dengan database locking.
- **Lab 03 Database Transaction & Distributed Transaction Boundary** → berfokus pada *atomic multi-step database changes*, *transaction boundary*, *partial failure* dengan external systems, serta pengantar *event, retry, saga, compensation, dan outbox*.
- **Lab 04 Caching** → berfokus pada *cache stampede, stale reads, and consistency*.
- **Lab 12 Optimistic Locking** → berfokus pada *concurrent modification/version conflict*.

Tujuan repository adalah *progressive learning*, bukan menyelesaikan semua masalah dalam satu lab.

### No Hand-Waving

Every lab solves a real problem with production-grade patterns:

- Timeouts are explicit, never "instant"
- Errors are wrapped, never ignored
- Context propagates through every I/O call
- Resources clean up via `defer`

### Code Is Documentation

- Function names describe the **why**, not just the what
- Error messages guide operators to root causes
- Comments explain trade-offs and limitations

### Test What Matters

- Unit tests verify pure logic
- Integration tests verify system boundaries
- Concurrent tests expose race conditions
- Benchmarks measure solution overhead

## Architecture

```
software-engineer-lab/
├── cmd/              # Entry points (api, worker, mock-provider)
├── internal/         # Business domains (order, payment, etc.)
├── pkg/              # Shared utilities (database, observability, resilience)
├── labs/             # Isolated failure labs
├── migrations/       # Database schema evolution
├── docs/             # Architecture and engineering notes
└── Makefile
```

## Domain

The mini production system models a simple e-commerce flow:

```
Customer creates order
        ↓
Inventory is reserved
        ↓
Payment is processed
        ↓
Order becomes paid
        ↓
Notification is emitted
```

### Domain Entities

| Entity | Purpose |
|--------|---------|
| Order | Aggregate root for an order |
| OrderItem | Line item in an order |
| Payment | Payment record (idempotent) |
| InventoryItem | Stock for a product |
| Wallet | User balance |
| WalletTransaction | Audit log for wallet operations |
| Notification | Outgoing notification record |

### Layer Separation

```
domain/        # Entity definitions and interfaces
application/   # Service orchestration
repository/    # Database access (PostgreSQL)
transport/     # HTTP handlers (future)
infrastructure/ # Cross-cutting (db, observability, resilience)
```

## Labs

| Lab | Problem |
|-----|---------|
| 01 | Duplicate requests creating duplicate payments |
| 02 | Database index query optimization |
| 03 | Partial failures leaving inconsistent state |
| 04 | Cache stampede and stale reads |
| 05 | Race conditions corrupting inventory counts |
| 06 | API versioning - breaking changes and backward compatibility |
| 07 | Observability |
| 08 | Database Isolation Level: non-repeatable read, phantom read, serializable conflict, SELECT FOR UPDATE deadlock prevention |
| 09 | Code Review: identifying bugs before production |
| 10 | Project Estimation: breakdown, uncertainty, risk, spikes, ranges, effort vs duration, and contingency |
| 11 | Lock contention causing delays |
| 12 | Optimistic locking for concurrent updates |
| 13 | Deadlocks hanging the system |
| 14 | Outbox Pattern: transactional outbox deep dive, retry, DLQ, idempotency, recovery |
| 15 | Retry storms amplifying failures |
| 16 | Circuit Breaker: cascading failures from downstream |
| 17 | Rate Limiting: token bucket, client limits, resource protection |

## How to Run

### Prerequisites

- Go 1.25+
- Docker and Docker Compose
- `make`

### Setup

```bash
make infra-up        # Start PostgreSQL, Redis, Prometheus, Grafana
make migrate-up      # Run database migrations
make run             # Start API on :8080
```

### Run a Specific Lab

```bash
cd labs/01-idempotency
go test -v ./...
```

### Run Lab 05 (Race Condition)

```bash
make lab-05-test        # Unit tests
make lab-05-test-race   # Tests with race detector
make lab-05-vet         # Vet
make lab-05-fmt         # Format
make lab-05-integration # Integration tests (requires DB)
```

Note: Lab 05 is a nested Go module. Run from repository root with `make lab-05-test` or from the lab directory directly with `cd labs/05-race-condition && go test ./...`. Running `go test ./...` from the root will NOT test Lab 05.

## Testing

```bash
make test          # Unit tests
make test-race     # Tests with race detector
make lint          # go vet
```

## Production Engineering Topics

### Concurrency
- Race conditions, locking (optimistic vs pessimistic)
- Deadlock prevention
- Goroutine leaks

### Reliability
- Idempotency, retries with backoff
- Circuit breakers, timeouts, cancellations

### Distributed Systems
- Two-phase commit alternatives
- Event ordering and duplication
- Consistency models, distributed locks

### Performance
- Connection pool sizing, query optimization
- Result caching, backpressure

### Observability
- Structured logging, metrics and tracing
- Error correlation, health checks

## Roadmap

### Phase 1: Core Labs
- [x] 01 Idempotency
- [x] 02 Database Index
- [ ] 03 Database Transaction
- [x] 04 Caching
- [x] 05 Race Condition
- [x] 06 API Versioning
- [x] 07 Observability
- [x] 08 Database Isolation Level
- [x] 09 Code Review
- [x] 10 Project Estimation
- [x] 11 Pessimistic Locking
- [x] 12 Optimistic Locking
- [x] 13 Deadlock
- [x] 14 Outbox Pattern
- [x] 15 Retry
- [x] 16 Circuit Breaker
- [x] 17 Rate Limiting

### Phase 2: Extended Topics
- Bulkhead pattern
- Saga pattern
- Leader election
- Feature flags
- Graceful shutdown

## What This Is Not

- ❌ A Go tutorial (you need basic Go knowledge)
- ❌ A framework (no magic, explicit code)
- ❌ Production-ready as-is (educational focus)
- ❌ Claiming to solve everything (what's left to fail is documented)

## References

- [Go Proverbs](https://go-proverbs.github.io/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Martin Fowler's Patterns](https://martinfowler.com/categories.html)