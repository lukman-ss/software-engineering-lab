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

// TTLWithJitter adds random jitter to prevent synchronized expiration (stampede).
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}

// CounterRepository counts calls for stampede demonstration.
type CounterRepository struct {
	callCount atomic.Int64
	mu        sync.Mutex
	blockCh   chan struct{}
}

func NewCounterRepository() *CounterRepository {
	return &CounterRepository{}
}

// queryDelay is the artificial delay before each DB query returns.
// This ensures a realistic time-of-check-to-time-of-use window for cache stampede demonstrations.
const queryDelay = 5 * time.Millisecond

func (r *CounterRepository) GetDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	// Check for blocking (thread-safe)
	r.mu.Lock()
	ch := r.blockCh
	r.mu.Unlock()

	if ch != nil {
		select {
		case <-ch:
		case <-ctx.Done():
			return Dashboard{}, ctx.Err()
		}
	}

	// Simulate DB query latency so concurrent cache misses accumulate
	// before any single goroutine finishes and writes to cache.
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

// Block makes the repository wait before returning, for testing coordination.
func (r *CounterRepository) Block() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blockCh = make(chan struct{})
}

// Unblock releases all blocked repository calls.
func (r *CounterRepository) Unblock() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.blockCh != nil {
		close(r.blockCh)
		r.blockCh = nil
	}
}

// --- VERSI BROKEN (STAMPEDE) ---

// BrokenStampedeService demonstrates cache stampede vulnerability.
// Multiple concurrent cache misses result in parallel DB queries.
// INTENTIONALLY BROKEN FOR LAB DEMONSTRATION.
// Do not use as production code.
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

	// 1. Check cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
			return d, nil
		}
		// If unmarshal fails, we still fall through to DB
	}

	// 2. Cache miss -> parallel DB queries (NO protection)
	// In real system: 100 concurrent requests → 100 DB queries
	d, err := s.repo.GetDashboard(ctx, 1, branchID, time.Now())
	if err != nil {
		return Dashboard{}, err
	}

	// 3. Set cache
	data, err := json.Marshal(d)
	if err != nil {
		return Dashboard{}, fmt.Errorf("marshal: %w", err)
	}
	if err := s.cache.Set(ctx, key, string(data), s.ttl); err != nil {
		// Log but don't fail
		fmt.Printf("warn: cache set failed: %v\n", err)
	}

	return d, nil
}

// --- VERSI PROTECTED (SINGLEFLIGHT) ---

// ProtectedStampedeService protects against cache stampede using singleflight.
// Concurrent requests for same key share a single DB query result.
type ProtectedStampedeService struct {
	cache  CacheInterface
	repo   *CounterRepository
	flight singleflight.Group
	ttl    time.Duration
}

func NewProtectedStampedeService(cache CacheInterface, repo *CounterRepository) *ProtectedStampedeService {
	return &ProtectedStampedeService{cache: cache, repo: repo, ttl: 5 * time.Minute}
}

func (s *ProtectedStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	// 1. First cache check
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
			return d, nil
		}
		// If unmarshal fails, fall through to singleflight
	}

	// 2. Cache miss -> use singleflight with DoChan for context-awareness
	// DoChan returns a channel that delivers the result, allowing
	// callers to also select on ctx.Done() for proper cancellation.
	ch := s.flight.DoChan(key, func() (interface{}, error) {
		// Double-check cache inside singleflight gate
		// This handles the window where another goroutine may have populated cache
		// between our first check and acquiring the singleflight slot
		cached, err := s.cache.Get(ctx, key)
		if err == nil && cached != "" {
			var d Dashboard
			if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
				return d, nil
			}
		}

		// Only one goroutine reaches here - it fetches from DB
		d, err := s.repo.GetDashboard(ctx, 1, branchID, time.Now())
		if err != nil {
			return Dashboard{}, err
		}

		// Populate cache
		data, err := json.Marshal(d)
		if err != nil {
			return Dashboard{}, fmt.Errorf("marshal: %w", err)
		}
		if setErr := s.cache.Set(ctx, key, string(data), s.ttl); setErr != nil {
			fmt.Printf("warn: cache set failed: %v\n", setErr)
		}
		return d, nil
	})

	// 3. Wait with context awareness
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
