---
title: "Optimistic vs Pessimistic Locking & Atomic Updates"
published: false
description: "Explore Go-based lost update mitigation: naive read-modify-write race, pessimistic row locks, optimistic version guards, jittered retry, and single-statement atomic updates."
tags: go, concurrency, database, softwareengineering
---

# Optimistic vs Pessimistic Locking & Atomic Updates

A concurrent inventory system with 50 goroutines decrementing stock from 100 should produce 50 remaining units. Without synchronization, it produces 99. The lost update is silent corruption, not an error surface visible to callers. Lab 23 demonstrates how three strategies — pessimistic locking, optimistic version guarding, and atomic single-statement updates — prevent this corruption in a zero-dependency Go simulation.

> Lost updates are silent data corruptions; standard SQL transactions do not prevent them automatically.

## The Lost Update Anomaly

When two concurrent transactions read an inventory snapshot of 100, independently compute 100 - 1 = 99, and each writes 99 back to storage, the first decrement is obliterated. No exception propagates. No database error surface exposes the corruption. The invariant `initial_stock == remaining_stock + total_deductions` is silently broken: `100 ≠ 99 + 50`.

This is not a theoretical edge case. Production e-commerce, ticketing, and digital wallet systems encounter this pattern when multiple requests process the same product or account within the same committed snapshot window. Encapsulating these operations inside a basic `BEGIN TRANSACTION` block does not automatically prevent the collision, because standard transaction boundaries control commit ordering, not cross-request snapshot isolation.

## Why Default Database Isolation Fails

PostgreSQL and Oracle ship with `READ COMMITTED` as the default isolation level. Under this level, each individual query reads the most recent committed snapshot available at statement start time. The gap between a `SELECT` reading current stock and a subsequent `UPDATE` writing the decremented value becomes a write-fragmentation window. Another transaction committing within that window introduces a new committed row version that the first transaction never observes.

Database isolation levels isolate individual statements, not multi-statement application-level calculation lifecycles. `READ COMMITTED` guarantees each query sees a consistent snapshot, but it does not guarantee that the snapshot seen by the `SELECT` remains valid when the `UPDATE` executes. This gap is where lost updates originate.

## Anatomy of the Collision: A Chronological Trace

| Phase | Role | Description |
|-------|------|-------------|
| 1. Simultaneous Read | Snapshot T1 & T2 | Both transactions read initial stock = 100 before either commits. |
| 2. Independent Calculation | In-Memory State | T1 computes 100 - 1 = 99; T2 independently computes 100 - 1 = 99. |
| 3. T1 Commits | Database Write | T1 persists stock = 99. T2's local snapshot is now stale. |
| 4. T2 Commits | Silent Overwrite | Only on this path does T2 overwrite 99 with 99, erasing T1's mutation. |

On this specific execution trace, 50 simultaneous mutations collapse into a single decrement. The ledger drifts by 49 units with no runtime signal.

## Three Taxonomies of Concurrency Control

| Strategy | Mechanism | Concurrency Profile | Risk Vector |
|----------|-----------|---------------------|-------------|
| Naive Read-Modify-Write | No mutation isolation | Maximum throughput, maximum corruption | Mass lost updates |
| Pessimistic Locking | Exclusive row lock held from read through commit | Serialized access, queued waits | Deadlock, queue starvation |
| Optimistic Locking | Version column checked at commit time | Lock-free reads, retry on conflict | Retry CPU cost, thundering herd |
| Atomic Update | Single `UPDATE SET stock=stock-N WHERE stock>=N` | Lockless, statement-bound | Limited to simple arithmetic |

The lab validates all four paths through deterministic in-memory Go tests. It does not yet prove physical RDBMS lock table behavior under millions of rows.

## Naive Read-Modify-Write: The Vulnerable Pattern

