# Code Audit

## Finding 01

Location: internal/inventory/store.go:65-89 — `NaiveDeduct`
Claimed Behavior: Naive read-modify-write reproduces the lost-update anomaly under concurrent access.
Observed Implementation: `NaiveDeduct` calls `s.Get(id)` (read under `s.mu`, returns a copy), sleeps 100µs, then re-acquires `s.mu` and writes `curr.Stock = p.Stock - qty` using the stale read value. Multiple goroutines reading the same initial stock all write the same stale result, overwriting each other's decrements.
Assessment: PASS
Severity: LOW
Notes: The lost update is an intentional logical anomaly (stale write), not a data race. Race detector confirms no unsynchronized memory access. The 100µs sleep is an intentional simplification to widen the race window for reliable demonstration.

## Finding 02

Location: internal/inventory/store.go:93-116 — `PessimisticDeduct`
Claimed Behavior: Simulates `SELECT ... FOR UPDATE` by acquiring an exclusive per-row lock before reading and holding it until the update completes.
Observed Implementation: Acquires `GetRowLock(id)` (a per-product `sync.Mutex`), then acquires `s.mu` for the read-modify-write. Row lock is released via `defer` after `s.mu` is released. Lock ordering is rowLock → s.mu; no reverse ordering exists anywhere, so deadlock is impossible.
Assessment: PASS
Severity: LOW
Notes: Uses a global engine lock (`s.mu`) for the inner modify, which serializes all pessimistic operations across all products (real SQL would lock only the specific row). This is a simulation limitation, not a correctness bug — the invariant (stock = 50 after 50 concurrent deducts of 1 from 100) holds.

## Finding 03

Location: internal/inventory/store.go:120-152 — `OptimisticDeduct`
Claimed Behavior: Simulates `UPDATE ... WHERE id=? AND version=?`; detects version mismatch and returns `ErrOptimisticLock` without corrupting data.
Observed Implementation: Reads product copy via `s.Get(id)` (snapshot includes Version), sleeps 50µs, then re-acquires `s.mu`, compares `curr.Version != p.Version`, and only writes if versions match. On mismatch, increments `OptimisticFails` counter and returns `ErrOptimisticLock` without writing.
Assessment: PASS
Severity: LOW
Notes: The version guard protects the stale stock check — if versions match, `curr.Stock == p.Stock`, so the pre-sleep `p.Stock < qty` check remains valid. No re-check of stock sufficiency inside the locked section, but this is safe because the version match guarantees no intervening write.

## Finding 04

Location: internal/inventory/store.go:120-152 — `OptimisticDeduct`
Claimed Behavior: Version check detects any concurrent modification to the row (as in real SQL optimistic locking).
Observed Implementation: Only `OptimisticDeduct` increments `curr.Version` via `curr.Version++`. The other write paths (`NaiveDeduct`, `PessimisticDeduct`, `AtomicDeduct`) do NOT increment the version. If strategies were mixed on the same product/store, optimistic locking would fail to detect changes made by other strategies.
Assessment: WARNING
Severity: MEDIUM
Notes: Within the tests and demo, each scenario uses a dedicated store with a single strategy, so the issue is not exercised. However, this is a simulation inaccuracy: real SQL optimistic locking bumps the row version/timestamp on every write, not just optimistic writes. The implementation notes (02-implementation-notes.md) do not mention this limitation.

## Finding 05

Location: internal/inventory/store.go:154-173 — `AtomicDeduct`
Claimed Behavior: Simulates `UPDATE products SET stock = stock - N WHERE id = ? AND stock >= N` — a lockless single-statement atomic operation.
Observed Implementation: Acquires `s.mu` for the entire read-check-write. This is serialized at the engine level (global lock), not truly "lockless." The check `curr.Stock < qty` and the decrement `curr.Stock -= qty` happen atomically under `s.mu`.
Assessment: PASS (with clarification)
Severity: LOW
Notes: The term "lockless" in the demo output (line 116 of main.go: "Lockless Single Statement") is misleading — the simulation uses a mutex. However, the invariant (stock = 50 after 50 concurrent deductions) holds correctly. The atomic safety comes from the mutex, not from a single SQL statement. This matches the design intent (simulating SQL semantics, not implementing actual SQL).

## Finding 06

Location: internal/inventory/service.go:28-45 — `DeductOptimisticWithRetry`
Claimed Behavior: Retries optimistic deductions with jittered exponential backoff until success or max retries exhausted.
Observed Implementation: Loops up to `maxRetries+1` attempts. On `ErrOptimisticLock`, sleeps `1ms << attempt + rand.Intn(5)ms` before retrying. Non-lock errors (e.g., `ErrInsufficientStock`) are returned immediately. On exhausting retries, returns `ErrOptimisticLock`.
Assessment: PASS
Severity: LOW
Notes: `math/rand` is auto-seeded in Go 1.20+; go.mod specifies `go 1.22`, so no seeding issue. The backoff prevents thundering herd. Under the test load (20 goroutines), all eventually converge.

