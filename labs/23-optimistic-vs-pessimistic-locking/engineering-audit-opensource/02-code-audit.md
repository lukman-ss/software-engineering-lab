# Code Audit

## Overview

Module: github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking
Go version: 1.22
Implementation files:
- internal/inventory/model.go (17 lines)
- internal/inventory/store.go (173 lines)
- internal/inventory/service.go (49 lines)
- cmd/demo/main.go (121 lines)

Architecture: single in-memory Store with per-row mutexes simulating SQL row locks; Service layer exposing deduction strategies.

## Finding 1: Naive read-modify-write lost update simulation

Location: internal/inventory/store.go:63-89
Claimed Behavior: NaiveDeduct reads stock, sleeps (calculation window), then writes using the stale calculated value, reproducing lost update.
Observed Implementation: `Get` returns a Product value snapshot `p`. A 100us sleep elapses. Under `s.mu`, the code re-reads `curr` but writes `curr.Stock = p.Stock - qty` (stale value, not `curr.Stock - qty`). The `NaivelyDrawn` atomic counter is incremented per call.
Assessment: PASS
Severity: -
Notes: Correct simulation of the lost-update anomaly. Stale write guarantees overwritten updates. Demo confirmed stock=99, not 50.

## Finding 2: Pessimistic locking (row lock) correctness

Location: internal/inventory/store.go:92-116
Claimed Behavior: PessimisticDeduct acquires per-row mutex before read and holds until update completes, serializing concurrent writers.
Observed Implementation: `GetRowLock(id)` fetches/creates a per-row `sync.Mutex` from `rowLocks` map (creation guarded by `s.mu`). The row lock is held for the duration of read-check-decrement under `s.mu`. Stock check and decrement are atomic relative to other pessimistic callers.
Assessment: PASS
Severity: -
Notes: 50 concurrent goroutines each decrementing 1 from 100 yields exactly 50. Demo and test both confirm.

## Finding 3: Optimistic locking version guard

Location: internal/inventory/store.go:118-152
Claimed Behavior: OptimisticDeduct reads current version, sleeps, re-checks version under lock; if version changed returns ErrOptimisticLock.
Observed Implementation: `Get` returns snapshot with Version. After 50us sleep, `s.mu` is acquired; `curr.Version != p.Version` triggers `OptimisticFails++` and returns `ErrOptimisticLock`. On match, `curr.Stock -= qty; curr.Version++`.
Assessment: PASS
Severity: -
Notes: Version check and increment are atomic under `s.mu`. Demo confirms 1 success, 19 conflicts rejected, stock=99.

## Finding 4: Atomic single-statement conditional update

Location: internal/inventory/store.go:154-173
Claimed Behavior: AtomicDeduct checks `stock >= qty` and decrements under a single mutex-hold, modeling `UPDATE ... SET stock = stock - N WHERE stock >= N`.
Observed Implementation: Acquires `s.mu`, checks existence and `curr.Stock < qty`, then `curr.Stock -= qty`. Returns ErrInsufficientStock if check fails.
Assessment: PASS
Severity: -
Notes: 50 concurrent decrements of 1 from 100 yields exactly 50. Demo confirms.

## Finding 5: Optimistic retry with jittered backoff

Location: internal/inventory/service.go:28-45
Claimed Behavior: DeductOptimisticWithRetry retries on ErrOptimisticLock up to maxRetries with jittered exponential backoff.
Observed Implementation: Loop `attempt 0..maxRetries` (inclusive). Returns nil on success; returns non-ErrOptimisticLock errors immediately; returns ErrOptimisticLock after final attempt. Backoff: `1<<attempt ms + rand(5) ms`.
Assessment: PASS
Severity: -
Notes: Bounded retry prevents infinite loop. Demo confirms convergence (20 successes, stock=80).

## Finding 6: Concurrency safety of shared state

Location: internal/inventory/store.go (all methods)
Claimed Behavior: All shared state (products map, rowLocks map, counters) is protected.
Observed Implementation: `s.mu` guards `products` and `rowLocks` maps. Per-row `rowLocks[id]` mutexes guard pessimistic row operations. `NaivelyDrawn`, `Pessimistically`, `Optimistically`, `OptimisticFails`, `Atomically` are `int64` updated via `atomic.AddInt64`. No field is written without protection.
Assessment: PASS
Severity: -
Notes: Race detector run passed with zero warnings (see test audit).

## Finding 7: Error handling and input validation

Location: internal/inventory/model.go, internal/inventory/store.go (all methods)
Claimed Behavior: Invalid quantity returns ErrInvalidQuantity; missing product returns ErrNotFound; insufficient stock returns ErrInsufficientStock; version mismatch returns ErrOptimisticLock.
Observed Implementation: every Deduct function checks `qty <= 0` first. All return correct error sentinels. No panics or unhandled errors.
Assessment: PASS
Severity: -
Notes: Error types defined in model.go; consistently used across store and service.

## Finding 8: Demo output authenticity

Location: cmd/demo/main.go
Claimed Behavior: Demo executes real concurrent goroutines and prints real measured results (final stock, conflict counts, elapsed time).
Observed Implementation: Demo runs actual goroutines, uses real `store.Get` calls, reads real counter values (`store.Optimistically`, `store.OptimisticFails`), and prints `time.Since(start)`. No hardcoded output strings. Verified by running `go run ./cmd/demo` which produced variable results consistent with code logic.
Assessment: PASS
Severity: -
Notes: Output matches execution (stock=99 naive, stock=50 pessimistic, 1 success/19 conflicts optimistic, 20 optimistic-with-retry, stock=50 atomic). No fake numbers.
