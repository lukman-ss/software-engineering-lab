# Code Audit

Target Lab: labs/23-optimistic-vs-pessimistic-locking

## Finding 1: In-Memory Storage Engine Emulation of Database Primitives
Location: `internal/inventory/store.go:9-51`
Claimed Behavior: Emulates row-level locking (`SELECT ... FOR UPDATE`), optimistic version validation (`WHERE version = ?`), and atomic statement execution (`UPDATE ... SET stock = stock - ? WHERE stock >= ?`) within pure Go standard library.
Observed Implementation:
- Per-row mutexes stored in `rowLocks map[int]*sync.Mutex` managed thread-safely via store mutex `s.mu`.
- State access and version checking strictly guarded by mutexes and atomic counters.
- Clean separation of storage engine internals and application business service.
Assessment: PASS
Severity: LOW
Notes: Appropriate in-memory modeling without external database dependencies; adheres strictly to YAGNI and portability requirements.

## Finding 2: Naive Read-Modify-Write Lost Update Manifestation
Location: `internal/inventory/store.go:65-89`
Claimed Behavior: Concurrent callers experience lost updates when executing separate read and write phases without locking.
Observed Implementation:
- Reads row state via `s.Get(id)`, yields control (`time.Sleep(100 * time.Microsecond)`), then acquires lock to write computed stale value back.
- Goroutines concurrently calculate based on stale stock snapshot, reliably demonstrating the lost update anomaly without triggering data races on Go memory (memory access is synchronized, logic is unsynchronized).
Assessment: PASS
Severity: LOW
Notes: Thread-safe memory access combined with logical race window is correctly implemented.

## Finding 3: Pessimistic Row Locking
Location: `internal/inventory/store.go:93-116`
Claimed Behavior: Acquires row-level exclusivity prior to reading and holds through modification, serializing access per ID.
Observed Implementation:
- Acquires `rowLock := s.GetRowLock(id)` then `rowLock.Lock()`, deferring unlock.
- Checks stock sufficiency, decrements atomically with respect to other row workers.
- Returns `ErrInsufficientStock` when stock is inadequate.
Assessment: PASS
Severity: LOW
Notes: Properly reflects `SELECT ... FOR UPDATE` isolation semantics.

## Finding 4: Optimistic Version Verification & Jittered Retry
Location: `internal/inventory/store.go:120-152` and `internal/inventory/service.go:28-45`
Claimed Behavior: Checks version matching upon commit/write. Rejects mismatched versions with `ErrOptimisticLock`. Service retries using exponential backoff with jitter.
Observed Implementation:
- Stores version field; compares `curr.Version != p.Version` under store lock; increments version on successful write.
- Retry loop in `DeductOptimisticWithRetry` performs up to `maxRetries` with `(1<<attempt)*time.Millisecond + rand.Intn(5)*time.Millisecond`.
Assessment: PASS
Severity: LOW
Notes: Fully captures optimistic concurrency control (OCC) commit validation and exponential backoff retry pattern.

## Finding 5: Atomic Conditional Decrement
Location: `internal/inventory/store.go:155-173`
Claimed Behavior: Atomic decrement conditional on sufficient stock within a single operation.
Observed Implementation:
- Acquires store mutex, validates `curr.Stock >= qty`, decrements, unlocks.
- Simulates SQL single-statement evaluation `UPDATE ... WHERE stock >= qty`.
Assessment: PASS
Severity: LOW
Notes: Correctly models single-statement SQL behavior.
