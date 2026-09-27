# Optimistic vs Pessimistic Locking & Atomic Updates

## Problem

When two concurrent transactions read the same data, compute a new value independently, and write back sequentially, the second write silently overwrites the first. This is the **lost update** anomaly — a corruption that produces no error.

Example: Transaction 1 reads stock=100, calculates 99. Transaction 2 reads stock=100 (before T1 commits), calculates 99. Both write 99. Final stock = 99, but the correct result is 98.

This anomaly is not a bug in application logic alone — it is a consequence of the **read-modify-write** pattern executing without concurrency control. Under PostgreSQL's default **READ COMMITTED** isolation level, this scenario is explicitly possible: each statement re-evaluates against the latest committed row version, so a write from another concurrent transaction can be silently overwritten.

## Why This Matters

Lost updates violate business invariants — for inventory, `stock = initial - total_deductions` no longer holds. Impact:

- Financial reporting inconsistencies
- Inventory discrepancies (phantom stock, phantom shortages)
- Double-booking (seats, tickets, reservations)
- Real financial loss from over-deducted or under-deducted resources

A transaction by itself does **not** prevent this. Isolation level READ COMMITTED, which is the default on PostgreSQL and Oracle, still permits lost updates on application-level read-modify-write sequences. Concurrency control must be explicit.

## Mental Model

```text
Time →
T1:  [READ stock=100] → [CALC 100-1=99] → [WRITE 99] → [COMMIT]
T2:  [READ stock=100] ————————→ [CALC 100-1=99] → [WRITE 99] → [COMMIT]

Result: 99 ❌   Correct: 99 → wait, actually 98 ❌
```

The core problem: T2 reads before T1 writes, then T2 writes a value computed from the old snapshot. T1's write is overwritten and lost.

## Core Concept

Three strategies address this, each eliminating the race window differently:

### 1. Pessimistic Locking

`SELECT ... FOR UPDATE` acquires a **row-level exclusive lock** held until transaction commit or rollback. Other transactions attempting to lock the same row block until the lock is released.

```sql
BEGIN;
SELECT stock FROM products WHERE id=1 FOR UPDATE;
-- row now locked; no other transaction can modify it
UPDATE products SET stock = stock - 1 WHERE id = 1;
COMMIT;
-- lock released
```

- **Mechanism**: Block concurrent writers on the same row.
- **Use case**: High-conflict, business-critical data (wallet balances, limited stock, seat reservations).
- **Cost**: Reduced concurrency; deadlock risk; must keep transactions short.

### 2. Optimistic Locking

A version column (integer or timestamp) is stored on the row. The update validates that the version has not changed since the read:

```sql
BEGIN;
SELECT stock, version FROM products WHERE id = 1;  -- version = 5
UPDATE products
  SET stock = 99, version = 6
  WHERE id = 1 AND version = 5;   -- version guard
-- if 0 rows affected → another transaction modified the row → conflict
COMMIT;
```

- **Mechanism**: Proceed without blocking; detect conflict at commit by checking `affected_rows == 0`.
- **Use case**: Read-heavy, low-conflict scenarios (user profiles, CMS content).
- **Requirement**: Application-level retry or conflict response (409).

### 3. Atomic Single-Statement Update

Eliminate the read-modify-write window entirely by performing the check and update in a single statement:

```sql
UPDATE products SET stock = stock - 1 WHERE id = 1 AND stock >= 1;
-- check affected_rows == 1 for confirmation
```

- **Mechanism**: Database evaluates `WHERE stock >= 1` and applies `stock = stock - 1` atomically within one statement.
- **Use case**: Simple in-place arithmetic (counters, quantity decrement).
- **Scope limit**: Works for single-row, single-condition mutations. Multi-row invariants require transactions or locking.

## Architecture

The lab uses an in-memory store that simulates SQL storage-engine semantics:

