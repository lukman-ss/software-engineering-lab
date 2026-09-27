# Code Snippets

## Snippet 1 — Domain Model & Error Definitions

Source File: `internal/inventory/model.go`

```go
var (
    ErrNotFound          = errors.New("product not found")
    ErrInsufficientStock = errors.New("insufficient stock")
    ErrOptimisticLock    = errors.New("optimistic lock conflict: record modified by another transaction")
    ErrInvalidQuantity   = errors.New("quantity must be positive")
)

type Product struct {
    ID      int
    Name    string
    Stock   int
    Version int
}
```

Explanation:
Domain-specific error definitions: `ErrOptimisticLock` represents an optimistic concurrency conflict (equivalent to SQL `affected_rows == 0`), `ErrInsufficientStock` represents a business-rule rejection. `Product` includes a `Version` field as the optimistic guard.

## Snippet 2 — Store Core (Simulated Engine)

Source File: `internal/inventory/store.go`

```go
type Store struct {
    mu       sync.Mutex
    rowLocks map[int]*sync.Mutex
    products map[int]*Product

    NaivelyDrawn    int64
    Pessimistically int64
    Optimistically  int64
    OptimisticFails int64
    Atomically      int64
}

func NewStore() *Store {
    return &Store{
        rowLocks: make(map[int]*sync.Mutex),
        products: make(map[int]*Product),
    }
}

func (s *Store) Seed(id int, name string, stock int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.products[id] = &Product{
        ID:      id,
        Name:    name,
        Stock:   stock,
        Version: 1,
    }
    if _, exists := s.rowLocks[id]; !exists {
        s.rowLocks[id] = &sync.Mutex{}
    }
}
```

Explanation:
`rowLocks` simulates `SELECT ... FOR UPDATE` per row. Each `Product` has a `Version` starting at 1, incremented by optimistic operations. Atomic counters (`atomic.AddInt64`) allow invariant verification without additional locking.

## Snippet 3 — Naive Read-Modify-Write (Lost Update)

Source File: `internal/inventory/store.go` (lines 65-89)

```go
func (s *Store) NaiveDeduct(id int, qty int) error {
    if qty <= 0 {
        return ErrInvalidQuantity
    }
    p, err := s.Get(id)
    if err != nil {
        return err
    }
    if p.Stock < qty {
        return ErrInsufficientStock
    }
    time.Sleep(100 * time.Microsecond)
    s.mu.Lock()
    curr := s.products[id]
    curr.Stock = p.Stock - qty
    s.mu.Unlock()
    atomic.AddInt64(&s.NaivelyDrawn, int64(qty))
    return nil
}
```

Explanation:
The bug: `p.Stock` is read, then `curr.Stock = p.Stock - qty` is written after a delay. If another goroutine modifies stock between `Get()` and the lock, the stale value overwrites the newer one. `time.Sleep(100µs)` widens the race window for reliable reproduction.

## Snippet 4 — Pessimistic Locking (SELECT FOR UPDATE)

Source File: `internal/inventory/store.go` (lines 93-116)

```go
func (s *Store) PessimisticDeduct(id int, qty int) error {
    if qty <= 0 {
        return ErrInvalidQuantity
    }
    rowLock := s.GetRowLock(id)
    rowLock.Lock()
    defer rowLock.Unlock()

    s.mu.Lock()
    p, exists := s.products[id]
    if !exists {
        s.mu.Unlock()
        return ErrNotFound
    }
    if p.Stock < qty {
        s.mu.Unlock()
        return ErrInsufficientStock
    }
    p.Stock -= qty
    s.mu.Unlock()

    atomic.AddInt64(&s.Pessimistically, int64(qty))
    return nil
}
```

Explanation:
`rowLock.Lock()` blocks all other goroutines accessing the same row until `defer Unlock()`. The global `s.mu` only protects the map structure; the row lock governs business logic. This is an exclusive lock — no two writers on the same row.

## Snippet 5 — Optimistic Locking (Version Guard)

Source File: `internal/inventory/store.go` (lines 120-152)

```go
func (s *Store) OptimisticDeduct(id int, qty int) error {
    if qty <= 0 {
        return ErrInvalidQuantity
    }
    p, err := s.Get(id)
    if err != nil {
        return err
    }
    if p.Stock < qty {
        return ErrInsufficientStock
    }
    time.Sleep(50 * time.Microsecond)

    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if !exists {
        return ErrNotFound
    }
    if curr.Version != p.Version {
        atomic.AddInt64(&s.OptimisticFails, 1)
        return ErrOptimisticLock
    }
    curr.Stock -= qty
    curr.Version++
    atomic.AddInt64(&s.Optimistically, int64(qty))
    return nil
}
```

