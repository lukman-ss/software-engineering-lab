# Code Audit

Target Lab: `labs/23-optimistic-vs-pessimistic-locking`

## Finding 1

Location: `internal/inventory/store.go:65-89` (`NaiveDeduct`)
Claimed Behavior: Unsynchronized read-modify-write pattern that triggers lost updates under concurrent access.
Observed Implementation: Reads stock via `Get(id)`, sleeps 100µs to simulate processing window, then acquires global `s.mu` to overwrite `curr.Stock = p.Stock - qty` using the stale read value.
Assessment: PASS
Severity: LOW
Notes: Correctly demonstrates the lost update anomaly in a controlled manner for educational purposes.

## Finding 2

Location: `internal/inventory/store.go:93-116` (`PessimisticDeduct`)
Claimed Behavior: Granular row-level locking simulating `SELECT ... FOR UPDATE`.
Observed Implementation: Fetches/creates a per-row `sync.Mutex` via `GetRowLock(id)`, locks the row mutex before reading, holds it through check and stock decrement. Releases `s.mu` during return.
Assessment: PASS
Severity: LOW
Notes: Properly isolates row locks. Map access to `rowLocks` and `products` is protected by `s.mu`.

## Finding 3

Location: `internal/inventory/store.go:120-152` (`OptimisticDeduct`)
Claimed Behavior: Version-guarded update simulating `UPDATE ... WHERE id = ? AND version = ?`.
Observed Implementation: Reads version via `Get(id)` without row lock, sleeps 50µs, acquires `s.mu`, verifies `curr.Version == p.Version`. If mismatched, increments `OptimisticFails` and returns `ErrOptimisticLock`. If matched, updates stock and increments version.
Assessment: PASS
Severity: LOW
Notes: Correctly implements optimistic version verification and state protection.

## Finding 4

Location: `internal/inventory/service.go:28-45` (`DeductOptimisticWithRetry`)
Claimed Behavior: Optimistic locking with retry and jittered exponential backoff.
Observed Implementation: Loops up to `maxRetries`. On `ErrOptimisticLock`, calculates sleep duration `(1<<attempt)*ms + rand(0..5ms)` and retries. Returns non-optimistic errors immediately.
Assessment: PASS
Severity: LOW
Notes: Clean implementation of backoff retry convergence.

## Finding 5

Location: `internal/inventory/store.go:155-173` (`AtomicDeduct`)
Claimed Behavior: Atomic single-statement update simulating `UPDATE ... SET stock = stock - N WHERE stock >= N`.
Observed Implementation: Performs check and decrement within `s.mu` lock in a single step without holding lock across read phase or external sleeps.
Assessment: PASS
Severity: LOW
Notes: Accurately models database statement-level atomic execution semantics.

## Finding 6

Location: `internal/inventory/store.go:66-68`, `94-96`, `121-123`, `156-158`
Claimed Behavior: Input validation for quantities.
Observed Implementation: Validates `qty <= 0` returning `ErrInvalidQuantity` across all deduction methods.
Assessment: PASS
Severity: LOW
Notes: Prevents negative or zero deductions across all strategies.
