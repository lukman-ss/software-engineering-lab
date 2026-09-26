# Code Audit

## Overview

Source files reviewed: 4 (`model.go`, `service.go`, `store.go`, `cmd/demo/main.go`)
Test files reviewed: 1 (`tests/locking_test.go`)
Module: `github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking`
Go version: `go 1.22` (go.mod), go vet passes, `go test -race` passes with zero warnings.

---

## Finding 1

Location: `internal/inventory/store.go:65-89` (`NaiveDeduct`)
Claimed Behavior: Reproduces the classic lost-update (read-modify-write race) anomaly by reading stock, releasing the lock, simulating a computation delay, then writing back a stale value.
Observed Implementation: Calls `s.Get(id)` (which locks `s.mu`, reads a copy, releases lock), sleeps 100µs, then re-acquires `s.mu` and writes `curr.Stock = p.Stock - qty` using the stale copy `p`. The `NaivelyDrawn` counter is updated via `atomic.AddInt64`.
Assessment: PASS
Severity: LOW
Notes: The stale-write is intentional and correctly reproduces the anomaly. A nil-pointer guard on `curr` is absent but not exploitable since no delete method exists and `Get()` already validated existence. The `NaivelyDrawn` counter tracks attempted deductions; the lost-update manifests because the stock only decreases by 1 regardless of how many concurrent goroutines write 99.

---

## Finding 2

Location: `internal/inventory/store.go:93-116` (`PessimisticDeduct`)
Claimed Behavior: Simulates `SELECT ... FOR UPDATE` by acquiring an exclusive per-row mutex before reading and holding it until the update completes, preventing lost updates.
Observed Implementation: Calls `s.GetRowLock(id)` to obtain a per-row `*sync.Mutex`, acquires it, defers its release, then acquires `s.mu`, validates existence and stock, mutates `p.Stock -= qty`, releases `s.mu`, and increments `Pessimistically` via `atomic.AddInt64`.
Assessment: PASS
Severity: LOW
Notes: Lock ordering is `rowLock` → `s.mu`. No reverse ordering (`s.mu` → `rowLock`) exists anywhere in the codebase, so deadlock is impossible. The per-row granularity allows concurrent access to different product IDs. `Get` returns a copy and `Seed` always stores a non-nil `*Product`, so the `p.Stock -= qty` mutation is safe.

---

## Finding 3

Location: `internal/inventory/store.go:120-152` (`OptimisticDeduct`)
Claimed Behavior: Simulates `UPDATE ... WHERE id=? AND version=?` by reading the product (including version), sleeping, then re-reading and comparing versions under the engine lock. A mismatch returns `ErrOptimisticLock`; a match applies the update and increments the version.
Observed Implementation: Calls `s.Get(id)` (returns a copy with current version), sleeps 50µs, acquires `s.mu`, reads `curr` from the live map, compares `curr.Version != p.Version`, increments `OptimisticFails` and returns `ErrOptimisticLock` on mismatch, otherwise decrements stock and increments version under the lock.
Assessment: PASS
Severity: LOW
Notes: The version check and update are atomic under `s.mu`, so there is no TOCTOU window within the write phase. No data corruption is possible because failed version checks reject the write before mutation. The pattern correctly models the SQL `WHERE version = ?` guard.

---

## Finding 4

Location: `internal/inventory/store.go:155-173` (`AtomicDeduct`)
Claimed Behavior: Simulates an atomic single-statement conditional update (`UPDATE ... WHERE stock >= qty`) by performing the stock check and decrement under a single lock acquisition.
Observed Implementation: Acquires `s.mu` upfront, validates existence and `curr.Stock >= qty`, decrements `curr.Stock -= qty`, and increments `Atomically` via `atomic.AddInt64`.
Assessment: PASS
Severity: LOW
Notes: The entire read-check-write is serialized by `s.mu`, eliminating the lost-update window. This is the simplest correct strategy for in-process simulation.

---

## Finding 5

