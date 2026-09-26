package inventory

import (
	"sync"
	"sync/atomic"
	"time"
)

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

func (s *Store) GetRowLock(id int) *sync.Mutex {
	s.mu.Lock()
	defer s.mu.Unlock()
	lock, exists := s.rowLocks[id]
	if !exists {
		lock = &sync.Mutex{}
		s.rowLocks[id] = lock
	}
	return lock
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

// NaiveDeduct performs an unsynchronized read-modify-write across separate statements,
// reproducing the classic lost-update concurrency anomaly.
func (s *Store) NaiveDeduct(id int, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	// 1. Read phase
	p, err := s.Get(id)
	if err != nil {
		return err
	}
	if p.Stock < qty {
		return ErrInsufficientStock
	}

	// Artificial yield/pause simulating application calculation window
	time.Sleep(100 * time.Microsecond)

	// 2. Write phase using stale calculation
	s.mu.Lock()
	curr := s.products[id]
	curr.Stock = p.Stock - qty
	s.mu.Unlock()

	atomic.AddInt64(&s.NaivelyDrawn, int64(qty))
	return nil
}

// PessimisticDeduct simulates SELECT ... FOR UPDATE by acquiring an exclusive row lock
// before reading and holding it until the update completes.
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

// OptimisticDeduct simulates UPDATE ... WHERE id = ? AND version = ?
// and checks if rows were affected (version matched).
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

	// Artificial computation delay
	time.Sleep(50 * time.Microsecond)

	s.mu.Lock()
	defer s.mu.Unlock()
	curr, exists := s.products[id]
	if !exists {
		return ErrNotFound
	}

	// Compare version guard
	if curr.Version != p.Version {
		atomic.AddInt64(&s.OptimisticFails, 1)
		return ErrOptimisticLock
	}

	curr.Stock -= qty
	curr.Version++
	atomic.AddInt64(&s.Optimistically, int64(qty))
	return nil
}

// AtomicDeduct simulates UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?
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
