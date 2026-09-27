# Code Audit

## Finding 1
Location: internal/inventory/store.go:65-89 (`NaiveDeduct`)
Claimed Behavior: Reproduces lost-update anomaly via unsynchronized read-modify-write.
Observed Implementation: Reads `Product` (copy) via `Get`, sleeps 100us, then under `s.mu` writes `p.Stock - qty` (stale copy) back to the live `curr.Stock` pointer.
Assessment: PASS
Severity: LOW
Notes: Logic correctly simulates lost update. Race safety around map/counter access: all shared-state access occurs under `s.mu` and `atomic.*`. No data race.

## Finding 2
Location: internal/inventory/store.go:93-116 (`PessimisticDeduct`)
Claimed Behavior: `SELECT ... FOR UPDATE` row exclusivity prevents conflicts.
Observed Implementation: Acquires per-row `sync.Mutex` via `GetRowLock`, then reads-modifies under `s.mu`.
Assessment: PASS
Severity: LOW
Notes: Correct exclusive lock pattern. No deadlock risk: row lock acquired first, `s.mu` acquired/released inside, row lock released last via `defer`.

## Finding 3
Location: internal/inventory/store.go:120-152 (`OptimisticDeduct`)
Claimed Behavior: Version guard rejects stale writes; non-corrupting.
Observed Implementation: Reads copy + version via `Get`, sleeps 50us, then under `s.mu` compares `curr.Version != p.Version`. Mismatch returns `ErrOptimisticLock` and increments `OptimisticFails`. Match applies decrement and increments version.
Assessment: PASS
Severity: LOW
Notes: Faithful to `UPDATE ... WHERE id=? AND version=?`. Version increment +1 on success is correct.

## Finding 4
Location: internal/inventory/service.go:28-45 (`DeductOptimisticWithRetry`)
Claimed Behavior: Retry loop with jittered exponential backoff converges.
Observed Implementation: Loops up to `maxRetries`, retries only on `ErrOptimisticLock`, sleeps `(1<<attempt)ms + rand(0,5)ms`.
Assessment: PASS
Severity: LOW
Notes: Standard backoff+jitter pattern. Returns `ErrOptimisticLock` after exhausting retries.

## Finding 5
Location: internal/inventory/store.go:154-173 (`AtomicDeduct`)
Claimed Behavior: Single-statement `UPDATE ... WHERE stock >= qty`.
Observed Implementation: Atomically under `s.mu`: check existence, check `stock < qty`, decrement.
Assessment: PASS
Severity: LOW
Notes: Correct conditional decrement semantics.

## Finding 6
Location: internal/inventory/store.go:14
Claimed Behavior: Counters (`NaivelyDrawn`, etc.) are concurrency-safe.
Observed Implementation: All counter mutations use `atomic.AddInt64`.
Assessment: PASS
Severity: LOW
Notes: No non-atomic counter access found.

## Finding 7
Location: internal/inventory/store.go:42-51 (`GetRowLock`)
Claimed Behavior: Per-row lock exists for each product ID.
Observed Implementation: Lazily creates row mutex in map under `s.mu`. Returns pointer; caller acquires it.
Assessment: WARNING
Severity: LOW
Notes: Lazy creation is safe under `s.mu`, but `Seed` also creates the lock. Redundant creation paths both guarded. Acceptable.

## Finding 8
Location: internal/inventory/model.go, service.go, store.go
Claimed Behavior: Error handling for boundary inputs.
Observed Implementation: All deduct methods check `qty <= 0` and return `ErrInvalidQuantity`. `ErrNotFound` and `ErrInsufficientStock` handled.
Assessment: PASS
Severity: LOW
Notes: Boundary checks present but not exercised by tests (see test audit).

## Finding 9
Location: internal/inventory/store.go:53-61 (`Get`)
Claimed Behavior: Returns a copy of the Product.
Observed Implementation: `return *p, nil` copies struct.
Assessment: PASS
Severity: LOW
Notes: Copy prevents external mutation of internal state.