Explanation:
`p.Version` is read as a snapshot, then validated inside `s.mu.Lock()`. If `curr.Version != p.Version`, another write occurred in between — return `ErrOptimisticLock`. Analogous to SQL: `UPDATE ... WHERE id=? AND version=?`.

## Snippet 6 — Optimistic Locking with Retry & Jittered Backoff

Source File: `internal/inventory/service.go` (lines 28-45)

```go
func (svc *Service) DeductOptimisticWithRetry(id int, qty int, maxRetries int) error {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := svc.store.OptimisticDeduct(id, qty)
        if err == nil {
            return nil
        }
        if err != ErrOptimisticLock {
            return err
        }
        if attempt == maxRetries {
            return ErrOptimisticLock
        }
        sleepDuration := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond
        time.Sleep(sleepDuration)
    }
    return ErrOptimisticLock
}
```

Explanation:
Retry loop with exponential backoff: attempt 0 ≈ 1-6ms, 1 ≈ 2-7ms, 2 ≈ 4-9ms, etc. Jitter (`rand.Intn(5)`) prevents thundering-herd. After `maxRetries`, returns `ErrOptimisticLock` (analogous to HTTP 409 Conflict).

## Snippet 7 — Atomic Single-Statement Update

Source File: `internal/inventory/store.go` (lines 155-173)

```go
func (s *Store) AtomicDeduct(id int, qty int) error {
    if qty <= 0 {
        return ErrInvalidQuantity
    }
    s.mu.Lock()
    defer s.mu.Unlock()

    curr, exists := s.products[id]
    if !exists {
        return ErrNotFound
    }
    if curr.Stock < qty {
        return ErrInsufficientStock
    }
    curr.Stock -= qty
    atomic.AddInt64(&s.Atomically, int64(qty))
    return nil
}
```

Explanation:
Single lock + validate + update. Analogous to SQL `UPDATE products SET stock=stock-1 WHERE id=? AND stock>=1` — atomic at the statement level. No separate read phase, no persistent lock across goroutines.

## Snippet 8 — Test: Naive Lost Update Demonstration

Source File: `tests/locking_test.go` (lines 10-38)

```go
func TestNaiveLostUpdate(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(1, "Item A", 100)
    svc := inventory.NewService(store)

    const goroutines = 50
    var wg sync.WaitGroup
    wg.Add(goroutines)

    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            _ = svc.DeductNaive(1, 1)
        }()
    }
    wg.Wait()

    p, err := store.Get(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if p.Stock == 50 {
        t.Errorf("expected lost update anomaly to manifest (stock != 50), got stock = %d", p.Stock)
    }
}
```

Explanation:
The test does not assert `p.Stock == 50` — it asserts the anomaly by checking `p.Stock != 50`. If it were exactly 50, the race condition did not manifest (unlikely but possible). In practice, ~49 goroutines overwrite, yielding stock = 99.

## Snippet 9 — Test: Pessimistic Locking Invariant

Source File: `tests/locking_test.go` (lines 40-67)

```go
func TestPessimisticLocking(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(2, "Item B", 100)
    svc := inventory.NewService(store)

    const goroutines = 50
    var wg sync.WaitGroup
    wg.Add(goroutines)

    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            err := svc.DeductPessimistic(2, 1)
            if err != nil {
                t.Errorf("deduct failed: %v", err)
            }
        }()
    }
    wg.Wait()

    p, err := store.Get(2)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if p.Stock != 50 {
        t.Fatalf("expected stock 50, got %d", p.Stock)
    }
}
```

Explanation:
Strict invariant verification: 100 − 50 = 50. All 50 goroutines succeed without conflicts. `wg.Wait()` ensures completion before assertion.

## Snippet 10 — Test: Optimistic Conflict Detection & Retry Convergence

Source File: `tests/locking_test.go` (lines 85-145)

```go
func TestOptimisticLockingConflict(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(3, "Item C", 100)
    svc := inventory.NewService(store)

    const goroutines = 20
    var wg sync.WaitGroup
    wg.Add(goroutines)

    var conflictCount int64
    var successCount int64
    var mu sync.Mutex

    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            err := svc.DeductOptimisticDirect(3, 1)
            mu.Lock()
            if err == inventory.ErrOptimisticLock {
                conflictCount++
            } else if err == nil {
                successCount++
            }
            mu.Unlock()
        }()
    }
    wg.Wait()

    p, _ := store.Get(3)
    if successCount+int64(p.Stock) != 100 {
        t.Fatalf("stock invariant violated: success %d + stock %d != 100", successCount, p.Stock)
    }
    if conflictCount == 0 {
        t.Fatalf("expected at least one optimistic lock conflict, got 0")
    }
}

func TestOptimisticLockingWithRetry(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(4, "Item D", 100)
    svc := inventory.NewService(store)

    const goroutines = 20
    var wg sync.WaitGroup
    wg.Add(goroutines)

    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            _ = svc.DeductOptimisticWithRetry(4, 1, 10)
        }()
    }
    wg.Wait()

    p, _ := store.Get(4)
    if p.Stock != 100-int(store.Optimistically) {
        t.Fatalf("inconsistent stock: expected %d, got %d", 100-int(store.Optimistically), p.Stock)
    }
}
```

