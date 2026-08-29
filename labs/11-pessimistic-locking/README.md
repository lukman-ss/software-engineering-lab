# Lab 11 — Pessimistic Locking

Pessimistic locking prevents concurrency conflicts by acquiring an exclusive lock on the data before reading and modifying it (`SELECT ... FOR UPDATE`).

## 1. Pessimistic vs. Optimistic Locking

| Metric | Optimistic Locking | Pessimistic Locking (`FOR UPDATE`) |
|--------|--------------------|------------------------------------|
| **Assumption** | Conflicts are rare | Conflicts are frequent or costly |
| **Locking Mechanism** | Version check at commit time (`WHERE version = v`) | Row-level exclusive lock during transaction |
| **Concurrency Behavior** | High concurrency; updates fail on conflict and require retry | Serialized access on hot rows; waiting queue forms |
| **Best Used For** | Low contention, UI forms, read-heavy workloads | High contention, financial balances, inventory reservation |

## 2. Lock Duration & Contention

Holding a pessimistic lock for too long (e.g., executing slow queries, HTTP calls, or lengthy processing inside the transaction) drastically impacts system performance:
- **Waiting Requests**: Concurrent goroutines queue up waiting for the lock to release, causing latency spikes.
- **Throughput Drop**: Transactions block sequentially, lowering overall system throughput.
- **Connection Pool Saturation**: Active transactions tie up database connections.

---

## Navigasi

- **Previous**: [Lab 10 — Rate Limiting](../10-rate-limiting/)
- **Next**: [Lab 12 — Optimistic Locking](../12-optimistic-locking/)