```text
┌──────────────────────────────────────────────────────┐
│  cmd/demo/main.go   (CLI demo runner)                │
│  tests/locking_test.go  (concurrency tests)           │
│                                                      │
│  ┌─────────────────────────────────────────┐         │
│  │ internal/inventory/                     │         │
│  │ ├── model.go   (Product struct, errors) │         │
│  │ ├── store.go   (simulated DB engine)    │         │
│  │ └── service.go (business operations)    │         │
│  └─────────────────────────────────────────┘         │
└──────────────────────────────────────────────────────┘
```

The `Store` models row-level locking, versioned rows, and atomic updates using Go's `sync.Mutex`. This is an in-memory simulation — it reproduces the logic of each strategy rather than connecting to a real RDBMS. Artificial micro-delays (100µs in the naive path, 50µs in the optimistic path) widen the race window so the anomaly is reproducible.

## Implementation

### Store — Simulated Engine (`internal/inventory/store.go`)

- `rowLocks map[int]*sync.Mutex` — one mutex per product row, modeling `SELECT ... FOR UPDATE`.
- `products map[int]*Product` — table rows, each with `Stock` and `Version` fields.
- `NaiveDeduct` — read, sleep, then write using the stale snapshot. Not safe.
- `PessimisticDeduct` — acquire the row lock first, validate, decrement. Safe.
- `OptimisticDeduct` — read with version snapshot, validate version under lock, return `ErrOptimisticLock` on mismatch.
- `AtomicDeduct` — single lock-validate-update. Safe and lockless from the caller's perspective.

### Service Layer (`internal/inventory/service.go`)

- `DeductNaive` — direct wrapper of the unsafe path.
- `DeductPessimistic` — direct wrapper of the safe blocked path.
- `DeductOptimisticDirect` — optimistic without retry.
- `DeductOptimisticWithRetry(id, qty, maxRetries)` — retry loop with jittered exponential backoff: `(1<<attempt)ms + rand(0..4)ms`.
- `DeductAtomic` — direct wrapper of the atomic path.

### Model & Errors (`internal/inventory/model.go`)

```go
type Product struct {
    ID      int
    Name    string
    Stock   int
    Version int
}

var (
    ErrNotFound          = errors.New("product not found")
    ErrInsufficientStock = errors.New("insufficient stock")
    ErrOptimisticLock    = errors.New("optimistic lock conflict: record modified by another transaction")
    ErrInvalidQuantity   = errors.New("quantity must be positive")
)
```

The `Version` field supports optimistic locking. `ErrOptimisticLock` represents a commit-time conflict (equivalent to SQL `affected_rows == 0`). `ErrInsufficientStock` is a business-rule rejection, not corruption.

## Code Walkthrough

### Naive Read-Modify-Write (the bug)

```go
// internal/inventory/store.go, NaiveDeduct
func (s *Store) NaiveDeduct(id int, qty int) error {
    p, err := s.Get(id)                         // 1. READ snapshot
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }

    time.Sleep(100 * time.Microsecond)           // widened race window

    s.mu.Lock()
    curr := s.products[id]
    curr.Stock = p.Stock - qty                   // 2. WRITE from stale snapshot
    s.mu.Unlock()
    return nil
}
```

The read and write are separate steps. Between them, another goroutine may commit a different value. The write then overwrites that newer value with a stale calculation. No Go data race occurs (memory access is mutex-protected); the anomaly is at the logic level.

### Pessimistic Locking (block then proceed)

```go
// internal/inventory/store.go, PessimisticDeduct
func (s *Store) PessimisticDeduct(id int, qty int) error {
    rowLock := s.GetRowLock(id)
    rowLock.Lock()                              // block other writers on this row
    defer rowLock.Unlock()

    s.mu.Lock()
    p, exists := s.products[id]
    if !exists { s.mu.Unlock(); return ErrNotFound }
    if p.Stock < qty { s.mu.Unlock(); return ErrInsufficientStock }
    p.Stock -= qty
    s.mu.Unlock()
    return nil
}
```

The row lock serializes all writers on this row. The global `s.mu` only guards the map structure; the row lock governs business logic.