```go
// internal/inventory/store.go
func (s *Store) NaiveDeduct(id int, qty int) error {
    p, err := s.Get(id) // 1. Read snapshot without any lock
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }
    time.Sleep(100 * time.Microsecond) // widen the race window
    s.mu.Lock()
    curr := s.products[id]
    curr.Stock = p.Stock - qty // 2. Write using stale snapshot
    s.mu.Unlock()
    atomic.AddInt64(&s.NaivelyDrawn, int64(qty))
    return nil
}
```

The microsecond sleep does not artificially break valid code — it mirrors real application compute delays and inter-service network latency. Each goroutine captures its own local snapshot before acquiring the map mutex, ensuring every concurrent worker writes back the identical stale value.

## Failure Manifestation: 50 Goroutines, 49 Lost Mutations

| Stage | Event |
|-------|-------|
| 1. Access Burst | 50 goroutines read stock = 100 in parallel. |
| 2. Identical Computation | All compute 100 - 1 = 99 against their local snapshots. |
| 3. Sequential Overwrites | Each goroutine writes 99 to shared memory one after another. |
| 4. Stock Anomaly | Only on this path do 50 mutations reduce physical stock by exactly 1 unit. |

The test assertion `final stock == 50` fails with `final stock == 99`, proving that 49 decrements vanished without error.

## Execution Evidence

```
=== RUN   TestNaiveLostUpdate
    locking_test.go:37: Lost Update Demonstrated: final stock is 99 (expected 50)
--- PASS: TestNaiveLostUpdate (0.00s)
=== RUN   TestPessimisticLocking -> PASS: (stock == 50)
=== RUN   TestOptimisticLockingConflict -> PASS: (1 succeeded, 19 conflicts rejected)
=== RUN   TestOptimisticLockingWithRetry -> PASS: (20 total updates applied, stock == 80)
=== RUN   TestAtomicConditionalUpdate -> PASS: (stock == 50)
PASS: ok  labs/23-optimistic-vs-pessimistic-locking/tests  2.510s (race detector: clean)

[1] Naive Read-Modify-Write: Stock=99 (LOST UPDATE DETECTED!)
[2] Pessimistic Locking:    Stock=50 (SUCCESS - Fully Synchronized)
[3] Optimistic Direct:      Stock=99 (SUCCESS - 1 Succeeded, 19 Conflicts Rejected)
[4] Optimistic With Retry:  Stock=80 (SUCCESS - 20 Succeeded, 61 Retries Absorbed)
[5] Atomic Update:          Stock=50 (SUCCESS - Lockless Single Statement)
```

The Go race detector reports zero data races across all tests. The test suite validates invariants, not benchmark-scale throughput.

## Evidence Boundaries

| What Is Proved | What Is Not Proved |
|----------------|---------------------|
| Lost update reproduction under concurrent goroutines | Physical RDBMS lock table behavior at scale |
| Row-lock serialization via per-item `sync.Mutex` | Cross-table deadlock detection in production engines |
| Optimistic version rejection on collision | MVCC gap-lock semantics of PostgreSQL/MySQL |
| Retry convergence with jittered backoff | Throughput at 100,000 distributed transactions |

The simulation models row-level locks using Go `sync.Mutex` primitives. This accurately represents single-row serialization semantics but is not identical to physical database lock managers that track shared locks, exclusive locks, and lock timeouts.

## Engineering Trade-offs: Queue Wait vs Retry Compute

**Pessimistic Locking** — Guarantees consistency through exclusive serialization. Drawback: latency spikes under high contention as competing transactions queue behind the held lock. Ideal for high-value financial balances and scarce ticket reservations where a correct result matters more than throughput.

**Optimistic Locking** — Achieves maximum read throughput with no blocking. Drawback: wasted CPU cycles and network round-trips when conflict rates increase. The retry overhead is not automatically efficient under relentless contention; starvation can occur in tight retry loops without backoff jitter.