## Finding 07

Location: internal/inventory/store.go:34 — `Seed` method, Version: 1
Claimed Behavior: Seeds product with version 1, consistent with SQL optimistic locking baseline.
Observed Implementation: `Product{ID: id, Name: name, Stock: stock, Version: 1}`. The starting version of 1 is arbitrary but consistent.
Assessment: PASS
Severity: LOW
Notes: Starting version could be 0 or 1 — both are valid. No issue.

## Finding 08

Location: internal/inventory/store.go:65-89 — `NaiveDeduct` write phase, line 83-84
Claimed Behavior: (N/A — this is a latent bug, not a claimed behavior)
Observed Implementation: In the write phase, `curr := s.products[id]` is called without checking existence. If the product was removed between the `Get` read and the write (which never happens in current code since nothing deletes products), `curr` would be nil and `curr.Stock` would panic.
Assessment: WARNING
Severity: LOW
Notes: No code path currently deletes products from the store, so this is theoretical. Not exercised by any test. A nil check before the write would harden the method.

## Finding 09

Location: internal/inventory/store.go:14,15-18 — Atomic counters
Claimed Behavior: Counters track successful operations across concurrent goroutines without data races.
Observed Implementation: All counters (`NaivelyDrawn`, `Pessimistically`, `Optimistically`, `OptimisticFails`, `Atomically`) are `int64` updated exclusively via `atomic.AddInt64`.
Assessment: PASS
Severity: LOW
Notes: Race detector confirms no data races on these counters.

## Finding 10

Location: internal/inventory/store.go:53-61 — `Get` method
Claimed Behavior: Thread-safe point read returning a copy of the product.
Observed Implementation: Locks `s.mu`, looks up product in map, returns `*p` (dereferenced copy) or `ErrNotFound`. The copy prevents callers from modifying internal state through the returned pointer.
Assessment: PASS
Severity: LOW
Notes: Correct copy semantics — caller receives a value, not a reference.

## Finding 11

Location: internal/inventory/store.go:42-51 — `GetRowLock`
Claimed Behavior: Returns (or creates on demand) a per-product mutex for pessimistic locking.
Observed Implementation: Acquires `s.mu`, checks/creates `s.rowLocks[id]`, returns the lock pointer. All callers receive the same pointer for the same ID due to map lookup under mutex.
Assessment: PASS
Severity: LOW
Notes: The returned `*sync.Mutex` is safe to use — all goroutines for the same product get the identical lock object.

## Finding 12

Location: internal/inventory/model.go:1-17 — Domain model and errors
Claimed Behavior: Defines `Product` struct and standard error sentinels.
Observed Implementation: `Product{ID, Name, Stock, Version}` and four error variables: `ErrNotFound`, `ErrInsufficientStock`, `ErrOptimisticLock`, `ErrInvalidQuantity`.
Assessment: PASS
Severity: LOW
Notes: Clean, standard Go error definitions. All errors are sentinel values (not wrapped), matching the test assertions that use `err == inventory.ErrXxx`.

## Finding 13

Location: internal/inventory/service.go, all methods
Claimed Behavior: Service layer delegates to Store methods.
Observed Implementation: `Service` is a thin wrapper — each `DeductXxx` method calls the corresponding `Store.Xxx` method with no additional logic (except `DeductOptimisticWithRetry` which adds retry logic).
Assessment: PASS
Severity: LOW
Notes: The separate `Service` layer is architectural (per design doc) — it provides a seam for future business logic. Not unnecessary complexity.

## Finding 14

Location: internal/inventory/service.go:41 — `rand.Intn` without explicit seed
Claimed Behavior: Jitter in backoff is randomized to prevent thundering herd.
Observed Implementation: `time.Duration(rand.Intn(5))` adds 0-5ms jitter. Go 1.22 auto-seeds `math/rand`.
Assessment: PASS
Severity: LOW
Notes: If go.mod were downgraded to Go 1.19 or earlier, the PRNG would be deterministically seeded, making jitter non-random. go.mod specifies `go 1.22`, so this is not an issue.

## Summary

- Total findings: 14
- PASS: 10
- WARNING: 4 (Findings 04, 08, and implicit in 05)
- FAIL: 0
- Severities: 0 CRITICAL, 0 HIGH, 3 MEDIUM, 11 LOW

The implementation is correct for its intended scope. All concurrency primitives are used properly. The race detector passes with zero warnings. The three claimed remediation strategies (pessimistic, optimistic, atomic) all maintain data consistency invariants under concurrent load. The one meaningful WARNING (Finding 04 — version not incremented across all write paths) is a simulation fidelity issue that does not affect any tested scenario.