### Optimistic Locking (detect then reject)

```go
// internal/inventory/store.go, OptimisticDeduct
func (s *Store) OptimisticDeduct(id int, qty int) error {
    p, err := s.Get(id)                         // snapshot with version
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }

    time.Sleep(50 * time.Microsecond)           // widened conflict window

    s.mu.Lock()
    defer s.mu.Unlock()
    curr := s.products[id]
    if curr.Version != p.Version {              // version guard (WHERE version = ?)
        s.OptimisticFails++
        return ErrOptimisticLock                // 0-rows-affected equivalent
    }
    curr.Stock -= qty
    curr.Version++
    s.Optimistically++
    return nil
}
```

If another transaction committed between the read and the validation, `curr.Version` differs from `p.Version`. The conflict is returned as `ErrOptimisticLock`. No data corruption occurs — the write is safe to reject.

### Optimistic with Retry (converge)

```go
// internal/inventory/service.go, DeductOptimisticWithRetry
func (svc *Service) DeductOptimisticWithRetry(id int, qty int, maxRetries int) error {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := svc.store.OptimisticDeduct(id, qty)
        if err == nil { return nil }
        if err != ErrOptimisticLock { return err }         // non-retryable
        if attempt == maxRetries { return ErrOptimisticLock }
        sleepDuration := time.Duration(1<<attempt)*time.Millisecond +
            time.Duration(rand.Intn(5))*time.Millisecond   // jittered backoff
        time.Sleep(sleepDuration)
    }
    return ErrOptimisticLock
}
```

Exponential backoff (`1 << attempt` ms) plus jitter (`rand.Intn(5)` ms) prevents the thundering-herd problem. After `maxRetries` failed attempts, the conflict is returned and the caller must decide: fail the request (409 Conflict) or escalate.

### Atomic Single-Statement (no read-modify-write gap)

```go
// internal/inventory/store.go, AtomicDeduct
func (s *Store) AtomicDeduct(id int, qty int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if !exists { return ErrNotFound }
    if curr.Stock < qty { return ErrInsufficientStock }
    curr.Stock -= qty                                      // check + update atomically
    s.Atomically++
    return nil
}
```

No intermediate read phase is exposed to the caller. The check (`stock >= qty`) and the mutation (`stock -= qty`) happen within a single locked critical section, modeling `UPDATE ... WHERE stock >= N`.

## What the Tests Prove

Two invariants are tested across all safe strategies:

1. **Correctness invariant**: After N concurrent decrements of 1 from initial stock S, final stock = S − N (when all succeed).
2. **Corruption-free failure**: When a strategy cannot safely proceed (insufficient stock, version mismatch), it returns an error — it never writes a corrupted value.

| Test | Scenario | Assertion | Result |
|------|----------|-----------|--------|
| `TestNaiveLostUpdate` | 50 goroutines, DeductNaive(1, 1), stock 100 | `Stock != 50` (anomaly detected) | PASS (stock = 99) |
| `TestPessimisticLocking` | 50 goroutines, DeductPessimistic(2, 1), stock 100 | `Stock == 50` (invariant held) | PASS |
| `TestPessimisticLockingInsufficientStock` | stock=2, deduct 2 then 1 | second returns `ErrInsufficientStock` | PASS |
| `TestOptimisticLockingConflict` | 20 goroutines, DeductOptimisticDirect(3, 1), stock 100 | conflicts detected (`> 0`); `successCount + Stock == 100` | PASS (1 success, 19 conflicts) |
| `TestOptimisticLockingWithRetry` | 20 goroutines, retry max=10, stock 100 | all converge; `Stock == 100 - Optimistically` | PASS |
| `TestAtomicConditionalUpdate` | 50 goroutines, DeductAtomic(5, 1), stock 100 | `Stock == 50` | PASS |

All tests pass, and `go test -race ./...` reports zero data races. The race detector confirms memory-safety; the test assertions confirm logic-safety.

## Recovery / Rollback