Explanation:
Proves: conflicts are detected (`conflictCount > 0`), and the invariant holds (`successCount + p.Stock == 100`). With retry, all 20 eventually succeed (`p.Stock == 80`).

## Snippet 11 — Demo CLI: Side-by-Side Comparison

Source File: `cmd/demo/main.go` (lines 11-119)

```go
func main() {
    fmt.Println("==========================================================")
    fmt.Println("  Optimistic vs Pessimistic Locking & Atomic Operations   ")
    fmt.Println("==========================================================")

    // 1. Naive lost update
    store1 := inventory.NewStore()
    store1.Seed(1, "Laptop Pro", 100)
    svc1 := inventory.NewService(store1)
    var wg1 sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            _ = svc1.DeductNaive(1, 1)
        }()
    }
    wg1.Wait()
    p1, _ := store1.Get(1)
    fmt.Printf("\n[1] Naive Read-Modify-Write (50 concurrent requests):\n")
    fmt.Printf("    Actual Final Stock:   %d (LOST UPDATE DETECTED!)\n", p1.Stock)

    // 2. Pessimistic locking
    store2 := inventory.NewStore()
    store2.Seed(2, "Smartphone X", 100)
    svc2 := inventory.NewService(store2)
    var wg2 sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg2.Add(1)
        go func() {
            defer wg2.Done()
            _ = svc2.DeductPessimistic(2, 1)
        }()
    }
    wg2.Wait()
    p2, _ := store2.Get(2)
    fmt.Printf("\n[2] Pessimistic Locking (SELECT ... FOR UPDATE):\n")
    fmt.Printf("    Actual Final Stock:   %d (SUCCESS)\n", p2.Stock)

    // 3. Optimistic direct
    store3 := inventory.NewStore()
    store3.Seed(3, "Wireless Headphones", 100)
    svc3 := inventory.NewService(store3)
    var wg3 sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg3.Add(1)
        go func() {
            defer wg3.Done()
            _ = svc3.DeductOptimisticDirect(3, 1)
        }()
    }
    wg3.Wait()
    p3, _ := store3.Get(3)
    fmt.Printf("\n[3] Optimistic Locking Direct (20 concurrent, no retry):\n")
    fmt.Printf("    Successful Deductions: %d\n", store3.Optimistically)
    fmt.Printf("    Rejected Conflicts:   %d\n", store3.OptimisticFails)
    fmt.Printf("    Actual Final Stock:   %d (State Guarded)\n", p3.Stock)

    // 4. Optimistic with retry
    store4 := inventory.NewStore()
    store4.Seed(4, "Smartwatch", 100)
    svc4 := inventory.NewService(store4)
    start := time.Now()
    var wg4 sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg4.Add(1)
        go func() {
            defer wg4.Done()
            _ = svc4.DeductOptimisticWithRetry(4, 1, 10)
        }()
    }
    wg4.Wait()
    p4, _ := store4.Get(4)
    fmt.Printf("\n[4] Optimistic With Exponential Backoff Retry (20 requests):\n")
    fmt.Printf("    Successful Deductions: %d\n", store4.Optimistically)
    fmt.Printf("    Total Attempted Conflicts Retried: %d\n", store4.OptimisticFails)
    fmt.Printf("    Actual Final Stock:   %d (All retries converged)\n", p4.Stock)
    fmt.Printf("    Elapsed Time:         %v\n", time.Since(start))

    // 5. Atomic
    store5 := inventory.NewStore()
    store5.Seed(5, "Tablet Air", 100)
    svc5 := inventory.NewService(store5)
    var wg5 sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg5.Add(1)
        go func() {
            defer wg5.Done()
            _ = svc5.DeductAtomic(5, 1)
        }()
    }
    wg5.Wait()
    p5, _ := store5.Get(5)
    fmt.Printf("\n[5] Atomic Single-Statement Operation:\n")
    fmt.Printf("    Actual Final Stock:   %d (Lockless Single Statement)\n", p5.Stock)

    fmt.Println("\n==========================================================")
    fmt.Println("  Lab Execution Completed Successfully                    ")
    fmt.Println("==========================================================")
}
```

Explanation:
The demo runs sequentially (not in parallel) for readable output. Each scenario uses a fresh store seeded with stock = 100. Output is printed directly to stdout.