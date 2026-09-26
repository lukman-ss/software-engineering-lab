# Code Audit

## Finding 1

Location: store.go:65-89 NaiveDeduct
Claimed Behavior: Unsynchronized read-modify-write reproduces lost update.
Observed Implementation: Read via Get, 100µs sleep, write stale `p.Stock - qty` under `s.mu`. Stale overwrite confirmed.
Assessment: PASS
Severity: LOW
Notes: No mutex on read path is intentional. Counter `NaivelyDrawn` incremented even on overwritten writes — correct for demonstrating anomaly.

## Finding 2

Location: store.go:93-116 PessimisticDeduct
Claimed Behavior: SELECT ... FOR UPDATE row exclusivity.
Observed Implementation: Per-row `sync.Mutex` held across read-check-write. Correct serialization. `GetRowLock` mutex-protected map access.
Assessment: PASS
Severity: LOW
Notes: Lock ordering single-resource, no deadlock path.

## Finding 3

Location: store.go:120-152 OptimisticDeduct
Claimed Behavior: UPDATE ... WHERE version=? guard.
Observed Implementation: Snapshot read, 50µs delay, compare `curr.Version != p.Version` under `s.mu`, increment version on success. Stock guard pre-check only;-but write-phase does not recheck `Stock < qty`. Stale snapshot with sufficient stock could drive stock negative if concurrent drains interleave between snapshot and commit.
Assessment: WARNING
Severity: MEDIUM
Notes: No oversell test exercises concurrent drain-to-zero for optimistic path. Atomic path does guard correctly.

## Finding 4

Location: store.go:155-173 AtomicDeduct
Claimed Behavior: UPDATE ... WHERE stock >= qty single-statement.
Observed Implementation: Check-and-decrement under single `s.mu` critical section. Correct.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 5

Location: service.go:28-45 DeductOptimisticWithRetry
Claimed Behavior: Retry with jittered exponential backoff.
Observed Implementation: Loop `attempt <= maxRetries`, return non-conflict errors immediately, sleep `(1<<attempt)ms + rand(0-5)ms`. Uses unseeded global `math/rand` — deterministic jitter sequence, functionally fine.
Assessment: PASS
Severity: LOW
Notes: `rand.Intn` on shared global source under concurrency is safe in Go 1.22+.

## Finding 6

Location: store.go:28-51 Seed / GetRowLock / Get
Claimed Behavior: Thread-safe simulated engine.
Observed Implementation: All map accesses under `s.mu`. `NaiveDeduct` missing-ID path: `curr := s.products[id]` nil deref panic if ID absent after successful Get — unreachable in practice (no delete API), no test covers it.
Assessment: WARNING
Severity: LOW
Notes: Minor robustness gap, no delete path exists.

## Finding 7

Location: store.go counters + service.go error propagation
Claimed Behavior: Error propagation, cleanup.
Observed Implementation: All `ErrNotFound / ErrInsufficientStock / ErrInvalidQuantity / ErrOptimisticLock` propagated correctly. Mutexes use defer. No leaks.
Assessment: PASS
Severity: LOW
Notes: `DeductOptimisticWithRetry` correctly stops retry on non-conflict errors.