- **Pessimistic / deadlock**: Databases like PostgreSQL and Oracle detect deadlocks automatically and abort one transaction. The application must catch the error and retry. Mitigation: consistent lock ordering, short transactions, never hold locks across network I/O.
- **Optimistic / conflict**: A detected conflict returns `ErrOptimisticLock` (equivalent to SQL `affected_rows == 0`). The application retries with backoff, or surfaces a 409 Conflict to the caller. No partial write occurs — the transaction never commits a stale value.
- **Atomic / insufficient stock**: Returns `ErrInsufficientStock`. This is a business-rule rejection, not corruption. The row is unchanged.

## Production Considerations

- **Pessimistic locking** is appropriate for high-conflict, business-critical resources. Never hold the lock across external calls (HTTP, payment gateways) — doing so extends the transaction, increases deadlock probability, and destroys throughput.
- **Optimistic locking** suits read-heavy, low-conflict workloads. Retry logic is mandatory — without it, conflicts become failed requests. For very high write frequency, prefer a 64-bit integer or timestamp version column to reduce overflow risk.
- **Atomic updates** are the simplest correct solution for counters and quantity decrements, provided no complex validation is required before the update. They do not replace transactions when invariants span multiple rows or tables.
- **Isolation levels alone are not sufficient**: READ COMMITTED (the default) permits lost updates on application-level read-modify-write. REPEATABLE READ and SERIALIZABLE prevent some anomalies but introduce serialization errors (e.g., PostgreSQL's `could not serialize access due to concurrent update`) that also require retry handling.

## Common Mistakes

1. **Assuming a transaction prevents lost updates** — Under READ COMMITTED, it does not. Application-level read-modify-write still races.
2. **Holding locks across network I/O** — Extends the critical section, increasing contention and deadlock likelihood.
3. **Ignoring `affected_rows == 0`** — An optimistic update that affects zero rows is a conflict. Treating it as success silently skips the deduction.
4. **Naive read-modify-write without a guard** — `load() → modify() → save()` is the classic lost-update pattern.
5. **Reaching for Redis before database-native locking** — If the data is in a relational database, `SELECT ... FOR UPDATE` or `advisory locks` are simpler and correct. External locks add complexity and are only justified for cross-database or cross-process resources.

## Case Study

The demo (`go run ./cmd/demo`) runs five side-by-side scenarios, each starting from stock = 100:

```text
[1] Naive Read-Modify-Write (50 concurrent requests):
    Initial Stock: 100
    Expected Final Stock: 50
    Actual Final Stock:   99 (LOST UPDATE DETECTED!)

[2] Pessimistic Locking (SELECT ... FOR UPDATE):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Fully Synchronized)

[3] Optimistic Locking Direct (20 concurrent requests, no retry):
    Initial Stock: 100
    Successful Deductions: 1
    Rejected Conflicts:   19
    Actual Final Stock:   99 (SUCCESS - State Guarded, Zero Corruption)

[4] Optimistic Locking With Exponential Backoff Retry (20 requests):
    Initial Stock: 100
    Successful Deductions: 20
    Total Attempted Conflicts Retried: 61
    Actual Final Stock:   80 (SUCCESS - All retries eventually converged)
    Elapsed Time:         ~50ms

[5] Atomic Single-Statement Operation (UPDATE ... WHERE stock >= qty):
    Initial Stock: 100
    Actual Final Stock:   50 (SUCCESS - Lockless Single Statement)
```

Scenario [1] demonstrates the bug: 50 decrements yield a final stock of 99 instead of 50 (49 concurrent overwrites clobbered the result). Scenarios [2] and [5] maintain the invariant exactly (100 − 50 = 50). Scenario [3] shows optimistic detection without retry — only 1 write succeeds, 19 conflicts are rejected, and the invariant `successCount + Stock == 100` holds. Scenario [4] shows the retry loop converging: all 20 requests eventually succeed, with 61 total conflict retries absorbed.

> Note: Concurrency scheduling in Go is nondeterministic. Exact figures (final stock in naive/optimistic, retry counts in retry) may vary between runs. The invariants and anomaly behavior are consistent; the specific numbers are illustrative.

## Checklist

- [ ] Identify the operation: read-modify-write or simple in-place mutation?
- [ ] Assess conflict frequency: high, moderate, or low?
- [ ] Choose strategy:
  - High conflict or critical data → Pessimistic (`SELECT FOR UPDATE`)
  - Low conflict, read-heavy → Optimistic (version guard + retry)
  - Simple counter/quantity → Atomic single-statement `UPDATE`
- [ ] Never hold a pessimistic lock across external I/O
- [ ] Implement retry with exponential backoff for optimistic conflicts
- [ ] Always check `affected_rows` (or error result) to detect conflicts
- [ ] Run the race detector: `go test -race ./...`
- [ ] Verify the invariant: `FinalStock = InitialStock − TotalSuccessfulDeductions`

## Key Takeaways

1. Lost update is silent data corruption: two transactions read the same row, compute independently, and the second write overwrites the first — no error is raised.
2. A transaction alone does not prevent it. READ COMMITTED (default on PostgreSQL and Oracle) permits lost updates on application-level read-modify-write sequences.
3. Pessimistic locking (`SELECT FOR UPDATE`) prevents conflicts by blocking concurrent writers on the same row until the lock holder commits.
4. Optimistic locking detects conflicts at commit time via a version guard in the WHERE clause; `affected_rows == 0` signals a conflict that must be retried or surfaced as an error.
5. Atomic single-statement updates (`UPDATE ... SET stock = stock - N WHERE stock >= N`) eliminate the read-modify-write window entirely for simple in-place mutations.
6. Choose the strategy by conflict frequency and data criticality — there is no universal solution.
7. Deadlocks are a risk of pessimistic locking; databases auto-detect and abort one transaction. Mitigate with consistent lock ordering and short transactions.
8. Optimistic strategies require retry logic. Without it, conflicts become failed requests.
9. Isolation level alone is not sufficient — correct query design (lock, version guard, or atomic statement) is required for consistency.
10. The test suite proves both the anomaly (naive loses updates) and the invariant (safe strategies maintain `InitialStock − TotalSuccessfulDeductions`).
11. This lab is an in-memory simulation (`sync.Mutex`), not a real RDBMS. MySQL documentation was not directly accessible (HTTP 403); MySQL-related behavior is based on the SQL standard and secondary sources.
12. Error-path coverage is partial: `ErrInvalidQuantity`, `ErrNotFound`, optimistic retry exhaustion, and concurrent overdraft are implemented but not individually tested by the automated suite.

## Sources

- PostgreSQL 18 Documentation, "13.2 Transaction Isolation" and "13.3 Explicit Locking" — defines READ COMMITTED behavior, `SELECT FOR UPDATE`, snapshot isolation, and deadlock auto-detection.
- Oracle 19c Concepts, "Data Concurrency and Consistency (Chapter 10)" — row locks (TX), isolation levels, lost update under READ COMMITTED, ORA-08177.
- Wikipedia, "Optimistic Concurrency Control" — OCC by Kung & Robinson (1981), three-phase model.
- Wikipedia, "Pessimistic Offline Lock" — application-level pessimistic locking pattern.
- Wikipedia, "Write-write conflict (Lost Update)" — formal definition of the anomaly.
- Martin Fowler, "Optimistic Offline Lock" — pre-commit validation pattern.
- Martin Fowler, "Pessimistic Offline Lock" — conflict prevention by lock acquisition.
- Microsoft Learn, "Handling Concurrency Conflicts — EF Core" — version-guard pattern, `DbUpdateConcurrencyException`, `affected_rows == 0` semantics.

> MySQL 8.0 documentation was not directly accessible during research (HTTP 403 from dev.mysql.com). Claims about MySQL/InnoDB are based on SQL-standard behavior and secondary sources; the in-memory lab does not require a MySQL instance.
