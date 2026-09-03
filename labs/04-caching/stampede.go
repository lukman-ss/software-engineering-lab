package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

var fixedNow = time.Now()

// TTLWithJitter adds positive random jitter to prevent synchronized expiration (stampede).
// Returns base + random[0, maxJitter). Never reduces TTL below base.
// Upper bound is exclusive: maxJitter is not included [base, base + maxJitter).
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}

// CounterRepository counts calls for stampede demonstration.
// Provides deterministic blocking/coordinating capabilities for tests.
type CounterRepository struct {
	callCount atomic.Int64
	mu        sync.Mutex
	blockCh   chan struct{}
	enteredCh chan struct{} // buffered, signals that a goroutine has entered GetDashboard
}

func NewCounterRepository() *CounterRepository {
	return &CounterRepository{}
}

// Block sets up blocking behavior: all GetDashboard calls will block until Unblock is called.
// Returns immediately. Signals via enteredCh when a call enters GetDashboard.
func (r *CounterRepository) Block() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blockCh = make(chan struct{})
	r.enteredCh = make(chan struct{}, 100)
}

// WaitUntilEntered blocks until a goroutine signals it has entered GetDashboard.
// Call after Block() - returns when a goroutine enters the repository.
func (r *CounterRepository) WaitUntilEntered() {
	r.mu.Lock()
	ch := r.enteredCh
	r.mu.Unlock()
	if ch != nil {
		<-ch
	}
}

// Unblock releases all blocked GetDashboard calls.
func (r *CounterRepository) Unblock() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.blockCh != nil {
		close(r.blockCh)
		r.blockCh = nil
	}
}

// queryDelay is the artificial delay before each DB query returns.
const queryDelay = 5 * time.Millisecond

func (r *CounterRepository) GetDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	// Signal entry into repository (for deterministic test coordination)
	r.mu.Lock()
	if r.enteredCh != nil {
		select {
		case r.enteredCh <- struct{}{}:
		default:
		}
	}
	// Check for blocking (thread-safe)
	ch := r.blockCh
	r.mu.Unlock()

	if ch != nil {
		select {
		case <-ch:
		case <-ctx.Done():
			return Dashboard{}, ctx.Err()
		}
	}

	select {
	case <-ctx.Done():
		return Dashboard{}, ctx.Err()
	case <-time.After(queryDelay):
	}

	newCount := r.callCount.Add(1)
	return Dashboard{InvoiceCountToday: int(newCount)}, nil
}

func (r *CounterRepository) CallCount() int64 {
	return r.callCount.Load()
}

// --- VERSI BROKEN (STAMPEDE) ---

type BrokenStampedeService struct {
	cache CacheInterface
	repo  *CounterRepository
	ttl   time.Duration
}

func NewBrokenStampedeService(cache CacheInterface, repo *CounterRepository) *BrokenStampedeService {
	return &BrokenStampedeService{cache: cache, repo: repo, ttl: 5 * time.Minute}
}

func (s *BrokenStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
			return d, nil
		}
	}

	d, err := s.repo.GetDashboard(ctx, 1, branchID, fixedNow)
	if err != nil {
		return Dashboard{}, err
	}

	data, err := json.Marshal(d)
	if err != nil {
		return Dashboard{}, fmt.Errorf("marshal: %w", err)
	}
	if err := s.cache.Set(ctx, key, string(data), s.ttl); err != nil {
		fmt.Printf("warn: cache set failed: %v\n", err)
	}

	return d, nil
}

// --- VERSI PROTECTED (SINGLEFLIGHT) ---

type ProtectedStampedeService struct {
	cache       CacheInterface
	repo        *CounterRepository
	flight      singleflight.Group
	ttl         time.Duration
	onWaitEntry func(key string)
}

func NewProtectedStampedeService(cache CacheInterface, repo *CounterRepository) *ProtectedStampedeService {
	return &ProtectedStampedeService{cache: cache, repo: repo, ttl: 5 * time.Minute}
}

func (s *ProtectedStampedeService) SetOnWaitEntry(fn func(key string)) {
	s.onWaitEntry = fn
}

func (s *ProtectedStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
			return d, nil
		}
	}

	ch := s.flight.DoChan(key, func() (interface{}, error) {
		rebuildCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()

		cached, err := s.cache.Get(rebuildCtx, key)
		if err == nil && cached != "" {
			var d Dashboard
			if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
				return d, nil
			}
		}

		d, err := s.repo.GetDashboard(rebuildCtx, 1, branchID, fixedNow)
		if err != nil {
			return Dashboard{}, err
		}

		data, err := json.Marshal(d)
		if err != nil {
			return Dashboard{}, fmt.Errorf("marshal: %w", err)
		}
		if setErr := s.cache.Set(rebuildCtx, key, string(data), s.ttl); setErr != nil {
			fmt.Printf("warn: cache set failed: %v\n", setErr)
		}
		return d, nil
	})

	if s.onWaitEntry != nil {
		s.onWaitEntry(key)
	}

	select {
	case <-ctx.Done():
		return Dashboard{}, ctx.Err()
	case result := <-ch:
		if result.Err != nil {
			return Dashboard{}, result.Err
		}
		return result.Val.(Dashboard), nil
	}
}