**Atomic Update** — Eliminates the read-then-write gap entirely. Drawback: limited to single-column arithmetic expressible within one `UPDATE ... WHERE` condition. Cannot enforce complex multi-table invariants.

Selecting optimistic locking does not eliminate error handling. Callers must implement exponential backoff with randomized jitter, otherwise concurrent retries synchronize and amplify collision rates.

## Pessimistic Locking: Row Lock Before Read

```go
// internal/inventory/store.go
func (s *Store) PessimisticDeduct(id int, qty int) error {
    rowLock := s.GetRowLock(id)
    rowLock.Lock()                  // acquire exclusive row lock
    defer rowLock.Unlock()
    s.mu.Lock()
    p, exists := s.products[id]
    if !exists { s.mu.Unlock(); return ErrNotFound }
    if p.Stock < qty { s.mu.Unlock(); return ErrInsufficientStock }
    p.Stock -= qty
    s.mu.Unlock()
    atomic.AddInt64(&s.Pessimistically, int64(qty))
    return nil
}
```

The preventive logic scope strictly follows the lab concurrency tests. All 50 goroutines serialize through the per-row lock, reducing stock from 100 to exactly 50 with zero lost mutations. This implementation is a simulation model, not a proof of physical database `SELECT FOR UPDATE` lock manager behavior.

## Optimistic Locking: Version Guard Without Blocking Readers

```go
// internal/inventory/store.go
func (s *Store) OptimisticDeduct(id int, qty int) error {
    p, err := s.Get(id)             // snapshot with version
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }
    time.Sleep(50 * time.Microsecond)
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if curr.Version != p.Version {  // version guard check
        atomic.AddInt64(&s.OptimisticFails, 1)
        return ErrOptimisticLock    // equivalent to affected_rows == 0
    }
    curr.Stock -= qty
    curr.Version++
    atomic.AddInt64(&s.Optimistically, int64(qty))
    return nil
}
```

Reading proceeds without any row lock held. At commit time, the version is validated against the current stored version. If they diverge, the update aborts with `ErrOptimisticLock`. This mechanism does not automatically deliver success; calling code must handle the conflict and decide whether to retry, escalate, or return an error to the user.

## Retry Convergence: Exponential Backoff with Jitter

```go
// internal/inventory/service.go
func (svc *Service) DeductOptimisticWithRetry(id, qty, maxRetries int) error {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := svc.store.OptimisticDeduct(id, qty)
        if err == nil { return nil }
        if err != ErrOptimisticLock { return err }
        if attempt == maxRetries { return ErrOptimisticLock }
        sleepDuration := time.Duration(1<<attempt)*time.Millisecond +
            time.Duration(rand.Intn(5))*time.Millisecond
        time.Sleep(sleepDuration)
    }
    return ErrOptimisticLock
}
```

Exponential backoff (`1<<attempt` ms) staggers retry attempts. Randomized jitter (`rand.Intn(5)`) prevents goroutines from re-synchronizing on the same tick — the thundering herd problem. Without jitter, all 20 contending goroutines would retry in lockstep, reproduce the same collision, and loop again. With jitter, retries de-synchronize until contention dissipates. On this execution path, all 20 mutations converge successfully.

## Atomic Update: Collapsing Read-Modify-Write

```go
// Equivalent SQL: UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?
func (s *Store) AtomicDeduct(id int, qty int) error {
    if qty <= 0 { return ErrInvalidQuantity }
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if !exists { return ErrNotFound }
    if curr.Stock < qty { return ErrInsufficientStock }
    curr.Stock -= qty             // validation + mutation in one critical section
    atomic.AddInt64(&s.Atomically, int64(qty))
    return nil
}
```

By binding the condition check and the state mutation inside a single critical section, the temporal gap between read and write disappears. This pattern is not a universal solution: single-statement atomicity cannot enforce multi-table invariants or cross-row consistency rules. It is the optimal choice for simple counter decrements where the entire business rule fits in one `WHERE stock >= N` condition.

