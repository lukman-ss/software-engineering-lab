## Finding 1

Location: internal/inventory/store.go:65-88 (NaiveDeduct method)
Claimed Behavior: Naive read-modify-write should demonstrate lost update anomaly under concurrency.
Observed Implementation: 
- Reads product under s.mu lock, releases lock, sleeps 100μs, then re-acquires s.mu to write updated stock.
- The sleep creates a window where other goroutines can read the same stock value, leading to lost updates.
- Uses atomic.AddInt64 to increment NaivelyDrawn counter (for observation).
Assessment: PASS
Severity: 
Notes: Correctly implements the naive pattern that loses updates. The separate read and write phases under different lock instances create the necessary race condition.

## Finding 2

Location: internal/inventory/store.go:91-115 (PessimisticDeduct method)
Claimed Behavior: Simulates SELECT ... FOR UPDATE by acquiring exclusive row lock before read and holding until update completes.
Observed Implementation:
- Acquires per-row mutex via GetRowLock(id) before any database access.
- Locks s.mu to read product, check stock, update stock, then unlocks s.mu.
- Defers rowLock unlock after s.mu unlock.
- Increments Pessimistically counter atomically.
Assessment: PASS
Severity: 
Notes: Correctly implements pessimistic locking. Row lock is held throughout the critical section, preventing concurrent modifications to the same row. No deadlock risk as locks are always acquired in consistent order (rowLock then s.mu).

## Finding 3

Location: internal/inventory/store.go:118-152 (OptimisticDeduct method)
Claimed Behavior: Simulates UPDATE ... WHERE id = ? AND version = ? and checks if rows were affected (version matched).
Observed Implementation:
- Fetches current product (with version) under s.mu lock in Get().
- Sleeps 50μs to simulate computation.
- Re-acquires s.mu lock, compares fetched version with current version.
- If version matches, decrements stock, increments version, updates Optimistically counter.
- If version mismatch, increments OptimisticFails counter and returns ErrOptimisticLock.
Assessment: PASS
Severity: 
Notes: Correctly implements optimistic locking with version guard. The version check is performed under s.mu lock, ensuring no concurrent modification can occur between check and update. However, note that s.mu is a global mutex, making the version guard redundant for correctness (but still correct).

## Finding 4

Location: internal/inventory/store.go:154-172 (AtomicDeduct method)
Claimed Behavior: Simulates UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ? in a single atomic statement.
Observed Implementation:
- Locks s.mu, loads product, checks stock >= qty, decrements stock, updates Atomically counter.
- All operations performed while holding s.mu lock.
Assessment: PASS
Severity: 
Notes: Correctly implements atomic conditional update. The check and update are performed atomically under s.mu, preventing lost updates. Equivalent to a database-level atomic statement with WHERE clause.

## Finding 5

Location: internal/inventory/store.go:10-19 (Store struct fields)
Claimed Behavior: Metrics counters track usage of each strategy for observation.
Observed Implementation:
- Six int64 fields: NaivelyDrawn, Pessimistically, Optimistically, OptimisticFails, Atomically, plus mutex and maps.
- All counters updated via atomic.AddInt64.
Assessment: PASS
Severity: 
Notes: Counters are correctly updated atomically. No data race on increments.

## Finding 6

Location: internal/inventory/service.go:28-45 (DeductOptimisticWithRetry method)
Claimed Behavior: Implements optimistic locking with exponential backoff retry logic.
Observed Implementation:
- Loops attempt from 0 to maxRetries inclusive.
- Calls OptimisticDeduct; returns nil on success.
- Returns immediately on non-OptimisticLock errors.
- On OptimisticLock error, continues if attempts remain, else returns ErrOptimisticLock.
- Sleep duration: time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond.
Assessment: PASS
Severity: 
Notes: Correct retry logic with bounded attempts and jittered exponential backoff. Will eventually succeed under contention if maxRetries sufficient.

## Finding 7

Location: internal/inventory/store.go:42-50 (GetRowLock method)
Claimed Behavior: Returns a mutex for row-level locking, creating if necessary.
Observed Implementation:
- Locks s.mu to access rowLocks map.
- If lock missing, creates new *sync.Mutex and stores it.
- Returns lock pointer while holding s.mu, then unlocks via defer.
Assessment: PASS
Severity: 
Notes: Thread-safe creation of per-row mutexes. No risk of returning nil lock.

## Finding 8

Location: internal/inventory/store.go:21-26 (NewStore constructor)
Claimed Behavior: Initializes store with empty product and row lock maps.
Observed Implementation:
- Returns &Store{rowLocks: make(map[int]*sync.Mutex), products: make(map[int]*Product)}.
Assessment: PASS
Severity: 
Notes: Correct zero-value initialization.

## Finding 9

Location: internal/inventory/model.go:12-16 (Product struct)
Claimed Behavior: Represents inventory item with stock and version for optimistic locking.
Observed Implementation:
- Fields: ID, Name, Stock, Version.
Assessment: PASS
Severity: 
Notes: Simple struct matches domain needs.

## Finding 10

Location: internal/inventory/model.go:5-10 (Error variables)
Claimed Behavior: Defines domain-specific error values for error handling.
Observed Implementation:
- ErrNotFound, ErrInsufficientStock, ErrOptimisticLock, ErrInvalidQuantity.
Assessment: PASS
Severity: 
Notes: Proper sentinel errors for callers to check via errors.Is.

## Finding 11

Location: internal/inventory/store.go:53-61 (Get method)
Claimed Behavior: Retrieves product by ID, returns ErrNotFound if missing.
Observed Implementation:
- Locks s.mu, looks up product, returns copy or error.
Assessment: PASS
Severity: 
Returns copy to prevent caller from mutating internal state unintentionally (though copy is shallow; but Product has no pointers so fine).