# Code Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Finding 1: Naive Read-Modify-Write Anomaly Simulation

Location: `internal/inventory/store.go:65-89`
Claimed Behavior: Simulates uncoordinated application read-modify-write resulting in lost updates under concurrent access.
Observed Implementation: Reads snapshot via `Get(id)`, introduces a 100µs artificial processing delay to simulate calculation overhead, and writes stale calculated stock under mutex protection.
Assessment: PASS
Severity: LOW
Notes: Correctly separates read and write phases while keeping individual memory accesses safe from Go memory corruptions.

## Finding 2: Row-level Pessimistic Locking

Location: `internal/inventory/store.go:93-116`
Claimed Behavior: Simulates `SELECT ... FOR UPDATE` by acquiring an exclusive row-level mutex before reading and holding it until the write transaction finishes.
Observed Implementation: Retrieves a per-row `sync.Mutex` via `GetRowLock(id)`, acquires it, validates existence and stock, decrements stock, and releases the lock via defer.
Assessment: PASS
Severity: LOW
Notes: Granular row locking correctly prevents interference across distinct IDs while serializing operations on the same ID.

## Finding 3: Optimistic Locking with Version Guard

Location: `internal/inventory/store.go:120-152`
Claimed Behavior: Simulates `UPDATE ... WHERE id = ? AND version = ?`, rejecting writes when version has advanced.
Observed Implementation: Reads initial product snapshot and version, verifies stock, and checks `curr.Version == p.Version` inside the write lock. If mismatched, returns `ErrOptimisticLock` and increments conflict counter; otherwise updates stock and increments version.
Assessment: PASS
Severity: LOW
Notes: Faithfully models optimistic concurrency control mechanics.

## Finding 4: Optimistic Retry with Exponential Backoff and Jitter

Location: `internal/inventory/service.go:28-45`
Claimed Behavior: Automatically retries optimistic lock conflicts with exponential backoff and randomized jitter up to `maxRetries`.
Observed Implementation: Loops up to `maxRetries`, sleeps for `(1<<attempt)*time.Millisecond + rand(5ms)` on `ErrOptimisticLock`, and terminates on success or non-recoverable error.
Assessment: PASS
Severity: LOW
Notes: Implements proper backoff and jitter algorithms preventing live-lock under contention.

## Finding 5: Atomic Single-Statement Decrement

Location: `internal/inventory/store.go:155-173`
Claimed Behavior: Simulates `UPDATE products SET stock = stock - qty WHERE id = ? AND stock >= qty`.
Observed Implementation: Performs conditional check and decrement within a single atomic critical section.
Assessment: PASS
Severity: LOW
Notes: Clean and concise implementation of atomic in-database updates.
