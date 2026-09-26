# Code Audit

## Finding 1

Location: `internal/inventory/store.go:53-61`
Claimed Behavior: Safe concurrent reading of product state.
Observed Implementation: `Get(id)` acquires `s.mu.Lock()` and returns a copy of `Product` (`*p`), ensuring thread-safe reads of the map and struct contents.
Assessment: PASS
Severity: LOW
Notes: Correct map synchronization and return by value prevents callers from mutating internal store pointer state directly.

## Finding 2

Location: `internal/inventory/store.go:65-89`
Claimed Behavior: Replicate naive unsynchronized read-modify-write lost update anomaly.
Observed Implementation: `NaiveDeduct` calls `Get(id)` to read state (releasing mutex), sleeps 100 microseconds to simulate calculation window, then acquires `s.mu.Lock()` to write updated stock based on stale read state (`p.Stock - qty`).
Assessment: PASS
Severity: LOW
Notes: Accurately replicates lost update behavior seen in applications reading data outside atomic database locks.

## Finding 3

Location: `internal/inventory/store.go:93-116`
Claimed Behavior: Pessimistic locking row-level exclusivity (`SELECT ... FOR UPDATE`).
Observed Implementation: `GetRowLock(id)` retrieves/creates a row `sync.Mutex`. `PessimisticDeduct` acquires `rowLock.Lock()`, then inside `s.mu.Lock()` checks existence/stock and mutates `p.Stock`.
Assessment: PASS
Severity: LOW
Notes: Granular row mutex ensures exclusive access across concurrently operating goroutines per product ID without global table locking.

## Finding 4

Location: `internal/inventory/store.go:120-152`
Claimed Behavior: Optimistic locking version guard check (`UPDATE ... WHERE id = ? AND version = ?`).
Observed Implementation: Reads product state via `Get(id)`, pauses 50 microseconds, acquires `s.mu.Lock()`, checks if `curr.Version != p.Version`. If modified, increments `OptimisticFails` counter and returns `ErrOptimisticLock`. Otherwise updates `Stock` and increments `Version`.
Assessment: PASS
Severity: LOW
Notes: Version guard accurately simulates Optimistic Concurrency Control (OCC).

## Finding 5

Location: `internal/inventory/service.go:28-45`
Claimed Behavior: Optimistic concurrency control retry loop with jittered exponential backoff.
Observed Implementation: Retries up to `maxRetries` upon receiving `ErrOptimisticLock`. Calculates `(1<<attempt)*ms + rand(0..5ms)` backoff delay before re-attempting operation. Non-optimistic errors fail fast immediately.
Assessment: PASS
Severity: LOW
Notes: Idiomatic retry loop prevents thundering herd problem.

## Finding 6

Location: `internal/inventory/store.go:155-173`
Claimed Behavior: Atomic single-statement update (`UPDATE ... WHERE stock >= qty`).
Observed Implementation: Checks stock and updates state within a single mutex-protected block (`s.mu.Lock()`), incrementing `Atomically` counter.
Assessment: PASS
Severity: LOW
Notes: Effectively models database statement-level atomic execution where read and write occur in a single non-interleavable operation.