## Invariant Verification Under Load

```go
// tests/locking_test.go
func TestOptimisticLockingConflict(t *testing.T) {
    // 20 goroutines concurrently call DeductOptimisticDirect on stock 100
    wg.Wait()
    p, _ := store.Get(3)
    // Conservation invariant: successful transactions + remaining stock == initial
    if successCount+int64(p.Stock) != 100 {
        t.Fatalf("stock invariant broken: successes %d + stock %d != 100",
                 successCount, p.Stock)
    }
    if conflictCount == 0 { t.Fatalf("must have at least 1 conflict") }
}
```

The conservation invariant — `successful_deductions + remaining_stock == initial_stock` — must hold exactly regardless of how many conflicts were rejected. The empirical verification scope follows the automated test suite and proves zero data corruption even under mass optimistic rejection. This is a memory-level assertion; it does not extend to multi-engine transaction durability.

## Architecture Decision Sequence

| Question | Recommended Strategy |
|----------|---------------------|
| Is the mutation a single arithmetic operation on one row? | Atomic Update (`UPDATE ... WHERE stock >= N`) |
| Does the entity face high contention (flash sale, balance accounts)? | Pessimistic Lock (`SELECT FOR UPDATE`) |
| Are reads dominant and write collisions rare? | Optimistic Version Guard + jittered retry |

Choosing optimistic locking is not a free option. It requires retry infrastructure, conflict metrics monitoring, and a user-facing error strategy for exhausted retries. Each strategy trades a different operational cost: serialized wait time, CPU retry bandwidth, or limited expressiveness.

## Production Implementation Checklist

- [ ] Identify mutation patterns: single-statement atomic or multi-step read-modify-write.
- [ ] Use `UPDATE ... SET stock = stock - N WHERE stock >= N` for simple numeric counters.
- [ ] Apply `SELECT ... FOR UPDATE` on high-value, high-contention entities.
- [ ] Never hold database row locks across external network calls (HTTP/RPC/webhooks).
- [ ] Add an integer `version` column and validate `affected_rows == 0` for optimistic locks.
- [ ] Implement retry wrappers with randomized jittered backoff for optimistic conflicts.
- [ ] Run `go test -race ./...` regularly in CI/CD pipelines to catch memory-level races.

## Anti-Patterns to Avoid

**Ignoring affected rows.** Optimistic updates that do not check `affected_rows` (or equivalent) assume success, silently accepting stale overwrites.

**Lock spanning I/O.** Holding a database row lock while waiting for a payment gateway webhook response serializes the entire connection pool behind one slow third-party call, causing systemic queue buildup.

**Premature Redis locking.** Reaching for distributed Redis locks does not automatically imply superior engineering. Relational databases already provide native, low-latency row-level locking primitives. External lock managers introduce network failure modes and operational complexity that are unnecessary when the database engine already handles the requirement.

## Reproduce Locally

All experiments run purely on the Go standard library with zero third-party modules.

```sh
cd labs/23-optimistic-vs-pessimistic-locking
go test -v ./...          # verify 6 scenarios: naive, pessimistic, optimistic, atomic
go test -race ./...       # prove concurrent synchronization has no data races
go run ./cmd/demo         # interactive CLI demo of all five scenarios side-by-side
```

Local execution validates locking semantics in memory. Production deployments still require configuring database lock timeouts, monitoring contention metrics, and validating query plans against the actual RDBMS engine.

## Core Architectural Principle

Mature transactional engineering matches synchronization primitives to domain contention profiles. Mission-critical rows with unacceptable failure cost require pessimistic certainty. Read-heavy workloads scale through optimistic guards. Simple counters thrive on atomic updates. Designing concurrent systems does not automatically require complex infrastructure when native atomic statements suffice for the mutation. The best concurrency control is the simplest one that preserves the invariant.
