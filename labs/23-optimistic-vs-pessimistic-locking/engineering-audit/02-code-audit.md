# Code Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Finding 1: In-Memory Row Locking Simulation

Location: `internal/inventory/store.go:42-51, 91-116`
Claimed Behavior: Simulates `SELECT ... FOR UPDATE` row exclusivity with per-row locks.
Observed Implementation: `Store` uses `map[int]*sync.Mutex` with synchronized retrieval and acquisition in `PessimisticDeduct`. Dedicated mutex per ID ensures no cross-product contention.
Assessment: PASS
Severity: LOW
Notes: Concurrency safe and properly separates row lock lifecycle from storage map mutex.

## Finding 2: Optimistic Concurrency Control with Version Guard

Location: `internal/inventory/store.go:120-152` and `internal/inventory/service.go:28-45`
Claimed Behavior: Simulates optimistic locking (`UPDATE ... WHERE id = ? AND version = ?`), detects conflicts when version deviates, and supports backoff retry.
Observed Implementation: Version check `curr.Version != p.Version` inside critical section increments `OptimisticFails` counter and returns `ErrOptimisticLock`. `DeductOptimisticWithRetry` applies jittered exponential backoff `(1<<attempt)*time.Millisecond + rand(5)ms`.
Assessment: PASS
Severity: LOW
Notes: Correctly models SQL version check semantics and retry handling.

## Finding 3: Atomic Conditional Update

Location: `internal/inventory/store.go:154-173`
Claimed Behavior: Simulates single-statement atomic update (`UPDATE ... SET stock = stock - N WHERE stock >= N`).
Observed Implementation: Checks stock condition and deducts in a single mutex block, simulating SQL database engine-level row-lock and single statement execution.
Assessment: PASS
Severity: LOW
Notes: Correctly captures lock-free application level pattern that relies on DB atomic statement guarantees.

## Finding 4: Input Validation and Error Propagation

Location: `internal/inventory/store.go:66, 94, 121, 156`
Claimed Behavior: Validates input quantities and returns structured domain errors.
Observed Implementation: Returns `ErrInvalidQuantity` for `qty <= 0`, `ErrNotFound` for non-existent products, and `ErrInsufficientStock` when balance is inadequate.
Assessment: PASS
Severity: LOW
Notes: Standard error handling follows idiomatic Go patterns.
