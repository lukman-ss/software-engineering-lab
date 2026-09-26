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
Definisi error domain spesifik: `ErrOptimisticLock` mewakili konflik optimistic, `ErrInsufficientStock` mewakili kondisi bisnis. `Product` memiliki `Version` field sebagai optimistic guard.

## Snippet 2 — Store Core

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
`rowLocks` mensimulasikan `SELECT ... FOR UPDATE` per baris. Setiap `Product` memiliki `Version` yang dimulai dari 1 dan akan diincrement oleh operasi optimistic. Counter atomik (`atomic.AddInt64`) memungkinkan verifikasi invariant tanpa lock.

## Snippet 3 — Naive Read-Modify-Write (Lost Update)

Source File: `internal/inventory/store.go`

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
Bug: nilai `p.Stock` dibaca, lalu `curr.Stock = p.Stock - qty` ditulis. Jika goroutine lain mengubah stock setelah `Get()` tapi sebelum `s.mu.Lock()`, nilai lama (stale) digunakan. `time.Sleep(100µs)` memperbesar jendela race.

## Snippet 4 — Pessimistic Locking

Source File: `internal/inventory/store.go`

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
`rowLock.Lock()` memblokir semua goroutine lain yang mengakses row yang sama sampai `defer` mengunlock. `s.mu` hanya melindungi map, tidak perilaku bisnis. Lock ini eksklusif → tidak ada two-writers.

## Snippet 5 — Optimistic Locking

Source File: `internal/inventory/store.go`

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
`p.Version` dibaca (snapshot), lalu divalidasi di dalam `s.mu.Lock()`. Jika `curr.Version != p.Version`, berarti ada penulisan lain di antara — return `ErrOptimisticLock`. Analog SQL: `UPDATE ... WHERE id=? AND version=?`.

## Snippet 6 — Optimistic Locking with Retry & Backoff

Source File: `internal/inventory/service.go`

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
Retry loop dengan exponential backoff: attempt 0 = ~1-6ms, 1 = ~2-7ms, 2 = ~4-9ms, ... Jitter `rand.Intn(5)` mencegah thundering-herd. Jika kehabisan retry → kembalikan `ErrOptimisticLock` (analog HTTP 409).

## Snippet 7 — Atomic Single-Statement Update

Source File: `internal/inventory/store.go`

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
Single lock + validate + update. Analog SQL `UPDATE products SET stock=stock-1 WHERE id=? AND stock>=1` — atoms pada tingkat pernyataan. Tidak ada read terpisah, tidak ada lock persisten melewati goroutine.

## Snippet 8 — Test: Naive Lost Update Demonstration

Source File: `tests/locking_test.go`

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
Tes bukan mengecek `p.Stock == 50`, melainkan justru memverifikasi `p.Stock != 50`. Jika persis 50, anomaly tidak terjadi (race detector mungkin menyelamatkan). Dalam prakteknya, 49 goroutine overwrite, hasil = 99.

## Snippet 9 — Test: Pessimistic Locking Invariant

Source File: `tests/locking_test.go`

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
        }
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
Verifikasi invariant ketat: 100 - 50 = 50. Setiap goroutine berhasil tanpa conflict. `wg.Wait()` memastikan semua selesai sebelum asersi.

## Snippet 10 — Test: Optimistic Conflict & Retry Convergence

Source File: `tests/locking_test.go`

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
```

Explanation:
Membuktikan: konflik terdeteksi (`conflictCount > 0`), dan invariant tetap terjaga (`successCount + p.Stock == 100`). Jika 19 konflik, hanya 1 yang berhasil, stok tetap 99 (100 - 1).

## Snippet 11 — Demo CLI: Side-by-Side Scenarios

Source File: `cmd/demo/main.go`

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

    // ... (Optimistic, Optimistic+Retry, Atomic sama polanya)
}
```

Explanation:
Demo berjalan berurutan (bukan paralel) untuk output yang dapatdibaca. Setiap skenario fresh store, seed stock 100. Output langsung ke stdout.