Location: `internal/inventory/service.go:28-45` (`DeductOptimisticWithRetry`)
Claimed Behavior: Retries optimistic deductions with jittered exponential backoff on `ErrOptimisticLock`; returns other errors immediately; returns `ErrOptimisticLock` if `maxRetries` is exhausted.
Observed Implementation: Loops `attempt` from 0 to `maxRetries` (inclusive), calls `OptimisticDeduct`, returns nil on success, returns non-OptimisticLock errors immediately, returns `ErrOptimisticLock` on the final attempt, and sleeps `(1 << attempt) * 1ms + rand.Intn(5) * 1ms` between retries.
Assessment: PASS
Severity: LOW
Notes: The loop bound `attempt <= maxRetries` means `maxRetries + 1` total attempts (correct: initial + retries). `math/rand.Intn` is safe for concurrent use in Go 1.20+ (go.mod specifies `go 1.22`). Only `ErrOptimisticLock` triggers a retry; `ErrInsufficientStock`, `ErrNotFound`, and `ErrInvalidQuantity` short-circuit the loop, which is correct.

---

## Finding 6

Location: `internal/inventory/store.go:42-51` (`GetRowLock`)
Claimed Behavior: Returns (or lazily creates) a per-row mutex for a given product ID, enabling row-level locking.
Observed Implementation: Acquires `s.mu`, checks `s.rowLocks[id]`, creates a new `*sync.Mutex` if absent, stores it, and returns it. The mutex itself is then acquired by the caller (`PessimisticDeduct`).
Assessment: PASS with WARNING
Severity: LOW
Notes: The lazy creation is safe under `s.mu`. However, `GetRowLock` can create a lock for a non-existent product ID, which is a minor resource leak in the `rowLocks` map. Not exploitable for data corruption, but it allows unbounded growth of `rowLocks` if called with arbitrary IDs. The method is only called from `PessimisticDeduct`, which validates the product under `s.mu` afterward, so the practical impact is negligible.

---

## Finding 7

Location: `internal/inventory/store.go:28-40` (`Seed`)
Claimed Behavior: Inserts or overwrites a product with initial stock and version 1, creating a per-row lock if one does not already exist.
Observed Implementation: Acquires `s.mu`, creates `s.products[id] = &Product{...}`, and conditionally creates `s.rowLocks[id]` if absent.
Assessment: PASS
Severity: LOW
Notes: Re-seeding the same ID overwrites the product (including resetting version to 1) while preserving the existing row lock. This is acceptable for a test/demo utility.

---

## Finding 8

Location: `internal/inventory/store.go:53-61` (`Get`)
Claimed Behavior: Returns a snapshot copy of a product under the engine lock.
Observed Implementation: Acquires `s.mu`, retrieves `s.products[id]`, returns `*p` (dereferenced copy) or `ErrNotFound`.
Assessment: PASS
Severity: LOW
Notes: Returning a copy ensures callers cannot mutate internal state through the returned reference. This is critical for `OptimisticDeduct` and `NaiveDeduct`, which compare stale copies against live state.

---

## Finding 9

Location: `internal/inventory/store.go:14-19` (`Store` struct counters)
Claimed Behavior: All counters (`NaivelyDrawn`, `Pessimistically`, `Optimistically`, `OptimisticFails`, `Atomically`) track successful or failed operations across goroutines.
Observed Implementation: Each counter is `int64` and mutated exclusively via `atomic.AddInt64`. No direct reads/writes outside atomic operations.
Assessment: PASS
Severity: LOW
Notes: All counter access is atomic. The `go test -race` run confirms no data races on these fields.

---

## Finding 10

Location: `internal/inventory/model.go:5-10` (error definitions)
Claimed Behavior: Domain errors are defined as sentinel values for direct comparison.
Observed Implementation: `errors.New` creates four sentinel errors: `ErrNotFound`, `ErrInsufficientStock`, `ErrOptimisticLock`, `ErrInvalidQuantity`.
Assessment: PASS
Severity: LOW
Notes: Sentinel errors are compared via `==` in both the service and tests. This is idiomatic Go for domain-level error codes. Using `errors.Is` would also work but is unnecessary since these are not wrapped.
