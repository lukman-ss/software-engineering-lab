# Code Snippets

## Snippet 1 — Product Model

Source File: `internal/inventory/model.go`
Purpose: Domain model dan error definitions.

```go
type Product struct {
	ID      int
	Name    string
	Stock   int
	Version int
}

var (
	ErrNotFound          = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOptimisticLock    = errors.New("optimistic lock conflict: record modified by another transaction")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
)
```

---

## Snippet 2 — Store (In-Memory Engine)

Source File: `internal/inventory/store.go`
Purpose: Simulated storage engine dengan row-level locks, version guard, atomic counters.

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

func (s *Store) Get(id int) (Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, exists := s.products[id]
	if !exists {
		return Product{}, ErrNotFound
	}
	return *p, nil
}
```

---

## Snippet 3 — Naive Deduct (Lost Update Path)

Source File: `internal/inventory/store.go` (lines 65-89)
Purpose: Unsynchronized read-modify-write reproducendo lost update anomaly.

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

	time.Sleep(100 * time.Microsecond) // simulating application calculation window

	s.mu.Lock()
	curr := s.products[id]
	curr.Stock = p.Stock - qty // STALE WRITE
	s.mu.Unlock()

	atomic.AddInt64(&s.NaivelyDrawn, int64(qty))
	return nil
}
```

---

## Snippet 4 — Pessimistic Deduct (Row-Level Lock)

Source File: `internal/inventory/store.go` (lines 93-116)
Purpose: `SELECT ... FOR UPDATE` simulasi dengan per-row mutex.

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

---

## Snippet 5 — Optimistic Deduct (Version Guard)

Source File: `internal/inventory/store.go` (lines 120-152)
Purpose: `UPDATE ... WHERE id = ? AND version = ?` simulasi dengan version check.

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

	time.Sleep(50 * time.Microsecond) // simulating computation delay

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

---

## Snippet 6 — Optimistic Retry (Jittered Exponential Backoff)

Source File: `internal/inventory/service.go` (lines 28-45)
Purpose: Optimistic locking dengan retry + jittered exponential backoff.

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

---

## Snippet 7 — Atomic Deduct (Single-Statement)

Source File: `internal/inventory/store.go` (lines 155-173)
Purpose: `UPDATE ... SET stock = stock - N WHERE stock >= N` simulasi dengan statement-level atomicity.

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

---

## Snippet 8 — Test: Naive Lost Update

Source File: `tests/locking_test.go` (lines 10-38)
Purpose: Demonstrasi lost update dengan 50 concurrent goroutines.

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
	t.Logf("Lost Update Demonstrated: 50 deduct calls occurred, but final stock is %d (expected 50 under proper locking)", p.Stock)
}
```

---

## Snippet 9 — Test: Optimistic Locking Conflict

Source File: `tests/locking_test.go` (lines 85-121)
Purpose: Verifikasi optimisitik conflict detection dan state invariant.

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
	t.Logf("Optimistic Locking Conflict Demonstrated: %d succeeded, %d conflicts rejected", successCount, conflictCount)
}
```

---

## Snippet 10 — Demo Runner

Source File: `cmd/demo/main.go` (lines 11-120)
Purpose: CLI runner menampilkan 5 skenario side-by-side.

```go
func main() {
	fmt.Println("==============================================================")
	fmt.Println("  Optimistic vs Pessimistic Locking & Atomic Operations   ")
	fmt.Println("==============================================================")

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
	// 3. Optimistic locking direct
	// 4. Optimistic locking with retry
	// 5. Atomic single-statement update
	// ... (same pattern as above)
}
```

---

## Snippet 11 — Service Layer

Source File: `internal/inventory/service.go`
Purpose: Business logic layer wrapping store methods, adding retry logic.

```go
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (svc *Service) DeductNaive(id int, qty int) error {
	return svc.store.NaiveDeduct(id, qty)
}

func (svc *Service) DeductPessimistic(id int, qty int) error {
	return svc.store.PessimisticDeduct(id, qty)
}

func (svc *Service) DeductOptimisticDirect(id int, qty int) error {
	return svc.store.OptimisticDeduct(id, qty)
}

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

func (svc *Service) DeductAtomic(id int, qty int) error {
	return svc.store.AtomicDeduct(id, qty)
}
```