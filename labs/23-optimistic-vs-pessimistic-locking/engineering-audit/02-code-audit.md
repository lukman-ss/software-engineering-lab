# Code Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Finding 1

Location: `internal/inventory/store.go:65-89` (NaiveDeduct)
Claimed Behavior: Simulate uncoordinated read-modify-write resulting in lost update anomaly.
Observed Implementation: Reads stock under `Get(id)`, artificially sleeps to open concurrency window, then writes back `p.Stock - qty` based on stale snapshot.
Assessment: PASS
Severity: LOW
Notes: Properly models unsynchronized application-level read-modify-write cycles.

## Finding 2

Location: `internal/inventory/store.go:93-116` (PessimisticDeduct)
Claimed Behavior: Simulate `SELECT ... FOR UPDATE` using exclusive row-level mutex.
Observed Implementation: Fetches row-specific mutex via `GetRowLock(id)`, locks before reading state, validates stock constraint, and writes update before unlocking.
Assessment: PASS
Severity: LOW
Notes: Thread-safe, protects row boundary correctly, and prevents concurrent mutation during read-modify-write.

## Finding 3

Location: `internal/inventory/store.go:120-152` (OptimisticDeduct)
Claimed Behavior: Simulate optimistic concurrency control via version guard (`WHERE version = @v`).
Observed Implementation: Reads item version without row lock, pauses, then under store mutex verifies whether version matches snapshot. Increments version upon success, returns `ErrOptimisticLock` upon conflict.
Assessment: PASS
Severity: LOW
Notes: Invariant `curr.Version != p.Version` strictly rejects concurrent interleaved modifications.

## Finding 4

Location: `internal/inventory/service.go:28-45` (DeductOptimisticWithRetry)
Claimed Behavior: Retry optimistic lock conflicts with jittered exponential backoff.
Observed Implementation: Loops up to `maxRetries`, catches `ErrOptimisticLock`, calculates exponential backoff delay with random jitter `(1<<attempt)*time.Millisecond + rand.Intn(5)*time.Millisecond`, and returns final error if attempts exhausted.
Assessment: PASS
Severity: LOW
Notes: Correctly handles non-conflict errors immediately without spurious retries.

## Finding 5

Location: `internal/inventory/store.go:155-173` (AtomicDeduct)
Claimed Behavior: Simulate atomic conditional single-statement update (`UPDATE ... WHERE stock >= qty`).
Observed Implementation: Evaluates stock constraint and decrements in a single synchronized critical section mimicking database engine in-place page latch/write.
Assessment: PASS
Severity: LOW
Notes: Clean, minimal, zero extraneous lock overhead.

## Finding 6

Location: `cmd/demo/main.go:1-121`
Claimed Behavior: Demonstrates all 5 scenarios with live execution and metrics.
Observed Implementation: Runs 5 distinct test scenarios concurrently using goroutines and waitgroups, logging actual stock counters and conflict statistics.
Assessment: PASS
Severity: LOW
Notes: Real runnable execution, zero mocked outputs.
